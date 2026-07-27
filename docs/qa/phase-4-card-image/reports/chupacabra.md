# Phase 4 Card-Image QA Report

## Batch

- chupacabra
- leonardo
- loki
- muhammad-ali
- pandora

Source archive: `unmatched-bundle-13.zip`

QA canonical branch: `phase-4b-worker-d-latest`

Canonical branch tip inspected before report persistence: `b9ae31c3b1a958e34bfb507d695cbd14650b9ed6`

Operating mode: independent read-only QA of card-image evidence versus canonical Git documentation. No canonical fighter/card manifests were modified by this QA.

## Verdict

**FAIL**

## Summary

The batch contained 5 readable fighter ZIPs and 59 unique card images. All nested archives passed ZIP integrity checks, every image decoded successfully, no zero-byte or truncated image was found, no duplicate binary image was found, and no action-card image was missing from the archive manifests.

All five fighter archive quantity totals reconcile with the canonical fixed 30-card decks. No quantity failure or canonical printed metadata mismatch was found.

Two material canonical semantic discrepancies were identified:

1. **P1 — Loki / Looking for Trouble:** the normalized Git representation does not preserve the printed condition and replacement order. In particular, the opponent's current defense must first return to their hand before Loki inspects and selects the replacement, meaning the returned card may legally be selected again. The current representation selects from the hand before the old defense is returned and lacks the printed `if the opponent played a card` gate.
2. **P2 — Chupacabra / Natural Toughness:** the printed effect is an **AFTER COMBAT** return-to-hand effect, while Git represents it as `cleanup_replacement`, producing a narrower timing mismatch with possible edge-case consequences.

Three fighters (`pandora`, `muhammad-ali`, `leonardo`) have no material P1/P2 card discrepancy. Several archive-only `card_id`/display-name slug mismatches were found for Loki, Pandora, and Leonardo; canonical Git identities are correct, so these are P3 evidence-package issues rather than corpus corrections.

Batch totals:

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| PASS | 0 |
| PASS_WITH_QUALIFICATIONS | 3 |
| FAIL | 2 |
| BLOCKED | 0 |
| Total unique card images | 59 |
| Images successfully inspected | 59 |
| Unreadable images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Duplicate binary images | 0 |
| Quantity failures | 0 |
| Canonical printed-metadata discrepancies | 0 |
| Archive card-ID mapping discrepancies | 9 cards |
| Material semantic discrepancies | 2 cards |

## Fighter Results

### chupacabra

Canonical paths:

- deck: `docs/cards/phase-4b/chupacabra.yaml`
- fighter: `docs/fighters/phase-4b/chupacabra.yaml`

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: **PASS**
- manifest readable: **PASS**
- images readable: **11/11**
- unique images: **11**
- manifest card entries: **11**
- duplicate images: **0**
- missing images: **0**
- extra images: **0**

#### Canonical manifest comparison

- canonical construction: `fixed`
- archive quantity sum: **30**
- Git available pool: **30**
- Git game deck: **30**
- unique archive cards: **11**
- unique Git definitions: **11**
- per-definition quantity comparison: **PASS**

#### Unique-card image completeness

All 11 action-card definitions have exactly one readable image and one archive manifest entry. Physical copies are represented through `quantity`; no copy duplication was detected in the evidence package.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Blood in the Air | PASS | PASS | PASS | PASS |
| Feeding | PASS | PASS | PASS | PASS |
| Ambush | PASS | PASS | PASS | PASS |
| The More They Struggle | PASS | PASS | PASS | PASS |
| Wounded Beast | PASS | PASS | PASS | PASS |
| Ravenous Lunge | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |
| Tooth and Tail | PASS | PASS | PASS | PASS |
| Natural Toughness | PASS | PASS | **FAIL** | **FAIL** |
| Traveler of the Night | PASS | PASS | PASS | PASS |
| Unsettle | PASS | PASS | PASS | PASS |

#### Discrepancies

##### CHUPACABRA-001

- fighter: `chupacabra`
- card: `Natural Toughness`
- Git location: `docs/cards/phase-4b/chupacabra.yaml`, `natural-toughness`
- severity: **P2**
- expected / image fact: the printed effect uses **AFTER COMBAT** timing; if Chupacabra lost the combat, the card returns to the player's hand.
- observed / Git representation: `trigger: cleanup_replacement` with a move from the resolution zone to hand.
- evidence/reasoning: the return is printed as an After Combat effect, not merely a cleanup disposition rule. Delaying the zone transition to cleanup can alter the hand state during the remainder of the After Combat window and therefore may change edge-case legal behavior.
- expected semantic correction: preserve the printed **AFTER COMBAT** timing while ensuring the card does not later enter discard; an explicit after-combat resolution-to-hand transition or another representation with equivalent timing is required.

#### Integration requirements confirmed

- `D-REQ-001` — `Tooth and Tail` directly supports dependent staged resolution and operation-result binding: recovery depends on how many selected adjacent fighters were actually damaged.

The fighter-level/shared integration status may include other requirements, but this report claims direct visual confirmation only where the card images themselves provide evidence.

#### Verdict

**FAIL**

Reason: complete and quantity-correct card evidence, but one material timing-level semantic discrepancy exists for `Natural Toughness`.

---

### leonardo

Canonical paths:

- deck: `docs/cards/phase-4b/leonardo.yaml`
- fighter: `docs/fighters/phase-4b/leonardo.yaml`

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: **PASS**
- manifest readable: **PASS**
- images readable: **12/12**
- unique images: **12**
- manifest card entries: **12**
- duplicate images: **0**
- missing images: **0**
- extra images: **0**

#### Canonical manifest comparison

- canonical construction: `fixed`
- archive quantity sum: **30**
- Git available pool: **30**
- Git game deck: **30**
- unique archive cards: **12**
- unique Git definitions: **12**
- per-definition quantity comparison: **PASS**

#### Unique-card image completeness

All 12 canonical action-card definitions are represented by exactly one readable image. No physical-copy duplication exists in the evidence package.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Katana | PASS | PASS | PASS | PASS |
| For Sensei | PASS | PASS | PASS | PASS |
| Turtle Power! | PASS | PASS | PASS | PASS |
| Fearless Leader | PASS | PASS | PASS | PASS |
| Quick Strike | PASS | PASS | PASS | PASS |
| Spatial Awareness | PASS | PASS | PASS | PASS |
| Eat, Sleep, and Breathe Ninjutsu | PASS | PASS* | PASS | PASS* |
| Wise Beyond His Years | PASS | PASS | PASS | PASS |
| I Have a Plan | PASS | PASS | PASS | PASS |
| Heroes in a Half Shell | PASS | PASS | PASS | PASS |
| Protective Father | PASS | PASS | PASS | PASS |
| Leonardo Leads | PASS | PASS | PASS | PASS |

`*` indicates an evidence-archive mapping/name issue only; canonical Git is correct.

#### Discrepancies

##### LEONARDO-001

- fighter: `leonardo`
- card: `Eat, Sleep, and Breathe Ninjutsu`
- severity: **P3**
- expected: canonical printed title `Eat, Sleep, and Breathe Ninjutsu`; canonical Git ID `eat-sleep-and-breathe-ninjutsu`.
- observed: archive uses `eat-sleep-and-breath-ninjutsu` and the archive display name also uses `breath` rather than `breathe`.
- evidence/reasoning: the card image itself clearly shows `Breathe`; Git matches the physical component. This is an archive manifest typo, not a canonical corpus error.
- expected correction: evidence archive mapping/name only. No Git correction is indicated.

#### Project normalization qualification

`For Sensei` confirms that both combat fighters may move up to 1 space. The printed card does not specify which of the two moves is resolved first. Git explicitly labels its sequential-order chooser as a medium-confidence `project_normalization`, rather than presenting it as a printed fact or official ruling. This is an appropriate qualification, not a transcription discrepancy.

#### Integration requirements confirmed

- `D-REQ-001` — directly supported by `For Sensei` and `Eat, Sleep, and Breathe Ninjutsu`, both of which require staged/multi-fighter position resolution whose later legal choices can depend on earlier movement/placement.

#### Verdict

**PASS_WITH_QUALIFICATIONS**

All printed card gameplay data match Git. Qualification: archive spelling/ID issue plus the already-explicit `For Sensei` project ordering normalization and shared integration requirement.

---

### loki

Canonical paths:

- deck: `docs/cards/phase-4b/loki.yaml`
- fighter: `docs/fighters/phase-4b/loki.yaml`

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: **PASS**
- manifest readable: **PASS**
- images readable: **11/11**
- unique images: **11**
- manifest card entries: **11**
- duplicate images: **0**
- missing images: **0**
- extra images: **0**

#### Canonical manifest comparison

- canonical construction: `fixed`
- archive quantity sum: **30**
- Git available pool: **30**
- Git game deck: **30**
- unique archive cards: **11**
- unique Git definitions: **11**
- per-definition quantity comparison: **PASS**

TRICK cards correctly have no printed BOOST in canonical Git.

#### Unique-card image completeness

All 11 unique action-card definitions have one readable image and one archive manifest entry. No binary duplicate images exist.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Trick: Baldr's Downfall | PASS | PASS* | PASS | PASS* |
| Looking for Trouble | PASS | PASS | **FAIL** | **FAIL** |
| Shapeshifter | PASS | PASS* | PASS | PASS* |
| Ragnarök | PASS | PASS | PASS | PASS |
| Malicious Flyting | PASS | PASS | PASS | PASS |
| Underhanded | PASS | PASS | PASS | PASS |
| Trick: Freyja's Rescue | PASS | PASS* | PASS | PASS* |
| God of Mischief | PASS | PASS | PASS | PASS |
| Laevateinn | PASS | PASS | PASS | PASS |
| Trick: Sindri's Bet | PASS | PASS* | PASS | PASS* |
| Trick: Svadilfari's Lure | PASS | PASS* | PASS | PASS* |

`*` indicates archive `card_id` spelling/slug mismatch only.

#### Discrepancies

##### LOKI-001

- fighter: `loki`
- card: `Looking for Trouble`
- Git location: `docs/cards/phase-4b/loki.yaml`, `looking-for-trouble`
- severity: **P1**
- expected / image fact:
  1. the effect applies **if the opponent played a card**;
  2. return that current combat card to the opponent's hand;
  3. Loki looks at the resulting hand;
  4. Loki chooses a card from that hand for the opponent to play in the same combat.
- observed / Git representation: the normalized effect immediately requests a mandatory replacement card from `combat_opponent_hand`, then performs `REPLACE_COMBAT_CARD`; the rule says the replaced card is returned to the original hand as part of the replacement operation.
- evidence/reasoning:
  - the explicit printed condition that the opponent actually played a combat card is missing from the normalized effect;
  - Git constructs the replacement choice domain before the old defense is returned to hand;
  - therefore the returned defense is excluded from the selection domain even though the printed procedure makes it part of the inspected hand and allows it to be chosen again;
  - the printed return -> inspect -> choose order is gameplay-significant and must be preserved.
- expected semantic correction: model a conditional staged current-combat replacement: verify that the opponent played a card -> return it to hand -> expose the required hand inspection to Loki -> choose a legal replacement from the resulting hand, including the returned card -> commit that selected card into the same combat without restarting the Attack action.

##### LOKI-002

- fighter: `loki`
- cards: five archive mappings
- severity: **P3**
- expected canonical IDs:
  - `trick-baldrs-downfall`
  - `shapeshifter`
  - `trick-freyjas-rescue`
  - `trick-sindris-bet`
  - `trick-svadilfaris-lure`
- observed archive IDs:
  - `trick-baldr-s-downfall`
  - `shapershifter`
  - `trick-freyja-s-rescue`
  - `trick-sindri-s-bet`
  - `trick-svadilfari-s-lure`
- additional observed archive metadata issue: `Shapeshifter` is displayed as `Shapershifter` in the evidence manifest.
- evidence/reasoning: every image identifies the intended card unambiguously; canonical Git identities are correct.
- expected correction: evidence archive mapping only. No canonical Git correction is indicated.

#### Integration requirements confirmed

- `D-REQ-001` — `Trick: Freyja's Rescue` requires capture of both current combat values before applying the two-value swap.
- `D-REQ-017` — `Trick: Svadilfari's Lure` directly confirms hidden-information selection semantics: Loki selects a concrete card instance from another player's hand without prior permission to inspect non-TRICK identities; only the selected card is then revealed and discarded.

`D-REQ-004` belongs primarily to Loki's fighter-level dynamic Maneuver movement rule and is not independently established by this card-only evidence package.

The fighter manifest's multiplayer non-combat TRICK recipient chooser remains correctly labeled as a medium-confidence `project_normalization`; this QA found no card image that promotes that normalization to a printed or official fact.

#### Verdict

**FAIL**

Reason: one P1 gameplay-changing replacement-flow error on `Looking for Trouble`.

---

### muhammad-ali

Canonical paths:

- deck: `docs/cards/phase-4b/muhammad-ali.yaml`
- fighter: `docs/fighters/phase-4b/muhammad-ali.yaml`

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: **PASS**
- manifest readable: **PASS**
- images readable: **13/13**
- unique images: **13**
- manifest card entries: **13**
- duplicate images: **0**
- missing images: **0**
- extra images: **0**

#### Canonical manifest comparison

- canonical construction: `fixed`
- archive quantity sum: **30**
- Git available pool: **30**
- Git game deck: **30**
- unique archive cards: **13**
- unique Git definitions: **13**
- per-definition quantity comparison: **PASS**

#### Unique-card image completeness

All 13 canonical action-card definitions have exactly one readable image. No extra/copy images were found.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| The Greatest | PASS | PASS | PASS | PASS |
| Stronger Than the Skill | PASS | PASS | PASS | PASS |
| Ali Shuffle | PASS | PASS | PASS | PASS |
| Jab | PASS | PASS | PASS | PASS |
| Champion of the World | PASS | PASS | PASS | PASS |
| Close and Clinch | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Stick and Move | PASS | PASS | PASS | PASS |
| Fancy Footwork | PASS | PASS | PASS | PASS |
| Hard to Be Humble | PASS | PASS | PASS | PASS |
| Rope-a-Dope | PASS | PASS | PASS | PASS |
| Answer the Bell | PASS | PASS | PASS | PASS |
| Louisville Lip | PASS | PASS | PASS | PASS |

The cards' Float-dependent branches are consistently normalized through `ali_stance: float`; this is compatible with the fighter manifest's explicit Float/Sting state model.

#### Discrepancies

None.

#### Integration requirements confirmed

- `D-REQ-006` — `Close and Clinch` directly confirms a turn-scoped leave-space restriction applying to Muhammad Ali and the opposing combat fighter while still allowing ordinary defeat/removal semantics.

`D-REQ-005` concerns Ali's Float attack-reach ability, not action-card text, and therefore is not independently visually confirmed by this card-only batch.

#### Verdict

**PASS_WITH_QUALIFICATIONS**

All 13 accessible action-card images match canonical Git. Qualification is integration-level only; no card transcription discrepancy was found.

---

### pandora

Canonical paths:

- deck: `docs/cards/phase-4b/pandora.yaml`
- fighter: `docs/fighters/phase-4b/pandora.yaml`

#### Archive integrity

- archive integrity: **PASS**
- nested ZIP readable: **PASS**
- manifest readable: **PASS**
- images readable: **12/12**
- unique images: **12**
- manifest card entries: **12**
- duplicate images: **0**
- missing images: **0**
- extra images: **0**

#### Canonical manifest comparison

- canonical construction: `fixed`
- archive quantity sum: **30**
- Git available pool: **30**
- Git game deck: **30**
- unique archive action cards: **12**
- unique Git action-card definitions: **12**
- per-definition quantity comparison: **PASS**

Canonical Git separately defines seven Pandora Miseries as `external_definitions`. They do not count toward the 30-card action deck and their absence from this action-card evidence ZIP is therefore not treated as a missing action-card failure.

#### Unique-card image completeness

All 12 canonical action-card definitions have exactly one readable image and one archive manifest entry.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Hera's Curiosity | PASS | PASS* | PASS | PASS* |
| Aphrodite's Beauty | PASS | PASS* | PASS | PASS* |
| Hindsight | PASS | PASS | PASS | PASS |
| Celestial Raiments | PASS | PASS | PASS | PASS |
| Offering to the Gods | PASS | PASS | PASS | PASS |
| Guided by the Fates | PASS | PASS | PASS | PASS |
| Divine Intervention | PASS | PASS | PASS | PASS |
| Zeus's Mischief | PASS | PASS* | PASS | PASS* |
| Forged by Hephaestus | PASS | PASS | PASS | PASS |
| Spite | PASS | PASS | PASS | PASS |
| Malice | PASS | PASS | PASS | PASS |
| Feint | PASS | PASS | PASS | PASS |

`*` indicates archive-only card-ID punctuation/slug differences.

#### Discrepancies

##### PANDORA-001

- fighter: `pandora`
- cards: `Hera's Curiosity`, `Aphrodite's Beauty`, `Zeus's Mischief`
- severity: **P3**
- expected canonical IDs:
  - `heras-curiosity`
  - `aphrodites-beauty`
  - `zeuss-mischief`
- observed archive IDs:
  - `hera-s-curiosity`
  - `aphrodite-s-beauty`
  - `zeus-s-mischief`
- evidence/reasoning: printed card identities unambiguously match Git; only the evidence archive slug convention differs.
- expected correction: evidence archive mapping only. No Git correction is indicated.

No printed gameplay discrepancy was found among the 12 action-card images.

#### Integration requirements confirmed

- `D-REQ-001` — action-card evidence directly supports ordered/dependent stages, especially `Guided by the Fates`, where the top deck card is discarded before its BOOST is consumed, and ordered multi-step resolution such as `Zeus's Mischief`.

`D-REQ-002` concerns the fighter ability's repeated Pandora's Box reveal/resolve/continue/forced-stop procedure. The seven Misery components and the ability card/rules component were not included in this evidence ZIP, so this report does **not** claim that `D-REQ-002` is visually confirmed by the current card corpus.

#### Verdict

**PASS_WITH_QUALIFICATIONS**

All accessible action-card images match Git. Qualification: Misery/ability components fall outside this ZIP, while pair-level integration remains partial.

## Findings

### Material findings

#### LOKI-001 — Looking for Trouble replacement flow

- fighter: `loki`
- card: `Looking for Trouble`
- severity: **P1**
- expected: only when the opponent has played a combat card, return that card to their hand, inspect the resulting hand, then select a replacement from that resulting hand. The just-returned card remains a legal candidate.
- observed: Git requests a replacement from the pre-return hand and lacks the explicit opponent-played-card gate; return of the current defense is encoded as part of the later replacement operation.
- evidence/reasoning: this loses a legal candidate and changes behavior for an undefended attack. It also reverses printed semantic order. This is gameplay-changing.

#### CHUPACABRA-001 — Natural Toughness timing

- fighter: `chupacabra`
- card: `Natural Toughness`
- severity: **P2**
- expected: return to hand in the **AFTER COMBAT** window if Chupacabra lost the combat.
- observed: Git models return as `cleanup_replacement`.
- evidence/reasoning: the resulting final destination is the same, but hand state during later After Combat processing can differ. This is an edge-case timing discrepancy.

### Non-material evidence-package findings

#### LOKI-002 — archive IDs/naming

Severity: **P3**

Five Loki archive identifiers differ from canonical Git; one archive display name misspells `Shapeshifter` as `Shapershifter`. Physical image identities and canonical Git are unambiguous and correct.

#### PANDORA-001 — archive IDs

Severity: **P3**

Three apostrophe-derived Pandora archive IDs differ from canonical Git IDs. No gameplay or canonical metadata error results.

#### LEONARDO-001 — archive spelling/ID

Severity: **P3**

The archive uses `breath` where the printed card and canonical Git use `breathe` for `Eat, Sleep, and Breathe Ninjutsu`.

## Corpus-Level Observations

1. **Archive technical integrity is strong.** All 59 images decode, all nested ZIPs pass integrity checks, and there are no binary duplicates, zero-byte files, manifest-orphan images, or action-card manifest entries without images.

2. **Action-deck construction reconciles cleanly.** Every fighter in this batch is a fixed 30-card deck and every archive quantity sum equals the Git `available_pool_count` and `game_deck_count`. This is an observed result for this batch only and is not treated as a universal Unmatched deck-size rule.

3. **One image per unique definition is respected.** Repeated physical copies are represented by `quantity` rather than duplicated image files.

4. **Archive manifests are useful mappings but are not authoritative.** Loki, Pandora, and Leonardo demonstrate why identity must ultimately be verified from the printed component and canonical Git rather than blindly trusting archive `card_id` spelling.

5. **External gameplay definitions remain distinct from action-card instances.** Pandora's seven Miseries are correctly excluded from the 30-card action deck. Their absence from this action-card archive does not constitute a missing ordinary card image.

6. **Structured normalization is generally faithful.** Staged effects, bindings, costs, branches, `MOVE`/`PLACE`, hidden-information distinctions, value layers, and state-based conditions successfully preserve printed semantics on almost all cards. The two material findings are specifically about semantic ordering/timing, not a general failure of the normalized model.

7. **Shared requirements are not transcription failures by themselves.** Directly card-supported integration requirements confirmed in this batch include `D-REQ-001`, `D-REQ-006`, and `D-REQ-017`. A card can remain visually/transcription-correct while requiring future shared engine support.

8. **Project normalizations are acceptable when explicitly marked.** Leonardo's `For Sensei` move-order authority remains an explicit medium-confidence project normalization because the printed card does not resolve the ordering authority. Loki's multiplayer non-combat TRICK-recipient chooser is likewise not promoted to an official printed fact by this batch.

9. **Fighter abilities outside the evidence ZIP are not claimed as visually verified.** This specifically limits direct image confirmation of Ali's Float attack-reach behavior, Loki's Maneuver movement modifier, and Pandora's Box repeat/stop procedure.

## Final Assessment

The batch is **not ready for unconditional Phase 4 card-image corpus acceptance** because two canonical semantic discrepancies remain:

- a **P1** gameplay-changing replacement-flow error for Loki's `Looking for Trouble`;
- a **P2** After Combat versus cleanup timing discrepancy for Chupacabra's `Natural Toughness`.

The underlying evidence bundle itself is complete and technically healthy for all 59 action-card images, and all five decks reconcile quantitatively. `pandora`, `muhammad-ali`, and `leonardo` are clean of material card discrepancies, with only documented qualifications or evidence-package metadata issues.

Required corpus follow-up is therefore narrow and localized: reconcile the two semantic findings above without treating the P3 archive slug/name differences as canonical Git defects.

Batch verdict: **FAIL**.
