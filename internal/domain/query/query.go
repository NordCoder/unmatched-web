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
	TypeDisposition     Type = "disposition"
)

type Visibility string

const (
	Public       Visibility = "public"
	OwnerPrivate Visibility = "owner_private"
	Opaque       Visibility = "opaque"
)

type ValueSpec struct {
	Type       Type                 `json:"type"`
	Element    *ValueSpec           `json:"element,omitempty"`
	Fields     map[string]ValueSpec `json:"fields,omitempty"`
	Visibility Visibility           `json:"visibility"`
	Optional   bool                 `json:"optional,omitempty"`
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

type StatePathSchema struct {
	Pattern []string
	Spec    ValueSpec
}

var statePaths = []StatePathSchema{
	{[]string{"match_id"}, ValueSpec{Type: TypeString, Visibility: Public}},
	{[]string{"revision"}, ValueSpec{Type: TypeNumber, Visibility: Public}},
	{[]string{"event_sequence"}, ValueSpec{Type: TypeNumber, Visibility: Public}},
	{[]string{"lifecycle"}, ValueSpec{Type: TypeString, Visibility: Public}},
	{[]string{"fighters", "*", "controller_player_id"}, ValueSpec{Type: TypePlayer, Visibility: Public}},
	{[]string{"fighters", "*", "state", "mode"}, ValueSpec{Type: TypeString, Visibility: Public}},
	{[]string{"fighters", "*", "state", "paid"}, ValueSpec{Type: TypeBool, Visibility: Public}},
	{[]string{"fighters", "*", "state", "marked"}, ValueSpec{Type: TypeBool, Visibility: Public}},
	{[]string{"players", "*", "private_zones"}, ValueSpec{Type: TypeObject, Fields: map[string]ValueSpec{}, Visibility: OwnerPrivate}},
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
		if !rv.IsValid() || (rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map && rv.Kind() != reflect.String) {
			return nil, fmt.Errorf("count requires collection")
		}
		return rv.Len(), nil
	case Sum:
		v, err := one(e, c)
		if err != nil {
			return nil, err
		}
		xs, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("sum requires list")
		}
		total := 0.0
		for _, x := range xs {
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
	assertVis := e.Visibility
	switch e.Kind {
	case Literal:
		s := literalSpec(e.Value, Public)
		if e.ValueType != "" {
			if !literalAssertion(s.Type, e.ValueType) {
				return ValueSpec{}, fmt.Errorf("literal declared %s but is %s", e.ValueType, s.Type)
			}
			s.Type = e.ValueType
		}
		if assertVis != "" {
			s.Visibility = maxVis(s.Visibility, assertVis)
		}
		return s, nil
	case Reference:
		if len(e.Path) == 0 {
			return ValueSpec{}, fmt.Errorf("reference requires path")
		}
		var s ValueSpec
		var ok bool
		rest := e.Path
		switch e.Source {
		case State:
			s, ok = stateSpec(e.Path)
			rest = nil
		case Captured:
			s, ok = env.Captured[e.Path[0]]
			rest = e.Path[1:]
		case Result:
			s, ok = env.Results[e.Path[0]]
			rest = e.Path[1:]
		case Choice:
			s, ok = env.Choices[e.Path[0]]
			rest = e.Path[1:]
		case Input:
			s, ok = env.Input[e.Path[0]]
			rest = e.Path[1:]
		default:
			return ValueSpec{}, fmt.Errorf("unknown source %q", e.Source)
		}
		if !ok {
			return ValueSpec{}, fmt.Errorf("unknown reference %s:%v", e.Source, e.Path)
		}
		var err error
		s, err = Traverse(s, rest)
		if err != nil {
			return ValueSpec{}, fmt.Errorf("reference %s:%v: %w", e.Source, e.Path, err)
		}
		if e.ValueType != "" && !exactType(s.Type, e.ValueType) {
			return ValueSpec{}, fmt.Errorf("reference declared %s but schema is %s", e.ValueType, s.Type)
		}
		if assertVis != "" {
			if visRank(assertVis) < visRank(s.Visibility) {
				return ValueSpec{}, fmt.Errorf("reference cannot downgrade %s to %s", s.Visibility, assertVis)
			}
			s.Visibility = maxVis(s.Visibility, assertVis)
		}
		return s, nil
	case List:
		if len(e.Args) == 0 {
			return ValueSpec{Type: TypeList, Element: &ValueSpec{Type: TypeAny, Visibility: visibilityOrPublic(assertVis)}, Visibility: visibilityOrPublic(assertVis)}, nil
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
			if !sameShape(first, s) {
				return ValueSpec{}, fmt.Errorf("list element type mismatch")
			}
			first.Visibility = maxVis(first.Visibility, s.Visibility)
		}
		v := maxVis(first.Visibility, visibilityOrPublic(assertVis))
		return ValueSpec{Type: TypeList, Element: &first, Visibility: v}, nil
	case Object:
		fields := map[string]ValueSpec{}
		v := visibilityOrPublic(assertVis)
		for _, k := range sorted(e.Fields) {
			s, err := Infer(e.Fields[k], env)
			if err != nil {
				return ValueSpec{}, err
			}
			fields[k] = s
			v = maxVis(v, s.Visibility)
		}
		return ValueSpec{Type: TypeObject, Fields: fields, Visibility: v}, nil
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
		if !sameShape(s[0], s[1]) && !scalarCompatible(s[0].Type, s[1].Type) {
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
		if !sameShape(s[1], s[2]) {
			return ValueSpec{}, fmt.Errorf("if branch type mismatch")
		}
		r := s[1]
		r.Visibility = maxVis(s[0].Visibility, maxVis(s[1].Visibility, s[2].Visibility))
		return r, nil
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

func Traverse(s ValueSpec, path []string) (ValueSpec, error) {
	for _, p := range path {
		if s.Fields == nil {
			return ValueSpec{}, fmt.Errorf("cannot traverse %s through %s", p, s.Type)
		}
		n, ok := s.Fields[p]
		if !ok {
			return ValueSpec{}, fmt.Errorf("field %q not in schema", p)
		}
		n.Visibility = maxVis(s.Visibility, n.Visibility)
		s = n
	}
	return s, nil
}
func ValidateValue(v any, s ValueSpec) error {
	if v == nil {
		if s.Optional || s.Type == TypeAny {
			return nil
		}
		return fmt.Errorf("nil does not match %s", s.Type)
	}
	switch s.Type {
	case TypeAny:
		return nil
	case TypeBool:
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("expected bool")
		}
	case TypeNumber:
		if _, ok := number(v); !ok {
			return fmt.Errorf("expected number")
		}
	case TypeString, TypePlayer, TypeFighter, TypeDisposition:
		if _, ok := v.(string); !ok {
			return fmt.Errorf("expected string-like %s", s.Type)
		}
	case TypeList:
		xs, ok := v.([]any)
		if !ok {
			return fmt.Errorf("expected list")
		}
		if s.Element != nil {
			for i, x := range xs {
				if err := ValidateValue(x, *s.Element); err != nil {
					return fmt.Errorf("element %d: %w", i, err)
				}
			}
		}
	case TypeObject, TypeOperationResult:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object")
		}
		for k, fs := range s.Fields {
			x, exists := m[k]
			if !exists {
				if fs.Optional {
					continue
				}
				return fmt.Errorf("missing field %s", k)
			}
			if err := ValidateValue(x, fs); err != nil {
				return fmt.Errorf("field %s: %w", k, err)
			}
		}
	default:
		return fmt.Errorf("unknown type %s", s.Type)
	}
	return nil
}
func Validate(e Expr) error {
	_, err := Infer(e, TypeEnv{Captured: map[string]ValueSpec{}, Results: map[string]ValueSpec{}, Choices: map[string]ValueSpec{}, Input: map[string]ValueSpec{}})
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
func ReferencesBinding(e Expr, s Source, name string) bool {
	for _, x := range References(e, s) {
		if x == name {
			return true
		}
	}
	return false
}
func OperationResultSpec(value ValueSpec, vis Visibility) ValueSpec {
	return ValueSpec{Type: TypeOperationResult, Visibility: vis, Fields: map[string]ValueSpec{"applied": {Type: TypeBool, Visibility: vis}, "disposition": {Type: TypeDisposition, Visibility: vis}, "code": {Type: TypeString, Visibility: vis, Optional: true}, "value": withOptional(value, true)}}
}

func stateSpec(path []string) (ValueSpec, bool) {
	for _, x := range statePaths {
		if match(x.Pattern, path) {
			return cloneSpec(x.Spec), true
		}
	}
	return ValueSpec{}, false
}
func match(pattern, path []string) bool {
	if len(pattern) != len(path) {
		return false
	}
	for i, p := range pattern {
		if p != "*" && p != path[i] {
			return false
		}
	}
	return true
}
func resolve(e Expr, c Context) (any, error) {
	var root any
	path := append([]string(nil), e.Path...)
	switch e.Source {
	case State:
		b, _ := json.Marshal(c.State)
		_ = json.Unmarshal(b, &root)
	case Captured:
		root = decode(c.Captured[path[0]])
		path = path[1:]
	case Result:
		root = decode(c.Results[path[0]])
		path = path[1:]
	case Choice:
		root = decode(c.Choices[path[0]])
		path = path[1:]
	case Input:
		root = decode(c.Input[path[0]])
		path = path[1:]
	default:
		return nil, fmt.Errorf("unknown source %q", e.Source)
	}
	for _, p := range path {
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
	switch x := v.(type) {
	case bool:
		return ValueSpec{Type: TypeBool, Visibility: vis}
	case string:
		return ValueSpec{Type: TypeString, Visibility: vis}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		return ValueSpec{Type: TypeNumber, Visibility: vis}
	case []any:
		if len(x) == 0 {
			return ValueSpec{Type: TypeList, Element: &ValueSpec{Type: TypeAny, Visibility: vis}, Visibility: vis}
		}
		e := literalSpec(x[0], vis)
		return ValueSpec{Type: TypeList, Element: &e, Visibility: vis}
	case map[string]any:
		f := map[string]ValueSpec{}
		v2 := vis
		for k, y := range x {
			s := literalSpec(y, vis)
			f[k] = s
			v2 = maxVis(v2, s.Visibility)
		}
		return ValueSpec{Type: TypeObject, Fields: f, Visibility: v2}
	default:
		return ValueSpec{Type: TypeAny, Visibility: vis}
	}
}
func sameShape(a, b ValueSpec) bool {
	if !scalarCompatible(a.Type, b.Type) {
		return false
	}
	if a.Type == TypeList {
		return a.Element != nil && b.Element != nil && sameShape(*a.Element, *b.Element)
	}
	if a.Type == TypeObject || a.Type == TypeOperationResult {
		if len(a.Fields) != len(b.Fields) {
			return false
		}
		for k, x := range a.Fields {
			y, ok := b.Fields[k]
			if !ok || !sameShape(x, y) {
				return false
			}
		}
	}
	return true
}
func scalarCompatible(a, b Type) bool { return a == b || a == TypeAny || b == TypeAny }
func exactType(a, b Type) bool        { return a == b }
func literalAssertion(actual, declared Type) bool {
	return actual == declared || (actual == TypeString && (declared == TypePlayer || declared == TypeFighter))
}
func visibilityOrPublic(v Visibility) Visibility {
	if v == "" {
		return Public
	}
	return v
}
func visRank(v Visibility) int {
	switch v {
	case OwnerPrivate:
		return 1
	case Opaque:
		return 2
	default:
		return 0
	}
}
func maxVis(a, b Visibility) Visibility {
	if visRank(b) > visRank(a) {
		return b
	}
	return a
}
func cloneSpec(s ValueSpec) ValueSpec {
	r := s
	if s.Element != nil {
		x := cloneSpec(*s.Element)
		r.Element = &x
	}
	if s.Fields != nil {
		r.Fields = map[string]ValueSpec{}
		for k, v := range s.Fields {
			r.Fields[k] = cloneSpec(v)
		}
	}
	return r
}
func withOptional(s ValueSpec, o bool) ValueSpec { s.Optional = o; return s }
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
	switch x := v.(type) {
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint64:
		return float64(x), true
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case json.Number:
		n, err := x.Float64()
		return n, err == nil
	}
	return 0, false
}
func decode(b json.RawMessage) any {
	if b == nil {
		return nil
	}
	var v any
	_ = json.Unmarshal(b, &v)
	return v
}
func sorted[V any](m map[string]V) []string {
	xs := make([]string, 0, len(m))
	for k := range m {
		xs = append(xs, k)
	}
	sort.Strings(xs)
	return xs
}

func MaxVisibility(a, b Visibility) Visibility { return maxVis(a, b) }
func VisibilityRank(v Visibility) int          { return visRank(v) }
func SameShape(a, b ValueSpec) bool            { return sameShape(a, b) }
