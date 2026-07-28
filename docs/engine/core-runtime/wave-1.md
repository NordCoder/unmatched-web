# Core Runtime Wave 1 — Running match host

## Candidate scope

This Wave 1 implementation provides an in-memory, server-authoritative host around
the Lead-owned `model.GameState` and `contracts.RulesEngine` boundary.

Delivered behavior:

- immutable definition bundles and pinned `DefinitionRef` identity;
- deterministic injected allocation for match, player, fighter, card, action,
  procedure, and event instances;
- `CreateMatch` and `JoinMatch` with durable Principal → Player authority bindings;
- command-envelope validation, required expected revision, and coarse idempotency;
- one atomic in-memory event batch per accepted command;
- replay from ordered events into identical hosted state;
- generic Maneuver, Scheme, and Attack procedure shells;
- RulesEngine invocation, persisted pending interaction, owner-authorized choice,
  and resume of the exact serialized `ProcedureRef`;
- player-specific projection orchestration with owner-only pending prompts;
- transport-neutral projection encoding.

## Event and transaction boundary

An accepted command produces one `EventBatch`. Every event in the batch shares the
next match revision and has an increasing sequence. The store atomically checks:

1. global command-id idempotency;
2. expected current match revision;
3. batch metadata consistency;
4. append plus stored command result.

A duplicate command with the same principal and envelope returns the original
serialized result. Reuse of a command ID with different input is rejected.

## Rules boundary

Core recognizes only its own lifecycle envelopes. Events returned by Rules are
stored and ordered as authoritative history, but Core does not interpret fighter,
card, or battlefield identities and does not implement card-stage semantics.

Wave 1 intentionally does not select network transport, PostgreSQL libraries,
snapshots, authentication technology, or public schema generation. Those remain
Wave 2 decisions.

## I1 evidence

`tests/integration/core/wave1_test.go` covers the complete checkpoint:

```text
register synthetic definitions
→ create match and first authority binding
→ join mirror fighter and allocate distinct repeated instances
→ reject stale revision
→ start Scheme procedure
→ persist owner-only PendingInteraction
→ reject another player's choice
→ resume exact ProcedureRef
→ append Rules event and complete action
→ return identical result for duplicate command
→ replay full event stream with zero state divergence
```
