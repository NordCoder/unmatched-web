package persistence

import (
	"context"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

func (s *MemoryStore) LoadEventsForCommand(ctx context.Context, lease CommandLease, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	if lease.backend != s {
		return nil, ErrInvalidLease
	}
	return s.LoadEvents(ctx, matchID)
}

func (s *MemoryStore) ResolveAuthorityForCommand(ctx context.Context, lease CommandLease, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	if lease.backend != s {
		return "", false, ErrInvalidLease
	}
	return s.ResolveAuthorityContext(ctx, matchID, principalID)
}

var _ EventStore = (*MemoryStore)(nil)
