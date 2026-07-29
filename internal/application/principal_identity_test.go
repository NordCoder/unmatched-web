package application

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func TestPrincipalCanonicalizationUnifiesMemoryAuthorityOperations(t *testing.T) {
	ctx := context.Background()
	store := persistence.NewMemoryStore()
	host := newPrincipalTestHost(t, store, "principal-memory")

	ownerDecomposed := model.PrincipalID("principal-owner-e\u0301")
	ownerComposed := model.PrincipalID("principal-owner-é")
	create := principalCreateCommand(t, "principal-create")
	created, err := host.Execute(ctx, ownerDecomposed, create)
	if err != nil {
		t.Fatalf("create with decomposed principal: %v", err)
	}
	assertPrincipalProjection(t, ctx, host, created.MatchID, ownerDecomposed, created.PlayerID)
	assertPrincipalProjection(t, ctx, host, created.MatchID, ownerComposed, created.PlayerID)

	if playerID, ok := store.ResolveAuthority(created.MatchID, ownerComposed); !ok || playerID != created.PlayerID {
		t.Fatalf("canonical owner authority missing: player=%q ok=%v", playerID, ok)
	}
	if _, ok := store.ResolveAuthority(created.MatchID, ownerDecomposed); ok {
		t.Fatal("memory adapter persisted the original decomposed owner as an alias")
	}

	replayedCreate, err := host.Execute(ctx, ownerComposed, create)
	if err != nil || !reflect.DeepEqual(created, replayedCreate) {
		t.Fatalf("NFC-equivalent create did not replay the original result: result=%+v err=%v", replayedCreate, err)
	}
	if events := store.Events(created.MatchID); len(events) != 1 {
		t.Fatalf("equivalent create appended events: %d", len(events))
	}

	joinerDecomposed := model.PrincipalID("principal-joiner-e\u0301")
	joinerComposed := model.PrincipalID("principal-joiner-é")
	one := uint64(1)
	join := principalCommand(t, "principal-join", CommandJoinMatch, created.MatchID, "", &one,
		JoinMatchPayload{FighterDefinition: "principal-fighter"})
	joined, err := host.Execute(ctx, joinerDecomposed, join)
	if err != nil {
		t.Fatalf("join with decomposed principal: %v", err)
	}
	assertPrincipalProjection(t, ctx, host, created.MatchID, joinerDecomposed, joined.PlayerID)
	assertPrincipalProjection(t, ctx, host, created.MatchID, joinerComposed, joined.PlayerID)

	if playerID, ok := store.ResolveAuthority(created.MatchID, joinerComposed); !ok || playerID != joined.PlayerID {
		t.Fatalf("canonical joiner authority missing: player=%q ok=%v", playerID, ok)
	}
	if _, ok := store.ResolveAuthority(created.MatchID, joinerDecomposed); ok {
		t.Fatal("memory adapter persisted the original decomposed joiner as an alias")
	}

	replayedJoin, err := host.Execute(ctx, joinerComposed, join)
	if err != nil || !reflect.DeepEqual(joined, replayedJoin) {
		t.Fatalf("NFC-equivalent join did not replay the original result: result=%+v err=%v", replayedJoin, err)
	}
	if events := store.Events(created.MatchID); len(events) != 2 {
		t.Fatalf("equivalent join appended events: %d", len(events))
	}

	if _, err := host.ProjectContext(ctx, created.MatchID, "principal-genuinely-different"); CodeOf(err) != CodeUnauthorized {
		t.Fatalf("different normalized principal was authorized: %v", err)
	}
}

func TestInvalidPrincipalFailsBeforeAuthorityStorageAccess(t *testing.T) {
	tests := []struct {
		name      string
		principal model.PrincipalID
		code      string
	}{
		{name: "empty", principal: "", code: CodeUnauthorized},
		{name: "invalid UTF-8", principal: model.PrincipalID(string([]byte{0xff})), code: CodeInvalidCommand},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			memory := persistence.NewMemoryStore()
			recording := &principalIngressRecordingStore{EventStore: memory}
			host := newPrincipalTestHost(t, recording, "principal-invalid")

			if _, err := host.Execute(ctx, test.principal, principalCreateCommand(t, "principal-invalid-create")); CodeOf(err) != test.code {
				t.Fatalf("unexpected Execute error: %v", err)
			}
			if recording.acquireCalls != 0 {
				t.Fatalf("invalid principal reached command acquisition %d times", recording.acquireCalls)
			}
			if events := memory.Events("principal-invalid-match-1"); len(events) != 0 {
				t.Fatalf("invalid principal created state: %d events", len(events))
			}

			if _, err := host.ProjectContext(ctx, "principal-invalid-match-1", test.principal); CodeOf(err) != test.code {
				t.Fatalf("unexpected ProjectContext error: %v", err)
			}
			if recording.loadCalls != 0 || recording.resolveCalls != 0 {
				t.Fatalf("invalid principal reached projection storage: loads=%d resolves=%d", recording.loadCalls, recording.resolveCalls)
			}
		})
	}
}

func newPrincipalTestHost(t *testing.T, store persistence.EventStore, idPrefix string) *Host {
	t.Helper()
	registry := coreruntime.NewMemoryDefinitionRegistry()
	fighter := coreruntime.FighterDefinition{ID: "principal-fighter"}
	bundle := coreruntime.DefinitionBundle{
		Ref: model.DefinitionRef{
			RulesetVersion:     "principal-rules/v1",
			CapabilityRegistry: "principal-capabilities/v1",
			FighterManifestDigests: map[string]string{
				string(fighter.ID): coreruntime.FighterDefinitionDigest(fighter),
			},
		},
		Fighters: map[model.DefinitionID]coreruntime.FighterDefinition{fighter.ID: fighter},
	}
	if err := registry.Register("principal-test@v1", bundle); err != nil {
		t.Fatalf("register principal test definitions: %v", err)
	}
	return NewHost(registry, coreruntime.NewSequenceIDProvider(idPrefix), store, &principalTestRules{})
}

func principalCreateCommand(t *testing.T, commandID model.CommandID) contracts.Command {
	t.Helper()
	return principalCommand(t, commandID, CommandCreateMatch, "", "", nil, CreateMatchPayload{
		DefinitionKey: "principal-test@v1", FighterDefinition: "principal-fighter",
	})
}

func principalCommand(
	t *testing.T,
	commandID model.CommandID,
	commandType string,
	matchID model.MatchID,
	actorPlayerID model.PlayerID,
	expectedRevision *uint64,
	payload any,
) contracts.Command {
	t.Helper()
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal principal test payload: %v", err)
	}
	return contracts.Command{
		ID: commandID, SchemaVersion: CommandSchemaV1, Type: commandType,
		MatchID: matchID, ActorPlayerID: actorPlayerID,
		ExpectedRevision: expectedRevision, Payload: encoded,
	}
}

func assertPrincipalProjection(
	t *testing.T,
	ctx context.Context,
	host *Host,
	matchID model.MatchID,
	principal model.PrincipalID,
	wantPlayerID model.PlayerID,
) {
	t.Helper()
	projection, err := host.ProjectContext(ctx, matchID, principal)
	if err != nil {
		t.Fatalf("project principal %q: %v", principal, err)
	}
	if projection.PlayerID != wantPlayerID {
		t.Fatalf("projection resolved player %q, want %q", projection.PlayerID, wantPlayerID)
	}
}

type principalIngressRecordingStore struct {
	persistence.EventStore
	acquireCalls int
	loadCalls    int
	resolveCalls int
}

func (s *principalIngressRecordingStore) AcquireCommand(ctx context.Context, identity persistence.CommandIdentity) (persistence.CommandLease, persistence.CommandRecord, bool, error) {
	s.acquireCalls++
	return s.EventStore.AcquireCommand(ctx, identity)
}

func (s *principalIngressRecordingStore) LoadEvents(ctx context.Context, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	s.loadCalls++
	return s.EventStore.LoadEvents(ctx, matchID)
}

func (s *principalIngressRecordingStore) ResolveAuthorityContext(ctx context.Context, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	s.resolveCalls++
	return s.EventStore.ResolveAuthorityContext(ctx, matchID, principalID)
}

type principalTestRules struct{}

func (*principalTestRules) Resolve(model.GameState, contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted}, nil
}

func (*principalTestRules) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}

func (*principalTestRules) Project(model.GameState, model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}

var _ persistence.EventStore = (*principalIngressRecordingStore)(nil)
var _ contracts.RulesEngine = (*principalTestRules)(nil)
