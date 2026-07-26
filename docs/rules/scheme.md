# Scheme action

## Legality and declaration

`SCHEME-001` — To take a Scheme action, the active player chooses a scheme card from hand that at least one undefeated fighter is authorized to resolve.  
Source: CORE p. 8; banner restrictions: CORE p. 6.

`SCHEME-002` — The active player declares the specific undefeated fighter resolving the scheme before resolving its effect.  
Source: CORE p. 8.

`SCHEME-003` — A scheme whose only authorized fighter is defeated cannot be played. An `ANY` scheme may be resolved by any undefeated fighter the player controls.  
Source: CORE pp. 6, 8.

`SCHEME-004` — If no scheme card in hand can be legally resolved by an undefeated authorized fighter, Scheme is not a legal action choice.  
Source: normalized action-legality consequence of CORE p. 8.

## Resolution

`SCHEME-010` — The chosen scheme card is played face up. Its effect is resolved, then the card is moved to that player's discard pile unless a specific effect replaces that destination.  
Source: CORE p. 8; replacement exceptions deferred to character rules.

`SCHEME-011` — Core effect semantics from `effect-resolution-baseline.md` apply to scheme effects: mandatory unless optional wording says otherwise; resolve as much as possible if part is impossible; effect-specific rules may override general rules.  
Source: CORE pp. 1, 13.

## Discard pile

`SCHEME-020` — Each player maintains a separate discard pile.  
Source: CORE p. 8.

`SCHEME-021` — The discard pile is public information. Either player may inspect either discard pile at any time.  
Source: CORE p. 8.

`SCHEME-022` — Played action cards normally enter their owner's discard pile after they have finished resolving. Combat cards use the later combat cleanup step rather than being discarded immediately on play.  
Source: CORE pp. 8, 12.

## Later extensions

Character-specific scheme prerequisites, ongoing schemes, missions, cards that end a turn, and schemes that remain in play are not generalized here. They require explicit Phase 2/4 mechanics and provenance.