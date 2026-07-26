# Bonus attacks

Bonus attacks are a rules-defined composite combat procedure. They are neither ordinary extra actions nor a second attack declaration that re-checks all initial targeting requirements.

## Trigger and action identity

`BONUS-001` — A bonus attack is triggered by an `AFTER COMBAT` effect of the main combat. The main combat resolves normally until that triggering effect is reached.  
Source: REF10 Card Effects: Bonus Attacks.

`BONUS-002` — If the bonus-attack-triggering effect is canceled before it resolves, the bonus attack does not occur.  
Source: REF10 Card Effects: Bonus Attacks; cancellation semantics.

`BONUS-003` — The bonus attack is part of the same top-level `ATTACK` action as the main combat. It does not consume another ordinary action opportunity and does not create a generic action boundary between the main and bonus combat.  
Source: REF10 Card Effects: Bonus Attacks; normalized action consequence.

`BONUS-004` — Therefore generic start/end-of-action winner checks are not inserted between the main and bonus combat. The top-level Attack action reaches its end checkpoint only after the composite bonus-attack sequence and all remaining stages belonging to that action are complete.  
Source: CORE p. 14 + BONUS-003.

## Participants

`BONUS-010` — The bonus combat uses the **same attacker and defender** as the main combat.  
Source: REF10 Card Effects: Bonus Attacks.

`BONUS-011` — The defender gets a new defense decision for the bonus combat and may play a new legal defense card or decline to defend.  
Source: REF10 Card Effects: Bonus Attacks.

`BONUS-012` — Do not re-select the target merely because the defender moved after the main combat. The bonus attack remains against the same defender even if that fighter would no longer be a legal target for a newly declared ordinary attack (for example, because it is now out of range).  
Source: REF10 Bonus Attack procedure; RULES-HUB General Rulings/Card Interactions: targeting legality is checked for the initial attack, not re-run to retarget the bonus attack.

`BONUS-013` — If either the attacker or defender was defeated during the main combat, the bonus attack does not occur.  
Source: REF10 Official Errata/Card Effects: Bonus Attacks; current Rules Hub official ruling index.

This corrects older interpretations that allowed a bonus attack to continue after one combat participant was defeated.

## Bonus attack card semantics

`BONUS-020` — The bonus attack has its own displayed name, combat value, and effects as defined by the source card's bonus-attack section. For combat resolution, treat that section as the attack-card definition for the nested bonus combat.  
Source: REF10 Card Effects: Bonus Attacks.

`BONUS-021` — The bonus section uses the same BOOST value as the source card unless the published component explicitly says otherwise.  
Source: REF10 Card Effects: Bonus Attacks.

`BONUS-022` — Effects belonging only to the main attack portion do not automatically repeat in the bonus combat. Effects printed for the bonus section resolve according to their own timing.  
Source: REF10 Bonus Attack procedure; normalized component distinction.

## Nested combat pipeline

`BONUS-030` — Once triggered, the bonus combat resolves using the ordinary combat stages that apply to it:

```text
BONUS_COMBAT_START
  -> defender commits/declines defense
  -> reveal bonus attack definition + defense card
  -> IMMEDIATELY
  -> DURING_COMBAT
  -> APPLY_COMBAT_DAMAGE
  -> AFTER_COMBAT
  -> cleanup applicable to the bonus combat
  -> return to the outer attack-resolution continuation
```

Source: REF10: resolve the bonus attack like a normal attack; CORE combat pipeline.

`BONUS-031` — The implementation must preserve the fact that the bonus attack definition is derived from an already-played source card. It must not create a second physical copy of that card in a zone merely to model the nested combat.  
Source: normalized digital representation requirement.

`BONUS-032` — A defense card played specifically for the bonus combat is a distinct played card and follows normal defense-card cleanup for that combat.  
Source: ordinary combat procedure applied to the bonus combat.

## Battlefield items and per-combat attachments

`BONUS-040` — A combat item/effect used for the main combat does not automatically carry into the bonus combat merely because both combats belong to one Attack action. Each combat applies its own item eligibility/use procedure.  
Source: RULES-HUB bonus-attack/battlefield-item ruling.

`BONUS-041` — If an eligible combat item was not used for the main combat, it may be available for the bonus combat if all source-specific conditions are satisfied.  
Source: RULES-HUB bonus-attack/battlefield-item ruling.

## Re-entrancy guard

`BONUS-050` — A bonus attack does not recursively produce another bonus attack unless an authoritative source explicitly creates that possibility. The generic bonus-attack procedure must not infer recursion from ordinary wording.  
Source: source-policy non-invention rule; current published bonus-attack structure.

## Engine representation

A bonus attack should be represented as nested combat state inside the parent Attack action, for example:

```yaml
action:
  type: attack
  stage: resolving_bonus_combat
  primary_combat: ...
  bonus_combat:
    attacker: fighter_id
    defender: fighter_id
    attack_definition_source: card_instance_id
    bonus_definition_id: ...
    defense_commitment: ...
  resume_point: outer_after_combat_queue
```

The exact runtime schema may differ. The invariant is that reconnect can distinguish the top-level action from its nested combat.
