# Core effect-resolution baseline

This file records only effect semantics explicitly present in the modern Core Rules. Phase 2 will build the complete event/choice/effect model and reconcile all global rulings.

## General precedence

`EFFECT-001` — When a specific effect contradicts a general rule, apply the specific effect for its stated scope.  
Source: CORE p. 1.

`EFFECT-002` — An effect is mandatory unless its wording makes it optional, such as by using `may`.  
Source: CORE p. 13.

`EFFECT-003` — If an optional effect is being resolved, the controlling player chooses whether to perform it at that resolution point.  
Source: CORE p. 13.

`EFFECT-004` — An instruction allowing an amount `up to N` permits choosing zero.  
Source: CORE p. 13.

## Partial resolution

`EFFECT-010` — If only part of an effect can be resolved, resolve every resolvable part and ignore the impossible remainder unless the effect itself establishes a dependency or different instruction.  
Source: CORE p. 13.

`EFFECT-011` — Failure of one independent clause does not automatically cancel other independently resolvable clauses of the same effect.  
Source: normalized consequence of EFFECT-010.

The boundary between “independent clause,” prerequisite, cost, and conditional consequence is Phase 2 work and must be derived from rulings/card wording rather than guessed.

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

`EFFECT-031` — Core rules do not by themselves define every cancellation identity/scope needed by all fighters. For example, rulings distinguish card effects from character abilities and address cards with no printed effects. Those details are Phase 2/4 requirements.  
Source: REF10 major ruling “Cancelling Card Effects” and character clarifications.

## Simultaneous combat effects

`EFFECT-040` — Within a combat timing window, if attacker and defender effects would resolve at the same time, defender-controlled effects have priority.  
Source: CORE p. 11.

`EFFECT-041` — If one player controls multiple effects that would resolve at the same time, that player chooses their order.  
Source: CORE p. 11.

These are combat-specific ordering rules. Do not automatically extend them to every simultaneous non-combat event until Phase 2 establishes a general model.

## Effects from defeated fighters

`EFFECT-050` — A played combat card remains scheduled for its later combat effects even if the fighter who played it is defeated during that combat.  
Source: CORE p. 11.

`EFFECT-051` — Removal of the defeated fighter can make a later effect fail a positional condition or lose a valid target; in that case apply partial-resolution rules rather than pretending the fighter remains on the battlefield.  
Source: CORE pp. 11, 14.

## Movement and placement

`EFFECT-060` — Effect-granted movement follows generic movement rules unless the effect overrides them. When moving an opposing fighter, interpret friendly/opposing occupancy from that fighter owner's perspective.  
Source: CORE p. 7.

`EFFECT-061` — Placement does not trace a movement path. A successful placement results in the fighter occupying an empty destination.  
Source: CORE p. 7.

`EFFECT-062` — The later official placement ruling permits some wording to select an occupied space and then have placement fail. This choice-vs-result distinction is intentionally not generalized until Phase 2.  
Source: REF10 major ruling “Placement”; current Rules Hub freshness index.

## Phase 2 handoff

The following are explicitly **not closed by this document**:

- costs versus effects;
- nested/triggered effects;
- event snapshots and re-checking conditions;
- effect ownership/control when cards change hands;
- `after attacking`, `start of turn`, `end of turn`, defeat triggers, and other non-card timing queues;
- hidden/public choices and when options are locked;
- `End the turn`;
- dormant players;
- bonus attacks;
- generalized extra/free actions;
- replacement effects and destination replacement;
- cancellation of whole cards versus individual effects.

Each must become explicit before the overall specification reaches `developer-ready`.