# Attack and combat resolution

This file defines the generic competitive Attack action and the core combat pipeline.

## Attack legality and target declaration

`ATK-001` — An Attack action begins by declaring one of the active player's fighters as attacker and one opposing fighter as target.  
Source: CORE p. 10.

`ATK-002` — Attack is illegal unless the active player has both (a) a valid target for some attacker and (b) an attack-capable card in hand that the chosen attacker is authorized to use.  
Source: CORE p. 10.

`ATK-003` — Every fighter may make a melee attack against an opposing fighter in an adjacent space, regardless of zones. A melee-only fighter has no generic attack range beyond adjacency.  
Source: CORE p. 10.

`ATK-004` — A ranged fighter may additionally attack an opposing fighter whose space shares a zone with the attacker's space, regardless of adjacency.  
Source: CORE p. 10.

`ATK-005` — Generic attack legality therefore holds when the chosen target is adjacent to the attacker, or when the attacker is ranged and the two spaces share at least one zone, subject to card/banner and special-rule restrictions.  
Source: normalized predicate from ATK-003/004.

Whether an adjacent attack by a ranged fighter needs a distinct `melee`/`ranged` classification for later effects is not inferred here; Phase 2/character rules must source any semantic distinction beyond target legality.

## Combat-card selection

`ATK-010` — After target declaration, the attacker chooses one legal attack or versatile card from hand, authorized for the attacking fighter, and commits it face down.  
Source: CORE p. 10.

`ATK-011` — The defender then chooses either (a) one legal defense or versatile card from hand authorized for the defending fighter, committed face down, or (b) no defense card. Defense is optional under the generic rule.  
Source: CORE p. 10.

`ATK-012` — A versatile card is legal as either attack or defense and is considered both an attack card and a defense card for effects that inspect its card type.  
Source: CORE p. 10.

`ATK-013` — After both players have completed their choices, the committed combat cards are revealed simultaneously. Hidden selection/resume protocol is a Phase 2 engine concern; the semantic requirement is that the defender does not know the attack card identity when making the normal defense choice.  
Source: CORE p. 10; implementation consequence of face-down commitment and simultaneous reveal.

## Core combat timing

`COMBAT-001` — Once cards are revealed, generic combat resolves in this ordered pipeline:

1. `IMMEDIATELY` effects;
2. `DURING COMBAT` effects;
3. apply combat damage and determine the combat result;
4. `AFTER COMBAT` card effects;
5. cleanup played action cards;
6. effects defined as happening `after attacking` or equivalent attack-level additional effects.

Source: CORE pp. 11–12.

`COMBAT-002` — Effects may originate from played combat cards or from other rules/components such as character abilities. Their timing label/window determines when they participate.  
Source: CORE p. 11.

`COMBAT-003` — When attacker and defender effects are due in the same timing window, the defender resolves their due effect(s) before the attacker resolves theirs.  
Source: CORE p. 11.

`COMBAT-004` — When the same player controls multiple effects due at the same time, that player chooses their resolution order.  
Source: CORE p. 11.

`COMBAT-005` — Defeat of a fighter during combat does not by itself erase effects on that fighter's already-played combat card. Those effects still attempt to resolve in their scheduled windows, although conditions or targets may no longer exist because the defeated fighter has been removed from the battlefield.  
Source: CORE p. 11; defeat removal: CORE p. 14.

The exact generalized event queue, nested triggers, and ordering between non-card effects are expanded in Phase 2.

## Combat damage

`COMBAT-010` — Let `A` be the final effective value of the played attack card when combat damage is applied. Let `D` be the final effective value of the played defense card, or `0` if the defender played no defense card. Generic combat damage dealt to the defender is `max(A - D, 0)`.  
Source: CORE p. 12; non-negative normalization follows the explicit rule that excess defense never damages the attacker.

`COMBAT-011` — The defender never deals generic combat damage back to the attacker merely because `D > A`.  
Source: CORE p. 12.

`COMBAT-012` — Apply combat damage to the defending fighter's health before resolving `AFTER COMBAT` effects. If that fighter reaches zero health, normal defeat/removal happens before those later effects.  
Source: CORE pp. 12, 14.

`COMBAT-013` — Damage produced by card/ability effects is not combat damage unless a specific rule explicitly defines it as such.  
Source: CORE p. 13.

## Winning the combat

`COMBAT-020` — The attacker wins the combat exactly when the attack deals at least 1 combat damage to the defender.  
Source: CORE p. 13.

`COMBAT-021` — The defender wins the combat exactly when the defender takes 0 combat damage from the attack itself, even if an effect deals other damage during the same combat.  
Source: CORE p. 13.

`COMBAT-022` — Under the generic rule, the combat therefore has exactly one winner: attacker for positive combat damage, defender otherwise.  
Source: normalized consequence of COMBAT-020/021.

## Cleanup and action completion

`COMBAT-030` — After `AFTER COMBAT` card effects resolve, all played action cards still associated with the combat are placed in their respective discard piles during cleanup, unless a specific effect replaces that destination.  
Source: CORE p. 12.

`COMBAT-031` — An effect that happens `after attacking` is not an `AFTER COMBAT` card effect. It occurs after combat cleanup as the final generic step of the Attack action.  
Source: CORE p. 12.

`COMBAT-032` — Game-end checking follows `defeat-and-game-end.md`; the modern generic rules do not stop the Attack action at the instant a hero first reaches zero health.  
Source: CORE p. 14; REF10 official erratum “Winning the Game”.