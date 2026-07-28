# Command and Event Contract

## Status

```text
status: normative-foundation
parent_issue: #19
correction_round_1: #32
correction_round_2: #35
architecture: architecture-contract.md
state_model: state-model.md
fixture_contract: deterministic-fixture-contract.md
```

This contract defines the language-neutral transaction boundary between authenticated principals, command processing, immutable idempotency records, authority records and persisted gameplay events.

## 1. Identity vocabulary

- `principal_id` — authenticated external subject derived from trusted request context;
- `player_instance_id` — opaque match-scoped gameplay identity;
- `session_id` / `client_instance_id` — operational transport metadata;
- `command_id` — client-generated retry identity scoped with the authenticated principal;
- `match_id` — server-allocated match identity.

Clients cannot establish authority by supplying a principal or player identity. Existing-seat commands are authorized only after the server resolves an active durable principal-to-player binding.

## 2. Command envelope

```yaml
CommandRequest:
  command_id:
  command_schema_version:
  type:
  payload: {}
  match_id:        # omitted when forbidden by the command schema
  actor_player_id: # omitted when forbidden by the command schema
  expected_revision:
  client_metadata:
    client_instance_id:
    submitted_at:
    trace_id:
```

`principal_id` is attached by the trusted authentication boundary and is not a payload field. Client metadata is excluded from gameplay validation, request fingerprinting and gameplay ordering.

| Scope | `match_id` | `actor_player_id` | `expected_revision` |
| --- | --- | --- | --- |
| `CreateMatch` | omitted | omitted | omitted |
| `JoinMatch` | required | omitted | required |
| existing-seat commands | required | required | required |

For schema v1, omitted and JSON `null` are distinct. A forbidden or absent optional field is omitted; `null` is accepted only where the command schema explicitly permits it.

## 3. Lifecycle authorization

### 3.1 `CreateMatch`

1. parse the envelope enough to identify schema, type and command ID;
2. authenticate the principal;
3. normalize the request and compute its fingerprint;
4. perform idempotency lookup for `(principal_id, command_id)`;
5. validate creation policy and payload;
6. allocate match/player IDs through deterministic or injected providers;
7. produce `MatchCreated` and first-seat `PlayerJoined`;
8. atomically persist authority binding, immutable command result, event batch and match head.

No pre-existing match or player instance is required.

### 3.2 `JoinMatch`

1. parse and authenticate;
2. normalize target match, expected revision and payload;
3. perform idempotency lookup;
4. load the match and validate revision/lifecycle/join policy;
5. reject an incompatible existing principal or seat binding;
6. allocate the new match-scoped player identity and seat;
7. produce `PlayerJoined`;
8. atomically persist authority binding, immutable command result, event batch and match head.

The joining request does not contain `actor_player_id`.

### 3.3 Existing-seat commands

1. parse and authenticate;
2. normalize and fingerprint;
3. perform idempotency lookup;
4. resolve exact `(principal_id, match_id, actor_player_id)` authority;
5. validate revision, lifecycle and pending-interaction authority;
6. validate legality, ownership, visibility, costs and command-specific invariants;
7. produce either one deterministic rejection or one accepted event batch;
8. persist the result atomically.

A client-supplied `actor_player_id` is only a claimed runtime reference until step 4 succeeds.

## 4. Canonical request identity

The idempotency key is:

```text
(principal_id, command_id)
```

The immutable record stores a versioned fingerprint over:

```yaml
NormalizedCommandIdentity:
  fingerprint_schema_version:
  principal_id:
  lifecycle_scope: create_match | join_match | existing_seat
  match_id: <value-or-explicit-absent-marker>
  actor_player_id: <value-or-explicit-absent-marker>
  command_schema_version:
  type:
  normalized_payload: {}
  expected_revision_policy: absent | exact
  expected_revision: <value-or-explicit-absent-marker>
```

Normalization:

- rejects unknown or duplicate fields;
- applies schema-defined defaults before hashing;
- preserves array order;
- normalizes all strings and object keys to Unicode NFC;
- rejects invalid Unicode and NFC-created duplicate keys;
- permits only schema-declared safe-range integers;
- uses RFC 8785 JCS bytes followed by SHA-256;
- consumes no RNG and mutates no state.

The precise canonicalization contract and vectors are in `deterministic-fixture-contract.md` and `fixtures/foundation-v1.json`.

## 5. Immutable idempotency record

```yaml
CommandIdempotencyRecord:
  principal_id:
  command_id:
  fingerprint_schema_version:
  request_fingerprint:
  lifecycle_scope:
  match_id:
  actor_player_id:
  result_schema_version:
  result:
  committed_at_operational:
```

Exactly one immutable record may exist per `(principal_id, command_id)`.

### 5.1 Same fingerprint

Return the original durable semantic result:

- accepted results return the same allocated IDs, revision and event range;
- rejected results return the same rejection;
- no validation, execution, RNG, ID allocation or new event batch occurs;
- delivery may be marked `duplicate`.

### 5.2 Different fingerprint: derived collision

When the key exists but the presented fingerprint differs:

1. return a derived `DUPLICATE_CONFLICT` response;
2. do not execute or validate the conflicting request beyond collision detection;
3. do not create a second command-result or idempotency record;
4. do not modify the original record;
5. disclose none of the original payload, fingerprint, result or private data;
6. emit no gameplay event and change no revision or sequence;
7. classify repeated conflicting retries identically.

This response is outside the durable-terminal-result rule because the command identity is already occupied by another normalized request. An optional security audit entry is operational only and is not a command result or deterministic evidence.

### 5.3 Results and non-results

After authentication and successful canonical normalization, a new unoccupied command identity records exactly one accepted or deterministic rejected result.

The following create no command-result record:

- authentication failure without trusted principal identity;
- structurally unparseable input without a canonical fingerprint;
- infrastructure or transaction failure before commit;
- `DUPLICATE_CONFLICT`, which is derived from the occupied key.

## 6. Response classes

```yaml
CommandResponse:
  response_class: durable_result | derived_collision
  command_id:
  delivery_status: first | duplicate | collision
  semantic_status: accepted | rejected
  match_id:
  actor_player_id:
  accepted_revision:
  event_sequence_range:
  rejection_code:
  projection_revision:
  allocated_runtime_ids: {}
```

Allowed combinations:

- `durable_result / first` — newly committed accepted or rejected result;
- `durable_result / duplicate` — replay of the immutable original result;
- `derived_collision / collision` — non-persisted `DUPLICATE_CONFLICT`.

A derived collision never embeds the original semantic result.

## 7. Command families

Lifecycle:

```text
CreateMatch
JoinMatch
SelectFighter
SelectBattlefield
ConfirmSetup
Concede
```

Turn/action:

```text
StartManeuver
ChooseManeuverBoost
ChooseMovementPath
PlayScheme
DeclareAttack
ChooseDefense
```

Interaction:

```text
SubmitChoice
SubmitCommittedChoice
DeclineOptionalChoice
```

Gameplay commands reference runtime instance IDs or revision-scoped legal-action IDs, never arbitrary rules expressions or client-computed damage.

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

Descriptors are projection-safe facts. The server regenerates or validates legality on acceptance.

## 9. Event envelope

```yaml
DomainEvent:
  event_schema_version:
  event_id:
  match_id:
  sequence:
  revision:
  event_type:
  caused_by_command_id:
  parent_event_id: # omitted when absent
  source_ref:      # omitted when absent
  ruleset_version:
  public_payload: {}
  private_payloads_by_player: {}
```

Every exact fixture event supplies every mandatory field. Empty payloads are explicit `{}`. Optional fields are omitted rather than set to `null` unless their schema permits `null`.

Private payloads are keyed by authorized `player_instance_id`. Principal/session identities are not gameplay payload owners.

## 10. Event batches

```yaml
EventBatch:
  command_id:
  previous_revision:
  next_revision:
  first_sequence:
  last_sequence:
  events: []
```

One accepted command commits one ordered atomic batch together with its command result and match head. Lifecycle batches also commit authority changes.

Every event is explicit. Reducers cannot emit additional persisted events or infer undocumented state changes as side effects.

## 11. Event families

Setup:

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

Turn/action:

```text
TurnStarted
ActionPermissionGranted
ActionStarted
ActionCostPaid
ActionCompleted
TurnEndRequested
TurnEnded
```

Card/zone:

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

Effect/choice:

```text
EffectQueued
EffectStageStarted
EffectStageChanged
EffectDequeued
InteractionOpened
ChoiceSubmitted
CommittedChoicesRevealed
InteractionClosed
EffectResolved
EffectCancelled
DelayedObligationCreated
DelayedObligationExpired
```

Other existing combat, movement, health, random and game-result families remain permitted under versioned schemas. `ConnectionStateChanged` is reserved for a gameplay battlefield connection/path and never represents network presence.

## 12. Choice and replay sufficiency

A choice command identifies the interaction, choice and ordered selected values. The server validates ownership, cardinality, legal domain and visibility.

Persisted events must contain enough information to replay every state transition. For the normative snapshot-tail fixture this includes:

- selected values in authorized private payload;
- complete card source/destination positions and resulting face state;
- interaction closure;
- effect stage/status change;
- queue removal;
- explicit `ActionCompleted`;
- envelope-driven revision and history cursor progression.

## 13. Rejections

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

Errors must not leak hidden state or prior conflicting request data.

## 14. Event application

```text
apply(GameState, DomainEvent) -> GameState
```

The reducer is pure and deterministic. It cannot query services, consume RNG, depend on wall-clock time, inspect assets, branch on fighter/card IDs or create implicit persisted events.

## 15. Compatibility and required evidence

Command, fingerprint, result and event schemas are independently versioned. Migrations preserve semantic replay or old readers remain available.

Normative artifacts:

```text
docs/engine/fixtures/schema-v1.json
docs/engine/fixtures/foundation-v1.json
docs/engine/fixtures/foundation-v1-transition-audit.json
```

They prove:

- lifecycle creation without pre-existing runtime identity;
- same-key/same-fingerprint replay;
- same-key/different-fingerprint derived collision with one immutable record;
- duplicate-key rejection and NFC/JCS vectors;
- isolated provider resets;
- full event envelopes;
- snapshot plus contiguous tail equality;
- explicit action completion and projection equality.
