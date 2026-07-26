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
├── sets/                   # Phase 3 complete
│   ├── README.md
│   ├── registry.yaml
│   ├── mechanics-index.md
│   ├── source-bibliography.md
│   └── phase-3-validation.md
│
├── rulings/
│   ├── global-rulings.md
│   └── ambiguity-register.md
│
├── sources/
│   ├── source-policy.md
│   └── source-registry.md
│
├── fighters/               # Phase 4
├── cards/                  # Phase 4
├── battlefields/           # Phase 5
├── research-plan.md
└── specification-readiness.md
```

## Current boundary

Phases 1–3 now establish:

- the competitive core rules;
- deterministic timing/effect/choice semantics;
- global rulings and ambiguity containment;
- the exhaustive released product/content work queue;
- fighter-to-release and battlefield-to-release membership;
- authoritative entry points for known set-specific mechanics;
- explicit handling of reprints, addenda, supplements and announced-but-blocked content.

This still does **not** make the actual fighter roster `developer-ready`. Phase 4 must transcribe and validate fighter/deck content against the Phase 2 effect framework. Phase 5 must turn each supported battlefield into a verified graph.

The structure may evolve when the corpus demonstrates a real need. Do not introduce abstractions only to make the documentation look uniform.