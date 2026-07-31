package game

import (
	"fmt"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func (m *Match) resolveCombat() error {
	if m.Combat == nil {
		return reject("combat_missing", "there is no active combat")
	}
	combat := m.Combat
	attacker := m.Fighters[combat.AttackerID]
	defender := m.Fighters[combat.DefenderID]
	attackDef := m.cardDefinition(combat.AttackCard)
	var defenseDef content.CardDefinition
	if combat.DefenseCard != nil {
		defenseDef = m.cardDefinition(*combat.DefenseCard)
	}

	attackEffectsCancelled := false
	defenseEffectsCancelled := false
	ignoreDefenseValue := false
	// The defender resolves immediate effects first at the same timing window.
	if combat.DefenseCard != nil && defenseDef.Effect == content.EffectCancelOpposing {
		attackEffectsCancelled = true
	}
	if !attackEffectsCancelled {
		switch attackDef.Effect {
		case content.EffectCancelOpposing:
			defenseEffectsCancelled = true
		case content.EffectCancelIgnoreValue:
			defenseEffectsCancelled = true
			ignoreDefenseValue = true
		}
	}

	attackValue := attackDef.Value
	defenseValue := 0
	if combat.DefenseCard != nil && !ignoreDefenseValue {
		defenseValue = defenseDef.Value
	}
	if !attackEffectsCancelled {
		attackValue = m.duringCombatValue(combat.AttackerPlayerID, attacker, combat.AttackCard, attackValue)
	}
	if combat.DefenseCard != nil && !defenseEffectsCancelled {
		defenseValue = m.duringCombatValue(combat.DefenderPlayerID, defender, *combat.DefenseCard, defenseValue)
	}
	if !attackEffectsCancelled && attackDef.Effect == content.EffectAmbush {
		attackValue += m.resolveAmbush(combat.DefenderPlayerID)
	}

	combat.AttackCardRevealed = true
	combat.DefenseCardRevealed = true
	message := fmt.Sprintf("%s revealed %s (%d)", m.fighterName(attacker), attackDef.Name, attackValue)
	if combat.DefenseCard != nil {
		message += fmt.Sprintf(" against %s (%d)", defenseDef.Name, defenseValue)
	} else {
		message += " with no defense"
	}
	m.log("combat_revealed", message, map[string]any{"attack_card": attackDef.Name, "attack_value": attackValue, "defense_card": defenseDef.Name, "defense_value": defenseValue})

	damage := attackValue - defenseValue
	if damage < 0 {
		damage = 0
	}
	dealt := m.damage(defender.ID, damage, "combat")
	if dealt > 0 {
		m.log("combat_damage", fmt.Sprintf("%s took %d combat damage", m.fighterName(defender), dealt), map[string]any{"fighter_id": defender.ID, "amount": dealt})
	}

	// Combat cards always move to their owners' discard piles after reveal.
	m.Players[combat.AttackerPlayerID].Discard = append(m.Players[combat.AttackerPlayerID].Discard, combat.AttackCard)
	if combat.DefenseCard != nil {
		m.Players[combat.DefenderPlayerID].Discard = append(m.Players[combat.DefenderPlayerID].Discard, *combat.DefenseCard)
	}
	m.Combat = nil
	if m.Phase == PhaseEnded {
		return nil
	}

	attackerWon := attackValue > defenseValue
	defenderWon := defenseValue >= attackValue
	if combat.DefenseCard != nil && !defenseEffectsCancelled {
		m.applyAfterCombat(combat.DefenderPlayerID, defender, attacker, *combat.DefenseCard, defenderWon, 0)
	}
	if m.Phase == PhaseEnded {
		return nil
	}
	if !attackEffectsCancelled {
		m.applyAfterCombat(combat.AttackerPlayerID, attacker, defender, combat.AttackCard, attackerWon, dealt)
	}
	if m.Phase == PhaseEnded {
		return nil
	}
	attackerDeck := m.Registry.Decks[m.Players[combat.AttackerPlayerID].DeckDefinitionID]
	if attackerDeck.Ability == content.AbilityPostAttackRetreat && !attacker.Defeated {
		options := m.destinations(attacker, 2, false, false)
		for i := range options {
			options[i].ID = "retreat:" + options[i].ID
			options[i].Label = "Retreat: " + options[i].Label
		}
		m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptMove, OwnerID: combat.AttackerPlayerID, Message: "After attacking, you may move the attacking fighter up to 2 spaces.", Options: options}})
	}
	if len(m.Pending) == 0 {
		m.finishAction()
	}
	return nil
}

func (m *Match) duringCombatValue(playerID string, fighter *Fighter, card CardInstance, current int) int {
	definition := m.cardDefinition(card)
	if definition.Effect != content.EffectMomentousShift {
		return current
	}
	if m.TurnStartSpaces != nil && m.TurnStartSpaces[fighter.ID] != "" && m.TurnStartSpaces[fighter.ID] != fighter.SpaceID {
		return 5
	}
	return current
}

func (m *Match) resolveAmbush(opponentID string) int {
	hand := m.Players[opponentID].Hand
	if len(hand) == 0 {
		return 0
	}
	index := int(m.random() % uint64(len(hand)))
	card := hand[index]
	m.Players[opponentID].Hand = append(hand[:index:index], hand[index+1:]...)
	m.Players[opponentID].Discard = append(m.Players[opponentID].Discard, card)
	definition := m.cardDefinition(card)
	m.log("random_discard", fmt.Sprintf("%s discarded %s at random", m.Players[opponentID].Name, definition.Name), map[string]any{"card": definition.Name, "boost": definition.Boost})
	return definition.Boost
}

func (m *Match) applyAfterCombat(playerID string, own, opposing *Fighter, card CardInstance, won bool, damageToOpposing int) {
	definition := m.cardDefinition(card)
	switch definition.Effect {
	case content.EffectDrawDamage:
		m.draw(playerID, damageToOpposing)
	case content.EffectDrawTwo:
		m.draw(playerID, 2)
	case content.EffectDrawReturnSidekick:
		m.draw(playerID, 1)
		m.queueReturnSidekick(playerID)
	case content.EffectRegroup:
		if won {
			m.draw(playerID, 2)
		} else {
			m.draw(playerID, 1)
		}
	case content.EffectAdjacentOpponentDamage:
		m.damageAdjacentOpponents(playerID, own, 1, "card_effect")
	case content.EffectDrawIfAdjacent:
		if own != nil && opposing != nil && !own.Defeated && !opposing.Defeated && m.adjacent(own.SpaceID, opposing.SpaceID) {
			m.draw(playerID, 1)
		}
	case content.EffectSavagery:
		if won {
			m.damageAdjacentOpponents(playerID, own, 1, "card_effect")
		}
	case content.EffectPostCombatMoveFive:
		if own != nil && !own.Defeated {
			options := m.destinations(own, 5, true, false)
			for i := range options {
				options[i].ID = "post-move:" + options[i].ID
				options[i].Label = "Move: " + options[i].Label
			}
			m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptMove, OwnerID: playerID, Message: "Move your combat fighter up to 5 spaces; it may move through opposing fighters.", Options: options}})
		}
	case content.EffectPostCombatPlaceZone:
		m.queuePlaceInZone(playerID, own)
	case content.EffectSkirmish:
		if won {
			m.queueSkirmish(playerID, own, opposing)
		}
	}
	if m.Phase == PhaseEnded {
		m.Pending = nil
	}
}

func (m *Match) damageAdjacentOpponents(playerID string, source *Fighter, amount int, reason string) {
	if source == nil || source.Defeated {
		return
	}
	targets := []string{}
	for _, fighter := range m.Fighters {
		if fighter.OwnerID != playerID && !fighter.Defeated && m.adjacent(source.SpaceID, fighter.SpaceID) {
			targets = append(targets, fighter.ID)
		}
	}
	sort.Strings(targets)
	for _, fighterID := range targets {
		m.damage(fighterID, amount, reason)
		if m.Phase == PhaseEnded {
			return
		}
	}
}

func (m *Match) queueReturnSidekick(playerID string) {
	player := m.Players[playerID]
	hero := m.Fighters[player.HeroID]
	if hero == nil || hero.Defeated {
		return
	}
	defeated := []*Fighter{}
	deck := m.Registry.Decks[player.DeckDefinitionID]
	for _, fighterID := range player.FighterIDs {
		fighter := m.Fighters[fighterID]
		if fighter.Defeated && deck.Fighters[fighter.DefinitionID].Role == "sidekick" {
			defeated = append(defeated, fighter)
		}
	}
	if len(defeated) == 0 {
		return
	}
	heroSpace, _ := m.space(hero.SpaceID)
	empty := []string{}
	for _, space := range m.Registry.Battlefield.Spaces {
		if sharesZone(heroSpace, space) && m.occupied(space.ID) == nil {
			empty = append(empty, space.ID)
		}
	}
	sort.Strings(empty)
	options := []Option{}
	for _, fighter := range defeated {
		for _, spaceID := range empty {
			options = append(options, Option{ID: "return:" + fighter.ID + ":" + spaceID, Label: "Return " + m.fighterName(fighter) + " at " + spaceID, FighterID: fighter.ID, Destination: spaceID})
		}
	}
	m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptReturnSidekick, OwnerID: playerID, Message: "Return one defeated sidekick to an empty space in your hero's zone?", Options: options}})
}

func (m *Match) queuePlaceInZone(playerID string, fighter *Fighter) {
	if fighter == nil || fighter.Defeated {
		return
	}
	origin, _ := m.space(fighter.SpaceID)
	options := []Option{}
	for _, space := range m.Registry.Battlefield.Spaces {
		if space.ID != fighter.SpaceID && sharesZone(origin, space) && m.occupied(space.ID) == nil {
			options = append(options, Option{ID: "place:" + fighter.ID + ":" + space.ID, Label: "Place at " + space.ID, FighterID: fighter.ID, Destination: space.ID})
		}
	}
	sort.Slice(options, func(i, j int) bool { return options[i].Destination < options[j].Destination })
	if len(options) > 0 {
		m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptPlace, OwnerID: playerID, Message: "Place your combat fighter in an empty space in its zone.", Options: options}})
	}
}

func (m *Match) queueSkirmish(playerID string, first, second *Fighter) {
	options := []Option{}
	for _, fighter := range []*Fighter{first, second} {
		if fighter == nil || fighter.Defeated {
			continue
		}
		for _, option := range m.destinations(fighter, 2, false, false) {
			option.ID = "skirmish:" + option.ID
			option.Label = "Move " + m.fighterName(fighter) + " to " + option.Destination
			options = append(options, option)
		}
	}
	if len(options) > 0 {
		m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptSkirmish, OwnerID: playerID, Message: "Move either fighter in the combat up to 2 spaces.", Options: dedupeOptions(options)}})
	}
}
