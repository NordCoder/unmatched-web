// Package game implements the typed ordinary two-player Unmatched kernel used
// by the launch playable slice. It depends on semantic definitions, never on
// fighter or card identity dispatch.
package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
)

type Phase string

const (
	PhaseWaiting Phase = "waiting_for_player"
	PhaseActive  Phase = "active"
	PhaseEnded   Phase = "ended"
)

type CommandType string

const (
	CommandManeuver CommandType = "maneuver"
	CommandScheme   CommandType = "scheme"
	CommandAttack   CommandType = "attack"
	CommandChoose   CommandType = "choose"
)

type Command struct {
	Type             CommandType `json:"type"`
	ExpectedRevision uint64      `json:"expected_revision"`
	FighterID        string      `json:"fighter_id,omitempty"`
	TargetID         string      `json:"target_id,omitempty"`
	SecondaryID      string      `json:"secondary_id,omitempty"`
	CardID           string      `json:"card_id,omitempty"`
	BoostCardID      string      `json:"boost_card_id,omitempty"`
	Path             []string    `json:"path,omitempty"`
	SecondaryPath    []string    `json:"secondary_path,omitempty"`
	Choice           string      `json:"choice,omitempty"`
}

type RuleError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *RuleError) Error() string { return e.Code + ": " + e.Message }

func reject(code, format string, args ...any) error {
	return &RuleError{Code: code, Message: fmt.Sprintf(format, args...)}
}

func ErrorCode(err error) string {
	var rule *RuleError
	if errors.As(err, &rule) {
		return rule.Code
	}
	return "internal_error"
}

type CardInstance struct {
	ID           string `json:"id"`
	DefinitionID string `json:"definition_id"`
}

type Fighter struct {
	ID           string `json:"id"`
	DefinitionID string `json:"definition_id"`
	OwnerID      string `json:"owner_id"`
	Health       int    `json:"health"`
	MaxHealth    int    `json:"max_health"`
	SpaceID      string `json:"space_id,omitempty"`
	Defeated     bool   `json:"defeated"`
}

type Player struct {
	ID               string
	Seat             int
	Name             string
	DeckDefinitionID string
	HeroID           string
	FighterIDs       []string
	Deck             []CardInstance
	Hand             []CardInstance
	Discard          []CardInstance
	ActionsRemaining int
	TurnStartSpaces  map[string]string
}

type PromptKind string

const (
	PromptManeuverBoost  PromptKind = "maneuver_boost"
	PromptManeuverMove   PromptKind = "maneuver_move"
	PromptDefense        PromptKind = "defense"
	PromptStealDiscard   PromptKind = "steal_discard"
	PromptMove           PromptKind = "move"
	PromptPlace          PromptKind = "place"
	PromptSkirmish       PromptKind = "skirmish"
	PromptReturnSidekick PromptKind = "return_sidekick"
	PromptIsolatedDraw   PromptKind = "isolated_draw"
	PromptDiscardDown    PromptKind = "discard_down"
)

type Option struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	FighterID   string   `json:"fighter_id,omitempty"`
	CardID      string   `json:"card_id,omitempty"`
	Destination string   `json:"destination,omitempty"`
	Path        []string `json:"path,omitempty"`
}

type Prompt struct {
	Kind    PromptKind `json:"kind"`
	OwnerID string     `json:"owner_id"`
	Message string     `json:"message"`
	Options []Option   `json:"options,omitempty"`
}

type pendingTask struct {
	Prompt
	data map[string]string
}

type Combat struct {
	AttackerID          string
	DefenderID          string
	AttackerPlayerID    string
	DefenderPlayerID    string
	AttackCard          CardInstance
	DefenseCard         *CardInstance
	AttackCardRevealed  bool
	DefenseCardRevealed bool
}

type Event struct {
	At      time.Time      `json:"at"`
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

type Match struct {
	ID              string
	Revision        uint64
	Phase           Phase
	Registry        *content.Registry
	Players         map[string]*Player
	PlayerOrder     []string
	Fighters        map[string]*Fighter
	CurrentPlayer   string
	Pending         []pendingTask
	Combat          *Combat
	WinnerID        string
	Events          []Event
	seed            uint64
	cardSequence    uint64
	endTurnQueued   bool
	TurnStartSpaces map[string]string
}

type CardView struct {
	ID           string             `json:"id"`
	DefinitionID string             `json:"definition_id"`
	Effect       content.EffectKind `json:"effect"`
	Name         string             `json:"name"`
	Type         content.CardType   `json:"type"`
	Value        int                `json:"value"`
	Boost        int                `json:"boost"`
	UsableBy     []string           `json:"usable_by"`
}

type FighterView struct {
	ID           string             `json:"id"`
	DefinitionID string             `json:"definition_id"`
	Name         string             `json:"name"`
	OwnerID      string             `json:"owner_id"`
	Health       int                `json:"health"`
	MaxHealth    int                `json:"max_health"`
	SpaceID      string             `json:"space_id,omitempty"`
	Defeated     bool               `json:"defeated"`
	AttackType   content.AttackType `json:"attack_type"`
}

type PlayerView struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Seat             int        `json:"seat"`
	DeckDefinitionID string     `json:"deck_definition_id"`
	HeroID           string     `json:"hero_id"`
	Health           int        `json:"hero_health"`
	DeckCount        int        `json:"deck_count"`
	HandCount        int        `json:"hand_count"`
	DiscardCount     int        `json:"discard_count"`
	ActionsRemaining int        `json:"actions_remaining"`
	Hand             []CardView `json:"hand,omitempty"`
}

type SpaceView struct {
	ID      string       `json:"id"`
	X       int          `json:"x"`
	Y       int          `json:"y"`
	Zones   []string     `json:"zones"`
	Fighter *FighterView `json:"fighter,omitempty"`
}

type EdgeView struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PendingView struct {
	Kind      PromptKind `json:"kind"`
	OwnerID   string     `json:"owner_id"`
	OwnerName string     `json:"owner_name"`
	Message   string     `json:"message"`
	Options   []Option   `json:"options,omitempty"`
}

type CombatView struct {
	AttackerID        string    `json:"attacker_id"`
	DefenderID        string    `json:"defender_id"`
	AttackCard        *CardView `json:"attack_card,omitempty"`
	DefenseCard       *CardView `json:"defense_card,omitempty"`
	WaitingForDefense bool      `json:"waiting_for_defense"`
}

type LegalView struct {
	CanManeuver   bool                `json:"can_maneuver"`
	SchemeCards   []string            `json:"scheme_cards,omitempty"`
	AttackCards   map[string][]string `json:"attack_cards_by_fighter,omitempty"`
	AttackTargets map[string][]string `json:"attack_targets_by_fighter,omitempty"`
}

type View struct {
	MatchID         string        `json:"match_id"`
	Revision        uint64        `json:"revision"`
	Phase           Phase         `json:"phase"`
	BattlefieldID   string        `json:"battlefield_id"`
	BattlefieldName string        `json:"battlefield_name"`
	CurrentPlayerID string        `json:"current_player_id,omitempty"`
	ViewingPlayerID string        `json:"viewing_player_id"`
	Players         []PlayerView  `json:"players"`
	Fighters        []FighterView `json:"fighters"`
	Spaces          []SpaceView   `json:"spaces"`
	Edges           []EdgeView    `json:"edges"`
	Pending         *PendingView  `json:"pending,omitempty"`
	Combat          *CombatView   `json:"combat,omitempty"`
	Legal           LegalView     `json:"legal"`
	WinnerID        string        `json:"winner_id,omitempty"`
	WinnerName      string        `json:"winner_name,omitempty"`
	Events          []Event       `json:"events"`
}
