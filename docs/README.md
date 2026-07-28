# Documentation map

This directory is the authoritative specification workspace for `unmatched-web`.

## Status model

Every specification area uses one of these states:

- `planned` — scope identified, research not complete;
- `researching` — source collection or reconciliation in progress;
- `drafted` — rules/data written but not fully cross-checked;
- `verified` — checked against the required source hierarchy;
- `developer-ready` — complete enough to implement without unresolved semantic questions;
- `blocked` — missing authoritative information prevents completion.

`developer-ready` is intentionally stronger than “documented.” See [`specification-readiness.md`](specification-readiness.md).

## Workstreams

| Area | Status | Primary output |
| --- | --- | --- |
| Source policy and provenance | verified | `sources/` |
| Competitive core rules | **verified — Phase 1 gate passed** | `rules/` |
| Timing, choices, effects, rulings | **verified — Phase 2 gate passed** | `mechanics/`, `rulings/` |
| Set/release registry and set-specific mechanics | **verified — Phase 3 gate passed** | `sets/` |
| Fighter/deck schema stress test | **verified — Phase 4A gate passed** | `fighters/phase-4a/`, `cards/phase-4a/` |
| Complete fighter/deck corpus | **verified — final Phase 4 gate passed** | `fighters/`, `cards/`, Phase 4 validation |
| Shared runtime requirements consolidation | next — Phase 4C | `mechanics/runtime-requirements/` |
| Battlefields | planned — MVP graph may proceed in parallel | `battlefields/` |
| Cross-validation and ambiguity closure | planned | validation reports / rulings |
| Engine-facing contract | planned | stable schemas + mechanic taxonomy |
| Adventures/co-op extension | deferred | separate future workstream |

## Current specification corpus

```text
docs/
├── rules/                  # Phase 1 complete
├── mechanics/              # Phase 2 complete; Phase 4C consolidation next
├── sets/                   # Phase 3 complete
├── fighters/               # complete 74-fighter Phase 4 corpus
│   ├── README.md
│   ├── schema.md
│   ├── phase-4a-mechanics.md
│   ├── phase-4a-validation.md
│   ├── phase-4-final-validation.md
│   ├── phase-4-final-validation.json
│   ├── phase-4a/           # ten representative fighter manifests
│   └── phase-4b/           # remaining 64 reconciled fighter manifests
│
├── cards/
│   ├── README.md
│   ├── phase-4a/           # ten representative card manifests
│   └── phase-4b/           # remaining 64 reconciled card manifests
├── phase-4b/               # Worker A/B/C/D reports and requirement definitions
├── qa/
│   └── phase-4-card-image/ # immutable reports, correction reports and integration QA
├── rulings/
│   ├── global-rulings.md
│   └── ambiguity-register.md
├── sources/
│   ├── source-policy.md
│   └── source-registry.md
├── battlefields/           # Phase 5
├── research-plan.md
└── specification-readiness.md
```

## Current boundary

Phases 1–4 now establish:

- competitive core rules and deterministic timing/choice semantics;
- the authoritative release/set work queue and source hierarchy;
- a canonical corpus of **74 competitive fighter identities**;
- **74 fighter manifests** and **74 paired card manifests**;
- **926 unique action-card definitions** representing **2214 available copies**;
- reconciled fixed and constructed deck rules;
- explicit auxiliary decks, external definitions and non-card component boundaries;
- **52 resolvable corpus-proven shared-runtime requirements**;
- preserved evidence, semantics, integration and policy status distinctions;
- reproducible machine validation with zero errors and zero warnings.

This still does **not** make the project globally `developer-ready`. Phase 4C must consolidate owner requirements into canonical generic runtime capabilities. Phase 5 must provide independently verified battlefield graphs. Phases 6–7 must close launch-scope ambiguities and freeze the developer-facing state/command/event/choice/persistence contract.

Engine foundation and one vertical-slice battlefield may proceed in parallel where the relevant contracts are already stable. Full-roster implementation must not bypass Phase 4C by adding fighter-ID conditionals or character-specific core-engine branches.
