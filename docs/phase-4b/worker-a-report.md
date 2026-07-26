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

Current reference:

- https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf

Official rulings archive:

- https://docs.google.com/document/d/13b-FbPq_vuqcc3IokeHvQ2ctJaDNZZuUaZmt4uft5h0/

Published UmDb states that `/umdb/...` covers officially published Unmatched decks. Community/fan `/decks/...`, patch collections, reskins and rebalances were explicitly excluded from canonical provenance.

## Source-hardening pass — 2026-07-27

A second source pass was performed after the initial Worker A handoff.

### Provenance hardening

Weak secondary canonical links were removed/replaced where present, including:

- Geektopia fighter/deck pages;
- Gridbeast fighter pages;
- UnmatchedArena fighter pages;
- generic The Unmatched Club hero-index pages used as evidence rather than discovery.

They were replaced by the strongest available combination of:

- official Restoration Games product/rulebook material;
- official Mondo release material for Deadpool;
- current Unmatched Reference;
- official Rulings Archive;
- current published UmDb deck/card records.

The Unmatched Club remains useful as a discovery interface to individual sourced rulings, but is not required as canonical fighter/card provenance in the hardened manifests.

### Current UmDb reconciliation

All **17/17** assigned deck manifests were rechecked against the current published UmDb namespace.

Checks included:

- published pool and game-deck quantity;
- unique-card count;
- card quantities;
- character/`usable_by` ownership where UmDb exposes it;
- card type;
- printed combat value;
- BOOST;
- aggregate Attack / Versatile / Defense / Scheme distribution;
- source-sensitive card effects;
- rejection of similarly named fan patches and `UnPatched` records.

Result: **PASS — 17/17 decks reconcile by quantity and current UmDb type distribution.**

### Confirmed corrections made by the hardening pass

1. **Alice — `Manxome Foe`:** normalized as discarding the top deck card and applying that discarded card's BOOST value, rather than as a separate reveal operation.
2. **Alice — `Claws That Catch`:** corrected from Versatile to **Attack**; Alice now reconciles to current UmDb totals (7 Attack / 17 Versatile / 2 Defense / 4 Scheme).
3. **Bigfoot — `Loner by Nature`:** end-of-turn draw corrected to **optional** (`may draw`).
4. **Buffy — setup order:** sidekick selection and 25+5 deck construction moved to the current pre-shuffle/draw setup stage rather than after hero placement.
5. **Spike — Shadow lifecycle:** current Reference semantics record that, once all three Shadows are on the board, the start-turn ability may move an existing Shadow instead of requiring an unused token.
6. **Yennenga — `VOLLEY`:** current ruling recorded that the bonus attack does not occur if **either the attacker or defender** of the parent attack was defeated.
7. **Angel — `Wisdom of Ages`:** corrected from Versatile to **Attack**; Angel now reconciles to 12 Attack / 13 Versatile / 2 Defense / 3 Scheme.
8. **Medusa — `Snipe`:** corrected from Attack to **Versatile**; Medusa now reconciles to 6 Attack / 17 Versatile / 3 Defense / 4 Scheme.

No correction was imported from a fan deck or balance patch.

## Fighter results after hardening

| Fighter | Result | Primary note |
| --- | --- | --- |
| `alice` | verified | Public Big/Small state and all 30 cards reconcile with current UmDb. |
| `king-arthur` | blocked | `The Holy Grail` assigns an exact health value; current shared model has recovery but no absolute health assignment. Evidence is sufficient. |
| `medusa` | verified | Multi-Harpy topology and deck normalize; `Snipe` type corrected during hardening. |
| `sinbad` | verified | Voyage tag + discard-zone counting fit existing metadata/state queries. |
| `robin-hood` | verified | Outlaw topology and post-attack movement fit the existing model. |
| `bigfoot` | verified | End-turn zone condition fits existing model; optional draw corrected during hardening. |
| `robert-muldoon` | blocked | Positioned traps require token placement/return plus movement-entry interruption. Evidence is sufficient. |
| `invisible-man` | blocked | Positioned Fog semantics and `Vanish` undefeated off-board/dormant lifetime need generic support. Evidence is sufficient. |
| `jekyll-and-hyde` | verified | Public form state and form-tagged cards fit the existing model. |
| `buffy` | blocked | 35→30 construction is verified. The claimed alternate Buffy permission on `Right-hand Man` could not be substantiated by current official rulebook/Reference/UmDb evidence and is therefore not encoded. |
| `willow` | blocked | Dark Willow state fits; `Resurrect` requires exact health assignment after returning a defeated fighter. Evidence is sufficient. |
| `spike` | blocked | Positioned Shadow tokens and conditional blind-BOOST multiplication need generic support. Evidence is sufficient. |
| `angel` | verified | Angel/Faith topology and losing-attack ability fit; card-type distribution reconciled after correcting `Wisdom of Ages`. |
| `little-red-riding-hood` | blocked | Basket/discard-symbol state is representable, but `What Big Ears You Have` conditionally changes legal card mode/play permission. Evidence is sufficient. |
| `beowulf` | verified | Rage counter and per-damage-effect gain semantics fit existing resources/events. |
| `deadpool` | blocked | Published deck intentionally uses external-world/physical predicates and a per-card melee/ranged override outside the deterministic shared model. Official Mondo release provenance + complete published UmDb corpus are recorded. |
| `yennenga` | blocked | Damage redirection and VOLLEY are sourced; `Stallion Charge` still needs movement-path capture. |

Verified: **8** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `angel`, `beowulf`.  
Blocked: **9** — `king-arthur`, `robert-muldoon`, `invisible-man`, `buffy`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`.

Of the nine blocked fighters:

- **8 are blocked by shared model / digital-policy gaps**, not missing deck evidence;
- **1 (`buffy`) has a remaining component-level source gap** for the claimed alternate-user clause on `Right-hand Man`.

## Quantity and structure validation

| Fighter | Published pool | Game deck | Current UmDb reconciliation |
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
| Little Red Riding Hood | 30 | 30 | PASS; Basket is non-action setup/reference component |
| Beowulf | 30 | 30 | PASS |
| Deadpool | 30 | 30 | PASS: 30 distinct one-copy cards |
| Yennenga | 30 | 30 | PASS |

**Quantity validation: PASS (17/17).**  
**Current UmDb type-distribution validation: PASS (17/17).**

`usable_by` was checked against fighter topology. Conditional or uncertain permissions are not silently widened. In particular, `Right-hand Man` remains canonically Xander-owned until primary component evidence establishes any alternate Buffy play permission.

## Generic schema/effect extension proposals

### 1. Exact health assignment

Affected: King Arthur — `The Holy Grail`; Willow — `Resurrect`.  
Evidence: published UmDb plus official/current reference material.  
Need: `SET_HEALTH(target, value)` or equivalent exact-health operation distinct from recovery/damage.  
Blocker: **yes**.

### 2. Positioned special components/tokens

Affected: Robert Muldoon traps, Invisible Man Fog, Spike Shadows.  
Evidence: official set rules/current Reference + published UmDb.  
Need: reusable position-sensitive component instances with `PLACE_TOKEN`, `MOVE_TOKEN`, `RETURN_TOKEN_TO_SUPPLY`, token-space selectors/predicates and ownership.  
Blocker: **yes**.

### 3. Movement-entry interruption

Affected: Robert Muldoon traps.  
Evidence: official InGen vs. Raptors rules/current Reference.  
Need: movement-step event plus protected `STOP_CURRENT_MOVEMENT`.  
Blocker: **yes**.

### 4. Temporary undefeated off-board / dormant fighter state

Affected: Invisible Man — `Vanish`.  
Evidence: published deck plus current Reference Dormant Player rules.  
Need: undefeated `off_board` location/state with source-owned return timing and legal placement domain.  
Blocker: **yes**.

### 5. Conditional card user/type/attack-mode permission

Affected: Little Red Riding Hood — `What Big Ears You Have`; Deadpool — `Xavier Institute Faculty`. Buffy would also use this model if primary evidence later establishes the claimed `Right-hand Man` alternate-user clause.  
Need: declarative play-permission clauses evaluated at commit time, able to vary legal combat mode/type or eligible fighter without mutating printed metadata.  
Blocker: **yes** for Little Red and Deadpool; Buffy currently remains source-blocked first.

### 6. Blind BOOST result transformation

Affected: Spike — `Always Surprising`.  
Need: transform hook over a resolved blind-BOOST amount before application.  
Blocker: **yes**.

### 7. Movement-path capture/query

Affected: Yennenga — `Stallion Charge`.  
Need: movement resolution records the ordered traversed-space path and exposes selectors over fighters/spaces crossed by that move.  
Blocker: **yes**.

### 8. External-world / physical-component predicates

Affected: Deadpool — multiple cards and fighter ability.  
Evidence: official Mondo release provenance plus complete current published UmDb deck.  
Need: central online-platform policy / compatibility layer for real-name, food, spoken/noise, ownership, mirror, sleeve/card-writing, clothing, wager and subjective-colour predicates.  
Blocker: **yes**.

## Remaining source gaps / limitations

### Blocking source gap

- **Buffy — `Right-hand Man`:** current official Buffy rulebook, current Reference and current published UmDb establish a Xander card, value 2, BOOST 3, and the adjacency-based value-6 effect. They do **not** substantiate the community-transcribed parenthetical that Buffy may also play the card while adjacent to Xander. The parenthetical has been removed from canonical semantics pending a primary component scan/photo or publisher ruling.

### Non-blocking provenance limitation

- **Deadpool:** official Mondo material establishes first-party release provenance and current published UmDb supplies the complete 30-card corpus. A first-party online full card-text/rulebook dump was not located. This is weaker primary card-level provenance than the normal Restoration rulebook sets, but it does not prevent deck/card reconciliation; Deadpool is blocked by digital semantics, not missing card identity/quantity data.

No other unresolved source gap blocks assigned fighter/deck reconstruction.

## Files and scope

Fighter manifests: **17** under `docs/fighters/phase-4b/`.  
Deck manifests: **17** under `docs/cards/phase-4b/`.  
Report: `docs/phase-4b/worker-a-report.md`.

Pre-report hardening Base→Head audit still contains exactly the 35 Worker A-owned files. No shared schema, mechanics, rules, set registry, ambiguity register, research plan or README file was changed.

## Worker 4B-A Handoff

Branch: `phase-4b-worker-a-classics`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head immediately before this hardened report update: `804a976b92c719095836a7bbed8d41eef235464b`  
Exact final Head: read after this final report write and supplied in the delivery response.  
Assigned fighters: **17**  
Verified: **8** — `alice`, `medusa`, `sinbad`, `robin-hood`, `bigfoot`, `jekyll-and-hyde`, `angel`, `beowulf`  
Blocked: **9** — `king-arthur`, `robert-muldoon`, `invisible-man`, `buffy`, `willow`, `spike`, `little-red-riding-hood`, `deadpool`, `yennenga`  
Quantity validation: **PASS (17/17)**  
Current UmDb structure/type validation: **PASS (17/17)**  
Hardening corrections: **8 confirmed content corrections**  
Schema-extension proposals: **8 generic extension families**; no shared schema file modified  
Blocking source gaps: **1** — Buffy `Right-hand Man` alternate-user clause  
Model/policy-blocked fighters: **8**  
Files in Worker A scope: **35**  
Merge to main: **not performed**
