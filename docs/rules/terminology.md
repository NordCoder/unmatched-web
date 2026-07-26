# Terminology

This file defines terms used by the normalized competitive core specification.

## Actors and fighters

`TERM-010` — A **player** is one participant controlling one character in standard two-player play.  
Source: CORE pp. 1–3.

`TERM-011` — A **character** is the complete playable selection associated with a character card, action-card deck, hero, optional sidekick(s), health trackers, and any character-specific components.  
Source: CORE pp. 2–3; normalized term for implementation.

`TERM-012` — A **fighter** is any hero or sidekick controlled by a player.  
Source: CORE p. 2.

`TERM-013` — A **hero** is the character's primary fighter. In the ordinary core model it is represented by a miniature and has starting health identified by the character card.  
Source: CORE p. 2.

`TERM-014` — A **sidekick** is any non-hero fighter belonging to the character. A character may have zero, one, or multiple sidekicks.  
Source: CORE p. 2.

`TERM-015` — A sidekick with no health dial has 1 starting health unless set-specific rules explicitly provide otherwise.  
Source: CORE p. 3; exception mechanism: TERM-001.

`TERM-016` — **Friendly fighters** are fighters controlled by the same player. Team play extends friendly status to a teammate's fighters. **Opposing fighters** are fighters controlled by an opponent.  
Source: CORE pp. 7, 18–19.

## Battlefield

`TERM-020` — A **battlefield** is the graph-like play area containing spaces, adjacency connections, zones, starting spaces, and possibly set-specific features.  
Source: CORE pp. 3–4; normalized representation.

`TERM-021` — A **space** is one fighter position on the battlefield. At most one fighter may occupy a space.  
Source: CORE p. 4.

`TERM-022` — Two spaces are **adjacent** when the battlefield has a direct connection between them.  
Source: CORE p. 4.

`TERM-023` — A **zone** is the set of all spaces carrying the same zone color. The spaces need not be contiguous. A multicolored space belongs to every represented zone.  
Source: CORE p. 4.

## Cards and card data

`TERM-030` — An **action card** belongs to a character's action-card deck and has a card type, banner/user restriction, name, optional combat value, optional effect text, BOOST value, and deck/copy metadata.  
Source: CORE p. 5.

`TERM-031` — Core card types are **attack**, **defense**, **scheme**, and **versatile**. A versatile card may be used for attack or defense and counts as both for effects that inspect card type.  
Source: CORE pp. 5, 10.

`TERM-032` — A card's **banner** determines which fighter(s) may play/resolving it. `ANY` permits any of the player's fighters.  
Source: CORE p. 6.

`TERM-033` — If a banner refers to a group of interchangeable sidekicks, the player chooses which surviving token is resolving that card. Other defeated sidekicks of the same type do not prevent a surviving one from resolving it.  
Source: CORE p. 6.

`TERM-034` — A card's **printed value** is the value physically printed on the card before any modification. It remains distinct from the effective value used in combat.  
Source: CORE p. 13.

`TERM-035` — A card's **BOOST value** is a separate numeric property used when a rule/effect permits that card to BOOST something. Discarding a card for a normal Maneuver BOOST does not resolve the card's printed effect.  
Source: CORE p. 8.

`TERM-036` — A player's **deck** is their face-down draw pile. A player's **hand** is private held cards. A player's **discard pile** is face up and public to both players at all times.  
Source: CORE pp. 6, 8.

## Actions, combat and health

`TERM-040` — An **action** is one unit of the active player's action economy. Core action types are Maneuver, Scheme and Attack.  
Source: CORE p. 5.

`TERM-041` — **Combat** is the resolution started by an Attack action after attacker/target declaration and combat-card selection. Its core windows are `IMMEDIATELY`, `DURING COMBAT`, combat damage, `AFTER COMBAT`, cleanup, and attack-level additional effects.  
Source: CORE pp. 10–12.

`TERM-042` — **Damage** is a health-reducing game result. **Combat damage** is specifically the non-negative damage produced by comparing the played attack value against the played defense value (or zero when undefended). Damage caused directly by card/ability effects is not combat damage.  
Source: CORE pp. 12–13.

`TERM-043` — **Health** is the remaining life of a fighter. Starting health is the maximum recoverable health unless a specific effect says otherwise.  
Source: CORE p. 14.

`TERM-044` — A fighter is **defeated** when its health is reduced to zero. A defeated fighter is removed from the battlefield immediately and cannot recover health under the general rule.  
Source: CORE p. 14.

## Movement vocabulary

`TERM-050` — **Move** follows normal movement/path rules unless the granting effect explicitly changes them.  
Source: CORE p. 7.

`TERM-051` — **Place** relocates a fighter without traversing a path. The resulting occupied destination must be an empty space under the core rule. Selection behavior for wording that does not say `empty`/`other` is a Phase 2 ruling topic.  
Source: CORE p. 7; REF10 p. 2 for the later placement ruling.

## Normative language

`TERM-060` — Core effects are mandatory unless they explicitly use optional wording such as `may`.  
Source: CORE p. 13.

`TERM-061` — `Up to N` permits any amount from zero through N, including zero.  
Source: CORE p. 13.

`TERM-062` — The specification uses **active player** for the player whose turn is currently resolving, even where the physical rulebook says “the player whose turn it is.”  
Source: normalized implementation term; GAMEEND rules preserve the official semantics.