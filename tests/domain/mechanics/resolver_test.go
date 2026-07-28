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

func TestOpaqueResumeAndProjectionAllowList(t *testing.T) {
	e := must(t, choiceDefinition())
	s := base()
	secret := raw(map[string]any{"secret": "secret-a"})
	proc := model.ProcedureRef{ID: "p", Kind: "choice", Bindings: map[string]json.RawMessage{effects.StateBinding: secret}}
	s.Action = &model.ActionState{ID: "a", Kind: "SCHEME", ActorID: "p1", Status: "WAITING", Procedure: &proc}
	s.Combat = &model.CombatState{ID: "c", AttackerID: "alice", DefenderID: "bob", Stage: "DEFENSE", Procedure: &proc}
	projection, err := e.Project(s, "p1")
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range [][]byte{[]byte("secret-a"), []byte(effects.StateBinding), []byte("bindings")} {
		if bytes.Contains(projection, forbidden) {
			t.Fatalf("projection leak %q: %s", forbidden, projection)
		}
	}
	first, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "start", Procedure: model.ProcedureRef{ID: "p", Kind: "choice", Bindings: map[string]json.RawMessage{"actor_player_id": raw("p1")}}})
	if err != nil || first.Status != contracts.ResolutionPending {
		t.Fatalf("pending: %+v %v", first, err)
	}
	if bytes.Contains(first.PendingInteraction.LegalDomain, []byte("secret-")) {
		t.Fatalf("domain leak: %s", first.PendingInteraction.LegalDomain)
	}
	h := firstHandle(t, first)
	out, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "resume", Procedure: first.PendingInteraction.ResumeProcedure, InteractionID: first.PendingInteraction.ID, Choice: raw(h)})
	if err != nil || out.Status != contracts.ResolutionCompleted {
		t.Fatalf("resume: %+v %v", out, err)
	}
	for _, event := range out.Events {
		if bytes.Contains(event.PublicPayload, []byte("secret-")) {
			t.Fatalf("event leak: %s", event.PublicPayload)
		}
	}
}

func TestCostBeforeChoicePersistsOnce(t *testing.T) {
	d := effects.Definition{ID: "cost-choice", Kind: "linear_stages", RuleID: "FX-032", CapabilityIDs: []string{"CAP-001", "CAP-002"}, Stages: []effects.Stage{{
		ID: "pay", Costs: []operations.Definition{{ID: "pay", Kind: operations.SetFighterState, Cost: true, CauseKind: "cost", ResultBinding: "paid", Arguments: map[string]query.Expr{"fighter": fighter("alice"), "key": lit("paid"), "value": lit(true)}}},
		Choice:     &effects.Choice{Kind: "select", Binding: "picked", Visibility: query.Opaque, Owner: player("p1"), Domain: list(lit("left"), lit("right")), Prompt: map[string]any{"message": "choose"}, EmptyDomain: effects.EmptyReject, ValueType: query.TypeString},
		Operations: []operations.Definition{{ID: "after", Kind: operations.EmitEvent, CauseKind: "effect", Arguments: map[string]query.Expr{"event_type": lit("test.after"), "payload": obj("ok", lit(true))}}},
	}}}
	e := must(t, d)
	s := base()
	s.Fighters["alice"] = model.RuntimeObject{DefinitionID: "alice", State: map[string]any{}}
	first, err := e.Resolve(s, contracts.ResolutionInput{CommandID: "start", Procedure: model.ProcedureRef{ID: "p", Kind: d.ID, SourceRef: "card:cost"}})
	if err != nil || countType(first.Events, "rules.fighter_state_set") != 1 {
		t.Fatalf("first: %+v %v", first, err)
	}
	persisted := s
	for _, ev := range first.Events {
		persisted, err = e.ApplyEvent(persisted, ev)
		if err != nil {
			t.Fatal(err)
		}
	}
	if persisted.Fighters["alice"].State["paid"] != true {
		t.Fatal("cost not reduced")
	}
	resume, err := e.Resolve(persisted, contracts.ResolutionInput{CommandID: "resume", Procedure: first.PendingInteraction.ResumeProcedure, InteractionID: first.PendingInteraction.ID, Choice: raw(firstHandle(t, first))})
	if err != nil || countType(resume.Events, "rules.fighter_state_set") != 0 || countType(resume.Events, "test.after") != 1 {
		t.Fatalf("resume: %+v %v", resume, err)
	}
}

func choiceDefinition() effects.Definition {
	return effects.Definition{ID: "choice", Kind: "linear_stages", RuleID: "FX-060", CapabilityIDs: []string{"CAP-001", "CAP-002"}, Stages: []effects.Stage{{ID: "choose", Choice: &effects.Choice{Kind: "select", Binding: "picked", Visibility: query.Opaque, Owner: player("p1"), Domain: list(lit("secret-a"), lit("secret-b")), Prompt: map[string]any{"message": "choose"}, EmptyDomain: effects.EmptyReject, ValueType: query.TypeString}, Operations: []operations.Definition{{ID: "bind", Kind: operations.BindValue, ResultBinding: "selected", Arguments: map[string]query.Expr{"value": ref(query.Choice, query.TypeString, "picked")}}}}}}
}
func firstHandle(t *testing.T, out contracts.ResolutionOutcome) string {
	t.Helper()
	var d struct {
		Options []struct {
			Handle string `json:"handle"`
		} `json:"options"`
	}
	if err := json.Unmarshal(out.PendingInteraction.LegalDomain, &d); err != nil || len(d.Options) == 0 {
		t.Fatalf("domain: %v %s", err, out.PendingInteraction.LegalDomain)
	}
	return d.Options[0].Handle
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
func player(v string) query.Expr {
	return query.Expr{Kind: query.Literal, Value: v, ValueType: query.TypePlayer}
}
func fighter(v string) query.Expr {
	return query.Expr{Kind: query.Literal, Value: v, ValueType: query.TypeFighter}
}
func ref(s query.Source, t query.Type, p ...string) query.Expr {
	return query.Expr{Kind: query.Reference, Source: s, Path: p, ValueType: t}
}
func list(v ...query.Expr) query.Expr { return query.Expr{Kind: query.List, Args: v} }
func obj(k string, v query.Expr) query.Expr {
	return query.Expr{Kind: query.Object, Fields: map[string]query.Expr{k: v}}
}
func raw(v any) json.RawMessage { b, _ := json.Marshal(v); return b }
func countType(es []contracts.DomainEvent, t string) int {
	n := 0
	for _, e := range es {
		if e.Type == t {
			n++
		}
	}
	return n
}
