# Phase 4 Card-Image QA Report

## Batch

- bullseye
- elektra
- ghost-rider
- luke-cage
- moon-knight

## Verdict

FAIL

## Summary

Independent read-only card-corpus QA was completed against the current tip of `phase-4b-worker-b-licensed` before report persistence (`f383a639e980d9e753c742b20b7e085d3163502f`). The external batch contained exactly five fighter ZIPs. All nested ZIPs passed integrity/CRC checks. All 59 manifest entries had exactly one corresponding image, all 59 images decoded and were visually inspected, no zero-byte files, missing images, extra images, or binary duplicate images were found.

All five archive quantity sums reconcile with canonical Git construction: Elektra 20/20; Bullseye, Ghost Rider, Luke Cage, and Moon Knight 30/30. No quantity failure was found.

The corpus is not clean: material canonical Git discrepancies were found in every fighter. The principal issues are nine P1 card-type mismatches, two P1 optional/mandatory semantic mismatches, and one P2 defeat-condition mismatch. Several archive card IDs also differ from canonical Git IDs only by apostrophe/slug normalization; those are P3 mapping issues rather than gameplay errors.

## Fighter Results

### elektra

- archive integrity: PASS — nested ZIP readable; manifest readable; 10/10 images readable; 10 unique images; 10 manifest entries; 0 duplicates; 0 missing; 0 extra.
- canonical manifest comparison: quantity PASS; fixed construction; archive sum 20; Git available pool 20; Git game deck 20; 10 archive definitions and 10 Git definitions.
- unique-card image completeness: PASS — 10/10.
- printed card identity/content verification: 10/10 images inspected. Mystic Assassin, The Fist, Intercept, Snakeroot Clan, Sai, Ninjitsu, Whirlwind, Mesmerize, and Cloaked in Shadow match canonical metadata and normalized semantics. Hands of Red has a semantic mismatch.
- discrepancies: `B06-ELE-001` (P1), Hands of Red mandatory-vs-optional return behavior.
- integration qualification: fighter manifest requires `B-REQ-003` for Elektra's resurrection/off-board lifecycle, but the supplied evidence contains card faces rather than the fighter ability component; this requirement is therefore not claimed as visually confirmed by this batch.
- verdict: FAIL.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Mystic Assassin | PASS | PASS | PASS | PASS |
| The Fist | PASS | PASS | PASS | PASS |
| Hands of Red | PASS | PASS | FAIL P1 | FAIL |
| Intercept | PASS | PASS | PASS | PASS |
| Snakeroot Clan | PASS | PASS | PASS | PASS |
| Sai | PASS | PASS | PASS | PASS |
| Ninjitsu | PASS | PASS | PASS | PASS |
| Whirlwind | PASS | PASS | PASS | PASS |
| Mesmerize | PASS | PASS | PASS | PASS |
| Cloaked in Shadow | PASS | PASS | PASS | PASS |

### bullseye

- archive integrity: PASS — nested ZIP readable; manifest readable; 12/12 images readable; 12 unique images; 12 manifest entries; 0 duplicates; 0 missing; 0 extra.
- canonical manifest comparison: quantity PASS; fixed construction; archive sum 30; Git available pool 30; Git game deck 30; 12 archive definitions and 12 Git definitions.
- unique-card image completeness: PASS — 12/12.
- printed card identity/content verification: 12/12 images inspected. Eleven cards have matching gameplay semantics. Ricochet has a narrower Git defeat predicate than the printed condition. Two archive IDs use non-canonical apostrophe slugging.
- discrepancies: `B06-BUL-001` (P2), Ricochet; `B06-BUL-P3-001` (P3), archive card-ID normalization.
- integration qualification: fighter manifest requires `B-REQ-008` for graph-distance attack legality; this is fighter ability behavior and is not directly printed on the supplied action cards.
- verdict: FAIL.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| World's Greatest Assassin | PASS | QUALIFIED P3 | PASS | QUALIFIED |
| I Never Miss | PASS | PASS | PASS | PASS |
| For My Next Trick | PASS | PASS | PASS | PASS |
| I Planned To Be Here | PASS | PASS | PASS | PASS |
| Ricochet | PASS | PASS | FAIL P2 | FAIL |
| Master Strategist | PASS | PASS | PASS | PASS |
| Right Between The Eyes | PASS | PASS | PASS | PASS |
| I'm Better And I'll Prove It | PASS | QUALIFIED P3 | PASS | QUALIFIED |
| Arrogant But Effective | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Tactical Retreat | PASS | PASS | PASS | PASS |
| Study The Target | PASS | PASS | PASS | PASS |

### ghost-rider

- archive integrity: PASS — nested ZIP readable; manifest readable; 12/12 images readable; 12 unique images; 12 manifest entries; 0 duplicates; 0 missing; 0 extra.
- canonical manifest comparison: quantity PASS; fixed construction; archive sum 30; Git available pool 30; Git game deck 30; 12 archive definitions and 12 Git definitions.
- unique-card image completeness: PASS — 12/12.
- printed card identity/content verification: 12/12 images inspected. Printed values, BOOST values, restrictions, timing, Hellfire costs/dependencies, movement, damage, action gain, and other effect semantics match. Six canonical Git card types do not match the physical cards.
- discrepancies: six P1 card-type mismatches (`B06-GR-001` through `B06-GR-006`).
- verdict: FAIL.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Spirit of Vengeance | PASS | PASS | PASS | PASS |
| Control the Demon | PASS | FAIL P1 | PASS | FAIL |
| Penance Stare | PASS | FAIL P1 | PASS | FAIL |
| I Finally Escaped Hell | PASS | PASS | PASS | PASS |
| I Brought the Devil With Me | PASS | FAIL P1 | PASS | FAIL |
| Blaze of Glory | PASS | FAIL P1 | PASS | FAIL |
| The Wicked Will Burn | PASS | PASS | PASS | PASS |
| Deal With the Devil | PASS | FAIL P1 | PASS | FAIL |
| Stoke the Flames | PASS | PASS | PASS | PASS |
| Feint | PASS | FAIL P1 | PASS | FAIL |
| Hell Rides With Me | PASS | PASS | PASS | PASS |
| Chains of Hellfire | PASS | PASS | PASS | PASS |

### luke-cage

- archive integrity: PASS — nested ZIP readable; manifest readable; 13/13 images readable; 13 unique images; 13 manifest entries; 0 duplicates; 0 missing; 0 extra.
- canonical manifest comparison: quantity PASS; fixed construction; archive sum 30; Git available pool 30; Git game deck 30; 13 archive definitions and 13 Git definitions.
- unique-card image completeness: PASS — 13/13.
- printed card identity/content verification: 13/13 images inspected. Restrictions (ANY/Luke/Misty), printed values, BOOSTs, forced and optional movement, defender replacement, `instead`, end-turn, placement, damage and recycle semantics match. Three canonical Git card types differ from physical cards. One archive ID uses non-canonical apostrophe slugging.
- discrepancies: three P1 card-type mismatches; one P3 archive mapping mismatch.
- verdict: FAIL.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Sweet Christmas! | PASS | PASS | PASS | PASS |
| Commanding Impact | PASS | PASS | PASS | PASS |
| Get Paid | PASS | PASS | PASS | PASS |
| Hero For Hire | PASS | PASS | PASS | PASS |
| Power Man | PASS | FAIL P1 | PASS | FAIL |
| Still Standing | PASS | FAIL P1 | PASS | FAIL |
| Pushback | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| Skin Like Titanium | PASS | PASS | PASS | PASS |
| Got My Back? | PASS | PASS | PASS | PASS |
| Trash Talk | PASS | PASS | PASS | PASS |
| Daughter of the Dragon | PASS | FAIL P1 | PASS | FAIL |
| Where's My Money? | PASS | QUALIFIED P3 | PASS | QUALIFIED |

### moon-knight

- archive integrity: PASS — nested ZIP readable; manifest readable; 12/12 images readable; 12 unique images; 12 manifest entries; 0 duplicates; 0 missing; 0 extra.
- canonical manifest comparison: quantity PASS; fixed construction; archive sum 30; Git available pool 30; Git game deck 30; 12 archive definitions and 12 Git definitions.
- unique-card image completeness: PASS — 12/12.
- printed card identity/content verification: 12/12 images inspected. Eleven cards preserve printed gameplay semantics. That's the Part I Like loses the printed optionality around the complete hidden-information effect. Four archive IDs use non-canonical apostrophe slugging.
- discrepancies: `B06-MK-001` (P1), optionality/hidden-information behavior; P3 archive mapping normalization.
- verdict: FAIL.

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| I'm Not Real | PASS | QUALIFIED P3 | PASS | QUALIFIED |
| That's Why I Always Win | PASS | QUALIFIED P3 | PASS | QUALIFIED |
| Good Enough For Us | PASS | PASS | PASS | PASS |
| That's the Part I Like | PASS | QUALIFIED P3 | FAIL P1 | FAIL |
| Fist of Khonshu | PASS | PASS | PASS | PASS |
| Past and Present Intermingle | PASS | PASS | PASS | PASS |
| A Totally Sane Thing to Do | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Let Your Insanity Guide You | PASS | PASS | PASS | PASS |
| We're All in This Together | PASS | QUALIFIED P3 | PASS | QUALIFIED |
| Madness Will Keep You Alive | PASS | PASS | PASS | PASS |
| Travelers of the Night | PASS | PASS | PASS | PASS |

## Findings

### B06-ELE-001 — Hands of Red mandatory return modeled as optional

- fighter: `elektra`
- card: `hands-of-red`
- severity: P1
- expected: printed effect requires returning a defeated Hand to a space in Elektra's zone when a legal defeated Hand exists; the physical card does not say `may`.
- observed: canonical Git models the `hand_fighter` choice as `optional: true`, allowing the controller to skip a legal return.
- evidence/reasoning: optional versus mandatory behavior changes legal gameplay. The semantic representation should make target selection/return mandatory when a legal defeated Hand exists, while naturally allowing no operation if no legal target exists.

### B06-BUL-001 — Ricochet defeat condition narrowed

- fighter: `bullseye`
- card: `ricochet`
- severity: P2
- expected: printed condition tests whether `the opposing fighter was not defeated`.
- observed: canonical Git uses `opposing_combat_fighter_defeated_by_this_attack: false`.
- evidence/reasoning: the Git predicate is narrower because it encodes defeat specifically by the attack. An edge case where the opposing fighter becomes defeated during the combat through another cause can therefore produce behavior different from the printed condition. The normalized predicate should test the opposing fighter's defeat state/outcome at the relevant timing without imposing an unsupported cause restriction.

### B06-GR-001 — Control the Demon type

- fighter: `ghost-rider`
- card: `control-the-demon`
- severity: P1
- expected: VERSATILE 0, BOOST 1.
- observed: canonical Git type is `attack`; value and BOOST are correct.
- evidence/reasoning: card type controls attack/defense legality. Expected semantic correction: `versatile`.

### B06-GR-002 — Penance Stare type

- fighter: `ghost-rider`
- card: `penance-stare`
- severity: P1
- expected: VERSATILE 3, BOOST 2.
- observed: canonical Git type is `attack`.
- evidence/reasoning: physical card can be used as either attack or defense. Expected semantic correction: `versatile`.

### B06-GR-003 — I Brought the Devil With Me type

- fighter: `ghost-rider`
- card: `i-brought-the-devil-with-me`
- severity: P1
- expected: ATTACK 3, BOOST 2.
- observed: canonical Git type is `versatile`.
- evidence/reasoning: Git incorrectly permits defense use. Expected semantic correction: `attack`.

### B06-GR-004 — Blaze of Glory type

- fighter: `ghost-rider`
- card: `blaze-of-glory`
- severity: P1
- expected: ATTACK 2, BOOST 3.
- observed: canonical Git type is `versatile`.
- evidence/reasoning: Git incorrectly permits defense use. Expected semantic correction: `attack`.

### B06-GR-005 — Deal With the Devil type

- fighter: `ghost-rider`
- card: `deal-with-the-devil`
- severity: P1
- expected: DEFENSE 2, BOOST 1.
- observed: canonical Git type is `versatile`.
- evidence/reasoning: Git incorrectly permits attack use. Expected semantic correction: `defense`.

### B06-GR-006 — Feint type

- fighter: `ghost-rider`
- card: `feint`
- severity: P1
- expected: VERSATILE 2, BOOST 1.
- observed: canonical Git type is `defense`.
- evidence/reasoning: Git incorrectly prevents attack use. Expected semantic correction: `versatile`.

### B06-LC-001 — Power Man type

- fighter: `luke-cage`
- card: `power-man`
- severity: P1
- expected: DEFENSE 2, BOOST 1.
- observed: canonical Git type is `attack`.
- evidence/reasoning: legal combat role is reversed. Expected semantic correction: `defense`.

### B06-LC-002 — Still Standing type

- fighter: `luke-cage`
- card: `still-standing`
- severity: P1
- expected: ATTACK 4, BOOST 2.
- observed: canonical Git type is `versatile`.
- evidence/reasoning: Git incorrectly permits defense use. Expected semantic correction: `attack`.

### B06-LC-003 — Daughter of the Dragon type

- fighter: `luke-cage`
- card: `daughter-of-the-dragon`
- severity: P1
- expected: VERSATILE 2, BOOST 2.
- observed: canonical Git type is `defense`.
- evidence/reasoning: Git incorrectly prevents attack use. Expected semantic correction: `versatile`.

### B06-MK-001 — That's the Part I Like loses optionality

- fighter: `moon-knight`
- card: `thats-the-part-i-like`
- severity: P1
- expected: after winning combat, the controller **may** begin the effect. If used, look at the top three cards of the opponent's deck, discard one, and return the other two in any order.
- observed: canonical Git executes `LOOK_AT` whenever the combat is won, with no preceding optional choice, and then requires the discard/reorder sequence.
- evidence/reasoning: this changes both player choice and the hidden-information boundary. The entire inspect/discard/reorder package should be gated behind a public optional decision; once chosen, the printed follow-up operations are mandatory.

### P3 archive card-ID normalization findings

These do not alter printed gameplay and do not invalidate the corresponding physical image identity, but the archive manifest does not use the canonical Git card ID exactly.

- `bullseye`: `world-s-greatest-assassin` -> `worlds-greatest-assassin`
- `bullseye`: `i-m-better-and-i-ll-prove-it` -> `im-better-and-ill-prove-it`
- `luke-cage`: `where-s-my-money` -> `wheres-my-money`
- `moon-knight`: `i-m-not-real` -> `im-not-real`
- `moon-knight`: `that-s-why-i-always-win` -> `thats-why-i-always-win`
- `moon-knight`: `that-s-the-part-i-like` -> `thats-the-part-i-like`
- `moon-knight`: `we-re-all-in-this-together` -> `were-all-in-this-together`

Expected correction is either canonical card IDs in the evidence manifest or an explicit card-level alias mapping. No new fighter/card identity should be inferred from punctuation-only slug differences.

## Corpus-Level Observations

- Fighters received: 5.
- Fighters fully checked: 5.
- PASS: 0.
- PASS WITH QUALIFICATIONS: 0.
- FAIL: 5.
- BLOCKED: 0.
- Total unique card images: 59.
- Images successfully inspected: 59.
- Unreadable images: 0.
- Missing images: 0.
- Extra images: 0.
- Binary duplicate images: 0.
- Quantity failures: 0.
- Canonical Git card-type mismatches: 9 P1 findings.
- Other semantic discrepancies: 2 P1 findings and 1 P2 finding.
- Archive card-ID normalization mismatches: 7 P3 findings.
- Clean fighters without P1/P2 discrepancies: none.
- Real blockers: none.

The archive manifest was used only for image/card/quantity mapping. Printed facts and printed semantics were established by visual inspection of every physical card image and compared to the normalized Git representation. Literal prose equality was not required. Structured stages, choices, bindings, costs, dependencies, and normalized operations were accepted whenever they preserved the legal printed behavior.

Non-card fighter abilities were not claimed as visually confirmed merely because the fighter manifest documents them. In particular, Elektra's `B-REQ-003` and Bullseye's `B-REQ-008` remain documented integration requirements, but the supplied card-image evidence does not itself prove the corresponding fighter ability components.

The batch demonstrates an important corpus-level failure mode: quantity/source validation can pass while physical card type transcription remains materially wrong. Card type must therefore remain part of image-level QA rather than being inferred from existing canonical documentation.

## Final Assessment

The evidence archive itself is technically complete and usable: all five fighter ZIPs are readable, every manifest card has one decodable unique image, every image was inspected, and all deck quantities reconcile with canonical construction.

The current canonical Git corpus represented by the tested branch tip is **not acceptable as a fully verified Phase 4 card-image corpus for this batch**. Every fighter has at least one material P1/P2 discrepancy. Before the batch can receive a clean card-image QA verdict, the canonical representations need independent correction/reconciliation for the findings above and then a fresh verification against the same physical evidence.

No canonical manifests, schemas, source documentation, or other QA reports were modified by this QA persistence step.