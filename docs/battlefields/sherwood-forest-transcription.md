# Sherwood Forest Battlefield Transcription

## Selection

```text
battlefield_id: sherwood-forest
status: EVIDENCE_ACQUISITION_REQUIRED
branch: phase-5-mvp-battlefield
parent_issue: #18
base_sha: 106ae552ce597cde954c0a1b22374ef446974ce2
```

Sherwood Forest is the first battlefield graph target because:

- the integrated set registry places it in `robin-hood-vs-bigfoot`;
- Robin Hood and Bigfoot are both `verified` with `verification.integration: ready`;
- the product registry records no set-specific battlefield mechanics;
- using the fighters and battlefield from one product reduces setup/evidence mismatch in the first vertical slice.

This selection does **not** establish graph topology. No space, zone or edge may be entered from memory.

## Evidence gate

Before `sherwood-forest.yaml` is created, collect:

1. a high-resolution full-board reference where every space boundary, connection, zone color and starting marker is readable;
2. authoritative product/rule references establishing the battlefield identity and any setup rule;
3. evidence metadata sufficient to distinguish source image, edition and crop;
4. an evidence-use policy: metadata/digest in Git, binary storage outside Git unless separately authorized.

The chosen image must support zoomed edge-by-edge review. A board thumbnail or catalog preview is insufficient.

## Transcription procedure

### Pass 1 — spaces and zones

- enumerate every visible space once;
- assign temporary visual labels;
- record every printed zone membership;
- record all starting markers;
- mark any unreadable region as blocked rather than inferred.

### Pass 2 — stable IDs and coordinates

- assign stable `sf-NN` IDs in a documented visual traversal order;
- attach optional normalized UI anchor coordinates;
- verify that IDs do not encode transient occupants or current graph assumptions.

### Pass 3 — edges

- inspect every space boundary at high zoom;
- record each printed connection exactly once as a normalized pair;
- record explicit non-connections where visual proximity could cause ambiguity;
- do not infer symmetry from appearance: represent the printed edge, then normalize its bidirectional rule.

### Pass 4 — independent QA

A reviewer who did not perform the transcription must compare:

- every source space against one manifest ID;
- every source edge against one edge record;
- every manifest edge against visible evidence;
- every zone membership and starting marker;
- graph connectivity and any suspicious articulation/disconnected region.

The QA report must issue one of:

```text
PASS
PASS_WITH_NONBLOCKING_QUALIFICATIONS
FAIL
BLOCKED
```

## Planned outputs

```text
docs/battlefields/sherwood-forest.yaml
docs/battlefields/sherwood-forest-validation.md
tests/fixtures/battlefields/sherwood-forest.json
scripts/validate_battlefield.py
```

## Current blockers

```text
high_quality_reference_image: not yet registered
graph_transcription: not started
independent_visual_qa: not started
```

No gameplay implementation may use a guessed Sherwood Forest graph while these blockers remain.
