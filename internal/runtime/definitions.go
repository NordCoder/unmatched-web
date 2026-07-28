package runtime

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

// FighterDefinition is the minimum construction data Core needs in Wave 1.
// CardDefinitions contains one entry per physical card instance, so repeated
// definition IDs intentionally allocate distinct runtime CardIDs.
type FighterDefinition struct {
	ID              model.DefinitionID
	CardDefinitions []model.DefinitionID
}

// DefinitionBundle pins the immutable corpus identity and construction data
// used for the lifetime of a match.
type DefinitionBundle struct {
	Ref      model.DefinitionRef
	Fighters map[model.DefinitionID]FighterDefinition
}

type DefinitionRegistry interface {
	Resolve(key string) (DefinitionBundle, bool)
	ResolveRef(ref model.DefinitionRef) (DefinitionBundle, bool)
}

type MemoryDefinitionRegistry struct {
	mu      sync.RWMutex
	bundles map[string]DefinitionBundle
}

func NewMemoryDefinitionRegistry() *MemoryDefinitionRegistry {
	return &MemoryDefinitionRegistry{bundles: make(map[string]DefinitionBundle)}
}

func (r *MemoryDefinitionRegistry) Register(key string, bundle DefinitionBundle) error {
	if key == "" {
		return fmt.Errorf("definition key is required")
	}
	if bundle.Ref.RulesetVersion == "" || bundle.Ref.CapabilityRegistry == "" {
		return fmt.Errorf("definition identity requires ruleset and capability registry versions")
	}
	if len(bundle.Fighters) == 0 {
		return fmt.Errorf("definition bundle requires at least one fighter")
	}
	for id, fighter := range bundle.Fighters {
		if id == "" || fighter.ID != id {
			return fmt.Errorf("fighter definition key %q does not match definition ID %q", id, fighter.ID)
		}
		for _, cardID := range fighter.CardDefinitions {
			if cardID == "" {
				return fmt.Errorf("fighter %q contains an empty card definition ID", id)
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bundles[key]; exists {
		return fmt.Errorf("definition key %q is already registered", key)
	}
	r.bundles[key] = cloneDefinitionBundle(bundle)
	return nil
}

func (r *MemoryDefinitionRegistry) Resolve(key string) (DefinitionBundle, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	bundle, ok := r.bundles[key]
	if !ok {
		return DefinitionBundle{}, false
	}
	return cloneDefinitionBundle(bundle), true
}

func (r *MemoryDefinitionRegistry) ResolveRef(ref model.DefinitionRef) (DefinitionBundle, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, bundle := range r.bundles {
		if reflect.DeepEqual(bundle.Ref, ref) {
			return cloneDefinitionBundle(bundle), true
		}
	}
	return DefinitionBundle{}, false
}

func cloneDefinitionBundle(bundle DefinitionBundle) DefinitionBundle {
	cloned := DefinitionBundle{
		Ref:      cloneDefinitionRef(bundle.Ref),
		Fighters: make(map[model.DefinitionID]FighterDefinition, len(bundle.Fighters)),
	}
	for id, fighter := range bundle.Fighters {
		cards := append([]model.DefinitionID(nil), fighter.CardDefinitions...)
		cloned.Fighters[id] = FighterDefinition{ID: fighter.ID, CardDefinitions: cards}
	}
	return cloned
}

func cloneDefinitionRef(ref model.DefinitionRef) model.DefinitionRef {
	ref.FighterManifestDigests = cloneStringMap(ref.FighterManifestDigests)
	ref.CardManifestDigests = cloneStringMap(ref.CardManifestDigests)
	return ref
}

func cloneStringMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
