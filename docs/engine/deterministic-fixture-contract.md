# Deterministic Engine Fixture Contract

## Status

```text
status: normative-foundation
schema_id: unmatched.engine.fixture/v1
parent_issue: #19
correction_round_1: #32
correction_round_2: #35
```

This document defines the language-neutral fixture contract shared by Core Runtime, Rules Mechanics and Lead integration.

Normative machine-readable artifacts:

```text
docs/engine/fixtures/schema-v1.json
docs/engine/fixtures/foundation-v1.json
docs/engine/fixtures/foundation-v1-transition-audit.json
```

The JSON artifacts are normative. Prose summarizes their execution and canonicalization rules.

## 1. Suite and case execution

A suite contains an ordered `cases` array. Order is reporting order only.

Every case is isolated:

- no inherited deterministic state;
- no inherited authority records;
- no inherited idempotency records;
- no inherited snapshots;
- no shared ID, RNG or choice cursor;
- providers reset to the case-declared arrays;
- provider underflow, unexpected allocation and unused required entries fail validation.

Each case declares its complete initial source, authority/idempotency seeds, providers, steps and expected final evidence.

## 2. Strict parsing

Before semantic validation:

1. decode UTF-8 without BOM;
2. reject malformed JSON, invalid Unicode and unpaired surrogates;
3. reject duplicate object keys;
4. normalize every object key and string value to Unicode NFC;
5. reject any duplicate key created by NFC normalization;
6. validate against `fixtures/schema-v1.json`.

Unknown fields in suite, case, initial source, state, event, projection and audit objects are rejected.

Absent optional fields are omitted. JSON `null` is a distinct explicit value and is legal only where the schema permits it. Array order is always significant.

## 3. Canonical bytes and hashes

Canonical bytes are:

```text
strict parsed JSON
→ Unicode NFC keys and string values
→ RFC 8785 JSON Canonicalization Scheme
→ UTF-8 bytes
→ SHA-256
```

Fixture v1 permits only schema-declared integers in the exact safe range. Floating point, NaN and infinity are forbidden.

Worker-local `sort_keys` approximations are not normative. RFC 8785 defines key ordering, escaping and number serialization.

Hash syntax:

```text
sha256:<64 lowercase hexadecimal digits>
```

`foundation-v1.json` includes:

- NFC-equivalent input vectors with identical bytes/hash;
- escapable and non-ASCII string bytes/hash;
- duplicate-key rejection evidence.

## 4. Initial sources

### 4.1 `new_match`

Contains complete lifecycle inputs, authority/idempotency seeds and deterministic providers needed to create the initial state through commands/events.

### 4.2 `canonical_state`

Contains:

```text
state
state_revision
last_event_sequence
state_hash
authority_records
idempotency_records
providers
```

The declared revision and sequence must equal the corresponding state fields. The recomputed state hash must match.

### 4.3 `snapshot_plus_tail`

Contains:

```text
snapshot:
  state
  revision
  last_event_sequence
  state_hash
ordered_tail_events
authority_records
idempotency_records
providers
expected_final_state
expected_final_hash
```

The first tail sequence must be `snapshot.last_event_sequence + 1`; all following sequences are contiguous.

## 5. Idempotency evidence

The fixture distinguishes:

```text
durable_result / first
durable_result / duplicate
derived_collision / collision
```

For the same key and fingerprint, the original durable result is returned without execution.

For the same key and a different fingerprint:

- response is derived `DUPLICATE_CONFLICT`;
- execution count delta is zero;
- command-result record count delta is zero;
- gameplay event count delta is zero;
- original record is unchanged;
- repeated collision requests remain derived and non-disclosing.

## 6. Exact event mode

Each event contains:

```text
event_schema_version
event_id
match_id
sequence
revision
event_type
caused_by_command_id
ruleset_version
public_payload
private_payloads_by_player
```

`parent_event_id` and `source_ref` are omitted when absent and required by event schemas where applicable. Empty payloads are explicit `{}`.

Fixture v1 defines exact payload families for its event types. No field is inferred from prose.

## 7. Snapshot-tail transition semantics

Envelope application updates:

```text
/revision
/event_sequence
/history_cursor/last_event_sequence
```

Fixture event payloads authorize these additional changes:

| Event | Authorized state changes |
| --- | --- |
| `ChoiceSubmitted` | selected values of the identified pending interaction |
| `CardMoved` | identified card zone, position and resulting face state |
| `InteractionClosed` | identified pending interaction becomes `null` |
| `EffectStageChanged` | identified effect stage and status |
| `EffectDequeued` | identified effect is removed from the queue |
| `ActionCompleted` | identified action state becomes `null` |

No other state path may change for the normative case.

`foundation-v1-transition-audit.json` records, for every event:

```text
sequence
event_type
pre_state_hash
changed_json_pointers
post_state_hash
```

Each `pre_state_hash` equals the previous entry's `post_state_hash`. The final audit hash equals both the suite expected final hash and the uninterrupted execution hash.

## 8. Projection evidence

Viewer projections are recomputed from the final deterministic state and authority context. Operational presence is excluded.

Each expected projection supplies:

- complete projection object;
- projection hash;
- forbidden values for that viewer.

## 9. Conformance checks

A conforming artifact validator rejects:

- unknown schema/version or field;
- duplicate keys before or after NFC normalization;
- invalid Unicode;
- non-integer numeric values;
- canonicalization/hash mismatch;
- case state/provider leakage;
- command fingerprint mismatch;
- duplicate re-execution;
- collision persistence or original-result disclosure;
- event envelope omission;
- sequence gap/reordering;
- undeclared state mutation;
- transition-audit discontinuity;
- final state/projection mismatch;
- presence data in deterministic evidence.

An executable gameplay runner is not part of this correction Candidate. Static artifact parsing, schema validation, canonical hash recomputation and transition application are the required evidence.
