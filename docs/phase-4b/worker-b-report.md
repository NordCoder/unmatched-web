# Phase 4B Worker B reconciliation report

## 1. Worker identity

- Repository: `NordCoder/unmatched-web`
- Branch: `phase-4b-worker-b-licensed`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: 18 assigned licensed fighters from Worker B
- Reconciliation pre-report head: `5331ca7854b775a277be24ba204eea74f2011170`
- Shared files changed by Worker B: none
- Phase 4A Black Panther changed: no
- Merge to `main`: not performed

The reconciliation contract separates evidence confidence, deterministic gameplay semantics, shared-engine integration capability, and policy status. A missing generic engine capability therefore produces `partial`, not `blocked`, when gameplay behavior is already deterministic.

## 2. Final status matrix

| Fighter | Status | Evidence | Semantics | Integration | Policy |
| --- | --- | --- | --- | --- | --- |
| daredevil | verified | verified | verified | ready | not_applicable |
| elektra | partial | verified | verified | requires_shared_extension | not_applicable |
| bullseye | partial | verified | verified | requires_shared_extension | not_applicable |
| ghost-rider | verified | verified | verified | ready | not_applicable |
| luke-cage | verified | verified | verified | ready | not_applicable |
| moon-knight | verified | verified | verified | ready | not_applicable |
| dr-ellie-sattler | partial | verified | verified | requires_shared_extension | not_applicable |
| t-rex | partial | qualified | qualified | requires_shared_extension | not_applicable |
| houdini | partial | verified | verified | requires_shared_extension | not_applicable |
| genie | verified | verified | verified | ready | not_applicable |
| cloak-and-dagger | verified | verified | verified | ready | not_applicable |
| ms-marvel | partial | verified | verified | requires_shared_extension | not_applicable |
| squirrel-girl | partial | qualified | qualified | requires_shared_extension | not_applicable |
| black-widow | verified | verified | verified | ready | not_applicable |
| winter-soldier | verified | verified | verified | ready | not_applicable |
| doctor-strange | verified | verified | verified | ready | not_applicable |
| she-hulk | partial | verified | verified | requires_shared_extension | not_applicable |
| spider-man | partial | verified | verified | requires_shared_extension | not_applicable |

Totals:
- Verified: **9**
- Partial: **9**
- Blocked: **0**

## 3. B-REQ requirements

### B-REQ-001 — Multi-space fighter footprint

- Family: `fighter_presence_and_occupancy`
- Affected: `t-rex`
- Requirement: represent one fighter occupying an ordered/oriented two-space footprint, including legal movement/rotation, placement, occupied-space queries, attack origins/range, zone/adjacency evaluation, and turn-start footprint history.
- Why shared: the behavior changes generic occupancy, movement, placement, targeting and historical predicates; it cannot be a T. Rex-only state object.

### B-REQ-002 — Small-fighter shared occupancy and propagated damage

- Family: `fighter_presence_and_occupancy`
- Affected: `squirrel-girl`
- Requirement: support small-fighter occupancy class, shared-space compatibility/capacity, pass-through, same-space adjacency, per-token fighter identity, and same-type co-located damage propagation with source/recipient provenance.
- Damage normalization: every Squirrel that receives propagated damage is a damaged fighter and contributes its received amount to total damage from the originating effect; provenance preserves primary versus propagated recipients.

### B-REQ-003 — Off-board without defeat lifecycle

- Family: `fighter_presence_and_occupancy`
- Affected: `elektra`
- Requirement: a fighter may leave battlefield presence while explicitly remaining undefeated, preserve identity/ownership/state off-board, and later return through source-defined resurrection setup.

### B-REQ-004 — Card-used-as-BOOST event and history

- Family: `boost_pipeline`
- Affected: `houdini`
- Requirement: BOOST resolution emits a source-card event with card instance, controller, boosted target/context, disposition and reusable history so BOOSTED WITH effects and later selectors can refer to the exact card used.

### B-REQ-005 — Field-only disclosure of hidden card information

- Family: `search_randomness_and_disclosure`
- Affected: `spider-man`
- Requirement: disclose a selected immutable/effective field such as printed numeric value to a defined audience at attack-commit/pre-defense timing while preserving the hidden card's identity and other fields.

### B-REQ-006 — Positioned battlefield token instances

- Family: `battlefield_components_and_paths`
- Affected: `dr-ellie-sattler`
- Requirement: reusable battlefield token instances/multisets whose location is supply or `space_ref`, including same-space multiplicity, selectors/counts and move/return operations. Required for five reusable Insight tokens.

### B-REQ-007 — Operation-cause provenance

- Family: `identity_history_and_provenance`
- Affected: `houdini`
- Requirement: operation-result events, especially discard, retain affected card instance, owner, actor, causing effect/source and disposition so a later trigger can distinguish opponent-effect discard from costs, cleanup and own effects.

### B-REQ-008 — Graph-distance attack legality

- Family: `movement_targeting_and_combat_legality`
- Affected: `bullseye`, `ms-marvel`
- Requirement: reusable per-fighter attack-legality policy with maximum graph distance and zone bypass, without storing attack rules as ad-hoc runtime state.

### B-REQ-009 — Action-permission spending as a cost

- Family: `resources_actions_and_turn_control`
- Affected: `ms-marvel`, `she-hulk`
- Requirement: consume one legal spendable action permission as a source-defined cost while respecting normal, gained, free and restricted action permissions. This must not be implemented as `actions_remaining -= 1`.

### B-REQ-010 — Private-zone reactive effect cancellation

- Family: `resolution_and_choices`
- Affected: `houdini`
- Requirement: a card remaining in a private hand may react to a pending opponent effect, reveal itself as payment/permission where required, and cancel the bound triggering effect instance as a whole. This preserves `A Magician Never Reveals His Secrets` without merely preventing one nested `LOOK_AT` operation.

Requirements count: **10**.

## 4. Evidence-qualified interpretations

### T. Rex — `Momentous Shift` overlap

- Interpretation id: `t-rex-momentous-shift-overlap`
- Status: `project_normalization`
- Confidence: `medium`
- Deterministic project behavior: if turn-start footprint is `{A, B}` and current footprint is `{B, C}`, the condition is **not** satisfied because T. Rex remains in starting space `B`. The condition becomes true only when current and turn-start footprints share no space.
- Evidence qualification: published sources establish the two-space fighter and card wording, but no captured authoritative exact-case ruling resolves overlapping footprints. Current secondary FAQ gives the opposite interpretation.
- Replacement condition: `replace_if_authoritative_ruling_found`.

### Squirrel Girl — Go Nuts `empty` correction

- Interpretation id: `squirrel-go-nuts-empty-correction`
- Status: `project_normalization`
- Confidence: `medium`
- Deterministic project behavior: Go Nuts may use an adjacent destination permitted by small-fighter occupancy even when compatibly occupied; the printed `empty` restriction is not enforced.
- Evidence qualification: official Teen Spirit rules contain `empty`; current reference/index wording omits it and secondary FAQ calls the printed restriction an error. An exact authoritative correcting ruling has not been captured.
- Replacement condition: `replace_if_authoritative_ruling_found`.

### Squirrel Girl — propagated damage counting

- Interpretation id: `squirrel-propagated-damage-counting`
- Status: `project_normalization`
- Confidence: `high`
- Official evidence: an official ruling confirms that each same-type small fighter in the shared space takes equal propagated damage; four-small-fighter capacity and shared-space adjacency are also confirmed.
- Deterministic project behavior: all recipients that actually take propagated damage count as damaged fighters and contribute to total damage caused by the originating effect; provenance distinguishes the primary target from propagated recipients.
- Replacement condition: `replace_if_authoritative_ruling_found` for any future ruling that specifically changes downstream attribution/counting.

Evidence-qualified fighter pairs: **2** (`t-rex`, `squirrel-girl`).
Project-normalization records: **3**.

## 5. Validation matrix

| Fighter | Quantity | Pool / game deck | Usable-by | References | Source coverage | Semantic structure | Integration | Fan content |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| daredevil | PASS | 22 / 22 | PASS | PASS | verified | PASS | ready | PASS |
| elektra | PASS | 20 / 20 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |
| bullseye | PASS | 30 / 30 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |
| ghost-rider | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| luke-cage | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| moon-knight | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| dr-ellie-sattler | PASS | 30 / 30 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |
| t-rex | PASS | 30 / 30 | PASS | PASS | qualified | PASS | requires_shared_extension | PASS |
| houdini | PASS | 30 / 30 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |
| genie | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| cloak-and-dagger | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| ms-marvel | PASS | 30 / 30 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |
| squirrel-girl | PASS | 30 / 30 | PASS | PASS | qualified | PASS | requires_shared_extension | PASS |
| black-widow | PASS | 31 / 31 | PASS | PASS | verified | PASS | ready | PASS |
| winter-soldier | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| doctor-strange | PASS | 30 / 30 | PASS | PASS | verified | PASS | ready | PASS |
| she-hulk | PASS | 30 / 30 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |
| spider-man | PASS | 30 / 30 | PASS | PASS | verified | PASS | requires_shared_extension | PASS |

Quantity reconciliation: **18/18 PASS**.

Structural reconciliation checks represented by this matrix are documentation-contract checks, not a claim that a repository-side YAML parser or CI suite was executed during the Connector-only reconciliation pass.

## 6. Source gaps

No source gap prevents deterministic Worker B implementation.

Two narrow exact-case authority gaps remain qualified rather than blocked because explicit reversible project behavior exists:
1. T. Rex `Momentous Shift` when the current two-space footprint overlaps the turn-start footprint.
2. Squirrel Girl Go Nuts removal of the printed `empty` restriction.

There are no unresolved deck-list, card-quantity, ordinary card-text, or policy blockers in Worker B scope. No fan `/decks/...` record is used as canonical evidence.

## 7. Files

Worker-owned corpus: **37 files**.
- 18 fighter manifests under `docs/fighters/phase-4b/`.
- 18 card manifests under `docs/cards/phase-4b/`.
- `docs/phase-4b/worker-b-report.md`.

Scope comparison against Authorized Base contains only these 37 Worker B artifacts. Shared schema/mechanics/rules/registry files remain untouched. Black Panther remains untouched.

## Worker 4B-B Reconciliation Handoff

- Branch: `phase-4b-worker-b-licensed`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Exact final Head: emitted by the external GitHub handoff after this report commit; a Git commit cannot embed its own SHA in its tracked contents without changing that SHA.
- Assigned: **18**
- Verified: **9** — `daredevil`, `ghost-rider`, `luke-cage`, `moon-knight`, `genie`, `cloak-and-dagger`, `black-widow`, `winter-soldier`, `doctor-strange`
- Partial: **9** — `elektra`, `bullseye`, `dr-ellie-sattler`, `t-rex`, `houdini`, `ms-marvel`, `squirrel-girl`, `she-hulk`, `spider-man`
- Blocked: **0**
- Quantity validation: **18/18 PASS**
- Requirements count: **10**
- Evidence-qualified fighter pairs: **2** (`t-rex`, `squirrel-girl`)
- Project-normalization items: **3**
- Policy blockers: **0**
- Files changed in Worker B scope: **37**
- Shared files changed: **0**
- Black Panther changed: **no**
- Merge to `main`: **not performed**
