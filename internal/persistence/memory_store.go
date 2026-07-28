package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

var (
	ErrRevisionConflict  = errors.New("revision conflict")
	ErrCommandConflict   = errors.New("command identity conflict")
	ErrAuthorityConflict = errors.New("authority binding conflict")
	ErrInvalidLease      = errors.New("invalid command lease")
)

const AuthorityActive = "ACTIVE"

type CommandKey struct {
	PrincipalID model.PrincipalID
	CommandID   model.CommandID
}

type CommandRecord struct {
	Fingerprint  []byte
	Result       []byte
	MatchID      model.MatchID
	ErrorCode    string
	ErrorMessage string
}

type CommandLease struct {
	key   CommandKey
	token uint64
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
	MatchID   model.MatchID
	Batch     contracts.EventBatch
	Result    []byte
	Authority *AuthorityRecord
}

type EventStore interface {
	AcquireCommand(ctx context.Context, key CommandKey, fingerprint []byte) (lease CommandLease, record CommandRecord, duplicate bool, err error)
	Commit(lease CommandLease, request CommitRequest) (CommandRecord, error)
	RejectCommand(lease CommandLease, matchID model.MatchID, code, message string) (CommandRecord, error)
	AbortCommand(lease CommandLease)
	Events(matchID model.MatchID) []contracts.DomainEvent
	ResolveAuthority(matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool)
}

type stream struct {
	revision uint64
	events   []contracts.DomainEvent
}

type commandEntry struct {
	fingerprint []byte
	record      *CommandRecord
	done        chan struct{}
	token       uint64
}

type authorityKey struct {
	matchID     model.MatchID
	principalID model.PrincipalID
}

type playerAuthorityKey struct {
	matchID  model.MatchID
	playerID model.PlayerID
}

type MemoryStore struct {
	mu                sync.RWMutex
	streams           map[model.MatchID]*stream
	commands          map[CommandKey]*commandEntry
	authorities       map[authorityKey]AuthorityRecord
	playerAuthorities map[playerAuthorityKey]model.PrincipalID
	nextToken         uint64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		streams:           make(map[model.MatchID]*stream),
		commands:          make(map[CommandKey]*commandEntry),
		authorities:       make(map[authorityKey]AuthorityRecord),
		playerAuthorities: make(map[playerAuthorityKey]model.PrincipalID),
	}
}

func (s *MemoryStore) AcquireCommand(ctx context.Context, key CommandKey, fingerprint []byte) (CommandLease, CommandRecord, bool, error) {
	if key.PrincipalID == "" || key.CommandID == "" || len(fingerprint) == 0 {
		return CommandLease{}, CommandRecord{}, false, ErrInvalidLease
	}
	for {
		s.mu.Lock()
		entry, exists := s.commands[key]
		if !exists {
			s.nextToken++
			entry = &commandEntry{
				fingerprint: append([]byte(nil), fingerprint...),
				done:        make(chan struct{}),
				token:       s.nextToken,
			}
			s.commands[key] = entry
			lease := CommandLease{key: key, token: entry.token}
			s.mu.Unlock()
			return lease, CommandRecord{}, false, nil
		}
		if !bytes.Equal(entry.fingerprint, fingerprint) {
			s.mu.Unlock()
			return CommandLease{}, CommandRecord{}, false, ErrCommandConflict
		}
		if entry.record != nil {
			record := cloneCommandRecord(*entry.record)
			s.mu.Unlock()
			return CommandLease{}, record, true, nil
		}
		done := entry.done
		s.mu.Unlock()

		select {
		case <-done:
			// The owner committed, rejected, or aborted. Re-check the key.
		case <-ctx.Done():
			return CommandLease{}, CommandRecord{}, false, ctx.Err()
		}
	}
}

func (s *MemoryStore) Commit(lease CommandLease, request CommitRequest) (CommandRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, err := s.pendingEntryLocked(lease)
	if err != nil {
		return CommandRecord{}, err
	}
	if request.Batch.CommandID != lease.key.CommandID {
		return CommandRecord{}, fmt.Errorf("%w: batch command does not match lease", ErrInvalidLease)
	}

	current := s.streams[request.MatchID]
	var currentRevision uint64
	if current != nil {
		currentRevision = current.revision
	}
	if currentRevision != request.Batch.PreviousRevision {
		return CommandRecord{}, fmt.Errorf("%w: expected %d, current %d", ErrRevisionConflict, request.Batch.PreviousRevision, currentRevision)
	}
	if request.Batch.NextRevision != request.Batch.PreviousRevision+1 {
		return CommandRecord{}, fmt.Errorf("invalid event batch revision transition %d -> %d", request.Batch.PreviousRevision, request.Batch.NextRevision)
	}
	if len(request.Batch.Events) == 0 {
		return CommandRecord{}, errors.New("event batch is empty")
	}
	for _, event := range request.Batch.Events {
		if event.MatchID != request.MatchID || event.Revision != request.Batch.NextRevision || event.CausedByCommand != request.Batch.CommandID {
			return CommandRecord{}, errors.New("event metadata does not match batch")
		}
	}
	if request.Authority != nil {
		if err := s.validateAuthorityLocked(lease, request.MatchID, *request.Authority); err != nil {
			return CommandRecord{}, err
		}
	}

	if current == nil {
		current = &stream{}
		s.streams[request.MatchID] = current
	}
	current.events = append(current.events, cloneEvents(request.Batch.Events)...)
	current.revision = request.Batch.NextRevision
	if request.Authority != nil {
		record := cloneAuthorityRecord(*request.Authority)
		s.authorities[authorityKey{matchID: record.MatchID, principalID: record.PrincipalID}] = record
		s.playerAuthorities[playerAuthorityKey{matchID: record.MatchID, playerID: record.PlayerID}] = record.PrincipalID
	}

	record := CommandRecord{
		Fingerprint: append([]byte(nil), entry.fingerprint...),
		Result:      append([]byte(nil), request.Result...),
		MatchID:     request.MatchID,
	}
	s.completeEntryLocked(entry, record)
	return cloneCommandRecord(record), nil
}

func (s *MemoryStore) RejectCommand(lease CommandLease, matchID model.MatchID, code, message string) (CommandRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, err := s.pendingEntryLocked(lease)
	if err != nil {
		return CommandRecord{}, err
	}
	if code == "" {
		return CommandRecord{}, errors.New("rejection code is required")
	}
	record := CommandRecord{
		Fingerprint:  append([]byte(nil), entry.fingerprint...),
		MatchID:      matchID,
		ErrorCode:    code,
		ErrorMessage: message,
	}
	s.completeEntryLocked(entry, record)
	return cloneCommandRecord(record), nil
}

func (s *MemoryStore) AbortCommand(lease CommandLease) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.commands[lease.key]
	if !ok || entry.token != lease.token || entry.record != nil {
		return
	}
	delete(s.commands, lease.key)
	close(entry.done)
}

func (s *MemoryStore) Events(matchID model.MatchID) []contracts.DomainEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	current := s.streams[matchID]
	if current == nil {
		return nil
	}
	return cloneEvents(current.events)
}

func (s *MemoryStore) ResolveAuthority(matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.authorities[authorityKey{matchID: matchID, principalID: principalID}]
	if !ok || record.Status != AuthorityActive {
		return "", false
	}
	return record.PlayerID, true
}

func (s *MemoryStore) pendingEntryLocked(lease CommandLease) (*commandEntry, error) {
	entry, ok := s.commands[lease.key]
	if !ok || entry.token != lease.token || entry.record != nil {
		return nil, ErrInvalidLease
	}
	return entry, nil
}

func (s *MemoryStore) validateAuthorityLocked(lease CommandLease, matchID model.MatchID, record AuthorityRecord) error {
	if record.MatchID != matchID || record.PrincipalID != lease.key.PrincipalID || record.EstablishedByCommandID != lease.key.CommandID {
		return fmt.Errorf("%w: authority metadata does not match command", ErrAuthorityConflict)
	}
	if record.PlayerID == "" || record.Seat <= 0 || record.BindingVersion == 0 || record.Status != AuthorityActive {
		return fmt.Errorf("%w: authority record is incomplete", ErrAuthorityConflict)
	}
	principalKey := authorityKey{matchID: record.MatchID, principalID: record.PrincipalID}
	if existing, ok := s.authorities[principalKey]; ok && existing.PlayerID != record.PlayerID {
		return fmt.Errorf("%w: principal is already bound to another player", ErrAuthorityConflict)
	}
	playerKey := playerAuthorityKey{matchID: record.MatchID, playerID: record.PlayerID}
	if existing, ok := s.playerAuthorities[playerKey]; ok && existing != record.PrincipalID {
		return fmt.Errorf("%w: player is already bound to another principal", ErrAuthorityConflict)
	}
	return nil
}

func (s *MemoryStore) completeEntryLocked(entry *commandEntry, record CommandRecord) {
	cloned := cloneCommandRecord(record)
	entry.record = &cloned
	close(entry.done)
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
				payloads[playerID] = append([]byte(nil), payload...)
			}
			event.PrivatePayloads = payloads
		}
		result[index] = event
	}
	return result
}
