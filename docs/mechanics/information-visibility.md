# Information visibility

A server-authoritative implementation must preserve the tabletop distinction between public state, private state, committed hidden choices, and temporary reveals. Visibility is a gameplay property, not only a UI concern.

## Visibility classes

`INFO-001` — The engine distinguishes at least:

```text
PUBLIC
OWNER_PRIVATE
PLAYER_PRIVATE(viewer)
COMMITTED_HIDDEN
TEMPORARILY_REVEALED
```

Source: normalized requirement from CORE card-zone and combat procedures.

`INFO-002` — Public event history may record information that was legitimately revealed at an earlier time even after the underlying card returns to a hidden zone. Historical knowledge is not the same as current authoritative location knowledge.  
Source: normalized information-model consequence.

## Hands

`INFO-010` — A player's hand-card identities are private from opponents unless a rule explicitly reveals them or grants another player permission to look at them.  
Source: CORE hidden combat selection and card effects that explicitly allow looking at/revealing hand cards; normalized competitive implication.

`INFO-011` — A player's own client may receive their complete current hand identities. Opponents must not receive those identities through ordinary game-state synchronization.  
Source: INFO-010; normalized online-play requirement.

`INFO-012` — Hand size is public game-state metadata for the digital implementation. It mirrors directly observable tabletop state and is required to apply/verify public hand-limit and card-count interactions without disclosing identities.  
Source: normalized tabletop-observability policy; card identities remain private under INFO-010.

## Decks and discard piles

`INFO-020` — Deck order/card identities are hidden except where a rule grants reveal/look access.  
Source: CORE decks are shuffled and kept face down; reveal/look mechanics.

`INFO-021` — A discard pile is public. Cards are face up, and either player may inspect the discard pile.  
Source: CORE p. 8.

`INFO-022` — Shuffling a hidden zone invalidates positional knowledge of cards in that zone unless a rule explicitly preserves information. The event that a particular card was seen historically may remain known, but its post-shuffle position is not.  
Source: normalized consequence of shuffle/randomization.

## Combat commitments

`INFO-030` — The attacker's committed attack card remains hidden from the defender until the combat reveal step.  
Source: CORE p. 10.

`INFO-031` — The defender's committed defense card remains hidden until the same reveal step. The server must represent commitment separately from public reveal.  
Source: CORE p. 10.

`INFO-032` — Combat cards reveal simultaneously after the attack and defense commitments (or explicit no-defense decision) are complete.  
Source: CORE p. 10.

`INFO-033` — A reconnecting opponent before reveal must not receive the identity of an already committed hidden combat card. A reconnecting owner may receive enough private state to restore their own committed choice without being allowed to change it.  
Source: normalized reconnect requirement from INFO-030–032.

## Reveal versus look at

`INFO-040` — `REVEAL` makes the specified card identity public to all players for the source-defined duration. If the effect does not change the card's destination, it returns/remains where the effect specifies after the reveal.  
Source: REF10 Card Effects: Revealing Cards.

`INFO-041` — `LOOK_AT` grants identity knowledge only to the player(s) authorized by the source. It is not automatically equivalent to public reveal.  
Source: published `look at` effects and contrast with REF10 reveal rule; normalized distinction.

`INFO-042` — A choice based on privately viewed information remains owned by the authorized viewer; the server may validate the choice using hidden authoritative state without exposing the full inspected set to unauthorized clients.  
Source: normalized online-play requirement.

## Random hidden selection

`INFO-050` — When a rule requires random selection from a hidden zone, the authoritative engine performs/verifies the random selection without exposing the full hidden zone. The selected card becomes public only if its resulting destination/effect makes it public (for example, entering the face-up discard pile).  
Source: normalized secure implementation requirement; destination visibility from CORE.

## Logs and spectator/replay policy

`INFO-060` — The authoritative event log must distinguish private payloads from public payloads. A replay/spectator endpoint must not reveal information earlier than it became public in the original game unless an explicit post-game full-information mode is implemented as a separate product policy.  
Source: normalized privacy/replay requirement.

## Out of scope

Table-talk, voluntary communication of private information, screenshots, and anti-cheating policy are social/platform concerns rather than deterministic Unmatched rules. They must not change the engine's authoritative visibility classification.
