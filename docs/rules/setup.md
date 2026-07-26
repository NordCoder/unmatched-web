# Two-player setup

This file defines the generic setup path for a normal two-player competitive match. Character-specific pre-game choices and altered setup order are layered later.

## Standard setup sequence

`SETUP-001` — Resolve the two participating characters and the battlefield before fighter placement.  
Source: CORE p. 3.

`SETUP-002` — Each player receives the components required by their character: action cards, character card, hero miniature, sidekick tokens if any, health trackers, and character-specific components.  
Source: CORE p. 3.

`SETUP-003` — Initialize each fighter to its printed starting health. A sidekick without a health dial uses the generic 1-health rule unless set-specific rules override it.  
Source: CORE p. 3.

`SETUP-004` — Build the character's starting draw pile according to that character's deck-construction rules, shuffle it, and draw the starting hand. The generic core starting hand is 5 cards.  
Source: CORE p. 3. Character-specific deck construction is Phase 3/4 content.

`SETUP-005` — A generic fixed-deck character therefore enters placement with a shuffled face-down deck, a 5-card hand, an empty discard pile, initialized fighter health, and all required off-board components initialized.  
Source: normalized state derived from CORE p. 3.

## Player order and starting slots

`SETUP-010` — The physical Core Rules assign starting slot 1 and the first turn to the younger player, and slot 2 to the older player.  
Source: CORE p. 3.

`SETUP-011` — The engine must represent the resolved `player_order`/`starting_player` explicitly and must not require age as persistent game-state data. A digital lobby policy (for example, random assignment or mutually agreed assignment) is a product-layer choice unless strict tabletop setup replication is required. Once assignment is resolved, gameplay semantics are identical: player 1 occupies starting slot 1 and acts first.  
Source: implementation normalization of SETUP-010; this does not change post-assignment game rules.

## Hero placement

`SETUP-020` — Player 1 places their hero on battlefield starting slot 1. Player 2 then places their hero on starting slot 2.  
Source: CORE p. 3.

## Sidekick placement

`SETUP-030` — Immediately after placing a hero, that player's sidekicks are placed one at a time into separate unoccupied spaces sharing at least one zone with the hero's starting space.  
Source: CORE p. 3 plus FIELD-001 occupancy invariant.

`SETUP-031` — If the hero's starting space belongs to multiple zones, each sidekick may be placed in any unoccupied space belonging to any of those zones; different sidekicks need not use the same one of those zones.  
Source: CORE p. 3.

`SETUP-032` — If no legal unoccupied space exists in the hero's starting zone set for a sidekick, that sidekick may instead be placed in any unoccupied battlefield space.  
Source: CORE p. 3 plus FIELD-001.

## Setup completion

`SETUP-040` — After both players have completed required placement and all character-specific setup hooks, the player assigned to starting slot 1 begins the first turn.  
Source: CORE p. 3.

## Extension point for later characters

The publisher rulings include a more detailed ordering for characters with pre-game choices (for example Buffy, Geralt of Rivia, Yennefer & Triss, and Alice). That generalized setup pipeline is intentionally deferred to Phase 2/3 because it introduces ordered hooks around character selection, deck construction, hand draw, hero placement and sidekick placement.

`SETUP-050` — Character-specific setup may insert or replace steps, but no implementation may invent their ordering. Each such hook must be sourced and registered before that character becomes supported.  
Source: REF10 major ruling “Setup Order”; TERM-001.