# Rules Mechanics Wave 1 — corrected I1 resolver foundation

Status: correction Candidate for Issue #29 after the QA verdict on `d8788cc1428431750a87e40d87b25598db44e985`.

Capabilities: `CAP-001`, `CAP-002`, `CAP-003`, `CAP-004`, `CAP-018`.

## Runtime boundary

The package is a deterministic domain interpreter. It receives immutable `model.GameState` plus a serializable `model.ProcedureRef` and returns explicit `contracts.DomainEvent` values or a pending interaction. It imports no application, runtime, persistence, transport, filesystem, network, SQL, wall clock, or adapter package. Dispatch is by expression, operation, stage, checkpoint, and capability kind; fighter/card identity is data only.

`RulesEngine.Project` is an explicit allow-list. It never serializes `ProcedureRef`, `ResolverState`, `PendingInteraction`, raw action/combat procedure bindings, or arbitrary `Turn` maps. Core remains responsible for projecting the owner-only pending-interaction DTO.

Core authenticates and authorizes `SubmitChoice` before calling Rules. Rules validates the interaction identity and the authoritative opaque handle captured in the serialized procedure; it does not accept caller identity through arbitrary rule `Context`.

## Serializable program counter

`effects.State` version 2 persists:

```text
stage cursor
phase: enter | costs | choice | operations | checkpoint | complete
cost index
operation index
captured bindings
operation results
accepted choices
pending authoritative option map
checkpoint frame and ordered queue
```

A choice resumes at `operations`, not at stage entry. Costs are evaluated transactionally against a clone and committed to the outcome only when the whole cost set succeeds. A pending outcome therefore contains the exact pre-choice events plus a post-cost resume procedure. Core must persist those events atomically with the interaction-pending event; the corresponding shared-contract request is Issue #30 comment `5108015325`.

## Typed definition language

Expressions carry closed kinds and statically inferred value/visibility specifications:

```text
any, bool, number, string, player_ref, fighter_ref,
list<T>, object, operation_result

public, owner_private, opaque
```

Load-time validation enforces:

- operator arity and operand/result types;
- typed current-state and binding references;
- operation argument schemas and unknown-argument rejection;
- cost eligibility;
- dependency availability and unique result/choice bindings;
- choice owner/domain types, visibility enum, and explicit empty-domain policy;
- public-event rejection when an argument carries owner-private or opaque data;
- declared capability existence and dependency closure;
- ordered queue entry validity.

Unsupported or untyped definitions fail loading instead of reaching a runtime type assertion.

## Checkpoints, queues, and cancellation

A checkpoint frame stores a stable queue sorted by:

```text
priority descending
source order ascending
stable queue ID ascending
```

Queue entries contain typed generic operations and survive JSON serialization. Cancellation marks only queued, cancelable entries in the explicit scope. The resolver then executes remaining entries in deterministic order and emits provenance-bearing cancellation and resolution events. Cancellation never rolls back an already applied operation.

## Event and history contract

Every Rules event uses `rules-event/v1`:

```json
{
  "schema": "rules-event/v1",
  "data": {},
  "provenance": {
    "operation_instance_id": "...",
    "source_ref": "...",
    "cause_kind": "...",
    "participants": []
  }
}
```

The history ledger indexes the exact envelope produced by the resolver. It no longer assumes arbitrary payloads contain undocumented `cause_kind` or `participants` fields. `Engine.ApplyEvent` is the worker-owned deterministic reducer for currently implemented Rules mutation events; Lead integration must connect the accepted shared replay port.

## Rule mapping and evidence

| Rule / capability | Evaluation and ordering | Visibility / replay evidence |
|---|---|---|
| `FX-010 / CAP-001` | Source-ordered stages with serialized phase/cursors | ordered fixture; JSON restore; deterministic rerun |
| `FX-020 / CAP-018` | current-state reads happen when reached; captured bindings remain immutable | current-versus-prior-result fixture |
| `FX-031 / CAP-001` | impossible independent operation is recorded and later independent work continues | partial-resolution fixture |
| `FX-032 / CAP-001` | costs are transactional; explicit dependency gates only its dependent operation | failed dependency and cost-before-choice tests |
| `FX-060 / CAP-002` | owner domain captured once; opaque projection exposes handles only | no-context Core-compatible resume and projection leak regression |
| `CAP-003` | ordered serialized queue with scoped cancellation | checkpoint fixture and encode/decode test |
| `CAP-004` | source, cause, operation, and participant indexes derive from actual events | resolver-produced provenance regression |

Machine-readable acceptance assets:

```text
tests/fixtures/mechanics/schema-v1.json
tests/fixtures/mechanics/i1.json
tests/domain/mechanics/fixture_runner_test.go
```

The fixture runner proves ordered stages, correct-time conditions, prior-result binding, dependency semantics, independent partial resolution, private pause, procedure JSON restore, exactly-once post-cost resume, checkpoint ordering/cancellation, explicit events, deterministic rerun, event reduction, and final state assertions.

## Remaining Lead-owned integration item

Worker-owned findings are corrected. Complete cross-line I1 remains gated on the additive Lead/Core change requested in Issue #30 comment `5108015325`:

1. persist `ResolutionOutcome.Events` for pending outcomes in the same atomic batch as `core.interaction.pending`;
2. connect a deterministic Rules-event reducer to Core replay;
3. optionally add an explicit authority-bound submitting-player field if future non-actor-owned interactions require Rules-side identity revalidation.
