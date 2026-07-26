# Phase 4B Worker B research report

## Worker 4B-B Handoff

Branch: `phase-4b-worker-b-licensed`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported by the final worker handoff after the final commit; a commit cannot embed its own SHA without changing that SHA.  
Assigned fighters: 18  
Fully verified at worker boundary: 11 — `daredevil`, `ghost-rider`, `luke-cage`, `moon-knight`, `genie`, `cloak-and-dagger`, `ms-marvel`, `black-widow`, `winter-soldier`, `doctor-strange`, `she-hulk`  
Blocked at worker boundary: 7 — `elektra`, `bullseye`, `dr-ellie-sattler`, `t-rex`, `houdini`, `squirrel-girl`, `spider-man`  
Quantity validation: PASS — 18/18 published deck quantities reconcile with their declared construction.  
Shared semantic/control files changed: none.  
Black Panther rewritten: no.

### Follow-up research result

A second targeted source pass was performed for all seven previously blocked fighters. It used official Restoration rulebooks/product material first, current UmDb published records second, and current ruling indexes/FAQ material for clarifications. Community `unmatched.cards/decks/...` pages were explicitly excluded.

The follow-up closed almost every factual/card-corpus gap:

- `elektra`: the first-life removal/non-defeat state, next-turn resurrection, 7→9 health profile, discard-to-deck shuffle, Hand return rules, interaction with revival effects, remaining actions and exhaustion timing are now independently supported. Remaining blocker is semantic only.
- `bullseye`: all 12 definitions, quantities, types, values and effects are recovered. The only unverified printed field is the BOOST value of `Ricochet`; it remains `null` rather than inferred.
- `dr-ellie-sattler`: all 14 action-card definitions now reconcile to 30 with complete normalized card metadata. Follow-up rulings prove five reusable Insight tokens and that identical Insight tokens can stack in one space. The remaining blocker is generic token-instance battlefield positioning.
- `t-rex`: two-space footprint, rotation/pathing, one-way override, secret-passage prohibition, two-space melee reach, Momentous Shift interaction, one-fighter identity and end-turn draw ordering are now supported. Remaining blocker is semantic only.
- `houdini`: all 13 definitions now have quantities, types, printed values/BOOST values and normalized published effects. Remaining blocker is the generic event/history required by `BOOSTED WITH` and `The Big Reveal`.
- `squirrel-girl`: official small-fighter rules are now supplemented by the current ruling correcting the printed `empty space` restriction. The ability is optional; when all eight Squirrels are on the battlefield it repositions one; occupied-space legality follows small-fighter limits. Remaining blocker is semantic only.
- `spider-man`: the full 12-definition/30-card deck is now recovered. The official Brains and Brawn rulebook directly proves Spidey-Sense timing and several card definitions. Remaining blocker is semantic only: disclosing only one printed field of a still-hidden committed attack card.

### Status split after follow-up

Semantic/runtime blockers:

- `elektra` — off-board but explicitly not defeated;
- `dr-ellie-sattler` — finite battlefield-token instances with same-type stacking and return-to-supply lifecycle;
- `t-rex` — multi-space fighter footprint/orientation/pathing;
- `houdini` — effects on the card consumed as a BOOST source plus source-card history;
- `squirrel-girl` — small-fighter shared occupancy, pass-through, adjacency and damage propagation/attribution;
- `spider-man` — pre-defense disclosure of one scalar field from a hidden committed card.

Residual source/evidence blocker:

- `bullseye` — only `Ricochet.boost` remains unverified.

### Non-standard deck sizes

- `daredevil`: 22 published action cards.
- `elektra`: 20 published action cards.
- `black-widow`: 31 published action cards; `The Moscow Protocol` begins in hand before the remaining cards are shuffled and the normal starting draw occurs.
- All other assigned fighters have a published 30-card action deck.

No implementation assumption that a fixed deck has 30 cards is introduced.

### Battlefield-item review

Hell's Kitchen, Redemption Row, Teen Spirit, For King and Country, and Brains and Brawn include set-level battlefield items. The assigned fighter decks do not own those components as fighter resources or private card zones. Their effects remain set/battlefield semantics for integration.

Brains and Brawn's official rulebook confirms that scheme items consume an action without counting as a scheme card, while combat items are attached when a combat card is played and modify that combat under their own timing/cancellation rules. No shared rules file was edited by this worker.

### Generic semantic extension proposals

#### 1. Multi-space fighter footprint and orientation

Affected: `t-rex`.

Evidence: current published T. Rex data plus current rulings.

Why insufficient: current fighter occupancy and `MOVE`/`PLACE` semantics assume one occupied battlefield space. T. Rex occupies up to two, rotates the extended base during movement, may originate movement from either occupied space, ignores one-way direction, cannot use secret passages and remains one fighter for triggers.

Proposed extension: fighter `footprint`/orientation state with occupied-space set, source/destination footprint validation, rotation/path semantics and movement overrides.

Blocked until integration: yes.

#### 2. Small-fighter shared occupancy and local damage propagation

Affected: `squirrel-girl`.

Evidence: official Teen Spirit special rules plus current ruling correcting the printed `empty space` restriction.

Why insufficient: up to four small fighters may share a space in addition to another ordinary/opposing fighter; small fighters and opposing fighters may pass through each other; co-located fighters are adjacent; same-type small fighters share incoming damage; current ruling also makes an occupied space a legal Squirrel placement destination subject to those limits.

Proposed extension: fighter `occupancy_class`, per-space occupancy compatibility/capacity, shared-space adjacency/pass-through rules and generic same-type damage fan-out with source attribution semantics.

Blocked until integration: yes.

#### 3. Off-board without defeat

Affected: `elektra`.

Evidence: published Elektra ability and current rulings/FAQ.

Why insufficient: first-life Elektra and all Hand leave battlefield presence while Elektra is explicitly not defeated. `DEFEAT` is wrong and `RETURN_FIGHTER` only models the later return half.

Proposed extension: `REMOVE_FIGHTER` or equivalent lifecycle transition to an explicit off-board/not-defeated state, preserving fighter identity and ownership for later return.

Blocked until integration: yes.

#### 4. Card-used-as-BOOST source event and history

Affected: `houdini`, especially `BOOSTED WITH` clauses and `The Big Reveal`.

Evidence: official Houdini vs. The Genie material, published UmDb cards and current rulings.

Why insufficient: `BOOST` can apply a value, but the model has no event belonging to the card consumed as the BOOST source and no generic historical selector for cards used to BOOST one specific effect/combat card.

Proposed extension: `card_used_as_boost` / `boost_source_resolved` event containing source card instance, controller, boosted target/context and disposition, plus explicit captured source-card history where a later effect needs it.

Blocked until integration: yes.

#### 5. Pre-defense scalar disclosure from a hidden committed card

Affected: `spider-man`.

Evidence: official Brains and Brawn rulebook.

Why insufficient: after an opponent commits an attack against Spider-Man but before Spider-Man chooses a defense, only the attack card's printed banner value is announced. Ordinary `REVEAL` exposes card identity and leaks too much information.

Proposed extension: reusable attack-commit/pre-defense-choice timing hook plus `DISCLOSE_CARD_FIELD(card, printed_value)` or equivalent field-level disclosure that does not alter card visibility.

Blocked until integration: yes.

#### 6. Battlefield token-instance positioning

Affected: `dr-ellie-sattler`; expected to be reusable for other finite battlefield component pools.

Evidence: current Sattler rules/rulings establish five Insight token instances, placement after fighter movement, return to supply, coexistence with other components and same-type Insight stacking.

Why insufficient: `token_pool` represents the finite supply count, but current persistent-state types cannot represent individual physical token identities transitioning between supply and battlefield spaces, especially when multiple identical tokens occupy the same space. A single `space_ref` or a mathematical set of spaces would lose multiplicity.

Proposed extension: generic battlefield token instances with stable token IDs, current location (`supply` or `space_ref`), visibility, placement legality, same-space multiplicity, selectors/counts and return-to-supply operation.

Blocked until integration: yes.

### Source-gap detail

#### `bullseye`

The published deck structure is now almost fully recovered. Quantities reconcile to 30 and current deck statistics reconcile the physical type distribution. Exact component/card sources and current references support all printed metadata currently populated. One field remains intentionally unresolved: `Ricochet` BOOST. No fan copy or aggregate guess was used to fill it.

#### Resolved former gaps

- `dr-ellie-sattler`: former per-card metadata gap resolved; deck is `verified`.
- `houdini`: former combat-card type gap resolved; printed deck metadata complete. Only semantic BOOST-source handling remains.
- `spider-man`: former 11/30 corpus gap resolved; full 30-card deck is represented and quantity-reconciled.
- `elektra`, `t-rex`, `squirrel-girl`: rule/ruling evidence strengthened enough that their blockers are no longer missing-rule questions.

### Ownership/card-zone review

Black Panther Phase 4A remains the ownership/card-zone reference. None of Worker B's assigned fighters requires ownership transfer. Black Widow mission acquisition changes only card location; Houdini BOOST cards retain their owner while moving through BOOST/disposition contexts. No foreign-card ownership semantics were added.

### Validation summary

- Assigned fighter ID coverage: PASS — 18/18 fighter manifests and 18/18 card manifests.
- Worker-owned path scope: PASS.
- Shared semantic/control files untouched: PASS.
- Black Panther untouched: PASS.
- Fixed-deck quantity sums: PASS — 18/18.
- `usable_by` targets: PASS for populated card records.
- Resource/source semantics: explicit; Sattler's previously lossy single-`space_ref` token representation was removed and converted into a precise integration blocker.
- Hidden/public semantics: explicit where relevant; Spider-Man's field-only disclosure remains a named blocker rather than using over-broad `REVEAL`.
- Fan/community `unmatched.cards/decks/...` data imported into official corpus: NO.
- Battlefield-item fighter ownership invented: NO.
- Shared generic semantic extensions implemented by worker: NO; all six are proposals only.

### Files created / worker-owned scope

37 files total:

- 18 fighter manifests under `docs/fighters/phase-4b/` for exactly the assigned Worker B fighter IDs;
- 18 card manifests under `docs/cards/phase-4b/` for exactly the assigned Worker B fighter IDs;
- `docs/phase-4b/worker-b-report.md`.

No merge to `main` was performed.