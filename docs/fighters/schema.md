# Fighter and deck manifest schema

**Schema maturity:** Phase 4A verified. Phase 4B may add corpus-proven fields, but existing semantic distinctions must not be collapsed.

## 1. Fighter manifest

```yaml
schema_version: 1
id: stable-fighter-id
display_name: Display Name
set_ids: [canonical-release-id]
status: verified | partial | blocked

movement: 2

topology:
  type: single_hero | hero_with_sidekicks | multiple_heroes | selectable_hero
  loss_rule: primary_hero_defeated | all_heroes_defeated
  fighters:
    - id: stable-piece-id
      display_name: ...
      role: hero | sidekick | selectable_hero
      count: 1
      attack_type: melee | ranged
      starting_health: 10 | source_defined_by_role | null
      starts_on_board: true
      summonable: false

setup:
  hooks: []
  state_initialization: []

ability:
  id: ...
  controller: owner
  effects: []

resources:
  - id: ...
    kind: token_pool | card_zone | counter | external_definition_pool
    visibility: public | private
    initial: ...

persistent_state:
  - key: ...
    type: boolean | integer | enum | fighter_ref | fighter_ref_set | space_ref | card_ref | card_ref_set
    reset: never | turn_start | turn_end | action_start | action_end | combat_start | combat_end

deck:
  manifest: ../cards/<id>.yaml
  construction: fixed | choose_groups | source_defined

rulings: []
sources: []
```

### Required invariants

`FIGHTER-SCHEMA-001` — every independently damageable fighter has its own stable piece/instance identity and health semantics. A visual token group is not sufficient.

`FIGHTER-SCHEMA-002` — `hero`, `sidekick`, `multiple_heroes`, and setup-selected roles are gameplay semantics, not UI labels. The loss rule is explicit rather than inferred from the number of health dials.

`FIGHTER-SCHEMA-003` — summonable pieces remain canonical fighter definitions even when they start off-board. Reserve availability is represented separately from battlefield presence.

`FIGHTER-SCHEMA-004` — setup choices use the Phase 2 pending-choice model. A client may not silently choose gear, hero identity, sidekick, or other configuration locally.

`FIGHTER-SCHEMA-005` — historical state needed by effects (for example action ordinal, the turn-start space, whether an attack happened this turn, or source-defined parent context) is persisted explicitly; it must not be reconstructed from current board state.

## 2. Deck manifest

```yaml
schema_version: 1
fighter_id: stable-fighter-id
status: verified | partial | blocked

construction:
  kind: fixed | choose_groups
  available_pool_count: 30
  game_deck_count: 30
  rules: []

card_zones:
  - id: deck
    owner: self
  - id: hand
    owner: self
  - id: discard
    owner: self

cards:
  - id: stable-card-definition-id
    name: display name
    quantity: 2
    inclusion: base | choose_group | fixed
    usable_by: [fighter-piece-ids] | any
    type: attack | defense | versatile | scheme
    printed_value: 3 | null
    boost: 2 | null
    tags: []
    effects:
      - id: ...
        trigger: ...
        controller: ...
        optional: false
        condition: null
        choices: []
        costs: []
        operations: []
        cancellation: cancelable | protected | source_defined

external_definitions: []
validation:
  available_pool_count: ...
  game_deck_count: ...
  quantity_sum: ...

sources: []
```

### Card identity and instances

`CARD-SCHEMA-001` — a card **definition** is distinct from a physical/game **card instance**. Quantity creates multiple instances sharing a definition.

`CARD-SCHEMA-002` — a card instance must retain at least:

```yaml
instance_id: ...
definition_id: ...
owner_player: ...
current_zone: ...
zone_controller: ...
visibility: ...
```

Ownership is immutable unless an authoritative rule explicitly changes ownership. Location/control/use authority is mutable.

This distinction is required by Black Panther: an opponent-owned card can be face-up in Black Panther's Vibranium Suit, usable by Black Panther only for BOOST, and later return to the **owner's** discard pile.

`CARD-SCHEMA-003` — a custom card zone is a first-class zone, not an arbitrary list attached to a fighter. It declares visibility, who may move/use cards from it, replacement/disposition rules, and source lifetime.

`CARD-SCHEMA-004` — deck size is data. Do not hard-code `30`. Geralt has a 36-card available pool but constructs a 30-card game deck; published fixed decks also vary elsewhere in the product line.

`CARD-SCHEMA-005` — `external_definitions` are gameplay definitions referenced by cards/abilities but not ordinary action-card instances in the game deck. Phase 4A examples include bonus-attack definitions and Wayward Sisters spells.

`CARD-SCHEMA-006` — the engine must distinguish immutable definition facts from runtime value layers:

```text
printed_value_base
effective_printed_value
current_combat_value
boost_base
effective_boost_value
```

A source can modify one layer without rewriting the others.

## 3. Normalized effect records

Phase 4 card effects target `docs/mechanics/effect-model.md`.

Current generic operations include:

- card/zone: `DRAW`, `DISCARD`, `REVEAL`, `LOOK_AT`, `MOVE_CARD`, `SHUFFLE`, `REORDER`;
- health/damage: `DEAL_EFFECT_DAMAGE`, `APPLY_COMBAT_DAMAGE`, `PREVENT_DAMAGE`, `REDIRECT_DAMAGE`, `RECOVER`, `DEFEAT`, `RETURN_FIGHTER`;
- position: `MOVE`, `PLACE`, `SWAP`, `SUMMON`;
- actions/resources/state: `GAIN_ACTION`, `GAIN_RESOURCE`, `SPEND_RESOURCE`, `SET_STATE`, `CHANGE_STATE`, `PREVENT_OPERATION`;
- combat/value: `ADD_VALUE`, `SET_VALUE`, `SET_PRINTED_VALUE`, `IGNORE_VALUE`, `ADD_BOOST_VALUE`, `BOOST`, `CANCEL_EFFECTS`;
- control/composites: `REQUEST_CHOICE`, `TRIGGER_COMPOSITE`, `REPLACE_COMBAT_CARD`, `END_TURN`.

A manifest may introduce a new operation only when corpus evidence proves the existing vocabulary cannot preserve the rule. The new semantic must then be promoted/documented generically, not left as an unexplained character-named branch.

## 4. Source representation vs normalized representation

Do not copy long card text into the corpus. For each card retain:

1. factual printed metadata (name, quantity, user, type, value, BOOST, tags);
2. normalized semantic operations;
3. source URL(s) and, where relevant, ruling IDs.

A future UI may present separately licensed/localized display text. Engine correctness must not depend on scraping prose at runtime.

## 5. Validation contract

A fighter/deck pair passes Phase 4 validation only when:

- all card quantities reconcile with construction rules;
- every `usable_by` target exists in fighter topology;
- every referenced resource/zone/state exists;
- every effect maps to a documented generic primitive/composite or an explicitly justified custom mechanic;
- source-sensitive interpretations link provenance;
- known official rulings are attached or explicitly deferred;
- no fan-patch data is mixed with published data.