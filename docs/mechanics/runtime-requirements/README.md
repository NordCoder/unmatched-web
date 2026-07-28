# Shared Runtime Requirements

## Status

```text
phase: 4C
status: parallel-research
base_sha: 106ae552ce597cde954c0a1b22374ef446974ce2
parent_issue: #12
owner_aliases_expected: 52
```

This directory consolidates corpus-proven `A/B/C/D-REQ-*` aliases into canonical generic runtime capabilities.

## Invariant

Published fighter/card behavior remains in manifests. Shared gameplay behavior is implemented through generic capabilities. A canonical requirement ID must never become a disguised fighter/card special case.

Core runtime code may dispatch by operation or registered capability type. It may not branch on fighter or card identity.

## Worker ownership

| Worker | Issue | Branch | Owned fragment | Owned narrative |
|---|---:|---|---|---|
| Resolution, choices, history | #13 | `phase-4c-worker-1-resolution` | `fragments/resolution.yaml` | `resolution.md` |
| Presence and objects | #14 | `phase-4c-worker-2-presence` | `fragments/fighter-presence.yaml` | `fighter-presence.md` |
| Movement and combat legality | #15 | `phase-4c-worker-3-movement-combat` | `fragments/movement-combat.yaml` | `movement-combat.md` |
| Damage, health and modifiers | #16 | `phase-4c-worker-4-damage-modifiers` | `fragments/damage-modifiers.yaml` | `damage-modifiers.md` |
| Cards, resources and actions | #17 | `phase-4c-worker-5-cards-actions` | `fragments/cards-actions.yaml` | `cards-actions.md` |

A worker must not edit another worker's fragment or narrative. Cross-domain candidates are recorded in the worker output as transfer/proposed-overlap dispositions and resolved centrally during integration.

## Canonical requirement schema

```yaml
id: RUNTIME-...
title:
domain:
aliases:
  - A-REQ-...
affected_fighters: []
established_behavior:
  summary:
  invariants: []
  source_boundaries: []
runtime_contract:
  inputs: []
  preconditions: []
  transition:
  outputs: []
  failure_behavior:
state_requirements: []
command_event_implications: []
choice_implications: []
visibility_implications: []
persistence_implications: []
test_scenarios:
  - id:
    given: []
    when:
    then: []
status: draft | deterministic | blocked
```

## Alias dispositions

Every owner alias must receive exactly one disposition:

- `canonical` — mapped to one canonical requirement;
- `covered_by_existing` — already fully expressed by a frozen generic contract;
- `split` — source alias contains multiple independent canonical capabilities, with every part enumerated;
- `blocked` — deterministic runtime behavior cannot yet be defined, with a concrete blocker;
- `out_of_scope` — not a shared runtime requirement, with the destination artifact named.

Aliases may not disappear from the registry. A wording rename is not completion.

## Integration rules

1. Read the source owner definition and every manifest that references it.
2. Preserve established published behavior before proposing abstractions.
3. Merge aliases only when state, timing, choice, visibility and persistence semantics agree.
4. Split aliases when one owner definition bundled independent capabilities.
5. Prefer a reusable typed capability over an extension handler.
6. Record every material boundary and test it using at least two affected fighter/card examples when the corpus provides them.
7. Keep unresolved policy decisions explicit.

## Extension-handler gate

A handler is permitted only when all conditions hold:

- the generic operation/capability model cannot represent the behavior without losing deterministic semantics;
- the handler contract has typed inputs, events, state, choices, visibility and persistence behavior;
- core code does not test fighter/card IDs;
- registration occurs through data/capability metadata;
- independent QA approves the exception.

## Phase 4C gate

- 52/52 owner aliases have one complete disposition;
- no alias maps to multiple competing implementations;
- canonical requirements are deterministic for launch scope;
- state, command/event, choice, visibility and persistence implications are explicit;
- required fixtures are enumerated;
- no undocumented character-specific semantics remain;
- fresh independent QA passes before merge to `main`.
