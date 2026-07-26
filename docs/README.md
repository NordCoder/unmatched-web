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
| Competitive core rules | planned | `rules/` |
| Timing, choices, effects, rulings | planned | `mechanics/`, `rulings/` |
| Set registry and set-specific rules | planned | `sets/` |
| Fighters and decks | planned | `fighters/`, `cards/` |
| Battlefields | planned | `battlefields/` |
| Cross-validation and ambiguity closure | planned | validation reports / rulings |
| Engine-facing contract | planned | stable schemas + mechanic taxonomy |
| Adventures/co-op extension | deferred | separate future workstream |

## Planned structure

```text
docs/
├── sources/
│   ├── source-policy.md
│   └── source-registry.md
├── rules/
│   ├── terminology.md
│   ├── setup.md
│   ├── turn-structure.md
│   ├── maneuver.md
│   ├── combat.md
│   ├── schemes.md
│   ├── exhaustion.md
│   └── game-end.md
├── mechanics/
│   ├── timing.md
│   ├── effects.md
│   ├── choices.md
│   ├── movement-vs-placement.md
│   ├── bonus-attacks.md
│   ├── multiple-heroes.md
│   ├── summoning.md
│   ├── ongoing-schemes.md
│   └── battlefield-mechanics.md
├── sets/
├── fighters/
├── cards/
├── battlefields/
├── rulings/
├── research-plan.md
└── specification-readiness.md
```

The structure may evolve when the corpus demonstrates a real need. Do not introduce abstractions only to make the documentation look uniform.
