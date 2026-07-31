package content

import "testing"

func TestLoadLaunchRegistry(t *testing.T) {
	registry, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Decks) != 2 {
		t.Fatalf("decks=%d", len(registry.Decks))
	}
	if got := len(registry.Battlefield.Spaces); got != 30 {
		t.Fatalf("spaces=%d", got)
	}
	if got := len(registry.Battlefield.Zones); got != 7 {
		t.Fatalf("zones=%d", got)
	}
	if got := len(registry.Battlefield.Edges); got != 39 {
		t.Fatalf("edges=%d", got)
	}
	starts := map[string]int{}
	for _, space := range registry.Battlefield.Spaces {
		if space.StartingSeat != 0 {
			starts[space.ID] = space.StartingSeat
		}
	}
	if starts["s20"] != 1 || starts["s19"] != 2 || len(starts) != 2 {
		t.Fatalf("starts=%v", starts)
	}
	for id, deck := range registry.Decks {
		total := 0
		for _, card := range deck.Cards {
			total += card.Quantity
		}
		if total != 30 {
			t.Fatalf("deck %s total=%d", id, total)
		}
	}
}
