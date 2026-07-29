package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

func assertPrincipalIdentitySurvivesReconstruction(t *testing.T, ctx context.Context, pool *pgxpool.Pool, registry coreruntime.DefinitionRegistry) {
	t.Helper()
	store := newPostgresStore(t, pool)
	host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("principal-nfc"), store, &postgresCompletedRules{})

	ownerDecomposed := model.PrincipalID("principal-owner-e\u0301")
	ownerComposed := model.PrincipalID("principal-owner-é")
	create := command(
		"principal-nfc-create", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
	)
	created, err := host.Execute(ctx, ownerDecomposed, create)
	if err != nil {
		t.Fatalf("create with decomposed principal: %v", err)
	}
	assertPostgresPrincipalProjection(t, ctx, host, created.MatchID, ownerDecomposed, created.PlayerID)
	assertPostgresPrincipalProjection(t, ctx, host, created.MatchID, ownerComposed, created.PlayerID)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_authority_bindings WHERE match_id = $1 AND principal_id = $2", 1,
		created.MatchID, ownerComposed)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_authority_bindings WHERE match_id = $1 AND principal_id = $2", 0,
		created.MatchID, ownerDecomposed)

	one := uint64(1)
	joinerDecomposed := model.PrincipalID("principal-joiner-e\u0301")
	joinerComposed := model.PrincipalID("principal-joiner-é")
	join := command(
		"principal-nfc-join", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"},
	)
	joined, err := host.Execute(ctx, joinerDecomposed, join)
	if err != nil {
		t.Fatalf("join with decomposed principal: %v", err)
	}
	assertPostgresPrincipalProjection(t, ctx, host, created.MatchID, joinerDecomposed, joined.PlayerID)
	assertPostgresPrincipalProjection(t, ctx, host, created.MatchID, joinerComposed, joined.PlayerID)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_authority_bindings WHERE match_id = $1 AND principal_id = $2", 1,
		created.MatchID, joinerComposed)
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_authority_bindings WHERE match_id = $1 AND principal_id = $2", 0,
		created.MatchID, joinerDecomposed)

	reconstructed := application.NewHost(
		registry, coreruntime.NewSequenceIDProvider("principal-nfc-unused"), newPostgresStore(t, pool), &postgresCompletedRules{},
	)
	assertPostgresPrincipalProjection(t, ctx, reconstructed, created.MatchID, ownerDecomposed, created.PlayerID)
	assertPostgresPrincipalProjection(t, ctx, reconstructed, created.MatchID, ownerComposed, created.PlayerID)
	assertPostgresPrincipalProjection(t, ctx, reconstructed, created.MatchID, joinerDecomposed, joined.PlayerID)
	assertPostgresPrincipalProjection(t, ctx, reconstructed, created.MatchID, joinerComposed, joined.PlayerID)

	replayedCreate, err := reconstructed.Execute(ctx, ownerComposed, create)
	if err != nil || !reflect.DeepEqual(created, replayedCreate) {
		t.Fatalf("equivalent owner principal did not replay create: result=%+v err=%v", replayedCreate, err)
	}
	replayedJoin, err := reconstructed.Execute(ctx, joinerComposed, join)
	if err != nil || !reflect.DeepEqual(joined, replayedJoin) {
		t.Fatalf("equivalent joiner principal did not replay join: result=%+v err=%v", replayedJoin, err)
	}
	if events := countMatchEvents(t, ctx, pool, created.MatchID); events != 2 {
		t.Fatalf("equivalent principal duplicate appended events: %d", events)
	}

	if _, err := reconstructed.ProjectContext(ctx, created.MatchID, "principal-genuinely-different"); application.CodeOf(err) != application.CodeUnauthorized {
		t.Fatalf("different normalized principal was authorized: %v", err)
	}
}

func assertPostgresPrincipalProjection(
	t *testing.T,
	ctx context.Context,
	host *application.Host,
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
