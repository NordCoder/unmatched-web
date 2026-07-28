// Package contracts defines the Lead-owned boundary between Core Runtime and
// Rules Mechanics. The contracts are intentionally small and may evolve through
// explicit Lead decisions while implementation proceeds.
package contracts

import (
	"encoding/json"

	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

type Command struct {
	ID               model.CommandID `json:"command_id"`
	SchemaVersion    string          `json:"command_schema_version"`
	Type             string          `json:"type"`
	MatchID          model.MatchID   `json:"match_id,omitempty"`
	ActorPlayerID    model.PlayerID  `json:"actor_player_id,omitempty"`
	ExpectedRevision *uint64         `json:"expected_revision,omitempty"`
	Payload          json.RawMessage `json:"payload"`
}

type DomainEvent struct {
	SchemaVersion   string                             `json:"event_schema_version"`
	ID              model.EventID                      `json:"event_id"`
	MatchID         model.MatchID                      `json:"match_id"`
	Sequence        uint64                             `json:"sequence"`
	Revision        uint64                             `json:"revision"`
	Type            string                             `json:"event_type"`
	CausedByCommand model.CommandID                    `json:"caused_by_command_id"`
	ParentEventID   model.EventID                      `json:"parent_event_id,omitempty"`
	SourceRef       string                             `json:"source_ref,omitempty"`
	RulesetVersion  string                             `json:"ruleset_version"`
	PublicPayload   json.RawMessage                    `json:"public_payload"`
	PrivatePayloads map[model.PlayerID]json.RawMessage `json:"private_payloads_by_player"`
}

type EventBatch struct {
	CommandID        model.CommandID `json:"command_id"`
	PreviousRevision uint64          `json:"previous_revision"`
	NextRevision     uint64          `json:"next_revision"`
	Events           []DomainEvent   `json:"events"`
}

type ResolutionInput struct {
	CommandID     model.CommandID            `json:"command_id"`
	Procedure     model.ProcedureRef         `json:"procedure"`
	InteractionID model.InteractionID        `json:"interaction_id,omitempty"`
	Choice        json.RawMessage            `json:"choice,omitempty"`
	Context       map[string]json.RawMessage `json:"context,omitempty"`
}

type ResolutionStatus string

const (
	ResolutionCompleted ResolutionStatus = "COMPLETED"
	ResolutionPending   ResolutionStatus = "PENDING_INTERACTION"
	ResolutionRejected  ResolutionStatus = "REJECTED"
)

type ResolutionOutcome struct {
	Status             ResolutionStatus          `json:"status"`
	Events             []DomainEvent             `json:"events,omitempty"`
	PendingInteraction *model.PendingInteraction `json:"pending_interaction,omitempty"`
	Procedure          *model.ProcedureRef        `json:"procedure,omitempty"`
	RejectionCode      string                     `json:"rejection_code,omitempty"`
	Diagnostics        map[string]string          `json:"diagnostics,omitempty"`
}

// RulesEngine is pure from Core Runtime's perspective: it receives an immutable
// state value plus a serializable procedure input and returns explicit events or
// a pending interaction. Persistence and transport never enter this boundary.
type RulesEngine interface {
	Resolve(state model.GameState, input ResolutionInput) (ResolutionOutcome, error)
	LegalActions(state model.GameState, playerID model.PlayerID) ([]json.RawMessage, error)
	Project(state model.GameState, playerID model.PlayerID) (json.RawMessage, error)
}
