# Battlefield graph schema

The playable slice treats a battlefield as source-backed data. Gameplay never
infers movement or range from artwork at runtime.

## Required fields

```yaml
schema_version: 1
id: stable-battlefield-id
display_name: Human name
zones: [stable-zone-id]
spaces:
  - id: stable-space-id
    x: 0..100
    y: 0..100
    zones: [stable-zone-id]
    starting_seat: 1 # optional
edges:
  - from: stable-space-id
    to: stable-space-id
    type: normal
special_connections: []
components: []
rules: []
sources: []
validation: {}
```

Space IDs are independent of artwork coordinates. `x` and `y` are presentation
coordinates only and have no rules authority. Every ordinary edge is stored once
as an undirected connection. Runtime adjacency is built symmetrically.

## Validation invariants

- every space ID is non-empty and unique;
- every space belongs to one or more registered zones;
- each starting seat appears exactly once;
- edge endpoints exist and differ;
- no undirected edge is duplicated in either orientation;
- the graph is connected;
- runtime JSON and the test fixture are byte-identical;
- the canonical launch-map space, zone, start and edge inventory is exact;
- source-backed visual transcription receives independent final QA.

The schema leaves extension arrays for future battlefields, but the launch slice
does not implement directional connections, doors, secret passages, items or
set-specific battlefield rules.
