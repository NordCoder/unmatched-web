// Package query implements the closed, pure expression language used by rules.
package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

type Kind string

const (
	Literal          Kind = "literal"
	Reference        Kind = "reference"
	List             Kind = "list"
	Object           Kind = "object"
	And              Kind = "and"
	Or               Kind = "or"
	Not              Kind = "not"
	Equal            Kind = "equal"
	Count            Kind = "count"
	Sum              Kind = "sum"
	If               Kind = "if"
	HighestThreshold Kind = "highest_met_threshold"
)

type Source string

const (
	State    Source = "state"
	Captured Source = "captured"
	Result   Source = "result"
	Choice   Source = "choice"
	Input    Source = "input"
)

type Expr struct {
	Kind   Kind            `json:"kind"`
	Value  any             `json:"value,omitempty"`
	Source Source          `json:"source,omitempty"`
	Path   []string        `json:"path,omitempty"`
	Args   []Expr          `json:"args,omitempty"`
	Fields map[string]Expr `json:"fields,omitempty"`
}
type Context struct {
	State                             model.GameState
	Captured, Results, Choices, Input map[string]json.RawMessage
}

func Eval(e Expr, c Context) (any, error) {
	switch e.Kind {
	case Literal:
		return e.Value, nil
	case Reference:
		return resolve(e, c)
	case List:
		r := make([]any, len(e.Args))
		for i, a := range e.Args {
			v, err := Eval(a, c)
			if err != nil {
				return nil, err
			}
			r[i] = v
		}
		return r, nil
	case Object:
		r := map[string]any{}
		keys := sorted(e.Fields)
		for _, k := range keys {
			v, err := Eval(e.Fields[k], c)
			if err != nil {
				return nil, err
			}
			r[k] = v
		}
		return r, nil
	case Not:
		v, err := one(e, c)
		if err != nil {
			return nil, err
		}
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("not requires bool")
		}
		return !b, nil
	case And, Or:
		for _, a := range e.Args {
			v, err := Eval(a, c)
			if err != nil {
				return nil, err
			}
			b, ok := v.(bool)
			if !ok {
				return nil, fmt.Errorf("%s requires bool", e.Kind)
			}
			if e.Kind == And && !b {
				return false, nil
			}
			if e.Kind == Or && b {
				return true, nil
			}
		}
		return e.Kind == And, nil
	case Equal:
		a, b, err := two(e, c)
		if err != nil {
			return nil, err
		}
		return reflect.DeepEqual(a, b), nil
	case Count:
		v, err := one(e, c)
		if err != nil {
			return nil, err
		}
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map && rv.Kind() != reflect.String {
			return nil, fmt.Errorf("count requires collection")
		}
		return rv.Len(), nil
	case Sum:
		v, err := one(e, c)
		if err != nil {
			return nil, err
		}
		items, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("sum requires list")
		}
		total := 0.0
		for _, x := range items {
			n, ok := number(x)
			if !ok {
				return nil, fmt.Errorf("sum requires numbers")
			}
			total += n
		}
		return total, nil
	case If:
		if len(e.Args) != 3 {
			return nil, fmt.Errorf("if requires 3 args")
		}
		v, err := Eval(e.Args[0], c)
		if err != nil {
			return nil, err
		}
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("if condition requires bool")
		}
		if b {
			return Eval(e.Args[1], c)
		}
		return Eval(e.Args[2], c)
	case HighestThreshold:
		if len(e.Args) != 2 {
			return nil, fmt.Errorf("highest threshold requires value and list")
		}
		v, err := Eval(e.Args[0], c)
		if err != nil {
			return nil, err
		}
		n, ok := number(v)
		if !ok {
			return nil, fmt.Errorf("threshold value requires number")
		}
		raw, err := Eval(e.Args[1], c)
		if err != nil {
			return nil, err
		}
		xs, ok := raw.([]any)
		if !ok {
			return nil, fmt.Errorf("thresholds require list")
		}
		best := 0.0
		found := false
		for _, x := range xs {
			t, ok := number(x)
			if !ok {
				return nil, fmt.Errorf("threshold requires number")
			}
			if t <= n && (!found || t > best) {
				best = t
				found = true
			}
		}
		if !found {
			return nil, nil
		}
		return best, nil
	default:
		return nil, fmt.Errorf("unknown expression kind %q", e.Kind)
	}
}
func Validate(e Expr) error {
	_, ok := map[Kind]bool{Literal: true, Reference: true, List: true, Object: true, And: true, Or: true, Not: true, Equal: true, Count: true, Sum: true, If: true, HighestThreshold: true}[e.Kind]
	if !ok {
		return fmt.Errorf("unknown expression kind %q", e.Kind)
	}
	if e.Kind == Reference && (e.Source == "" || len(e.Path) == 0) {
		return fmt.Errorf("reference requires source and path")
	}
	for _, a := range e.Args {
		if err := Validate(a); err != nil {
			return err
		}
	}
	for _, a := range e.Fields {
		if err := Validate(a); err != nil {
			return err
		}
	}
	return nil
}
func References(e Expr, s Source) []string {
	var out []string
	if e.Kind == Reference && e.Source == s && len(e.Path) > 0 {
		out = append(out, e.Path[0])
	}
	for _, a := range e.Args {
		out = append(out, References(a, s)...)
	}
	for _, a := range e.Fields {
		out = append(out, References(a, s)...)
	}
	sort.Strings(out)
	return out
}
func resolve(e Expr, c Context) (any, error) {
	var root any
	switch e.Source {
	case State:
		data, _ := json.Marshal(c.State)
		_ = json.Unmarshal(data, &root)
	case Captured:
		root = decode(c.Captured[e.Path[0]])
		e.Path = e.Path[1:]
	case Result:
		root = decode(c.Results[e.Path[0]])
		e.Path = e.Path[1:]
	case Choice:
		root = decode(c.Choices[e.Path[0]])
		e.Path = e.Path[1:]
	case Input:
		root = decode(c.Input[e.Path[0]])
		e.Path = e.Path[1:]
	default:
		return nil, fmt.Errorf("unknown source %q", e.Source)
	}
	for _, p := range e.Path {
		m, ok := root.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("cannot access %q", p)
		}
		root = m[p]
	}
	return root, nil
}
func decode(r json.RawMessage) any {
	if len(r) == 0 {
		return nil
	}
	var v any
	_ = json.Unmarshal(r, &v)
	return v
}
func one(e Expr, c Context) (any, error) {
	if len(e.Args) != 1 {
		return nil, fmt.Errorf("%s requires 1 arg", e.Kind)
	}
	return Eval(e.Args[0], c)
}
func two(e Expr, c Context) (any, any, error) {
	if len(e.Args) != 2 {
		return nil, nil, fmt.Errorf("%s requires 2 args", e.Kind)
	}
	a, err := Eval(e.Args[0], c)
	if err != nil {
		return nil, nil, err
	}
	b, err := Eval(e.Args[1], c)
	return a, b, err
}
func number(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}
func sorted(m map[string]Expr) []string {
	r := make([]string, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
