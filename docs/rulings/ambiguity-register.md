# Ambiguity register

**Last reviewed:** 2026-07-26  
**Purpose:** prevent unresolved gameplay questions from being silently answered by implementation code.

Severity follows `docs/specification-readiness.md`:

- `P0` — can change winner/game end or irrecoverably corrupt state;
- `P1` — changes legal action, hidden information, targeting, damage, card order/ownership, fighter defeat or combat outcome;
- `P2` — narrower edge case with deterministic containment/deferment;
- `P3` — documentation/product-policy quality only.

Status meanings:

- `open` — authoritative answer still required;
- `deferred-content` — framework is sufficient, but a fighter/set cannot be marked developer-ready until its own authority closes the question;
- `policy` — not a gameplay-rule ambiguity; project policy is explicitly documented;
- `closed` — resolved in the current corpus.

## Open/deferred items

| ID | Severity | Status | Question | Containment / next authority |
| --- | --- | --- | --- | --- |
| `AMB-001` | P1 | `deferred-content` | When two different players have genuinely simultaneous **non-combat** triggered effects and no specific source establishes order, who orders/resolves them? | Do not invent a global active-player rule. During Phase 4B, each affected interaction must cite an official ruling/set rule or remain blocked. Combat already has defender-first authority. Phase 4A's simultaneous Lodge choice is not this case: its source explicitly requires all choices to lock before reveal. |
| `AMB-002` | P1 | `deferred-content` | Which start/end-of-turn character abilities continue or become meaningful while their controller is dormant/off-board? | Global dormant restrictions are closed; fighter-specific lifecycle behavior must be verified from its rulings before that fighter is supported. No Phase 4A representative requires an unresolved dormant exception. |
| `AMB-003` | P1 | `deferred-content` | Exact return/revival health/location semantics differ among mechanics that return defeated/off-board fighters. | `RETURN_FIGHTER` is a distinct primitive. Dracula's `Baptism of Blood` is resolved for the representative corpus by its own source-defined Sister return; other return mechanics remain per-content work. |
| `AMB-004` | P1 | `deferred-content` | New 2026 Stars & Stripes action grants may have use-it/lose-it or action-type-specific semantics beyond the generic gained/free-action model. | Phase 4B must transcribe current official set/card rules before those fighters are developer-ready. |
| `AMB-005` | P2 | `deferred-content` | Some future cards may use dependency/cost wording whose grammar is not captured by the representative corpus. | Phase 4A validated explicit costs/dependencies on multiple decks; preserve source wording/provenance and extend generically only when later corpus evidence requires it. |
| `AMB-006` | P2 | `deferred-content` | Large-fighter swaps/placements and special occupancy can require orientation/space-selection rules beyond ordinary one-space relocation. | Battlefield/fighter-specific authority in Phase 4B/5; ordinary relocation is already deterministic. |
| `AMB-007` | P2 | `deferred-content` | Does a specific special cleanup/replacement happen when `End the turn` redirects to Cleanup? | Global rule says Cleanup occurs. Phase 4A closes the Wayward representative case: Bubbling Brew can replace the cleanup discard, while post-Cleanup after-attacking spell opportunity is skipped. Other custom cleanup sources remain per-content. |
| `AMB-008` | P3 | `policy` | May players voluntarily disclose private hand information or communicate inferred information outside engine actions? | Out of deterministic rules scope. Server visibility remains defined by `information-visibility.md`; social/anti-cheat policy can be decided later. |

## Closed ambiguities

| ID | Former issue | Resolution |
| --- | --- | --- |
| `CLOSED-001` | Placement destination must be empty at selection time? | No, not generically. When wording omits `empty/other`, an occupied space can be selected and placement can then fail. See `PLACE-020–022`. |
| `CLOSED-002` | Does bonus attack continue if a participant dies in first combat? | No if either attacker or defender was defeated. See `BONUS-013`. |
| `CLOSED-003` | Is bonus attack a new action? | No. It is nested combat inside the original Attack action. See `BONUS-003–004`. |
| `CLOSED-004` | Does `End the turn` skip Cleanup? | No. It redirects to Cleanup, then End-of-Turn. See `ENDTURN-001–008`. |
| `CLOSED-005` | Is a free action still an action? | Yes for action-counting/checkpoint semantics, while not consuming ordinary action budget. See `ACTION-020–023`. |
| `CLOSED-006` | Does attempting to cancel an effectless card mean its effects were canceled? | No under the official specific cancellation ruling. See `CANCEL-020–021`. |
| `CLOSED-007` | Can the engine continue while another player choice is pending? | No. Persist the choice and semantic resume point. See `CHOICE-001–003`. |
| `CLOSED-008` | Are `move` and `place` interchangeable? | No. They are distinct operations with different path and selection/success semantics. |
| `CLOSED-009` | Do non-draw deck operations trigger exhaustion on an empty deck? | Not generically. Exhaustion is tied to a failed required draw; blind BOOST is explicitly not a draw. |
| `CLOSED-010` | Can card ownership be inferred from its current player-controlled zone? | No. Phase 4A Black Panther proves immutable owner and current zone/use authority must be separate. |
| `CLOSED-011` | Is printed value always the immutable number stored on the card definition? | No. Phase 4A Sherlock proves a source may temporarily change the effective printed-value layer while immutable base metadata remains preserved. |
| `CLOSED-012` | Must all BOOST choices in an `up to N` sequence be locked before resolving any? | No. Phase 4A Black Panther proves source timing can make repetitions sequential because an ability resolves after each BOOST and changes state before the next choice. |

## Gate impact

There are currently **no open P0 items**.

There are also **no unresolved P1 items required to execute the Phase 4A representative fighter/deck corpus** under its documented source rules.

The remaining P1 items are `deferred-content`: they explicitly block only affected Phase 4B/5 content until its authoritative source is transcribed. The framework does not invent missing behavior.