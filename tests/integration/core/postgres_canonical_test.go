package core_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

func assertCanonicalIdentitySurvivesReconstruction(t *testing.T, ctx context.Context, pool *pgxpool.Pool, registry coreruntime.DefinitionRegistry) {
	t.Helper()
	store := newPostgresStore(t, pool)
	host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("canonical"), store, &postgresCompletedRules{})
	create := command(
		"canonical-create", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
	)
	create.Payload = json.RawMessage(`{"fighter_definition_id":"mirror","definition_key":"synthetic@v1"}`)
	created, err := host.Execute(ctx, "principal-canonical-a", create)
	if err != nil {
		t.Fatalf("create canonical match: %v", err)
	}

	reconstructed := application.NewHost(registry, coreruntime.NewSequenceIDProvider("canonical-unused"), newPostgresStore(t, pool), &postgresCompletedRules{})
	reordered := create
	reordered.Payload = json.RawMessage(`{"definition_key":"synthetic@v1","fighter_definition_id":"mirror"}`)
	replayed, err := reconstructed.Execute(ctx, "principal-canonical-a", reordered)
	if err != nil || !reflect.DeepEqual(created, replayed) {
		t.Fatalf("reordered payload did not replay original result: result=%+v err=%v", replayed, err)
	}

	zero := uint64(0)
	explicitRevision := reordered
	explicitRevision.ExpectedRevision = &zero
	beforeCreateEvents := countMatchEvents(t, ctx, pool, created.MatchID)
	if _, err := reconstructed.Execute(ctx, "principal-canonical-a", explicitRevision); application.CodeOf(err) != application.CodeCommandConflict {
		t.Fatalf("explicit revision did not differ from absent revision policy: %v", err)
	}
	if afterCreateEvents := countMatchEvents(t, ctx, pool, created.MatchID); afterCreateEvents != beforeCreateEvents {
		t.Fatalf("absence-policy conflict mutated events: before=%d after=%d", beforeCreateEvents, afterCreateEvents)
	}

	one := uint64(1)
	joined, err := reconstructed.Execute(ctx, "principal-canonical-b", command(
		"canonical-join", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("join canonical match: %v", err)
	}

	two := uint64(2)
	defaulted := contracts.Command{
		ID:               "canonical-defaults",
		SchemaVersion:    application.CommandSchemaV1,
		Type:             application.CommandStartAction,
		MatchID:          created.MatchID,
		ActorPlayerID:    created.PlayerID,
		ExpectedRevision: &two,
		Payload:          json.RawMessage(`{"kind":"MANEUVER"}`),
	}
	defaultResult, err := reconstructed.Execute(ctx, "principal-canonical-a", defaulted)
	if err != nil {
		t.Fatalf("execute defaulted command: %v", err)
	}
	explicitDefaults := defaulted
	explicitDefaults.Payload = json.RawMessage(`{"context":{},"source_ref":"","kind":"MANEUVER"}`)
	defaultReplay, err := application.NewHost(
		registry, coreruntime.NewSequenceIDProvider("canonical-default-unused"), newPostgresStore(t, pool), &postgresCompletedRules{},
	).Execute(ctx, "principal-canonical-a", explicitDefaults)
	if err != nil || !reflect.DeepEqual(defaultResult, defaultReplay) {
		t.Fatalf("schema defaults did not replay original result: result=%+v err=%v", defaultReplay, err)
	}

	three := uint64(3)
	nfcCommand := contracts.Command{
		ID:               "canonical-nfc",
		SchemaVersion:    application.CommandSchemaV1,
		Type:             application.CommandStartAction,
		MatchID:          created.MatchID,
		ActorPlayerID:    created.PlayerID,
		ExpectedRevision: &three,
		Payload:          json.RawMessage(`{"kind":"MANEUVER","source_ref":"e\u0301","context":{"label":"e\u0301","b":2,"a":1}}`),
	}
	nfcResult, err := reconstructed.Execute(ctx, "principal-canonical-a", nfcCommand)
	if err != nil {
		t.Fatalf("execute decomposed NFC command: %v", err)
	}
	nfcRetry := nfcCommand
	nfcRetry.Payload = json.RawMessage(`{"context":{"a":1,"label":"é","b":2},"source_ref":"é","kind":"MANEUVER"}`)
	nfcReplay, err := application.NewHost(
		registry, coreruntime.NewSequenceIDProvider("canonical-nfc-unused"), newPostgresStore(t, pool), &postgresCompletedRules{},
	).Execute(ctx, "principal-canonical-a", nfcRetry)
	if err != nil || !reflect.DeepEqual(nfcResult, nfcReplay) {
		t.Fatalf("NFC-equivalent payload did not replay original result: result=%+v err=%v", nfcReplay, err)
	}

	beforeEvents := countMatchEvents(t, ctx, pool, created.MatchID)
	conflict := nfcCommand
	conflict.Payload = json.RawMessage(`{"kind":"ATTACK","source_ref":"é","context":{"a":1,"b":2,"label":"é"}}`)
	if _, err := reconstructed.Execute(ctx, "principal-canonical-a", conflict); application.CodeOf(err) != application.CodeCommandConflict {
		t.Fatalf("genuinely different canonical request did not conflict: %v", err)
	}
	assertQueryCount(t, ctx, pool,
		"SELECT count(*) FROM core_command_results WHERE principal_id = $1 AND command_id = $2", 1,
		"principal-canonical-a", nfcCommand.ID)
	if afterEvents := countMatchEvents(t, ctx, pool, created.MatchID); afterEvents != beforeEvents {
		t.Fatalf("derived conflict mutated events: before=%d after=%d", beforeEvents, afterEvents)
	}

	four := uint64(4)
	invalidPayloads := []struct {
		id      model.CommandID
		payload json.RawMessage
	}{
		{id: "canonical-unknown", payload: json.RawMessage(`{"kind":"MANEUVER","unknown":true}`)},
		{id: "canonical-duplicate", payload: json.RawMessage(`{"kind":"MANEUVER","kind":"ATTACK"}`)},
		{id: "canonical-nfc-duplicate", payload: json.RawMessage(`{"kind":"MANEUVER","context":{"e\u0301":1,"é":2}}`)},
		{id: "canonical-unsafe-integer", payload: json.RawMessage(`{"kind":"MANEUVER","context":{"value":9007199254740992}}`)},
	}
	for _, invalid := range invalidPayloads {
		candidate := contracts.Command{
			ID:               invalid.id,
			SchemaVersion:    application.CommandSchemaV1,
			Type:             application.CommandStartAction,
			MatchID:          created.MatchID,
			ActorPlayerID:    created.PlayerID,
			ExpectedRevision: &four,
			Payload:          invalid.payload,
		}
		if _, err := reconstructed.Execute(ctx, "principal-canonical-a", candidate); application.CodeOf(err) != application.CodeInvalidCommand {
			t.Fatalf("invalid payload %q was not rejected before persistence: %v", invalid.id, err)
		}
		assertNoDurableCommand(t, ctx, pool, "principal-canonical-a", invalid.id)
	}

	if joined.PlayerID == "" {
		t.Fatal("canonical setup omitted joined player")
	}
}
