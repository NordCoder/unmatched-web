package game

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func (m *Match) startManeuver(playerID string) error {
	player, err := m.requireCurrent(playerID)
	if err != nil {
		return err
	}
	m.draw(playerID, 1)
	if m.Phase == PhaseEnded {
		return nil
	}
	player.ActionsRemaining--
	options := []Option{{ID: "boost:none", Label: "Do not BOOST movement"}}
	for _, card := range player.Hand {
		definition := m.cardDefinition(card)
		options = append(options, Option{ID: "boost:" + card.ID, CardID: card.ID, Label: fmt.Sprintf("BOOST with %s (+%d)", definition.Name, definition.Boost)})
	}
	m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptManeuverBoost, OwnerID: playerID, Message: "Choose an optional BOOST card. The resulting movement value applies to each fighter.", Options: options}})
	m.log("maneuver_started", fmt.Sprintf("%s drew a card and began a maneuver", player.Name), map[string]any{"player_id": playerID})
	return nil
}

func (m *Match) resolveManeuverBoost(playerID string, option Option) error {
	player := m.Players[playerID]
	movement := m.Registry.Decks[player.DeckDefinitionID].Movement
	if option.CardID != "" {
		boost, ok := findCard(player.Hand, option.CardID)
		if !ok {
			return reject("card_not_in_hand", "BOOST card is no longer in hand")
		}
		definition := m.cardDefinition(boost)
		movement += definition.Boost
		player.Hand, boost, _ = removeCard(player.Hand, boost.ID)
		player.Discard = append(player.Discard, boost)
		m.log("boost_used", fmt.Sprintf("%s discarded %s to BOOST movement", player.Name, definition.Name), map[string]any{"card": definition.Name, "movement": movement})
	}
	remaining := make([]string, 0, len(player.FighterIDs))
	for _, fighterID := range player.FighterIDs {
		if fighter := m.Fighters[fighterID]; fighter != nil && !fighter.Defeated {
			remaining = append(remaining, fighterID)
		}
	}
	m.queueManeuverMove(playerID, movement, remaining)
	return nil
}

func (m *Match) queueManeuverMove(playerID string, movement int, remaining []string) {
	options := []Option{{ID: "maneuver:done", Label: "Finish maneuver"}}
	for _, fighterID := range remaining {
		fighter := m.Fighters[fighterID]
		for _, option := range m.destinations(fighter, movement, false, false) {
			if option.Destination == fighter.SpaceID {
				continue
			}
			option.ID = encodeManeuverOption(fighter.ID, option.Path)
			option.Label = fmt.Sprintf("%s → %s", m.fighterName(fighter), option.Destination)
			options = append(options, option)
		}
	}
	m.Pending = append([]pendingTask{{
		Prompt: Prompt{Kind: PromptManeuverMove, OwnerID: playerID, Message: "Move any remaining fighter, or finish the maneuver.", Options: dedupeOptions(options)},
		data:   map[string]string{"movement": strconv.Itoa(movement), "remaining": strings.Join(remaining, ",")},
	}}, m.Pending...)
}

func encodeManeuverOption(fighterID string, path []string) string {
	return "maneuver:" + fighterID + ":" + strings.Join(path, ".")
}

func dedupeOptions(options []Option) []Option {
	seen := map[string]bool{}
	result := make([]Option, 0, len(options))
	for _, option := range options {
		if !seen[option.ID] {
			seen[option.ID] = true
			result = append(result, option)
		}
	}
	return result
}

func (m *Match) resolveManeuverMove(playerID string, option Option, data map[string]string) error {
	if option.ID == "maneuver:done" {
		m.log("maneuver_completed", fmt.Sprintf("%s completed a maneuver", m.Players[playerID].Name), nil)
		return nil
	}
	movement, err := strconv.Atoi(data["movement"])
	if err != nil || movement <= 0 {
		return reject("invalid_procedure", "maneuver movement context is invalid")
	}
	remaining := strings.Split(data["remaining"], ",")
	allowed := false
	next := make([]string, 0, len(remaining)-1)
	for _, fighterID := range remaining {
		if fighterID == option.FighterID {
			allowed = true
			continue
		}
		if fighterID != "" {
			next = append(next, fighterID)
		}
	}
	if !allowed {
		return reject("invalid_fighter", "fighter has already moved during this maneuver")
	}
	fighter := m.Fighters[option.FighterID]
	if fighter == nil || fighter.OwnerID != playerID || fighter.Defeated {
		return reject("invalid_fighter", "maneuver fighter is invalid")
	}
	if err := m.validatePath(fighter, option.Path, movement, false); err != nil {
		return err
	}
	m.move(fighter, option.Path)
	if len(next) == 0 {
		m.log("maneuver_completed", fmt.Sprintf("%s completed a maneuver", m.Players[playerID].Name), nil)
		return nil
	}
	m.queueManeuverMove(playerID, movement, next)
	return nil
}

func (m *Match) startAttack(playerID string, command Command) error {
	player, err := m.requireCurrent(playerID)
	if err != nil {
		return err
	}
	attacker := m.Fighters[command.FighterID]
	defender := m.Fighters[command.TargetID]
	if attacker == nil || attacker.OwnerID != playerID || attacker.Defeated {
		return reject("invalid_attacker", "attacker is not an available controlled fighter")
	}
	if defender == nil || defender.OwnerID == playerID || defender.Defeated {
		return reject("invalid_defender", "defender must be a living opposing fighter")
	}
	card, ok := findCard(player.Hand, command.CardID)
	if !ok {
		return reject("card_not_in_hand", "attack card is not in the player's hand")
	}
	definition := m.cardDefinition(card)
	if definition.Type != content.CardAttack && definition.Type != content.CardVersatile {
		return reject("illegal_card_type", "%s cannot be used to attack", definition.Name)
	}
	if !usableBy(definition, attacker) {
		return reject("fighter_cannot_use_card", "%s cannot use %s", m.fighterName(attacker), definition.Name)
	}
	attackType := m.fighterDefinition(attacker).AttackType
	legal := false
	if attackType == content.AttackMelee {
		legal = m.adjacent(attacker.SpaceID, defender.SpaceID)
	} else {
		legal = m.zonesOverlap(attacker.SpaceID, defender.SpaceID)
	}
	if !legal {
		return reject("illegal_attack", "target is outside the attacker's legal range")
	}

	player.Hand, card, _ = removeCard(player.Hand, card.ID)
	player.ActionsRemaining--
	m.Combat = &Combat{AttackerID: attacker.ID, DefenderID: defender.ID, AttackerPlayerID: playerID, DefenderPlayerID: defender.OwnerID, AttackCard: card}
	options := []Option{{ID: "none", Label: "Take the attack without a defense card"}}
	for _, candidate := range m.Players[defender.OwnerID].Hand {
		candidateDef := m.cardDefinition(candidate)
		if (candidateDef.Type == content.CardDefense || candidateDef.Type == content.CardVersatile) && usableBy(candidateDef, defender) {
			options = append(options, Option{ID: candidate.ID, CardID: candidate.ID, Label: candidateDef.Name + fmt.Sprintf(" (%d)", candidateDef.Value)})
		}
	}
	m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptDefense, OwnerID: defender.OwnerID, Message: "Choose an optional defense card. Both combat cards remain hidden until this choice.", Options: options}})
	m.log("attack_declared", fmt.Sprintf("%s attacked %s", m.fighterName(attacker), m.fighterName(defender)), map[string]any{"attacker_id": attacker.ID, "defender_id": defender.ID})
	return nil
}

func (m *Match) choose(playerID string, command Command) error {
	if len(m.Pending) == 0 {
		return reject("no_pending_choice", "there is no pending choice")
	}
	task := m.Pending[0]
	if task.OwnerID != playerID {
		return reject("unauthorized_choice", "the pending choice belongs to another player")
	}
	choice := normalizeChoice(command.Choice)
	var selected *Option
	for i := range task.Options {
		if task.Options[i].ID == choice {
			selected = &task.Options[i]
			break
		}
	}
	if selected == nil {
		return reject("illegal_choice", "choice %q is outside the legal domain", choice)
	}
	m.Pending = m.Pending[1:]
	var err error
	switch task.Kind {
	case PromptManeuverBoost:
		err = m.resolveManeuverBoost(playerID, *selected)
	case PromptManeuverMove:
		err = m.resolveManeuverMove(playerID, *selected, task.data)
	case PromptDefense:
		err = m.resolveDefense(playerID, *selected)
	case PromptStealDiscard:
		err = m.resolveSteal(playerID, *selected, task.data)
	case PromptMove, PromptPlace:
		err = m.resolveMoveChoice(*selected)
	case PromptSkirmish:
		err = m.resolveSkirmishChoice(*selected)
	case PromptReturnSidekick:
		err = m.resolveReturnSidekick(*selected)
	case PromptIsolatedDraw:
		err = m.resolveIsolatedDraw(playerID, *selected)
	case PromptDiscardDown:
		err = m.resolveDiscardDown(playerID, *selected)
	default:
		err = reject("unknown_prompt", "unsupported prompt kind %q", task.Kind)
	}
	if err != nil {
		m.Pending = append([]pendingTask{task}, m.Pending...)
		return err
	}
	if m.Phase != PhaseEnded && len(m.Pending) == 0 {
		m.continueAfterPrompts()
	}
	return nil
}

func (m *Match) resolveDefense(playerID string, option Option) error {
	if m.Combat == nil || m.Combat.DefenderPlayerID != playerID {
		return reject("combat_missing", "defense no longer belongs to an active combat")
	}
	if option.ID != "none" {
		card, ok := findCard(m.Players[playerID].Hand, option.CardID)
		if !ok {
			return reject("card_not_in_hand", "defense card is no longer in hand")
		}
		definition := m.cardDefinition(card)
		defender := m.Fighters[m.Combat.DefenderID]
		if (definition.Type != content.CardDefense && definition.Type != content.CardVersatile) || !usableBy(definition, defender) {
			return reject("illegal_defense", "selected card cannot defend this fighter")
		}
		m.Players[playerID].Hand, card, _ = removeCard(m.Players[playerID].Hand, card.ID)
		m.Combat.DefenseCard = &card
	}
	return m.resolveCombat()
}

func (m *Match) resolveSteal(playerID string, option Option, data map[string]string) error {
	ownerID := data["scheme_owner"]
	if option.ID == "decline" {
		m.draw(ownerID, 1)
		m.log("scheme_choice", fmt.Sprintf("%s declined to discard", m.Players[playerID].Name), nil)
		return nil
	}
	card, ok := findCard(m.Players[playerID].Hand, option.CardID)
	if !ok {
		return reject("card_not_in_hand", "selected discard is no longer in hand")
	}
	m.Players[playerID].Hand, card, _ = removeCard(m.Players[playerID].Hand, card.ID)
	m.Players[playerID].Discard = append(m.Players[playerID].Discard, card)
	m.log("card_discarded", fmt.Sprintf("%s discarded a card for the scheme", m.Players[playerID].Name), map[string]any{"player_id": playerID})
	return nil
}

func (m *Match) resolveMoveChoice(option Option) error {
	fighter := m.Fighters[option.FighterID]
	if fighter == nil || fighter.Defeated {
		return reject("invalid_fighter", "move target is no longer available")
	}
	if len(option.Path) > 0 {
		m.move(fighter, option.Path)
	} else if option.Destination != "" && option.Destination != fighter.SpaceID {
		if m.occupied(option.Destination) != nil {
			return reject("illegal_movement", "placement destination is occupied")
		}
		from := fighter.SpaceID
		fighter.SpaceID = option.Destination
		m.log("fighter_placed", fmt.Sprintf("%s was placed from %s to %s", m.fighterName(fighter), from, fighter.SpaceID), map[string]any{"fighter_id": fighter.ID, "destination": fighter.SpaceID})
	}
	return nil
}

func (m *Match) resolveSkirmishChoice(option Option) error { return m.resolveMoveChoice(option) }

func (m *Match) resolveReturnSidekick(option Option) error {
	if option.ID == "decline" {
		return nil
	}
	fighter := m.Fighters[option.FighterID]
	if fighter == nil || !fighter.Defeated {
		return reject("invalid_fighter", "sidekick is no longer defeated")
	}
	if m.occupied(option.Destination) != nil {
		return reject("illegal_placement", "return destination is occupied")
	}
	fighter.Defeated = false
	fighter.Health = fighter.MaxHealth
	fighter.SpaceID = option.Destination
	m.log("fighter_returned", fmt.Sprintf("%s returned to %s", m.fighterName(fighter), option.Destination), map[string]any{"fighter_id": fighter.ID, "destination": option.Destination})
	return nil
}

func (m *Match) resolveIsolatedDraw(playerID string, option Option) error {
	if option.ID == "draw" {
		m.draw(playerID, 1)
		m.log("ability_used", fmt.Sprintf("%s drew a card while isolated", m.Players[playerID].Name), nil)
	}
	return nil
}

func (m *Match) resolveDiscardDown(playerID string, option Option) error {
	card, ok := findCard(m.Players[playerID].Hand, option.CardID)
	if !ok {
		return reject("card_not_in_hand", "discard choice is no longer in hand")
	}
	m.Players[playerID].Hand, card, _ = removeCard(m.Players[playerID].Hand, card.ID)
	m.Players[playerID].Discard = append(m.Players[playerID].Discard, card)
	m.log("hand_limit_discard", fmt.Sprintf("%s discarded to the hand limit", m.Players[playerID].Name), nil)
	return nil
}

func (m *Match) finishAction() {
	if m.Phase == PhaseEnded || len(m.Pending) > 0 {
		return
	}
	player := m.Players[m.CurrentPlayer]
	if player.ActionsRemaining > 0 {
		return
	}
	m.queueEndTurn()
}

func (m *Match) continueAfterPrompts() {
	if m.Phase == PhaseEnded {
		return
	}
	if m.endTurnQueued {
		if len(m.Players[m.CurrentPlayer].Hand) > 7 {
			m.queueDiscardPrompt(m.CurrentPlayer)
			return
		}
		m.switchTurn()
		return
	}
	m.finishAction()
}

func (m *Match) queueEndTurn() {
	if m.endTurnQueued || m.Phase == PhaseEnded {
		return
	}
	m.endTurnQueued = true
	playerID := m.CurrentPlayer
	player := m.Players[playerID]
	deck := m.Registry.Decks[player.DeckDefinitionID]
	if deck.Ability == content.AbilityIsolatedEndDraw && m.isolatedFromOpponents(m.Fighters[player.HeroID]) {
		m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptIsolatedDraw, OwnerID: playerID, Message: "Your hero shares no zone with an opposing fighter. Draw one card?", Options: []Option{{ID: "draw", Label: "Draw 1 card"}, {ID: "decline", Label: "Do not draw"}}}})
	}
	if len(player.Hand) > 7 {
		m.queueDiscardPrompt(playerID)
	}
	if len(m.Pending) == 0 {
		m.switchTurn()
	}
}

func (m *Match) queueDiscardPrompt(playerID string) {
	for _, task := range m.Pending {
		if task.Kind == PromptDiscardDown {
			return
		}
	}
	options := make([]Option, 0, len(m.Players[playerID].Hand))
	for _, card := range m.Players[playerID].Hand {
		options = append(options, Option{ID: card.ID, CardID: card.ID, Label: "Discard " + m.cardDefinition(card).Name})
	}
	m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptDiscardDown, OwnerID: playerID, Message: "Discard down to the 7-card hand limit.", Options: options}})
}

func (m *Match) switchTurn() {
	previous := m.CurrentPlayer
	next := m.opponent(previous)
	m.Players[previous].ActionsRemaining = 0
	m.CurrentPlayer = next
	m.Players[next].ActionsRemaining = 2
	m.endTurnQueued = false
	m.captureTurnStart(next)
	m.log("turn_started", fmt.Sprintf("%s begins a turn", m.Players[next].Name), map[string]any{"player_id": next})
}

func (m *Match) attackTargets(fighter *Fighter) []string {
	if fighter == nil || fighter.Defeated {
		return nil
	}
	result := []string{}
	attackType := m.fighterDefinition(fighter).AttackType
	for _, candidate := range m.Fighters {
		if candidate.OwnerID == fighter.OwnerID || candidate.Defeated {
			continue
		}
		legal := m.adjacent(fighter.SpaceID, candidate.SpaceID)
		if attackType == content.AttackRanged {
			legal = m.zonesOverlap(fighter.SpaceID, candidate.SpaceID)
		}
		if legal {
			result = append(result, candidate.ID)
		}
	}
	sort.Strings(result)
	return result
}
