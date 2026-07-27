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
| Annie Christmas | `docs/fighters/phase-4b/annie-christmas.yaml` | `docs/cards/phase-4b/annie-christmas.yaml` | 30 | verified | health-relative static bonus, defender replacement, damage-preserving health floor |
| Dr. Jill Trent | `docs/fighters/phase-4b/dr-jill-trent.yaml` | `docs/cards/phase-4b/dr-jill-trent.yaml` | 30 | verified | public gadget enum state, printed-value comparisons |
| Golden Bat | `docs/fighters/phase-4b/golden-bat.yaml` | `docs/cards/phase-4b/golden-bat.yaml` | 30 | verified | maneuver history, turn-start space snapshot |
| Nikola Tesla | `docs/fighters/phase-4b/nikola-tesla.yaml` | `docs/cards/phase-4b/nikola-tesla.yaml` | 30 | verified | charged-coil resource, permissive declared discharge |
| Oda Nobunaga | `docs/fighters/phase-4b/oda-nobunaga.yaml` | `docs/cards/phase-4b/oda-nobunaga.yaml` | 30 | verified | independent 6-health sidekicks, flanking, team-friendly semantics, per-instance history |
| Tomoe Gozen | `docs/fighters/phase-4b/tomoe-gozen.yaml` | `docs/cards/phase-4b/tomoe-gozen.yaml` | 30 | verified | relocation transition events, declaration-time attack override |
| William Shakespeare | `docs/fighters/phase-4b/shakespeare.yaml` | `docs/cards/phase-4b/shakespeare.yaml` | 30 | verified | persistent ordered Line, full completion effects, cleanup timing |
| Hamlet | `docs/fighters/phase-4b/hamlet.yaml` | `docs/cards/phase-4b/hamlet.yaml` | 30 | verified | public question enum state, conditional self-damage choice, damage history |
| Titania | `docs/fighters/phase-4b/titania.yaml` | `docs/cards/phase-4b/titania.yaml` | 30 + 6 Glamours | verified | shuffled external card pool, active-source lifetime, pre-defense interrupt |
| Ciri | `docs/fighters/phase-4b/ciri.yaml` | `docs/cards/phase-4b/ciri.yaml` | 30 | verified | discard-derived Source count, deck-wide cancellation protection, search/shuffle |
| Ancient Leshen | `docs/fighters/phase-4b/ancient-leshen.yaml` | `docs/cards/phase-4b/ancient-leshen.yaml` | 30 | verified | per-fighter movement/history, summon/return Wolves, temporary non-defeat off-board state |
| Eredin | `docs/fighters/phase-4b/eredin.yaml` | `docs/cards/phase-4b/eredin.yaml` | 30 | verified | derived Enraged state, Rider costs, team-friendly defender replacement, declaration/type overrides |
| Philippa | `docs/fighters/phase-4b/philippa.yaml` | `docs/cards/phase-4b/philippa.yaml` | 30 | verified | ordered hidden-information choices, exact health assignment, combat-card replacement, ongoing movement override |

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
- **Authority:** official `Slings and Arrows` set rules / current Shakespeare rulings.
- **Gap:** current card-zone semantics do not define ordered Cleanup append, syllable metrics, exact/overflow completion, or a completion-effect channel distinct from ordinary card effects.
- **Recommended model:** ordered sequence zone + Cleanup append + definition-level metric + exact/overflow threshold completion.
- **Completion record:** full normalized effect record, including controller choices, conditions, dependencies and operations.
- **Timing:** ordinary After Combat -> source-specific pre-Line hooks such as `Again` -> Cleanup append -> threshold evaluation.
- **Cleanup invariant:** after a completion effect resolves, discard only card instances still present in the Line. A card moved out by its completion effect, e.g. `Deceive`, must not be discarded again from a stale snapshot.
- **Integration requirement:** required.

### C-EXT-002 — auxiliary shuffled gameplay-card deck

- **Affected:** Titania Glamours.
- **Authority:** official `Slings and Arrows` set rules.
- **Gap:** six runtime card instances need hidden shuffled order, one active face-up card, public discard, explicit returns, no automatic reshuffle, and effects active only while that runtime instance is face-up.
- **Recommended model:** generic auxiliary gameplay-card deck with instance identity, ownership, visibility, zones and source-lifetime activation.
- **Integration requirement:** required.

### C-EXT-003 — pre-commit combat play-mode and participant override

- **Affected:** Tomoe `Witness My Last Battle`; Eredin `Foul Purpose` / `Implacable` / `Portal Defense`; Titania `Glamour of Jealousy`; Oda `Spring the Trap`.
- **Authority:** published components plus official set rules/rulings.
- **Gap:** these mechanics alter legal combat participation, range/type, visibility or defender identity outside the ordinary combat-card effect pipeline. In team play, replacement defender fighter/controller can differ from the controller/owner of the defense card already committed to the combat.
- **Recommended model:** declaration-time `combat_play_rules` plus `BEFORE_COMBAT_CARD_COMMIT` / same-combat participant replacement. Combat state must track participant fighter identity/controller separately from committed combat-card instance owner/controller.
- **Integration requirement:** required.

### C-EXT-004 — permissive declared resource discharge

- **Affected:** Nikola Tesla coil-discharge effects, including fixed `discharge both` branches such as `The Alternating Current`.
- **Authority:** official `Tales to Amaze` set rules; official Rulings Archive verdict mirrored at `https://www.the-unmatched.club/tools/disputes/a1c7dc96-d154-4253-80f1-d86015292f9e`; Unmatched Reference v10; published deck facts from UmDb.
- **Verified:** declared two-coil discharge is legal with `0` or `1` charged coils; it is not downgraded to a one-coil branch/tier; an underfunded dependent effect may intentionally fizzle.
- **Recommended integration behavior:** consume available charged coils up to the declared amount. Therefore `1 charged + declare 2 -> 0 charged`; preserve declared branch/tier; no lower-tier fallback; source-defined dependent effect may fizzle.
- **Evidence qualifier:** confidence **high** from composed official semantics. No exact-case publisher sentence literally stating `1 charged + declare 2 -> 0 charged` was found.
- **Recommended model:** `declared_resource_discharge` plus tiered specialization; declaration legality is independent of current resource sufficiency, partial resource consumption is explicit, dependent resolution is source-defined.
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
- **Gap:** existing lookup/reveal/move operations do not define whole-zone search by predicate and hidden-information/post-search semantics.
- **Recommended model:** `SEARCH(zone, filter, viewer, reveal_selected, destination, post_search)`; searched deck applies authoritative post-search shuffle.
- **Searching Strike:** search own deck for a Source card, reveal it, move it to hand, shuffle the searched deck.
- **Integration requirement:** required.

### C-EXT-007 — temporary non-defeat battlefield removal with scheduled return

- **Affected:** Ancient Leshen `Vanish Into Murder`.
- **Authority:** published component content plus official `The Witcher — Steel & Silver` rules/rulings.
- **Gap:** `DEFEAT` corrupts defeat/loss semantics; `RETURN_FIGHTER` models only the return half.
- **Recommended model:** `REMOVE_FROM_BATTLEFIELD(target, defeat:false, preserve_health:true)` plus persisted scheduled-return context consumed by `RETURN_FIGHTER`. The scheduled context is authoritative; no duplicate fighter-local boolean mirror.
- **Integration requirement:** required.
- **Targeted evidence check:** Ancient Leshen × dormant-player ordering below.

### C-EXT-008 — exact health assignment

- **Affected:** Philippa `Backup Plan`.
- **Authority:** published component content / current Philippa ruling index.
- **Gap:** `RECOVER` cannot lower health; this effect sets health to exactly `5`.
- **Recommended model:** `SET_HEALTH(target, value, source_limits)` distinct from damage and recovery.
- **Integration requirement:** required.

### C-EXT-009 — withdrawn; existing relocation transition semantics are sufficient

- **Affected:** Tomoe Gozen special ability.
- **Reason withdrawn:** `docs/mechanics/movement-and-placement.md` already requires `left_space`, `entered_space`, `left_zone`, `entered_zone`, and `removed_from_board` transition facts and explicitly covers Tomoe-class interactions.
- **Integration action:** consume existing `left_zone` / removal facts and filter to the actual opposing hero transition.
- **Integration requirement:** no new primitive.

### C-EXT-010 — runtime-fighter-keyed historical state

- **Affected:** Oda `Momentous Shift`; Ancient Leshen `Planted Feet`; future repeated fighter definitions/instances.
- **Gap:** scalar `space_ref` / `boolean` cannot represent independent historical values for multiple runtime fighter instances.
- **Recommended model:** persistent maps keyed by runtime fighter instance, e.g. `fighter_instance_space_map`, `fighter_instance_boolean_map`, with source-defined reset boundaries.
- **Integration requirement:** required.

### C-EXT-011 — fighter-specific base attributes and continuous source-lifetime modifiers

- **Affected:** Ancient Leshen/Wolf move values; Eredin Enraged move 3; Philippa `Polymorphy` move 5; Titania `Glamour of Invisibility` movement permission.
- **Authority:** official set rules / published ongoing sources.
- **Gap:** permanent `SET_STATE` is wrong for a modifier that disappears when its condition/source stops applying, and one character-level move value cannot represent Ancient Leshen 1 / Wolves 3.
- **Recommended model:** definition-level base attributes plus derived/effective attribute and rule-permission modifiers tied to condition/source lifetime. Source removal removes the modifier without rollback mutation.
- **Integration requirement:** required.

### C-EXT-012 — damage-preserving health floor

- **Affected:** Annie Christmas `Mississippi Queen`.
- **Authority:** published card text plus official Rulings Archive entry mirrored at `https://www.the-unmatched.club/tools/disputes/b8f0f25e-e75c-4edd-a6dd-17706ddc6d01`.
- **Gap:** `PREVENT_DAMAGE` changes whether damage was dealt. Mississippi Queen does not prevent the damage event; it prevents Annie's resulting health from being reduced below 1.
- **Recommended model:** damage application can carry a source-lifetime health floor: compute/record damage as dealt normally, then clamp resulting health to `minimum_health`. Damage-dealt observers continue to see the original dealt amount.
- **Integration requirement:** required.

### C-EXT-013 — ordered dependent effect stages and empty choice domains

- **Affected:** Philippa `Spymaster's Ruse` / `Do My Bidding`; Titania `What Fools These Mortals Be`; Annie `Bottom Dealing`; Jill `Insightful Deduction`; Golden Bat `Sight Beyond Sight`; Ciri `Bane of the Aen Elle`; Hamlet `Method in the Madness`; Shakespeare return/completion choices; similar effects.
- **Authority:** `FX-011`, `FX-030` plus published component ordering.
- **Gap:** a flat `choices -> operations` record cannot faithfully express `operation -> inspect/capture result -> choice using that information -> dependent operation`. Existing manifests use several ad-hoc spellings (`choices_after`, `followup`, `nested_choice`).
- **Recommended model:** generic ordered effect stages. Each stage may contain choices/costs/operations and may read captures from completed earlier stages.
- **Empty-domain invariant:** if a required dependent choice has no legal options, create no pending choice; skip the impossible dependent branch under ordinary partial-resolution rules and continue independently resolvable later stages.
- **Integration requirement:** required.

## Existing global semantics relied upon

- Oda's independent 6-health Honor Guards use ordinary independent runtime fighter identities/health.
- `your fighters` means fighters controlled by the effect/card controller; `friendly fighters` can include a teammate in team play. Manifests preserve this distinction.
- Jill `active_gadget` and Hamlet `question_state` are each a single authoritative public enum state; their physical tokens are UI representations, not duplicate mutable resources.
- Golden Bat / Hamlet turn-history checks use explicit historical state under `FX-021`.
- `cannot leave current space` uses existing `PLACE-040` and blocks ordinary maneuver movement, MOVE and PLACE but not defeat removal.
- Leshen Wolf summon/return uses existing `SUMMON` / `RETURN_FIGHTER`; only `Vanish Into Murder` needs C-EXT-007.
- Witcher ongoing schemes reuse the Phase 4A ongoing-source lifecycle. `max_active: 1` is replacement capacity: a new ongoing scheme discards the current one; it is not a play-legality restriction. Discard conditions are checked at controller turn end only.
- Philippa `Do My Bidding` uses the Phase 4A Dracula `REPLACE_COMBAT_CARD` pattern and resumes the same combat.
- Tomoe zone-leave detection uses existing relocation transition facts; C-EXT-009 is withdrawn.

## Required integration fixtures

These are deterministic cross-mechanic fixtures, not unresolved rule questions.

- **Shakespeare × `END_TURN`:** end-turn request -> Cleanup -> current combat card enters Line -> exact-10 completion resolves if reached -> gained actions cannot be spent because turn control is already ending.
- **Shakespeare `Deceive` completion:** completion moves its instance to hand -> subsequent Line cleanup discards only instances still in Line.
- **Witcher ongoing replacement:** second ongoing scheme replaces/discards current active scheme and becomes sole active scheme.
- **Continuous modifiers:** returning an Eredin Rider removes Enraged move 3; discarding `Polymorphy` removes move 5; replacing/discarding `Glamour of Invisibility` removes its movement permission.
- **Team defender replacement:** Oda `Spring the Trap` and Eredin `Portal Defense` can select a teammate's friendly fighter; replacement fighter/controller and already-played defense-card owner/controller remain distinct.
- **Annie health floor:** damage that would pass below 1 still counts as dealt while resulting health is clamped at 1; damage-dealt observers must see the dealt amount.
- **Ordered choices:** `Spymaster's Ruse` resolves opponent selection -> opponent-selected reveal set -> controller-selected discard set -> draw; an empty dependent domain must not create an unresolvable pending choice.

## Targeted evidence / integration follow-up

- **Ancient Leshen × dormant player:** if both Wolves are defeated and `Vanish Into Murder` removes Leshen, the player becomes dormant at the end-of-action checkpoint. The card schedules return + draw at the next turn start, while global dormant rules restrict dormant draws. Current shared docs explicitly defer character-specific dormant/start-turn interactions; establish exact ordering before encoding an exception or suppressing the draw.
- **Hamlet `The Readiness Is All` × failed placement:** global placement semantics permit a selected occupied destination to fail. Exact dependency of the subsequent `cannot leave that space this turn` restriction on placement success is not yet source-resolved.
- **Eredin `Icy Guile` self-defeat:** printed wording permits defeating `a Red Rider` and does not say `another`; whether the currently attacking Red Rider may select itself remains an exact-ruling follow-up. Do not add an inferred exclusion.

## Evidence qualifier

- **Tesla exact `1 -> 0` wording:** no exact-case publisher sentence was found; C-EXT-004 retains the high-confidence composed-official recommendation as a non-blocking qualifier.

## Source notes

- Tales to Amaze competitive card metadata is complete in published UmDb; no separate current publisher-hosted per-hero card-text PDF was located in this pass.
- Ciri `Searching Strike` post-search shuffle is covered by the current general search/check-deck ruling.
- No fan `/decks/...` balance-patch data was imported.

## Worker 4B-C Handoff

Branch: `phase-4b-worker-c-modern`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: exact branch tip supplied in the external handoff; a persisted report cannot self-contain the SHA of the commit that contains itself  
Assigned fighters: 13  
Verified: 13  
Blocked: none  
Quantity validation: **PASS — 390/390 action cards; 13/13 decks**  
Integration requirements: C-EXT-001 through C-EXT-008, C-EXT-010 through C-EXT-013  
Withdrawn extension: C-EXT-009 — existing relocation transition semantics are sufficient  
Required fixtures: Shakespeare `END_TURN`/Line cleanup; Witcher ongoing replacement; continuous modifier removal; team defender replacement; Annie health floor; ordered dependent choices  
Targeted follow-up: Ancient Leshen dormant ordering; Hamlet failed placement dependency; Eredin Icy Guile self-defeat  
Evidence qualifier: Tesla C-EXT-004 exact-case `1 -> 0` wording not found; recommended behavior remains high-confidence composed official semantics  
Files created: 27 — 13 fighter manifests, 13 deck manifests, this report  
Shared semantic files changed: **none**  
Phase 4A manifests rewritten: **none**  
Merged to `main`: **no**
