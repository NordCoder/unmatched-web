# Phase 4 Card-Image QA Report

## Batch

- angel
- beowulf
- deadpool
- little-red-riding-hood
- spike

Source evidence archive: `unmatched-bundle-04.zip`.

Canonical repository: `NordCoder/unmatched-web`.

Canonical branch used for this batch: `phase-4b-worker-a-classics`.

Branch tip observed immediately before persistence: `c2eddc19e2183cb50dc51ba9ad5528df9b5c8b7d`.

QA mode: independent, read-only with respect to canonical fighter/card data. No manifests, schemas, mechanics, sources, or other QA reports were modified.

## Verdict

**FAIL**

## Summary

The batch contained five readable fighter ZIPs: Spike, Angel, Little Red Riding Hood, Beowulf, and Deadpool. All nested ZIPs opened successfully, all 78 supplied card images decoded successfully, and no zero-byte files, corrupted images, binary duplicates, orphan images, or manifest entries without files were found.

Four fighters had complete image coverage relative to their archive manifests. Little Red Riding Hood had a corpus-completeness defect: the archive collapsed the three canonical `Into the Woods` printed-symbol variants into one wolf-symbol image/manifest entry with quantity 3, while Git correctly models three distinct quantity-1 definitions (`wolf`, `pelt`, `knife`). Two canonical physical variants therefore lacked image evidence.

All five decks are fixed 30-card constructions in canonical Git. Aggregate archive quantity sums were 30 for every fighter. The batch contains 78 supplied unique images versus 80 canonical card definitions because of the two missing Little Red `Into the Woods` variants.

All 78 supplied unique card images were visually inspected. The QA found 12 material findings: **11 P1** gameplay-changing discrepancies and **1 P2** evidence/completeness issue. It also found **9 P3** archive-ID or non-gameplay printed/social preservation issues. Beowulf was the only clean fighter.

Batch totals:

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully processed | 5 |
| Fighters with complete canonical image coverage | 4 |
| PASS | 1 |
| PASS WITH FINDINGS | 0 |
| FAIL | 4 |
| BLOCKED | 0 |
| Supplied unique card images | 78 |
| Images successfully inspected | 78 |
| Canonical unique definitions represented by batch fighters | 80 |
| Unreadable images | 0 |
| Technical missing images | 0 |
| Missing canonical variant images | 2 |
| Duplicate binary images | 0 |
| Quantity failures | 1 fighter |
| P1 findings | 11 |
| P2 findings | 1 |
| P3 findings | 9 |
| Gameplay metadata discrepancies | 5 |
| Gameplay semantic discrepancies | 6 |

## Fighter Results

### angel

Canonical deck manifest: `docs/cards/phase-4b/angel.yaml`  
Canonical fighter manifest: `docs/fighters/phase-4b/angel.yaml`

- archive integrity: **PASS**
- canonical manifest comparison: **FAIL** due to printed metadata and semantics mismatches
- unique-card image completeness: **PASS**, 12/12 supplied definitions represented
- printed card identity/content verification: **FAIL**
- discrepancies: 4 material P1 findings
- verdict: **FAIL**

Archive/count evidence:

| Check | Result |
|---|---|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12/12 |
| Unique images | 12 |
| Manifest card entries | 12 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 12 |
| Construction | fixed |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Angelus, Scourge of Europe | PASS | PASS | PASS | PASS |
| Five by Five | PASS | PASS | PASS | PASS |
| Disengage | PASS | PASS | PASS | PASS_WITH_ERRATA_CONTEXT |
| Cursed with a Soul | PASS | FAIL | PASS | FAIL |
| Wisdom of Ages | PASS | PASS | FAIL | FAIL |
| Brooding | PASS | FAIL | PASS | FAIL |
| The Rogue Slayer | PASS | FAIL | PASS | FAIL |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Haunted by the Faces | PASS | PASS | PASS | PASS |
| Killer of the Dead | PASS | PASS | PASS | PASS |

`Disengage` visibly uses “Move”, while canonical Git represents the corrected effect as `PLACE` and explicitly cites a card-specific `disengage_errata` ruling source. That difference was therefore not treated as a transcription discrepancy.

Integration requirements confirmed: none.

### beowulf

Canonical deck manifest: `docs/cards/phase-4b/beowulf.yaml`  
Canonical fighter manifest: `docs/fighters/phase-4b/beowulf.yaml`

- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS**, 12/12
- printed card identity/content verification: **PASS**
- discrepancies: none
- verdict: **PASS**

Archive/count evidence:

| Check | Result |
|---|---|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12/12 |
| Unique images | 12 |
| Manifest card entries | 12 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 12 |
| Construction | fixed |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Remnant of Valor | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Epic Poem | PASS | PASS | PASS | PASS |
| The Ancient Heirloom | PASS | PASS | PASS | PASS |
| No Contest Expecteth | PASS | PASS | PASS | PASS |
| The War-King | PASS | PASS | PASS | PASS |
| The Equal of Grendel | PASS | PASS | PASS | PASS |
| Vigor and Courage | PASS | PASS | PASS | PASS |
| Golden Drinking Horn | PASS | PASS | PASS | PASS |
| Hot for the Battle | PASS | PASS | PASS | PASS |
| Fatal Struggle | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |

The canonical `the-ancient-heirloom` interpretation is compatible with the printed image. Git records printed-order handling as a high-confidence derived interpretation rather than presenting it as an official printed fact.

Integration requirements confirmed: none.

### deadpool

Canonical deck manifest: `docs/cards/phase-4b/deadpool.yaml`  
Canonical fighter manifest: `docs/fighters/phase-4b/deadpool.yaml`

- archive integrity: **PASS**
- canonical manifest comparison: **FAIL** due to printed metadata/semantic mismatches
- unique-card image completeness: **PASS**, 30/30
- printed card identity/content verification: **FAIL**
- discrepancies: 4 material P1 findings plus 6 P3 findings
- verdict: **FAIL**

Canonical Git separately marks Deadpool `status: blocked`, `policy: blocked`, and `integration: requires_shared_extension` because physical/social card behavior has no settled digital-adaptation policy. That repository status did not prevent independent image QA and was not itself treated as a transcription failure.

Archive/count evidence:

| Check | Result |
|---|---|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 30/30 |
| Unique images | 30 |
| Manifest card entries | 30 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 30 |
| Unique Git definitions | 30 |
| Construction | fixed |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Chimichanga Break! | PASS | PASS | PASS | PASS |
| Eat Me | PASS | PASS | PASS | PASS |
| And For My Next Move... | PASS | PASS | PASS | PASS |
| 3 of Hearts | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| Gimme Gimme Chimichanga | PASS | PASS | PASS | PASS |
| Dumpster Divin' Deadpool | PASS | PASS | FAIL | FAIL |
| Non-Retinal Scan Access to Danger Room | PASS | PASS | PASS | PASS |
| Time out time out time out! | PASS | PASS | P3 social-instruction omission | PASS_WITH_P3 |
| Holy Mackerel! | PASS | PASS | P3 social-instruction omission | PASS_WITH_P3 |
| Sweeet! | PASS | PASS | FAIL | FAIL |
| Wanna bet? | PASS | PASS | PASS | PASS |
| Rob's Pouch & Shoe Emporium | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| Push to Teleport | PASS | PASS | PASS | PASS |
| Cha-Ching! | PASS | PASS | PASS | PASS |
| I Always Get The Last Word | PASS | PASS | PASS | PASS |
| Klunkin' Heads | PASS | PASS | PASS | PASS |
| Deadpool™ Merc For Hire, LLC | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| Exploding Card! | PASS | PASS | PASS | PASS |
| Passwords | PASS | PASS | PASS | PASS |
| Gaze of Stone | PASS | FAIL | PASS/qualified | FAIL |
| Excuse me while I grow some limbs. | PASS | PASS | PASS | PASS |
| Transit Card | PASS | PASS | FAIL | FAIL |
| Faint | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Super Feint | PASS | PASS | PASS | PASS |
| Call Me | PASS | PASS | PASS | PASS |
| I'm Not Wearing Pants | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| They Have An Amazing Buffet | PASS | PASS | PASS | PASS |
| Underrated Super Heroes | PASS | PASS | PASS | PASS |
| Xavier Institute Faculty | PASS | PASS | PASS | PASS |

Integration requirements confirmed:

- `A-REQ-008`: Xavier Institute Faculty directly confirms card-specific melee/ranged combat legality.
- `A-REQ-015`: multiple cards directly confirm dependence on real names, food, spoken/noise actions, ownership, mirrors, sleeves, writing on a physical card, clothing, wagers, and subjective board colour.

### little-red-riding-hood

Canonical deck manifest: `docs/cards/phase-4b/little-red-riding-hood.yaml`  
Canonical fighter manifest: `docs/fighters/phase-4b/little-red-riding-hood.yaml`

- archive integrity: **FAIL for canonical corpus completeness**, although the ZIP/files themselves are technically healthy
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **FAIL**: archive supplies 12 unique images while Git has 14 definitions; two `Into the Woods` variants lack image evidence
- printed card identity/content verification: **FAIL**
- discrepancies: 1 P1, 1 P2, 2 P3
- verdict: **FAIL**

Archive/count evidence:

| Check | Result |
|---|---|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12/12 |
| Unique images | 12 |
| Manifest card entries | 12 |
| Duplicate images | 0 |
| Technical manifest→image missing | 0 |
| Canonical physical variant images missing | 2 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 14 |
| Construction | fixed |

Git models `Into the Woods` as three distinct quantity-1 definitions with printed symbols `wolf`, `pelt`, and `knife`. The archive instead maps one wolf-symbol image as a single quantity-3 entry. Aggregate deck quantity remains 30, but physical definition coverage is incomplete.

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| What Large Hands You Have | PASS | PASS | PASS | PASS |
| The Wolf's Skin | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| What a Terrible Big Mouth You Have | PASS | PASS | PASS | PASS |
| What Big Ears You Have | PASS | PASS | PASS | PASS |
| What Big Eyes You Have | PASS | PASS | PASS | PASS |
| What's That in My Basket? | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| Never Leave the Path | PASS | PASS | PASS | PASS |
| Into the Woods — wolf | PASS | FAIL collapsed mapping | PASS | FAIL/P2 |
| Into the Woods — pelt | MISSING | — | — | UNVERIFIED |
| Into the Woods — knife | MISSING | — | — | UNVERIFIED |
| Long Have I Sought You | PASS | PASS | PASS | PASS |
| A Grimm Tale | PASS | PASS | FAIL | FAIL |
| Stones in the Belly | PASS | PASS | PASS | PASS |
| Once Upon a Time | PASS | PASS | PASS | PASS |

Integration requirements confirmed:

- `A-REQ-008`: `What Big Ears You Have` confirms a conditional defense mode distinct from immutable printed attack type.
- `A-REQ-009`: supplied cards directly confirm Basket symbols and wild-symbol-dependent card behavior. The separate Basket reference/setup component was not included in this card-image ZIP, so that auxiliary physical component itself was not visually certified by this QA.

### spike

Canonical deck manifest: `docs/cards/phase-4b/spike.yaml`  
Canonical fighter manifest: `docs/fighters/phase-4b/spike.yaml`

- archive integrity: **PASS**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **PASS**, 12/12
- printed card identity/content verification: **FAIL**
- discrepancies: 2 P1 plus 1 P3
- verdict: **FAIL**

Archive/count evidence:

| Check | Result |
|---|---|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12/12 |
| Unique images | 12 |
| Manifest card entries | 12 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 12 |
| Construction | fixed |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Leap Away | PASS | PASS | PASS | PASS |
| The Rush | PASS | PASS | PASS | PASS |
| The Sight | PASS | FAIL | PASS | FAIL |
| Let's Dance | PASS | P3 archive ID mismatch | PASS | PASS_WITH_P3 |
| Seek the Shadows | PASS | PASS | PASS | PASS |
| Empathy | PASS | PASS | PASS | PASS |
| Always Surprising | PASS | PASS | PASS | PASS |
| Arrogance | PASS | PASS | FAIL | FAIL |
| Bloody Hell! | PASS | PASS | PASS | PASS |

Integration requirements confirmed:

- `A-REQ-003`: Shadow tokens/components and Shadow-sensitive predicates are directly confirmed by the card images.
- `A-REQ-007`: Always Surprising directly confirms blind BOOST plus a conditional transform of the resolved BOOST amount.

## Findings

### B04-SPI-001 — The Sight printed BOOST mismatch

- fighter: `spike`
- card: `The Sight`
- expected: canonical immutable metadata should record printed BOOST **2**
- observed: supplied physical image shows BOOST **2**; Git `docs/cards/phase-4b/spike.yaml` records `boost: 3`
- severity: **P1**
- evidence/reasoning: printed BOOST is direct component metadata and affects maneuver BOOST and any effect reading the card's BOOST value

### B04-SPI-002 — Arrogance effect semantics are incorrect

- fighter: `spike`
- card: `Arrogance`
- expected: printed DURING COMBAT behavior permits discarding Spike's entire hand; if performed, this card's value becomes **6 instead**
- observed: Git models selection of a hand card for `BOOST`, followed by moving/reordering the remaining hand onto the bottom of the deck
- severity: **P1**
- evidence/reasoning: the Git representation changes the cost/action, card destinations, resulting value calculation, and replacement semantics; this is not an equivalent normalization

### B04-SPI-003 — Let's Dance archive ID mismatch

- fighter: `spike`
- card: `Let's Dance`
- expected: evidence mapping should use canonical `lets-dance`
- observed: archive uses `let-s-dance`
- severity: **P3**
- evidence/reasoning: identity mapping inconsistency; gameplay content itself matched

### B04-ANG-001 — Wisdom of Ages adds a nonexistent conditional extra draw

- fighter: `angel`
- card: `Wisdom of Ages`
- expected: AFTER COMBAT, draw **1 card**
- observed: Git draws 2 after a win and 1 otherwise
- severity: **P1**
- evidence/reasoning: supplied card text is an unconditional one-card draw; Git introduces a result-dependent second card

### B04-ANG-002 — Cursed with a Soul printed BOOST mismatch

- fighter: `angel`
- card: `Cursed with a Soul`
- expected: printed BOOST **3**
- observed: Git records BOOST **2**
- severity: **P1**
- evidence/reasoning: direct physical metadata mismatch

### B04-ANG-003 — Brooding printed BOOST mismatch

- fighter: `angel`
- card: `Brooding`
- expected: printed BOOST **2**
- observed: Git records BOOST **1**
- severity: **P1**
- evidence/reasoning: direct physical metadata mismatch

### B04-ANG-004 — The Rogue Slayer printed BOOST mismatch

- fighter: `angel`
- card: `The Rogue Slayer`
- expected: printed BOOST **3**
- observed: Git records BOOST **2**
- severity: **P1**
- evidence/reasoning: direct physical metadata mismatch

### B04-LRRH-001 — Into the Woods physical variants collapsed in archive

- fighter: `little-red-riding-hood`
- card: `Into the Woods`
- expected: three physically distinct symbol variants, one each for `wolf`, `pelt`, and `knife`, matching canonical definitions `into-the-woods-wolf`, `into-the-woods-pelt`, and `into-the-woods-knife`, each quantity 1
- observed: archive contains one wolf-symbol image/manifest entry with quantity 3; pelt and knife physical variants are absent
- severity: **P2**
- evidence/reasoning: aggregate deck quantity is correct but two canonical physical component definitions cannot be visually certified

### B04-LRRH-002 — A Grimm Tale checks the wrong Basket item

- fighter: `little-red-riding-hood`
- card: `A Grimm Tale`
- expected: recover 4 when the effective Basket item is **wolf**; otherwise recover 2
- observed: Git grants the 4-health branch for **pelt** and 2 for non-pelt
- severity: **P1**
- evidence/reasoning: wrong conditional branch changes recovery amount in normal gameplay

### B04-LRRH-003 — The Wolf's Skin archive ID mismatch

- fighter: `little-red-riding-hood`
- card: `The Wolf's Skin`
- expected: canonical `the-wolfs-skin`
- observed: archive uses `the-wolf-s-skin`
- severity: **P3**
- evidence/reasoning: mapping-only discrepancy

### B04-LRRH-004 — What's That in My Basket? archive ID mismatch

- fighter: `little-red-riding-hood`
- card: `What's That in My Basket?`
- expected: canonical `whats-that-in-my-basket`
- observed: archive uses `what-s-that-in-my-basket`
- severity: **P3**
- evidence/reasoning: mapping-only discrepancy

### B04-DP-001 — Gaze of Stone printed BOOST mismatch

- fighter: `deadpool`
- card: `Gaze of Stone`
- expected: supplied physical evidence shows printed BOOST **4**; canonical card metadata must agree with the verified physical component or the evidence image must be shown to be the wrong component
- observed: Git records BOOST **1**
- severity: **P1**
- evidence/reasoning: direct component metadata conflict; external mirror/Medusa behavior under `A-REQ-015` does not reconcile the printed BOOST mismatch

### B04-DP-002 — Dumpster Divin' Deadpool makes the five-card operation optional by count

- fighter: `deadpool`
- card: `Dumpster Divin' Deadpool`
- expected: shuffle **5 cards** from the discard pile into the deck, then recover 1 health; if fewer than five exist, ordinary partial-resolution rules may limit what is possible, but the player should not voluntarily choose 0–4 when five are available
- observed: Git choice allows `min: 0, max: 5`
- severity: **P1**
- evidence/reasoning: the canonical representation creates player choice that is absent from the printed instruction and changes legal resolution

### B04-DP-003 — Sweeet! uses PLACE instead of MOVE

- fighter: `deadpool`
- card: `Sweeet!`
- expected: first **move** to a space in a yellow-ish zone, then **move** to a different space in a yellow-ish zone
- observed: Git performs two `PLACE` operations
- severity: **P1**
- evidence/reasoning: the project's normalized effect model explicitly distinguishes `MOVE` (movement/path rules apply) from `PLACE` (no path). The source uses MOVE, so the normalized behavior is not equivalent

### B04-DP-004 — Transit Card uses PLACE instead of MOVE

- fighter: `deadpool`
- card: `Transit Card`
- expected: AFTER COMBAT, **move** Deadpool to any space in his zone
- observed: Git uses `PLACE` into a chosen zone space
- severity: **P1**
- evidence/reasoning: canonical effect model treats MOVE and PLACE as materially different operations; no card-specific erratum was documented for Transit Card

### B04-DP-005 — Time out time out time out! omits printed social instruction

- fighter: `deadpool`
- card: `Time out time out time out!`
- expected: preserve or explicitly classify the source-defined instruction to call “time out” before the combat-card replacement
- observed: Git models the combat-card replacement but omits the spoken instruction
- severity: **P3**
- evidence/reasoning: no gameplay-state discrepancy was found, but the physical/social printed instruction is part of the component and belongs under the unresolved Deadpool digital-adaptation policy

### B04-DP-006 — Holy Mackerel! omits printed “go fish” instruction

- fighter: `deadpool`
- card: `Holy Mackerel!`
- expected: retain or explicitly classify the printed spoken instruction when the named-card guess fails
- observed: Git correctly models the state-changing draw but omits the spoken instruction
- severity: **P3**
- evidence/reasoning: non-gameplay social/component fidelity issue under `A-REQ-015`

### B04-DP-007 — 3 of Hearts archive ID mismatch

- fighter: `deadpool`
- card: `3 of Hearts`
- expected: canonical `three-of-hearts`
- observed: archive uses `3-of-hearts`
- severity: **P3**
- evidence/reasoning: mapping-only discrepancy

### B04-DP-008 — I'm Not Wearing Pants archive ID mismatch

- fighter: `deadpool`
- card: `I'm Not Wearing Pants`
- expected: canonical `im-not-wearing-pants`
- observed: archive uses `i-m-not-wearing-pants`
- severity: **P3**
- evidence/reasoning: mapping-only discrepancy

### B04-DP-009 — Rob's Pouch & Shoe Emporium archive ID mismatch

- fighter: `deadpool`
- card: `Rob's Pouch & Shoe Emporium`
- expected: canonical `robs-pouch-and-shoe-emporium`
- observed: archive uses `rob-s-pouch-and-shoe-emporium`
- severity: **P3**
- evidence/reasoning: mapping-only discrepancy

### B04-DP-010 — Deadpool™ Merc For Hire, LLC archive ID mismatch

- fighter: `deadpool`
- card: `Deadpool™ Merc For Hire, LLC`
- expected: canonical `deadpool-merc-for-hire-llc`
- observed: archive uses `deadpooltm-merc-for-hire-llc`
- severity: **P3**
- evidence/reasoning: mapping-only discrepancy

## Material Findings Index

| ID | Severity | Fighter | Card | Finding |
|---|---|---|---|---|
| B04-SPI-001 | P1 | spike | The Sight | Image BOOST 2; Git 3 |
| B04-SPI-002 | P1 | spike | Arrogance | Printed whole-hand discard/set-to-6 behavior modeled incorrectly |
| B04-ANG-001 | P1 | angel | Wisdom of Ages | Printed draw 1 replaced by win-sensitive 2/1 draw |
| B04-ANG-002 | P1 | angel | Cursed with a Soul | Image BOOST 3; Git 2 |
| B04-ANG-003 | P1 | angel | Brooding | Image BOOST 2; Git 1 |
| B04-ANG-004 | P1 | angel | The Rogue Slayer | Image BOOST 3; Git 2 |
| B04-LRRH-001 | P2 | little-red-riding-hood | Into the Woods | Three symbol variants collapsed to one image; two variants absent |
| B04-LRRH-002 | P1 | little-red-riding-hood | A Grimm Tale | 4-heal branch keyed to pelt instead of wolf |
| B04-DP-001 | P1 | deadpool | Gaze of Stone | Image BOOST 4; Git 1 |
| B04-DP-002 | P1 | deadpool | Dumpster Divin' Deadpool | Git permits choosing 0–5 instead of requiring five when available |
| B04-DP-003 | P1 | deadpool | Sweeet! | Two printed MOVE operations normalized as PLACE |
| B04-DP-004 | P1 | deadpool | Transit Card | Printed MOVE normalized as PLACE |

## Corpus-Level Observations

1. **Archive technical integrity was strong.** All five nested ZIPs were readable and all 78 supplied images decoded. No zero-byte files, corrupt images, duplicate binaries, orphan images, or manifest-image reference failures were observed.

2. **Archive manifests are useful mapping evidence but are not authoritative.** The Little Red `Into the Woods` collapse demonstrates why quantity totals alone cannot establish unique physical-definition completeness. Aggregate quantity 30 passed while two canonical printed variants were absent.

3. **Canonical deck sizes must be checked through construction data, not assumed.** These five fighters all happen to be fixed 30-card decks, and each archive quantity sum was 30, but the QA treated this as manifest evidence rather than a universal Unmatched rule.

4. **Printed metadata errors are material.** Five direct BOOST mismatches were found: Spike/The Sight; Angel/Cursed with a Soul; Angel/Brooding; Angel/The Rogue Slayer; Deadpool/Gaze of Stone.

5. **Normalized semantics must preserve source verbs and dependencies.** The Deadpool `Sweeet!` and `Transit Card` findings are material specifically because the project effect model defines MOVE and PLACE as different operations. Likewise, Spike/Arrogance and Little Red/A Grimm Tale are not merely structural differences; their resulting legal behavior differs from the printed cards.

6. **Card-specific errata must not be mistaken for transcription errors.** Angel/Disengage visibly says Move, but canonical Git deliberately models PLACE and cites `disengage_errata`; QA therefore accepted that normalized behavior.

7. **Shared integration requirements can be visually confirmed without becoming transcription failures.** This batch directly confirms `A-REQ-003`, `A-REQ-007`, `A-REQ-008`, `A-REQ-009`, and `A-REQ-015` where the relevant behavior appears on supplied card components. Engine/policy gaps do not by themselves invalidate an otherwise accurate card transcription.

8. **Deadpool remains a special policy case.** The physical cards visibly depend on real-world/social predicates and actions. The repository's `A-REQ-015` digital-adaptation-policy blocker is therefore substantiated by the evidence. The policy blocker is separate from the four independent P1 corpus mismatches found for Deadpool in this QA.

9. **Fighter abilities were not treated as visually certified unless present on card evidence.** Fighter manifests were used for topology/resources/setup context only. This report does not claim visual verification of fighter cards/ability panels not contained in the evidence archive.

10. **Clean fighter:** `beowulf` had complete archive coverage, correct quantities, matching printed metadata, and semantically equivalent normalized effects for every supplied unique card.

## Final Assessment

This batch is **not ready to be accepted as a fully verified Phase 4 card-image corpus**.

Reasons:

- Spike has two P1 canonical card discrepancies.
- Angel has four P1 canonical card discrepancies.
- Little Red Riding Hood has one P1 semantic discrepancy and lacks two canonical physical `Into the Woods` variant images.
- Deadpool has four P1 canonical card discrepancies in addition to its separately documented digital-adaptation policy qualification.
- Only Beowulf passed all card-image checks without material findings.

The evidence archive itself is technically well formed, and all supplied images were inspectable. The failure is primarily one of **canonical transcription/normalization accuracy** plus the localized Little Red evidence-completeness gap, not ZIP corruption or broad archive failure.

No fixes were applied to canonical manifests as part of this QA.