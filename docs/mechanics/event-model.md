# Event and timing model

This document turns the rulebook's procedural timing into stable engine checkpoints. It extends, and where noted refines, `docs/rules/combat.md`, `turn-structure.md`, and `defeat-and-game-end.md`.

## Resolution vocabulary

`TIMING-001` — An **event** is a rules-significant occurrence that can create zero or more effects to resolve, e.g. `TURN_START`, `ATTACK_DECLARED`, `COMBAT_AFTER`, `FIGHTER_DEFEATED`.  
Source: normalized engine terminology derived from CORE timing labels and character/ruling triggers.

`TIMING-002` — A **window** is an ordered point at which effects with the matching timing may resolve. Combat has explicit `IMMEDIATELY`, `DURING_COMBAT`, `APPLY_COMBAT_DAMAGE`, `AFTER_COMBAT`, `CLEANUP`, and `ADDITIONAL_EFFECTS` stages.  
Source: CORE pp. 11–12.

`TIMING-003` — A **checkpoint** is a mandatory rules evaluation not represented as a card effect, e.g. the modern winner check at the start/end of an action and the dormant-player check at the end of an action.  
Source: CORE p. 14; REF10 Major Rulings: Dormant Players.

`TIMING-004` — Resolving an effect may emit new events, but the engine must preserve the source-defined order and may not arbitrarily interleave later independent effects into the middle of the current effect.  
Source: CORE p. 13 partial-resolution/order semantics; RULES-HUB global/card interaction rulings (finish the current effect before proceeding to the next unless a rule explicitly creates replacement/interruption semantics).

## Generic two-player turn skeleton

`TIMING-010` — A normal turn is represented by the following semantic sequence:

```text
TURN_START
  -> START_OF_TURN window
  -> ACTION_BOUNDARY_START checkpoint
  -> ACTION_1 declaration/resolution
  -> ACTION_BOUNDARY_END checkpoint
  -> ACTION_BOUNDARY_START checkpoint
  -> ACTION_2 declaration/resolution
  -> ACTION_BOUNDARY_END checkpoint
  -> END_OF_TURN window
  -> HAND_LIMIT checkpoint
  -> TURN_TRANSFER
```

Effects may add actions or end the turn; those modifications are defined in `action-accounting.md` and `end-turn-and-dormancy.md`.  
Source: CORE pp. 5, 14 plus normalized lifecycle.

`TIMING-011` — The modern winner predicate is evaluated at every `ACTION_BOUNDARY_START` and `ACTION_BOUNDARY_END`. It is not replaced by an immediate game-over mutation when a hero reaches zero health.  
Source: CORE p. 14; REF10 Official Errata: Winning the Game; `docs/rules/defeat-and-game-end.md`.

`TIMING-012` — The dormant predicate is evaluated at `ACTION_BOUNDARY_END` after the current action has finished.  
Source: REF10 Major Rulings: Dormant Players.

`TIMING-013` — `START_OF_TURN` effects occur before the active player declares their first action. `END_OF_TURN` effects occur after the final action has resolved and before the hand-limit cleanup unless a specific rule states otherwise.  
Source: CORE p. 5 and printed character timing; RULES-HUB examples such as She-Hulk/Philippa; normalized lifecycle.

## Action lifecycle

`TIMING-020` — A declared action is one semantic action even when its resolution contains multiple internal stages, effects, choices, or a bonus attack.  
Source: CORE action procedures; REF10 Card Effects: Bonus Attacks.

`TIMING-021` — The start-of-action winner checkpoint occurs before accepting the declaration/inputs of that action.  
Source: CORE p. 14; normalized implementation consequence.

`TIMING-022` — The end-of-action winner and dormant checkpoints occur only after every non-skipped stage that belongs to that action has finished.  
Source: CORE pp. 12, 14; REF10 Dormant Players; normalized implementation consequence.

## Attack/combat pipeline

`TIMING-030` — An ordinary Attack action resolves in this sequence:

```text
ATTACK_DECLARED
  -> attacker commits attack card
  -> defender chooses whether/how to defend
  -> cards reveal simultaneously
  -> IMMEDIATELY
  -> DURING_COMBAT
  -> APPLY_COMBAT_DAMAGE
  -> AFTER_COMBAT_CARD_EFFECTS
  -> AFTER_COMBAT_OTHER_EFFECTS
  -> CLEANUP
  -> ADDITIONAL_EFFECTS ("after attacking")
  -> ACTION_BOUNDARY_END
```

Source: CORE pp. 10–12; RULES-HUB general ruling that non-card after-combat abilities resolve after card After Combat effects.

`TIMING-031` — `AFTER COMBAT` card effects are not the same timing as effects that happen `after attacking`; the latter occur only in `ADDITIONAL_EFFECTS`, after Cleanup.  
Source: CORE p. 12.

`TIMING-032` — A played combat card remains eligible to resolve its later scheduled card effects if its fighter is defeated during combat. Conditions/targets are evaluated against the actual current state, so positional requirements may fail after the fighter is removed.  
Source: CORE p. 11.

`TIMING-033` — Within the same combat timing window, defender-controlled effects resolve before attacker-controlled effects. If one player has multiple effects at that same combat timing, that player chooses their order.  
Source: CORE p. 11.

`TIMING-034` — Static/inherent modifiers are not one-shot queued effects when their source says they continuously apply; their condition is re-evaluated as needed while relevant.  
Source: RULES-HUB General Rulings: inherent/static/floating bonuses.

## Effect ordering inside a source

`TIMING-040` — Card/effect instructions are followed in printed/source order unless the effect explicitly presents multiple options whose order the controller may choose or another rule changes the ordering.  
Source: CORE p. 13; RULES-HUB General Rulings: card effects are generally followed in printed order.

`TIMING-041` — When one effect presents multiple decisions that must be made for that effect, those decisions are locked before executing the chosen outcomes, unless the text explicitly makes a later choice depend on an earlier resolved result.  
Source: RULES-HUB General Rulings: multiple decisions are chosen before executing them; conservative exception for explicit sequential wording.

## Cross-player simultaneous non-combat effects

`TIMING-050` — Phase 2 does **not** invent a universal cross-player ordering rule for every future non-combat trigger. Combat has an explicit defender-first rule. Some specific rulings give the active player tie-breaking authority for particular simultaneous decisions. A future fighter/set that creates a materially ambiguous simultaneous non-combat interaction must cite its own authority or enter the ambiguity register before that content is `developer-ready`.  
Source: source-policy non-invention rule; CORE p. 11 is combat-scoped; current global rulings are interaction-specific.

This is deliberate: a software engine may have a generic queue implementation, but the queue's ordering policy must come from gameplay authority, not from implementation convenience.
