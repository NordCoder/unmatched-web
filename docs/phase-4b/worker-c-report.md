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
| Nikola Tesla | `docs/fighters/phase-4b/nikola-tesla.yaml` | `docs/cards/phase-4b/nikola-tesla.yaml` | 30 | verified | charged-coil resource, permissive tier declaration, C-EXT-004 |
| Oda Nobunaga | `docs/fighters/phase-4b/oda-nobunaga.yaml` | `docs/cards/phase-4b/oda-nobunaga.yaml` | 30 | verified | two independent 6-health sidekicks, flanking |
| Tomoe Gozen | `docs/fighters/phase-4b/tomoe-gozen.yaml` | `docs/cards/phase-4b/tomoe-gozen.yaml` | 30 | verified | zone-leave event semantics, declaration-time attack override |
| William Shakespeare | `docs/fighters/phase-4b/shakespeare.yaml` | `docs/cards/phase-4b/shakespeare.yaml` | 30 | verified | persistent ordered Line, syllable threshold, completion effects |
| Hamlet | `docs/fighters/phase-4b/hamlet.yaml` | `docs/cards/phase-4b/hamlet.yaml` | 30 | verified | turn-start state choice, damage history |
| Titania | `docs/fighters/phase-4b/titania.yaml` | `docs/cards/phase-4b/titania.yaml` | 30 + 6 Glamours | verified | shuffled external card pool, active/discard lifecycle, pre-defense interrupt |
| Ciri | `docs/fighters/phase-4b/ciri.yaml` | `docs/cards/phase-4b/ciri.yaml` | 30 | verified | discard-derived Source count, highest threshold only, ongoing scheme |
| Ancient Leshen | `docs/fighters/phase-4b/ancient-leshen.yaml` | `docs/cards/phase-4b/ancient-leshen.yaml` | 30 | verified | summon/return Wolves, temporary non-defeat off-board state, ongoing scheme |
| Eredin | `docs/fighters/phase-4b/eredin.yaml` | `docs/cards/phase-4b/eredin.yaml` | 30 | verified | derived Enraged state, returnable Riders, declaration/type overrides, ongoing scheme |
| Philippa | `docs/fighters/phase-4b/philippa.yaml` | `docs/cards/phase-4b/philippa.yaml` | 30 | verified | action history, exact health assignment, combat-card replacement, ongoing scheme |

## Quantity validation

**PASS — 390/390 action cards; 13/13 deck quantities reconcile.**

- Titania's six Glamours are external gameplay cards and are not counted in her 30-card action deck.
- Ciri's 13 Source-tagged copies remain part of her ordinary 30-card deck.
- No Tales to Amaze enemy, initiative, or scenario cards are included.
- Every explicit `usable_by` target exists in the corresponding fighter topology.

## Proposed reusable semantic extensions

These are integration requirements, not fighter research blockers. Source facts are verified unless an evidence qualifier is stated. Worker C did not modify shared semantic files.

### C-EXT-001 — ordered persistent card sequence with threshold completion

- **Affected:** Shakespeare Line mechanics.
- **Authority:** official `Slings and Arrows` set rules: `https://restorationgames.com/wp-content/uploads/2024/03/UM-SaA-Set-Rules.pdf`.
- **Gap:** current `card_zone` / `MOVE_CARD` semantics do not define ordered cleanup append, sequence metrics, exact/overflow thresholds, or completion effects distinct from ordinary card effects.
- **Recommended model:** ordered card zone + append-at-cleanup semantics + definition-level sequence metric + exact/overflow threshold completion.
- **Integration requirement:** required.

### C-EXT-002 — auxiliary shuffled gameplay-card deck

- **Affected:** Titania Glamours.
- **Authority:** official `Slings and Arrows` set rules.
- **Gap:** Titania needs six runtime card instances with hidden shuffled order, one active public card, public discard, explicit returns, and no automatic reshuffle.
- **Recommended model:** generic auxiliary gameplay-card deck with normal instance identity, ownership, zones, visibility, and source-defined reshuffle policy.
- **Integration requirement:** required.

### C-EXT-003 — pre-commit combat play-mode and participant override

- **Affected:** Tomoe `Witness My Last Battle`; Eredin `Foul Purpose` / `Implacable`; Titania `Glamour of Jealousy`.
- **Authority:** published components plus official set rules/rulings.
- **Gap:** these mechanics alter legal combat participation or play mode before the ordinary combat-card commit/reveal pipeline.
- **Recommended model:** declaration-time `combat_play_rules` plus a `BEFORE_COMBAT_CARD_COMMIT` interrupt capable of replacing a participant and continuing the same combat declaration.
- **Integration requirement:** required.

### C-EXT-004 — permissive tiered resource discharge

- **Affected:** Nikola Tesla coil-discharge effects.
- **Authority:** official `Tales to Amaze` set rules (`https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Set-rules_EN_Light.pdf`); official Rulings Archive verdict mirrored at `https://www.the-unmatched.club/tools/disputes/a1c7dc96-d154-4253-80f1-d86015292f9e`; Unmatched Reference v10; published deck facts from UmDb.
- **Verified:** tier `2` may be declared with `0` or `1` charged coils; selected tier remains `2`; no tier-1 fallback occurs; the tier-2 dependent effect may intentionally fizzle.
- **Recommended integration behavior:** discharge available charged coils up to the declared amount. Therefore `1 charged + declare 2 -> 0 charged`; selected tier remains `2`; tier-2 dependent effect may fizzle.
- **Evidence qualifier:** confidence **high** from composed official semantics. Official rules define an available coil discharge as `charged -> discharged`, while the official ruling independently permits the underfunded tier-2 declaration/fizzle. No exact-case publisher ruling literally stating `1 charged + declare 2 -> 0 charged` was found.
- **Recommended model:** generic declared resource tier whose legality is independent of current resource sufficiency; consume available resource up to the declared amount; preserve declared tier identity; permit source-defined fizzle without lower-tier fallback.
- **Integration requirement:** required. **Research status:** verified.

### C-EXT-005 — filtered zone aggregate and highest-met threshold

- **Affected:** Ciri Source cards and cancellation protection.
- **Authority:** official `The Witcher — Steel & Silver` set rules: `https://restorationgames.com/wp-content/uploads/2024/10/UM-W-SaS-Set-Rules.pdf`.
- **Gap:** Source effects derive a discard-zone count and resolve only the highest threshold currently met.
- **Recommended model:** `COUNT_CARDS(zone, filter)` plus `RESOLVE_HIGHEST_MET_THRESHOLD(metric, ordered_thresholds)`.
- **Integration requirement:** required.

### C-EXT-006 — search card zone by definition predicate

- **Affected:** Ciri `Searching Strike` and future published search effects.
- **Authority:** published Ciri component content / UmDb cross-check.
- **Gap:** existing lookup/reveal/move operations do not define a legal whole-zone search by predicate and disposition semantics.
- **Recommended model:** generic `SEARCH(zone, filter, viewer, reveal_selected, destination, post_search_order_rule)` with no implicit shuffle.
- **Integration requirement:** required.

### C-EXT-007 — temporary non-defeat battlefield removal with scheduled return

- **Affected:** Ancient Leshen `Vanish Into Murder`.
- **Authority:** published component content plus official `The Witcher — Steel & Silver` rules/rulings.
- **Gap:** `DEFEAT` would incorrectly trigger defeat/loss semantics; `RETURN_FIGHTER` models only the return half.
- **Recommended model:** `REMOVE_FROM_BATTLEFIELD(target, defeat:false, preserve_health:true)` plus persisted scheduled return context consumed by `RETURN_FIGHTER`.
- **Integration requirement:** required.

### C-EXT-008 — exact health assignment

- **Affected:** Philippa `Backup Plan`.
- **Authority:** published component content / current Philippa ruling index.
- **Gap:** `RECOVER` cannot lower health; this effect sets health to exactly `5`.
- **Recommended model:** `SET_HEALTH(target, value, source_limits)` distinct from damage and recovery.
- **Integration requirement:** required.

### C-EXT-009 — first-class fighter-left-zone event with cause context

- **Affected:** Tomoe Gozen special ability.
- **Authority:** official `Sun's Origin` set rules: `https://restorationgames.com/wp-content/uploads/2024/01/UM-SO-Set-Rules.pdf`.
- **Gap:** Tomoe triggers on actual opposing-hero zone exits caused by movement, placement, or removal, including multiple transitions, but not when Tomoe's own movement merely changes the zone relationship.
- **Recommended model:** `FIGHTER_LEFT_ZONE(fighter, zone, cause, source_operation)` emitted for each actual observed-fighter transition.
- **Integration requirement:** required.

## Mechanics covered by existing Phase 4A semantics

- Oda's independent 6-health Honor Guards.
- Jill's public gadget state.
- Golden Bat / Hamlet turn-history checks.
- Leshen Wolf summon/return; only `Vanish Into Murder` needs C-EXT-007.
- Derived Eredin Enraged and Ciri Source count.
- Witcher ongoing schemes.
- Philippa `Do My Bidding` via `REPLACE_COMBAT_CARD`.

## Remaining integration ambiguity

- **Shakespeare × `END_TURN`:** general rulings jump to the current action Cleanup Step, while Shakespeare's Line uses cleanup-time append/completion. Integration should include an explicit fixture for ordering. No contradictory source was found.

## Source notes

- Tales to Amaze competitive card metadata is complete in published UmDb; no separate current publisher-hosted per-hero card-text PDF was located in this pass.
- Tesla's exact `1 -> 0` wording is a non-blocking evidence qualifier only; C-EXT-004 records the recommended high-confidence integration behavior.
- No fan `/decks/...` balance-patch data was imported.

## Worker 4B-C Handoff

Branch: `phase-4b-worker-c-modern`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Assigned fighters: 13  
Verified: 13  
Blocked: none  
Quantity validation: **PASS — 390/390 action cards; 13/13 decks**  
Integration requirements: C-EXT-001 through C-EXT-009  
Remaining integration ambiguity: Shakespeare × `END_TURN` fixture only  
Evidence qualifier: Tesla C-EXT-004 exact-case `1 -> 0` wording not found; recommended behavior is high-confidence composed official semantics  
Files created: 27 — 13 fighter manifests, 13 deck manifests, this report  
Shared semantic files changed: **none**  
Phase 4A manifests rewritten: **none**  
Merged to `main`: **no**
