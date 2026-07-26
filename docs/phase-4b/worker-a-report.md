# Phase 4B Worker A research report

Scope: early/classic + retired licensed fighters assigned to Worker A.

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Base verification: branch was identical to Authorized Base before Worker A writes (`ahead=0`, `behind=0`).

## Source discipline

Canonical evidence order used by this worker:

1. official Restoration Games / IELLO / Mondo product, rulebook, set-rule, release and errata material;
2. current Unmatched Reference plus the publisher-backed Unmatched Rulings Archive;
3. published UmDb `/umdb/...` records for official deck/card metadata and effect transcription;
4. secondary indexes only for discovery/cross-checking, never as canonical evidence when a stronger source is available.

Current Reference:

- https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf

Official Rulings Archive:

- https://docs.google.com/document/d/13b-FbPq_vuqcc3IokeHvQ2ctJaDNZZuUaZmt4uft5h0/

Community/fan `/decks/...`, patch collections, reskins, `UnPatched` records and rebalances were excluded from canonical provenance.

## Source-hardening pass — 2026-07-27

All **17/17** assigned deck manifests were rechecked against the current published UmDb namespace.

Checks covered:

- published pool and game-deck quantity;
- unique-card count;
- card quantities;
- `usable_by` / fighter ownership;
- card type;
- printed combat value;
- BOOST;
- aggregate Attack / Versatile / Defense / Scheme distribution;
- source-sensitive card effects.

Weak secondary canonical links were removed where present, including Geektopia, Gridbeast, UnmatchedArena and generic The Unmatched Club hero-index pages. Hardened manifests use the strongest available combination of official product/rulebook/release material, current Reference, official Rulings Archive and published UmDb.

Result: **PASS — 17/17 decks reconcile by quantity and current UmDb type distribution.**

### Confirmed corrections

1. **Alice — `Manxome Foe`:** discard the top deck card and apply that discarded card's BOOST value; no separate reveal operation.
2. **Alice — `Claws That Catch`:** corrected from Versatile to **Attack**; totals now reconcile to 7 Attack / 17 Versatile / 2 Defense / 4 Scheme.
3. **Bigfoot — `Loner by Nature`:** end-of-turn draw corrected to optional (`may draw`).
4. **Buffy — setup order:** sidekick selection and 25+5 deck construction occur before shuffle/draw and before hero placement under the current setup-order ruling.
5. **Spike — Shadow lifecycle:** when all three Shadows are already on the board, the start-turn ability may move an existing Shadow.
6. **Yennenga — `VOLLEY`:** the bonus attack does not occur if either participant of the parent attack was defeated.
7. **Angel — `Wisdom of Ages`:** corrected from Versatile to **Attack**; totals now reconcile to 12 Attack / 13 Versatile / 2 Defense / 3 Scheme.
8. **Medusa — `Snipe`:** corrected from Attack to **Versatile**; totals now reconcile to 6 Attack / 17 Versatile / 3 Defense / 4 Scheme.
9. **Buffy — `Right-hand Man`:** final owner verification confirms the card is **Xander-only**. Adjacency to Buffy changes its combat value to 6; it does **not** grant Buffy permission to play the card. The earlier alternate-user hypothesis and its source blocker were removed completely.

No correction was imported from a fan deck or balance patch.

## Fighter results

| Fighter | Result | Primary note |
| --- | --- | --- |
| `alice` | verified | Big/Small state and all 30 cards reconcile with current UmDb. |
| `king-arthur` | blocked | `The Holy Grail` requires exact health assignment not present in the shared engine model. |
| `medusa` | verified | Multi-Harpy topology and deck normalize; `Snipe` type corrected. |
| `sinbad` | verified | Voyage tag + discard-zone counting fit existing metadata/state queries. |
| `robin-hood` | verified | Outlaw topology and post-attack movement fit the existing model. |
| `bigfoot` | verified | End-turn zone condition fits; optional draw corrected. |
| `robert-muldoon` | blocked | Positioned traps require generic token positioning plus movement-entry interruption. |
| `invisible-man` | blocked | Fog positioning and `Vanish` undefeated off-board/dormant lifetime need generic engine support. |
| `jekyll-and-hyde` | verified | Public form state and form-tagged cards fit the existing model. |
| `buffy` | verified | 35→30 construction, setup order and Xander-only `Right-hand Man` semantics are resolved. |
| `willow` | blocked | `Resurrect` requires exact health assignment after returning a defeated fighter. |
| `spike` | blocked | Positioned Shadows and conditional blind-BOOST multiplication need generic support. |
| `angel` | verified | Angel/Faith topology and deck reconcile after correcting `Wisdom of Ages`. |
| `little-red-riding-hood` | blocked | `What Big Ears You Have` conditionally changes legal combat-card mode. |
| `beowulf` | verified | Rage counter and per-damage-effect gain semantics fit existing resources/events. |
| `deadpool` | blocked | One card needs dynamic melee/ranged attack mode; multiple published effects also require an explicit digital-adaptation policy for external/physical predicates. |
| `yennenga` | blocked | `Stallion Charge` requires movement-path capture/query. |

Verified: **9** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `buffy`, `angel`, `beowulf`.  
Blocked: **8** — `king-arthur`, `robert-muldoon`, `invisible-man`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`.

## Quantity and structure validation

| Fighter | Published pool | Game deck | Reconciliation |
| --- | ---: | ---: | --- |
| Alice | 30 | 30 | PASS |
| King Arthur | 30 | 30 | PASS |
| Medusa | 30 | 30 | PASS |
| Sinbad | 30 | 30 | PASS |
| Robin Hood | 30 | 30 | PASS |
| Bigfoot | 30 | 30 | PASS |
| Robert Muldoon / InGen | 30 | 30 | PASS |
| Invisible Man | 30 | 30 | PASS |
| Jekyll & Hyde | 30 | 30 | PASS |
| Buffy | 35 | 30 | PASS: 25 base + selected Giles/Xander five-card group |
| Willow | 30 | 30 | PASS |
| Spike | 30 | 30 | PASS |
| Angel | 30 | 30 | PASS |
| Little Red Riding Hood | 30 | 30 | PASS; Basket is a non-action setup/reference component |
| Beowulf | 30 | 30 | PASS |
| Deadpool | 30 | 30 | PASS: 30 distinct one-copy cards |
| Yennenga | 30 | 30 | PASS |

**Quantity validation: PASS (17/17).**  
**Current UmDb structure/type validation: PASS (17/17).**

`usable_by` was checked against fighter topology. `Right-hand Man` is Xander-only; adjacency to Buffy changes only its combat value.

# Remaining blockers

## 1. Evidence blockers

**None.**

No assigned fighter remains blocked because a required published rule/card fact is unknown or unconfirmed.

### Non-blocking provenance limitation

- **Deadpool:** official Mondo material establishes first-party release provenance and current published UmDb supplies the complete 30-card corpus. A first-party online full card-text/rulebook dump was not located. This is weaker primary card-level provenance than normal Restoration rulebook sets, but it does not prevent deck/card reconciliation.

## 2. Engine blockers

### Exact health assignment

Affected:

- King Arthur — `The Holy Grail`;
- Willow — `Resurrect`.

Need: `SET_HEALTH(target, value)` or equivalent exact-health operation distinct from recovery/damage.

### Positioned special components/tokens

Affected:

- Robert Muldoon — traps;
- Invisible Man — Fog;
- Spike — Shadows.

Need: reusable position-sensitive component instances with generic placement, movement, return-to-supply, token-space selectors/predicates and ownership.

### Movement-entry interruption

Affected: Robert Muldoon traps.

Need: movement-step event plus protected `STOP_CURRENT_MOVEMENT` or equivalent interruption primitive.

### Temporary undefeated off-board / dormant fighter state

Affected: Invisible Man — `Vanish`.

Need: undefeated `off_board` fighter location/state with source-owned return timing and legal placement domain.

### Conditional combat-card mode / attack mode

Affected:

- Little Red Riding Hood — `What Big Ears You Have`;
- Deadpool — `Xavier Institute Faculty`.

Need: declarative play-permission/combat-mode clauses that can vary legal card mode or melee/ranged attack mode without mutating printed metadata.

### Blind BOOST result transformation

Affected: Spike — `Always Surprising`.

Need: transform hook over a resolved blind-BOOST amount before it is applied.

### Movement-path capture/query

Affected: Yennenga — `Stallion Charge`.

Need: movement resolution records the ordered traversed-space path and exposes selectors over fighters/spaces crossed by that specific move.

## 3. Other blockers

### Deadpool digital-adaptation policy

Deadpool has genuine published effects depending on external/physical conditions such as real player names, food, spoken/noise actions, set ownership, mirror presence, sleeve/card-writing state, clothing, wager and subjective board colour.

This is not an evidence problem and cannot be solved only by adding deterministic board-state primitives. The project needs one central policy for external-world effects, for example whether each predicate is automatically supported, manually confirmed by players, substituted with a digital equivalent, or disabled in online play.

After that policy exists, the engine can expose generic `external_predicate`, `external_action` and/or player-confirmation primitives without fighter-specific hacks.

## Files and scope

Fighter manifests: **17** under `docs/fighters/phase-4b/`.  
Deck manifests: **17** under `docs/cards/phase-4b/`.  
Report: `docs/phase-4b/worker-a-report.md`.

Worker A owns exactly these **35** files. No shared schema, mechanics, rules, set registry, ambiguity register, research plan or README file is modified by this branch.

## Worker 4B-A Handoff

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head immediately before this cleanup report update: `13ceacf1ac786d091e41bf159be180e660265cc6`  
Exact final Head: read after this final report write and supplied in the delivery response.  
Assigned fighters: **17**  
Verified: **9** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `buffy`, `angel`, `beowulf`  
Blocked: **8** — `king-arthur`, `robert-muldoon`, `invisible-man`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`  
Quantity validation: **PASS (17/17)**  
Current UmDb structure/type validation: **PASS (17/17)**  
Blocking evidence gaps: **0**  
Engine-blocked fighters: **8** (Deadpool also has a digital-policy blocker)  
Generic extension families: **8**; no shared schema file modified  
Files in Worker A scope: **35**  
Merge to main: **not performed**
