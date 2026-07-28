package persistence

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const leaseCleanupTimeout = 5 * time.Second

type PostgresStore struct {
	pool *pgxpool.Pool
}

type postgresLease struct {
	store       *PostgresStore
	conn        *pgxpool.Conn
	commandLock int64
	matchLock   *int64

	mu       sync.Mutex
	released bool
}

func NewPostgresStore(pool *pgxpool.Pool) (*PostgresStore, error) {
	if pool == nil {
		return nil, fmt.Errorf("PostgreSQL pool is required")
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *PostgresStore) AcquireCommand(ctx context.Context, identity CommandIdentity) (CommandLease, CommandRecord, bool, error) {
	if err := validateCommandIdentity(identity); err != nil {
		return CommandLease{}, CommandRecord{}, false, err
	}
	key := identity.Key
	fingerprint := identity.Fingerprint
	matchID := identity.MatchID
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return CommandLease{}, CommandRecord{}, false, fmt.Errorf("acquire PostgreSQL connection: %w", err)
	}
	leaseState := &postgresLease{
		store:       s,
		conn:        conn,
		commandLock: advisoryKey("command", string(key.PrincipalID)+"\x00"+string(key.CommandID)),
	}
	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", leaseState.commandLock); err != nil {
		leaseState.release(ctx)
		return CommandLease{}, CommandRecord{}, false, fmt.Errorf("lock command identity: %w", err)
	}

	record, found, err := readCommandRecord(ctx, conn, key)
	if err != nil {
		leaseState.release(ctx)
		return CommandLease{}, CommandRecord{}, false, err
	}
	if found {
		leaseState.release(ctx)
		if !bytes.Equal(record.Fingerprint, fingerprint) {
			return CommandLease{}, CommandRecord{}, false, ErrCommandConflict
		}
		return CommandLease{}, record, true, nil
	}

	if matchID != "" {
		matchLock := advisoryKey("match", string(matchID))
		if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", matchLock); err != nil {
			leaseState.release(ctx)
			return CommandLease{}, CommandRecord{}, false, fmt.Errorf("lock match writer: %w", err)
		}
		leaseState.matchLock = &matchLock
	}

	token, err := randomToken()
	if err != nil {
		leaseState.release(ctx)
		return CommandLease{}, CommandRecord{}, false, err
	}
	if _, err := conn.Exec(ctx, `
INSERT INTO core_command_reservations (
    principal_id, command_id, fingerprint_schema_version, request_fingerprint,
    reservation_token, lifecycle_scope, match_id, actor_player_id, reserved_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, clock_timestamp())
ON CONFLICT (principal_id, command_id) DO UPDATE SET
    fingerprint_schema_version = EXCLUDED.fingerprint_schema_version,
    request_fingerprint = EXCLUDED.request_fingerprint,
    reservation_token = EXCLUDED.reservation_token,
    lifecycle_scope = EXCLUDED.lifecycle_scope,
    match_id = EXCLUDED.match_id,
    actor_player_id = EXCLUDED.actor_player_id,
    reserved_at = EXCLUDED.reserved_at`,
		key.PrincipalID, key.CommandID, FingerprintSchemaV1, fingerprint, token,
		identity.Scope, matchID, identity.ActorPlayerID,
	); err != nil {
		leaseState.release(ctx)
		return CommandLease{}, CommandRecord{}, false, fmt.Errorf("reserve command identity: %w", err)
	}

	return CommandLease{
		key:           key,
		token:         token,
		matchID:       matchID,
		actorPlayerID: identity.ActorPlayerID,
		scope:         identity.Scope,
		backend:       leaseState,
	}, CommandRecord{}, false, nil
}

func (s *PostgresStore) Commit(ctx context.Context, lease CommandLease, request CommitRequest) (CommandRecord, error) {
	leaseState, err := s.leaseState(lease)
	if err != nil {
		return CommandRecord{}, err
	}
	if err := validateCommitRequest(lease, request); err != nil {
		return CommandRecord{}, err
	}
	if err := leaseState.ensureMatchLock(ctx, request.MatchID); err != nil {
		return CommandRecord{}, err
	}
	definitionJSON, err := json.Marshal(request.DefinitionRef)
	if err != nil {
		return CommandRecord{}, fmt.Errorf("encode pinned definition reference: %w", err)
	}
	previousRevisionDB, err := safeInt64(request.Batch.PreviousRevision)
	if err != nil {
		return CommandRecord{}, err
	}
	nextRevisionDB, err := safeInt64(request.Batch.NextRevision)
	if err != nil {
		return CommandRecord{}, err
	}
	firstSequence := request.Batch.Events[0].Sequence
	lastSequence := request.Batch.Events[len(request.Batch.Events)-1].Sequence
	firstSequenceDB, err := safeInt64(firstSequence)
	if err != nil {
		return CommandRecord{}, err
	}
	lastSequenceDB, err := safeInt64(lastSequence)
	if err != nil {
		return CommandRecord{}, err
	}

	tx, err := leaseState.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return CommandRecord{}, fmt.Errorf("begin accepted-command transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	fingerprint, err := verifyReservation(ctx, tx, lease)
	if err != nil {
		return CommandRecord{}, err
	}

	headRevision, headSequence, headDefinition, headExists, err := loadMatchHeadForUpdate(ctx, tx, request.MatchID)
	if err != nil {
		return CommandRecord{}, err
	}
	if !headExists {
		if request.Batch.PreviousRevision != 0 || request.Batch.Events[0].Sequence != 1 {
			return CommandRecord{}, fmt.Errorf("%w: new match must begin at revision 0 and sequence 1", ErrRevisionConflict)
		}
	} else {
		if headRevision != request.Batch.PreviousRevision {
			return CommandRecord{}, fmt.Errorf("%w: expected %d, current %d", ErrRevisionConflict, request.Batch.PreviousRevision, headRevision)
		}
		if request.Batch.Events[0].Sequence != headSequence+1 {
			return CommandRecord{}, fmt.Errorf("event sequence must continue at %d", headSequence+1)
		}
		var storedRef model.DefinitionRef
		if err := json.Unmarshal(headDefinition, &storedRef); err != nil {
			return CommandRecord{}, fmt.Errorf("decode stored definition reference: %w", err)
		}
		if !reflect.DeepEqual(storedRef, request.DefinitionRef) {
			return CommandRecord{}, fmt.Errorf("pinned definition reference changed")
		}
	}

	if _, err := tx.Exec(ctx, `
INSERT INTO core_command_results (
    principal_id, command_id, fingerprint_schema_version, request_fingerprint,
    lifecycle_scope, match_id, actor_player_id,
    result_schema_version, result_payload, error_code, error_message
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, '', '')`,
		lease.key.PrincipalID, lease.key.CommandID, FingerprintSchemaV1, fingerprint,
		lease.scope, request.MatchID, lease.actorPlayerID, CommandResultSchemaV1, request.Result,
	); err != nil {
		return CommandRecord{}, classifyPostgresError("insert command result", err)
	}

	if request.Authority != nil {
		if err := validateAuthority(lease, request.MatchID, *request.Authority); err != nil {
			return CommandRecord{}, err
		}
		bindingVersionDB, err := safeInt64(request.Authority.BindingVersion)
		if err != nil {
			return CommandRecord{}, err
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO core_authority_bindings (
    match_id, player_id, principal_id, seat, binding_version, status,
    established_by_command_id
) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			request.Authority.MatchID, request.Authority.PlayerID, request.Authority.PrincipalID,
			request.Authority.Seat, bindingVersionDB, request.Authority.Status,
			request.Authority.EstablishedByCommandID,
		); err != nil {
			return CommandRecord{}, classifyPostgresError("insert authority binding", err)
		}
	}

	if _, err := tx.Exec(ctx, `
INSERT INTO core_event_batches (
    match_id, revision, previous_revision, command_id,
    first_sequence, last_sequence, event_count
) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		request.MatchID, nextRevisionDB, previousRevisionDB,
		request.Batch.CommandID, firstSequenceDB, lastSequenceDB, len(request.Batch.Events),
	); err != nil {
		return CommandRecord{}, classifyPostgresError("insert event batch", err)
	}

	for _, event := range request.Batch.Events {
		eventSequenceDB, err := safeInt64(event.Sequence)
		if err != nil {
			return CommandRecord{}, err
		}
		eventRevisionDB, err := safeInt64(event.Revision)
		if err != nil {
			return CommandRecord{}, err
		}
		privatePayloads, err := json.Marshal(event.PrivatePayloads)
		if err != nil {
			return CommandRecord{}, fmt.Errorf("encode event private payloads: %w", err)
		}
		if len(privatePayloads) == 0 || string(privatePayloads) == "null" {
			privatePayloads = []byte(`{}`)
		}
		publicPayload := event.PublicPayload
		if len(publicPayload) == 0 {
			publicPayload = []byte(`{}`)
		}
		if !json.Valid(publicPayload) || !json.Valid(privatePayloads) {
			return CommandRecord{}, fmt.Errorf("event payload is not valid JSON")
		}
		if _, err := tx.Exec(ctx, `
INSERT INTO core_events (
    match_id, sequence, event_id, revision, event_schema_version, event_type,
    caused_by_command_id, parent_event_id, source_ref, ruleset_version,
    public_payload, private_payloads_by_player
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			event.MatchID, eventSequenceDB, event.ID, eventRevisionDB, event.SchemaVersion,
			event.Type, event.CausedByCommand, event.ParentEventID, event.SourceRef,
			event.RulesetVersion, string(publicPayload), string(privatePayloads),
		); err != nil {
			return CommandRecord{}, classifyPostgresError("insert event", err)
		}
	}

	if headExists {
		tag, err := tx.Exec(ctx, `
UPDATE core_match_heads
SET revision = $2, event_sequence = $3, updated_at = clock_timestamp()
WHERE match_id = $1 AND revision = $4 AND event_sequence = $5`,
			request.MatchID, nextRevisionDB, lastSequenceDB,
			previousRevisionDB, int64(headSequence),
		)
		if err != nil {
			return CommandRecord{}, classifyPostgresError("advance match head", err)
		}
		if tag.RowsAffected() != 1 {
			return CommandRecord{}, ErrRevisionConflict
		}
	} else {
		if _, err := tx.Exec(ctx, `
INSERT INTO core_match_heads (match_id, revision, event_sequence, definition_ref)
VALUES ($1, $2, $3, $4)`,
			request.MatchID, nextRevisionDB, lastSequenceDB, string(definitionJSON),
		); err != nil {
			return CommandRecord{}, classifyPostgresError("insert match head", err)
		}
	}

	if err := deleteReservation(ctx, tx, lease); err != nil {
		return CommandRecord{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return CommandRecord{}, classifyPostgresError("commit accepted command", err)
	}
	leaseState.release(ctx)
	return CommandRecord{
		Fingerprint: append([]byte(nil), fingerprint...),
		Result:      append([]byte(nil), request.Result...),
		MatchID:     request.MatchID,
	}, nil
}

func (s *PostgresStore) RejectCommand(ctx context.Context, lease CommandLease, matchID model.MatchID, code, message string) (CommandRecord, error) {
	leaseState, err := s.leaseState(lease)
	if err != nil {
		return CommandRecord{}, err
	}
	if code == "" {
		return CommandRecord{}, fmt.Errorf("rejection code is required")
	}
	if lease.matchID != "" && matchID != lease.matchID {
		return CommandRecord{}, fmt.Errorf("%w: rejected match does not match reserved match", ErrInvalidLease)
	}
	tx, err := leaseState.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return CommandRecord{}, fmt.Errorf("begin rejected-command transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()

	fingerprint, err := verifyReservation(ctx, tx, lease)
	if err != nil {
		return CommandRecord{}, err
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO core_command_results (
    principal_id, command_id, fingerprint_schema_version, request_fingerprint,
    lifecycle_scope, match_id, actor_player_id,
    result_schema_version, result_payload, error_code, error_message
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, ''::bytea, $9, $10)`,
		lease.key.PrincipalID, lease.key.CommandID, FingerprintSchemaV1, fingerprint,
		lease.scope, matchID, lease.actorPlayerID, CommandResultSchemaV1, code, message,
	); err != nil {
		return CommandRecord{}, classifyPostgresError("insert rejected command result", err)
	}
	if err := deleteReservation(ctx, tx, lease); err != nil {
		return CommandRecord{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return CommandRecord{}, classifyPostgresError("commit rejected command", err)
	}
	leaseState.release(ctx)
	return CommandRecord{
		Fingerprint:  append([]byte(nil), fingerprint...),
		MatchID:      matchID,
		ErrorCode:    code,
		ErrorMessage: message,
	}, nil
}

func (s *PostgresStore) AbortCommand(ctx context.Context, lease CommandLease) error {
	leaseState, err := s.leaseState(lease)
	if err != nil {
		return err
	}
	_, deleteErr := leaseState.conn.Exec(ctx, `
DELETE FROM core_command_reservations
WHERE principal_id = $1 AND command_id = $2 AND reservation_token = $3`,
		lease.key.PrincipalID, lease.key.CommandID, lease.token,
	)
	leaseState.release(ctx)
	if deleteErr != nil {
		return fmt.Errorf("abort command reservation: %w", deleteErr)
	}
	return nil
}

func (s *PostgresStore) LoadEvents(ctx context.Context, matchID model.MatchID) ([]contracts.DomainEvent, error) {
	rows, err := s.pool.Query(ctx, `
SELECT event_schema_version, event_id, sequence, revision, event_type,
       caused_by_command_id, parent_event_id, source_ref, ruleset_version,
       public_payload, private_payloads_by_player
FROM core_events
WHERE match_id = $1
ORDER BY sequence`, matchID)
	if err != nil {
		return nil, fmt.Errorf("load match events: %w", err)
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
			return nil, fmt.Errorf("scan match event: %w", err)
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
		return nil, fmt.Errorf("iterate match events: %w", err)
	}
	return events, nil
}

func (s *PostgresStore) ResolveAuthorityContext(ctx context.Context, matchID model.MatchID, principalID model.PrincipalID) (model.PlayerID, bool, error) {
	var playerID model.PlayerID
	err := s.pool.QueryRow(ctx, `
SELECT player_id
FROM core_authority_bindings
WHERE match_id = $1 AND principal_id = $2 AND status = $3`,
		matchID, principalID, AuthorityActive,
	).Scan(&playerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("resolve principal authority: %w", err)
	}
	return playerID, true, nil
}

func (s *PostgresStore) leaseState(lease CommandLease) (*postgresLease, error) {
	leaseState, ok := lease.backend.(*postgresLease)
	if !ok || leaseState == nil || leaseState.store != s || lease.key.PrincipalID == "" || lease.key.CommandID == "" || lease.token == "" {
		return nil, ErrInvalidLease
	}
	leaseState.mu.Lock()
	released := leaseState.released
	leaseState.mu.Unlock()
	if released {
		return nil, ErrInvalidLease
	}
	return leaseState, nil
}

func (l *postgresLease) ensureMatchLock(ctx context.Context, matchID model.MatchID) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.released {
		return ErrInvalidLease
	}
	key := advisoryKey("match", string(matchID))
	if l.matchLock != nil {
		if *l.matchLock != key {
			return fmt.Errorf("%w: lease already owns another match lock", ErrInvalidLease)
		}
		return nil
	}
	if _, err := l.conn.Exec(ctx, "SELECT pg_advisory_lock($1)", key); err != nil {
		return fmt.Errorf("lock generated match writer: %w", err)
	}
	l.matchLock = &key
	return nil
}

func (l *postgresLease) release(ctx context.Context) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.released {
		return
	}
	l.released = true

	cleanupContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), leaseCleanupTimeout)
	defer cancel()
	healthy := true
	if l.matchLock != nil {
		var unlocked bool
		if err := l.conn.QueryRow(cleanupContext, "SELECT pg_advisory_unlock($1)", *l.matchLock).Scan(&unlocked); err != nil || !unlocked {
			healthy = false
		}
	}
	var unlocked bool
	if err := l.conn.QueryRow(cleanupContext, "SELECT pg_advisory_unlock($1)", l.commandLock).Scan(&unlocked); err != nil || !unlocked {
		healthy = false
	}
	if healthy {
		l.conn.Release()
		return
	}
	raw := l.conn.Hijack()
	_ = raw.Close(cleanupContext)
}

type commandRecordQuerier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readCommandRecord(ctx context.Context, query commandRecordQuerier, key CommandKey) (CommandRecord, bool, error) {
	var record CommandRecord
	err := query.QueryRow(ctx, `
SELECT request_fingerprint, result_payload, match_id, error_code, error_message
FROM core_command_results
WHERE principal_id = $1 AND command_id = $2`,
		key.PrincipalID, key.CommandID,
	).Scan(&record.Fingerprint, &record.Result, &record.MatchID, &record.ErrorCode, &record.ErrorMessage)
	if errors.Is(err, pgx.ErrNoRows) {
		return CommandRecord{}, false, nil
	}
	if err != nil {
		return CommandRecord{}, false, fmt.Errorf("read command result: %w", err)
	}
	return record, true, nil
}

func verifyReservation(ctx context.Context, tx pgx.Tx, lease CommandLease) ([]byte, error) {
	var fingerprint []byte
	err := tx.QueryRow(ctx, `
SELECT request_fingerprint
FROM core_command_reservations
WHERE principal_id = $1 AND command_id = $2 AND reservation_token = $3
FOR UPDATE`,
		lease.key.PrincipalID, lease.key.CommandID, lease.token,
	).Scan(&fingerprint)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrInvalidLease
	}
	if err != nil {
		return nil, fmt.Errorf("verify command reservation: %w", err)
	}
	return fingerprint, nil
}

func deleteReservation(ctx context.Context, tx pgx.Tx, lease CommandLease) error {
	tag, err := tx.Exec(ctx, `
DELETE FROM core_command_reservations
WHERE principal_id = $1 AND command_id = $2 AND reservation_token = $3`,
		lease.key.PrincipalID, lease.key.CommandID, lease.token,
	)
	if err != nil {
		return fmt.Errorf("delete command reservation: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return ErrInvalidLease
	}
	return nil
}

func loadMatchHeadForUpdate(ctx context.Context, tx pgx.Tx, matchID model.MatchID) (uint64, uint64, []byte, bool, error) {
	var revision, sequence int64
	var definitionRef []byte
	err := tx.QueryRow(ctx, `
SELECT revision, event_sequence, definition_ref
FROM core_match_heads
WHERE match_id = $1
FOR UPDATE`, matchID).Scan(&revision, &sequence, &definitionRef)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, nil, false, nil
	}
	if err != nil {
		return 0, 0, nil, false, fmt.Errorf("load match head: %w", err)
	}
	if revision < 0 || sequence < 0 {
		return 0, 0, nil, false, fmt.Errorf("stored match head contains a negative cursor")
	}
	return uint64(revision), uint64(sequence), definitionRef, true, nil
}

func validateAuthority(lease CommandLease, matchID model.MatchID, record AuthorityRecord) error {
	if record.MatchID != matchID || record.PrincipalID != lease.key.PrincipalID || record.EstablishedByCommandID != lease.key.CommandID {
		return fmt.Errorf("%w: authority metadata does not match command", ErrAuthorityConflict)
	}
	if record.PlayerID == "" || record.Seat <= 0 || record.BindingVersion == 0 || record.Status != AuthorityActive {
		return fmt.Errorf("%w: authority record is incomplete", ErrAuthorityConflict)
	}
	return nil
}

func classifyPostgresError(operation string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch {
		case pgErr.Code == "23505" && strings.HasPrefix(pgErr.ConstraintName, "core_authority"):
			return fmt.Errorf("%w: %s", ErrAuthorityConflict, pgErr.Message)
		case pgErr.Code == "23505" && pgErr.ConstraintName == "core_command_results_pkey":
			return fmt.Errorf("%w: %s", ErrCommandConflict, pgErr.Message)
		case pgErr.Code == "23505" && (pgErr.ConstraintName == "core_match_heads_pkey" || pgErr.ConstraintName == "core_event_batches_pkey"):
			return fmt.Errorf("%w: %s", ErrRevisionConflict, pgErr.Message)
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}

func advisoryKey(namespace, value string) int64 {
	sum := sha256.Sum256([]byte(namespace + "\x00" + value))
	return int64(binary.BigEndian.Uint64(sum[:8]))
}

func randomToken() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("allocate command reservation token: %w", err)
	}
	return hex.EncodeToString(raw[:]), nil
}

func safeInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("cursor %d exceeds PostgreSQL BIGINT", value)
	}
	return int64(value), nil
}
