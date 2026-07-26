# Phase 2 validation

**Phase:** 2 — Timing, choices, effects and global rulings  
**Validation date:** 2026-07-26  
**Result:** PASS for the global resolution framework  
**Overall project readiness:** NOT YET `developer-ready`; content/set/battlefield phases remain.

## Gate under test

The Phase 2 gate requires:

1. every core timing window to be ordered;
2. a deterministic rule for pausing/resuming when player input is required;
3. global errata/rulings to be reconciled with the current Core Rules;
4. unknown interactions to be explicit rather than guessed.

## Coverage matrix

| Area | Result | Specification |
| --- | --- | --- |
| turn/action checkpoints | PASS | `event-model.md` |
| combat timing and nested effect ordering | PASS | `event-model.md` |
| pending choices and reconnect | PASS | `choices-and-resume.md` |
| optional/useless choices | PASS | `choices-and-resume.md` |
| effect categories and ordered operations | PASS | `effect-model.md` |
| partial resolution/dependency representation | PASS | `effect-model.md` |
| cancellation identity/scope | PASS | `cancellation.md` |
| gained/free actions | PASS | `action-accounting.md` |
| move/place distinction and failed placement | PASS | `movement-and-placement.md` |
| corrected bonus-attack procedure | PASS | `bonus-attacks.md` |
| `End the turn` control flow | PASS | `end-turn-and-dormancy.md` |
| dormant-player lifecycle | PASS | `end-turn-and-dormancy.md` |
| hidden/public information classes | PASS | `information-visibility.md` |
| generalized setup extension order | PASS | `setup-hooks.md` |
| global rulings indexed | PASS | `../rulings/global-rulings.md` |
| unresolved questions explicit | PASS | `../rulings/ambiguity-register.md` |

## Fixture A — persisted hidden combat choice

State:

```yaml
phase: ATTACK_CARD_SELECTION
active_player: P1
attacker: H1
defender: H2
```

Trace:

1. P1 chooses a legal attack card (`CHOICE-030`).
2. Server removes/marks the card from ordinary hand availability and stores a `COMMITTED_HIDDEN` combat commitment (`INFO-030`).
3. Public state tells P2 that an attack is pending but does not disclose the card identity.
4. Resolution pauses at the defense choice (`CHOICE-001–003`).
5. P2 disconnects and reconnects.
6. Server reconstructs the same defense choice and legal defense options from authoritative state; P1's committed card remains hidden (`INFO-033`).
7. P2 commits a defense card or explicitly declines.
8. Both commitments reveal simultaneously (`INFO-032`).
9. Combat resumes at `IMMEDIATELY`, not from the start of Attack.

**Result:** deterministic pause/reconnect without information leak or duplicate card choice.

## Fixture B — cancellation cannot travel backward in time

Scenario: defender and attacker both have `IMMEDIATELY` effects; the attacker's later effect cancels effects on the opposing combat card.

1. Enter `IMMEDIATELY` window (`TIMING-030`).
2. Defender effect resolves first (`TIMING-033`).
3. Apply its state changes fully.
4. Attacker cancellation resolves next.
5. Cancel remaining eligible effects on the defender's card, but do not undo the already resolved defender effect (`CANCEL-030–031`).

**Result:** resolution follows official defender-first combat ordering and causal effect resolution.

## Fixture C — selectable occupied placement that fails

Effect wording: place fighter `F` on a space in a specified zone; wording does **not** say `empty` or `other`.

1. Choice generator finds all spaces satisfying the textual zone constraint (`PLACE-020`).
2. An occupied qualifying space is included as a selectable option.
3. Player selects that occupied space.
4. Placement operation checks ordinary occupancy and fails.
5. Fighter remains in its original space (`PLACE-020`, `PLACE-031`).
6. Resolution continues to independently resolvable later operations (`FX-031`).

**Result:** choice legality and operation success are not incorrectly collapsed.

## Fixture D — bonus attack after movement

Main Attack produces a bonus attack; an After Combat effect moves the defender out of ordinary range before the bonus resolves. Neither combat participant is defeated.

1. Main combat resolves through its After Combat queue.
2. Bonus trigger resolves (`BONUS-001`).
3. Same attacker and defender are retained (`BONUS-010`).
4. Do not run a new ordinary target/range declaration check (`BONUS-012`).
5. Defender gets a new defense choice (`BONUS-011`).
6. Resolve nested bonus combat (`BONUS-030`).
7. Do not run top-level winner checkpoint between the two combats (`BONUS-003–004`).
8. Return to outer Attack continuation and eventually reach one end-of-action checkpoint.

Variant: if either participant was defeated during the main combat, step 2 produces no bonus combat (`BONUS-013`).

**Result:** corrected official bonus-attack semantics are representable without treating the bonus as an extra action.

## Fixture E — `End the turn` during an Attack

Assume an effect resolves during the action and instructs `End the turn` while later effects/additional effects remain scheduled.

1. Set end-turn control flag and stop unresolved effect processing (`ENDTURN-001–002`).
2. Jump to Cleanup.
3. Perform ordinary and fighter-specific Cleanup behavior that applies (`ENDTURN-003`).
4. Skip post-Cleanup Additional Effects that would otherwise happen after attacking (`ENDTURN-004`).
5. Reach action end; run winner and dormant checkpoints (`ENDTURN-006`).
6. Process applicable End-of-Turn effects and hand limit (`ENDTURN-007–008`).
7. Discard unused normal/gained/free action opportunities with the ending turn (`ENDTURN-005`).

**Result:** `End the turn` is not misimplemented as either `actions_remaining = 0` or an immediate function return that skips Cleanup.

## Fixture F — gained/free action boundaries

1. P1 starts with two ordinary action opportunities.
2. An effect grants one unrestricted action (`ACTION-010`).
3. After resolving the current action, P1 must still consume all remaining required action opportunities unless another rule ends the turn (`ACTION-002`, `ACTION-011`).
4. A separate source creates a free action.
5. Free action consumes no ordinary budget but creates a genuine action start/end boundary (`ACTION-020–022`).
6. Winner/dormant checkpoints run for that free action as normal.

**Result:** action identity is independent of action-budget cost.

## Fixture G — dormant transition

During P1's action, all P2 fighters leave the battlefield but one undefeated fighter might later return.

Case 1 — fighter returns before action end:

1. Intermediate board presence becomes zero.
2. Do not immediately mark P2 dormant (`DORMANT-002`).
3. Fighter returns before the action ends.
4. At end-of-action check P2 has a fighter on board; P2 is not dormant.

Case 2 — no fighter returns:

1. At action end P2 has zero fighters on board.
2. Mark/derive P2 as dormant (`DORMANT-001`).
3. P2 cannot take normal actions/draw/discard on the dormant turn and is excluded from opponent-selection semantics as specified (`DORMANT-003–006`).
4. If a later authoritative effect returns a fighter, re-evaluate board presence at the defined checkpoint and dormancy can clear (`DORMANT-007`).

**Result:** dormancy is a checkpoint-derived lifecycle state, not immediate elimination.

## Fixture H — setup choice survives reconnect

Character uses a pre-game deck/configuration choice.

1. Setup reaches `CHARACTER_CONFIGURATION`/`DECK_CONSTRUCTION` hook (`HOOK-001`, `HOOK-020–021`).
2. Create a normal pending private/public choice according to the character source (`HOOK-011`).
3. Persist selection and source provenance.
4. Only after legal construction completes does the engine shuffle and draw the starting hand.
5. Reconnect before completion restores the same setup stage and legal choice domain.

**Result:** pre-game configuration uses the same deterministic choice/resume infrastructure as gameplay.

## Authority reconciliation

Phase 2 explicitly incorporates the material global changes/rulings discovered during research:

- modern action-boundary winner rule (carried from Phase 1);
- corrected bonus-attack defeat behavior;
- dormant players;
- `End the turn` → Cleanup → End-of-Turn;
- corrected placement selection-versus-success behavior;
- generalized setup order;
- effectless-card cancellation distinction;
- free actions remain actions;
- multiple decisions locked before execution;
- reveal/blind-BOOST semantics;
- battlefield-item effects as effects on the combat card.

See `docs/rulings/global-rulings.md` for provenance classification.

## Unknown-interaction validation

The ambiguity register contains no open `P0` item. Remaining `P1` questions are explicitly `deferred-content`: they concern fighter/set-specific interactions for content that has not yet passed Phase 4 transcription. The framework provides a place to express the result, but does not invent the missing ruling.

Examples include cross-player simultaneous non-combat triggers and special dormant/off-board lifecycle abilities.

## Gate decision

**PASS — Phase 2 is complete for the global resolution framework.**

The documentation now specifies:

- ordered core lifecycle/timing checkpoints;
- resumable player choices;
- hidden/public information boundaries;
- normalized effect categories and primitives;
- cancellation identity;
- normal/gained/free action accounting;
- move/place/failure semantics;
- bonus attacks;
- `End the turn`;
- dormancy;
- setup extension hooks;
- an explicit process for content-specific unresolved rulings.

This does **not** make actual fighter decks developer-ready. Phase 3 must establish the exhaustive current set/release registry, and Phase 4 must transcribe and validate fighter/card behavior against this framework.
