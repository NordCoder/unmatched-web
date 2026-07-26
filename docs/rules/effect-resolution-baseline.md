# Core effect-resolution baseline

This file records effect semantics explicitly present in the modern Core Rules. Phase 2 has now expanded this baseline into the complete global event/choice/effect framework under `docs/mechanics/` and `docs/rulings/`.

## General precedence

`EFFECT-001` — When a specific effect contradicts a general rule, apply the specific effect for its stated scope.  
Source: CORE p. 1.

`EFFECT-002` — An effect is mandatory unless its wording makes it optional, such as by using `may`.  
Source: CORE p. 13.

`EFFECT-003` — If an optional effect is being resolved, the controlling player chooses whether to perform it at that resolution point.  
Source: CORE p. 13.

`EFFECT-004` — An instruction allowing an amount `up to N` permits choosing zero.  
Source: CORE p. 13.

Phase 2 refinement: see `docs/mechanics/choices-and-resume.md` for explicit choice ownership, locking, visibility and resume semantics.

## Partial resolution

`EFFECT-010` — If only part of an effect can be resolved, resolve every resolvable part and ignore the impossible remainder unless the effect itself establishes a dependency or different instruction.  
Source: CORE p. 13.

`EFFECT-011` — Failure of one independent clause does not automatically cancel other independently resolvable clauses of the same effect.  
Source: normalized consequence of EFFECT-010.

Phase 2 refinement: `docs/mechanics/effect-model.md` distinguishes conditions, explicit costs/prerequisites, choices, ordered operations and dependency edges such as `if you do`.

## Values

`EFFECT-020` — `printed value` means the value printed on the card and ignores temporary or permanent changes from other effects.  
Source: CORE p. 13.

`EFFECT-021` — If a card's value is instructed to be ignored, treat that value as 0 when calculating combat damage, even if other effects would otherwise modify it.  
Source: CORE p. 13.

`EFFECT-022` — The engine must therefore distinguish at least `printed_value` from the effective/current value used at combat-damage calculation.  
Source: normalized requirement from EFFECT-020/021.

## Cancellation

`EFFECT-030` — If an effect is canceled, that canceled effect is not resolved.  
Source: CORE p. 13.

`EFFECT-031` — Card effects, character abilities, battlefield effects and other effect-bearing sources must not be treated as one undifferentiated cancellation scope.  
Source: REF10 cancellation rulings; Phase 2 normalization.

Phase 2 refinement: see `docs/mechanics/cancellation.md`, including effectless-card cancellation and attached combat-item effects.

## Simultaneous combat effects

`EFFECT-040` — Within a combat timing window, if attacker and defender effects would resolve at the same time, defender-controlled effects have priority.  
Source: CORE p. 11.

`EFFECT-041` — If one player controls multiple effects that would resolve at the same combat timing, that player chooses their order.  
Source: CORE p. 11.

These are combat-specific ordering rules. Phase 2 deliberately does not invent a universal cross-player non-combat ordering rule; see `docs/mechanics/event-model.md` and `docs/rulings/ambiguity-register.md`.

## Effects from defeated fighters

`EFFECT-050` — A played combat card remains scheduled for its later combat effects even if the fighter who played it is defeated during that combat.  
Source: CORE p. 11.

`EFFECT-051` — Removal of the defeated fighter can make a later effect fail a positional condition or lose a valid target; in that case apply partial-resolution rules rather than pretending the fighter remains on the battlefield.  
Source: CORE pp. 11, 14.

## Movement and placement

`EFFECT-060` — Effect-granted movement follows generic movement rules unless the effect overrides them. When moving an opposing fighter, interpret friendly/opposing occupancy from that fighter owner's perspective.  
Source: CORE p. 7.

`EFFECT-061` — Placement does not trace a movement path. A successful placement relocates the fighter directly to the selected destination under the applicable occupancy policy.  
Source: CORE p. 7; Phase 2 placement ruling normalization.

`EFFECT-062` — Selection legality and placement success are distinct. When wording does not require `empty` or `other`, an occupied space may be selectable and the subsequent placement may fail, leaving the fighter where it was.  
Source: REF10 Major Rulings: Placement; current Rules Hub official correction.

Phase 2 refinement: see `docs/mechanics/movement-and-placement.md`.

## Phase 2 closure

The former Phase 2 handoff topics are now covered by the following authoritative specification layers:

- costs/dependencies and normalized operations → `docs/mechanics/effect-model.md`;
- event windows, nested effects and conditions → `docs/mechanics/event-model.md`;
- hidden/public choices and resume points → `docs/mechanics/choices-and-resume.md`, `information-visibility.md`;
- cancellation identity/scope → `docs/mechanics/cancellation.md`;
- `End the turn` and dormant players → `docs/mechanics/end-turn-and-dormancy.md`;
- bonus attacks → `docs/mechanics/bonus-attacks.md`;
- gained/free actions → `docs/mechanics/action-accounting.md`;
- setup hooks → `docs/mechanics/setup-hooks.md`;
- unresolved content-specific interactions → `docs/rulings/ambiguity-register.md`.

Phase 2 does not claim that every future fighter-specific interaction is already resolved. It defines the global framework and requires unresolved content-specific cases to remain explicit blockers rather than implementation guesses.
