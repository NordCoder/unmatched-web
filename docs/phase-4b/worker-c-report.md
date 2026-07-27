# Phase 4B Worker C Reconciliation Report

## 1. Worker identity

- Repository: `NordCoder/unmatched-web`
- Branch: `phase-4b-worker-c-modern`
- Authorized base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: 13 modern mechanic-heavy competitive fighters
- Shared semantic files changed: **none**
- Phase 4A manifests changed: **none**
- Merge to `main`: **no**
- Tales to Amaze scope: competitive player-hero behavior only; cooperative enemy/minion/initiative/threat/scenario logic excluded.

## 2. Status matrix

All 13 pairs are `partial` after dependency reconciliation. This is not a blanket evidence downgrade: every pair has deterministic published behavior, but each requires at least one shared runtime capability absent from the current shared gameplay model. Only four pairs contain evidence-qualified edge-case normalizations.

| Fighter | Status | Evidence | Semantics | Integration | Policy | Requirements |
| --- | --- | --- | --- | --- | --- | --- |
| `annie-christmas` | partial | verified | verified | requires_shared_extension | ready | C-REQ-012, C-REQ-013 |
| `dr-jill-trent` | partial | verified | verified | requires_shared_extension | ready | C-REQ-013 |
| `golden-bat` | partial | verified | verified | requires_shared_extension | ready | C-REQ-013 |
| `nikola-tesla` | partial | qualified | qualified | requires_shared_extension | ready | C-REQ-004, C-REQ-013 |
| `oda-nobunaga` | partial | verified | verified | requires_shared_extension | not_applicable | C-REQ-003, C-REQ-010, C-REQ-013 |
| `tomoe-gozen` | partial | verified | verified | requires_shared_extension | not_applicable | C-REQ-003, C-REQ-013 |
| `shakespeare` | partial | verified | verified | requires_shared_extension | not_applicable | C-REQ-001, C-REQ-013 |
| `hamlet` | partial | qualified | qualified | requires_shared_extension | not_applicable | C-REQ-013 |
| `titania` | partial | verified | verified | requires_shared_extension | not_applicable | C-REQ-002, C-REQ-003, C-REQ-011, C-REQ-013 |
| `ciri` | partial | verified | verified | requires_shared_extension | not_applicable | C-REQ-005, C-REQ-006, C-REQ-013 |
| `ancient-leshen` | partial | qualified | qualified | requires_shared_extension | not_applicable | C-REQ-007, C-REQ-010, C-REQ-011, C-REQ-013 |
| `eredin` | partial | qualified | qualified | requires_shared_extension | not_applicable | C-REQ-003, C-REQ-011, C-REQ-013 |
| `philippa` | partial | verified | verified | requires_shared_extension | not_applicable | C-REQ-008, C-REQ-011, C-REQ-013 |

Summary: **Verified 0 / Partial 13 / Blocked 0**. Evidence/semantics are verified for 9 pairs and qualified for 4; no pair has unresolved deterministic behavior.

## 3. C-REQ requirements

```yaml
requirements:
  - id: C-REQ-001
    family: card_zones_and_auxiliary_systems
    severity: required
    affects: [shakespeare]
    established_rule: Shakespeare combat cards enter an ordered public Line at Cleanup; exactly 10 syllables resolves the newest card's completion effect, over 10 clears without completion, and cleanup discards only instances still in the Line.
    missing_capability: Ordered persistent gameplay-card sequence with cleanup append, sequence metric, exact/overflow threshold lifecycle, and a completion-effect channel distinct from ordinary card effects.
    proposed_generic_contract: First-class ordered card zone; append at source-defined lifecycle point; evaluate metric; resolve full staged completion effect at exact threshold; clear current zone contents after resolution without stale-snapshot discards.
    evidence_status: verified

  - id: C-REQ-002
    family: card_zones_and_auxiliary_systems
    severity: required
    affects: [titania]
    established_rule: Six Glamour cards form a shuffled auxiliary gameplay deck with hidden order, one active face-up card, public discard, explicit returns, and no automatic reshuffle.
    missing_capability: Runtime auxiliary gameplay-card definitions/instances and zones outside the ordinary action deck.
    proposed_generic_contract: Auxiliary card definitions instantiate normal owned card instances into explicit zones with visibility, ordering, source lifetime, movement and discard semantics.
    evidence_status: verified

  - id: C-REQ-003
    family: movement_targeting_and_combat_legality
    severity: required
    affects: [oda-nobunaga, tomoe-gozen, titania, eredin]
    established_rule: Published effects can alter legal targets, card visibility/type, or the combat participant before or during the commit pipeline; team-play defender replacement can cross fighter/card controllers.
    missing_capability: Declaration-time combat legality and same-combat participant replacement with independent fighter/controller and card owner/controller identities.
    proposed_generic_contract: Structured combat play rules plus pre-commit/same-combat participant replacement; combat state preserves fighter identity/controller separately from committed card instance owner/controller.
    evidence_status: verified

  - id: C-REQ-004
    family: resources_actions_and_turn_control
    severity: required
    affects: [nikola-tesla]
    established_rule: Tesla may legally declare a two-coil discharge with zero or one charged coils; the declared tier is not downgraded and a dependent full-discharge effect may fizzle.
    missing_capability: Resource declaration legality independent from available quantity, partial consumption up to declared amount, and source-defined dependent-resolution behavior.
    proposed_generic_contract: declared_resource_discharge(resource, declared, consume_available, outcome_by_declared_tier, dependent_requires_full_declared) with no implicit lower-tier fallback.
    evidence_status: verified

  - id: C-REQ-005
    family: derived_attributes_and_modifiers
    severity: required
    affects: [ciri]
    established_rule: Source count is derived from Source-tagged cards in discard and Source cards resolve only the highest threshold currently met.
    missing_capability: Filtered card-zone aggregate plus ordered highest-met threshold resolution.
    proposed_generic_contract: Structured count(selector) expression over authoritative card instances plus highest_met_threshold(metric, ordered_thresholds) resolving exactly one tier.
    evidence_status: verified

  - id: C-REQ-006
    family: search_randomness_and_disclosure
    severity: required
    affects: [ciri]
    established_rule: Searching Strike searches the deck by card property, reveals the selected card, moves it to hand, then shuffles the searched deck.
    missing_capability: Whole-zone search by structured predicate with viewer/disclosure, destination and post-search disposition.
    proposed_generic_contract: search_card_zone(zone, filter, viewer, reveal_selected, destination, post_search) with authoritative post-search shuffle.
    evidence_status: verified

  - id: C-REQ-007
    family: fighter_presence_and_occupancy
    severity: required
    affects: [ancient-leshen]
    established_rule: Vanish Into Murder removes a living Leshen from the battlefield without defeat and schedules its return at the next controller turn start while preserving health.
    missing_capability: Temporary non-defeat battlefield absence with persistent scheduled return context.
    proposed_generic_contract: remove_from_battlefield(defeat:false, preserve_health:true, return_trigger) plus return_fighter consuming the persisted scheduled context.
    evidence_status: verified

  - id: C-REQ-008
    family: damage_and_health
    severity: required
    affects: [philippa]
    established_rule: Backup Plan sets Philippa's health to exactly 5 and may therefore lower health as well as raise it.
    missing_capability: Exact health assignment distinct from damage and recovery.
    proposed_generic_contract: set_health(target, value, source_limits) with explicit health-bound semantics and no damage/recovery event substitution.
    evidence_status: verified

  - id: C-REQ-010
    family: identity_history_and_provenance
    severity: required
    affects: [oda-nobunaga, ancient-leshen]
    established_rule: Effects inspect historical state independently for repeated runtime fighter instances such as two Honor Guards or two Wolves.
    missing_capability: Persistent historical values keyed by runtime fighter instance identity.
    proposed_generic_contract: Typed fighter-instance-keyed maps with source-defined reset boundaries and deterministic instance lookup.
    evidence_status: verified

  - id: C-REQ-011
    family: derived_attributes_and_modifiers
    severity: required
    affects: [titania, ancient-leshen, eredin, philippa]
    established_rule: Fighters may have definition-specific base attributes and continuous permissions/value overrides that disappear immediately when their source/condition ceases to apply.
    missing_capability: Fighter-definition base attributes plus source-lifetime derived attribute/rule modifiers without rollback mutation.
    proposed_generic_contract: effective attribute/permission layer evaluated from base definition plus active continuous modifiers keyed to condition/source lifetime.
    evidence_status: verified

  - id: C-REQ-012
    family: damage_and_health
    severity: required
    affects: [annie-christmas]
    established_rule: Mississippi Queen keeps Annie at or above 1 health while incoming damage still counts as dealt for combat/result/trigger observers.
    missing_capability: Health-result floor distinct from damage prevention.
    proposed_generic_contract: Damage application records the dealt amount normally, then applies a source-lifetime minimum-health clamp to the resulting health value.
    evidence_status: verified

  - id: C-REQ-013
    family: resolution_and_choices
    severity: required
    affects: [annie-christmas, dr-jill-trent, golden-bat, nikola-tesla, oda-nobunaga, tomoe-gozen, shakespeare, hamlet, titania, ciri, ancient-leshen, eredin, philippa]
    established_rule: Published effects frequently require operation -> result binding -> informed choice/cost -> dependent continuation, including empty dependent domains and cross-combat captured state.
    missing_capability: Ordered effect stages with typed choices, operation-result bindings, later-stage references, and deterministic empty-domain handling.
    proposed_generic_contract: effects[].stages[] where each stage has id/when/choices/costs/operations; choices use id/owner/visibility/optional/domain/bind; operations use bind_result; later stages read ref; empty required dependent domains create no pending choice and skip only the impossible dependent continuation.
    evidence_status: verified
```

`C-EXT-009` is **withdrawn and not promoted to a C-REQ**. Existing relocation transition facts (`left_zone`, removal, etc.) are sufficient for Tomoe's ability; adding a second parallel fighter-left-zone event model would duplicate shared semantics.

## 4. Interpretations and evidence qualifications

```yaml
interpretations:
  - id: ancient-leshen-dormant-return-draw
    affects: [ancient-leshen, vanish-into-murder, dormant-player]
    status: derived
    confidence: high
    behavior: If Vanish leaves the controller with no fighters, dormancy starts at the normal end-of-action checkpoint. The scheduled next-turn return still resolves; after a successful Leshen return the no-fighter dormant condition no longer applies and the card's subsequent draw 1 resolves.
    source_refs: [leshen-deck-index, dormant-rule-changes, worker-c-normalization-leshen]
    replacement_condition: Replace if an exact publisher/designer ruling explicitly suppresses or reorders the draw after a dormant Leshen return.

  - id: hamlet-readiness-failed-placement
    affects: [hamlet, the-readiness-is-all]
    status: project_normalization
    confidence: medium-high
    behavior: The cannot-leave restriction depends on successful placement into the selected destination. If a legal occupied-space selection makes PLACE fail, the fighter remains in its prior space and no cannot-leave restriction attaches to that prior space.
    source_refs: [hamlet-deck-index, placement-rule-changes, worker-c-normalization-hamlet]
    replacement_condition: Replace if an exact publisher/designer ruling states that the restriction attaches after failed placement or defines another dependency.

  - id: eredin-icy-guile-self-defeat
    affects: [eredin, icy-guile, red-rider]
    status: derived
    confidence: high
    behavior: The currently attacking Red Rider is a legal Icy Guile cost target because the published effect supplies no another exclusion. If it defeats itself, IGNORE_VALUE resolves and the already-started combat continues unless the defeat ended the game.
    source_refs: [eredin-deck-index, core-combat-continuation, worker-c-normalization-eredin]
    replacement_condition: Replace if an exact publisher/designer ruling excludes the current combat fighter or terminates this combat after self-defeat.

  - id: tesla-underfunded-two-coil-resource-transition
    affects: [nikola-tesla, charged-coils, coil-discharge-effects]
    status: derived
    confidence: high
    behavior: With exactly one charged coil, a legal declaration to discharge two consumes the available charged coil (1 -> 0), keeps the two-coil branch/tier selected, never falls back to one coil, and permits a source-dependent full-discharge effect to fizzle.
    source_refs: [tta-set-rules-tesla, tesla-underfunded-ruling, worker-c-normalization]
    replacement_condition: Replace if an exact publisher/designer ruling explicitly defines a different one-charged-plus-declare-two resource transition.
```

**Unresolved behavior choices: none.** These four interpretations are deterministic Worker C integration inputs with explicit evidence qualifiers, not blockers.

## 5. Validation

- Quantity: **PASS — 390/390 action cards across 13 fighters**.
- Titania: six Glamour auxiliary gameplay cards are separate runtime instances and are **not** included in her 30-card action deck.
- `usable_by`: **PASS** for all 13 pairs.
- References: **PASS**; runtime-state/resource/zone references are explicit or linked to a C-REQ.
- Source coverage: **PASS** for 9 pairs; **QUALIFIED** for Tesla, Hamlet, Ancient Leshen and Eredin only because of the four interpretations above.
- Semantic structure: **PASS**; dependent resolution uses staged effects with explicit bindings/references; legacy `choices_after`, `nested_choice`, `followup` and ad-hoc continuation branches are removed.
- Integration: **requires_shared_extension** for all 13 pairs, with per-pair dependencies listed in the status matrix.
- Fan content: **PASS**; no fan `/decks/...` content imported.

## 6. Source gaps

No source gap blocks deterministic gameplay behavior.

- Tesla `1 charged + declare 2 -> 0`: no exact-case publisher sentence; high-confidence derived normalization.
- Ancient Leshen dormant return + draw: no exact-case publisher sentence; high-confidence derived normalization.
- Hamlet failed placement dependency: no exact-case publisher sentence; medium-high project normalization.
- Eredin Icy Guile self-defeat: no exact-case publisher sentence; high-confidence derived normalization.
- Tales to Amaze competitive action-card corpus is covered through published UmDb records plus official set/product/ruling material; no separate current publisher-hosted per-hero card-text PDF was located.

## 7. Files

Reconciliation changed exactly **27 Worker C files**:

- 13 fighter manifests under `docs/fighters/phase-4b/`;
- 13 deck manifests under `docs/cards/phase-4b/`;
- `docs/phase-4b/worker-c-report.md`.

No shared semantic file and no Phase 4A manifest was changed.

## Worker 4B-C Reconciliation Handoff

Branch: `phase-4b-worker-c-modern`  
Exact final Head: **supplied in the external handoff because a persisted file cannot contain the SHA of the commit that contains that same file**  
Assigned: **13**  
Verified: **0**  
Partial: **13**  
Blocked: **0**  
Quantity result: **PASS — 390/390 action cards; 13/13 decks; Titania Glamours external to her 30-card action deck**  
Requirements count: **12 active C-REQs** (`001-008`, `010-013`); `C-EXT-009` withdrawn  
Qualified interpretations: **4** — Ancient Leshen dormant return/draw; Hamlet failed placement dependency; Eredin Icy Guile self-defeat; Tesla underfunded two-coil resource transition  
Unresolved behavior choices: **none**  
Files changed: **27 Worker C files**  
Shared files changed: **none**  
Phase 4A manifests changed: **none**  
Merged to `main`: **no**
