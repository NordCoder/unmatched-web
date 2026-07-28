# Persistence and Reconnect Contract

## Status

```text
status: normative-foundation
parent_issue: #19
correction_round_1: #32
correction_round_2: #35
architecture: architecture-contract.md
state_model: state-model.md
command_contract: command-event-contract.md
fixture_contract: deterministic-fixture-contract.md
```

## 1. Durable recovery model

A match is recoverable from:

```text
pinned definition versions
+ latest approved deterministic snapshot
+ contiguous persisted gameplay event tail
+ immutable command idempotency records
+ durable principal-to-player authority records
```

Operational presence is rebuilt separately and is not replay input.

## 2. Atomic accepted-command transaction

An accepted existing-seat command atomically commits:

```yaml
command_idempotency_record:
event_batch:
match_head:
snapshot_policy_result:
```

`CreateMatch` and `JoinMatch` additionally commit required authority records in the same transaction.

No success or new deterministic projection is published before durable commit. A failed transaction exposes no partial result, event batch or seat binding.

## 3. Rejected commands

A newly submitted command with an unoccupied `(principal_id, command_id)` records one deterministic rejection after authentication and canonical normalization. It emits no gameplay event.

A same-key/same-fingerprint retry returns the immutable rejection without re-evaluation.

Authentication failures, unparseable requests and pre-commit infrastructure failures are not command results.

## 4. Idempotency namespace and collision behavior

```yaml
CommandIdempotencyRecord:
  principal_id:
  command_id:
  fingerprint_schema_version:
  request_fingerprint:
  lifecycle_scope:
  match_id:
  actor_player_id:
  result_schema_version:
  result:
    semantic_status: accepted | rejected
    accepted_revision:
    event_sequence_range:
    rejection_code:
    allocated_runtime_ids: {}
  committed_at_operational:
```

The store enforces exactly one immutable record per `(principal_id, command_id)`.

A different presented fingerprint does not create another record. The server derives `DUPLICATE_CONFLICT` from the occupied key and the mismatch:

- original record unchanged;
- conflicting request not executed;
- no command-result insert/update;
- no gameplay event, revision or sequence delta;
- no original payload, fingerprint, result or private data disclosed;
- repeated conflicts produce the same response class.

A separate security audit may record collision occurrence, but it is operational, non-authoritative and excluded from deterministic state, replay, projections and fixture hashes.

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

The registry enforces exact principal/match/player authority and atomic lifecycle seat assignment. It is protected application/security state, not gameplay state.

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

- payload contains deterministic `GameState` only;
- authority, principal, session and presence data are excluded;
- hashes use NFC-normalized RFC 8785 bytes and SHA-256;
- snapshot creation does not mutate gameplay state;
- private gameplay state remains access controlled.

## 7. Event log

Events are append-only. The store enforces:

- unique `(match_id, sequence)`;
- unique `event_id`;
- monotonic contiguous sequence assignment;
- one revision transition per accepted batch;
- retained event/payload schema versions;
- private payload access policy.

Operational connection changes append no gameplay event.

## 8. Randomness and choice persistence

Random outcomes, shuffles and selected values are persisted before dependent state is considered durable. Replay never regenerates past randomness or choices.

## 9. Pending interactions

Persist enough deterministic state to restore:

- interaction identity and owner;
- prompt/choice schema and legal domain;
- submitted hidden values;
- effect/procedure stage;
- queue/checkpoint state;
- exact resume cursor.

Reconnect returns the same logical interaction.

## 10. Reconnect

```yaml
ReconnectRequest:
  match_id:
  last_seen_revision:
  last_seen_event_sequence:
  client_instance_id:
```

The server authenticates the principal, resolves authority, reconstructs state from snapshot plus tail, derives the authorized projection and then registers operational presence.

Presence changes never alter revision, sequence, legality, ownership, pending interaction, state hash or Mechanics input.

## 11. Concurrency

The command processor provides single-writer semantics per match.

Idempotency lookup precedes execution. New requests serialize against revision. Lifecycle transactions serialize seat allocation. Two commands cannot advance the same revision.

## 12. Replay

Replay:

1. validates snapshot hash;
2. requires the first tail sequence to equal `snapshot.event_sequence + 1`;
3. rejects gaps, duplicates, reordering and unknown event schemas;
4. applies every explicit event in order;
5. verifies state and viewer projection hashes.

No implicit mutation or fighter/card-specific dispatch is allowed.

## 13. Fixture evidence

Normative files:

```text
docs/engine/fixtures/schema-v1.json
docs/engine/fixtures/foundation-v1.json
docs/engine/fixtures/foundation-v1-transition-audit.json
```

The fixture suite resets state, authority, idempotency, IDs, RNG and choices before every case. The transition audit lists each sequence, event type, pre-hash, changed JSON pointers and post-hash. Its final hash must equal the suite final-state hash.

## 14. Recovery failures

A match is quarantined rather than repaired by guessing when a snapshot, event sequence, definition, private payload owner, pending interaction or authority record is missing or invalid.
