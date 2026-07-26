# Research and specification plan

## Goal

Produce a **necessary and sufficient competitive Unmatched specification for implementation**: a developer should be able to build the game-state machine, effect resolver, legal-action generator, fighter/deck loader, and battlefield model without guessing at rules.

The work is ordered by semantic dependency, not by release date. A large card catalog is not useful if timing, choices, combat resolution, or precedence are still ambiguous.

---

## Phase 0 — Foundation and provenance

**Status:** complete in the initial documentation commit.

### Work

1. Establish the source hierarchy.
2. Record canonical entry points and freshness expectations.
3. Define how conflicts, errata, fan sources, and missing data are handled.
4. Define the `developer-ready` gate.

### Gate

- Every future rule/data record can name its provenance and authority level.
- No fan deck can accidentally enter the official corpus.

---

## Phase 1 — Formalize competitive core rules

**Status:** **complete — gate passed 2026-07-26.**

### Output

`docs/rules/` contains implementation-oriented rules for:

- player/character/hero/sidekick/fighter semantics;
- battlefield spaces, adjacency and zones;
- setup and starting player;
- turn/action economy;
- Maneuver, movement, BOOST, drawing and exhaustion;
- Scheme;
- Attack legality and combat resolution;
- defeat, game-end checks and hand limit;
- separately scoped multiplayer deltas.

Every normative rule has a stable rule ID and source reference.

### Gate

A minimal vanilla two-player match can be simulated from documentation alone, including exhaustion and all generic two-player game-end paths.

**Gate result:** PASS. See [`rules/phase-1-validation.md`](rules/phase-1-validation.md).

---

## Phase 2 — Timing, choices, effects and global rulings

**Status:** **complete — gate passed 2026-07-26.**

This phase establishes the global resolution framework that fighter/card content must target.

### Completed work

1. Built event/window/checkpoint semantics for turns, actions and combat.
2. Defined deterministic pause/resume for player choices and reconnect.
3. Separated public/private/committed-hidden/revealed information.
4. Defined normalized effects with trigger, condition, choices, dependencies, ordered operations and provenance.
5. Established an extensible primitive taxonomy.
6. Reconciled cancellation, move/place, gained/free actions, bonus attacks, `End the turn`, dormant players and setup hooks.
7. Indexed global rulings separately from normalized engine semantics.
8. Added a severity-based ambiguity register.

### Gate

- Every core timing window is ordered: **PASS**.
- Pausing/resuming for player input is deterministic: **PASS**.
- Global errata/rulings are reconciled with current Core Rules: **PASS**.
- Unknown interactions are explicit rather than guessed: **PASS**.

See [`mechanics/phase-2-validation.md`](mechanics/phase-2-validation.md).

---

## Phase 3 — Canonical set and release registry

**Status:** **complete — gate passed 2026-07-26.**

Phase 3 establishes the exhaustive work queue for fighter/deck and battlefield research.

### Completed work

1. Created [`sets/registry.yaml`](sets/registry.yaml) with canonical IDs, release/edition metadata, licensing/availability classification, fighters, battlefields, components, set mechanics and source provenance.
2. Registered **25 released primary product records** through Stars & Stripes, while preserving historical single-fighter releases and product lineage.
3. Registered gameplay supplements outside the primary-set model:
   - TMNT Shredder & Krang Hero Decks;
   - Nova High Battlefield Mat.
4. Registered Hellboy as `announced/blocked` without importing unconfirmed community roster data.
5. Added [`sets/mechanics-index.md`](sets/mechanics-index.md), mapping special rule families to authoritative official sources.
6. Added [`sets/source-bibliography.md`](sets/source-bibliography.md) for current catalog, historical rulebooks, retired licensed content and freshness checks.
7. Preserved edition/addendum provenance:
   - Cobble & Fog Vault reprint/clarifications;
   - Robin Hood vs. Bigfoot Vault return;
   - Bruce Lee original solo lineage and 2025 return;
   - Stars & Stripes White House Secret Passages addendum.
8. Kept Adventures enemy/scenario logic deferred while registering its competitive-compatible heroes/battlefields.

### Freshness decisions

- **Stars & Stripes:** released; separate Secret Passages addendum is first-class authority. Complete cards still require Phase 4 physical/authoritative verification if public normalized databases lag.
- **TMNT:** 2025 crowdfunding/release lineage with 2026 Restoration general availability; Turtles are competitive-compatible and Shredder/Krang Hero Decks are separate competitive content.
- **Nova High:** official battlefield existence is verified; exact graph/equivalence is deferred to Phase 5.
- **Hellboy:** official future four-character product is announced, but no final playable corpus is currently published; remains blocked.

### Gate

No known released competitive fighter or battlefield is absent from the registry, and every identified set-specific mechanic has an authoritative official entry point.

**Gate result:** PASS. See [`sets/phase-3-validation.md`](sets/phase-3-validation.md).

Passing Phase 3 does not make fighter/card or battlefield data `developer-ready`; it defines their exhaustive later-phase queues.

---

## Phase 4A — Representative fighter/deck stress-test corpus

**Status:** planned.

Before transcribing the full catalog, model a deliberately difficult sample that exercises different engine behaviors. Initial candidates:

- Achilles — sidekick-death state changes / bonus attacks;
- Bloody Mary — action-count-dependent behavior / bonus attacks;
- Sun Wukong — summonable sidekicks;
- Sherlock Holmes — cancellation restrictions and information effects;
- Dracula — start-of-turn ability and multiple sidekicks;
- Raptors — multiple heroes;
- Wayward Sisters — multiple heroes, cauldron/spell resource model;
- Geralt of Rivia — pre-game deck construction / gear selection;
- Yennefer & Triss — setup-time hero selection;
- one Marvel fighter with non-standard resources/components.

### Work per fighter

Capture:

- fighter stats and attack type;
- hero/sidekick/multiple-hero topology;
- health and movement;
- ability and setup rules;
- resources/tokens/state;
- deck-construction rules;
- every card's quantity, user, type, printed value, BOOST and normalized effects;
- applicable errata and official rulings;
- source provenance for every non-trivial interpretation.

### Gate

The representative sample can be expressed without ad-hoc undocumented engine behavior. Any mechanic that cannot be normalized becomes an explicit custom-mechanic requirement rather than an implicit exception.

---

## Phase 4B — Complete fighter and deck corpus

**Status:** planned.

### Work

Expand the Phase 4A schema to every released competitive fighter in `sets/registry.yaml`.

Use UmDb as the primary normalized deck index, but verify set-specific behavior and disputed wording against higher-authority sources. Never import `/decks/...` fan decks into the official corpus; published UmDb entries live under `/umdb/...`.

### Validation

For every fighter:

- expected card count reconciles with deck construction;
- quantities sum correctly;
- all card users are valid fighters;
- every effect maps to a known mechanic primitive or explicit custom mechanic;
- all referenced resources are defined;
- all known official fighter/card rulings are linked.

### Gate

Every released competitive fighter is `verified` or has an explicit `blocked` record describing the missing authoritative data.

---

## Phase 5 — Battlefield corpus

**Status:** planned.

Battlefields require more than images or names; the engine needs graph data.

### Work per battlefield

Capture:

- stable space/node IDs;
- zone membership;
- undirected adjacency;
- directional/one-way edges;
- starting-space assignments;
- special connections such as secret passages;
- doors, high ground, battlefield items, conveyor or other map components;
- large-fighter restrictions/edge markings where relevant;
- setup rules and official rulings.

UmDb's board catalog is an index, not proof of graph correctness. Graphs must be independently verified against authoritative battlefield material.

### Gate

Every supported battlefield has a deterministic graph representation that can answer legal movement/targeting questions without image interpretation at runtime.

---

## Phase 6 — Cross-validation and ambiguity closure

**Status:** planned.

### Work

1. Run consistency checks across core rules, mechanics, fighter cards and battlefields.
2. Maintain the living ambiguity register.
3. Resolve each ambiguity using the source hierarchy.
4. Mark genuinely unresolved cases as blockers rather than inventing behavior.
5. Verify that later errata/addenda supersede older wording without destroying provenance.
6. Sample full turns and combats involving stress-test fighters and special battlefields.

### Gate

There are no unresolved **P0/P1 semantic ambiguities** affecting legal moves, hidden information, choice ownership, combat outcome, fighter defeat, or game end for the intended launch scope.

---

## Phase 7 — Freeze the developer-facing rules contract

**Status:** planned.

### Output

Publish stable documentation schemas and engine semantics for:

- `GameState` and immutable identifiers;
- legal client/server commands;
- event/timing windows;
- pending interactions and choices;
- effect operations and custom handlers;
- fighter state/resources;
- deck construction and card instances;
- battlefield graph;
- deterministic event/audit history;
- persistence/reconnect implications.

This phase documents **what the engine must represent**, not a premature language/framework architecture.

### Gate

All requirements in [`specification-readiness.md`](specification-readiness.md) pass. Only then should implementation planning become authoritative.

---

## Phase 8 — Adventures/cooperative extension

**Status:** deferred.

After competitive play is stable, separately model villains, minions, initiative, objectives, threat/scenario state and other Adventures-specific behavior. Competitive heroes from Adventures releases belong in Phases 3–6; cooperative enemy logic does not.