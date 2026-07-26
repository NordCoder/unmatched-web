# Phase 4B Worker C report

**Branch:** `phase-4b-worker-c-modern`  
**Authorized base:** `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
**Scope:** modern mechanic-heavy competitive fighters only  
**Assigned fighters:** 13

The branch was verified identical to the Authorized Base before transcription. Worker C did not edit shared schema, mechanics, rules, set registry, ambiguity-register, README, or Phase 4A manifests.

For `tales-to-amaze`, only competitive player-hero behavior is in scope. Cooperative enemies, minions, initiative cards, threat/scenario state, and other Adventures enemy behavior are excluded.

## Corpus status

| Fighter | Fighter manifest | Deck manifest | Action cards | Status | Primary stress dimensions |
| --- | --- | --- | ---: | --- | --- |
| Annie Christmas | `docs/fighters/phase-4b/annie-christmas.yaml` | `docs/cards/phase-4b/annie-christmas.yaml` | 30 | verified | health-relative static bonus, defender replacement |
| Dr. Jill Trent | `docs/fighters/phase-4b/dr-jill-trent.yaml` | `docs/cards/phase-4b/dr-jill-trent.yaml` | 30 | verified | public gadget state machine, printed-value comparisons |
| Golden Bat | `docs/fighters/phase-4b/golden-bat.yaml` | `docs/cards/phase-4b/golden-bat.yaml` | 30 | verified | maneuver history, turn-start space snapshot |
| Nikola Tesla | `docs/fighters/phase-4b/nikola-tesla.yaml` | `docs/cards/phase-4b/nikola-tesla.yaml` | 30 | verified | charged-coil resource, permissive declared discharge, C-EXT-004 |
| Oda Nobunaga | `docs/fighters/phase-4b/oda-nobunaga.yaml` | `docs/cards/phase-4b/oda-nobunaga.yaml` | 30 | verified | independent 6-health sidekicks, flanking, per-instance history |
| Tomoe Gozen | `docs/fighters/phase-4b/tomoe-gozen.yaml` | `docs/cards/phase-4b/tomoe-gozen.yaml` | 30 | verified | relocation transition events, declaration-time attack override |
| William Shakespeare | `docs/fighters/phase-4b/shakespeare.yaml` | `docs/cards/phase-4b/shakespeare.yaml` | 30 | verified | persistent ordered Line, full completion effects, cleanup timing |
| Hamlet | `docs/fighters/phase-4b/hamlet.yaml` | `docs/cards/phase-4b/hamlet.yaml` | 30 | verified | turn-start state choice, conditional self-damage choice, damage history |
| Titania | `docs/fighters/phase-4b/titania.yaml` | `docs/cards/phase-4b/titania.yaml` | 30 + 6 Glamours | verified | shuffled external card pool, active-source lifetime, pre-defense interrupt |
| Ciri | `docs/fighters/phase-4b/ciri.yaml` | `docs/cards/phase-4b/ciri.yaml` | 30 | verified | discard-derived Source count, deck-wide cancellation protection, search/shuffle |
| Ancient Leshen | `docs/fighters/phase-4b/ancient-leshen.yaml` | `docs/cards/phase-4b/ancient-leshen.yaml` | 30 | verified | per-fighter movement, summon/return Wolves, temporary non-defeat off-board state |
| Eredin | `docs/fighters/phase-4b/eredin.yaml` | `docs/cards/phase-4b/eredin.yaml` | 30 | verified | derived Enraged state, explicit Rider costs, declaration/type overrides |
| Philippa | `docs/fighters/phase-4b/philippa.yaml` | `docs/cards/phase-4b/philippa.yaml` | 30 | verified | exact health assignment, combat-card replacement, continuous movement override |

## Quantity validation

**PASS — 390/390 action cards; 13/13 deck quantities reconcile.**

- Titania's six Glamours are external gameplay cards and are not counted in her 30-card action deck.
- Ciri's 13 Source-tagged copies remain part of her ordinary 30-card deck.
- No Tales to Amaze enemy, initiative, or scenario cards are included.
- Every explicit `usable_by` target exists in the corresponding fighter topology.

## Proposed reusable semantic extensions

These are integration requirements, not fighter research blockers. Source facts are verified unless an evidence qualifier is stated. Worker C did not modify shared semantic files.

### C-EXT-001 — ordered persistent card sequence with full completion-effect records

- **Affected:** Shakespeare Line mechanics.
- **Authority:** official `Slings and Arrows` set rules.
- **Gap:** current `card_zone` / `MOVE_CARD` semantics do not define ordered cleanup append, sequence metrics, exact/overflow thresholds, or a completion-effect channel distinct from ordinary card effects.
- **Recommended model:** ordered sequence zone + append-at-cleanup semantics + definition-level sequence metric + exact/overflow threshold completion.
- **Completion record:** a completion effect is a full normalized effect record and may contain controller choices, conditions, dependencies and operations; it is not a bare operation list.
- **Timing:** ordinary After Combat processing completes; source-specific pre-Line hooks such as `Again` resolve; Cleanup appends the played combat card; then the Line threshold is evaluated.
- **Integration requirement:** required.

### C-EXT-002 — auxiliary shuffled gameplay-card deck

- **Affected:** Titania Glamours.
- **Authority:** official `Slings and Arrows` set rules.
- **Gap:** Titania needs six runtime card instances with hidden shuffled order, one active face-up card, public discard, explicit returns, no automatic reshuffle, and effects active only while that runtime instance is face-up.
- **Recommended model:** generic auxiliary gameplay-card deck with normal instance identity, ownership, zones, visibility and source-lifetime activation.
- **Integration requirement:** required.

### C-EXT-003 — pre-commit combat play-mode and participant override

- **Affected:** Tomoe `Witness My Last Battle`; Eredin `Foul Purpose` / `Implacable`; Titania `Glamour of Jealousy`.
- **Authority:** published components plus official set rules/rulings.
- **Gap:** these mechanics alter legal combat participation, range/type, visibility or defender identity before the ordinary combat-card commit/reveal pipeline.
- **Recommended model:** declaration-time `combat_play_rules` plus a `BEFORE_COMBAT_CARD_COMMIT` interrupt capable of replacing a participant and continuing the same combat declaration.
- **Integration requirement:** required.

### C-EXT-004 — permissive declared resource discharge

- **Affected:** Nikola Tesla coil-discharge effects, including tiered effects and fixed `discharge both` branches such as `The Alternating Current`.
- **Authority:** official `Tales to Amaze` set rules; official Rulings Archive verdict mirrored at `https://www.the-unmatched.club/tools/disputes/a1c7dc96-d154-4253-80f1-d86015292f9e`; Unmatched Reference v10; published deck facts from UmDb.
- **Verified:** a declared two-coil discharge remains a legal declaration with `0` or `1` charged coils; it is not downgraded to a one-coil branch/tier; an underfunded dependent effect may intentionally fizzle.
- **Recommended integration behavior:** consume available charged coils up to the declared amount. Therefore `1 charged + declare 2 -> 0 charged`; preserve the declared branch/tier; do not resolve a lower tier as fallback; if the dependent effect requires full declared discharge, it may fizzle.
- **Evidence qualifier:** confidence **high** from composed official semantics. Official rules define an individual coil discharge as `charged -> discharged`; the official ruling independently permits underfunded two-coil declaration/fizzle. No exact-case publisher sentence literally stating `1 charged + declare 2 -> 0 charged` was found.
- **Recommended model:** `declared_resource_discharge` plus a tiered specialization; declaration legality can be independent of current resource sufficiency, resource consumption is partial up to the declared amount, and dependent resolution is source-defined.
- **Integration requirement:** required. **Research status:** verified.

### C-EXT-005 — filtered zone aggregate and highest-met threshold

- **Affected:** Ciri Source cards and cancellation protection.
- **Authority:** official `The Witcher — Steel & Silver` set rules.
- **Gap:** Source effects derive a discard-zone count and resolve only the highest threshold currently met.
- **Recommended model:** `COUNT_CARDS(zone, filter)` plus `RESOLVE_HIGHEST_MET_THRESHOLD(metric, ordered_thresholds)`.
- **Integration requirement:** required.

### C-EXT-006 — search card zone by definition predicate

- **Affected:** Ciri `Searching Strike` and future published search effects.
- **Authority:** published Ciri component content plus current general ruling that checking/searching a deck is followed by shuffling it.
- **Gap:** existing lookup/reveal/move operations do not define a legal whole-zone search by predicate and its hidden-information/post-search disposition semantics.
- **Recommended model:** `SEARCH(zone, filter, viewer, reveal_selected, destination, post_search)`; when a deck is searched, apply the authoritative post-search shuffle rule.
- **Searching Strike:** search own deck for a Source card, reveal the selected card, move it to hand, then shuffle the searched deck.
- **Integration requirement:** required.

### C-EXT-007 — temporary non-defeat battlefield removal with scheduled return

- **Affected:** Ancient Leshen `Vanish Into Murder`.
- **Authority:** published component content plus official `The Witcher — Steel & Silver` rules/rulings.
- **Gap:** `DEFEAT` would incorrectly trigger defeat/loss semantics; `RETURN_FIGHTER` models only the return half.
- **Recommended model:** `REMOVE_FROM_BATTLEFIELD(target, defeat:false, preserve_health:true)` plus persisted scheduled return context consumed by `RETURN_FIGHTER`.
- **Integration requirement:** required.
- **Targeted integration evidence check:** see Ancient Leshen × dormant-player ordering below; do not infer a new dormant exception without authority.

### C-EXT-008 — exact health assignment

- **Affected:** Philippa `Backup Plan`.
- **Authority:** published component content / current Philippa ruling index.
- **Gap:** `RECOVER` cannot lower health; this effect sets health to exactly `5`.
- **Recommended model:** `SET_HEALTH(target, value, source_limits)` distinct from damage and recovery.
- **Integration requirement:** required.

### C-EXT-009 — withdrawn; existing relocation transition semantics are sufficient

- **Affected:** Tomoe Gozen special ability.
- **Reason withdrawn:** `docs/mechanics/movement-and-placement.md` already requires relocation operations to emit `left_space`, `entered_space`, `left_zone`, `entered_zone`, and `removed_from_board` facts and explicitly cites Tomoe-class interactions. A second parallel `FIGHTER_LEFT_ZONE` event primitive would duplicate existing semantics.
- **Integration action:** consume existing `left_zone` / removal transition facts and filter to the actual opposing hero transition. Tomoe moving herself does not emit a `left_zone` event for the opposing hero.
- **Integration requirement:** no new primitive.

### C-EXT-010 — runtime-fighter-keyed historical state

- **Affected:** Oda `Momentous Shift`; Ancient Leshen `Planted Feet`; future characters with repeated fighter definitions/instances.
- **Gap:** schema scalar `space_ref` / `boolean` cannot represent independent historical values for Oda plus two Honor Guards, or Leshen plus two Wolves.
- **Recommended model:** persistent maps keyed by runtime fighter instance, e.g. `fighter_instance_space_map` and `fighter_instance_boolean_map`, with source-defined reset boundaries.
- **Integration requirement:** required.

### C-EXT-011 — fighter-specific base attributes and continuous source-lifetime modifiers

- **Affected:** Ancient Leshen/Wolf move values; Eredin Enraged move 3; Philippa `Polymorphy` move 5; Titania `Glamour of Invisibility` movement permission.
- **Authority:** official set rules / published ongoing sources.
- **Gap:** permanent-looking `SET_STATE` is incorrect for a modifier that disappears when its condition/source stops applying, and a single character-level move value cannot represent Ancient Leshen 1 / Wolves 3.
- **Recommended model:** definition-level fighter base attributes plus derived/effective attribute and rule-permission modifiers whose lifetime is tied to the active condition/source. Removing the source automatically removes the modifier; no rollback mutation is required.
- **Integration requirement:** required.

## Mechanics covered by existing Phase 4A/global semantics

- Oda's independent 6-health Honor Guards use ordinary independent runtime fighter identities/health.
- Jill's public gadget state uses ordinary public enum state.
- Golden Bat / Hamlet turn-history checks use explicit historical state under `FX-021`.
- `cannot leave current space` restrictions use existing `PLACE-040` and therefore block ordinary maneuver movement, MOVE and PLACE but not defeat removal.
- Leshen Wolf summon/return uses existing `SUMMON` / `RETURN_FIGHTER`; only `Vanish Into Murder` needs C-EXT-007.
- Witcher ongoing schemes reuse the Phase 4A ongoing-source lifecycle. `max_active: 1` is replacement capacity: playing a second ongoing scheme discards the current active scheme; it is not a play-legality restriction. Source discard conditions are checked at the end of the controller's turn only.
- Philippa `Do My Bidding` uses existing `REPLACE_COMBAT_CARD` and resumes the same combat.
- Tomoe zone-leave detection uses existing relocation transition facts; C-EXT-009 is withdrawn.

## Required integration fixtures

These are deterministic cross-mechanic fixtures, not unresolved rule questions.

- **Shakespeare × `END_TURN`:** existing `ENDTURN-003` explicitly states Cleanup still occurs and identifies Shakespeare Line processing as a canonical cleanup behavior. Fixture should assert: end-turn request -> Cleanup -> current combat card enters Line -> exact-10 completion resolves if reached -> gained actions cannot subsequently be spent because turn control is already ending.
- **Witcher ongoing replacement:** when a second ongoing scheme is played, current active scheme goes to discard before/while the new one becomes the sole active scheme, according to the shared Witcher lifecycle.
- **Continuous modifiers:** returning an Eredin Rider removes Enraged move 3 immediately; discarding `Polymorphy` removes move 5 immediately; discarding/replacing `Glamour of Invisibility` removes its movement permission immediately.

## Targeted evidence / integration follow-up

- **Ancient Leshen × dormant player:** if both Wolves are defeated and `Vanish Into Murder` removes the Leshen, the player becomes dormant at the end-of-action checkpoint. The card schedules `place the Leshen ... then draw 1` at the start of the next turn, while the global dormant rule says dormant players cannot draw. Current shared docs explicitly defer character-specific dormant/start-turn lifecycle interactions. Integration must establish the exact ordering/source authority before encoding a special dormant exception or suppressing the card's draw.
- **Tesla exact `1 -> 0` wording:** no exact-case publisher sentence was found; C-EXT-004 retains the high-confidence composed-official recommendation as a non-blocking evidence qualifier.

## Source notes

- Tales to Amaze competitive card metadata is complete in published UmDb; no separate current publisher-hosted per-hero card-text PDF was located in this pass.
- Ciri `Searching Strike` post-search shuffle is now covered by the current general search/check-deck ruling and is no longer an open source gap.
- No fan `/decks/...` balance-patch data was imported.

## Worker 4B-C Handoff

Branch: `phase-4b-worker-c-modern`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: exact branch tip supplied in the external handoff; a persisted report cannot self-contain the SHA of the commit that contains itself  
Assigned fighters: 13  
Verified: 13  
Blocked: none  
Quantity validation: **PASS — 390/390 action cards; 13/13 decks**  
Integration requirements: C-EXT-001 through C-EXT-008, C-EXT-010 and C-EXT-011  
Withdrawn extension: C-EXT-009 — existing relocation transition semantics are sufficient  
Required fixtures: Shakespeare × `END_TURN`; Witcher ongoing replacement; continuous modifier removal  
Targeted follow-up: Ancient Leshen × dormant start-turn return/draw ordering  
Evidence qualifier: Tesla C-EXT-004 exact-case `1 -> 0` wording not found; recommended behavior remains high-confidence composed official semantics  
Files created: 27 — 13 fighter manifests, 13 deck manifests, this report  
Shared semantic files changed: **none**  
Phase 4A manifests rewritten: **none**  
Merged to `main`: **no**
