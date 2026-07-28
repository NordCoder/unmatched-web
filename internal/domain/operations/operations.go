// Package operations provides typed generic mutation requests and result bindings.
package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
	"sort"
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
}
type Record struct {
	Applied bool   `json:"applied"`
	Code    string `json:"code,omitempty"`
	Value   any    `json:"value,omitempty"`
}
type Patch struct {
	FighterID model.FighterID
	Key       string
	Value     any
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
type Registry struct{ handlers map[string]Handler }

func Default() *Registry {
	r := &Registry{map[string]Handler{}}
	r.handlers[BindValue] = bind
	r.handlers[EmitEvent] = emit
	r.handlers[SetFighterState] = set
	r.handlers[Impossible] = impossible
	return r
}
func (r *Registry) Has(k string) bool { _, ok := r.handlers[k]; return ok }
func (r *Registry) Kinds() []string {
	xs := make([]string, 0, len(r.handlers))
	for k := range r.handlers {
		xs = append(xs, k)
	}
	sort.Strings(xs)
	return xs
}
func (r *Registry) Execute(c Context, d Definition) (Result, error) {
	h, ok := r.handlers[d.Kind]
	if !ok {
		return Result{}, fmt.Errorf("unknown operation kind %q", d.Kind)
	}
	return h(c, d)
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
	if _, ok := d.Arguments["payload"]; ok {
		p, e = arg(c, d, "payload")
		if e != nil {
			return Result{}, e
		}
	}
	raw, _ := json.Marshal(p)
	return Result{Record: Record{Applied: true, Value: p}, Events: []contracts.DomainEvent{{Type: t, SourceRef: c.SourceRef, PublicPayload: raw}}}, nil
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
	p := map[string]any{"fighter_instance_id": fs, "key": ks, "value": v}
	raw, _ := json.Marshal(p)
	return Result{Record: Record{Applied: true, Value: p}, Events: []contracts.DomainEvent{{Type: "rules.fighter_state_set", SourceRef: c.SourceRef, PublicPayload: raw}}, Patches: []Patch{{model.FighterID(fs), ks, v}}}, nil
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
