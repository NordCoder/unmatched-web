# Deterministic Engine Fixture Contract

## Status

```text
status: normative-foundation
schema_id: unmatched.engine.fixture/v2
parent_issue: #19
correction_round_1: #32
correction_round_2: #35
stabilization_round_3: #36
```

This contract defines the language-neutral deterministic evidence shared by Core Runtime, Rules Mechanics and Lead integration.

## Normative artifacts

```text
docs/engine/fixtures/schema-v2.json
docs/engine/fixtures/foundation-v2.json
docs/engine/fixtures/foundation-v2-transition-audit.json
scripts/validate_engine_foundation.py
tests/architecture/test_engine_foundation.py
```

The previous v1 JSON artifacts remain historical correction evidence. They are not the acceptance Candidate after stabilization round 3.

## 1. Required evidence categories

A conforming suite contains exactly one independently executable case for each category:

```text
create
join
idempotency
reconnect
replay
```

The committed validator fails when a category is absent, duplicated or replaced by prose-only evidence.

The cases prove:

- `CreateMatch` requires no pre-existing match or player runtime identity;
- `JoinMatch` authenticates a principal and creates a new match-scoped player identity;
- same key and same fingerprint returns the original durable result without execution;
- same key and different fingerprint returns derived `DUPLICATE_CONFLICT` without a second command-result record;
- reconnect changes operational presence but preserves deterministic state hash and pending interaction;
- snapshot plus contiguous events reaches the exact final state;
- viewer projections match the final state and forbidden hidden values are absent.

## 2. Closed schema surfaces

Every normative object boundary used by the suite is closed. Unknown fields are rejected in:

- suite, definition, state and case objects;
- command envelopes and command-kind payloads;
- normalized request identities and payloads;
- event envelopes and event-kind payloads;
- pending interactions and card state;
- projections, private state and transition audit entries.

JSON Schema supplies structural validation. `scripts/validate_engine_foundation.py` supplies cross-field semantic validation that JSON Schema cannot express reliably, including case-kind requirements, command/payload correlation, definition references, collision behavior and replay transitions.

## 3. Pinned semantic inputs

Every case references a committed `definition_id`.

Each definition entry contains:

```text
ruleset
capability registry identity
setup definition
synthetic action/effect definitions
canonical SHA-256 digest
```

States, events and projections repeat the definition identity. A missing or mismatched definition reference is a validation failure. Required semantics may not exist only in prose or validator source code.

## 4. Strict parsing and canonical bytes

Before semantic validation:

1. decode UTF-8 JSON;
2. reject duplicate object keys;
3. normalize all object keys and string values to Unicode NFC;
4. reject duplicate keys created by NFC normalization;
5. reject unknown fields through the closed schema;
6. serialize the supported JSON value subset using NFC-normalized RFC 8785-compatible canonical bytes;
7. hash the UTF-8 bytes using SHA-256.

The normative suite contains NFC-equivalence, escaping/non-ASCII and duplicate-key rejection vectors.

## 5. Idempotency model

The authoritative store contains one immutable command-result record per:

```text
(principal_id, command_id)
```

For the same fingerprint, the original durable result is returned without execution, allocation or new gameplay events.

For a different fingerprint:

- response class is `derived_collision`;
- rejection code is `DUPLICATE_CONFLICT`;
- command-result record delta is zero;
- gameplay-event delta is zero;
- the original record remains unchanged and undisclosed.

## 6. Reconnect boundary

Operational presence fields are forbidden in deterministic states and hashes.

The reconnect case binds two different operational client states to one unchanged deterministic state entry. The exact pending interaction ID and state hash must be preserved.

## 7. Exact replay

The replay case contains a pinned snapshot and contiguous event sequence. Each event has a closed envelope and definition identity.

The committed transition audit records:

```text
sequence
event type
pre-state hash
exact changed JSON pointers
post-state hash
```

The validator applies every event through the generic fixture reducer, rejects sequence gaps and undeclared mutations, and requires the final state and projection hashes to match.

## 8. Executable conformance gate

Run:

```bash
python -m unittest discover -s tests/architecture -p "test_engine_foundation.py"
python scripts/validate_engine_foundation.py
```

Regression tests include mutations for:

- missing JoinMatch or reconnect evidence;
- unknown nested command fields;
- missing or mismatched definition pins;
- collision persistence or durable collision response;
- actor identity supplied to JoinMatch;
- reconnect hash drift;
- event sequence gaps;
- undeclared replay state changes;
- transition hash corruption;
- hidden-information exposure;
- operational presence inside deterministic state;
- duplicate keys before and after NFC normalization.

A handoff claim is not acceptance evidence unless these checks pass in the exact-head `engine-foundation` workflow.
