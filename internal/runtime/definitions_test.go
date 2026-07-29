package runtime

import (
	"strings"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

func TestDefinitionRegistryRejectsIdentityContentConflicts(t *testing.T) {
	registry := NewMemoryDefinitionRegistry()
	bundle := testBundle("rules/v1", []model.DefinitionID{"a", "a", "b"})
	if err := registry.Register("one", bundle); err != nil {
		t.Fatalf("register valid bundle: %v", err)
	}

	conflicting := bundle
	conflicting.Fighters = map[model.DefinitionID]FighterDefinition{
		"fighter": {ID: "fighter", CardDefinitions: []model.DefinitionID{"a", "b"}},
	}
	if err := registry.Register("conflicting", conflicting); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("conflicting construction data accepted under pinned identity: %v", err)
	}

	if err := registry.Register("same-ref", bundle); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("duplicate definition reference accepted: %v", err)
	}
}

func TestDefinitionRegistryRequiresExactDigestCoverage(t *testing.T) {
	registry := NewMemoryDefinitionRegistry()
	bundle := testBundle("rules/v1", []model.DefinitionID{"a", "b"})
	delete(bundle.Ref.CardManifestDigests, "b")
	if err := registry.Register("missing-card", bundle); err == nil || !strings.Contains(err.Error(), "exactly cover") {
		t.Fatalf("missing card digest accepted: %v", err)
	}

	bundle = testBundle("rules/v1", []model.DefinitionID{"a", "b"})
	bundle.Ref.FighterManifestDigests["fighter"] = "sha256:wrong"
	if err := registry.Register("wrong-fighter", bundle); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("wrong fighter construction digest accepted: %v", err)
	}
}

func TestResolveRefUsesDeterministicOneToOneIndex(t *testing.T) {
	registry := NewMemoryDefinitionRegistry()
	first := testBundle("rules/v1", []model.DefinitionID{"a"})
	second := testBundle("rules/v2", []model.DefinitionID{"b"})
	if err := registry.Register("first", first); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register("second", second); err != nil {
		t.Fatal(err)
	}
	for iteration := 0; iteration < 100; iteration++ {
		resolved, ok := registry.ResolveRef(second.Ref)
		if !ok || resolved.Ref.RulesetVersion != "rules/v2" || resolved.Fighters["fighter"].CardDefinitions[0] != "b" {
			t.Fatalf("iteration %d resolved wrong bundle: %+v ok=%v", iteration, resolved, ok)
		}
	}
}

func testBundle(ruleset string, cards []model.DefinitionID) DefinitionBundle {
	fighter := FighterDefinition{ID: "fighter", CardDefinitions: append([]model.DefinitionID(nil), cards...)}
	cardDigests := make(map[string]string)
	for _, card := range cards {
		cardDigests[string(card)] = "sha256:" + string(card)
	}
	return DefinitionBundle{
		Ref: model.DefinitionRef{
			RulesetVersion: ruleset, CapabilityRegistry: "capabilities/v1",
			FighterManifestDigests: map[string]string{"fighter": FighterDefinitionDigest(fighter)},
			CardManifestDigests:    cardDigests,
		},
		Fighters: map[model.DefinitionID]FighterDefinition{"fighter": fighter},
	}
}
