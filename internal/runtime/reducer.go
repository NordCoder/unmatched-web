package runtime

import (
	"encoding/json"
	"fmt"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

const (
	EventMatchCreated       = "core.match.created"
	EventPlayerJoined       = "core.player.joined"
	EventActionStarted      = "core.action.started"
	EventInteractionPending = "core.interaction.pending"
	EventActionCompleted    = "core.action.completed"
)

type AuthorityBinding struct {
	PrincipalID model.PrincipalID `json:"principal_id"`
	PlayerID    model.PlayerID    `json:"player_id"`
}

type HostedState struct {
	Game        model.GameState
	Authorities map[model.PrincipalID]model.PlayerID
}

type MatchCreatedPayload struct {
	MatchID       model.MatchID       `json:"match_id"`
	DefinitionRef model.DefinitionRef `json:"definition_ref"`
	Lifecycle     model.Lifecycle     `json:"lifecycle"`
	PlayerID      model.PlayerID      `json:"player_id"`
}

type MatchCreatedPrivatePayload struct {
	State     model.GameState  `json:"state"`
	Authority AuthorityBinding `json:"authority"`
}

type PlayerJoinedPayload struct {
	PlayerID           model.PlayerID    `json:"player_id"`
	Seat               int               `json:"seat"`
	FighterInstanceIDs []model.FighterID `json:"fighter_instance_ids"`
	Lifecycle          model.Lifecycle   `json:"lifecycle"`
}

type PlayerJoinedPrivatePayload struct {
	Player    model.PlayerState                       `json:"player"`
	Fighters  map[model.FighterID]model.RuntimeObject `json:"fighters"`
	Cards     map[model.CardID]model.RuntimeObject    `json:"cards"`
	Lifecycle model.Lifecycle                         `json:"lifecycle"`
	Authority AuthorityBinding                        `json:"authority"`
}

type ActionStartedPublicPayload struct {
	ActionID model.ActionID `json:"action_id"`
	Kind     string         `json:"kind"`
	ActorID  model.PlayerID `json:"actor_player_id"`
	Status   string         `json:"status"`
}

type ActionStartedPayload struct {
	Action    model.ActionState  `json:"action"`
	Procedure model.ProcedureRef `json:"procedure"`
}

type InteractionPendingPublicPayload struct {
	InteractionID model.InteractionID `json:"interaction_id"`
	OwnerPlayerID model.PlayerID      `json:"owner_player_id"`
	Kind          string              `json:"kind"`
	Visibility    string              `json:"visibility"`
}

type InteractionPendingPayload struct {
	Interaction model.PendingInteraction `json:"interaction"`
}

type ActionCompletedPayload struct {
	ActionID model.ActionID `json:"action_id"`
}

func Replay(events []contracts.DomainEvent) (HostedState, error) {
	state := HostedState{Authorities: make(map[model.PrincipalID]model.PlayerID)}
	for _, event := range events {
		var err error
		state, err = Apply(state, event)
		if err != nil {
			return HostedState{}, err
		}
	}
	return state, nil
}

func Apply(state HostedState, event contracts.DomainEvent) (HostedState, error) {
	if state.Authorities == nil {
		state.Authorities = make(map[model.PrincipalID]model.PlayerID)
	}

	switch event.Type {
	case EventMatchCreated:
		var public MatchCreatedPayload
		if err := json.Unmarshal(event.PublicPayload, &public); err != nil {
			return HostedState{}, fmt.Errorf("decode match-created public event: %w", err)
		}
		var payload MatchCreatedPrivatePayload
		if err := decodePrivate(event.PrivatePayloads, public.PlayerID, &payload); err != nil {
			return HostedState{}, fmt.Errorf("decode match-created private event: %w", err)
		}
		state.Game = payload.State
		state.Game.DefinitionRef = cloneDefinitionRef(state.Game.DefinitionRef)
		state.Game.Players = clonePlayers(state.Game.Players)
		state.Game.Fighters = cloneRuntimeObjects(state.Game.Fighters)
		state.Game.Cards = cloneCardObjects(state.Game.Cards)
		state.Authorities[payload.Authority.PrincipalID] = payload.Authority.PlayerID
	case EventPlayerJoined:
		var public PlayerJoinedPayload
		if err := json.Unmarshal(event.PublicPayload, &public); err != nil {
			return HostedState{}, fmt.Errorf("decode player-joined public event: %w", err)
		}
		var payload PlayerJoinedPrivatePayload
		if err := decodePrivate(event.PrivatePayloads, public.PlayerID, &payload); err != nil {
			return HostedState{}, fmt.Errorf("decode player-joined private event: %w", err)
		}
		if state.Game.Players == nil {
			state.Game.Players = make(map[model.PlayerID]model.PlayerState)
		}
		state.Game.Players[payload.Player.ID] = clonePlayer(payload.Player)
		mergeRuntimeObjects(state.Game.Fighters, payload.Fighters)
		mergeRuntimeObjects(state.Game.Cards, payload.Cards)
		state.Game.Lifecycle = payload.Lifecycle
		state.Authorities[payload.Authority.PrincipalID] = payload.Authority.PlayerID
	case EventActionStarted:
		var public ActionStartedPublicPayload
		if err := json.Unmarshal(event.PublicPayload, &public); err != nil {
			return HostedState{}, fmt.Errorf("decode action-started public event: %w", err)
		}
		var payload ActionStartedPayload
		if err := decodePrivate(event.PrivatePayloads, public.ActorID, &payload); err != nil {
			return HostedState{}, fmt.Errorf("decode action-started private event: %w", err)
		}
		action := payload.Action
		procedure := cloneProcedure(payload.Procedure)
		action.Procedure = &procedure
		state.Game.Action = &action
		state.Game.Resolver.ActiveProcedure = &procedure
		state.Game.Resolver.PendingInteraction = nil
	case EventInteractionPending:
		var public InteractionPendingPublicPayload
		if err := json.Unmarshal(event.PublicPayload, &public); err != nil {
			return HostedState{}, fmt.Errorf("decode interaction-pending public event: %w", err)
		}
		var payload InteractionPendingPayload
		if err := decodePrivate(event.PrivatePayloads, public.OwnerPlayerID, &payload); err != nil {
			return HostedState{}, fmt.Errorf("decode interaction-pending private event: %w", err)
		}
		interaction := cloneInteraction(payload.Interaction)
		procedure := cloneProcedure(interaction.ResumeProcedure)
		state.Game.Resolver.PendingInteraction = &interaction
		state.Game.Resolver.ActiveProcedure = &procedure
		if state.Game.Action != nil {
			state.Game.Action.Status = "WAITING_FOR_CHOICE"
			state.Game.Action.Procedure = &procedure
		}
	case EventActionCompleted:
		var payload ActionCompletedPayload
		if err := json.Unmarshal(event.PublicPayload, &payload); err != nil {
			return HostedState{}, fmt.Errorf("decode action-completed event: %w", err)
		}
		if state.Game.Action != nil && state.Game.Action.ID == payload.ActionID {
			state.Game.Action.Status = "COMPLETED"
			state.Game.Action.Procedure = nil
		}
		state.Game.Resolver.ActiveProcedure = nil
		state.Game.Resolver.PendingInteraction = nil
	default:
		// Rules events are authoritative history. Core advances their ordered
		// envelope without interpreting fighter/card semantics.
		state.Game.Resolver.History = append(state.Game.Resolver.History, event.Type)
	}

	state.Game.Revision = event.Revision
	state.Game.EventSequence = event.Sequence
	return state, nil
}

func decodePrivate(payloads map[model.PlayerID]json.RawMessage, playerID model.PlayerID, target any) error {
	raw, ok := payloads[playerID]
	if !ok {
		return fmt.Errorf("private payload for player %q is missing", playerID)
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return err
	}
	return nil
}

func clonePlayers(source map[model.PlayerID]model.PlayerState) map[model.PlayerID]model.PlayerState {
	if source == nil {
		return make(map[model.PlayerID]model.PlayerState)
	}
	result := make(map[model.PlayerID]model.PlayerState, len(source))
	for id, player := range source {
		result[id] = clonePlayer(player)
	}
	return result
}

func clonePlayer(player model.PlayerState) model.PlayerState {
	player.FighterInstanceIDs = append([]model.FighterID(nil), player.FighterInstanceIDs...)
	player.ActionPermissionIDs = append([]string(nil), player.ActionPermissionIDs...)
	player.SubmittedChoiceIDs = append([]string(nil), player.SubmittedChoiceIDs...)
	if player.PrivateZones != nil {
		zones := make(map[string][]model.CardID, len(player.PrivateZones))
		for key, cards := range player.PrivateZones {
			zones[key] = append([]model.CardID(nil), cards...)
		}
		player.PrivateZones = zones
	}
	if player.Resources != nil {
		resources := make(map[string]int, len(player.Resources))
		for key, value := range player.Resources {
			resources[key] = value
		}
		player.Resources = resources
	}
	return player
}

func cloneRuntimeObjects(source map[model.FighterID]model.RuntimeObject) map[model.FighterID]model.RuntimeObject {
	result := make(map[model.FighterID]model.RuntimeObject, len(source))
	for id, object := range source {
		result[id] = cloneRuntimeObject(object)
	}
	return result
}

func cloneCardObjects(source map[model.CardID]model.RuntimeObject) map[model.CardID]model.RuntimeObject {
	result := make(map[model.CardID]model.RuntimeObject, len(source))
	for id, object := range source {
		result[id] = cloneRuntimeObject(object)
	}
	return result
}

func mergeRuntimeObjects[K ~string](target map[K]model.RuntimeObject, source map[K]model.RuntimeObject) {
	for id, object := range source {
		target[id] = cloneRuntimeObject(object)
	}
}

func cloneRuntimeObject(object model.RuntimeObject) model.RuntimeObject {
	if object.State != nil {
		state := make(map[string]any, len(object.State))
		for key, value := range object.State {
			state[key] = value
		}
		object.State = state
	}
	return object
}

func cloneProcedure(procedure model.ProcedureRef) model.ProcedureRef {
	if procedure.Bindings != nil {
		bindings := make(map[string]json.RawMessage, len(procedure.Bindings))
		for key, value := range procedure.Bindings {
			bindings[key] = append(json.RawMessage(nil), value...)
		}
		procedure.Bindings = bindings
	}
	return procedure
}

func cloneInteraction(interaction model.PendingInteraction) model.PendingInteraction {
	interaction.Prompt = append(json.RawMessage(nil), interaction.Prompt...)
	interaction.LegalDomain = append(json.RawMessage(nil), interaction.LegalDomain...)
	interaction.ResumeProcedure = cloneProcedure(interaction.ResumeProcedure)
	return interaction
}
