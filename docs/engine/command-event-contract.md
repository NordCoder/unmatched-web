# Command and Event Contract

## Status

```text
status: draft-foundation
parent_issue: #19
correction_issue: #32
architecture: architecture-contract.md
state_model: state-model.md
fixture_contract: deterministic-fixture-contract.md
```

This contract defines the language-neutral transaction boundary between authenticated external principals, the authoritative command handler, durable authority/idempotency records and persisted match history.

## 1. Identity vocabulary

- `principal_id` — authenticated external subject derived from trusted request context;
- `player_instance_id` — opaque match-scoped gameplay identity;
- `session_id` / `client_instance_id` — operational transport metadata;
- `command_id` — client-generated retry identity scoped with the authenticated principal;
- `match_id` — server-allocated match identity.

A client cannot establish authority by sending a `principal_id` or `player_instance_id`. The server authenticates the principal and validates the durable principal-to-player binding.

## 2. Common command request

```yaml
CommandRequest:
  command_id:
  command_schema_version:
  type:
  payload: {}
  match_id: null
  actor_player_id: null
  expected_revision: null
  client_metadata:
    client_instance_id:
    submitted_at:
    trace_id:
```

`principal_id` is not a payload field; it is attached by the trusted authentication boundary.

Client metadata is diagnostic and excluded from gameplay validation, normalized request identity and gameplay ordering.

Lifecycle-specific requirements are:

| Scope | `match_id` | `actor_player_id` | `expected_revision` |
| --- | --- | --- | --- |
| `CreateMatch` | absent | absent | absent |
| `JoinMatch` | required | absent | required |
| Existing-seat commands | required | required | required |

A schema version must define whether absence is represented by omission or explicit `null`. Canonical normalization converts the accepted representation into the same explicit absent marker.

## 3. Lifecycle authorization

### 3.1 CreateMatch

`CreateMatch` establishes a new match and first seat. It requires an authenticated principal but no pre-existing match or player instance.

Validation/transaction order:

1. parse the envelope sufficiently to identify schema version, type and command ID;
2. authenticate the principal;
3. normalize the request and compute the fingerprint;
4. lookup `(principal_id, command_id)` idempotency record;
5. validate creation policy and normalized payload;
6. allocate `match_id` and first `player_instance_id` through deterministic/injected providers;
7. produce `MatchCreated` and first-seat `PlayerJoined` events;
8. atomically persist authority binding, idempotency result, event batch and match head.

A successful result returns the allocated match/player identities. They are result data, not client claims.

### 3.2 JoinMatch

`JoinMatch` authenticates an external principal and assigns a new match-scoped player instance before membership-only validation can apply.

Validation/transaction order:

1. parse and authenticate the principal;
2. normalize target `match_id`, expected revision and payload;
3. lookup idempotency record;
4. load match and validate revision/lifecycle/join policy;
5. reject an incompatible existing principal or seat binding;
6. allocate `player_instance_id` and seat;
7. produce `PlayerJoined` event;
8. atomically persist authority binding, idempotency result, event batch and match head.

The joining client does not provide `actor_player_id`.

### 3.3 Existing-seat commands

For selection, setup, gameplay, interaction and concession commands:

1. parse and authenticate the principal;
2. normalize the request and compute the fingerprint;
3. lookup idempotency record;
4. resolve exact `(principal_id, match_id, actor_player_id)` authority binding;
5. validate expected revision;
6. validate lifecycle and pending-interaction authority;
7. validate legal action/target domain, source, zone, ownership, visibility, costs and command-specific invariants;
8. produce a deterministic accepted batch or deterministic rejection;
9. atomically persist the result.

A client-supplied `actor_player_id` is a claimed runtime reference until step 4 succeeds.

## 4. Canonical normalization and request fingerprint

### 4.1 Fingerprint input

The idempotency store key is:

```text
(principal_id, command_id)
```

Its record contains `fingerprint_schema_version` and a fingerprint over:

```yaml
NormalizedCommandIdentity:
  fingerprint_schema_version:
  principal_id:
  lifecycle_scope: create_match | join_match | existing_seat
  match_id: <value-or-explicit-absent>
  actor_player_id: <value-or-explicit-absent>
  command_schema_version:
  type:
  normalized_payload: {}
  expected_revision_policy: absent | exact
  expected_revision: <value-or-explicit-absent>
```

Excluded fields:

- client timestamps;
- transport/session/client identifiers;
- trace IDs;
- retry count;
- network address;
- authentication token bytes;
- server receipt time.

The initial fingerprint is lowercase `sha256` over canonical UTF-8 JSON bytes as defined in `deterministic-fixture-contract.md`. Fingerprint schema changes require a new explicit version and compatibility reader; existing records are never silently rehashed under different rules.

### 4.2 Payload normalization

Normalization occurs before gameplay validation and must be deterministic:

- reject unknown fields unless the command schema explicitly permits them;
- apply schema-defined defaults before hashing;
- normalize enums and identifiers to their canonical exact representation;
- preserve array order where semantically meaningful;
- canonicalize unordered sets as sorted unique arrays only when the schema declares set semantics;
- do not resolve hidden/current game facts into the payload fingerprint;
- do not consume RNG or mutate state.

Structurally unparseable or unauthenticated requests cannot form a principal-scoped durable idempotency record.

## 5. Idempotency record and duplicate behavior

```yaml
CommandIdempotencyRecord:
  principal_id:
  command_id:
  fingerprint_schema_version:
  request_fingerprint:
  lifecycle_scope:
  match_id: null
  actor_player_id: null
  result_schema_version:
  result:
  committed_at_operational:
```

`committed_at_operational` is audit metadata and is excluded from deterministic gameplay state.

For the same `(principal_id, command_id)`:

### Same fingerprint

Return the exact durable semantic result of the first committed request:

- accepted result returns the same allocated IDs, accepted revision and event sequence range;
- rejected result returns the same rejection code and projection revision where applicable;
- no validation, state mutation, RNG use or second event batch occurs;
- response status may be represented as `duplicate`, but the embedded semantic result is identical.

### Different fingerprint

Return `DUPLICATE_CONFLICT`:

- do not execute the incoming request;
- do not expose the prior request payload or private result;
- do not return an unrelated accepted result as if it belonged to the conflicting request;
- emit no gameplay event;
- keep the original idempotency record unchanged.

### Non-results

Infrastructure/transaction failures before durable commit are not command results and create no idempotency record. Retrying the same request may execute once.

Authentication failure and structurally unparseable schema failure are not durably recorded in the principal command namespace because no trusted normalized principal/request identity exists. They may appear in protected operational audit logs.

After authentication and successful canonical normalization, every terminal accepted or deterministic rejected result is durably recorded. This includes membership, stale revision, lifecycle, legality, ownership, hidden-information and cost rejections. No such rejection is optionally re-evaluated under the same command ID.

## 6. Command result

```yaml
CommandResult:
  result_schema_version:
  command_id:
  delivery_status: first | duplicate
  semantic_status: accepted | rejected
  match_id: null
  actor_player_id: null
  accepted_revision: null
  event_sequence_range: null
  rejection_code: null
  projection_revision: null
  allocated_runtime_ids: {}
```

`delivery_status: duplicate` does not change the semantic result. A conflict uses a rejected semantic result with `DUPLICATE_CONFLICT` and no prior private payload.

## 7. Initial command catalog

### Match lifecycle

```text
CreateMatch
JoinMatch
SelectFighter
SelectBattlefield
ConfirmSetup
Concede
```

Only `CreateMatch` and `JoinMatch` use lifecycle envelopes without an existing actor player instance. All other commands require a validated seat binding.

### Turn/action

```text
StartManeuver
ChooseManeuverBoost
ChooseMovementPath
PlayScheme
DeclareAttack
ChooseDefense
```

### Interaction

```text
SubmitChoice
SubmitCommittedChoice
DeclineOptionalChoice
```

Each gameplay command references runtime instance IDs or revision-scoped legal-action descriptor IDs. It does not send arbitrary rules expressions or computed damage/legality claims from the client.

## 8. Legal-action descriptors

```yaml
LegalActionDescriptor:
  legal_action_id:
  revision:
  actor_player_id:
  type:
  source_instance_ref:
  permission_id:
  required_inputs: []
  target_domain:
  cost_preview:
  ui_metadata:
```

The application layer validates principal-to-player authority before requesting descriptors. `legal_action_id` is scoped to revision and actor. The server regenerates or validates legality on command acceptance.

UI metadata contains only projection-safe labels, icons or definition references. It is not executable behavior and cannot reveal hidden reasons an alternative is illegal.

## 9. Event envelope

```yaml
DomainEvent:
  event_schema_version:
  event_id:
  match_id:
  sequence:
  revision:
  type:
  caused_by_command_id:
  parent_event_id:
  source_ref:
  ruleset_version:
  public_payload: {}
  private_payloads: {}
```

Private payloads are keyed by authorized `player_instance_id` values and stored durably. Principal/session identities are not gameplay payload owners and need not appear in events.

A public replay/export does not concatenate private payloads by default.

## 10. Event batch and lifecycle atomicity

One accepted command produces an ordered batch:

```yaml
EventBatch:
  command_id:
  previous_revision:
  next_revision:
  first_sequence:
  last_sequence:
  events: []
  terminal_projection_hints: []
```

The batch is atomically persisted with the command result and match head. Lifecycle commands additionally persist required authority-record changes in the same application transaction.

Internal deterministic continuation should resolve in the same transaction until external input is required. A separate system continuation identity is allowed only when the match cannot be externally observed between logically atomic steps and the continuation remains replay/idempotency safe.

## 11. Event families

### Match/setup

```text
MatchCreated
PlayerJoined
FighterSelected
BattlefieldSelected
DeckConstructed
DeckShuffled
StartingHandDrawn
SetupPlacementEstablished
MatchStarted
```

`MatchCreated` and `PlayerJoined` contain match-scoped runtime identities and public gameplay facts. They do not expose external principal identity or session data. The corresponding authority bindings are application/security records committed atomically with the event batch.

### Turn/action

```text
TurnStarted
ActionPermissionGranted
ActionStarted
ActionCostPaid
ActionCompleted
TurnEndRequested
TurnEnded
```

### Card/zone

```text
CardMoved
CardsDrawn
CardsDiscarded
CardRevealed
CardVisibilityChanged
ZoneReordered
CardAttached
CardControlChanged
```

### Movement/battlefield

```text
MovementPathCommitted
FighterMovedStep
FighterPlaced
MovementInterrupted
ComponentDeployed
ComponentMoved
ConnectionStateChanged
SpaceLocked
```

`ConnectionStateChanged` means a gameplay battlefield connection/path state change. It never means network/session connectivity.

### Combat

```text
AttackDeclared
DefenseCommitted
NoDefenseChosen
CombatCardsRevealed
CombatParticipantReplaced
CombatCardReplaced
CombatValueChanged
CombatDamageCalculated
CombatEnded
```

### Damage/health

```text
DamageProposed
DamagePrevented
DamageRedirected
DamageApplied
HealthAssigned
HealthRecovered
FighterDefeated
```

### Effect/choice

```text
EffectQueued
EffectStageStarted
InteractionOpened
ChoiceSubmitted
CommittedChoicesRevealed
InteractionClosed
EffectResolved
EffectCancelled
DelayedObligationCreated
DelayedObligationExpired
```

### Random/audit/result

```text
RandomResultEstablished
GameEnded
MatchQuarantined
```

Names are provisional where Phase 4C has not frozen payload semantics. Envelope, identity, replay, visibility and idempotency rules are normative.

## 12. Choice commands

A choice command contains:

```yaml
interaction_instance_id:
choice_id:
selected_values: []
```

The server checks current interaction identity, actor ownership, cardinality/distinctness, authoritative legal domain, required order, visibility and optional-decline permission.

A choice result event persists selected runtime instance IDs and any authorized hidden payload. Reconnect resumes the same interaction ID and cursor; connectivity creates no replacement interaction.

## 13. Movement and attack commands

The client submits an ordered path of stable space IDs for movement. The server validates the complete path and may resolve it stepwise to permit entry-trigger interruptions.

`MOVE` traverses a path; `PLACE` establishes a destination without implied traversal.

`DeclareAttack` references attacker, defender, attack card and legal action ID. The server derives attack type, range policy, target legality, card usability and cost. The client does not submit damage or claim range.

`ChooseDefense` references the current combat interaction and one legal defense card instance or an explicit no-defense choice.

## 14. Rejection model

Stable categories include:

```text
INVALID_SCHEMA
NOT_AUTHENTICATED
NOT_MATCH_MEMBER
PRINCIPAL_ALREADY_BOUND
SEAT_UNAVAILABLE
STALE_REVISION
MATCH_NOT_ACTIVE
NOT_COMMAND_OWNER
NO_PENDING_INTERACTION
ILLEGAL_ACTION
ILLEGAL_TARGET
ILLEGAL_CARD
INSUFFICIENT_COST
HIDDEN_INFORMATION_VIOLATION
DUPLICATE_CONFLICT
MATCH_QUARANTINED
```

Player-facing details must not leak hidden state or prior conflicting request payloads. Protected diagnostic detail belongs in operational logs.

`INVALID_SCHEMA` may be durably recorded only when authentication succeeded and canonical normalization completed sufficiently to establish the exact request fingerprint. Unparseable input is not recorded in the command namespace.

## 15. Event application

```text
apply(GameState, DomainEvent) -> GameState
```

The reducer is pure and deterministic. It cannot query external services, authority/presence registries, consume RNG, depend on wall-clock time, inspect binary assets, branch on fighter/card IDs or emit additional persisted events as an implicit side effect.

Event generation and effect resolution decide which events exist; application only applies recorded results.

## 16. Compatibility and persistence

- command schemas are versioned at the trust boundary;
- fingerprint schemas are versioned independently and retained for old idempotency records;
- command-result schemas are durable and versioned;
- event schemas are versioned for replay;
- definition/capability versions are stored with matches/events;
- migrations preserve semantic replay or retain old readers;
- renaming or reinterpreting a durable event without migration is incompatible.

## 17. Required deterministic evidence

Fixtures under `deterministic-fixture-contract.md` must prove:

- create and join require no pre-existing runtime player identity;
- later commands require exact principal-to-seat binding;
- same command ID and same fingerprint return the same durable result without duplicate events;
- same command ID and different fingerprint return `DUPLICATE_CONFLICT` without execution;
- deterministic rejections are stable under retry;
- stale revision and illegal commands emit no gameplay event;
- event sequences are total and gap-free;
- persisted batches reproduce expected state hashes;
- private payloads never appear in unauthorized projections/errors;
- reconnect resumes the exact pending interaction;
- random results and shuffled order are replayed, not regenerated;
- fighter/card behavior dispatches through generic operations rather than identity checks.
