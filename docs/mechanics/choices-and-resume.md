# Player choices and resume semantics

Online play requires every player decision to be represented explicitly. A rules engine must be able to persist a game while resolution is waiting for input, reconnect a client, reproduce the exact legal options, and resume from the same semantic point.

## Choice object

`CHOICE-001` — Any rule/effect that requires player input creates an explicit pending choice owned by exactly one player (or a source-defined set/order of players).  
Source: normalized implementation requirement from CORE action/effect procedures.

Minimum semantic fields:

```yaml
choice_id: stable_runtime_id
owner: player_id
source: rule_or_effect_instance
visibility: public | private | mixed
prompt_kind: choose_option | choose_card | choose_fighter | choose_space | order_effects | yes_no | other
legal_options: deterministic snapshot_or_generator
min_selections: integer
max_selections: integer
ordered: true | false
optional: true | false
resume_point: semantic_checkpoint
```

`CHOICE-002` — While a required pending choice is unresolved, the engine may accept only commands that are legal for that pending interaction (plus non-gameplay transport operations such as reconnect). It must not continue gameplay resolution speculatively.  
Source: normalized server-authoritative consequence.

`CHOICE-003` — On reconnect, the engine must reconstruct the same choice owner, public/private prompt information, legal option set under the applicable locking rule, and resume point.  
Source: normalized persistence requirement.

## Optional effects

`CHOICE-010` — `may` creates a choice at the effect's resolution point. The controlling player may decline the optional effect.  
Source: CORE p. 13.

`CHOICE-011` — `up to N` permits selecting zero where zero is otherwise legal.  
Source: CORE p. 13.

`CHOICE-012` — A valid choice may be selected even when its eventual effect is useless or resolves to no state change, unless the text requires a successful result.  
Source: RULES-HUB General Rulings: a player may choose a valid option that has no effect.

This means legality and usefulness are separate concepts.

## Choice locking

`CHOICE-020` — When one effect presents multiple decisions/options that are to be executed as part of that effect, the decisions are locked before executing the outcomes, unless explicit sequential wording makes a later decision depend on an earlier resolved state.  
Source: RULES-HUB General Rulings: multiple decisions are made before executing them.

Example normalized shape:

```yaml
choices:
  - choose: effect_A
  - choose: effect_B
constraints:
  different: true
execution_order: controller_choice
```

The engine must not execute `effect_A`, expose its result, and only then allow selection of `effect_B` when the gameplay rule requires both decisions first.

`CHOICE-021` — If one player controls multiple effects at the same combat timing, ordering those effects is itself a player choice and must be resolved before the selected order executes.  
Source: CORE p. 11.

## Hidden combat choices

`CHOICE-030` — Attack-card selection is private until reveal. The attacker commits a legal card before the defender chooses whether/how to defend.  
Source: CORE p. 10.

`CHOICE-031` — Defense is optional. If defending, the defender commits a legal defense card privately; otherwise the defender explicitly passes the defense choice.  
Source: CORE p. 10.

`CHOICE-032` — After both combat choices are complete, the committed cards are revealed simultaneously. Neither player may revise a committed card because of information learned only at reveal.  
Source: CORE p. 10; normalized commitment consequence.

The server therefore needs a committed-but-not-yet-public card state distinct from hand and revealed combat state.

## Choice legality versus operation success

`CHOICE-040` — A destination/target can be a legal **choice** while the subsequent operation fails, when an authoritative ruling explicitly distinguishes selection legality from operation success. Placement is the canonical example; see `movement-and-placement.md`.  
Source: REF10 Major Rulings: Placement; RULES-HUB Rule Changes & Errata.

`CHOICE-041` — If a mandatory operation has no valid selectable object/space at all, skip the impossible part and continue resolving the rest of the effect according to partial-resolution rules.  
Source: CORE p. 13; RULES-HUB General Rulings: no valid move/place space.

## Information boundary

`CHOICE-050` — The choice payload sent to a client must contain only information that player is entitled to know at that point. A legal-option generator may internally use authoritative hidden state without exposing it.  
Source: normalized online-play requirement; see `information-visibility.md`.

## Determinism requirement

For the same authoritative game state and same pending source, independent implementations must derive the same legal choice domain. UI presentation may differ; gameplay legality may not.
