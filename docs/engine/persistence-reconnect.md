# Persistence and Reconnect Contract

## Status

```text
status: draft-foundation
parent_issue: #19
correction_issue: #32
architecture: architecture-contract.md
state_model: state-model.md
command_contract: command-event-contract.md
fixture_contract: deterministic-fixture-contract.md
```

## 1. Durability model

A match is operationally recoverable from:

```text
immutable definition-version references
+ latest approved deterministic GameState snapshot
+ ordered persisted gameplay events after the snapshot sequence
+ command idempotency records
+ durable principal-to-player authority records
```

The deterministic replay model consumes definitions, snapshot and gameplay events. Command authorization/reconnect additionally consumes authority records. Operational session/connectivity presence is not required for recovery and is rebuilt separately.

Snapshots are optimization artifacts. Persisted events establish authoritative gameplay changes and random/choice outcomes. Authority and idempotency records are application/security state, not gameplay events.

## 2. Atomic accepted-command transaction

For one accepted existing-seat command, the persistence layer atomically commits:

```yaml
command_idempotency_record:
  principal_id:
  command_id:
  request_fingerprint:
  result:
event_batch:
  first_sequence:
  last_sequence:
  events: []
match_head:
  revision:
  event_sequence:
  lifecycle:
snapshot_policy_result:
```

For `CreateMatch` and `JoinMatch`, the same transaction additionally commits the new or changed `MatchAuthorityRecord` entries required by the lifecycle result.

No success response or new deterministic projection is published before every required record is durable. A transaction failure exposes no accepted result, no partial event batch and no partial principal/seat binding.

Retrying the same `(principal_id, command_id)` is safe under `command-event-contract.md`.

## 3. Rejected commands and durable idempotency

Rejected commands produce no gameplay event.

After authentication and successful canonical normalization, every terminal deterministic rejection is persisted in the command idempotency store with its exact request fingerprint and semantic result. This includes membership, binding, revision, lifecycle, legality, ownership, visibility and cost failures.

A same-ID/same-fingerprint retry returns the same durable rejection without re-evaluation. A same-ID/different-fingerprint request returns `DUPLICATE_CONFLICT` and never executes.

The following are not durable command results:

- authentication failure without a trusted principal identity;
- structurally unparseable input that cannot produce a canonical fingerprint;
- infrastructure/transaction failure before commit.

Protected operational audit logs may record such failures, but they are separate from match replay history and client-visible command results.

## 4. Command idempotency record

```yaml
CommandIdempotencyRecord:
  principal_id:
  command_id:
  fingerprint_schema_version:
  request_fingerprint:
  lifecycle_scope:
  match_id: null
  actor_player_id: null
  result_schema_version:
  result:
    semantic_status: accepted | rejected
    accepted_revision: null
    event_sequence_range: null
    rejection_code: null
    allocated_runtime_ids: {}
  committed_at_operational:
```

The store enforces unique `(principal_id, command_id)`. `committed_at_operational` is audit metadata and does not participate in gameplay replay or deterministic fixture hashes.

A result record cannot be overwritten by a request with another fingerprint. Fingerprint schema readers are retained for old records.

## 5. Match authority registry

```yaml
MatchAuthorityRecord:
  match_id:
  player_instance_id:
  principal_id:
  seat:
  binding_version:
  status: ACTIVE | REVOKED
  established_by_command_id:
```

The registry enforces:

- one active binding per match seat/player instance;
- no client-selected reassignment of a player instance;
- exact principal/match/player validation for existing-seat commands;
- atomic lifecycle binding with the corresponding gameplay event batch;
- protected access separate from public/private gameplay projections.

Principal identity is not required in gameplay events or deterministic snapshots. Replay reconstructs match/player state; the registry reconstructs external authority.

## 6. Snapshot contract

```yaml
Snapshot:
  match_id:
  revision:
  event_sequence:
  definition_ref:
  state_schema_version:
  state_payload:
  state_hash:
  created_by_policy:
```

Requirements:

- `state_payload` is deterministic `GameState` only;
- principal IDs, sessions, client instances, connectivity and last-seen metadata are excluded;
- the hash is computed over canonical serialization from `deterministic-fixture-contract.md`;
- snapshot creation cannot mutate gameplay state;
- a snapshot is accepted only after reducer/replay validation policy succeeds;
- old snapshots may be removed only when a newer recoverable snapshot and required events are durable;
- snapshots retain private gameplay state and require access controls.

## 7. Event log

Events are append-only for a match lineage. Historical gameplay events are not edited in place.

The event store enforces:

- unique `(match_id, sequence)`;
- unique `event_id`;
- monotonic contiguous sequence assignment;
- one revision transition per accepted command batch;
- payload/schema version retention;
- private payload access policy.

Principal/session/connectivity changes do not append gameplay events. A gameplay battlefield event named `ConnectionStateChanged` concerns a rules-defined battlefield connection/path, never network presence.

## 8. Randomness persistence

```text
resolver requests random result
→ server RNG establishes outcome
→ RandomResultEstablished event is persisted
→ dependent events use the persisted result
```

For shuffles, persistence stores resulting card-instance order or an equivalent complete permutation. Reconnect/replay never calls RNG to recover a past outcome.

## 9. Pending interaction persistence

Opening an interaction persists enough deterministic state to restore:

- interaction ID and source effect;
- owning player/group;
- prompt/choice schema;
- legal domain or deterministic query context;
- cardinality, optionality and ordering requirements;
- private submitted values for committed choices;
- exact resolver resume cursor.

A server restart or client reconnect returns the same logical interaction. It cannot repeat prior draws, reveals, damage, choices or random outcomes.

## 10. Reconnect request and authorization

A reconnect request supplies transport-neutral references:

```yaml
ReconnectRequest:
  match_id:
  last_seen_revision:
  last_seen_event_sequence:
  client_instance_id:
```

The external principal is derived from trusted authentication context, not from a payload field. The server:

1. authenticates the principal;
2. resolves an active `MatchAuthorityRecord` for the match;
3. rejects an unbound principal without revealing private match state;
4. loads deterministic `GameState` from snapshot plus ordered event tail;
5. derives the authorized player projection and legal actions/pending interaction;
6. registers the current session/client in the operational presence registry.

Response:

```yaml
ReconnectState:
  authoritative_revision:
  event_sequence:
  viewer_player_id:
  player_projection:
  legal_actions_or_pending_interaction:
  public_history_delta_optional:
  resync_required:
  operational_presence_optional:
```

The deterministic projection at a revision is equivalent regardless of whether it was reached live, after process restart or by reconnect. Operational presence is a separately composed field and is excluded from deterministic equality/hashes.

## 11. Operational connection state

Connection/disconnection belongs to `OperationalPresenceRegistry`, outside event-sourced `GameState`.

Presence changes:

- do not change match revision or event sequence;
- do not change deterministic state hash;
- do not transfer gameplay ownership;
- do not auto-decline an optional effect;
- do not select a default card/target;
- do not end a turn;
- do not regenerate a legal domain;
- do not expose opponent private state;
- are unavailable to deterministic Mechanics queries.

The registry may contain online state, active sessions, client-instance IDs and last-seen timestamps. It may be rebuilt after process restart.

A timeout/forfeit feature is a product policy layer. When enabled, it must issue an explicit authorized command whose accepted result is persisted as gameplay events. Socket loss alone is never a gameplay transition.

## 12. Concurrent commands

The match command processor provides single-writer semantics per match.

Concurrent requests first pass principal-scoped idempotency lookup. New requests are serialized by match revision. A stale command is rejected unless a future explicit rebase policy proves its normalized request remains identical and legal; the foundation assumes no automatic rebasing.

Only the owner of the current pending interaction may advance it, except registered simultaneous committed-choice groups where each actor submits independently.

Lifecycle transactions serialize seat allocation and authority binding so two joins cannot acquire the same seat.

## 13. Replay

Replay loads exact definition/capability versions, starts from an approved snapshot or initial deterministic state and applies events in sequence.

Replay validation compares:

- expected revision and sequence;
- canonical state hashes at checkpoints;
- pending interaction identity and resume cursor;
- public/private viewer projections;
- game result.

Gameplay replay does not consume principal authority or operational presence. An authorized replay service may combine replayed state with authority records only to decide which private projection a caller may receive.

Replay tools expose private payloads only under explicit protected authorization.

## 14. Schema evolution

Every durable command idempotency result, event and snapshot includes a schema version. Fingerprint schema version is retained separately.

A deployment may read an existing match only when it has:

- compatible readers/reducers for stored event/snapshot versions; or
- a verified deterministic migration preserving event meaning and declared state hashes.

Authority-record migrations preserve principal/player binding meaning. Definition changes do not migrate active matches by default.

## 15. Recovery failures

If deterministic replay or snapshot verification fails, the match enters quarantined handling and gameplay commands are blocked while evidence is preserved.

The system must not:

- skip an unknown event;
- fabricate a missing random result;
- discard a pending interaction;
- repair card order from current canonical definitions;
- continue from an unverified partial state;
- guess a missing principal-to-seat binding.

A missing/corrupt authority record is an authorization/recovery blocker, not permission to expose a player projection or create a replacement player identity.

## 16. Required persistence and reconnect evidence

Normative fixtures must prove:

- crash after accepted transaction commit but before response: retry returns duplicate accepted result;
- crash before commit: retry executes once;
- create/join authority binding is atomic with lifecycle events;
- same ID/different fingerprint returns `DUPLICATE_CONFLICT`;
- deterministic rejection is stable under same-request retry;
- disconnect/reconnect leaves revision, state hash and pending interaction unchanged;
- reconnect after shuffle restores the same order without RNG;
- replay from initial state equals snapshot plus ordered tail;
- player A/B projections remain authorization-specific after replay;
- unauthorized reconnect cannot access another player's private projection;
- two concurrent commands cannot both advance the same revision;
- ended match restores the same result and rejects gameplay commands.
