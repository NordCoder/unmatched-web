CREATE TABLE IF NOT EXISTS core_match_heads (
    match_id TEXT PRIMARY KEY,
    revision BIGINT NOT NULL CHECK (revision >= 0),
    event_sequence BIGINT NOT NULL CHECK (event_sequence >= 0),
    definition_ref JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);

CREATE TABLE IF NOT EXISTS core_command_reservations (
    principal_id TEXT NOT NULL,
    command_id TEXT NOT NULL,
    fingerprint_schema_version TEXT NOT NULL,
    request_fingerprint BYTEA NOT NULL,
    reservation_token TEXT NOT NULL,
    lifecycle_scope TEXT NOT NULL CHECK (lifecycle_scope IN ('create_match', 'join_match', 'existing_seat')),
    match_id TEXT NOT NULL DEFAULT '',
    actor_player_id TEXT NOT NULL DEFAULT '',
    reserved_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (principal_id, command_id)
);

CREATE TABLE IF NOT EXISTS core_command_results (
    principal_id TEXT NOT NULL,
    command_id TEXT NOT NULL,
    fingerprint_schema_version TEXT NOT NULL,
    request_fingerprint BYTEA NOT NULL,
    lifecycle_scope TEXT NOT NULL CHECK (lifecycle_scope IN ('create_match', 'join_match', 'existing_seat')),
    match_id TEXT NOT NULL DEFAULT '',
    actor_player_id TEXT NOT NULL DEFAULT '',
    result_schema_version TEXT NOT NULL,
    result_payload BYTEA NOT NULL DEFAULT ''::bytea,
    error_code TEXT NOT NULL DEFAULT '',
    error_message TEXT NOT NULL DEFAULT '',
    committed_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (principal_id, command_id),
    CHECK (
        (error_code = '' AND octet_length(result_payload) > 0)
        OR
        (error_code <> '' AND octet_length(result_payload) = 0)
    )
);

CREATE TABLE IF NOT EXISTS core_authority_bindings (
    match_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    principal_id TEXT NOT NULL,
    seat INTEGER NOT NULL CHECK (seat > 0),
    binding_version BIGINT NOT NULL CHECK (binding_version > 0),
    status TEXT NOT NULL CHECK (status IN ('ACTIVE', 'REVOKED')),
    established_by_command_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (match_id, principal_id),
    CONSTRAINT core_authority_player_unique UNIQUE (match_id, player_id),
    CONSTRAINT core_authority_seat_unique UNIQUE (match_id, seat)
);

CREATE TABLE IF NOT EXISTS core_event_batches (
    match_id TEXT NOT NULL,
    revision BIGINT NOT NULL CHECK (revision > 0),
    previous_revision BIGINT NOT NULL CHECK (previous_revision >= 0),
    command_id TEXT NOT NULL,
    first_sequence BIGINT NOT NULL CHECK (first_sequence > 0),
    last_sequence BIGINT NOT NULL CHECK (last_sequence >= first_sequence),
    event_count INTEGER NOT NULL CHECK (event_count > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    PRIMARY KEY (match_id, revision),
    CHECK (revision = previous_revision + 1),
    CHECK (event_count = last_sequence - first_sequence + 1)
);

CREATE TABLE IF NOT EXISTS core_events (
    match_id TEXT NOT NULL,
    sequence BIGINT NOT NULL CHECK (sequence > 0),
    event_id TEXT NOT NULL UNIQUE,
    revision BIGINT NOT NULL CHECK (revision > 0),
    event_schema_version TEXT NOT NULL,
    event_type TEXT NOT NULL,
    caused_by_command_id TEXT NOT NULL,
    parent_event_id TEXT NOT NULL DEFAULT '',
    source_ref TEXT NOT NULL DEFAULT '',
    ruleset_version TEXT NOT NULL,
    public_payload JSONB NOT NULL,
    private_payloads_by_player JSONB NOT NULL,
    PRIMARY KEY (match_id, sequence),
    FOREIGN KEY (match_id, revision)
        REFERENCES core_event_batches(match_id, revision)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX IF NOT EXISTS core_events_match_revision_idx
    ON core_events(match_id, revision, sequence);

CREATE INDEX IF NOT EXISTS core_authority_principal_idx
    ON core_authority_bindings(principal_id, match_id)
    WHERE status = 'ACTIVE';

CREATE OR REPLACE FUNCTION core_reject_immutable_change()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION '% is append-only', TG_TABLE_NAME USING ERRCODE = '55000';
END;
$$;

DROP TRIGGER IF EXISTS core_command_results_immutable ON core_command_results;
CREATE TRIGGER core_command_results_immutable
BEFORE UPDATE OR DELETE ON core_command_results
FOR EACH ROW EXECUTE FUNCTION core_reject_immutable_change();

DROP TRIGGER IF EXISTS core_event_batches_immutable ON core_event_batches;
CREATE TRIGGER core_event_batches_immutable
BEFORE UPDATE OR DELETE ON core_event_batches
FOR EACH ROW EXECUTE FUNCTION core_reject_immutable_change();

DROP TRIGGER IF EXISTS core_events_immutable ON core_events;
CREATE TRIGGER core_events_immutable
BEFORE UPDATE OR DELETE ON core_events
FOR EACH ROW EXECUTE FUNCTION core_reject_immutable_change();
