package game

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func activeMatch(t *testing.T, id string) (*Match, string, string) {
	t.Helper()
	registry, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	match, p1, err := New(id, registry, "robin-hood")
	if err != nil {
		t.Fatal(err)
	}
	p2, err := match.Join("bigfoot")
	if err != nil {
		t.Fatal(err)
	}
	return match, p1, p2
}

func TestProjectionHidesHandsAndDefense(t *testing.T) {
	match, p1, p2 := activeMatch(t, "privacy")
	placeForCombat(match, p1, p2)
	attack := giveCard(t, match, p1, "a-hunters-eye")
	defense := giveCard(t, match, p2, "its-just-your-imagination")
	beforeP2 := len(match.Players[p2].Hand)
	if err := match.Apply(p1, Command{Type: CommandAttack, ExpectedRevision: match.Revision, FighterID: match.Players[p1].HeroID, TargetID: match.Players[p2].HeroID, CardID: attack.ID}); err != nil {
		t.Fatal(err)
	}
	p1View, _ := match.Project(p1)
	p2View, _ := match.Project(p2)
	if len(p1View.Players[1].Hand) != 0 {
		t.Fatal("opponent hand leaked to attacker")
	}
	encoded, err := json.Marshal(p1View)
	if err != nil {
		t.Fatal(err)
	}
	for _, card := range match.Players[p2].Hand {
		if !containsDefinition(p1View.Players[0].Hand, card.DefinitionID) && strings.Contains(string(encoded), `"`+card.DefinitionID+`"`) {
			t.Fatalf("opponent card identity leaked: %s", card.DefinitionID)
		}
	}
	if p1View.Pending == nil || len(p1View.Pending.Options) != 0 {
		t.Fatal("defense domain leaked to attacker")
	}
	if p2View.Pending == nil || len(p2View.Pending.Options) == 0 {
		t.Fatal("defender did not receive private defense domain")
	}
	if p1View.Combat.AttackCard != nil || p2View.Combat.AttackCard != nil {
		t.Fatal("attack card revealed before defense choice")
	}
	revision := match.Revision
	if err := match.Apply(p1, Command{Type: CommandChoose, ExpectedRevision: revision, Choice: defense.ID}); ErrorCode(err) != "unauthorized_choice" {
		t.Fatalf("unexpected attacker choice error: %v", err)
	}
	if match.Revision != revision {
		t.Fatal("rejected choice mutated revision")
	}
	if err := match.Apply(p2, Command{Type: CommandChoose, ExpectedRevision: revision, Choice: defense.ID}); err != nil {
		t.Fatal(err)
	}
	if len(match.Players[p2].Hand) != beforeP2-1 {
		t.Fatal("defense card was not consumed")
	}
	if got := match.Fighters[match.Players[p2].HeroID].Health; got != 14 {
		t.Fatalf("Bigfoot health=%d want 14", got)
	}
	if err := match.Apply(p1, Command{Type: CommandManeuver, ExpectedRevision: revision}); ErrorCode(err) != "stale_revision" {
		t.Fatalf("stale error=%v", err)
	}
}

func TestVisibleCardProjectionIncludesDeckIdentity(t *testing.T) {
	match, p1, _ := activeMatch(t, "card-art-projection")
	view, err := match.Project(p1)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Players[0].Hand) == 0 {
		t.Fatal("expected Robin Hood hand")
	}
	for _, card := range view.Players[0].Hand {
		if card.DefinitionID == "" || card.DeckDefinitionID != "robin-hood" {
			t.Fatalf("visible card lacks correct identity: %+v", card)
		}
	}
}

func containsDefinition(cards []CardView, definitionID string) bool {
	for _, card := range cards {
		if card.DefinitionID == definitionID {
			return true
		}
	}
	return false
}

func TestRejectedMovementAttackAndAuthorityAreAtomic(t *testing.T) {
	match, p1, _ := activeMatch(t, "rejects")
	if err := match.Apply(p1, Command{Type: CommandManeuver, ExpectedRevision: match.Revision}); err != nil {
		t.Fatal(err)
	}
	if err := match.Apply(p1, Command{Type: CommandChoose, ExpectedRevision: match.Revision, Choice: "boost:none"}); err != nil {
		t.Fatal(err)
	}
	revision := match.Revision
	if err := match.Apply(p1, Command{Type: CommandChoose, ExpectedRevision: revision, Choice: "not-a-legal-option"}); ErrorCode(err) != "illegal_choice" {
		t.Fatalf("choice error=%v", err)
	}
	if match.Revision != revision {
		t.Fatal("illegal movement choice changed revision")
	}
	if err := match.Apply(p1, Command{Type: CommandChoose, ExpectedRevision: revision, Choice: "maneuver:done"}); err != nil {
		t.Fatal(err)
	}
	attack := giveCard(t, match, p1, "a-hunters-eye")
	hero := match.Fighters[match.Players[p1].HeroID]
	opposing := match.Fighters[match.Players[match.opponent(p1)].HeroID]
	hero.SpaceID = "s01"
	opposing.SpaceID = "s09"
	handBefore := len(match.Players[p1].Hand)
	revision = match.Revision
	err := match.Apply(p1, Command{Type: CommandAttack, ExpectedRevision: revision, FighterID: hero.ID, TargetID: opposing.ID, CardID: attack.ID})
	if ErrorCode(err) != "illegal_attack" {
		t.Fatalf("attack error=%v", err)
	}
	if match.Revision != revision || len(match.Players[p1].Hand) != handBefore {
		t.Fatal("illegal attack partially mutated state")
	}
	err = match.Apply("intruder", Command{Type: CommandManeuver, ExpectedRevision: revision})
	if ErrorCode(err) != "unauthorized" || match.Revision != revision {
		t.Fatalf("authority error=%v revision=%d", err, match.Revision)
	}
}

func TestDeterministicCompleteMatch(t *testing.T) {
	first := runCompleteMatch(t, "complete-match")
	second := runCompleteMatch(t, "complete-match")
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	if string(a) != string(b) {
		t.Fatal("scripted match final projection is not deterministic")
	}
	if first.WinnerID != "player-2" || first.Phase != PhaseEnded {
		t.Fatalf("winner=%s phase=%s", first.WinnerID, first.Phase)
	}
	seen := map[string]bool{}
	for _, event := range first.Events {
		seen[event.Type] = true
	}
	for _, kind := range []string{"maneuver_started", "scheme_played", "attack_declared", "combat_revealed", "game_ended"} {
		if !seen[kind] {
			t.Fatalf("missing event %s", kind)
		}
	}
}

func runCompleteMatch(t *testing.T, id string) View {
	t.Helper()
	match, p1, p2 := activeMatch(t, id)
	resetCards(match, p1)
	resetCards(match, p2)
	placeForCombat(match, p1, p2)
	outlaw := firstFighter(match, p1, "outlaw")
	boost := giveCard(t, match, p1, "a-hunters-eye")
	mustApply(t, match, p1, Command{Type: CommandManeuver})
	mustChoose(t, match, p1, findPromptOption(t, match, func(o Option) bool { return o.CardID == boost.ID }).ID)
	option := findPromptOption(t, match, func(o Option) bool { return o.FighterID == outlaw.ID && len(o.Path) > 0 })
	mustChoose(t, match, p1, option.ID)
	secondMove := findPromptOption(t, match, func(o Option) bool {
		fighter := match.Fighters[o.FighterID]
		return fighter != nil && fighter.DefinitionID == "outlaw" && fighter.ID != outlaw.ID && len(o.Path) > 0
	})
	mustChoose(t, match, p1, secondMove.ID)
	mustChoose(t, match, p1, "maneuver:done")
	steal := giveCard(t, match, p1, "steal-from-the-rich")
	mustApply(t, match, p1, Command{Type: CommandScheme, CardID: steal.ID})
	mustChoose(t, match, p2, "decline")
	if match.CurrentPlayer != p2 {
		t.Fatal("turn did not advance to Bigfoot")
	}

	crash := giveCard(t, match, p2, "crash-through-the-trees")
	larger := giveCard(t, match, p2, "larger-than-life")
	bigfoot := match.Fighters[match.Players[p2].HeroID]
	robin := match.Fighters[match.Players[p1].HeroID]
	mustApply(t, match, p2, Command{Type: CommandScheme, CardID: crash.ID, FighterID: bigfoot.ID, Path: []string{"s21", "s20"}})
	defenders := giveCard(t, match, p1, "defenders-of-sherwood")
	mustApply(t, match, p2, Command{Type: CommandAttack, FighterID: bigfoot.ID, TargetID: robin.ID, CardID: larger.ID})
	mustChoose(t, match, p1, defenders.ID)
	if robin.Health != 10 {
		t.Fatalf("Robin health after defended attack=%d", robin.Health)
	}

	disarming := giveCard(t, match, p1, "disarming-shot")
	regroup := giveCard(t, match, p2, "regroup")
	mustApply(t, match, p1, Command{Type: CommandAttack, FighterID: robin.ID, TargetID: bigfoot.ID, CardID: disarming.ID})
	mustChoose(t, match, p2, regroup.ID)
	mustChoose(t, match, p1, findPromptOption(t, match, func(o Option) bool { return o.Destination == robin.SpaceID }).ID)
	mustApply(t, match, p1, Command{Type: CommandManeuver})
	mustChoose(t, match, p1, "boost:none")
	mustChoose(t, match, p1, "maneuver:done")
	for len(match.Pending) > 0 && match.Pending[0].Kind == PromptDiscardDown {
		mustChoose(t, match, p1, match.Pending[0].Options[0].ID)
	}
	if match.CurrentPlayer != p2 {
		t.Fatal("second turn did not advance")
	}

	larger2 := giveCard(t, match, p2, "larger-than-life")
	savagery := giveCard(t, match, p2, "savagery")
	mustApply(t, match, p2, Command{Type: CommandAttack, FighterID: bigfoot.ID, TargetID: robin.ID, CardID: larger2.ID})
	mustChoose(t, match, p1, "none")
	mustApply(t, match, p2, Command{Type: CommandAttack, FighterID: bigfoot.ID, TargetID: robin.ID, CardID: savagery.ID})
	mustChoose(t, match, p1, "none")
	view, err := match.Project(p2)
	if err != nil {
		t.Fatal(err)
	}
	return view
}

func mustApply(t *testing.T, m *Match, p string, c Command) {
	t.Helper()
	c.ExpectedRevision = m.Revision
	if err := m.Apply(p, c); err != nil {
		t.Fatalf("apply %s: %v", c.Type, err)
	}
}
func mustChoose(t *testing.T, m *Match, p, choice string) {
	t.Helper()
	mustApply(t, m, p, Command{Type: CommandChoose, Choice: choice})
}
func findPromptOption(t *testing.T, m *Match, predicate func(Option) bool) Option {
	t.Helper()
	if len(m.Pending) == 0 {
		t.Fatal("missing prompt")
	}
	for _, o := range m.Pending[0].Options {
		if predicate(o) {
			return o
		}
	}
	t.Fatalf("no matching option in %v", m.Pending[0].Options)
	return Option{}
}

func resetCards(m *Match, playerID string) {
	p := m.Players[playerID]
	all := append([]CardInstance{}, p.Deck...)
	all = append(all, p.Hand...)
	all = append(all, p.Discard...)
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	p.Deck = all
	p.Hand = nil
	p.Discard = nil
}
func giveCard(t *testing.T, m *Match, playerID, definitionID string) CardInstance {
	t.Helper()
	p := m.Players[playerID]
	for i, c := range p.Hand {
		if c.DefinitionID == definitionID {
			return p.Hand[i]
		}
	}
	for i, c := range p.Deck {
		if c.DefinitionID == definitionID {
			p.Deck = append(p.Deck[:i:i], p.Deck[i+1:]...)
			p.Hand = append(p.Hand, c)
			return c
		}
	}
	t.Fatalf("no card %s for %s", definitionID, playerID)
	return CardInstance{}
}
func firstFighter(m *Match, playerID, definitionID string) *Fighter {
	for _, id := range m.Players[playerID].FighterIDs {
		if m.Fighters[id].DefinitionID == definitionID {
			return m.Fighters[id]
		}
	}
	return nil
}
func placeForCombat(m *Match, p1, p2 string) {
	placements := map[string]string{m.Players[p1].HeroID: "s21", m.Players[p2].HeroID: "s22"}
	remote1 := []string{"s01", "s03", "s12", "s14"}
	i := 0
	for _, id := range m.Players[p1].FighterIDs {
		if id == m.Players[p1].HeroID {
			continue
		}
		placements[id] = remote1[i]
		i++
	}
	for _, id := range m.Players[p2].FighterIDs {
		if id != m.Players[p2].HeroID {
			placements[id] = "s08"
		}
	}
	for id, space := range placements {
		m.Fighters[id].SpaceID = space
	}
	m.captureTurnStart(m.CurrentPlayer)
}
