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
| Timing, choices, effects, rulings | planned | `mechanics/`, `rulings/` |
| Set registry and set-specific rules | planned | `sets/` |
| Fighters and decks | planned | `fighters/`, `cards/` |
| Battlefields | planned | `battlefields/` |
| Cross-validation and ambiguity closure | planned | validation reports / rulings |
| Engine-facing contract | planned | stable schemas + mechanic taxonomy |
| Adventures/co-op extension | deferred | separate future workstream |

## Current rules corpus

```text
docs/rules/
├── README.md
├── terminology.md
├── battlefield.md
├── setup.md
├── turn-structure.md
├── maneuver.md
├── scheme.md
├── combat.md
├── effect-resolution-baseline.md
├── defeat-and-game-end.md
├── multiplayer-deltas.md
└── phase-1-validation.md
```

Phase 1 is complete for its declared **vanilla competitive two-player** scope. This does not make the full project `developer-ready`: Phase 2 still has to formalize pending choices, generalized timing/effect semantics, global rulings, extra/free actions, cancellation, placement edge cases, bonus attacks, and related mechanics before real fighter decks can safely become executable data.

## Planned remaining structure

```text
docs/
├── sources/
│   ├── source-policy.md
│   └── source-registry.md
├── rules/                 # Phase 1 complete
├── mechanics/             # Phase 2
├── sets/                  # Phase 3
├── fighters/              # Phase 4
├── cards/                 # Phase 4
├── battlefields/          # Phase 5
├── rulings/               # Phase 2 onward
├── research-plan.md
└── specification-readiness.md
```

The structure may evolve when the corpus demonstrates a real need. Do not introduce abstractions only to make the documentation look uniform.