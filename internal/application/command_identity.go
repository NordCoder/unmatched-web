package application

import (
	"crypto/sha256"
	"fmt"
	"unicode/utf8"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	"golang.org/x/text/unicode/norm"
)

func normalizeCommandRequest(principal model.PrincipalID, command contracts.Command) (model.PrincipalID, contracts.Command, []byte, error) {
	if err := validateEnvelope(principal, command); err != nil {
		return "", contracts.Command{}, nil, err
	}

	normalizedPrincipal, err := normalizeIdentifier("principal ID", string(principal))
	if err != nil {
		return "", contracts.Command{}, nil, err
	}
	commandID, err := normalizeIdentifier("command ID", string(command.ID))
	if err != nil {
		return "", contracts.Command{}, nil, err
	}
	schemaVersion, err := normalizeIdentifier("command schema version", command.SchemaVersion)
	if err != nil {
		return "", contracts.Command{}, nil, err
	}
	commandType, err := normalizeIdentifier("command type", command.Type)
	if err != nil {
		return "", contracts.Command{}, nil, err
	}
	matchID, err := normalizeOptionalIdentifier("match ID", string(command.MatchID))
	if err != nil {
		return "", contracts.Command{}, nil, err
	}
	actorPlayerID, err := normalizeOptionalIdentifier("actor player ID", string(command.ActorPlayerID))
	if err != nil {
		return "", contracts.Command{}, nil, err
	}

	payloadValue, err := parseCanonicalJSON(command.Payload)
	if err != nil {
		return "", contracts.Command{}, nil, opError(CodeInvalidCommand, "command payload: "+err.Error())
	}
	payloadObject, ok := payloadValue.(map[string]any)
	if !ok {
		return "", contracts.Command{}, nil, opError(CodeInvalidCommand, "command payload must be a JSON object")
	}
	normalizedPayload, err := normalizeCommandPayload(commandType, payloadObject)
	if err != nil {
		return "", contracts.Command{}, nil, err
	}

	var revisionValue any = absentValue()
	revisionPolicy := "absent"
	var normalizedRevision *uint64
	if command.ExpectedRevision != nil {
		if *command.ExpectedRevision > uint64(maxCanonicalInteger) {
			return "", contracts.Command{}, nil, opError(CodeInvalidCommand, "expected revision exceeds the safe integer range")
		}
		revision := *command.ExpectedRevision
		normalizedRevision = &revision
		revisionPolicy = "exact"
		revisionValue = int64(revision)
	}

	normalizedCommand := contracts.Command{
		ID:               model.CommandID(commandID),
		SchemaVersion:    schemaVersion,
		Type:             commandType,
		MatchID:          model.MatchID(matchID),
		ActorPlayerID:    model.PlayerID(actorPlayerID),
		ExpectedRevision: normalizedRevision,
		Payload:          canonicalJSON(normalizedPayload),
	}
	normalizedPrincipalID := model.PrincipalID(normalizedPrincipal)

	identity := map[string]any{
		"fingerprint_schema_version": persistence.FingerprintSchemaV1,
		"principal_id":               normalizedPrincipal,
		"lifecycle_scope":            string(commandScope(commandType)),
		"match_id":                   optionalIdentityValue(matchID),
		"actor_player_id":            optionalIdentityValue(actorPlayerID),
		"command_schema_version":     schemaVersion,
		"type":                       commandType,
		"normalized_payload":         normalizedPayload,
		"expected_revision_policy":   revisionPolicy,
		"expected_revision":          revisionValue,
	}
	canonicalIdentity := canonicalJSON(identity)
	fingerprint := sha256.Sum256(canonicalIdentity)
	return normalizedPrincipalID, normalizedCommand, fingerprint[:], nil
}

func normalizeIdentifier(name, value string) (string, error) {
	if value == "" {
		return "", opError(CodeInvalidCommand, name+" is required")
	}
	if !utf8.ValidString(value) {
		return "", opError(CodeInvalidCommand, name+" contains invalid Unicode")
	}
	return norm.NFC.String(value), nil
}

func normalizeOptionalIdentifier(name, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	if !utf8.ValidString(value) {
		return "", opError(CodeInvalidCommand, name+" contains invalid Unicode")
	}
	return norm.NFC.String(value), nil
}

func normalizeCommandPayload(commandType string, payload map[string]any) (map[string]any, error) {
	switch commandType {
	case CommandCreateMatch:
		if err := rejectUnknownFields(payload, "definition_key", "fighter_definition_id"); err != nil {
			return nil, err
		}
		definitionKey, err := requiredStringField(payload, "definition_key")
		if err != nil {
			return nil, err
		}
		fighterDefinition, err := requiredStringField(payload, "fighter_definition_id")
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"definition_key":        definitionKey,
			"fighter_definition_id": fighterDefinition,
		}, nil
	case CommandJoinMatch:
		if err := rejectUnknownFields(payload, "fighter_definition_id"); err != nil {
			return nil, err
		}
		fighterDefinition, err := requiredStringField(payload, "fighter_definition_id")
		if err != nil {
			return nil, err
		}
		return map[string]any{"fighter_definition_id": fighterDefinition}, nil
	case CommandStartAction:
		if err := rejectUnknownFields(payload, "kind", "source_ref", "context"); err != nil {
			return nil, err
		}
		kind, err := requiredStringField(payload, "kind")
		if err != nil {
			return nil, err
		}
		sourceRef := ""
		if value, exists := payload["source_ref"]; exists {
			var ok bool
			sourceRef, ok = value.(string)
			if !ok {
				return nil, invalidPayloadField("source_ref", "must be a string")
			}
		}
		contextValue := map[string]any{}
		if value, exists := payload["context"]; exists {
			var ok bool
			contextValue, ok = value.(map[string]any)
			if !ok {
				return nil, invalidPayloadField("context", "must be an object")
			}
		}
		return map[string]any{
			"kind":       kind,
			"source_ref": sourceRef,
			"context":    contextValue,
		}, nil
	case CommandSubmitChoice:
		if err := rejectUnknownFields(payload, "interaction_id", "choice"); err != nil {
			return nil, err
		}
		interactionID, err := requiredStringField(payload, "interaction_id")
		if err != nil {
			return nil, err
		}
		choice, exists := payload["choice"]
		if !exists {
			return nil, invalidPayloadField("choice", "is required")
		}
		return map[string]any{
			"interaction_id": interactionID,
			"choice":         choice,
		}, nil
	default:
		return nil, opError(CodeInvalidCommand, "unsupported command type")
	}
}

func rejectUnknownFields(payload map[string]any, allowedFields ...string) error {
	allowed := make(map[string]struct{}, len(allowedFields))
	for _, field := range allowedFields {
		allowed[field] = struct{}{}
	}
	for field := range payload {
		if _, ok := allowed[field]; !ok {
			return invalidPayloadField(field, "is not allowed")
		}
	}
	return nil
}

func requiredStringField(payload map[string]any, field string) (string, error) {
	value, exists := payload[field]
	if !exists {
		return "", invalidPayloadField(field, "is required")
	}
	text, ok := value.(string)
	if !ok {
		return "", invalidPayloadField(field, "must be a string")
	}
	return text, nil
}

func invalidPayloadField(field, message string) error {
	return opError(CodeInvalidCommand, fmt.Sprintf("command payload field %q %s", field, message))
}

func optionalIdentityValue(value string) any {
	if value == "" {
		return absentValue()
	}
	return value
}

func absentValue() map[string]any {
	return map[string]any{"absent": true}
}
