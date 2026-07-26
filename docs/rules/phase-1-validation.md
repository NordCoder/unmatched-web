# Phase 1 validation

**Phase:** 1 — Formalize competitive core rules  
**Validation date:** 2026-07-26  
**Result:** PASS for the declared vanilla two-player core subset  
**Overall project readiness:** NOT YET `developer-ready`; Phases 2–7 remain.

## Gate under test

The Phase 1 gate from `docs/research-plan.md` requires that a minimal two-fighter match using only vanilla core behavior can be simulated from documentation alone, including exhaustion and every generic two-player game-end path.

“Vanilla” here means:

- one hero per player;
- optional ordinary sidekicks are allowed if they use only core health/movement rules;
- no fighter-specific setup replacement;
- no bonus attacks, ongoing schemes, summoned fighters, multiple heroes, special resources, battlefield mechanics, or custom timing windows;
- cards may have no effects, or only effects whose semantics are already fully described by the Phase 1 baseline.

This scope prevents special-character rulings from being mistaken for missing core rules.

## Coverage checklist

| Required semantic area | Result | Canonical Phase 1 document |
| --- | --- | --- |
| player / character / hero / sidekick / fighter terminology | PASS | `terminology.md` |
| spaces, adjacency, zones and occupancy | PASS | `battlefield.md` |
| standard two-player setup | PASS | `setup.md` |
| turn and two-action economy | PASS | `turn-structure.md` |
| Maneuver draw / movement / BOOST | PASS | `maneuver.md` |
| exhaustion | PASS | `maneuver.md` |
| Scheme action | PASS | `scheme.md` |
| melee/ranged attack legality | PASS | `combat.md` |
| hidden combat-card commitment and optional defense | PASS | `combat.md` |
| combat timing and defender-first same-window ordering | PASS | `combat.md`, `effect-resolution-baseline.md` |
| combat damage and combat winner | PASS | `combat.md` |
| fighter defeat/removal/recovery | PASS | `defeat-and-game-end.md` |
| modern start/end-of-action winner checks | PASS | `defeat-and-game-end.md` |
| hand limit | PASS | `turn-structure.md` |
| generic effect baseline / partial resolution | PASS | `effect-resolution-baseline.md` |
| multiplayer rules kept separate from two-player invariants | PASS | `multiplayer-deltas.md` |
| normative provenance | PASS | all Phase 1 files |

## Deterministic vanilla simulation fixture

The fixture below is synthetic. It exists only to test whether the documented state transitions are sufficient; it does not reproduce a published fighter or deck.

### Initial fixture

```yaml
players:
  - id: P1
    hero:
      id: H1
      health: 10
      attack_type: melee
    move: 2
  - id: P2
    hero:
      id: H2
      health: 10
      attack_type: melee
    move: 2

starting_player: P1
starting_hands: 5_each
special_abilities: none
special_battlefield_rules: none
```

The battlefield graph contains legal numbered starting spaces and ordinary adjacency/zones.

### Setup trace

1. Initialize both fighters to starting health (`SETUP-003`).
2. Shuffle each legal fixed deck and draw 5 (`SETUP-004`).
3. Place H1 at starting slot 1 and H2 at starting slot 2 (`SETUP-020`).
4. Resolve any ordinary sidekick placement if present (`SETUP-030–032`).
5. Begin P1's turn (`SETUP-040`).

No unspecified gameplay decision remains once player order and player choices are supplied.

### Maneuver trace

P1 declares Maneuver as action 1.

1. Draw exactly one card (`MAN-010`).
2. Optionally discard one card for the normal movement BOOST (`MAN-030–033`).
3. Move each chosen fighter at most once, sequentially, through a legal adjacency path (`MAN-020–025`).
4. Return to the turn action-selection state with one normal action remaining (`TURN-001–004`).

The legality of every destination is computable from the graph and current occupancy.

### Ordinary combat trace

Assume H1 is now adjacent to H2, P1 has a legal attack card with effective value 4, and P2 chooses a legal defense card with effective value 2.

1. P1 declares H1 attacking H2 (`ATK-001–005`).
2. P1 privately commits the attack card (`ATK-010`).
3. P2 privately commits the defense card or could choose no defense (`ATK-011`).
4. Reveal both cards (`ATK-013`).
5. Resolve `IMMEDIATELY` effects. Fixture has none (`COMBAT-001`).
6. Resolve `DURING COMBAT` effects. Fixture has none (`COMBAT-001`).
7. Combat damage is `max(4 - 2, 0) = 2`; H2 health becomes 8 (`COMBAT-010–012`).
8. Attacker is the combat winner because positive combat damage was dealt (`COMBAT-020`).
9. Resolve `AFTER COMBAT` effects. Fixture has none.
10. Cleanup played cards (`COMBAT-030`).
11. Resolve attack-level additional effects. Fixture has none (`COMBAT-031`).
12. End the action and run the winner check (`GAMEEND-001`, `GAMEEND-020`). Neither hero is defeated, so play continues.

No rule beyond Phase 1 is needed for this trace.

## Exhaustion trace

Assume P1 begins a later Maneuver with an empty deck and controls H1 plus one ordinary sidekick S1.

1. Maneuver still begins legally (`MAN-013`).
2. The mandatory draw cannot be fulfilled (`MAN-010–011`).
3. One missing draw causes 2 damage to **each** P1 fighter: H1 and S1 (`MAN-011`).
4. Any fighter reduced to zero is immediately defeated and removed (`DEFEAT-001–002`).
5. The Maneuver continues through any still-applicable movement step; hero defeat alone does not instantly terminate the action (`GAMEEND-004`).
6. At the end of the Maneuver, run the winner check (`GAMEEND-001`).

This is sufficient to simulate exhaustion without reshuffling or inventing an early game-over shortcut.

## Two-player game-end path validation

### Path A — exactly one hero is defeated

At a start/end-of-action check, if exactly one hero is defeated, the player controlling the undefeated hero wins (`GAMEEND-002`, `GAMEEND-006`). This result is independent of whose turn is currently resolving.

### Path B — both heroes defeated during one action

Both defeated fighters are removed when their health reaches zero (`DEFEAT-002`). The current action finishes. At the required boundary, both heroes are defeated, so the player whose turn is resolving wins (`GAMEEND-003–004`, `GAMEEND-006`).

### Path C — active player's hero defeated but opponent survives

The action continues to its boundary. At the winner check exactly one hero is defeated, so the opposing player wins even though it is not their turn (`GAMEEND-002`, `GAMEEND-006`). The implementation must derive this from the boundary winner predicate rather than from an immediate-loss trigger at the instant damage is dealt.

### Path D — opponent's hero defeated but active player survives

The action continues to its boundary. At the winner check exactly one hero is defeated, so the active player wins (`GAMEEND-002`, `GAMEEND-006`).

### Path E — defeat already exists at the start of an action

Because the modern rule checks at both start and end of any action, an implementation must run the same winner predicate before beginning the next action (`GAMEEND-001`). This prevents a defeated-hero state from incorrectly entering another action when a previous special sequence reaches a legal boundary.

## Hand-limit validation

A hand may exceed seven during the turn (`TURN-020`). After the turn's actions and applicable end-of-turn resolution, discard chosen cards until exactly seven remain (`TURN-021–022`). There is no immediate discard merely because a draw temporarily creates an eighth card.

## Provenance validation

Phase 1 intentionally normalizes current rules rather than preserving old generic wording as executable semantics.

The most important case is the winner rule:

- modern Core Rules: winner is checked at start/end of actions;
- Unmatched Reference v10 records the corresponding official erratum;
- older “game ends immediately” generic text is retained only as superseded provenance.

This distinction is explicitly documented in `defeat-and-game-end.md`.

## Phase 2 carryover register

These questions are **not required for the vanilla Phase 1 fixture** and therefore do not fail this gate, but they are blockers for broader fighter support:

| Topic | Why Phase 2 is required |
| --- | --- |
| pending player choices / resume points | online server must pause and resume deterministic resolution |
| exact hidden-information boundaries | combat, hand inspection, secret choices |
| generalized simultaneous non-combat effects | core defender-first rule is combat-specific |
| placement selection vs successful placement | later official ruling refines core wording |
| `End the turn` | later global ruling changes what still resolves |
| dormant players | relevant after elimination in multiplayer/team contexts |
| bonus attacks | special combat/action continuation semantics |
| extra/free actions | must define action count and game-end boundaries precisely |
| cancellation identity/scope | needed by fighter-specific cards and abilities |
| effect costs / prerequisites / dependencies | partial resolution alone is not a complete effect language |
| character-specific setup hooks | pre-game choices have an official generalized order |

## Gate decision

**PASS — Phase 1 is complete for its declared scope.**

There is no unresolved semantic question required to simulate a vanilla competitive two-player match from setup through normal actions, combat, exhaustion, hand-limit cleanup, fighter defeat, and the current generic game-end rule.

This does **not** mean the repository is ready for implementation of the real character roster. Phase 2 must close the timing/choice/effect model before fighter/deck transcription can safely become executable data.