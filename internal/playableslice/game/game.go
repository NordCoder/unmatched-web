package game

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"

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
	startBySeat := map[int]string{}
	for _, space := range m.Registry.Battlefield.Spaces {
		if space.StartingSeat != 0 {
			startBySeat[space.StartingSeat] = space.ID
		}
	}
	occupied := map[string]bool{}
	for _, playerID := range m.PlayerOrder {
		player := m.Players[playerID]
		deck := m.Registry.Decks[player.DeckDefinitionID]
		player.ActionsRemaining = 0
		player.TurnStartSpaces = map[string]string{}
		player.Deck = m.buildDeck(deck)
		m.shuffle(player.Deck)
		start := startBySeat[player.Seat]
		if start == "" {
			return fmt.Errorf("missing starting space for seat %d", player.Seat)
		}
		placements := m.startingPlacement(playerID, start, deck, occupied)
		for _, definition := range sortedFighters(deck.Fighters) {
			for i := 1; i <= definition.Count; i++ {
				instanceID := fmt.Sprintf("%s-%s", playerID, definition.ID)
				if definition.Count > 1 {
					instanceID = fmt.Sprintf("%s-%s-%d", playerID, definition.ID, i)
				}
				spaceID, ok := placements[instanceID]
				if !ok {
					return fmt.Errorf("no starting placement for %s", instanceID)
				}
				fighter := &Fighter{ID: instanceID, DefinitionID: definition.ID, OwnerID: playerID, Health: definition.StartingHealth, MaxHealth: definition.StartingHealth, SpaceID: spaceID}
				m.Fighters[instanceID] = fighter
				player.FighterIDs = append(player.FighterIDs, instanceID)
				occupied[spaceID] = true
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

func (m *Match) startingPlacement(playerID, heroStart string, deck content.DeckDefinition, occupied map[string]bool) map[string]string {
	result := map[string]string{}
	definitions := sortedFighters(deck.Fighters)
	available := m.nearestSpaces(heroStart, 6)
	cursor := 0
	for _, definition := range definitions {
		for i := 1; i <= definition.Count; i++ {
			instanceID := fmt.Sprintf("%s-%s", playerID, definition.ID)
			if definition.Count > 1 {
				instanceID = fmt.Sprintf("%s-%s-%d", playerID, definition.ID, i)
			}
			if definition.Role == "hero" && !occupied[heroStart] {
				result[instanceID] = heroStart
				occupied[heroStart] = true
				continue
			}
			for cursor < len(available) && occupied[available[cursor]] {
				cursor++
			}
			if cursor >= len(available) {
				return result
			}
			result[instanceID] = available[cursor]
			occupied[available[cursor]] = true
			cursor++
		}
	}
	return result
}

func (m *Match) nearestSpaces(start string, limit int) []string {
	adj := m.adjacency()
	visited := map[string]bool{start: true}
	queue := []string{start}
	result := []string{start}
	for len(queue) > 0 && len(result) < limit {
		current := queue[0]
		queue = queue[1:]
		next := append([]string(nil), adj[current]...)
		sort.Strings(next)
		for _, space := range next {
			if visited[space] {
				continue
			}
			visited[space] = true
			queue = append(queue, space)
			result = append(result, space)
			if len(result) == limit {
				break
			}
		}
	}
	return result
}

func (m *Match) buildDeck(deck content.DeckDefinition) []CardInstance {
	definitions := make([]content.CardDefinition, 0, len(deck.Cards))
	for _, definition := range deck.Cards {
		definitions = append(definitions, definition)
	}
	sort.Slice(definitions, func(i, j int) bool { return definitions[i].ID < definitions[j].ID })
	result := make([]CardInstance, 0, 30)
	for _, definition := range definitions {
		for i := 1; i <= definition.Quantity; i++ {
			m.cardSequence++
			result = append(result, CardInstance{ID: fmt.Sprintf("card-%03d", m.cardSequence), DefinitionID: definition.ID})
		}
	}
	return result
}

func (m *Match) shuffle(cards []CardInstance) {
	for i := len(cards) - 1; i > 0; i-- {
		j := int(m.random() % uint64(i+1))
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func (m *Match) random() uint64 {
	x := m.seed
	if x == 0 {
		x = 0x9e3779b97f4a7c15
	}
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	m.seed = x
	return x
}

func hashSeed(value string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(value))
	return h.Sum64()
}

func (m *Match) Apply(playerID string, command Command) error {
	if _, ok := m.Players[playerID]; !ok {
		return reject("unauthorized", "principal is not bound to this match")
	}
	if command.ExpectedRevision != m.Revision {
		return reject("stale_revision", "expected revision %d, current revision %d", command.ExpectedRevision, m.Revision)
	}
	if m.Phase == PhaseEnded {
		return reject("match_ended", "the match already has a winner")
	}
	if m.Phase != PhaseActive {
		return reject("match_not_active", "the second player has not joined")
	}
	var err error
	if len(m.Pending) > 0 && command.Type != CommandChoose {
		return reject("pending_choice", "resolve the pending choice before another action")
	}
	switch command.Type {
	case CommandManeuver:
		err = m.startManeuver(playerID)
	case CommandScheme:
		err = m.playScheme(playerID, command)
	case CommandAttack:
		err = m.startAttack(playerID, command)
	case CommandChoose:
		err = m.choose(playerID, command)
	default:
		err = reject("unknown_command", "unsupported command type %q", command.Type)
	}
	if err != nil {
		return err
	}
	m.Revision++
	return nil
}

func (m *Match) requireCurrent(playerID string) (*Player, error) {
	if m.CurrentPlayer != playerID {
		return nil, reject("not_current_player", "it is not this player's turn")
	}
	player := m.Players[playerID]
	if player.ActionsRemaining <= 0 {
		return nil, reject("no_actions_remaining", "the player has no actions remaining")
	}
	return player, nil
}

func (m *Match) draw(playerID string, count int) {
	player := m.Players[playerID]
	for i := 0; i < count; i++ {
		if len(player.Deck) == 0 {
			m.log("exhaustion", fmt.Sprintf("%s exhausted their deck", player.Name), map[string]any{"player_id": playerID})
			for _, fighterID := range player.FighterIDs {
				fighter := m.Fighters[fighterID]
				if !fighter.Defeated {
					m.damage(fighterID, 2, "exhaustion")
				}
			}
			continue
		}
		card := player.Deck[0]
		player.Deck = player.Deck[1:]
		player.Hand = append(player.Hand, card)
	}
}

func (m *Match) damage(fighterID string, amount int, source string) int {
	fighter := m.Fighters[fighterID]
	if fighter == nil || fighter.Defeated || amount <= 0 {
		return 0
	}
	if amount > fighter.Health {
		amount = fighter.Health
	}
	fighter.Health -= amount
	if fighter.Health == 0 {
		fighter.Defeated = true
		fighter.SpaceID = ""
		m.log("fighter_defeated", fmt.Sprintf("%s was defeated", m.fighterName(fighter)), map[string]any{"fighter_id": fighter.ID, "source": source})
	}
	m.checkWinner()
	return amount
}

func (m *Match) heal(fighterID string, amount int) int {
	fighter := m.Fighters[fighterID]
	if fighter == nil || fighter.Defeated || amount <= 0 {
		return 0
	}
	before := fighter.Health
	fighter.Health += amount
	if fighter.Health > fighter.MaxHealth {
		fighter.Health = fighter.MaxHealth
	}
	return fighter.Health - before
}

func (m *Match) checkWinner() {
	if m.Phase == PhaseEnded {
		return
	}
	defeatedHeroes := []string{}
	for _, playerID := range m.PlayerOrder {
		if hero := m.Fighters[m.Players[playerID].HeroID]; hero != nil && hero.Defeated {
			defeatedHeroes = append(defeatedHeroes, playerID)
		}
	}
	if len(defeatedHeroes) == 0 {
		return
	}
	winner := ""
	if len(defeatedHeroes) == 1 {
		winner = m.opponent(defeatedHeroes[0])
	} else {
		winner = m.opponent(m.CurrentPlayer)
	}
	m.WinnerID = winner
	m.Phase = PhaseEnded
	m.Pending = nil
	m.Combat = nil
	m.log("game_ended", fmt.Sprintf("%s wins", m.Players[winner].Name), map[string]any{"winner_id": winner})
}

func (m *Match) opponent(playerID string) string {
	for _, candidate := range m.PlayerOrder {
		if candidate != playerID {
			return candidate
		}
	}
	return ""
}

func (m *Match) captureTurnStart(playerID string) {
	m.TurnStartSpaces = map[string]string{}
	for fighterID, fighter := range m.Fighters {
		m.TurnStartSpaces[fighterID] = fighter.SpaceID
	}
	player := m.Players[playerID]
	player.TurnStartSpaces = map[string]string{}
	for _, fighterID := range player.FighterIDs {
		player.TurnStartSpaces[fighterID] = m.Fighters[fighterID].SpaceID
	}
}

func (m *Match) log(eventType, message string, data map[string]any) {
	m.Events = append(m.Events, Event{At: time.Unix(0, int64(len(m.Events)+1)), Type: eventType, Message: message, Data: data})
	if len(m.Events) > 100 {
		m.Events = append([]Event(nil), m.Events[len(m.Events)-100:]...)
	}
}

func (m *Match) fighterName(fighter *Fighter) string {
	if fighter == nil {
		return "unknown fighter"
	}
	deck := m.Registry.Decks[m.Players[fighter.OwnerID].DeckDefinitionID]
	if definition, ok := deck.Fighters[fighter.DefinitionID]; ok {
		return definition.Name
	}
	return fighter.DefinitionID
}

func (m *Match) cardDefinition(card CardInstance) content.CardDefinition {
	return m.Registry.Cards[card.DefinitionID]
}

func removeCard(cards []CardInstance, id string) ([]CardInstance, CardInstance, bool) {
	for i, card := range cards {
		if card.ID == id {
			return append(cards[:i:i], cards[i+1:]...), card, true
		}
	}
	return cards, CardInstance{}, false
}

func findCard(cards []CardInstance, id string) (CardInstance, bool) {
	for _, card := range cards {
		if card.ID == id {
			return card, true
		}
	}
	return CardInstance{}, false
}

func (m *Match) fighterDefinition(fighter *Fighter) content.FighterDefinition {
	return m.Registry.Decks[m.Players[fighter.OwnerID].DeckDefinitionID].Fighters[fighter.DefinitionID]
}

func usableBy(card content.CardDefinition, fighter *Fighter) bool {
	for _, id := range card.UsableBy {
		if id == "any" || id == fighter.DefinitionID {
			return true
		}
	}
	return false
}

func sharesZone(a, b content.SpaceDefinition) bool {
	for _, za := range a.Zones {
		for _, zb := range b.Zones {
			if za == zb {
				return true
			}
		}
	}
	return false
}

func (m *Match) space(id string) (content.SpaceDefinition, bool) {
	for _, space := range m.Registry.Battlefield.Spaces {
		if space.ID == id {
			return space, true
		}
	}
	return content.SpaceDefinition{}, false
}

func (m *Match) adjacency() map[string][]string {
	result := make(map[string][]string)
	for _, edge := range m.Registry.Battlefield.Edges {
		result[edge.From] = append(result[edge.From], edge.To)
		result[edge.To] = append(result[edge.To], edge.From)
	}
	return result
}

func (m *Match) adjacent(a, b string) bool {
	for _, next := range m.adjacency()[a] {
		if next == b {
			return true
		}
	}
	return false
}

func (m *Match) occupied(spaceID string) *Fighter {
	for _, fighter := range m.Fighters {
		if !fighter.Defeated && fighter.SpaceID == spaceID {
			return fighter
		}
	}
	return nil
}

func (m *Match) validatePath(fighter *Fighter, path []string, max int, allowOpposingThrough bool) error {
	if fighter == nil || fighter.Defeated {
		return reject("invalid_fighter", "fighter is unavailable")
	}
	if len(path) > max {
		return reject("illegal_movement", "path length %d exceeds movement %d", len(path), max)
	}
	current := fighter.SpaceID
	for i, destination := range path {
		if _, ok := m.space(destination); !ok {
			return reject("illegal_movement", "unknown destination %q", destination)
		}
		if !m.adjacent(current, destination) {
			return reject("illegal_movement", "%s and %s are not adjacent", current, destination)
		}
		occupant := m.occupied(destination)
		if occupant != nil {
			if i == len(path)-1 {
				return reject("illegal_movement", "destination %s is occupied", destination)
			}
			if occupant.OwnerID != fighter.OwnerID && !allowOpposingThrough {
				return reject("illegal_movement", "cannot move through opposing fighter at %s", destination)
			}
		}
		current = destination
	}
	return nil
}

func (m *Match) move(fighter *Fighter, path []string) {
	if len(path) == 0 {
		return
	}
	from := fighter.SpaceID
	fighter.SpaceID = path[len(path)-1]
	m.log("fighter_moved", fmt.Sprintf("%s moved from %s to %s", m.fighterName(fighter), from, fighter.SpaceID), map[string]any{"fighter_id": fighter.ID, "path": append([]string(nil), path...)})
}

func (m *Match) destinations(fighter *Fighter, max int, allowOpposingThrough bool, sameZoneOnly bool) []Option {
	if fighter == nil || fighter.Defeated {
		return nil
	}
	type node struct {
		id   string
		path []string
	}
	queue := []node{{fighter.SpaceID, nil}}
	best := map[string]int{fighter.SpaceID: 0}
	options := []Option{{ID: "stay", Label: "Do not move", FighterID: fighter.ID, Destination: fighter.SpaceID}}
	origin, _ := m.space(fighter.SpaceID)
	adj := m.adjacency()
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if len(current.path) == max {
			continue
		}
		nexts := append([]string(nil), adj[current.id]...)
		sort.Strings(nexts)
		for _, next := range nexts {
			path := append(append([]string(nil), current.path...), next)
			if distance, seen := best[next]; seen && distance <= len(path) {
				continue
			}
			occupant := m.occupied(next)
			if occupant != nil && occupant.OwnerID != fighter.OwnerID && !allowOpposingThrough {
				continue
			}
			best[next] = len(path)
			if occupant == nil {
				if destination, ok := m.space(next); ok && (!sameZoneOnly || sharesZone(origin, destination)) {
					options = append(options, Option{ID: fighter.ID + ":" + next, Label: next, FighterID: fighter.ID, Destination: next, Path: path})
				}
			}
			queue = append(queue, node{next, path})
		}
	}
	return options
}

func (m *Match) zonesOverlap(a, b string) bool {
	sa, oka := m.space(a)
	sb, okb := m.space(b)
	return oka && okb && sharesZone(sa, sb)
}

func (m *Match) isolatedFromOpponents(fighter *Fighter) bool {
	if fighter == nil || fighter.Defeated {
		return false
	}
	for _, other := range m.Fighters {
		if other.OwnerID == fighter.OwnerID || other.Defeated {
			continue
		}
		if m.zonesOverlap(fighter.SpaceID, other.SpaceID) {
			return false
		}
	}
	return true
}

func normalizeChoice(value string) string { return strings.TrimSpace(value) }
