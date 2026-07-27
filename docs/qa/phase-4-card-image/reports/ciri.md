# Phase 4 Card-Image QA Report

## Batch

- ciri
- geralt-of-rivia
- hamlet
- titania
- wayward-sisters

## Verdict

FAIL

## Summary

Independent read-only Phase 4 card-image QA was completed for a five-fighter evidence batch containing `hamlet`, `titania`, `wayward-sisters`, `geralt-of-rivia`, and `ciri`.

The outer archive and all five nested fighter ZIPs were readable and passed CRC/integrity checks. All 62 supplied unique card images decoded successfully and were visually inspected; no supplied image was skipped as a sample, unreadable, zero-byte, truncated, or binary-duplicated.

Canonical documentation was compared against the repository state appropriate to each fighter:

- `wayward-sisters`, `geralt-of-rivia`: Phase 4A manifests on `main`.
- `hamlet`, `titania`, `ciri`: Phase 4B Worker C manifests on `phase-4b-worker-c-modern`.

The batch-level verdict is **FAIL** because Titania has one gameplay-changing semantic discrepancy in `Gift Of The Fair Folk` and the evidence archive omits all six canonical auxiliary Glamour card images. The other four fighters have no P1/P2 card-data discrepancies. Hamlet and Geralt each have only a P3 archive `card_id` normalization mismatch; Ciri has integration qualifications only.

### Batch metrics

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters with every canonical gameplay-card image available | 4 |
| Supplied unique card images | 62 |
| Images successfully inspected | 62/62 |
| Unreadable images | 0 |
| Corrupted images | 0 |
| Zero-byte images | 0 |
| Duplicate binary images | 0 |
| Canonical gameplay-card images missing | 6 |
| Quantity/construction failures | 1 evidence-completeness failure |
| P1 gameplay discrepancies | 1 |
| P2 evidence discrepancies | 1 |
| P3 metadata discrepancies | 2 |

## Fighter Results

### hamlet

- archive integrity: **PASS**
- canonical manifest comparison: **PASS WITH QUALIFICATIONS**
- unique-card image completeness: **12/12 supplied, 12/12 canonical action-card definitions**
- printed card identity/content verification: **all 12 images inspected; no P1/P2 mismatch**
- discrepancies: one P3 archive card-ID normalization mismatch; `The Readiness Is All` remains subject to an already documented project normalization which the image does not contradict
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive / construction counts

| Check | Result |
|---|---:|
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
| Unique Git definitions | 12 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Blood Will Have Blood | PASS | PASS | PASS | PASS |
| Cruel To Be Kind | PASS | PASS | PASS | PASS |
| Maddening Insight | PASS | PASS | PASS | PASS |
| Method In The Madness | PASS | PASS | PASS | PASS |
| Nothing Either Good Or Bad | PASS | PASS | PASS | PASS |
| Outrageous Fortune | PASS | PASS | PASS | PASS |
| The Ghost | PASS | PASS | PASS | PASS |
| The Play's The Thing | PASS | P3 | PASS | PASS |
| The Readiness Is All | PASS | PASS | PASS* | PASS |
| The Rest Is Silence | PASS | PASS | PASS | PASS |
| To Sleep, Perchance To Dream | PASS | PASS | PASS | PASS |
| Uncertain Doom | PASS | PASS | PASS | PASS |

`The Readiness Is All` is marked `PASS*` because canonical Git deliberately carries a project normalization for failed PLACE: the cannot-leave restriction is attached only after successful placement. The physical card does not contradict that interpretation; this remains a documented evidence/integration qualification, not a transcription discrepancy.

#### Integration requirements confirmed

- `C-REQ-013` — staged resolution / operation-result dependency is directly supported by the printed card corpus.

---

### titania

- archive integrity: **FAIL relative to canonical gameplay corpus completeness**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **12/12 ordinary action-card images supplied, but 6 canonical auxiliary Glamour cards are absent**
- printed card identity/content verification: **12 supplied action-card faces all inspected; one P1 semantic mismatch found**
- discrepancies: `TITANIA-001` P1 and `TITANIA-002` P2
- verdict: **FAIL**

#### Archive / construction counts

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Supplied images readable | 12/12 |
| Unique supplied images | 12 |
| Archive manifest entries | 12 |
| Duplicate images | 0 |
| Missing vs archive manifest | 0 |
| Missing vs canonical gameplay corpus | 6 |
| Extra images | 0 |
| Archive action-card quantity sum | 30 |
| Git action-card available pool | 30 |
| Git game deck | 30 |
| Unique archive action cards | 12 |
| Unique Git action-card definitions | 12 |
| Canonical auxiliary Glamour definitions | 6 |
| Canonical gameplay-card definitions total | 18 |
| Construction | fixed 30-card action deck + separate 6-card auxiliary Glamour system |

The six canonical auxiliary gameplay-card definitions absent from the archive are:

- Glamour of Greed
- Glamour of Invisibility
- Glamour of Jealousy
- Glamour of Love
- Glamour of Rhyme
- Glamour of Sleep

They must not be counted as ordinary action-deck copies. Canonical Git intentionally models them as auxiliary gameplay-card definitions/instances in a separate Glamour deck lifecycle.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| A Momentary Glance | PASS | PASS | PASS | PASS |
| As Wise As Beautiful | PASS | PASS | PASS | PASS |
| But A Dream | PASS | PASS | PASS | PASS |
| Fairy Song | PASS | PASS | PASS | PASS |
| Gift Of The Fair Folk | PASS | PASS | FAIL | FAIL |
| Met By Moonlight | PASS | PASS | PASS | PASS |
| Parting Gift | PASS | PASS | PASS | PASS |
| Protection Of The Fairy Woods | PASS | PASS | PASS | PASS |
| Queen Of The Fairies | PASS | PASS | PASS | PASS |
| The Moon Looks Down | PASS | PASS | PASS | PASS |
| What Fools These Mortals Be | PASS | PASS | PASS | PASS |
| Whisked Away | PASS | PASS | PASS | PASS |
| Glamour of Greed | MISSING | — | — | NOT VERIFIED |
| Glamour of Invisibility | MISSING | — | — | NOT VERIFIED |
| Glamour of Jealousy | MISSING | — | — | NOT VERIFIED |
| Glamour of Love | MISSING | — | — | NOT VERIFIED |
| Glamour of Rhyme | MISSING | — | — | NOT VERIFIED |
| Glamour of Sleep | MISSING | — | — | NOT VERIFIED |

#### Integration requirements confirmed

- `C-REQ-002` — supplied action cards directly reference the active/discarded Glamour system and therefore corroborate the auxiliary-card-zone mechanic, although the six Glamour faces themselves were not supplied.
- `C-REQ-013` — dependent discard/cost → consequence structures are directly visible on cards such as `Fairy Song` and `Whisked Away`.

`C-REQ-003` and `C-REQ-011` were not claimed as independently image-confirmed in this QA because the missing Glamour cards carry relevant behavior.

---

### wayward-sisters

- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **12/12 supplied and complete for the ordinary action deck**
- printed card identity/content verification: **all 12 images inspected; no discrepancy**
- discrepancies: none
- verdict: **PASS**

#### Archive / construction counts

| Check | Result |
|---|---:|
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
| Unique Git action-card definitions | 12 |
| External spell definitions | 4, correctly not counted as ordinary cards |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| All-Seeing Familiar | PASS | PASS | PASS | PASS |
| Curious Familiar | PASS | PASS | PASS | PASS |
| Double, Double | PASS | PASS | PASS | PASS |
| Fire Burn And Cauldron Bubble | PASS | PASS | PASS | PASS |
| Hurly-Burly | PASS | PASS | PASS | PASS |
| Pricking Of My Thumbs | PASS | PASS | PASS | PASS |
| Prophecy | PASS | PASS | PASS | PASS |
| Something Wicked This Way Comes | PASS | PASS | PASS | PASS |
| The Stars Align | PASS | PASS | PASS | PASS |
| Toil And Trouble | PASS | PASS | PASS | PASS |
| Unnatural Remedy | PASS | PASS | PASS | PASS |
| Ward | PASS | PASS | PASS | PASS |

Sensitive distinctions checked successfully include:

- `All-Seeing Familiar`: private LOOK AT semantics, not public REVEAL.
- `Hurly-Burly`: reveal → opponent disposition choice → conditional combat-value addition.
- `Something Wicked This Way Comes`: two different spells without ordinary ingredient requirement/consumption semantics.
- `Unnatural Remedy`: up to three cards from the cauldron and a recovery choice for each moved card.
- `Prophecy`: LOOK_AT / selection / reorder semantics rather than ordinary DRAW semantics.

Wayward Sisters spells are external gameplay definitions rather than action-card instances and were therefore correctly excluded from the one-image-per-action-card corpus count.

---

### geralt-of-rivia

- archive integrity: **PASS**
- canonical manifest comparison: **PASS WITH QUALIFICATIONS**
- unique-card image completeness: **15/15 definitions in the 36-card available pool**
- printed card identity/content verification: **all 15 images inspected; no P1/P2 mismatch**
- discrepancies: one P3 archive card-ID normalization mismatch
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive / construction counts

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 15/15 |
| Unique images | 15 |
| Manifest card entries | 15 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum / available pool | 36 |
| Git available pool | 36 |
| Git base quantity | 24 |
| Selectable gear quantity | 12 |
| Selected gear cards per game | 6 |
| Git game deck | 30 |
| Unique Git definitions | 15 |
| Construction | choose_groups |

The non-standard Geralt construction was explicitly validated rather than treated as a 30-card archive expectation. Canonical construction is:

- 24 base cards always included;
- choose exactly one Potion definition, both copies;
- choose exactly one Armor definition, both copies;
- choose exactly one Sword definition, both copies;
- 6 selected Gear cards are added to the base 24 to form a 30-card game deck.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Annoying Tune | PASS | PASS | PASS | PASS |
| Damn, You're Ugly | PASS | P3 | PASS | PASS |
| Disciplined Duelist | PASS | PASS | PASS | PASS |
| GEAR: Armor of the Forgotten Wolf | PASS | PASS | PASS | PASS |
| GEAR: Blizzard | PASS | PASS | PASS | PASS |
| GEAR: Sword of Silver | PASS | PASS | PASS | PASS |
| GEAR: Sword of Steel | PASS | PASS | PASS | PASS |
| GEAR: Tawny Owl | PASS | PASS | PASS | PASS |
| GEAR: Wolf Medallion | PASS | PASS | PASS | PASS |
| Igni | PASS | PASS | PASS | PASS |
| Plot Twist | PASS | PASS | PASS | PASS |
| Rend | PASS | PASS | PASS | PASS |
| Riposte | PASS | PASS | PASS | PASS |
| Witcher Senses | PASS | PASS | PASS | PASS |
| Yrden | PASS | PASS | PASS | PASS |

---

### ciri

- archive integrity: **PASS**
- canonical manifest comparison: **PASS WITH QUALIFICATIONS**
- unique-card image completeness: **11/11**
- printed card identity/content verification: **all 11 images inspected; no image-to-Git discrepancy**
- discrepancies: none
- verdict: **PASS_WITH_QUALIFICATIONS** due integration requirements only

#### Archive / construction counts

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
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
| Source-tagged copies | 13 |
| Construction | fixed |

Visual inspection specifically confirmed that `Channel The Source` itself does **not** carry the Source card icon; Source symbols are referenced by its effect text. Therefore canonical Git correctly omits a `source` tag on that definition.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Bane Of The Aen Elle | PASS | PASS | PASS | PASS |
| Blink | PASS | PASS | PASS | PASS |
| Channel The Source | PASS | PASS | PASS | PASS |
| Child Of The Elder Blood | PASS | PASS | PASS | PASS |
| Lion Cub Of Cintra | PASS | PASS | PASS | PASS |
| Parry | PASS | PASS | PASS | PASS |
| Pushed To The Brink | PASS | PASS | PASS | PASS |
| Searching Strike | PASS | PASS | PASS | PASS |
| The Lady Of Space And Time | PASS | PASS | PASS | PASS |
| Unicorn Ally | PASS | PASS | PASS | PASS |
| Zireael | PASS | PASS | PASS | PASS |

Sensitive cases checked successfully include:

- `Bane Of The Aen Elle`: the movement result precedes and constrains the later choice of an opposing fighter actually moved through.
- Source-threshold cards: canonical modeling preserves the source-defined highest-met-threshold behavior.
- `Searching Strike`: search by Source property, reveal selected card, move it to hand; the post-search shuffle is a governing search-rule behavior rather than text that must be literally present on the card face.

#### Integration requirements confirmed

- `C-REQ-005` — filtered Source count / highest-met-threshold behavior.
- `C-REQ-006` — structured search/reveal/destination/post-search behavior for `Searching Strike`.
- `C-REQ-013` — staged dependent continuation, especially visible in `Bane Of The Aen Elle`.

## Findings

### TITANIA-001 — Gift Of The Fair Folk target cardinality

- fighter: `titania`
- card: `Gift Of The Fair Folk`
- severity: **P1**
- expected: the printed movement option requires choosing **two fighters** and moving each chosen fighter up to 2 spaces.
- observed: canonical Git models the `move_two` branch with a fighter selector using `max: 2`, which structurally permits selecting fewer than two fighters.
- evidence/reasoning: this changes the legal choice set. `up to 2 spaces` modifies movement distance for each selected fighter; it does not change the printed cardinality from exactly two fighters to up to two fighters.
- expected semantic correction: represent the branch as a choice of exactly two distinct legal fighters, then move each chosen fighter up to 2 spaces. Generic impossible-operation handling may still apply where global rules require it, but the normal source-defined choice must not become optional/fewer-target selection.

### TITANIA-002 — missing six auxiliary Glamour card images

- fighter: `titania`
- card: `Glamour of Greed`, `Glamour of Invisibility`, `Glamour of Jealousy`, `Glamour of Love`, `Glamour of Rhyme`, `Glamour of Sleep`
- severity: **P2 evidence issue**
- expected: one image for every canonical unique gameplay-card definition represented by the evidence corpus, including Titania's six first-class auxiliary Glamour definitions.
- observed: the archive contains only the 12 ordinary action-card definitions; all six Glamour faces are absent.
- evidence/reasoning: canonical Git explicitly defines the six Glamours as auxiliary gameplay cards instantiated into a separate Glamour deck/lifecycle. They are not ordinary action-deck copies and must not be added to the 30-card action deck, but their printed faces are still material gameplay evidence.
- expected semantic/evidence correction: provide one image for each of the six Glamour definitions while preserving their auxiliary-card status outside the ordinary 30-card action deck.

### HAMLET-001 — archive card ID normalization

- fighter: `hamlet`
- card: `The Play's The Thing`
- severity: **P3**
- expected: archive mapping uses canonical `card_id: the-plays-the-thing`.
- observed: archive mapping uses `the-play-s-the-thing`.
- evidence/reasoning: printed identity and gameplay mapping are unambiguous, so this is metadata/presentation rather than gameplay corruption.
- expected correction: normalize the archive mapping ID only; no gameplay-semantic change.

### GERALT-001 — archive card ID normalization

- fighter: `geralt-of-rivia`
- card: `Damn, You're Ugly`
- severity: **P3**
- expected: archive mapping uses canonical `card_id: damn-youre-ugly`.
- observed: archive mapping uses `damn-you-re-ugly`.
- evidence/reasoning: printed identity and gameplay mapping are unambiguous.
- expected correction: normalize the archive mapping ID only; no gameplay-semantic change.

## Corpus-Level Observations

1. **The supplied archive was technically healthy.** All five nested ZIPs opened successfully; all 62 supplied images decoded; no zero-byte, corrupted, truncated, or binary-duplicate image was found.

2. **One image per unique card definition was correctly used for ordinary action decks.** Physical copy quantity was stored in manifests rather than represented by duplicated image files.

3. **Deck size must remain data-driven.** Geralt demonstrates why a QA worker must not require archive quantity to equal 30: his canonical available pool is 36 physical cards while legal pre-game selection constructs a 30-card game deck.

4. **External gameplay definitions are not automatically card images.** Wayward Sisters spells are rules/gameplay definitions but not action-card instances and therefore were correctly excluded from ordinary card-image completeness.

5. **Auxiliary physical gameplay cards are different from non-card external definitions.** Titania's Glamours are actual auxiliary card definitions/instances. Omitting their images creates a real evidence-completeness gap even though her ordinary action deck is complete at 30 cards.

6. **Normalized Git semantics were compared behaviorally, not by literal prose equality.** Structured stages, bindings, choices, costs, MOVE/PLACE distinctions, LOOK_AT/REVEAL distinctions, threshold behavior, and source-defined construction were accepted when they preserved printed legal behavior.

7. **Project normalizations were not treated as publisher rulings.** Hamlet's failed-placement interpretation for `The Readiness Is All` was kept as an explicit qualification because the supplied image does not settle the failed-PLACE edge case and does not contradict the documented normalization.

8. **Integration requirements are not transcription failures by themselves.** Ciri's `C-REQ-005`, `C-REQ-006`, `C-REQ-013`, Hamlet's `C-REQ-013`, and Titania's confirmed auxiliary/dependency requirements remain implementation-contract qualifications rather than evidence mismatches where printed behavior and normalized Git behavior agree.

## Final Assessment

This batch is **not fully acceptable as a clean Phase 4 card-image corpus in its current state**.

Four fighters are materially clean with respect to supplied physical card faces:

- `hamlet` — no P1/P2 card discrepancy; one P3 archive-ID issue plus an already documented edge-case normalization qualification.
- `wayward-sisters` — clean PASS.
- `geralt-of-rivia` — no P1/P2 discrepancy; one P3 archive-ID issue.
- `ciri` — no image-to-Git discrepancy; integration requirements only.

`titania` prevents a batch PASS for two independent reasons:

1. `Gift Of The Fair Folk` has a **P1 gameplay-changing target-cardinality mismatch** between printed text and canonical normalized representation.
2. The evidence archive omits all **six canonical auxiliary Glamour card images**, creating a P2 evidence-completeness gap.

No repository correction was performed by this QA worker. The findings above localize the required semantic/evidence corrections while preserving Git as the canonical documentation state under review.
