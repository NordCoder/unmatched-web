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

type Type string

const (
	TypeAny             Type = "any"
	TypeBool            Type = "bool"
	TypeNumber          Type = "number"
	TypeString          Type = "string"
	TypePlayer          Type = "player_ref"
	TypeFighter         Type = "fighter_ref"
	TypeList            Type = "list"
	TypeObject          Type = "object"
	TypeOperationResult Type = "operation_result"
)

type Visibility string

const (
	Public       Visibility = "public"
	OwnerPrivate Visibility = "owner_private"
	Opaque       Visibility = "opaque"
)

type ValueSpec struct {
	Type       Type       `json:"type"`
	Element    *ValueSpec `json:"element,omitempty"`
	Visibility Visibility `json:"visibility"`
}
type Expr struct {
	Kind       Kind            `json:"kind"`
	Value      any             `json:"value,omitempty"`
	Source     Source          `json:"source,omitempty"`
	Path       []string        `json:"path,omitempty"`
	Args       []Expr          `json:"args,omitempty"`
	Fields     map[string]Expr `json:"fields,omitempty"`
	ValueType  Type            `json:"value_type,omitempty"`
	Visibility Visibility      `json:"visibility,omitempty"`
}
type Context struct {
	State                             model.GameState
	Captured, Results, Choices, Input map[string]json.RawMessage
}
type TypeEnv struct{ Captured, Results, Choices, Input map[string]ValueSpec }

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
		for _, k := range sorted(e.Fields) {
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
		if len(e.Args) < 2 {
			return nil, fmt.Errorf("%s requires at least 2 args", e.Kind)
		}
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
		if !rv.IsValid() {
			return nil, fmt.Errorf("count requires collection")
		}
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

func Infer(e Expr, env TypeEnv) (ValueSpec, error) {
	vis := e.Visibility
	if vis == "" {
		vis = Public
	}
	switch e.Kind {
	case Literal:
		return literalSpec(e.Value, vis), nil
	case Reference:
		if len(e.Path) == 0 {
			return ValueSpec{}, fmt.Errorf("reference requires path")
		}
		var spec ValueSpec
		var ok bool
		switch e.Source {
		case State:
			spec, ok = stateSpec(e)
		case Captured:
			spec, ok = env.Captured[e.Path[0]]
		case Result:
			spec, ok = env.Results[e.Path[0]]
		case Choice:
			spec, ok = env.Choices[e.Path[0]]
		case Input:
			spec, ok = env.Input[e.Path[0]]
		default:
			return ValueSpec{}, fmt.Errorf("unknown source %q", e.Source)
		}
		if !ok {
			if e.ValueType == "" {
				return ValueSpec{}, fmt.Errorf("untyped reference %s:%v", e.Source, e.Path)
			}
			spec = ValueSpec{Type: e.ValueType, Visibility: vis}
		}
		if len(e.Path) > 1 && e.ValueType != "" {
			spec.Type = e.ValueType
			spec.Element = nil
		} else if e.ValueType != "" && !compatible(spec.Type, e.ValueType) {
			return ValueSpec{}, fmt.Errorf("reference declared %s but source is %s", e.ValueType, spec.Type)
		}
		if visRank(vis) > visRank(spec.Visibility) {
			spec.Visibility = vis
		}
		return spec, nil
	case List:
		if len(e.Args) == 0 {
			return ValueSpec{Type: TypeList, Element: &ValueSpec{Type: TypeAny, Visibility: vis}, Visibility: vis}, nil
		}
		first, err := Infer(e.Args[0], env)
		if err != nil {
			return ValueSpec{}, err
		}
		for _, a := range e.Args[1:] {
			s, err := Infer(a, env)
			if err != nil {
				return ValueSpec{}, err
			}
			if !compatible(first.Type, s.Type) {
				return ValueSpec{}, fmt.Errorf("list element type mismatch %s/%s", first.Type, s.Type)
			}
			first.Visibility = maxVis(first.Visibility, s.Visibility)
		}
		return ValueSpec{Type: TypeList, Element: &first, Visibility: first.Visibility}, nil
	case Object:
		v := vis
		for _, a := range e.Fields {
			s, err := Infer(a, env)
			if err != nil {
				return ValueSpec{}, err
			}
			v = maxVis(v, s.Visibility)
		}
		return ValueSpec{Type: TypeObject, Visibility: v}, nil
	case Not:
		s, err := inferArgs(e, env, 1, 1)
		if err != nil {
			return ValueSpec{}, err
		}
		if s[0].Type != TypeBool {
			return ValueSpec{}, fmt.Errorf("not requires bool")
		}
		return ValueSpec{Type: TypeBool, Visibility: s[0].Visibility}, nil
	case And, Or:
		s, err := inferArgs(e, env, 2, -1)
		if err != nil {
			return ValueSpec{}, err
		}
		v := Public
		for _, x := range s {
			if x.Type != TypeBool {
				return ValueSpec{}, fmt.Errorf("%s requires bool", e.Kind)
			}
			v = maxVis(v, x.Visibility)
		}
		return ValueSpec{Type: TypeBool, Visibility: v}, nil
	case Equal:
		s, err := inferArgs(e, env, 2, 2)
		if err != nil {
			return ValueSpec{}, err
		}
		if !compatible(s[0].Type, s[1].Type) {
			return ValueSpec{}, fmt.Errorf("equal type mismatch %s/%s", s[0].Type, s[1].Type)
		}
		return ValueSpec{Type: TypeBool, Visibility: maxVis(s[0].Visibility, s[1].Visibility)}, nil
	case Count:
		s, err := inferArgs(e, env, 1, 1)
		if err != nil {
			return ValueSpec{}, err
		}
		if s[0].Type != TypeList && s[0].Type != TypeString && s[0].Type != TypeObject {
			return ValueSpec{}, fmt.Errorf("count requires collection")
		}
		return ValueSpec{Type: TypeNumber, Visibility: s[0].Visibility}, nil
	case Sum:
		s, err := inferArgs(e, env, 1, 1)
		if err != nil {
			return ValueSpec{}, err
		}
		if s[0].Type != TypeList || s[0].Element == nil || s[0].Element.Type != TypeNumber {
			return ValueSpec{}, fmt.Errorf("sum requires list<number>")
		}
		return ValueSpec{Type: TypeNumber, Visibility: s[0].Visibility}, nil
	case If:
		s, err := inferArgs(e, env, 3, 3)
		if err != nil {
			return ValueSpec{}, err
		}
		if s[0].Type != TypeBool {
			return ValueSpec{}, fmt.Errorf("if condition requires bool")
		}
		if !compatible(s[1].Type, s[2].Type) {
			return ValueSpec{}, fmt.Errorf("if branch type mismatch")
		}
		return ValueSpec{Type: s[1].Type, Element: s[1].Element, Visibility: maxVis(s[0].Visibility, maxVis(s[1].Visibility, s[2].Visibility))}, nil
	case HighestThreshold:
		s, err := inferArgs(e, env, 2, 2)
		if err != nil {
			return ValueSpec{}, err
		}
		if s[0].Type != TypeNumber || s[1].Type != TypeList || s[1].Element == nil || s[1].Element.Type != TypeNumber {
			return ValueSpec{}, fmt.Errorf("highest threshold requires number,list<number>")
		}
		return ValueSpec{Type: TypeNumber, Visibility: maxVis(s[0].Visibility, s[1].Visibility)}, nil
	default:
		return ValueSpec{}, fmt.Errorf("unknown expression kind %q", e.Kind)
	}
}
func Validate(e Expr) error {
	_, err := Infer(e, TypeEnv{map[string]ValueSpec{}, map[string]ValueSpec{}, map[string]ValueSpec{}, map[string]ValueSpec{}})
	return err
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
		var exists bool
		root, exists = m[p]
		if !exists {
			return nil, fmt.Errorf("path segment %q not found", p)
		}
	}
	return root, nil
}
func stateSpec(e Expr) (ValueSpec, bool) {
	if e.ValueType != "" {
		return ValueSpec{Type: e.ValueType, Visibility: visibilityOrPublic(e.Visibility)}, true
	}
	if len(e.Path) == 1 {
		switch e.Path[0] {
		case "match_id":
			return ValueSpec{Type: TypeString, Visibility: Public}, true
		case "revision", "event_sequence":
			return ValueSpec{Type: TypeNumber, Visibility: Public}, true
		case "lifecycle":
			return ValueSpec{Type: TypeString, Visibility: Public}, true
		}
	}
	return ValueSpec{}, false
}
func inferArgs(e Expr, env TypeEnv, min, max int) ([]ValueSpec, error) {
	if len(e.Args) < min || (max >= 0 && len(e.Args) > max) {
		return nil, fmt.Errorf("%s invalid arity", e.Kind)
	}
	r := make([]ValueSpec, len(e.Args))
	for i, a := range e.Args {
		s, err := Infer(a, env)
		if err != nil {
			return nil, err
		}
		r[i] = s
	}
	return r, nil
}
func literalSpec(v any, vis Visibility) ValueSpec {
	t := TypeAny
	switch v.(type) {
	case bool:
		t = TypeBool
	case string:
		t = TypeString
	case int, int32, int64, float32, float64, json.Number:
		t = TypeNumber
	case []any:
		t = TypeList
	case map[string]any:
		t = TypeObject
	case nil:
		t = TypeAny
	}
	return ValueSpec{Type: t, Visibility: vis}
}
func compatible(a, b Type) bool {
	return a == b || a == TypeAny || b == TypeAny || (a == TypeString && (b == TypePlayer || b == TypeFighter)) || (b == TypeString && (a == TypePlayer || a == TypeFighter))
}
func maxVis(a, b Visibility) Visibility {
	if visRank(a) >= visRank(b) {
		return a
	}
	return b
}
func visRank(v Visibility) int {
	switch v {
	case Opaque:
		return 2
	case OwnerPrivate:
		return 1
	default:
		return 0
	}
}
func visibilityOrPublic(v Visibility) Visibility {
	if v == "" {
		return Public
	}
	return v
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
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
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
