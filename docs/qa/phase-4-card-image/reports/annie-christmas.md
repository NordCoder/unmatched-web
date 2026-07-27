# Phase 4 Card-Image QA Report

## Batch

- annie-christmas
- doctor-strange
- dr-jill-trent
- she-hulk
- spider-man

## Verdict

FAIL

## Summary

Independent Phase 4 card-image QA completed for one five-fighter evidence batch. The external archive and all five nested fighter ZIPs were readable. Every supplied unique card image decoded successfully and was visually inspected. No missing, zero-byte, corrupt, or duplicate-binary card images were found.

Canonical comparison used the current branch tips captured during QA:

- `phase-4b-worker-b-licensed` -> `f383a639e980d9e753c742b20b7e085d3163502f`
- `phase-4b-worker-c-modern` -> `1f4fd973de66ab7995c0b2213e7dcda4163ff7c8`

Batch totals:

- fighters received: 5
- fighters fully checked: 5
- total unique card images: 60
- total images successfully inspected: 60
- unreadable images: 0
- missing images: 0
- duplicate images: 0
- quantity failures: 0
- material semantic discrepancies: 6
  - P1: 4
  - P2: 2
- non-gameplay metadata/mapping findings: 2 P3 findings covering 3 metadata facts

Fighter verdicts:

- `spider-man`: PASS
- `doctor-strange`: PASS WITH FINDINGS
- `she-hulk`: FAIL
- `annie-christmas`: FAIL
- `dr-jill-trent`: FAIL

## Fighter Results

### annie-christmas

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/annie-christmas.yaml`
- fighter manifest: `docs/fighters/phase-4b/annie-christmas.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL due to material semantic findings
- unique-card image completeness: PASS, 12/12 unique images present and readable
- printed card identity/content verification: all 12 card images inspected
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 12
- unique Git definitions: 12
- construction: fixed
- discrepancies: 1 P1, 1 P2, 1 P3
- verdict: FAIL

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| A Few More Pearls | PASS | PASS | PASS | PASS |
| Captain's Orders | PASS | P3 | FAIL P1 | FAIL |
| Keep Your Hands to Yourself | PASS | PASS | PASS | PASS |
| Lagniappe | PASS | PASS | PASS | PASS |
| Mississippi Queen | PASS | PASS | PASS | PASS |
| Striking Beauty | PASS | PASS | PASS | PASS |
| Better Together | PASS | PASS | PASS | PASS |
| Long Shot | PASS | PASS | PASS | PASS |
| Quite a Pair | PASS | PASS | PASS | PASS |
| Slick Talker | PASS | PASS | PASS | PASS |
| Bottom Dealing | PASS | PASS | P2 | FAIL |
| The Turn and the River | PASS | PASS | PASS | PASS |

Integration requirements directly confirmed by supplied card images:

- `C-REQ-012`: `Mississippi Queen` confirms the need for a health-result floor distinct from damage prevention. The printed card prevents Annie's health from dropping below 1 for the turn while the normalized model intentionally preserves damage dealt for downstream observers.
- `C-REQ-013`: cards such as `Bottom Dealing` and `Quite a Pair` confirm the need for ordered stages, result bindings, informed later choices, and dependent continuation.

### doctor-strange

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/doctor-strange.yaml`
- fighter manifest: `docs/fighters/phase-4b/doctor-strange.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS for gameplay; minor identity/presentation findings only
- unique-card image completeness: PASS, 11/11 unique images present and readable
- printed card identity/content verification: all 11 card images inspected
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 11
- unique Git definitions: 11
- construction: fixed
- discrepancies: 1 P3 finding covering printed-title punctuation and archive card-id normalization
- verdict: PASS WITH FINDINGS

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Cloak of Levitation | PASS | PASS | PASS | PASS |
| Steadfast Disciple | PASS | PASS | PASS | PASS |
| Master of Kamar-Taj | PASS | PASS | PASS | PASS |
| Bolts of Balthakk | PASS | PASS | PASS | PASS |
| Seven Suns of Cinnibus | PASS | PASS | PASS | PASS |
| The Mists of Munnopor | PASS | PASS | PASS | PASS |
| The Rings of Raggadorr | PASS | PASS | PASS | PASS |
| The Winds of Watoomb | PASS | PASS | PASS | PASS |
| Eye of Agamotto | PASS | PASS | PASS | PASS |
| No, Really, I'm a Doctor | PASS | P3 | PASS | P3 |

Integration requirements confirmed from supplied card images: none.

### dr-jill-trent

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/dr-jill-trent.yaml`
- fighter manifest: `docs/fighters/phase-4b/dr-jill-trent.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL due to material semantic findings
- unique-card image completeness: PASS, 13/13 unique images present and readable
- printed card identity/content verification: all 13 card images inspected
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 13
- unique Git definitions: 13
- construction: fixed
- discrepancies: 1 P1, 1 P2
- verdict: FAIL

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Hypnotist | PASS | PASS | PASS | PASS |
| Battle of Wits | PASS | PASS | PASS | PASS |
| Utility Belt | PASS | PASS | FAIL P1 | FAIL |
| Stasis Diffuser | PASS | PASS | PASS | PASS |
| Gyroscopic Jetpack | PASS | PASS | PASS | PASS |
| Indestructible Cloth | PASS | PASS | PASS | PASS |
| Laser Pen | PASS | PASS | PASS | PASS |
| Insightful Deduction | PASS | PASS | P2 | FAIL |
| Caught Red-Handed | PASS | PASS | PASS | PASS |
| Energizing Spray | PASS | PASS | PASS | PASS |
| Helpful Assistant | PASS | PASS | PASS | PASS |
| Sisters in Arms | PASS | PASS | PASS | PASS |
| Ace Fighter | PASS | PASS | PASS | PASS |

Integration requirements directly confirmed by supplied card images:

- `C-REQ-013`: `Insightful Deduction` directly requires reveal -> result binding -> dependent selection -> dependent ordering.

### she-hulk

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/she-hulk.yaml`
- fighter manifest: `docs/fighters/phase-4b/she-hulk.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL due to two gameplay-changing semantic findings
- unique-card image completeness: PASS, 12/12 unique images present and readable
- printed card identity/content verification: all 12 card images inspected
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 12
- unique Git definitions: 12
- construction: fixed
- discrepancies: 2 P1
- verdict: FAIL

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Lady Justice | PASS | PASS | PASS | PASS |
| The Defense Rests | PASS | PASS | PASS | PASS |
| Cease and Desist | PASS | PASS | PASS | PASS |
| Legalese | PASS | PASS | FAIL P1 | FAIL |
| Nerve Cluster Strike | PASS | PASS | FAIL P1 | FAIL |
| The Savage She-Hulk | PASS | PASS | PASS | PASS |
| Green Energy | PASS | PASS | PASS | PASS |
| Sensational | PASS | PASS | PASS | PASS |
| Omega-Level Threat | PASS | PASS | PASS | PASS |
| Double Jeopardy | PASS | PASS | PASS | PASS |
| Jennifer Walters, Esq. | PASS | PASS | PASS | PASS |
| Leap Toward | PASS | PASS | PASS | PASS |

Integration requirements directly confirmed by supplied card images:

- `B-REQ-009`: `The Savage She-Hulk` visibly permits spending an additional action during combat to make the card value 9, confirming that a legal action permission must be spendable as a source-defined cost rather than modeled as a raw counter decrement.

### spider-man

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/spider-man.yaml`
- fighter manifest: `docs/fighters/phase-4b/spider-man.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: PASS, 12/12 unique images present and readable
- printed card identity/content verification: all 12 card images inspected
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 12
- unique Git definitions: 12
- construction: fixed
- discrepancies: none
- verdict: PASS

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Momentous Shift | PASS | PASS | PASS | PASS |
| Disarming Shot | PASS | PASS | PASS | PASS |
| Snark | PASS | PASS | PASS | PASS |
| Friendly Neighborhood Spider-Man | PASS | PASS | PASS | PASS |
| Counter-Attack | PASS | PASS | PASS | PASS |
| Web Shooters | PASS | PASS | PASS | PASS |
| Right in the Face! | PASS | PASS | PASS | PASS |
| Spider-Sense Tingling! | PASS | PASS | PASS | PASS |
| Swinging Kick | PASS | PASS | PASS | PASS |
| Thwip! | PASS | PASS | PASS | PASS |
| Wall Crawler | PASS | PASS | PASS | PASS |
| With Great Power | PASS | PASS | PASS | PASS |

The fighter/deck pair declares `B-REQ-005`, but that requirement concerns Spider-Man's Spidey-Sense fighter ability rather than any card image in this evidence archive. Because no fighter ability component was supplied, this QA does not claim visual confirmation of `B-REQ-005`.

## Findings

### SH-P1-001 — Legalese mandatory discard lost

- fighter: `she-hulk`
- card: `Legalese`
- severity: P1
- expected: once discard mode is selected, each player with at least one card in hand must discard exactly one card; only an actually empty hand can contribute zero cards.
- observed: Git models each discard choice with `min: 0`, `max: 1`, and an `up_to_1_limited_by_hand_size` domain, permitting a player with cards in hand to voluntarily choose zero.
- evidence/reasoning: the printed card uses mandatory discard wording rather than `may` or `up to`. The current normalized choice broadens legal behavior and can change gameplay outcomes.
- expected semantic correction: require exactly one selected discard from each non-empty hand, with zero only as an unavoidable empty-hand partial resolution.

### SH-P1-002 — Nerve Cluster Strike restriction too narrow

- fighter: `she-hulk`
- card: `Nerve Cluster Strike`
- severity: P1
- expected: the opposing combat fighter may not leave their current space for the rest of the turn, except for removal caused by defeat or another source-defined exception.
- observed: Git explicitly prevents only `MOVE` and `PLACE` operations.
- evidence/reasoning: the project movement model treats `SWAP` as a distinct relocation composite, not as two ordinary placements. Therefore a MOVE/PLACE-only prohibition does not preserve the full printed cannot-leave-space restriction.
- expected semantic correction: model a general source-lifetime cannot-leave-space restriction that applies to MOVE, PLACE, SWAP, and equivalent non-defeat relocation semantics.

### AC-P1-001 — Captain's Orders omits both `empty` destination constraints

- fighter: `annie-christmas`
- card: `Captain's Orders`
- severity: P1
- expected: Annie is placed in an empty space in her zone, then another friendly fighter is placed in an empty space in Annie's zone.
- observed: Git uses `any_space_in_annie_zone` for both placements without encoding `empty` in either destination selector.
- evidence/reasoning: under the project's placement contract, explicit `empty` changes the selectable destination domain. Without it, an occupied space can remain selectable and merely cause placement failure later, which is not equivalent to the printed rule.
- expected semantic correction: both placement choices must exclude occupied spaces before selection.

### AC-P2-001 — Bottom Dealing public reveal knowledge may be lost

- fighter: `annie-christmas`
- card: `Bottom Dealing`
- severity: P2
- expected: the bottom card is publicly revealed and its final top/bottom disposition remains trackable public information as required by the physical resolution.
- observed: Git represents the subsequent top/bottom destination choice as `visibility: private`.
- evidence/reasoning: the project's `REVEAL` semantics make the card identity public. A private disposition choice can erase information that the physical card resolution exposes to both players, affecting hidden-information boundaries and edge-case knowledge tracking.
- expected semantic correction: preserve public knowledge of the revealed card's final top/bottom destination, or explicitly encode an equivalent public-knowledge result while keeping implementation internals private.

### AC-P3-001 — Captain's Orders evidence mapping ID differs from canonical ID

- fighter: `annie-christmas`
- card: `Captain's Orders`
- severity: P3
- expected: archive/evidence mapping references the existing canonical card definition ID `captains-orders`.
- observed: archive manifest uses `captain-s-orders`.
- evidence/reasoning: this is a non-gameplay mapping inconsistency and does not create a new card definition.
- expected semantic correction: normalize the evidence mapping to `captains-orders`.

### JT-P1-001 — Utility Belt exact movement modeled as `up to`

- fighter: `dr-jill-trent`
- card: `Utility Belt`
- severity: P1
- expected: when the movement mode is selected, Jill moves exactly 1 space if that operation can legally resolve.
- observed: Git uses `MOVE` with `max_distance: 1`, which admits zero-space movement under the project's general `up to N` movement semantics.
- evidence/reasoning: the printed option says `move Jill Trent 1 space`, not `up to 1 space`. Exact and maximum-distance movement have different legal behaviors.
- expected semantic correction: encode exact-one-space movement; if no legal one-space destination exists, the selected operation is impossible rather than voluntarily resolving as zero movement.

### JT-P2-001 — Insightful Deduction revealed-card disposition represented as private

- fighter: `dr-jill-trent`
- card: `Insightful Deduction`
- severity: P2
- expected: after the opponent's top cards are publicly revealed, the chosen bottom card and resulting top order remain public/trackable to the extent exposed by physical resolution.
- observed: Git represents both `bottom-card` and `top-order` choices with private visibility.
- evidence/reasoning: the identities were already publicly revealed. Treating the later disposition/order as wholly private risks losing known-card information and changing hidden-information semantics.
- expected semantic correction: retain public knowledge of which revealed card goes to the bottom and the resulting order of the remaining revealed cards, or explicitly encode equivalent public-information persistence.

### DS-P3-001 — Doctor Strange title punctuation and evidence mapping normalization

- fighter: `doctor-strange`
- card: `No, Really, I'm a Doctor`
- severity: P3
- expected: canonical display metadata preserves printed title punctuation, and evidence mapping uses the canonical card definition ID.
- observed: Git display name is `No Really, I'm a Doctor` without the comma after `No`; archive mapping uses `no-really-i-m-a-doctor` while Git canonical ID is `no-really-im-a-doctor`.
- evidence/reasoning: both differences are non-gameplay identity/presentation issues. Printed semantics, type, BOOST, quantity, and effect behavior match.
- expected semantic correction: preserve the printed title punctuation in display metadata and normalize evidence mapping to the existing canonical ID.

## Corpus-Level Observations

1. The physical evidence corpus for this batch is technically complete: 60 unique card definitions correspond to 60 successfully decoded and inspected images. There were no binary duplicate images and no physical-copy duplication of repeated card quantities.

2. Quantity and fixed-deck construction are internally consistent for all five fighters. Every archive quantity sum is 30 and every canonical available/game deck count is 30. No quantity discrepancy was found.

3. Generic/familiar cards were not skipped. Every unique image, including ordinary shared-style effects, was visually checked rather than inferred from filenames or archive manifests.

4. The archive manifests were treated as evidence mappings rather than authority. This exposed two card-id normalization issues without incorrectly treating the archive IDs as new canonical definitions.

5. Several findings demonstrate why literal prose equivalence is not required but semantic precision is. The material failures are not formatting differences; they concern mandatory versus optional choices, exact versus `up to` movement, selection-domain restrictions, relocation scope, and public-information persistence.

6. Shared integration requirements can still be valid even when a transcription discrepancy exists elsewhere in the same fighter. This batch directly confirms `B-REQ-009`, `C-REQ-012`, and `C-REQ-013` from supplied card images.

7. `B-REQ-005` was not visually confirmed because it belongs to Spider-Man's fighter ability rather than a supplied card component. No claim was made beyond the accessible card evidence.

8. No local evidence blocker remained. Every fighter ZIP, manifest, and card image required for this batch-level card-corpus check was accessible and readable.

## Final Assessment

The batch is **not ready for unconditional acceptance as a canonical Phase 4 card-image corpus** because three fighters contain material normalized-semantics discrepancies.

- `spider-man` is clean and suitable for the corpus as checked.
- `doctor-strange` is gameplay-clean but carries minor P3 identity/presentation normalization findings.
- `she-hulk` requires correction of two P1 gameplay semantics (`Legalese`, `Nerve Cluster Strike`).
- `annie-christmas` requires correction of one P1 placement-domain issue and review/correction of one P2 hidden-information issue, plus one P3 evidence-mapping normalization.
- `dr-jill-trent` requires correction of one P1 exact-movement issue and review/correction of one P2 revealed-information disposition issue.

No files were modified as part of the independent QA itself. This report records findings only; canonical fighter/card manifests were not changed.
