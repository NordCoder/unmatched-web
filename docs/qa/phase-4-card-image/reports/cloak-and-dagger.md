# Phase 4 Card-Image QA Report

## Batch

- cloak-and-dagger
- dr-ellie-sattler
- genie
- houdini
- t-rex

## Verdict

PASS WITH FINDINGS

## Summary

Independent Phase 4 card-image QA completed for all five fighters in this batch.

- Repository: `NordCoder/unmatched-web`
- Canonical fighter/deck source branch used for this batch: `phase-4b-worker-b-licensed`
- Canonical branch tip observed immediately before persistence: `8a10d08cb3d3a074e8b5b0fd89cc21fddeba5544`
- Fighters fully checked: 5/5
- Total unique card images: 62
- Images successfully decoded and visually inspected: 62/62
- Unreadable images: 0
- Missing images: 0
- Extra images: 0
- Binary duplicate images: 0
- Quantity failures: 0
- P1 findings: 0
- P2 findings: 0
- P3 findings: 5
- Gameplay-semantic discrepancies: 0
- Blockers: 0

All 62 physical card components matched the canonical Git corpus for printed card identity, quantity, fighter restriction, card type, printed combat value, BOOST, trigger/timing, effect ordering, optionality, conditions, targets, movement/placement semantics, draw/discard/reveal/look-at distinctions, cancellation behavior, and relevant dependencies.

The five findings are non-gameplay `card_id` slug mismatches in the external evidence manifests. They do not change printed identity or gameplay semantics. Three fighters also retain already-documented integration qualifications, and T. Rex retains the repository's pre-existing qualified project normalization for the `Momentous Shift` two-space-overlap edge case.

## Fighter Results

### cloak-and-dagger

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/cloak-and-dagger.yaml`
- fighter manifest: `docs/fighters/phase-4b/cloak-and-dagger.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: PASS — 12 manifest entries, 12 images, 12/12 readable, 0 missing, 0 extra, 0 duplicates
- quantity validation: PASS — archive sum 30; Git available pool 30; Git game deck 30; 12 archive definitions; 12 Git definitions
- printed card identity/content verification: PASS — 12/12 unique cards inspected
- discrepancies: none
- integration requirements confirmed: none
- verdict: PASS

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Commanding Impact | PASS | PASS | PASS | PASS |
| Perfect Balance | PASS | PASS | PASS | PASS |
| Channel the Dark | PASS | PASS | PASS | PASS |
| Darkforce Dimension | PASS | PASS | PASS | PASS |
| Into Darkness | PASS | PASS | PASS | PASS |
| Into the Void | PASS | PASS | PASS | PASS |
| Lightforce Barrage | PASS | PASS | PASS | PASS |
| Living Shadow | PASS | PASS | PASS | PASS |
| Chosen Fate | PASS | PASS | PASS | PASS |
| The Living Light | PASS | PASS | PASS | PASS |
| Traverse the Darkforce | PASS | PASS | PASS | PASS |

#### Dependency-sensitive observations

- `Channel the Dark`: the printed effect requires the placement branch to resolve before the gained action. Git correctly gates the action on successful placement.
- `Living Shadow`: the printed swap must actually occur before Cloak replaces Dagger as the combat participant and the card becomes value 4. Git preserves that dependency.
- `Chosen Fate`: Git preserves the relationship between damage dealt to one friendly hero and recovery of the other hero, with the recovery amount derived from the actual damage result rather than flattening the operations into unrelated effects.

### dr-ellie-sattler

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/dr-ellie-sattler.yaml`
- fighter manifest: `docs/fighters/phase-4b/dr-ellie-sattler.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS WITH P3 MAPPING FINDINGS
- unique-card image completeness: PASS — 14 manifest entries, 14 images, 14/14 readable, 0 missing, 0 extra, 0 duplicates
- quantity validation: PASS — archive sum 30; Git available pool 30; Git game deck 30; 14 archive definitions; 14 Git definitions
- printed card identity/content verification: PASS — 14/14 unique cards inspected
- discrepancies: 2 P3 archive `card_id` slug mismatches; no gameplay discrepancy
- integration requirements confirmed: `B-REQ-006`
- verdict: PASS_WITH_QUALIFICATIONS

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Regroup | PASS | PASS | PASS | PASS |
| You Never Had Control, That's the Illusion | PASS | PASS* | PASS | PASS* |
| Woman Inherits the Earth | PASS | PASS | PASS | PASS |
| Violently, If Necessary | PASS | PASS | PASS | PASS |
| Sexism in Survival Situations | PASS | PASS | PASS | PASS |
| Chaotician | PASS | PASS | PASS | PASS |
| Life Finds a Way | PASS | PASS | PASS | PASS |
| The Future Ex-Mrs. Malcolm | PASS | PASS | PASS | PASS |
| Must Go Faster | PASS | PASS | PASS | PASS |
| Lock the Doors! | PASS | PASS | PASS | PASS |
| Hey! Hey! Hey! | PASS | PASS | PASS | PASS |
| The Concept of Attraction | PASS | PASS | PASS | PASS |
| I Think We're Back in Business | PASS | PASS* | PASS | PASS* |

`*` Printed identity, quantity, type, value, BOOST, restriction, and semantics match canonical Git. Qualification concerns only the external evidence-manifest `card_id`.

#### Integration requirement evidence

`B-REQ-006` is directly supported by the physical card corpus: multiple cards count, place, consume, or return Insight tokens on battlefield spaces. This confirms the need for positioned reusable battlefield-token instances. The fighter ability itself was not claimed as visually verified from these card images.

### genie

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/genie.yaml`
- fighter manifest: `docs/fighters/phase-4b/genie.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS WITH P3 MAPPING FINDINGS
- unique-card image completeness: PASS — 12 manifest entries, 12 images, 12/12 readable, 0 missing, 0 extra, 0 duplicates
- quantity validation: PASS — archive sum 30; Git available pool 30; Git game deck 30; 12 archive definitions; 12 Git definitions
- printed card identity/content verification: PASS — 12/12 unique cards inspected
- discrepancies: 2 P3 archive `card_id` slug mismatches; no gameplay discrepancy
- integration requirements confirmed: none
- verdict: PASS_WITH_QUALIFICATIONS

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Careful What You Wish For | PASS | PASS | PASS | PASS |
| Back in the Lamp | PASS | PASS | PASS | PASS |
| I Am Freed | PASS | PASS | PASS | PASS |
| I Grant You... Death | PASS | PASS | PASS | PASS |
| Imprisoned Wrath | PASS | PASS | PASS | PASS |
| I've Made Sultans Out of Less | PASS | PASS* | PASS | PASS* |
| Prisoner's Torment | PASS | PASS* | PASS | PASS* |
| This Is No Parlor Trick | PASS | PASS | PASS | PASS |
| Wishing for More Wishes | PASS | PASS | PASS | PASS |
| Your Wish Is My Command | PASS | PASS | PASS | PASS |
| Three Wishes | PASS | PASS | PASS | PASS |

`*` Printed identity and gameplay data are correct; only the archive slug differs from canonical Git ID.

#### Dependency-sensitive observation

`I Am Freed` was checked specifically for the printed `Then`. The Git effect preserves the printed operation order (place, then adjacent damage) without inventing an `if you do` dependency. This is consistent with the project's normalized effect model, which distinguishes ordered operations from explicit dependency grammar.

### houdini

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/houdini.yaml`
- fighter manifest: `docs/fighters/phase-4b/houdini.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS
- unique-card image completeness: PASS — 13 manifest entries, 13 images, 13/13 readable, 0 missing, 0 extra, 0 duplicates
- quantity validation: PASS — archive sum 30; Git available pool 30; Git game deck 30; 13 archive definitions; 13 Git definitions
- printed card identity/content verification: PASS — 13/13 unique cards inspected
- discrepancies: none
- integration requirements confirmed: `B-REQ-004`, `B-REQ-007`, `B-REQ-010`
- verdict: PASS_WITH_QUALIFICATIONS

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| For My Next Trick | PASS | PASS | PASS | PASS |
| An Illusion of My Own Design | PASS | PASS | PASS | PASS |
| Flourish | PASS | PASS | PASS | PASS |
| Grand Escape | PASS | PASS | PASS | PASS |
| Misdirection | PASS | PASS | PASS | PASS |
| Sleight of Hand | PASS | PASS | PASS | PASS |
| Smoke and Mirrors | PASS | PASS | PASS | PASS |
| The Big Reveal | PASS | PASS | PASS | PASS |
| Vanishing Act | PASS | PASS | PASS | PASS |
| A Magician Never Reveals His Secrets | PASS | PASS | PASS | PASS |
| All Part of the Show | PASS | PASS | PASS | PASS |
| And the Beautiful Bess! | PASS | PASS | PASS | PASS |
| Set the Stage | PASS | PASS | PASS | PASS |

#### Integration requirement evidence

- `B-REQ-004` — directly confirmed by printed `BOOSTED WITH` effects and cards such as `The Big Reveal`; implementation must retain card-instance BOOST events/history.
- `B-REQ-007` — directly confirmed by `And the Beautiful Bess!`, whose printed behavior distinguishes a card discarded from hand because of an opponent's effect from unrelated discard causes.
- `B-REQ-010` — directly confirmed by `A Magician Never Reveals His Secrets`; a card remaining in private hand may react to an opponent effect, reveal itself, and cancel the triggering effect instance.

#### Dependency-sensitive observation

`Sleight of Hand` was checked as a combat-card replacement composite. The Git structure correctly preserves the sequence: optionally return this combat card to hand, choose a different legal combat card only after the return succeeded, then replace and resolve the new combat card normally.

### t-rex

- branch: `phase-4b-worker-b-licensed`
- deck manifest: `docs/cards/phase-4b/t-rex.yaml`
- fighter manifest: `docs/fighters/phase-4b/t-rex.yaml`
- archive integrity: PASS
- canonical manifest comparison: PASS WITH QUALIFICATIONS
- unique-card image completeness: PASS — 11 manifest entries, 11 images, 11/11 readable, 0 missing, 0 extra, 0 duplicates
- quantity validation: PASS — archive sum 30; Git available pool 30; Git game deck 30; 11 archive definitions; 11 Git definitions
- printed card identity/content verification: PASS — 11/11 unique cards inspected
- discrepancies: 1 P3 archive `card_id` slug mismatch; no new gameplay discrepancy
- integration requirements confirmed: `B-REQ-001`
- verdict: PASS_WITH_QUALIFICATIONS

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Momentous Shift | PASS | PASS | PASS† | PASS† |
| Commanding Impact | PASS | PASS | PASS | PASS |
| Terrifying Roar | PASS | PASS | PASS | PASS |
| Closer Than She Appears | PASS | PASS | PASS | PASS |
| Thrash | PASS | PASS | PASS | PASS |
| Reckless Lunge | PASS | PASS | PASS | PASS |
| When Dinosaurs Ruled the Earth | PASS | PASS | PASS | PASS |
| 65 Million Years of Gut Instinct | PASS | PASS | PASS | PASS |
| Ripples in the Water | PASS | PASS | PASS | PASS |
| You're Just Making Her Angry | PASS | PASS* | PASS | PASS* |
| 15,000 Pounds of Muscle | PASS | PASS | PASS | PASS |

`†` The image verifies the published `Momentous Shift` wording but does not independently resolve the special two-space-fighter overlap edge case. Canonical Git explicitly marks that edge behavior as a medium-confidence `project_normalization`, with replacement if an authoritative ruling is later found. This is an existing qualification, not a new transcription discrepancy.

`*` Printed identity and gameplay facts match; only the external evidence-manifest slug differs.

#### Integration requirement evidence

`B-REQ-001` is confirmed at card-corpus level for effects that require T. Rex movement, placement, zone evaluation, and turn-start position/history. The card faces do not by themselves prove every underlying large-fighter board rule, so this confirmation is limited to the requirement's direct relevance to printed card behavior.

## Findings

### B07-SAT-P3-01 — archive card_id slug mismatch

- fighter: `dr-ellie-sattler`
- card: `You Never Had Control, That's the Illusion`
- expected: canonical Git ID `you-never-had-control-thats-the-illusion`
- observed: archive manifest ID `you-never-had-control-that-s-the-illusion`
- severity: P3
- evidence/reasoning: the physical card title, quantity, fighter restriction, card type, combat value, BOOST, and printed effect all map unambiguously to the canonical Git definition. The only difference is apostrophe/contraction slug normalization in the external evidence manifest.
- expected correction: no Git/gameplay correction; external archive mapping should use the canonical Git ID.

### B07-SAT-P3-02 — archive card_id slug mismatch

- fighter: `dr-ellie-sattler`
- card: `I Think We're Back in Business`
- expected: canonical Git ID `i-think-were-back-in-business`
- observed: archive manifest ID `i-think-we-re-back-in-business`
- severity: P3
- evidence/reasoning: physical card identity and gameplay facts match Git exactly; only apostrophe slug normalization differs.
- expected correction: no Git/gameplay correction; align the external evidence mapping with canonical Git ID.

### B07-GEN-P3-01 — archive card_id slug mismatch

- fighter: `genie`
- card: `I've Made Sultans Out of Less`
- expected: canonical Git ID `ive-made-sultans-out-of-less`
- observed: archive manifest ID `i-ve-made-sultans-out-of-less`
- severity: P3
- evidence/reasoning: printed identity, quantity, metadata, and effect semantics match the canonical definition. Difference is limited to slug formation.
- expected correction: no Git/gameplay correction; align external archive mapping.

### B07-GEN-P3-02 — archive card_id slug mismatch

- fighter: `genie`
- card: `Prisoner's Torment`
- expected: canonical Git ID `prisoners-torment`
- observed: archive manifest ID `prisoner-s-torment`
- severity: P3
- evidence/reasoning: printed card and canonical Git semantics match; archive slug splits the possessive apostrophe differently.
- expected correction: no Git/gameplay correction; align external archive mapping.

### B07-TRX-P3-01 — archive card_id slug mismatch

- fighter: `t-rex`
- card: `You're Just Making Her Angry`
- expected: canonical Git ID `youre-just-making-her-angry`
- observed: archive manifest ID `you-re-just-making-her-angry`
- severity: P3
- evidence/reasoning: printed title, quantity, restriction, type, value, BOOST, and effect semantics all match Git. Difference is only external slug normalization.
- expected correction: no Git/gameplay correction; align external archive mapping.

### Existing qualification — T. Rex Momentous Shift overlap

- fighter: `t-rex`
- card: `Momentous Shift`
- expected: repository should distinguish printed fact from project interpretation where published evidence does not resolve the two-space-overlap exact case
- observed: Git does so; the interpretation is explicitly marked `project_normalization`, confidence `medium`, with a replacement condition
- severity: warning / qualification, not a discrepancy
- evidence/reasoning: physical card evidence proves the printed wording but cannot adjudicate whether an overlapping start/current two-space footprint counts as having left the starting space. QA therefore does not elevate the project's selected deterministic interpretation into a printed fact.

### Integration-only qualifications

These are not transcription failures.

- `dr-ellie-sattler`: `B-REQ-006` — positioned reusable battlefield token instances.
- `houdini`: `B-REQ-004` — card-used-as-BOOST event/history.
- `houdini`: `B-REQ-007` — operation-cause provenance.
- `houdini`: `B-REQ-010` — private-zone reactive effect cancellation.
- `t-rex`: `B-REQ-001` — multi-space fighter footprint support.

## Corpus-Level Observations

1. The outer archive contained exactly five readable fighter ZIPs. All fighters resolved unambiguously to canonical IDs in the supplied Phase 4 mapping.
2. The expected invariant `1 unique card definition = 1 image = 1 archive manifest entry` held for all 62 definitions.
3. Physical-copy multiplicity was represented through `quantity`; no fighter duplicated image files merely to represent multiple card copies.
4. Every image decoded successfully. No zero-byte, truncated, corrupt, or unreadable card image was encountered.
5. SHA/binary duplicate checking found no duplicate image payloads.
6. Every archive manifest entry had a corresponding image, and every image had a manifest entry.
7. All five decks are fixed 30-card constructions in the canonical Git manifests. Each archive quantity sum was 30 and matched both `available_pool_count` and `game_deck_count`.
8. Generic cards such as `Feint`, `Regroup`, and `Commanding Impact` were visually inspected rather than assumed correct from familiarity.
9. Semantic verification was not based on literal prose equality. Printed text was interpreted into gameplay behavior and compared against the normalized staged/operation representation in Git.
10. Special attention was given to trigger windows (`IMMEDIATELY`, `DURING COMBAT`, `AFTER COMBAT`), optional vs mandatory language, effect order, `then`/dependency behavior, MOVE vs PLACE, hand visibility, cancellation, BOOST use/history, and combat-participant replacement.
11. No P1 or P2 discrepancy was found. The five P3 findings are confined to evidence-manifest identifier normalization and do not alter gameplay.
12. Existing `requires` entries were treated as integration requirements rather than transcription errors when the physical cards supported the documented gameplay semantics.
13. The T. Rex `Momentous Shift` overlap item remains correctly classified as a project normalization rather than an official printed ruling.
14. Fighter abilities were not claimed as visually proven unless directly represented by card components. This report is primarily card-corpus evidence.

## Final Assessment

This batch is suitable for the Phase 4 card-image corpus.

All 62 accessible unique physical card images were successfully inspected and matched the canonical repository corpus for gameplay-relevant printed facts and semantics. Deck construction and per-definition quantities are complete and consistent. There are no missing, unreadable, corrupt, duplicate, or unmatched card images and no material P1/P2 findings.

The batch verdict is **PASS WITH FINDINGS** solely because:

- five external archive `card_id` values use non-canonical apostrophe/contraction slug normalization;
- `dr-ellie-sattler`, `houdini`, and `t-rex` retain already-documented shared integration requirements;
- `t-rex` retains the already-documented evidence qualification for the `Momentous Shift` overlap exact case.

No canonical fighter or card manifest correction is indicated by this QA.
