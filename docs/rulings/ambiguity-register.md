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
| `AMB-001` | P1 | `deferred-content` | When two different players have genuinely simultaneous **non-combat** triggered effects and no specific source establishes order, who orders/resolves them? | Do not invent a global active-player rule. During Phase 4, each affected interaction must cite an official ruling/set rule or remain blocked. Combat already has defender-first authority. |
| `AMB-002` | P1 | `deferred-content` | Which start/end-of-turn character abilities continue or become meaningful while their controller is dormant/off-board? | Global dormant restrictions are closed; fighter-specific lifecycle behavior must be verified from its rulings before that fighter is supported. |
| `AMB-003` | P1 | `deferred-content` | Exact return/revival health/location semantics differ among mechanics that return defeated/off-board fighters. | `RETURN_FIGHTER` is a distinct primitive. Each returning mechanic must provide its authoritative health/destination rule; do not apply `RECOVER` generically. |
| `AMB-004` | P1 | `deferred-content` | New 2026 Stars & Stripes action grants may have use-it/lose-it or action-type-specific semantics beyond the generic gained/free-action model. | Phase 3/4 must transcribe current official set/card rules before those fighters are developer-ready. |
| `AMB-005` | P2 | `deferred-content` | Some future cards may use dependency/cost wording whose grammar is not captured by the generic `if you do` examples. | Preserve source text/provenance and classify per card; add a primitive only when corpus evidence requires it. |
| `AMB-006` | P2 | `deferred-content` | Large-fighter swaps/placements and special occupancy can require orientation/space-selection rules beyond ordinary one-space relocation. | Battlefield/fighter-specific authority in Phase 4/5; ordinary relocation is already deterministic. |
| `AMB-007` | P2 | `deferred-content` | Does a specific special cleanup/replacement happen when `End the turn` redirects to Cleanup? | Global rule says Cleanup occurs. Each custom cleanup rule must state whether/how it participates; Shakespeare/Wayward Sisters are known examples. |
| `AMB-008` | P3 | `policy` | May players voluntarily disclose private hand information or communicate inferred information outside engine actions? | Out of deterministic rules scope. Server visibility remains defined by `information-visibility.md`; social/anti-cheat policy can be decided later. |

## Closed Phase 2 ambiguities

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

## Phase 2 gate impact

There are currently **no open P0 items**.

The P1 items above are `deferred-content`: they do not make the global event/choice/effect framework ambiguous, but they explicitly block any affected fighter/set from reaching `developer-ready` until authoritative content-specific resolution is recorded.

This distinction is intentional. Phase 2 defines a complete extension framework; it does not pretend that fighter-specific interactions can be resolved before the fighter/card corpus is researched.
