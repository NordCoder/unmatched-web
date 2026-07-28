package persistence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

var (
	ErrRevisionConflict = errors.New("revision conflict")
	ErrCommandConflict  = errors.New("command identity conflict")
)

type CommandRecord struct {
	Fingerprint []byte
	Result      []byte
	MatchID     model.MatchID
}

type CommitRequest struct {
	MatchID     model.MatchID
	Fingerprint []byte
	Batch       contracts.EventBatch
	Result      []byte
}

type EventStore interface {
	LookupCommand(commandID model.CommandID) (CommandRecord, bool)
	Commit(request CommitRequest) (record CommandRecord, duplicate bool, err error)
	Events(matchID model.MatchID) []contracts.DomainEvent
}

type stream struct {
	revision uint64
	events   []contracts.DomainEvent
}

type MemoryStore struct {
	mu       sync.RWMutex
	streams  map[model.MatchID]*stream
	commands map[model.CommandID]CommandRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		streams:  make(map[model.MatchID]*stream),
		commands: make(map[model.CommandID]CommandRecord),
	}
}

func (s *MemoryStore) LookupCommand(commandID model.CommandID) (CommandRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.commands[commandID]
	return cloneCommandRecord(record), ok
}

func (s *MemoryStore) Commit(request CommitRequest) (CommandRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.commands[request.Batch.CommandID]; ok {
		if !bytes.Equal(existing.Fingerprint, request.Fingerprint) {
			return CommandRecord{}, false, ErrCommandConflict
		}
		return cloneCommandRecord(existing), true, nil
	}

	current := s.streams[request.MatchID]
	var currentRevision uint64
	if current != nil {
		currentRevision = current.revision
	}
	if currentRevision != request.Batch.PreviousRevision {
		return CommandRecord{}, false, fmt.Errorf("%w: expected %d, current %d", ErrRevisionConflict, request.Batch.PreviousRevision, currentRevision)
	}
	if request.Batch.NextRevision != request.Batch.PreviousRevision+1 {
		return CommandRecord{}, false, fmt.Errorf("invalid event batch revision transition %d -> %d", request.Batch.PreviousRevision, request.Batch.NextRevision)
	}
	if len(request.Batch.Events) == 0 {
		return CommandRecord{}, false, errors.New("event batch is empty")
	}
	for _, event := range request.Batch.Events {
		if event.MatchID != request.MatchID || event.Revision != request.Batch.NextRevision || event.CausedByCommand != request.Batch.CommandID {
			return CommandRecord{}, false, errors.New("event metadata does not match batch")
		}
	}

	if current == nil {
		current = &stream{}
		s.streams[request.MatchID] = current
	}
	current.events = append(current.events, cloneEvents(request.Batch.Events)...)
	current.revision = request.Batch.NextRevision

	record := CommandRecord{
		Fingerprint: append([]byte(nil), request.Fingerprint...),
		Result:      append([]byte(nil), request.Result...),
		MatchID:     request.MatchID,
	}
	s.commands[request.Batch.CommandID] = record
	return cloneCommandRecord(record), false, nil
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

func cloneCommandRecord(record CommandRecord) CommandRecord {
	record.Fingerprint = append([]byte(nil), record.Fingerprint...)
	record.Result = append([]byte(nil), record.Result...)
	return record
}

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
