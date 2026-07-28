# Engine Architecture Contract

## Status

```text
status: draft-foundation
branch: engine-foundation
base_sha: 106ae552ce597cde954c0a1b22374ef446974ce2
parent_issue: #19
implementation_authority: partial
```

This document freezes the language-neutral boundary of the competitive Unmatched runtime. It is authoritative for engine foundation work where it states a mandatory invariant. It does not select a programming language, framework, database or transport.

The first vertical slice is data-loaded Robin Hood versus Bigfoot on Sherwood Forest after the battlefield graph passes #18.

## 1. Scope

The runtime is responsible for:

- loading immutable game definitions;
- creating and advancing a two-player match;
- validating commands against authoritative state;
- generating deterministic events;
- resolving effects, timing windows and choices;
- persisting random outcomes and audit history;
- projecting only authorized information to each player;
- restoring a match after process restart or client reconnect;
- detecting exhaustion, defeat and game end.

The runtime is not responsible for:

- interpreting battlefield artwork during play;
- deriving rules from card image assets;
- trusting client-computed legality;
- cooperative Adventures enemy/scenario behavior;
- changing published behavior through balance patches;
- embedding fighter/card IDs in core control flow.

## 2. Architectural invariant

```text
Command
→ validate actor, version and legality against GameState
→ emit one or more Events
→ append Events in a total order
→ apply Events to GameState
→ enqueue/resolve Effects and timing checkpoints
→ pause on a typed PendingInteraction when input is required
→ emit player-specific Projections
```

A command never mutates state directly. An event is the only persisted statement that authoritative state changed or a persisted random/choice result was established.

## 3. Definitions versus runtime instances

### 3.1 Immutable definitions

Definitions are loaded from versioned canonical data and are immutable for the lifetime of a match:

- `GameDefinition`
- `FighterDefinition`
- `CardDefinition`
- `BattlefieldDefinition`
- generic capability definitions
- source-defined setup and construction definitions

Every definition has a stable string `definition_id` and a content/version identity sufficient to reproduce a match with the same rules corpus.

### 3.2 Runtime instances

Runtime objects have opaque instance IDs unique inside a match:

- `player_instance_id`
- `fighter_instance_id`
- `card_instance_id`
- `component_instance_id`
- `effect_instance_id`
- `interaction_instance_id`
- `combat_instance_id`
- `action_instance_id`

A definition ID is never used as a runtime instance ID. Multiple physical copies of one card definition produce distinct card instances. Multiple sidekicks with one definition produce distinct fighter instances.

### 3.3 Asset separation

A card definition may reference an optional asset-registry key. Binary artwork is not gameplay state and must not affect legality, resolution, replay hashes or definition identity unless an explicit asset-registry version is being audited.

## 4. Authoritative aggregates

### 4.1 GameDefinition

Contains the immutable match rules package:

```text
ruleset_version
fighter_definitions
card_definitions
battlefield_definition
deck_construction_results
capability_registry_version
setup_definition
```

Pre-game random or player-selected construction results are persisted as events and then referenced by the effective match definition/state.

### 4.2 GameState

The canonical state aggregate contains at minimum:

```text
match_id
definition_version
revision
phase
players
fighters
cards
components
battlefield
turn_state
action_state
combat_state
pending_interaction
effect_queue
checkpoint_stack
random_state_ref
history_cursor
game_result
```

`revision` increases monotonically after each accepted command transaction. A command must declare the revision it was based on or an idempotency key that safely resolves duplicate delivery.

### 4.3 PlayerState

Contains authoritative player-owned runtime state, not a client projection:

```text
player_instance_id
seat
connection_status
owned_fighter_instances
hand_zone
private_choice_state
action_permissions
resources
```

Connection status cannot change gameplay ownership or choice authority.

### 4.4 FighterInstance

Contains identity and mutable fighter state:

```text
fighter_instance_id
definition_id
owner_player_id
role
health
presence_state
position_or_footprint
statuses
resources
source_bound_modifiers
```

`presence_state` is distinct from defeat. A fighter may be off-board, dormant, in reserve or pending placement while still undefeated if a canonical capability permits it.

### 4.5 CardInstance

Contains:

```text
card_instance_id
definition_id
immutable_owner_player_id
current_controller_player_id
zone
zone_position
visibility_state
face_state
attached_to
instance_state
```

Ownership, control, current zone, visibility and permission to use the card are independent fields/concepts.

### 4.6 BattlefieldState

Contains only runtime deltas from the immutable graph:

```text
fighter_occupancy
component_occupancy
path_components
temporary_connections
space_locks
battlefield_item_state
```

Adjacency and zones come from `BattlefieldDefinition`; runtime state may enable, disable or annotate connections only through canonical capabilities.

## 5. Commands

### 5.1 Command envelope

Every command contains:

```text
command_id
match_id
actor_player_id
expected_revision
command_type
payload
submitted_at_client_optional
```

`submitted_at_client` is diagnostic only and never orders gameplay.

### 5.2 Validation order

The server validates:

1. envelope and schema;
2. authenticated actor and match membership;
3. idempotency / duplicate command ID;
4. expected revision policy;
5. current phase and pending-interaction ownership;
6. legal action/target domain generated from current state;
7. command-specific costs and prerequisites;
8. hidden-information constraints;
9. deterministic payload normalization.

A rejected command emits no gameplay event. Operational rejection diagnostics are not part of authoritative match history unless security/audit policy separately records them.

### 5.3 Initial command families

The contract must support at least:

```text
CreateMatch
JoinMatch
SelectFighter
SelectBattlefield
ConfirmSetup
StartTurnAction
ChooseManeuverBoost
ChooseMovementPath
PlayScheme
DeclareAttack
ChooseDefense
SubmitInteractionChoice
Concede
```

Concrete payload schemas belong in `command-event-contract.md`.

## 6. Events

### 6.1 Event envelope

Every persisted event contains:

```text
event_id
match_id
sequence
caused_by_command_id
parent_event_id optional
source_instance_ref optional
event_type
public_payload
private_payloads_by_player optional
ruleset_version
```

`sequence` is the authoritative total order. Wall-clock timestamps are metadata and cannot determine gameplay order.

### 6.2 Event properties

Events must be:

- deterministic;
- serializable;
- replayable;
- sufficient to restore authoritative state together with an approved snapshot;
- explicit about private payload ownership;
- explicit about random and choice results;
- stable enough for audit/debug tooling.

An event may describe a request to open a timing window or interaction, but a player decision is persisted as a separate result event.

## 7. Effect and checkpoint resolver

### 7.1 Effect instances

A normalized effect definition produces an `EffectInstance` containing:

```text
effect_instance_id
source_definition_id
source_instance_id
controller_player_id
trigger_context
captured_bindings
stages
current_stage
status
parent_effect_id optional
```

Captured bindings are immutable historical facts. References to current state are reevaluated only when the definition explicitly calls for a current/derived value.

### 7.2 Queues and stacks

The runtime maintains:

- an ordered event application stream;
- an effect-resolution queue;
- a checkpoint/timing-window stack;
- at most one externally actionable `PendingInteraction` at the match boundary unless a canonical simultaneous-choice capability explicitly creates committed independent submissions.

Internal deterministic operations may continue until:

- the queue is empty;
- a checkpoint requires a choice;
- a simultaneous committed choice group is incomplete;
- the game ends;
- a validation invariant fails and the match is quarantined for diagnostics.

### 7.3 PendingInteraction

A pending interaction contains:

```text
interaction_instance_id
interaction_type
owning_player_or_group
visibility
source_effect_id
prompt_definition
legal_domain_snapshot_or_query_contract
cardinality
optional
submitted_choices
resume_token
created_at_revision
```

The entire interaction is serializable. Reconnect returns the same authorized prompt and legal domain unless an accepted event has legitimately invalidated it; invalidation must itself be deterministic and recorded.

## 8. Turn, action and combat state machines

### 8.1 TurnState

Tracks:

```text
active_player
action_ordinal
action_permissions
turn_start_snapshot
turn_history_refs
end_turn_requested
```

Action permissions are typed runtime objects. They can represent ordinary actions, gained actions, free/restricted actions, immediate obligations and expiring permissions. They are not reducible to one integer counter.

### 8.2 ActionState

Tracks the current action type, actor, source, costs, stage, movement/combat context and child effects.

### 8.3 CombatState

Tracks:

```text
combat_instance_id
attacking_fighter_instance_id
defending_fighter_instance_id
attack_type
attack_card_instance_id
defense_card_instance_id optional
reveal_state
current_window
printed_values_at_required_windows
effective_values
damage_result
participant_replacement_history
```

Participant/card replacement updates the current combat references through events while preserving prior history. Printed, effective and captured historical values remain distinct.

## 9. Legal-action generator

The engine exposes an authoritative query:

```text
GetLegalActions(player_id, game_state_revision) -> LegalActionSet
```

Each legal action descriptor includes a stable action kind, source, required choices, target-domain query/result, cost preview and visibility-safe UI metadata.

The client may cache or render this set but cannot expand it. The server regenerates or validates legality on command acceptance.

The legal-action generator must be pure with respect to authoritative input state: it cannot mutate state, consume RNG or reveal unauthorized information.

## 10. Deterministic randomness

### 10.1 Rule

Randomness is requested by a deterministic operation and resolved into a persisted result event before dependent state changes.

### 10.2 Requirements

- the random algorithm/version and seed material are server-controlled;
- an accepted random result is never recomputed during reconnect or replay;
- shuffles persist the resulting card-instance order or an equivalent replay-sufficient permutation result;
- random choices persist the selected instance IDs;
- tests may inject a deterministic random-result provider;
- projections reveal only the part of a random result authorized by game rules.

The random generator's internal mutable state may be snapshotted for forward execution, but replay correctness depends on persisted outcomes, not on regenerating the same stream.

## 11. Visibility projections

The server stores one authoritative state and derives a projection for:

```text
player A
player B
spectator/admin policy optional
```

A projection contains only information authorized at the current revision. Hidden card identities, private choices and private random results are omitted or replaced by opaque counts/instance placeholders as required.

Projection code cannot change state. Projection tests must prove that adding a private field to authoritative state does not automatically expose it.

Committed-hidden simultaneous choices are stored privately per actor until the reveal condition is satisfied by an event.

## 12. Persistence, reconnect and replay

### 12.1 Persistence model

The authoritative persistence unit is:

```text
latest approved snapshot
+ ordered events after snapshot sequence
+ immutable definition version references
```

Snapshots are optimization artifacts. Events remain the audit source for changes after the snapshot.

### 12.2 Atomic command transaction

An accepted command transaction must atomically persist:

- idempotency result;
- emitted event batch;
- resulting match revision;
- any updated snapshot pointer required by policy.

No client receives a successful command result before the authoritative event batch is durable.

### 12.3 Reconnect

Reconnect loads authoritative state at the latest revision, authenticates the player, and emits:

- current player projection;
- current legal action set or pending interaction;
- monotonic revision/sequence cursor;
- optional public history summary.

A reconnect never restarts an effect, repeats a draw/shuffle or changes a choice domain merely because a connection was lost.

### 12.4 Replay

Replay applies persisted events against the exact referenced definition/capability versions. Replay may stop at any event sequence for debugging and projection inspection.

## 13. Generic capability boundary

A fighter/card manifest may declare generic capabilities and parameter data. The core engine dispatches by operation/capability type, never by fighter/card identity.

A new extension handler requires all of:

1. evidence that existing generic operations cannot represent published behavior;
2. a frozen handler contract with typed input/output and state boundary;
3. no fighter/card ID condition inside shared core logic;
4. deterministic event, choice, visibility and persistence behavior;
5. reusable applicability beyond a single textual card where objectively possible;
6. explicit registration in the capability registry;
7. independent QA approval.

A handler name that merely disguises a fighter-specific branch does not satisfy this policy.

## 14. Failure and quarantine behavior

A violated internal invariant must not be repaired by guessing. The match enters a diagnostic quarantine state that blocks further gameplay commands while preserving:

- last durable revision;
- triggering command/event IDs;
- resolver stack/queue diagnostics;
- definition versions;
- projection-safe player message.

Operational retries may repeat transport or transaction work only through command idempotency. They may not emit a second gameplay result.

## 15. Vertical-slice acceptance

The foundation is ready for implementation only after related contracts specify deterministic behavior for:

- create/join and seat authority;
- fighter/battlefield selection;
- setup and starting hands;
- Maneuver, movement and BOOST;
- Scheme;
- Attack, optional defense and all combat windows;
- pending choices;
- damage and defeat;
- exhaustion and game end;
- disconnect/reconnect;
- Robin Hood's optional post-attack movement;
- Bigfoot's optional zone-based end-turn draw.

Sherwood Forest graph data is a prerequisite for movement/zone legality tests, not for writing the language-neutral state/command/event contracts.

## 16. Foundation gate

This contract can advance from `draft-foundation` when:

- Phase 4C maps all launch-scope owner requirements to canonical capabilities;
- state, command, event, choice, RNG, visibility and persistence companion documents agree;
- the deterministic fixture format is defined;
- the runtime stack ADR is accepted;
- independent architecture QA finds no unresolved P0/P1 launch-scope ambiguity;
- initial code proves definitions are loaded as data and core logic contains no fighter/card ID branches.
