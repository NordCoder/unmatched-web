# `End the turn` and dormant players

These are global control-flow rulings. Both can invalidate a naïve engine that assumes every started effect/action reaches every ordinary stage.

## `End the turn`

`ENDTURN-001` — When an `End the turn` instruction resolves, stop processing remaining effects for the current action at that point and proceed to the action's Cleanup stage.  
Source: REF10 Major Rulings: End the Turn; RULES-HUB Rule Changes & Errata.

`ENDTURN-002` — `End the turn` is therefore a control-flow operation, not merely `actions_remaining = 0`. It suppresses unresolved effect processing that the ruling says to stop.  
Source: ENDTURN-001; normalized implementation consequence.

`ENDTURN-003` — Cleanup still occurs. Content-specific cleanup transformations still apply when they are defined as Cleanup behavior. Shakespeare's line processing and Wayward Sisters' cauldron transfer are canonical published clarifications.  
Source: REF10 Shakespeare Rules Clarification; Wayward Sisters Rules Clarification.

`ENDTURN-004` — Attack `ADDITIONAL EFFECTS` that would occur only after Cleanup are skipped when `End the turn` has already redirected resolution through Cleanup to the End-of-Turn step.  
Source: REF10 Wayward Sisters clarification.

`ENDTURN-005` — Unused ordinary, gained, or free action opportunities are not taken after `End the turn`; control advances through the required cleanup/end-of-turn process instead.  
Source: REF10 Shakespeare clarification; Major Ruling: End the Turn.

`ENDTURN-006` — After Cleanup, the current action reaches its end boundary. Run the normal end-of-action winner and dormant checkpoints before completing turn transfer/end-of-turn processing.  
Source: CORE p. 14 winner check + REF10 End-the-Turn routing; normalized checkpoint consequence.

`ENDTURN-007` — The normal End-of-Turn stage still exists after an `End the turn` instruction. Any end-of-turn effect must be evaluated according to its own source and the current state rather than silently skipped merely because the turn was ended by an effect.  
Source: REF10 Major Ruling wording: proceed to the End of Turn Step.

`ENDTURN-008` — Hand-limit cleanup remains an end-of-turn responsibility unless a more specific source overrides it.  
Source: CORE p. 5 + ENDTURN-007.

### Required control-flow state

```yaml
turn_control:
  end_turn_requested: true
  abort_remaining_current_effects: true
  next_stage: cleanup
```

The field names are illustrative; the semantics are normative.

## Dormant players

`DORMANT-001` — At the end of an action, if a player has **no fighters on the battlefield**, that player is dormant. This can be true even if one or more of that player's fighters are undefeated but currently off-board.  
Source: REF10 Major Rulings: Dormant Players.

`DORMANT-002` — Dormancy is evaluated at the end-of-action checkpoint, not immediately at every intermediate removal/off-board transition. A fighter can leave and return during the same action without forcing an intermediate dormant turn state.  
Source: REF10 Dormant Players end-of-action check; normalized consequence.

`DORMANT-003` — A dormant player does not take normal actions on their turn.  
Source: REF10 Major Rulings: Dormant Players.

`DORMANT-004` — A dormant player does not draw cards and does not discard cards merely as part of ordinary turn/action procedures while dormant.  
Source: REF10 Major Rulings: Dormant Players.

`DORMANT-005` — A dormant player cannot be selected as `an opponent` for effects that require choosing an opponent, and is excluded from effects that operate on `all opponents`, under the global dormant ruling.  
Source: REF10 Major Rulings: Dormant Players.

`DORMANT-006` — An undefeated fighter that is off the battlefield and belongs to a dormant player cannot take damage while off-board under the global dormant ruling.  
Source: REF10 Major Rulings: Dormant Players.

`DORMANT-007` — Dormancy is a derived lifecycle status, not permanent elimination. Re-evaluate it at the defined checkpoint; if the player again has a fighter on the battlefield, the player is no longer dormant.  
Source: normalized consequence of DORMANT-001 and return-to-board mechanics.

## Dormancy versus defeat

`DORMANT-010` — `dormant` and `defeated/eliminated` are different states. A dormant player's hero may still be undefeated, and later effects may return that player's fighter(s) to the board.  
Source: REF10 Dormant Players and fighter return mechanics.

`DORMANT-011` — Do not use `hero defeated` as the dormant predicate. Use actual board presence at the end-of-action checkpoint.  
Source: DORMANT-001.

## Character-specific lifecycle hooks

The global ruling does not by itself settle every fighter-specific `start of turn`, `end of turn`, off-board, transformation, or resurrection interaction for a dormant player. Those interactions must be linked to fighter-specific authority during Phase 4 and are tracked in `docs/rulings/ambiguity-register.md` when not yet verified.

This limitation does not change the global dormant predicate or restrictions above.
