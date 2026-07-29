package core_test

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

func assertSaturatedPoolProgress(t *testing.T, ctx context.Context, pool *pgxpool.Pool, registry coreruntime.DefinitionRegistry) {
	t.Helper()
	config := pool.Config().Copy()
	config.MaxConns = 2
	config.MinConns = 0
	smallPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("open saturated test pool: %v", err)
	}
	t.Cleanup(smallPool.Close)

	setupStore := newPostgresStore(t, smallPool)
	setupHost := application.NewHost(registry, coreruntime.NewSequenceIDProvider("saturated-setup"), setupStore, &postgresCompletedRules{})
	created, err := setupHost.Execute(ctx, "principal-saturated-a", command(
		"saturated-create", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("create saturated match: %v", err)
	}
	one := uint64(1)
	if _, err := setupHost.Execute(ctx, "principal-saturated-b", command(
		"saturated-join", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"},
	)); err != nil {
		t.Fatalf("join saturated match: %v", err)
	}

	two := uint64(2)
	action := command(
		"saturated-action", application.CommandStartAction, created.MatchID, created.PlayerID, &two,
		application.StartActionPayload{Kind: application.ActionManeuver},
	)
	rules := &countingCompletedRules{}
	ownerStore := newAcquireBarrierStore(newPostgresStore(t, smallPool), action.ID)
	ownerHost := application.NewHost(registry, coreruntime.NewSequenceIDProvider("saturated-owner"), ownerStore, rules)
	waiterHost := application.NewHost(registry, coreruntime.NewSequenceIDProvider("saturated-waiter"), newPostgresStore(t, smallPool), rules)

	ownerResult := make(chan commandExecution, 1)
	waiterResult := make(chan commandExecution, 1)
	go func() {
		result, err := ownerHost.Execute(ctx, "principal-saturated-a", action)
		ownerResult <- commandExecution{result: result, err: err}
	}()
	select {
	case <-ownerStore.acquired:
	case <-ctx.Done():
		t.Fatal("owner did not acquire command lease")
	}
	go func() {
		result, err := waiterHost.Execute(ctx, "principal-saturated-a", action)
		waiterResult <- commandExecution{result: result, err: err}
	}()
	deadline := time.Now().Add(5 * time.Second)
	for smallPool.Stat().AcquiredConns() < 2 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if acquired := smallPool.Stat().AcquiredConns(); acquired != 2 {
		t.Fatalf("duplicate did not saturate the two-connection pool: acquired=%d", acquired)
	}
	ownerStore.unblockLease()

	first := receiveExecution(t, ownerResult, "owner")
	second := receiveExecution(t, waiterResult, "waiter")
	if first.err != nil || second.err != nil {
		t.Fatalf("saturated duplicate failed: owner=%v waiter=%v", first.err, second.err)
	}
	if !reflect.DeepEqual(first.result, second.result) {
		t.Fatalf("saturated duplicate changed result\nowner: %#v\nwaiter: %#v", first.result, second.result)
	}
	if calls := rules.calls.Load(); calls != 1 {
		t.Fatalf("Rules executed %d times under saturation", calls)
	}
}

type countingCompletedRules struct {
	calls atomic.Int32
}

func (r *countingCompletedRules) Resolve(model.GameState, contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	r.calls.Add(1)
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted}, nil
}

func (*countingCompletedRules) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}

func (*countingCompletedRules) Project(model.GameState, model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}

type acquireBarrierStore struct {
	delegate persistence.EventStore
	target   model.CommandID
	acquired chan struct{}
	release  chan struct{}
	once     sync.Once
	unblock  sync.Once
}

func newAcquireBarrierStore(delegate persistence.EventStore, target model.CommandID) *acquireBarrierStore {
	return &acquireBarrierStore{
		delegate: delegate,
		target:   target,
		acquired: make(chan struct{}),
		release:  make(chan struct{}),
	}
}

func (s *acquireBarrierStore) AcquireCommand(ctx context.Context, identity persistence.CommandIdentity) (persistence.CommandLease, persistence.CommandRecord, bool, error) {
	lease, record, duplicate, err := s.delegate.AcquireCommand(ctx, identity)
	if err != nil || duplicate || identity.Key.CommandID != s.target {
		return lease, record, duplicate, err
	}
	s.once.Do(func() { close(s.acquired) })
	select {
	case <-s.release:
		return lease, record, duplicate, nil
	case <-ctx.Done():
		_ = s.delegate.AbortCommand(context.WithoutCancel(ctx), lease)
		return persistence.CommandLease{}, persistence.CommandRecord{}, false, ctx.Err()
	}
}

func (s *acquireBarrierStore) Commit(ctx context.Context, lease persistence.CommandLease, request persistence.CommitRequest) (persistence.CommandRecord, error) {
	return s.delegate.Commit(ctx, lease, request)
}

func (s *acquireBarrierStore) RejectCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID, code, message string) (persistence.CommandRecord, error) {
	return s.delegate.RejectCommand(ctx, lease, matchID, code, message)
}

func (s *acquireBarrierStore) AbortCommand(ctx context.Context, lease persistence.CommandLease) error {
	return s.delegate.AbortCommand(ctx, lease)
}

func (s *acquireBarrierStore) LoadEvents(ctx context.Context, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	return s.delegate.LoadEvents(ctx, matchID)
}

func (s *acquireBarrierStore) ResolveAuthorityContext(ctx context.Context, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	return s.delegate.ResolveAuthorityContext(ctx, matchID, principalID)
}

func (s *acquireBarrierStore) LoadEventsForCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	return s.delegate.LoadEventsForCommand(ctx, lease, matchID)
}

func (s *acquireBarrierStore) ResolveAuthorityForCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	return s.delegate.ResolveAuthorityForCommand(ctx, lease, matchID, principalID)
}

func (s *acquireBarrierStore) unblockLease() {
	s.unblock.Do(func() { close(s.release) })
}

type commandExecution struct {
	result application.CommandResult
	err    error
}

func receiveExecution(t *testing.T, channel <-chan commandExecution, name string) commandExecution {
	t.Helper()
	select {
	case result := <-channel:
		return result
	case <-time.After(10 * time.Second):
		t.Fatalf("%s command did not complete", name)
		return commandExecution{}
	}
}

var _ persistence.EventStore = (*acquireBarrierStore)(nil)
var _ contracts.RulesEngine = (*countingCompletedRules)(nil)
