# Global rulings registry

**Last reviewed:** 2026-07-26

This registry records global rulings that materially refine the Core Rules and links each ruling to the normalized specification that implements it. It is an index, not a replacement for the underlying rule documents.

Authority labels:

- `official_erratum` — explicitly identified as official errata in REF10/current publisher material;
- `official_ruling` — publisher/designer ruling represented in the Rulings Archive/REF10;
- `official_specific` — official ruling exists for a narrower interaction from which only that supported scope is normalized;
- `normalized_core` — deterministic consequence of authoritative rules, not new publisher wording.

## Registry

| ID | Topic | Normalized rule | Authority | Primary source |
| --- | --- | --- | --- | --- |
| `GR-001` | modern winner check | check at start/end of actions; if both heroes defeated, active-turn player wins | `official_erratum` | CORE p.14; REF10 Official Errata |
| `GR-002` | bonus attack after defeat | no bonus attack if attacker or defender was defeated during first combat | `official_erratum` / `official_ruling` | REF10 Bonus Attacks |
| `GR-003` | cancellation of effectless card | playing generic cancellation against a card with no cancelable printed effects does not mean that card's effects were canceled | `official_specific` | REF10 Cancelling Card Effects / King Arthur ruling |
| `GR-004` | dormant player | board-presence predicate checked at end of action; dormant restrictions apply | `official_ruling` | REF10 Major Rulings: Dormant Players |
| `GR-005` | `End the turn` | stop unresolved effects, proceed to Cleanup, then End-of-Turn | `official_ruling` | REF10 Major Rulings; current Rules Hub |
| `GR-006` | placement correction | occupied space may be selectable when text does not require `empty/other`, then placement fails under occupancy | `official_ruling` | REF10 Placement; current Rules Hub correction |
| `GR-007` | generalized setup order | character configuration precedes ordinary shuffle/draw; placement/post-placement order standardized | `official_ruling` | REF10 Setup Order |
| `GR-008` | multiple decisions | when an effect gives multiple decisions, lock the required decisions before executing outcomes unless wording is explicitly sequential | `official_ruling` | current Rules Hub, linked ruling corpus |
| `GR-009` | no valid move/place destination | skip the impossible operation and continue independently resolvable parts | `official_ruling` + core partial-resolution rule | CORE p.13; Rules Hub |
| `GR-010` | non-card After Combat ability ordering | card `AFTER COMBAT` effects resolve before abilities that trigger after combat/attacking when the ruling distinguishes those layers | `official_ruling` | Rules Hub ruling index |
| `GR-011` | gained actions | a gained action is not voluntarily skippable | `official_ruling` / core-consistent | CORE mandatory actions; Rules Hub |
| `GR-012` | reveal semantics | revealed cards are shown to all players and return/remain in source location unless effect changes disposition | `official_ruling` / reference | REF10 Revealing Cards |
| `GR-013` | blind BOOST | blind BOOST is not a draw; empty deck contributes 0 rather than exhaustion damage | `official_ruling` / reference | REF10 Blind BOOST |
| `GR-014` | returned defeated fighter | generic return-to-battlefield effects restore a defeated fighter according to the applicable return ruling, not ordinary recovery | `official_specific` | current Rules Hub; must be checked per returning mechanic |
| `GR-015` | `cannot leave space` | prevents move/place relocation but defeat removal still occurs | `official_ruling` | current Rules Hub |
| `GR-016` | free actions count as actions | a free action does not spend ordinary budget but is still an action for action-counting/checkpoint rules | `official_specific` generalized only to action identity | REF10 Bloody Mary clarification |
| `GR-017` | static modifiers | continuously applicable bonuses/modifiers are re-evaluated rather than consumed as one-shot queued effects | `official_ruling` / interaction corpus | current Rules Hub |
| `GR-018` | combat item attachment | chosen combat item effect becomes an effect on the played combat card for that combat and participates in card-effect ordering/cancellation | `official_ruling` / set reference | REF10 Battlefield Items |

## Normalized locations

- `GR-001` → `docs/rules/defeat-and-game-end.md`, `docs/mechanics/event-model.md`
- `GR-002` → `docs/mechanics/bonus-attacks.md`
- `GR-003`, `GR-018` → `docs/mechanics/cancellation.md`
- `GR-004` → `docs/mechanics/end-turn-and-dormancy.md`
- `GR-005` → `docs/mechanics/end-turn-and-dormancy.md`
- `GR-006`, `GR-009`, `GR-015` → `docs/mechanics/movement-and-placement.md`
- `GR-007` → `docs/mechanics/setup-hooks.md`
- `GR-008` → `docs/mechanics/choices-and-resume.md`, `event-model.md`
- `GR-010`, `GR-017` → `docs/mechanics/event-model.md`
- `GR-011`, `GR-016` → `docs/mechanics/action-accounting.md`
- `GR-012` → `docs/mechanics/effect-model.md`, `information-visibility.md`
- `GR-013` → `docs/mechanics/effect-model.md`

## Scope discipline

Some global statements in community reference material are editorial abstractions from one or more official rulings. This project adopts only the behavior supported by the cited ruling scope. During fighter/card transcription, a more specific official ruling may refine one of these records without invalidating unrelated global semantics.

When that happens:

1. preserve this global ruling ID;
2. add the narrower fighter/card ruling separately;
3. document precedence;
4. update validation fixtures if behavior changes.
