package game

import (
	"reflect"
	"testing"
)

type maneuverProjectionState struct {
	Revision         uint64
	CurrentPlayer    string
	ActionsRemaining int
	Deck             []CardInstance
	Hand             []CardInstance
	Discard          []CardInstance
	PendingCount     int
	EventCount       int
	Fighters         map[string]Fighter
}

func captureManeuverProjectionState(match *Match, playerID string) maneuverProjectionState {
	player := match.Players[playerID]
	fighters := make(map[string]Fighter, len(match.Fighters))
	for id, fighter := range match.Fighters {
		fighters[id] = *fighter
	}
	return maneuverProjectionState{
		Revision:         match.Revision,
		CurrentPlayer:    match.CurrentPlayer,
		ActionsRemaining: player.ActionsRemaining,
		Deck:             append([]CardInstance(nil), player.Deck...),
		Hand:             append([]CardInstance(nil), player.Hand...),
		Discard:          append([]CardInstance(nil), player.Discard...),
		PendingCount:     len(match.Pending),
		EventCount:       len(match.Events),
		Fighters:         fighters,
	}
}

func TestManeuverLegalProjectionCarriesBaseMovementAndAuthoritativePaths(t *testing.T) {
	match, p1, _ := activeMatch(t, "maneuver-projection")
	defeated := firstFighter(match, p1, "outlaw")
	if defeated == nil {
		t.Fatal("missing Outlaw fixture")
	}
	defeated.Defeated = true
	defeated.SpaceID = ""

	before := captureManeuverProjectionState(match, p1)
	view, err := match.Project(p1)
	if err != nil {
		t.Fatal(err)
	}
	action := view.Legal.ManeuverAction
	if !view.Legal.CanManeuver || action == nil {
		t.Fatalf("missing Maneuver legal action: %+v", view.Legal)
	}
	wantMovement := match.Registry.Decks[match.Players[p1].DeckDefinitionID].Movement
	if action.BaseMovement != wantMovement {
		t.Fatalf("base movement=%d want %d", action.BaseMovement, wantMovement)
	}

	living := 0
	projectedOptions := 0
	for _, fighterID := range match.Players[p1].FighterIDs {
		fighter := match.Fighters[fighterID]
		options, projected := action.DestinationsByFighter[fighterID]
		if fighter == nil || fighter.Defeated {
			if projected {
				t.Fatalf("defeated or missing fighter %s received a movement domain", fighterID)
			}
			continue
		}
		living++
		if !projected {
			t.Fatalf("living fighter %s is missing from the movement domain", fighterID)
		}
		for _, option := range options {
			projectedOptions++
			if option.FighterID != fighterID {
				t.Fatalf("option fighter=%q want %q", option.FighterID, fighterID)
			}
			if option.Destination == "" || option.Destination == fighter.SpaceID {
				t.Fatalf("invalid projected destination for %s: %+v", fighterID, option)
			}
			if len(option.Path) == 0 || len(option.Path) > wantMovement {
				t.Fatalf("invalid projected path length for %s: %+v", fighterID, option)
			}
			if option.Path[len(option.Path)-1] != option.Destination {
				t.Fatalf("path %v does not terminate at %s", option.Path, option.Destination)
			}
			if err := match.validatePath(fighter, option.Path, wantMovement, false); err != nil {
				t.Fatalf("projected path is not authoritative for %s: %+v: %v", fighterID, option, err)
			}
		}
	}
	if len(action.DestinationsByFighter) != living {
		t.Fatalf("projected fighter domains=%d want %d", len(action.DestinationsByFighter), living)
	}
	if projectedOptions == 0 {
		t.Fatal("expected at least one base-movement destination")
	}

	after := captureManeuverProjectionState(match, p1)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("projection mutated match state\nbefore=%+v\nafter=%+v", before, after)
	}
}

func TestManeuverLegalProjectionIsPrivateAndLimitedToIdleActionWindow(t *testing.T) {
	match, p1, p2 := activeMatch(t, "maneuver-projection-window")

	opponentView, err := match.Project(p2)
	if err != nil {
		t.Fatal(err)
	}
	if opponentView.Legal.CanManeuver || opponentView.Legal.ManeuverAction != nil {
		t.Fatal("non-current player received private Maneuver domains")
	}

	mustApply(t, match, p1, Command{Type: CommandManeuver})
	pendingView, err := match.Project(p1)
	if err != nil {
		t.Fatal(err)
	}
	if pendingView.Pending == nil || pendingView.Pending.Kind != PromptManeuverBoost {
		t.Fatalf("expected Maneuver BOOST prompt, got %+v", pendingView.Pending)
	}
	if pendingView.Legal.CanManeuver || pendingView.Legal.ManeuverAction != nil {
		t.Fatal("Maneuver domains remained available while a choice was pending")
	}

	noActions, current, _ := activeMatch(t, "maneuver-projection-no-actions")
	noActions.Players[current].ActionsRemaining = 0
	noActionView, err := noActions.Project(current)
	if err != nil {
		t.Fatal(err)
	}
	if noActionView.Legal.CanManeuver || noActionView.Legal.ManeuverAction != nil {
		t.Fatal("Maneuver domains remained available with no actions remaining")
	}
}
