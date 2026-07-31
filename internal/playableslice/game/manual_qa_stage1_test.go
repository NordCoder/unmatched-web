package game

import (
	"reflect"
	"strings"
	"testing"
)

func TestManeuverDrawsExactlyOneCardAcrossMultipleFighterMoves(t *testing.T) {
	match, playerID, _ := activeMatch(t, "maneuver-single-draw")
	resetCards(match, playerID)
	player := match.Players[playerID]
	handBefore := len(player.Hand)
	deckBefore := len(player.Deck)
	eventsBefore := len(match.Events)

	mustApply(t, match, playerID, Command{Type: CommandManeuver})
	if len(player.Hand) != handBefore+1 || len(player.Deck) != deckBefore-1 {
		t.Fatalf("maneuver draw changed hand/deck by more than one: hand=%d want=%d deck=%d want=%d", len(player.Hand), handBefore+1, len(player.Deck), deckBefore-1)
	}

	maneuverEvents := 0
	for _, event := range match.Events[eventsBefore:] {
		if event.Type != "maneuver_started" {
			continue
		}
		maneuverEvents++
		if event.Data["player_id"] != playerID || !strings.Contains(event.Message, "drew a card") {
			t.Fatalf("maneuver draw audit is incomplete: %+v", event)
		}
	}
	if maneuverEvents != 1 {
		t.Fatalf("maneuver emitted %d start/draw events, want 1", maneuverEvents)
	}

	mustChoose(t, match, playerID, "boost:none")
	moved := map[string]bool{}
	for len(moved) < 2 {
		option := findPromptOption(t, match, func(option Option) bool {
			return option.FighterID != "" && len(option.Path) > 0 && !moved[option.FighterID]
		})
		mustChoose(t, match, playerID, option.ID)
		moved[option.FighterID] = true
		if len(player.Hand) != handBefore+1 || len(player.Deck) != deckBefore-1 {
			t.Fatalf("moving fighters drew extra cards: moved=%d hand=%d deck=%d", len(moved), len(player.Hand), len(player.Deck))
		}
	}
	mustChoose(t, match, playerID, "maneuver:done")
	if len(player.Hand) != handBefore+1 || len(player.Deck) != deckBefore-1 {
		t.Fatalf("finishing maneuver drew extra cards: hand=%d deck=%d", len(player.Hand), len(player.Deck))
	}
}

func TestRobinRetreatMayMoveThroughFriendlyButNotOpposingFighters(t *testing.T) {
	t.Run("friendly fighter", func(t *testing.T) {
		match, robinPlayer, bigfootPlayer := activeMatch(t, "retreat-friendly")
		robin := match.Fighters[match.Players[robinPlayer].HeroID]
		outlaw := firstFighter(match, robinPlayer, "outlaw")
		bigfoot := match.Fighters[match.Players[bigfootPlayer].HeroID]
		stageOnePlaceFighters(match, robinPlayer, bigfootPlayer, map[string]string{
			robin.ID:   "s20",
			outlaw.ID:  "s21",
			bigfoot.ID: "s22",
		})

		stageOneResolveRobinAttack(t, match, robinPlayer, bigfootPlayer, bigfoot)
		option := stageOneFindDestination(match.Pending[0].Options, "s27")
		if option == nil || !reflect.DeepEqual(option.Path, []string{"s21", "s27"}) {
			t.Fatalf("retreat did not traverse friendly fighter: options=%v", match.Pending[0].Options)
		}
	})

	t.Run("opposing fighter", func(t *testing.T) {
		match, robinPlayer, bigfootPlayer := activeMatch(t, "retreat-opposing")
		robin := match.Fighters[match.Players[robinPlayer].HeroID]
		bigfoot := match.Fighters[match.Players[bigfootPlayer].HeroID]
		jackalope := firstFighter(match, bigfootPlayer, "jackalope")
		stageOnePlaceFighters(match, robinPlayer, bigfootPlayer, map[string]string{
			robin.ID:     "s20",
			bigfoot.ID:   "s21",
			jackalope.ID: "s22",
		})

		stageOneResolveRobinAttack(t, match, robinPlayer, bigfootPlayer, jackalope)
		if option := stageOneFindDestination(match.Pending[0].Options, "s27"); option != nil {
			t.Fatalf("retreat illegally traversed opposing fighter: %+v", *option)
		}
	})
}

func TestRobinRetreatAllowsZeroMovement(t *testing.T) {
	match, robinPlayer, bigfootPlayer := activeMatch(t, "retreat-zero")
	robin := match.Fighters[match.Players[robinPlayer].HeroID]
	bigfoot := match.Fighters[match.Players[bigfootPlayer].HeroID]
	stageOnePlaceFighters(match, robinPlayer, bigfootPlayer, map[string]string{
		robin.ID:   "s20",
		bigfoot.ID: "s22",
	})

	stageOneResolveRobinAttack(t, match, robinPlayer, bigfootPlayer, bigfoot)
	stay := stageOneFindDestination(match.Pending[0].Options, robin.SpaceID)
	if stay == nil || len(stay.Path) != 0 {
		t.Fatalf("retreat omitted the zero-space option: options=%v", match.Pending[0].Options)
	}
	before := robin.SpaceID
	mustChoose(t, match, robinPlayer, stay.ID)
	if robin.SpaceID != before {
		t.Fatalf("zero-space retreat moved Robin Hood from %s to %s", before, robin.SpaceID)
	}
}

func TestRangedAttackUsesEveryZoneOnMultiZoneSpace(t *testing.T) {
	match, robinPlayer, bigfootPlayer := activeMatch(t, "multi-zone-range")
	robin := match.Fighters[match.Players[robinPlayer].HeroID]
	bigfoot := match.Fighters[match.Players[bigfootPlayer].HeroID]
	stageOnePlaceFighters(match, robinPlayer, bigfootPlayer, map[string]string{robin.ID: "s06"})

	cases := []struct {
		name    string
		spaceID string
		legal   bool
	}{
		{name: "orange", spaceID: "s02", legal: true},
		{name: "green", spaceID: "s23", legal: true},
		{name: "light gray", spaceID: "s10", legal: true},
		{name: "unshared yellow", spaceID: "s11", legal: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			bigfoot.SpaceID = testCase.spaceID
			got := stageOneContains(match.attackTargets(robin), bigfoot.ID)
			if got != testCase.legal {
				t.Fatalf("ranged target at %s legal=%v want=%v", testCase.spaceID, got, testCase.legal)
			}
		})
	}
}

func stageOneResolveRobinAttack(t *testing.T, match *Match, robinPlayer, bigfootPlayer string, target *Fighter) {
	t.Helper()
	robin := match.Fighters[match.Players[robinPlayer].HeroID]
	attack := giveCard(t, match, robinPlayer, "a-hunters-eye")
	mustApply(t, match, robinPlayer, Command{Type: CommandAttack, FighterID: robin.ID, TargetID: target.ID, CardID: attack.ID})
	mustChoose(t, match, bigfootPlayer, "none")
	if len(match.Pending) == 0 || match.Pending[0].Kind != PromptMove {
		t.Fatalf("Robin Hood retreat prompt missing: %+v", match.Pending)
	}
}

func stageOnePlaceFighters(match *Match, robinPlayer, bigfootPlayer string, overrides map[string]string) {
	fallback := []string{"s01", "s04", "s05", "s08", "s12", "s14", "s16"}
	fighterIDs := append([]string(nil), match.Players[robinPlayer].FighterIDs...)
	fighterIDs = append(fighterIDs, match.Players[bigfootPlayer].FighterIDs...)
	for index, fighterID := range fighterIDs {
		match.Fighters[fighterID].SpaceID = fallback[index]
	}
	for fighterID, spaceID := range overrides {
		match.Fighters[fighterID].SpaceID = spaceID
	}
	match.CurrentPlayer = robinPlayer
	match.Players[robinPlayer].ActionsRemaining = 2
	match.Players[bigfootPlayer].ActionsRemaining = 0
	match.Pending = nil
	match.endTurnQueued = false
	match.captureTurnStart(robinPlayer)
}

func stageOneFindDestination(options []Option, destination string) *Option {
	for index := range options {
		if options[index].Destination == destination {
			return &options[index]
		}
	}
	return nil
}

func stageOneContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
