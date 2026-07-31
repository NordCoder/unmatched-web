// Package content owns the launch-slice registration boundary. It is the only
// production package allowed to translate fighter or card identities into typed
// gameplay semantics.
package content

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	phase4bcards "github.com/NordCoder/unmatched-web/docs/cards/phase-4b"
	phase4bfighters "github.com/NordCoder/unmatched-web/docs/fighters/phase-4b"
)

type CardType string

const (
	CardAttack    CardType = "attack"
	CardDefense   CardType = "defense"
	CardVersatile CardType = "versatile"
	CardScheme    CardType = "scheme"
)

type AttackType string

const (
	AttackMelee  AttackType = "melee"
	AttackRanged AttackType = "ranged"
)

// EffectKind is a closed semantic vocabulary for the ordinary mechanics used
// by the launch decks. The game kernel dispatches on these mechanics, never on
// card identities.
type EffectKind string

const (
	EffectNone                   EffectKind = "none"
	EffectStealFromRich          EffectKind = "steal_from_rich"
	EffectDrawDamage             EffectKind = "draw_damage"
	EffectDrawTwo                EffectKind = "draw_two"
	EffectCancelIgnoreValue      EffectKind = "cancel_ignore_value"
	EffectDrawReturnSidekick     EffectKind = "draw_return_sidekick"
	EffectCancelOpposing         EffectKind = "cancel_opposing"
	EffectRegroup                EffectKind = "regroup"
	EffectAdjacentOpponentDamage EffectKind = "adjacent_opponent_damage"
	EffectDrawIfAdjacent         EffectKind = "draw_if_adjacent"
	EffectAmbush                 EffectKind = "ambush"
	EffectSavagery               EffectKind = "savagery"
	EffectMoveThroughFive        EffectKind = "move_through_five"
	EffectHorns                  EffectKind = "horns"
	EffectPostCombatMoveFive     EffectKind = "post_combat_move_five"
	EffectPostCombatPlaceZone    EffectKind = "post_combat_place_zone"
	EffectSkirmish               EffectKind = "skirmish"
	EffectMomentousShift         EffectKind = "momentous_shift"
)

type FighterAbility string

const (
	AbilityPostAttackRetreat FighterAbility = "post_attack_retreat"
	AbilityIsolatedEndDraw   FighterAbility = "isolated_end_draw"
)

type CardDefinition struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Quantity int        `json:"quantity"`
	UsableBy []string   `json:"usable_by"`
	Type     CardType   `json:"type"`
	Value    int        `json:"value"`
	Boost    int        `json:"boost"`
	Effect   EffectKind `json:"effect"`
}

type FighterDefinition struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Role           string     `json:"role"`
	Count          int        `json:"count"`
	AttackType     AttackType `json:"attack_type"`
	StartingHealth int        `json:"starting_health"`
}

type DeckDefinition struct {
	ID       string                       `json:"id"`
	Name     string                       `json:"name"`
	Movement int                          `json:"movement"`
	HeroID   string                       `json:"hero_id"`
	Ability  FighterAbility               `json:"ability"`
	Fighters map[string]FighterDefinition `json:"fighters"`
	Cards    map[string]CardDefinition    `json:"cards"`
}

type SpaceDefinition struct {
	ID           string   `json:"id"`
	X            int      `json:"x"`
	Y            int      `json:"y"`
	Zones        []string `json:"zones"`
	StartingSeat int      `json:"starting_seat,omitempty"`
}

type EdgeDefinition struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type BattlefieldDefinition struct {
	SchemaVersion int               `json:"schema_version"`
	ID            string            `json:"id"`
	DisplayName   string            `json:"display_name"`
	Zones         []string          `json:"zones"`
	Spaces        []SpaceDefinition `json:"spaces"`
	Edges         []EdgeDefinition  `json:"edges"`
}

type Registry struct {
	Decks       map[string]DeckDefinition
	Cards       map[string]CardDefinition
	Battlefield BattlefieldDefinition
}

//go:embed data/sherwood-forest.json
var sherwoodForestJSON []byte

func Load() (*Registry, error) {
	robin, err := parseDeck("robin-hood", "Robin Hood", phase4bfighters.RobinHood, phase4bcards.RobinHood)
	if err != nil {
		return nil, fmt.Errorf("load robin-hood manifests: %w", err)
	}
	bigfoot, err := parseDeck("bigfoot", "Bigfoot", phase4bfighters.Bigfoot, phase4bcards.Bigfoot)
	if err != nil {
		return nil, fmt.Errorf("load bigfoot manifests: %w", err)
	}
	var battlefield BattlefieldDefinition
	if err := json.Unmarshal(sherwoodForestJSON, &battlefield); err != nil {
		return nil, fmt.Errorf("decode Sherwood Forest: %w", err)
	}
	if err := validateBattlefield(battlefield); err != nil {
		return nil, fmt.Errorf("validate Sherwood Forest: %w", err)
	}
	cards := make(map[string]CardDefinition)
	for _, deck := range []DeckDefinition{robin, bigfoot} {
		for id, card := range deck.Cards {
			if existing, ok := cards[id]; ok {
				if existing.Name != card.Name || existing.Type != card.Type || existing.Value != card.Value || existing.Boost != card.Boost || existing.Effect != card.Effect {
					return nil, fmt.Errorf("shared card %q has inconsistent definitions", id)
				}
				continue
			}
			cards[id] = card
		}
	}
	return &Registry{
		Decks: map[string]DeckDefinition{robin.ID: robin, bigfoot.ID: bigfoot},
		Cards: cards, Battlefield: battlefield,
	}, nil
}

func parseDeck(id, name string, fighterYAML, cardYAML []byte) (DeckDefinition, error) {
	movement, fighters, ability, err := parseFighterManifest(fighterYAML)
	if err != nil {
		return DeckDefinition{}, err
	}
	cards, err := parseCardManifest(cardYAML)
	if err != nil {
		return DeckDefinition{}, err
	}
	quantity := 0
	for _, card := range cards {
		quantity += card.Quantity
	}
	if quantity != 30 {
		return DeckDefinition{}, fmt.Errorf("deck %s has %d cards, want 30", id, quantity)
	}
	heroID := ""
	for fighterID, fighter := range fighters {
		if fighter.Role == "hero" {
			if heroID != "" {
				return DeckDefinition{}, errors.New("multiple heroes in fighter manifest")
			}
			heroID = fighterID
		}
	}
	if heroID == "" {
		return DeckDefinition{}, errors.New("fighter manifest has no hero")
	}
	return DeckDefinition{ID: id, Name: name, Movement: movement, HeroID: heroID, Ability: ability, Fighters: fighters, Cards: cards}, nil
}

func parseFighterManifest(raw []byte) (int, map[string]FighterDefinition, FighterAbility, error) {
	movement := 0
	fighters := make(map[string]FighterDefinition)
	ability := FighterAbility("")
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "movement:"):
			movement, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(trimmed, "movement:")))
		case strings.HasPrefix(line, "    - {id:"):
			fields := parseInlineMap(strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
			count, _ := strconv.Atoi(fields["count"])
			health, _ := strconv.Atoi(fields["starting_health"])
			fighter := FighterDefinition{
				ID: fields["id"], Name: unquote(fields["display_name"]), Role: fields["role"], Count: count,
				AttackType: AttackType(fields["attack_type"]), StartingHealth: health,
			}
			if fighter.ID != "" {
				fighters[fighter.ID] = fighter
			}
		case strings.HasPrefix(line, "  id:"):
			ability = abilityForID(strings.TrimSpace(strings.TrimPrefix(trimmed, "id:")))
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, nil, "", err
	}
	if movement <= 0 || len(fighters) == 0 || ability == "" {
		return 0, nil, "", fmt.Errorf("incomplete fighter manifest: movement=%d fighters=%d ability=%q", movement, len(fighters), ability)
	}
	return movement, fighters, ability, nil
}

func parseCardManifest(raw []byte) (map[string]CardDefinition, error) {
	lines := strings.Split(string(raw), "\n")
	cards := make(map[string]CardDefinition)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if !strings.HasPrefix(line, "  - ") {
			continue
		}
		trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "- "))
		fields := map[string]string{}
		if strings.HasPrefix(trimmed, "{") {
			fields = parseInlineMap(trimmed)
		} else if strings.HasPrefix(trimmed, "id:") {
			fields["id"] = strings.TrimSpace(strings.TrimPrefix(trimmed, "id:"))
			for j := i + 1; j < len(lines) && !strings.HasPrefix(lines[j], "  - "); j++ {
				if !strings.HasPrefix(lines[j], "    ") || strings.HasPrefix(lines[j], "      ") {
					continue
				}
				key, value, ok := strings.Cut(strings.TrimSpace(lines[j]), ":")
				if ok {
					fields[key] = strings.TrimSpace(value)
				}
			}
		}
		card, err := cardFromFields(fields)
		if err != nil {
			return nil, fmt.Errorf("card manifest line %d: %w", i+1, err)
		}
		if _, exists := cards[card.ID]; exists {
			return nil, fmt.Errorf("duplicate card id %q", card.ID)
		}
		cards[card.ID] = card
	}
	if len(cards) == 0 {
		return nil, errors.New("card manifest contains no cards")
	}
	return cards, nil
}

func cardFromFields(fields map[string]string) (CardDefinition, error) {
	id := unquote(fields["id"])
	quantity, err := strconv.Atoi(fields["quantity"])
	if err != nil || quantity <= 0 {
		return CardDefinition{}, fmt.Errorf("card %q has invalid quantity", id)
	}
	value := 0
	if raw := fields["printed_value"]; raw != "" && raw != "null" {
		value, err = strconv.Atoi(raw)
		if err != nil {
			return CardDefinition{}, fmt.Errorf("card %q has invalid printed value", id)
		}
	}
	boost, err := strconv.Atoi(fields["boost"])
	if err != nil || boost < 0 {
		return CardDefinition{}, fmt.Errorf("card %q has invalid boost", id)
	}
	cardType := CardType(unquote(fields["type"]))
	if cardType != CardAttack && cardType != CardDefense && cardType != CardVersatile && cardType != CardScheme {
		return CardDefinition{}, fmt.Errorf("card %q has unsupported type %q", id, cardType)
	}
	return CardDefinition{
		ID: id, Name: unquote(fields["name"]), Quantity: quantity, UsableBy: parseList(fields["usable_by"]),
		Type: cardType, Value: value, Boost: boost, Effect: effectForCard(id),
	}, nil
}

func parseInlineMap(raw string) map[string]string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")
	result := make(map[string]string)
	start, depth := 0, 0
	quoted := false
	flush := func(end int) {
		part := strings.TrimSpace(raw[start:end])
		if key, value, ok := strings.Cut(part, ":"); ok {
			result[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	for i, r := range raw {
		switch r {
		case '"':
			quoted = !quoted
		case '[', '{':
			if !quoted {
				depth++
			}
		case ']', '}':
			if !quoted {
				depth--
			}
		case ',':
			if !quoted && depth == 0 {
				flush(i)
				start = i + 1
			}
		}
	}
	flush(len(raw))
	return result
}

func parseList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "any" {
		return []string{"any"}
	}
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		result = append(result, unquote(strings.TrimSpace(part)))
	}
	return result
}

func unquote(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		if decoded, err := strconv.Unquote(value); err == nil {
			return decoded
		}
	}
	return value
}

func abilityForID(id string) FighterAbility {
	switch id {
	case "outlaw-mobility":
		return AbilityPostAttackRetreat
	case "large-and-elusive":
		return AbilityIsolatedEndDraw
	default:
		return ""
	}
}

func effectForCard(id string) EffectKind {
	switch id {
	case "steal-from-the-rich":
		return EffectStealFromRich
	case "disarming-shot":
		return EffectDrawDamage
	case "piercing-shot":
		return EffectDrawTwo
	case "highway-robbery":
		return EffectCancelIgnoreValue
	case "defenders-of-sherwood":
		return EffectDrawReturnSidekick
	case "feint", "its-just-your-imagination":
		return EffectCancelOpposing
	case "regroup":
		return EffectRegroup
	case "wily-fighting":
		return EffectAdjacentOpponentDamage
	case "snark":
		return EffectDrawIfAdjacent
	case "ambush":
		return EffectAmbush
	case "savagery":
		return EffectSavagery
	case "crash-through-the-trees":
		return EffectMoveThroughFive
	case "jackalope-horns":
		return EffectHorns
	case "hoax":
		return EffectPostCombatMoveFive
	case "disengage":
		return EffectPostCombatPlaceZone
	case "skirmish":
		return EffectSkirmish
	case "momentous-shift":
		return EffectMomentousShift
	default:
		return EffectNone
	}
}

func validateBattlefield(field BattlefieldDefinition) error {
	if field.ID != "sherwood-forest" || len(field.Spaces) != 30 || len(field.Zones) != 7 || len(field.Edges) != 39 {
		return fmt.Errorf("unexpected identity/counts: id=%q spaces=%d zones=%d edges=%d", field.ID, len(field.Spaces), len(field.Zones), len(field.Edges))
	}
	spaces := make(map[string]SpaceDefinition, len(field.Spaces))
	starts := map[int]int{}
	for _, space := range field.Spaces {
		if space.ID == "" || len(space.Zones) == 0 {
			return fmt.Errorf("space has empty identity or zones")
		}
		if _, duplicate := spaces[space.ID]; duplicate {
			return fmt.Errorf("duplicate space %q", space.ID)
		}
		spaces[space.ID] = space
		if space.StartingSeat != 0 {
			starts[space.StartingSeat]++
		}
	}
	if starts[1] != 1 || starts[2] != 1 || spaces["s20"].StartingSeat != 1 || spaces["s19"].StartingSeat != 2 {
		return fmt.Errorf("starting spaces must be seat 1 at s20 and seat 2 at s19")
	}
	adjacency := make(map[string]map[string]struct{}, len(spaces))
	seen := make(map[string]struct{})
	for _, edge := range field.Edges {
		if edge.From == edge.To {
			return fmt.Errorf("self edge at %q", edge.From)
		}
		if _, ok := spaces[edge.From]; !ok {
			return fmt.Errorf("edge references unknown space %q", edge.From)
		}
		if _, ok := spaces[edge.To]; !ok {
			return fmt.Errorf("edge references unknown space %q", edge.To)
		}
		a, b := edge.From, edge.To
		if a > b {
			a, b = b, a
		}
		key := a + "\x00" + b
		if _, duplicate := seen[key]; duplicate {
			return fmt.Errorf("duplicate undirected edge %s-%s", a, b)
		}
		seen[key] = struct{}{}
		if adjacency[a] == nil {
			adjacency[a] = map[string]struct{}{}
		}
		if adjacency[b] == nil {
			adjacency[b] = map[string]struct{}{}
		}
		adjacency[a][b] = struct{}{}
		adjacency[b][a] = struct{}{}
	}
	visited := map[string]bool{}
	queue := []string{field.Spaces[0].ID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		for next := range adjacency[current] {
			if !visited[next] {
				queue = append(queue, next)
			}
		}
	}
	if len(visited) != len(spaces) {
		return fmt.Errorf("battlefield is disconnected: reached %d/%d", len(visited), len(spaces))
	}
	return nil
}
