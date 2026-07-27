# Phase 4B Worker A — Reconciliation Contract v1

## 1. Worker identity

- Repository: `NordCoder/unmatched-web`
- Branch: `phase-4b-worker-a-classics`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: early/classic + retired licensed
- Assigned fighters: **17**
- Owned files: **35**
- Shared schema/mechanics/rules/set files modified: **none**
- Merge to `main`: **not performed**

Reconciliation policy:

- `verified` = published behavior is established and the manifest needs no Worker-A-specific shared extension or unresolved policy decision;
- `partial` = deterministic behavior is established but a shared runtime capability is still required;
- `blocked` = deterministic implementation still depends on unresolved evidence or project policy.

## 2. Final status matrix

| Fighter | Status | Evidence | Semantics | Integration | Policy | Requirements |
| --- | --- | --- | --- | --- | --- | --- |
| `alice` | verified | verified | verified | ready | not_applicable | — |
| `king-arthur` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-006 |
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
| `yennenga` | partial | verified | verified | requires_shared_extension | not_applicable | A-REQ-010, A-REQ-011 |

Totals:

- Verified: **8**
- Partial: **8**
- Blocked: **1**

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
  established_rule: Traps, Fog and Shadows are persistent battlefield components with identity, ownership, position and source-specific lifecycle/occupancy rules.
  missing_capability: Shared runtime does not yet define reusable positioned non-fighter component instances and generic placement/movement/query operations.
  proposed_generic_contract: Component instances have instance_id, definition/type, owner/controller, lifecycle location, board space, occupancy/stacking policy and generic PLACE_COMPONENT/MOVE_COMPONENT/select/count semantics.
  evidence_status: verified

- id: A-REQ-004
  family: movement_targeting_and_combat_legality
  severity: P0
  affects:
    - robert-muldoon/ingen-traps
  established_rule: An opposing fighter entering a trap space immediately stops that movement, takes 1 damage, and the trap goes to the box; source rulings also cover effect placement into a trap space.
  missing_capability: Shared movement resolution has no generic protected mid-move/entry interruption contract tied to space-entry/placement events.
  proposed_generic_contract: Emit entry events per resolved relocation step and allow protected source reactions to interrupt the active movement while preserving the reached space and subsequent ordered effects.
  evidence_status: verified

- id: A-REQ-005
  family: fighter_presence_and_occupancy
  severity: P0
  affects:
    - invisible-man/vanish
  established_rule: Vanish removes an undefeated Invisible Man from the board until the start of his next turn, then returns him to a chosen space; Dormant Player rules apply while no fighter is on board.
  missing_capability: Shared fighter presence model lacks temporary undefeated off-board state plus source-owned scheduled return integrated with dormant turn handling.
  proposed_generic_contract: Fighter presence distinguishes on_board, defeated and off_board_undefeated; scheduled source hooks can restore presence before/within the correct turn-start lifecycle without treating the fighter as defeated.
  evidence_status: verified

- id: A-REQ-006
  family: damage_and_health
  severity: P0
  affects:
    - king-arthur/the-holy-grail
    - willow/resurrect
  established_rule: Holy Grail assigns Arthur's current health exactly to 8 when its condition is met; Resurrect returns a defeated fighter with exactly 3 health.
  missing_capability: RECOVER is insufficient for exact assignment, and revive needs an atomic source-defined health result rather than an invalid intermediate health state.
  proposed_generic_contract: Support exact SET_HEALTH for undefeated fighters and atomic RETURN_FIGHTER with an exact-health rule, preserving defeat/recovery trigger distinctions.
  evidence_status: verified

- id: A-REQ-007
  family: boost_pipeline
  severity: P1
  affects:
    - spike/always-surprising
  established_rule: Always Surprising performs a blind BOOST and doubles the resolved BOOST amount when Drusilla is on a Shadow.
  missing_capability: Shared BOOST pipeline cannot transform the resolved blind-BOOST amount between reveal/read and application.
  proposed_generic_contract: BOOST exposes a resolved amount value and supports source-defined transforms before that amount is added to combat value.
  evidence_status: verified

- id: A-REQ-008
  family: movement_targeting_and_combat_legality
  severity: P0
  affects:
    - little-red-riding-hood/what-big-ears-you-have
    - deadpool/xavier-institute-faculty
  established_rule: Legal combat-card/attack mode can depend on source state or the card itself rather than immutable fighter/card printed type alone.
  missing_capability: Shared legality model does not expose declarative conditional card modes or card-specific melee/ranged attack permissions.
  proposed_generic_contract: Preserve immutable printed metadata while deriving legal play/attack modes from scoped declarative permissions evaluated at card commit.
  evidence_status: verified

- id: A-REQ-009
  family: card_zones_and_auxiliary_systems
  severity: P0
  affects:
    - little-red-riding-hood/basket-reference
    - little-red-riding-hood/basket-dependent-cards
  established_rule: Little Red's Basket is not one of the 30 action cards but is a real reference object placed in the ordered discard state and can be the current top object whose symbol determines Basket behavior; wild resolves as exactly one item for a source resolution.
  missing_capability: Standard card-zone runtime assumes ordinary card instances and has no explicit non-action auxiliary object that participates in ordered-zone top-state without affecting deck/card counts.
  proposed_generic_contract: Ordered zones may contain typed auxiliary objects distinct from action-card instances; top-object metadata is queryable and a wild-derived attribute can be bound once per source resolution.
  evidence_status: verified

- id: A-REQ-010
  family: damage_and_health
  severity: P0
  affects:
    - yennenga/shield-of-the-archers-ability
  established_rule: Pending damage to Yennenga may be divided among eligible Archers in her zone, no Archer may receive more than remaining health, and all unallocated damage remains on Yennenga with the original damage identity/type.
  missing_capability: REDIRECT_DAMAGE models recipient replacement, not partial allocation of one pending event across multiple recipients plus residual damage.
  proposed_generic_contract: Pending damage supports an allocation stage with eligible recipients, per-recipient caps, residual original target and preserved source/type/event identity.
  evidence_status: verified

- id: A-REQ-011
  family: battlefield_components_and_paths
  severity: P0
  affects:
    - yennenga/stallion-charge
  established_rule: Stallion Charge damages each opposing fighter whose space Yennenga moved through during that specific movement.
  missing_capability: Shared MOVE result does not expose the ordered traversed-space path for later source stages.
  proposed_generic_contract: MOVE may bind a resolution result containing origin, ordered traversed_spaces and destination; later selectors can query distinct fighter instances whose spaces occurred in that path.
  evidence_status: verified

- id: A-REQ-014
  family: fighter_presence_and_occupancy
  severity: P1
  affects:
    - buffy/setup-selected-sidekick
  established_rule: Buffy chooses Giles or Xander before deck construction; only the chosen sidekick and matching five-card package participate in that match.
  missing_capability: Shared setup/runtime contract does not explicitly guarantee that an unselected fighter definition is absent from active, defeated and targetable match rosters.
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
  proposed_generic_contract: First decide a central adaptation policy per external predicate/action class; only then expose generic EXTERNAL_PREDICATE/EXTERNAL_INSTRUCTION confirmation and persistence primitives.
  evidence_status: verified
```

Requirements count: **11**.

Mirror-safety is **not** a new Worker A requirement: shared `FIGHTER-SCHEMA-001` and `CARD-SCHEMA-001/002` already require definition/instance separation, stable fighter instance identity, immutable owner and runtime card instance identity. Reconciled local definition IDs are interpreted within their fighter/deck manifest scope; runtime references must resolve to player-owned instances under those existing invariants.

## 4. Evidence uncertainties / project normalizations

### Qualified evidence

`deadpool` is the only qualified-evidence fighter:

- first-party Mondo material establishes the official release;
- published UmDb supplies the complete 30-card corpus used here;
- a first-party online full card-text/component dump was not located.

This is a provenance-quality limitation, not an unknown deterministic card fact. Deadpool is blocked by **policy**, not by this evidence limitation.

### Interpretations / normalizations

- `beowulf/the-ancient-heirloom`: `derived`, high confidence. Its two optional DURING COMBAT clauses are staged in printed order; either or both may resolve. The deck manifest carries the replacement condition if a future card-specific authoritative ruling states otherwise.
- Mirror-safe scoping: `authoritative` project schema behavior. Definition IDs are manifest-local; gameplay targets/state use runtime fighter/card instances and owner scope required by the shared schema.
- No fan balance patch, UnPatched record or `/decks/...` community deck is used as canonical evidence.

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

Reconciliation structural checks:

- all 17 fighter manifests and all 17 deck manifests declare `status`, `verification`, structured `sources`, reconciliation `validation`, and `reconciliation_contract: phase-4b-v1`;
- non-trivial reconciled effects use staged resolution with explicit controller/cancellation and structured choices;
- ordinary procedural `REQUEST_CHOICE` / procedural `external_definitions` have been removed;
- `external_definitions` remains only where it is a real gameplay definition (`yennenga/volley`);
- standard `deck`, `hand`, `discard` are not redundantly declared as custom zones;
- requirement-backed operations are tagged with their `A-REQ` instead of free-form per-card blocker prose.

## 6. Source gaps

Blocking source gaps: **0**.

Qualified non-blocking source gap: **1** — Deadpool lacks a located first-party online full card-text/component dump; official release provenance plus published UmDb card corpus remain sufficient to reconstruct deterministic printed semantics.

Unresolved policy items: **1** — `A-REQ-015` Deadpool digital adaptation policy.

## 7. Files owned

Worker A owns and reconciled exactly:

- 17 files under `docs/fighters/phase-4b/` for the assigned IDs;
- 17 files under `docs/cards/phase-4b/` for the assigned IDs;
- `docs/phase-4b/worker-a-report.md`.

Files changed by reconciliation: **35**.

No shared schema, mechanics, rules, rulings, set registry, research-plan, README, or Phase 4A manifest was changed.

## Worker 4B-A Reconciliation Handoff

Branch: `phase-4b-worker-a-classics`  
Exact final Head: **supplied in the delivery response after this report write; a commit cannot self-embed its own resulting SHA**  
Assigned count: **17**  
Verified count: **8**  
Partial count: **8**  
Blocked count: **1**  
Quantity validation: **PASS (17/17)**  
Requirements count: **11**  
Unresolved evidence items: **1 qualified, non-blocking provenance limitation (Deadpool); 0 blocking evidence gaps**  
Unresolved policy items: **1 — Deadpool digital adaptation policy (`A-REQ-015`)**  
Files changed: **35**  
Merge to main: **not performed**
