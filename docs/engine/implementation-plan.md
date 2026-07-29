# Engine implementation plan

## Status

```text
owner_policy: coarse implementation first
lead: ChatGPT Web Lead
integration_branch: engine-runtime-integration
worker_branches:
  core: engine-core-runtime
  rules: engine-rules-mechanics
quality_mode: milestone-based
rules_mode: strict semantics and provenance
```

The project now moves from foundation refinement to implementation. The existing
foundation is a starting contract, not a claim that every implementation detail
is already settled. The Lead resolves cross-line conflicts while work proceeds.

## Delivery principles

1. Build in large coherent waves, not micro-PRs for every type or ambiguity.
2. Working end-to-end behavior is preferred over premature framework polish.
3. Internal APIs may be reshaped between milestones by the Lead.
4. Core Runtime must not contain fighter/card-specific dispatch or interpret card text.
5. Rules Mechanics must not access persistence, transport, filesystem or client state.
6. Published rule behavior remains strict: rule IDs, timing, dependencies, visibility,
   state transitions, uncertainty and provenance must be explicit.
7. An unsupported fighter is reported as partial/unsupported; it is not implemented
   through an identity-specific shortcut.
8. QA is performed at cross-line milestones and before promotion to main, not after
   every small internal commit.

## Lead-owned boundary

The Lead owns:

```text
internal/domain/model/**
internal/domain/contracts/**
tests/integration/crossline/**
docs/engine/implementation-plan.md
engine-runtime-integration
```

Workers request shared changes with a concrete failing use case. The Lead may accept,
reshape or defer the change and publishes the decision on the integration branch.
Unrelated work continues while the conflict is being resolved.

# Core Runtime plan

## Core Wave 1 — Running match host

Build one usable in-memory authoritative runtime:

- absorb the useful Go/server scaffold from historical PR #24;
- definition registry and immutable definition identity;
- CreateMatch, JoinMatch and player authority bindings;
- runtime instance allocation for players, fighters and cards;
- command dispatcher, expected revision and idempotency;
- atomic in-memory event batches and reducer replay;
- turn/action shells for Maneuver, Scheme and Attack;
- invocation and resume of the shared RulesEngine port;
- player projection orchestration and pending-interaction delivery;
- broad integration tests rather than production-grade adapters.

Exit milestone `I1`:

```text
create match
join match
load synthetic definitions
start an action
Rules returns a pending interaction
Core stores it
submit choice
resume the same procedure
replay reaches identical state
```

## Core Wave 2 — Durable online runtime

Turn the in-memory host into a recoverable service:

- PostgreSQL event, snapshot, authority and idempotency repositories;
- transaction boundary for command result + events + match head;
- snapshot policy and snapshot-plus-tail reconstruction;
- reconnect service preserving pending interactions;
- legal-action and player-projection application APIs;
- minimal HTTP/WebSocket transport selected pragmatically by the Core worker;
- health/readiness and development Compose integration;
- failure-path and restart integration tests.

Exit milestone `I2`:

```text
match survives process restart
pending choice survives reconnect
private projection remains private
accepted command is durable before acknowledgement
```

## Core Wave 3 — Playable product slice

Assemble the first real browser-playable match:

- load canonical Robin Hood and Bigfoot data;
- load approved battlefield topology when available;
- setup, starting hands and player projections;
- command/query transport used by the TypeScript client;
- game lifecycle through exhaustion, defeat and game end;
- basic observability and operational cleanup;
- no performance or horizontal-scaling optimization unless required by correctness.

Exit milestone `I3`:

```text
Robin Hood vs Bigfoot is playable end-to-end as data
full match replay has zero divergence
no launch identity occurs in core control flow
```

# Rules Mechanics plan

## Rules Wave 1 — Semantic language and resolver

Implement the deterministic language used by all card/fighter definitions:

- closed selector/query/expression AST;
- serializable staged procedures and stage cursor;
- immutable captured bindings and explicit current-state reads;
- typed operations and operation-result bindings;
- pending interactions, ownership, visibility and resume;
- effect queue, timing checkpoints and scoped cancellation;
- event/history/provenance ledger;
- rule-definition validator and fixture runner;
- detailed rule registry mapping rule IDs and capability IDs to implementation.

Rules are documented before or together with implementation. Each rule entry states:

```text
semantic rule ID
source/provenance
trigger or evaluation point
inputs and authoritative state read
ordering and dependency behavior
choices and owner
operations/events produced
visibility
serialization/replay requirements
known uncertainty or policy decision
fixtures proving the rule
```

Exit milestone `I1` is shared with Core and proves pause/serialize/resume/replay.

## Rules Wave 2 — Complete ordinary Unmatched kernel

Implement the complete base game carefully, not merely enough for one fighter:

- setup ordering and lifecycle hooks;
- turn start/end and two-action accounting;
- Maneuver, movement, BOOST and exhaustion;
- Scheme play and ordered effect resolution;
- Attack declaration, defense, card reveal and combat values;
- immediately/during/after combat timing windows;
- damage, prevention, recovery, defeat and game-end checkpoints;
- card ownership/control/zones/visibility;
- adjacency, zones, path movement and ordinary placement;
- partial resolution, costs, prerequisites and `if you do` dependencies;
- ordinary ongoing, replacement and cancellation semantics;
- reconnect-safe pending choices at every interruptible point.

The rule documents under `docs/mechanics/**` remain semantic authority. Where code
reveals a conflict or gap, the Rules worker records the exact case and the Lead chooses
one canonical behavior. Uncertain source behavior remains explicitly marked uncertain.

Exit milestone `I2`:

```text
Maneuver, Scheme and defended/undefended Attack work generically
all base timing windows are explicit
exhaustion, defeat and game end work
hidden defense and choices do not leak
```

## Rules Wave 3 — Corpus capability expansion

Expand generic capabilities according to corpus pressure rather than set order:

- off-board/dormant fighters and scheduled return;
- occupancy classes, multi-space footprints and shared spaces;
- battlefield components, traps, paths and entry interruption;
- combat participant/card replacement;
- damage allocation, propagation, health floors and exact assignment;
- source-bound continuous modifiers and value layers;
- ordinary, blind and dynamic BOOST;
- authoritative randomness and paid rerolls;
- typed action credits and restricted immediate actions;
- resources, delayed obligations and multi-actor payments;
- auxiliary/ordered zones, structured search and disclosure;
- hand-size policy and other corpus-wide derived rules.

Conformance goal:

```text
all competitive fighter/card manifests parse
all deterministic requirements map to typed capabilities
no unknown operation/query/procedure kind
no fighter/card-ID dispatch
all pending procedures serialize
all accepted random and private-choice results replay
```

Deadpool and genuinely policy-dependent behavior remain explicit boundaries rather
than invented semantics.

# Milestone integration

Only three major cross-line gates are planned:

| Gate | Meaning |
| --- | --- |
| `I1` | Core host + staged resolver + pending choice + replay |
| `I2` | ordinary full game kernel + persistence/reconnect/private projections |
| `I3` | data-loaded playable launch slice and representative advanced mechanics |

Between these gates, workers may commit freely within owned paths. The Lead integrates
compile-safe increments when useful and resolves contract conflicts directly.
