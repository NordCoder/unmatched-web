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

func TestProjectionNeverSerializesCanonicalFighterState(t *testing.T) {
	e := must(t, choiceDefinition())
	s := base()
	s.Players["p2"] = model.PlayerState{ID: "p2"}
	s.Fighters["alice"] = model.RuntimeObject{DefinitionID: "alice", OwnerID: "p1", ControllerID: "p1", State: map[string]any{"secret": "fighter-secret-a"}}
	s.Fighters["bob"] = model.RuntimeObject{DefinitionID: "bob", OwnerID: "p2", ControllerID: "p2", State: map[string]any{"secret": "fighter-secret-b"}}
	for _, viewer := range []model.PlayerID{"p1", "p2"} {
		projection, err := e.Project(s, viewer)
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range [][]byte{[]byte("fighter-secret-a"), []byte("fighter-secret-b")} {
			if bytes.Contains(projection, forbidden) {
				t.Fatalf("projection leak for %s %q: %s", viewer, forbidden, projection)
			}
		}
		var decoded struct {
			Fighters []map[string]any `json:"fighters"`
		}
		if err = json.Unmarshal(projection, &decoded); err != nil {
			t.Fatal(err)
		}
		if len(decoded.Fighters) != 2 {
			t.Fatalf("fighter DTO count %d", len(decoded.Fighters))
		}
		for _, fighter := range decoded.Fighters {
			if _, exposed := fighter["state"]; exposed {
				t.Fatalf("canonical fighter state exposed: %v", fighter)
			}
		}
	}
}

func TestRuntimeBindingContractsRejectUntrustedData(t *testing.T) {
	payloadSpec := query.ValueSpec{Type: query.TypeObject, Visibility: query.Public, Fields: map[string]query.ValueSpec{
		"safe": {Type: query.TypeString, Visibility: query.Public},
	}}
	capturedDef := effects.Definition{
		ID: "captured-contract", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"},
		CapturedBindings: map[string]query.ValueSpec{"payload": payloadSpec},
		Stages: []effects.Stage{{ID: "emit", Operations: []operations.Definition{{
			ID: "emit", Kind: operations.EmitEvent,
			Arguments: map[string]query.Expr{"event_type": lit("captured.public"), "payload": ref(query.Captured, query.TypeObject, "payload")},
		}}}},
	}
	e := must(t, capturedDef)
	cases := []struct {
		name     string
		bindings map[string]json.RawMessage
	}{
		{name: "unknown binding", bindings: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "ok"}), "extra": raw(true)}},
		{name: "missing binding", bindings: map[string]json.RawMessage{}},
		{name: "wrong nested type", bindings: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": true})}},
		{name: "undeclared nested field", bindings: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "ok", "secret": "leak"})}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "bad", Procedure: model.ProcedureRef{ID: "p", Kind: capturedDef.ID, Bindings: tc.bindings}})
			if err != nil {
				t.Fatal(err)
			}
			if out.Status != contracts.ResolutionRejected || out.RejectionCode != "invalid_captured_bindings" || len(out.Events) != 0 {
				t.Fatalf("unexpected outcome %+v", out)
			}
		})
	}
	good, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "good", Procedure: model.ProcedureRef{ID: "p", Kind: capturedDef.ID, Bindings: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "ok"})}}})
	if err != nil || good.Status != contracts.ResolutionCompleted || countType(good.Events, "captured.public") != 1 {
		t.Fatalf("valid captured binding failed: %+v %v", good, err)
	}
	for _, ev := range good.Events {
		if bytes.Contains(ev.PublicPayload, []byte("secret")) || bytes.Contains(ev.PublicPayload, []byte("leak")) {
			t.Fatalf("untrusted field reached event: %s", ev.PublicPayload)
		}
	}

	inputDef := capturedDef
	inputDef.ID = "input-contract"
	inputDef.CapturedBindings = nil
	inputDef.InputBindings = map[string]query.ValueSpec{"payload": payloadSpec}
	inputDef.Stages[0].Operations[0].Arguments["payload"] = ref(query.Input, query.TypeObject, "payload")
	inputEngine := must(t, inputDef)
	inputCases := []struct {
		name    string
		context map[string]json.RawMessage
	}{
		{name: "unknown input", context: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "ok"}), "extra": raw(true)}},
		{name: "missing input", context: map[string]json.RawMessage{}},
		{name: "wrong input type", context: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": true})}},
		{name: "undeclared input field", context: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "ok", "secret": "leak"})}},
	}
	for _, tc := range inputCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := inputEngine.Resolve(base(), contracts.ResolutionInput{CommandID: "bad-input", Procedure: model.ProcedureRef{ID: "p", Kind: inputDef.ID}, Context: tc.context})
			if err != nil {
				t.Fatal(err)
			}
			if out.Status != contracts.ResolutionRejected || out.RejectionCode != "invalid_input_bindings" || len(out.Events) != 0 {
				t.Fatalf("unexpected input outcome %+v", out)
			}
		})
	}

	resumeDef := effects.Definition{
		ID: "input-resume-contract", Kind: "linear_stages", RuleID: "FX-060", CapabilityIDs: []string{"CAP-001", "CAP-002"},
		InputBindings: map[string]query.ValueSpec{"payload": payloadSpec},
		Stages: []effects.Stage{{ID: "choose", Choice: &effects.Choice{
			Kind: "select", Binding: "picked", Visibility: query.Opaque, Owner: player("p1"), Domain: list(lit("left"), lit("right")), EmptyDomain: effects.EmptyReject, ValueType: query.TypeString,
		}, Operations: []operations.Definition{{ID: "emit", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("input.resume"), "payload": ref(query.Input, query.TypeObject, "payload")}}}}},
	}
	resumeEngine := must(t, resumeDef)
	transportPayload := raw(map[string]any{"safe": "ok"})
	first, err := resumeEngine.Resolve(base(), contracts.ResolutionInput{CommandID: "input-start", Procedure: model.ProcedureRef{ID: "p", Kind: resumeDef.ID, Bindings: map[string]json.RawMessage{"payload": transportPayload}}, Context: map[string]json.RawMessage{"payload": transportPayload}})
	if err != nil || first.Status != contracts.ResolutionPending {
		t.Fatalf("duplicated Core context did not open interaction: %+v %v", first, err)
	}
	mismatch, err := resumeEngine.Resolve(base(), contracts.ResolutionInput{CommandID: "input-mismatch", Procedure: model.ProcedureRef{ID: "other", Kind: resumeDef.ID, Bindings: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "captured"})}}, Context: map[string]json.RawMessage{"payload": raw(map[string]any{"safe": "input"})}})
	if err != nil {
		t.Fatal(err)
	}
	if mismatch.Status != contracts.ResolutionRejected || mismatch.RejectionCode != "binding_mismatch" {
		t.Fatalf("mismatched duplicated context accepted: %+v", mismatch)
	}
	resumed, err := resumeEngine.Resolve(base(), contracts.ResolutionInput{CommandID: "input-resume", Procedure: first.PendingInteraction.ResumeProcedure, InteractionID: first.PendingInteraction.ID, Choice: raw(firstHandle(t, first))})
	if err != nil || resumed.Status != contracts.ResolutionCompleted || countType(resumed.Events, "input.resume") != 1 {
		t.Fatalf("captured input fallback failed: %+v %v", resumed, err)
	}
}

func TestWave1RejectsOpenBindingsAndMultiSelect(t *testing.T) {
	publicList := query.ValueSpec{Type: query.TypeList, Element: &query.ValueSpec{Type: query.TypeString, Visibility: query.Public}, Visibility: query.Public}
	cases := []effects.Definition{
		{ID: "open-external", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, CapturedBindings: map[string]query.ValueSpec{"payload": {Type: query.TypeAny, Visibility: query.Public}}, Stages: []effects.Stage{{ID: "s"}}},
		choiceFromCaptured("multi-unsupported", publicList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptyReject, Multi: true}, nil),
		choiceFromCaptured("multi-empty-unsupported", publicList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptyBindEmpty, Multi: true}, nil),
	}
	for _, d := range cases {
		if _, err := rules.New([]effects.Definition{d}); err == nil {
			t.Fatalf("accepted %s", d.ID)
		}
	}
}
