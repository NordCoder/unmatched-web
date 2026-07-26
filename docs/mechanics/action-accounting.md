# Action accounting: normal, gained and free actions

This document separates three concepts that must not be conflated: the turn's action budget, an action as a rules-significant object, and an effect that resolves without consuming the ordinary budget.

## Normal action budget

`ACTION-001` — A normal turn begins with 2 required action opportunities unless a rule/effect changes the count. Ordinary Maneuver, Scheme, and Attack actions consume one ordinary action opportunity.  
Source: CORE p. 5.

`ACTION-002` — Ordinary available actions cannot be voluntarily skipped.  
Source: CORE p. 5; RULES-HUB General Rulings confirms gained actions cannot simply be declined.

## Gaining actions

`ACTION-010` — `gain 1 action` increases the active player's remaining ordinary action budget by one unless the source restricts the gained action.  
Source: published abilities/cards; RULES-HUB General Rulings: gained actions cannot be skipped.

`ACTION-011` — A gained unrestricted action is still a normal rules action for all action-boundary semantics: start/end winner checks, action triggers, and dormant checks after resolution.  
Source: CORE p. 14 plus RULES-HUB/REF10 clarification that free/gained actions remain actions; normalized consequence.

`ACTION-012` — If a source grants an action with restrictions (for example, only a particular action type), those restrictions belong to that granted-action token/permission rather than globally modifying the player's normal action choices.  
Source: normalized extensibility requirement; exact restricted-action content is set/fighter-specific and must retain provenance.

## Free actions

`ACTION-020` — A source that explicitly lets something be performed `as a free action` creates an action that does not consume one of the ordinary action opportunities granted by the turn.  
Source: set-specific use of `free action`; contrast with ordinary actions in CORE and battlefield-item rules.

`ACTION-021` — A free action is nevertheless still an **action** for rules that count or check actions.  
Source: REF10 Bloody Mary Rules Clarification: free actions are still considered actions.

`ACTION-022` — Therefore the modern start/end-of-action winner checks apply to a free action as to any other action.  
Source: CORE p. 14 + ACTION-021.

`ACTION-023` — A free action must retain any source-defined type/eligibility restriction. `Free` describes budget cost; it does not imply permission to ignore normal legality or fighter/banner requirements.  
Source: source-specific precedence; normalized requirement.

## Bonus attacks are not extra actions

`ACTION-030` — A bonus attack is an internal continuation of the original Attack action, triggered from the first combat's After Combat effect. It does not consume an ordinary action and does not create a separate top-level action boundary.  
Source: REF10 Card Effects: Bonus Attacks; see `bonus-attacks.md`.

`ACTION-031` — Do not run the generic start/end-of-action winner checkpoint between the main and bonus combat; the Attack action has not yet ended.  
Source: ACTION-030 + CORE p. 14 action-boundary rule.

## State model

A future engine should represent action permissions explicitly rather than only an integer:

```yaml
action_budget:
  ordinary_remaining: 2
  granted:
    - id: ...
      cost: 1 | 0
      allowed_types: [maneuver, scheme, attack] | source_defined
      source: ...
```

This prevents fighter-specific `free attack`, `play a scheme as a free action`, or future restricted grants from being incorrectly treated as generic `actions_remaining += 1`.

## Deferred set-specific refinement

Stars & Stripes (2026) introduces new action-specific grants whose exact use-it/lose-it semantics must be taken from that set's authoritative rules/card corpus during Phase 3/4. Phase 2 deliberately provides the representation slot without importing community explanations as canonical gameplay authority.
