// Package effects defines serializable staged procedures, choices and cancellation state.
package effects

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"sort"
)

const StateBinding = "__rules_procedure_state_v1"

type Choice struct {
	Kind, Binding, Visibility string
	Owner, Domain             query.Expr
	Prompt                    any
}
type Stage struct {
	ID                string
	Condition         *query.Expr
	Prerequisites     []query.Expr
	Costs, Operations []operations.Definition
	Choice            *Choice
	Checkpoint        string
}
type Definition struct {
	ID, Kind, RuleID string
	CapabilityIDs    []string
	Stages           []Stage
}
type Pending struct {
	ID                        model.InteractionID `json:"id"`
	Owner                     model.PlayerID      `json:"owner"`
	Kind, Visibility, Binding string
	Options                   map[string]json.RawMessage `json:"options"`
	Domain                    json.RawMessage            `json:"domain"`
	Prompt                    json.RawMessage            `json:"prompt"`
}
type State struct {
	DefinitionID string                     `json:"definition_id"`
	Cursor       int                        `json:"cursor"`
	Phase        string                     `json:"phase"`
	Captured     map[string]json.RawMessage `json:"captured,omitempty"`
	Results      map[string]json.RawMessage `json:"results,omitempty"`
	Choices      map[string]json.RawMessage `json:"choices,omitempty"`
	Pending      *Pending                   `json:"pending,omitempty"`
	Status       string                     `json:"status"`
}

func Decode(r model.ProcedureRef) (State, error) {
	raw, ok := r.Bindings[StateBinding]
	if !ok {
		return State{r.Kind, 0, "enter", copyMap(r.Bindings), map[string]json.RawMessage{}, map[string]json.RawMessage{}, nil, "running"}, nil
	}
	var s State
	if err := json.Unmarshal(raw, &s); err != nil {
		return s, err
	}
	if s.Captured == nil {
		s.Captured = map[string]json.RawMessage{}
	}
	if s.Results == nil {
		s.Results = map[string]json.RawMessage{}
	}
	if s.Choices == nil {
		s.Choices = map[string]json.RawMessage{}
	}
	return s, nil
}
func Encode(r model.ProcedureRef, s State, stage string) (model.ProcedureRef, error) {
	raw, e := json.Marshal(s)
	if e != nil {
		return r, e
	}
	r.Stage = stage
	r.Bindings = copyMap(r.Bindings)
	r.Bindings[StateBinding] = raw
	return r, nil
}
func Handle(id model.ProcedureID, stage string, i int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s\x00%s\x00%d", id, stage, i)))
	return hex.EncodeToString(h[:12])
}
func Validate(d Definition, r *operations.Registry) error {
	if d.ID == "" || d.Kind != "linear_stages" || len(d.Stages) == 0 {
		return fmt.Errorf("invalid procedure definition")
	}
	results := map[string]bool{}
	choices := map[string]bool{}
	ids := map[string]bool{}
	check := func(e query.Expr) error {
		if err := query.Validate(e); err != nil {
			return err
		}
		for _, x := range query.References(e, query.Result) {
			if !results[x] {
				return fmt.Errorf("result %q unavailable", x)
			}
		}
		for _, x := range query.References(e, query.Choice) {
			if !choices[x] {
				return fmt.Errorf("choice %q unavailable", x)
			}
		}
		return nil
	}
	for _, s := range d.Stages {
		if s.ID == "" || ids[s.ID] {
			return fmt.Errorf("duplicate/empty stage")
		}
		ids[s.ID] = true
		if s.Condition != nil {
			if e := check(*s.Condition); e != nil {
				return e
			}
		}
		for _, e := range s.Prerequisites {
			if x := check(e); x != nil {
				return x
			}
		}
		if s.Choice != nil {
			if choices[s.Choice.Binding] {
				return fmt.Errorf("duplicate choice")
			}
			if e := check(s.Choice.Owner); e != nil {
				return e
			}
			if e := check(s.Choice.Domain); e != nil {
				return e
			}
			choices[s.Choice.Binding] = true
		}
		ops := append(append([]operations.Definition{}, s.Costs...), s.Operations...)
		for _, o := range ops {
			if !r.Has(o.Kind) {
				return fmt.Errorf("unknown operation %q", o.Kind)
			}
			if o.Dependency != nil && !results[o.Dependency.Binding] {
				return fmt.Errorf("dependency unavailble")
			}
			for _, e := range o.Arguments {
				if x := check(e); x != nil {
					return x
				}
			}
			if o.ResultBinding != "" {
				if results[o.ResultBinding] {
					return fmt.Errorf("duplicate result")
				}
				results[o.ResultBinding] = true
			}
		}
	}
	_, e := json.Marshal(d)
	return e
}

type Queued struct {
	ID, Scope      string
	Priority, Order int
	Cancelable      bool
	Canceled        bool
}
type Queue []Queued

func (q *Queue) Cancel(scope string) []string {
	var out []string
	for i := range *q {
		if (*q)[i].Scope == scope && (*q)[i].Cancelable && !(*q)[i].Canceled {
			(*q)[i].Canceled = true
			out = append(out, (*q)[i].ID)
		}
	}
	sort.Strings(out)
	return out
}
func copyMap(m map[string]json.RawMessage) map[string]json.RawMessage {
	r := map[string]json.RawMessage{}
	for k, v := range m {
		if k != StateBinding {
			r[k] = append(json.RawMessage(nil), v...)
		}
	}
	return r
}
