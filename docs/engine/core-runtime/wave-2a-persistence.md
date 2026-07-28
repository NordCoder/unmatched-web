# Core Runtime Wave 2A — PostgreSQL persistence

## Scope

This Candidate implements Issue #42 only. PostgreSQL becomes the durable source
of truth for match heads, ordered event batches, command idempotency, authority
bindings and pinned definition references. Snapshots, recovery quarantine,
reconnect and network transport remain in Issues #43–#45.

## Technology choice

The adapter uses `pgx/v5` and `pgxpool` directly. SQL migrations are embedded in
the server module and applied by `persistence.ApplyMigrations`. No ORM, Redis,
broker or distributed lease system is introduced.

## Durable tables

- `core_match_heads` stores the accepted revision, event cursor and pinned
  `DefinitionRef` for each match.
- `core_event_batches` stores one immutable ordered batch per accepted command.
- `core_events` stores the append-only event envelopes and public/private JSON
  payloads.
- `core_command_results` stores one immutable accepted or deterministic rejected
  result per `(principal_id, command_id)`.
- `core_command_reservations` stores operational in-flight claims. A stale claim
  can be replaced after the owning PostgreSQL session disappears because it is
  not a terminal command result.
- `core_authority_bindings` stores protected Principal → Player authority outside
  gameplay state and gameplay events.

Command results, event batches and events reject update/delete operations through
append-only database triggers.

## Single-writer and idempotency model

A command uses a dedicated pooled PostgreSQL connection for its execution lease.
The adapter acquires session advisory locks in one fixed order:

```text
command identity lock
→ match writer lock
```

The command lock ensures that only one process executes a normalized
`(principal_id, command_id)`. A same-fingerprint retry waits and then receives the
original immutable result. A different fingerprint receives a derived conflict
without changing or disclosing the original result.

For commands targeting an existing match, the match lock is acquired before state
loading, Rules invocation or runtime ID allocation. `CreateMatch` acquires the same
lock immediately after its generated match ID is known and before the durable
transaction. This provides one writer per match and prevents two commands from
committing the same revision or generated match identity. PostgreSQL releases both
locks if the process or connection dies.

## Accepted-command transaction

One PostgreSQL transaction performs:

```text
verify command reservation
+ verify current match head/revision/event cursor
+ insert immutable command result
+ insert authority binding when applicable
+ insert ordered event batch
+ insert every event
+ advance or create the match head
+ remove the operational reservation
```

No result is returned by `Host.Execute` until this transaction commits. Any
failure rolls back every insert and head mutation. Infrastructure failures abort
the reservation and create no durable command result.

A deterministic rejection uses a smaller transaction that inserts only the
immutable rejected result and removes the reservation. It emits no gameplay
event and changes no match revision.

## Local development

Start PostgreSQL:

```bash
docker compose up -d postgres
```

Copy `.env.example` to a local untracked `.env`, then load `DATABASE_URL` and the
pool settings into the process environment. Application startup should:

```go
pool, err := persistence.OpenPostgres(ctx, config)
if err != nil {
    // fail startup
}
if err := persistence.ApplyMigrations(ctx, pool); err != nil {
    // fail startup
}
store, err := persistence.NewPostgresStore(pool)
```

The migration runner is idempotent and serializes concurrent migrators with a
PostgreSQL transaction advisory lock.

## Integration tests

`tests/integration/core/postgres_test.go` uses `CORE_TEST_DATABASE_URL` when it is
set. Otherwise it starts an isolated `postgres:17-alpine` container with Docker.
Each run creates a unique schema, applies migrations and drops the schema during
cleanup.

The test covers:

```text
create + join durable match
store and Host reconstruction
same-request duplicate returns the original result
conflicting duplicate has zero durable mutation
deterministic rejection survives reconstruction
late transaction failure rolls back result, batch, event and head
concurrent duplicate across two stores invokes Rules once
command IDs remain principal-scoped inside one match
Host does not acknowledge while commit is blocked
```

## Deferred hardening

- Snapshot policy, snapshot hashes, tail recovery and quarantine are Issue #43.
- Principal reconnect and private projection recovery are Issue #44.
- HTTP/WebSocket adapters, production authentication and dependency-backed
  readiness are Issue #45.
- Horizontal ownership, distributed leases and performance tuning are outside
  the current single-process plus PostgreSQL target.
