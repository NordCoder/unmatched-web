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

### Inputs

- current Core Rules;
- official errata and general rulings;
- later set rules where they clarify a core term retroactively.

### Output

`docs/rules/` now contains implementation-oriented rules for:

- component and actor model: player, character, hero, sidekick, fighter;
- battlefield spaces, adjacency, zones, starting spaces;
- setup and starting player;
- turn and action economy;
- Maneuver, movement, BOOST and drawing;
- Scheme action;
- Attack legality, melee/ranged targeting and defense;
- combat reveal and ordered resolution;
- `IMMEDIATELY`, `DURING COMBAT`, damage and `AFTER COMBAT`;
- defeat, game-end checks, hand limit and exhaustion;
- two-player rules plus separately scoped free-for-all/team deltas.

Each normative rule has a stable rule ID and source reference.

### Gate

A minimal two-fighter match with only vanilla cards can be simulated from documentation alone, including exhaustion and every generic two-player game-end path.

**Gate result:** PASS. See [`rules/phase-1-validation.md`](rules/phase-1-validation.md).

The gate deliberately excludes fighter/set-specific mechanics that require Phase 2 timing, choice, cancellation, bonus-attack, extra-action, placement, or other global rulings. Passing Phase 1 does not make the overall project `developer-ready`.

---

## Phase 2 — Timing, choices, effects and global rulings

**Status:** **complete — gate passed 2026-07-26.**

This phase establishes the global resolution framework that fighter/card content must target.

### Completed work

1. Built an event/window/checkpoint model for turns, actions and combat.
2. Defined deterministic pause/resume semantics for player choices, including reconnect during hidden combat and setup choices.
3. Separated public, private, committed-hidden and temporarily revealed information.
4. Defined a normalized effect representation that preserves trigger, condition, choice, cost/dependency, ordered operations and provenance.
5. Established an extensible primitive taxonomy for cards, health, relocation, actions/resources/state, combat values, cancellation and control flow.
6. Reconciled cancellation semantics, including effectless cards and attached battlefield-item effects.
7. Formalized normal/gained/free action accounting.
8. Formalized move versus place, corrected occupied-space placement selection, failed placement and swap semantics.
9. Formalized corrected bonus-attack behavior as nested combat inside one Attack action.
10. Formalized current `End the turn` and dormant-player rulings.
11. Added generalized setup hooks for pre-game configuration/deck construction/placement.
12. Indexed global rulings separately from normalized engine semantics.
13. Created a severity-based ambiguity register so unresolved content-specific interactions cannot be silently guessed.

### Output

Primary documents are under [`mechanics/`](mechanics/) and [`rulings/`](rulings/):

- `mechanics/event-model.md`;
- `mechanics/choices-and-resume.md`;
- `mechanics/effect-model.md`;
- `mechanics/cancellation.md`;
- `mechanics/action-accounting.md`;
- `mechanics/movement-and-placement.md`;
- `mechanics/bonus-attacks.md`;
- `mechanics/end-turn-and-dormancy.md`;
- `mechanics/information-visibility.md`;
- `mechanics/setup-hooks.md`;
- `rulings/global-rulings.md`;
- `rulings/ambiguity-register.md`.

The primitive taxonomy remains extensible. A future fighter may require an explicit custom mechanic, but it must declare the same trigger/choice/state/cancellation/provenance semantics rather than hiding behavior in implementation code.

### Gate

- Every core timing window is ordered: **PASS**.
- A rule exists for pausing resolution when another player must choose something: **PASS**.
- Global errata/rulings have been reconciled with the current Core Rules: **PASS**.
- Unknown interactions are explicitly marked; none are silently guessed: **PASS**.

See [`mechanics/phase-2-validation.md`](mechanics/phase-2-validation.md).

Remaining P1 ambiguities in the register are explicitly `deferred-content`: they block only affected fighter/set content until Phase 3/4 supplies authoritative rules. There are no open P0 ambiguities in the current global framework.

Passing Phase 2 does not make the published fighter roster `developer-ready`; it makes the global resolution model ready to receive authoritative set/fighter/card data.

---

## Phase 3 — Canonical set and release registry

**Status:** planned.

### Work

For every officially released competitive-compatible set:

1. record canonical set ID, name, release/edition metadata and license status;
2. link official product/rulebook/set-rule sources;
3. enumerate fighters, battlefields and special components;
4. capture set-level mechanics and setup changes;
5. capture addenda and later errata as separate provenance records.

The registry must include standalone fighters and competitive heroes originating in Adventures products.

### Explicit freshness cases

- **Stars & Stripes:** track the separate Secret Passages addendum and any incomplete public deck databases.
- **Hellboy:** remain `blocked/unreleased` until authoritative release material exists.
- Revised/reprinted products must not silently overwrite older wording; preserve edition/source metadata.

### Gate

No released competitive fighter or battlefield is absent from the registry, and every set-specific mechanic has at least one authoritative source.

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

Expand the Phase 4A schema to every released competitive fighter.

Use UmDb as the primary normalized deck index, but verify set-specific behavior and disputed wording against higher-authority sources. Never import `/decks/...` fan decks into the official corpus; published UmDb entries live under `/umdb/...`.

### Validation

For every fighter:

- expected card count reconciles with deck construction;
- quantities sum correctly;
- all card users are valid fighters;
- every effect maps to a known mechanic primitive or an explicitly documented custom mechanic;
- all resources referenced by cards are defined;
- all known official fighter/card rulings are linked.

### Gate

Every released competitive fighter is `verified` or has an explicit `blocked` record describing the missing authoritative data.

---

## Phase 5 — Battlefield corpus

**Status:** planned.

Battlefields require more than images or names; the engine needs graph data.

### Work per battlefield

Capture:

- spaces as stable node IDs;
- zone membership;
- undirected adjacency;
- directional/one-way edges where applicable;
- starting-space assignments;
- special connections such as secret passages;
- doors, high ground, portals, item locations or other battlefield components;
- large-fighter restrictions/edge markings where relevant;
- setup rules and official rulings.

UmDb's board catalog is an index, not sufficient proof of graph correctness. Graphs must be manually or independently cross-validated against authoritative battlefield material.

### Gate

Every supported battlefield has a deterministic graph representation that can answer legal movement/targeting questions without image interpretation at runtime.

---

## Phase 6 — Cross-validation and ambiguity closure

**Status:** planned.

### Work

1. Run consistency checks across core rules, mechanics, fighter cards and battlefields.
2. Build a living ambiguity register.
3. Resolve each ambiguity using the source hierarchy.
4. Mark genuinely unresolved cases as implementation blockers rather than inventing behavior.
5. Verify that later errata/addenda supersede older wording without destroying provenance.
6. Sample full turns and combats involving the stress-test fighters and special battlefields.

### Gate

There are no unresolved **P0/P1 semantic ambiguities** affecting legal moves, hidden information, choice ownership, combat outcome, fighter defeat, or game end for the intended launch scope.

---

## Phase 7 — Freeze the developer-facing rules contract

**Status:** planned.

### Output

Publish stable documentation schemas and engine semantics for:

- `GameState` and immutable identifiers;
- action/command legality;
- event/timing windows;
- pending interactions and player choices;
- effect operations and custom handlers;
- fighter state/resources;
- deck construction and card instances;
- battlefield graph;
- deterministic resolution and audit/event history;
- persistence/reconnect requirements implied by hidden information and pending choices.

This phase documents **what the engine must represent**, not a premature language/framework architecture.

### Gate

All requirements in [`specification-readiness.md`](specification-readiness.md) pass. Only then should implementation planning become authoritative.

---

## Phase 8 — Adventures/cooperative extension

**Status:** deferred.

After competitive play is stable, separately model villains, minions, initiative, objectives, threat/scenario state and other Adventures-specific behavior. Competitive heroes from Adventures releases belong in Phases 3–6; cooperative enemy logic does not.
