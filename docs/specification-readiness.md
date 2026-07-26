# Specification readiness gate

## Purpose

Implementation should begin only when the rules corpus is detailed enough that engine behavior is derived from documentation rather than invented during coding.

This document defines `developer-ready` for the first target: **competitive two-player Unmatched**, with the data model capable of later free-for-all/team and Adventures extensions.

## 1. Core rules gate

All must pass:

- terminology has stable definitions (`player`, `character`, `hero`, `sidekick`, `fighter`, `space`, `zone`, etc.);
- setup is deterministic given fighter/set/battlefield and player ordering;
- turn/action economy is explicit;
- Maneuver, BOOST, movement and draw semantics are explicit;
- attack legality is explicit for melee and ranged fighters;
- defense is optional/required exactly where the rules specify;
- complete combat timing is ordered;
- effect precedence/cancellation rules are documented;
- fighter defeat and every game-end checkpoint are explicit;
- exhaustion and hand-limit behavior are explicit;
- all normative rules retain provenance.

## 2. Choice/resolution gate

The documentation must identify when resolution pauses for player input.

Examples include:

- choosing an attack/defense card;
- choosing whether to defend;
- choosing a fighter, card, opponent or space;
- choosing ordering when a rule grants ordering authority;
- pre-game deck/setup selections;
- optional effects (`may`).

For each choice type, specify:

```yaml
owner: player identifier
visibility: public | private
legal_options: deterministic rule
optional: true | false
resume_point: timing/event identifier
```

This is required for online play, reconnect and server-authoritative resolution.

## 3. Fighter gate

Every supported fighter must define:

- canonical ID, display name and set/edition;
- character structure (single hero, hero + sidekick(s), multiple heroes, etc.);
- attack type(s);
- starting health for every independently tracked fighter;
- movement value;
- special ability;
- setup/pre-game choices;
- tokens/resources/state machines;
- defeat/summon/off-board behavior if non-standard;
- deck-construction rule;
- linked official rulings/errata;
- provenance.

## 4. Card gate

Every card needed by a supported fighter must define:

```yaml
id: stable_id
name: display_name
quantity: integer
usable_by: [fighter_ids]
type: attack | defense | versatile | scheme | special
printed_value: integer | null
boost: integer | null
effects:
  - timing: ...
    condition: ...
    choices: ...
    operations: ...
sources: [...]
```

The exact schema may change during research; the semantic requirements may not be discarded.

Validation requirements:

- quantities reconcile with deck construction;
- no effect refers to an undefined fighter/resource/state;
- each effect maps to a documented primitive or explicit custom mechanic;
- effect order is preserved;
- any later errata is linked and normalized without deleting original provenance.

## 5. Battlefield gate

Every supported battlefield must define a graph, not only an image:

```yaml
spaces:
  - id: ...
    zones: [...]
    starting_slot: ...
edges:
  - from: ...
    to: ...
    kind: normal | one_way | door | special
special_connections: [...]
special_components: [...]
sources: [...]
```

The graph must support deterministic answers for:

- adjacency;
- legal path traversal;
- zone membership;
- ranged targeting;
- starting placement;
- special passages/connections;
- special battlefield rules;
- large-fighter restrictions when relevant.

## 6. Rulings and ambiguity gate

Maintain an ambiguity/ruling record with severity:

- `P0` — can change winner/game end or make state irrecoverably inconsistent;
- `P1` — changes legal action, hidden information, target, damage, card ownership/order, fighter defeat or combat outcome;
- `P2` — narrower edge case with deterministic workaround;
- `P3` — wording/documentation quality only.

Before implementation planning:

- unresolved P0 = 0;
- unresolved P1 = 0 for the initially supported fighter/battlefield roster;
- P2/P3 items may remain only when they are explicitly documented and cannot affect the intended launch subset.

## 7. Corpus completeness gate

For the intended launch scope:

- every set is registered;
- every fighter is registered;
- every required action card is present;
- every battlefield is present;
- every set addendum/erratum is linked;
- known official rulings are indexed;
- data gaps are explicit rather than silently approximated.

A globally unreleased or externally incomplete fighter may remain `blocked` without blocking an MVP, provided it is outside the declared supported roster.

## 8. Engine-contract gate

The specification must be sufficient to define, without choosing a programming language/framework:

- `GameState` facts that must persist;
- legal client/server commands;
- deterministic legal-action generation;
- event/timing sequence;
- pending interactions and resume behavior;
- hidden/public information boundaries;
- effect operation model;
- custom mechanic extension point;
- card/fighter/battlefield definitions;
- deterministic event/audit history required for reconnect/debugging.

## Final readiness statement

Implementation planning may become authoritative only when a review can truthfully state:

> For the declared initial roster and battlefields, two independent implementers can derive the same legal actions and state transitions from the specification, without inventing gameplay rules.
