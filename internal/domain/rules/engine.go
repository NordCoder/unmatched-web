// Package rules implements the concrete pure staged RulesEngine used by Core.
package rules

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/NordCoder/unmatched-web/internal/domain/capabilities"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"sort"
)

type Engine struct {
	defs map[string]effects.Definition
	ops  *operations.Registry
}

func New(defs []effects.Definition) (*Engine, error) {
	ops := operations.Default()
	if err := capabilities.Validate(capabilities.Wave1()); err != nil {
		return nil, err
	}
	e := &Engine{map[string]effects.Definition{}, ops}
	for _, d := range defs {
		if err := effects.Validate(d, ops); err != nil {
			return nil, fmt.Errorf("definition %s: %w", d.ID, err)
		}
		if _, ok := e.defs[d.ID]; ok {
			return nil, fmt.Errorf("duplicate definition")
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
		actor := decode(in.Context["actor_player_id"])
		if actor != string(ps.Pending.Owner) {
			return reject("wrong_choice_owner"), nil
		}
		var h string
		if json.Unmarshal(in.Choice, &h) != nil {
			return reject("invalid_choice"), nil
		}
		v, ok := ps.Pending.Options[h]
		if !ok {
			return reject("choice_outside_domain"), nil
		}
		ps.Choices[ps.Pending.Binding] = v
		events = append(events, meta("rules.choice_accepted", in.Procedure.SourceRef, map[string]any{"interaction_instance_id": ps.Pending.ID, "binding": ps.Pending.Binding}))
		ps.Pending = nil
		ps.Status = "running"
	}
	for ps.Cursor < len(d.Stages) {
		s := d.Stages[ps.Cursor]
		ctx := query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}
		if s.Condition != nil {
			v, x := query.Eval(*s.Condition, ctx)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			if v != true {
				events = append(events, meta("rules.stage_skipped", in.Procedure.SourceRef, map[string]any{"stage": s.ID, "reason": "condition_false"}))
				ps.Cursor++
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
			events = append(events, meta("rules.stage_skipped", in.Procedure.SourceRef, map[string]any{"stage": s.ID, "reason": "prerequisite_failed"}))
			ps.Cursor++
			continue
		}
		costOK := true
		for i, o := range s.Costs {
			r, x := e.exec(working, ctx, in, s.ID, i, o)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			if !r.Record.Applied {
				costOK = false
				break
			}
			if x = operations.Apply(&working, r.Patches); x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			events = append(events, r.Events...)
			bind(&ps, o, r.Record)
		}
		if !costOK {
			events = append(events, meta("rules.stage_skipped", in.Procedure.SourceRef, map[string]any{"stage": s.ID, "reason": "cost_unpaid"}))
			ps.Cursor++
			continue
		}
		if s.Choice != nil {
			if _, ok := ps.Choices[s.Choice.Binding]; !ok {
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
				opts := map[string]json.RawMessage{}
				public := make([]any, 0, len(items))
				for i, item := range items {
					raw, _ := json.Marshal(item)
					h := effects.Handle(in.Procedure.ID, s.ID, i)
					opts[h] = raw
					if s.Choice.Visibility == "opaque" {
						public = append(public, map[string]any{"handle": h})
					} else {
						public = append(public, map[string]any{"handle": h, "value": item})
					}
				}
				domain, _ := json.Marshal(map[string]any{"options": public})
				prompt, _ := json.Marshal(s.Choice.Prompt)
				pid := model.InteractionID(id(string(in.Procedure.ID), s.ID, "interaction"))
				ps.Pending = &effects.Pending{ID: pid, Owner: model.PlayerID(owner), Kind: s.Choice.Kind, Visibility: s.Choice.Visibility, Binding: s.Choice.Binding, Options: opts, Domain: domain, Prompt: prompt}
				ps.Status = "pending"
				ref, _ := effects.Encode(in.Procedure, ps, s.ID)
				pending := &model.PendingInteraction{ID: pid, OwnerPlayerID: model.PlayerID(owner), Kind: s.Choice.Kind, Visibility: s.Choice.Visibility, Prompt: prompt, LegalDomain: domain, ResumeProcedure: ref}
				normalize(events, state, in)
				return contracts.ResolutionOutcome{Status: contracts.ResolutionPending, Events: events, PendingInteraction: pending, Procedure: &ref}, nil
			}
		}
		for i, o := range s.Operations {
			if o.Dependency != nil {
				var rec operations.Record
				_ = json.Unmarshal(ps.Results[o.Dependency.Binding], &rec)
				if o.Dependency.RequireApplied && !rec.Applied {
					events = append(events, meta("rules.operation_skipped", in.Procedure.SourceRef, map[string]any{"operation": o.ID, "reason": "dependency_failed"}))
					continue
				}
			}
			r, x := e.exec(working, query.Context{State: working, Captured: ps.Captured, Results: ps.Results, Choices: ps.Choices, Input: in.Context}, in, s.ID, i, o)
			if x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			if x = operations.Apply(&working, r.Patches); x != nil {
				return contracts.ResolutionOutcome{}, x
			}
			events = append(events, r.Events...)
			bind(&ps, o, r.Record)
			events = append(events, meta("rules.operation_result", in.Procedure.SourceRef, map[string]any{"operation": o.ID, "binding": o.ResultBinding, "applied": r.Record.Applied, "code": r.Record.Code}))
		}
		if s.Checkpoint != "" {
			events = append(events, meta("rules.checkpoint_opened", in.Procedure.SourceRef, map[string]any{"checkpoint": s.Checkpoint}))
		}
		ps.Cursor++
	}
	ps.Status = "completed"
	ref, _ := effects.Encode(in.Procedure, ps, "completed")
	normalize(events, state, in)
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted, Events: events, Procedure: &ref}, nil
}
func (e *Engine) exec(s model.GameState, q query.Context, in contracts.ResolutionInput, stage string, i int, o operations.Definition) (operations.Result, error) {
	return e.ops.Execute(operations.Context{State: s, Query: q, SourceRef: in.Procedure.SourceRef}, o)
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
func (e *Engine) Project(s model.GameState, p model.PlayerID) (json.RawMessage, error) {
	player, ok := s.Players[p]
	if !ok {
		return nil, fmt.Errorf("player not found")
	}
	fighters := make([]model.FighterID, 0, len(s.Fighters))
	for id := range s.Fighters {
		fighters = append(fighters, id)
	}
	sort.Slice(fighters, func(i, j int) bool { return fighters[i] < fighters[j] })
	v := map[string]any{"match_id": s.MatchID, "revision": s.Revision, "lifecycle": s.Lifecycle, "definition_ref": s.DefinitionRef, "self": player, "fighter_ids": fighters, "turn": s.Turn, "action": s.Action, "combat": s.Combat}
	if x := s.Resolver.PendingInteraction; x != nil && x.OwnerPlayerID == p {
		v["pending_interaction"] = x
	}
	return json.Marshal(v)
}
func bind(s *effects.State, o operations.Definition, r operations.Record) {
	if o.ResultBinding != "" {
		raw, _ := json.Marshal(r)
		s.Results[o.ResultBinding] = raw
	}
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
func meta(t, src string, p any) contracts.DomainEvent {
	raw, _ := json.Marshal(p)
	return contracts.DomainEvent{Type: t, SourceRef: src, PublicPayload: raw}
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
func decode(r json.RawMessage) any { var v any; _ = json.Unmarshal(r, &v); return v }
func id(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:16])
}
