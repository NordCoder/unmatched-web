package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func (h *Host) buildBatch(command contracts.Command, prepared preparedCommand) (contracts.EventBatch, coreruntime.HostedState, error) {
	previousRevision := prepared.previous.Game.Revision
	previousSequence := prepared.previous.Game.EventSequence
	nextRevision := previousRevision + 1
	oldIDs := make(map[model.EventID]model.EventID)
	for index := range prepared.events {
		oldID := prepared.events[index].ID
		newID := model.EventID(h.ids.Next("event"))
		prepared.events[index].ID = newID
		if oldID != "" {
			oldIDs[oldID] = newID
		}
	}
	state := prepared.previous
	for index := range prepared.events {
		event := &prepared.events[index]
		event.SchemaVersion = "core-event/v1"
		event.MatchID = prepared.matchID
		event.Sequence = previousSequence + uint64(index) + 1
		event.Revision = nextRevision
		event.CausedByCommand = command.ID
		event.RulesetVersion = prepared.ruleset
		if replacement, ok := oldIDs[event.ParentEventID]; ok {
			event.ParentEventID = replacement
		}
		if event.PublicPayload == nil {
			event.PublicPayload = json.RawMessage(`{}`)
		}
		var err error
		state, err = coreruntime.Apply(state, *event)
		if err != nil {
			return contracts.EventBatch{}, coreruntime.HostedState{}, internalError("apply prepared event", err)
		}
	}
	return contracts.EventBatch{
		CommandID: command.ID, PreviousRevision: previousRevision,
		NextRevision: nextRevision, Events: prepared.events,
	}, state, nil
}

func (h *Host) loadAndAuthorize(ctx context.Context, lease persistence.CommandLease, principal model.PrincipalID, command contracts.Command, requireActor bool) (coreruntime.HostedState, model.PlayerID, error) {
	if command.MatchID == "" {
		return coreruntime.HostedState{}, "", opError(CodeInvalidCommand, "match ID is required")
	}
	state, err := h.stateForCommand(ctx, lease, command.MatchID)
	if err != nil {
		return coreruntime.HostedState{}, "", err
	}
	if err := checkExpectedRevision(command, state.Game.Revision); err != nil {
		return coreruntime.HostedState{}, "", err
	}
	playerID, ok, err := h.resolveAuthorityForCommand(ctx, lease, command.MatchID, principal)
	if err != nil {
		return state, "", internalError("resolve principal authority", err)
	}
	if !ok {
		return state, "", opError(CodeUnauthorized, "principal is not bound to this match")
	}
	if requireActor && command.ActorPlayerID == "" {
		return coreruntime.HostedState{}, "", opError(CodeInvalidCommand, "actor player ID is required")
	}
	if command.ActorPlayerID != "" && command.ActorPlayerID != playerID {
		return coreruntime.HostedState{}, "", opError(CodeUnauthorized, "actor player ID does not match principal authority")
	}
	return state, playerID, nil
}

func (h *Host) project(state coreruntime.HostedState, playerID model.PlayerID) (PlayerProjection, error) {
	view, err := h.rules.Project(state.Game, playerID)
	if err != nil {
		return PlayerProjection{}, internalError("project player state", err)
	}
	legalActions, err := h.rules.LegalActions(state.Game, playerID)
	if err != nil {
		return PlayerProjection{}, internalError("query legal actions", err)
	}
	projection := PlayerProjection{
		MatchID: state.Game.MatchID, PlayerID: playerID,
		Revision: state.Game.Revision, EventSequence: state.Game.EventSequence,
		Lifecycle: state.Game.Lifecycle, View: append(json.RawMessage(nil), view...),
		LegalActions: cloneRawSlice(legalActions),
	}
	if pending := state.Game.Resolver.PendingInteraction; pending != nil {
		projection.BlockedByInteraction = true
		if pending.OwnerPlayerID == playerID {
			projection.PendingInteraction = &ProjectedInteraction{
				ID: pending.ID, OwnerPlayerID: pending.OwnerPlayerID,
				Kind: pending.Kind, Visibility: pending.Visibility,
				Prompt:      append(json.RawMessage(nil), pending.Prompt...),
				LegalDomain: append(json.RawMessage(nil), pending.LegalDomain...),
			}
		}
	}
	return projection, nil
}

func (h *Host) instantiatePlayer(seat int, fighter coreruntime.FighterDefinition) (model.PlayerState, map[model.FighterID]model.RuntimeObject, map[model.CardID]model.RuntimeObject) {
	playerID := model.PlayerID(h.ids.Next("player"))
	fighterID := model.FighterID(h.ids.Next("fighter"))
	fighters := map[model.FighterID]model.RuntimeObject{
		fighterID: {DefinitionID: fighter.ID, OwnerID: playerID, ControllerID: playerID, State: make(map[string]any)},
	}
	cards := make(map[model.CardID]model.RuntimeObject, len(fighter.CardDefinitions))
	deck := make([]model.CardID, 0, len(fighter.CardDefinitions))
	for _, definitionID := range fighter.CardDefinitions {
		cardID := model.CardID(h.ids.Next("card"))
		deck = append(deck, cardID)
		cards[cardID] = model.RuntimeObject{DefinitionID: definitionID, OwnerID: playerID, ControllerID: playerID, State: make(map[string]any)}
	}
	player := model.PlayerState{
		ID: playerID, Seat: seat, AuthorityState: "ACTIVE",
		FighterInstanceIDs: []model.FighterID{fighterID},
		PrivateZones:       map[string][]model.CardID{"deck": deck},
		Resources:          make(map[string]int),
	}
	return player, fighters, cards
}

func (h *Host) findPinnedBundle(ref model.DefinitionRef) (coreruntime.DefinitionBundle, bool) {
	return h.definitions.ResolveRef(ref)
}

func validateEnvelope(principal model.PrincipalID, command contracts.Command) error {
	if principal == "" {
		return opError(CodeUnauthorized, "principal ID is required")
	}
	if command.ID == "" || command.SchemaVersion != CommandSchemaV1 || command.Type == "" {
		return opError(CodeInvalidCommand, "command ID, schema version, and type are required")
	}
	if len(command.Payload) == 0 || !json.Valid(command.Payload) {
		return opError(CoeInvalidCommand, "command payload must be valid JSON")
	}
	return nil
}

func checkExpectedRevision(command contracts.Command, current uint64) error {
	if command.ExpectedRevision == nil {
		return opError(CodeInvalidCommand, "expected revision is required")
	}
	if *command.ExpectedRevision != current {
		return opError(CodeRevisionConflict, fmt.Sprintf("expected revision %d, current %d", *command.ExpectedRevision, current))
	}
	return nil
}

func decodePayload(raw json.RawMessage, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return opError(CoeInvalidCommand, "command payload does not match its type")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return opError(CoeInvalidCommand, "command payload contains trailing data")
	}
	return nil
}

func decodeResult(raw []byte) (CommandResult, error) {
	var result CommandResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return CommandResult{}, internalError("decode stored command result", err)
	}
	return result, nil
}

func cloneRawMap(source map[string]json.RawMessage) map[string]json.RawMessage {
	if source == nil {
		return nil
	}
	result := make(map[string]json.RawMessage, len(source))
	for key, value := range source {
		result[key] = append(json.RawMessage(nil), value...)
	}
	return result
}

func cloneRawSlice(source []json.RawMessage) []json.RawMessage {
	result := make([]json.RawMessage, len(source))
	for index, value := range source {
		result[index] = append(json.RawMessage(nil), value...)
	}
	return result
}

func mustJSON(value any) json.RawMessage {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return encoded
}

func opError(code, message string) error { return &OperationalError{Code: code, Message: message} }

func internalError(operation string, err error) error {
	return &OperationalError{Code: CodeInternal, Message: operation + ": " + err.Error()}
}
