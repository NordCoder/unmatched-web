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
| Complete fighter/deck corpus | planned — Phase 4B | `fighters/`, `cards/` |
| Battlefields | planned | `battlefields/` |
| Cross-validation and ambiguity closure | planned | validation reports / rulings |
| Engine-facing contract | planned | stable schemas + mechanic taxonomy |
| Adventures/co-op extension | deferred | separate future workstream |

## Current specification corpus

```text
docs/
├── rules/                  # Phase 1 complete
├── mechanics/              # Phase 2 complete; effect taxonomy extended by Phase 4A
├── sets/                   # Phase 3 complete
├── fighters/               # Phase 4A complete; Phase 4B pending
│   ├── README.md
│   ├── schema.md
│   ├── phase-4a-mechanics.md
│   ├── phase-4a-validation.md
│   └── phase-4a/
│       ├── achilles.yaml
│       ├── bloody-mary.yaml
│       ├── sun-wukong.yaml
│       ├── sherlock-holmes.yaml
│       ├── dracula.yaml
│       ├── raptors.yaml
│       ├── wayward-sisters.yaml
│       ├── geralt-of-rivia.yaml
│       ├── yennefer-triss.yaml
│       └── black-panther.yaml
│
├── cards/
│   ├── README.md
│   └── phase-4a/           # full normalized card manifests for the same 10 fighters
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

Phases 1–4A now establish:

- the competitive core rules;
- deterministic timing/effect/choice semantics;
- global rulings and ambiguity containment;
- the exhaustive released product/content work queue;
- authoritative set-mechanic entry points;
- a verified fighter/deck manifest schema;
- fixed and pre-game-constructed deck modeling;
- multi-hero/selectable-hero/summonable-sidekick topology;
- first-class alternate card zones and immutable card ownership;
- external definitions such as bonus attacks and Wayward spells;
- corpus-proven generic primitives for damage prevention/redirection, printed/BOOST value modification, operation prevention, card reordering, captured parent context and combat-card replacement.

This still does **not** make the full fighter roster `developer-ready`: Phase 4B must apply the validated schema to every released competitive fighter in `sets/registry.yaml`. Phase 5 must turn each supported battlefield into a verified graph.

The structure may evolve when the full corpus demonstrates a real need. New abstractions must be justified by published gameplay semantics rather than introduced only for uniformity.