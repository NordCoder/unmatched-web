package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

var (
	ErrRevisionConflict  = errors.New("revision conflict")
	ErrCommandConflict   = errors.New("command identity conflict")
	ErrAuthorityConflict = errors.New("authority binding conflict")
	ErrInvalidLease      = errors.New("invalid command lease")
)

const (
	AuthorityActive       = "ACTIVE"
	FingerprintSchemaV1   = "core-command-fingerprint/v1"
	CommandResultSchemaV1 = "core-command-result/v1"
)

type CommandKey struct {
	PrincipalID model.PrincipalID
	CommandID   model.CommandID
}

type CommandScope string

const (
	CommandScopeCreateMatch  CommandScope = "create_match"
	CommandScopeJoinMatch    CommandScope = "join_match"
	CommandScopeExistingSeat CommandScope = "existing_seat"
)

type CommandIdentity struct {
	Key           CommandKey
	Fingerprint   []byte
	MatchID       model.MatchID
	ActorPlayerID model.PlayerID
	Scope         CommandScope
}

type CommandRecord struct {
	Fingerprint  []byte
	Result       []byte
	MatchID      model.MatchID
	ErrorCode    string
	ErrorMessage string
}

type CommandLease struct {
	key           CommandKey
	token         string
	matchID       model.MatchID
	actorPlayerID model.PlayerID
	scope         CommandScope
	backend       any
}

type AuthorityRecord struct {
	MatchID                model.MatchID
	PlayerID               model.PlayerID
	PrincipalID            model.PrincipalID
	Seat                   int
	BindingVersion         uint64
	Status                 string
	EstablishedByCommandID model.CommandID
}

type CommitRequest struct {
	MatchID       model.MatchID
	DefinitionRef model.DefinitionRef
	Batch         contracts.EventBatch
	Result        []byte
	Authority     *AuthorityRecord
}

type EventStore interface {
	AcquireCommand(ctx context.Context, identity CommandIdentity) (lease CommandLease, record CommandRecord, duplicate bool, err error)
	Commit(ctx context.Context, lease CommandLease, request CommitRequest) (CommandRecord, error)
	RejectCommand(ctx context.Context, lease CommandLease, matchID model.MatchID, code, message string) (CommandRecord, error)
	AbortCommand(ctx context.Context, lease CommandLease) error
	LoadEvents(ctx context.Context, matchID model.MatchID) ([]contracts.DomainEvent, error)
	ResolveAuthorityContext(ctx context.Context, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error)
	LoadEventsForCommand(ctx context.Context, lease CommandLease, matchID model.MatchID) ([]contracts.DomainEvent, error)
	ResolveAuthorityForCommand(ctx context.Context, lease CommandLease, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error)
}

func cloneCommandRecord(record CommandRecord) CommandRecord {
	record.Fingerprint = append([]byte(nil), record.Fingerprint...)
	record.Result = append([]byte(nil), record.Result...)
	return record
}

func cloneAuthorityRecord(record AuthorityRecord) AuthorityRecord { return record }

func cloneEvents(events []contracts.DomainEvent) []contracts.DomainEvent {
	result := make([]contracts.DomainEvent, len(events))
	for index, event := range events {
		event.PublicPayload = append([]byte(nil), event.PublicPayload...)
		if event.PrivatePayloads != nil {
			payloads := make(map[model.PlayerID]json.RawMessage, len(event.PrivatePayloads))
			for playerID, payload := range event.PrivatePayloads {
				payloads[playerID] = append(json.RawMessage(nil), payload...)
			}
			event.PrivatePayloads = payloads
		}
		result[index] = event
	}
	return result
}
