// Package effects defines serializable staged procedures, choices, checkpoints and cancellation state.
package effects

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/domain/capabilities"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
)

const StateBinding = "__rules_procedure_state_v2"

type Phase string

const (
	PhaseEnter      Phase = "enter"
	PhaseCosts      Phase = "costs"
	PhaseChoice     Phase = "choice"
	PhaseOperations Phase = "operations"
	PhaseCheckpoint Phase = "checkpoint"
	PhaseComplete   Phase = "complete"
)

type EmptyDomainPolicy string

const (
	EmptySkipStage EmptyDomainPolicy = "skip_stage"
	EmptyBindEmpty EmptyDomainPolicy = "bind_empty"
	EmptyReject    EmptyDomainPolicy = "reject"
	EmptyComplete  EmptyDomainPolicy = "complete_without_choice"
)

type Choice struct {
	Kind        string            `json:"kind"`
	Binding     string            `json:"binding"`
	Visibility  query.Visibility  `json:"visibility"`
	Owner       query.Expr        `json:"owner"`
	Domain      query.Expr        `json:"domain"`
	Prompt      any               `json:"prompt"`
	EmptyDomain EmptyDomainPolicy `json:"empty_domain"`
	ValueType   query.Type        `json:"value_type"`
}
type QueueDefinition struct {
	ID          string                  `json:"id"`
	Scope       string                  `json:"scope"`
	Priority    int                     `json:"priority"`
	SourceOrder int                     `json:"source_order"`
	Cancelable  bool                    `json:"cancelable"`
	Operations  []operations.Definition `json:"operations"`
}
type Cancellation struct {
	Scope string `json:"scope"`
}
type Stage struct {
	ID            string                  `json:"id"`
	Condition     *query.Expr             `json:"condition,omitempty"`
	Prerequisites []query.Expr            `json:"prerequisites,omitempty"`
	Costs         []operations.Definition `json:"costs,omitempty"`
	Choice        *Choice                 `json:"choice,omitempty"`
	Operations    []operations.Definition `json:"operations,omitempty"`
	Checkpoint    string                  `json:"checkpoint,omitempty"`
	Queue         []QueueDefinition       `json:"queue,omitempty"`
	Cancellations []Cancellation          `json:"cancellations,omitempty"`
}
type Definition struct {
	ID               string                     `json:"id"`
	Kind             string                     `json:"kind"`
	RuleID           string                     `json:"rule_id"`
	CapabilityIDs    []string                   `json:"capability_ids"`
	CapturedBindings map[string]query.ValueSpec `json:"captured_bindings,omitempty"`
	InputBindings    map[string]query.ValueSpec `json:"input_bindings,omitempty"`
	Stages           []Stage                    `json:"stages"`
}
type Pending struct {
	ID         model.InteractionID        `json:"id"`
	Owner      model.PlayerID             `json:"owner"`
	Kind       string                     `json:"kind"`
	Visibility query.Visibility           `json:"visibility"`
	Binding    string                     `json:"binding"`
	Options    map[string]json.RawMessage `json:"options"`
	Domain     json.RawMessage            `json:"domain"`
	Prompt     json.RawMessage            `json:"prompt"`
}
type QueueEntry struct {
	ID          string                  `json:"id"`
	Scope       string                  `json:"scope"`
	Priority    int                     `json:"priority"`
	SourceOrder int                     `json:"source_order"`
	Cancelable  bool                    `json:"cancelable"`
	Operations  []operations.Definition `json:"operations"`
	Status      string                  `json:"status"`
}
type CheckpointFrame struct {
	ID          string       `json:"id"`
	ResumeStage int          `json:"resume_stage"`
	Queue       []QueueEntry `json:"queue"`
	NextIndex   int          `json:"next_index"`
}
type State struct {
	DefinitionID   string                     `json:"definition_id"`
	Cursor         int                        `json:"cursor"`
	Phase          Phase                      `json:"phase"`
	CostIndex      int                        `json:"cost_index"`
	OperationIndex int                        `json:"operation_index"`
	Captured       map[string]json.RawMessage `json:"captured,omitempty"`
	Results        map[string]json.RawMessage `json:"results,omitempty"`
	Choices        map[string]json.RawMessage `json:"choices,omitempty"`
	Pending        *Pending                   `json:"pending,omitempty"`
	Checkpoint     *CheckpointFrame           `json:"checkpoint,omitempty"`
	Status         string                     `json:"status"`
}

func Decode(r model.ProcedureRef) (State, error) {
	raw, ok := r.Bindings[StateBinding]
	if !ok {
		return State{DefinitionID: r.Kind, Cursor: 0, Phase: PhaseEnter, Captured: copyMap(r.Bindings), Results: map[string]json.RawMessage{}, Choices: map[string]json.RawMessage{}, Status: "running"}, nil
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
func Validate(d Definition, r *operations.Registry, caps map[string]capabilities.Capability) error {
	if d.ID == "" || d.Kind != "linear_stages" || len(d.Stages) == 0 {
		return fmt.Errorf("invalid procedure definition")
	}
	if d.RuleID == "" {
		return fmt.Errorf("rule ID is required")
	}
	if err := capabilities.ValidateDeclaration(d.CapabilityIDs, caps); err != nil {
		return err
	}
	env := query.TypeEnv{Captured: copySpecs(d.CapturedBindings), Results: map[string]query.ValueSpec{}, Choices: map[string]query.ValueSpec{}, Input: copySpecs(d.InputBindings)}
	ids := map[string]bool{}
	for _, s := range d.Stages {
		if s.ID == "" || ids[s.ID] {
			return fmt.Errorf("duplicate/empty stage")
		}
		ids[s.ID] = true
		if s.Condition != nil {
			spec, err := query.Infer(*s.Condition, env)
			if err != nil {
				return fmt.Errorf("stage %s condition: %w", s.ID, err)
			}
			if spec.Type != query.TypeBool {
				return fmt.Errorf("stage %s condition must be bool", s.ID)
			}
		}
		for _, p := range s.Prerequisites {
			spec, err := query.Infer(p, env)
			if err != nil {
				return fmt.Errorf("stage %s prerequisite: %w", s.ID, err)
			}
			if spec.Type != query.TypeBool {
				return fmt.Errorf("stage %s prerequisite must be bool", s.ID)
			}
		}
		for _, o := range s.Costs {
			if !o.Cost {
				return fmt.Errorf("stage %s cost %s must declare cost=true", s.ID, o.ID)
			}
			spec, err := r.Validate(o, env)
			if err != nil {
				return err
			}
			if o.ResultBinding != "" {
				if _, exists := env.Results[o.ResultBinding]; exists {
					return fmt.Errorf("duplicate result %s", o.ResultBinding)
				}
				env.Results[o.ResultBinding] = spec
			}
		}
		if s.Choice != nil {
			ch := s.Choice
			if ch.Binding == "" || ch.Kind == "" {
				return fmt.Errorf("stage %s malformed choice", s.ID)
			}
			if _, exists := env.Choices[ch.Binding]; exists {
				return fmt.Errorf("duplicate choice %s", ch.Binding)
			}
			if ch.Visibility != query.Public && ch.Visibility != query.OwnerPrivate && ch.Visibility != query.Opaque {
				return fmt.Errorf("stage %s invalid visibility", s.ID)
			}
			if ch.EmptyDomain != EmptySkipStage && ch.EmptyDomain != EmptyBindEmpty && ch.EmptyDomain != EmptyReject && ch.EmptyDomain != EmptyComplete {
				return fmt.Errorf("stage %s requires empty-domain policy", s.ID)
			}
			owner, err := query.Infer(ch.Owner, env)
			if err != nil {
				return err
			}
			if owner.Type != query.TypePlayer && owner.Type != query.TypeString {
				return fmt.Errorf("choice owner must be player_ref")
			}
			domain, err := query.Infer(ch.Domain, env)
			if err != nil {
				return err
			}
			if domain.Type != query.TypeList {
				return fmt.Errorf("choice domain must be list")
			}
			valueType := ch.ValueType
			if valueType == "" && domain.Element != nil {
				valueType = domain.Element.Type
			}
			if valueType == "" {
				valueType = query.TypeAny
			}
			env.Choices[ch.Binding] = query.ValueSpec{Type: valueType, Visibility: ch.Visibility}
		}
		for _, o := range s.Operations {
			if o.Cost {
				return fmt.Errorf("operation %s belongs in costs", o.ID)
			}
			if o.Dependency != nil {
				if _, ok := env.Results[o.Dependency.Binding]; !ok {
					return fmt.Errorf("dependency unavailable: %s", o.Dependency.Binding)
				}
			}
			spec, err := r.Validate(o, env)
			if err != nil {
				return err
			}
			if o.ResultBinding != "" {
				if _, exists := env.Results[o.ResultBinding]; exists {
					return fmt.Errorf("duplicate result %s", o.ResultBinding)
				}
				env.Results[o.ResultBinding] = spec
			}
		}
		if s.Checkpoint == "" && (len(s.Queue) > 0 || len(s.Cancellations) > 0) {
			return fmt.Errorf("stage %s queue/cancellation requires checkpoint", s.ID)
		}
		qids := map[string]bool{}
		for _, q := range s.Queue {
			if q.ID == "" || qids[q.ID] {
				return fmt.Errorf("stage %s duplicate queue ID", s.ID)
			}
			if len(q.Operations) == 0 {
				return fmt.Errorf("queued effect %s has no operations", q.ID)
			}
			for _, operation := range q.Operations {
				if operation.Cost || operation.ResultBinding != "" {
					return fmt.Errorf("queued effect %s operations cannot be costs or bind results", q.ID)
				}
				if _, err := r.Validate(operation, env); err != nil {
					return err
				}
			}
			qids[q.ID] = true
		}
	}
	_, e := json.Marshal(d)
	return e
}
func BuildCheckpoint(id string, resume int, defs []QueueDefinition) *CheckpointFrame {
	q := make([]QueueEntry, len(defs))
	for i, d := range defs {
		q[i] = QueueEntry{ID: d.ID, Scope: d.Scope, Priority: d.Priority, SourceOrder: d.SourceOrder, Cancelable: d.Cancelable, Operations: append([]operations.Definition(nil), d.Operations...), Status: "queued"}
	}
	sort.SliceStable(q, func(i, j int) bool {
		if q[i].Priority != q[j].Priority {
			return q[i].Priority > q[j].Priority
		}
		if q[i].SourceOrder != q[j].SourceOrder {
			return q[i].SourceOrder < q[j].SourceOrder
		}
		return q[i].ID < q[j].ID
	})
	return &CheckpointFrame{ID: id, ResumeStage: resume, Queue: q}
}
func (f *CheckpointFrame) Cancel(scope string) []string {
	var out []string
	for i := range f.Queue {
		q := &f.Queue[i]
		if q.Scope == scope && q.Cancelable && q.Status == "queued" {
			q.Status = "canceled"
			out = append(out, q.ID)
		}
	}
	sort.Strings(out)
	return out
}
func (f *CheckpointFrame) Next() *QueueEntry {
	for f.NextIndex < len(f.Queue) {
		i := f.NextIndex
		f.NextIndex++
		if f.Queue[i].Status == "queued" {
			f.Queue[i].Status = "resolved"
			return &f.Queue[i]
		}
	}
	return nil
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
func copySpecs(m map[string]query.ValueSpec) map[string]query.ValueSpec {
	r := map[string]query.ValueSpec{}
	for k, v := range m {
		r[k] = v
	}
	return r
}
