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

`verified` in this Worker A corpus means **published content/evidence and fighter-specific normalization are verified to the level recorded here**. It does **not** mean the current shared DSL/runtime is developer-ready for every selector, binding, batch operation or cross-fighter interaction. Shared integration-contract gaps are listed separately below.

## Fighter results

| Fighter | Result | Primary note |
| --- | --- | --- |
| `alice` | verified | Big/Small state and 30-card deck reconcile. |
| `king-arthur` | blocked | `The Holy Grail` requires exact health assignment; `Prophecy` choice dependency is now explicit. |
| `medusa` | verified | Multi-Harpy topology and deck reconcile. |
| `sinbad` | verified | Voyage/discard-state semantics sourced and normalized. |
| `robin-hood` | verified | Outlaw topology, Hit and Run and multiplayer `Steal From the Rich` flow recorded. |
| `bigfoot` | verified | End-turn zone condition and optional draw recorded. |
| `robert-muldoon` | blocked | Positioned trap lifecycle plus movement-entry interruption need shared support; trap lifecycle is consistently `available → board → box`. |
| `invisible-man` | blocked | Fog positioning/dynamic movement adjacency and `Vanish` undefeated off-board lifetime need shared support; ordered `Vanish` resolution is now explicit. |
| `jekyll-and-hyde` | verified | Form state and card legality are sourced; choice dependencies are explicit. |
| `buffy` | verified | 35→30 construction, setup stages, Xander-only `Right-hand Man`, `Cartwheel Kick`, and all-or-none `Slayer's Strength` semantics are resolved. |
| `willow` | blocked | `Resurrect` requires exact post-return health assignment; `Love and Loss` printed order is now preserved. |
| `spike` | blocked | Positioned Shadows and conditional blind-BOOST multiplication need shared support. |
| `angel` | verified | Angel/Faith topology and card distribution reconcile. |
| `little-red-riding-hood` | blocked | Conditional combat-card mode plus Basket-as-zone-component support are missing from the shared model. |
| `beowulf` | verified | Rage is one canonical resource; `Ancient Heirloom` printed clause order is explicit. |
| `deadpool` | blocked | Dynamic attack mode plus digital-adaptation policy are unresolved; ordinary deterministic card semantics have been hardened. |
| `yennenga` | blocked | `Stallion Charge` path capture and partial pending-damage allocation need shared support; `Surprise Volley`/`Shield Formation` are corrected. |

Verified: **9** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `buffy`, `angel`, `beowulf`.  
Blocked: **8** — `king-arthur`, `robert-muldoon`, `invisible-man`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`.

# Completed normalization corrections

## Source-hardening corrections

The source-hardening pass corrected:

- Alice `Manxome Foe` discard semantics;
- Alice `Claws That Catch` type;
- Bigfoot optional end-turn draw;
- Buffy setup order;
- Spike Shadow reuse when all three are already on board;
- Yennenga VOLLEY defeat condition;
- Angel `Wisdom of Ages` type;
- Medusa `Snipe` type;
- Buffy `Right-hand Man`: Xander-only; adjacency changes value only.

No correction came from fan balance patches.

## First logic-audit corrections

### Regroup replacement semantics

All **11 Worker A decks containing published `Regroup`** use mutually exclusive branches:

- won combat → `DRAW 2`;
- did not win → `DRAW 1`.

This preserves the published `draw 2 instead` semantics rather than representing it as two independent draw effects.

Affected: Alice, King Arthur, Medusa, Sinbad, Robin Hood, Bigfoot, InGen, Buffy, Willow, Spike, Angel.

### Buffy

- `Cartwheel Kick`: only movement is optional through `up to 3`; adjacent damage still resolves.
- `Insight`: explicitly chooses an opponent before hand inspection/discard.
- setup uses `CHARACTER_CONFIGURATION` and `DECK_CONSTRUCTION`.

### Robin Hood

`Steal From the Rich` records the multiplayer flow: draw 1 → choose opponent → chosen opponent privately discards or declines → owner draws one more if they decline.

### Willow

`Hacker` records the actual top-or-bottom disposition after the private look/choice.

### Jekyll & Hyde

- `Strange Case` declares its adjacent-fighter choice.
- `Calming Research` declares draw-count and keep-one choices.

### Duplicate mutable state cleanup

- Beowulf Rage exists only as `resources.rage`.
- Muldoon trap availability derives from token lifecycle.
- Spike Shadow availability derives from token lifecycle.

### Yennenga damage allocation

The previous plain `REDIRECT_DAMAGE` approximation was removed. Yennenga requires partial allocation of one pending damage event across eligible Archers, with per-Archer health caps and residual damage staying on Yennenga without changing damage type/event identity.

# Second logic-audit corrections

## Yennenga

### `Surprise Volley`

Corrected from a false AFTER COMBAT/VOLLEY interpretation to its published IMMEDIATELY semantics:

- optionally return a defeated Archer to the opposing fighter's zone;
- if returned, that Archer becomes the attacker in the **current combat**;
- if no Archer is returned, gain 1 action;
- absence of a defeated Archer must not create an empty mandatory choice that blocks the gain-action branch.

### `Shield Formation`

Corrected card identity and timing:

- canonical name: `Shield Formation`;
- timing: `IMMEDIATELY`;
- combat opponent may discard one card;
- only if they do not discard does Yennenga's owner attempt to return a defeated Archer to Yennenga's zone.

### Multiplayer relation

`Pin the Prey` now scopes its discard to the current combat opponent rather than an arbitrary opponent.

## Buffy — `Slayer's Strength`

The manifest no longer allows an arbitrary subset of adjacent fighters.

Published semantics are represented as one optional all-or-none effect choice:

- if used, attempt to move **all** fighters adjacent to Buffy;
- then deal 1 damage only to fighters actually moved.

## Willow — `Love and Loss`

The card is now one ordered Scheme resolution:

1. draw 2;
2. if Dark Willow, discard the top 2 cards;
3. choose a sidekick in Willow's zone;
4. deal 3 damage to that sidekick.

The DARK clause is no longer an independently reorderable same-timing Scheme effect.

## Invisible Man

### `Vanish`

`Vanish` is now one ordered resolution:

1. recover 1;
2. remove Invisible Man without defeating him;
3. schedule source-owned return at the start of the next turn;
4. if played as action 1, end the turn **after** those preceding steps.

The off-board/return primitives remain engine-blocked, but `END_TURN` can no longer preempt the earlier printed steps through same-timing reordering.

### Choice closure

Added explicit definitions for:

- `Coded Notes` choose-two-and-order dependency;
- `Confound` opponent-discard-or-owner-Fog-movement dependency.

## Deadpool

### `Excuse me while I grow some limbs.`

Corrected to replacement-style branches:

- won → combat opponent discards 2;
- did not win → combat opponent discards 1.

It is no longer represented as two independent discard effects.

### `They Have An Amazing Buffet`

Preserves printed order:

1. recover 2;
2. then test whether Deadpool is at full health;
3. if yes, Deadpool takes 2 damage.

### Multiplayer relation

Combat-card uses of `opponent` are explicitly scoped to the current combat opponent where relevant, distinct from Scheme effects that intentionally choose any opponent.

## Beowulf — `The Ancient Heirloom`

The two DURING COMBAT clauses are represented as one ordered source effect rather than independently reorderable same-timing effects:

1. optionally spend 2 Rage to make value 5;
2. optionally spend 1 Rage to BOOST;
3. either or both clauses may be used.

`Golden Drinking Horn` also records precommitted distinct choices plus owner-selected resolution order; shared choice serialization still needs to formalize this centrally.

## Robert Muldoon

Removed stale `supply` terminology from `Remote Detonation` and the proposed component semantics.

Canonical trap lifecycle is consistently:

`available → board → box`

A returned/triggered trap is not reusable from the box. `Remote Detonation` returns its selected trap to the box and participates in Muldoon's draw-on-trap-return ability.

## King Arthur — `Prophecy`

Closed the dangling `prophecy_keep_two` reference:

- look at top 4;
- choose exactly 2 for hand;
- leave the other 2 on top;
- owner chooses the order of those remaining 2.

## Little Red Riding Hood

### Setup hook

The old ad-hoc `before_shuffle` stage was replaced with canonical `BEFORE_STARTING_HAND` for the start-of-game Basket initialization.

### Wild Basket binding

A wild Basket now binds to exactly **one** effective item for the current source resolution. It must not independently satisfy multiple Basket-item branches in the same resolution.

`Never Leave the Path` therefore resolves exactly one of its wolf / knife / pelt branches when the Basket is wild.

### Basket runtime gap

The Basket is not an action card but physically occupies the ordered discard zone and can be the top discard object. Worker A records this source fact but does not invent a shared runtime primitive for it.

# 1. Evidence blockers

**None.**

No assigned fighter remains blocked because a required published rule/card fact is unknown or unconfirmed.

### Non-blocking provenance limitation

**Deadpool:** official Mondo material establishes first-party release provenance and published UmDb supplies the complete 30-card corpus. A first-party online full card-text/rulebook dump was not located. This is weaker primary card-level provenance than ordinary Restoration rulebook sets, but it does not block deck/card reconstruction.

# 2. Engine blockers

## Exact health semantics

Affected:
- King Arthur — `The Holy Grail`;
- Willow — `Resurrect`.

Need:

- exact health assignment on an undefeated fighter;
- atomic return/revive with source-defined exact health, avoiding an invalid intermediate health state.

## Positioned special components

Affected:
- Muldoon — traps;
- Invisible Man — Fog;
- Spike — Shadows.

Need reusable component instances with identity, owner/controller, lifecycle/location, space, placement/movement rules, occupancy/stacking policy and selectors.

Invisible Man additionally needs scoped dynamic adjacency: Fog spaces are mutually adjacent for Invisible Man's **movement only**.

## Movement-entry interruption

Affected: Muldoon traps.

Need movement/placement entry events plus protected interruption semantics capable of stopping a move immediately when a trap is entered. Opposing fighters can also trigger traps when placed into the trap space by an effect.

## Temporary undefeated off-board / dormancy

Affected: Invisible Man — `Vanish`.

Need undefeated `off_board` state, source-owned return timing and deterministic interaction between return hooks and dormant-player turn handling.

## Conditional combat-card / attack mode

Affected:
- Little Red — `What Big Ears You Have`;
- Deadpool — `Xavier Institute Faculty`.

Need declarative play-mode/attack-mode permissions without mutating immutable printed metadata.

## Blind BOOST result transformation

Affected: Spike — `Always Surprising`.

Need a transform hook over the resolved blind-BOOST amount before it is applied.

## Movement-path capture/query

Affected: Yennenga — `Stallion Charge`.

Need an ordered traversed-space path for a specific movement resolution and selectors over fighters/spaces crossed by that move.

## Partial pending-damage allocation

Affected: Yennenga — `Shield of the Archers` fighter ability.

Need allocation of one pending damage event among multiple recipients with per-recipient caps, residual damage on the original target, and preservation of original damage type/event identity.

## Non-action component in an ordered card zone

Affected: Little Red's Basket.

Need a runtime component that:

- is not counted as an action card;
- can occupy the discard pile as a real ordered-zone object;
- can be the current top discard object;
- exposes its printed wild symbol to Basket resolution;
- remains distinguishable from normal card instances for deck-count/exhaustion rules.

# 3. Other blockers

## Deadpool digital-adaptation policy

Published Deadpool effects depend on external/physical facts such as player names, food, spoken/noise actions, set ownership, mirror presence, sleeve/card-writing state, clothing, wagers and subjective board colour.

The project needs one central policy deciding whether each external predicate/action is:

- automatically supported;
- manually confirmed by players;
- replaced with a documented digital equivalent; or
- unsupported/disabled for online play.

Only after that policy is chosen should the shared engine expose generic external-predicate/player-confirmation primitives.

# Cross-cutting integration-contract gaps

These are **not evidence gaps** and are intentionally not patched by Worker A because they belong in shared schema/mechanics files owned by the Phase 4B orchestrator.

## Runtime identity / namespaces

Definition IDs such as `feint`, state keys such as `form`, and fighter IDs such as `alice` are only locally unambiguous. Integration needs explicit scoping/runtime identity so mirror matches and different decks using the same local IDs cannot collide.

At minimum distinguish:

- fighter definition vs fighter instance;
- card definition vs card instance;
- owner/player scope;
- fighter-local state/resource scope.

## Typed selectors, expressions and result binding

Worker manifests use operands such as:

- `discarded_card_boost`;
- `combat_damage_taken`;
- `count_other_voyage_cards_in_discard`;
- `chosen_draw_count`;
- `cards_drawn_by_this_effect`;
- `source_found_zone`.

The shared model needs typed operation outputs, bindings, expressions and selectors instead of implementation-defined free-form strings.

## Typed player relations

The runtime must distinguish at least:

- `combat_opponent` — the other player participating in the current combat;
- `chosen_opponent` — any opponent selected by a Scheme/effect.

These collapse to the same player in 1v1 but not in 3–4 player games.

## Choice contract / hidden information

Need a normative shape for chooser, choice domain, visibility, output binding, reconnect state and dependent consequences.

Choice bundles also need central rules for:

- precommitting selected options/parameters before resolution where required;
- owner-selected resolution order where the published/ruling semantics grant it;
- preserving private information until the correct reveal boundary.

## Batch / multi-target ordering

Operations such as `MOVE each_friendly_fighter` and damage to multiple fighters need shared rules for:

- sequential vs simultaneous resolution;
- controller-selected order where applicable;
- transition/trigger generation between elements;
- state re-evaluation after each resolved element.

This becomes critical once traps, transformation-on-damage and other cross-fighter triggers compose.

## Setup-selected active roster

Buffy chooses exactly one sidekick. Generic selectors must not treat the unselected Giles/Xander definition as active, defeated or targetable for that match.

## Deterministic randomness / replay

Random discard and similar effects need server-authoritative random results recorded in match history so reconnect/replay never re-rolls a resolved event.

## Effect defaults and control-flow vocabulary

Shared schema should explicitly define defaults for omitted effect fields and either promote or replace manifest constructs such as:

- `branches`;
- `normalization_blocker`;
- source-defined composites;
- dynamic destination/result identifiers;
- ordered optional clauses.

# Cross-fighter deterministic fixtures required before developer-ready

At minimum:

- Muldoon trap × effect movement/placement;
- Muldoon trap × Yennenga `Stallion Charge` path;
- Fog movement adjacency × traps;
- Yennenga partial allocation × multi-target/simultaneous damage;
- Yennenga `Surprise Volley` combat-participant replacement × Momentous Shift/history-sensitive effects;
- Invisible Man `Vanish` × dormant turn/return timing;
- Little Red Basket × discard manipulation;
- Little Red wild Basket × multi-branch card resolution;
- Buffy selectable sidekick × generic fighter selectors;
- Buffy `Slayer's Strength` × partial movement impossibility;
- Beowulf `Ancient Heirloom` × value-set/BOOST ordering;
- `Regroup` × near-empty deck/exhaustion;
- mirror matches to prove ID/state isolation;
- 3–4 player combat-opponent vs chosen-opponent isolation.

# Files and scope

Worker A owns exactly **35 files**:

- 17 fighter manifests under `docs/fighters/phase-4b/`;
- 17 deck manifests under `docs/cards/phase-4b/`;
- this report.

No shared schema, mechanics, rules, set registry, ambiguity register, research plan or README file is modified by this branch.

# Worker 4B-A Handoff

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Exact final Head: read after this final report write and supplied in the delivery response.  
Assigned fighters: **17**  
Verified: **9**  
Blocked: **8**  
Quantity validation: **PASS (17/17)**  
Current UmDb structure/type validation: **PASS (17/17)**  
Blocking evidence gaps: **0**  
Engine blocker families: **9**  
Other policy blockers: **1** — Deadpool digital adaptation  
Cross-cutting integration-contract gaps: **documented above; shared files intentionally untouched**  
Files in Worker A scope: **35**  
Merge to main: **not performed**
