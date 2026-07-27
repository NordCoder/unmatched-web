# Phase 4 Card-Image QA Report

## Batch

- bigfoot
- bruce-lee
- dracula
- raptors
- robert-muldoon

Filename anchor: `bigfoot` (lexicographically first fighter ID in this batch).

## Verdict

**FAIL**

The batch is technically complete and fully inspectable, and all five fighter archives reconcile to the expected deck quantities. However, three P1 gameplay-changing discrepancies were found in canonical Git card data:

1. `bruce-lee / one-inch-punch` narrows a printed `adjacent fighter` target to `adjacent_opposing_fighters`.
2. `raptors / ambush` is printed as ATTACK but canonical Git stores it as `versatile`.
3. `raptors / eaten-alive` is printed as VERSATILE but canonical Git stores it as `attack`.

Two additional P3 archive-manifest card-ID mapping mismatches were found. No archive corruption, unreadable images, missing images, duplicate binary card images, quantity failures, or blockers were found.

## Summary

Independent Phase 4 Card-Image QA was performed read-only against the physical card-image evidence batch and the canonical repository state.

Repository: `NordCoder/unmatched-web`

Canonical refs used during QA:

| Scope | Branch | Verified tip during QA |
|---|---|---|
| Phase 4A fighters | `main` | `f40c3f9970da24ab4e17ff51fc75fb5d7080b239` |
| Phase 4B Worker A fighters | `phase-4b-worker-a-classics` | `840a03572bf8b5912b03f32e5ae2a48cc8fadc2b` |
| Phase 4B Worker D fighters | `phase-4b-worker-d-latest` | `b9ae31c3b1a958e34bfb507d695cbd14650b9ed6` |

Immediately before persistence, `main` was rechecked and still resolved to `f40c3f9970da24ab4e17ff51fc75fb5d7080b239`.

Outer archive contained exactly five fighter ZIP files:

- `bigfoot.zip`
- `bruce-lee.zip`
- `dracula.zip`
- `raptors.zip`
- `robert-muldoon.zip`

All nested ZIP files opened successfully. Every fighter archive contained a readable `manifest.yaml`. Every card image decoded successfully and was visually inspected; filename and archive manifest were used only for mapping, not as authoritative printed evidence.

### Batch corpus counts

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| Total unique card images | 63 |
| Images successfully inspected | 63 |
| Unreadable images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Duplicate binary images | 0 |
| Corrupted images | 0 |
| Corrupted nested ZIPs | 0 |
| Quantity failures | 0 |
| P1 findings | 3 |
| P2 findings | 0 |
| P3 findings | 2 |

## Fighter Results

### bigfoot

Canonical source:

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/bigfoot.yaml`
- fighter manifest: `docs/fighters/phase-4b/bigfoot.yaml`
- construction: `fixed`
- canonical available pool: 30
- canonical game deck: 30

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: PASS
- manifest readable: PASS
- images readable: 11/11
- unique images: 11
- manifest card entries: 11
- duplicate binary images: 0
- missing images: 0
- extra images: 0

#### Canonical manifest comparison

- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 11
- unique Git definitions: 11
- quantity validation: **PASS**

#### Unique-card image completeness

All 11 unique physical card definitions were present exactly once as images. Physical copy counts are represented through archive/Git quantities rather than duplicate image files.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Savagery | PASS | PASS | PASS | PASS |
| It's Just Your Imagination | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Jackalope Horns | PASS | PASS | PASS | PASS |
| Crash Through the Trees | PASS | PASS | PASS | PASS |
| Hoax | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Disengage | PASS | PASS | PASS | PASS |
| Larger Than Life | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |

High-risk semantics explicitly checked:

- `Jackalope Horns`: the printed effect moves the Jackalope, then deals 2 damage to any one fighter adjacent to the Jackalope. Git preserves the second operation as a dependent later stage with an adjacent-fighter choice.
- `Disengage`: legal behavior is represented as placement into an empty space in the combat fighter's zone, rather than ordinary movement traversal.
- `Regroup`: normal outcome draws 1 card; a won combat draws 2 cards instead. Git represents the two outcomes as mutually exclusive combat-result branches rather than cumulative draws.
- `Momentous Shift`: value modification is conditional on the combat fighter having started the turn in a different space.

#### Discrepancies

None.

#### Integration requirements confirmed

None directly required by these card images.

#### Verdict

**PASS**

---

### bruce-lee

Canonical source:

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/bruce-lee.yaml`
- fighter manifest: `docs/fighters/phase-4b/bruce-lee.yaml`
- construction: `fixed`
- canonical available pool: 30
- canonical game deck: 30

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: PASS
- manifest readable: PASS
- images readable: 17/17
- unique images: 17
- manifest card entries: 17
- duplicate binary images: 0
- missing images: 0
- extra images: 0

#### Canonical manifest comparison

- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 17
- unique Git definitions: 17
- quantity validation: **PASS**

#### Unique-card image completeness

All 17 unique physical card definitions were present and inspected.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Little Dragon | PASS | PASS | PASS | PASS |
| Taste of Blood | PASS | PASS | PASS | PASS |
| Be Like Water | PASS | PASS | PASS | PASS |
| Jeet Kune Do: High Straight Lead | PASS | PASS | PASS | PASS |
| Jeet Kune Do: Downward Side Kick | PASS | PASS | PASS | PASS |
| Jeet Kune Do: Wrist Lock | PASS | PASS | PASS | PASS |
| Jeet Kune Do: Intercepting Fist | PASS | PASS | PASS | PASS |
| Jeet Kune Do: Corkscrew Finger Jab | PASS | PASS | PASS | PASS |
| Jeet Kune Do: Short Lead Hook | PASS | PASS | PASS | PASS |
| Bring It On | PASS | PASS | PASS | PASS |
| One-Inch Punch | PASS | PASS | FAIL | FAIL |
| Nunchaku | PASS | PASS | PASS | PASS |
| "HOO! WHAAAAAA!" | PASS | P3 archive mapping | PASS | PASS_WITH_QUALIFICATION |

Printed `BRUCE LEE` restriction versus normalized `usable_by: any` on several generic cards was reviewed. Bruce Lee has single-hero topology, so this normalization does not expand legal users within this fighter and was not classified as a gameplay-changing discrepancy.

High-risk semantics explicitly checked:

- `High Straight Lead`: the condition refers to either combat fighter starting the turn in a different space; Git preserves this broader condition.
- `Intercepting Fist`: IMMEDIATELY cancellation and AFTER COMBAT action gain remain distinct effects with correct timing.
- `Short Lead Hook`: swap is optional and action gain remains a separate after-combat effect.
- `Bring It On`: uses placement of an opposing fighter, followed by action gain.
- `Nunchaku`: turn-scoped attack-value bonus and gained action are preserved.

#### Discrepancies

##### P1 — `One-Inch Punch` legal target narrowed

- fighter: `bruce-lee`
- card: `one-inch-punch`
- severity: **P1**
- Git location: `docs/cards/phase-4b/bruce-lee.yaml`, `one-inch-punch`
- expected from physical card: choose **an adjacent fighter**, deal 2 damage to that fighter; if that fighter is defeated by this card, return this card to hand instead of discarding it.
- observed in Git: choice domain is normalized as `adjacent_opposing_fighters`.
- evidence/reasoning: the printed physical card does not restrict the target to an opposing fighter. The Git selector therefore removes legal friendly targets. In multiplayer/team contexts this changes the legal action space and is gameplay-significant.
- expected semantic correction: target domain should include any adjacent fighter, while retaining the exact defeat dependency and return-to-hand replacement behavior.

##### P3 — archive mapping ID differs for `"HOO! WHAAAAAA!"`

- fighter: `bruce-lee`
- card: `"HOO! WHAAAAAA!"`
- severity: **P3**
- expected canonical card ID: `hoo-whaaaaa`
- observed archive manifest ID: `hoo-whaaaaaa`
- evidence/reasoning: the evidence manifest contains one extra `a`. Printed identity and gameplay semantics are otherwise matched correctly.
- expected correction: evidence/archive mapping should use the canonical Git card ID. No gameplay manifest correction is implied by this P3 finding.

#### Integration requirements confirmed

None.

#### Verdict

**FAIL**

Reason: one material P1 target-domain discrepancy.

---

### dracula

Canonical source:

- branch: `main`
- deck manifest: `docs/cards/phase-4a/dracula.yaml`
- fighter manifest: `docs/fighters/phase-4a/dracula.yaml`
- construction: `fixed`
- canonical available pool: 30
- canonical game deck: 30

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: PASS
- manifest readable: PASS
- images readable: 13/13
- unique images: 13
- manifest card entries: 13
- duplicate binary images: 0
- missing images: 0
- extra images: 0

#### Canonical manifest comparison

- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 13
- unique Git definitions: 13
- quantity validation: **PASS**

#### Unique-card image completeness

All 13 unique physical card definitions were present and inspected.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Ambush | PASS | PASS | PASS | PASS |
| Dash | PASS | PASS | PASS | PASS |
| Exploit | PASS | PASS | PASS | PASS |
| Feeding Frenzy | PASS | PASS | PASS | PASS |
| Beastform | PASS | PASS | PASS | PASS |
| Do My Bidding | PASS | PASS | PASS | PASS |
| Look Into My Eyes | PASS | PASS | PASS | PASS |
| Mistform | PASS | PASS | PASS | PASS |
| Prey Upon | PASS | PASS | PASS | PASS |
| Baptism of Blood | PASS | PASS | PASS | PASS |
| Thirst for Sustenance | PASS | PASS | PASS | PASS |
| Ravening Seduction | PASS | PASS | PASS | PASS |

High-risk semantics explicitly checked:

- `Ambush`: opponent random-discard occurs first; the discarded card's BOOST contributes to this card's value.
- `Beastform`: controller may discard any number of cards; combat value increases by the number actually discarded.
- `Do My Bidding`: current opposing attack returns to the opponent hand, Dracula's controller looks at that hand and chooses a replacement attack/versatile card. Git correctly models this as replacement of the current combat card, not a second Attack action.
- `Prey Upon`: Dracula recovery is tied to damage actually dealt to adjacent opposing fighters, preserving the damage-to-healing dependency.
- `Baptism of Blood`: a defeated Sister is returned to play with source-defined placement/health behavior; this is not ordinary recovery of a defeated fighter.
- `Ravening Seduction`: the target is moved first and damage is calculated from Sisters adjacent after movement.

#### Discrepancies

None.

#### Integration requirements confirmed

None directly required by these card images.

#### Verdict

**PASS**

---

### raptors

Canonical source:

- branch: `main`
- deck manifest: `docs/cards/phase-4a/raptors.yaml`
- fighter manifest: `docs/fighters/phase-4a/raptors.yaml`
- construction: `fixed`
- canonical available pool: 30
- canonical game deck: 30

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: PASS
- manifest readable: PASS
- images readable: 11/11
- unique images: 11
- manifest card entries: 11
- duplicate binary images: 0
- missing images: 0
- extra images: 0

#### Canonical manifest comparison

- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 11
- unique Git definitions: 11
- quantity validation: **PASS**

#### Unique-card image completeness

All 11 unique physical card definitions were present and inspected.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Disengage | PASS | PASS | PASS / errata-aware normalization | PASS_WITH_QUALIFICATION |
| Ambush | PASS | FAIL | PASS | FAIL |
| Coordinated Attack Pattern | PASS | PASS | PASS | PASS |
| Working Things Out | PASS | PASS | PASS | PASS |
| They Remember | PASS | PASS | PASS | PASS |
| Clever Girl | PASS | PASS | PASS | PASS |
| Pack Hunters | PASS | PASS | PASS | PASS |
| Eaten Alive | PASS | FAIL | PASS | FAIL |
| Decoy | PASS | PASS | PASS | PASS |
| Eviscerate | PASS | PASS | PASS | PASS |

`Disengage` was reviewed separately because the physical printing uses movement wording while canonical Git normalizes the resolved behavior to placement into an empty space in the combat fighter's zone. This is retained as an errata-aware/project normalization and was not classified as a card transcription discrepancy.

Other semantics checked:

- `Pack Hunters`: damage after a won combat depends on friendly Raptors adjacent to the opposing combat fighter.
- `Decoy`: only another undefeated Raptor may be placed; the representation must not revive a defeated Raptor.
- `Coordinated Attack Pattern`: choose one undefeated Raptor as the anchor and optionally place each other undefeated Raptor into its zone.
- `Working Things Out`: each undefeated Raptor may move up to 3 spaces and the card grants an action after resolving the movement effect.

#### Discrepancies

##### P1 — `Ambush` card type incorrect

- fighter: `raptors`
- card: `ambush`
- severity: **P1**
- Git location: `docs/cards/phase-4a/raptors.yaml`, `ambush`
- expected from physical card: **ATTACK**, printed value 2, BOOST 3.
- observed in Git: `type: versatile`, printed value 2, BOOST 3.
- evidence/reasoning: the physical component is Attack-only. `versatile` incorrectly permits the card to be committed as a defense, changing legal gameplay.
- expected semantic correction: canonical card type must be `attack`; value, BOOST and during-combat random-discard/BOOST effect remain unchanged.

##### P1 — `Eaten Alive` card type incorrect

- fighter: `raptors`
- card: `eaten-alive`
- severity: **P1**
- Git location: `docs/cards/phase-4a/raptors.yaml`, `eaten-alive`
- expected from physical card: **VERSATILE**, printed value 4, BOOST 2.
- observed in Git: `type: attack`, printed value 4, BOOST 2.
- evidence/reasoning: the physical component can legally be committed on attack or defense. `attack` removes its legal defensive use.
- expected semantic correction: canonical card type must be `versatile`; existing value, BOOST and after-combat effect remain unchanged.

#### Integration requirements confirmed

None.

#### Verdict

**FAIL**

Reason: two independent P1 card-type discrepancies.

---

### robert-muldoon

Canonical source:

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/robert-muldoon.yaml`
- fighter manifest: `docs/fighters/phase-4b/robert-muldoon.yaml`
- construction: `fixed`
- canonical available pool: 30
- canonical game deck: 30

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: PASS
- manifest readable: PASS
- images readable: 11/11
- unique images: 11
- manifest card entries: 11
- duplicate binary images: 0
- missing images: 0
- extra images: 0

#### Canonical manifest comparison

- archive quantity sum: 30
- Git available pool: 30
- Git game deck: 30
- unique archive cards: 11
- unique Git definitions: 11
- quantity validation: **PASS**

#### Unique-card image completeness

All 11 unique physical card definitions were present and inspected.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Leap Away | PASS | PASS | PASS | PASS |
| Rending Shot | PASS | PASS | PASS | PASS |
| Second Shot | PASS | PASS | PASS | PASS |
| Call for Backup | PASS | PASS | PASS | PASS |
| Shoot Her! | PASS | PASS | PASS | PASS |
| Tactical Advance | PASS | PASS | PASS | PASS |
| I've Hunted Most Things That Can Hunt You | PASS | P3 archive mapping | PASS | PASS_WITH_QUALIFICATION |
| They Should All Be Destroyed | PASS | PASS | PASS | PASS |
| Remote Detonation | PASS | PASS | PASS / source-dependent follow-up qualification | PASS_WITH_QUALIFICATION |

High-risk semantics explicitly checked:

- `Second Shot`: optional BOOST from a chosen hand card is preserved.
- `Call for Backup`: choose exactly two distinct effects from the three printed options. Placement is `up to 3` traps; defeated Workers may be returned; draw option draws 2. Git's staged choice representation preserves two distinct selected options.
- `They Should All Be Destroyed`: combat value depends on Muldoon traps adjacent to the opposing fighter.
- `Remote Detonation`: physical card directly confirms choosing a trap in a zone containing an InGen Worker, damaging each opposing fighter adjacent to it, and returning that trap to the box.

`Remote Detonation` qualification: Git also materializes a subsequent draw when the trap returns to the box. That draw is derived from Muldoon's source-defined trap ability rather than printed on the `Remote Detonation` card image itself. The card image therefore confirms the trap-return event that invokes the integration behavior, but the fighter ability text itself was not present in this evidence ZIP and was not claimed as visually confirmed by this QA.

#### Discrepancies

##### P3 — archive mapping ID differs for `I've Hunted Most Things That Can Hunt You`

- fighter: `robert-muldoon`
- card: `I've Hunted Most Things That Can Hunt You`
- severity: **P3**
- expected canonical card ID: `ive-hunted-most-things-that-can-hunt-you`
- observed archive manifest ID: `i-ve-hunted-most-things-that-can-hunt-you`
- evidence/reasoning: archive mapping introduces an extra separator after `i`; physical identity and gameplay data otherwise match canonical Git.
- expected correction: evidence/archive mapping should use the canonical Git ID. No gameplay manifest correction is implied.

#### Integration requirements confirmed

**`A-REQ-003` confirmed directly by card images.**

The physical card corpus requires persistent trap-component behavior:

- `Remote Detonation` selects and returns a trap.
- `Call for Backup` places up to three traps.
- `They Should All Be Destroyed` counts traps adjacent to the opposing fighter.

Canonical Git marks these card operations as requiring `A-REQ-003`.

`A-REQ-004` is not marked as image-confirmed by this QA. It is associated with interruption of opposing movement when a trap triggers, which belongs to the fighter ability/rules layer rather than the card images present in this ZIP.

#### Verdict

**PASS_WITH_QUALIFICATIONS**

No P1/P2 card discrepancy. Qualifications are the shared trap-component integration dependency, source-dependent draw behavior associated with trap return, and one P3 archive mapping mismatch.

## Findings

### Material findings

#### F-01 — Bruce Lee / One-Inch Punch target domain

- fighter: `bruce-lee`
- card: `one-inch-punch`
- severity: **P1**
- expected: damage target is any adjacent fighter.
- observed: Git restricts the target to `adjacent_opposing_fighters`.
- impact: removes otherwise legal targets; gameplay-changing.
- evidence/reasoning: physical card wording was read directly from the provided image; Git structured choice domain was compared against it.

#### F-02 — Raptors / Ambush type

- fighter: `raptors`
- card: `ambush`
- severity: **P1**
- expected: ATTACK 2, BOOST 3.
- observed: Git stores `type: versatile`, value 2, BOOST 3.
- impact: incorrectly enables defensive use.
- evidence/reasoning: card type icon/printed component was visually checked rather than inferred from filename or archive manifest.

#### F-03 — Raptors / Eaten Alive type

- fighter: `raptors`
- card: `eaten-alive`
- severity: **P1**
- expected: VERSATILE 4, BOOST 2.
- observed: Git stores `type: attack`, value 4, BOOST 2.
- impact: incorrectly removes defensive use.
- evidence/reasoning: card type icon/printed component was visually checked directly.

### Non-material findings and warnings

#### F-04 — Bruce Lee evidence mapping mismatch

- fighter: `bruce-lee`
- card: `"HOO! WHAAAAAA!"`
- severity: **P3**
- expected archive mapping ID: `hoo-whaaaaa`
- observed: `hoo-whaaaaaa`
- impact: mapping/metadata only; printed card and gameplay semantics match.

#### F-05 — Robert Muldoon evidence mapping mismatch

- fighter: `robert-muldoon`
- card: `I've Hunted Most Things That Can Hunt You`
- severity: **P3**
- expected archive mapping ID: `ive-hunted-most-things-that-can-hunt-you`
- observed: `i-ve-hunted-most-things-that-can-hunt-you`
- impact: mapping/metadata only; printed card and gameplay semantics match.

#### W-01 — Robert Muldoon source-dependent trap return draw

- fighter: `robert-muldoon`
- card: `remote-detonation`
- severity: qualification / evidence boundary, not a discrepancy
- observation: the card image itself confirms trap return but does not print the follow-up draw that Git resolves via Muldoon's trap ability.
- reasoning: this QA intentionally does not claim visual confirmation of a fighter ability component not present in the evidence ZIP.

#### W-02 — Raptors Disengage MOVE/PLACE wording

- fighter: `raptors`
- card: `disengage`
- severity: qualification / normalization, not a discrepancy
- observation: physical card wording uses movement language; Git models the resolved behavior as placement to an empty space in the Raptor's zone.
- reasoning: the normalization is treated as an errata-aware/source-defined gameplay representation rather than literal transcription, consistent with the QA contract's instruction to compare resulting legal behavior rather than require prose equality.

## Corpus-Level Observations

1. **Technical corpus integrity is clean.** All five nested archives, all five manifests, and all 63 unique card images were readable. There are no zero-byte images, corrupted images, missing image mappings, extra images, or duplicate binary card images in this batch.

2. **Quantity construction is clean.** Every fighter in this batch reconciles to 30 physical deck copies and fixed construction in canonical Git. No fighter in this batch uses a non-standard choose-pool/external-definition deck construction.

3. **Visual inspection was exhaustive.** All 63 unique images were inspected. Generic definitions such as Feint, Regroup, Momentous Shift, Ambush, and Disengage were not skipped merely because their effects are familiar.

4. **Archive manifests were treated as mapping evidence only.** The two P3 ID mismatches demonstrate why archive identifiers cannot be treated as canonical. Printed image identity plus canonical Git remains authoritative for QA comparison.

5. **Structured normalization generally held up well.** Staged effects, combat-result branches, optional operations, card replacement, actual-damage capture, trap component operations, and return-vs-recover distinctions generally preserved the printed gameplay behavior even when Git structure differed substantially from card prose.

6. **Card type is a high-risk transcription field.** Both material Raptors findings are swapped ATTACK/VERSATILE classifications. These errors are small in YAML surface area but directly change when the card may legally be played.

7. **Target-domain normalization is also high risk.** Bruce Lee's `One-Inch Punch` shows that changing `fighter` to `opposing fighter` is not a harmless normalization. Relation words must be preserved exactly when they determine legal target sets.

8. **Integration requirements are not transcription failures.** Muldoon's `A-REQ-003` requirement is directly justified by card components that manipulate trap tokens. This is an engine/integration dependency, not a mismatch between printed evidence and canonical card semantics.

9. **Fighter-ability evidence boundaries were respected.** The QA did not claim that fighter abilities themselves were visually verified when the relevant ability card/component was absent from the evidence ZIP.

## Final Assessment

The batch is **technically complete and usable as evidence**, but it is **not clean enough for unconditional Phase 4 corpus acceptance** because three gameplay-changing P1 discrepancies remain in canonical Git.

Per-fighter assessment:

| Fighter | Verdict | Material findings |
|---|---|---:|
| `bigfoot` | PASS | 0 |
| `bruce-lee` | FAIL | 1 P1 |
| `dracula` | PASS | 0 |
| `raptors` | FAIL | 2 P1 |
| `robert-muldoon` | PASS_WITH_QUALIFICATIONS | 0 P1/P2 |

Clean of material P1/P2 discrepancies:

- `bigfoot`
- `dracula`
- `robert-muldoon`

Material correction is required before the batch can be considered fully reconciled:

- restore `One-Inch Punch` to an unrestricted adjacent-fighter target domain;
- classify Raptors `Ambush` as ATTACK;
- classify Raptors `Eaten Alive` as VERSATILE.

No blocker prevented complete QA of any fighter in this batch.
