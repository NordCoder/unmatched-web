# Engine Mechanics Analysis

## Status

```text
status: analysis-draft
branch: engine-mechanics-analysis
issue: #25
engine_foundation_base: bc075caef767aa1ccf616260bca6c5d48c239734
canonical_corpus: 106ae552ce597cde954c0a1b22374ef446974ce2
phase4c_assignment: be78417eefeb23cd334c796ba42f192272b8c9a2
fighters: 74
unique_action_card_definitions: 926
owner_requirement_aliases: 52
```

This work turns the normalized competitive fighter/card corpus into an implementation-oriented runtime design. It complements the canonical Phase 4C requirement dispositions; it does not replace them or edit fighter/card manifests.

## Corpus authority

The analysis is grounded in:

- the Phase 4 final corpus validation: 74 fighter manifests, 74 card manifests and 926 unique action-card definitions;
- the normalized effect vocabulary in `docs/mechanics/effect-model.md`;
- the 52 corpus-proven `A/B/C/D-REQ-*` integration requirements;
- the engine definition/instance, command/event, choice, visibility and persistence contracts;
- representative stress tests from the most mechanically demanding manifests.

The owner requirements are important because ordinary cards are already expressible using the baseline operations. The 52 aliases identify the places where a generic runtime needs additional state, timing, legality, visibility, history or procedure semantics.

## Primary conclusion

The engine should be a **typed deterministic rules interpreter**, not a collection of fighter/card handlers.

```text
immutable definitions
        ↓
pure selectors / expressions
        ↓
serializable procedures and stages
        ↓
typed operations
        ↓
persisted events and operation results
        ↓
authoritative state + player projections
```

Core code may dispatch by:

- operation kind;
- query/expression kind;
- timing window/checkpoint;
- generic capability kind;
- registered typed composite procedure.

Core code must not dispatch by fighter or card identity.

## What the corpus actually makes difficult

The corpus does not primarily demand hundreds of unrelated operations. Most published effects are combinations of a stable operation vocabulary such as draw, discard, reveal, move, place, damage, recover, change state, change value, gain action and cancel effects.

The difficult dimensions are cross-cutting:

1. **Control flow** — ordered stages, optional branches, explicit dependencies, repeat/stop procedures and empty-domain handling.
2. **Timing** — attack declaration, pre-defense, combat windows, cleanup, turn checkpoints and protected nested reactions.
3. **Bindings/history** — later stages use exact prior cards, fighters, random results, movement paths, damage results or historical events.
4. **Visibility** — private choices, field-only disclosure, opaque hidden-card handles and delayed reveal.
5. **Identity/lifecycle** — definition versus runtime instance, off-board undefeated fighters, components, auxiliary cards and repeated sidekicks.
6. **Derived legality** — attack modes, range, movement costs, hand policies and continuous permissions that change without mutating base definitions.
7. **Transactional results** — damage allocation, exact health assignment, partial resource discharge, action credits and cancellation scope.
8. **Persistence** — every paused choice, procedure cursor, accepted random result, obligation and source-bound modifier must survive reconnect.

## Stress-test coverage

| Mechanic class | Representative corpus stress test | Engine implication |
|---|---|---|
| Resumable procedure | Pandora Misery reveal/resolve/stop loop | persisted procedure cursor and forced-stop predicate |
| Multi-space occupancy | T. Rex | footprint-aware movement, targeting, adjacency and turn-start history |
| Shared-space occupancy | Squirrel Girl | occupancy classes/capacity plus propagated-damage provenance |
| Same-combat rebinding | Yennenga Surprise Volley | combat identity survives attacker replacement |
| Ordered persistent zone | Shakespeare Line | ordered zone metric, completion channel and cleanup against live membership |
| Dynamic randomness/BOOST | Krang | persisted random results, paid rerolls and unified BOOST resolution |
| Hidden selection | Loki / Spider-Man | opaque handles and field-only disclosure |
| Component/path lifecycle | Muldoon, John Henry, Shredder | stable component instances anchored to spaces or paths |
| Off-board without defeat | Invisible Man, Elektra, Ancient Leshen | presence state independent from defeat with scheduled return |
| Damage transaction | Yennenga, Annie Christmas, Squirrel Girl | allocation, provenance and health-result modifiers |
| Action/obligation model | Ms. Marvel, Stars and Stripes, Shredder | typed action credits and durable delayed obligations |
| External policy | Deadpool | explicit policy boundary, not an invented runtime behavior |

## Runtime capability registry

`capability-map.json` maps every one of the 52 owner aliases to exactly one primary implementation capability. Cross-domain relationships are recorded separately and may repeat without changing primary ownership.

The current map contains 20 capabilities grouped around:

- staged resolution and interactions;
- timing/cancellation and provenance;
- randomness;
- fighter presence and occupancy;
- battlefield components and relocation;
- combat legality and rebinding;
- damage, modifiers and BOOST;
- ordered/searchable card zones;
- actions, resources, obligations and hand policy;
- pure derived queries;
- external adaptation policy.

## Design restrictions

The runtime definition format must remain deliberately non-Turing-complete:

- no arbitrary loops; use bounded typed procedures such as repeat/stop or repeat-for-each;
- no arbitrary code expressions; use a closed pure expression/query AST;
- no mutable global variables; state belongs to typed runtime instances;
- no direct database/network/filesystem access from domain resolution;
- no dynamic invocation by fighter/card ID;
- no hidden side effects outside emitted events;
- no recomputation of accepted random or private choice results during replay.

A custom handler is acceptable only after the existing extension-handler gate proves that a generic capability cannot preserve deterministic semantics. The handler must still expose typed inputs, outputs, state, interactions, visibility and events.

## Artifacts

- `capability-map.json` — machine-readable 52/52 capability coverage.
- `runtime-blueprint.md` — proposed kernel, state machines and package boundaries.
- `implementation-waves.md` — dependency-safe delivery order from launch kernel to full corpus.
- `scripts/validate_engine_mechanics.py` — coverage and structural validator.

## Relationship to Phase 4C

Phase 4C remains authoritative for the final wording and disposition of each source alias. This analysis provides an implementation convergence model:

- worker fragments may refine or split semantic requirements;
- the capability map must then be updated without losing any alias;
- a canonical Phase 4C requirement may map to one capability or a typed combination of capabilities;
- capability IDs must not become disguised fighter-specific branches.

## Merge gate

- all 52 aliases have exactly one primary capability mapping;
- capability dependencies are acyclic;
- each capability states required state, queries, operations, events, interactions and persistence behavior;
- no fighter/card identity dispatch is proposed;
- Phase 4C reconciliation confirms semantic compatibility;
- fresh architecture QA passes.
