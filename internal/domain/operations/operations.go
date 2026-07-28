// Package operations provides typed generic mutation requests and result bindings.
package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
)

const (
	BindValue       = "bind_value"
	EmitEvent       = "emit_event"
	SetFighterState = "set_fighter_state"
	Impossible      = "impossible"
)

type Dependency struct {
	Binding        string `json:"binding"`
	RequireApplied bool   `json:"require_applied"`
}
type Definition struct {
	ID            string                `json:"id"`
	Kind          string                `json:"kind"`
	Arguments     map[string]query.Expr `json:"arguments,omitempty"`
	ResultBinding string                `json:"result_binding,omitempty"`
	Dependency    *Dependency           `json:"dependency,omitempty"`
	Cost          bool                  `json:"cost,omitempty"`
	CauseKind     string                `json:"cause_kind,omitempty"`
}
type Record struct {
	Applied bool   `json:"applied"`
	Code    string `json:"code,omitempty"`
	Value   any    `json:"value,omitempty"`
}
type Patch struct {
	FighterID model.FighterID `json:"fighter_instance_id"`
	Key       string          `json:"key"`
	Value     any             `json:"value"`
}
type Result struct {
	Record  Record
	Events  []contracts.DomainEvent
	Patches []Patch
}
type Context struct {
	State     model.GameState
	Query     query.Context
	SourceRef string
}
type Handler func(Context, Definition) (Result, error)
type ArgSpec struct {
	Type       query.Type
	Required   bool
	PublicOnly bool
}
type Schema struct {
	Arguments   map[string]ArgSpec
	Result      query.ValueSpec
	CostAllowed bool
	Handler     Handler
}
type Registry struct{ entries map[string]Schema }

type Provenance struct {
	OperationID  string   `json:"operation_instance_id"`
	SourceRef    string   `json:"source_ref,omitempty"`
	CauseKind    string   `json:"cause_kind"`
	Participants []string `json:"participants,omitempty"`
}
type EventEnvelope struct {
	Schema     string     `json:"schema"`
	Data       any        `json:"data"`
	Provenance Provenance `json:"provenance"`
}

func Default() *Registry {
	r := &Registry{entries: map[string]Schema{}}
	r.entries[BindValue] = Schema{Arguments: map[string]ArgSpec{"value": {Type: query.TypeAny, Required: true}}, Result: query.ValueSpec{Type: query.TypeAny, Visibility: query.OwnerPrivate}, CostAllowed: false, Handler: bind}
	r.entries[EmitEvent] = Schema{Arguments: map[string]ArgSpec{"event_type": {Type: query.TypeString, Required: true, PublicOnly: true}, "payload": {Type: query.TypeObject, PublicOnly: true}}, Result: query.ValueSpec{Type: query.TypeObject, Visibility: query.Public}, CostAllowed: false, Handler: emit}
	r.entries[SetFighterState] = Schema{Arguments: map[string]ArgSpec{"fighter": {Type: query.TypeFighter, Required: true, PublicOnly: true}, "key": {Type: query.TypeString, Required: true, PublicOnly: true}, "value": {Type: query.TypeAny, Required: true, PublicOnly: true}}, Result: query.ValueSpec{Type: query.TypeOperationResult, Visibility: query.Public}, CostAllowed: true, Handler: set}
	r.entries[Impossible] = Schema{Arguments: map[string]ArgSpec{"code": {Type: query.TypeString, PublicOnly: true}}, Result: query.ValueSpec{Type: query.TypeOperationResult, Visibility: query.Public}, CostAllowed: true, Handler: impossible}
	return r
}
func (r *Registry) Has(k string) bool { _, ok := r.entries[k]; return ok }
func (r *Registry) Kinds() []string {
	xs := make([]string, 0, len(r.entries))
	for k := range r.entries {
		xs = append(xs, k)
	}
	sort.Strings(xs)
	return xs
}
func (r *Registry) Schema(k string) (Schema, bool) { s, ok := r.entries[k]; return s, ok }
func (r *Registry) Validate(d Definition, env query.TypeEnv) (query.ValueSpec, error) {
	s, ok := r.entries[d.Kind]
	if !ok {
		return query.ValueSpec{}, fmt.Errorf("unknown operation kind %q", d.Kind)
	}
	if d.ID == "" {
		return query.ValueSpec{}, fmt.Errorf("operation ID is required")
	}
	if d.Cost && !s.CostAllowed {
		return query.ValueSpec{}, fmt.Errorf("operation %q cannot be a cost", d.Kind)
	}
	for name, spec := range s.Arguments {
		_, present := d.Arguments[name]
		if spec.Required && !present {
			return query.ValueSpec{}, fmt.Errorf("operation %q missing %q", d.ID, name)
		}
	}
	actualArgs := map[string]query.ValueSpec{}
	for name, e := range d.Arguments {
		spec, known := s.Arguments[name]
		if !known {
			return query.ValueSpec{}, fmt.Errorf("operation %q has unknown argument %q", d.ID, name)
		}
		actual, err := query.Infer(e, env)
		if err != nil {
			return query.ValueSpec{}, fmt.Errorf("operation %s argument %s: %w", d.ID, name, err)
		}
		if spec.Type != query.TypeAny && !typeCompatible(actual.Type, spec.Type) {
			return query.ValueSpec{}, fmt.Errorf("operation %s argument %s requires %s, got %s", d.ID, name, spec.Type, actual.Type)
		}
		if spec.PublicOnly && actual.Visibility != query.Public {
			return query.ValueSpec{}, fmt.Errorf("operation %s argument %s cannot expose %s data", d.ID, name, actual.Visibility)
		}
		actualArgs[name] = actual
	}
	if d.Kind == BindValue {
		return actualArgs["value"], nil
	}
	return s.Result, nil
}
func (r *Registry) Execute(c Context, d Definition) (Result, error) {
	s, ok := r.entries[d.Kind]
	if !ok {
		return Result{}, fmt.Errorf("unknown operation kind %q", d.Kind)
	}
	return s.Handler(c, d)
}
func Apply(s *model.GameState, ps []Patch) error {
	for _, p := range ps {
		o, ok := s.Fighters[p.FighterID]
		if !ok {
			return fmt.Errorf("fighter %q not found", p.FighterID)
		}
		if o.State == nil {
			o.State = map[string]any{}
		}
		o.State[p.Key] = p.Value
		s.Fighters[p.FighterID] = o
	}
	return nil
}
func ApplyEvent(s model.GameState, e contracts.DomainEvent) (model.GameState, error) {
	var env EventEnvelope
	if err := json.Unmarshal(e.PublicPayload, &env); err != nil || env.Schema != "rules-event/v1" {
		return s, nil
	}
	if e.Type != "rules.fighter_state_set" {
		return s, nil
	}
	b, err := json.Marshal(env.Data)
	if err != nil {
		return s, err
	}
	var p Patch
	if err = json.Unmarshal(b, &p); err != nil {
		return s, err
	}
	if err = Apply(&s, []Patch{p}); err != nil {
		return s, err
	}
	return s, nil
}
func bind(c Context, d Definition) (Result, error) {
	v, e := arg(c, d, "value")
	return Result{Record: Record{Applied: e == nil, Value: v}}, e
}
func emit(c Context, d Definition) (Result, error) {
	v, e := arg(c, d, "event_type")
	if e != nil {
		return Result{}, e
	}
	t, ok := v.(string)
	if !ok || t == "" {
		return Result{}, fmt.Errorf("event_type must be string")
	}
	p := any(map[string]any{})
	if _, ok = d.Arguments["payload"]; ok {
		p, e = arg(c, d, "payload")
		if e != nil {
			return Result{}, e
		}
	}
	ev, err := event(t, c, d, p, nil)
	if err != nil {
		return Result{}, err
	}
	return Result{Record: Record{Applied: true, Value: p}, Events: []contracts.DomainEvent{ev}}, nil
}
func set(c Context, d Definition) (Result, error) {
	f, e := arg(c, d, "fighter")
	if e != nil {
		return Result{}, e
	}
	k, e := arg(c, d, "key")
	if e != nil {
		return Result{}, e
	}
	v, e := arg(c, d, "value")
	if e != nil {
		return Result{}, e
	}
	fs, fok := f.(string)
	ks, kok := k.(string)
	if !fok || !kok {
		return Result{}, fmt.Errorf("fighter and key must be strings")
	}
	if _, ok := c.State.Fighters[model.FighterID(fs)]; !ok {
		return Result{Record: Record{Code: "fighter_not_found"}}, nil
	}
	p := Patch{FighterID: model.FighterID(fs), Key: ks, Value: v}
	ev, err := event("rules.fighter_state_set", c, d, p, []string{fs})
	if err != nil {
		return Result{}, err
	}
	return Result{Record: Record{Applied: true, Value: p}, Events: []contracts.DomainEvent{ev}, Patches: []Patch{p}}, nil
}
func impossible(c Context, d Definition) (Result, error) {
	code := "impossible"
	if _, ok := d.Arguments["code"]; ok {
		v, e := arg(c, d, "code")
		if e != nil {
			return Result{}, e
		}
		if s, ok := v.(string); ok {
			code = s
		}
	}
	return Result{Record: Record{Code: code}}, nil
}
func event(t string, c Context, d Definition, data any, participants []string) (contracts.DomainEvent, error) {
	cause := d.CauseKind
	if cause == "" {
		cause = "effect"
	}
	env := EventEnvelope{Schema: "rules-event/v1", Data: data, Provenance: Provenance{OperationID: d.ID, SourceRef: c.SourceRef, CauseKind: cause, Participants: participants}}
	raw, err := json.Marshal(env)
	if err != nil {
		return contracts.DomainEvent{}, err
	}
	return contracts.DomainEvent{Type: t, SourceRef: c.SourceRef, PublicPayload: raw}, nil
}
func MetaEvent(t, source, operation, cause string, data any, participants []string) (contracts.DomainEvent, error) {
	if cause == "" {
		cause = "resolver"
	}
	env := EventEnvelope{Schema: "rules-event/v1", Data: data, Provenance: Provenance{OperationID: operation, SourceRef: source, CauseKind: cause, Participants: participants}}
	raw, err := json.Marshal(env)
	return contracts.DomainEvent{Type: t, SourceRef: source, PublicPayload: raw}, err
}
func arg(c Context, d Definition, n string) (any, error) {
	e, ok := d.Arguments[n]
	if !ok {
		return nil, fmt.Errorf("operation %q missing %q", d.ID, n)
	}
	return query.Eval(e, c.Query)
}
func DeterministicID(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:16])
}
func typeCompatible(a, b query.Type) bool {
	return a == b || a == query.TypeAny || b == query.TypeAny || (a == query.TypeString && (b == query.TypePlayer || b == query.TypeFighter)) || (b == query.TypeString && (a == query.TypePlayer || a == query.TypeFighter))
}
