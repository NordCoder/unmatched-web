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
	defs        map[string]effects.Definition
	ops         *operations.Registry
	caps        map[string]capabilities.Capability
	choiceSpecs map[string]map[string]query.ValueSpec
}
type attempt struct {
	definition  operations.Definition
	result      operations.Result
	resultEvent contracts.DomainEvent
}

func New(defs []effects.Definition) (*Engine, error) {
	ops := operations.Default()
	caps := capabilities.Wave1()
	if err := capabilities.Validate(caps); err != nil {
		return nil, err
	}
	e := &Engine{defs: map[string]effects.Definition{}, ops: ops, caps: caps, choiceSpecs: map[string]map[string]query.ValueSpec{}}
	for _, d := range defs {
		if err := effects.Validate(d, ops, caps); err != nil {
			return nil, fmt.Errorf("definition %s: %w", d.ID, err)
		}
		if _, ok := e.defs[d.ID]; ok {
			return nil, fmt.Errorf("duplicate definition %s", d.ID)
		}
		e.defs[d.ID] = d
		e.choiceSpecs[d.ID] = map[string]query.ValueSpec{}
		for _, s := range d.Stages {
			if s.Choice != nil {
				sp, err := effects.ChoiceSpec(d, s.Choice.Binding, ops)
				if err != nil {
					return nil, err
				}
				e.choiceSpecs[d.ID][s.Choice.Binding] = sp
			}
		}
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
	working, err := operations.CloneState(state)
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
		var decoded any
		if err = json.Unmarshal(value, &decoded); err != nil {
			return reject("invalid_choice_value"), nil
		}
		if err = query.ValidateValue(decoded, ps.Pending.ValueSpec); err != nil {
			return reject("invalid_choice_value"), nil
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
			candidate, x := operations.CloneState(working)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			tempResults := copyRaw(ps.Results)
			attempts := make([]attempt, 0, len(s.Costs))
			paid := true
			for i, o := range s.Costs {
				a, x := e.attempt(candidate, query.Context{State: candidate, Captured: ps.Captured, Results: tempResults, Choices: ps.Choices, Input: in.Context}, in, o, "cost")
				if x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				attempts = append(attempts, a)
				bindRaw(tempResults, o, a.result.Record)
				ps.CostIndex = i + 1
				if a.result.Record.Disposition != operations.DispositionApplied {
					paid = false
					break
				}
				if x = operations.Apply(&candidate, a.result.Patches); x != nil {
					return contracts.ResolutionOutcome{}, x
				}
			}
			if !paid {
				for i := range attempts {
					if attempts[i].result.Record.Disposition == operations.DispositionApplied {
						attempts[i].result.Record.Disposition = operations.DispositionRolledBackCost
						attempts[i].result.Record.Applied = false
						attempts[i].resultEvent = e.resultMeta(in, attempts[i].definition, attempts[i].result.Record, "cost")
					}
					bindRaw(ps.Results, attempts[i].definition, attempts[i].result.Record)
					events = append(events, attempts[i].resultEvent)
				}
				events = append(events, e.meta(in, "rules.stage_skipped", "stage:"+s.ID, map[string]any{"stage": s.ID, "reason": "cost_unpaid"}))
				advance(&ps)
				continue
			}
			working = candidate
			ps.Results = tempResults
			for _, a := range attempts {
				events = append(events, a.result.Events...)
				events = append(events, a.resultEvent)
			}
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
			spec := e.choiceSpecs[d.ID][s.Choice.Binding]
			if len(items) == 0 {
				switch s.Choice.EmptyDomain {
				case effects.EmptySkipStage:
					events = append(events, e.meta(in, "rules.stage_skipped", "choice:"+s.Choice.Binding, map[string]any{"stage": s.ID, "reason": "empty_domain"}))
					advance(&ps)
					continue
				case effects.EmptyBindDefault:
					if s.Choice.Default == nil {
						return contracts.ResolutionOutcome{}, fmt.Errorf("missing validated default")
					}
					v, x := query.Eval(*s.Choice.Default, ctx)
					if x != nil {
						return contracts.ResolutionOutcome{}, x
					}
					if x = query.ValidateValue(v, spec); x != nil {
						return contracts.ResolutionOutcome{}, x
					}
					ps.Choices[s.Choice.Binding] = mustRaw(v)
					ps.Phase = effects.PhaseOperations
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
				if x = query.ValidateValue(item, spec); x != nil {
					return contracts.ResolutionOutcome{}, fmt.Errorf("choice domain item %d: %w", i, x)
				}
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
			ps.Pending = &effects.Pending{ID: pid, Owner: model.PlayerID(owner), Kind: s.Choice.Kind, Visibility: s.Choice.Visibility, Binding: s.Choice.Binding, ValueSpec: spec, Options: opts, Domain: domain, Prompt: prompt}
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
					if o.Dependency.RequireApplied && rec.Disposition != operations.DispositionApplied {
						skip := operations.Record{Disposition: operations.DispositionSkippedDependency, Code: "dependency_failed"}
						bind(&ps, o, skip)
						events = append(events, e.meta(in, "rules.operation_skipped", o.ID, map[string]any{"operation": o.ID, "reason": "dependency_failed"}))
						events = append(events, e.resultMeta(in, o, skip, "ordinary"))
						ps.OperationIndex++
						continue
					}
				}
				a, x := e.attempt(working, query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}, in, o, "ordinary")
				if x != nil {
					return contracts.ResolutionOutcome{}, x
				}
				if a.result.Record.Disposition == operations.DispositionApplied {
					if x = operations.Apply(&working, a.result.Patches); x != nil {
						return contracts.ResolutionOutcome{}, x
					}
					events = append(events, a.result.Events...)
				}
				bind(&ps, o, a.result.Record)
				events = append(events, a.resultEvent)
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
						events = append(events, e.meta(in, "rules.queued_effect_outcome", qid, map[string]any{"checkpoint": s.Checkpoint, "effect": qid, "disposition": operations.DispositionCanceled}))
					}
				}
				for {
					q := ps.Checkpoint.Next()
					if q == nil {
						break
					}
					applied, total := 0, len(q.Operations)
					for _, o := range q.Operations {
						a, x := e.attempt(working, query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}, in, o, "queue")
						if x != nil {
							return contracts.ResolutionOutcome{}, x
						}
						if a.result.Record.Disposition == operations.DispositionApplied {
							applied++
							if x = operations.Apply(&working, a.result.Patches); x != nil {
								return contracts.ResolutionOutcome{}, x
							}
							events = append(events, a.result.Events...)
						}
						events = append(events, a.resultEvent)
					}
					disp := operations.DispositionNotApplied
					if applied == total {
						disp = operations.DispositionApplied
					} else if applied > 0 {
						disp = operations.DispositionPartial
					}
					ps.Checkpoint.Finish(q.ID, disp)
					events = append(events, e.meta(in, "rules.queued_effect_outcome", q.ID, map[string]any{"checkpoint": s.Checkpoint, "effect": q.ID, "priority": q.Priority, "source_order": q.SourceOrder, "disposition": disp}))
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

func (e *Engine) attempt(s model.GameState, q query.Context, in contracts.ResolutionInput, o operations.Definition, mode string) (attempt, error) {
	res, err := e.ops.Execute(operations.Context{State: s, Query: q, SourceRef: in.Procedure.SourceRef}, o)
	if err != nil {
		return attempt{}, err
	}
	return attempt{definition: o, result: res, resultEvent: e.resultMeta(in, o, res.Record, mode)}, nil
}
func (e *Engine) resultMeta(in contracts.ResolutionInput, o operations.Definition, r operations.Record, mode string) contracts.DomainEvent {
	return e.meta(in, "rules.operation_result", o.ID, map[string]any{"operation": o.ID, "binding": o.ResultBinding, "applied": r.Disposition == operations.DispositionApplied, "disposition": r.Disposition, "code": r.Code, "mode": mode})
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
type projection struct {
	MatchID       model.MatchID                           `json:"match_id"`
	Revision      uint64                                  `json:"revision"`
	EventSequence uint64                                  `json:"event_sequence"`
	Lifecycle     model.Lifecycle                         `json:"lifecycle"`
	Player        model.PlayerState                       `json:"player"`
	Fighters      map[model.FighterID]model.RuntimeObject `json:"fighters"`
	Action        *projectedAction                        `json:"action,omitempty"`
	Combat        *projectedCombat                        `json:"combat,omitempty"`
}

func (e *Engine) Project(s model.GameState, p model.PlayerID) (json.RawMessage, error) {
	player, ok := s.Players[p]
	if !ok {
		return nil, fmt.Errorf("player not found")
	}
	out := projection{MatchID: s.MatchID, Revision: s.Revision, EventSequence: s.EventSequence, Lifecycle: s.Lifecycle, Player: player, Fighters: s.Fighters}
	if s.Action != nil {
		out.Action = &projectedAction{ID: s.Action.ID, Kind: s.Action.Kind, ActorID: s.Action.ActorID, Status: s.Action.Status}
	}
	if s.Combat != nil {
		out.Combat = &projectedCombat{ID: s.Combat.ID, AttackerID: s.Combat.AttackerID, DefenderID: s.Combat.DefenderID, Stage: s.Combat.Stage}
	}
	return json.Marshal(out)
}
func (e *Engine) meta(in contracts.ResolutionInput, t, operation string, data any) contracts.DomainEvent {
	ev, _ := operations.MetaEvent(t, in.Procedure.SourceRef, operation, "resolver", data, nil)
	return ev
}
func advance(ps *effects.State) {
	ps.Cursor++
	ps.Phase = effects.PhaseEnter
	ps.CostIndex = 0
	ps.OperationIndex = 0
	ps.Checkpoint = nil
}
func bind(ps *effects.State, o operations.Definition, r operations.Record) {
	if o.ResultBinding != "" {
		ps.Results[o.ResultBinding] = mustRaw(r)
	}
}
func bindRaw(m map[string]json.RawMessage, o operations.Definition, r operations.Record) {
	if o.ResultBinding != "" {
		m[o.ResultBinding] = mustRaw(r)
	}
}
func mustRaw(v any) json.RawMessage { b, _ := json.Marshal(v); return b }
func copyRaw(m map[string]json.RawMessage) map[string]json.RawMessage {
	r := map[string]json.RawMessage{}
	for k, v := range m {
		r[k] = append(json.RawMessage(nil), v...)
	}
	return r
}
func reject(code string) contracts.ResolutionOutcome {
	return contracts.ResolutionOutcome{Status: contracts.ResolutionRejected, RejectionCode: code}
}
func id(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:16])
}
func normalize(events []contracts.DomainEvent, state model.GameState, in contracts.ResolutionInput) {
	for i := range events {
		events[i].SchemaVersion = "rules-event/v1"
		events[i].ID = model.EventID(id(string(in.CommandID), fmt.Sprint(i), events[i].Type))
		events[i].MatchID = state.MatchID
		events[i].Sequence = state.EventSequence + uint64(i) + 1
		events[i].Revision = state.Revision + 1
		events[i].CausedByCommand = in.CommandID
		events[i].RulesetVersion = state.DefinitionRef.RulesetVersion
		if events[i].PrivatePayloads == nil {
			events[i].PrivatePayloads = map[model.PlayerID]json.RawMessage{}
		}
	}
}
func SortedEventTypes(events []contracts.DomainEvent) []string {
	r := make([]string, len(events))
	for i, e := range events {
		r[i] = e.Type
	}
	sort.Strings(r)
	return r
}
