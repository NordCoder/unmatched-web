# Phase 4B Worker A — canonical handoff

Scope: early/classic + retired licensed fighters assigned to Worker A.

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Base verification: branch was identical to Authorized Base before Worker A writes (`ahead=0`, `behind=0`).

## Source discipline

Canonical evidence order:

1. official Restoration Games / IELLO / Mondo product, rulebook, release, errata and set-rule material;
2. current Unmatched Reference and publisher-backed Rulings Archive;
3. published UmDb `/umdb/...` records for official deck/card metadata and normalized card text;
4. secondary indexes only for discovery/cross-checking.

Community/fan `/decks/...`, reskins, `UnPatched` records and balance patches are excluded from canonical provenance.

Current Reference:
- https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf

Official Rulings Archive:
- https://docs.google.com/document/d/13b-FbPq_vuqcc3IokeHvQ2ctJaDNZZuUaZmt4uft5h0/

## Validation summary

- Assigned fighters: **17**
- Fighter manifests: **17**
- Deck manifests: **17**
- Quantity validation: **PASS (17/17)**
- Current UmDb quantity/type reconciliation: **PASS (17/17)**
- Blocking evidence gaps: **0**
- Fan-patch contamination: **none found**

`verified` in this Worker A corpus means **published content/evidence and the fighter-specific normalized representation are verified to the level recorded here**. It does **not** mean the current shared DSL/runtime is developer-ready for every referenced selector, binding, batch operation or cross-fighter interaction. Those integration-contract gaps are listed separately below and belong to the Phase 4B orchestrator/shared-model pass.

## Fighter results

| Fighter | Result | Primary note |
| --- | --- | --- |
| `alice` | verified | Big/Small state and 30-card deck reconcile. |
| `king-arthur` | blocked | `The Holy Grail` requires exact health assignment. |
| `medusa` | verified | Multi-Harpy topology and deck reconcile. |
| `sinbad` | verified | Voyage/discard-state semantics sourced and normalized. |
| `robin-hood` | verified | Outlaw topology, Hit and Run and multiplayer `Steal From the Rich` flow recorded. |
| `bigfoot` | verified | End-turn zone condition and optional draw recorded. |
| `robert-muldoon` | blocked | Positioned trap lifecycle plus movement interruption need shared support. |
| `invisible-man` | blocked | Fog positioning/dynamic movement adjacency and `Vanish` off-board lifetime need shared support. |
| `jekyll-and-hyde` | verified | Form state and card legality are sourced; choice dependencies are now explicit in the manifest. |
| `buffy` | verified | 35→30 construction, setup stages and Xander-only `Right-hand Man` semantics are resolved. |
| `willow` | blocked | `Resurrect` requires exact post-return health assignment. |
| `spike` | blocked | Positioned Shadows and blind-BOOST result transformation need shared support. |
| `angel` | verified | Angel/Faith topology and card distribution reconcile. |
| `little-red-riding-hood` | blocked | `What Big Ears You Have` conditionally changes legal combat-card mode. |
| `beowulf` | verified | Rage is represented as one canonical resource counter. |
| `deadpool` | blocked | Dynamic attack mode plus digital-adaptation policy are unresolved. |
| `yennenga` | blocked | Movement-path capture and partial pending-damage allocation need shared support. |

Verified: **9** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `buffy`, `angel`, `beowulf`.  
Blocked: **8** — `king-arthur`, `robert-muldoon`, `invisible-man`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`.

## Normalization cleanup completed

The post-research logic audit produced the following manifest-level corrections that were safe to make without editing shared semantic files.

### Regroup replacement semantics

All **11 Worker A decks containing published `Regroup`** now encode mutually exclusive combat-result branches:

- won combat → `DRAW 2`;
- did not win → `DRAW 1`.

The previous `DRAW 1` plus a second conditional `DRAW 1` representation was removed because the published card says `draw 2 instead`; separate draw effects can change exhaustion/event accounting even when the normal hand result happens to be the same.

Affected: Alice, King Arthur, Medusa, Sinbad, Robin Hood, Bigfoot, InGen, Buffy, Willow, Spike, Angel.

### Buffy

- `Cartwheel Kick`: removed effect-level optionality. `MOVE up to 3` already permits zero movement; the following adjacent damage still resolves.
- `Insight`: explicitly chooses an opponent before looking at and selecting a card from that opponent's hand.
- setup is split into canonical `CHARACTER_CONFIGURATION` and `DECK_CONSTRUCTION` hooks.
- `Right-hand Man` remains Xander-only; adjacency to Buffy changes only its combat value.

### Robin Hood

`Steal From the Rich` now records the full multiplayer flow:

1. draw 1;
2. owner chooses an opponent;
3. chosen opponent privately chooses whether to discard one card;
4. if they do not discard, owner draws one more.

### Willow

`Hacker` now records the actual top-or-bottom disposition after the private look/choice rather than ending after `REQUEST_CHOICE`.

### Jekyll & Hyde

- `Strange Case` now declares the adjacent-fighter choice before using `chosen_fighter`.
- `Calming Research` now declares the 0..3 draw-count choice and the subsequent keep-one choice over cards drawn by that effect. Shared integration still needs a typed result-binding convention for identifiers such as `chosen_draw_count` and `cards_drawn_by_this_effect`.

### Setup hook normalization

Worker-local setup stages now use the canonical shared hook names where applicable:

- Alice → `POST_PLACEMENT_CONFIGURATION`;
- Invisible Man → `POST_PLACEMENT_CONFIGURATION`;
- Buffy → `CHARACTER_CONFIGURATION` + `DECK_CONSTRUCTION`.

### Duplicate mutable state cleanup

Removed duplicate counters that represented the same gameplay resource twice:

- Beowulf Rage now lives only in `resources.rage`;
- Muldoon trap availability is derived from trap-token instance lifecycle, not a separate `traps_available` integer;
- Spike Shadow availability is derived from Shadow-token instance lifecycle, not a separate `shadows_available` integer.

Muldoon traps explicitly require lifecycle `available → board → box`; traps moved to the box are not reusable.

### Yennenga damage allocation

The previous generic `REDIRECT_DAMAGE` operation was removed from `Shield of the Archers` because it cannot faithfully represent the published rule. Yennenga needs partial allocation of one pending damage event across zero or more eligible Archers, capped by each Archer's remaining health, with the unallocated remainder still applied to Yennenga without changing the original damage type.

## 1. Evidence blockers

**None.**

No assigned fighter remains blocked because a required published rule/card fact is unknown or unconfirmed.

### Non-blocking provenance limitation

**Deadpool:** official Mondo material establishes first-party release provenance and published UmDb supplies the complete 30-card corpus. A first-party online full card-text/rulebook dump was not located. This is weaker primary card-level provenance than ordinary Restoration rulebook sets, but it does not block deck/card reconstruction.

## 2. Engine blockers

### Exact health semantics

Affected:
- King Arthur — `The Holy Grail`;
- Willow — `Resurrect`.

Need two related but distinct semantics:

- exact health assignment on an undefeated fighter;
- atomic return/revive with source-defined exact health, so a resurrected fighter does not pass through an invalid intermediate health state.

### Positioned special components

Affected:
- Muldoon — traps;
- Invisible Man — Fog;
- Spike — Shadows.

Need reusable token/component instances with identity, owner/controller, lifecycle/location, space, placement/movement rules, occupancy/stacking policy and selectors.

Invisible Man additionally needs scoped dynamic adjacency: Fog spaces are mutually adjacent for Invisible Man's **movement only**, not for attacks or generic adjacency.

### Movement-entry interruption

Affected: Muldoon traps.

Need movement/placement entry events plus protected interruption semantics capable of stopping a move immediately when a trap is entered. Official rules also trigger traps when an opposing fighter is placed into the trap space by another effect/player.

### Temporary undefeated off-board / dormancy

Affected: Invisible Man — `Vanish`.

Need undefeated `off_board` state, source-owned return timing and a deterministic interaction between the return hook and dormant-player turn handling.

### Conditional combat-card / attack mode

Affected:
- Little Red Riding Hood — `What Big Ears You Have`;
- Deadpool — `Xavier Institute Faculty`.

Need declarative play-mode/attack-mode permissions without mutating immutable printed metadata.

### Blind BOOST result transformation

Affected: Spike — `Always Surprising`.

Need a transform hook over the resolved blind-BOOST amount before it is applied.

### Movement-path capture/query

Affected: Yennenga — `Stallion Charge`.

Need an ordered traversed-space path for a specific movement resolution and selectors over fighters/spaces crossed by that move.

### Partial pending-damage allocation

Affected: Yennenga — `Shield of the Archers`.

Need allocation of one pending damage event among multiple recipients with per-recipient caps, residual damage on the original target, and preservation of the original damage type/event identity.

## 3. Other blockers

### Deadpool digital-adaptation policy

Published Deadpool effects depend on external/physical facts such as player names, food, spoken/noise actions, set ownership, mirror presence, sleeve/card-writing state, clothing, wagers and subjective board colour.

The project needs one central policy deciding whether each external predicate/action is:

- automatically supported;
- manually confirmed by players;
- replaced with a documented digital equivalent; or
- unsupported/disabled for online play.

Only after that policy is chosen should the shared engine expose generic external-predicate/player-confirmation primitives.

## Cross-cutting integration-contract gaps

These are **not evidence gaps** and are intentionally not patched by Worker A because they belong in shared schema/mechanics files owned by the Phase 4B orchestrator.

### Runtime identity / namespaces

Definition IDs such as `feint`, state keys such as `form`, and fighter IDs such as `alice` are only locally unambiguous. Integration needs explicit scoping/runtime identity so mirror matches and different decks using the same card title cannot collide.

At minimum the shared runtime should distinguish:

- fighter definition vs fighter instance;
- card definition vs card instance;
- owner/player scope;
- fighter-local state/resource scope.

### Typed selectors, expressions and result binding

Worker manifests necessarily use operands such as:

- `discarded_card_boost`;
- `combat_damage_taken`;
- `count_other_voyage_cards_in_discard`;
- `chosen_draw_count`;
- `cards_drawn_by_this_effect`;
- `source_found_zone`.

The shared model needs a typed convention for operation outputs, bindings, expressions and selectors rather than relying on free-form strings interpreted by implementation code.

### Choice contract / hidden information

The shared model needs a normative shape for chooser, choice domain, visibility, output binding, reconnect state and dependent consequences. Worker A now makes the most important multiplayer choices explicit, but the generic serialization contract still needs to be promoted centrally.

### Batch / multi-target ordering

Operations such as `MOVE each_friendly_fighter` and damage to multiple fighters need shared rules for:

- sequential vs simultaneous resolution;
- controller-selected order where applicable;
- transition/trigger generation between elements;
- state re-evaluation after each resolved element.

This becomes critical once traps, transformation-on-damage and other cross-fighter triggers compose.

### Setup-selected active roster

Buffy chooses exactly one sidekick. Integration must make it impossible for generic selectors to treat the unselected Giles/Xander definition as an active, defeated or targetable fighter for the match.

### Deterministic randomness / replay

Random discard and similar effects need server-authoritative random resolution recorded in match history so reconnect/replay does not re-roll an already resolved random event.

### Effect defaults and control-flow vocabulary

Shared schema should explicitly define defaults for omitted effect fields and either promote or replace manifest constructs such as `branches`, `normalization_blocker`, dynamic destination/result identifiers and source-defined composite metadata.

### Cross-fighter deterministic fixtures

Before Phase 4B becomes developer-ready, integration should add fixtures for at least:

- Muldoon trap × effect movement/placement;
- Muldoon trap × Yennenga movement path;
- Fog movement adjacency × traps;
- Yennenga partial allocation × multi-target/simultaneous damage;
- Invisible Man `Vanish` × dormant turn/return timing;
- Little Red Basket × discard manipulation;
- Buffy selectable sidekick × generic fighter selectors;
- `Regroup` × near-empty deck/exhaustion;
- mirror matches to prove ID/state isolation.

## Files and scope

Worker A owns exactly **35 files**:

- 17 fighter manifests under `docs/fighters/phase-4b/`;
- 17 deck manifests under `docs/cards/phase-4b/`;
- this report.

No shared schema, mechanics, rules, set registry, ambiguity register, research plan or README file is modified by this branch.

## Worker 4B-A Handoff

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head immediately before this report consolidation: `360d39d696dc9fedc7d6ba1f9b89e895a3c00f5c`  
Exact final Head: read after this final report write and supplied in the delivery response.  
Assigned fighters: **17**  
Verified: **9**  
Blocked: **8**  
Quantity validation: **PASS (17/17)**  
Current UmDb structure/type validation: **PASS (17/17)**  
Blocking evidence gaps: **0**  
Engine blocker families: **8**  
Other policy blockers: **1** — Deadpool digital adaptation  
Cross-cutting integration-contract gaps: **documented above; shared files intentionally untouched**  
Files in Worker A scope: **35**  
Merge to main: **not performed**
