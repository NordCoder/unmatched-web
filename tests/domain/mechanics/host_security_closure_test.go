package mechanics_test

import (
	"encoding/json"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
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

func TestTrustedActorHostBindingIsValidatedAndSurvivesResume(t *testing.T) {
	d := effects.Definition{
		ID: "trusted-host-binding", Kind: "linear_stages", RuleID: "FX-060", CapabilityIDs: []string{"CAP-001", "CAP-002"},
		Stages: []effects.Stage{{ID: "choose", Choice: &effects.Choice{
			Kind: "select", Binding: "picked", Visibility: query.Opaque, Owner: player("p1"), Domain: list(lit("left"), lit("right")), EmptyDomain: effects.EmptyReject, ValueType: query.TypeString,
		}}},
	}
	e := must(t, d)
	procedure := model.ProcedureRef{ID: "host-procedure", Kind: d.ID, Bindings: map[string]json.RawMessage{"actor_player_id": raw("p1")}}
	first, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "host-start", Procedure: procedure})
	if err != nil || first.Status != contracts.ResolutionPending {
		t.Fatalf("trusted actor binding did not open interaction: %+v %v", first, err)
	}
	var persistedActor string
	if err = json.Unmarshal(first.PendingInteraction.ResumeProcedure.Bindings["actor_player_id"], &persistedActor); err != nil || persistedActor != "p1" {
		t.Fatalf("trusted actor binding was not persisted: actor=%q err=%v", persistedActor, err)
	}
	resumed, err := e.Resolve(base(), contracts.ResolutionInput{
		CommandID: "host-resume", Procedure: first.PendingInteraction.ResumeProcedure,
		InteractionID: first.PendingInteraction.ID, Choice: raw(firstHandle(t, first)),
	})
	if err != nil || resumed.Status != contracts.ResolutionCompleted {
		t.Fatalf("trusted actor binding did not survive resume: %+v %v", resumed, err)
	}

	invalid, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "host-invalid", Procedure: model.ProcedureRef{ID: "bad", Kind: d.ID, Bindings: map[string]json.RawMessage{"actor_player_id": raw(true)}}})
	if err != nil || invalid.Status != contracts.ResolutionRejected || invalid.RejectionCode != "invalid_captured_bindings" {
		t.Fatalf("invalid host actor accepted: %+v %v", invalid, err)
	}
	mismatch, err := e.Resolve(base(), contracts.ResolutionInput{
		CommandID: "host-mismatch", Procedure: procedure,
		Context: map[string]json.RawMessage{"actor_player_id": raw("p2")},
	})
	if err != nil || mismatch.Status != contracts.ResolutionRejected || mismatch.RejectionCode != "invalid_input_bindings" {
		t.Fatalf("untrusted context actor accepted: %+v %v", mismatch, err)
	}

	reserved := d
	reserved.ID = "reserved-host-redefinition"
	reserved.CapturedBindings = map[string]query.ValueSpec{"actor_player_id": {Type: query.TypePlayer, Visibility: query.Public}}
	if _, err := rules.New([]effects.Definition{reserved}); err == nil {
		t.Fatal("definition redefined reserved actor_player_id binding")
	}
}
