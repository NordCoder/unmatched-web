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
	out, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "resume", Procedure: first.PendingInteraction.ResumeProcedure, InteractionID: first.PendingInteraction.ID, Choice: raw(firstHandle(t, first))})
	if err != nil || out.Status != contracts.ResolutionCompleted {
		t.Fatalf("resume: %+v %v", out, err)
	}
}

func TestCostBeforeChoicePersistsOnce(t *testing.T) {
	d := costChoiceDefinition()
	e := must(t, d)
	s := base()
	s.Fighters["alice"] = model.RuntimeObject{DefinitionID: "alice", State: map[string]any{}}
	first, err := e.Resolve(s, contracts.ResolutionInput{CommandID: "start", Procedure: model.ProcedureRef{ID: "p", Kind: d.ID, SourceRef: "card:cost"}})
	if err != nil || countType(first.Events, "rules.fighter_state_set") != 1 || countType(first.Events, "rules.operation_result") != 1 {
		t.Fatalf("first: %+v %v", first, err)
	}
	persisted := reduceAll(t, e, s, first.Events)
	if persisted.Fighters["alice"].State["paid"] != true {
		t.Fatal("cost not reduced")
	}
	resume, err := e.Resolve(persisted, contracts.ResolutionInput{CommandID: "resume", Procedure: first.PendingInteraction.ResumeProcedure, InteractionID: first.PendingInteraction.ID, Choice: raw(firstHandle(t, first))})
	if err != nil || countType(resume.Events, "rules.fighter_state_set") != 0 || countType(resume.Events, "test.after") != 1 {
		t.Fatalf("resume: %+v %v", resume, err)
	}
}

func TestFailedCostRollsBackAndReportsDisposition(t *testing.T) {
	d := effects.Definition{ID: "failed-cost", Kind: "linear_stages", RuleID: "FX-032", CapabilityIDs: []string{"CAP-001"}, Stages: []effects.Stage{{ID: "pay", Costs: []operations.Definition{{ID: "set", Kind: operations.SetFighterState, Cost: true, ResultBinding: "setr", Arguments: map[string]query.Expr{"fighter": fighter("alice"), "key": lit("paid"), "value": lit(true)}}, {ID: "fail", Kind: operations.Impossible, Cost: true, ResultBinding: "failr"}}, Operations: []operations.Definition{{ID: "never", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("never"), "payload": obj("x", lit(true))}}}}}}
	e := must(t, d)
	s := base()
	s.Fighters["alice"] = model.RuntimeObject{DefinitionID: "alice", State: map[string]any{}}
	out, err := e.Resolve(s, contracts.ResolutionInput{CommandID: "c", Procedure: model.ProcedureRef{ID: "p", Kind: d.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if countType(out.Events, "rules.fighter_state_set") != 0 {
		t.Fatal("rolled-back mutation event escaped")
	}
	disps := resultDispositions(t, out.Events)
	if !contains(disps, string(operations.DispositionRolledBackCost)) || !contains(disps, string(operations.DispositionNotApplied)) {
		t.Fatalf("dispositions %v", disps)
	}
	after := reduceAll(t, e, s, out.Events)
	if _, ok := after.Fighters["alice"].State["paid"]; ok {
		t.Fatal("failed cost mutated state")
	}
}

func TestQueuedEffectsUseRealDisposition(t *testing.T) {
	d := effects.Definition{
		ID: "queue", Kind: "linear_stages", RuleID: "FX-010",
		CapabilityIDs: []string{"CAP-001", "CAP-002", "CAP-003"},
		Stages: []effects.Stage{{
			ID: "q", Checkpoint: "during",
			Cancellations: []effects.Cancellation{{Scope: "cancel"}},
			Queue: []effects.QueueDefinition{
				{ID: "none", Scope: "x", Priority: 3, Operations: []operations.Definition{{ID: "impossible", Kind: operations.Impossible}}},
				{ID: "mixed", Scope: "x", Priority: 2, Operations: []operations.Definition{
					{ID: "miss", Kind: operations.Impossible},
					{ID: "hit", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("queue.hit"), "payload": obj("ok", lit(true))}},
				}},
				{ID: "canceled", Scope: "cancel", Priority: 1, Cancelable: true, Operations: []operations.Definition{
					{ID: "bad", Kind: operations.EmitEvent, Arguments: map[string]query.Expr{"event_type": lit("queue.canceled.bad"), "payload": obj("ok", lit(true))}},
				}},
			},
		}},
	}
	e := must(t, d)
	out, err := e.Resolve(base(), contracts.ResolutionInput{CommandID: "q", Procedure: model.ProcedureRef{ID: "p", Kind: d.ID}})
	if err != nil {
		t.Fatal(err)
	}
	outcomes := queueOutcomes(t, out.Events)
	if outcomes["none"] != string(operations.DispositionNotApplied) || outcomes["mixed"] != string(operations.DispositionPartial) || outcomes["canceled"] != string(operations.DispositionCanceled) {
		t.Fatalf("outcomes %v", outcomes)
	}
	if countType(out.Events, "queue.canceled.bad") != 0 {
		t.Fatal("canceled effect executed")
	}
}

func choiceDefinition() effects.Definition {
	return effects.Definition{
		ID: "choice", Kind: "linear_stages", RuleID: "FX-060",
		CapabilityIDs: []string{"CAP-001", "CAP-002"},
		Stages: []effects.Stage{{
			ID: "choose",
			Choice: &effects.Choice{
				Kind: "select", Binding: "picked", Visibility: query.Opaque,
				Owner: player("p1"), Domain: list(lit("secret-a"), lit("secret-b")),
				Prompt: map[string]any{"message": "choose"}, EmptyDomain: effects.EmptyReject,
				ValueType: query.TypeString,
			},
			Operations: []operations.Definition{{
				ID: "bind", Kind: operations.BindValue, ResultBinding: "selected",
				Arguments: map[string]query.Expr{"value": ref(query.Choice, query.TypeString, "picked")},
			}},
		}},
	}
}

func costChoiceDefinition() effects.Definition {
	return effects.Definition{
		ID: "cost-choice", Kind: "linear_stages", RuleID: "FX-032",
		CapabilityIDs: []string{"CAP-001", "CAP-002"},
		Stages: []effects.Stage{{
			ID: "pay",
			Costs: []operations.Definition{{
				ID: "pay", Kind: operations.SetFighterState, Cost: true, CauseKind: "cost", ResultBinding: "paid",
				Arguments: map[string]query.Expr{"fighter": fighter("alice"), "key": lit("paid"), "value": lit(true)},
			}},
			Choice: &effects.Choice{
				Kind: "select", Binding: "picked", Visibility: query.Opaque,
				Owner: player("p1"), Domain: list(lit("left"), lit("right")),
				Prompt: map[string]any{"message": "choose"}, EmptyDomain: effects.EmptyReject,
				ValueType: query.TypeString,
			},
			Operations: []operations.Definition{{
				ID: "after", Kind: operations.EmitEvent, CauseKind: "effect",
				Arguments: map[string]query.Expr{"event_type": lit("test.after"), "payload": obj("ok", lit(true))},
			}},
		}},
	}
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
	return model.GameState{
		MatchID:       "m",
		DefinitionRef: model.DefinitionRef{RulesetVersion: "v1", CapabilityRegistry: "wave1"},
		Revision:      7, EventSequence: 10, Lifecycle: model.LifecycleActive,
		Players:  map[model.PlayerID]model.PlayerState{"p1": {ID: "p1"}},
		Fighters: map[model.FighterID]model.RuntimeObject{}, Cards: map[model.CardID]model.RuntimeObject{},
		Components: map[model.ComponentID]model.RuntimeObject{}, Turn: map[string]any{},
	}
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
func countType(es []contracts.DomainEvent, typ string) int {
	n := 0
	for _, e := range es {
		if e.Type == typ {
			n++
		}
	}
	return n
}
func reduceAll(t *testing.T, e *rules.Engine, s model.GameState, events []contracts.DomainEvent) model.GameState {
	t.Helper()
	var err error
	for _, ev := range events {
		s, err = e.ApplyEvent(s, ev)
		if err != nil {
			t.Fatal(err)
		}
	}
	return s
}
func resultDispositions(t *testing.T, events []contracts.DomainEvent) []string {
	t.Helper()
	var out []string
	for _, ev := range events {
		if ev.Type != "rules.operation_result" {
			continue
		}
		var env operations.EventEnvelope
		_ = json.Unmarshal(ev.PublicPayload, &env)
		m, _ := env.Data.(map[string]any)
		if d, ok := m["disposition"].(string); ok {
			out = append(out, d)
		}
	}
	return out
}
func queueOutcomes(t *testing.T, events []contracts.DomainEvent) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, ev := range events {
		if ev.Type != "rules.queued_effect_outcome" {
			continue
		}
		var env operations.EventEnvelope
		_ = json.Unmarshal(ev.PublicPayload, &env)
		m, _ := env.Data.(map[string]any)
		effect, _ := m["effect"].(string)
		disp, _ := m["disposition"].(string)
		out[effect] = disp
	}
	return out
}
func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
