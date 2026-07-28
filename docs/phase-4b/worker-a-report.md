# Phase 4B Worker A — Reconciliation Contract v1

## 1. Worker identity

- Repository: `NordCoder/unmatched-web`
- Branch: `phase-4b-worker-a-classics`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: early/classic + retired licensed
- Assigned fighters: **17**
- Owned files: **35**
- Shared schema/mechanics/rules/rulings/set files modified: **none**
- Merge to `main`: **not performed**

Status policy:

- `verified` — evidence and deterministic semantics are verified, integration is ready, and policy is not blocked;
- `partial` — deterministic behavior is established, but a shared capability or qualified/reversible normalization remains;
- `blocked` — deterministic implementation still depends on unresolved evidence or an unresolved project policy decision.

Engine-only gaps therefore remain `partial`, not `blocked`.

## 2. Final status matrix

| Fighter | Status | Evidence | Semantics | Integration | Policy | Requirements |
| --- | --- | --- | --- | --- | --- | --- |
| `alice` | verified | verified | verified | ready | not_applicable | — |
| `king-arthur` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-006, A-REQ-013 |
| `medusa` | verified | verified | verified | ready | not_applicable | — |
| `sinbad` | verified | verified | verified | ready | not_applicable | — |
| `robin-hood` | verified | verified | verified | ready | not_applicable | — |
| `bigfoot` | verified | verified | verified | ready | not_applicable | — |
| `robert-muldoon` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-003, A-REQ-004 |
| `invisible-man` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-003, A-REQ-005 |
| `jekyll-and-hyde` | verified | verified | verified | ready | not_applicable | — |
| `buffy` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-014 |
| `willow` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-006 |
| `spike` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-003, A-REQ-007 |
| `angel` | verified | verified | verified | ready | not_applicable | — |
| `little-red-riding-hood` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-008, A-REQ-009 |
| `beowulf` | verified | verified | verified | ready | not_applicable | — |
| `deadpool` | blocked | qualified | verified | requires_shared_extension | blocked | A-REQ-008, A-REQ-015 |
| `yennenga` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-010, A-REQ-011, A-REQ-012 |

Totals: **8 verified / 8 partial / 1 blocked**.

The audited engine-only cases are intentionally not blocked: King Arthur exact health/search, Muldoon traps, Invisible Man Fog/Vanish, Willow exact-health return, Spike Shadows/blind BOOST transform, Little Red Basket, and Yennenga damage/path/combat-replacement behavior are deterministic. Deadpool alone remains blocked because the digital-adaptation policy is not decided.

## 3. A-REQ integration requirements

```yaml
- id: A-REQ-003
  family: battlefield_components_and_paths
  severity: P0
  affects:
    - robert-muldoon/ingen-traps
    - robert-muldoon/remote-detonation
    - invisible-man/fog
    - invisible-man/fog-dependent-cards
    - spike/shadows
    - spike/shadow-dependent-cards
  established_rule: Traps, Fog and Shadows are persistent positioned battlefield components with identity, ownership and source-specific lifecycle/occupancy behavior.
  missing_capability: Shared runtime does not define reusable positioned non-fighter component instances and generic placement/movement/query operations.
  proposed_generic_contract: Component instances carry instance_id, definition/type, owner/controller, lifecycle location, board space and occupancy policy, with generic place/move/select/count operations.
  evidence_status: verified

- id: A-REQ-004
  family: movement_targeting_and_combat_legality
  severity: P0
  affects:
    - robert-muldoon/ingen-traps
  established_rule: An opposing fighter entering a trap space immediately interrupts that relocation, takes 1 damage, the trap goes to the box, and the trap-return draw resolves; authoritative rulings also cover effect placement into a trap space.
  missing_capability: Shared relocation has no protected mid-move/entry interruption contract tied to resolved entry events.
  proposed_generic_contract: Emit entry events for resolved movement/placement transitions and allow protected reactions to interrupt the active movement while preserving the reached space and ordered follow-up operations.
  evidence_status: verified

- id: A-REQ-005
  family: fighter_presence_and_occupancy
  severity: P0
  affects:
    - invisible-man/vanish
  established_rule: Vanish removes an undefeated Invisible Man from the board until the start of his next turn, then returns him to a chosen space; dormant-player behavior applies while no fighter is present.
  missing_capability: Shared presence model lacks temporary undefeated off-board state plus source-owned scheduled return integrated with dormant turns.
  proposed_generic_contract: Presence distinguishes on_board, defeated and off_board_undefeated; source hooks may restore the same fighter instance at the defined lifecycle point without converting absence into defeat.
  evidence_status: verified

- id: A-REQ-006
  family: damage_and_health
  severity: P0
  affects:
    - king-arthur/the-holy-grail
    - willow/resurrect
  established_rule: The Holy Grail assigns Arthur's current health exactly to 8 when its condition is met; Resurrect returns a defeated fighter with exactly 3 health.
  missing_capability: RECOVER is insufficient for exact assignment, and revive needs an atomic source-defined health result rather than an invalid intermediate state.
  proposed_generic_contract: Support exact SET_HEALTH for living fighters and atomic RETURN_FIGHTER with explicit exact/start-health semantics while preserving defeat/recovery event distinctions.
  evidence_status: verified

- id: A-REQ-007
  family: boost_pipeline
  severity: P1
  affects:
    - spike/always-surprising
  established_rule: Always Surprising performs a blind BOOST and doubles the resolved BOOST amount when Drusilla is on a Shadow.
  missing_capability: Shared BOOST processing cannot transform the resolved blind-BOOST amount between reveal/read and application.
  proposed_generic_contract: BOOST exposes a bound resolved amount and supports source-defined transforms before that amount changes combat value.
  evidence_status: verified

- id: A-REQ-008
  family: movement_targeting_and_combat_legality
  severity: P0
  affects:
    - little-red-riding-hood/what-big-ears-you-have
    - deadpool/xavier-institute-faculty
  established_rule: Legal combat-card/attack mode can depend on source state or the card itself rather than immutable printed fighter/card type alone.
  missing_capability: Shared legality does not expose declarative conditional card modes or card-specific melee/ranged permissions.
  proposed_generic_contract: Preserve immutable printed metadata while deriving legal play/attack modes from scoped declarative permissions evaluated at commit/declaration time.
  evidence_status: verified

- id: A-REQ-009
  family: card_zones_and_auxiliary_systems
  severity: P0
  affects:
    - little-red-riding-hood/basket-reference
    - little-red-riding-hood/basket-dependent-cards
  established_rule: Little Red's Basket is not one of the 30 action cards but is a real auxiliary object in ordered discard state; its top-object symbol determines Basket behavior and wild binds to exactly one effective item per source resolution.
  missing_capability: Standard card zones assume ordinary card instances and cannot represent a non-action auxiliary object that participates in ordered-zone top state without affecting deck counts.
  proposed_generic_contract: Ordered zones may contain typed auxiliary objects distinct from action-card instances; top-object metadata is queryable and wild-derived values are bound once per resolving source.
  evidence_status: verified

- id: A-REQ-010
  family: damage_and_health
  severity: P0
  affects:
    - yennenga/shield-of-the-archers-ability
  established_rule: Pending damage to Yennenga may be divided among eligible Archers in her zone; no Archer may receive more than remaining health and all residual damage remains on Yennenga with the original damage identity/type.
  missing_capability: REDIRECT_DAMAGE cannot represent partial allocation of one pending event across multiple recipients plus a residual original target.
  proposed_generic_contract: Pending damage supports an allocation stage with eligible recipients, per-recipient caps, residual original target and preserved source/type/event identity.
  evidence_status: verified

- id: A-REQ-011
  family: battlefield_components_and_paths
  severity: P0
  affects:
    - yennenga/stallion-charge
  established_rule: Stallion Charge damages each opposing fighter whose space Yennenga moved through during that specific movement.
  missing_capability: Shared MOVE result does not expose the traversed-space path for later source stages.
  proposed_generic_contract: MOVE may bind origin, ordered traversed_spaces and destination; later selectors can query distinct fighter instances whose spaces occurred in that path.
  evidence_status: verified

- id: A-REQ-012
  family: movement_targeting_and_combat_legality
  severity: P0
  affects:
    - yennenga/surprise-volley
  established_rule: Surprise Volley may return a defeated Archer and, if it returns, that Archer replaces Yennenga as the attacker in the already-active combat rather than starting a new attack.
  missing_capability: Shared combat state has no explicit same-combat participant replacement preserving combat identity, defender, committed cards and the current timing window.
  proposed_generic_contract: Provide a generic REPLACE_COMBAT_PARTICIPANT operation keyed by active combat instance and role; preserve combat/card/controller provenance while subsequent stages resolve against the replacement fighter instance.
  evidence_status: verified

- id: A-REQ-013
  family: search_randomness_and_disclosure
  severity: P1
  affects:
    - king-arthur/the-lady-of-the-lake
  established_rule: The Lady of the Lake searches Arthur's deck or discard for Excalibur, moves the found card to hand, and shuffles the deck when the deck was searched.
  missing_capability: Shared card-zone model has no whole-zone structured search operation with source-zone/result binding, visibility and post-search disposition.
  proposed_generic_contract: Search one or more zones with a structured card predicate and explicit viewer/disclosure policy; bind the found card instance, source zone and searched-zone facts so later stages can move the card and apply authoritative post-search shuffle behavior.
  evidence_status: verified

- id: A-REQ-014
  family: fighter_presence_and_occupancy
  severity: P1
  affects:
    - buffy/setup-selected-sidekick
  established_rule: Buffy chooses Giles or Xander before deck construction; only the chosen sidekick and matching five-card package participate in that match.
  missing_capability: Shared setup/runtime does not explicitly guarantee that an unselected fighter definition is absent from active, defeated and targetable match rosters.
  proposed_generic_contract: Match setup materializes an active fighter roster from canonical definitions; unselected definitions remain out_of_play and cannot satisfy fighter selectors or revival effects.
  evidence_status: verified

- id: A-REQ-015
  family: digital_adaptation_policy
  severity: P0
  affects:
    - deadpool/maximum-effort
    - deadpool/physical-social-card-effects
  established_rule: Published Deadpool behavior intentionally depends on external facts/actions including real names, food, noises, set ownership, mirrors, sleeves, writing on a physical card, clothing, wagers and subjective board colour.
  missing_capability: No project policy defines which physical/social predicates are automatic, player-confirmed, digitally substituted or unsupported online.
  proposed_generic_contract: First choose a central adaptation policy per external predicate/action class; only then expose generic external-predicate/instruction confirmation and persistence primitives.
  evidence_status: verified
```

Requirements count: **13**.

Orchestrator consolidation notes:

- `A-REQ-012` is the same generic problem family as Worker C `C-REQ-003` same-combat participant replacement; these should converge into one shared combat-participant contract.
- `A-REQ-013` is the same generic problem family as Worker C `C-REQ-006` structured card-zone search; these should converge into one shared search/disclosure contract.
- `A-REQ-005` is semantically compatible with Worker B `B-REQ-003` and Worker C `C-REQ-007` undefeated off-board lifecycle.
- `A-REQ-006` exact health assignment overlaps Worker C `C-REQ-008` and should converge at integration.
- `A-REQ-003` positioned component instances overlap Worker B `B-REQ-006`; source-specific lifecycle/occupancy constraints remain data, not separate engines.

Mirror-safety is not a new Worker A requirement: existing fighter/card definition-instance contracts already require runtime instance identity and immutable ownership. Local definition IDs in this corpus must resolve through the owning player's runtime instances, including repeated sidekicks and mirror matches.

## 4. Evidence uncertainties / project normalizations

### Qualified evidence

`deadpool` is the only qualified-evidence fighter:

- first-party Mondo material establishes the official release;
- published UmDb supplies the complete 30-card corpus used by this worker;
- no first-party online full card-text/component dump was located.

This is a provenance-quality limitation, not an unknown deterministic card fact. Deadpool remains blocked by policy, not by this evidence limitation.

### Interpretations / normalizations

- `beowulf/the-ancient-heirloom`: `derived`, high confidence. The two optional DURING COMBAT clauses resolve in printed order; either or both may be used. The deck manifest contains the replacement condition for a future exact authoritative ruling.
- Mirror-safe scoping: authoritative shared-schema behavior. Definition IDs are manifest-local; runtime state/targets resolve to fighter/card instances under owner scope.
- No fan balance patch, UnPatched record or community `/decks/...` record is canonical evidence.

No additional project-selected gameplay behavior was introduced during reconciliation. In particular, no Deadpool digital substitute/policy was invented.

## 5. Validation matrix

| Fighter | Quantity | usable_by | references | source coverage | semantic structure | integration | fan content |
| --- | --- | --- | --- | --- | --- | --- | --- |
| alice | pass | pass | pass | pass | pass | pass | none |
| king-arthur | pass | pass | pass | pass | pass | requires_extensions | none |
| medusa | pass | pass | pass | pass | pass | pass | none |
| sinbad | pass | pass | pass | pass | pass | pass | none |
| robin-hood | pass | pass | pass | pass | pass | pass | none |
| bigfoot | pass | pass | pass | pass | pass | pass | none |
| robert-muldoon | pass | pass | pass | pass | pass | requires_extensions | none |
| invisible-man | pass | pass | pass | pass | pass | requires_extensions | none |
| jekyll-and-hyde | pass | pass | pass | pass | pass | pass | none |
| buffy | pass | pass | pass | pass | pass | requires_extensions | none |
| willow | pass | pass | pass | pass | pass | requires_extensions | none |
| spike | pass | pass | pass | pass | pass | requires_extensions | none |
| angel | pass | pass | pass | pass | pass | pass | none |
| little-red-riding-hood | pass | pass | pass | pass | pass | requires_extensions | none |
| beowulf | pass | pass | pass | pass | pass | pass | none |
| deadpool | pass | pass | pass | qualified | pass | blocked | none |
| yennenga | pass | pass | pass | pass | pass | requires_extensions | none |

Quantity validation: **PASS (17/17)**.

Structural reconciliation checks:

- every fighter/deck pair declares unified `status`, `verification`, structured `sources` and `validation`;
- non-trivial dependent effects use `stages` with explicit choices and operation-result bindings;
- ordinary procedural `choices_after`, `nested_choice`, `followup`, free-form `normalization_blocker`, and procedural `REQUEST_CHOICE` are not used as escape hatches;
- `external_definitions` remains only for a real gameplay definition (`yennenga/volley`);
- standard `deck`, `hand`, `discard` are not redundantly declared as custom zones;
- engine-only gaps use shared `A-REQ` references rather than repeated prose blockers;
- Buffy materializes only the selected sidekick into the active roster; the unselected definition remains out of play;
- Little Red's Basket remains distinct from an ordinary action-card instance while participating in ordered discard state;
- Yennenga movement path, pending-damage allocation and current-combat attacker replacement are separate typed concerns;
- Deadpool external predicates remain explicit and policy-blocked.

## 6. Source gaps

Blocking source gaps: **0**.

Qualified non-blocking source gap: **1** — Deadpool lacks a located first-party online full card-text/component dump; official release provenance plus published UmDb remain sufficient to reconstruct deterministic printed semantics.

Unresolved policy items: **1** — Deadpool digital adaptation policy (`A-REQ-015`).

## 7. Files owned

Worker A owns and reconciled exactly:

- 17 assigned files under `docs/fighters/phase-4b/`;
- 17 matching files under `docs/cards/phase-4b/`;
- `docs/phase-4b/worker-a-report.md`.

Files changed relative to Authorized Base in Worker A scope: **35**.

No shared schema, mechanics, rules, rulings, set registry, research plan, README or Phase 4A manifest was changed.

## Worker 4B-A Reconciliation Handoff

Branch: `phase-4b-worker-a-classics`  
Exact final Head: **reported externally after the final repository write; a tracked file cannot self-embed the SHA of the commit that contains that value**  
Assigned count: **17**  
Verified count: **8**  
Partial count: **8**  
Blocked count: **1**  
Quantity validation: **PASS (17/17)**  
Requirements count: **13**  
Unresolved evidence items: **1 qualified, non-blocking provenance limitation (Deadpool); 0 blocking evidence gaps**  
Unresolved policy items: **1 — Deadpool digital adaptation policy (`A-REQ-015`)**  
Files changed: **35**  
Merge to main: **not performed**
