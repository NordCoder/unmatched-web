# Phase 4 Card-Image QA Report

## Batch

- achilles
- bloody-mary
- daredevil
- sun-wukong
- yennenga

## Verdict

FAIL

## Summary

Independent read-only Phase 4 card-image QA was completed for one external batch containing five fighter ZIP archives. All nested ZIP files were readable, all manifests were readable, all expected image files decoded successfully, and every unique card image in the batch was visually inspected against the canonical Git documentation state for the fighter's mapped branch.

Batch integrity was clean: 57/57 unique card images were inspected; there were no unreadable, missing, extra, zero-byte, corrupted, or binary-duplicate images. Quantity validation passed for all five fighters, including Daredevil's non-standard fixed 22-card construction. Two fighters were clean (`achilles`, `bloody-mary`), while material discrepancies were found for `sun-wukong`, `yennenga`, and `daredevil`.

Canonical branch tips used during QA and rechecked before finalization:

- `main`: `f40c3f9970da24ab4e17ff51fc75fb5d7080b239`
- `phase-4b-worker-a-classics`: `840a03572bf8b5912b03f32e5ae2a48cc8fadc2b`
- `phase-4b-worker-b-licensed`: `f383a639e980d9e753c742b20b7e085d3163502f`

Batch totals:

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| PASS | 2 |
| PASS_WITH_QUALIFICATIONS | 0 |
| FAIL | 3 |
| BLOCKED | 0 |
| Total unique card images | 57 |
| Total images successfully inspected | 57 |
| Unreadable images | 0 |
| Missing images | 0 |
| Duplicate images | 0 |
| Quantity failures | 0 |

Material discrepancy totals:

- P1 gameplay metadata: 5
- P1 semantic: 1
- P2 semantic: 1
- P3 evidence-mapping metadata: 2
- Total P1/P2 findings: 7

## Fighter Results

### achilles

- branch: `main`
- deck manifest: `docs/cards/phase-4a/achilles.yaml`
- fighter manifest: `docs/fighters/phase-4a/achilles.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: PASS, 12/12
- printed card identity/content verification: PASS, all 12 unique cards checked
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 12
- unique Git card definitions: 12
- discrepancies: none
- verdict: PASS

`Relentless Assault` was correctly treated as an external bonus-attack definition rather than an additional physical action-card instance.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Brothers in Arms | PASS | PASS | PASS | PASS |
| Battle Frenzy | PASS | PASS | PASS | PASS |
| The Day of Your Doom | PASS | PASS | PASS | PASS |
| Test for Weakness | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Wily Fighting | PASS | PASS | PASS | PASS |
| Blessed by Hermes | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Battle Hardened | PASS | PASS | PASS | PASS |
| Achilles' Heel | PASS | PASS | PASS | PASS |
| Under Achilles' Helm | PASS | PASS | PASS | PASS |
| Spear Throw | PASS | PASS | PASS | PASS |

### bloody-mary

- branch: `main`
- deck manifest: `docs/cards/phase-4a/bloody-mary.yaml`
- fighter manifest: `docs/fighters/phase-4a/bloody-mary.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: PASS, 13/13
- printed card identity/content verification: PASS, all 13 unique cards checked
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 13
- unique Git card definitions: 13
- discrepancies: none
- verdict: PASS

`Bloody Reprise` was correctly treated as an external bonus attack rather than an additional action-card instance. Third-action checks and the historical `turn_start_space` dependency used by Jump Scare were consistent with the fighter manifest's persisted state model.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Bloody Requiem | PASS | PASS | PASS | PASS |
| Speak Three Times | PASS | PASS | PASS | PASS |
| Ghostly Touch | PASS | PASS | PASS | PASS |
| Out of the Mirror | PASS | PASS | PASS | PASS |
| Evade | PASS | PASS | PASS | PASS |
| Mirror Image | PASS | PASS | PASS | PASS |
| Closer Than She Appears | PASS | PASS | PASS | PASS |
| Stolen Memories | PASS | PASS | PASS | PASS |
| Broken Glass | PASS | PASS | PASS | PASS |
| Trick of the Light | PASS | PASS | PASS | PASS |
| Infinity Mirror | PASS | PASS | PASS | PASS |
| Jump Scare | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |

### sun-wukong

- branch: `main`
- deck manifest: `docs/cards/phase-4a/sun-wukong.yaml`
- fighter manifest: `docs/fighters/phase-4a/sun-wukong.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL due to one semantic-scope discrepancy
- unique-card image completeness: PASS, 12/12
- printed card identity/content verification: completed for all 12 unique cards
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 12
- unique Git card definitions: 12
- discrepancies: SW-P2-001, SW-P3-001
- verdict: FAIL

`Tricked You` was correctly treated as an external bonus attack rather than a physical action-card definition.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Infinite Strikes | PASS | PASS | PASS | PASS |
| Ruyi Jingo Bang | PASS | PASS | PASS | PASS |
| Taunting Laughter | PASS | PASS | PASS | PASS |
| Ox Form | PASS | PASS | PASS | PASS |
| 72 Transformations | PASS | P3 | PASS | QUALIFIED |
| Sly Monkey | PASS | PASS | PASS | PASS |
| Wily Fighting | PASS | PASS | PASS | PASS |
| Bewilderment | PASS | PASS | P2 | FAIL |
| Golden Chain Mail | PASS | PASS | PASS | PASS |
| Tortoise Form | PASS | PASS | PASS | PASS |
| Fiery Eyes That See | PASS | PASS | PASS | PASS |
| Phoenix Form | PASS | PASS | PASS | PASS |

#### Integration requirements confirmed

None.

### yennenga

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/yennenga.yaml`
- fighter manifest: `docs/fighters/phase-4b/yennenga.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL due to three fighter-restriction discrepancies
- unique-card image completeness: PASS, 12/12
- printed card identity/content verification: completed for all 12 unique cards
- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 12
- unique Git card definitions: 12
- discrepancies: YEN-P1-001, YEN-P1-002, YEN-P1-003
- verdict: FAIL

`Volley` was correctly treated as an external bonus-attack definition rather than a physical action card.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Rain of Arrows | PASS | PASS | PASS | PASS |
| Surprise Volley | PASS | FAIL | PASS | FAIL |
| Shield Formation | PASS | PASS | PASS | PASS |
| Master of the Hunt | PASS | PASS | PASS | PASS |
| One With the Land | PASS | PASS | PASS | PASS |
| Stallion Charge | PASS | PASS | PASS | PASS |
| Jaws of the Beast | PASS | PASS | PASS | PASS |
| Point Blank | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Divide and Conquer | PASS | FAIL | PASS | FAIL |
| Pin the Prey | PASS | FAIL | PASS | FAIL |
| Skirmish | PASS | PASS | PASS | PASS |

#### Integration requirements confirmed

- `A-REQ-011`: directly confirmed by `Stallion Charge`; movement resolution must retain traversed-space/fighter information for follow-up damage.
- `A-REQ-012`: directly confirmed by `Surprise Volley`; a defeated Archer may return and replace the current attacker within an already-running combat.

`A-REQ-010` was not counted as card-image-confirmed because it belongs to Yennenga's fighter ability rather than a card component present in this evidence ZIP.

### daredevil

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/daredevil.yaml`
- fighter manifest: `docs/fighters/phase-4b/daredevil.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL due to two card-type errors and one movement-semantic error
- unique-card image completeness: PASS, 8/8
- printed card identity/content verification: completed for all 8 unique cards
- archive quantity sum: 22
- Git available pool: 22
- Git game deck: 22
- unique archive cards: 8
- unique Git card definitions: 8
- discrepancies: DD-P1-001, DD-P1-002, DD-P1-003, DD-P3-001
- verdict: FAIL

Daredevil's non-standard fixed 22-card deck construction was handled correctly; the QA did not assume a universal 30-card deck size.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Through Adversity | PASS | PASS | PASS | PASS |
| Man Without Fear | PASS | PASS | PASS | PASS |
| Son of a Boxer | PASS | FAIL | PASS | FAIL |
| Feint | PASS | FAIL | PASS | FAIL |
| Take a Knee | PASS | PASS | PASS | PASS |
| Grappling Hook | PASS | PASS | FAIL | FAIL |
| Breather | PASS | PASS | PASS | PASS |
| Devil of Hell's Kitchen | PASS | P3 | PASS | QUALIFIED |

#### Integration requirements confirmed

None.

## Findings

### SW-P2-001 — Bewilderment damage-prevention scope

- fighter: `sun-wukong`
- card: `Bewilderment`
- severity: P2
- expected: prevention applies specifically to combat damage as printed on the physical card.
- observed: Git represents the effect as `PREVENT_DAMAGE` with `scope: all_damage_during_combat_window`.
- evidence/reasoning: the normalized project effect model distinguishes combat damage from non-combat effect damage. A broad `all_damage_during_combat_window` scope can suppress non-combat damage merely because it occurs during combat timing, which is broader than the printed behavior. The neighboring `Golden Chain Mail` representation already identifies combat damage specifically.
- expected semantic correction: scope prevention to the combat-damage application event / combat damage only, not arbitrary damage occurring in the combat window.

### SW-P3-001 — 72 Transformations evidence ID mismatch

- fighter: `sun-wukong`
- card: `72 Transformations`
- severity: P3
- expected: archive evidence mapping uses canonical Git card ID `seventy-two-transformations`.
- observed: archive manifest uses `72-transformations`.
- evidence/reasoning: card image, printed name, quantity, metadata and semantics map unambiguously to the canonical definition; gameplay is unaffected.
- expected correction: normalize the archive/evidence card ID. No Git gameplay correction is implied.

### YEN-P1-001 — Surprise Volley fighter restriction

- fighter: `yennenga`
- card: `Surprise Volley`
- severity: P1
- expected: printed restriction is ANY, allowing either Yennenga or an Archer to use the card.
- observed: Git uses `usable_by: [yennenga]`.
- evidence/reasoning: the physical card's restriction icon permits any friendly fighter in the deck, while the Git representation excludes both Archer sidekicks.
- expected semantic correction: represent Surprise Volley as usable by any fighter in Yennenga's deck.

### YEN-P1-002 — Divide and Conquer fighter restriction

- fighter: `yennenga`
- card: `Divide and Conquer`
- severity: P1
- expected: printed restriction is ARCHER.
- observed: Git uses `usable_by: any`.
- evidence/reasoning: Git permits Yennenga herself to commit a combat card printed for Archer only.
- expected semantic correction: restrict the card to `archer`.

### YEN-P1-003 — Pin the Prey fighter restriction

- fighter: `yennenga`
- card: `Pin the Prey`
- severity: P1
- expected: printed restriction is ARCHER.
- observed: Git uses `usable_by: any`.
- evidence/reasoning: Git permits Yennenga herself to use an Archer-only card.
- expected semantic correction: restrict the card to `archer`.

### DD-P1-001 — Son of a Boxer card type

- fighter: `daredevil`
- card: `Son of a Boxer`
- severity: P1
- expected: printed card type is DEFENSE, printed value 3.
- observed: Git uses `type: versatile`, printed value 3.
- evidence/reasoning: Versatile illegally permits the card to be committed while attacking, changing legal gameplay behavior.
- expected semantic correction: preserve the printed Defense-only type.

### DD-P1-002 — Feint card type

- fighter: `daredevil`
- card: `Feint`
- severity: P1
- expected: printed card type is VERSATILE, printed value 2.
- observed: Git uses `type: defense`.
- evidence/reasoning: Git incorrectly removes the printed ability to play this Feint while attacking.
- expected semantic correction: represent the card as Versatile.

### DD-P1-003 — Grappling Hook MOVE versus PLACE

- fighter: `daredevil`
- card: `Grappling Hook`
- severity: P1
- expected: after combat, Daredevil may MOVE up to 2 spaces under normal movement semantics.
- observed: Git allows choosing any space in Daredevil's zone and performs `PLACE`.
- evidence/reasoning: MOVE and PLACE are distinct project primitives. MOVE is path/distance constrained; PLACE is not. The Git representation therefore changes both the relocation mechanism and the set of legal destinations.
- expected semantic correction: represent a MOVE of Daredevil with maximum distance 2 rather than zone-based placement.

### DD-P3-001 — Devil of Hell's Kitchen evidence ID mismatch

- fighter: `daredevil`
- card: `Devil of Hell's Kitchen`
- severity: P3
- expected: archive evidence mapping uses canonical Git card ID `devil-of-hells-kitchen`.
- observed: archive manifest uses `devil-of-hell-s-kitchen`.
- evidence/reasoning: image identity, quantity, printed card metadata and gameplay semantics map unambiguously; gameplay is unaffected.
- expected correction: normalize the archive/evidence card ID. No Git gameplay correction is implied.

## Corpus-Level Observations

1. Archive integrity was strong across the entire batch: all five nested ZIP files were readable, every image decoded, and no zero-byte, corrupted, missing, extra, or binary-duplicate card images were found.
2. The batch respected the one-image-per-unique-definition model. Physical copy count was represented in manifests rather than by duplicate image files.
3. Quantity validation passed for every fighter. The QA explicitly handled non-standard construction and did not assume that all decks contain 30 cards; Daredevil correctly validated at 22 cards.
4. External gameplay definitions were not counted as physical action-card copies. This was relevant for Achilles (`Relentless Assault`), Bloody Mary (`Bloody Reprise`), Sun Wukong (`Tricked You`), and Yennenga (`Volley`).
5. All generic/repeated card identities such as Feint, Skirmish and Wily Fighting were still visually inspected rather than inferred from familiarity.
6. The main material error classes in this batch were fighter restrictions, card type normalization, and relocation semantics. These are exactly the classes where relying on filenames or approximate prose equivalence would have missed gameplay-changing differences.
7. Shared integration requirements can be confirmed by physical cards without making the card transcription itself fail. For Yennenga, `A-REQ-011` and `A-REQ-012` were directly supported by card evidence.
8. Fighter ability behavior was not claimed as image-confirmed when the evidence archive contained only card components. In particular, Yennenga's `A-REQ-010` ability requirement was not promoted to card-image evidence.

## Final Assessment

The batch is technically complete and usable as evidence: all 57 unique card images are readable, all manifests and quantities are structurally coherent, and no archive-level blockers exist.

However, the batch cannot be considered fully reconciled with the canonical repository corpus because seven material P1/P2 discrepancies were found across three fighters:

- `sun-wukong`: one P2 semantic-scope discrepancy on `Bewilderment`;
- `yennenga`: three P1 fighter-restriction discrepancies;
- `daredevil`: two P1 card-type discrepancies and one P1 MOVE-versus-PLACE semantic discrepancy.

`achilles` and `bloody-mary` are clean against the inspected card evidence. `sun-wukong`, `yennenga`, and `daredevil` require canonical documentation corrections before the batch can receive an overall PASS.

Final batch verdict: **FAIL**.
