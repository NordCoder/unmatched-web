package mechanics_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/history"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"github.com/NordCoder/unmatched-web/internal/domain/rules"
)

func TestCheckpointOrderingCancellationAndSerialization(t *testing.T) {
	frame := effects.BuildCheckpoint("during", 2, []effects.QueueDefinition{
		{ID: "late", Scope: "x", Priority: 1, SourceOrder: 2, Cancelable: true, Operations: []operations.Definition{{ID: "l", Kind: operations.Impossible}}},
		{ID: "first", Scope: "x", Priority: 2, SourceOrder: 5, Operations: []operations.Definition{{ID: "f", Kind: operations.Impossible}}},
		{ID: "middle", Scope: "y", Priority: 1, SourceOrder: 1, Cancelable: true, Operations: []operations.Definition{{ID: "m", Kind: operations.Impossible}}},
	})
	if !reflect.DeepEqual(frame.Cancel("x"), []string{"late"}) {
		t.Fatal(frame)
	}
	state := effects.State{DefinitionID: "d", Phase: effects.PhaseCheckpoint, Checkpoint: frame, Status: "running", Captured: map[string]json.RawMessage{}, Results: map[string]json.RawMessage{}, Choices: map[string]json.RawMessage{}}
	ref, err := effects.Encode(model.ProcedureRef{ID: "p", Kind: "d"}, state, "s")
	if err != nil {
		t.Fatal(err)
	}
	restored, err := effects.Decode(ref)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for {
		q := restored.Checkpoint.Next()
		if q == nil {
			break
		}
		got = append(got, q.ID)
	}
	if !reflect.DeepEqual(got, []string{"first", "middle"}) {
		t.Fatalf("order %v", got)
	}
}

func TestTypedValidationRejectsUnsafeDefinitions(t *testing.T) {
	cases := []effects.Definition{
		{ID: "missing-policy", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001", "CAP-002"}, Stages: []effects.Stage{{ID: "s", Choice: &effects.Choice{Kind: "x", Binding: "x", Visibility: query.Opaque, Owner: player("p1"), Domain: list(lit("x")), ValueType: query.TypeString}}}},
		{ID: "private-event", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, CapturedBindings: map[string]query.ValueSpec{"secret": {Type: query.TypeObject, Visibility: query.OwnerPrivate}}, Stages: []effects.Stage{{ID: "s", Operations: []operations.Definition{{ID: "emit", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("bad"), "payload": ref(query.Captured, query.TypeObject, "secret")}}}}}},
		{ID: "unknown-cap", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-999"}, Stages: []effects.Stage{{ID: "s"}}},
		{ID: "untyped-ref", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, Stages: []effects.Stage{{ID: "s", Condition: &query.Expr{Kind: query.Reference, Source: query.State, Path: []string{"fighters", "alice", "state", "mode"}}}}},
	}
	for _, d := range cases {
		if _, err := rules.New([]effects.Definition{d}); err == nil {
			t.Fatalf("accepted %s", d.ID)
		}
	}
}

func TestActualEventsBuildProvenanceLedger(t *testing.T) {
	d := effects.Definition{ID: "prov", Kind: "linear_stages", RuleID: "FX-021", CapabilityIDs: []string{"CAP-001", "CAP-004"}, Stages: []effects.Stage{{ID: "s", Operations: []operations.Definition{{ID: "set", Kind: operations.SetFighterState, CauseKind: "combat", Arguments: map[string]query.Expr{"fighter": fighter("alice"), "key": lit("marked"), "value": lit(true)}}}}}}
	e := must(t, d)
	s := base()
	s.Fighters["alice"] = model.RuntimeObject{DefinitionID: "alice", State: map[string]any{}}
	out, err := e.Resolve(s, contracts.ResolutionInput{CommandID: "c", Procedure: model.ProcedureRef{ID: "p", Kind: d.ID, SourceRef: "card:source"}})
	if err != nil {
		t.Fatal(err)
	}
	entries := history.Build(out.Events).ByParticipant("alice")
	if len(entries) != 1 || entries[0].CauseKind != "combat" || entries[0].SourceRef != "card:source" || entries[0].OperationID != "set" {
		t.Fatalf("entries %+v", entries)
	}
}
