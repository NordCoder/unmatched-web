package game

import (
	"reflect"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func TestSetupPlacesSidekicksInHeroZonesDeterministically(t *testing.T) {
	first, _, _ := activeMatch(t, "legal-setup")
	second, _, _ := activeMatch(t, "legal-setup")

	firstPlacements := fighterPlacements(first)
	if !reflect.DeepEqual(firstPlacements, fighterPlacements(second)) {
		t.Fatalf("setup placements are not deterministic:\nfirst=%v\nsecond=%v", firstPlacements, fighterPlacements(second))
	}

	starts := make(map[int]string)
	for _, space := range first.Registry.Battlefield.Spaces {
		if space.StartingSeat != 0 {
			starts[space.StartingSeat] = space.ID
		}
	}
	occupied := make(map[string]string)
	for _, playerID := range first.PlayerOrder {
		player := first.Players[playerID]
		hero := first.Fighters[player.HeroID]
		if hero.SpaceID != starts[player.Seat] {
			t.Fatalf("%s hero starts at %s, want %s", player.Name, hero.SpaceID, starts[player.Seat])
		}
		heroSpace, _ := first.space(hero.SpaceID)
		deck := first.Registry.Decks[player.DeckDefinitionID]
		for _, fighterID := range player.FighterIDs {
			fighter := first.Fighters[fighterID]
			if previous, duplicate := occupied[fighter.SpaceID]; duplicate {
				t.Fatalf("%s and %s both occupy %s", previous, fighterID, fighter.SpaceID)
			}
			occupied[fighter.SpaceID] = fighterID
			if deck.Fighters[fighter.DefinitionID].Role == "hero" {
				continue
			}
			sidekickSpace, _ := first.space(fighter.SpaceID)
			if !sharesZone(heroSpace, sidekickSpace) {
				t.Fatalf("%s starts at %s outside %s's zone at %s", fighterID, fighter.SpaceID, player.HeroID, hero.SpaceID)
			}
		}
	}
}

func TestSetupFailureIsAtomic(t *testing.T) {
	registry, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	candidate := *registry
	candidate.Battlefield = registry.Battlefield
	candidate.Battlefield.Spaces = append([]content.SpaceDefinition(nil), registry.Battlefield.Spaces...)
	for i := range candidate.Battlefield.Spaces {
		if candidate.Battlefield.Spaces[i].ID == "s20" {
			candidate.Battlefield.Spaces[i].Zones = []string{"gray"}
		} else {
			candidate.Battlefield.Spaces[i].Zones = []string{"orange"}
		}
	}

	match, playerID, err := New("atomic-setup", &candidate, "robin-hood")
	if err != nil {
		t.Fatal(err)
	}
	seedBefore := match.seed
	eventsBefore := append([]Event(nil), match.Events...)

	if _, err := match.Join("bigfoot"); err == nil {
		t.Fatal("setup unexpectedly succeeded without legal Robin Hood sidekick spaces")
	}
	player := match.Players[playerID]
	if match.Phase != PhaseWaiting || match.Revision != 0 || match.CurrentPlayer != "" {
		t.Fatalf("failed setup changed match lifecycle: phase=%s revision=%d current=%q", match.Phase, match.Revision, match.CurrentPlayer)
	}
	if len(match.PlayerOrder) != 1 || len(match.Players) != 1 || len(match.Fighters) != 0 {
		t.Fatalf("failed setup left partial participants or fighters: order=%v players=%d fighters=%d", match.PlayerOrder, len(match.Players), len(match.Fighters))
	}
	if len(player.Deck) != 0 || len(player.Hand) != 0 || len(player.Discard) != 0 || len(player.FighterIDs) != 0 || player.HeroID != "" {
		t.Fatalf("failed setup mutated seat one: %+v", player)
	}
	if match.seed != seedBefore || match.cardSequence != 0 {
		t.Fatalf("failed setup consumed deterministic state: seed=%d want=%d cards=%d", match.seed, seedBefore, match.cardSequence)
	}
	if !reflect.DeepEqual(match.Events, eventsBefore) {
		t.Fatalf("failed setup changed events: before=%v after=%v", eventsBefore, match.Events)
	}
}

func fighterPlacements(match *Match) map[string]string {
	result := make(map[string]string, len(match.Fighters))
	for fighterID, fighter := range match.Fighters {
		result[fighterID] = fighter.SpaceID
	}
	return result
}
