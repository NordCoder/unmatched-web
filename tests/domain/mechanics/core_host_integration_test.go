package mechanics_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"github.com/NordCoder/unmatched-web/internal/domain/rules"
	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func TestConcreteRulesActorOwnerWorksForBothCoreSeats(t *testing.T) {
	registry := coreruntime.NewMemoryDefinitionRegistry()
	bundle := hostIntegrationBundle()
	if err := registry.Register("host-integration@v1", bundle); err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	engine, err := rules.New([]effects.Definition{hostIntegrationRulesDefinition()})
	if err != nil {
		t.Fatalf("construct concrete rules: %v", err)
	}
	store := persistence.NewMemoryStore()
	host := application.NewHost(registry, coreruntime.NewSequenceIDProvider("host"), store, engine)
	ctx := context.Background()

	created, err := host.Execute(ctx, "principal-a", hostIntegrationCommand(
		"create", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "host-integration@v1", FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("create match: %v", err)
	}
	one := uint64(1)
	joined, err := host.Execute(ctx, "principal-b", hostIntegrationCommand(
		"join", application.CommandJoinMatch, created.MatchID, "", &one,
		application.JoinMatchPayload{FighterDefinition: "mirror"},
	))
	if err != nil {
		t.Fatalf("join match: %v", err)
	}

	exerciseCoreActorChoice(t, ctx, host, created.MatchID, "principal-a", created.PlayerID, 2, "first")
	exerciseCoreActorChoice(t, ctx, host, created.MatchID, "principal-b", joined.PlayerID, 4, "second")
}

func exerciseCoreActorChoice(
	t *testing.T,
	ctx context.Context,
	host *application.Host,
	matchID model.MatchID,
	principal model.PrincipalID,
	actor model.PlayerID,
	startRevision uint64,
	commandPrefix string,
) {
	t.Helper()
	started, err := host.Execute(ctx, principal, hostIntegrationCommand(
		commandPrefix+"-start", application.CommandStartAction, matchID, actor, &startRevision,
		application.StartActionPayload{Kind: application.ActionScheme, SourceRef: "card:host-integration"},
	))
	if err != nil {
		t.Fatalf("Core StartAction rejected actor %q: %v", actor, err)
	}
	projected := started.Projection.PendingInteraction
	if projected == nil {
		t.Fatalf("concrete Rules did not open interaction for actor %q: %+v", actor, started)
	}
	if projected.OwnerPlayerID != actor {
		t.Fatalf("projected owner %q, want actor %q", projected.OwnerPlayerID, actor)
	}

	state, err := host.State(matchID)
	if err != nil {
		t.Fatalf("load pending state: %v", err)
	}
	pending := state.Game.Resolver.PendingInteraction
	if pending == nil {
		t.Fatal("pending interaction was not persisted")
	}
	if pending.OwnerPlayerID != actor {
		t.Fatalf("persisted interaction owner %q, want %q", pending.OwnerPlayerID, actor)
	}
	var persistedActor model.PlayerID
	if err = json.Unmarshal(pending.ResumeProcedure.Bindings[effects.TrustedActorBinding], &persistedActor); err != nil {
		t.Fatalf("decode persisted actor: %v", err)
	}
	if persistedActor != actor {
		t.Fatalf("persisted actor %q, want %q", persistedActor, actor)
	}

	resumeRevision := startRevision + 1
	completed, err := host.Execute(ctx, principal, hostIntegrationCommand(
		commandPrefix+"-resume", application.CommandSubmitChoice, matchID, actor, &resumeRevision,
		application.SubmitChoicePayload{
			InteractionID: projected.ID,
			Choice:        hostIntegrationChoice(t, projected.LegalDomain),
		},
	))
	if err != nil {
		t.Fatalf("Core SubmitChoice could not resume actor %q: %v", actor, err)
	}
	if completed.Projection.PendingInteraction != nil || completed.Projection.BlockedByInteraction {
		t.Fatalf("action remained pending after actor %q choice: %+v", actor, completed)
	}
}

func hostIntegrationRulesDefinition() effects.Definition {
	return effects.Definition{
		ID: application.ActionScheme, Kind: "linear_stages", RuleID: "FX-060", CapabilityIDs: []string{"CAP-001", "CAP-002"},
		Stages: []effects.Stage{{ID: "choose", Choice: &effects.Choice{
			Kind: "select", Binding: "picked", Visibility: query.Opaque,
			Owner: effects.ActorOwner(),
			Domain: query.Expr{Kind: query.List, Args: []query.Expr{
				{Kind: query.Literal, Value: "left"},
				{Kind: query.Literal, Value: "right"},
			}},
			EmptyDomain: effects.EmptyReject, ValueType: query.TypeString,
		}}},
	}
}

func hostIntegrationBundle() coreruntime.DefinitionBundle {
	fighter := coreruntime.FighterDefinition{ID: "mirror"}
	return coreruntime.DefinitionBundle{
		Ref: model.DefinitionRef{
			RulesetVersion:         "host-integration-rules/v1",
			CapabilityRegistry:     "host-integration-capabilities/v1",
			FighterManifestDigests: map[string]string{"mirror": coreruntime.FighterDefinitionDigest(fighter)},
			CardManifestDigests:    map[string]string{},
		},
		Fighters: map[model.DefinitionID]coreruntime.FighterDefinition{"mirror": fighter},
	}
}

func hostIntegrationCommand(id, kind string, matchID model.MatchID, actor model.PlayerID, revision *uint64, payload any) contracts.Command {
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

func hostIntegrationChoice(t *testing.T, domain json.RawMessage) json.RawMessage {
	t.Helper()
	var projected struct {
		Options []struct {
			Handle string `json:"handle"`
		} `json:"options"`
	}
	if err := json.Unmarshal(domain, &projected); err != nil {
		t.Fatalf("decode legal domain: %v", err)
	}
	if len(projected.Options) == 0 || projected.Options[0].Handle == "" {
		t.Fatalf("legal domain has no opaque handle: %s", domain)
	}
	encoded, err := json.Marshal(projected.Options[0].Handle)
	if err != nil {
		t.Fatalf("encode choice handle: %v", err)
	}
	return encoded
}
