# Deterministic Engine Fixture Contract

## Status

```text
status: normative-foundation
schema_id: unmatched.engine.fixture/v1
parent_issue: #19
correction_issue: #32
```

This document defines the language-neutral fixture contract shared by Core Runtime, Rules Mechanics and Lead integration. It is normative for deterministic command, event, replay, reconnect and projection evidence. The examples use synthetic definitions and do not bypass Phase 4C or battlefield gates.

## 1. Required suite shape

A fixture is UTF-8 YAML or JSON with one top-level suite:

```yaml
fixture_schema: unmatched.engine.fixture/v1
fixture_id:
purpose:
definition_bundle: {}
id_providers: {}
authority_records: []
initial_source: {}
random_inputs: []
cases: []
```

A runner must reject an unknown `fixture_schema`. `unmatched.engine.fixture/v1` is immutable after acceptance; semantic changes require a new version and an old reader or verified migration.

## 2. Canonical serialization and hashes

Fingerprint, state and projection hashes use SHA-256 over canonical JSON bytes:

1. UTF-8 without BOM;
2. Unicode NFC;
3. object keys sorted by Unicode code-point order;
4. no insignificant whitespace;
5. arrays preserve declared order;
6. integers use minimal base-10 form;
7. floating-point values are forbidden;
8. unordered semantic sets are sorted and deduplicated before serialization;
9. absent and `null` remain distinct when the schema permits both.

Hash syntax is `sha256:<64 lowercase hex digits>`. YAML source text is never hashed directly.

## 3. Definitions and deterministic providers

Every referenced definition is pinned by a digest. Inline synthetic definitions are allowed when their canonical JSON and digest are both declared.

```yaml
definition_bundle:
  ruleset_version: rules-v1
  capability_registry_version: cap-v1
  inline_definitions:
    battlefield:
      canonical_json: '{"definition_id":"synthetic-field","spaces":["s1","s2"],"zones":{"z1":["s1","s2"]}}'
      digest: sha256:5314539850defcd6d40cd8d82b9543fb5fed27f61f263cc14a0baa01bd4bf4f5
    fighter_a:
      canonical_json: '{"definition_id":"test-a","move":2,"sidekick_count":0}'
      digest: sha256:6761035c5a248a131289908fa42fb4de4e359859aa21ad99ea95552a724aa634
    fighter_b:
      canonical_json: '{"definition_id":"test-b","move":2,"sidekick_count":0}'
      digest: sha256:38c48554eb82623c7c98ee27aa7fd777352560c5a0b620601c08c80a585dadfb
    cards:
      canonical_json: '{"definitions":[{"boost":1,"definition_id":"test-card","type":"scheme","value":null}]}'
      digest: sha256:61477f9138f58cf2044ebca7abc0ee16709dba1fcba7ead7372b270e3bee7985
    setup:
      canonical_json: '{"definition_id":"synthetic-setup-v1","seat_count":2}'
      digest: sha256:4053365a5acb7d126c1d688bdc4388d687df6b0296116774f763ef78936bfc68
```

Binary asset paths and image availability are excluded from definition and gameplay hashes.

Typed IDs are consumed in declared order:

```yaml
id_providers:
  match_ids: [match-0001]
  player_ids: [player-0001, player-0002]
  card_ids: [card-0001, card-secret-a, card-secret-b]
  event_ids: [event-0001, event-0002, event-0003, event-0013, event-0014, event-0015, event-0016, event-0017]
  action_ids: [action-0001]
  effect_ids: [effect-0001]
  interaction_ids: [interaction-0001]
```

A runner fails on provider underflow, unexpected allocation or cross-type reuse. It must not replace declared IDs with random UUIDs.

## 4. Authority, presence and initial source

Durable authority is separate from deterministic `GameState`:

```yaml
authority_records:
  - match_id: match-0001
    player_instance_id: player-0001
    principal_id: principal-alice
    seat: 1
    binding_version: 1
    status: ACTIVE
```

Connectivity is operational input:

```yaml
presence_change:
  principal_id:
  player_instance_id:
  online:
  client_instance_ids: []
```

Presence never enters reducers, snapshots, replay, state hashes or Mechanics views. A deterministic projection hash excludes presence; an optional `PresenceProjection` is asserted separately.

`initial_source.kind` is exactly one of:

```text
empty_match_host
canonical_state
snapshot_plus_tail
```

For canonical state or snapshots, the fixture contains the complete canonical JSON and its hash.

## 5. Command, result and event assertions

A command step declares trusted authentication context separately from the client request:

```yaml
step: command
authenticated_principal_id:
request:
  command_id:
  command_schema_version:
  type:
  match_id: null
  actor_player_id: null
  expected_revision: null
  payload: {}
expected_normalized_identity_json:
expected_request_fingerprint:
expect:
  delivery_status: first | duplicate
  semantic_status: accepted | rejected
  rejection_code: null
  accepted_revision: null
  event_sequence_range: null
  allocated_runtime_ids: {}
  authority_changes: []
  events: []
  state_checkpoint: null
```

The authenticated principal is not a payload field. A same-fingerprint duplicate returns the exact durable semantic result without validator, reducer, RNG or ID-provider consumption. Reusing the same `(principal_id, command_id)` with a different fingerprint returns a first-delivery rejected result with `DUPLICATE_CONFLICT`; the original record is unchanged and no prior private result is exposed.

Expected events use exact mode:

```yaml
event_schema_version:
event_id:
match_id:
sequence:
revision:
type:
caused_by_command_id:
parent_event_id: null
source_ref: null
ruleset_version:
public_payload: {}
private_payloads: {}
```

Rejected commands and operational presence changes emit no gameplay events.

## 6. State, interaction, projection and replay assertions

```yaml
state_checkpoint:
  revision:
  event_sequence:
  state_schema_version:
  canonical_state_json:
  state_hash:
  pending_interaction: null
```

A symbolic or placeholder hash is invalid. Connectivity/session data must not occur in `canonical_state_json`.

```yaml
projection_checkpoint:
  viewer:
    match_id:
    player_instance_id:
    authority_version:
  projection_schema_version:
  canonical_projection_json:
  projection_hash:
  forbidden_json_paths: []
  forbidden_values: []
```

Reconnect asserts the bound player, unchanged state hash and unchanged pending interaction, with zero revision, event and RNG deltas.

Snapshot-tail replay applies contiguous events against the pinned definitions. It fails on sequence gaps, unknown schemas, hash mismatch, missing private payload ownership or unserializable pending state.

## 7. Normative correction suite

### 7.1 CreateMatch without pre-existing runtime identity

Trusted principal: `principal-alice`.

Request:

```yaml
command_id: command-create-0001
command_schema_version: unmatched.command/v1
type: CreateMatch
payload:
  ruleset_version: rules-v1
```

`match_id`, `actor_player_id` and `expected_revision` are absent.

Canonical normalized identity:

```json
{"actor_player_id":{"absent":true},"command_schema_version":"unmatched.command/v1","expected_revision":{"absent":true},"expected_revision_policy":"absent","fingerprint_schema_version":"unmatched.command-fingerprint/v1","lifecycle_scope":"create_match","match_id":{"absent":true},"normalized_payload":{"ruleset_version":"rules-v1"},"principal_id":"principal-alice","type":"CreateMatch"}
```

```text
request_fingerprint: sha256:f24699d76aa7d2bc54ad4b7634ee385fbc06c07b027f650a9e560a555d184db5
```

Expected durable result:

```yaml
delivery_status: first
semantic_status: accepted
match_id: match-0001
actor_player_id: player-0001
accepted_revision: 1
event_sequence_range: [1, 2]
allocated_runtime_ids:
  match_id: match-0001
  player_instance_id: player-0001
authority_changes:
  - match_id: match-0001
    player_instance_id: player-0001
    principal_id: principal-alice
    seat: 1
    binding_version: 1
    status: ACTIVE
events:
  - event_id: event-0001
    sequence: 1
    revision: 1
    type: MatchCreated
    caused_by_command_id: command-create-0001
  - event_id: event-0002
    sequence: 2
    revision: 1
    type: PlayerJoined
    caused_by_command_id: command-create-0001
    public_payload:
      player_instance_id: player-0001
      seat: 1
```

Canonical state checkpoint:

```json
{"action":null,"battlefield":null,"cards":{},"combat":null,"components":{},"definition_ref":{"capability_registry_version":"cap-v1","ruleset_version":"rules-v1"},"event_sequence":2,"fighters":{},"game_result":null,"lifecycle":"WAITING_FOR_PLAYERS","match_id":"match-0001","players":{"player-0001":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":1,"submitted_hidden_choice_ids":[]}},"random":{"algorithm_version":"test-v1","generator_state_ref":null,"last_random_result_sequence":null},"resolver":{"checkpoint_stack":[],"delayed_obligations":[],"effect_instances":{},"pending_interaction":null,"queue":[],"simultaneous_choice_groups":{}},"revision":1,"turn":null}
```

```text
state_hash: sha256:46f4ce8b506587a242191f60e14f70401e1efda096ae77763910f34c5c6a8ded
```

### 7.2 Same ID duplicate and different-request conflict

Retry the exact request above under `principal-alice`.

```yaml
delivery_status: duplicate
semantic_status: accepted
accepted_revision: 1
event_sequence_range: [1, 2]
gameplay_events_delta: 0
id_allocations_delta: 0
rng_calls_delta: 0
```

Then reuse `command-create-0001` with `ruleset_version: rules-v2`.

Canonical conflicting identity:

```json
{"actor_player_id":{"absent":true},"command_schema_version":"unmatched.command/v1","expected_revision":{"absent":true},"expected_revision_policy":"absent","fingerprint_schema_version":"unmatched.command-fingerprint/v1","lifecycle_scope":"create_match","match_id":{"absent":true},"normalized_payload":{"ruleset_version":"rules-v2"},"principal_id":"principal-alice","type":"CreateMatch"}
```

```text
request_fingerprint: sha256:8efb8648eb9cd756b8999f19dc695a77b4d3e015dcc98a54e58c4cb2b60511b9
```

Expected:

```yaml
delivery_status: first
semantic_status: rejected
rejection_code: DUPLICATE_CONFLICT
gameplay_events_delta: 0
state_hash_unchanged: sha256:46f4ce8b506587a242191f60e14f70401e1efda096ae77763910f34c5c6a8ded
prior_result_payload_exposed: false
```

### 7.3 JoinMatch without pre-existing runtime player identity

Trusted principal: `principal-bob`.

```yaml
command_id: command-join-0001
command_schema_version: unmatched.command/v1
type: JoinMatch
match_id: match-0001
expected_revision: 1
payload: {}
```

`actor_player_id` is absent.

Canonical normalized identity:

```json
{"actor_player_id":{"absent":true},"command_schema_version":"unmatched.command/v1","expected_revision":1,"expected_revision_policy":"exact","fingerprint_schema_version":"unmatched.command-fingerprint/v1","lifecycle_scope":"join_match","match_id":"match-0001","normalized_payload":{},"principal_id":"principal-bob","type":"JoinMatch"}
```

```text
request_fingerprint: sha256:04befa0a7045ed4b101f6587d4a90b69c0e1c5e963b49dafeaa60116c2bafd29
```

Expected:

```yaml
delivery_status: first
semantic_status: accepted
actor_player_id: player-0002
accepted_revision: 2
event_sequence_range: [3, 3]
authority_changes:
  - match_id: match-0001
    player_instance_id: player-0002
    principal_id: principal-bob
    seat: 2
    binding_version: 1
    status: ACTIVE
events:
  - event_id: event-0003
    sequence: 3
    revision: 2
    type: PlayerJoined
    caused_by_command_id: command-join-0001
    public_payload:
      player_instance_id: player-0002
      seat: 2
```

Canonical state checkpoint:

```json
{"action":null,"battlefield":null,"cards":{},"combat":null,"components":{},"definition_ref":{"capability_registry_version":"cap-v1","ruleset_version":"rules-v1"},"event_sequence":3,"fighters":{},"game_result":null,"lifecycle":"SELECTION","match_id":"match-0001","players":{"player-0001":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":1,"submitted_hidden_choice_ids":[]},"player-0002":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":2,"submitted_hidden_choice_ids":[]}},"random":{"algorithm_version":"test-v1","generator_state_ref":null,"last_random_result_sequence":null},"resolver":{"checkpoint_stack":[],"delayed_obligations":[],"effect_instances":{},"pending_interaction":null,"queue":[],"simultaneous_choice_groups":{}},"revision":2,"turn":null}
```

```text
state_hash: sha256:07d3ef287c7a0a0b5669c7c27fcad60c2bfe735e1da0172702e8389deac649bf
```

### 7.4 Disconnect and reconnect preserve deterministic state

Approved snapshot/checkpoint:

```json
{"action":{"action_instance_id":"action-0001","action_type":"SCHEME","actor_player_id":"player-0001","stage":"AWAITING_CHOICE"},"battlefield":null,"cards":{"card-0001":{"current_controller_player_id":"player-0001","definition_id":"test-card","face_state":"HIDDEN","immutable_owner_player_id":"player-0001","zone_position":0,"zone_ref":{"owner_scope_ref":"player-0001","subzone_id":null,"zone_type":"deck"}},"card-secret-a":{"current_controller_player_id":"player-0001","definition_id":"test-card","face_state":"HIDDEN","immutable_owner_player_id":"player-0001","zone_position":0,"zone_ref":{"owner_scope_ref":"player-0001","subzone_id":null,"zone_type":"hand"}},"card-secret-b":{"current_controller_player_id":"player-0001","definition_id":"test-card","face_state":"HIDDEN","immutable_owner_player_id":"player-0001","zone_position":1,"zone_ref":{"owner_scope_ref":"player-0001","subzone_id":null,"zone_type":"hand"}}},"combat":null,"components":{},"definition_ref":{"capability_registry_version":"cap-v1","ruleset_version":"rules-v1"},"event_sequence":12,"fighters":{},"game_result":null,"lifecycle":"ACTIVE","match_id":"match-0001","players":{"player-0001":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":1,"submitted_hidden_choice_ids":[]},"player-0002":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":2,"submitted_hidden_choice_ids":[]}},"random":{"algorithm_version":"test-v1","generator_state_ref":null,"last_random_result_sequence":null},"resolver":{"checkpoint_stack":[],"delayed_obligations":[],"effect_instances":{"effect-0001":{"current_stage":"choose-one","status":"PAUSED"}},"pending_interaction":{"choice_schema":{"cardinality":{"max":1,"min":1},"optional":false},"created_revision":7,"interaction_instance_id":"interaction-0001","legal_domain":["card-secret-a","card-secret-b"],"owner_scope":"player-0001","resume_cursor":{"effect_instance_id":"effect-0001","stage":"choose-one"},"source_effect_id":"effect-0001","submitted_values":{},"visibility":"OWNER_PRIVATE"},"queue":["effect-0001"],"simultaneous_choice_groups":{}},"revision":7,"turn":{"action_ordinal":1,"active_player_id":"player-0001","turn_number":1}}
```

```text
state_hash: sha256:6fb9a4159186fb8060e0efb7bbf8c15c34fb0b8d83f452288dd6d1051fc81ea3
```

Operational sequence:

```yaml
- presence_change:
    principal_id: principal-alice
    player_instance_id: player-0001
    online: false
    client_instance_ids: []
- reconnect:
    authenticated_principal_id: principal-alice
    match_id: match-0001
    last_seen_revision: 7
    last_seen_event_sequence: 12
- presence_change:
    principal_id: principal-alice
    player_instance_id: player-0001
    online: true
    client_instance_ids: [client-a2]
```

Required assertions:

```yaml
bound_player_instance_id: player-0001
revision_delta: 0
gameplay_events_delta: 0
rng_calls_delta: 0
state_hash_before: sha256:6fb9a4159186fb8060e0efb7bbf8c15c34fb0b8d83f452288dd6d1051fc81ea3
state_hash_after: sha256:6fb9a4159186fb8060e0efb7bbf8c15c34fb0b8d83f452288dd6d1051fc81ea3
pending_interaction_id_before: interaction-0001
pending_interaction_id_after: interaction-0001
resume_cursor_equal: true
legal_domain_equal: true
```

Player A projection:

```json
{"event_sequence":12,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0002":{"hand_count":0}},"pending_interaction":{"cardinality":{"max":1,"min":1},"interaction_instance_id":"interaction-0001","legal_options":["card-secret-a","card-secret-b"],"optional":false,"owner":"player-0001","prompt":"Select one card","submitted_state":{}},"public_history":[],"public_state":{"action_type":"SCHEME","active_player_id":"player-0001"},"revision":7,"self_private_state":{"hand":["card-secret-a","card-secret-b"]},"viewer_player_id":"player-0001"}
```

```text
projection_hash: sha256:09c8b4cf6368bf1248cae72b33a1ad31d63a636453faa31cfc923fec2d0d1dae
```

Player B projection:

```json
{"event_sequence":12,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0001":{"hand_count":2}},"pending_interaction":{"interaction_instance_id":"interaction-0001","owner":"player-0001","waiting":true},"public_history":[],"public_state":{"action_type":"SCHEME","active_player_id":"player-0001"},"revision":7,"self_private_state":{"hand":[]},"viewer_player_id":"player-0002"}
```

```text
projection_hash: sha256:8e26321ca209975ac5d1925d632715b6d4a9933e20da50a281cb5681cea91a9f
forbidden_values: [card-secret-a, card-secret-b]
forbidden_json_paths:
  - $.pending_interaction.legal_options
```

### 7.5 Snapshot plus ordered event tail

Start from the revision 7 / sequence 12 snapshot above. Apply one accepted command batch at revision 8:

```yaml
- event_id: event-0013
  sequence: 13
  revision: 8
  type: ChoiceSubmitted
  caused_by_command_id: command-choice-0001
  private_payloads:
    player-0001:
      interaction_instance_id: interaction-0001
      selected_values: [card-secret-a]
- event_id: event-0014
  sequence: 14
  revision: 8
  type: CardMoved
  caused_by_command_id: command-choice-0001
  public_payload:
    card_instance_id: card-secret-a
    from_zone: hand
    to_zone: discard
- event_id: event-0015
  sequence: 15
  revision: 8
  type: CardMoved
  caused_by_command_id: command-choice-0001
  public_payload:
    card_instance_id: card-0001
    from_zone: deck
    to_zone: discard
- event_id: event-0016
  sequence: 16
  revision: 8
  type: InteractionClosed
  caused_by_command_id: command-choice-0001
  public_payload:
    interaction_instance_id: interaction-0001
- event_id: event-0017
  sequence: 17
  revision: 8
  type: EffectResolved
  caused_by_command_id: command-choice-0001
  source_ref: effect-0001
```

Expected final canonical state:

```json
{"action":null,"battlefield":null,"cards":{"card-0001":{"current_controller_player_id":"player-0001","definition_id":"test-card","face_state":"REVEALED","immutable_owner_player_id":"player-0001","zone_position":1,"zone_ref":{"owner_scope_ref":"player-0001","subzone_id":null,"zone_type":"discard"}},"card-secret-a":{"current_controller_player_id":"player-0001","definition_id":"test-card","face_state":"REVEALED","immutable_owner_player_id":"player-0001","zone_position":0,"zone_ref":{"owner_scope_ref":"player-0001","subzone_id":null,"zone_type":"discard"}},"card-secret-b":{"current_controller_player_id":"player-0001","definition_id":"test-card","face_state":"HIDDEN","immutable_owner_player_id":"player-0001","zone_position":0,"zone_ref":{"owner_scope_ref":"player-0001","subzone_id":null,"zone_type":"hand"}}},"combat":null,"components":{},"definition_ref":{"capability_registry_version":"cap-v1","ruleset_version":"rules-v1"},"event_sequence":17,"fighters":{},"game_result":null,"lifecycle":"ACTIVE","match_id":"match-0001","players":{"player-0001":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":1,"submitted_hidden_choice_ids":[]},"player-0002":{"action_permission_ids":[],"authority_state":"ACTIVE","fighter_instance_ids":[],"private_zones":[],"resources":{},"seat":2,"submitted_hidden_choice_ids":[]}},"random":{"algorithm_version":"test-v1","generator_state_ref":null,"last_random_result_sequence":null},"resolver":{"checkpoint_stack":[],"delayed_obligations":[],"effect_instances":{"effect-0001":{"current_stage":"complete","status":"RESOLVED"}},"pending_interaction":null,"queue":[],"simultaneous_choice_groups":{}},"revision":8,"turn":{"action_ordinal":1,"active_player_id":"player-0001","turn_number":1}}
```

```text
state_hash: sha256:02b0c0af497e7d4acdd112fcf8878a1d155adea9799c04f912cfae0c600da42f
rng_calls_delta: 0
```

Player A final projection:

```json
{"event_sequence":17,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0002":{"hand_count":0}},"pending_interaction":null,"public_history":[{"card_instance_id":"card-secret-a","event_type":"CardMoved","to_zone":"discard"}],"public_state":{"action_type":null,"active_player_id":"player-0001","discard":["card-secret-a","card-0001"]},"revision":8,"self_private_state":{"hand":["card-secret-b"]},"viewer_player_id":"player-0001"}
```

```text
projection_hash: sha256:fff64892162fbbe7ba41f877250a9876cf07db9ff423d459d2ce9795f1c8b3d3
```

Player B final projection:

```json
{"event_sequence":17,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0001":{"hand_count":1}},"pending_interaction":null,"public_history":[{"card_instance_id":"card-secret-a","event_type":"CardMoved","to_zone":"discard"}],"public_state":{"action_type":null,"active_player_id":"player-0001","discard":["card-secret-a","card-0001"]},"revision":8,"self_private_state":{"hand":[]},"viewer_player_id":"player-0002"}
```

```text
projection_hash: sha256:f9b090fe791678079bc70a660265718c5132b8d3b93fcd510bd7aefe9d8c0dea
forbidden_values: [card-secret-b]
```

Uninterrupted execution from the same revision 7 state and snapshot-plus-tail replay must produce these exact final state and projection hashes.

## 8. Runner conformance

A conforming runner fails on:

- unknown schema/version;
- definition digest mismatch;
- unexpected ID or RNG consumption;
- command fingerprint mismatch;
- duplicate re-execution;
- event sequence gaps or reordering;
- result/event/state/projection mismatch;
- authority binding mismatch;
- presence data in deterministic state;
- changed pending interaction on reconnect;
- private data leakage;
- snapshot-tail divergence.

Until an executable runner exists, validation is static: canonical examples, recomputed hashes and cross-document consistency. No handoff may claim executable fixture validation without Candidate-bound runner evidence.
