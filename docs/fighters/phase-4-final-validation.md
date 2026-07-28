# Phase 4 Final Validation

## Verdict

**PASS — final Phase 4 gate closed.**

```text
integration_branch: phase-4-final-integration
validated_corpus_head: f9b718c7518e9cfaabe73ecad6a235dc83ae6ea9
workflow: Phase 4 Corpus Validation
workflow_run: 30316612762
job: 90143389177
machine_result: docs/fighters/phase-4-final-validation.json
errors: 0
warnings: 0
```

The machine result was produced after fresh independent correction QA, controlled integration of all five correction scopes, YAML syntax normalization and owner-requirement parser normalization. Later commits that add this record and update roadmap prose do not alter fighter/card corpus semantics.

## Integrated corpus

| Metric | Result |
|---|---:|
| Canonical competitive fighter identities | 74 |
| Fighter manifests | 74 |
| Paired card manifests | 74 |
| Unique action-card definitions | 926 |
| Available action-card copies | 2214 |
| Owner requirement definitions | 52 |
| Referenced owner requirements | 52 |
| Unresolved requirement references | 0 |

Canonical identity checks passed:

- Bruce Lee occurs once as one canonical fighter lineage;
- `yennefer-and-triss` did not replace canonical `yennefer-triss`;
- `jekyll-hyde` did not replace canonical `jekyll-and-hyde`;
- `little-red` did not replace canonical `little-red-riding-hood`;
- `dr-sattler` did not replace canonical `dr-ellie-sattler`.

## Structural validation

The integrated validator passed all of the following:

- duplicate-key-aware YAML parse for every fighter/card manifest;
- one unique canonical fighter ID per identity;
- exact fighter/card manifest pairing;
- unique card IDs inside each fighter manifest;
- positive integer quantities and valid card types;
- `usable_by` references restricted to fighters in the local topology;
- fighter manifest deck paths resolving to the paired card manifest;
- valid fixed/constructed/choose-group count relationships;
- all referenced `A/B/C/D-REQ-*` IDs resolving to owner report definitions;
- status/verification consistency for Phase 4B manifests;
- no evidence archive alias promoted to a canonical identity.

The validator also audits local choice/operation-result references. No warnings remained in the final run.

## Deck validation

Every manifest's action-card quantity sum matches its declared available pool.

The required nonstandard constructions were preserved:

| Fighter | Available pool | Game deck | Construction |
|---|---:|---:|---|
| Daredevil | 22 | 22 | fixed |
| Elektra | 20 | 20 | fixed |
| Black Widow | 31 | 31 | fixed |
| Geralt of Rivia | 36 | 30 | choose groups |
| Buffy | 35 | 30 | choose matching sidekick group |

Additional boundaries remain explicit:

- Titania's Glamours are outside the ordinary action deck;
- Wayward Sisters spells are external definitions, not action-card copies;
- Pandora's Miseries remain a separate auxiliary system where declared;
- bonus attacks and other external gameplay definitions are not included in action-deck totals;
- non-card components are not counted as card definitions.

## Physical-card correction QA

Fresh independent QA passed all five correction scopes before integration:

| Scope | Correction head | Verdict |
|---|---|---|
| Phase 4A | `59dcc54d973b243399f86da8aae55beae6cacc05` | PASS |
| Worker A | `ef67563ed8fef7f935b03f26bf800b3125367272` | PASS |
| Worker B | `85dd4c8d1db943ac60b1ded8f6955aec661ab472` | PASS |
| Worker C | `9c3e76c12d4ef6abb25d4405a386f47c03fb045e` | PASS |
| Worker D | `a449fa38d79e5cccd14d61517e1dcd0e3c0f5f98` | PASS |

```text
unresolved physical-card correction P1: 0
unresolved physical-card correction P2: 0
scope violations: 0
foreign manifests changed: 0
```

Detailed evidence is recorded in [`../qa/phase-4-card-image/correction-integration-qa.md`](../qa/phase-4-card-image/correction-integration-qa.md). The fifteen original physical-card reports remain immutable evidence.

## Verification and policy state

The corpus intentionally does not mark every fighter `verified`.

| Status | Count |
|---|---:|
| verified | 28 |
| partial | 45 |
| blocked | 1 |

`partial` means published behavior is deterministic but one or more shared runtime capabilities or explicit evidence qualifications remain. It is not an invitation to add fighter-specific core-engine logic.

Deadpool is the only blocked fighter. Its published physical/social effects require a material digital-adaptation policy decision; no substitute behavior was invented during integration.

## Preserved qualifications and normalizations

- Loki multiplayer behavior remains explicitly identified as a project normalization where authoritative source coverage is insufficient.
- Titania's Glamours remain gameplay-relevant auxiliary cards even where physical-image evidence is incomplete.
- Krang's die remains a non-card component and is not declared action-card-verified without corresponding evidence.
- Little Red evidence-package aliases remain evidence-only; the three canonical card variants and canonical fighter ID were not renamed.
- unreadable or missing physical evidence remains a provenance qualification rather than being silently converted into `verified` physical evidence.
- official errata/rulings, non-card fighter abilities and external component semantics continue to outrank physical-card wording where the source policy requires it.

## Phase 4 gate

| Gate condition | Result |
|---|---|
| All 74 competitive fighters integrated | PASS |
| All fighter/card manifests structurally parseable | PASS |
| All deck constructions reconciled | PASS |
| Physical-card P1/P2 corrections closed | PASS |
| All owner requirements resolvable by ID | PASS |
| Evidence-only aliases separated | PASS |
| Verification/policy distinctions preserved | PASS |

**Phase 4B is complete.**

This verdict completes the canonical fighter/card corpus. It does not declare the overall project developer-ready. The next required gate is Phase 4C: consolidation of the 52 owner requirements into canonical generic runtime capabilities. Battlefield graph work, engine architecture and stable engine foundation may proceed in parallel where their contracts are already deterministic.
