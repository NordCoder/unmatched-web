# Rules Mechanics Wave 1 — sound I1 resolver foundation

Status: manual Wave 1 safety-closure Candidate for Issue #46 / PR #40, continuing exact Head `a7021602a8e8cab2b7920eac09098031186c9fad`.

Capabilities: `CAP-001`, `CAP-002`, `CAP-003`, `CAP-004`, `CAP-018`.

## Runtime boundary

The Rules package remains a deterministic domain interpreter. It receives immutable `model.GameState` and a serializable `model.ProcedureRef`, then returns explicit events, a pending interaction, or a deterministic rejection. It imports no adapter, storage, transport, clock, network, or fighter/card-specific dispatch.

`RulesEngine.Project` is an allow-list and never serializes procedure bindings, resolver internals, pending authoritative options, raw action/combat procedures, or canonical `RuntimeObject.State`. Fighter projection uses a sorted explicit public DTO containing only identity, definition, owner, and controller metadata.

## Authoritative state references

State references are accepted only through the engine-owned `StatePathSchema` registry. Each registered pattern fixes:

```text
path pattern
structural ValueSpec
visibility
allowed traversal
```

Definition-provided `ValueType` and `Visibility` are assertions only. Unknown paths fail definition loading; private paths cannot be downgraded to public; arbitrary traversal through maps or resolver/private-zone state is forbidden. Wave 1 intentionally uses a small allow-list covering the committed fixtures rather than pretending to model the entire future `GameState`.

## Recursive structural typing

`ValueSpec` carries:

```text
type
element
fields
visibility
optional
```

Nested references traverse the schema one segment at a time. Missing fields and terminal type overrides fail at load time. Operation result bindings expose a common structural record:

```text
applied: bool
disposition: disposition
code: string?
value: operation-specific ValueSpec?
```

`set_fighter_state`, for example, exposes a typed patch object whose nested `value` type is derived from the operation argument. Runtime trust boundaries enforce external captured/input values before evaluation. Binding sets are closed: unknown and missing bindings, wrong nested values, undeclared object fields, open `any` contracts, and conflicting captured/input schemas are rejected. Validated values are canonical-copied before evaluator use and procedure serialization. Current Core transport duplication of action context into both procedure bindings and `ResolutionInput.Context` is validated for equality and preserved across pause/resume.

## Choice and empty-domain semantics

Wave 1 is deliberately **single-select only**. Single-choice value type comes from `domain.Element`; `Choice.ValueType` is assertion-only. Choice visibility is never less restrictive than domain visibility.

Supported policies:

```text
reject
bind_default
skip_stage
complete_without_choice
```

`Choice.Multi` and `bind_empty` fail definition loading as explicitly unsupported. Multi-select cardinality, payload, ordering, duplicate handling, and list binding are deferred to a separate capability. `bind_default` requires an exact typed default. `skip_stage` and `complete_without_choice` do not establish a binding; conservative dataflow validation rejects any continuing expression that assumes the binding exists.

## Serializable procedure and pure replay

The stage, phase, cost index, operation index, captured values, result records, accepted choices, pending interaction, and checkpoint queue remain serialized in `ProcedureRef`.

`ApplyEvent` deep-clones the complete input state before applying a rules mutation. No output map, slice, raw message, or nested object aliases the input. Repeated reduction from identical bytes is deterministic.

The fixture runner now proves the complete path:

```text
pristine initial bytes
-> initial resolve
-> reduce every initial event
-> serialize/restore pending procedure
-> resume
-> reduce every resume event
-> compare direct, serialized, full replay, and repeated final states
```

It also asserts that the original fixture object remains byte-identical. A dedicated pause/resume fixture mutates fighter state only after choice acceptance; a negative proof omits that resume mutation and requires replay divergence.

## Unified operation dispositions

Ordinary operations, costs, and queued operations use one internal attempt primitive and one result schema. Authoritative dispositions are:

```text
applied
not_applied
skipped_dependency
rolled_back_cost
canceled
partial
```

Every attempted operation emits one `rules.operation_result` event. Dependencies consume disposition rather than an ad hoc boolean.

Costs execute against an isolated candidate state. When any cost fails, mutation events and patches do not commit; previously successful candidate attempts become `rolled_back_cost`, the failed attempt remains `not_applied`, and the stage emits explicit `cost_unpaid` behavior.

Queue lifecycle is:

```text
queued -> executing -> applied | partial | not_applied | canceled
```

`Next()` only marks an entry `executing`. The resolver computes the real outcome after all contained attempts and emits `rules.queued_effect_outcome`; an impossible queue operation can no longer be reported as successfully resolved.

## Event and provenance contract

Rules events use `rules-event/v1` with explicit source, cause, operation instance, participants, and typed data. History indexes are rebuilt from this actual envelope. Mutation events are emitted only for committed state changes; standardized result events remain visible for failed, skipped, rolled-back, canceled, and partial attempts.

## Acceptance evidence

Focused regressions cover:

- unknown/private state path rejection and public-event taint prevention;
- recursive nested result paths and terminal type assertions;
- choice domain type, visibility, typed default, explicit multi-select/bind-empty rejection, and presence dataflow;
- strict external captured/input binding validation, including adversarial undeclared fields and Core context duplication/resume;
- fighter projection denial for arbitrary canonical state;
- reducer non-aliasing and repeated deterministic reduction;
- full initial-plus-resume event replay from pristine bytes, including a resume-only mutation divergence proof;
- ordinary, cost, and queue disposition equivalence;
- rolled-back costs, impossible and partial queue effects, and canceled non-execution;
- existing opaque projection/resume, ordered stages, dependencies, provenance, and checkpoint ordering.

Machine-readable evidence:

```text
tests/fixtures/mechanics/schema-v1.json
tests/fixtures/mechanics/i1.json
tests/domain/mechanics/fixture_runner_test.go
```

## Remaining Lead-owned integration item

Issue #30 comment `5108015325` remains separate from this correction: Core must persist Rules events returned with a pending outcome and connect the accepted Rules reducer during Core replay. No Core- or Lead-owned path is changed here.
