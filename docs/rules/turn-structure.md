# Turn structure and action economy

## Active turn

`TURN-001` — On a normal turn, the active player must resolve exactly 2 actions unless an effect changes the available action count.  
Source: CORE p. 5.

`TURN-002` — A normal action cannot be voluntarily skipped.  
Source: CORE p. 5.

`TURN-003` — The active player may choose any legal sequence of two core actions. The same action type may be chosen twice.  
Source: CORE p. 5.

`TURN-004` — Core action types are `MANEUVER`, `SCHEME`, and `ATTACK`.  
Source: CORE p. 5.

## Action legality

`TURN-010` — The engine must determine action legality before accepting an action declaration. An action whose required inputs cannot be supplied is not legal.  
Source: normalized consequence of the action procedures.

`TURN-011` — `ATTACK` is illegal when the player has no legal attack card for any legal attacking fighter/target combination. In particular, the Core Rules explicitly prohibit taking Attack without an attack card in hand or a valid target.  
Source: CORE p. 10.

`TURN-012` — `SCHEME` requires a scheme card that can be legally resolved by an undefeated authorized fighter.  
Source: CORE p. 8.

`TURN-013` — Generic `MANEUVER` remains available even when the deck is empty; the mandatory draw attempt is then handled by exhaustion.  
Source: CORE p. 6.

Extra/free actions and effects that end a turn are Phase 2 topics. They must not be generalized from fighter text before their global rulings are documented.

## End of turn

`TURN-020` — A player may temporarily hold more than 7 cards during their turn.  
Source: CORE p. 5.

`TURN-021` — After the player's actions and all other end-of-turn-relevant effects are resolved, if the active player's hand contains more than 7 cards, that player chooses and discards enough cards to reach exactly 7.  
Source: CORE p. 5.

`TURN-022` — The generic hand-limit check applies at the end of that player's own turn, not immediately when the hand grows above 7.  
Source: CORE p. 5.

`TURN-023` — After end-of-turn resolution and the hand-limit step, turn control passes to the opponent unless the game has ended.  
Source: CORE p. 5; winner checks: `defeat-and-game-end.md`.

## Minimal state required

A rules engine must at least be able to represent:

```text
active_player
actions_remaining
turn_number or equivalent ordering state
hand sizes and actual private hand contents
pending end-of-turn effects/choices
```

The exact event model and resume points are Phase 2.