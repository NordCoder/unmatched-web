package core_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

func assertAuthorityRollbackIsAtomic(t *testing.T, ctx context.Context, pool *pgxpool.Pool, registry coreruntime.DefinitionRegistry, definitionRef model.DefinitionRef) {
	t.Helper()
	store := newPostgresStore(t, pool)
	host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("authority-rollback"), store, &postgresCompletedRules{})
	created, err := host.Execute(ctx, "principal-authority-owner", command(
		"authority-rollback-create", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("create authority rollback match: %v", err)
	}

	var revision, sequence int64
	var duplicateEventID model.EventID
	if err := pool.QueryRow(ctx,
		"SELECT revision, event_sequence FROM core_match_heads WHERE match_id = $1", created.MatchID,
	).Scan(&revision, &sequence); err != nil {
		t.Fatalf("read rollback head: %v", err)
	}
	if err := pool.QueryRow(ctx,
		"SELECT event_id FROM core_events WHERE match_id = $1 ORDER BY sequence LIMIT 1", created.MatchID,
	).Scan(&duplicateEventID); err != nil {
		t.Fatalf("read rollback event ID: %v", err)
	}

	key := persistence.CommandKey{PrincipalID: "principal-authority-rollback", CommandID: "authority-rollback-late-failure"}
	lease, _, duplicate, err := store.AcquireCommand(ctx, persistence.CommandIdentity{
		Key: key, Fingerprint: []byte("authority-rollback-fingerprint"), MatchID: created.MatchID,
		Scope: persistence.CommandScopeExistingSeat,
	})
	if err != nil || duplicate {
		t.Fatalf("reserve authority rollback command: duplicate=%v err=%v", duplicate, err)
	}
	nextRevision := uint64(revision + 1)
	nextSequence := uint64(sequence + 1)
	_, err = store.Commit(ctx, lease, persistence.CommitRequest{
		MatchID:       created.MatchID,
		DefinitionRef: definitionRef,
		Batch: contracts.EventBatch{
			CommandID:        key.CommandID,
			PreviousRevision: uint64(revision),
			NextRevision:     nextRevision,
			Events: []contracts.DomainEvent{{
				SchemaVersion: "core-event/v1", ID: duplicateEventID, MatchID: created.MatchID,
				Sequence: nextSequence, Revision: nextRevision, Type: "core.test.authority_rollback",
				CausedByCommand: key.CommandID, RulesetVersion: definitionRef.RulesetVersion,
				PublicPayload: json.RawMessage(`{}`), PrivatePayloads: map[model.PlayerID]json.RawMessage{},
			}},
		},
		Result: []byte(`{"accepted":true}`),
		Authority: &persistence.AuthorityRecord{
			MatchID: created.MatchID, PlayerID: "player-authority-rollback", PrincipalID: key.PrincipalID,
			Seat: 2, BindingVersion: 1, Status: persistence.AuthorityActive, EstablishedByCommandID: key.CommandID,
		},
	})
	if err == nil {
		t.Fatal("late event failure unexpectedly committed authority")
	}
	if err := store.AbortCommand(ctx, lease); err != nil {
		t.Fatalf("abort authority rollback command: %v", err)
	}

	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_authority_bindings WHERE match_id = $1 AND principal_id = $2", 0,
		created.MatchID, key.PrincipalID)
	assertNoDurableCommand(t, ctx, pool, key.PrincipalID, key.CommandID)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_event_batches WHERE match_id = $1 AND revision = $2", 0,
		created.MatchID, int64(nextRevision))
	var currentRevision, currentSequence int64
	if err := pool.QueryRow(ctx,
		"SELECT revision, event_sequence FROM core_match_heads WHERE match_id = $1", created.MatchID,
	).Scan(&currentRevision, &currentSequence); err != nil {
		t.Fatalf("read rollback head after failure: %v", err)
	}
	if currentRevision != revision || currentSequence != sequence {
		t.Fatalf("authority-bearing failure advanced head: before=%d/%d after=%d/%d", revision, sequence, currentRevision, currentSequence)
	}
}

func assertNoDurableCommand(t *testing.T, ctx context.Context, pool *pgxpool.Pool, principal model.PrincipalID, commandID model.CommandID) {
	t.Helper()
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_command_results WHERE principal_id = $1 AND command_id = $2", 0,
		principal, commandID)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_command_reservations WHERE principal_id = $1 AND command_id = $2", 0,
		principal, commandID)
}

func countMatchEvents(t *testing.T, ctx context.Context, pool *pgxpool.Pool, matchID model.MatchID) int64 {
	t.Helper()
	var count int64
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM core_events WHERE match_id = $1", matchID).Scan(&count); err != nil {
		t.Fatalf("count match events: %v", err)
	}
	return count
}

func (s *commitBarrierStore) LoadEventsForCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	return s.delegate.LoadEventsForCommand(ctx, lease, matchID)
}

func (s *commitBarrierStore) ResolveAuthorityForCommand(ctx context.Context, lease persistence.CommandLease, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	return s.delegate.ResolveAuthorityForCommand(ctx, lease, matchID, principalID)
}
