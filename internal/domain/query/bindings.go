package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// ValidateClosedSpec rejects open-ended external data contracts. External
// captured/input bindings must have a complete recursive shape before runtime.
func ValidateClosedSpec(s ValueSpec) error {
	if s.Type == TypeAny {
		return fmt.Errorf("type any is not allowed at an external binding boundary")
	}
	switch s.Type {
	case TypeList:
		if s.Element == nil {
			return fmt.Errorf("list element spec is required")
		}
		if err := ValidateClosedSpec(*s.Element); err != nil {
			return fmt.Errorf("list element: %w", err)
		}
	case TypeObject, TypeOperationResult:
		for name, field := range s.Fields {
			if err := ValidateClosedSpec(field); err != nil {
				return fmt.Errorf("field %s: %w", name, err)
			}
		}
	case TypeBool, TypeNumber, TypeString, TypePlayer, TypeFighter, TypeDisposition:
	default:
		return fmt.Errorf("unknown type %s", s.Type)
	}
	return nil
}

// ValidateBindings enforces a closed recursive runtime contract and returns
// canonical independent copies suitable for evaluator use and serialization.
func ValidateBindings(actual map[string]json.RawMessage, specs map[string]ValueSpec) (map[string]json.RawMessage, error) {
	if actual == nil {
		actual = map[string]json.RawMessage{}
	}
	if specs == nil {
		specs = map[string]ValueSpec{}
	}
	for name := range actual {
		if _, ok := specs[name]; !ok {
			return nil, fmt.Errorf("unknown binding %s", name)
		}
	}
	names := make([]string, 0, len(specs))
	for name := range specs {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make(map[string]json.RawMessage, len(actual))
	for _, name := range names {
		spec := specs[name]
		raw, ok := actual[name]
		if !ok {
			if spec.Optional {
				continue
			}
			return nil, fmt.Errorf("missing binding %s", name)
		}
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, fmt.Errorf("binding %s: invalid JSON: %w", name, err)
		}
		if err := validateClosedValue(value, spec); err != nil {
			return nil, fmt.Errorf("binding %s: %w", name, err)
		}
		canonical, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("binding %s: canonicalize: %w", name, err)
		}
		out[name] = canonical
	}
	return out, nil
}

func validateClosedValue(v any, s ValueSpec) error {
	if v == nil {
		if s.Optional {
			return nil
		}
		return fmt.Errorf("nil does not match %s", s.Type)
	}
	switch s.Type {
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
		if s.Element == nil {
			return fmt.Errorf("list element spec is required")
		}
		for i, x := range xs {
			if err := validateClosedValue(x, *s.Element); err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}
		}
	case TypeObject, TypeOperationResult:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object")
		}
		for name := range m {
			if _, declared := s.Fields[name]; !declared {
				return fmt.Errorf("undeclared field %s", name)
			}
		}
		for name, field := range s.Fields {
			x, exists := m[name]
			if !exists {
				if field.Optional {
					continue
				}
				return fmt.Errorf("missing field %s", name)
			}
			if err := validateClosedValue(x, field); err != nil {
				return fmt.Errorf("field %s: %w", name, err)
			}
		}
	default:
		return fmt.Errorf("unsupported external binding type %s", s.Type)
	}
	return nil
}

func ExactSpec(a, b ValueSpec) bool { return reflect.DeepEqual(a, b) }
