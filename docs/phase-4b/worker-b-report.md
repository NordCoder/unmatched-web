# Phase 4B Worker B research report

## Worker 4B-B Handoff

Branch: `phase-4b-worker-b-licensed`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported by the final worker handoff after the commit is created; a commit cannot embed its own SHA without changing that SHA.  
Assigned fighters: 18  
Fully verified at worker boundary: 11 — `daredevil`, `ghost-rider`, `luke-cage`, `moon-knight`, `genie`, `cloak-and-dagger`, `ms-marvel`, `black-widow`, `winter-soldier`, `doctor-strange`, `she-hulk`  
Blocked at worker boundary: 7 — `elektra`, `bullseye`, `dr-ellie-sattler`, `t-rex`, `houdini`, `squirrel-girl`, `spider-man`  
Quantity validation: FAIL — 17/18 per-definition deck sums reconcile; `spider-man` is intentionally blocked with only 11/30 physical cards independently recovered from allowed published/static evidence.  
Shared semantic/control files changed: none.  
Black Panther rewritten: no.

### Status split

Fighter manifests blocked on runtime semantics: `elektra`, `t-rex`, `squirrel-girl`, `spider-man`.

Deck manifests blocked on evidence/normalization completeness: `bullseye`, `dr-ellie-sattler`, `houdini`, `spider-man`.

A blocked status does not mean the printed character/deck is unknown. It means this worker could not prove every required field from the allowed source hierarchy or could not express a required rule without inventing a character-specific semantic primitive.

### Non-standard deck sizes

- `daredevil`: 22 published action cards.
- `elektra`: 20 published action cards.
- `black-widow`: 31 published action cards. `The Moscow Protocol` starts in hand before the remaining deck is shuffled and five cards are drawn.
- All other assigned fighters have a published 30-card action deck.

No implementation assumption that a fixed deck has 30 cards is introduced.

### Battlefield-item review

Hell's Kitchen, Redemption Row, Teen Spirit, For King and Country, and Brains and Brawn all include set-level battlefield items. The assigned fighter decks do not own those components as fighter resources or custom card zones. Their item effects therefore remain set/battlefield semantics for integration; no fighter manifest silently duplicates the set rule.

Brains and Brawn's official rulebook confirms the generic item timing model: scheme items are used as an action but are not scheme cards; combat items add their token value during combat and are returned to the box after combat. No shared rules file was edited here.

### Special topology/resource review

- `moon-knight`: one physical fighter/health track with persistent `active_identity`; the mandatory identity cycle is represented with existing state operations.
- `cloak-and-dagger`: two heroes; loss requires both heroes defeated.
- `squirrel-girl`: eight off-board summonable Squirrels; no Squirrel starts on the battlefield.
- `dr-ellie-sattler`: five public Insight tokens, returned to supply when removed.
- `ghost-rider`: public Hellfire counter, initial/cap 5.
- `black-widow`: setup-selected mission card and reveal-until-mission acquisition; no ownership-changing custom zone is needed.
- `t-rex`: large-fighter geometry is blocked pending a reusable footprint model.
- `elektra`: first defeat creates an off-board/not-defeated state before next-turn resurrection; current vocabulary cannot represent the removal half faithfully.

### Schema/effect extension proposals

#### 1. Multi-space fighter footprint and orientation

Affected: `t-rex`.

Evidence:
- official Jurassic Park Sattler vs. T. Rex set/rules material;
- current published/secondary T. Rex indexes.

Why current semantics are insufficient: `topology.fighters`, `MOVE`, and `PLACE` assume a fighter occupies one battlefield space. T. Rex may occupy two spaces, moves a large base with head/tail geometry, derives paths from its occupied footprint, ignores one-way direction, and cannot use secret passages.

Proposed generic extension: a fighter-level `footprint` model with occupied-space set, orientation, source/destination footprint validation, path-origin semantics, and movement-rule overrides. This should be reusable for any future multi-space fighter.

Blocked until integration: yes.

#### 2. Small-fighter shared occupancy and local damage fan-out

Affected: `squirrel-girl`.

Evidence:
- current Squirrel Girl published/secondary rules indexes.

Why current semantics are insufficient: up to four small fighters may share a space, including with another fighter, and damage to one Squirrel is applied to every other Squirrel in that same space. Ordinary occupancy and single-target damage cannot encode this.

Proposed generic extension: fighter `occupancy_class`, per-space occupancy capacity/compatibility, and a generic damage propagation rule keyed by co-located fighter class/identity.

Blocked until integration: yes.

#### 3. Off-board without defeat

Affected: `elektra`.

Evidence:
- current Elektra rules/rulings index.

Why current semantics are insufficient: the first time Elektra would be defeated, Elektra and all Hand pieces are removed from the battlefield but Elektra is not defeated; at the start of the next turn she returns with a new health profile and Hand placement constraints. `DEFEAT` is semantically wrong and `RETURN_FIGHTER` only models the return half.

Proposed generic extension: `REMOVE_FIGHTER` (or equivalent lifecycle transition) into an explicit off-board/not-defeated state, preserving ownership/identity and allowing a later `RETURN_FIGHTER`.

Blocked until integration: yes.

#### 4. Card-used-as-BOOST source event

Affected: `houdini` deck (`BOOSTED WITH` clauses, especially `The Big Reveal`).

Evidence:
- current published Houdini card-counter data;
- official Houdini vs. The Genie rulebook/product material.

Why current semantics are insufficient: the model can execute `BOOST`, but it has no event whose source is the card consumed as the BOOST. Houdini cards can resolve their own effects specifically when used to BOOST, and `The Big Reveal` can later return one such card.

Proposed generic extension: a `card_used_as_boost` / `boost_source_resolved` event carrying source card, boosted action/combat, controller, and normal disposition; allow effects on the BOOST source card and later selectors over cards used to BOOST.

Blocked until integration: deck yes; fighter ability itself is representable.

#### 5. Pre-defense scalar disclosure from a hidden combat card

Affected: `spider-man`.

Evidence:
- official Brains and Brawn rulebook;
- current Spider-Man rules index.

Why current semantics are insufficient: after the attack card is committed but before Spider-Man chooses a defense card, the attacker must disclose only the printed numeric value. `REVEAL` exposes the card itself, which leaks too much hidden information, and the current trigger vocabulary does not name this attack-commit / pre-defense-choice boundary.

Proposed generic extension: an attack declaration timing hook plus `DISCLOSE_CARD_FIELD` (or equivalent scalar disclosure) that reveals a selected printed field without changing card visibility.

Blocked until integration: yes.

### Evidence blockers / source gaps

#### `bullseye` deck
The full 12-definition quantity index reconciles to 30, and independently published records were recovered for `World's Greatest Assassin` and `Feint`. Per-card type/value/BOOST/effect metadata for the remaining definitions was not independently retrievable from an allowed current published source during this pass. No community `/decks/...` entry was used to fill it.

#### `dr-ellie-sattler` deck
The 14-definition quantity index reconciles to 30. Published UmDb records independently prove `You Never Had Control, That's the Illusion` and `Feint`; matching fan translations/reskins exposed by search were rejected. Remaining printed metadata is therefore blocked rather than copied from `/decks/...`.

#### `houdini` deck
All 13 definitions and quantities reconcile to 30, and current card-counter text proves values/BOOSTs/effects. The static text extraction exposes combat card type as an image, so unsupported type fields are left null. `BOOSTED WITH` normalization is additionally blocked by proposal 4.

#### `spider-man` deck
The official character ability and several official card records are proven, including `Web Shooters` and the Brains and Brawn rulebook example for `Thwip!`. The current allowed static sources did not expose the full official 30-card definition index. Search repeatedly surfaced community `/decks/...` Spider-Man builds with conflicting stats; all were rejected. The manifest intentionally contains only independently proven official definitions and fails per-definition quantity reconciliation.

### Ownership/card-zone review

Black Panther Phase 4A was used as the ownership/card-zone reference. None of Worker B's assigned fighters requires transferring ownership of a card. Black Widow mission acquisition changes only location (deck/revealed/hand) and preserves owner. Winter Soldier's BRAINWASHED marker is an effect category, not a storage zone. Houdini BOOST cards retain normal ownership while moving through boost/disposition contexts.

### Validation summary

- Assigned fighter ID coverage: PASS (18/18 fighter manifests, 18/18 card manifests).
- Worker-owned path scope: PASS.
- Shared semantic/control files untouched: PASS.
- Black Panther untouched: PASS.
- Fixed-deck counts: 17/18 per-definition PASS; Spider-Man BLOCKED.
- `usable_by` targets: PASS for all populated card records.
- Operation vocabulary: PASS for normalized operations; five reusable semantic gaps are isolated above rather than added ad hoc.
- Resource references: PASS (`hellfire`, `insight-supply`, summon pools/state are declared in fighter manifests).
- Hidden/public semantics: explicit where relevant (Black Widow mission search, public token pools, Spider-Man scalar-disclosure blocker).
- Fan/community `/decks/...` imported into official corpus: NO.
- Battlefield-item fighter ownership invented: NO.

### Files created

- `docs/fighters/phase-4b/daredevil.yaml`
- `docs/fighters/phase-4b/elektra.yaml`
- `docs/fighters/phase-4b/bullseye.yaml`
- `docs/fighters/phase-4b/ghost-rider.yaml`
- `docs/fighters/phase-4b/luke-cage.yaml`
- `docs/fighters/phase-4b/moon-knight.yaml`
- `docs/fighters/phase-4b/dr-ellie-sattler.yaml`
- `docs/fighters/phase-4b/t-rex.yaml`
- `docs/fighters/phase-4b/houdini.yaml`
- `docs/fighters/phase-4b/genie.yaml`
- `docs/fighters/phase-4b/cloak-and-dagger.yaml`
- `docs/fighters/phase-4b/ms-marvel.yaml`
- `docs/fighters/phase-4b/squirrel-girl.yaml`
- `docs/fighters/phase-4b/black-widow.yaml`
- `docs/fighters/phase-4b/winter-soldier.yaml`
- `docs/fighters/phase-4b/doctor-strange.yaml`
- `docs/fighters/phase-4b/she-hulk.yaml`
- `docs/fighters/phase-4b/spider-man.yaml`
- `docs/cards/phase-4b/daredevil.yaml`
- `docs/cards/phase-4b/elektra.yaml`
- `docs/cards/phase-4b/bullseye.yaml`
- `docs/cards/phase-4b/ghost-rider.yaml`
- `docs/cards/phase-4b/luke-cage.yaml`
- `docs/cards/phase-4b/moon-knight.yaml`
- `docs/cards/phase-4b/dr-ellie-sattler.yaml`
- `docs/cards/phase-4b/t-rex.yaml`
- `docs/cards/phase-4b/houdini.yaml`
- `docs/cards/phase-4b/genie.yaml`
- `docs/cards/phase-4b/cloak-and-dagger.yaml`
- `docs/cards/phase-4b/ms-marvel.yaml`
- `docs/cards/phase-4b/squirrel-girl.yaml`
- `docs/cards/phase-4b/black-widow.yaml`
- `docs/cards/phase-4b/winter-soldier.yaml`
- `docs/cards/phase-4b/doctor-strange.yaml`
- `docs/cards/phase-4b/she-hulk.yaml`
- `docs/cards/phase-4b/spider-man.yaml`
- `docs/phase-4b/worker-b-report.md`
