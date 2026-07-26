# Phase 4B Worker B research report

## Worker 4B-B Handoff

Branch: `phase-4b-worker-b-licensed`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported by the final worker handoff after the final commit; a commit cannot embed its own SHA without changing that SHA.  
Assigned fighters: 18  
Rule/card research complete: 18/18.  
Fully verified against the current semantic model: 12 — `daredevil`, `bullseye`, `ghost-rider`, `luke-cage`, `moon-knight`, `genie`, `cloak-and-dagger`, `ms-marvel`, `black-widow`, `winter-soldier`, `doctor-strange`, `she-hulk`.  
Blocked only on generic semantic/runtime extensions: 6 — `elektra`, `dr-ellie-sattler`, `t-rex`, `houdini`, `squirrel-girl`, `spider-man`.  
Residual source/evidence blockers: none.  
Quantity validation: PASS — 18/18 published deck quantities reconcile with their declared construction.  
Shared semantic/control files changed: none.  
Black Panther rewritten: no.

### Follow-up research result

A targeted second source pass was performed for the seven initially blocked fighters. Official Restoration rulebooks/product material were preferred, followed by current UmDb published records and current ruling/FAQ indexes. Community `unmatched.cards/decks/...` pages were explicitly excluded.

All factual/card-corpus gaps are now closed:

- `elektra`: first-life removal/non-defeat state, next-turn resurrection, 7→9 health profile, discard-to-deck shuffle and Hand return rules are supported. Remaining blocker is semantic only.
- `bullseye`: all 12 definitions, quantities, types, values, BOOST values and effects are recovered. Follow-up direct component imagery supplied during research visibly confirms `Ricochet` as value 3, BOOST 2, versatile, x3. No evidence blocker remains.
- `dr-ellie-sattler`: all 14 action-card definitions reconcile to 30 with normalized printed metadata. Rulings prove five reusable Insight tokens and same-space Insight stacking. Remaining blocker is generic token-instance battlefield positioning.
- `t-rex`: two-space footprint, rotation/pathing, one-way override, secret-passage prohibition, extended melee reach, one-fighter identity and relevant timing are supported. Remaining blocker is semantic only.
- `houdini`: all 13 definitions have quantities, types, printed values, BOOST values and normalized published effects. Remaining blocker is the generic event/history required by `BOOSTED WITH` and `The Big Reveal`.
- `squirrel-girl`: official small-fighter rules plus the current ruling correcting the printed `empty space` restriction are captured. Remaining blocker is semantic only.
- `spider-man`: the full 12-definition/30-card deck is represented. The official Brains and Brawn rulebook proves Spidey-Sense timing. Remaining blocker is semantic only: field-level disclosure from a hidden committed attack card.

### What `blocked` now means

The remaining six `blocked` statuses do **not** mean the game rules are unknown or ambiguous. Research has established the intended behavior. They mean the current shared fighter/effect schema cannot faithfully encode the mechanic without introducing a character-specific primitive. Those generic model decisions belong to the orchestration/integration stage.

### Non-standard deck sizes

- `daredevil`: 22 published action cards.
- `elektra`: 20 published action cards.
- `black-widow`: 31 published action cards; `The Moscow Protocol` begins in hand before the remaining cards are shuffled and the normal starting draw occurs.
- All other assigned fighters have a published 30-card action deck.

No implementation assumption that a fixed deck has 30 cards is introduced.

### Battlefield-item review

Hell's Kitchen, Redemption Row, Teen Spirit, For King and Country, and Brains and Brawn include set-level battlefield items. The assigned fighter decks do not own those components as fighter resources or private card zones. Their effects remain set/battlefield semantics for integration.

Brains and Brawn's official rulebook confirms that scheme items consume an action without counting as a scheme card, while combat items modify combat under their own timing/cancellation rules. No shared rules file was edited by this worker.

### Generic semantic extension proposals

#### 1. Multi-space fighter footprint and orientation

Affected: `t-rex`.

Known rule: T. Rex is one fighter whose extended base may occupy two battlefield spaces simultaneously; movement/pathing and placement operate on that footprint, with special one-way and secret-passage rules.

Current-model gap: ordinary fighter occupancy plus `MOVE`/`PLACE` assumes one occupied space.

Proposed extension: reusable fighter `footprint`/orientation state with occupied-space set, source/destination footprint validation, rotation/path semantics and movement overrides.

#### 2. Small-fighter shared occupancy and local damage propagation

Affected: `squirrel-girl`.

Known rule: multiple Squirrels can share a space under small-fighter capacity rules, may coexist with an ordinary/opposing fighter, use special shared-space/pass-through semantics, and damage to one Squirrel propagates to other Squirrels in the same space.

Current-model gap: normal occupancy and single-target damage cannot represent that behavior faithfully.

Proposed extension: fighter `occupancy_class`, per-space compatibility/capacity, shared-space adjacency/pass-through rules and generic same-type damage fan-out with preserved source attribution.

#### 3. Off-board without defeat

Affected: `elektra`.

Known rule: on the first would-be defeat, Elektra and all Hand pieces leave battlefield presence while Elektra is explicitly not defeated; she later returns under resurrection setup rules.

Current-model gap: `DEFEAT` is semantically wrong and `RETURN_FIGHTER` models only the return half.

Proposed extension: `REMOVE_FIGHTER` or equivalent lifecycle state for explicit off-board/not-defeated presence, preserving identity/ownership for later return.

#### 4. Card-used-as-BOOST source event and history

Affected: `houdini`, especially `BOOSTED WITH` clauses and `The Big Reveal`.

Known rule: effects can belong to the card consumed as the BOOST source, and later effects can refer back to a specific card used to BOOST a specific action/combat card.

Current-model gap: `BOOST` applies a value but exposes no generic source-card event/history relationship.

Proposed extension: `card_used_as_boost` / `boost_source_resolved` event containing source card instance, controller, boosted target/context and disposition, plus captured source-card history where later selection requires it.

#### 5. Pre-defense scalar disclosure from a hidden committed card

Affected: `spider-man`.

Known rule: after an opponent commits an attack against Spider-Man but before Spider-Man chooses defense, only the attack card's printed numeric value is announced; card identity remains hidden.

Current-model gap: ordinary `REVEAL` exposes the whole card and leaks too much information.

Proposed extension: reusable attack-commit/pre-defense-choice timing hook plus field-level disclosure such as `DISCLOSE_CARD_FIELD(card, printed_value)` that does not alter card visibility.

#### 6. Battlefield token-instance positioning

Affected: `dr-ellie-sattler`; expected to be reusable for other finite battlefield component pools.

Known rule: Sattler owns five physical Insight tokens that move between supply and battlefield spaces; multiple Insight tokens may occupy the same space.

Current-model gap: `token_pool` represents finite supply count, but the current persistent-state vocabulary cannot faithfully represent multiple physical token instances and same-space multiplicity.

Proposed extension: generic battlefield token instances with stable IDs or equivalent multiset semantics, current location (`supply` or `space_ref`), placement legality, same-space multiplicity, selectors/counts and return-to-supply operations.

### Ownership/card-zone review

Black Panther Phase 4A remains the ownership/card-zone reference. None of Worker B's assigned fighters requires ownership transfer. Black Widow mission acquisition changes only card location; Houdini BOOST cards retain their owner while moving through BOOST/disposition contexts. No foreign-card ownership semantics were added.

### Validation summary

- Assigned fighter ID coverage: PASS — 18/18 fighter manifests and 18/18 card manifests.
- Rule/card evidence coverage: PASS — no unresolved factual/source gaps remain.
- Fixed-deck quantity sums: PASS — 18/18.
- Bullseye final evidence gap: CLOSED — `Ricochet.boost = 2` from direct component imagery supplied during follow-up research.
- Worker-owned path scope: PASS.
- Shared semantic/control files untouched: PASS.
- Black Panther untouched: PASS.
- Fan/community `unmatched.cards/decks/...` data imported into official corpus: NO.
- Shared generic semantic extensions implemented by worker: NO; all six are proposals for orchestration/integration.
- Battlefield-item fighter ownership invented: NO.

### Files created / worker-owned scope

37 files total:

- 18 fighter manifests under `docs/fighters/phase-4b/` for exactly the assigned Worker B fighter IDs;
- 18 card manifests under `docs/cards/phase-4b/` for exactly the assigned Worker B fighter IDs;
- `docs/phase-4b/worker-b-report.md`.

No merge to `main` was performed.
