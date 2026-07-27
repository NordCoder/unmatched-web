# Phase 4 Card-Image QA Report

## Batch

- black-panther
- black-widow
- ms-marvel
- squirrel-girl
- winter-soldier

Evidence archive: `unmatched-bundle-08.zip`

QA source state used during the completed verification:

- `main`: `f40c3f9970da24ab4e17ff51fc75fb5d7080b239`
- `phase-4b-worker-b-licensed`: `f383a639e980d9e753c742b20b7e085d3163502f`

Persistence target was refreshed before this report was written; `main` tip immediately before the write was `43854bdbb4a8b021a671209a5637bedc7ddb60ba`.

## Verdict

**FAIL**

## Summary

The batch contained exactly five fighter ZIP archives and all five were fully inspected. Every nested ZIP opened successfully, every `manifest.yaml` was readable, and all 64 unique card images decoded and were visually inspected. There were no zero-byte/corrupted images, no binary duplicate images, no missing image files, no extra image files, and no quantity/deck-construction failures.

The physical card corpus nevertheless exposed material canonical-manifest discrepancies in four fighters:

- `black-panther`: 6 gameplay-changing card-type transcription errors.
- `ms-marvel`: 4 gameplay-changing card-type transcription errors plus one archive-to-canonical card-ID mapping mismatch.
- `black-widow`: 4 gameplay-changing card-type transcription errors, one forced-discard cardinality error, one narrower relocation-lock semantic gap, and archive-to-canonical card-ID mapping mismatches.
- `winter-soldier`: 1 exact-cardinality semantic error on `Reprogram`.
- `squirrel-girl`: no card-corpus discrepancy; shared integration/evidence qualifications remain.

Batch totals:

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| PASS | 0 |
| PASS_WITH_QUALIFICATIONS | 1 |
| FAIL | 4 |
| BLOCKED | 0 |
| Total unique card images | 64 |
| Images successfully inspected | 64 |
| Unreadable images | 0 |
| Missing images | 0 |
| Duplicate binary images | 0 |
| Extra images | 0 |
| Quantity failures | 0 |
| P1 findings | 16 |
| P2 findings | 1 |
| P3 card-level mapping mismatches | 5 |

## Fighter Results

### black-panther

- branch: `main`
- deck manifest: `docs/cards/phase-4a/black-panther.yaml`
- fighter manifest: `docs/fighters/phase-4a/black-panther.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **PASS — 13/13**
- printed card identity/content verification: **FAIL — six card-type mismatches; printed effect semantics otherwise matched**
- discrepancies: 6 P1
- verdict: **FAIL**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13 / 13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 13 |
| Construction | fixed |

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Analyze and Adjust | PASS | FAIL — type | PASS | FAIL |
| Ancestral Insight | PASS | FAIL — type | PASS | FAIL |
| Anti-Metal Claws | PASS | FAIL — type | PASS | FAIL |
| Cat-Like Reflexes | PASS | PASS | PASS | PASS |
| Evade | PASS | FAIL — type | PASS | FAIL |
| Feint | PASS | PASS | PASS | PASS |
| Nanotriage Processor | PASS | FAIL — type | PASS | FAIL |
| Regroup | PASS | PASS | PASS | PASS |
| Stalking Panther | PASS | PASS | PASS | PASS |
| Tactical Remote Scanning | PASS | PASS | PASS | PASS |
| Wakanda Forever! | PASS | PASS | PASS | PASS |
| Vibranium Shockwave | PASS | FAIL — type | PASS | FAIL |
| Microweave Mesh | PASS | PASS | PASS | PASS |

Integration requirements confirmed from card images: **None**.

### black-widow

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/black-widow.yaml`
- fighter manifest: `docs/fighters/phase-4b/black-widow.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **PASS — 15/15**
- printed card identity/content verification: **FAIL**
- discrepancies: 5 P1, 1 P2, 4 P3 card-ID mappings
- verdict: **FAIL**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 15 / 15 |
| Unique images | 15 |
| Manifest card entries | 15 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 31 |
| Git available pool | 31 |
| Git game deck | 31 |
| Unique Git definitions | 15 |
| Construction | fixed |

The 31-card deck is intentional; it was not treated as a quantity failure.

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Widow's Kiss | PASS | FAIL — type + archive card_id | FAIL | FAIL |
| Widow's Sting | PASS | FAIL — archive card_id | PASS | FAIL |
| Widow's Bite | PASS | FAIL — archive card_id | PASS | FAIL |
| Acting Director of S.H.I.E.L.D. | PASS | PASS | PASS | PASS |
| Caught in a Web | PASS | PASS | PASS | PASS |
| Double Identity | PASS | FAIL — type | PASS | FAIL |
| Fake Out | PASS | FAIL — type | PASS | FAIL |
| Widow's Line | PASS | FAIL — archive card_id | PASS | FAIL |
| Feint | PASS | FAIL — type | PASS | FAIL |
| Life Model Decoy | PASS | PASS | PASS | PASS |
| The Budapest Gambit | PASS | PASS | PASS | PASS |
| The Firenze Agenda | PASS | PASS | PASS | PASS |
| The Kinshasa Directive | PASS | PASS | FAIL | FAIL |
| The Madripoor Sanction | PASS | PASS | PASS | PASS |
| The Moscow Protocol | PASS | PASS | PASS | PASS |

The Moscow Protocol image also confirms the special setup instruction to start the game with that card in hand; the fighter manifest's setup hook is consistent with the printed component.

Integration requirements confirmed from card images: **None**.

### ms-marvel

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/ms-marvel.yaml`
- fighter manifest: `docs/fighters/phase-4b/ms-marvel.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **PASS — 11/11**
- printed card identity/content verification: **FAIL — four type mismatches; staged effect semantics otherwise matched**
- discrepancies: 4 P1, 1 P3 card-ID mapping
- verdict: **FAIL**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 11 / 11 |
| Unique images | 11 |
| Manifest card entries | 11 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 11 |
| Construction | fixed |

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Big Wind Up | PASS | PASS | PASS | PASS |
| Easy Peasy | PASS | PASS | PASS | PASS |
| Embiggen | PASS | FAIL — type | PASS | FAIL |
| Fangirl | PASS | PASS | PASS | PASS |
| Feint | PASS | FAIL — type | PASS | FAIL |
| Friends and Family | PASS | PASS | PASS | PASS |
| Gyro and Fries | PASS | PASS | PASS | PASS |
| I'm Not Touching You | PASS | FAIL — type + archive card_id | PASS | FAIL |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Shrink! Shrink! Shrink! | PASS | PASS | PASS | PASS |
| Slingshot | PASS | FAIL — type | PASS | FAIL |

Printed/staged semantics that were specifically checked and found consistent include:

- `Easy Peasy`: draw resolves before the hand-size gate for adjacent damage.
- `Shrink! Shrink! Shrink!`: correct distinction between the no-shared-zone branch that performs both effects and the shared-zone branch that chooses one.
- `Momentous Shift`: dependency on position relative to turn-start state.
- `Friends and Family`: additional-action expenditure is a dependency/cost before drawing up to seven.

Integration requirements confirmed from card images:

- `B-REQ-009`: directly confirmed by `Friends and Family`, which explicitly permits spending an additional action to draw until seven.
- `B-REQ-008`: not independently visually confirmed because it is fighter-ability/range behavior and the evidence ZIP is a card corpus, not a fighter-ability component corpus.

### squirrel-girl

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/squirrel-girl.yaml`
- fighter manifest: `docs/fighters/phase-4b/squirrel-girl.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS for accessible card corpus**
- unique-card image completeness: **PASS — 13/13**
- printed card identity/content verification: **PASS**
- discrepancies: none in card corpus
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13 / 13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique Git definitions | 13 |
| Construction | fixed |

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Bite of Steel | PASS | PASS | PASS | PASS |
| Call of the Mild | PASS | PASS | PASS | PASS |
| Dash | PASS | PASS | PASS | PASS |
| Eat Nuts | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Fuzzball Special | PASS | PASS | PASS | PASS |
| Get 'Em Tippy-Toe! | PASS | PASS | PASS | PASS |
| Horde of Squirrels | PASS | PASS | PASS | PASS |
| Kick Butts | PASS | PASS | PASS | PASS |
| Nutwork of Spies | PASS | PASS | PASS | PASS |
| Squirgility | PASS | PASS | PASS | PASS |
| Squirmish | PASS | PASS | PASS | PASS |
| Unbeatable Squirrel Girl | PASS | PASS | PASS | PASS |

The physical cards confirm the relevant card-side use of multiple Squirrel fighter instances: squirrel-count predicates, summoning, movement, placement, and special movement behavior.

Integration requirements confirmed from card images:

- `B-REQ-002`: confirmed by several card effects that need runtime support for multiple Squirrel fighters, summoning, movement, placement, occupancy and count predicates.

Qualification: the fighter manifest also contains project normalizations for `Go Nuts` destination occupancy and propagated Squirrel damage counting. The supplied evidence archive contains action cards rather than the fighter ability component, so those fighter-level printed facts were not independently visually verified by this batch. This is an evidence qualification, not a card transcription failure.

### winter-soldier

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/winter-soldier.yaml`
- fighter manifest: `docs/fighters/phase-4b/winter-soldier.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **PASS — 12/12**
- printed card identity/content verification: **FAIL — one exact-cardinality discrepancy**
- discrepancies: 1 P1
- verdict: **FAIL**

#### Archive / construction evidence

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12 / 12 |
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

#### Per-card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| A Boy Named Bucky | PASS | PASS | PASS | PASS |
| Bionic Arm | PASS | PASS | PASS | PASS |
| Born in the Barracks | PASS | PASS | PASS | PASS |
| Complete the Mission | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Manipulation | PASS | PASS | PASS | PASS |
| Marksman | PASS | PASS | PASS | PASS |
| Programmed to Kill | PASS | PASS | PASS | PASS |
| Reflex Memories | PASS | PASS | PASS | PASS |
| Reprogram | PASS | PASS | FAIL | FAIL |
| Wily Fighting | PASS | PASS | PASS | PASS |
| Without Remorse | PASS | PASS | PASS | PASS |

All BRAINWASHED/star effects were visually readable and otherwise matched the normalized manifest, including opponent-controlled optional draws, self-damage, random discard followed by draw, and `A Boy Named Bucky` suppressing BRAINWASHED/star clauses for the rest of the turn without suppressing ordinary card effects.

Integration requirements confirmed from card images: **None**.

## Findings

### B08-BP-001 — Analyze and Adjust card type

- fighter: `black-panther`
- card: `Analyze and Adjust`
- Git location: `docs/cards/phase-4a/black-panther.yaml`, definition `analyze-and-adjust`
- expected / image fact: red ATTACK icon, combat value 3, BOOST 3
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: combat-card type determines whether the card can be committed as attack, defense, or either; this is gameplay-changing.
- expected correction: represent the card as `attack`.

### B08-BP-002 — Ancestral Insight card type

- fighter: `black-panther`
- card: `Ancestral Insight`
- expected / image fact: purple VERSATILE icon, combat value 4, BOOST 1
- observed in Git: `type: attack`
- severity: **P1**
- evidence/reasoning: printed combat use differs from the canonical metadata.
- expected correction: represent the card as `versatile`.

### B08-BP-003 — Anti-Metal Claws card type

- fighter: `black-panther`
- card: `Anti-Metal Claws`
- expected / image fact: purple VERSATILE icon, combat value 1, BOOST 2
- observed in Git: `type: attack`
- severity: **P1**
- evidence/reasoning: printed combat use differs from the canonical metadata.
- expected correction: represent the card as `versatile`.

### B08-BP-004 — Evade card type

- fighter: `black-panther`
- card: `Evade`
- expected / image fact: blue DEFENSE icon, combat value 3, BOOST 2
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: a defense-only card is currently represented as usable on both sides of combat.
- expected correction: represent the card as `defense`.

### B08-BP-005 — Nanotriage Processor card type

- fighter: `black-panther`
- card: `Nanotriage Processor`
- expected / image fact: purple VERSATILE icon, combat value 2, BOOST 2
- observed in Git: `type: defense`
- severity: **P1**
- evidence/reasoning: printed combat use differs from canonical metadata.
- expected correction: represent the card as `versatile`.

### B08-BP-006 — Vibranium Shockwave card type

- fighter: `black-panther`
- card: `Vibranium Shockwave`
- expected / image fact: red ATTACK icon, combat value 2, BOOST 2
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: card legality in combat is broadened incorrectly.
- expected correction: represent the card as `attack`.

### B08-MM-001 — Embiggen card type

- fighter: `ms-marvel`
- card: `Embiggen`
- expected / image fact: purple VERSATILE icon, combat value 3, BOOST 3
- observed in Git: `type: attack`
- severity: **P1**
- evidence/reasoning: printed combat use differs from canonical metadata.
- expected correction: represent the card as `versatile`.

### B08-MM-002 — Feint card type

- fighter: `ms-marvel`
- card: `Feint`
- expected / image fact: purple VERSATILE icon, combat value 2, BOOST 1
- observed in Git: `type: defense`
- severity: **P1**
- evidence/reasoning: canonical metadata incorrectly removes attack-side use.
- expected correction: represent the card as `versatile`.

### B08-MM-003 — I'm Not Touching You card type

- fighter: `ms-marvel`
- card: `I'm Not Touching You`
- expected / image fact: red ATTACK icon, combat value 4, BOOST 1
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: canonical metadata incorrectly permits defense-side use.
- expected correction: represent the card as `attack`.

### B08-MM-004 — Slingshot card type

- fighter: `ms-marvel`
- card: `Slingshot`
- expected / image fact: blue DEFENSE icon, combat value 3, BOOST 2
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: canonical metadata incorrectly permits attack-side use.
- expected correction: represent the card as `defense`.

### B08-MM-005 — I'm Not Touching You archive mapping

- fighter: `ms-marvel`
- card: `I'm Not Touching You`
- expected: archive evidence mapping targets canonical ID `im-not-touching-you`
- observed: archive manifest uses `i-m-not-touching-you`
- severity: **P3**
- evidence/reasoning: the physical card name is unambiguous and Git has the intended definition; mismatch is in evidence mapping rather than gameplay semantics.
- expected correction: map the evidence image to canonical `im-not-touching-you`.

### B08-BW-001 — Double Identity card type

- fighter: `black-widow`
- card: `Double Identity`
- expected / image fact: blue DEFENSE icon, combat value 3, ANY restriction, BOOST 2
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: canonical metadata incorrectly permits attack-side use.
- expected correction: represent the card as `defense`.

### B08-BW-002 — Fake Out card type

- fighter: `black-widow`
- card: `Fake Out`
- expected / image fact: red ATTACK icon, combat value 1, ANY restriction, BOOST 1
- observed in Git: `type: versatile`
- severity: **P1**
- evidence/reasoning: canonical metadata incorrectly permits defense-side use.
- expected correction: represent the card as `attack`.

### B08-BW-003 — Feint card type

- fighter: `black-widow`
- card: `Feint`
- expected / image fact: purple VERSATILE icon, combat value 2, ANY restriction, BOOST 1
- observed in Git: `type: defense`
- severity: **P1**
- evidence/reasoning: canonical metadata incorrectly removes attack-side use.
- expected correction: represent the card as `versatile`.

### B08-BW-004 — Widow's Kiss card type

- fighter: `black-widow`
- card: `Widow's Kiss`
- expected / image fact: purple VERSATILE icon, combat value 4, Black Widow restriction, BOOST 2
- observed in Git: `type: attack`
- severity: **P1**
- evidence/reasoning: printed combat use differs from canonical metadata.
- expected correction: represent the card as `versatile`.

### B08-BW-005 — The Kinshasa Directive discard cardinality

- fighter: `black-widow`
- card: `The Kinshasa Directive`
- expected / image fact: on mission success, the chosen opponent discards **2 cards**; if fewer cards are actually available, ordinary partial-resolution constraints may limit what can happen, but the opponent is not voluntarily choosing 0–2.
- observed in Git: card-set choice uses `count_rule: up_to_2_limited_by_hand_size`, `min: 0`, `max: 2`.
- severity: **P1**
- evidence/reasoning: with at least two cards in hand, canonical semantics permit choosing zero or one even though the printed card requires two.
- expected correction: require exactly `min(2, available hand size)` rather than a voluntary `0..2` selection.

### B08-BW-006 — Widow's Kiss relocation lock scope

- fighter: `black-widow`
- card: `Widow's Kiss`
- expected / image fact: the opposing fighter “may not leave their space for the rest of the turn.”
- observed in Git: prevention is modeled specifically for `MOVE` and `PLACE`.
- severity: **P2**
- evidence/reasoning: the project's effect model defines `SWAP` as a separate relocation primitive. A general leave-space prohibition is semantically broader than only MOVE/PLACE and can create an edge case if another relocation primitive moves the affected fighter.
- expected correction: represent a general cannot-leave-current-space restriction for the duration, or explicitly cover all relocation operations that can cause the fighter to leave that space.

### B08-BW-007 — Widow apostrophe archive mappings

- fighter: `black-widow`
- cards: `Widow's Bite`, `Widow's Kiss`, `Widow's Line`, `Widow's Sting`
- expected canonical IDs: `widows-bite`, `widows-kiss`, `widows-line`, `widows-sting`
- observed archive IDs: `widow-s-bite`, `widow-s-kiss`, `widow-s-line`, `widow-s-sting`
- severity: **P3** per affected card (4 card-level mapping mismatches)
- evidence/reasoning: physical names are unambiguous and the canonical definitions are identifiable; mismatch is evidence mapping, not gameplay content.
- expected correction: map each evidence entry to the corresponding canonical ID.

### B08-WS-001 — Reprogram exact-cardinality semantics

- fighter: `winter-soldier`
- card: `Reprogram`
- expected / image fact: “Choose 3 cards in your discard pile and shuffle them into your deck.”
- observed in Git: the card-set choice permits `min: 0`, `max: 3`.
- severity: **P1**
- evidence/reasoning: when at least three cards are available, the printed effect requires choosing three; canonical semantics allow voluntarily choosing fewer or none.
- expected correction: select exactly three when three are available, with ordinary partial resolution only when fewer than three are actually available.

## Corpus-Level Observations

1. **Archive integrity is strong.** All five fighter ZIPs were technically sound. No corruption, truncation, zero-byte images, missing evidence images or binary duplicate card images were found.
2. **The one-image-per-unique-definition convention is satisfied.** Copies are represented through manifest quantity rather than duplicate images.
3. **All deck quantity checks passed.** This includes Black Widow's non-standard 31-card deck; no universal 30-card assumption was applied.
4. **Printed image evidence materially improved QA beyond source-index checks.** Fourteen card-type transcription errors were visible directly from the printed combat-type icons despite otherwise plausible names, values and effects.
5. **Normalized staged effect structures were not penalized for differing from printed prose.** They were compared by legal gameplay behavior, including timing, ordering, dependency, optionality, targeting, movement versus placement, and hidden-information semantics.
6. **Exact cardinality is a recurring risk.** `The Kinshasa Directive` and `Reprogram` demonstrate that translating “discard/choose N” into an `up_to_N` domain changes legal player choice even if generic partial resolution could handle insufficient available cards.
7. **Generic relocation locks need operation-taxonomy awareness.** Printed “may not leave their space” language should not accidentally become only MOVE/PLACE prevention if the engine also models SWAP or another relocation primitive.
8. **Archive manifests are useful mappings, not authoritative card truth.** Five evidence `card_id` slugs disagreed with canonical IDs while the physical card identity itself remained unambiguous.
9. **Shared requirements were treated separately from transcription failures.** `B-REQ-002` is directly supported by Squirrel Girl card evidence, and `B-REQ-009` is directly supported by Ms. Marvel's `Friends and Family`. These do not themselves make the cards transcription failures.
10. **Fighter-ability evidence boundaries were preserved.** The Squirrel Girl `Go Nuts` normalization and Ms. Marvel range ability were not claimed as visually confirmed by a card-only archive.

## Final Assessment

This batch is **not ready to be considered a fully reconciled Phase 4 card-image corpus** in its current canonical documentation state.

The evidence archive itself is complete and technically usable: all 64 unique physical card images are present, readable and correctly organized at the one-image-per-unique-card-definition level, and all five deck quantity/construction checks pass.

However, four of the five fighters contain material canonical discrepancies. The most widespread issue is combat-card type transcription: 14 cards have ATTACK/DEFENSE/VERSATILE metadata that conflicts with the printed component. In addition, Black Widow's `The Kinshasa Directive` and Winter Soldier's `Reprogram` lose mandatory cardinality, and Widow's Kiss has a narrower operation-level relocation representation than its printed “may not leave their space” rule.

`Squirrel Girl` is the only fighter in this batch without a material card-corpus discrepancy and is therefore clean at the card-image level, subject to the explicitly documented integration/evidence qualifications.

No canonical manifests or other repository files were modified as part of the QA itself; this file only persists the independent QA findings.