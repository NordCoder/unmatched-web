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
	EmptySkipStage   EmptyDomainPolicy = "skip_stage"
	EmptyBindEmpty   EmptyDomainPolicy = "bind_empty"
	EmptyReject      EmptyDomainPolicy = "reject"
	EmptyComplete    EmptyDomainPolicy = "complete_without_choice"
	EmptyBindDefault EmptyDomainPolicy = "bind_default"
)

type Choice struct {
	Kind        string            `json:"kind"`
	Binding     string            `json:"binding"`
	Visibility  query.Visibility  `json:"visibility"`
	Owner       query.Expr        `json:"owner"`
	Domain      query.Expr        `json:"domain"`
	Prompt      any               `json:"prompt"`
	EmptyDomain EmptyDomainPolicy `json:"empty_domain"`
	ValueType   query.Type        `json:"value_type,omitempty"`
	Default     *query.Expr       `json:"default,omitempty"`
	Multi       bool              `json:"multi,omitempty"`
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
	ValueSpec  query.ValueSpec            `json:"value_spec"`
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
	Status      operations.Disposition  `json:"status"`
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
	raw, err := json.Marshal(s)
	if err != nil {
		return r, err
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
	for si, s := range d.Stages {
		if s.ID == "" || ids[s.ID] {
			return fmt.Errorf("duplicate/empty stage")
		}
		ids[s.ID] = true
		if s.Condition != nil {
			sp, err := query.Infer(*s.Condition, env)
			if err != nil {
				return fmt.Errorf("stage %s condition: %w", s.ID, err)
			}
			if sp.Type != query.TypeBool {
				return fmt.Errorf("stage %s condition must be bool", s.ID)
			}
		}
		for _, p := range s.Prerequisites {
			sp, err := query.Infer(p, env)
			if err != nil {
				return fmt.Errorf("stage %s prerequisite: %w", s.ID, err)
			}
			if sp.Type != query.TypeBool {
				return fmt.Errorf("stage %s prerequisite must be bool", s.ID)
			}
		}
		for _, o := range s.Costs {
			if !o.Cost {
				return fmt.Errorf("stage %s cost %s must declare cost=true", s.ID, o.ID)
			}
			sp, err := r.Validate(o, env)
			if err != nil {
				return err
			}
			if o.ResultBinding != "" {
				if _, exists := env.Results[o.ResultBinding]; exists {
					return fmt.Errorf("duplicate result %s", o.ResultBinding)
				}
				env.Results[o.ResultBinding] = sp
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
			if ch.EmptyDomain != EmptySkipStage && ch.EmptyDomain != EmptyBindEmpty && ch.EmptyDomain != EmptyReject && ch.EmptyDomain != EmptyComplete && ch.EmptyDomain != EmptyBindDefault {
				return fmt.Errorf("stage %s requires empty-domain policy", s.ID)
			}
			owner, err := query.Infer(ch.Owner, env)
			if err != nil {
				return err
			}
			if owner.Type != query.TypePlayer {
				return fmt.Errorf("choice owner must be player_ref")
			}
			domain, err := query.Infer(ch.Domain, env)
			if err != nil {
				return err
			}
			if domain.Type != query.TypeList || domain.Element == nil {
				return fmt.Errorf("choice domain must be typed list")
			}
			value := *domain.Element
			if ch.Multi {
				value = query.ValueSpec{Type: query.TypeList, Element: domain.Element, Visibility: domain.Visibility}
			}
			if ch.ValueType != "" && ch.ValueType != value.Type {
				return fmt.Errorf("choice %s declares %s but domain yields %s", ch.Binding, ch.ValueType, value.Type)
			}
			if query.VisibilityRank(ch.Visibility) < query.VisibilityRank(domain.Visibility) {
				return fmt.Errorf("choice %s cannot weaken %s domain to %s", ch.Binding, domain.Visibility, ch.Visibility)
			}
			value.Visibility = query.MaxVisibility(domain.Visibility, ch.Visibility)
			guaranteed := false
			switch ch.EmptyDomain {
			case EmptyReject:
				guaranteed = true
			case EmptyBindDefault:
				if ch.Default == nil {
					return fmt.Errorf("choice %s bind_default requires default", ch.Binding)
				}
				ds, err := query.Infer(*ch.Default, env)
				if err != nil {
					return err
				}
				if !query.SameShape(ds, value) {
					return fmt.Errorf("choice %s default type mismatch", ch.Binding)
				}
				if query.VisibilityRank(ds.Visibility) > query.VisibilityRank(value.Visibility) {
					value.Visibility = ds.Visibility
				}
				guaranteed = true
			case EmptyBindEmpty:
				if !ch.Multi || value.Type != query.TypeList {
					return fmt.Errorf("choice %s bind_empty requires multi list choice", ch.Binding)
				}
				guaranteed = true
			case EmptySkipStage, EmptyComplete:
				if readsChoiceFrom(d.Stages[si:], ch.Binding) {
					return fmt.Errorf("choice %s may be absent on continuing path", ch.Binding)
				}
			}
			if guaranteed {
				env.Choices[ch.Binding] = value
			}
		}
		for _, o := range s.Operations {
			if o.Cost {
				return fmt.Errorf("operation %s belongs in costs", o.ID)
			}
			if o.Dependency != nil {
				dep, ok := env.Results[o.Dependency.Binding]
				if !ok {
					return fmt.Errorf("dependency unavailable: %s", o.Dependency.Binding)
				}
				disp, err := query.Traverse(dep, []string{"disposition"})
				if err != nil || disp.Type != query.TypeDisposition {
					return fmt.Errorf("dependency %s lacks disposition", o.Dependency.Binding)
				}
			}
			sp, err := r.Validate(o, env)
			if err != nil {
				return err
			}
			if o.ResultBinding != "" {
				if _, exists := env.Results[o.ResultBinding]; exists {
					return fmt.Errorf("duplicate result %s", o.ResultBinding)
				}
				env.Results[o.ResultBinding] = sp
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
			for _, o := range q.Operations {
				if o.Cost || o.ResultBinding != "" {
					return fmt.Errorf("queued effect %s operations cannot be costs or bind results", q.ID)
				}
				if _, err := r.Validate(o, env); err != nil {
					return err
				}
			}
			qids[q.ID] = true
		}
	}
	_, err := json.Marshal(d)
	return err
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
			q.Status = operations.DispositionCanceled
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
			f.Queue[i].Status = "executing"
			return &f.Queue[i]
		}
	}
	return nil
}
func (f *CheckpointFrame) Finish(id string, d operations.Disposition) {
	for i := range f.Queue {
		if f.Queue[i].ID == id {
			f.Queue[i].Status = d
			return
		}
	}
}
func copyMap(m map[string]json.RawMessage) map[string]json.RawMessage {
	r := map[string]json.RawMessage{}
	for k, v := range m {
		r[k] = append(json.RawMessage(nil), v...)
	}
	return r
}
func copySpecs(m map[string]query.ValueSpec) map[string]query.ValueSpec {
	r := map[string]query.ValueSpec{}
	for k, v := range m {
		b, _ := json.Marshal(v)
		var x query.ValueSpec
		_ = json.Unmarshal(b, &x)
		r[k] = x
	}
	return r
}
func readsChoiceFrom(stages []Stage, name string) bool {
	for _, s := range stages {
		if s.Condition != nil && query.ReferencesBinding(*s.Condition, query.Choice, name) {
			return true
		}
		for _, p := range s.Prerequisites {
			if query.ReferencesBinding(p, query.Choice, name) {
				return true
			}
		}
		for _, o := range append(append([]operations.Definition{}, s.Costs...), s.Operations...) {
			for _, a := range o.Arguments {
				if query.ReferencesBinding(a, query.Choice, name) {
					return true
				}
			}
		}
		for _, q := range s.Queue {
			for _, o := range q.Operations {
				for _, a := range o.Arguments {
					if query.ReferencesBinding(a, query.Choice, name) {
						return true
					}
				}
			}
		}
	}
	return false
}

// ChoiceSpec returns the validator-derived structural value specification for a guaranteed choice.
func ChoiceSpec(d Definition, binding string, r *operations.Registry) (query.ValueSpec, error) {
	env := query.TypeEnv{Captured: copySpecs(d.CapturedBindings), Results: map[string]query.ValueSpec{}, Choices: map[string]query.ValueSpec{}, Input: copySpecs(d.InputBindings)}
	for _, s := range d.Stages {
		for _, o := range s.Costs {
			sp, err := r.Validate(o, env)
			if err != nil {
				return query.ValueSpec{}, err
			}
			if o.ResultBinding != "" {
				env.Results[o.ResultBinding] = sp
			}
		}
		if s.Choice != nil {
			ch := s.Choice
			domain, err := query.Infer(ch.Domain, env)
			if err != nil {
				return query.ValueSpec{}, err
			}
			if domain.Element == nil {
				return query.ValueSpec{}, fmt.Errorf("choice %s has no element", ch.Binding)
			}
			value := *domain.Element
			if ch.Multi {
				value = query.ValueSpec{Type: query.TypeList, Element: domain.Element, Visibility: domain.Visibility}
			}
			value.Visibility = query.MaxVisibility(domain.Visibility, ch.Visibility)
			if ch.Binding == binding {
				return value, nil
			}
			if ch.EmptyDomain == EmptyReject || ch.EmptyDomain == EmptyBindDefault || ch.EmptyDomain == EmptyBindEmpty {
				env.Choices[ch.Binding] = value
			}
		}
		for _, o := range s.Operations {
			sp, err := r.Validate(o, env)
			if err != nil {
				return query.ValueSpec{}, err
			}
			if o.ResultBinding != "" {
				env.Results[o.ResultBinding] = sp
			}
		}
	}
	return query.ValueSpec{}, fmt.Errorf("unknown choice binding %s", binding)
}
