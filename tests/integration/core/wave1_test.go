package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func TestWave1PendingChoiceResumeAndReplay(t *testing.T) {
	registry := coreruntime.NewMemoryDefinitionRegistry()
	bundle := syntheticBundle()
	if err := registry.Register("synthetic@v1", bundle); err != nil {
		t.Fatalf("register definitions: %v", err)
	}

	store := persistence.NewMemoryStore()
	rules := &fakeRulesEngine{}
	host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("i1"), store, rules)
	ctx := context.Background()

	create := command("create-1", application.CommandCreateMatch, "", "", nil, application.CreateMatchPayload{
		DefinitionKey: "synthetic@v1", FighterDefinition: "mirror",
	})
	created, err := host.Execute(ctx, "principal-a", create)
	if err != nil {
		t.Fatalf("create match: %v", err)
	}
	if created.Revision != 1 || created.Projection.Lifecycle != model.LifecycleWaitingForPlayers {
		t.Fatalf("unexpected create result: %+v", created)
	}
	if player, ok := store.ResolveAuthority(created.MatchID, "principal-a"); !ok || player != created.PlayerID {
		t.Fatalf("first authority record missing: player=%q ok=%v", player, ok)
	}
	initialEventCount := len(store.Events(created.MatchID))
	duplicate, err := host.Execute(ctx, "principal-a", create)
	if err != nil {
		t.Fatalf("repeat create command: %v", err)
	}
	if !reflect.DeepEqual(created, duplicate) || len(store.Events(created.MatchID)) != initialEventCount {
		t.Fatal("duplicate command did not return the original result without duplicate events")
	}

	one := uint64(1)
	joined, err := host.Execute(ctx, "principal-b", command(
		"join-1", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("join match: %v", err)
	}
	if joined.Revision != 2 || joined.Projection.Lifecycle != model.LifecycleActive {
		t.Fatalf("unexpected join result: %+v", joined)
	}
	if player, ok := store.ResolveAuthority(created.MatchID, "principal-b"); !ok || player != joined.PlayerID {
		t.Fatalf("second authority record missing: player=%q ok=%v", player, ok)
	}

	stateAfterJoin, err := host.State(created.MatchID)
	if err != nil {
		t.Fatalf("read joined state: %v", err)
	}
	assertDistinctInstances(t, stateAfterJoin.Game)
	for _, player := range stateAfterJoin.Game.Players {
		if player.AuthorityState != "ACTIVE" {
			t.Fatalf("gameplay authority state is not ACTIVE: %+v", player)
		}
	}
	if !reflect.DeepEqual(stateAfterJoin.Game.DefinitionRef, bundle.Ref) {
		t.Fatal("match did not pin complete definition identity")
	}

	stale := uint64(1)
	_, err = host.Execute(ctx, "principal-a", command(
		"stale-action", application.CommandStartAction, created.MatchID, created.PlayerID, &stale,
		application.StartActionPayload{Kind: application.ActionScheme, SourceRef: "card-instance"},
	))
	if application.CodeOf(err) != application.CodeRevisionConflict {
		t.Fatalf("expected revision conflict, got %v", err)
	}
	// Deterministic rejections are immutable command results.
	_, err = host.Execute(ctx, "principal-a", command(
		"stale-action", application.CommandStartAction, created.MatchID, created.PlayerID, &stale,
		application.StartActionPayload{Kind: application.ActionScheme, SourceRef: "card-instance"},
	))
	if application.CodeOf(err) != application.CodeRevisionConflict {
		t.Fatalf("expected replayed revision conflict, got %v", err)
	}

	two := uint64(2)
	started, err := host.Execute(ctx, "principal-a", command(
		"action-1", application.CommandStartAction, created.MatchID, created.PlayerID, &two,
		application.StartActionPayload{Kind: application.ActionScheme, SourceRef: "card-instance"},
	))
	if err != nil {
		t.Fatalf("start action: %v", err)
	}
	if started.Revision != 3 || started.Projection.PendingInteraction == nil {
		t.Fatalf("action did not persist a pending interaction: %+v", started)
	}
	pending := started.Projection.PendingInteraction

	otherProjection, err := host.Project(created.MatchID, "principal-b")
	if err != nil {
		t.Fatalf("project other player: %v", err)
	}
	if !otherProjection.BlockedByInteraction || otherProjection.PendingInteraction != nil {
		t.Fatal("non-owner projection leaked pending interaction details")
	}

	three := uint64(3)
	_, err = host.Execute(ctx, "principal-b", command(
		"unrelated-action", application.CommandStartAction, created.MatchID, joined.PlayerID, &three,
		application.StartActionPayload{Kind: application.ActionManeuver},
	))
	if application.CodeOf(err) != application.CodePendingInteraction {
		t.Fatalf("expected unrelated command to be blocked, got %v", err)
	}

	_, err = host.Execute(ctx, "principal-b", command(
		"choice-wrong-owner", application.CommandSubmitChoice, created.MatchID, joined.PlayerID, &three,
		application.SubmitChoicePayload{InteractionID: pending.ID, Choice: json.RawMessage(`{"option":"left"}`)},
	))
	if application.CodeOf(err) != application.CodeUnauthorized {
		t.Fatalf("expected unauthorized choice, got %v", err)
	}

	completedCommand := command(
		"choice-1", application.CommandSubmitChoice, created.MatchID, created.PlayerID, &three,
		application.SubmitChoicePayload{InteractionID: pending.ID, Choice: json.RawMessage(`{"option":"left"}`)},
	)
	completed, err := host.Execute(ctx, "principal-a", completedCommand)
	if err != nil {
		t.Fatalf("submit choice: %v", err)
	}
	if completed.Revision != 4 || completed.Projection.PendingInteraction != nil {
		t.Fatalf("choice did not complete action: %+v", completed)
	}
	if rules.startProcedure.ID == "" || !reflect.DeepEqual(rules.startProcedure, rules.resumeProcedure) {
		t.Fatalf("RulesEngine did not resume the exact serialized procedure: start=%+v resume=%+v", rules.startProcedure, rules.resumeProcedure)
	}

	for _, event := range store.Events(created.MatchID) {
		encoded, marshalErr := json.Marshal(event)
		if marshalErr != nil {
			t.Fatalf("encode event: %v", marshalErr)
		}
		for _, forbidden := range [][]byte{[]byte("principal-a"), []byte("principal-b"), []byte("principal_id")} {
			if bytes.Contains(encoded, forbidden) {
				t.Fatalf("gameplay event %q leaked authority data %q: %s", event.Type, forbidden, encoded)
			}
		}
	}

	beforeDuplicate := len(store.Events(created.MatchID))
	repeatedCompletion, err := host.Execute(ctx, "principal-a", completedCommand)
	if err != nil {
		t.Fatalf("repeat choice command: %v", err)
	}
	if !reflect.DeepEqual(completed, repeatedCompletion) || len(store.Events(created.MatchID)) != beforeDuplicate {
		t.Fatal("duplicate choice produced a second event batch")
	}

	conflictingChoice := completedCommand
	conflictingChoice.Payload = json.RawMessage(`{"interaction_id":"interaction-1","choice":{"option":"right"}}`)
	_, err = host.Execute(ctx, "principal-a", conflictingChoice)
	if application.CodeOf(err) != application.CodeCommandConflict {
		t.Fatalf("expected conflicting command reuse to fail, got %v", err)
	}

	live, err := host.State(created.MatchID)
	if err != nil {
		t.Fatalf("load live state: %v", err)
	}
	replayed, err := coreruntime.Replay(store.Events(created.MatchID))
	if err != nil {
		t.Fatalf("full replay: %v", err)
	}
	if !reflect.DeepEqual(live, replayed) {
		t.Fatalf("replay diverged\nlive: %#v\nreplayed: %#v", live, replayed)
	}
	if live.Game.EventSequence != uint64(len(store.Events(created.MatchID))) {
		t.Fatalf("sequence %d does not match durable event count", live.Game.EventSequence)
	}
	if len(live.Game.Resolver.History) != 1 || live.Game.Resolver.History[0] != "rules.synthetic.choice_accepted" {
		t.Fatalf("Rules event was not retained as ordered authoritative history: %+v", live.Game.Resolver.History)
	}
}

func syntheticBundle() coreruntime.DefinitionBundle {
	fighter := coreruntime.FighterDefinition{
		ID: "mirror", CardDefinitions: []model.DefinitionID{"strike", "strike", "scheme"},
	}
	return coreruntime.DefinitionBundle{
		Ref: model.DefinitionRef{
			RulesetVersion: "synthetic-rules/v1", CapabilityRegistry: "synthetic-capabilities/v1",
			FighterManifestDigests: map[string]string{"mirror": coreruntime.FighterDefinitionDigest(fighter)},
			CardManifestDigests:    map[string]string{"strike": "sha256:strike", "scheme": "sha256:scheme"},
		},
		Fighters: map[model.DefinitionID]coreruntime.FighterDefinition{"mirror": fighter},
	}
}

func newHost(t *testing.T, ids coreruntime.IDProvider, rules contracts.RulesEngine) (*application.Host, *persistence.MemoryStore) {
	t.Helper()
	registry := coreruntime.NewMemoryDefinitionRegistry()
	if err := registry.Register("synthetic@v1", syntheticBundle()); err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	store := persistence.NewMemoryStore()
	return application.NewHost(registry, ids, store, rules), store
}

func createAndJoin(t *testing.T, host *application.Host) (application.CommandResult, application.CommandResult) {
	t.Helper()
	created, err := host.Execute(context.Background(), "principal-a", command(
		"create-1", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("create match: %v", err)
	}
	one := uint64(1)
	joined, err := host.Execute(context.Background(), "principal-b", command(
		"join-1", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("join match: %v", err)
	}
	return created, joined
}

func assertDistinctInstances(t *testing.T, state model.GameState) {
	t.Helper()
	if len(state.Players) != 2 || len(state.Fighters) != 2 || len(state.Cards) != 6 {
		t.Fatalf("unexpected instance counts: players=%d fighters=%d cards=%d", len(state.Players), len(state.Fighters), len(state.Cards))
	}
	seenCards := make(map[model.CardID]struct{})
	for _, player := range state.Players {
		if len(player.FighterInstanceIDs) != 1 {
			t.Fatalf("player does not own exactly one fighter: %+v", player)
		}
		for _, cardID := range player.PrivateZones["deck"] {
			if _, exists := seenCards[cardID]; exists {
				t.Fatalf("card runtime ID %q was reused", cardID)
			}
			seenCards[cardID] = struct{}{}
		}
	}
	strikeCount := 0
	for _, card := range state.Cards {
		if card.DefinitionID == "strike" {
			strikeCount++
		}
	}
	if strikeCount != 4 {
		t.Fatalf("expected four distinct strike copies across mirror fighters, got %d", strikeCount)
	}
}

func command(id string, kind string, matchID model.MatchID, actor model.PlayerID, revision *uint64, payload any) contracts.Command {
	encoded, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return contracts.Command{
		ID: model.CommandID(id), SchemaVersion: application.CommandSchemaV1,
		Type: kind, MatchID: matchID, ActorPlayerID: actor,
		ExpectedRevision: revision, Payload: encoded,
	}
}

type fakeRulesEngine struct {
	startProcedure  model.ProcedureRef
	resumeProcedure model.ProcedureRef
}

func (f *fakeRulesEngine) Resolve(_ model.GameState, input contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	if input.InteractionID == "" {
		f.startProcedure = input.Procedure
		var owner model.PlayerID
		if err := json.Unmarshal(input.Procedure.Bindings["actor_player_id"], &owner); err != nil {
			return contracts.ResolutionOutcome{}, err
		}
		return contracts.ResolutionOutcome{
			Status: contracts.ResolutionPending,
			PendingInteraction: &model.PendingInteraction{
				ID: "interaction-1", OwnerPlayerID: owner, Kind: "CHOOSE_TEST_OPTION",
				Visibility: "OWNER_ONLY", Prompt: json.RawMessage(`{"message":"choose"}`),
				LegalDomain:     json.RawMessage(`{"options":["left","right"]}`),
				ResumeProcedure: input.Procedure,
			},
		}, nil
	}
	f.resumeProcedure = input.Procedure
	return contracts.ResolutionOutcome{
		Status: contracts.ResolutionCompleted,
		Events: []contracts.DomainEvent{{
			Type: "rules.synthetic.choice_accepted", PublicPayload: json.RawMessage(`{"accepted":true}`),
		}},
	}, nil
}

func (f *fakeRulesEngine) LegalActions(state model.GameState, playerID model.PlayerID) ([]json.RawMessage, error) {
	if state.Resolver.PendingInteraction != nil {
		return nil, nil
	}
	return []json.RawMessage{json.RawMessage(`{"type":"StartAction","player_id":"` + string(playerID) + `"}`)}, nil
}

func (f *fakeRulesEngine) Project(_ model.GameState, playerID model.PlayerID) (json.RawMessage, error) {
	return json.Marshal(struct {
		PlayerID model.PlayerID `json:"player_id"`
	}{PlayerID: playerID})
}
