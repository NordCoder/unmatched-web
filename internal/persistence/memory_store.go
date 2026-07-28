package persistence

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

type stream struct {
	revision      uint64
	eventSequence uint64
	definitionRef model.DefinitionRef
	events        []contracts.DomainEvent
}

type commandEntry struct {
	fingerprint []byte
	record      *CommandRecord
	done        chan struct{}
	token       string
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

func (s *MemoryStore) AcquireCommand(ctx context.Context, identity CommandIdentity) (CommandLease, CommandRecord, bool, error) {
	key := identity.Key
	fingerprint := identity.Fingerprint
	if err := validateCommandIdentity(identity); err != nil {
		return CommandLease{}, CommandRecord{}, false, err
	}
	for {
		s.mu.Lock()
		entry, exists := s.commands[key]
		if !exists {
			s.nextToken++
			entry = &commandEntry{
				fingerprint: append([]byte(nil), fingerprint...),
				done:        make(chan struct{}),
				token:       strconv.FormatUint(s.nextToken, 10),
			}
			s.commands[key] = entry
			lease := CommandLease{key: key, token: entry.token, matchID: identity.MatchID, actorPlayerID: identity.ActorPlayerID, scope: identity.Scope, backend: s}
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
		case <-ctx.Done():
			return CommandLease{}, CommandRecord{}, false, ctx.Err()
		}
	}
}

func (s *MemoryStore) Commit(_ context.Context, lease CommandLease, request CommitRequest) (CommandRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, err := s.pendingEntryLocked(lease)
	if err != nil {
		return CommandRecord{}, err
	}
	if err := validateCommitRequest(lease, request); err != nil {
		return CommandRecord{}, err
	}

	current := s.streams[request.MatchID]
	var currentRevision, currentSequence uint64
	if current != nil {
		currentRevision = current.revision
		currentSequence = current.eventSequence
		if !reflect.DeepEqual(current.definitionRef, request.DefinitionRef) {
			return CommandRecord{}, errors.New("pinned definition reference changed")
		}
	}
	if currentRevision != request.Batch.PreviousRevision {
		return CommandRecord{}, fmt.Errorf("%w: expected %d, current %d", ErrRevisionConflict, request.Batch.PreviousRevision, currentRevision)
	}
	if request.Batch.Events[0].Sequence != currentSequence+1 {
		return CommandRecord{}, fmt.Errorf("event sequence must continue at %d", currentSequence+1)
	}
	if request.Authority != nil {
		if err := s.validateAuthorityLocked(lease, request.MatchID, *request.Authority); err != nil {
			return CommandRecord{}, err
		}
	}

	if current == nil {
		current = &stream{definitionRef: cloneDefinitionRef(request.DefinitionRef)}
		s.streams[request.MatchID] = current
	}
	current.events = append(current.events, cloneEvents(request.Batch.Events)...)
	current.revision = request.Batch.NextRevision
	current.eventSequence = request.Batch.Events[len(request.Batch.Events)-1].Sequence
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

func (s *MemoryStore) RejectCommand(_ context.Context, lease CommandLease, matchID model.MatchID, code, message string) (CommandRecord, error) {
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

func (s *MemoryStore) AbortCommand(_ context.Context, lease CommandLease) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if lease.backend != s {
		return ErrInvalidLease
	}
	entry, ok := s.commands[lease.key]
	if !ok || entry.token != lease.token || entry.record != nil {
		return nil
	}
	delete(s.commands, lease.key)
	close(entry.done)
	return nil
}

func (s *MemoryStore) LoadEvents(_ context.Context, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	return s.Events(matchID), nil
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

func (s *MemoryStore) ResolveAuthorityContext(_ context.Context, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	playerID, ok := s.ResolveAuthority(matchID, principalID)
	return playerID, ok, nil
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
	if lease.backend != s {
		return nil, ErrInvalidLease
	}
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

func validateCommitRequest(lease CommandLease, request CommitRequest) error {
	if len(request.Batch.Events) == 0 {
		return errors.New("event batch is empty")
	}
	if request.MatchID == "" || request.Batch.CommandID != lease.key.CommandID || request.MatchID != request.Batch.Events[0].MatchID {
		return fmt.Errorf("%w: commit identity does not match lease", ErrInvalidLease)
	}
	if lease.matchID != "" && request.MatchID != lease.matchID {
		return fmt.Errorf("%w: committed match does not match reserved match", ErrInvalidLease)
	}
	if request.Batch.NextRevision != request.Batch.PreviousRevision+1 {
		return fmt.Errorf("invalid event batch revision transition %d -> %d", request.Batch.PreviousRevision, request.Batch.NextRevision)
	}
	if len(request.Result) == 0 {
		return errors.New("command result is empty")
	}
	if request.DefinitionRef.RulesetVersion == "" || request.DefinitionRef.CapabilityRegistry == "" {
		return errors.New("pinned definition reference is incomplete")
	}
	var previousSequence uint64
	for index, event := range request.Batch.Events {
		if event.MatchID != request.MatchID || event.Revision != request.Batch.NextRevision || event.CausedByCommand != request.Batch.CommandID {
			return errors.New("event metadata does not match batch")
		}
		if event.ID == "" || event.SchemaVersion == "" || event.Type == "" || event.RulesetVersion == "" {
			return errors.New("event envelope is incomplete")
		}
		if index > 0 && event.Sequence != previousSequence+1 {
			return errors.New("event sequence inside batch is not contiguous")
		}
		previousSequence = event.Sequence
	}
	return nil
}

func validateCommandIdentity(identity CommandIdentity) error {
	if identity.Key.PrincipalID == "" || identity.Key.CommandID == "" || len(identity.Fingerprint) == 0 {
		return ErrInvalidLease
	}
	switch identity.Scope {
	case CommandScopeCreateMatch, CommandScopeJoinMatch, CommandScopeExistingSeat:
		return nil
	default:
		return fmt.Errorf("%w: unknown command scope", ErrInvalidLease)
	}
}

func cloneDefinitionRef(ref model.DefinitionRef) model.DefinitionRef {
	ref.FighterManifestDigests = cloneStringMap(ref.FighterManifestDigests)
	ref.CardManifestDigests = cloneStringMap(ref.CardManifestDigests)
	return ref
}

func cloneStringMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
