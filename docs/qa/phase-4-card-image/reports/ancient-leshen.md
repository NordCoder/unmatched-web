# Phase 4 Card-Image QA Report

## Batch

- ancient-leshen
- blackbeard
- eredin
- philippa
- yennefer-and-triss

## Verdict

**FAIL**

The batch archive is technically complete and all 58 unique card images were successfully inspected. All five deck constructions reconcile to their canonical 30-card game decks. One material gameplay-semantic discrepancy was found in `philippa / Do My Bidding`; the other 57 unique cards match canonical Git metadata and normalized printed gameplay semantics. Several non-gameplay slug/mapping differences and documented integration qualifications are recorded below.

## QA Scope And Repository Snapshot

Repository: `NordCoder/unmatched-web`

This report persists the already-completed independent QA; no card research was restarted during report persistence.

Canonical refs inspected during the QA run:

- `phase-4b-worker-c-modern` @ `1f4fd973de66ab7995c0b2213e7dcda4163ff7c8`
- `phase-4b-worker-d-latest` @ `b9ae31c3b1a958e34bfb507d695cbd14650b9ed6`
- `main` @ `f40c3f9970da24ab4e17ff51fc75fb5d7080b239` for the Phase 4A Yennefer/Triss pair

Persistence write base:

- `main` @ `2fcc915301d9a69157395fbcffd1e3de486c2e40`

The report itself does not modify any canonical fighter/card manifest, schema, source document, or other QA report.

## Summary

### Batch-level archive integrity

- outer batch ZIP readable: PASS
- nested fighter ZIPs found: 5
- expected fighter ZIP range (3-5): PASS
- nested ZIP CRC/integrity: PASS for all 5
- fighter manifests readable: PASS for all 5
- total unique card images: 58
- images successfully decoded/inspected: 58
- zero-byte images: 0
- corrupted/truncated images: 0
- duplicate binary images: 0
- missing images: 0
- extra images: 0
- manifest entries without image: 0
- images without manifest entry: 0

### Batch-level deck construction

| Fighter | Unique definitions | Archive quantity sum | Git available pool | Git game deck | Construction | Result |
|---|---:|---:|---:|---:|---|---|
| ancient-leshen | 11 | 30 | 30 | 30 | fixed | PASS |
| blackbeard | 12 | 30 | 30 | 30 | fixed | PASS |
| eredin | 11 | 30 | 30 | 30 | fixed | PASS |
| philippa | 13 | 30 | 30 | 30 | fixed | PASS |
| yennefer-and-triss | 11 | 30 | 30 | 30 | fixed | PASS |

Total: **58 unique definitions / 150 physical action-card copies represented by quantity**.

No fighter in this batch uses a non-standard choose-group construction or external gameplay definition that changes the action-deck count.

## Fighter Results

### ancient-leshen

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/ancient-leshen.yaml`
- fighter manifest: `docs/fighters/phase-4b/ancient-leshen.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: 11/11
- archive quantity sum: 30
- canonical available pool: 30
- canonical game deck: 30
- printed card identity/content verification: PASS for all 11 cards
- discrepancies: none
- verdict: **PASS_WITH_QUALIFICATIONS**

The qualification is integration/evidence-related, not a card transcription failure. The canonical pair documents shared requirements and a derived dormant-return interpretation around `Vanish Into Murder`.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Wily Fighting | PASS | PASS | PASS | PASS |
| Flock of Birds | PASS | PASS | PASS | PASS |
| Command the Forest | PASS | PASS | PASS | PASS |
| Primeval Slam | PASS | PASS | PASS | PASS |
| Nature Abounds | PASS | PASS | PASS | PASS |
| Primeval Guardian | PASS | PASS | PASS | PASS |
| Planted Feet | PASS | PASS | PASS | PASS |
| Strength of the Pack | PASS | PASS | PASS | PASS |
| Vanish Into Murder | PASS | PASS | PASS | PASS |
| Disturbing Howls | PASS | PASS | PASS | PASS |
| Harrying Strike | PASS | PASS | PASS | PASS |

#### Evidence notes

`Vanish Into Murder` visually confirms the gameplay sequence represented in Git: effect damage, temporary non-defeat removal of Ancient Leshen, scheduled return at the next controller turn start, preservation of health, and the subsequent draw after return. The image does not contradict the repository's explicitly qualified dormant-player interpretation.

`Planted Feet` and the two distinct Wolf fighter instances also support the need for runtime-fighter-instance historical movement state rather than a single fighter-definition boolean.

#### Integration requirements confirmed

- `C-REQ-007` — temporary non-defeat battlefield absence with scheduled return
- `C-REQ-010` — historical state keyed by runtime fighter instance
- `C-REQ-013` — ordered staged resolution with bindings/dependent continuation

`C-REQ-011` remains part of the canonical fighter/deck pair, but the action-card images alone are not sufficient evidence for the complete fighter-level derived-attribute contract.

---

### blackbeard

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/blackbeard.yaml`
- fighter manifest: `docs/fighters/phase-4b/blackbeard.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: 12/12
- archive quantity sum: 30
- canonical available pool: 30
- canonical game deck: 30
- printed card identity/content verification: PASS for all 12 cards
- discrepancies: two P3 archive-to-Git card-ID slug differences only
- verdict: **PASS_WITH_QUALIFICATIONS**

The qualification is the documented `D-REQ-016` shared-Treasury ransom arbitration requirement plus non-gameplay archive mapping differences.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Avast Ye! | PASS | PASS | PASS | PASS |
| A Brace of Primed Pistols | PASS | PASS* | PASS | PASS |
| Give No Quarter | PASS | PASS | PASS | PASS |
| Intimidating Visage | PASS | PASS | PASS | PASS |
| Queen Anne's Revenge | PASS | PASS* | PASS | PASS |
| Fearsome and Calculating | PASS | PASS | PASS | PASS |
| Parley | PASS | PASS | PASS | PASS |
| No Prey, No Pay | PASS | PASS | PASS | PASS |
| Show a Leg! | PASS | PASS | PASS | PASS |
| Plunder | PASS | PASS | PASS | PASS |
| Scourge of the Seven Seas | PASS | PASS | PASS | PASS |
| Light the Fuse | PASS | PASS | PASS | PASS |

`*` Metadata qualification is limited to archive mapping slug normalization; printed names, fighter restrictions, card types, printed values, BOOST values, quantities and gameplay behavior match Git.

#### Evidence notes

Visual inspection confirms the ransom amounts, timing windows and unpaid branches represented by Git. In particular:

- `Avast Ye!` has separate ransomable IMMEDIATELY and AFTER COMBAT clauses;
- `A Brace of Primed Pistols` selects an opposing fighter in Blackbeard's zone, then presents the two-doubloon ransom before the three-damage unpaid branch;
- `Queen Anne's Revenge` reduces its value by 2 for each Treasury doubloon paid and transferred;
- `Parley` uses a two-doubloon ransom before the unpaid cancellation branch;
- `Scourge of the Seven Seas` contains three separately ransomable effects rather than one aggregate payment;
- `Plunder` uses a deterministic Treasury-to-Blackbeard transfer and is not itself a multi-actor ransom window.

#### Integration requirements confirmed

- `D-REQ-016` — multi-actor payment arbitration over the shared Treasury, with the first accepted payment closing the corresponding ransom window and transferring the stated amount to Blackbeard

---

### eredin

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/eredin.yaml`
- fighter manifest: `docs/fighters/phase-4b/eredin.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: 11/11
- archive quantity sum: 30
- canonical available pool: 30
- canonical game deck: 30
- printed card identity/content verification: PASS for all 11 cards
- discrepancies: one P3 archive-to-Git card-ID slug difference only
- verdict: **PASS_WITH_QUALIFICATIONS**

The qualification is shared integration plus the canonical high-confidence derived interpretation for `Icy Guile` self-defeat; no printed-card transcription discrepancy was found.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Skirmish | PASS | PASS | PASS | PASS |
| Backhand | PASS | PASS | PASS | PASS |
| Brutal Strike | PASS | PASS | PASS | PASS |
| Foul Purpose | PASS | PASS | PASS | PASS |
| Unyielding Hordes | PASS | PASS | PASS | PASS |
| Portal Defense | PASS | PASS | PASS | PASS |
| Implacable | PASS | PASS | PASS | PASS |
| Icy Guile | PASS | PASS | PASS | PASS |
| Wild Hunt | PASS | PASS | PASS | PASS |
| Might of the Aen Elle | PASS | PASS* | PASS | PASS |
| Close for the Kill | PASS | PASS | PASS | PASS |

`*` Metadata qualification is an ID slug difference only.

#### Evidence notes

`Icy Guile` visually says a Red Rider rather than "another Red Rider", so the printed component does not exclude the current Red Rider combat fighter from the optional defeat cost. The repository's further conclusion that the already-started combat continues after such a self-defeat remains explicitly marked as a derived interpretation, which is appropriate and was not treated as a transcription discrepancy.

`Foul Purpose`, `Implacable`, and `Portal Defense` directly confirm that Eredin's corpus needs declaration-time legality/type changes and same-combat defender replacement semantics rather than simple value modifications.

#### Integration requirements confirmed

- `C-REQ-003` — declaration-time combat legality / same-combat participant replacement
- `C-REQ-011` — source-lifetime derived fighter/card attributes and permissions
- `C-REQ-013` — staged choices, bindings and dependent continuation

---

### philippa

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/philippa.yaml`
- fighter manifest: `docs/fighters/phase-4b/philippa.yaml`
- archive integrity: PASS
- canonical manifest comparison: FAIL on one card semantic restriction
- unique-card image completeness: 13/13
- archive quantity sum: 30
- canonical available pool: 30
- canonical game deck: 30
- printed card identity/content verification: 12 PASS / 1 FAIL
- discrepancies: one P1 semantic discrepancy; one P3 archive mapping difference
- verdict: **FAIL**

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Paralyzing Fetters | PASS | PASS | PASS | PASS |
| Owlform | PASS | PASS | PASS | PASS |
| Lightning Bolt | PASS | PASS | PASS | PASS |
| Chain Lightning | PASS | PASS | PASS | PASS |
| Regicide | PASS | PASS | PASS | PASS |
| Blinding Dust | PASS | PASS | PASS | PASS |
| **Do My Bidding** | PASS | PASS | **FAIL** | **FAIL** |
| Spellbreaker | PASS | PASS | PASS | PASS |
| Redanian Plot | PASS | PASS | PASS | PASS |
| Cunning | PASS | PASS | PASS | PASS |
| Polymorphy | PASS | PASS | PASS | PASS |
| Backup Plan | PASS | PASS | PASS | PASS |
| Spymaster's Ruse | PASS | PASS* | PASS | PASS |

#### Evidence notes

The rest of Philippa's corpus matches the normalized Git semantics, including:

- `Backup Plan` sets Philippa's health to exactly 5 rather than merely recovering toward 5;
- `Polymorphy` places Philippa and creates a source-lifetime movement-value override;
- `Spymaster's Ruse` makes the opponent choose cards to reveal before Philippa's controller chooses the discard subset;
- `Blinding Dust` binds the discarded card before using that card's BOOST value for damage;
- `Paralyzing Fetters` resets/locks the opposing combat card to its printed value for the current combat.

#### Integration requirements confirmed

- `C-REQ-008` — exact health assignment
- `C-REQ-011` — source-lifetime derived attribute modifier
- `C-REQ-013` — ordered resolution with private/public choices and dependencies

---

### yennefer-and-triss

- canonical branch: `main`
- requested archive/QA fighter ID: `yennefer-and-triss`
- current canonical Git fighter ID/file stem inspected by QA: `yennefer-triss`
- deck manifest: `docs/cards/phase-4a/yennefer-triss.yaml`
- fighter manifest: `docs/fighters/phase-4a/yennefer-triss.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS after unambiguous alias normalization
- unique-card image completeness: 11/11
- archive quantity sum: 30
- canonical available pool: 30
- canonical game deck: 30
- printed card identity/content verification: PASS for all 11 cards
- discrepancies: fighter/path naming drift plus one P3 card-ID punctuation slug difference
- verdict: **PASS_WITH_QUALIFICATIONS**

The supplied Phase 4 QA mapping names this fighter `yennefer-and-triss`, but the actual current Phase 4A canonical files inspected during QA use `yennefer-triss`. The identity is unambiguous and no second fighter was inferred or created.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Incinerate | PASS | PASS | PASS | PASS |
| Merigold's Hailstorm | PASS | PASS* | PASS | PASS |
| Echoing Blast | PASS | PASS | PASS | PASS |
| Portal to Anywhere | PASS | PASS | PASS | PASS |
| Paralyzing Fetters | PASS | PASS | PASS | PASS |
| Magical Barrier | PASS | PASS | PASS | PASS |
| Ball Lightning | PASS | PASS | PASS | PASS |
| Telepathy | PASS | PASS | PASS | PASS |
| Quick and Ready | PASS | PASS | PASS | PASS |
| Advisor to the King | PASS | PASS | PASS | PASS |
| Lodge of Sorceresses | PASS | PASS | PASS | PASS |

`*` Metadata qualification is an archive mapping punctuation/slug difference only.

#### Evidence notes

Visual inspection confirms the special semantics represented in Git:

- `Incinerate`: opponent may pay the source-defined two-card discard cost to ignore the card's value;
- `Portal to Anywhere`: uses PLACE rather than MOVE and then grants an action;
- `Telepathy`: opponent may discard a card or the owner's combat-card effects are canceled;
- `Advisor to the King`: ongoing scheme modifies BOOST values and has its own discard condition;
- `Lodge of Sorceresses`: simultaneous hidden commitment/reveal structure is required to preserve information boundaries.

#### Integration requirements confirmed

No explicit `X-REQ-*` identifier is attached to this Phase 4A pair in the inspected canonical manifests.

## Findings

### B12-PHILIPPA-001 — Do My Bidding replacement-card legality

- fighter: `philippa`
- card: `Do My Bidding`
- severity: **P1**
- Git location: `docs/cards/phase-4b/philippa.yaml`, `do-my-bidding` replacement choice
- expected: after returning the original attack and inspecting the opponent's hand, the chosen replacement must be an ATTACK/VERSATILE card that the **original attacking fighter can legally play in the already-declared combat**. Normal card usability and combat legality remain relevant because the printed effect instructs that opponent to play the chosen card instead.
- observed: the current normalized choice domain filters the viewed hand only by `type_in: [attack, versatile]`, then passes the selected card into `REPLACE_COMBAT_CARD`. The choice-domain record does not encode `usable_by` / original-attacker / current-combat legality.
- evidence/reasoning: the printed card says to choose an attack or versatile card **for the opponent to play instead**, not merely to substitute any card-shaped object. A hero-only card is therefore not a legal replacement when the original attacker is an incompatible sidekick, and equivalent source-defined attack-legality restrictions must also remain enforced. Exposing those cards as valid choices changes the legal gameplay decision space.
- expected semantic correction: constrain the replacement-card choice to cards legally playable by the original attacker in that current combat, including fighter usability and other relevant combat legality, while preserving same-combat replacement semantics.

This finding is gameplay-changing and makes the fighter/batch verdict FAIL even though the archive itself is complete.

### B12-EREDIN-P3-01 — Might of the Aen Elle card-ID slug

- fighter: `eredin`
- card: `Might of the Aen Elle`
- severity: **P3**
- expected canonical Git ID: `might-of-aen-elle`
- observed archive mapping ID: `might-of-the-aen-elle`
- evidence/reasoning: printed name, quantity, card type, value, BOOST and effect all match. Difference is mapping/slug normalization only.
- expected correction: align archive mapping to the canonical card definition ID; no gameplay correction.

### B12-PHILIPPA-P3-02 — Spymaster's Ruse card-ID slug

- fighter: `philippa`
- card: `Spymaster's Ruse`
- severity: **P3**
- expected canonical Git ID: `spymasters-ruse`
- observed archive mapping ID: `spymaster-s-ruse`
- evidence/reasoning: apostrophe-derived slug difference only; printed and semantic identity is unambiguous.
- expected correction: align archive mapping to canonical ID; no gameplay correction.

### B12-YT-P3-01 — fighter/path ID drift

- fighter: archive/QA ID `yennefer-and-triss`
- severity: **P3**
- expected according to the QA contract mapping: `yennefer-and-triss`
- observed current Phase 4A Git files during QA: `yennefer-triss`
- evidence/reasoning: both sources unambiguously describe the same Yennefer & Triss selectable-hero fighter pair. No competing fighter identity exists in the batch.
- expected correction: reconcile orchestration/archive alias mapping with the actual canonical Git ID/file stem. No gameplay correction.

### B12-YT-P3-02 — Merigold's Hailstorm card-ID slug

- fighter: `yennefer-and-triss` / Git `yennefer-triss`
- card: `Merigold's Hailstorm`
- severity: **P3**
- expected canonical Git ID: `merigolds-hailstorm`
- observed archive mapping ID: `merigold-s-hailstorm`
- evidence/reasoning: punctuation-to-slug difference only; card identity and gameplay semantics match.
- expected correction: align archive mapping to canonical ID; no gameplay correction.

### B12-BB-P3-01 — A Brace of Primed Pistols card-ID slug

- fighter: `blackbeard`
- card: `A Brace of Primed Pistols`
- severity: **P3**
- expected canonical Git ID: `brace-of-primed-pistols`
- observed archive mapping ID: `a-brace-of-primed-pistols`
- evidence/reasoning: canonical slug omits the initial article; printed and gameplay identity match.
- expected correction: align archive mapping to canonical ID; no gameplay correction.

### B12-BB-P3-02 — Queen Anne's Revenge card-ID slug

- fighter: `blackbeard`
- card: `Queen Anne's Revenge`
- severity: **P3**
- expected canonical Git ID: `queen-annes-revenge`
- observed archive mapping ID: `queen-anne-s-revenge`
- evidence/reasoning: apostrophe-to-slug normalization only; printed and gameplay identity match.
- expected correction: align archive mapping to canonical ID; no gameplay correction.

## Corpus-Level Observations

1. **Archive construction quality is clean.** The batch follows the intended one-definition/one-image rule. There are no duplicated physical-copy images, no corrupted files and no orphan image/manifest entries.

2. **Quantity data is fully reconciled.** All five archive manifests sum to 30 and match their corresponding canonical fixed constructions. No quantity mismatch was found at card or deck level.

3. **Printed metadata is highly consistent with Git.** Across 58 images, printed names, combat types, printed values, BOOST values, fighter restrictions and quantities match the canonical corpus except for non-gameplay ID-slug mapping differences listed above.

4. **Structured normalization generally preserves the printed semantics.** The QA did not require literal card-text equality. Ordered stages, choices, costs, bindings, PLACE/MOVE distinctions, ongoing-source lifetimes and resource-transfer branches were accepted where they preserve the physical card's legal behavior.

5. **Shared integration requirements are not transcription failures.** `C-REQ-*` / `D-REQ-*` dependencies confirmed by images were recorded separately. A card can therefore pass printed-semantic QA while the pair remains `partial` for runtime integration.

6. **Derived/project interpretations were not promoted to printed facts.** Ancient Leshen's dormant-return continuation and Eredin's Icy Guile self-defeat continuation remain explicitly qualified repository interpretations. The images support the underlying printed wording but do not independently establish every engine-level consequence.

7. **Replacement effects require legal-choice-domain modeling, not only replacement execution.** `Do My Bidding` demonstrates that a correct `REPLACE_COMBAT_CARD` primitive is insufficient if the prior selectable domain admits cards the original combat participant cannot legally play. This is the sole material semantic defect found in the batch.

## Material Findings

P1/P2 only:

| ID | Severity | Fighter | Card | Result |
|---|---|---|---|---|
| B12-PHILIPPA-001 | P1 | philippa | Do My Bidding | replacement choice lacks original-attacker/current-combat legality constraint |

No P2 findings.

## Clean Fighters

No P1/P2 material discrepancies:

- ancient-leshen
- blackbeard
- eredin
- yennefer-and-triss / canonical Git `yennefer-triss`

`philippa` has 12 clean unique definitions and one P1-failing definition.

## Blockers

None.

All 58 images were readable and inspectable, every fighter was identified unambiguously, and all required canonical Git sources were accessible during the QA run.

## Final Assessment

The batch is **technically suitable as physical card-image evidence**: archive structure, image completeness, binary integrity and deck quantities all pass, and 57 of 58 unique card definitions reconcile with canonical Git gameplay semantics.

The batch as a whole is **not ready for an unconditional Phase 4 card-image corpus PASS** because `philippa / Do My Bidding` currently has a P1 semantic gap in the canonical normalized representation: the replacement choice does not encode that the selected card must be legally playable by the original attacker in the current combat.

Final batch verdict: **FAIL**.

No repository correction was made by this QA worker; this report only records the evidence and expected semantic correction.