# Maneuver, movement, BOOST and exhaustion

## Maneuver sequence

`MAN-001` — A Maneuver resolves in this order: (1) mandatory draw attempt, then (2) optional fighter movement. Movement cannot occur before the draw attempt is resolved.  
Source: CORE p. 6.

## Mandatory draw and exhaustion

`MAN-010` — A normal Maneuver attempts to draw exactly 1 card from the top of the active player's deck into their hand.  
Source: CORE p. 6.

`MAN-011` — An empty deck is not reshuffled. When a player must draw a card that does not exist, every one of that player's fighters takes 2 damage for each missing draw.  
Source: CORE p. 6.

`MAN-012` — Generalized for an effect that requires `N` draws while only `K < N` cards are available: resolve the available draws and there are `N-K` missing draws; each fighter takes `2 × (N-K)` exhaustion damage. The exact ordering of multi-draw/exhaustion events and defeat-triggered effects is a Phase 2 timing topic.  
Source: CORE p. 6; arithmetic normalization.

`MAN-013` — Exhaustion does not make Maneuver illegal. Attempting the mandatory draw is still part of the action.  
Source: CORE p. 6.

## Optional movement

`MAN-020` — After the draw step, the player may move zero or more of their fighters.  
Source: CORE p. 7.

`MAN-021` — During a normal Maneuver, each fighter may be moved at most once. Fighters are resolved one at a time, in an order chosen by the active player; one fighter's movement must finish before the next begins.  
Source: CORE p. 7.

`MAN-022` — For each fighter moved, choose an integer distance from 0 through the current move value. Each traversed step must go to an adjacent space.  
Source: CORE p. 7.

`MAN-023` — A fighter may traverse spaces occupied by friendly fighters but may not finish on an occupied space.  
Source: CORE p. 7.

`MAN-024` — A fighter may neither traverse nor finish on a space occupied by an opposing fighter under generic movement rules.  
Source: CORE p. 7.

`MAN-025` — Movement of earlier fighters changes occupancy for later fighters in the same Maneuver. Legal paths are therefore evaluated against the current board state at the time that fighter moves.  
Source: normalized consequence of CORE p. 7's sequential movement.

## Normal Maneuver BOOST

`MAN-030` — Before moving fighters during a Maneuver, the active player may perform the normal movement BOOST once: discard 1 card from hand and add its BOOST value to the move value used for this Maneuver.  
Source: CORE pp. 7–8.

`MAN-031` — The discarded BOOST card's printed effect is ignored.  
Source: CORE p. 8.

`MAN-032` — A card may still be discarded for BOOST even if its banner's fighter is defeated and the card could no longer legally be played.  
Source: CORE p. 8.

`MAN-033` — The resulting boosted move value is available independently to each fighter moved during that Maneuver; movement points are not a shared pool.  
Source: CORE pp. 7–9; normalized from “each fighter” using the move value.

`MAN-034` — Effects may permit BOOSTing properties other than Maneuver movement. Those are not the normal Maneuver BOOST and must be modeled according to the granting effect.  
Source: CORE p. 8.

## Movement granted by effects

`MAN-040` — If an effect lets a player move an opposing fighter, normal movement constraints are interpreted from that fighter's owner's perspective unless the effect says otherwise.  
Source: CORE p. 7.

`MAN-041` — A `place` operation does not traverse adjacent spaces and therefore is not movement. Its successful destination must be empty under the generic core rule.  
Source: CORE p. 7.

`MAN-042` — Do not infer placement target-selection legality solely from `MAN-041`; later official rulings distinguish selecting a space from successfully placing into it. That distinction is reserved for Phase 2.  
Source: REF10 major ruling “Placement”.