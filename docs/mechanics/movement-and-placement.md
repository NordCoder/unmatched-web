# Movement, placement, swap and destination semantics

Movement and placement are distinct operations. The engine must not implement placement as `move with unlimited distance`, because path, occupancy, trigger and failure semantics differ.

## Move

`PLACE-001` — `MOVE` follows the normal movement rules unless the source explicitly overrides them: traverse legal connections, respect occupancy/path restrictions, and end in a legal destination.  
Source: CORE p. 7.

`PLACE-002` — When an effect moves an opposing fighter, movement legality is interpreted from that fighter owner's perspective.  
Source: CORE p. 7.

`PLACE-003` — `up to N spaces` permits moving zero spaces.  
Source: CORE p. 13; RULES-HUB General Rulings.

`PLACE-004` — Moving through a space means entering and then leaving it; merely ending on/entering a space is not by itself `moving through` for effects that use that wording.  
Source: RULES-HUB General Rulings.

## Place

`PLACE-010` — `PLACE` does not trace a movement path and does not consume movement distance. It relocates the fighter directly to the selected destination if placement succeeds.  
Source: CORE p. 7; RULES-HUB General Rulings.

`PLACE-011` — Under ordinary occupancy rules a successful placement cannot result in two ordinary fighters occupying the same space. Set/fighter mechanics such as small fighters may explicitly override ordinary occupancy and must be modeled as destination-policy exceptions.  
Source: CORE pp. 4, 7; later specific-rule precedence.

## Choice legality is not placement success

`PLACE-020` — For a placement effect that does **not** say `empty` or `other`, an occupied space can be a legal selection if it otherwise satisfies the wording. After selection, the placement operation fails because ordinary occupancy prevents placing the fighter there; the fighter remains where it was.  
Source: REF10 Major Rulings: Placement; RULES-HUB Rule Changes & Errata (correction of the prior ruling).

`PLACE-021` — If placement wording requires an `empty` space, occupied spaces are excluded from the choice domain before selection.  
Source: ordinary text constraint + PLACE-020 distinction.

`PLACE-022` — If placement wording requires an `other` space, the fighter's current space is excluded from the choice domain even if choosing it would otherwise create a failed placement.  
Source: ordinary text constraint + PLACE-020 distinction.

The engine therefore needs two separate functions:

```text
selectable_destinations(effect_text_constraints, state)
placement_succeeds(selected_destination, occupancy_policy, state)
```

Do not collapse them into one `legal_destination` predicate.

## No destination

`PLACE-030` — If a mandatory move/place/return operation has no valid selectable destination at all, skip that impossible part and continue resolving independently resolvable parts of the effect.  
Source: CORE p. 13; RULES-HUB General Rulings on no valid move/place space.

`PLACE-031` — Failure to place a fighter is not automatically equivalent to defeat, removal, or moving zero spaces; it is simply a failed placement unless the source says otherwise.  
Source: REF10 Placement ruling.

## Restrictions on leaving a space

`PLACE-040` — A restriction stating that a fighter `cannot leave their space` prevents ordinary maneuver movement, movement effects, and placement effects that would relocate that fighter. It does not prevent removal caused by defeat.  
Source: RULES-HUB General Rulings.

This requires movement/placement legality to query active restrictions before applying an operation.

## Swap

`PLACE-050` — `SWAP` is a distinct composite relocation, not two unrelated placements executed naively in sequence. Its intent is to exchange the fighters' occupied positions without the first relocation failing solely because the other fighter currently occupies the destination.  
Source: published swap effects; normalized atomicity requirement.

`PLACE-051` — Large fighters and special occupancy can require additional source-defined choices during a swap (which occupied space, facing/orientation, etc.). Those cases remain battlefield/fighter-specific and must be documented with the relevant content rather than generalized from ordinary one-space fighters.  
Source: RULES-HUB General Rulings on swapping with large fighters.

## Triggers caused by relocation

`PLACE-060` — An ability that triggers on a fighter `leaving`, `entering`, or changing zone may be triggered by movement, placement, removal, or another relocation mechanism when its authoritative wording/ruling says so. Therefore trigger generation must observe semantic transitions, not only Maneuver actions.  
Source: published abilities and rulings such as Tomoe Gozen; specific fighter applicability remains Phase 4 content.

The relocation operation should emit transition facts (`left_space`, `entered_space`, `left_zone`, `entered_zone`, `removed_from_board`) that fighter-specific triggers can consume later.
