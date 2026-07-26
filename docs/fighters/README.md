# Fighter corpus

This directory contains normalized, implementation-oriented fighter manifests.

## Status

- **Phase 4A:** **complete — gate passed 2026-07-26.**
- **Phase 4B:** complete released roster — planned.

Phase 4 consumes the release queue from `docs/sets/registry.yaml` and the semantic framework from `docs/mechanics/`.

## Separation of concerns

A fighter manifest describes:

- fighter topology (single hero, hero + sidekick(s), multiple heroes, selectable hero);
- attack type and starting health for each independently tracked fighter;
- movement;
- special ability;
- setup/pre-game configuration;
- resources, token pools and persistent state;
- deck construction rule;
- fighter-level rulings and provenance.

Action-card definitions live under `docs/cards/`. A fighter manifest links its deck rather than duplicating every card.

## Phase 4A representative set

The stress-test corpus intentionally targets different engine pressures:

| Fighter | Why it is in the sample |
| --- | --- |
| Achilles | sidekick-defeat state change, combat participant replacement, bonus attack |
| Bloody Mary | action ordinal/history, start-space snapshot, bonus attack |
| Sun Wukong | summonable sidekicks, reserve pool, damage prevention/redirection, bonus attack |
| Sherlock Holmes | protected card effects, face-up prediction, private information |
| Dracula | multiple sidekicks, start-turn conditional damage/draw, returned defeated sidekick |
| Raptors | multiple heroes, shared deck, multi-hero defeat semantics |
| Wayward Sisters | multiple heroes, alternate discard destination, ingredient resource, external spell pool |
| Geralt of Rivia | pre-game deck construction from a 36-card pool, gear categories, ongoing schemes |
| Yennefer & Triss | setup-time hero selection, role-dependent health/ability, simultaneous/private choices |
| Black Panther | foreign-card storage, card ownership vs control/zone, BOOST from non-hand zone |

All ten have verified fighter manifests and full normalized card manifests under `phase-4a/`. See:

- [`schema.md`](schema.md) — promoted fighter/deck schema;
- [`phase-4a-mechanics.md`](phase-4a-mechanics.md) — corpus-driven semantic findings;
- [`phase-4a-validation.md`](phase-4a-validation.md) — gate evidence.

## Phase 4A result

The sample did not require opaque character-specific execution code. Real cards extended the generic effect vocabulary with reusable semantics including damage prevention/redirection, effective printed/BOOST value changes, operation prevention, deck reordering, captured parent context and combat-card replacement.

Phase 4B must now apply the same schema and validation contract to every released competitive fighter from `docs/sets/registry.yaml`.

## Copyright/content policy

The repository stores factual card metadata and **normalized gameplay semantics**, not wholesale rulebook/card prose. Card names, quantities, values and BOOST values are indexed for deterministic implementation; effects are represented as operations/conditions/choices and retain links to the published source.

Official published-deck data may be discovered through UmDb only under `/umdb/...`. Fan decks under `/decks/...`, including community balance patches, are excluded from the authoritative corpus.