package application

import (
	"context"
	"encoding/json"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func (h *Host) prepareCreate(principal model.PrincipalID, command contracts.Command) (preparedCommand, error) {
	if command.MatchID != "" || command.ActorPlayerID != "" {
		return preparedCommand{}, opError(CodeInvalidCommand, "CreateMatch cannot supply match or actor IDs")
	}
	if command.ExpectedRevision != nil && *command.ExpectedRevision != 0 {
		return preparedCommand{}, opError(CodeRevisionConflict, "CreateMatch expected revision must be zero")
	}
	var payload CreateMatchPayload
	if err := decodePayload(command.Payload, &payload); err != nil {
		return preparedCommand{}, err
	}
	bundle, ok := h.definitions.Resolve(payload.DefinitionKey)
	if !ok {
		return preparedCommand{}, opError(CodeNotFound, "definition bundle is not registered")
	}
	fighter, ok := bundle.Fighters[payload.FighterDefinition]
	if !ok {
		return preparedCommand{}, opError(CodeNotFound, "fighter definition is not registered in the bundle")
	}

	matchID := model.MatchID(h.ids.Next("match"))
	player, fighters, cards := h.instantiatePlayer(1, fighter)
	state := model.GameState{
		MatchID: matchID, DefinitionRef: bundle.Ref,
		Lifecycle: model.LifecycleWaitingForPlayers,
		Players:   map[model.PlayerID]model.PlayerState{player.ID: player},
		Fighters:  fighters, Cards: cards,
		Components:  make(map[model.ComponentID]model.RuntimeObject),
		Battlefield: make(map[string]any), Turn: make(map[string]any),
		Resolver: model.ResolverState{}, Random: make(map[string]any), GameResult: make(map[string]any),
	}
	public := mustJSON(coreruntime.MatchCreatedPayload{
		MatchID: matchID, DefinitionRef: bundle.Ref, Lifecycle: state.Lifecycle, PlayerID: player.ID,
	})
	private := mustJSON(coreruntime.MatchCreatedPrivatePayload{State: state})
	return preparedCommand{
		matchID: matchID, playerID: player.ID, ruleset: bundle.Ref.RulesetVersion,
		events: []contracts.DomainEvent{{
			Type: coreruntime.EventMatchCreated, PublicPayload: public,
			PrivatePayloads: map[model.PlayerID]json.RawMessage{player.ID: private},
		}},
		authority: &persistence.AuthorityRecord{
			MatchID: matchID, PlayerID: player.ID, PrincipalID: principal, Seat: player.Seat,
			BindingVersion: 1, Status: persistence.AuthorityActive, EstablishedByCommandID: command.ID,
		},
	}, nil
}

func (h *Host) prepareJoin(ctx context.Context, principal model.PrincipalID, command contracts.Command) (preparedCommand, error) {
	if command.MatchID == "" {
		return preparedCommand{}, opError(CodeInvalidCommand, "match ID is required")
	}
	if command.ActorPlayerID != "" {
		return preparedCommand{}, opError(CodeInvalidCommand, "JoinMatch cannot supply an actor player ID")
	}
	if _, alreadyBound, err := h.store.ResolveAuthorityContext(ctx, command.MatchID, principal); err != nil {
		return preparedCommand{}, internalError("resolve join authority", err)
	} else if alreadyBound {
		return preparedCommand{}, opError(CodeInvalidCommand, "principal already joined this match")
	}
	state, err := h.StateContext(ctx, command.MatchID)
	if err != nil {
		return preparedCommand{}, err
	}
	if err := checkExpectedRevision(command, state.Game.Revision); err != nil {
		return preparedCommand{}, err
	}
	if state.Game.Resolver.PendingInteraction != nil {
		return preparedCommand{}, opError(CodePendingInteraction, "match is blocked by a pending interaction")
	}
	if len(state.Game.Players) >= 2 {
		return preparedCommand{}, opError(CodeInvalidCommand, "match already has two players")
	}
	var payload JoinMatchPayload
	if err := decodePayload(command.Payload, &payload); err != nil {
		return preparedCommand{}, err
	}
	bundle, ok := h.findPinnedBundle(state.Game.DefinitionRef)
	if !ok {
		return preparedCommand{}, opError(CodeNotFound, "pinned definition bundle is no longer available")
	}
	fighter, ok := bundle.Fighters[payload.FighterDefinition]
	if !ok {
		return preparedCommand{}, opError(CodeNotFound, "fighter definition is not registered in the pinned bundle")
	}
	player, fighters, cards := h.instantiatePlayer(2, fighter)
	publicEvent := coreruntime.PlayerJoinedPayload{
		PlayerID: player.ID, Seat: player.Seat, FighterInstanceIDs: player.FighterInstanceIDs,
		Lifecycle: model.LifecycleActive,
	}
	privateEvent := coreruntime.PlayerJoinedPrivatePayload{
		Player: player, Fighters: fighters, Cards: cards, Lifecycle: model.LifecycleActive,
	}
	return preparedCommand{
		matchID: command.MatchID, playerID: player.ID, previous: state,
		ruleset: state.Game.DefinitionRef.RulesetVersion,
		events: []contracts.DomainEvent{{
			Type: coreruntime.EventPlayerJoined, PublicPayload: mustJSON(publicEvent),
			PrivatePayloads: map[model.PlayerID]json.RawMessage{player.ID: mustJSON(privateEvent)},
		}},
		authority: &persistence.AuthorityRecord{
			MatchID: command.MatchID, PlayerID: player.ID, PrincipalID: principal, Seat: player.Seat,
			BindingVersion: 1, Status: persistence.AuthorityActive, EstablishedByCommandID: command.ID,
		},
	}, nil
}
