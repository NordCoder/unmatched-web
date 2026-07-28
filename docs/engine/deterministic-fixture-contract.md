# Deterministic Engine Fixture Contract

## Status

```text
status: normative-foundation
schema_id: unmatched.engine.fixture/v1
parent_issue: #19
correction_issue: #32
architecture: architecture-contract.md
state_model: state-model.md
command_contract: command-event-contract.md
persistence_contract: persistence-reconnect.md
visibility_contract: visibility-projections.md
```

This document defines the language-neutral fixture format shared by Core Runtime, Rules Mechanics and Lead integration. A runner may be implemented in Go or another tool, but it must consume these semantics without creating a worker-local variant.

The examples are contract evidence, not final Robin Hood/Bigfoot or Sherwood Forest gameplay fixtures. They deliberately use synthetic definitions so they do not bypass Phase 4C or battlefield gates.

## 1. Goals

A fixture must be sufficient to verify:

- pinned definitions and capability versions;
- deterministic runtime IDs;
- lifecycle principal-to-player binding;
- normalized command fingerprints and command-ID reuse;
- accepted event batches and rejected results;
- persisted random/choice outcomes;
- revision, sequence and state-hash checkpoints;
- pending interaction serialization/resume;
- disconnect/reconnect invariance;
- snapshot plus ordered event-tail replay;
- viewer-specific projections and hidden-information assertions.

Fixtures never depend on wall-clock time, process-global randomness, network state or hidden local files.

## 2. File and suite shape

Normative authoring format is UTF-8 YAML or JSON. Every fixture file contains one suite:

```yaml
fixture_schema: unmatched.engine.fixture/v1
fixture_id: stable-string
purpose: stable-string
definition_bundle: {}
id_providers: {}
authority_records: []
initial_source: {}
random_inputs: []
cases: []
canonical_checkpoints: {}
```

A runner must reject an unknown `fixture_schema`. Additive optional fields require a schema revision when their interpretation can affect execution or expected evidence.

## 3. Canonical serialization

All fingerprints and deterministic hashes use canonical JSON bytes with these rules:

1. UTF-8 without BOM;
2. Unicode strings normalized to NFC;
3. object keys sorted by Unicode code-point order;
4. no insignificant whitespace;
5. arrays preserve declared order;
6. integers use base-10 minimal representation;
7. floating-point values are forbidden in deterministic engine fixtures;
8. booleans are lowercase JSON `true`/`false`;
9. absent and `null` are distinct when the schema declares both;
10. unordered semantic sets are represented as sorted unique arrays before serialization;
11. map keys that are runtime IDs are serialized as ordinary sorted strings;
12. line endings in authoring files do not affect canonical bytes.

Hash syntax:

```text
sha256:<64 lowercase hexadecimal digits>
```

The hash is SHA-256 over the canonical JSON bytes of the declared object, not over YAML source text.

## 4. Evolution rules

- `unmatched.engine.fixture/v1` is immutable once accepted.
- A semantic interpretation change requires `/v2` rather than silently changing a runner.
- A runner may accept unknown authoring-only metadata only when it is explicitly excluded from execution and hashes.
- Durable command/event/snapshot/fingerprint schema versions remain independently declared inside fixtures.
- A fixture migration must preserve old expected semantics or retain an old reader.

## 5. Definition bundle

```yaml
definition_bundle:
  ruleset_version:
  capability_registry_version:
  fighter_manifests:
    fighter-id: sha256:...
  card_manifests:
    corpus-or-fighter-id: sha256:...
  battlefield_manifest: sha256:...
  setup_definition_ref:
  deck_construction_result_ref:
```

Every definition used by state, commands or events must be covered by a pinned digest or an inline synthetic definition block whose canonical digest is declared.

Binary asset paths and image availability are excluded from fixture identity and gameplay hashes.

## 6. Deterministic ID providers

```yaml
id_providers:
  match_ids: [match-0001]
  player_ids: [player-0001, player-0002]
  fighter_ids: []
  card_ids: [card-0001, card-secret-a, card-secret-b]
  event_ids: [event-0001, event-0002]
  effect_ids: [effect-0001]
  interaction_ids: [interaction-0001]
  action_ids: [action-0001]
```

IDs are consumed in declared order by their typed allocation sites. A fixture runner must fail on underflow, unexpected allocation or cross-type reuse. It must not replace fixture IDs with random UUIDs.

## 7. Authority and operational presence

Durable authority setup is declared separately from deterministic state:

```yaml
authority_records:
  - match_id: match-0001
    player_instance_id: player-0001
    principal_id: principal-alice
    seat: 1
    binding_version: 1
    status: ACTIVE
```

Operational presence steps use:

```yaml
presence_change:
  principal_id:
  player_instance_id:
  online:
  client_instance_ids: []
```

Presence changes are never reduced into `GameState`. Unless a case explicitly tests a separate `PresenceProjection`, they are excluded from deterministic state/projection hashes.

## 8. Initial source

A fixture chooses exactly one:

```yaml
initial_source:
  kind: empty_match_host
```

or:

```yaml
initial_source:
  kind: canonical_state
  state_schema_version:
  canonical_state_json: '{}'
  state_hash: sha256:...
```

or:

```yaml
initial_source:
  kind: snapshot_plus_tail
  snapshot: {}
  tail_events: []
```

For `canonical_state` and snapshot payloads, `canonical_state_json` is the complete canonical JSON object used for hashing. YAML convenience views may be present but are non-authoritative if they differ.

## 9. Command step

```yaml
- step: command
  authenticated_principal_id:
  request:
    command_id:
    command_schema_version:
    type:
    match_id: null
    actor_player_id: null
    expected_revision: null
    payload: {}
  expected_normalized_identity_json: '{}'
  expected_request_fingerprint: sha256:...
  expect:
    delivery_status: first | duplicate
    semantic_status: accepted | rejected
    rejection_code: null
    accepted_revision: null
    event_sequence_range: null
    allocated_runtime_ids: {}
    events: []
    authority_changes: []
    state_checkpoint: null
```

The authenticated principal is test harness context, not a command payload field.

For a duplicate step, the runner must prove no validator, reducer, RNG provider or ID provider was consumed again. An instrumentation assertion may be used where available.

## 10. Event expectations

Each expected event may assert:

```yaml
- event_schema_version:
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

Omitted event fields are not wildcards unless the fixture explicitly declares an assertion mode. Normative correction fixtures use exact mode.

Rejected commands, duplicate conflicts and operational presence steps expect no gameplay event batch.

## 11. Random and choice inputs

A fixture may provide deterministic random results:

```yaml
random_inputs:
  - request_ref: random-request-0001
    result_event:
      type: RandomResultEstablished
      public_payload: {}
      private_payloads: {}
```

The runner must consume the declared persisted outcome. Replay must use the event and must not call the random provider again.

Choice commands identify the existing `interaction_instance_id` and selected values. Expected events must persist the chosen runtime IDs or authorized private payloads.

## 12. State checkpoints

```yaml
state_checkpoint:
  revision:
  event_sequence:
  state_schema_version:
  canonical_state_json: '{}'
  state_hash: sha256:...
  pending_interaction:
    interaction_instance_id:
    resume_cursor: {}
```

A runner computes and compares the hash. A symbolic or placeholder hash is invalid.

Connectivity/session data must not appear in `canonical_state_json`.

## 13. Projection checkpoints

```yaml
projection_checkpoint:
  viewer:
    match_id:
    player_instance_id:
    authority_version:
  projection_schema_version:
  canonical_projection_json: '{}'
  projection_hash: sha256:...
  forbidden_json_paths: []
```

The projection hash covers deterministic `PlayerProjection` only. `PresenceProjection` is asserted separately when required.

Forbidden paths are evaluated after serialization. A runner fails if a forbidden path or value appears anywhere in the viewer output, error detail or public history covered by the case.

## 14. Reconnect step

```yaml
- step: reconnect
  authenticated_principal_id:
  match_id:
  last_seen_revision:
  last_seen_event_sequence:
  expect:
    bound_player_instance_id:
    state_hash_unchanged_from:
    projection_checkpoint:
    pending_interaction_unchanged_from:
    rng_calls_delta: 0
    gameplay_events_delta: 0
    revision_delta: 0
```

A preceding `presence_change` may set the player offline and online. The deterministic expectations remain identical.

## 15. Snapshot-tail replay step

```yaml
- step: replay
  snapshot:
    revision:
    event_sequence:
    canonical_state_json: '{}'
    state_hash: sha256:...
  tail_events: []
  expect:
    final_state_checkpoint:
    projections: []
    rng_calls_delta: 0
```

Replay applies tail events in contiguous sequence against exact pinned definitions. It fails on a gap, unknown event schema, hash mismatch or unserializable pending state.

## 16. Complete correction example suite

The following suite is normative evidence for the four Architecture QA blockers. It uses synthetic definitions and four connected cases.

```yaml
fixture_schema: unmatched.engine.fixture/v1
fixture_id: foundation-correction-v1
purpose: lifecycle-idempotency-presence-reconnect-replay-projection
definition_bundle:
  ruleset_version: rules-v1
  capability_registry_version: cap-v1
  fighter_manifests:
    test-a: sha256:3333333333333333333333333333333333333333333333333333333333333333
    test-b: sha256:4444444444444444444444444444444444444444444444444444444444444444
  card_manifests:
    test: sha256:2222222222222222222222222222222222222222222222222222222222222222
  battlefield_manifest: sha256:1111111111111111111111111111111111111111111111111111111111111111
  setup_definition_ref: synthetic-setup-v1
  deck_construction_result_ref: deck-result-1
id_providers:
  match_ids: [match-0001]
  player_ids: [player-0001, player-0002]
  event_ids: [event-0001, event-0002, event-0003]
initial_source:
  kind: empty_match_host
cases:
  - case_id: create-match-without-runtime-player
  - case_id: duplicate-and-conflict
  - case_id: join-match-without-runtime-player
  - case_id: disconnect-reconnect-and-snapshot-tail
```

### 16.1 CreateMatch without pre-existing match/player identity

Authenticated context:

```yaml
principal_id: principal-alice
```

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

Expected fingerprint:

```text
sha256:f24699d76aa7d2bc54ad4b7634ee385fbc06c07b027f650a9e560a555d184db5
```

Expected accepted result:

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

Expected state hash after applying the batch:

```text
sha256:4ba1fb4009d37abdb79af9cf99c4dccc41c01ae8e681c7e8ec2266b6847839fd
```

### 16.2 Same ID/same request and same ID/different request

Retry the exact request above under the same authenticated principal.

Expected:

```yaml
delivery_status: duplicate
semantic_status: accepted
accepted_revision: 1
event_sequence_range: [1, 2]
allocated_runtime_ids:
  match_id: match-0001
  player_instance_id: player-0001
gameplay_events_delta: 0
id_allocations_delta: 0
rng_calls_delta: 0
```

Then reuse `command-create-0001` with payload `ruleset_version: rules-v2`.

Canonical conflicting identity:

```json
{"actor_player_id":{"absent":true},"command_schema_version":"unmatched.command/v1","expected_revision":{"absent":true},"expected_revision_policy":"absent","fingerprint_schema_version":"unmatched.command-fingerprint/v1","lifecycle_scope":"create_match","match_id":{"absent":true},"normalized_payload":{"ruleset_version":"rules-v2"},"principal_id":"principal-alice","type":"CreateMatch"}
```

Conflicting fingerprint:

```text
sha256:8efb8648eb9cd756b8999f19dc695a77b4d3e015dcc98a54e58c4cb2b60511b9
```

Expected:

```yaml
delivery_status: first
semantic_status: rejected
rejection_code: DUPLICATE_CONFLICT
gameplay_events_delta: 0
state_hash_unchanged: sha256:4ba1fb4009d37abdb79af9cf99c4dccc41c01ae8e681c7e8ec2266b6847839fd
prior_result_payload_exposed: false
```

### 16.3 JoinMatch without pre-existing runtime player identity

Authenticated context:

```yaml
principal_id: principal-bob
```

Request:

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

Expected fingerprint:

```text
sha256:04befa0a7045ed4b101f6587d4a90b69c0e1c5e963b49dafeaa60116c2bafd29
```

Expected result:

```yaml
semantic_status: accepted
actor_player_id: player-0002
accepted_revision: 2
event_sequence_range: [3, 3]
allocated_runtime_ids:
  player_instance_id: player-0002
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
    public_payload:
      player_instance_id: player-0002
      seat: 2
state_hash: sha256:23a1f8d7111875d1a610e1c605eebccaf871b033894472d4306cbb841de548ae
```

### 16.4 Disconnect/reconnect invariant

The synthetic active match checkpoint has:

```yaml
revision: 7
event_sequence: 12
pending_interaction:
  interaction_instance_id: interaction-0001
  owner: player-0001
  legal_domain: [card-secret-a, card-secret-b]
  resume_cursor:
    effect_instance_id: effect-0001
    stage: choose-one
state_hash: sha256:0b90ce7f10e770a10df44a4868fe0ea49944885349b44fde52f40ba1ee67ba65
```

Presence steps:

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
state_hash_before: sha256:0b90ce7f10e770a10df44a4868fe0ea49944885349b44fde52f40ba1ee67ba65
state_hash_after: sha256:0b90ce7f10e770a10df44a4868fe0ea49944885349b44fde52f40ba1ee67ba65
pending_interaction_before: interaction-0001
pending_interaction_after: interaction-0001
resume_cursor_equal: true
legal_domain_equal: true
```

Player A pending projection canonical JSON and hash:

```json
{"event_sequence":12,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0002":{"hand_count":0}},"pending_interaction":{"cardinality":{"max":1,"min":1},"interaction_instance_id":"interaction-0001","legal_options":["card-secret-a","card-secret-b"],"optional":false,"owner":"player-0001","prompt":"Select one card","submitted_state":{}},"public_history":[],"public_state":{"action_type":"SCHEME","active_player_id":"player-0001"},"revision":7,"self_private_state":{"hand":["card-secret-a","card-secret-b"]},"viewer_player_id":"player-0001"}
```

```text
sha256:09c8b4cf6368bf1248cae72b33a1ad31d63a636453faa31cfc923fec2d0d1dae
```

Player B pending projection canonical JSON and hash:

```json
{"event_sequence":12,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0001":{"hand_count":2}},"pending_interaction":{"interaction_instance_id":"interaction-0001","owner":"player-0001","waiting":true},"public_history":[],"public_state":{"action_type":"SCHEME","active_player_id":"player-0001"},"revision":7,"self_private_state":{"hand":[]},"viewer_player_id":"player-0002"}
```

```text
sha256:8e26321ca209975ac5d1925d632715b6d4a9933e20da50a281cb5681cea91a9f
```

Forbidden assertions for player B:

```yaml
forbidden_values:
  - secret-a
  - secret-b
  - card-secret-a
  - card-secret-b
forbidden_json_paths:
  - $.pending_interaction.legal_options
  - $.self_private_state.opponent_hand
```

### 16.5 Snapshot plus ordered tail

Take an approved deterministic snapshot at revision 7, sequence 12, hash:

```text
sha256:0b90ce7f10e770a10df44a4868fe0ea49944885349b44fde52f40ba1ee67ba65
```

Apply exact tail:

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
  type: InteractionClosed
  caused_by_command_id: command-choice-0001
  public_payload:
    interaction_instance_id: interaction-0001
```

The fixture's synthetic reducer also moves `card-secret-a` and `card-0001` to public discard and resolves `effect-0001` as part of the same accepted batch.

Expected final checkpoint:

```yaml
revision: 8
event_sequence: 14
pending_interaction: null
state_hash: sha256:ef58e80da4a6a0fb821dc80dc4b1c71333dfa7e78dbd7b90ed22d653e96c473c
rng_calls_delta: 0
```

Player A final projection:

```json
{"event_sequence":14,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0002":{"hand_count":0}},"pending_interaction":null,"public_history":[{"card_instance_id":"card-secret-a","event_type":"CardMoved","to_zone":"discard"}],"public_state":{"action_type":null,"active_player_id":"player-0001","discard":["card-secret-a","card-0001"]},"revision":8,"self_private_state":{"hand":["card-secret-b"]},"viewer_player_id":"player-0001"}
```

```text
sha256:0dfba59b3f47f335e798f68b733ca8872afe687e92ebc3f6b7ddb690d6ec3d9d
```

Player B final projection:

```json
{"event_sequence":14,"legal_actions":[],"lifecycle":"ACTIVE","match_id":"match-0001","opponent_hidden_state_summary":{"player-0001":{"hand_count":1}},"pending_interaction":null,"public_history":[{"card_instance_id":"card-secret-a","event_type":"CardMoved","to_zone":"discard"}],"public_state":{"action_type":null,"active_player_id":"player-0001","discard":["card-secret-a","card-0001"]},"revision":8,"self_private_state":{"hand":[]},"viewer_player_id":"player-0002"}
```

```text
sha256:d0853cafdf4c2efaf1775377b69c3b44e44d68ee491d2a98dc051a659a765903
```

The same final hashes must result from uninterrupted execution and snapshot-plus-tail replay.

## 17. Runner conformance

A conforming runner must fail on:

- unknown fixture/schema version;
- definition digest mismatch;
- unexpected ID or RNG consumption;
- command fingerprint mismatch;
- duplicate request re-execution;
- event sequence gap/reordering;
- accepted/rejected result mismatch;
- state/projection hash mismatch;
- authority binding mismatch;
- presence data appearing in deterministic state;
- missing or changed pending interaction on reconnect;
- private value leakage;
- snapshot-tail divergence.

Until a runner exists, reviewers validate schema completeness, canonical examples and cross-document consistency. No document may claim executable fixture validation without an actual runner and Candidate-bound evidence.
