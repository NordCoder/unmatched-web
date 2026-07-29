package core_test

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
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

func TestPostgresPersistenceCandidate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool := openIsolatedPostgres(t, ctx)
	if err := persistence.ApplyMigrations(ctx, pool); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if err := persistence.ApplyMigrations(ctx, pool); err != nil {
		t.Fatalf("reapply migrations: %v", err)
	}

	registry := coreruntime.NewMemoryDefinitionRegistry()
	bundle := syntheticBundle()
	if err := registry.Register("synthetic@v1", bundle); err != nil {
		t.Fatalf("register definitions: %v", err)
	}

	store1 := newPostgresStore(t, pool)
	host1 := application.NewHost(registry, coreruntime.NewSequenceIDProvider("pg-a"), store1, &fakeRulesEngine{})
	create := command("pg-create-1", application.CommandCreateMatch, "", "", nil, application.CreateMatchPayload{
		DefinitionKey: "synthetic@v1", FighterDefinition: "mirror",
	})
	created, err := host1.Execute(ctx, "principal-pg-a", create)
	if err != nil {
		t.Fatalf("create durable match: %v", err)
	}

	store2 := newPostgresStore(t, pool)
	host2 := application.NewHost(registry, coreruntime.NewSequenceIDProvider("unused-duplicate"), store2, &fakeRulesEngine{})
	duplicateCreate, err := host2.Execute(ctx, "principal-pg-a", create)
	if err != nil {
		t.Fatalf("repeat create after store reconstruction: %v", err)
	}
	if !reflect.DeepEqual(created, duplicateCreate) {
		t.Fatalf("reconstructed duplicate changed result\nfirst: %#v\nduplicate: %#v", created, duplicateCreate)
	}
	events, err := store2.LoadEvents(ctx, created.MatchID)
	if err != nil {
		t.Fatalf("load events after create: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("duplicate create appended events: %d", len(events))
	}

	one := uint64(1)
	join := command("pg-join-1", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"})
	joined, err := host2.Execute(ctx, "principal-pg-b", join)
	if err != nil {
		t.Fatalf("join durable match: %v", err)
	}

	store3 := newPostgresStore(t, pool)
	host3 := application.NewHost(registry, coreruntime.NewSequenceIDProvider("unused-restart"), store3, &fakeRulesEngine{})
	state, err := host3.StateContext(ctx, created.MatchID)
	if err != nil {
		t.Fatalf("reopen durable match: %v", err)
	}
	if state.Game.Revision != 2 || state.Game.EventSequence != 2 || len(state.Game.Players) != 2 {
		t.Fatalf("unexpected reconstructed state: %+v", state.Game)
	}
	playerID, ok, err := store3.ResolveAuthorityContext(ctx, created.MatchID, "principal-pg-b")
	if err != nil {
		t.Fatalf("resolve durable authority: %v", err)
	}
	if !ok || playerID != joined.PlayerID {
		t.Fatalf("durable authority mismatch: player=%q ok=%v", playerID, ok)
	}
	duplicateJoin, err := host3.Execute(ctx, "principal-pg-b", join)
	if err != nil || !reflect.DeepEqual(joined, duplicateJoin) {
		t.Fatalf("duplicate join after reconstruction changed result: result=%+v err=%v", duplicateJoin, err)
	}

	conflictingCreate := create
	conflictingCreate.Payload = json.RawMessage(`{"definition_key":"synthetic@v1","fighter_definition_id":"other"}`)
	_, err = host3.Execute(ctx, "principal-pg-a", conflictingCreate)
	if application.CodeOf(err) != application.CodeCommandConflict {
		t.Fatalf("expected derived command conflict, got %v", err)
	}
	assertDatabaseCount(t, ctx, pool, "core_command_results", 2)
	assertDatabaseCount(t, ctx, pool, "core_events", 2)

	stale := uint64(1)
	staleCommand := command("pg-stale-1", application.CommandStartAction, created.MatchID, created.PlayerID, &stale,
		application.StartActionPayload{Kind: application.ActionScheme})
	_, err = host3.Execute(ctx, "principal-pg-a", staleCommand)
	if application.CodeOf(err) != application.CodeRevisionConflict {
		t.Fatalf("expected durable revision rejection, got %v", err)
	}
	host4 := application.NewHost(registry, coreruntime.NewSequenceIDProvider("unused-rejection"), newPostgresStore(t, pool), &fakeRulesEngine{})
	_, err = host4.Execute(ctx, "principal-pg-a", staleCommand)
	if application.CodeOf(err) != application.CodeRevisionConflict {
		t.Fatalf("expected reconstructed rejection, got %v", err)
	}
	assertDatabaseCount(t, ctx, pool, "core_command_results", 3)
	assertDatabaseCount(t, ctx, pool, "core_events", 2)

	t.Run("failed transaction exposes no partial acceptance", func(t *testing.T) {
		assertRollbackIsAtomic(t, ctx, pool, store3, created.MatchID, bundle.Ref)
	})

	t.Run("same request executes once across stores", func(t *testing.T) {
		rules := newBlockingRulesEngine()
		storeA := newPostgresStore(t, pool)
		storeB := newPostgresStore(t, pool)
		hostA := application.NewHost(registry, coreruntime.NewSequenceIDProvider("pg-dupe"), storeA, rules)
		hostB := application.NewHost(registry, coreruntime.NewSequenceIDProvider("pg-dupe"), storeB, rules)
		two := uint64(2)
		action := command("pg-action-1", application.CommandStartAction, created.MatchID, created.PlayerID, &two,
			application.StartActionPayload{Kind: application.ActionManeuver})

		type execution struct {
			result application.CommandResult
			err    error
		}
		first := make(chan execution, 1)
		second := make(chan execution, 1)
		go func() {
			result, err := hostA.Execute(ctx, "principal-pg-a", action)
			first <- execution{result: result, err: err}
		}()
		select {
		case <-rules.started:
		case <-ctx.Done():
			t.Fatal("first command did not reach Rules")
		}
		go func() {
			result, err := hostB.Execute(ctx, "principal-pg-a", action)
			second <- execution{result: result, err: err}
		}()
		select {
		case result := <-second:
			t.Fatalf("duplicate returned before the original committed: %+v", result)
		case <-time.After(150 * time.Millisecond):
		}
		defer rules.unblock()
		rules.unblock()
		firstResult := <-first
		secondResult := <-second
		if firstResult.err != nil || secondResult.err != nil {
			t.Fatalf("concurrent duplicate failed: first=%v second=%v", firstResult.err, secondResult.err)
		}
		if !reflect.DeepEqual(firstResult.result, secondResult.result) {
			t.Fatalf("concurrent duplicate changed result\nfirst: %#v\nsecond: %#v", firstResult.result, secondResult.result)
		}
		if calls := rules.calls.Load(); calls != 1 {
			t.Fatalf("Rules executed %d times", calls)
		}
	})

	t.Run("command IDs remain principal-scoped inside one match", func(t *testing.T) {
		store := newPostgresStore(t, pool)
		host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("pg-scope"), store, &postgresCompletedRules{})
		scopeCreated, err := host.Execute(ctx, "principal-scope-a", command(
			"pg-scope-create", application.CommandCreateMatch, "", "", nil,
			application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
		))
		if err != nil {
			t.Fatalf("create principal-scope match: %v", err)
		}
		one := uint64(1)
		scopeJoined, err := host.Execute(ctx, "principal-scope-b", command(
			"pg-scope-join", application.CommandJoinMatch, scopeCreated.MatchID, "", &one,
			application.JoinMatchPayload{FighterDefinition: "mirror"},
		))
		if err != nil {
			t.Fatalf("join principal-scope match: %v", err)
		}

		two := uint64(2)
		sharedID := "pg-shared-action-id"
		if _, err := host.Execute(ctx, "principal-scope-a", command(
			sharedID, application.CommandStartAction, scopeCreated.MatchID, scopeCreated.PlayerID, &two,
			application.StartActionPayload{Kind: application.ActionManeuver},
		)); err != nil {
			t.Fatalf("first principal command with shared ID: %v", err)
		}
		three := uint64(3)
		if _, err := host.Execute(ctx, "principal-scope-b", command(
			sharedID, application.CommandStartAction, scopeCreated.MatchID, scopeJoined.PlayerID, &three,
			application.StartActionPayload{Kind: application.ActionManeuver},
		)); err != nil {
			t.Fatalf("second principal command with shared ID: %v", err)
		}
		assertQueryCount(t, ctx, pool,
			"SELECT count(*) FROM core_event_batches WHERE match_id = $1 AND command_id = $2", 2,
			scopeCreated.MatchID, sharedID)
	})

	t.Run("host waits for durable commit", func(t *testing.T) {
		store := newPostgresStore(t, pool)
		barrier := newCommitBarrierStore(store)
		host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("pg-barrier"), barrier, &fakeRulesEngine{})
		createBarrier := command("pg-barrier-create", application.CommandCreateMatch, "", "", nil,
			application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"})
		done := make(chan error, 1)
		go func() {
			_, err := host.Execute(ctx, "principal-barrier", createBarrier)
			done <- err
		}()
		select {
		case <-barrier.started:
		case <-ctx.Done():
			t.Fatal("command did not reach durable commit")
		}
		select {
		case err := <-done:
			t.Fatalf("host acknowledged before commit completed: %v", err)
		case <-time.After(150 * time.Millisecond):
		}
		defer barrier.unblock()
		barrier.unblock()
		if err := <-done; err != nil {
			t.Fatalf("command failed after commit release: %v", err)
		}
	})
}

func assertRollbackIsAtomic(t *testing.T, ctx context.Context, pool *pgxpool.Pool, store *persistence.PostgresStore, matchID model.MatchID, definitionRef model.DefinitionRef) {
	t.Helper()
	var revision, sequence int64
	var duplicateEventID model.EventID
	if err := pool.QueryRow(ctx,
		"SELECT revision, event_sequence FROM core_match_heads WHERE match_id = $1", matchID,
	).Scan(&revision, &sequence); err != nil {
		t.Fatalf("read head before rollback test: %v", err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT event_id FROM core_events WHERE match_id = $1 ORDER BY sequence LIMIT 1", matchID,
	).Scan(&duplicateEventID); err != nil {
		t.Fatalf("read duplicate event ID: %v", err)
	}

	key := persistence.CommandKey{PrincipalID: "principal-rollback", CommandID: "pg-rollback-1"}
	lease, _, duplicate, err := store.AcquireCommand(ctx, persistence.CommandIdentity{
		Key: key, Fingerprint: []byte("rollback-fingerprint"), MatchID: matchID,
		Scope: persistence.CommandScopeExistingSeat,
	})
	if err != nil || duplicate {
		t.Fatalf("reserve rollback command: duplicate=%v err=%v", duplicate, err)
	}
	nextRevision := uint64(revision + 1)
	nextSequence := uint64(sequence + 1)
	event := contracts.DomainEvent{
		SchemaVersion:   "core-event/v1",
		ID:              duplicateEventID,
		MatchID:         matchID,
		Sequence:        nextSequence,
		Revision:        nextRevision,
		Type:            "core.test.rollback",
		CausedByCommand: key.CommandID,
		RulesetVersion:  definitionRef.RulesetVersion,
		PublicPayload:   json.RawMessage(`{}`),
		PrivatePayloads: map[model.PlayerID]json.RawMessage{},
	}
	_, err = store.Commit(ctx, lease, persistence.CommitRequest{
		MatchID:       matchID,
		DefinitionRef: definitionRef,
		Batch: contracts.EventBatch{
			CommandID: key.CommandID, PreviousRevision: uint64(revision),
			NextRevision: nextRevision, Events: []contracts.DomainEvent{event},
		},
		Result: []byte(`{"accepted":true}`),
	})
	if err == nil {
		t.Fatal("duplicate event ID unexpectedly committed")
	}
	if err := store.AbortCommand(ctx, lease); err != nil {
		t.Fatalf("abort failed transaction: %v", err)
	}

	var currentRevision, currentSequence int64
	if err := pool.QueryRow(ctx,
		"SELECT revision, event_sequence FROM core_match_heads WHERE match_id = $1", matchID,
	).Scan(&currentRevision, &currentSequence); err != nil {
		t.Fatalf("read head after rollback test: %v", err)
	}
	if currentRevision != revision || currentSequence != sequence {
		t.Fatalf("failed transaction advanced head: before=%d/%d after=%d/%d", revision, sequence, currentRevision, currentSequence)
	}
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_command_results WHERE principal_id = $1 AND command_id = $2", 0,
		key.PrincipalID, key.CommandID)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_event_batches WHERE match_id = $1 AND revision = $2", 0,
		matchID, int64(nextRevision))
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_command_reservations WHERE principal_id = $1 AND command_id = $2", 0,
		key.PrincipalID, key.CommandID)
}

type postgresCompletedRules struct{}

func (*postgresCompletedRules) Resolve(model.GameState, contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted}, nil
}

func (*postgresCompletedRules) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}

func (*postgresCompletedRules) Project(model.GameState, model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}

type blockingRulesEngine struct {
	startedOnce sync.Once
	releaseOnce sync.Once
	started     chan struct{}
	release     chan struct{}
	calls       atomic.Int32
}

func newBlockingRulesEngine() *blockingRulesEngine {
	return &blockingRulesEngine{started: make(chan struct{}), release: make(chan struct{})}
}

func (r *blockingRulesEngine) Resolve(_ model.GameState, _ contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	r.calls.Add(1)
	r.startedOnce.Do(func() { close(r.started) })
	<-r.release
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted}, nil
}

func (r *blockingRulesEngine) unblock() {
	r.releaseOnce.Do(func() { close(r.release) })
}

func (r *blockingRulesEngine) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}

func (r *blockingRulesEngine) Project(model.GameState, model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}

type commitBarrierStore struct {
	delegate    persistence.EventStore
	started     chan struct{}
	release     chan struct{}
	once        sync.Once
	releaseOnce sync.Once
}

func newCommitBarrierStore(delegate persistence.EventStore) *commitBarrierStore {
	return &commitBarrierStore{delegate: delegate, started: make(chan struct{}), release: make(chan struct{})}
}

func (s *commitBarrierStore) AcquireCommand(ctx context.Context, identity persistence.CommandIdentity) (persistence.CommandLease, persistence.CommandRecord, bool, error) {
	return s.delegate.AcquireCommand(ctx, identity)
}

func (s *commitBarrierStore) Commit(ctx context.Context, lease persistence.CommandLease, request persistence.CommitRequest) (persistence.CommandRecord, error) {
	s.once.Do(func() { close(s.started) })
	select {
	case <-s.release:
	case <-ctx.Done():
		return persistence.CommandRecord{}, ctx.Err()
	}
	return s.delegate.Commit(ctx, lease, request)
}

func (s *commitBarrierStore) unblock() {
	s.releaseOnce.Do(func() { close(s.release) })
}

func (s *commitBarrierStore) RejectCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID, code, message string) (persistence.CommandRecord, error) {
	return s.delegate.RejectCommand(ctx, lease, matchID, code, message)
}

func (s *commitBarrierStore) AbortCommand(ctx context.Context, lease persistence.CommandLease) error {
	return s.delegate.AbortCommand(ctx, lease)
}

func (s *commitBarrierStore) LoadEvents(ctx context.Context, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	return s.delegate.LoadEvents(ctx, matchID)
}

func (s *commitBarrierStore) ResolveAuthorityContext(ctx context.Context, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	return s.delegate.ResolveAuthorityContext(ctx, matchID, principalID)
}

func newPostgresStore(t *testing.T, pool *pgxpool.Pool) *persistence.PostgresStore {
	t.Helper()
	store, err := persistence.NewPostgresStore(pool)
	if err != nil {
		t.Fatalf("create PostgreSQL store: %v", err)
	}
	return store
}

func openIsolatedPostgres(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()
	baseURL := os.Getenv("CORE_TEST_DATABASE_URL")
	var stopContainer func()
	if baseURL == "" {
		baseURL, stopContainer = startPostgresContainer(t, ctx)
	}
	if stopContainer != nil {
		t.Cleanup(stopContainer)
	}

	adminPool, err := persistence.OpenPostgres(ctx, persistence.PostgresConfig{
		DatabaseURL: baseURL, MaxConns: 4, MinConns: 0, ConnectTimeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("open PostgreSQL admin pool: %v", err)
	}
	t.Cleanup(adminPool.Close)
	schema := "core_test_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	if _, err := adminPool.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatalf("create isolated schema: %v", err)
	}
	t.Cleanup(func() {
		cleanupContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = adminPool.Exec(cleanupContext, "DROP SCHEMA IF EXISTS "+schema+" CASCADE")
	})

	pool, err := persistence.OpenPostgres(ctx, persistence.PostgresConfig{
		DatabaseURL: baseURL, MaxConns: 12, MinConns: 0,
		ConnectTimeout: 10 * time.Second, SearchPath: schema,
	})
	if err != nil {
		t.Fatalf("open isolated PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func startPostgresContainer(t *testing.T, ctx context.Context) (string, func()) {
	t.Helper()
	docker, err := exec.LookPath("docker")
	if err != nil {
		if os.Getenv("CI") != "" {
			t.Fatal("Docker is required for PostgreSQL integration tests in CI")
		}
		t.Skip("Docker is unavailable and CORE_TEST_DATABASE_URL is not set")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate PostgreSQL test port: %v", err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("release PostgreSQL test port: %v", err)
	}
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatalf("parse PostgreSQL test address: %v", err)
	}

	name := "unmatched-core-test-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	command := exec.CommandContext(ctx, docker, "run", "--rm", "-d", "--name", name,
		"-e", "POSTGRES_DB=unmatched", "-e", "POSTGRES_USER=unmatched",
		"-e", "POSTGRES_PASSWORD=unmatched-test-only",
		"-p", "127.0.0.1:"+port+":5432", "postgres:17-alpine")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v\n%s", err, output)
	}
	containerID := strings.TrimSpace(string(output))
	stop := func() {
		cleanupContext, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_ = exec.CommandContext(cleanupContext, docker, "rm", "-f", containerID).Run()
	}

	baseURL := "postgres://unmatched:unmatched-test-only@" + address + "/unmatched?sslmode=disable"
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		pool, openErr := persistence.OpenPostgres(ctx, persistence.PostgresConfig{
			DatabaseURL: baseURL, MaxConns: 2, MinConns: 0, ConnectTimeout: time.Second,
		})
		if openErr == nil {
			pool.Close()
			return baseURL, stop
		}
		time.Sleep(250 * time.Millisecond)
	}
	logs, _ := exec.CommandContext(ctx, docker, "logs", containerID).CombinedOutput()
	stop()
	t.Fatalf("PostgreSQL container did not become ready:\n%s", logs)
	return "", nil
}

func assertDatabaseCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, table string, expected int64) {
	t.Helper()
	allowed := map[string]bool{
		"core_command_results": true,
		"core_events":          true,
	}
	if !allowed[table] {
		t.Fatalf("unsupported count table %q", table)
	}
	assertQueryCount(t, ctx, pool, "SELECT count(*) FROM "+table, expected)
}

func assertQueryCount(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string, expected int64, arguments ...any) {
	t.Helper()
	var actual int64
	if err := pool.QueryRow(ctx, query, arguments...).Scan(&actual); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if actual != expected {
		t.Fatalf("unexpected count for %q: got %d want %d", query, actual, expected)
	}
}

var _ persistence.EventStore = (*commitBarrierStore)(nil)
var _ contracts.RulesEngine = (*blockingRulesEngine)(nil)
