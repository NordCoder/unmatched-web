# Timing, choices and effect semantics

**Phase:** 2 — Timing, choices, effects and global rulings  
**Status:** verified — Phase 2 gate passed 2026-07-26  
**Last verified:** 2026-07-26

This directory refines the Phase 1 core rules into a deterministic resolution contract suitable for an online rules engine. It describes **gameplay semantics**, not a programming-language architecture.

## Source baseline

- `CORE` — current [Unmatched Core Rules](https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Core-rules_EN_Light.pdf).
- `REF10` — [Unmatched Reference v10.0 (October 2025)](https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf). It explicitly states that its errata/rulings originate with the publisher and links the Rulings Archive; editorial synthesis remains lower authority than the underlying ruling.
- `ARCHIVE` — [Unmatched Rulings Archive](https://docs.google.com/document/d/13b-FbPq_vuqcc3IokeHvQ2ctJaDNZZuUaZmt4uft5h0/), the publisher/designer ruling compilation identified by REF10 as considered official by Restoration Games.
- `RULES-HUB` — [The Unmatched Club Rules Hub](https://www.the-unmatched.club/rules), a current secondary index. Its dispute entries distinguish rulings backed by the Rulings Archive; it is used for discovery/freshness, not to silently override higher-authority material.

The authority hierarchy remains the one in `docs/sources/source-policy.md`. Where a secondary summary uses looser wording than `CORE`/official errata, the higher-authority exact rule wins; for example, the canonical game-end predicate remains the official `start or end of any action` wording.

## Stable rule-ID prefixes

| Prefix | Scope |
| --- | --- |
| `TIMING-*` | event windows, ordering, checkpoints and queues |
| `CHOICE-*` | choice ownership, legality, locking and resume |
| `FX-*` | normalized effect structure and operation semantics |
| `CANCEL-*` | cancellation scope and effect identity |
| `ACTION-*` | gained/free actions and action accounting |
| `PLACE-*` | move/place/swap and destination semantics |
| `BONUS-*` | bonus-attack composite semantics |
| `ENDTURN-*` | `End the turn` control flow |
| `DORMANT-*` | dormant-player lifecycle |
| `INFO-*` | public/private/revealed information |
| `HOOK-*` | setup and lifecycle extension hooks |

IDs are semantic identifiers. Do not recycle them for unrelated meanings.

## Documents

1. [`event-model.md`](event-model.md) — turn/action/combat event model and effect ordering.
2. [`choices-and-resume.md`](choices-and-resume.md) — server-authoritative pending choices and resume points.
3. [`effect-model.md`](effect-model.md) — effect structure, conditions, dependencies and normalized primitive taxonomy.
4. [`cancellation.md`](cancellation.md) — card-effect cancellation and non-card effects.
5. [`action-accounting.md`](action-accounting.md) — gained/free actions and mandatory action use.
6. [`movement-and-placement.md`](movement-and-placement.md) — movement, placement, failure, swap and destination selection.
7. [`bonus-attacks.md`](bonus-attacks.md) — corrected bonus-attack semantics.
8. [`end-turn-and-dormancy.md`](end-turn-and-dormancy.md) — modern `End the turn` and dormant-player rulings.
9. [`information-visibility.md`](information-visibility.md) — public/private state and temporary reveal semantics.
10. [`setup-hooks.md`](setup-hooks.md) — generalized pre-game/setup extension order.
11. [`../rulings/global-rulings.md`](../rulings/global-rulings.md) — normalized global-ruling registry.
12. [`../rulings/ambiguity-register.md`](../rulings/ambiguity-register.md) — unresolved or deliberately deferred questions.
13. [`phase-2-validation.md`](phase-2-validation.md) — Phase 2 gate evidence.

## Core design consequence

A digital implementation must not treat resolution as a single synchronous `resolveEffect()` call. Effects can require choices owned by different players and therefore must be resumable from a stable semantic checkpoint.

Conceptually:

```text
event occurs
  -> collect eligible effects
  -> establish order
  -> resolve next effect
       -> evaluate condition
       -> request/lock required choices
       -> execute ordered operations
       -> emit resulting events
  -> continue queue
  -> reach action/turn checkpoint
```

The implementation may use any architecture that preserves these semantics.
