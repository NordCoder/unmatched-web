package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func (h *Host) prepareStartAction(ctx context.Context, principal model.PrincipalID, command contracts.Command) (preparedCommand, error) {
	state, playerID, err := h.loadAndAuthorize(principal, command, true)
	if err != nil {
		return preparedCommand{}, err
	}
	if state.Game.Lifecycle != model.LifecycleActive {
		return preparedCommand{}, opError(CodeInvalidCommand, "match is not active")
	}
	if state.Game.Resolver.PendingInteraction != nil {
		return preparedCommand{}, opError(CodePendingInteraction, "resolve the pending interaction before starting another action")
	}
	var payload StartActionPayload
	if err := decodePayload(command.Payload, &payload); err != nil {
		return preparedCommand{}, err
	}
	if payload.Kind != ActionManeuver && payload.Kind != ActionScheme && payload.Kind != ActionAttack {
		return preparedCommand{}, opError(CodeInvalidCommand, "action kind must be MANEUVER, SCHEME, or ATTACK")
	}
	procedure := model.ProcedureRef{
		ID:        model.ProcedureID(h.ids.Next("procedure")),
		Kind:      payload.Kind,
		Stage:     "START",
		SourceRef: payload.SourceRef,
		Bindings:  cloneRawMap(payload.Context),
	}
	if procedure.Bindings == nil {
		procedure.Bindings = make(map[string]json.RawMessage)
	}
	procedure.Bindings["actor_player_id"] = mustJSON(playerID)
	action := model.ActionState{
		ID:      model.ActionID(h.ids.Next("action")),
		Kind:    payload.Kind,
		ActorID: playerID,
		Status:  "RESOLVING",
	}
	started := contracts.DomainEvent{
		Type: coreruntime.EventActionStarted,
		PublicPayload: mustJSON(coreruntime.ActionStartedPublicPayload{
			ActionID: action.ID, Kind: action.Kind, ActorID: action.ActorID, Status: action.Status,
		}),
		PrivatePayloads: map[model.PlayerID]json.RawMessage{
			playerID: mustJSON(coreruntime.ActionStartedPayload{Action: action, Procedure: procedure}),
		},
	}
	outcome, err := h.rules.Resolve(state.Game, contracts.ResolutionInput{
		CommandID: command.ID,
		Procedure: procedure,
		Context:   cloneRawMap(payload.Context),
	})
	if err != nil {
		return preparedCommand{}, internalError("resolve action", err)
	}
	events, err := h.eventsForOutcome(action.ID, procedure, outcome)
	if err != nil {
		return preparedCommand{}, err
	}
	return preparedCommand{
		matchID: command.MatchID, playerID: playerID, previous: state,
		ruleset: state.Game.DefinitionRef.RulesetVersion,
		events:  append([]contracts.DomainEvent{started}, events...),
	}, nil
}

func (h *Host) prepareSubmitChoice(ctx context.Context, principal model.PrincipalID, command contracts.Command) (preparedCommand, error) {
	state, playerID, err := h.loadAndAuthorize(principal, command, true)
	if err != nil {
		return preparedCommand{}, err
	}
	var payload SubmitChoicePayload
	if err := decodePayload(command.Payload, &payload); err != nil {
		return preparedCommand{}, err
	}
	pending := state.Game.Resolver.PendingInteraction
	if pending == nil {
		return preparedCommand{}, opError(CodeInvalidChoice, "match has no pending interaction")
	}
	if pending.ID != payload.InteractionID {
		return preparedCommand{}, opError(CodeInvalidChoice, "interaction ID does not match the pending interaction")
	}
	if pending.OwnerPlayerID != playerID {
		return preparedCommand{}, opError(CodeUnauthorized, "pending interaction belongs to another player")
	}
	if len(payload.Choice) == 0 || !json.Valid(payload.Choice) {
		return preparedCommand{}, opError(CodeInvalidChoice, "choice must be valid JSON")
	}
	procedure := pending.ResumeProcedure
	outcome, err := h.rules.Resolve(state.Game, contracts.ResolutionInput{
		CommandID:     command.ID,
		Procedure:     procedure,
		InteractionID: pending.ID,
		Choice:        append(json.RawMessage(nil), payload.Choice...),
	})
	if err != nil {
		return preparedCommand{}, internalError("resume action", err)
	}
	if state.Game.Action == nil {
		return preparedCommand{}, internalError("resume action", errors.New("pending interaction has no active action"))
	}
	events, err := h.eventsForOutcome(state.Game.Action.ID, procedure, outcome)
	if err != nil {
		return preparedCommand{}, err
	}
	return preparedCommand{
		matchID: command.MatchID, playerID: playerID, previous: state,
		ruleset: state.Game.DefinitionRef.RulesetVersion, events: events,
	}, nil
}

func (h *Host) eventsForOutcome(actionID model.ActionID, procedure model.ProcedureRef, outcome contracts.ResolutionOutcome) ([]contracts.DomainEvent, error) {
	switch outcome.Status {
	case contracts.ResolutionPending:
		if outcome.PendingInteraction == nil {
			return nil, internalError("resolve action", errors.New("pending outcome omitted interaction"))
		}
		interaction := *outcome.PendingInteraction
		if interaction.ResumeProcedure.ID == "" {
			interaction.ResumeProcedure = procedure
		}
		if interaction.ResumeProcedure.ID != procedure.ID {
			return nil, internalError("resolve action", errors.New("pending outcome changed procedure identity"))
		}
		return []contracts.DomainEvent{{
			Type: coreruntime.EventInteractionPending,
			PublicPayload: mustJSON(coreruntime.InteractionPendingPublicPayload{
				InteractionID: interaction.ID, OwnerPlayerID: interaction.OwnerPlayerID,
				Kind: interaction.Kind, Visibility: interaction.Visibility,
			}),
			PrivatePayloads: map[model.PlayerID]json.RawMessage{
				interaction.OwnerPlayerID: mustJSON(coreruntime.InteractionPendingPayload{Interaction: interaction}),
			},
		}}, nil
	case contracts.ResolutionCompleted:
		events := append([]contracts.DomainEvent(nil), outcome.Events...)
		events = append(events, contracts.DomainEvent{
			Type:          coreruntime.EventActionCompleted,
			PublicPayload: mustJSON(coreruntime.ActionCompletedPayload{ActionID: actionID}),
		})
		return events, nil
	case contracts.ResolutionRejected:
		code := outcome.RejectionCode
		if code == "" {
			code = "unspecified"
		}
		return nil, opError(CodeRulesRejected, code)
	default:
		return nil, internalError("resolve action", fmt.Errorf("unknown resolution status %q", outcome.Status))
	}
}
