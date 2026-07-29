package game

import (
	"fmt"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func (m *Match) playScheme(playerID string, command Command) error {
	player, err := m.requireCurrent(playerID)
	if err != nil {
		return err
	}
	card, ok := findCard(player.Hand, command.CardID)
	if !ok {
		return reject("card_not_in_hand", "scheme card is not in the player's hand")
	}
	definition := m.cardDefinition(card)
	if definition.Type != content.CardScheme {
		return reject("illegal_card_type", "%s is not a scheme", definition.Name)
	}
	fighter, err := m.resolveUsableFighter(playerID, command.FighterID, definition)
	if err != nil {
		return err
	}

	var hornsTarget *Fighter
	switch definition.Effect {
	case content.EffectStealFromRich:
		// No additional target is required in this two-player slice.
	case content.EffectMoveThroughFive:
		if err := m.validatePath(fighter, command.Path, 5, true); err != nil {
			return err
		}
	case content.EffectHorns:
		if err := m.validatePath(fighter, command.Path, 5, true); err != nil {
			return err
		}
		finalSpace := fighter.SpaceID
		if len(command.Path) > 0 {
			finalSpace = command.Path[len(command.Path)-1]
		}
		hornsTarget, err = m.validateHornsTarget(fighter, finalSpace, command.TargetID)
		if err != nil {
			return err
		}
	default:
		return reject("unsupported_effect", "scheme %s has unsupported typed effect %q", definition.Name, definition.Effect)
	}

	player.Hand, card, _ = removeCard(player.Hand, card.ID)
	player.Discard = append(player.Discard, card)
	player.ActionsRemaining--
	m.log("scheme_played", fmt.Sprintf("%s played %s", player.Name, definition.Name), map[string]any{"player_id": playerID, "card": definition.Name})
	switch definition.Effect {
	case content.EffectStealFromRich:
		m.draw(playerID, 1)
		if m.Phase == PhaseEnded {
			return nil
		}
		opponentID := m.opponent(playerID)
		options := []Option{{ID: "decline", Label: "Decline; opponent draws one card"}}
		for _, opponentCard := range m.Players[opponentID].Hand {
			cardDef := m.cardDefinition(opponentCard)
			options = append(options, Option{ID: opponentCard.ID, CardID: opponentCard.ID, Label: "Discard " + cardDef.Name})
		}
		m.Pending = append(m.Pending, pendingTask{Prompt: Prompt{Kind: PromptStealDiscard, OwnerID: opponentID, Message: "Discard one card, or decline and let the opponent draw one.", Options: options}, data: map[string]string{"scheme_owner": playerID}})
	case content.EffectMoveThroughFive:
		m.move(fighter, command.Path)
	case content.EffectHorns:
		m.move(fighter, command.Path)
		if hornsTarget != nil {
			m.damage(hornsTarget.ID, 2, "scheme")
		}
	}
	if len(m.Pending) == 0 {
		m.finishAction()
	}
	return nil
}

func (m *Match) validateHornsTarget(jackalope *Fighter, finalSpace, requestedTargetID string) (*Fighter, error) {
	candidates := m.adjacentLivingFighters(finalSpace, jackalope.ID)
	if len(candidates) == 0 {
		if requestedTargetID != "" {
			return nil, reject("illegal_target", "Jackalope Horns has no adjacent fighter at the destination")
		}
		return nil, nil
	}
	if requestedTargetID == "" {
		return nil, reject("illegal_target", "Jackalope Horns requires an adjacent fighter target")
	}
	for _, candidate := range candidates {
		if candidate.ID == requestedTargetID {
			return candidate, nil
		}
	}
	return nil, reject("illegal_target", "selected fighter is not adjacent to Jackalope's destination")
}

func (m *Match) adjacentLivingFighters(spaceID, excludedFighterID string) []*Fighter {
	result := make([]*Fighter, 0)
	for _, fighter := range m.Fighters {
		if fighter.ID == excludedFighterID || fighter.Defeated {
			continue
		}
		if m.adjacent(spaceID, fighter.SpaceID) {
			result = append(result, fighter)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (m *Match) resolveUsableFighter(playerID, requested string, card content.CardDefinition) (*Fighter, error) {
	if requested != "" {
		fighter := m.Fighters[requested]
		if fighter == nil || fighter.OwnerID != playerID || fighter.Defeated {
			return nil, reject("invalid_fighter", "selected fighter is not available")
		}
		if !usableBy(card, fighter) {
			return nil, reject("fighter_cannot_use_card", "%s cannot use %s", m.fighterName(fighter), card.Name)
		}
		return fighter, nil
	}
	for _, fighterID := range m.Players[playerID].FighterIDs {
		fighter := m.Fighters[fighterID]
		if !fighter.Defeated && usableBy(card, fighter) {
			return fighter, nil
		}
	}
	return nil, reject("fighter_cannot_use_card", "no living fighter can use %s", card.Name)
}
