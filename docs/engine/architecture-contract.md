# Engine Architecture Contract

## Status

```text
status: draft-foundation
branch: engine-foundation
base_sha: 106ae552ce597cde954c0a1b22374ef446974ce2
parent_issue: #19
correction_issue: #32
implementation_authority: partial
```

This document freezes the language-neutral boundary of the competitive Unmatched runtime. It is authoritative where it states a mandatory invariant. It does not select a concrete framework, authentication provider, transport schema technology or database library.

The first vertical slice is data-loaded Robin Hood versus Bigfoot on Sherwood Forest after the battlefield graph passes #18.

## 1. Scope

The runtime is responsible for:

- loading immutable game definitions;
- creating and advancing a two-player match;
- binding authenticated principals to match-scoped player instances;
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
- treating transport sessions as gameplay state;
- cooperative Adventures enemy/scenario behavior;
- changing published behavior through balance patches;
- embedding fighter/card IDs in core control flow.

## 2. Architectural invariant

For commands that require an existing match seat:

```text
authenticated Principal
→ resolve durable Principal-to-PlayerInstance binding
→ validate normalized command identity and idempotency
→ validate revision and legality against GameState
→ emit one or more Events
→ append Events in a total order
→ apply Events to GameState
→ enqueue/resolve Effects and timing checkpoints
→ pause on a typed PendingInteraction when input is required
→ emit player-specific deterministic Projections
→ optionally compose non-gameplay operational presence
```

`CreateMatch` and `JoinMatch` use the lifecycle-specific authorization order in section 5. They do not require a pre-existing `PlayerInstanceID`.

A command never mutates `GameState` directly. A persisted event is the only statement that deterministic gameplay state changed or that a persisted random/choice result was established. Principal bindings, idempotency records and operational connection presence are separate application/security records with the ownership rules below.

## 3. Identity domains

### 3.1 External principal identity

`PrincipalID` identifies an authenticated external subject. It is derived from trusted authentication/session context, not accepted as an authoritative client payload field. This contract does not choose an authentication provider or token format.

A principal may have multiple transport sessions or client instances. Session and client identifiers are operational metadata and never substitute for `PrincipalID` or `PlayerInstanceID`.

### 3.2 Immutable definitions

Definitions are loaded from versioned canonical data and remain immutable for the lifetime of a match:

- `GameDefinition`;
- `FighterDefinition`;
- `CardDefinition`;
- `BattlefieldDefinition`;
- generic capability definitions;
- source-defined setup and construction definitions.

Every definition has a stable `definition_id` and a content/version identity sufficient to reproduce a match with the same rules corpus.

### 3.3 Runtime instances

Runtime objects have opaque instance IDs unique inside a match:

- `player_instance_id`;
- `fighter_instance_id`;
- `card_instance_id`;
- `component_instance_id`;
- `effect_instance_id`;
- `interaction_instance_id`;
- `combat_instance_id`;
- `action_instance_id`.

A definition ID is never used as a runtime instance ID. Multiple physical copies of one card definition produce distinct card instances. Repeated fighter definitions and mirror matches remain instance-scoped.

### 3.4 Principal-to-seat authority

The durable `MatchAuthorityRegistry` maps:

```text
(match_id, player_instance_id) ↔ principal_id
```

The registry belongs to the application/security boundary, not to Mechanics and not to deterministic `GameState`. A principal binding cannot be inferred from a client-supplied player ID.

`CreateMatch` atomically creates the match, first `PlayerInstanceID`, seat and authority binding. `JoinMatch` atomically creates the joining `PlayerInstanceID`, assigns the seat and creates the authority binding. Later commands authenticate a principal and require the registry binding to match the command's `match_id` and `actor_player_id`.

### 3.5 Asset separation

A card definition may reference an optional asset-registry key. Binary artwork is not gameplay state and must not affect legality, resolution, replay hashes or definition identity unless an explicit asset-registry version is being audited.

## 4. Authoritative aggregates and operational registries

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

### 4.2 Deterministic GameState

The canonical deterministic aggregate contains at minimum:

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

`revision` increases monotonically after each accepted command transaction. Connectivity, transport sessions, last-seen timestamps and external principal identifiers are excluded from `GameState`, deterministic reducers, snapshots, replay hashes and Mechanics views.

### 4.3 PlayerState

Contains match-scoped gameplay state, not external identity or connection state:

```text
player_instance_id
seat
owned_fighter_instances
hand_zone
private_choice_state
action_permissions
resources
```

Gameplay ownership and pending-choice authority are represented by `PlayerInstanceID` and deterministic state. They do not change when a client disconnects.

### 4.4 Operational presence

A match host may maintain an `OperationalPresenceRegistry` containing session/client references, online status and last-seen metadata. It is outside deterministic state and may be rebuilt from active connections.

Presence changes:

- emit no gameplay event;
- change no match revision or event sequence;
- affect no legal action, timing window, pending interaction or state hash;
- are unavailable to Mechanics;
- may be exposed only through an explicitly authorized operational projection.

A timeout/forfeit product policy, when enabled, must enter gameplay through an explicit authenticated command and persisted event; connection loss alone is never that command.

### 4.5 FighterInstance

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

### 4.6 CardInstance

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

Ownership, control, current zone, visibility and permission to use the card are independent concepts.

### 4.7 BattlefieldState

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

## 5. Commands, lifecycle authorization and idempotency

### 5.1 Common transport-neutral fields

Every command request has:

```text
command_id
command_schema_version
command_type
payload
client_metadata optional
```

The authenticated `principal_id` is supplied by trusted request context. Client timestamps and client-instance metadata are diagnostic only.

Lifecycle-specific fields are:

| Command scope | `match_id` | `actor_player_id` | `expected_revision` |
| --- | --- | --- | --- |
| `CreateMatch` | absent; allocated by server | absent; allocated by server | absent |
| `JoinMatch` | required | absent; allocated by server | required |
| Existing-seat commands | required | required | required |

A transport may carry explicit `null` for an absent field only if its schema version defines that representation. Canonical normalization treats omitted and explicit-null forms according to that schema; implementations must not disagree silently.

### 5.2 CreateMatch authorization order

1. parse the envelope sufficiently to identify schema version and command ID;
2. authenticate the external principal;
3. normalize the lifecycle request and compute its idempotency fingerprint;
4. perform idempotency lookup in the principal lifecycle namespace;
5. validate creation policy and payload;
6. allocate `MatchID` and first `PlayerInstanceID` through deterministic/injected ID providers;
7. atomically persist authority binding, idempotency result, event batch and match head.

No existing match membership or runtime player identity is required.

### 5.3 JoinMatch authorization order

1. parse and authenticate the external principal;
2. normalize the request, including target `match_id` and expected revision;
3. perform idempotency lookup;
4. load the target match and validate lifecycle/revision/join policy;
5. prove the principal is not already bound incompatibly;
6. allocate a new `PlayerInstanceID`, assign a seat and create the authority binding;
7. atomically persist authority binding, idempotency result, event batch and match head.

Membership validation for the new player begins only after the binding is created. The joining principal cannot claim a client-chosen player instance.

### 5.4 Existing-seat command authorization order

1. parse and authenticate the external principal;
2. normalize the request and compute its fingerprint;
3. perform idempotency lookup;
4. resolve the authoritative `(principal_id, match_id, actor_player_id)` binding;
5. validate expected revision;
6. validate lifecycle, pending-interaction ownership and legal action/target domain;
7. validate costs, prerequisites and hidden-information constraints;
8. produce and atomically persist the result.

A client-supplied `actor_player_id` is only a claimed match-scoped actor reference. It becomes authoritative only after binding validation.

### 5.5 Normalized request identity

The idempotency key is `(principal_id, command_id)`. Its durable record contains a versioned fingerprint computed from canonical serialization of:

```text
fingerprint_schema_version
principal_id
lifecycle_scope
match_id or explicit absent marker
actor_player_id or explicit absent marker
command_schema_version
command_type
normalized_payload
expected_revision_policy
expected_revision or explicit absent marker
```

Client timestamps, trace IDs, transport headers, session IDs and retry counters are excluded.

Canonical serialization and hash rules are defined in `deterministic-fixture-contract.md`. The initial fingerprint algorithm is `sha256` over canonical UTF-8 JSON bytes.

For the same `(principal_id, command_id)`:

- the same fingerprint returns the same durable accepted or rejected result and emits no second event batch;
- a different fingerprint returns `DUPLICATE_CONFLICT`, returns no prior unrelated payload and never executes the conflicting request;
- an infrastructure failure before durable commit is not a command result and may be retried;
- unauthenticated or structurally unparseable input has no principal-scoped idempotency record;
- after authentication and successful canonical normalization, every terminal accepted or deterministic rejected result is durably recorded.

## 6. Events and atomicity

Every persisted event contains:

```text
event_id
match_id
sequence
revision
caused_by_command_id
parent_event_id optional
source_instance_ref optional
event_type
public_payload
private_payloads_by_player optional
ruleset_version
event_schema_version
```

Events are deterministic, serializable, replayable, explicit about random/choice results and private payload ownership, and sufficient to restore authoritative state together with an approved snapshot.

An accepted lifecycle transaction atomically persists the gameplay event batch, match head, idempotency record and required principal-to-player authority changes. No client receives success before all required records are durable.

Rejected commands emit no gameplay event. Their durable idempotency records are security/application records, not match replay events.

## 7. Effect and checkpoint resolver

A normalized effect definition produces a serializable `EffectInstance` containing source identity, controller, trigger context, captured bindings, ordered stages, current stage, status and optional parent reference.

Captured bindings are immutable historical facts. Current-state references are reevaluated only when a definition explicitly requests a current/derived value.

The runtime maintains an ordered event stream, effect-resolution queue, checkpoint/timing-window stack and at most one externally actionable `PendingInteraction` unless a registered simultaneous committed-choice capability creates independent submissions.

A pending interaction contains a stable interaction ID, owner scope, visibility, source effect, prompt/choice definition, legal domain or query contract, cardinality, optionality, submitted values, resume cursor and creation revision. It is fully serializable. Reconnect returns the same authorized prompt and domain unless a persisted event deterministically invalidated it.

## 8. Turn, action and combat state machines

`TurnState` preserves active player, action ordinal, typed action permissions, required turn-start history and end-turn state. Ordinary, gained, free, restricted and immediate actions use typed permissions rather than one undifferentiated counter.

`ActionState` preserves action identity, actor, source, costs, stage, movement/combat context and child effects.

`CombatState` preserves combat identity, current and historical participants/cards, attack type, reveal/window state, printed/effective/captured values, damage result and replacement history. Participant or card replacement updates current references through events while retaining prior history.

## 9. Legal-action boundary

The engine exposes:

```text
GetLegalActions(player_instance_id, game_state_revision) -> LegalActionSet
```

The caller must already possess a validated principal-to-player binding. The generator is pure, deterministic, revision-scoped, visibility-safe, unable to consume RNG or mutate state, and dispatches by operation/query/procedure/capability kind rather than fighter/card identity.

## 10. Deterministic randomness

Randomness is requested by a deterministic operation and resolved into a persisted result event before dependent state changes. Accepted random results are never recomputed during reconnect or replay. Shuffles persist the resulting card-instance order or an equivalent complete permutation; tests may inject deterministic outcomes.

## 11. Visibility and projection composition

The deterministic projector derives one authorized view from `GameState` and immutable definitions for a validated `ViewerAuthority` containing the match and player instance. It does not accept an unverified client player ID.

Hidden card identities, private choices and private random results are omitted or represented only by authorized opaque abstractions. Projection code is deny-by-default and cannot change state.

Operational presence, when exposed, is composed after deterministic projection:

```text
project(GameState, ViewerAuthority) -> PlayerProjection
projectPresence(OperationalPresenceRegistry, ViewerAuthority) -> PresenceProjection
compose(PlayerProjection, PresenceProjection) -> DeliveryEnvelope
```

`PresenceProjection` has no gameplay revision, is excluded from deterministic fixture hashes and cannot influence legal-action generation.

## 12. Persistence, reconnect and replay

A recoverable match requires:

```text
immutable definition-version references
+ latest approved deterministic snapshot
+ ordered gameplay events after snapshot sequence
+ command idempotency records
+ durable principal-to-player authority records
```

Snapshots contain deterministic `GameState`, including private gameplay state, but exclude principal identity, sessions and connectivity.

Reconnect:

1. authenticates the principal;
2. resolves the durable match authority binding;
3. restores `GameState` from snapshot plus ordered event tail;
4. derives the same authorized deterministic projection and pending interaction as live execution at that revision;
5. registers current session presence separately.

Reconnect never changes gameplay revision/hash, restarts an effect, repeats RNG, transfers choice ownership or regenerates a domain merely because connectivity changed.

Replay applies persisted events against exact definition/capability versions and may inspect player projections at any sequence. Mechanics replay does not require principal/session/presence records.

## 13. Generic capability boundary

A fighter/card manifest may declare generic capabilities and parameter data. Core dispatches by operation, query, procedure, checkpoint and capability kind, never fighter/card identity.

A new extension handler requires evidence that existing operations cannot represent the behavior, a typed deterministic contract, serializable choices/state, explicit visibility/persistence semantics, registry entry and independent QA approval. A character-named wrapper around identity dispatch is invalid.

## 14. Failure and quarantine behavior

A violated internal invariant is never repaired by guessing. The match is quarantined while preserving the last durable revision, triggering IDs, resolver diagnostics, definition versions and a projection-safe message.

Operational retries may repeat transport or transaction work only through command idempotency. They may not emit a second gameplay result.

## 15. Deterministic fixture authority

`docs/engine/deterministic-fixture-contract.md` is the normative language-neutral fixture contract. It freezes schema/versioning, canonical serialization, definition pins, deterministic IDs, lifecycle identity, idempotency collisions, accepted/rejected batches, random/choice outcomes, state hashes, pending resume, snapshot-tail replay and viewer projections.

Core Runtime and Rules Mechanics must consume the same fixture semantics. A runner may be implemented later, but it may not reinterpret the contract locally.

## 16. Vertical-slice and foundation gate

The foundation is ready for implementation only after related contracts specify deterministic behavior for create/join and seat authority, setup, Maneuver/BOOST, Scheme, Attack/Defense, pending choices, damage/defeat, exhaustion/game end, persistence/reconnect and launch fighter capabilities through generic operations.

Sherwood Forest graph data remains a prerequisite for movement/zone legality tests, not for the language-neutral contract.

This contract can advance from `draft-foundation` when:

- Phase 4C maps launch requirements to canonical capabilities;
- state, command/event, choice, RNG, visibility and persistence documents agree;
- the deterministic fixture contract is frozen;
- ADR 0001 remains accepted;
- independent architecture QA finds no unresolved P0/P1 launch ambiguity;
- initial code proves definitions are data-loaded and shared core contains no fighter/card-ID branches.
