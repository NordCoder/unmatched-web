# Authoritative State Model

## Status

```text
status: draft-foundation
parent_issue: #19
correction_issue: #32
architecture: architecture-contract.md
fixture_contract: deterministic-fixture-contract.md
```

This document defines the minimum deterministic gameplay state required to validate commands, resolve gameplay, persist/reconnect and derive player projections. It also defines the adjacent authority and presence records that are deliberately excluded from deterministic `GameState`. It is a logical schema, not a language-specific class layout.

## 1. State principles

1. The server owns one canonical deterministic `GameState` per match revision.
2. Definitions are immutable; runtime instances refer to them by ID.
3. Gameplay state changes only by applying persisted gameplay events.
4. Hidden information exists in authoritative state but is omitted from unauthorized projections.
5. Every paused resolver state is serializable.
6. Historical/captured values are stored separately from current derived values.
7. Off-board presence, defeat and elimination are distinct states.
8. Card ownership, control, zone, visibility and use authority are distinct.
9. External principal identity and transport connectivity are not gameplay state.
10. Deterministic snapshots, reducers, hashes and Mechanics views exclude authority and presence registries.

## 2. Identity and ownership domains

### 2.1 PrincipalID

`PrincipalID` identifies an authenticated external subject and is derived from trusted authentication context. It is not a runtime game instance and is never accepted as an authoritative client payload value.

### 2.2 PlayerInstanceID

`PlayerInstanceID` is an opaque match-scoped runtime identity. It owns gameplay state, fighters, private zones, actions and pending choices. A client may reference it only after the application layer validates the authenticated principal's durable binding.

### 2.3 MatchAuthorityRecord

The durable application/security registry stores:

```yaml
MatchAuthorityRecord:
  match_id:
  player_instance_id:
  principal_id:
  seat:
  binding_version:
  status: ACTIVE | REVOKED
  established_by_command_id:
```

This record is not inside `GameState`, is unavailable to Mechanics and is excluded from state hashes and gameplay replay. It is nevertheless durable because reconnect and command authorization require it.

Creation rules:

- `CreateMatch` atomically creates the match, first player instance, seat and authority record;
- `JoinMatch` atomically creates the joining player instance, seat and authority record;
- later commands require an active exact binding for `(principal_id, match_id, actor_player_id)`;
- a principal cannot claim or replace a player instance through payload data.

### 2.4 OperationalPresenceRecord

A host-local or distributed operational registry may store:

```yaml
OperationalPresenceRecord:
  match_id:
  principal_id:
  player_instance_id:
  session_ids: []
  client_instance_ids: []
  online:
  last_seen_at:
```

Presence is not durable gameplay truth. It emits no gameplay event, changes no revision or state hash and cannot alter ownership, legality, timing or a pending obligation. `last_seen_at` is wall-clock operational metadata and is never read by deterministic reducers.

## 3. Deterministic logical root

```yaml
GameState:
  match_id:
  definition_ref:
  revision:
  event_sequence:
  lifecycle:
  players: {}
  fighters: {}
  cards: {}
  components: {}
  battlefield:
  turn:
  action:
  combat:
  resolver:
  random:
  game_result:
```

Maps are keyed by runtime instance IDs. Serialized ordering must be canonical even when the implementation uses unordered maps. `GameState` contains no external principal IDs, session IDs, client connection flags or last-seen timestamps.

## 4. Match lifecycle

```text
CREATED
WAITING_FOR_PLAYERS
SELECTION
SETUP
ACTIVE
ENDED
QUARANTINED
```

Transitions are event-driven. `ENDED` is immutable for gameplay commands except separately authorized administrative/replay operations. `QUARANTINED` preserves the last durable gameplay state when an internal invariant fails.

Lifecycle state records match progress, not socket/session presence. A match remains in the same gameplay lifecycle when either player disconnects.

## 5. Definition reference

```yaml
definition_ref:
  ruleset_version:
  fighter_manifest_digests: {}
  card_manifest_digests: {}
  battlefield_manifest_digest:
  capability_registry_version:
  deck_construction_result_ref:
```

A replay must load the same definitions. Updating canonical repository data must not mutate an already-started match.

## 6. Player state

```yaml
PlayerState:
  player_instance_id:
  seat:
  authority_state: ACTIVE | CONCEDED | ELIMINATED
  fighter_instance_ids: []
  private_zones: []
  resources: {}
  action_permission_ids: []
  submitted_hidden_choice_ids: []
```

`authority_state` is deterministic gameplay authority within the match; it is not authentication or connection status. It changes only through gameplay events such as concession or elimination.

No `connection_state` or `connection_status` field belongs in `PlayerState`. Disconnect cannot transfer ownership, cancel obligations, auto-select a choice, end a turn or alter legal actions.

## 7. Fighter state

```yaml
FighterInstance:
  fighter_instance_id:
  definition_id:
  owner_player_id:
  role:
  ordinal_within_definition:
  health:
  defeat_state:
  presence:
  position:
  statuses: []
  resources: {}
  modifier_source_ids: []
```

### 7.1 Defeat state

```text
UNDEFEATED
DEFEATED
```

Game-specific elimination/loss is established at canonical checkpoints, not by an unconditional `health <= 0` shortcut that bypasses prevention, redirection, replacement or simultaneous consequences.

### 7.2 Board presence

```yaml
presence:
  state: ON_BOARD | OFF_BOARD | RESERVE | DORMANT | PENDING_PLACEMENT
  reason:
  transition_event_id:
```

This `presence` is deterministic fighter board/roster presence and is unrelated to client connectivity. Only `ON_BOARD` has a battlefield position. A non-board undefeated fighter may retain health, ownership, statuses and future placement obligations.

### 7.3 Position

```yaml
position:
  kind: single_space | footprint
  space_ids: []
  orientation: null
```

Occupancy validation uses capability and battlefield rules. A multi-space fighter has one fighter instance and one footprint, not cloned fighters.

## 8. Card state

```yaml
CardInstance:
  card_instance_id:
  definition_id:
  immutable_owner_player_id:
  current_controller_player_id:
  zone_ref:
  zone_position:
  face_state:
  visibility_state:
  attached_to_ref:
  instance_state: {}
```

### 8.1 Card zones

```yaml
zone_ref:
  zone_type: deck | hand | discard | resolving | combat | auxiliary | attached | custom_registered
  owner_scope_ref:
  subzone_id:
```

Ordered zones persist complete card-instance order. Unordered-zone serialization still uses stable ordering for deterministic hashes/diffs.

### 8.2 Visibility

Visibility does not change immutable ownership. A card may be opponent-owned, controlled by another player, face up in a public zone or privately visible to one actor.

## 9. Component state

```yaml
ComponentInstance:
  component_instance_id:
  definition_id:
  owner_scope_ref:
  lifecycle: SUPPLY | DEPLOYED | DORMANT | REMOVED
  anchor:
    kind: space | edge | fighter | card | none
    ref:
  state: {}
```

Components include traps, fog/shadow/insight markers, paths, doors, items and other registered non-card objects. They are never action-card instances.

## 10. Battlefield state

```yaml
BattlefieldState:
  definition_id:
  fighter_occupancy: {}
  component_occupancy: {}
  enabled_special_connections: {}
  temporary_connections: []
  space_locks: []
  battlefield_rule_state: {}
```

The immutable definition supplies spaces, zones and base edges. Runtime state records only occupancy and source-defined deltas.

## 11. Turn state

```yaml
TurnState:
  turn_id:
  active_player_id:
  turn_number:
  action_ordinal:
  action_permission_ids: []
  turn_start_snapshot_ref:
  turn_event_range:
  end_requested:
```

The turn-start snapshot captures only values explicitly required by published behavior. It is not a duplicate full state unless snapshot policy chooses that representation internally.

## 12. Action permissions

```yaml
ActionPermission:
  permission_id:
  owner_player_id:
  source_ref:
  action_types: []
  required_target_constraints: []
  cost_policy:
  expiry:
  immediate_obligation:
  consumed_by_action_id:
```

Ordinary actions are permissions too. Gained/free/restricted actions extend the same model rather than creating bespoke counters.

## 13. Action state

```yaml
ActionState:
  action_instance_id:
  action_type:
  actor_player_id:
  permission_id:
  source_ref:
  stage:
  paid_costs: []
  movement_context:
  scheme_context:
  attack_context:
  child_effect_ids: []
```

At most one top-level action is active. Bonus attacks/effects are children with explicit parent context and do not silently consume an ordinary permission unless required by their contract.

## 14. Combat state

```yaml
CombatState:
  combat_instance_id:
  parent_action_id:
  attacker_fighter_id:
  defender_fighter_id:
  attack_type:
  attack_card_id:
  defense_card_id:
  reveal_state:
  timing_window:
  captured_values: {}
  effective_values: {}
  replacement_history: []
  damage_resolution_ref:
```

`null` defense card represents an accepted no-defense path, not an absent unresolved choice. Participant/card replacement is represented by events and history entries. Current references may change; prior participants/cards remain auditable.

## 15. Resolver state

```yaml
ResolverState:
  effect_instances: {}
  queue: []
  checkpoint_stack: []
  pending_interaction:
  simultaneous_choice_groups: {}
  delayed_obligations: []
```

### 15.1 Effect instance

```yaml
EffectInstance:
  effect_instance_id:
  definition_ref:
  source_instance_ref:
  controller_player_id:
  parent_effect_id:
  trigger_context:
  captured_bindings: {}
  current_stage:
  status: QUEUED | RESOLVING | PAUSED | RESOLVED | CANCELLED
```

### 15.2 Pending interaction

```yaml
PendingInteraction:
  interaction_instance_id:
  owner_scope:
  visibility:
  source_effect_id:
  choice_schema:
  legal_domain:
  submitted_values: {}
  resume_cursor:
  created_revision:
```

One accepted choice command either advances this exact interaction or is rejected. Reconnect does not create a replacement interaction ID or domain merely because presence changed.

### 15.3 Delayed obligation

```yaml
DelayedObligation:
  obligation_id:
  source_ref:
  trigger_condition:
  expiry_condition:
  captured_context:
  status:
```

Obligations survive turn/action boundaries and reconnect until resolved, cancelled or expired by an event.

## 16. Random state

```yaml
RandomState:
  algorithm_version:
  generator_state_ref:
  last_random_result_sequence:
```

Authoritative replay consumes persisted random-result events. Generator state supports future execution only; it is not sufficient evidence of past outcomes.

## 17. Game result

```yaml
GameResult:
  status: IN_PROGRESS | WON | DRAWN | CONCEDED | ABORTED
  winner_player_ids: []
  reason:
  established_by_event_id:
  final_sequence:
```

Game-end checks execute at documented checkpoints. Once a final result event is applied, subsequent queued gameplay effects do not continue unless canonical ordering required them before the result event.

## 18. Derived state and projections

The following are queries/caches, not canonical mutable truth unless a capability persists a historical result:

- legal actions and legal targets;
- current attack range;
- current effective card value;
- zone-based adjacency/reach;
- hand-size limit;
- movement allowance;
- winner before the game-end checkpoint completes.

Caches must be revision-keyed and reproducible from `GameState` plus immutable definitions.

A deterministic `PlayerProjection` is derived from `GameState` and validated viewer authority. Optional operational presence is projected from `OperationalPresenceRegistry` and composed only after the deterministic projection. Presence fields are excluded from projection golden hashes unless a fixture explicitly declares a separate non-gameplay presence assertion.

## 19. Lifecycle transaction invariants

For `CreateMatch` and `JoinMatch`, one application transaction must atomically persist:

- deterministic lifecycle event batch;
- resulting match head/revision;
- command idempotency record;
- required `MatchAuthorityRecord` changes.

A transaction failure exposes neither a successful command result nor a partial authority binding. Gameplay replay reconstructs the player/seat state from events; authorization reconstructs principal bindings from the authority registry.

## 20. State invariants

At every durable gameplay revision:

- each runtime ID is unique within its domain;
- every instance references an existing definition;
- every player instance has one match seat;
- every card instance exists in exactly one zone/attachment location;
- every on-board fighter has a legal position/footprint;
- occupancy indexes agree with fighter/component anchors;
- active action/combat/effect references exist;
- at most one externally actionable interaction exists except registered simultaneous committed groups;
- submitted hidden choices are visible only to authorized projections;
- event sequence and revision are monotonic;
- ended matches have one durable result;
- no definition object has been mutated;
- no principal/session/connectivity field appears in deterministic state.

At the application/security boundary:

- every active authority record references an existing match and player instance;
- one active match seat cannot be bound to multiple principals;
- a conflicting principal/seat claim is rejected rather than repaired;
- operational presence may disappear without making deterministic state unrecoverable.

Violation of a deterministic invariant quarantines the match; it is never repaired by guessing.
