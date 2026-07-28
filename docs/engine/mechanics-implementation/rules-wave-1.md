# Rules Mechanics Wave 1 — semantic language and resolver

Status: implementation candidate for Issue #29. Capabilities: CAP-001, CAP-002, CAP-003 foundation, CAP-004, CAP-018.

## Boundary

The package is a deterministic domain interpreter. It reads `model.GameState` and a serializable `model.ProcedureRef`, and returns explicit `contracts.DomainEvent` values or one pending interaction. It does not import application, runtime, persistence, transport, filesystem, network, wall clock, or adapters. Dispatch is by expression, operation, and procedure kind; fighter and card IDs remain data.

## Stable semantic rules

| Rule | Trigger/evaluation | State and ordering | Choice/visibility | Events and replay | Fixtures |
|---|---|---|---|---|---|
| FX-010 / CAP-001 | A procedure begins or resumes through `RulesEngine.Resolve`. | Stages execute in source order. Captured bindings are immutable; `state` references re-read the current working state. | A choice suspends before stage operations and stores owner, legal domain, visibility, binding, and resume procedure. | Cursor, bindings, results, accepted choices, and pending state are JSON inside `ProcedureRef`. Event IDs are deterministic. | ordered stages; serialize/resume/replay |
| FX-020 / CAP-018 | A condition is evaluated when its stage is reached. | Current-state references observe prior operations in the same resolution. Captured references retain trigger-time facts. | None. | False conditions emit a deterministic skip record. | current versus captured state |
| FX-031 / CAP-001 | An ordinary operation is attempted. | An impossible independent operation records `applied=false`; later independent operations continue. | None. | Operation-result events expose disposition, not private bound values. | partial resolution |
| FX-032 / CAP-001 | A cost, prerequisite, or explicit dependency is reached. | Failed prerequisites skip the stage. An unpaid cost blocks its dependent stage. `require_applied` blocks only the declared dependent operation. | None. | Skip reasons are explicit and replayable. | failed cost and `if you do` dependency |
| FX-060 / CAP-002 | A private or opaque interaction is materialized. | The authoritative option map is captured at suspension. Submitted handles are validated server-side. | Opaque projection exposes stable handles only; the owner is checked on resume. | Accepted choice identity stays in serialized procedure state and is not copied into public provenance events. | hidden choice leak regression |
| CAP-003 foundation | A stage names a checkpoint or queued effects share a cancellation scope. | Checkpoints are source ordered. Cancellation marks only cancelable effects in the explicit scope; it is not rollback. | Reaction domains can use the same interaction type in later waves. | Checkpoint and cancellation disposition are explicit. | scoped cancellation |
| CAP-004 | Events are indexed after replay. | Index order is event sequence. Source, cause, and participant are distinct dimensions. | Indexes contain authoritative events, not client projections. | Indexes are rebuildable and never persistence authority. | cause/source/participant distinction |

## Definition validation

Loading rejects unknown expression or operation kinds, duplicate stages/bindings, forward result or choice references, unavailable dependencies, malformed choices, and unserializable definitions. The language has no arbitrary loops, reflection-based calls, embedded Go, or dynamic function names.

## Core integration

`rules.Engine` implements `contracts.RulesEngine`. Core supplies immutable state and persists returned events and pending interactions. `LegalActions` reads only a generic per-player contract. `Project` is deny-by-default: it returns public match fields and the requesting player's own state, but not card definitions or other players' private zones.

## Known Wave 1 limits

This candidate establishes the I1 resolver foundation. Ordinary game timing, zone mutation, damage, combat, movement, random procedures, and full-corpus definitions remain later waves. CAP-003 contains checkpoint/cancellation primitives but not the complete base-game timing state machine.
