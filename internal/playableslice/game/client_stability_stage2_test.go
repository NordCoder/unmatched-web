package game

import (
	"reflect"
	"sort"
	"testing"
)

func TestSchemeLegalProjectionCarriesAuthoritativePathsAndTargets(t *testing.T) {
	match, p1, p2 := activeMatch(t, "stage-2-scheme-legal")
	resetCards(match, p2)
	crash := giveCard(t, match, p2, "crash-through-the-trees")
	horns := giveCard(t, match, p2, "jackalope-horns")

	match.CurrentPlayer = p2
	match.Players[p1].ActionsRemaining = 0
	match.Players[p2].ActionsRemaining = 2

	bigfoot := match.Fighters[match.Players[p2].HeroID]
	jackalope := firstFighter(match, p2, "jackalope")
	robin := match.Fighters[match.Players[p1].HeroID]
	bigfoot.SpaceID = "s04"
	jackalope.SpaceID = "s02"
	robin.SpaceID = "s20"
	for _, fighterID := range match.Players[p1].FighterIDs {
		if fighterID != robin.ID {
			match.Fighters[fighterID].SpaceID = "s10"
		}
	}

	view, err := match.Project(p2)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Legal.SchemeCards) != 2 {
		t.Fatalf("scheme cards=%v", view.Legal.SchemeCards)
	}

	crashAction, ok := view.Legal.SchemeActions[crash.ID]
	if !ok || len(crashAction.Fighters) != 1 || crashAction.Fighters[0].FighterID != bigfoot.ID {
		t.Fatalf("unexpected Crash Through the Trees domain: %+v", crashAction)
	}
	if _, ok := findProjectedDestination(crashAction.Fighters[0].Destinations, "s03"); !ok {
		t.Fatal("Crash Through the Trees did not project an authoritative destination path")
	}

	hornsAction, ok := view.Legal.SchemeActions[horns.ID]
	if !ok || len(hornsAction.Fighters) != 1 || hornsAction.Fighters[0].FighterID != jackalope.ID {
		t.Fatalf("unexpected Jackalope Horns domain: %+v", hornsAction)
	}
	destination, ok := findProjectedDestination(hornsAction.Fighters[0].Destinations, "s03")
	if !ok {
		t.Fatal("Jackalope Horns did not project s03")
	}
	if !reflect.DeepEqual(destination.Path, []string{"s03"}) {
		t.Fatalf("path=%v want [s03]", destination.Path)
	}
	targets := append([]string(nil), hornsAction.Fighters[0].TargetsByDestination["s03"]...)
	sort.Strings(targets)
	wantTargets := []string{bigfoot.ID, robin.ID}
	sort.Strings(wantTargets)
	if !reflect.DeepEqual(targets, wantTargets) {
		t.Fatalf("targets=%v want %v", targets, wantTargets)
	}

	opponentView, err := match.Project(p1)
	if err != nil {
		t.Fatal(err)
	}
	if len(opponentView.Legal.SchemeCards) != 0 || len(opponentView.Legal.SchemeActions) != 0 {
		t.Fatal("non-current player received private legal scheme domains")
	}
}

func TestSchemeLegalProjectionIncludesZeroSpaceMovement(t *testing.T) {
	match, _, p2 := activeMatch(t, "stage-2-zero-move")
	resetCards(match, p2)
	horns := giveCard(t, match, p2, "jackalope-horns")
	match.CurrentPlayer = p2
	match.Players[p2].ActionsRemaining = 2

	view, err := match.Project(p2)
	if err != nil {
		t.Fatal(err)
	}
	action := view.Legal.SchemeActions[horns.ID]
	if len(action.Fighters) != 1 {
		t.Fatalf("fighters=%v", action.Fighters)
	}
	fighter := match.Fighters[action.Fighters[0].FighterID]
	stay, ok := findProjectedDestination(action.Fighters[0].Destinations, fighter.SpaceID)
	if !ok {
		t.Fatal("server legal domain omitted zero-space movement")
	}
	if len(stay.Path) != 0 {
		t.Fatalf("zero-space path=%v", stay.Path)
	}
}

func findProjectedDestination(options []Option, destination string) (Option, bool) {
	for _, option := range options {
		if option.Destination == destination {
			return option, true
		}
	}
	return Option{}, false
}
