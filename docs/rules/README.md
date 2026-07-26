# Competitive core rules

**Phase:** 1 — Formalize competitive core rules  
**Status:** verified; Phase 1 gate passed for vanilla two-player play  
**Last verified:** 2026-07-26

This directory normalizes the competitive Unmatched core rules into an implementation-oriented specification. It deliberately does not yet define the complete effect DSL, pending-interaction protocol, fighter-specific exceptions, or battlefield-specific mechanics; those belong to later phases.

## Normative source baseline

- `CORE` — [Unmatched Core Rules](https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Core-rules_EN_Light.pdf), the current standalone Core Rules distributed with the modern Core/Set Rules publishing model.
- `REF10` — [Unmatched Reference v10.0 (October 2025)](https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf), used to discover publisher errata/rulings and identify superseded generic wording. It is not a replacement for official rulebooks.
- `RULE-CHANGES` — [The Unmatched Club — Rule changes & errata](https://www.the-unmatched.club/rules/rule-changes), a freshness/index source for recent official changes. It is secondary under `docs/sources/source-policy.md` and must lead back to publisher/designer rulings before being promoted into the canonical corpus.

The normalized game-end rule in this directory follows the modern Core Rules and the corresponding official erratum indexed by `REF10`; older set books that say a hero defeat ends the game immediately are superseded for generic competitive play.

## Rule ID convention

Normative statements use stable IDs. IDs should not be recycled if wording changes.

| Prefix | Area |
| --- | --- |
| `TERM-*` | terminology and component semantics |
| `FIELD-*` | spaces, adjacency and zones |
| `SETUP-*` | two-player setup |
| `TURN-*` | turn/action economy and hand limit |
| `MAN-*` | Maneuver, movement, drawing and BOOST |
| `SCHEME-*` | Scheme action and discard visibility |
| `ATK-*` | attack legality and card selection |
| `COMBAT-*` | combat timing and combat result |
| `EFFECT-*` | core effect-resolution semantics |
| `DEFEAT-*` | health, recovery and defeat |
| `GAMEEND-*` | winner checks |
| `FFA-*` / `TEAM-*` | multiplayer deltas |

A source annotation such as `Source: CORE p. 12` is provenance, not an independent rule.

## Documents

1. [`terminology.md`](terminology.md)
2. [`battlefield.md`](battlefield.md)
3. [`setup.md`](setup.md)
4. [`turn-structure.md`](turn-structure.md)
5. [`maneuver.md`](maneuver.md)
6. [`scheme.md`](scheme.md)
7. [`combat.md`](combat.md)
8. [`effect-resolution-baseline.md`](effect-resolution-baseline.md)
9. [`defeat-and-game-end.md`](defeat-and-game-end.md)
10. [`multiplayer-deltas.md`](multiplayer-deltas.md)
11. [`phase-1-validation.md`](phase-1-validation.md)

## Precedence

`TERM-001` — A fighter/set/card/battlefield effect that explicitly contradicts a general rule overrides that general rule for the scope of that effect.  
Source: CORE p. 1.

The engine therefore cannot treat this directory as a closed list of unconditional invariants. Later content manifests may install explicit exceptions, but every exception must retain provenance.

## Phase boundary

Phase 1 defines enough semantics to simulate a vanilla two-player match. Phase 2 must still formalize:

- a complete event/timing model beyond the core windows;
- interruption, player-choice and resume semantics;
- global rulings such as dormant players and `End the turn`;
- exact cancellation scopes and effect identities;
- placement target-selection edge cases;
- generalized extra/free actions;
- hidden-information boundaries and server-authoritative choice handling.

No unresolved Phase 2 topic may be silently filled in by implementation code.