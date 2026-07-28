# Phase 4C Validation

## State

```text
status: PENDING_PARALLEL_WORKERS
base_sha: 106ae552ce597cde954c0a1b22374ef446974ce2
parent_issue: #12
expected_owner_aliases: 52
```

## Required worker verdicts

| Domain | Issue | Worker artifact | QA verdict |
|---|---:|---|---|
| Resolution, choices, history | #13 | `runtime-requirements/resolution.md` + fragment | pending |
| Fighter presence and objects | #14 | `runtime-requirements/fighter-presence.md` + fragment | pending |
| Movement and combat legality | #15 | `runtime-requirements/movement-combat.md` + fragment | pending |
| Damage, health and modifiers | #16 | `runtime-requirements/damage-modifiers.md` + fragment | pending |
| Cards, resources and actions | #17 | `runtime-requirements/cards-actions.md` + fragment | pending |

## Integration checks

The integration worker must verify:

- all source reports parse and expose the same 52 aliases confirmed by Phase 4 validation;
- every alias has one complete disposition;
- split aliases enumerate every resulting canonical requirement;
- no alias is owned by competing worker fragments;
- every canonical requirement names affected fighters and source aliases;
- runtime state and transition contracts are deterministic;
- command/event, choice, visibility and persistence implications agree across domains;
- test scenarios cover launch fighters and every capability family required by them;
- extension-handler proposals satisfy the restrictive exception gate;
- manifests can migrate to canonical IDs or a lossless alias map without changing published behavior.

## Blocking conditions

Phase 4C fails while any of the following remains:

```text
unresolved owner aliases > 0
multiply-mapped owner aliases > 0
canonical requirement without deterministic runtime contract > 0
launch-scope requirement status = blocked
undocumented fighter/card-ID core behavior > 0
independent QA scope violations > 0
```

A final verdict will replace this pending record only after five worker QA verdicts and full integrated validation.
