package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

type preparedCommand struct {
	matchID  model.MatchID
	playerID model.PlayerID
	previous coreruntime.HostedState
	ruleset  string
	events   []contracts.DomainEvent
}

func (h *Host) Execute(ctx context.Context, principal model.PrincipalID, command contracts.Command) (CommandResult, error) {
	if err := validateEnvelope(principal, command); err != nil {
		return CommandResult{}, err
	}
	fingerprint, err := commandFingerprint(principal, command)
	if err != nil {
		return CommandResult{}, internalError("fingerprint command", err)
	}
	if existing, ok := h.store.LookupCommand(command.ID); ok {
		if !bytes.Equal(existing.Fingerprint, fingerprint) {
			return CommandResult{}, opError(CodeCommandConflict, "command ID was already used with different input")
		}
		return decodeResult(existing.Result)
	}

	prepared, err := h.prepare(ctx, principal, command)
	if err != nil {
		return CommandResult{}, err
	}
	batch, state, err := h.buildBatch(command, prepared)
	if err != nil {
		return CommandResult{}, err
	}
	projection, err := h.project(state, prepared.playerID)
	if err != nil {
		return CommandResult{}, err
	}
	result := CommandResult{
		MatchID:    prepared.matchID,
		PlayerID:   prepared.playerID,
		Revision:   state.Game.Revision,
		Sequence:   state.Game.EventSequence,
		Projection: projection,
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return CommandResult{}, internalError("encode command result", err)
	}
	record, duplicate, err := h.store.Commit(persistence.CommitRequest{
		MatchID: prepared.matchID, Fingerprint: fingerprint, Batch: batch, Result: encoded,
	})
	if err != nil {
		switch {
		case errors.Is(err, persistence.ErrRevisionConflict):
			return CommandResult{}, opError(CodeRevisionConflict, err.Error())
		case errors.Is(err, persistence.ErrCommandConflict):
			return CommandResult{}, opError(CodeCommandConflict, err.Error())
		default:
			return CommandResult{}, internalError("commit event batch", err)
		}
	}
	if duplicate {
		return decodeResult(record.Result)
	}
	return result, nil
}

func (h *Host) State(matchID model.MatchID) (coreruntime.HostedState, error) {
	events := h.store.Events(matchID)
	if len(events) == 0 {
		return coreruntime.HostedState{}, opError(CodeNotFound, "match was not found")
	}
	state, err := coreruntime.Replay(events)
	if err != nil {
		return coreruntime.HostedState{}, internalError("replay match", err)
	}
	return state, nil
}

func (h *Host) Project(matchID model.MatchID, principal model.PrincipalID) (PlayerProjection, error) {
	state, err := h.State(matchID)
	if err != nil {
		return PlayerProjection{}, err
	}
	playerID, ok := state.Authorities[principal]
	if !ok {
		return PlayerProjection{}, opError(CodeUnauthorized, "principal is not bound to this match")
	}
	return h.project(state, playerID)
}

func (h *Host) prepare(ctx context.Context, principal model.PrincipalID, command contracts.Command) (preparedCommand, error) {
	switch command.Type {
	case CommandCreateMatch:
		return h.prepareCreate(principal, command)
	case CommandJoinMatch:
		return h.prepareJoin(principal, command)
	case CommandStartAction:
		return h.prepareStartAction(ctx, principal, command)
	case CommandSubmitChoice:
		return h.prepareSubmitChoice(ctx, principal, command)
	default:
		return preparedCommand{}, opError(CodeInvalidCommand, "unsupported command type")
	}
}
