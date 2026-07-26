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

**Status:** planned.

### Inputs

- current Core Rules;
- official errata and general rulings;
- later set rules where they clarify a core term retroactively.

### Output

Create `docs/rules/` with implementation-oriented rules for:

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
- two-player first, then free-for-all/team deltas.

Each normative rule must have a stable rule ID and source reference.

### Gate

A minimal two-fighter match with only vanilla cards can be simulated from documentation alone, including exhaustion and every game-end path.

---

## Phase 2 — Timing, choices, effects and global rulings

**Status:** planned.

This is the highest-risk phase for the eventual engine.

### Work

1. Build an event/timing model.
2. Specify simultaneous-effect ordering and effect precedence.
3. Distinguish operations that look similar but are mechanically different, e.g. `move` vs `place`, combat damage vs non-combat damage, printed value vs combat value.
4. Define pending-player choices and interruption/resume semantics.
5. Build an initial effect taxonomy from the Unmatched Reference and Rulings Archive.
6. Record official errata and global rulings separately from inferred implementation behavior.

### Target effect primitives

The taxonomy should test whether common effects can normalize into primitives such as:

`draw`, `discard`, `damage`, `heal`, `move`, `place`, `swap`, `gain_action`, `modify_value`, `cancel_effect`, `reveal`, `look_at_hand`, `choose_card`, `choose_fighter`, `choose_space`, `summon`, `defeat`, `transform_state`.

This list is a hypothesis, not a frozen design.

### Gate

- Every core timing window is ordered.
- A rule exists for pausing resolution when another player must choose something.
- Global errata/rulings have been reconciled with the current Core Rules.
- Unknown interactions are explicitly marked; none are silently guessed.

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

There are no unresolved **P0/P1 semantic ambiguities** affecting legal moves, hidden information, choice ownership, combat outcome, fighter defeat, or game end.

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
