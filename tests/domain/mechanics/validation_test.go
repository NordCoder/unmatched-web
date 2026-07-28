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
	frame := effects.BuildCheckpoint("during", 2, []effects.QueueDefinition{{ID: "late", Scope: "x", Priority: 1, SourceOrder: 2, Cancelable: true, Operations: []operations.Definition{{ID: "l", Kind: operations.Impossible}}}, {ID: "first", Scope: "x", Priority: 2, SourceOrder: 5, Operations: []operations.Definition{{ID: "f", Kind: operations.Impossible}}}, {ID: "middle", Scope: "y", Priority: 1, SourceOrder: 1, Cancelable: true, Operations: []operations.Definition{{ID: "m", Kind: operations.Impossible}}}})
	if !reflect.DeepEqual(frame.Cancel("x"), []string{"late"}) {
		t.Fatal(frame)
	}
	state := effects.State{DefinitionID: "d", Phase: effects.PhaseCheckpoint, Checkpoint: frame, Status: "running", Captured: map[string]json.RawMessage{}, Results: map[string]json.RawMessage{}, Choices: map[string]json.RawMessage{}}
	r, err := effects.Encode(model.ProcedureRef{ID: "p", Kind: "d"}, state, "s")
	if err != nil {
		t.Fatal(err)
	}
	restored, err := effects.Decode(r)
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
		if q.Status != "executing" {
			t.Fatalf("premature status %s", q.Status)
		}
	}
	if !reflect.DeepEqual(got, []string{"first", "middle"}) {
		t.Fatalf("order %v", got)
	}
}

func TestStatePathAuthorityRejectsSelfDeclaredPublic(t *testing.T) {
	cases := []effects.Definition{
		defWithCondition("unknown", query.Expr{Kind: query.Reference, Source: query.State, Path: []string{"fighters", "alice", "state", "unknown"}, ValueType: query.TypeBool}),
		{ID: "private-state", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, Stages: []effects.Stage{{ID: "s", Operations: []operations.Definition{{ID: "emit", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("bad"), "payload": obj("zones", query.Expr{Kind: query.Reference, Source: query.State, Path: []string{"players", "p2", "private_zones"}, ValueType: query.TypeObject, Visibility: query.Public})}}}}}},
	}
	for _, d := range cases {
		if _, err := rules.New([]effects.Definition{d}); err == nil {
			t.Fatalf("accepted %s", d.ID)
		}
	}
	known := defWithCondition("known", query.Expr{Kind: query.Equal, Args: []query.Expr{ref(query.State, query.TypeString, "fighters", "alice", "state", "mode"), lit("before")}})
	if _, err := rules.New([]effects.Definition{known}); err != nil {
		t.Fatalf("known path rejected: %v", err)
	}
}

func TestRecursiveReferenceTyping(t *testing.T) {
	baseDef := func(refExpr query.Expr) effects.Definition {
		return effects.Definition{ID: "nested", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, Stages: []effects.Stage{{ID: "s", Operations: []operations.Definition{{ID: "set", Kind: operations.SetFighterState, ResultBinding: "setr", Arguments: map[string]query.Expr{"fighter": fighter("alice"), "key": lit("mode"), "value": lit("after")}}, {ID: "emit", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("ok"), "payload": obj("v", refExpr)}}}}}}
	}
	badField := baseDef(ref(query.Result, query.TypeString, "setr", "value", "missing"))
	if _, err := rules.New([]effects.Definition{badField}); err == nil {
		t.Fatal("accepted nonexistent nested field")
	}
	badOverride := baseDef(ref(query.Result, query.TypeBool, "setr", "value", "value"))
	if _, err := rules.New([]effects.Definition{badOverride}); err == nil {
		t.Fatal("accepted nested type override")
	}
	good := baseDef(ref(query.Result, query.TypeString, "setr", "value", "value"))
	if _, err := rules.New([]effects.Definition{good}); err != nil {
		t.Fatalf("valid nested path rejected: %v", err)
	}
}

func TestChoiceTypingVisibilityAndPresence(t *testing.T) {
	privateList := query.ValueSpec{Type: query.TypeList, Element: &query.ValueSpec{Type: query.TypeString, Visibility: query.OwnerPrivate}, Visibility: query.OwnerPrivate}
	publicList := query.ValueSpec{Type: query.TypeList, Element: &query.ValueSpec{Type: query.TypeString, Visibility: query.Public}, Visibility: query.Public}
	cases := []effects.Definition{
		choiceFromCaptured("wrong-type", publicList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptyReject, ValueType: query.TypeFighter}, nil),
		choiceFromCaptured("weak-vis", privateList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptyReject, ValueType: query.TypeString}, nil),
		choiceFromCaptured("scalar-empty", publicList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptyBindEmpty, ValueType: query.TypeString}, nil),
		choiceFromCaptured("absent-read", publicList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptySkipStage, ValueType: query.TypeString}, []operations.Definition{{ID: "use", Kind: operations.BindValue, Arguments: map[string]query.Expr{"value": ref(query.Choice, query.TypeString, "pick")}}}),
	}
	for _, d := range cases {
		if _, err := rules.New([]effects.Definition{d}); err == nil {
			t.Fatalf("accepted %s", d.ID)
		}
	}
	defExpr := lit("fallback")
	good := choiceFromCaptured("default", publicList, effects.Choice{Kind: "x", Binding: "pick", Visibility: query.Public, Owner: player("p1"), Domain: ref(query.Captured, query.TypeList, "opts"), EmptyDomain: effects.EmptyBindDefault, ValueType: query.TypeString, Default: &defExpr}, []operations.Definition{{ID: "use", Kind: operations.BindValue, Arguments: map[string]query.Expr{"value": ref(query.Choice, query.TypeString, "pick")}}})
	if _, err := rules.New([]effects.Definition{good}); err != nil {
		t.Fatalf("valid default rejected: %v", err)
	}
}

func TestPureReducerDoesNotAliasInput(t *testing.T) {
	d := effects.Definition{ID: "pure", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, Stages: []effects.Stage{{ID: "s", Operations: []operations.Definition{{ID: "set", Kind: operations.SetFighterState, Arguments: map[string]query.Expr{"fighter": fighter("alice"), "key": lit("marked"), "value": lit(true)}}}}}}
	e := must(t, d)
	s := base()
	s.Fighters["alice"] = model.RuntimeObject{DefinitionID: "alice", State: map[string]any{"marked": false}}
	before, _ := json.Marshal(s)
	out, err := e.Resolve(s, contracts.ResolutionInput{CommandID: "c", Procedure: model.ProcedureRef{ID: "p", Kind: d.ID}})
	if err != nil {
		t.Fatal(err)
	}
	a := reduceAll(t, e, s, out.Events)
	afterInput, _ := json.Marshal(s)
	if !reflect.DeepEqual(before, afterInput) {
		t.Fatal("input mutated")
	}
	if a.Fighters["alice"].State["marked"] != true {
		t.Fatal("event was not reduced")
	}
	a.Fighters["alice"].State["marked"] = false
	if s.Fighters["alice"].State["marked"] != false {
		t.Fatal("output aliases input")
	}
	b := reduceAll(t, e, s, out.Events)
	c := reduceAll(t, e, s, out.Events)
	if b.Fighters["alice"].State["marked"] != true || c.Fighters["alice"].State["marked"] != true {
		t.Fatal("repeat reduction did not apply event")
	}
	bb, _ := json.Marshal(b)
	cb, _ := json.Marshal(c)
	if !reflect.DeepEqual(bb, cb) {
		t.Fatal("repeat reduction diverged")
	}
	inputAgain, _ := json.Marshal(s)
	if !reflect.DeepEqual(before, inputAgain) {
		t.Fatal("reduction or output mutation changed input")
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

func defWithCondition(id string, c query.Expr) effects.Definition {
	return effects.Definition{ID: id, Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001"}, Stages: []effects.Stage{{ID: "s", Condition: &c}}}
}
func choiceFromCaptured(id string, spec query.ValueSpec, ch effects.Choice, ops []operations.Definition) effects.Definition {
	return effects.Definition{ID: id, Kind: "linear_stages", RuleID: "FX-060", CapabilityIDs: []string{"CAP-001", "CAP-002"}, CapturedBindings: map[string]query.ValueSpec{"opts": spec}, Stages: []effects.Stage{{ID: "s", Choice: &ch, Operations: ops}}}
}
