// Package capabilities records the generic Wave 1 capability dependency graph.
package capabilities

import (
	"fmt"
	"sort"
)

type Capability struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	DependsOn []string `json:"depends_on,omitempty"`
}

func Wave1() map[string]Capability {
	return map[string]Capability{
		"CAP-001": {ID: "CAP-001", Name: "staged resolution procedures"},
		"CAP-002": {ID: "CAP-002", Name: "interactions and visibility", DependsOn: []string{"CAP-001"}},
		"CAP-003": {ID: "CAP-003", Name: "reaction and cancellation windows", DependsOn: []string{"CAP-001", "CAP-002"}},
		"CAP-004": {ID: "CAP-004", Name: "history and provenance ledger"},
		"CAP-018": {ID: "CAP-018", Name: "derived query expressions"},
	}
}
func Validate(m map[string]Capability) error {
	seen, active := map[string]bool{}, map[string]bool{}
	var visit func(string) error
	visit = func(id string) error {
		if active[id] {
			return fmt.Errorf("capability cycle at %s", id)
		}
		if seen[id] {
			return nil
		}
		c, ok := m[id]
		if !ok {
			return fmt.Errorf("unknown capability %s", id)
		}
		active[id] = true
		for _, d := range c.DependsOn {
			if err := visit(d); err != nil {
				return err
			}
		}
		delete(active, id)
		seen[id] = true
		return nil
	}
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
func ValidateDeclaration(ids []string, m map[string]Capability) error {
	if len(ids) == 0 {
		return fmt.Errorf("at least one capability is required")
	}
	decl := map[string]bool{}
	for _, id := range ids {
		if _, ok := m[id]; !ok {
			return fmt.Errorf("unknown capability %s", id)
		}
		decl[id] = true
	}
	for id := range decl {
		for _, dep := range m[id].DependsOn {
			if !decl[dep] {
				return fmt.Errorf("capability %s requires %s", id, dep)
			}
		}
	}
	return nil
}
