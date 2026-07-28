package application

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

type preparedCommand struct {
	matchID   model.MatchID
	playerID  model.PlayerID
	previous  coreruntime.HostedState
	ruleset   string
	events    []contracts.DomainEvent
	authority *persistence.AuthorityRecord
}

func (h *Host) Execute(ctx context.Context, principal model.PrincipalID, command contracts.Command) (result CommandResult, err error) {
	if err := validateEnvelope(principal, command); err != nil {
		return CommandResult{}, err
	}
	fingerprint, err := commandFingerprint(principal, command)
	if err != nil {
		return CommandResult{}, internalError("fingerprint command", err)
	}
	key := persistence.CommandKey{PrincipalID: principal, CommandID: command.ID}
	lease, existing, duplicate, err := h.store.AcquireCommand(ctx, key, fingerprint)
	if err != nil {
		if errors.Is(err, persistence.ErrCommandConflict) {
			return CommandResult{}, opError(CodeCommandConflict, "command ID was already used with different input")
		}
		return CommandResult{}, internalError("acquire command identity", err)
	}
	if duplicate {
		return decodeCommandRecord(existing)
	}

	finished := false
	defer func() {
		if !finished {
			h.store.AbortCommand(lease)
		}
	}()

	reject := func(cause error) (CommandResult, error) {
		var operational *OperationalError
		if !errors.As(cause, &operational) || operational.Code == CodeInternal {
			return CommandResult{}, cause
		}
		if _, rejectErr := h.store.RejectCommand(lease, command.MatchID, operational.Code, operational.Message); rejectErr != nil {
			return CommandResult{}, internalError("persist command rejection", rejectErr)
		}
		finished = true
		return CommandResult{}, cause
	}

	prepared, err := h.prepare(ctx, principal, command)
	if err != nil {
		return reject(err)
	}
	batch, state, err := h.buildBatch(command, prepared)
	if err != nil {
		return reject(err)
	}
	projection, err := h.project(state, prepared.playerID)
	if err != nil {
		return reject(err)
	}
	result = CommandResult{
		MatchID: prepared.matchID, PlayerID: prepared.playerID,
		Revision: state.Game.Revision, Sequence: state.Game.EventSequence,
		Projection: projection,
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return CommandResult{}, internalError("encode command result", err)
	}
	_, err = h.store.Commit(lease, persistence.CommitRequest{
		MatchID: prepared.matchID, Batch: batch, Result: encoded, Authority: prepared.authority,
	})
	if err != nil {
		switch {
		case errors.Is(err, persistence.ErrRevisionConflict):
			return reject(opError(CodeRevisionConflict, err.Error()))
		case errors.Is(err, persistence.ErrAuthorityConflict):
			return reject(opError(CodeUnauthorized, err.Error()))
		default:
			return CommandResult{}, internalError("commit event batch", err)
		}
	}
	finished = true
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
	playerID, ok := h.store.ResolveAuthority(matchID, principal)
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

func decodeCommandRecord(record persistence.CommandRecord) (CommandResult, error) {
	if record.ErrorCode != "" {
		return CommandResult{}, &OperationalError{Code: record.ErrorCode, Message: record.ErrorMessage}
	}
	return decodeResult(record.Result)
}
