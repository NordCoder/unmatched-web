# Persistence and Reconnect Contract

## Status

```text
status: draft-foundation
parent_issue: #19
```

## 1. Durability model

A match is recoverable from:

```text
immutable definition-version references
+ latest approved snapshot
+ ordered persisted events after the snapshot sequence
+ command idempotency records
```

Snapshots are optimization artifacts. Persisted events establish authoritative changes and random/choice outcomes.

## 2. Atomic accepted-command transaction

For one accepted command, the persistence layer atomically commits:

```yaml
command_result:
  command_id:
  status: accepted
  previous_revision:
  next_revision:
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

No success response or new projection is published before this transaction is durable.

A transaction failure returns no accepted result. Retrying the same `command_id` is safe.

## 3. Rejected commands

Rejected commands produce no gameplay event. The idempotency store may persist a safe rejection result to make duplicate delivery stable.

Protected operational audit logs may record security/validation details, but they are separate from match replay history and must not expose hidden information to clients.

## 4. Snapshot contract

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

- the state hash is computed over canonical serialization;
- snapshot creation cannot mutate gameplay state;
- a snapshot is accepted only after reducer/replay validation policy succeeds;
- old snapshots may be removed only when a newer recoverable snapshot and required events are durable;
- snapshots retain private authoritative state and require access controls.

## 5. Event log

Events are append-only for a match lineage. Corrections are new events or administrative lineage operations; historical gameplay events are not edited in place.

The event store enforces:

- unique `(match_id, sequence)`;
- unique `event_id`;
- monotonic contiguous sequence assignment;
- one durable result per accepted `command_id`;
- payload/schema version retention;
- private payload access policy.

## 6. Randomness persistence

A random operation follows:

```text
resolver requests random result
→ server RNG establishes outcome
→ RandomResultEstablished event is persisted
→ dependent events use the persisted result
```

For shuffles, persistence stores the resulting card-instance order or an equivalent complete permutation result. Reconnect/replay never calls RNG to recover a past shuffle.

## 7. Pending interaction persistence

Opening an interaction persists enough state to restore:

- interaction ID and source effect;
- owning player/group;
- prompt/choice schema;
- legal domain or deterministic domain-query context;
- cardinality, optionality and ordering requirements;
- private submitted values for simultaneous committed choices;
- exact resolver resume cursor.

A server restart after the interaction-open event returns the same logical interaction. It cannot repeat prior draws, reveals, damage or random outcomes.

## 8. Reconnect protocol

A reconnect request provides:

```yaml
match_id:
player_identity:
last_seen_revision:
last_seen_event_sequence:
client_instance_id:
```

After authentication, the server returns:

```yaml
ReconnectState:
  authoritative_revision:
  event_sequence:
  player_projection:
  legal_actions_or_pending_interaction:
  public_history_delta_optional:
  resync_required:
```

The server may send a full projection or compatible delta. The projection at a revision must be equivalent regardless of whether it was reached live or by reconnect.

## 9. Connection state

Connection/disconnection is operational state unless a published rule explicitly depends on it, which competitive Unmatched does not by default.

Disconnect cannot:

- transfer choice ownership;
- auto-decline an optional effect;
- select a default card/target;
- end the player's turn;
- regenerate a legal domain;
- expose the opponent's private state.

Timeout/forfeit policy is a product policy layer and must be represented by an explicit command/event when enabled.

## 10. Concurrent commands

The match command processor provides single-writer semantics per match.

Concurrent commands are serialized by accepted revision. A stale command is rejected unless a future explicit rebase policy proves its payload remains identical and legal; the foundation does not assume automatic rebasing.

Only the owner of the current pending interaction may advance it, except registered simultaneous committed-choice groups where each actor submits independently.

## 11. Replay

Replay loads the exact definition/capability versions, starts from an approved snapshot or initial state, and applies events in sequence.

Replay validation compares:

- expected revision/sequence;
- state hashes at checkpoints;
- pending interaction identity;
- public/private projection fixtures;
- game result.

Replay tools may expose private payloads only under an explicit authorized debugging mode.

## 12. Schema evolution

Every durable command result, event and snapshot includes a schema version.

A deployment may read an existing match only when it has:

- compatible readers/reducers for the stored event versions; or
- a verified deterministic migration preserving event meaning and state hashes at declared checkpoints.

Definition changes do not migrate active matches by default.

## 13. Recovery failures

If replay or snapshot verification fails, the match enters `QUARANTINED` operational state. The server preserves evidence and blocks gameplay commands.

It must not:

- skip an unknown event;
- fabricate a missing random result;
- discard a pending interaction;
- repair card order from current canonical definitions;
- continue from an unverified partial state.

## 14. Minimal persistence acceptance tests

- crash after event batch commit but before response: retry returns duplicate accepted result;
- crash before commit: retry executes once;
- reconnect during an optional choice: same interaction ID/domain returns;
- reconnect after shuffle: same order returns without RNG;
- replay from initial state equals replay from snapshot plus tail events;
- unauthorized reconnect cannot access the other player's private projection;
- two concurrent commands cannot both advance the same revision;
- ended match restores the same final result and rejects gameplay commands.
