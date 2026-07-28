package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	mu        sync.RWMutex
	bundles   map[string]DefinitionBundle
	keysByRef map[string]string
}

func NewMemoryDefinitionRegistry() *MemoryDefinitionRegistry {
	return &MemoryDefinitionRegistry{
		bundles:   make(map[string]DefinitionBundle),
		keysByRef: make(map[string]string),
	}
}

// FighterDefinitionDigest returns the deterministic identity of the exact
// construction data Core consumes, including repeated card copies and order.
func FighterDefinitionDigest(fighter FighterDefinition) string {
	encoded, err := json.Marshal(struct {
		ID    model.DefinitionID   `json:"definition_id"`
		Cards []model.DefinitionID `json:"card_definitions"`
	}{ID: fighter.ID, Cards: fighter.CardDefinitions})
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (r *MemoryDefinitionRegistry) Register(key string, bundle DefinitionBundle) error {
	if err := validateDefinitionBundle(key, bundle); err != nil {
		return err
	}
	refKey, err := definitionRefKey(bundle.Ref)
	if err != nil {
		return fmt.Errorf("encode definition reference: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.bundles[key]; exists {
		return fmt.Errorf("definition key %q is already registered", key)
	}
	if existingKey, exists := r.keysByRef[refKey]; exists {
		return fmt.Errorf("definition reference is already registered by key %q", existingKey)
	}
	r.bundles[key] = cloneDefinitionBundle(bundle)
	r.keysByRef[refKey] = key
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
	refKey, err := definitionRefKey(ref)
	if err != nil {
		return DefinitionBundle{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	key, ok := r.keysByRef[refKey]
	if !ok {
		return DefinitionBundle{}, false
	}
	return cloneDefinitionBundle(r.bundles[key]), true
}

func validateDefinitionBundle(key string, bundle DefinitionBundle) error {
	if key == "" {
		return fmt.Errorf("definition key is required")
	}
	if bundle.Ref.RulesetVersion == "" || bundle.Ref.CapabilityRegistry == "" {
		return fmt.Errorf("definition identity requires ruleset and capability registry versions")
	}
	if len(bundle.Fighters) == 0 {
		return fmt.Errorf("definition bundle requires at least one fighter")
	}
	if len(bundle.Ref.FighterManifestDigests) != len(bundle.Fighters) {
		return fmt.Errorf("fighter manifest digests must exactly cover registered fighters")
	}
	usedCards := make(map[model.DefinitionID]struct{})
	for id, fighter := range bundle.Fighters {
		if id == "" || fighter.ID != id {
			return fmt.Errorf("fighter definition key %q does not match definition ID %q", id, fighter.ID)
		}
		expectedDigest := FighterDefinitionDigest(fighter)
		if bundle.Ref.FighterManifestDigests[string(id)] != expectedDigest {
			return fmt.Errorf("fighter %q construction data does not match its pinned manifest digest", id)
		}
		for _, cardID := range fighter.CardDefinitions {
			if cardID == "" {
				return fmt.Errorf("fighter %q contains an empty card definition ID", id)
			}
			usedCards[cardID] = struct{}{}
		}
	}
	if len(bundle.Ref.CardManifestDigests) != len(usedCards) {
		return fmt.Errorf("card manifest digests must exactly cover referenced card definitions")
	}
	for cardID := range usedCards {
		if bundle.Ref.CardManifestDigests[string(cardID)] == "" {
			return fmt.Errorf("card definition %q is missing a pinned manifest digest", cardID)
		}
	}
	return nil
}

func definitionRefKey(ref model.DefinitionRef) (string, error) {
	encoded, err := json.Marshal(cloneDefinitionRef(ref))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
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
