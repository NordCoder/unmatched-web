package mechanics_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"github.com/NordCoder/unmatched-web/internal/domain/rules"
)

func TestWave1RejectsPrivateStateAndPrivateControlFlow(t *testing.T) {
	privateState := effects.Definition{
		ID: "private-state-condition", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"},
		Stages: []effects.Stage{{ID: "leak", Condition: &query.Expr{
			Kind: query.Equal,
			Args: []query.Expr{
				{Kind: query.Count, Args: []query.Expr{{Kind: query.Reference, Source: query.State, Path: []string{"players", "p2", "private_zones"}, ValueType: query.TypeObject, Visibility: query.OwnerPrivate}}},
				lit(0),
			},
		}}},
	}
	if _, err := rules.New([]effects.Definition{privateState}); err == nil {
		t.Fatal("private player zones remained readable through the generic state schema")
	}

	privateCondition := effects.Definition{
		ID: "private-captured-condition", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"},
		CapturedBindings: map[string]query.ValueSpec{"secret": {Type: query.TypeBool, Visibility: query.OwnerPrivate}},
		Stages:           []effects.Stage{{ID: "leak", Condition: &query.Expr{Kind: query.Reference, Source: query.Captured, Path: []string{"secret"}, ValueType: query.TypeBool, Visibility: query.OwnerPrivate}}},
	}
	if _, err := rules.New([]effects.Definition{privateCondition}); err == nil {
		t.Fatal("private-tainted stage condition was accepted")
	}

	privatePrerequisite := privateCondition
	privatePrerequisite.ID = "private-captured-prerequisite"
	privatePrerequisite.Stages[0].Condition = nil
	privatePrerequisite.Stages[0].Prerequisites = []query.Expr{{Kind: query.Reference, Source: query.Captured, Path: []string{"secret"}, ValueType: query.TypeBool, Visibility: query.OwnerPrivate}}
	if _, err := rules.New([]effects.Definition{privatePrerequisite}); err == nil {
		t.Fatal("private-tainted stage prerequisite was accepted")
	}
}

func TestActorOwnerSelectorUsesTrustedBindingAndSurvivesResume(t *testing.T) {
	d := actorChoiceDefinition("trusted-host-binding")
	e := must(t, d)
	procedure := actorProcedure("host-procedure", d.ID, "p1")
	first, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "host-start", Procedure: procedure})
	if err != nil || first.Status != contracts.ResolutionPending {
		t.Fatalf("trusted actor binding did not open interaction: %+v %v", first, err)
	}
	if first.PendingInteraction.OwnerPlayerID != "p1" {
		t.Fatalf("interaction owner %q, want p1", first.PendingInteraction.OwnerPlayerID)
	}
	var persistedActor string
	if err = json.Unmarshal(first.PendingInteraction.ResumeProcedure.Bindings[effects.TrustedActorBinding], &persistedActor); err != nil || persistedActor != "p1" {
		t.Fatalf("trusted actor binding was not persisted: actor=%q err=%v", persistedActor, err)
	}
	resumed, err := e.Resolve(base(), contracts.ResolutionInput{
		CommandID: "host-resume", Procedure: first.PendingInteraction.ResumeProcedure,
		InteractionID: first.PendingInteraction.ID, Choice: raw(firstHandle(t, first)),
	})
	if err != nil || resumed.Status != contracts.ResolutionCompleted {
		t.Fatalf("trusted actor binding did not survive resume: %+v %v", resumed, err)
	}
}

func TestActorOwnerSelectorRejectsMissingWrongAndInjectedActor(t *testing.T) {
	d := actorChoiceDefinition("actor-validation")
	e := must(t, d)

	missing, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "missing", Procedure: model.ProcedureRef{ID: "missing", Kind: d.ID}})
	if err != nil || missing.Status != contracts.ResolutionRejected || missing.RejectionCode != "invalid_captured_bindings" {
		t.Fatalf("missing trusted actor accepted: %+v %v", missing, err)
	}

	invalid, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "wrong", Procedure: model.ProcedureRef{ID: "wrong", Kind: d.ID, Bindings: map[string]json.RawMessage{effects.TrustedActorBinding: raw(true)}}})
	if err != nil || invalid.Status != contracts.ResolutionRejected || invalid.RejectionCode != "invalid_captured_bindings" {
		t.Fatalf("invalid host actor accepted: %+v %v", invalid, err)
	}

	empty, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "empty", Procedure: actorProcedure("empty", d.ID, "")})
	if err != nil || empty.Status != contracts.ResolutionRejected || empty.RejectionCode != "invalid_actor_owner" {
		t.Fatalf("empty host actor accepted: %+v %v", empty, err)
	}

	injected, err := e.Resolve(base(), contracts.ResolutionInput{
		CommandID: "injected", Procedure: model.ProcedureRef{ID: "injected", Kind: d.ID},
		Context: map[string]json.RawMessage{effects.TrustedActorBinding: raw("p1")},
	})
	if err != nil || injected.Status != contracts.ResolutionRejected || injected.RejectionCode != "invalid_captured_bindings" {
		t.Fatalf("context actor replaced missing trusted actor: %+v %v", injected, err)
	}

	procedure := actorProcedure("mismatch", d.ID, "p1")
	mismatch, err := e.Resolve(base(), contracts.ResolutionInput{
		CommandID: "mismatch", Procedure: procedure,
		Context: map[string]json.RawMessage{effects.TrustedActorBinding: raw("p2")},
	})
	if err != nil || mismatch.Status != contracts.ResolutionRejected || mismatch.RejectionCode != "invalid_input_bindings" {
		t.Fatalf("untrusted context actor accepted: %+v %v", mismatch, err)
	}
}

func TestWave1RejectsLegacyAndTaintedChoiceOwners(t *testing.T) {
	legacy := actorChoiceDefinition("legacy-owner")
	legacy.Stages[0].Choice.Owner = query.Expr{Kind: query.Literal, Value: "p1", ValueType: query.TypePlayer}

	privateOwner := actorChoiceDefinition("private-owner")
	privateOwner.CapturedBindings = map[string]query.ValueSpec{"candidate": {Type: query.TypePlayer, Visibility: query.OwnerPrivate}}
	privateOwner.Stages[0].Choice.Owner = query.Expr{Kind: query.Reference, Source: query.Captured, Path: []string{"candidate"}, ValueType: query.TypePlayer, Visibility: query.OwnerPrivate}

	opaqueOwner := actorChoiceDefinition("opaque-owner")
	opaqueOwner.CapturedBindings = map[string]query.ValueSpec{"candidate": {Type: query.TypePlayer, Visibility: query.Opaque}}
	opaqueOwner.Stages[0].Choice.Owner = query.Expr{Kind: query.Reference, Source: query.Captured, Path: []string{"candidate"}, ValueType: query.TypePlayer, Visibility: query.Opaque}

	arbitraryActorReference := actorChoiceDefinition("actor-reference-outside-selector")
	arbitraryActorReference.Stages[0].Operations = []operations.Definition{{
		ID: "leak-actor", Kind: operations.EmitEvent,
		Arguments: map[string]query.Expr{
			"event_type": lit("bad.actor"),
			"payload":    obj("actor", ref(query.Captured, query.TypePlayer, effects.TrustedActorBinding)),
		},
	}}

	for _, d := range []effects.Definition{legacy, privateOwner, opaqueOwner, arbitraryActorReference} {
		if _, err := rules.New([]effects.Definition{d}); err == nil {
			t.Fatalf("accepted unsupported owner for %s", d.ID)
		}
	}

	legacyJSON := []byte(`{"id":"legacy-wire","kind":"linear_stages","rule_id":"FX-060","capability_ids":["CAP-001","CAP-002"],"stages":[{"id":"choose","choice":{"kind":"select","binding":"picked","visibility":"opaque","owner":{"kind":"literal","value":"p1","value_type":"player_ref"},"domain":{"kind":"list","args":[{"kind":"literal","value":"left"}]},"empty_domain":"reject","value_type":"string"}}]}`)
	var decoded effects.Definition
	if err := json.Unmarshal(legacyJSON, &decoded); err != nil {
		t.Fatalf("legacy owner wire shape must remain readable for explicit rejection: %v", err)
	}
	if _, err := rules.New([]effects.Definition{decoded}); err == nil {
		t.Fatal("legacy owner expression from serialized definition was accepted")
	}
}

func TestActorOwnerSelectorCannotBeRedeclaredAndSerializesClosedContract(t *testing.T) {
	d := actorChoiceDefinition("owner-wire")
	encoded, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(encoded, []byte(`"owner_selector":"actor"`)) || bytes.Contains(encoded, []byte(`"owner":`)) {
		t.Fatalf("definition did not serialize closed actor selector: %s", encoded)
	}

	var restored effects.Definition
	if err = json.Unmarshal(encoded, &restored); err != nil {
		t.Fatal(err)
	}
	if !effects.IsActorOwner(restored.Stages[0].Choice.Owner) {
		t.Fatalf("actor selector did not survive definition serialization: %+v", restored.Stages[0].Choice.Owner)
	}

	captured := d
	captured.ID = "reserved-captured-redefinition"
	captured.CapturedBindings = map[string]query.ValueSpec{effects.TrustedActorBinding: {Type: query.TypePlayer, Visibility: query.Public}}
	input := d
	input.ID = "reserved-input-redefinition"
	input.InputBindings = map[string]query.ValueSpec{effects.TrustedActorBinding: {Type: query.TypePlayer, Visibility: query.Public}}
	for _, candidate := range []effects.Definition{captured, input} {
		if _, err := rules.New([]effects.Definition{candidate}); err == nil {
			t.Fatalf("definition redefined reserved actor binding: %s", candidate.ID)
		}
	}
}

func actorChoiceDefinition(id string) effects.Definition {
	return effects.Definition{
		ID: id, Kind: "linear_stages", RuleID: "FX-060", CapabilityIDs: []string{"CAP-001", "CAP-002"},
		Stages: []effects.Stage{{ID: "choose", Choice: &effects.Choice{
			Kind: "select", Binding: "picked", Visibility: query.Opaque,
			Owner: effects.ActorOwner(), Domain: list(lit("left"), lit("right")),
			EmptyDomain: effects.EmptyReject, ValueType: query.TypeString,
		}}},
	}
}
