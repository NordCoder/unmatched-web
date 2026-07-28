// Package capabilities records the generic Wave 1 capability dependency graph.
package capabilities

import (
	"fmt"
	"sort"
)

type Capability struct {
	ID, Name  string
	DependsOn []string
}

func Wave1() map[string]Capability {
	return map[string]Capability{"CAP-001": {"CAP-001", "staged resolution procedures", nil}, "CAP-002": {"CAP-002", "interactions and visibility", []string{"CAP-001"}}, "CAP-003": {"CAP-003", "reaction and cancellation windows", []string{"CAP-001", "CAP-002"}}, "CAP-004": {"CAP-004", "history and provenance ledger", nil}, "CAP-018": {"CAP-018", "derived query expressions", nil}}
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
