package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

const commandCleanupTimeout = 5 * time.Second

type preparedCommand struct {
	matchID   model.MatchID
	playerID  model.PlayerID
	previous  coreruntime.HostedState
	ruleset   string
	events    []contracts.DomainEvent
	authority *persistence.AuthorityRecord
}

func (h *Host) Execute(ctx context.Context, principal model.PrincipalID, command contracts.Command) (result CommandResult, err error) {
	principal, command, fingerprint, err := normalizeCommandRequest(principal, command)
	if err != nil {
		return CommandResult{}, err
	}
	identity := persistence.CommandIdentity{
		Key:           persistence.CommandKey{PrincipalID: principal, CommandID: command.ID},
		Fingerprint:   fingerprint,
		MatchID:       command.MatchID,
		ActorPlayerID: command.ActorPlayerID,
		Scope:         commandScope(command.Type),
	}
	lease, existing, duplicate, err := h.store.AcquireCommand(ctx, identity)
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
		if finished {
			return
		}
		cleanupContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), commandCleanupTimeout)
		defer cancel()
		if abortErr := h.store.AbortCommand(cleanupContext, lease); abortErr != nil && err == nil {
			err = internalError("abort command reservation", abortErr)
		}
	}()

	reject := func(cause error) (CommandResult, error) {
		var operational *OperationalError
		if !errors.As(cause, &operational) || operational.Code == CodeInternal {
			return CommandResult{}, cause
		}
		if _, rejectErr := h.store.RejectCommand(ctx, lease, command.MatchID, operational.Code, operational.Message); rejectErr != nil {
			return CommandResult{}, internalError("persist command rejection", rejectErr)
		}
		finished = true
		return CommandResult{}, cause
	}

	prepared, err := h.prepare(ctx, lease, principal, command)
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
	_, err = h.store.Commit(ctx, lease, persistence.CommitRequest{
		MatchID:       prepared.matchID,
		DefinitionRef: state.Game.DefinitionRef,
		Batch:         batch,
		Result:        encoded,
		Authority:     prepared.authority,
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
	return h.StateContext(context.Background(), matchID)
}

func (h *Host) StateContext(ctx context.Context, matchID model.MatchID) (coreruntime.HostedState, error) {
	events, err := h.store.LoadEvents(ctx, matchID)
	if err != nil {
		return coreruntime.HostedState{}, internalError("load match events", err)
	}
	return replayStoredState(events)
}

func (h *Host) stateForCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID) (coreruntime.HostedState, error) {
	events, err := h.store.LoadEventsForCommand(ctx, lease, matchID)
	if err != nil {
		return coreruntime.HostedState{}, internalError("load match events", err)
	}
	return replayStoredState(events)
}

func replayStoredState(events []contracts.DomainEvent) (coreruntime.HostedState, error) {
	if len(events) == 0 {
		return coreruntime.HostedState{}, opError(CodeNotFound, "match was not found")
	}
	state, err := coreruntime.Replay(events)
	if err != nil {
		return coreruntime.HostedState{}, internalError("replay match", err)
	}
	return state, nil
}

func (h *Host) resolveAuthorityForCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID, principal model.PrincipalID) (model.PlayerID, bool, error) {
	return h.store.ResolveAuthorityForCommand(ctx, lease, matchID, principal)
}

func (h *Host) Project(matchID model.MatchID, principal model.PrincipalID) (PlayerProjection, error) {
	return h.ProjectContext(context.Background(), matchID, principal)
}

func (h *Host) ProjectContext(ctx context.Context, matchID model.MatchID, principal model.PrincipalID) (PlayerProjection, error) {
	state, err := h.StateContext(ctx, matchID)
	if err != nil {
		return PlayerProjection{}, err
	}
	playerID, ok, err := h.store.ResolveAuthorityContext(ctx, matchID, principal)
	if err != nil {
		return PlayerProjection{}, internalError("resolve principal authority", err)
	}
	if !ok {
		return PlayerProjection{}, opError(CodeUnauthorized, "principal is not bound to this match")
	}
	return h.project(state, playerID)
}

func (h *Host) prepare(ctx context.Context, lease persistence.CommandLease, principal model.PrincipalID, command contracts.Command) (preparedCommand, error) {
	switch command.Type {
	case CommandCreateMatch:
		return h.prepareCreate(principal, command)
	case CommandJoinMatch:
		return h.prepareJoin(ctx, lease, principal, command)
	case CommandStartAction:
		return h.prepareStartAction(ctx, lease, principal, command)
	case CommandSubmitChoice:
		return h.prepareSubmitChoice(ctx, lease, principal, command)
	default:
		return preparedCommand{}, opError(CodeInvalidCommand, "unsupported command type")
	}
}

func commandScope(commandType string) persistence.CommandScope {
	switch commandType {
	case CommandCreateMatch:
		return persistence.CommandScopeCreateMatch
	case CommandJoinMatch:
		return persistence.CommandScopeJoinMatch
	default:
		return persistence.CommandScopeExistingSeat
	}
}

func decodeCommandRecord(record persistence.CommandRecord) (CommandResult, error) {
	if record.ErrorCode != "" {
		return CommandResult{}, &OperationalError{Code: record.ErrorCode, Message: record.ErrorMessage}
	}
	return decodeResult(record.Result)
}
