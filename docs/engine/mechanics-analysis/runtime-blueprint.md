# Mechanics-Driven Runtime Blueprint

## 1. Architectural shape

The full competitive corpus should run on four cooperating deterministic layers:

```text
state machines        turn / action / combat / timing checkpoints
rule interpreter      procedures / stages / interactions / operations
pure derivation       selectors / expressions / legality / modifiers
event authority       events / history / persistence / projections
```

These layers are distinct even when one source definition touches all of them. A card definition describes behavior; it does not own a private execution engine.

## 2. Definition model

### 2.1 Rule definitions

Immutable match definitions contain typed records:

```text
EffectDefinition
ProcedureDefinition
StageDefinition
InteractionDefinition
OperationDefinition
QueryExpression
ModifierDefinition
CapabilityDefinition
CompositeDefinition
```

Every record has a stable definition ID. Runtime execution creates opaque instance IDs. A definition ID may identify semantics or data but must never be used as a runtime object identity.

### 2.2 Closed typed AST

The definition language must be closed and versioned. Unknown kinds fail loading before a match starts.

It is not a general-purpose scripting language:

- expressions are pure and cannot mutate state;
- operations are the only domain mutation requests;
- loops exist only as bounded procedure kinds;
- no arbitrary function names or reflection-based invocation;
- no network, filesystem, SQL or wall-clock access;
- no fighter/card ID switch statements.

## 3. Resolver kernel

### 3.1 Procedure instance

A resolving effect creates a serializable procedure instance:

```text
procedure_instance_id
procedure_definition_id
source_instance_ref
controller_player_id
trigger_context
current_stage_id
captured_bindings
loop_state
child_procedure_refs
status
```

Captured bindings are historical facts. Current-state expressions are reevaluated only when their definition explicitly asks for current state.

### 3.2 Stage algorithm

For each stage:

1. Evaluate the stage condition using the pure query engine.
2. If false, emit/record a deterministic skip result and advance.
3. Validate prerequisites and costs.
4. Materialize the interaction domain, when player input is needed.
5. Persist the submitted choice result.
6. Execute operations in source order.
7. Persist operation result events and bind declared results.
8. Open any required timing checkpoints or nested procedures.
9. Advance, repeat, stop or complete according to the typed procedure definition.

An impossible independent operation is skipped while later independent operations continue. An explicit dependency edge prevents only its dependent continuation. This distinction must be represented in the definition rather than inferred from card prose at runtime.

### 3.3 Procedure kinds

The minimum non-general-purpose procedure set is:

- `linear_stages`;
- `branch_by_bound_value`;
- `repeat_for_each` with a materialized ordered subject list;
- `resumable_repeat_stop` with optional continue and forced-stop predicates;
- `multi_actor_arbitration` with deterministic close conditions;
- `composite_combat_procedure` for bonus attacks or same-combat replacement;
- `scheduled_checkpoint_procedure` for delayed turn obligations.

New procedure kinds require corpus evidence and architecture review.

## 4. Query and expression engine

The query engine is pure, deterministic and visibility-aware at its boundary.

### 4.1 Selectors

Selectors return runtime instance references, never definition IDs pretending to be instances.

Required selector dimensions include:

- owner/controller/relation;
- fighter role, presence and defeat state;
- card zone, position, tags and base/effective fields;
- space, zone, adjacency, graph distance and path network;
- component type/state/anchor;
- combat role and active combat identity;
- event source, actor, cause and turn/action context;
- historical movement path or captured operation result.

### 4.2 Expressions

The initial expression AST should support:

- literals and references;
- boolean composition;
- equality/comparison;
- `count`, `sum`, `min`, `max`;
- conditional expression;
- `highest_met_threshold`;
- field access on typed bound results;
- current-state versus captured-history references.

Expressions cannot consume randomness or emit interactions.

### 4.3 Domain materialization

A pending interaction stores either:

- an immutable legal-domain snapshot; or
- a query contract plus the revision at which reevaluation is allowed.

Private domains are projected using authorized fields or opaque stable handles. The server always validates the submitted handle against authoritative state.

## 5. Timing and checkpoints

Effects trigger from typed events and checkpoints, not string matching.

The minimum checkpoint model includes:

```text
turn_start
before_action
maneuver_declared
movement_step_entered
scheme_played
attack_declared_pre_defense
defense_commitment
combat_reveal
immediately_window
during_combat_window
combat_damage
post_damage
post_combat
combat_cleanup
end_action
end_turn
```

The exact ordering is defined by the core rules contract. Source capabilities may register eligible effects at a checkpoint but cannot invent an unordered parallel timing model.

A checkpoint stack records:

- checkpoint identity and parent context;
- active combat/action/procedure references;
- eligible effects/reactions;
- priority/ordering authority;
- canceled or already-resolved effects;
- resume location.

## 6. Cancellation and reactions

Cancellation targets an effect instance or an explicitly defined scope. It is not a generic rollback of already emitted state changes.

Supported scopes include:

- one pending effect instance;
- remaining effects of a committed combat card;
- a pending operation before application;
- a complete bound source procedure when the published rule requires whole-effect cancellation.

The resolver records cancellation as an event and marks the target state. Protected effects and cleanup procedures remain non-cancelable according to definition metadata.

Private-zone reactions are permissions generated from authoritative state. The client never scans its hand locally to create reaction legality.

## 7. History and provenance

Every operation result carries provenance sufficient for later rules:

```text
operation_instance_id
source_effect_instance_id
source_card_or_ability_ref
actor_player_id
controller_player_id
affected_instance_refs
cause_kind
parent_operation_ref optional
result_payload
disposition
turn_id
action_instance_id optional
combat_instance_id optional
```

Canonical history queries are derived from immutable events or rebuildable indexes. Separate event kinds distinguish:

- played card;
- discarded card;
- revealed/looked-at card;
- card used as BOOST;
- source-caused damage;
- propagated damage;
- movement path/space entry;
- resource/action permission changes.

## 8. Battlefield and fighter presence

### 8.1 Presence

`presence_state` is independent from health/defeat:

```text
on_board
off_board_undefeated
reserve
pending_placement
defeated
out_of_play
```

Only the active match roster can become targetable. Definitions excluded by setup remain `out_of_play` and cannot satisfy selectors or revival effects.

### 8.2 Position and footprint

A fighter position is a footprint value rather than one mandatory space:

```text
anchor_space
occupied_spaces[]
orientation optional
occupancy_class
```

Ordinary fighters use a one-space footprint. Large and small fighters reuse the same occupancy API with definition-supplied policies.

### 8.3 Components

Components are stable runtime instances:

```text
component_instance_id
definition_id
owner/controller
pool/state
anchor: supply | space_ref | path_ref | attached_ref
occupancy_policy
source_lifetime
```

Path absence may make a path component dormant without destroying its identity when source rules require that lifecycle.

## 9. Movement and placement

Movement is a procedure, not one final position assignment.

A movement instance records:

- origin footprint;
- chosen path;
- ordered entered spaces and traversed paths;
- movement context: Maneuver, effect move, forced move, placement or swap;
- movement-point cost after contextual modifiers;
- interruption status;
- destination footprint.

Each resolved entry can emit a checkpoint. An entry reaction may interrupt the remaining relocation while preserving the reached space and ordered follow-up effects.

Placement performs no path traversal but still validates occupancy, locks and destination policy. Failed placement produces a typed result for dependent stages.

## 10. Combat engine

Combat is one persistent aggregate from attack declaration through cleanup.

```text
combat_instance_id
attacking_fighter_instance_id
defending_fighter_instance_id
attack_type
selected_attack_mode
attack_card_instance_id
defense_card_instance_id optional
card owner/controller identities
current_window
captured printed/effective values
participant replacement history
damage transaction ref
```

### 10.1 Legality

Attack legality derives from:

- base fighter/card definitions;
- current footprint and graph/path network;
- selected play mode;
- source-lifetime modifiers;
- temporary locks and permissions;
- card visibility/commitment rules.

Printed type/value/BOOST metadata remains immutable. Effective attack type, printed-value view, combat value and BOOST result are separate derived fields.

### 10.2 Rebinding

Participant or card replacement updates references inside the same combat and emits history. It does not create a new Attack action unless a typed bonus-attack composite explicitly does so.

The replacement operation must declare which context is preserved:

- combat ID;
- opposing participant;
- committed cards;
- timing window;
- captured values;
- controller/card ownership provenance.

## 11. Damage and health transaction

Damage resolves through a typed transaction:

```text
DamageProposed
→ prevention/replacement window
→ allocation/redirection window
→ per-recipient application
→ health-result modifiers/clamps
→ defeat checks
→ damage/result observers
```

The transaction preserves original source, damage type, primary target and propagated/allocated recipients.

Required distinctions:

- damage versus direct defeat;
- damage prevention versus health floor;
- recovery versus exact health assignment;
- living exact `SET_HEALTH` versus atomic return/revive health;
- redirected/allocated damage versus newly created damage;
- amount dealt versus final health reduction after a clamp.

## 12. Modifier layer

Base definitions never mutate. Effective values are derived from base data plus active modifiers.

A modifier instance contains:

```text
modifier_kind
target selector or bound instance
value/permission rule
context predicate
source instance
lifetime/expiry
priority/composition rule
```

Modifier families include:

- movement value/cost;
- attack type/reach/target network;
- card play mode and commitment visibility;
- health-result floor;
- hand size;
- printed/effective/current combat values;
- BOOST value/result transform;
- action permissions and operation locks.

Composition order must be explicit by modifier family. There is no one universal arithmetic stacking rule.

## 13. Card zones and auxiliary objects

Zones are typed runtime containers with ownership, ordering, visibility and lifecycle policy.

They may contain:

- ordinary action-card instances;
- auxiliary gameplay-card instances;
- non-card auxiliary objects when the source component behaves as an ordered-zone object;
- public token rows with stable object identity.

Zone operations always work against live membership. Cleanup cannot discard a stale snapshot when a completion effect moved an object elsewhere.

Structured search records:

- zones searched;
- predicate;
- chooser/viewer policy;
- selected instance;
- disclosure;
- destination;
- post-search disposition such as shuffle.

## 14. BOOST pipeline

All BOOST use cases share one pipeline:

```text
source committed
→ source card/random/dynamic value resolved
→ result bound
→ source-defined transforms
→ amount applied to target context
→ source disposition
→ history event
```

Blind BOOST, ordinary hand BOOST and dynamic die-based BOOST differ only in source-resolution definitions. Accepted random values and the exact source card instance are persisted.

## 15. Turn actions, resources and obligations

Action accounting uses typed permissions, not one integer.

An action permission may specify:

```text
origin
allowed_action_kinds
must_be_next
required_target optional
spendable_as_cost
expires_at
restricted/free/ordinary classification
```

Resource declarations are separated from available quantity when published rules allow underfunded declaration. The chosen tier remains bound even when only the available amount is consumed.

Delayed obligations record subject player, lifetime, satisfaction event and expiry consequence. Multi-actor payment windows use deterministic arbitration; the first accepted legal payment closes the window and produces an atomic resource transfer.

## 16. Randomness

A random operation opens an authoritative random procedure. Results are events before dependent effects resolve.

Paid rerolls replace the bound current result through additional events while retaining complete history. Replay reads accepted results; it never attempts to reproduce them from wall-clock timing or client randomness.

## 17. Visibility and projections

Authoritative state may contain hidden identities and private choices. Projection code is deny-by-default.

The runtime distinguishes:

- public identity;
- owner-only identity;
- chooser-visible fields;
- opaque instance handle/card back;
- revealed field without full identity;
- committed-hidden choice pending later reveal.

Knowledge already granted by an event remains available according to the rules even if the card later changes zone; this is modeled as authorized knowledge/history, not client memory trust.

## 18. Suggested Go domain packages

The package names are proposals; dependency direction is normative.

```text
internal/domain/model
internal/domain/query
internal/domain/resolver
internal/domain/timing
internal/domain/operations
internal/domain/history
internal/domain/zones
internal/domain/battlefield
internal/domain/movement
internal/domain/combat
internal/domain/damage
internal/domain/modifiers
internal/domain/turn
internal/domain/random
internal/domain/projection
internal/domain/registry
```

`model` contains shared value types and aggregates. Feature packages may depend on `model` and pure `query` contracts but must avoid cyclic ownership. The resolver coordinates registered typed operations through interfaces defined inside the domain layer.

## 19. Definition loading and capability registration

Loading a match definition must validate:

- every operation/query/procedure kind is registered;
- all references and bindings are type-correct;
- stage references point backward unless a procedure explicitly declares loop state;
- no private field is requested by an unauthorized interaction projection;
- declared capability dependencies are present;
- no definition requests an unapproved custom handler;
- all operation results used later are explicitly bound;
- persistent procedures are serializable.

## 20. Extension-handler gate

A custom composite/handler may be registered only when:

1. a concrete corpus behavior cannot be represented by the closed generic model;
2. the failed generic representation and semantic loss are documented;
3. inputs, outputs, state, choices, timing, visibility and persistence are typed;
4. registration is by capability metadata, never fighter/card ID dispatch;
5. deterministic fixtures cover interruption, reconnect and replay;
6. independent architecture QA approves it.

At the current analysis stage, no source other than the unresolved Deadpool adaptation policy proves a need for arbitrary fighter-specific runtime code.
