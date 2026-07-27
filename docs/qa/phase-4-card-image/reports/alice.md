# Phase 4 Card-Image QA Report

## Batch

- alice
- king-arthur
- medusa
- robin-hood
- sinbad

Evidence archive: `unmatched-bundle-01.zip`

Canonical repository: `NordCoder/unmatched-web`

Canonical branch for every fighter in this batch: `phase-4b-worker-a-classics`

Branch tip immediately before report persistence: `840a03572bf8b5912b03f32e5ae2a48cc8fadc2b`

## Verdict

**FAIL**

## Summary

The batch contains five fighter evidence ZIPs and 69 unique card images. The outer archive and all nested fighter archives were readable, all 69 images decoded successfully, every archive manifest was readable, and there were no zero-byte images, binary duplicate images, missing images, extra images, manifest entries without images, or images without manifest entries.

Every unique card image was inspected visually rather than inferred from filename or archive manifest. Archive quantities were compared to the current canonical Phase 4B deck manifests on `phase-4b-worker-a-classics`. All five decks reconcile to their canonical fixed 30-card constructions: archive quantity sum 30, Git `available_pool_count` 30, and Git `game_deck_count` 30 for every fighter.

Two material gameplay discrepancies were found in canonical Git semantics:

1. `medusa / Winged Frenzy`: the printed return of a defeated Harpy is mandatory when an eligible defeated Harpy exists, while Git models the Harpy choice as optional.
2. `robin-hood / Defenders of Sherwood`: the printed return of a defeated Outlaw is mandatory when an eligible defeated Outlaw exists, while Git models the Outlaw choice as optional.

Both are **P1** because the current representation permits a legal choice that the physical card does not permit.

Three non-gameplay evidence-mapping findings were also found. These concern archive `card_id` slug normalization only and do not indicate printed metadata or gameplay differences in Git:

- Alice: `i-m-late-i-m-late` vs canonical `im-late-im-late`.
- Sinbad: six Voyage definitions use longer evidence slugs than the canonical Git IDs.
- Robin Hood: `a-hunter-s-eye` vs canonical `a-hunters-eye`.

King Arthur has no transcription discrepancy. The physical card evidence directly confirms the already-documented integration qualifications `A-REQ-006` and `A-REQ-013`.

## Batch Integrity

| Check | Result |
|---|---:|
| Fighter ZIPs received | 5 |
| Fighter ZIPs fully checked | 5 |
| Nested ZIPs readable | 5/5 |
| Manifests readable | 5/5 |
| Total unique card images | 69 |
| Images successfully inspected | 69 |
| Unreadable images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Duplicate binary images | 0 |
| Quantity failures | 0 |
| Printed metadata discrepancies in Git | 0 |
| Semantic discrepancies | 2 P1 |
| Archive card-ID mapping findings | 3 P3 findings affecting 8 cards |
| Blocked fighters | 0 |

Image formats were accepted based on successful decoding, not extension. In particular, `medusa/dash.webp` and `sinbad/leap-away.webp` decoded and were inspected normally.

## Canonical Git State Used

The final QA re-read the current branch ref before persistence. The relevant canonical deck blobs were:

| Fighter | Deck manifest | Blob SHA | Canonical state |
|---|---|---|---|
| alice | `docs/cards/phase-4b/alice.yaml` | `168dc889428e523a40f1dcda21de39dd0b8e67f3` | verified; fixed 30/30 |
| king-arthur | `docs/cards/phase-4b/king-arthur.yaml` | `3f77842baa4d6474e74b20c22525f37ef7f773d1` | partial; evidence/semantics verified; fixed 30/30 |
| medusa | `docs/cards/phase-4b/medusa.yaml` | `0c9740996da8d0faf38308a677711076f46dba70` | verified; fixed 30/30 |
| robin-hood | `docs/cards/phase-4b/robin-hood.yaml` | `a6255dc99decf658af7f1ba21c09a07969a1e142` | verified; fixed 30/30 |
| sinbad | `docs/cards/phase-4b/sinbad.yaml` | `cdcc9f960a0377f5ee3bde2da2dbef6a4033d688` | verified; fixed 30/30 |

Relevant fighter-manifest blobs were:

| Fighter | Fighter manifest | Blob SHA |
|---|---|---|
| alice | `docs/fighters/phase-4b/alice.yaml` | `0b44408b61d2571dc05f3463b5598aa4665078a0` |
| king-arthur | `docs/fighters/phase-4b/king-arthur.yaml` | `3db0196bcf6d08cd346754f66447f6dd325d238f` |
| medusa | `docs/fighters/phase-4b/medusa.yaml` | `95a2d3d38c71d4c52661582cb6b347c34bc5692c` |
| robin-hood | `docs/fighters/phase-4b/robin-hood.yaml` | `afc54604a0efc20f2139ee5b9a0447a2eb7f8a9f` |
| sinbad | `docs/fighters/phase-4b/sinbad.yaml` | `4faa8d855ee0f7274a6bdf41cc6b08bfe3ae5db5` |

## Fighter Results

### alice

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/alice.yaml`
- fighter manifest: `docs/fighters/phase-4b/alice.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS — 15/15**
- printed card identity/content verification: **PASS**
- discrepancies: one P3 archive `card_id` normalization finding; no material gameplay discrepancy
- integration requirements confirmed: none
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive and quantity evidence

| Metric | Result |
|---|---:|
| Images readable | 15/15 |
| Unique images | 15 |
| Manifest card entries | 15 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 15 |
| Unique Git definitions | 15 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Looking Glass | PASS | PASS | PASS | PASS |
| Snicker-Snack | PASS | PASS | PASS | PASS |
| O Frabjous Day! | PASS | PASS | PASS | PASS |
| The Other Side of the Mushroom | PASS | PASS | PASS | PASS |
| Eat Me | PASS | PASS | PASS | PASS |
| I'm Late, I'm Late | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Drink Me | PASS | PASS | PASS | PASS |
| Jaws That Bite | PASS | PASS | PASS | PASS |
| Mad as a Hatter | PASS | PASS | PASS | PASS |
| Manxome Foe | PASS | PASS | PASS | PASS |
| Claws That Catch | PASS | PASS | PASS | PASS |

The visual comparison confirmed the printed timing windows, quantities, fighter restrictions, combat types, printed values, BOOST values, Alice-size state changes, movement/place distinction, hand inspection/discard behavior, choice structure, and value-modification behavior represented by the canonical manifest.

### king-arthur

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/king-arthur.yaml`
- fighter manifest: `docs/fighters/phase-4b/king-arthur.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS — 16/16**
- printed card identity/content verification: **PASS**
- discrepancies: none
- integration requirements confirmed: `A-REQ-006`, `A-REQ-013`
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive and quantity evidence

| Metric | Result |
|---|---:|
| Images readable | 16/16 |
| Unique images | 16 |
| Manifest card entries | 16 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 16 |
| Unique Git definitions | 16 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Swift Strike | PASS | PASS | PASS | PASS |
| Bewilderment | PASS | PASS | PASS | PASS |
| Noble Sacrifice | PASS | PASS | PASS | PASS |
| Excalibur | PASS | PASS | PASS | PASS |
| The Aid of Morgana | PASS | PASS | PASS | PASS |
| Divine Intervention | PASS | PASS | PASS | PASS |
| The Holy Grail | PASS | PASS | PASS | PASS |
| The Lady of the Lake | PASS | PASS | PASS | PASS |
| Prophecy | PASS | PASS | PASS | PASS |
| Aid the Chosen One | PASS | PASS | PASS | PASS |
| Restless Spirits | PASS | PASS | PASS | PASS |
| Command the Storms | PASS | PASS | PASS | PASS |

The physical cards confirm the normalized representations of Merlin placement, moving all fighters, Prophecy's top-four look/select/reorder procedure, Restless Spirits targeting, Noble Sacrifice BOOST behavior, exact health assignment by The Holy Grail, and structured Excalibur search by The Lady of the Lake.

`A-REQ-006` is directly supported by The Holy Grail: under its printed condition, Arthur's health is assigned exactly to 8 rather than recovered by an ordinary amount.

`A-REQ-013` is directly supported by The Lady of the Lake: the effect searches Arthur's deck or discard for Excalibur, moves a found Excalibur to hand, and requires the deck-search/shuffle behavior represented by the structured search requirement.

### medusa

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/medusa.yaml`
- fighter manifest: `docs/fighters/phase-4b/medusa.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL — one P1 semantic discrepancy**
- unique-card image completeness: **PASS — 11/11**
- printed card identity/content verification: **FAIL on Winged Frenzy; 10 other cards PASS**
- discrepancies: `MEDUSA-WINGED-FRENZY-OPTIONALITY` P1
- integration requirements confirmed: none
- verdict: **FAIL**

#### Archive and quantity evidence

| Metric | Result |
|---|---:|
| Images readable | 11/11 |
| Unique images | 11 |
| Manifest card entries | 11 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 11 |
| Unique Git definitions | 11 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Dash | PASS | PASS | PASS | PASS |
| Gaze of Stone | PASS | PASS | PASS | PASS |
| A Momentary Glance | PASS | PASS | PASS | PASS |
| Hiss and Slither | PASS | PASS | PASS | PASS |
| The Hounds of Mighty Zeus | PASS | PASS | PASS | PASS |
| Clutching Claws | PASS | PASS | PASS | PASS |
| Winged Frenzy | PASS | PASS | **FAIL** | **FAIL** |
| Snipe | PASS | PASS | PASS | PASS |
| Second Shot | PASS | PASS | PASS | PASS |

### robin-hood

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/robin-hood.yaml`
- fighter manifest: `docs/fighters/phase-4b/robin-hood.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL — one P1 semantic discrepancy; one P3 evidence-ID warning**
- unique-card image completeness: **PASS — 11/11**
- printed card identity/content verification: **FAIL on Defenders of Sherwood; 10 other cards PASS**
- discrepancies: one P1 mandatory-vs-optional return discrepancy and one P3 archive `card_id` normalization finding
- integration requirements confirmed: none
- verdict: **FAIL**

#### Archive and quantity evidence

| Metric | Result |
|---|---:|
| Images readable | 11/11 |
| Unique images | 11 |
| Manifest card entries | 11 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 11 |
| Unique Git definitions | 11 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Wily Fighting | PASS | PASS | PASS | PASS |
| Ambush | PASS | PASS | PASS | PASS |
| Disarming Shot | PASS | PASS | PASS | PASS |
| Piercing Shot | PASS | PASS | PASS | PASS |
| Snark | PASS | PASS | PASS | PASS |
| A Hunter's Eye | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Steal from the Rich | PASS | PASS | PASS | PASS |
| Defenders of Sherwood | PASS | PASS | **FAIL** | **FAIL** |
| Highway Robbery | PASS | PASS | PASS | PASS |

### sinbad

- branch: `phase-4b-worker-a-classics`
- deck manifest: `docs/cards/phase-4b/sinbad.yaml`
- fighter manifest: `docs/fighters/phase-4b/sinbad.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS — 16/16**
- printed card identity/content verification: **PASS**
- discrepancies: one grouped P3 archive `card_id` normalization finding affecting six Voyage definitions; no material gameplay discrepancy
- integration requirements confirmed: none
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive and quantity evidence

| Metric | Result |
|---|---:|
| Images readable | 16/16 |
| Unique images | 16 |
| Manifest card entries | 16 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 16 |
| Unique Git definitions | 16 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Commanding Impact | PASS | PASS | PASS | PASS |
| Leap Away | PASS | PASS | PASS | PASS |
| Exploit | PASS | PASS | PASS | PASS |
| Toil and Danger | PASS | PASS | PASS | PASS |
| Voyage Home | PASS | PASS | PASS | PASS |
| Riches Beyond Compare | PASS | PASS | PASS | PASS |
| By Fortune and Fate | PASS | PASS | PASS | PASS |
| Voyage to the Island That Was a Whale | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Voyage to the Valley of the Giant Snakes | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Voyage to the Creature with Eyes Like Coals of Fire | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Voyage to the Cannibals with the Root of Madness | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Voyage to the City of the Man-Eating Apes | PASS | PASS; P3 ID mapping | PASS | PASS* |
| Voyage to the City of the King of Serendib | PASS | PASS; P3 ID mapping | PASS | PASS* |

The printed Voyage combat-value rule and all distinct Voyage AFTER COMBAT effects were checked independently. The canonical semantics correctly preserve recovery, private hand inspection, random discard, optional movement, effect damage, draw behavior, and Voyage Home's return of other Voyage cards from discard.

## Findings

### MEDUSA-WINGED-FRENZY-OPTIONALITY

- fighter: `medusa`
- card: `Winged Frenzy`
- severity: **P1**
- Git location: `docs/cards/phase-4b/medusa.yaml`, `winged-frenzy-effect / choose-defeated-harpy`
- expected: after the movement portion resolves, if at least one defeated Harpy exists, the player must return one eligible defeated Harpy to a legal space in Medusa's zone; the return is not a player-declinable option
- observed: Git correctly guards the stage with `selector_nonempty` for a defeated Harpy, but `choose-defeated-harpy` is represented as `optional: true`
- evidence/reasoning: the printed card conditionally instructs the return when a defeated Harpy exists and does not use `may`; therefore optionality changes the legal behavior rather than merely changing YAML structure
- expected semantic correction: retain the existing defeated-Harpy existence condition, but make selection/return mandatory whenever the eligible set is nonempty

### ROBIN-HOOD-DEFENDERS-OPTIONALITY

- fighter: `robin-hood`
- card: `Defenders of Sherwood`
- severity: **P1**
- Git location: `docs/cards/phase-4b/robin-hood.yaml`, `defenders-effect / choose-defeated-outlaw`
- expected: after drawing the card, if at least one defeated Outlaw exists, one eligible defeated Outlaw must be returned to a legal space in Robin Hood's zone
- observed: Git correctly guards the stage with a defeated-Outlaw existence condition, but `choose-defeated-outlaw` is represented as `optional: true`
- evidence/reasoning: the printed instruction does not grant a `may`; Git therefore permits a decline that the physical card does not permit
- expected semantic correction: retain the existing no-Outlaw guard, but require the selection and return whenever at least one eligible defeated Outlaw exists

### ALICE-IM-LATE-CARD-ID

- fighter: `alice`
- card: `I'm Late, I'm Late`
- severity: **P3**
- expected: evidence mapping resolves to canonical card ID `im-late-im-late`
- observed: archive manifest uses `i-m-late-i-m-late`
- evidence/reasoning: printed name, quantity, user restriction, type, value, BOOST and gameplay effect all map unambiguously to the canonical Git definition; only the evidence slug differs
- expected correction: normalize the evidence mapping or maintain an explicit card-ID alias; no canonical gameplay change is indicated

### SINBAD-VOYAGE-CARD-IDS

- fighter: `sinbad`
- cards:
  - `Voyage to the Island That Was a Whale`
  - `Voyage to the Valley of the Giant Snakes`
  - `Voyage to the Creature with Eyes Like Coals of Fire`
  - `Voyage to the Cannibals with the Root of Madness`
  - `Voyage to the City of the Man-Eating Apes`
  - `Voyage to the City of the King of Serendib`
- severity: **P3**
- expected canonical IDs:
  - `voyage-whale-island`
  - `voyage-giant-snakes`
  - `voyage-eyes-like-coals`
  - `voyage-cannibals-root-of-madness`
  - `voyage-man-eating-apes`
  - `voyage-king-of-serendib`
- observed: the evidence manifest uses longer The-Unmatched-Club-style slugs beginning `voyage-to-the-...`
- evidence/reasoning: all six printed card identities and gameplay facts map unambiguously to the canonical definitions; no printed quantity, restriction, type, value, BOOST or effect-semantic difference was found
- expected correction: normalize evidence mappings or explicitly maintain these aliases; no canonical gameplay change is indicated

### ROBIN-HOOD-HUNTERS-EYE-CARD-ID

- fighter: `robin-hood`
- card: `A Hunter's Eye`
- severity: **P3**
- expected: evidence mapping resolves to canonical card ID `a-hunters-eye`
- observed: archive manifest uses `a-hunter-s-eye`
- evidence/reasoning: printed name, Robin Hood restriction, attack value 5, BOOST 4 and quantity 3 match the canonical definition; only the evidence slug differs
- expected correction: normalize the evidence mapping or maintain an explicit alias; no canonical gameplay change is indicated

## Integration Requirements Confirmed

### A-REQ-006 — King Arthur / The Holy Grail

The physical `The Holy Grail` card confirms the need to distinguish exact health assignment from ordinary recovery. When its printed condition is satisfied, Arthur's current health becomes exactly 8. The existing Git representation using `SET_HEALTH` preserves this semantic correctly; this is an integration requirement, not a transcription failure.

### A-REQ-013 — King Arthur / The Lady of the Lake

The physical `The Lady of the Lake` card confirms the need for a structured multi-zone card search that binds whether Excalibur was found and from which zone, moves the found card to hand, and supports the source-defined post-search shuffle behavior. The existing Git semantics preserve the printed behavior; this is an integration requirement, not a transcription failure.

## Corpus-Level Observations

1. The evidence archive follows the required one-unique-definition-to-one-image model. Physical copies are represented by manifest `quantity`; duplicate physical copies are not represented as duplicate image files.
2. All five archive quantity sums are 30 and agree with the current canonical fixed deck constructions. No evidence in this batch suggests treating 30 as a universal project invariant; it is simply correct for these five fighters.
3. Archive manifests were useful as image-to-card mappings but were not treated as authoritative. Three slug-normalization findings demonstrate why card identity must ultimately be resolved against image evidence and canonical Git definitions.
4. Generic cards such as Feint, Regroup, Skirmish, Momentous Shift, Dash and other shared definitions were still inspected on their actual images; they were not automatically accepted based on familiarity.
5. File extension was not used as a validity criterion. WEBP images in the batch were decoded and visually checked successfully.
6. The QA compared resulting gameplay behavior rather than requiring literal copyrighted card prose in Git. Structured `stages`, choices, conditions, result bindings and operations were accepted where they preserved the printed timing, ordering, optionality, targets and dependencies.
7. The two material findings are the same semantic class: a source-mandatory conditional return was encoded as an optional player choice. This is precisely the type of distinction the normalized effect model is intended to preserve.
8. No fighter was blocked. Every image was readable enough to inspect, and every required canonical Git source was available.
9. No repository state was modified during the independent verification itself. This report persistence commit is documentation-only and does not correct either finding.

## Final Assessment

This evidence batch is technically complete and suitable as a Phase 4 physical card-image evidence corpus: all five fighter archives are intact, all 69 unique card components are present and readable, quantities reconcile, and no physical-card identity or printed metadata is missing.

The **canonical repository corpus is not fully acceptable as-is**, however, because two physical cards expose P1 semantic mismatches in the current normalized Git representation:

- Medusa — `Winged Frenzy`;
- Robin Hood — `Defenders of Sherwood`.

Those findings make the batch-level QA verdict **FAIL** until the canonical semantics are corrected and independently revalidated. Alice, King Arthur and Sinbad have no material card-semantic discrepancy in this batch. The P3 evidence-ID mismatches do not affect gameplay, and King Arthur's `A-REQ-006` / `A-REQ-013` qualifications are confirmed integration requirements rather than card transcription errors.

No fixes were applied by this QA worker.
