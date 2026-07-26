# Phase 4B Worker C report

**Branch:** `phase-4b-worker-c-modern`  
**Authorized base:** `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
**Corpus head immediately before this report commit:** `ea0aa34126ec5f14e15e2ea724c9fe497110eae7`  
**Scope:** modern mechanic-heavy competitive fighters only  
**Assigned fighters:** 13

The branch was verified identical to the Authorized Base before transcription. Worker C did not edit shared schema, mechanics, rules, set registry, ambiguity-register, README, or Phase 4A manifests.

For `tales-to-amaze`, only the four player hero decks and their competitive behavior are in scope. Mothman, Martian Invader, minions, initiative cards, threat/scenario state, and all other cooperative enemy behavior are intentionally absent.

## Corpus status

| Fighter | Fighter manifest | Deck manifest | Action cards | Status | Primary stress dimensions |
| --- | --- | --- | ---: | --- | --- |
| Annie Christmas | `docs/fighters/phase-4b/annie-christmas.yaml` | `docs/cards/phase-4b/annie-christmas.yaml` | 30 | verified | health-relative static bonus, defender replacement |
| Dr. Jill Trent | `docs/fighters/phase-4b/dr-jill-trent.yaml` | `docs/cards/phase-4b/dr-jill-trent.yaml` | 30 | verified | public gadget state machine, printed-value comparisons |
| Golden Bat | `docs/fighters/phase-4b/golden-bat.yaml` | `docs/cards/phase-4b/golden-bat.yaml` | 30 | verified | maneuver history, turn-start space snapshot |
| Nikola Tesla | `docs/fighters/phase-4b/nikola-tesla.yaml` | `docs/cards/phase-4b/nikola-tesla.yaml` | 30 | **blocked** | charged-coil resource; one narrow partial-discharge state-transition gap |
| Oda Nobunaga | `docs/fighters/phase-4b/oda-nobunaga.yaml` | `docs/cards/phase-4b/oda-nobunaga.yaml` | 30 | verified | two independent 6-health sidekicks, flanking |
| Tomoe Gozen | `docs/fighters/phase-4b/tomoe-gozen.yaml` | `docs/cards/phase-4b/tomoe-gozen.yaml` | 30 | verified | zone-leave event semantics, declaration-time attack override |
| William Shakespeare | `docs/fighters/phase-4b/shakespeare.yaml` | `docs/cards/phase-4b/shakespeare.yaml` | 30 | verified | persistent ordered Line, syllable threshold, completion effects |
| Hamlet | `docs/fighters/phase-4b/hamlet.yaml` | `docs/cards/phase-4b/hamlet.yaml` | 30 | verified | turn-start state choice, damage history |
| Titania | `docs/fighters/phase-4b/titania.yaml` | `docs/cards/phase-4b/titania.yaml` | 30 + 6 Glamours | verified | shuffled external card pool, active/discard lifecycle, pre-defense interrupt |
| Ciri | `docs/fighters/phase-4b/ciri.yaml` | `docs/cards/phase-4b/ciri.yaml` | 30 | verified | discard-derived Source count, highest threshold only, ongoing scheme |
| Ancient Leshen | `docs/fighters/phase-4b/ancient-leshen.yaml` | `docs/cards/phase-4b/ancient-leshen.yaml` | 30 | verified | summon/return Wolves, attack history, temporary non-defeat off-board state, ongoing scheme |
| Eredin | `docs/fighters/phase-4b/eredin.yaml` | `docs/cards/phase-4b/eredin.yaml` | 30 | verified | derived Enraged state, returnable Riders, declaration/type overrides, ongoing scheme |
| Philippa | `docs/fighters/phase-4b/philippa.yaml` | `docs/cards/phase-4b/philippa.yaml` | 30 | verified | action history, exact health assignment, combat-card replacement, ongoing scheme |

## Quantity validation

**PASS for published deck quantities.**

- 13 assigned decks × 30 action cards = **390 action-card instances**.
- Every deck manifest reconciles `available_pool_count: 30`, `game_deck_count: 30`, and `quantity_sum: 30`.
- Titania's six Glamours are external gameplay cards and are not counted in her 30-card action deck.
- Ciri has 13 physical action-card copies tagged `source`; these remain part of the ordinary 30-card deck.
- No Tales to Amaze enemy/initiative/scenario card is included in the 390-card corpus.
- Every explicit `usable_by` target is present in the corresponding fighter topology.

Tesla's **quantity, printed corpus, underfunded declaration legality, and dependent-effect fizzle behavior are verified**. The fighter/deck remains `blocked` only because one exact resource-state transition is still missing from accessible normative evidence: with exactly one charged coil, after declaring a two-coil discharge that cannot be fully paid, does that one available coil flip to discharged or remain charged?

## Proposed reusable semantic extensions

These are integration proposals only. Worker C did **not** change shared semantic files.

### C-EXT-001 — ordered persistent card sequence with threshold completion

- **Affected:** Shakespeare; all cards that enter the Line.
- **Authority:** official `Slings and Arrows` set rules: `https://restorationgames.com/wp-content/uploads/2024/03/UM-SaA-Set-Rules.pdf`.
- **Why current semantics are insufficient:** `card_zone` and `MOVE_CARD` preserve storage but do not define ordered append-at-cleanup, a per-definition numeric sequence metric, exact-threshold versus overflow behavior, or a completion-effect channel distinct from ordinary card effects/cancellation.
- **Proposed generic model:** ordered card zone + `APPEND_TO_SEQUENCE_ZONE`; definition-level sequence metric; `RESOLVE_SEQUENCE_THRESHOLD(exact, overflow)` composite; completion effects resolve from the completing definition and then source-defined cleanup occurs.
- **Integration blocker:** yes for executable Shakespeare Line semantics; source facts themselves are verified.

### C-EXT-002 — auxiliary shuffled gameplay-card deck

- **Affected:** Titania Glamours.
- **Authority:** official `Slings and Arrows` set rules.
- **Why current semantics are insufficient:** `external_definition_pool` describes definitions, but Titania needs six runtime card instances with shuffled hidden order, one public active card, a public discard, explicit return-to-bottom operations, and no automatic reshuffle.
- **Proposed generic model:** `auxiliary_card_deck` resource built from external definitions, with ordinary instance identity/zone ownership/visibility semantics and source-defined reshuffle policy.
- **Integration blocker:** yes for executable Titania Glamour lifecycle; source facts are verified.

### C-EXT-003 — pre-commit combat play-mode and participant override

- **Affected:** Tomoe Gozen `Witness My Last Battle`; Eredin `Foul Purpose`; Eredin `Implacable`; Titania `Glamour of Jealousy`.
- **Authority:** published component content and official set rules/rulings for the affected sets.
- **Why current semantics are insufficient:** current combat effects start after ordinary card commitment. These mechanics change legal target range, card visibility, card type, or the defending participant **before** the ordinary combat-card commit/reveal pipeline.
- **Proposed generic model:** declaration-time `combat_play_rules` evaluated before commit, plus a `BEFORE_COMBAT_CARD_COMMIT` interrupt window capable of replacing a participant and then continuing the same combat declaration.
- **Integration blocker:** yes for those cards; source facts are verified.

### C-EXT-004 — permissive tiered resource-discharge declaration

- **Affected:** Nikola Tesla `7 Hertz`, `Death Ray`, `Intense Experimentation`, `Lightning Storm`, `Polyphase Coils`, `Repulsion Blast`, `X-Ray Radiation`, and potentially other effects using the same discharge grammar.
- **Authority:** official `Tales to Amaze` set rules: `https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Set-rules_EN_Light.pdf`; official Rulings Archive verdict mirrored at `https://www.the-unmatched.club/tools/disputes/a1c7dc96-d154-4253-80f1-d86015292f9e`; Unmatched Reference v10 at `https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf`, which identifies the Rulings Archive as officially recognized by Restoration Games; printed deck facts from `https://unmatched.cards/umdb/decks/nikola-tesla`.
- **Confirmed semantics:** the player may declare the two-coil option with one or zero charged coils. That declaration is not automatically downgraded to the one-coil tier. The ruling permits the underfunded choice specifically so its dependent tier effect may intentionally fizzle.
- **Why current semantics are insufficient:** Phase 4A cost semantics validate a required cost before committing its dependent branch, while Tesla permits declaration independently of current resource sufficiency and can then fizzle the declared tier effect.
- **Proposed generic model:** a declared resource tier whose option legality can be independent of available resource; when the declared amount is not fully available, the dependent declared-tier effect can fizzle without resolving a lower tier. Resource mutation must remain source-defined rather than inferred from ordinary strict `SPEND_RESOURCE` semantics.
- **Integration blocker:** **yes, but narrowed to one source-blocked state transition only**: when exactly one coil is charged and two are declared, accessible normative evidence does not establish whether the available charged coil is nevertheless flipped to discharged. The `0 charged → declare 2` case has no resource-mutation ambiguity, and the dependent two-coil effect is known to be allowed to fizzle.

### C-EXT-005 — filtered zone aggregate and highest-met threshold

- **Affected:** Ciri Source cards and her cancellation-protection ability.
- **Authority:** official `The Witcher — Steel & Silver` set rules: `https://restorationgames.com/wp-content/uploads/2024/10/UM-W-SaS-Set-Rules.pdf`.
- **Why current semantics are insufficient:** Ciri repeatedly derives a count from definitions in discard and resolves **only the highest** threshold currently met.
- **Proposed generic model:** deterministic query `COUNT_CARDS(zone, filter)` usable by conditions/resources plus `RESOLVE_HIGHEST_MET_THRESHOLD(metric, ordered_thresholds)`.
- **Integration blocker:** yes for executable Source cards; source facts are verified.

### C-EXT-006 — search a card zone by definition predicate

- **Affected:** Ciri `Searching Strike` and future published search effects.
- **Authority:** published Ciri component content (`UmDb` normalized card fact cross-check).
- **Why current semantics are insufficient:** `LOOK_AT`, `REVEAL`, and `MOVE_CARD` do not by themselves specify a legal whole-zone search by predicate and the associated information/disposition semantics.
- **Proposed generic operation:** `SEARCH(zone, filter, viewer, reveal_selected, destination, post_search_order_rule)` with no implicit shuffle unless the source/rules say to shuffle.
- **Integration blocker:** yes for executable search effects; no source ambiguity identified here.

### C-EXT-007 — temporary non-defeat removal from battlefield with scheduled return

- **Affected:** Ancient Leshen `Vanish Into Murder`.
- **Authority:** published component content plus official `The Witcher — Steel & Silver` set rules/rulings.
- **Why current semantics are insufficient:** `DEFEAT` would incorrectly fire defeat/loss semantics and `RETURN_FIGHTER` only models the return half. The living Leshen must be off-board temporarily while preserving health and character identity.
- **Proposed generic model:** `REMOVE_FROM_BATTLEFIELD(target, defeat:false, preserve_health:true)` plus persisted scheduled return trigger/context consumed by `RETURN_FIGHTER`.
- **Integration blocker:** yes for this card; source facts are verified.

### C-EXT-008 — exact health assignment

- **Affected:** Philippa `Backup Plan`.
- **Authority:** published component content / current Philippa ruling index.
- **Why current semantics are insufficient:** `RECOVER` cannot lower health. The effect sets Philippa's health to exactly 5 whether that raises or lowers the current value.
- **Proposed generic operation:** `SET_HEALTH(target, value, source_limits)` distinct from damage and recovery so damage/recovery triggers are not fabricated.
- **Integration blocker:** yes for this card; source facts are verified.

### C-EXT-009 — first-class zone-leave event with cause context

- **Affected:** Tomoe Gozen special ability.
- **Authority:** official `Sun's Origin` set rules: `https://restorationgames.com/wp-content/uploads/2024/01/UM-SO-Set-Rules.pdf`.
- **Why current semantics are insufficient:** Tomoe triggers whenever an opposing hero itself leaves her zone, including movement, placement, removal and multiple transitions inside one effect, but does **not** trigger when Tomoe moving merely changes the zone relationship.
- **Proposed generic event:** `FIGHTER_LEFT_ZONE(fighter, zone, cause, source_operation)` emitted for each actual fighter transition, with sufficient cause context to distinguish movement of the observed fighter from movement of the zone-owning fighter.
- **Integration blocker:** yes for deterministic Tomoe ability execution; source facts are verified.

## Mechanics that fit the existing Phase 4A model

No shared extension is proposed for these because current semantics are sufficient:

- Oda's two Honor Guards: existing independently damageable fighter instances + explicit health semantics.
- Jill's gadget state: public enum/counter state + ordinary printed-value conditions and generic effects.
- Golden Bat and Hamlet historical checks: explicit persistent turn history under `FX-021`.
- Leshen Wolves: canonical summonable sidekick instances + existing `SUMMON`/`RETURN_FIGHTER`; only Vanish needs the separate non-defeat removal extension.
- Eredin Enraged and Ciri Source count: derived authoritative state, not manually editable UI counters.
- Witcher ongoing schemes: public card zone + continuous/replacement effects + source-defined discard conditions, matching Phase 4A Geralt/Yennefer & Triss.
- Philippa `Do My Bidding`: existing Phase 4A `REPLACE_COMBAT_CARD` composite.

## Ambiguities / blockers

1. **Nikola Tesla — narrow P1 source gap:** declaration of two coils with only one or zero charged is confirmed legal, it must not auto-downgrade to the one-coil tier, and the dependent declared-tier effect may intentionally fizzle. The only unresolved case is `1 charged + declare 2`: does the available charged coil flip to discharged, or is resource state left unchanged because the two-coil discharge cannot be fully completed?
2. **Shakespeare × `End the turn` — integration cross-check:** current general rulings say `END_TURN` jumps to the current action Cleanup Step. Shakespeare's Line is itself a cleanup destination/lifecycle rule. Integration should run an explicit cross-fighter fixture to confirm Line append/completion ordering when `END_TURN` occurs; no contradictory source was found.
3. **No cooperative inference:** Tales to Amaze enemy/scenario procedures were intentionally not used to fill gaps in competitive hero behavior.

## Source gaps

- Tales to Amaze competitive card metadata is complete in published UmDb, and Restoration explicitly confirms the four heroes are usable in ordinary Unmatched. A separate current official publisher-hosted per-hero card-text PDF was not located during this pass; normalized card facts therefore retain UmDb provenance plus current Rulings Archive cross-checks.
- Tesla's only remaining source gap is the resource-state result of `1 charged + declare 2`; declaration legality and dependent-effect fizzle are no longer open questions.
- No fan/community `/decks/...` balance patch was imported. Community/reference pages were used only as secondary cross-checks where an official set rule or published UmDb entry did not expose the same detail.

## Worker 4B-C Handoff

Branch: `phase-4b-worker-c-modern`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: emitted as the exact branch SHA in the final Worker C handoff response; a commit cannot self-contain its own SHA  
Assigned fighters: 13  
Verified: 12 — Annie Christmas, Dr. Jill Trent, Golden Bat, Oda Nobunaga, Tomoe Gozen, William Shakespeare, Hamlet, Titania, Ciri, Ancient Leshen, Eredin, Philippa  
Blocked: 1 — Nikola Tesla (`1 charged + declare 2` partial-discharge resource transition only)  
Quantity validation: **PASS — 390/390 action cards; 13/13 deck quantities reconcile**  
Schema-extension proposals: C-EXT-001 through C-EXT-009 above  
New ambiguity/blockers: Tesla single remaining partial-discharge state transition; Shakespeare/END_TURN integration fixture  
Source gaps: Tesla `1 charged + declare 2` resource mutation only; no separate publisher-hosted Tales competitive card-text corpus found  
Files created: 27 — 13 fighter manifests, 13 deck manifests, this report  
Shared semantic files changed: **none**  
Phase 4A manifests rewritten: **none**  
Merged to `main`: **no**