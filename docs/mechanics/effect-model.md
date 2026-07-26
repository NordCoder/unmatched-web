# Normalized effect model

The purpose of this model is to express published gameplay semantics without turning every card into custom code. It is a documentation contract; concrete serialization may evolve during fighter transcription.

## Effect categories

`FX-001` — A **triggered effect** resolves when its stated event/timing occurs, subject to its condition. Examples include `IMMEDIATELY`, `AFTER COMBAT`, `after attacking`, `start of turn`, and `end of turn`.  
Source: CORE pp. 11–13 and published special-ability timing.

`FX-002` — A **static/ongoing modifier** remains applicable while its source and condition remain active. It is re-evaluated as relevant state changes rather than consumed as a one-shot queued effect.  
Source: RULES-HUB General Rulings on inherent/static/floating bonuses; REF10 ongoing/special rules.

`FX-003` — A **replacement effect/rule** changes what an attempted operation does when the source uses replacement semantics such as `instead`. Exhaustion is the baseline example: a missing required draw is replaced by damage rather than becoming a successful draw plus damage.  
Source: CORE p. 6; normalized replacement terminology.

`FX-004` — A **composite mechanic** is a rules-defined procedure built from multiple ordinary stages but with special invariants. Bonus attacks and `End the turn` are composite mechanics and are documented separately.  
Source: REF10 Major/Card Effects rulings.

## Minimum effect representation

`FX-010` — A non-trivial effect record must preserve at least:

```yaml
id: stable_definition_id
source: card | ability | battlefield | rule
controller: rule_for_determining_player
trigger: event_or_window
condition: predicate | null
optional: true | false
choices: []
costs: []
operations: []
continuous: false
replacement_for: null
cancellation: cancelable | protected | source_defined
provenance: []
```

Source: normalized implementation requirement.

`FX-011` — The model must preserve the semantic distinction between **condition**, **cost/prerequisite**, **choice**, and **operation**. Flattening all text into an unordered list is invalid because it changes partial-resolution/dependency behavior.  
Source: CORE p. 13; published `if you do`/requirement patterns; normalized requirement.

## Conditions and state

`FX-020` — Unless a rule explicitly captures earlier state, an effect condition is evaluated when that effect/operation is reached for resolution against the then-current authoritative state.  
Source: CORE combat example and defeated-fighter condition behavior; normalized rule.

`FX-021` — When wording depends on historical state (`started this turn`, `attacked this turn`, etc.), the engine must preserve that history/snapshot as state rather than reconstruct it from current position.  
Source: CORE Momentous Shift example; REF10 character rules.

`FX-022` — Defeated fighters are removed immediately, so later conditions involving adjacency/zone/board presence observe the fighter as off-board unless a specific effect says otherwise.  
Source: CORE pp. 11, 14.

## Ordered operations and dependencies

`FX-030` — Operations within one effect resolve in their published/source-defined order unless a choice explicitly grants ordering authority.  
Source: CORE p. 13; RULES-HUB General Rulings on printed order.

`FX-031` — If an independent operation cannot resolve, skip that impossible operation and continue with independently resolvable operations.  
Source: CORE p. 13 partial-resolution rule.

`FX-032` — An explicit dependency such as `if you do` prevents the dependent consequence when the required prior action/cost was not performed. Do not apply generic partial resolution across an explicit dependency edge.  
Source: ordinary published effect grammar (e.g. optional payment/ability patterns); normalized dependency requirement.

`FX-033` — A **cost** is not merely any discard/damage/resource change appearing before another effect. Treat something as a cost only when the authoritative wording/ruling makes the consequence conditional on paying/performing it.  
Source: source-policy non-invention rule.

`FX-034` — Costs and prerequisites are validated before committing the dependent branch. Once a source-defined cost is paid, the engine records it as a state transition even if a later independent operation becomes impossible.  
Source: normalized transactional consequence; per-card wording/rulings may override.

## Normalized operation taxonomy

This is the Phase 2 baseline. It is intentionally extensible.

### Cards/zones

- `DRAW(count)` — draw cards; missing required draws invoke exhaustion replacement.
- `DISCARD(from, selector/count)` — move cards to discard without invoking exhaustion.
- `REVEAL(source, selector/count)` — temporarily make card identity public, then follow source-defined disposition.
- `LOOK_AT(source, selector/count, viewer)` — grant temporary private knowledge without automatically making it public.
- `MOVE_CARD(from_zone, to_zone, selector)` — explicit card-zone transition.
- `SHUFFLE(zone)` — randomize a source-defined zone.

### Fighters/health

- `DEAL_EFFECT_DAMAGE(target, amount)` — non-combat damage.
- `APPLY_COMBAT_DAMAGE(attacker, defender, amount)` — combat damage from value comparison.
- `RECOVER(target, amount)` — health recovery for an undefeated valid fighter.
- `DEFEAT(target)` — direct defeat when a source explicitly does so.
- `RETURN_FIGHTER(target, destination, health_rule)` — return/revive semantics; not equivalent to recover.

### Battlefield position

- `MOVE(target, distance/path_rule)` — movement rules apply unless overridden.
- `PLACE(target, destination_rule)` — no path; destination/placement semantics in `movement-and-placement.md`.
- `SWAP(a, b, rule)` — composite relocation requiring fighter/map-specific occupancy handling.
- `SUMMON(pool, destination_rule)` — composite sidekick-token placement per set/global summoning rules.

### Actions/resources/state

- `GAIN_ACTION(count)` — modifies the active player's required action budget; see `action-accounting.md`.
- `GAIN_RESOURCE(resource, amount)` / `SPEND_RESOURCE(resource, amount)`.
- `SET_STATE(key, value)` / `CHANGE_STATE(key, transition)`.

### Combat values/effects

- `ADD_VALUE(target, amount)` — additive modifier.
- `SET_VALUE(target, value)` — sets/fixes current value; later applicable modifiers depend on timing/source.
- `IGNORE_VALUE(target)` — treat value as zero for combat-damage calculation without changing printed value.
- `BOOST(target, card_or_blind_source)`.
- `CANCEL_EFFECTS(target_scope)` — see `cancellation.md`.

### Control flow

- `REQUEST_CHOICE(definition)` — suspend until valid choice response.
- `TRIGGER_COMPOSITE(kind)` — e.g. bonus attack.
- `END_TURN` — special control-flow operation defined in `end-turn-and-dormancy.md`.

`FX-040` — These primitive names are **our normalized semantics**, not claims that official rulebooks use this vocabulary. Each fighter/card manifest must retain source provenance and may introduce an explicit custom primitive when the corpus proves this taxonomy insufficient.

## Draw versus non-draw deck operations

`FX-050` — Exhaustion damage applies when a player needs to **draw** a card and cannot. Discarding/revealing/blind-boosting from an empty deck is not automatically a draw and does not trigger exhaustion unless the source says it is a draw.  
Source: CORE p. 6; REF10 Blind BOOST; RULES-HUB General Rulings.

## Reveal semantics

`FX-060` — `reveal` makes the revealed card visible to all players for the required duration. Unless the effect changes its destination, return it to the place it came from.  
Source: REF10 Card Effects: Revealing Cards.

## Custom mechanics

`FX-070` — A fighter-specific behavior may remain a custom composite/handler only after documenting why it cannot be represented faithfully by the current primitives. Custom behavior must still declare trigger, controller, choices, visibility, state transitions, cancellation behavior, and provenance.

The goal is not a perfect DSL. The goal is to make exceptions explicit rather than hidden in ad-hoc implementation branches.
