# Phase 4 Card-Image QA Report

## Batch

- buffy
- invisible-man
- jekyll-and-hyde
- sherlock-holmes
- willow

Filename anchor: `buffy` (lexicographically first fighter ID in this batch).

## Verdict

**FAIL**

The batch is technically complete and all 66 unique card images were readable and inspected. Quantity/construction validation passed for all five fighters. Three material P1 semantic discrepancies were found in canonical Git normalization: one for Jekyll & Hyde and two for Willow. Buffy also has one P3 archive-manifest mapping mismatch. No fighter was blocked.

## QA Scope and Evidence Snapshot

This report persists the completed independent Phase 4 Card-Image QA for `unmatched-bundle-03.zip`. No canonical fighter/card manifests were modified by this QA.

Canonical source snapshots used during the completed QA:

- `main`: `f40c3f9970da24ab4e17ff51fc75fb5d7080b239`
- `phase-4b-worker-a-classics`: `840a03572bf8b5912b03f32e5ae2a48cc8fadc2b`

Repository tips were re-checked before report persistence. At persistence time:

- `main`: `ae4a2c3a471638cdd64331bb98add1d4571cc0a0`
- `phase-4b-worker-a-classics`: `c2eddc19e2183cb50dc51ba9ad5528df9b5c8b7d`

The persistence step did not re-run or reinterpret the completed QA; the findings below are the evidence-backed result of that QA session.

## Summary

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| PASS | 1 |
| PASS_WITH_QUALIFICATIONS | 2 |
| FAIL | 2 |
| BLOCKED | 0 |
| Total unique card images | 66 |
| Images successfully inspected | 66 |
| Unreadable images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Duplicate binary images | 0 |
| Quantity failures | 0 |
| Metadata discrepancies | 1 (P3) |
| Semantic discrepancies | 3 (P1) |

All nested fighter ZIPs opened successfully. Every `manifest.yaml` was readable. Every image file was non-zero and successfully decoded. No manifest entry lacked an image and no image lacked a manifest entry. No duplicate binary image hashes were found.

## Fighter Results

### buffy

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/buffy.yaml`
- fighter manifest: `docs/fighters/phase-4b/buffy.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS** for quantities/construction and gameplay card data
- unique-card image completeness: **13/13 PASS**
- printed card identity/content verification: **13/13 inspected; no gameplay-semantic mismatch**
- discrepancies: one P3 archive mapping mismatch for `Slayer's Strength`
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13/13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 35 |
| Git available pool | 35 |
| Git game deck | 30 |
| Unique Git definitions | 13 |
| Construction | `choose_groups` |

The 35-card archive pool is correct. Canonical construction is 25 base cards plus exactly one 5-card sidekick package, producing a 30-card game deck. The archive therefore correctly contains the complete 35-card available pool rather than one selected 30-card match deck.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Cartwheel Kick | PASS | PASS | PASS | PASS |
| Daring Strike | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Insight | PASS | PASS | PASS | PASS |
| Military Knowledge | PASS | PASS | PASS | PASS |
| Mr. Pointy | PASS | PASS | PASS | PASS |
| Rapid Recovery | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Right-hand Man | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Slayer's Strength | PASS | P3 | PASS | PASS_WITH_QUALIFICATION |
| Swift Strike | PASS | PASS | PASS | PASS |
| Training | PASS | PASS | PASS | PASS |

`Slayer's Strength` printed behavior was checked in detail: the optional movement of adjacent fighters and subsequent damage to fighters actually moved agrees with the canonical semantics. The only mismatch is the archive mapping slug described in Findings.

#### Integration qualification

Canonical documentation records `A-REQ-014` for setup-selected sidekick/active-roster handling. Card faces corroborate distinct Giles and Xander card packages, but the setup-selection rule itself is not printed on the action cards, so this QA does not claim that the full setup ability is visually proven by the card-image archive.

---

### invisible-man

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/invisible-man.yaml`
- fighter manifest: `docs/fighters/phase-4b/invisible-man.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **14/14 PASS**
- printed card identity/content verification: **14/14 PASS**
- discrepancies: none
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 14/14 |
| Unique images | 14 |
| Manifest card entries | 14 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 14 |
| Construction | `fixed` |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Coded Notes | PASS | PASS | PASS | PASS |
| Confound | PASS | PASS | PASS | PASS |
| Covert Preparation | PASS | PASS | PASS | PASS |
| Dreaming of Revenge | PASS | PASS | PASS | PASS |
| Emerge From Mist | PASS | PASS | PASS | PASS |
| Impossible to See | PASS | PASS | PASS | PASS |
| Into Thin Air | PASS | PASS | PASS | PASS |
| Lurking | PASS | PASS | PASS | PASS |
| Reign of Terror | PASS | PASS | PASS | PASS |
| Rolling Fog | PASS | PASS | PASS | PASS |
| Slip Away | PASS | PASS | PASS | PASS |
| Step Lightly | PASS | PASS | PASS | PASS |
| Surprise Attack | PASS | PASS | PASS | PASS |
| Vanish | PASS | PASS | PASS | PASS |

Semantically sensitive checks included:

- `Confound`: opponent's discard is optional; fog relocation is the alternative when the opponent does not discard.
- `Covert Preparation`: owner moves one fog, then opponent moves a different fog.
- `Step Lightly`: damage becomes 3 instead of 1 while Invisible Man is on fog.
- `Slip Away`: fog relocation resolves first and Invisible Man is then placed in that fog's destination space.
- `Vanish`: temporary undefeated off-board state, not defeat, followed by scheduled return at the start of the next turn.
- `Impossible to See`: opposing combat-card value becomes 0 and card effects cannot subsequently change it during that combat.

#### Integration requirements confirmed

- `A-REQ-003`: fog tokens are persistent positioned battlefield components manipulated by multiple cards.
- `A-REQ-005`: `Vanish` establishes temporary undefeated off-board fighter presence with delayed return.

These are integration requirements, not transcription discrepancies.

---

### jekyll-and-hyde

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/jekyll-and-hyde.yaml`
- fighter manifest: `docs/fighters/phase-4b/jekyll-and-hyde.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL** due to one semantic ordering discrepancy
- unique-card image completeness: **13/13 PASS**
- printed card identity/content verification: **12 PASS / 1 FAIL**
- discrepancies: `JH-001` (P1)
- verdict: **FAIL**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13/13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 13 |
| Construction | `fixed` |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Calming Research | PASS | PASS | PASS | PASS |
| Distracted Triage | PASS | PASS | PASS | PASS |
| Duality of Man | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Forever Hyde | PASS | PASS | PASS | PASS |
| Madness Relents | PASS | PASS | PASS | PASS |
| Pure Evil | PASS | PASS | PASS | PASS |
| Recoiling Blow | PASS | PASS | PASS | PASS |
| Scientific Method | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Strange Case | PASS | PASS | FAIL | FAIL |
| Succumb to Compulsion | PASS | PASS | PASS | PASS |
| With Haste! | PASS | PASS | PASS | PASS |

No quantity/type/value/BOOST/restriction mismatch was found for `Strange Case`; the failure is specifically the hidden-information/decision ordering described in Findings.

---

### sherlock-holmes

- branch: `main`
- deck manifest: `docs/cards/phase-4a/sherlock-holmes.yaml`
- fighter manifest: `docs/fighters/phase-4a/sherlock-holmes.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **13/13 PASS**
- printed card identity/content verification: **13/13 PASS**
- discrepancies: none
- verdict: **PASS**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13/13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 13 |
| Construction | `fixed` |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Administer Aid | PASS | PASS | PASS | PASS |
| Confirm Suspicion | PASS | PASS | PASS | PASS |
| Counterpunch | PASS | PASS | PASS | PASS |
| Deduce Strategy | PASS | PASS | PASS | PASS |
| Education Never Ends | PASS | PASS | PASS | PASS |
| Elementary | PASS | PASS | PASS | PASS |
| Eliminate the Impossible | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Fixed Point in a Changing Age | PASS | PASS | PASS | PASS |
| Master of Disguise | PASS | PASS | PASS | PASS |
| Service Revolver | PASS | PASS | PASS | PASS |
| Study Methods | PASS | PASS | PASS | PASS |
| The Game is Afoot | PASS | PASS | PASS | PASS |

Sensitive Holmes-specific checks included:

- `Deduce Strategy` uses opposing printed combat-card value semantics derived from its BOOST, rather than a generic current-value overwrite.
- `Elementary` has source-specific face-up commitment/prediction timing and, on success, cancels effects while ignoring the opposing attack value.
- Holmes/Watson-only banners and ANY banners agree with canonical `usable_by` restrictions.
- Watson-only cards agree with fighter topology.

---

### willow

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/willow.yaml`
- fighter manifest: `docs/fighters/phase-4b/willow.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL** due to two semantic discrepancies
- unique-card image completeness: **13/13 PASS**
- printed card identity/content verification: **11 PASS / 2 FAIL**
- discrepancies: `WILLOW-001`, `WILLOW-002` (both P1)
- verdict: **FAIL**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13/13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 13 |
| Construction | `fixed` |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Black Magic | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Flayed Alive | PASS | PASS | PASS | PASS |
| Hacker | PASS | PASS | PASS | PASS |
| Knowledge of the Craft | PASS | PASS | PASS | PASS |
| Love & Loss | PASS | PASS | PASS | PASS |
| Meditation | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Rending Shot | PASS | PASS | PASS | PASS |
| Resurrect | PASS | PASS | FAIL | FAIL |
| Revoke | PASS | PASS | PASS | PASS |
| Swift Strike | PASS | PASS | PASS | PASS |
| When Good Magic Fails | PASS | PASS | FAIL | FAIL |

#### Integration requirements confirmed

- `A-REQ-006`: `Resurrect` directly establishes return of a defeated fighter with health set to exactly 3 rather than recovery by 3.

The requirement itself is confirmed even though the current canonical target-domain normalization has the separate P1 issue below.

## Findings

### JH-001 — `Strange Case` target selection occurs before reveal

- fighter: `jekyll-and-hyde`
- card: `Strange Case`
- severity: **P1**
- Git location: `docs/cards/phase-4b/jekyll-and-hyde.yaml`, `strange-case`
- expected: reveal the top card first; once its BOOST is known, select one adjacent fighter, deal damage equal to the revealed BOOST, then put the revealed card into hand.
- observed: canonical normalization selects/binds the adjacent fighter before `REVEAL`, then reveals the top card and uses its BOOST for damage.
- evidence/reasoning: printed resolution order establishes the reveal before the damage instruction. Choosing the target before the reveal changes the hidden-information boundary and can change the legal strategic decision when multiple adjacent fighters are available.
- expected semantic correction: preserve `REVEAL -> choose adjacent fighter -> damage using bound revealed BOOST -> move revealed card to hand`.

### WILLOW-001 — `Resurrect` target domain is too narrow

- fighter: `willow`
- card: `Resurrect`
- severity: **P1**
- Git location: `docs/cards/phase-4b/willow.yaml`, `resurrect`
- expected: select a **friendly fighter who has been defeated**, return/place that fighter in Willow's zone, and set its health to exactly 3.
- observed: canonical selector is restricted to a defeated fighter with `owner: self`.
- evidence/reasoning: the printed card says `friendly fighter`, which is broader than ownership by Willow's player in team play; teammate fighters are friendly. The Git selector therefore removes legal teammate targets.
- expected semantic correction: target domain should be defeated + friendly rather than defeated + self-owned, while preserving Willow-zone return and exact health 3.

### WILLOW-002 — `When Good Magic Fails` normalizes MOVE as PLACE

- fighter: `willow`
- card: `When Good Magic Fails`
- severity: **P1**
- Git location: `docs/cards/phase-4b/willow.yaml`, `when-good-magic-fails`
- expected: while Dark Willow, **move Willow to any space in her zone**, then discard the top card of the deck.
- observed: canonical normalization uses `PLACE` for Willow's relocation, followed by the discard.
- evidence/reasoning: MOVE and PLACE are distinct gameplay semantics. `PLACE` may bypass movement/path legality and movement-specific interactions. The printed card explicitly uses movement semantics.
- expected semantic correction: represent a source-defined MOVE to the selected space in Willow's zone, followed by the mandatory top-card discard.

### BUFFY-001 — archive card ID mismatch for `Slayer's Strength`

- fighter: `buffy`
- card: `Slayer's Strength`
- severity: **P3**
- expected: archive evidence mapping should use canonical card ID `slayers-strength`.
- observed: archive `manifest.yaml` uses `slayer-s-strength` for the corresponding image.
- evidence/reasoning: printed name, quantity, restriction, type, value, BOOST, and effect semantics all map unambiguously to canonical `Slayer's Strength`; only the evidence-manifest slug differs.
- expected correction: normalize the evidence mapping to canonical `slayers-strength`. No canonical gameplay-manifest correction is implied by this finding.

## Corpus-Level Observations

1. **Outer and nested archive integrity is clean.** Five fighter ZIPs were present, all opened correctly, and all fighters were unambiguously mapped to canonical IDs.
2. **The image corpus is one-image-per-unique-definition.** Across 66 unique card definitions, there were 66 readable images, no binary duplicates, no repeated physical-copy images, no missing images, and no extras.
3. **All deck quantity checks passed.** Four fighters use fixed 30-card construction. Buffy's non-standard 35-card available pool / 30-card match deck was correctly represented and was not incorrectly treated as a quantity failure.
4. **Archive manifests were treated only as evidence mappings.** Printed card images and canonical Git manifests were used to resolve identity/content; archive metadata was not treated as authoritative when evaluating semantics.
5. **Printed prose was compared semantically rather than by literal string equality.** Structured stages, choices, bound results, conditions, and integration requirements were accepted when they preserved the legal printed behavior.
6. **Shared engine requirements were not misclassified as transcription failures.** Invisible Man's fog/off-board mechanics, Buffy's setup roster qualification, and Willow's exact-health return requirement remain integration qualifications where applicable.
7. **Generic cards were also checked.** Feint, Regroup, Skirmish, Swift Strike and other familiar definitions were visually inspected rather than assumed correct from familiarity.
8. **Failures did not stop subsequent checking.** All 66 cards were inspected despite material findings in Jekyll & Hyde and Willow.

## Final Assessment

The physical evidence batch is **technically complete and suitable as a card-image evidence corpus**, but the canonical normalized corpus is **not fully reconciled with that evidence** because three P1 gameplay-semantic mismatches remain:

- `jekyll-and-hyde / Strange Case`: hidden-information/order mismatch;
- `willow / Resurrect`: incorrect fighter target domain;
- `willow / When Good Magic Fails`: MOVE incorrectly represented as PLACE.

Buffy additionally has one non-gameplay P3 evidence-manifest ID mismatch. Invisible Man and Sherlock Holmes contain no material card-data discrepancies; Invisible Man retains documented shared-integration qualifications. No real blockers or unreadable evidence remain.

**Batch verdict: FAIL.**
