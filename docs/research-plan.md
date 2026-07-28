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

### Gate

No known released competitive fighter or battlefield is absent from the registry, and every identified set-specific mechanic has an authoritative official entry point.

**Gate result:** PASS. See [`sets/phase-3-validation.md`](sets/phase-3-validation.md).

---

## Phase 4A — Representative fighter/deck stress-test corpus

**Status:** **complete — gate passed 2026-07-26.**

Phase 4A validates the fighter/deck schema and the Phase 2 effect framework against deliberately difficult published characters before transcribing the full roster.

### Representative corpus

- Achilles — sidekick-defeat persistent state, combat participant replacement, bonus attack;
- Bloody Mary — action ordinal/history, turn-start snapshot, bonus attack parent context;
- Sun Wukong — summonable sidekicks, reserve pool, damage prevention/redirection;
- Sherlock Holmes — cancellation protection, face-up pre-reveal prediction, effective printed values;
- Dracula — multiple sidekicks, damage-dependent ability, returned fighter, combat-card replacement;
- Raptors — multiple heroes with independent health and all-heroes loss condition;
- Wayward Sisters — multiple heroes, Cauldron card zone, ingredient tags and external spells;
- Geralt of Rivia — 36-card available pool, pre-game 30-card construction, gear and ongoing schemes;
- Yennefer & Triss — setup-selected hero/sidekick roles, role-dependent ability, simultaneous hidden choices;
- Black Panther — opponent-owned cards in a controlled storage zone and sequential BOOST/draw semantics.

### Completed work

1. Created [`fighters/schema.md`](fighters/schema.md) and promoted it to the verified Phase 4A manifest contract.
2. Created complete fighter manifests for all ten representatives under `fighters/phase-4a/`.
3. Created complete normalized action-card manifests for all ten representatives under `cards/phase-4a/`.
4. Reconciled every fixed deck to its published quantity and modeled Geralt's 36→30 construction explicitly.
5. Modeled bonus attacks and Wayward spells as `external_definitions`, not fake action-card instances.
6. Proved first-class card zones and separation of immutable ownership from current zone/use authority through Cauldron and Vibranium Suit.
7. Extended `mechanics/effect-model.md` only where the published corpus proved a generic semantic was missing:
   - `PREVENT_DAMAGE`;
   - `REDIRECT_DAMAGE`;
   - `PREVENT_OPERATION`;
   - `SET_PRINTED_VALUE`;
   - `ADD_BOOST_VALUE`;
   - `REORDER`;
   - captured parent context;
   - `REPLACE_COMBAT_CARD` composite.
8. Added [`fighters/phase-4a-mechanics.md`](fighters/phase-4a-mechanics.md) and [`fighters/phase-4a-validation.md`](fighters/phase-4a-validation.md).
9. Excluded community balance-patch `/decks/...` data where it reuses original character/card names but changes published gameplay values.

### Gate

The representative sample can be expressed without ad-hoc undocumented engine behavior. Any new mechanic discovered by the corpus is promoted to an explicit reusable semantic or composite.

**Gate result:** PASS. See [`fighters/phase-4a-validation.md`](fighters/phase-4a-validation.md).

Passing 4A validated the schema/framework; the complete-roster gate was closed separately by Phase 4B.

---

## Phase 4B — Complete fighter and deck corpus

**Status:** **complete — final Phase 4 gate passed 2026-07-28.**

### Completed work

1. Integrated the ten Phase 4A representatives and all four reconciled Phase 4B ownership groups into one canonical corpus.
2. Established exactly **74 unique competitive fighter identities**, including one canonical Bruce Lee lineage and no evidence-alias duplicates.
3. Integrated **74 fighter manifests** and **74 paired card manifests**.
4. Completed fresh independent QA for all five physical-card correction scopes before integration.
5. Reconciled fixed and constructed decks, including Daredevil 22, Elektra 20, Black Widow 31, Geralt 36→30 and Buffy 35→30.
6. Preserved action cards, auxiliary decks, external definitions, positioned components and non-card components as distinct data categories.
7. Preserved four-axis verification status rather than inflating every fighter to `verified`.
8. Added a reproducible machine validator and persisted its result in [`fighters/phase-4-final-validation.json`](fighters/phase-4-final-validation.json).

### Final validation

```text
canonical fighter identities: 74
fighter manifests: 74
card manifests: 74
unique action-card definitions: 926
available action-card copies: 2214
owner requirements defined/referenced: 52/52
unresolved requirement references: 0
validation errors: 0
validation warnings: 0
status totals: 28 verified / 45 partial / 1 blocked
```

Deadpool remains the single blocked fighter because the digital-adaptation policy is intentionally unresolved. Partial fighters have deterministic published behavior but still require shared runtime capabilities or retain explicit source qualifications.

### Gate

- all 74 competitive fighters present: **PASS**;
- all deck constructions reconciled: **PASS**;
- all card-image P1/P2 corrections closed: **PASS**;
- all owner requirements resolvable by ID: **PASS**;
- evidence-only aliases separated from canonical IDs: **PASS**;
- physical-evidence qualifications preserved: **PASS**.

**Gate result:** PASS. See [`fighters/phase-4-final-validation.md`](fighters/phase-4-final-validation.md) and [`qa/phase-4-card-image/correction-integration-qa.md`](qa/phase-4-card-image/correction-integration-qa.md).

Passing Phase 4 completes the canonical published fighter/card corpus. It does **not** make the whole project developer-ready because owner requirements still need consolidation into one generic runtime contract.

---

## Phase 4C — Shared runtime requirements consolidation

**Status:** next major phase.

### Goal

Consolidate the 52 corpus-proven `A/B/C/D-REQ-*` definitions into canonical generic runtime capabilities. Aliases may remain for provenance, but one gameplay mechanic must not acquire multiple engine implementations.

### Parallel workstreams

1. Resolution, choices, history, delayed obligations, reconnect and RNG persistence.
2. Fighter presence, occupancy, footprints and positioned battlefield objects.
3. Movement, targeting, attack legality and combat participant replacement.
4. Damage, health assignment and continuous/dynamic modifiers.
5. Card zones, auxiliary decks, resources and action permissions.

### Output

- canonical requirement registry with owner aliases;
- developer-facing state/command/event/choice implications;
- persistence and visibility requirements;
- deterministic test scenarios;
- migrated manifest references or a temporary complete alias map.

### Gate

- owner requirements deduplicated;
- launch-scope runtime behavior deterministic;
- no undocumented character-specific semantics;
- no separate implementation of the same generic capability;
- launch-roster requirements implementation-ready.

---

## Phase 5 — Battlefield corpus

**Status:** planned; one MVP battlefield graph may proceed in parallel with Phase 4C.

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

**Status:** planned; stable engine-foundation work may begin earlier where the relevant contract is already deterministic.

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

All requirements in [`specification-readiness.md`](specification-readiness.md) pass. Only then should full implementation planning become authoritative.

---

## Phase 8 — Adventures/cooperative extension

**Status:** deferred.

After competitive play is stable, separately model villains, minions, initiative, objectives, threat/scenario state and other Adventures-specific behavior. Competitive heroes from Adventures releases belong in Phases 3–6; cooperative enemy logic does not.
