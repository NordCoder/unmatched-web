# Multiplayer core deltas

The first implementation target remains two-player competitive Unmatched. This file records the generic free-for-all and team-play deltas contained in the modern Core Rules so that two-player assumptions are not accidentally promoted to universal invariants.

Full multiplayer support is **not** part of the Phase 1 implementation-ready subset: character-specific rulings, multiple-hero characters, elimination timing, and newer global rulings still require later reconciliation.

## Free-for-all

`FFA-001` — Generic free-for-all supports 3 or 4 players and uses a battlefield with enough numbered starting spaces for all players.  
Source: CORE p. 17.

`FFA-002` — Resolve player order and assign the corresponding numbered starting spaces before fighter placement. Each player places their hero on their assigned starting space and then places sidekicks according to the normal setup placement rule.  
Source: CORE p. 17; SETUP-020–032.

`FFA-003` — During a player's first turn only, that player may attack only (a) the next player in turn order or (b) a player who has already taken their first turn. After the first round, normal target ownership restrictions apply.  
Source: CORE p. 17.

`FFA-004` — In a combat effect that says `your opponent`, the phrase refers to the other player participating in that combat, not every opposing player in the match.  
Source: CORE p. 17.

`FFA-005` — In the published free-for-all baseline, a player is eliminated when their hero is defeated; remove that player's remaining sidekicks from the battlefield when elimination occurs.  
Source: CORE p. 17.

`FFA-006` — The generic free-for-all victory condition is to be the last player whose hero remains undefeated.  
Source: CORE p. 17.

### Reconciliation requirement

`FFA-005`/`FFA-006` are recorded as the published multiplayer baseline, but Phase 2/4 must explicitly reconcile them with:

- the modern start/end-of-action winner-check erratum used for two-player play;
- multiple-hero characters;
- effects that can defeat multiple heroes during one action;
- any official rulings that define elimination timing in multiplayer.

Until that reconciliation is complete, do not reuse the two-player `GAMEEND-*` algorithm unchanged for free-for-all.

## Team play

`TEAM-001` — Generic team play supports 3 or 4 players. With four players there are two teams of two players. With three players, one team consists of one player controlling two characters while the opposing team consists of two players controlling one character each.  
Source: CORE pp. 18–19.

`TEAM-002` — Fighters controlled by a teammate count as friendly fighters for movement and effects that distinguish friendly from opposing fighters. Each player otherwise controls their own fighters and hand unless a rule explicitly says otherwise.  
Source: CORE pp. 18–19.

`TEAM-003` — Generic four-player setup alternates teams in player order rather than seating all members of one team consecutively. The published order is Team A player 1, Team B player 1, Team A player 2, Team B player 2.  
Source: CORE p. 18.

`TEAM-004` — Each player places their hero on the corresponding numbered starting space and places sidekicks under the normal setup placement rule.  
Source: CORE p. 18; SETUP-020–032.

`TEAM-010` — Turns follow the same alternating team order used by setup and repeat cyclically.  
Source: CORE p. 18.

`TEAM-011` — Teammates may communicate, but each player's hand remains that player's own game information unless a rule or chosen product policy explicitly permits additional disclosure. The physical Core Rules do not create shared hands.  
Source: CORE pp. 18–19; normalized information boundary.

`TEAM-020` — When a player's hero is defeated, remove that hero. Under the generic team rules, the player may continue taking turns as long as that player still controls at least one undefeated sidekick.  
Source: CORE p. 19.

`TEAM-021` — When all fighters controlled by one player are defeated, that player is eliminated and no longer takes turns.  
Source: CORE p. 19.

`TEAM-022` — The generic team victory condition is defeating both heroes on the opposing team.  
Source: CORE p. 19.

### Three-player control model

`TEAM-030` — In a three-player team match, the solo player controlling two characters keeps those characters' decks, hands, fighters, health, and character state distinct unless a specific rule says otherwise.  
Source: CORE p. 19; normalized implementation consequence of controlling two characters.

`TEAM-031` — The exact action/turn ownership semantics for a single human controlling two character seats should be represented as two player/seat turns belonging to one account, rather than merging both characters into one action economy.  
Source: implementation normalization of the published alternating-character turn order.

## Phase boundary

The following multiplayer topics remain outside the Phase 1 gate:

- multiple-hero character elimination semantics;
- dormant-player rulings;
- team interactions involving hidden information and private choices;
- simultaneous elimination across multiple players;
- exact application of delayed game-end checks to free-for-all/team modes;
- set-specific battlefields designed for multiplayer.

`FFA-090` / `TEAM-090` — No multiplayer mode should be marked `developer-ready` until these later rulings are reconciled. The purpose of this file is to prevent two-player-only assumptions from contaminating the data model.