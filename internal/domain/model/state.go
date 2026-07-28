// Package model contains Lead-owned cross-line runtime state contracts.
//
// The model is intentionally broad. Core Runtime owns hosting and persistence;
// Rules Mechanics owns deterministic interpretation. Neither worker should add
// fighter- or card-specific fields here without a Lead contract decision.
package model

import "encoding/json"

type (
	MatchID       string
	PrincipalID   string
	PlayerID      string
	FighterID     string
	CardID        string
	ComponentID   string
	ActionID      string
	CombatID      string
	ProcedureID   string
	InteractionID string
	EventID       string
	CommandID     string
	DefinitionID  string
)

type DefinitionRef struct {
	RulesetVersion           string            `json:"ruleset_version"`
	CapabilityRegistry       string            `json:"capability_registry_version"`
	FighterManifestDigests   map[string]string `json:"fighter_manifest_digests,omitempty"`
	CardManifestDigests      map[string]string `json:"card_manifest_digests,omitempty"`
	BattlefieldManifestDigest string           `json:"battlefield_manifest_digest,omitempty"`
	SetupDefinitionDigest    string            `json:"setup_definition_digest,omitempty"`
}

type Lifecycle string

const (
	LifecycleCreated           Lifecycle = "CREATED"
	LifecycleWaitingForPlayers Lifecycle = "WAITING_FOR_PLAYERS"
	LifecycleSelection         Lifecycle = "SELECTION"
	LifecycleSetup             Lifecycle = "SETUP"
	LifecycleActive            Lifecycle = "ACTIVE"
	LifecycleEnded             Lifecycle = "ENDED"
	LifecycleQuarantined       Lifecycle = "QUARANTINED"
)

type PlayerState struct {
	ID                    PlayerID           `json:"player_instance_id"`
	Seat                  int                `json:"seat"`
	AuthorityState        string             `json:"authority_state"`
	FighterInstanceIDs    []FighterID        `json:"fighter_instance_ids,omitempty"`
	PrivateZones          map[string][]CardID `json:"private_zones,omitempty"`
	Resources             map[string]int     `json:"resources,omitempty"`
	ActionPermissionIDs   []string           `json:"action_permission_ids,omitempty"`
	SubmittedChoiceIDs    []string           `json:"submitted_hidden_choice_ids,omitempty"`
}

type RuntimeObject struct {
	DefinitionID DefinitionID       `json:"definition_id"`
	OwnerID      PlayerID           `json:"owner_player_id,omitempty"`
	ControllerID PlayerID           `json:"controller_player_id,omitempty"`
	State        map[string]any     `json:"state,omitempty"`
}

type ActionState struct {
	ID       ActionID       `json:"action_instance_id"`
	Kind     string         `json:"kind"`
	ActorID  PlayerID       `json:"actor_player_id"`
	Status   string         `json:"status"`
	Procedure *ProcedureRef `json:"procedure,omitempty"`
}

type CombatState struct {
	ID             CombatID       `json:"combat_instance_id"`
	AttackerID     FighterID      `json:"attacker_id"`
	DefenderID     FighterID      `json:"defender_id"`
	Stage          string         `json:"stage"`
	AttackCardID   CardID         `json:"attack_card_id,omitempty"`
	DefenseCardID  CardID         `json:"defense_card_id,omitempty"`
	Procedure      *ProcedureRef   `json:"procedure,omitempty"`
}

type ProcedureRef struct {
	ID         ProcedureID      `json:"procedure_instance_id"`
	Kind       string           `json:"kind"`
	Stage      string           `json:"stage"`
	SourceRef  string           `json:"source_ref,omitempty"`
	Bindings   map[string]json.RawMessage `json:"bindings,omitempty"`
}

type PendingInteraction struct {
	ID             InteractionID   `json:"interaction_instance_id"`
	OwnerPlayerID  PlayerID        `json:"owner_player_id"`
	Kind           string          `json:"kind"`
	Visibility     string          `json:"visibility"`
	Prompt         json.RawMessage `json:"prompt"`
	LegalDomain    json.RawMessage `json:"legal_domain"`
	ResumeProcedure ProcedureRef   `json:"resume_procedure"`
}

type ResolverState struct {
	ActiveProcedure   *ProcedureRef       `json:"active_procedure,omitempty"`
	PendingInteraction *PendingInteraction `json:"pending_interaction,omitempty"`
	Queue             []ProcedureRef      `json:"queue,omitempty"`
	History           []string            `json:"history,omitempty"`
}

type GameState struct {
	MatchID       MatchID                    `json:"match_id"`
	DefinitionRef DefinitionRef              `json:"definition_ref"`
	Revision      uint64                     `json:"revision"`
	EventSequence uint64                     `json:"event_sequence"`
	Lifecycle     Lifecycle                  `json:"lifecycle"`
	Players       map[PlayerID]PlayerState   `json:"players"`
	Fighters      map[FighterID]RuntimeObject `json:"fighters,omitempty"`
	Cards         map[CardID]RuntimeObject    `json:"cards,omitempty"`
	Components    map[ComponentID]RuntimeObject `json:"components,omitempty"`
	Battlefield   map[string]any             `json:"battlefield,omitempty"`
	Turn          map[string]any             `json:"turn,omitempty"`
	Action        *ActionState               `json:"action,omitempty"`
	Combat        *CombatState               `json:"combat,omitempty"`
	Resolver      ResolverState              `json:"resolver"`
	Random        map[string]any             `json:"random,omitempty"`
	GameResult    map[string]any             `json:"game_result,omitempty"`
}
