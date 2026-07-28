# Battlefield Graph Schema

## Status

```text
schema_version: 1
status: draft-mvp
branch: phase-5-mvp-battlefield
parent_issue: #18
```

The engine consumes structured battlefield definitions. It must never interpret battlefield artwork during a match.

## 1. Definition boundary

A `BattlefieldDefinition` contains immutable graph and setup data:

```yaml
schema_version: 1
id:
display_name:
set_ids: []
status: researching | verified | blocked
spaces: []
edges: []
special_connections: []
components: []
rules: []
sources: []
validation: {}
```

Runtime fighter/component occupancy belongs in `BattlefieldState`, not in the definition.

## 2. Spaces

```yaml
spaces:
  - id: sf-01
    zones: [zone-red]
    starting_slots:
      - player: 1
        fighter_slot: hero
    capacity:
      class: standard
      units: 1
    artwork_anchor:
      x_normalized: 0.0
      y_normalized: 0.0
    notes: []
```

### Required fields

- `id` — stable battlefield-local identifier, independent of artwork coordinates;
- `zones` — one or more canonical zone IDs printed for the space;
- `starting_slots` — zero or more source-defined setup assignments;
- `capacity` — optional static capacity metadata when the battlefield itself changes occupancy rules.

`artwork_anchor` is optional UI/evidence metadata. It cannot determine adjacency or legality.

A space ID must not encode a fighter, current occupant or transient setup result.

## 3. Edges

```yaml
edges:
  - id: sf-edge-001
    from: sf-01
    to: sf-02
    type: standard
    direction: bidirectional
    tags: []
    source_evidence:
      image_ref:
      observation:
```

### Rules

- one canonical edge record represents one printed connection;
- undirected printed adjacency uses `direction: bidirectional`;
- a normalized pair order is required for duplicate detection;
- directionality may be declared only when authoritative evidence supports it;
- visual proximity without a printed connection is not an edge;
- edge type does not imply a gameplay rule unless the schema/rules registry defines it.

## 4. Zones

Zones are identifiers attached to spaces. A space may belong to multiple zones.

The definition must preserve printed zone membership exactly. Zone color names may be normalized to stable IDs, but source-facing labels/colors remain in evidence metadata.

Ranged-zone legality is derived from shared zone membership plus current game rules. It is not represented as precomputed attacker-target pairs.

## 5. Starting slots

Starting information must distinguish:

- player/seat assignment;
- hero versus sidekick/group slot where the source specifies it;
- numbered setup markers;
- optional/alternate setup rules.

A battlefield definition may expose source-defined slots; the setup resolver determines which fighter instances occupy them.

## 6. Special connections

```yaml
special_connections:
  - id:
    kind: secret_passage | door | one_way | teleport | custom_registered
    endpoints: []
    enabled_by_default: true
    state_key: null
    rules: []
    source_evidence: []
```

Special connections are explicit objects when they have state, nonstandard traversal, asymmetric legality or separate source rules. They are not encoded as ordinary edges with undocumented tags.

## 7. Battlefield components

```yaml
components:
  - definition_id:
    kind: item | token | door | path_component | other_registered
    supply_count:
    setup:
    allowed_anchors:
      spaces: []
      edges: []
    rules: []
```

Components have immutable definitions and runtime instances. Non-card battlefield components must not be represented as action cards.

## 8. Rules

```yaml
rules:
  - id:
    trigger_or_query:
    established_behavior:
    capability_requirement:
    sources: []
```

Generic battlefield behavior belongs in shared capabilities. A battlefield rule supplies parameters and references; it does not justify a battlefield-ID branch in core engine code.

## 9. Sources and evidence

```yaml
sources:
  - id:
    authority:
    url_or_registry_ref:
    checked_at:
    supports: []
    image_digest_optional:
    storage_policy:
```

The repository may store metadata/digests without storing copyrighted source images. Every transcribed space, zone, edge and special rule must be traceable to evidence.

UmDb/secondary board catalogs are indexes only and cannot establish graph topology.

## 10. Validation block

```yaml
validation:
  yaml_parse: pass
  unique_space_ids: pass
  known_edge_endpoints: pass
  duplicate_edges: pass
  bidirectional_normalization: pass
  graph_connectivity: pass
  zone_membership_review: pass
  starting_slot_review: pass
  edge_by_edge_visual_qa:
    verdict: pass
    reviewer:
    evidence:
  qualifications: []
```

A graph may intentionally contain disconnected components only when authoritative rules explain how they function. The validator must otherwise reject disconnected ordinary movement graphs.

## 11. Machine fixture

The test fixture is a normalized derivative of the YAML definition and contains no source-image dependence.

It must support deterministic queries for:

- `adjacent(space_a, space_b)`;
- `neighbors(space)`;
- `shared_zones(space_a, space_b)`;
- shortest legal graph distance under a supplied traversal policy;
- starting slots;
- enabled special connections;
- legal anchors for battlefield components.

The fixture generator must sort space IDs, normalized edges, zones and component definitions to produce stable diffs and hashes.

## 12. MVP gate

The schema may be marked verified after Sherwood Forest proves that:

- every visible space/edge/zone can be represented without artwork interpretation at runtime;
- stable IDs remain usable for UI coordinates and event history;
- edge-by-edge independent QA can reference evidence precisely;
- no simple-board assumption prevents later directionality, doors, passages, high ground, items or large-fighter restrictions.
