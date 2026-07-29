# Core Runtime Wave 1 — Running match host

## Candidate scope

This Wave 1 implementation provides an in-memory, server-authoritative host around
the Lead-owned `model.GameState` and `contracts.RulesEngine` boundary.

Delivered behavior:

- immutable definition bundles whose pinned `DefinitionRef` exactly covers the
  fighter construction data and referenced card definitions;
- deterministic injected allocation for match, player, fighter, card, action,
  procedure, interaction, and event instances;
- `CreateMatch` and `JoinMatch` with Principal → Player authority records stored
  outside deterministic gameplay state and gameplay history;
- principal-scoped command idempotency, expected-revision handling, and one
  serialized execution per `(principal_id, command_id)`;
- one atomic in-memory event batch per accepted command;
- replay from ordered gameplay events into identical `GameState`;
- generic Maneuver, Scheme, and Attack procedure shells;
- RulesEngine invocation, persisted pending interaction, owner-authorized choice,
  and resume of the exact serialized `ProcedureRef`;
- player-specific projection orchestration with an explicit allow-list for the
  owner-visible pending interaction;
- transport-neutral projection encoding.

## Event, command, and authority transaction boundary

Before command preparation, the store reserves the principal-scoped command key.
A concurrent same-fingerprint retry waits for the first execution and receives its
immutable result. It cannot invoke Rules, allocate IDs, or construct another event
batch. A different fingerprint under the same key receives a derived conflict.

An accepted lifecycle command commits under one store lock:

1. the ordered gameplay event batch;
2. the match revision/head;
3. the immutable command result;
4. the Principal → Player authority record.

Authority records are application/security data. Principal IDs never enter
`GameState`, gameplay event payloads, Mechanics input, or gameplay replay.
`PlayerState.authority_state` remains deterministic gameplay authority and starts
as `ACTIVE`.

Deterministic operational rejections are retained as immutable command results.
Infrastructure failures abort the command reservation rather than creating false
durable evidence.

## Definition identity

Registration validates that:

- every fighter is covered by a digest derived from its exact construction input,
  including ordered repeated physical card copies;
- every referenced unique card definition has a pinned manifest digest;
- no extra fighter/card digest is present;
- one `DefinitionRef` resolves to exactly one registered bundle.

`ResolveRef` uses a deterministic content-key index rather than iterating a Go map.
A changed fighter/deck construction cannot be registered under an unchanged pinned
reference.

## Projection boundary

The authoritative `model.PendingInteraction` remains inside gameplay state because
it contains the serialized resume procedure. Delivery constructs a separate
`ProjectedInteraction` containing only:

```text
interaction_instance_id
owner_player_id
kind
visibility
prompt
legal_domain
```

`ResumeProcedure`, bindings, stage cursor, and other execution state are not
representable in the transport-facing projection. The non-owner receives only the
blocked/waiting state without private prompt or legal-domain data.

## Rules boundary

Core recognizes only its own lifecycle envelopes. Events returned by Rules are
stored and ordered as authoritative history, but Core does not interpret fighter,
card, or battlefield identities and does not implement card-stage semantics.

Wave 1 intentionally does not select network transport, PostgreSQL libraries,
snapshots, authentication technology, or public schema generation. Those remain
Wave 2 decisions.

## I1 and correction evidence

`tests/integration/core/wave1_test.go` covers:

```text
register content-pinned synthetic definitions
→ create match and atomically bind first authority record
→ join mirror fighter and allocate distinct repeated instances
→ reject and replay stale revision deterministically
→ start Scheme procedure
→ persist owner-only PendingInteraction
→ reject another player's choice
→ resume exact ProcedureRef
→ append Rules event and complete action
→ return identical result for duplicate command
→ replay gameplay events with zero state divergence
→ prove no principal identity occurs in gameplay events
```

Focused correction evidence additionally proves:

- encoded owner projection excludes resume procedure, procedure ID, bindings, and
  captured internal values;
- equal command IDs are independent across principals;
- concurrent same-principal duplicates invoke Rules once and allocate IDs once;
- conflicting construction data cannot reuse a pinned definition reference;
- definition references resolve through a deterministic one-to-one index.
