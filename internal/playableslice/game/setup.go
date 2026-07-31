package game

import (
	"fmt"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func New(matchID string, registry *content.Registry, seatOneDeckID string) (*Match, string, error) {
	if registry == nil {
		return nil, "", fmt.Errorf("registry is nil")
	}
	if matchID == "" {
		return nil, "", fmt.Errorf("match id is empty")
	}
	if _, ok := registry.Decks[seatOneDeckID]; !ok {
		return nil, "", fmt.Errorf("unknown deck %q", seatOneDeckID)
	}
	playerID := "player-1"
	match := &Match{
		ID: matchID, Phase: PhaseWaiting, Registry: registry, Players: map[string]*Player{}, Fighters: map[string]*Fighter{},
		PlayerOrder: []string{playerID}, seed: hashSeed(matchID),
	}
	match.Players[playerID] = &Player{ID: playerID, Seat: 1, Name: registry.Decks[seatOneDeckID].Name, DeckDefinitionID: seatOneDeckID}
	match.log("match_created", fmt.Sprintf("%s created the match", match.Players[playerID].Name), nil)
	return match, playerID, nil
}

func (m *Match) Join(deckID string) (string, error) {
	if m.Phase != PhaseWaiting {
		return "", reject("match_not_joinable", "match is not waiting for a second player")
	}
	if len(m.PlayerOrder) != 1 {
		return "", reject("match_not_joinable", "match already has two players")
	}
	if _, ok := m.Registry.Decks[deckID]; !ok {
		return "", reject("unknown_deck", "unknown deck %q", deckID)
	}
	if deckID == m.Players[m.PlayerOrder[0]].DeckDefinitionID {
		return "", reject("fixed_matchup", "the launch slice requires two different fixed decks")
	}
	playerID := "player-2"
	m.Players[playerID] = &Player{ID: playerID, Seat: 2, Name: m.Registry.Decks[deckID].Name, DeckDefinitionID: deckID}
	m.PlayerOrder = append(m.PlayerOrder, playerID)
	if err := m.setup(); err != nil {
		delete(m.Players, playerID)
		m.PlayerOrder = m.PlayerOrder[:1]
		return "", err
	}
	m.Revision++
	return playerID, nil
}

func (m *Match) setup() error {
	if len(m.PlayerOrder) != 2 {
		return fmt.Errorf("setup requires two players")
	}
	startBySeat := make(map[int]string)
	for _, space := range m.Registry.Battlefield.Spaces {
		if space.StartingSeat != 0 {
			startBySeat[space.StartingSeat] = space.ID
		}
	}
	placements, err := m.planStartingPlacements(startBySeat)
	if err != nil {
		return err
	}
	if err := m.validateStartingPlacementPlan(placements); err != nil {
		return err
	}

	m.Fighters = make(map[string]*Fighter)
	for _, playerID := range m.PlayerOrder {
		player := m.Players[playerID]
		deck := m.Registry.Decks[player.DeckDefinitionID]
		player.ActionsRemaining = 0
		player.TurnStartSpaces = map[string]string{}
		player.Deck = m.buildDeck(deck)
		m.shuffle(player.Deck)

		for _, definition := range sortedFighters(deck.Fighters) {
			for i := 1; i <= definition.Count; i++ {
				instanceID := fighterInstanceID(playerID, definition, i)
				spaceID := placements[instanceID]
				fighter := &Fighter{
					ID: instanceID, DefinitionID: definition.ID, OwnerID: playerID,
					Health: definition.StartingHealth, MaxHealth: definition.StartingHealth, SpaceID: spaceID,
				}
				m.Fighters[instanceID] = fighter
				player.FighterIDs = append(player.FighterIDs, instanceID)
				if definition.Role == "hero" {
					player.HeroID = instanceID
				}
			}
		}
		for i := 0; i < 5; i++ {
			m.draw(playerID, 1)
		}
	}

	m.CurrentPlayer = m.PlayerOrder[0]
	m.Players[m.CurrentPlayer].ActionsRemaining = 2
	m.captureTurnStart(m.CurrentPlayer)
	m.Phase = PhaseActive
	m.log("setup_completed", "Deterministic setup completed; seat 1 takes the first turn", map[string]any{"starting_player_id": m.CurrentPlayer})
	return nil
}

func (m *Match) planStartingPlacements(startBySeat map[int]string) (map[string]string, error) {
	placements := make(map[string]string)
	occupied := make(map[string]bool)

	// Reserve every hero start before assigning any sidekick. This prevents the
	// first seat's sidekicks from consuming the second seat's printed marker.
	for _, playerID := range m.PlayerOrder {
		player := m.Players[playerID]
		deck := m.Registry.Decks[player.DeckDefinitionID]
		start := startBySeat[player.Seat]
		if start == "" {
			return nil, fmt.Errorf("missing starting space for seat %d", player.Seat)
		}
		hero, ok := deck.Fighters[deck.HeroID]
		if !ok || hero.Role != "hero" || hero.Count != 1 {
			return nil, fmt.Errorf("deck %s must define exactly one hero", deck.ID)
		}
		if occupied[start] {
			return nil, fmt.Errorf("starting space %s is assigned to multiple seats", start)
		}
		placements[fighterInstanceID(playerID, hero, 1)] = start
		occupied[start] = true
	}

	for _, playerID := range m.PlayerOrder {
		player := m.Players[playerID]
		deck := m.Registry.Decks[player.DeckDefinitionID]
		heroStart := startBySeat[player.Seat]
		available, err := m.legalSidekickStarts(heroStart, occupied)
		if err != nil {
			return nil, err
		}
		required := 0
		for _, definition := range deck.Fighters {
			if definition.Role != "hero" {
				required += definition.Count
			}
		}
		if len(available) < required {
			return nil, fmt.Errorf(
				"not enough legal starting spaces for %s sidekicks: need %d, have %d",
				player.Name, required, len(available),
			)
		}
		cursor := 0
		for _, definition := range sortedFighters(deck.Fighters) {
			if definition.Role == "hero" {
				continue
			}
			for i := 1; i <= definition.Count; i++ {
				instanceID := fighterInstanceID(playerID, definition, i)
				placements[instanceID] = available[cursor]
				occupied[available[cursor]] = true
				cursor++
			}
		}
	}
	return placements, nil
}

func (m *Match) validateStartingPlacementPlan(placements map[string]string) error {
	for _, playerID := range m.PlayerOrder {
		deck := m.Registry.Decks[m.Players[playerID].DeckDefinitionID]
		for _, definition := range sortedFighters(deck.Fighters) {
			for i := 1; i <= definition.Count; i++ {
				instanceID := fighterInstanceID(playerID, definition, i)
				if placements[instanceID] == "" {
					return fmt.Errorf("starting placement plan omitted %s", instanceID)
				}
			}
		}
	}
	return nil
}

func (m *Match) legalSidekickStarts(heroStart string, occupied map[string]bool) ([]string, error) {
	heroSpace, ok := m.space(heroStart)
	if !ok {
		return nil, fmt.Errorf("unknown hero starting space %q", heroStart)
	}
	distances := m.spaceDistances(heroStart)
	result := make([]string, 0)
	for _, candidate := range m.Registry.Battlefield.Spaces {
		if occupied[candidate.ID] || !sharesZone(heroSpace, candidate) {
			continue
		}
		if _, reachable := distances[candidate.ID]; !reachable {
			continue
		}
		result = append(result, candidate.ID)
	}
	sort.Slice(result, func(i, j int) bool {
		left, right := distances[result[i]], distances[result[j]]
		if left != right {
			return left < right
		}
		return result[i] < result[j]
	})
	return result, nil
}

func (m *Match) spaceDistances(start string) map[string]int {
	adjacency := m.adjacency()
	distances := map[string]int{start: 0}
	queue := []string{start}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		next := append([]string(nil), adjacency[current]...)
		sort.Strings(next)
		for _, candidate := range next {
			if _, seen := distances[candidate]; seen {
				continue
			}
			distances[candidate] = distances[current] + 1
			queue = append(queue, candidate)
		}
	}
	return distances
}

func fighterInstanceID(playerID string, definition content.FighterDefinition, number int) string {
	if definition.Count == 1 {
		return fmt.Sprintf("%s-%s", playerID, definition.ID)
	}
	return fmt.Sprintf("%s-%s-%d", playerID, definition.ID, number)
}

func sortedFighters(values map[string]content.FighterDefinition) []content.FighterDefinition {
	result := make([]content.FighterDefinition, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Role != result[j].Role {
			return result[i].Role == "hero"
		}
		return result[i].ID < result[j].ID
	})
	return result
}
