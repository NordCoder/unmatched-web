package mechanics_test

import (
	"bytes"
	"encoding/json"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/history"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"github.com/NordCoder/unmatched-web/internal/domain/rules"
	"reflect"
	"testing"
)

func TestOrderedCurrentAndCapturedBindings(t *testing.T) {
	d := effects.Definition{ID: "ordered", Kind: "linear_stages", RuleID: "FX-010", CapabilityIDs: []string{"CAP-001", "CAP-018"}, Stages: []effects.Stage{{ID: "set", Operations: []operations.Definition{{ID: "set-mode", Kind: operations.SetFighterState, ResultBinding: "set_mode", Arguments: map[string]query.Expr{"fighter": lit("alice"), "key": lit("mode"), "value": lit("after")}}}}, {ID: "read", Condition: ptr(query.Expr{Kind: query.Equal, Args: []query.Expr{ref(query.State, "fighters", "alice", "state", "mode"), lit("after")}}), Operations: []operations.Definition{{ID: "emit", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("test.read"), "payload": query.Expr{Kind: query.Object, Fields: map[string]query.Expr{"captured": ref(query.Captured, "before"), "current": ref(query.State, "fighters", "alice", "state", "mode"), "prior": ref(query.Result, "set_mode", "value", "value")}}}}}}}}
	e := must(t, d)
	s := base()
	s.Fighters["alice"] = model.RuntimeObject{State: map[string]any{"mode": "before"}}
	o, err := e.Resolve(s, contracts.ResolutionInput{CommandID: "cmd", Procedure: model.ProcedureRef{ID: "p", Kind: "ordered", SourceRef: "card:x", Bindings: map[string]json.RawMessage{"before": raw("before")}}})
	if err != nil {
		t.Fatal(err)
	}
	if o.Status != contracts.ResolutionCompleted {
		t.Fatal(o.Status)
	}
	var p map[string]any
	for _, x := range o.Events {
		if x.Type == "test.read" {
			json.Unmarshal(x.PublicPayload, &p)
		}
	}
	want := map[string]any{"captured": "before", "current": "after", "prior": "after"}
	if !reflect.DeepEqual(p, want) {
		t.Fatalf("%#v", p)
	}
}
func TestDependencyPartialResolutionAndCost(t *testing.T) {
	d := effects.Definition{ID: "deps", Kind: "linear_stages", Stages: []effects.Stage{{ID: "cost", Costs: []operations.Definition{{ID: "cost", Kind: operations.Impossible, Cost: true}}, Operations: []operations.Definition{{ID: "blocked", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("bad")}}}}, {ID: "partial", Operations: []operations.Definition{{ID: "miss", Kind: operations.Impossible, ResultBinding: "miss"}, {ID: "dependent", Kind: operations.EmitEvent, Dependency: &operations.Dependency{Binding: "miss", RequireApplied: true}, Arguments: map[string]query.Expr{"event_type": lit("bad2")}}, {ID: "continue", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("good")}}}}}}
	o, err := must(t, d).Resolve(base(), contracts.ResolutionInput{CommandID: "c", Procedure: model.ProcedureRef{ID: "p", Kind: "deps"}})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, x := range o.Events {
		seen[x.Type] = true
	}
	if seen["bad"] || seen["bad2"] || !seen["good"] {
		t.Fatal(seen)
	}
}
func TestOpaqueChoiceResumeReplayAndNoPublicLeak(t *testing.T) {
	d := effects.Definition{ID: "choice", Kind: "linear_stages", Stages: []effects.Stage{{ID: "choose", Choice: &effects.Choice{Kind: "select", Binding: "picked", Visibility: "opaque", Owner: lit("p1"), Domain: query.Expr{Kind: query.List, Args: []query.Expr{lit("secret-a"), lit("secret-b")}}, Prompt: map[string]any{"message": "choose"}}, Operations: []operations.Definition{{ID: "bind", Kind: operations.BindValue, ResultBinding: "selected", Arguments: map[string]query.Expr{"value": ref(query.Choice, "picked")}}}}}}
	e := must(t, d)
	first, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "c", Procedure: model.ProcedureRef{ID: "p", Kind: "choice"}})
	if err != nil {
		t.Fatal(err)
	}
	if first.Status != contracts.ResolutionPending || bytes.Contains(first.PendingInteraction.LegalDomain, []byte("secret")) {
		t.Fatalf("%s", first.PendingInteraction.LegalDomain)
	}
	var domain struct {
		Options []struct {
			Handle string `json:"handle"`
		} `json:"options"`
	}
	json.Unmarshal(first.PendingInteraction.LegalDomain, &domain)
	in := contracts.ResolutionInput{CommandID: "c", Procedure: first.PendingInteraction.ResumeProcedure, InteractionID: first.PendingInteraction.ID, Choice: raw(domain.Options[0].Handle), Context: map[string]json.RawMessage{"actor_player_id": raw("p1")}}
	a, _ := e.Resolve(base(), in)
	b, _ := e.Resolve(base(), in)
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	if !bytes.Equal(aj, bj) {
		t.Fatal("replay divergence")
	}
	for _, x := range a.Events {
		if bytes.Contains(x.PublicPayload, []byte("secret")) {
			t.Fatalf("leak: %s", x.PublicPayload)
		}
	}
}
func TestHistoryCancellationValidationProjection(t *testing.T) {
	events := []contracts.DomainEvent{{ID: "e2", Sequence: 2, SourceRef: "b", PublicPayload: raw(map[string]any{"cause_kind": "propagated", "participants": []string{"f"}})}, {ID: "e1", Sequence: 1, SourceRef: "a", PublicPayload: raw(map[string]any{"cause_kind": "combat", "participants": []string{"f"}})}}
	l := history.Build(events)
	if len(l.ByParticipant("f")) != 2 || l.ByParticipant("f")[0].EventID != "e1" {
		t.Fatal(l)
	}
	q := effects.Queue{{ID: "p", Scope: "x", Cancelable: false}, {ID: "c", Scope: "x", Cancelable: true}}
	if !reflect.DeepEqual(q.Cancel("x"), []string{"c"}) {
		t.Fatal(q)
	}
	bad := effects.Definition{ID: "bad", Kind: "linear_stages", Stages: []effects.Stage{{ID: "s", Operations: []operations.Definition{{ID: "x", Kind: operations.BindValue, Arguments: map[string]query.Expr{"value": ref(query.Result, "later")}}}}}}
	if _, err := rules.New([]effects.Definition{bad}); err == nil {
		t.Fatal("forward ref accepted")
	}
	e := must(t, effects.Definition{ID: "noop", Kind: "linear_stages", Stages: []effects.Stage{{ID: "s"}}})
	s := base()
	s.Players["p2"] = model.PlayerState{ID: "p2", PrivateZones: map[string][]model.CardID{"hand": {"secret"}}}
	s.Cards["secret"] = model.RuntimeObject{DefinitionID: "hidden"}
	p, _ := e.Project(s, "p1")
	if bytes.Contains(p, []byte("secret")) || bytes.Contains(p, []byte("hidden")) {
		t.Fatalf("leak %s", p)
	}
}
func must(t *testing.T, d effects.Definition) *rules.Engine {
	t.Helper()
	e, err := rules.New([]effects.Definition{d})
	if err != nil {
		t.Fatal(err)
	}
	return e
}
func base() model.GameState {
	return model.GameState{MatchID: "m", DefinitionRef: model.DefinitionRef{RulesetVersion: "v1", CapabilityRegistry: "wave1"}, Revision: 7, EventSequence: 10, Lifecycle: model.LifecycleActive, Players: map[model.PlayerID]model.PlayerState{"p1": {ID: "p1"}}, Fighters: map[model.FighterID]model.RuntimeObject{}, Cards: map[model.CardID]model.RuntimeObject{}, Components: map[model.ComponentID]model.RuntimeObject{}, Turn: map[string]any{}}
}
func lit(v any) query.Expr { return query.Expr{Kind: query.Literal, Value: v} }
func ref(s query.Source, p ...string) query.Expr {
	return query.Expr{Kind: query.Reference, Source: s, Path: p}
}
func ptr(e query.Expr) *query.Expr { return &e }
func raw(v any) json.RawMessage    { b, _ := json.Marshal(v); return b }
