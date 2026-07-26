# Fighter corpus

This directory contains normalized, implementation-oriented fighter manifests.

## Status

- **Phase 4A:** representative stress-test corpus — in progress.
- **Phase 4B:** complete released roster — planned after the Phase 4A schema gate passes.

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

## Copyright/content policy

The repository stores factual card metadata and **normalized gameplay semantics**, not wholesale rulebook/card prose. Card names, quantities, values and BOOST values are indexed for deterministic implementation; effects are represented as operations/conditions/choices and retain links to the published source.

Official published-deck data may be discovered through UmDb only under `/umdb/...`. Fan decks under `/decks/...`, including community balance patches, are excluded from the authoritative corpus.

See [`schema.md`](schema.md) for the manifest contract.