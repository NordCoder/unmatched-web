# Cancellation and effect identity

Cancellation must target a defined effect scope. It is not equivalent to nullifying a card, changing its printed value, removing its combat role, or canceling every rule associated with the fighter.

## Baseline

`CANCEL-001` — When an effect is canceled, do not resolve that canceled effect.  
Source: CORE p. 13.

`CANCEL-002` — A normal instruction to cancel effects **on a card** targets cancelable effects belonging to that played card. It does not, merely by association, cancel the fighter's special ability, static battlefield rule, or other non-card effect.  
Source: CORE p. 11 distinguishes card/character effects; REF10 character clarifications repeatedly distinguish special abilities from card effects.

`CANCEL-003` — Canceling effects on a combat card does not by itself change the card's printed/combat value or stop the card from serving as the attack/defense card unless the canceling source also says to ignore/change the value.  
Source: CORE p. 13; Feint and value rules.

## Effect-bearing identity

`CANCEL-010` — The engine must identify cancellation targets as effect instances with a source, not only as card IDs. At minimum distinguish:

```text
played_card_effect
attached_card_effect
character_ability_effect
battlefield_effect
static_modifier
rule/replacement_effect
```

Source: normalized requirement from cancellation rulings.

`CANCEL-011` — An effect that has been added to a played card by another component can become part of that card's cancelable effect set when the authoritative rule says so. Combat-item effects are the canonical example.  
Source: REF10 Battlefield Effects: Battlefield Items; RULES-HUB General Rulings.

`CANCEL-012` — If the played card's effects are explicitly protected from cancellation, card-attached combat-item effects inherit that protection under the current ruling.  
Source: RULES-HUB General Rulings: combat items on cards whose effects cannot be canceled.

## Cards with no printed effects

`CANCEL-020` — A card with no printed effects is not considered to have had its effects canceled merely because the opponent played a generic `cancel effects on your opponent's card` instruction.  
Source: REF10 Major Rulings: Cancelling Card Effects; official King Arthur/Excalibur ruling referenced there.

`CANCEL-021` — Therefore downstream mechanics that ask whether the opposing card's **effects were canceled** must use the semantic cancellation result, not merely whether a cancellation instruction was played.  
Source: normalized consequence of CANCEL-020.

This distinction is required for King Arthur's BOOST/Excalibur interaction.

## Timing interaction

`CANCEL-030` — Cancellation itself resolves at the timing window printed on its source. It cannot retroactively undo an effect that has already fully resolved unless an explicit rule says so.  
Source: CORE combat timing/order; normalized causal consequence.

`CANCEL-031` — In ordinary combat, defender effects at a timing window resolve before attacker effects at that same window. Consequently, a defender's `IMMEDIATELY` effect that resolves before the attacker's cancellation effect is not retroactively canceled by that later effect.  
Source: CORE p. 11; RULES-HUB card-interaction ruling on defender `IMMEDIATELY` ordering.

## Protected effects

`CANCEL-040` — If a source explicitly says an effect cannot be canceled, cancellation instructions do not suppress it.  
Source: specific card/ability precedence under CORE p. 1.

`CANCEL-041` — Protection belongs to the protected effect/source scope stated by the content. Do not infer that the fighter or card is globally immune to all future cancellation unless the source says so.  
Source: source-policy non-invention rule.

## Engine requirement

A minimal cancellation record needs:

```yaml
cancellation:
  target_effect_ids: [...]
  scope: effects_on_opposing_card | explicit_effect | source_defined
  result: canceled | no_cancelable_effect | protected
```

This allows later effects to distinguish `canceled` from `a cancellation instruction was attempted`.
