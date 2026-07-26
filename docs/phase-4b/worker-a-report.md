# Phase 4B Worker A research report

Scope: early/classic + retired licensed fighters assigned to Worker A.

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Base verification: branch was identical to Authorized Base before Worker A writes (`ahead=0`, `behind=0`).

## Source discipline

Research used the repository source hierarchy:

1. official Restoration Games / IELLO product, rulebook and set-rule material where publicly available;
2. official/current rulings material and the current Rulings Archive/index where applicable;
3. published UmDb `/umdb/...` records for normalized card/deck facts;
4. secondary indexes only for cross-checking or retired-product discovery.

Community/fan `/decks/...` pages were explicitly excluded from canonical provenance. Search results for balance patches that reused published card names but changed values, quantities or effects were rejected rather than reconciled into the official corpus.

## Fighter results

| Fighter | Result | Primary note |
| --- | --- | --- |
| `alice` | verified | Public Big/Small state and size-dependent static modifiers fit the existing model. |
| `king-arthur` | blocked | `The Holy Grail` assigns an exact health value; current shared model has recovery but no absolute health assignment. |
| `medusa` | verified | Multi-Harpy topology, gaze trigger and deck normalize without a new primitive. |
| `sinbad` | verified | Voyage tag + discard-zone counting fits existing card metadata/state queries. |
| `robin-hood` | verified | Outlaw topology and post-attack movement fit the existing model. |
| `bigfoot` | verified | End-turn zone condition and Jackalope effects fit the existing model. |
| `robert-muldoon` | blocked | Positioned traps require token placement/return plus movement-entry interruption. |
| `invisible-man` | blocked | Positioned Fog semantics and `Vanish` temporary undefeated off-board lifetime need generic support. |
| `jekyll-and-hyde` | verified | Public form state and form-tagged cards fit the existing model. |
| `buffy` | blocked | 35-card published pool / 30-card constructed deck is representable, but `Right-hand Man` has conditional alternate Buffy play permission not expressible by static `usable_by`. |
| `willow` | blocked | Dark Willow state fits; `Resurrect` requires exact health assignment after returning a defeated fighter. |
| `spike` | blocked | Positioned Shadow tokens and conditional blind-BOOST multiplication need generic support. |
| `angel` | verified | Angel/Faith topology and losing-attack ability fit the existing model. |
| `little-red-riding-hood` | blocked | Basket/discard-symbol state is representable, but `What Big Ears You Have` conditionally changes legal card mode/play permission. |
| `beowulf` | verified | Rage counter and per-damage-effect gain semantics fit existing resources/events. |
| `deadpool` | blocked | Published deck intentionally uses external-world/physical predicates and a per-card ranged/melee override that the shared digital model does not define. |
| `yennenga` | blocked | VOLLEY fits the bonus-attack composite; `Stallion Charge` needs the path of one specific movement to be queryable after that move. |

Verified: **8** (`alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `angel`, `beowulf`).  
Blocked: **9** (`king-arthur`, `robert-muldoon`, `invisible-man`, `buffy`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`).

No fighter was marked verified merely to complete the batch; `blocked` means card/topology facts were researched but at least one published semantic cannot be represented faithfully by the current shared schema/effect vocabulary.

## Quantity validation

| Fighter | Published pool | Game deck | Reconciliation |
| --- | ---: | ---: | --- |
| Alice | 30 | 30 | PASS |
| King Arthur | 30 | 30 | PASS |
| Medusa | 30 | 30 | PASS |
| Sinbad | 30 | 30 | PASS |
| Robin Hood | 30 | 30 | PASS |
| Bigfoot | 30 | 30 | PASS |
| Robert Muldoon | 30 | 30 | PASS |
| Invisible Man | 30 | 30 | PASS |
| Jekyll & Hyde | 30 | 30 | PASS |
| Buffy | 35 | 30 | PASS: 25 base + selected 5-card Giles/Xander group |
| Willow | 30 | 30 | PASS |
| Spike | 30 | 30 | PASS |
| Angel | 30 | 30 | PASS |
| Little Red Riding Hood | 30 | 30 | PASS; Basket is a non-action setup/reference component |
| Beowulf | 30 | 30 | PASS |
| Deadpool | 30 | 30 | PASS: 30 distinct one-copy cards |
| Yennenga | 30 | 30 | PASS |

**Quantity validation: PASS for all 17 assigned fighters.**

`usable_by` was checked against each fighter topology. Conditional alternate permissions are not silently widened: Buffy's `Right-hand Man` remains Xander-owned with the conditional Buffy permission recorded separately as a blocker. Resources/zones/state references are defined in fighter/deck manifests or explicitly identified below as generic-model blockers.

## Generic schema/effect extension proposals

### 1. Exact health assignment

Affected: King Arthur — `The Holy Grail`; Willow — `Resurrect`.  
Source: published official deck semantics / published UmDb cross-check.  
Problem: `RECOVER` adds health subject to normal recovery semantics; these effects require assigning current health to an exact source value.  
Proposed generic model: `SET_HEALTH(target, value)` (or equivalent exact-health operation), explicitly distinct from recovery/damage.  
Blocker: **yes** for King Arthur and Willow.

### 2. Positioned special components/tokens

Affected: Robert Muldoon traps, Invisible Man Fog, Spike Shadows.  
Source: official Jurassic Park, Cobble & Fog, and Buffy set rules plus published deck records.  
Problem: resources can count tokens but the shared effect model lacks generic board-position lifecycle semantics.  
Proposed generic model: position-sensitive component instances with generic `PLACE_TOKEN`, `MOVE_TOKEN`, `RETURN_TOKEN_TO_SUPPLY`, token-space selectors/predicates and token ownership.  
Blocker: **yes** for Muldoon, Invisible Man and Spike.

### 3. Movement-entry interruption

Affected: Robert Muldoon traps.  
Source: official InGen vs Raptors rules.  
Problem: an opposing fighter entering a trap space must stop the currently resolving movement before damage/return/draw resolution.  
Proposed generic model: movement-step event with a protected `STOP_CURRENT_MOVEMENT` control operation usable by position-triggered effects.  
Blocker: **yes** for Muldoon.

### 4. Temporary undefeated off-board fighter state

Affected: Invisible Man — `Vanish`.  
Source: published deck/rulings.  
Problem: the hero is temporarily absent from the board without being defeated, then is placed at a later turn boundary.  
Proposed generic model: `off_board` undefeated fighter location/state with source-owned return timing and legal placement domain.  
Blocker: **yes** for Invisible Man.

### 5. Conditional card user/type/attack-mode permission

Affected: Buffy — `Right-hand Man`; Little Red Riding Hood — `What Big Ears You Have`; Deadpool — `Xavier Institute Faculty`.  
Source: published card semantics.  
Problem: static `usable_by` and static card type/attack mode cannot represent conditional alternate user permission, conditional Attack/Defense legality, or a card-specific ranged/melee override.  
Proposed generic model: declarative play-permission clauses evaluated at commit time, able to vary eligible fighter and legal combat mode/type without mutating immutable printed metadata.  
Blocker: **yes** for Buffy, Little Red and Deadpool.

### 6. Blind BOOST result transformation

Affected: Spike — `Always Surprising`.  
Source: published deck plus retired-set cross-check.  
Problem: a blind top-deck BOOST result is conditionally multiplied before being applied. Existing `BOOST`/`ADD_BOOST_VALUE` does not expose a reusable transform layer for that resolved BOOST amount.  
Proposed generic model: transform hook over a resolved BOOST-value event, e.g. multiply/set the pending BOOST amount before application.  
Blocker: **yes** for Spike.

### 7. Movement-path capture/query

Affected: Yennenga — `Stallion Charge`.  
Source: Battle of Legends Vol. 2 published deck.  
Problem: the post-move effect targets opposing fighters whose spaces Yennenga moved through during that specific movement. Destination/current adjacency is insufficient.  
Proposed generic model: movement resolution records an ordered traversed-space path and exposes selectors over fighters/spaces crossed by that move.  
Blocker: **yes** for Yennenga.

### 8. External-world / physical-component predicates

Affected: Deadpool — multiple cards and fighter ability.  
Source: published Deadpool deck.  
Problem: official semantics intentionally depend on player real names, physical sleeves/card writing, clothing, food/drink, noises/actions, set ownership, a mirror, post-game wager, and subjective board colour. These are neither deterministic game-state predicates nor fan rules.  
Proposed generic model: explicit `external_predicate` / `external_action` compatibility layer with online-platform policy (supported, user-confirmed, substituted, or disabled) chosen centrally; never silently rewrite these effects into ordinary board-state logic.  
Blocker: **yes** for Deadpool.

## Rulings / ambiguity / source gaps

- King Arthur: a secondary transcription conflicts with published UmDb on Feint/Regroup BOOST. Published UmDb and deck-structure cross-check support BOOST 1; the manifests use BOOST 1.
- Buffy licensed set: official set rules establish fighter/setup mechanics, while individual retired card facts are primarily recoverable from published UmDb and current archival indexes. Those records were cross-checked; community rebalances were rejected.
- Spike: retired-set card text was cross-checked against current archival/secondary transcriptions. Shadow mechanics remain blocked on generic position semantics rather than encoded as custom Spike logic.
- Deadpool: external/physical joke mechanics are genuine published semantics, not fan-patch data. They remain explicit blockers for a deterministic web engine.
- Little Red: the Basket is not counted as an action-deck card; it is retained as a setup/reference component tied to top-discard symbol semantics.
- No unresolved source gap prevents quantity reconciliation for any assigned fighter.

## Files created

Fighter manifests (17):

`docs/fighters/phase-4b/alice.yaml`, `king-arthur.yaml`, `medusa.yaml`, `sinbad.yaml`, `robin-hood.yaml`, `bigfoot.yaml`, `robert-muldoon.yaml`, `invisible-man.yaml`, `jekyll-and-hyde.yaml`, `buffy.yaml`, `willow.yaml`, `spike.yaml`, `angel.yaml`, `little-red-riding-hood.yaml`, `beowulf.yaml`, `deadpool.yaml`, `yennenga.yaml`.

Deck manifests (17):

`docs/cards/phase-4b/alice.yaml`, `king-arthur.yaml`, `medusa.yaml`, `sinbad.yaml`, `robin-hood.yaml`, `bigfoot.yaml`, `robert-muldoon.yaml`, `invisible-man.yaml`, `jekyll-and-hyde.yaml`, `buffy.yaml`, `willow.yaml`, `spike.yaml`, `angel.yaml`, `little-red-riding-hood.yaml`, `beowulf.yaml`, `deadpool.yaml`, `yennenga.yaml`.

Report: `docs/phase-4b/worker-a-report.md`.

Pre-report Base→Head scope audit found exactly the 34 assigned fighter/deck paths, all added; no shared schema, mechanics, rules, set registry, ambiguity register, research plan or README file was changed.

## Worker 4B-A Handoff

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head immediately before this report commit: `bcb18cd962e58c100e7b9038057c1682f56248f3`  
Exact final Head: the commit containing this report cannot cryptographically embed its own Git SHA; the immutable exact final branch Head must be read after this final write and is reported in the Worker A delivery response. No further repository writes are permitted after that read.  
Assigned fighters: **17**  
Verified: **8** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `angel`, `beowulf`  
Blocked: **9** — `king-arthur`, `robert-muldoon`, `invisible-man`, `buffy`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`  
Quantity validation: **PASS (17/17)**  
Schema-extension proposals: **8 generic extension families**, listed above; no shared schema file modified  
New ambiguity/blockers: exact health assignment; positioned special tokens; movement interruption; temporary undefeated off-board state; conditional play permissions/card modes; blind BOOST transform; movement-path capture; external-world semantics  
Source gaps: **none blocking deck quantity/card-identity reconciliation**; retired licensed material required published archival/secondary cross-checks as documented above  
Files created: **35** (17 fighter manifests + 17 deck manifests + this report)  
Merge to main: **not performed**
