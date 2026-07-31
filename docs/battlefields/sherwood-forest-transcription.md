# Sherwood Forest transcription record

## Scope

This record supports the playable Robin Hood vs Bigfoot slice. It transcribes
only data required for ordinary movement, adjacency, ranged-zone legality and
starting placement.

## Evidence

- Unmatched Picks map browser, full-board Sherwood Forest reference visual,
  checked 2026-07-29. Used for space inventory, zone membership, printed start
  markers and visible connections.
- The Unmatched Club Sherwood Forest map page, checked 2026-07-29. Used as an
  independent count cross-check: 30 spaces and 7 zones.
- Official Robin Hood vs Bigfoot product/rules material, checked 2026-07-29.
  Used for ordinary space, zone, movement, ranged-attack and starting-player
  rules.

No copyrighted board image is committed. The repository records source URLs and
the resulting graph only.

## Stable identities

Spaces are `s01` through `s30`. IDs follow the transcription record rather than
artwork labels or pixel coordinates. The seven normalized zones are:

```text
gray
light-gray
green
brown
light-green
orange
yellow
```

Starting positions:

```text
seat 1 → s20
seat 2 → s19
```

## Graph inventory

The canonical machine fixture contains:

```text
spaces: 30
zones: 7
undirected ordinary edges: 39
special connections: 0
components: 0
battlefield-specific rules: 0
```

Each visible connection is represented once. Runtime code derives the reverse
adjacency rather than storing a second edge.

## Review status

The implementor completed the source-backed transcription and automated
structural validation. Independent edge-by-edge visual QA remains part of the
single final playable-slice review required by Issue #51; this document does not
self-approve that gate.
