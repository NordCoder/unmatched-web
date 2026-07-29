package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/jackc/pgx/v5"
)

func (s *PostgresStore) LoadEventsForCommand(ctx context.Context, lease CommandLease, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	leaseState, err := s.leaseState(lease)
	if err != nil {
		return nil, err
	}
	if err := ensureLeaseMatch(ctx, leaseState, lease, matchID); err != nil {
		return nil, err
	}
	rows, err := leaseState.conn.Query(ctx, `
SELECT event_schema_version, event_id, sequence, revision, event_type,
       caused_by_command_id, parent_event_id, source_ref, ruleset_version,
       public_payload, private_payloads_by_player
FROM core_events
WHERE match_id = $1
ORDER BY sequence`, matchID)
	if err != nil {
		return nil, fmt.Errorf("load match events through command lease: %w", err)
	}
	defer rows.Close()

	var events []contracts.DomainEvent
	for rows.Next() {
		var event contracts.DomainEvent
		var sequence, revision int64
		var publicPayload, privatePayloads []byte
		if err := rows.Scan(
			&event.SchemaVersion, &event.ID, &sequence, &revision, &event.Type,
			&event.CausedByCommand, &event.ParentEventID, &event.SourceRef,
			&event.RulesetVersion, &publicPayload, &privatePayloads,
		); err != nil {
			return nil, fmt.Errorf("scan match event through command lease: %w", err)
		}
		if sequence < 0 || revision < 0 {
			return nil, fmt.Errorf("stored event contains a negative cursor")
		}
		event.MatchID = matchID
		event.Sequence = uint64(sequence)
		event.Revision = uint64(revision)
		event.PublicPayload = append(json.RawMessage(nil), publicPayload...)
		if err := json.Unmarshal(privatePayloads, &event.PrivatePayloads); err != nil {
			return nil, fmt.Errorf("decode event private payloads: %w", err)
		}
		if event.PrivatePayloads == nil {
			event.PrivatePayloads = make(map[model.PlayerID]json.RawMessage)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate match events through command lease: %w", err)
	}
	return events, nil
}

func (s *PostgresStore) ResolveAuthorityForCommand(ctx context.Context, lease CommandLease, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	leaseState, err := s.leaseState(lease)
	if err != nil {
		return "", false, err
	}
	if err := ensureLeaseMatch(ctx, leaseState, lease, matchID); err != nil {
		return "", false, err
	}
	var playerID model.PlayerID
	err = leaseState.conn.QueryRow(ctx, `
SELECT player_id
FROM core_authority_bindings
WHERE match_id = $1 AND principal_id = $2 AND status = $3`,
		matchID, principalID, AuthorityActive,
	).Scan(&playerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("resolve principal authority through command lease: %w", err)
	}
	return playerID, true, nil
}

func ensureLeaseMatch(ctx context.Context, leaseState *postgresLease, lease CommandLease, matchID model.MatchID) error {
	if matchID == "" || lease.matchID == "" || lease.matchID != matchID {
		return fmt.Errorf("%w: lease does not own the requested match", ErrInvalidLease)
	}
	return leaseState.ensureMatchLock(ctx, matchID)
}

var _ EventStore = (*PostgresStore)(nil)
