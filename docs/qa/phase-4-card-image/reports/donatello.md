# Phase 4 Card-Image QA Report

## Batch

- donatello
- john-henry
- michelangelo
- raphael
- rosie-the-riveter

Repository: `NordCoder/unmatched-web`

Canonical branch checked: `phase-4b-worker-d-latest`

Canonical branch tip immediately before report persistence: `b9ae31c3b1a958e34bfb507d695cbd14650b9ed6`

Evidence archive: `unmatched-bundle-14.zip`

## Verdict

**FAIL**

The batch is technically complete and quantity-correct, and all 61 unique card images were readable and visually inspected. Three fighters contain material gameplay-changing discrepancies in the canonical Git normalization: Michelangelo, Rosie the Riveter, and John Henry. Donatello and Raphael have no P1/P2 discrepancies but retain non-material evidence/metadata qualifications.

## Summary

The outer archive contained exactly five fighter ZIPs: Donatello, Michelangelo, Raphael, Rosie the Riveter, and John Henry. All nested ZIPs opened correctly. Every fighter archive contained a readable `manifest.yaml`; every manifest entry mapped to exactly one readable image; no zero-byte images, missing images, extra images, or binary duplicate images were found.

All five decks reconcile to their canonical fixed 30-card construction. Across the batch there are 61 unique card definitions/images, and all 61 were inspected directly rather than sampled.

Material findings:

1. Michelangelo / **Shell Insertion**: Git grants owner-selected resolution ordering not granted by the printed card.
2. Michelangelo / **Shell Insertion**: the printed bottom-card-to-gain-action dependency is not explicitly preserved as a cost/dependency edge.
3. Rosie the Riveter / **D-Day**: Git scales damage/recovery by the number of upgrades already active instead of applying the printed fixed 1 damage / 1 recovery conditional result.
4. John Henry / **Twelve-Pound Hammer**: printed optional `may place` behavior is normalized as mandatory.

Non-material P3 observations concern evidence/canonical card slugs and printed-name punctuation only.

## Batch Integrity Totals

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| PASS | 0 |
| PASS_WITH_QUALIFICATIONS | 2 |
| FAIL | 3 |
| BLOCKED | 0 |
| Total unique card images | 61 |
| Images successfully inspected | 61 |
| Unreadable images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Duplicate binary images | 0 |
| Quantity failures | 0 |
| Material semantic discrepancies | 4 P1 |
| Non-material metadata/evidence findings | 4 P3 |

## Fighter Results

### donatello

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/donatello.yaml`
- fighter manifest: `docs/fighters/phase-4b/donatello.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS (14/14)**
- printed card identity/content verification: **PASS with one P3 evidence-ID qualification**
- discrepancies: `DON-P3-01`
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive and quantity evidence

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
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Bo Staff | PASS | PASS | PASS | PASS |
| Donatello Does Machines | PASS | PASS | PASS | PASS |
| Electro Grenade | PASS | PASS | PASS | PASS |
| Heroes in a Half Shell | PASS | PASS | PASS | PASS |
| Party Wagon! | PASS | PASS | PASS | PASS |
| Quick Strike | PASS | PASS | PASS | PASS |
| Self Defense Grid | PASS | PASS | PASS | PASS |
| Shift Focus | PASS | PASS | PASS | PASS |
| Short Circuit | PASS | PASS | PASS | PASS |
| Smoke Bomb | PASS | PASS | PASS | PASS |
| The Future of Ninjutsu | PASS | QUALIFIED P3 | PASS | PASS_WITH_QUALIFICATION |
| Thinking Ahead | PASS | PASS | PASS | PASS |
| Turtle Power! | PASS | PASS | PASS | PASS |
| Untested Enhancements | PASS | PASS | PASS | PASS |

All timing, card type, printed combat value, BOOST, fighter restriction, and gameplay effect semantics matched the canonical Git representation.

#### Integration requirements confirmed

- `D-REQ-001`: directly evidenced by staged/result-dependent behavior including Party Wagon!, Donatello Does Machines, and The Future of Ninjutsu.

---

### john-henry

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/john-henry.yaml`
- fighter manifest: `docs/fighters/phase-4b/john-henry.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL on one semantic normalization**
- unique-card image completeness: **PASS (12/12)**
- printed card identity/content verification: **FAIL**
- discrepancies: `JH-P1-01`
- verdict: **FAIL**

#### Archive and quantity evidence

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
| Bring Down the Mountain | PASS | PASS | PASS | PASS |
| Cool Drink of Water | PASS | PASS | PASS | PASS |
| Deeds of Valor | PASS | PASS | PASS | PASS |
| Eighteen-Pound Hammer | PASS | PASS | PASS | PASS |
| Hear That Cold Steel Ring | PASS | PASS | PASS | PASS |
| Knock Them Silly | PASS | PASS | PASS | PASS |
| Larger Than Life | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Nine-Pound Hammer | PASS | PASS | PASS | PASS |
| Power Through | PASS | PASS | PASS | PASS |
| Striking Fire | PASS | PASS | PASS | PASS |
| Twelve-Pound Hammer | PASS | PASS | FAIL | FAIL |

#### Integration requirements confirmed

- `D-REQ-001`: Knock Them Silly requires state carried from `IMMEDIATELY` into later opposing AFTER COMBAT resolution.
- `D-REQ-003`: Hear That Cold Steel Ring prints a free action restricted to attacking.
- `D-REQ-004`: Power Through demonstrates non-standard movement-distance accounting for track-token spaces.
- `D-REQ-007`: Hammer conditions, Twelve-Pound Hammer destinations, and Power Through directly depend on track/path state.

---

### michelangelo

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/michelangelo.yaml`
- fighter manifest: `docs/fighters/phase-4b/michelangelo.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL on two semantic normalizations**
- unique-card image completeness: **PASS (12/12)**
- printed card identity/content verification: **FAIL**
- discrepancies: `MIC-P1-01`, `MIC-P1-02`, `MIC-P3-01`, `MIC-P3-02`
- verdict: **FAIL**

#### Archive and quantity evidence

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
| Back for Seconds | PASS | PASS | PASS | PASS |
| Boisterous Beatdown | PASS | PASS | PASS | PASS |
| Cowabunga!! | PASS | PASS | PASS | PASS |
| Guaranteed Delivery | PASS | PASS | PASS | PASS |
| Hard-Hitting Investigation | PASS | PASS | PASS | PASS |
| Heroes in a Half Shell | PASS | PASS | PASS | PASS |
| Hi-Yaaaaaah!! | PASS | PASS | PASS | PASS |
| Let's Go!! | PASS | QUALIFIED P3 | PASS | PASS_WITH_QUALIFICATION |
| Michelangelo Is a Party Dude!! | PASS | PASS | PASS | PASS |
| Nunchaku | PASS | PASS | PASS | PASS |
| Shell Insertion | PASS | PASS | FAIL | FAIL |
| Turtle Power!! | PASS | QUALIFIED P3 | PASS | PASS_WITH_QUALIFICATION |

#### Integration requirements confirmed

- `D-REQ-001`: Shell Insertion confirms the need for generic choose-N/staged resolution, although the current specific normalization has semantic errors.
- `D-REQ-014`: Cowabunga!! directly confirms per-turn history of other played cards.
- `D-REQ-013` is fighter hand-size policy and was not claimed as visually confirmed by the action-card evidence ZIP.

---

### raphael

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/raphael.yaml`
- fighter manifest: `docs/fighters/phase-4b/raphael.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS (12/12)**
- printed card identity/content verification: **PASS with one P3 evidence-ID qualification**
- discrepancies: `RAP-P3-01`
- verdict: **PASS_WITH_QUALIFICATIONS**

#### Archive and quantity evidence

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
| Batter Up! | PASS | PASS | PASS | PASS |
| Break Something | PASS | PASS | PASS | PASS |
| Crowd Control | PASS | PASS | PASS | PASS |
| Heroes in a Half Shell | PASS | PASS | PASS | PASS |
| Let's Do This! | PASS | QUALIFIED P3 | PASS | PASS_WITH_QUALIFICATION |
| Payback Time! | PASS | PASS | PASS | PASS |
| Raphael Is Cool But Rude | PASS | PASS | PASS | PASS |
| Relentless | PASS | PASS | PASS | PASS |
| Sai | PASS | PASS | PASS | PASS |
| Slapshot | PASS | PASS | PASS | PASS |
| Turtle Power! | PASS | PASS | PASS | PASS |
| Unbridled Rage | PASS | PASS | PASS | PASS |

#### Integration requirements confirmed

- `D-REQ-001`: Break Something directly confirms the need to preserve actual cancellation result between `IMMEDIATELY` and `AFTER COMBAT`; Unbridled Rage also confirms choose-many behavior.

---

### rosie-the-riveter

- branch: `phase-4b-worker-d-latest`
- deck manifest: `docs/cards/phase-4b/rosie-the-riveter.yaml`
- fighter manifest: `docs/fighters/phase-4b/rosie-the-riveter.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **FAIL on one semantic normalization**
- unique-card image completeness: **PASS (11/11)**
- printed card identity/content verification: **FAIL**
- discrepancies: `ROS-P1-01`
- verdict: **FAIL**

#### Archive and quantity evidence

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
| Unique Git definitions | 11 |
| Construction | fixed |

#### Card verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Arc Welding | PASS | PASS | PASS | PASS |
| D-Day | PASS | PASS | FAIL | FAIL |
| "E" Award | PASS | PASS | PASS | PASS |
| Full Metal | PASS | PASS | PASS | PASS |
| Hasty Repairs | PASS | PASS | PASS | PASS |
| Loose Lips Sink Ships! | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Overrun | PASS | PASS | PASS | PASS |
| Rapid Development | PASS | PASS | PASS | PASS |
| Technological Superiority | PASS | PASS | PASS | PASS |
| We Can Do It! | PASS | PASS | PASS | PASS |

#### Integration requirements confirmed

- `D-REQ-003`: D-Day and Rapid Development explicitly contain a restricted free action to attack.
- `D-REQ-012`: activate/deactivate/activate-all/count-active upgrade state is repeatedly evidenced by the cards.
- Fighter-level `D-REQ-004` and `D-REQ-005` concern fighter ability behavior not printed on these action-card images and therefore were not claimed as directly visually confirmed by this batch.

## Findings

### MIC-P1-01 — Shell Insertion resolution order

- fighter: `michelangelo`
- card: `Shell Insertion`
- Git location: `docs/cards/phase-4b/michelangelo.yaml`, `shell-insertion` / `resolution-order`
- expected: after choosing two different printed effects, resolve them in published printed order unless an authoritative source explicitly grants ordering authority.
- observed: Git adds a separate required `resolution-order` choice with domain `permutations_of_bound_selected_effects`, allowing the controller to reorder the chosen effects.
- severity: **P1**
- evidence/reasoning: the physical card says to choose two different effects but does not grant a second choice to determine their order. The repository's normalized effect model also states that operations/effects resolve in published/source-defined order unless a choice explicitly grants ordering authority. The current Git model therefore creates legal gameplay behavior that is not printed on the card.
- expected correction: remove unsubstantiated player-controlled permutation ordering and resolve the selected effects in their printed order unless an authoritative ruling proves otherwise.

### MIC-P1-02 — Shell Insertion bottom-card/action dependency

- fighter: `michelangelo`
- card: `Shell Insertion`
- Git location: `docs/cards/phase-4b/michelangelo.yaml`, `bottom_card_for_action` branch
- expected: putting a card from hand on the bottom of the deck is the prerequisite/cost for gaining the action; failure or inability to perform that first part must prevent the dependent action gain.
- observed: Git represents the selected hand-card choice, `MOVE_CARD`, and `GAIN_ACTION` as ordinary operations without an explicit cost/dependency edge.
- severity: **P1**
- evidence/reasoning: printed wording links the two clauses as a single dependent effect. The project's effect model explicitly requires preserving distinctions among condition, cost/prerequisite, choice, and operation, and warns against applying generic partial resolution across explicit dependency edges.
- expected correction: model the bottom-card operation as the required prerequisite/cost for `GAIN_ACTION`, with no action granted unless the prerequisite is actually performed.

### ROS-P1-01 — D-Day damage/recovery scaling

- fighter: `rosie-the-riveter`
- card: `D-Day`
- Git location: `docs/cards/phase-4b/rosie-the-riveter.yaml`, `d-day`
- expected: after gaining the restricted attack action and activating all upgrades, if the printed precondition that the upgrades were already active is satisfied, deal **1** damage to the opposing fighter and recover **1** health for Rosie.
- observed: Git snapshots `count_active` before `activate_all`, then uses that numeric count as `amount_from` for both damage and recovery.
- severity: **P1**
- evidence/reasoning: the physical card prints a fixed one-damage and one-recovery consequence behind a condition. It does not say to deal/recover one per already-active upgrade. The current representation can therefore produce 2/2, 3/3, etc., which changes gameplay.
- expected correction: preserve the source-defined precondition as a boolean/historical condition and, when satisfied, apply fixed 1 damage and fixed 1 recovery rather than scaling by active-upgrade count.

### JH-P1-01 — Twelve-Pound Hammer optional placement

- fighter: `john-henry`
- card: `Twelve-Pound Hammer`
- Git location: `docs/cards/phase-4b/john-henry.yaml`, `twelve-place`
- expected: when the four-different-paths condition is satisfied, the controller **may** place a fighter from the combat on a space with a deployed track token. Declining the effect is legal.
- observed: Git makes both fighter and destination choices `optional: false` and follows with unconditional `PLACE`.
- severity: **P1**
- evidence/reasoning: the physical card explicitly uses optional `may` wording. The Git normalization converts an optional relocation into a mandatory one whenever the condition is true.
- expected correction: make the entire placement effect optional; if declined, neither combat fighter is relocated.

### DON-P3-01 — The Future of Ninjutsu evidence card ID

- fighter: `donatello`
- card: `The Future of Ninjutsu`
- expected: evidence mapping should use canonical card ID `future-of-ninjutsu`.
- observed: evidence/archive manifest uses `the-future-of-ninjutsu` while Git uses `future-of-ninjutsu`.
- severity: **P3**
- evidence/reasoning: physical identity, quantity, metadata, and gameplay semantics are unambiguous and match Git; only the evidence slug differs.
- expected correction: normalize the evidence-pack mapping to the canonical ID. No canonical gameplay change is required.

### MIC-P3-01 — Let's Go!! evidence identity metadata

- fighter: `michelangelo`
- card: `Let's Go!!`
- expected: evidence mapping should use canonical ID `lets-go`; printed-name metadata should preserve the physical title if the canonical `name` field is intended as exact printed identity.
- observed: evidence uses `let-s-go`; Git uses `lets-go` and `Let's Go`, while the card image prints `LET'S GO!!`.
- severity: **P3**
- evidence/reasoning: gameplay data and card mapping are unambiguous; discrepancy affects slug/punctuation only.
- expected correction: evidence mapping to `lets-go`; optionally synchronize printed-name punctuation in canonical metadata according to project naming policy.

### MIC-P3-02 — Turtle Power!! punctuation

- fighter: `michelangelo`
- card: `Turtle Power!!`
- expected: exact printed-name metadata should preserve `Turtle Power!!` when the field represents physical title text.
- observed: Git stores `Turtle Power!`.
- severity: **P3**
- evidence/reasoning: one exclamation mark differs; type, value, BOOST, restriction, quantity, and effect are correct.
- expected correction: synchronize printed-name punctuation if exact title fidelity is required. No gameplay correction is required.

### RAP-P3-01 — Let's Do This! evidence card ID

- fighter: `raphael`
- card: `Let's Do This!`
- expected: evidence mapping should use canonical card ID `lets-do-this`.
- observed: evidence/archive manifest uses `let-s-do-this` while Git uses `lets-do-this`.
- severity: **P3**
- evidence/reasoning: physical identity, quantity, metadata, and semantics are unambiguous; only evidence slug differs.
- expected correction: normalize evidence-pack mapping to `lets-do-this`. No canonical gameplay correction is required.

## Corpus-Level Observations

1. **Archive construction is sound.** The outer batch contains the expected five independently readable fighter archives. Every unique card definition has exactly one image and one manifest entry; copies are represented through quantity rather than duplicate physical images.

2. **All five decks reconcile to canonical quantity.** Donatello 30/30, Michelangelo 30/30, Raphael 30/30, Rosie the Riveter 30/30, and John Henry 30/30. No quantity discrepancy or non-standard construction issue was found in this batch.

3. **The evidence set is fully visually inspectable.** All 61 unique images decoded successfully and were inspected. No `IMAGE_UNREADABLE` result was necessary.

4. **Filename/manifest identity cannot be treated as authoritative.** Several evidence card IDs normalize punctuation/apostrophes differently from canonical IDs. Direct image inspection allowed all such mappings to be resolved without ambiguity.

5. **Most canonical semantics are strong.** 57 of 61 unique card definitions did not expose a material semantic mismatch. The four material findings are localized to three cards: Shell Insertion, D-Day, and Twelve-Pound Hammer.

6. **The main error class is control-flow fidelity, not transcription of numbers/types.** No wrong quantity, combat type, combat value, BOOST, or fighter restriction was found. The material failures concern resolution order, dependency, conditional scaling, and optional-versus-mandatory behavior.

7. **Shared requirements remain legitimate qualifications rather than automatic failures.** Directly evidenced requirements include `D-REQ-001`, `D-REQ-003`, `D-REQ-004`, `D-REQ-007`, `D-REQ-012`, and `D-REQ-014`. Their presence is not itself a card transcription discrepancy.

8. **Fighter abilities not present in the card-image evidence were not claimed as visually verified.** In particular, Michelangelo `D-REQ-013` hand-size policy and Rosie fighter-level ability requirements outside the action-card corpus were not promoted to image-confirmed facts.

9. **No repository changes were made during the QA itself.** The QA was read-only; this file is the separate persistence action requested after completion.

## Final Assessment

The evidence archive itself is suitable as a Phase 4 physical card-image corpus: it is structurally complete, fully readable, free of duplicate/missing images, and quantity-consistent with Git for all five fighters.

The corresponding canonical repository state is **not ready to be accepted as fully verified for this batch** because three fighters contain material P1 semantic discrepancies:

- Michelangelo: two independent Shell Insertion control-flow errors;
- Rosie the Riveter: D-Day incorrectly scales a fixed conditional effect;
- John Henry: Twelve-Pound Hammer loses printed optionality.

Donatello and Raphael are clean with respect to material gameplay semantics, subject only to minor P3 evidence/metadata qualifications.

**Final batch assessment: FAIL — evidence corpus usable, canonical card semantics require correction before this batch can be considered fully Phase-4-verified.**
