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
| Set registry and set-specific rules | planned | `sets/` |
| Fighters and decks | planned | `fighters/`, `cards/` |
| Battlefields | planned | `battlefields/` |
| Cross-validation and ambiguity closure | planned | validation reports / rulings |
| Engine-facing contract | planned | stable schemas + mechanic taxonomy |
| Adventures/co-op extension | deferred | separate future workstream |

## Current specification corpus

```text
docs/
├── rules/                  # Phase 1 complete
│   ├── README.md
│   ├── terminology.md
│   ├── battlefield.md
│   ├── setup.md
│   ├── turn-structure.md
│   ├── maneuver.md
│   ├── scheme.md
│   ├── combat.md
│   ├── effect-resolution-baseline.md
│   ├── defeat-and-game-end.md
│   ├── multiplayer-deltas.md
│   └── phase-1-validation.md
│
├── mechanics/              # Phase 2 complete
│   ├── README.md
│   ├── event-model.md
│   ├── choices-and-resume.md
│   ├── effect-model.md
│   ├── cancellation.md
│   ├── action-accounting.md
│   ├── movement-and-placement.md
│   ├── bonus-attacks.md
│   ├── end-turn-and-dormancy.md
│   ├── information-visibility.md
│   ├── setup-hooks.md
│   └── phase-2-validation.md
│
├── rulings/
│   ├── global-rulings.md
│   └── ambiguity-register.md
│
├── sources/
│   ├── source-policy.md
│   └── source-registry.md
│
├── sets/                   # Phase 3
├── fighters/               # Phase 4
├── cards/                  # Phase 4
├── battlefields/           # Phase 5
├── research-plan.md
└── specification-readiness.md
```

## Current boundary

Phases 1–2 now define the global competitive resolution framework: core actions, action boundaries, timing windows, resumable choices, visibility, effect categories, cancellation, gained/free actions, movement versus placement, bonus attacks, `End the turn`, dormancy, and setup hooks.

This still does **not** make the actual fighter roster `developer-ready`. Phase 3 must establish the exhaustive current set/release registry and authoritative set-specific rules. Phase 4 must then map every fighter/card into the Phase 2 effect framework, retaining content-specific rulings and explicit blockers where authority is missing.

The structure may evolve when the corpus demonstrates a real need. Do not introduce abstractions only to make the documentation look uniform.
