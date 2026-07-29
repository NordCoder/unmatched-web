package game

import (
	"sort"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

func (m *Match) Project(playerID string) (View, error) {
	viewer, ok := m.Players[playerID]
	if !ok {
		return View{}, reject("unauthorized", "principal is not bound to this match")
	}
	view := View{
		MatchID: m.ID, Revision: m.Revision, Phase: m.Phase, BattlefieldID: m.Registry.Battlefield.ID, BattlefieldName: m.Registry.Battlefield.DisplayName,
		CurrentPlayerID: m.CurrentPlayer, ViewingPlayerID: playerID, WinnerID: m.WinnerID, Events: append([]Event(nil), m.Events...),
	}
	if m.WinnerID != "" {
		view.WinnerName = m.Players[m.WinnerID].Name
	}
	for _, id := range m.PlayerOrder {
		player := m.Players[id]
		hero := m.Fighters[player.HeroID]
		health := 0
		if hero != nil {
			health = hero.Health
		}
		pv := PlayerView{ID: id, Name: player.Name, Seat: player.Seat, DeckDefinitionID: player.DeckDefinitionID, HeroID: player.HeroID, Health: health, DeckCount: len(player.Deck), HandCount: len(player.Hand), DiscardCount: len(player.Discard), ActionsRemaining: player.ActionsRemaining}
		if id == playerID {
			for _, card := range player.Hand {
				pv.Hand = append(pv.Hand, m.cardView(card))
			}
			sort.Slice(pv.Hand, func(i, j int) bool {
				if pv.Hand[i].Type != pv.Hand[j].Type {
					return pv.Hand[i].Type < pv.Hand[j].Type
				}
				return pv.Hand[i].Name < pv.Hand[j].Name
			})
		}
		view.Players = append(view.Players, pv)
	}
	fighterIDs := make([]string, 0, len(m.Fighters))
	for fighterID := range m.Fighters {
		fighterIDs = append(fighterIDs, fighterID)
	}
	sort.Strings(fighterIDs)
	for _, fighterID := range fighterIDs {
		view.Fighters = append(view.Fighters, m.fighterView(m.Fighters[fighterID]))
	}
	for _, space := range m.Registry.Battlefield.Spaces {
		sv := SpaceView{ID: space.ID, X: space.X, Y: space.Y, Zones: append([]string(nil), space.Zones...)}
		if fighter := m.occupied(space.ID); fighter != nil {
			fv := m.fighterView(fighter)
			sv.Fighter = &fv
		}
		view.Spaces = append(view.Spaces, sv)
	}
	for _, edge := range m.Registry.Battlefield.Edges {
		view.Edges = append(view.Edges, EdgeView{From: edge.From, To: edge.To})
	}
	if len(m.Pending) > 0 {
		prompt := m.Pending[0].Prompt
		pv := &PendingView{Kind: prompt.Kind, OwnerID: prompt.OwnerID, OwnerName: m.Players[prompt.OwnerID].Name, Message: prompt.Message}
		if prompt.OwnerID == playerID {
			pv.Options = append([]Option(nil), prompt.Options...)
		}
		view.Pending = pv
	}
	if m.Combat != nil {
		cv := &CombatView{AttackerID: m.Combat.AttackerID, DefenderID: m.Combat.DefenderID, WaitingForDefense: true}
		if m.Combat.AttackCardRevealed {
			card := m.cardView(m.Combat.AttackCard)
			cv.AttackCard = &card
		}
		if m.Combat.DefenseCard != nil && m.Combat.DefenseCardRevealed {
			card := m.cardView(*m.Combat.DefenseCard)
			cv.DefenseCard = &card
		}
		view.Combat = cv
	}
	view.Legal = m.legalView(viewer)
	return view, nil
}

func (m *Match) cardView(card CardInstance) CardView {
	definition := m.cardDefinition(card)
	return CardView{ID: card.ID, DefinitionID: definition.ID, Effect: definition.Effect, Name: definition.Name, Type: definition.Type, Value: definition.Value, Boost: definition.Boost, UsableBy: append([]string(nil), definition.UsableBy...)}
}

func (m *Match) fighterView(fighter *Fighter) FighterView {
	definition := m.fighterDefinition(fighter)
	return FighterView{ID: fighter.ID, DefinitionID: fighter.DefinitionID, Name: definition.Name, OwnerID: fighter.OwnerID, Health: fighter.Health, MaxHealth: fighter.MaxHealth, SpaceID: fighter.SpaceID, Defeated: fighter.Defeated, AttackType: definition.AttackType}
}

func (m *Match) legalView(player *Player) LegalView {
	result := LegalView{
		SchemeActions: map[string]SchemeActionView{},
		AttackCards:   map[string][]string{},
		AttackTargets: map[string][]string{},
	}
	if m.Phase != PhaseActive || len(m.Pending) > 0 || m.CurrentPlayer != player.ID || player.ActionsRemaining <= 0 {
		return result
	}
	result.CanManeuver = true
	for _, card := range player.Hand {
		definition := m.cardDefinition(card)
		if definition.Type != content.CardScheme {
			continue
		}
		result.SchemeCards = append(result.SchemeCards, card.ID)
		result.SchemeActions[card.ID] = m.schemeActionView(player, definition)
	}
	for _, fighterID := range player.FighterIDs {
		fighter := m.Fighters[fighterID]
		if fighter.Defeated {
			continue
		}
		targets := m.attackTargets(fighter)
		if len(targets) == 0 {
			continue
		}
		for _, card := range player.Hand {
			definition := m.cardDefinition(card)
			if (definition.Type == content.CardAttack || definition.Type == content.CardVersatile) && usableBy(definition, fighter) {
				result.AttackCards[fighterID] = append(result.AttackCards[fighterID], card.ID)
			}
		}
		if len(result.AttackCards[fighterID]) > 0 {
			result.AttackTargets[fighterID] = targets
		}
	}
	sort.Strings(result.SchemeCards)
	return result
}

func (m *Match) schemeActionView(player *Player, card content.CardDefinition) SchemeActionView {
	result := SchemeActionView{}
	if card.Effect != content.EffectMoveThroughFive && card.Effect != content.EffectHorns {
		return result
	}
	for _, fighterID := range player.FighterIDs {
		fighter := m.Fighters[fighterID]
		if fighter == nil || fighter.Defeated || !usableBy(card, fighter) {
			continue
		}
		candidate := SchemeFighterActionView{FighterID: fighter.ID}
		for _, destination := range m.destinations(fighter, 5, true, false) {
			destination.Path = append([]string(nil), destination.Path...)
			candidate.Destinations = append(candidate.Destinations, destination)
			if card.Effect != content.EffectHorns {
				continue
			}
			targets := m.adjacentLivingFighters(destination.Destination, fighter.ID)
			if len(targets) == 0 {
				continue
			}
			if candidate.TargetsByDestination == nil {
				candidate.TargetsByDestination = map[string][]string{}
			}
			for _, target := range targets {
				candidate.TargetsByDestination[destination.Destination] = append(candidate.TargetsByDestination[destination.Destination], target.ID)
			}
		}
		result.Fighters = append(result.Fighters, candidate)
	}
	sort.Slice(result.Fighters, func(i, j int) bool { return result.Fighters[i].FighterID < result.Fighters[j].FighterID })
	return result
}
