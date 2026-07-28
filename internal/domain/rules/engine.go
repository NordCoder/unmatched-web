// Package rules implements the concrete pure staged RulesEngine used by Core.
package rules

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/domain/capabilities"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
)

type Engine struct {
	defs map[string]effects.Definition
	ops  *operations.Registry
	caps map[string]capabilities.Capability
}

func New(defs []effects.Definition) (*Engine, error) {
	ops := operations.Default()
	caps := capabilities.Wave1()
	if err := capabilities.Validate(caps); err != nil {
		return nil, err
	}
	e := &Engine{defs: map[string]effects.Definition{}, ops: ops, caps: caps}
	for _, d := range defs {
		if err := effects.Validate(d, ops, caps); err != nil {
			return nil, fmt.Errorf("definition %s: %w", d.ID, err)
		}
		if _, ok := e.defs[d.ID]; ok {
			return nil, fmt.Errorf("duplicate definition %s", d.ID)
		}
		e.defs[d.ID] = d
	}
	return e, nil
}

func (e *Engine) Resolve(state model.GameState, in contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	d, ok := e.defs[in.Procedure.Kind]
	if !ok {
		return reject("unknown_procedure"), nil
	}
	ps, err := effects.Decode(in.Procedure)
	if err != nil {
		return reject("invalid_procedure_state"), nil
	}
	working, err := clone(state)
	if err != nil {
		return contracts.ResolutionOutcome{}, err
	}
	var events []contracts.DomainEvent
	if ps.Pending != nil {
		if in.InteractionID != ps.Pending.ID {
			return reject("interaction_mismatch"), nil
		}
		var handle string
		if json.Unmarshal(in.Choice, &handle) != nil {
			return reject("invalid_choice"), nil
		}
		value, ok := ps.Pending.Options[handle]
		if !ok {
			return reject("choice_outside_domain"), nil
		}
		ps.Choices[ps.Pending.Binding] = append(json.RawMessage(nil), value...)
		ev, _ := operations.MetaEvent("rules.choice_accepted", in.Procedure.SourceRef, "choice:"+ps.Pending.Binding, "choice", map[string]any{"interaction_instance_id": ps.Pending.ID, "binding": ps.Pending.Binding}, nil)
		events = append(events, ev)
		ps.Pending = nil
		ps.Status = "running"
		ps.Phase = effects.PhaseOperations
	}
	for ps.Cursor < len(d.Stages) {
		s := d.Stages[ps.Cursor]
		ctx := query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}
		switch ps.Phase {
		case effects.PhaseEnter:
			if s.Condition != nil {
				v, x := query.Eval(*s.Condition, ctx)
				if x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				if v != true {
					events = append(events, e.meta(in, "rules.stage_skipped", "stage:"+s.ID, map[string]any{"stage": s.ID, "reason": "condition_false"}))
					advance(&ps)
					continue
				}
			}
			valid := true
			for _, p := range s.Prerequisites {
				v, x := query.Eval(p, ctx)
				if x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				if v != true {
					valid = false
					break
				}
			}
			if !valid {
				events = append(events, e.meta(in, "rules.stage_skipped", "stage:"+s.ID, map[string]any{"stage": s.ID, "reason": "prerequisite_failed"}))
				advance(&ps)
				continue
			}
			ps.Phase = effects.PhaseCosts
		case effects.PhaseCosts:
			candidate, x := clone(working)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			tempResults := copyRaw(ps.Results)
			var tempEvents []contracts.DomainEvent
			paid := true
			for i, o := range s.Costs {
				r, x := e.exec(candidate, query.Context{State: candidate, Captured: ps.Captured, Results: tempResults, Choices: ps.Choices, Input: in.Context}, in, o)
				if x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				if !r.Record.Applied {
					paid = false
					break
				}
				if x = operations.Apply(&candidate, r.Patches); x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				tempEvents = append(tempEvents, r.Events...)
				bindRaw(tempResults, o, r.Record)
				ps.CostIndex = i + 1
			}
			if !paid {
				events = append(events, e.meta(in, "rules.stage_skipped", "stage:"+s.ID, map[string]any{"stage": s.ID, "reason": "cost_unpaid"}))
				advance(&ps)
				continue
			}
			working = candidate
			ps.Results = tempResults
			events = append(events, tempEvents...)
			ps.Phase = effects.PhaseChoice
		case effects.PhaseChoice:
			if s.Choice == nil {
				ps.Phase = effects.PhaseOperations
				continue
			}
			if _, ok := ps.Choices[s.Choice.Binding]; ok {
				ps.Phase = effects.PhaseOperations
				continue
			}
			ownerV, x := query.Eval(s.Choice.Owner, ctx)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			owner, ok := ownerV.(string)
			if !ok {
				return contracts.ResolutionOutcome{}, fmt.Errorf("choice owner must be string")
			}
			domainV, x := query.Eval(s.Choice.Domain, ctx)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			items, ok := domainV.([]any)
			if !ok {
				return contracts.ResolutionOutcome{}, fmt.Errorf("choice domain must be list")
			}
			if len(items) == 0 {
				switch s.Choice.EmptyDomain {
				case effects.EmptySkipStage:
					events = append(events, e.meta(in, "rules.stage_skipped", "choice:"+s.Choice.Binding, map[string]any{"stage": s.ID, "reason": "empty_domain"}))
					advance(&ps)
					continue
				case effects.EmptyBindEmpty:
					ps.Choices[s.Choice.Binding] = json.RawMessage(`[]`)
					ps.Phase = effects.PhaseOperations
					continue
				case effects.EmptyComplete:
					ps.Phase = effects.PhaseOperations
					continue
				case effects.EmptyReject:
					return reject("empty_choice_domain"), nil
				}
			}
			opts := map[string]json.RawMessage{}
			projected := make([]any, 0, len(items))
			for i, item := range items {
				raw, _ := json.Marshal(item)
				h := effects.Handle(in.Procedure.ID, s.ID, i)
				opts[h] = raw
				if s.Choice.Visibility == query.Opaque {
					projected = append(projected, map[string]any{"handle": h})
				} else {
					projected = append(projected, map[string]any{"handle": h, "value": item})
				}
			}
			domain, _ := json.Marshal(map[string]any{"options": projected})
			prompt, _ := json.Marshal(s.Choice.Prompt)
			pid := model.InteractionID(id(string(in.Procedure.ID), s.ID, "interaction"))
			ps.Pending = &effects.Pending{ID: pid, Owner: model.PlayerID(owner), Kind: s.Choice.Kind, Visibility: s.Choice.Visibility, Binding: s.Choice.Binding, Options: opts, Domain: domain, Prompt: prompt}
			ps.Status = "pending"
			ref, _ := effects.Encode(in.Procedure, ps, s.ID)
			pending := &model.PendingInteraction{ID: pid, OwnerPlayerID: model.PlayerID(owner), Kind: s.Choice.Kind, Visibility: string(s.Choice.Visibility), Prompt: prompt, LegalDomain: domain, ResumeProcedure: ref}
			normalize(events, state, in)
			return contracts.ResolutionOutcome{Status: contracts.ResolutionPending, Events: events, PendingInteraction: pending, Procedure: &ref}, nil
		case effects.PhaseOperations:
			for ps.OperationIndex < len(s.Operations) {
				o := s.Operations[ps.OperationIndex]
				if o.Dependency != nil {
					var rec operations.Record
					if json.Unmarshal(ps.Results[o.Dependency.Binding], &rec) != nil {
						return contracts.ResolutionOutcome{}, fmt.Errorf("invalid dependency result")
					}
					if o.Dependency.RequireApplied && !rec.Applied {
						events = append(events, e.meta(in, "rules.operation_skipped", o.ID, map[string]any{"operation": o.ID, "reason": "dependency_failed"}))
						ps.OperationIndex++
						continue
					}
				}
				r, x := e.exec(working, query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}, in, o)
				if x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				if x = operations.Apply(&working, r.Patches); x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				events = append(events, r.Events...)
				bind(&ps, o, r.Record)
				events = append(events, e.meta(in, "rules.operation_result", o.ID, map[string]any{"operation": o.ID, "binding": o.ResultBinding, "applied": r.Record.Applied, "code": r.Record.Code}))
				ps.OperationIndex++
			}
			ps.Phase = effects.PhaseCheckpoint
		case effects.PhaseCheckpoint:
			if s.Checkpoint != "" {
				if ps.Checkpoint == nil {
					ps.Checkpoint = effects.BuildCheckpoint(s.Checkpoint, ps.Cursor+1, s.Queue)
					events = append(events, e.meta(in, "rules.checkpoint_opened", "checkpoint:"+s.Checkpoint, map[string]any{"checkpoint": s.Checkpoint}))
				}
				for _, c := range s.Cancellations {
					for _, qid := range ps.Checkpoint.Cancel(c.Scope) {
						events = append(events, e.meta(in, "rules.effect_canceled", qid, map[string]any{"checkpoint": s.Checkpoint, "effect": qid, "scope": c.Scope}))
					}
				}
				for {
					q := ps.Checkpoint.Next()
					if q == nil {
						break
					}
					for _, operation := range q.Operations {
						r, x := e.exec(working, query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}, in, operation)
						if x != nil {
							return contracts.ResolutionOutcome{}, x
						}
						if x = operations.Apply(&working, r.Patches); x != nil {
							return contracts.ResolutionOutcome{}, x
						}
						events = append(events, r.Events...)
					}
					events = append(events, e.meta(in, "rules.queued_effect_resolved", q.ID, map[string]any{"checkpoint": s.Checkpoint, "effect": q.ID, "priority": q.Priority, "source_order": q.SourceOrder}))
				}
				ps.Checkpoint = nil
			}
			advance(&ps)
		default:
			return contracts.ResolutionOutcome{}, fmt.Errorf("invalid procedure phase %q", ps.Phase)
		}
	}
	ps.Status = "completed"
	ps.Phase = effects.PhaseComplete
	ref, _ := effects.Encode(in.Procedure, ps, "completed")
	normalize(events, state, in)
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted, Events: events, Procedure: &ref}, nil
}
func (e *Engine) exec(s model.GameState, q query.Context, in contracts.ResolutionInput, o operations.Definition) (operations.Result, error) {
	return e.ops.Execute(operations.Context{State: s, Query: q, SourceRef: in.Procedure.SourceRef}, o)
}
func (e *Engine) ApplyEvent(s model.GameState, event contracts.DomainEvent) (model.GameState, error) {
	return operations.ApplyEvent(s, event)
}
func (e *Engine) LegalActions(s model.GameState, p model.PlayerID) ([]json.RawMessage, error) {
	if _, ok := s.Players[p]; !ok {
		return nil, fmt.Errorf("player not found")
	}
	raw := s.Turn["legal_actions_by_player"]
	b, _ := json.Marshal(raw)
	var m map[string][]json.RawMessage
	_ = json.Unmarshal(b, &m)
	return m[string(p)], nil
}

type projectedAction struct {
	ID      model.ActionID `json:"action_instance_id"`
	Kind    string         `json:"kind"`
	ActorID model.PlayerID `json:"actor_player_id"`
	Status  string         `json:"status"`
}
type projectedCombat struct {
	ID         model.CombatID  `json:"combat_instance_id"`
	AttackerID model.FighterID `json:"attacker_id"`
	DefenderID model.FighterID `json:"defender_id"`
	Stage      string          `json:"stage"`
}

func (e *Engine) Project(s model.GameState, p model.PlayerID) (json.RawMessage, error) {
	player, ok := s.Players[p]
	if !ok {
		return nil, fmt.Errorf("player not found")
	}
	fighters := make([]model.FighterID, 0, len(s.Fighters))
	for fid := range s.Fighters {
		fighters = append(fighters, fid)
	}
	sort.Slice(fighters, func(i, j int) bool { return fighters[i] < fighters[j] })
	v := map[string]any{"match_id": s.MatchID, "revision": s.Revision, "event_sequence": s.EventSequence, "lifecycle": s.Lifecycle, "definition_ref": s.DefinitionRef, "self": player, "fighter_ids": fighters}
	if s.Action != nil {
		v["action"] = projectedAction{s.Action.ID, s.Action.Kind, s.Action.ActorID, s.Action.Status}
	}
	if s.Combat != nil {
		v["combat"] = projectedCombat{s.Combat.ID, s.Combat.AttackerID, s.Combat.DefenderID, s.Combat.Stage}
	}
	return json.Marshal(v)
}
func (e *Engine) meta(in contracts.ResolutionInput, t, op string, data any) contracts.DomainEvent {
	ev, _ := operations.MetaEvent(t, in.Procedure.SourceRef, op, "resolver", data, nil)
	return ev
}
func bind(s *effects.State, o operations.Definition, r operations.Record) { bindRaw(s.Results, o, r) }
func bindRaw(dst map[string]json.RawMessage, o operations.Definition, r operations.Record) {
	if o.ResultBinding != "" {
		raw, _ := json.Marshal(r)
		dst[o.ResultBinding] = raw
	}
}
func advance(s *effects.State) {
	s.Cursor++
	s.Phase = effects.PhaseEnter
	s.CostIndex = 0
	s.OperationIndex = 0
	s.Pending = nil
	s.Checkpoint = nil
}
func normalize(es []contracts.DomainEvent, s model.GameState, in contracts.ResolutionInput) {
	for i := range es {
		es[i].SchemaVersion = "unmatched.domain-event/v1"
		es[i].ID = model.EventID(id(string(in.CommandID), string(in.Procedure.ID), fmt.Sprint(i), es[i].Type))
		es[i].MatchID = s.MatchID
		es[i].Sequence = s.EventSequence + uint64(i) + 1
		es[i].Revision = s.Revision + 1
		es[i].CausedByCommand = in.CommandID
		es[i].RulesetVersion = s.DefinitionRef.RulesetVersion
		if es[i].PublicPayload == nil {
			es[i].PublicPayload = []byte(`{}`)
		}
		if es[i].PrivatePayloads == nil {
			es[i].PrivatePayloads = map[model.PlayerID]json.RawMessage{}
		}
	}
}
func reject(c string) contracts.ResolutionOutcome {
	return contracts.ResolutionOutcome{Status: contracts.ResolutionRejected, RejectionCode: c}
}
func clone(s model.GameState) (model.GameState, error) {
	b, e := json.Marshal(s)
	if e != nil {
		return s, e
	}
	var x model.GameState
	e = json.Unmarshal(b, &x)
	return x, e
}
func copyRaw(src map[string]json.RawMessage) map[string]json.RawMessage {
	r := make(map[string]json.RawMessage, len(src))
	for k, v := range src {
		r[k] = append(json.RawMessage(nil), v...)
	}
	return r
}
func id(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:16])
}
