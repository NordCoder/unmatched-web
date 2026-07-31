package game

import "testing"

func TestJackalopeHornsMovesWithoutTargetWhenSelectorIsEmpty(t *testing.T) {
	match, p1, p2 := activeMatch(t, "horns-empty")
	jackalope, _ := prepareHornsPosition(t, match, p1, p2, false)
	card := giveCard(t, match, p2, "jackalope-horns")

	healthBefore := make(map[string]int)
	for fighterID, fighter := range match.Fighters {
		healthBefore[fighterID] = fighter.Health
	}
	revision := match.Revision
	actions := match.Players[p2].ActionsRemaining
	hand := len(match.Players[p2].Hand)
	discard := len(match.Players[p2].Discard)

	err := match.Apply(p2, Command{
		Type: CommandScheme, ExpectedRevision: revision,
		FighterID: jackalope.ID, CardID: card.ID, Path: []string{"s02"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if jackalope.SpaceID != "s02" {
		t.Fatalf("Jackalope space=%s want s02", jackalope.SpaceID)
	}
	if match.Revision != revision+1 || match.Players[p2].ActionsRemaining != actions-1 {
		t.Fatalf("successful Horns did not consume exactly one command/action")
	}
	if len(match.Players[p2].Hand) != hand-1 || len(match.Players[p2].Discard) != discard+1 {
		t.Fatal("successful Horns did not move its card from hand to discard")
	}
	for fighterID, want := range healthBefore {
		if got := match.Fighters[fighterID].Health; got != want {
			t.Fatalf("%s health=%d want %d when no Horns target exists", fighterID, got, want)
		}
	}
}

func TestJackalopeHornsRequiresValidTargetAtomicallyAndAllowsFriendlyTarget(t *testing.T) {
	match, p1, p2 := activeMatch(t, "horns-target")
	jackalope, bigfoot := prepareHornsPosition(t, match, p1, p2, true)
	card := giveCard(t, match, p2, "jackalope-horns")
	robin := match.Fighters[match.Players[p1].HeroID]

	revision := match.Revision
	actions := match.Players[p2].ActionsRemaining
	hand := len(match.Players[p2].Hand)
	discard := len(match.Players[p2].Discard)
	jackalopeSpace := jackalope.SpaceID
	bigfootHealth := bigfoot.Health

	assertRejected := func(targetID string) {
		t.Helper()
		err := match.Apply(p2, Command{
			Type: CommandScheme, ExpectedRevision: revision,
			FighterID: jackalope.ID, CardID: card.ID, Path: []string{"s02"}, TargetID: targetID,
		})
		if ErrorCode(err) != "illegal_target" {
			t.Fatalf("target %q error=%v", targetID, err)
		}
		if match.Revision != revision || match.Players[p2].ActionsRemaining != actions {
			t.Fatal("rejected Horns changed revision or actions")
		}
		if len(match.Players[p2].Hand) != hand || len(match.Players[p2].Discard) != discard {
			t.Fatal("rejected Horns consumed or discarded its card")
		}
		if jackalope.SpaceID != jackalopeSpace || bigfoot.Health != bigfootHealth {
			t.Fatal("rejected Horns moved Jackalope or dealt damage")
		}
	}
	assertRejected("")
	assertRejected(robin.ID)

	err := match.Apply(p2, Command{
		Type: CommandScheme, ExpectedRevision: revision,
		FighterID: jackalope.ID, CardID: card.ID, Path: []string{"s02"}, TargetID: bigfoot.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if jackalope.SpaceID != "s02" || bigfoot.Health != bigfootHealth-2 {
		t.Fatalf("friendly Horns target did not resolve: jackalope=%s bigfoot=%d", jackalope.SpaceID, bigfoot.Health)
	}
}

func prepareHornsPosition(t *testing.T, match *Match, p1, p2 string, adjacentFriendly bool) (*Fighter, *Fighter) {
	t.Helper()
	jackalope := firstFighter(match, p2, "jackalope")
	bigfoot := match.Fighters[match.Players[p2].HeroID]
	jackalope.SpaceID = "s01"
	if adjacentFriendly {
		bigfoot.SpaceID = "s03"
	} else {
		bigfoot.SpaceID = "s10"
	}

	remote := []string{"s12", "s13", "s14", "s15", "s16"}
	cursor := 0
	for _, fighterID := range match.Players[p1].FighterIDs {
		match.Fighters[fighterID].SpaceID = remote[cursor]
		cursor++
	}
	match.CurrentPlayer = p2
	match.Players[p1].ActionsRemaining = 0
	match.Players[p2].ActionsRemaining = 2
	match.Pending = nil
	match.endTurnQueued = false
	match.captureTurnStart(p2)
	return jackalope, bigfoot
}
