package application

import (
	"encoding/json"
	"errors"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

const (
	CommandSchemaV1 = "core-command/v1"

	CommandCreateMatch  = "CreateMatch"
	CommandJoinMatch    = "JoinMatch"
	CommandStartAction  = "StartAction"
	CommandSubmitChoice = "SubmitChoice"

	ActionManeuver = "MANEUVER"
	ActionScheme   = "SCHEME"
	ActionAttack   = "ATTACK"
)

const (
	CodeInvalidCommand     = "invalid_command"
	CodeNotFound           = "not_found"
	CodeUnauthorized       = "unauthorized"
	CodeRevisionConflict   = "revision_conflict"
	CodeCommandConflict    = "command_conflict"
	CodePendingInteraction = "pending_interaction"
	CodeInvalidChoice      = "invalid_choice"
	CodeRulesRejected      = "rules_rejected"
	CodeInternal           = "internal"
)

type OperationalError struct {
	Code    string
	Message string
}

func (e *OperationalError) Error() string { return e.Code + ": " + e.Message }

func CodeOf(err error) string {
	var operational *OperationalError
	if errors.As(err, &operational) {
		return operational.Code
	}
	return ""
}

type CreateMatchPayload struct {
	DefinitionKey     string             `json:"definition_key"`
	FighterDefinition model.DefinitionID `json:"fighter_definition_id"`
}

type JoinMatchPayload struct {
	FighterDefinition model.DefinitionID `json:"fighter_definition_id"`
}

type StartActionPayload struct {
	Kind      string                     `json:"kind"`
	SourceRef string                     `json:"source_ref,omitempty"`
	Context   map[string]json.RawMessage `json:"context,omitempty"`
}

type SubmitChoicePayload struct {
	InteractionID model.InteractionID `json:"interaction_id"`
	Choice        json.RawMessage     `json:"choice"`
}

// ProjectedInteraction is an explicit delivery allow-list. Internal resolver
// state, including ResumeProcedure and captured bindings, is intentionally not
// representable in this type.
type ProjectedInteraction struct {
	ID            model.InteractionID `json:"interaction_instance_id"`
	OwnerPlayerID model.PlayerID      `json:"owner_player_id"`
	Kind          string              `json:"kind"`
	Visibility    string              `json:"visibility"`
	Prompt        json.RawMessage     `json:"prompt"`
	LegalDomain   json.RawMessage     `json:"legal_domain"`
}

type PlayerProjection struct {
	MatchID              model.MatchID         `json:"match_id"`
	PlayerID             model.PlayerID        `json:"player_id"`
	Revision             uint64                `json:"revision"`
	EventSequence        uint64                `json:"event_sequence"`
	Lifecycle            model.Lifecycle       `json:"lifecycle"`
	View                 json.RawMessage       `json:"view"`
	LegalActions         []json.RawMessage     `json:"legal_actions"`
	PendingInteraction   *ProjectedInteraction `json:"pending_interaction,omitempty"`
	BlockedByInteraction bool                  `json:"blocked_by_interaction"`
}

type CommandResult struct {
	MatchID    model.MatchID    `json:"match_id"`
	PlayerID   model.PlayerID   `json:"player_id"`
	Revision   uint64           `json:"revision"`
	Sequence   uint64           `json:"event_sequence"`
	Projection PlayerProjection `json:"projection"`
}

type Host struct {
	definitions coreruntime.DefinitionRegistry
	ids         coreruntime.IDProvider
	store       persistence.EventStore
	rules       contracts.RulesEngine
}

func NewHost(
	definitions coreruntime.DefinitionRegistry,
	ids coreruntime.IDProvider,
	store persistence.EventStore,
	rules contracts.RulesEngine,
) *Host {
	return &Host{definitions: definitions, ids: ids, store: store, rules: rules}
}
