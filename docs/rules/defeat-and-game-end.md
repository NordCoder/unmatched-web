# Health, defeat and game end

This document follows the **modern Core Rules** game-end model. Older set rulebooks frequently state that defeating a hero ends the game immediately; that generic wording is superseded by the later official rule/erratum described below.

## Damage, recovery and defeat

`DEFEAT-001` — Damage reduces a fighter's health. When the fighter reaches zero health, that fighter is defeated.  
Source: CORE p. 14.

`DEFEAT-002` — A defeated fighter is removed from the battlefield immediately when defeat occurs.  
Source: CORE p. 14.

`DEFEAT-003` — A fighter may recover health only while undefeated under the generic rule. A defeated fighter cannot recover health unless a specific rule explicitly returns/revives it.  
Source: CORE p. 14; specific-effect precedence: TERM-001.

`DEFEAT-004` — Recovery cannot increase a fighter above its starting health unless a specific effect explicitly changes that limit.  
Source: CORE p. 14.

`DEFEAT-005` — Defeat/removal and game-end checking are separate operations. A hero may be defeated and removed while the current action still has unresolved effects.  
Source: CORE pp. 11, 14; modern winner check below.

## Modern two-player winner check

`GAMEEND-001` — At the **start or end of any action**, evaluate the generic two-player win condition.  
Source: CORE p. 14; REF10 official erratum “Winning the Game”.

`GAMEEND-002` — At such a check, if the opponent's hero is defeated and the active player's hero is not defeated, the active player wins.  
Source: CORE p. 14.

`GAMEEND-003` — If both heroes are defeated at the winner check, the player whose turn is currently resolving wins.  
Source: CORE p. 14; REF10 official erratum “Winning the Game”.

`GAMEEND-004` — Therefore a hero reaching zero health during an action does **not** by itself terminate the action. Finish the action according to its rules, including later combat effects that remain valid, then perform the required end-of-action winner check.  
Source: CORE pp. 11–14; normalization of GAMEEND-001.

`GAMEEND-005` — Because defeated fighters are removed immediately, unresolved effects must observe the post-removal battlefield state. The delayed winner check does not leave a defeated hero on the board.  
Source: CORE pp. 11, 14.

## Superseded generic wording

`GAMEEND-010` — Statements in older rulebooks such as “the game ends immediately when the opponent's hero is defeated” are not the canonical generic rule for this project. Preserve those rulebooks as historical/set provenance, but normalize generic competitive winner checks to GAMEEND-001–005 unless later official authority supersedes them again.  
Source: source-policy precedence; CORE p. 14; REF10 official erratum.

The current community Rules Hub also tracks this as an official core-rule change. Because that site is a secondary index, future changes discovered there must be traced back to a publisher/designer ruling before replacing the canonical rules.

## Action boundaries

For Phase 1, an action boundary is the completion of one Maneuver, Scheme, or Attack action including all rules-defined resolution attached to that action.

`GAMEEND-020` — For an Attack, the end-of-action winner check occurs only after the attack pipeline is complete, including combat cleanup and attack-level additional effects.  
Source: CORE p. 12 combined with GAMEEND-001.

The exact interaction with special wording such as `End the turn`, dormant players, bonus attacks, and nested/free actions belongs to Phase 2 and can refine what constitutes the relevant action boundary for those mechanics.

## Engine requirement

The engine must not implement hero defeat as an unconditional `game_over` mutation. Instead it must represent defeat immediately and run a separate winner-check operation at the rule-defined boundary.

Conceptually:

```text
apply_damage
  -> maybe defeat/remove fighter
  -> continue current action
  -> action boundary
  -> check_winner
```

This distinction is required for combats where both heroes can become defeated during the same action.