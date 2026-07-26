# Battlefield core rules

This document covers only the generic space/zone model. Doors, high ground, one-way paths, secret passages, items, portals, large-fighter markings, and other battlefield-specific behavior are Phase 3/5 content.

## Occupancy and adjacency

`FIELD-001` — A battlefield consists of spaces. Each space contains at most one fighter.  
Source: CORE p. 4.

`FIELD-002` — Two spaces are adjacent exactly when the battlefield defines a direct connection between them. Generic adjacency is symmetric unless a battlefield-specific rule explicitly replaces normal traversal semantics.  
Source: CORE p. 4; directional exceptions deferred.

`FIELD-003` — Adjacency is independent of zone membership. Two adjacent spaces may belong to different zones.  
Source: CORE p. 4.

## Zones

`FIELD-010` — Every zone is identified by a battlefield color/zone identifier. All spaces carrying that identifier belong to the same zone even if they are visually separated on the board.  
Source: CORE p. 4.

`FIELD-011` — A space may belong to multiple zones simultaneously.  
Source: CORE p. 4.

`FIELD-012` — Two fighters are in the same zone when the sets of zone identifiers on their spaces intersect.  
Source: CORE p. 4; normalized predicate.

`FIELD-013` — Zone membership does not itself create adjacency and does not provide movement connectivity. It is used by ranged targeting and effects that explicitly refer to zones.  
Source: CORE pp. 4, 10.

## Starting spaces

`FIELD-020` — A battlefield may designate numbered starting spaces. Standard two-player setup uses starting slots 1 and 2. Multiplayer formats require additional starting slots as specified in `multiplayer-deltas.md`.  
Source: CORE pp. 3, 17–18.

## Engine-facing predicates

For the generic core rules, a battlefield representation must support these deterministic queries:

```text
is_occupied(space)
occupant(space)
are_adjacent(space_a, space_b)
zones(space)
share_zone(space_a, space_b)
starting_space(slot)
```

`FIELD-030` — Runtime legality must be derived from battlefield data, not inferred from artwork.  
Source: implementation normalization required by the core semantics.

## Visual equivalence

`FIELD-040` — If a board is printed with an accessibility/pattern variant and a plain-color variant of the same battlefield, those variants have identical gameplay unless a set-specific rule says otherwise.  
Source: CORE p. 4.