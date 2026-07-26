# Source registry

**Last reviewed:** 2026-07-26  
**Status:** verified entry-point registry; exhaustive per-release bibliography completed in Phase 3.

This file records the stable source hierarchy and canonical external entry points used across the project. Detailed release-specific provenance now lives in:

- [`../sets/registry.yaml`](../sets/registry.yaml) — canonical release/content inventory;
- [`../sets/source-bibliography.md`](../sets/source-bibliography.md) — exhaustive Phase 3 release/source map;
- [`../sets/mechanics-index.md`](../sets/mechanics-index.md) — set-specific mechanic → authoritative source mapping;
- [`../sets/phase-3-validation.md`](../sets/phase-3-validation.md) — completeness evidence.

The authority and conflict policy remains [`source-policy.md`](source-policy.md).

## Tier 1 — current core rules and rulings

### Current competitive Core Rules

- URL: https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Core-rules_EN_Light.pdf
- Host: IELLO.
- Role: current generic competitive rules, including two-player and multiplayer deltas.
- Status: verified.

### Unmatched Rulings Archive

- URL: https://docs.google.com/document/d/13b-FbPq_vuqcc3IokeHvQ2ctJaDNZZuUaZmt4uft5h0/
- Role: living publisher/designer ruling and errata compilation.
- Authority note: Unmatched Reference identifies the archive as considered official by Restoration Games.
- Usage: individual imported rulings retain their own provenance rather than citing this registry indirectly.

### Official publisher/addendum material

- Restoration Games: https://restorationgames.com/
- IELLO: https://iellogames.com/

Official product pages, rulebooks, set rules, errata and addenda outrank secondary databases for gameplay semantics.

## Tier 2 — consolidated reference

### Unmatched Reference v10.0

- Version: October 2025.
- URL: https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf
- Role: consolidated index of official errata, major rulings, character rules, card effects, battlefield effects and set reference through its publication scope.
- Limitation: some text is editorial synthesis; use it for discovery/cross-checking and trace sensitive behavior back to official rulebooks/rulings.
- Freshness: 2026 releases and addenda require newer official sources.

## Tier 3 — normalized published-content indexes

### UmDb published decks

- URL: https://unmatched.cards/umdb/decks/stats
- Role: normalized published fighter/deck data for Phase 4.
- Important boundary: `/umdb/...` is the published database; `/decks/...` is fan content and is never imported into the official corpus.

### UmDb sets

- URL: https://unmatched.cards/umdb/sets
- Role: release/set cross-checking.

### UmDb battlefields

- URL: https://unmatched.cards/umdb/boards
- Role: battlefield discovery/cross-checking.
- Limitation: not proof of engine graph correctness; Phase 5 independently verifies nodes, adjacency, zones, starts and special connections.

## Tier 4 — secondary discovery indexes

### The Unmatched Club

- URL: https://www.the-unmatched.club/
- Role: current roster/rulings discovery and freshness cross-checks.
- Rule: follow cited rulings back to official/publisher authority where possible.

### BoardGameGeek Unmatched Game System index

- URL: https://boardgamegeek.com/wiki/page/thing%3A295564%3Amoreinfo
- Role: historical completeness cross-check, especially for retired licensed releases no longer represented in the current Restoration storefront.
- Rule: may identify a missing historical fighter/battlefield/product, but does not override an official rulebook or ruling.

## Phase 3 canonical outputs

Phase 3 replaced the former provisional set-source list with an exhaustive registry.

The current corpus distinguishes:

- released primary products and historical standalone fighters;
- official competitive-content supplements such as Shredder/Krang Hero Decks and Nova High;
- reprints/returns as edition metadata rather than duplicate gameplay records;
- Adventures competitive-compatible hero content from deferred cooperative enemy/scenario logic;
- announced but incomplete content such as Hellboy from the released corpus.

See [`../sets/README.md`](../sets/README.md) for the counting and identity policy.

## Freshness cases requiring explicit provenance

### Stars & Stripes

- Product: https://restorationgames.com/shop/unmatched-stars-and-stripes/
- Rulebook: https://restorationgames.com/wp-content/uploads/2025/08/UM-Stars-Stripes_Rulebook-FLAT.pdf
- White House Secret Passages addendum: https://restorationgames.com/wp-content/uploads/2026/03/UM-SnS-rules-addendum_color.pdf

The separate addendum is first-class authority; the project must not rewrite history as though the omitted Secret Passages rule appeared in the original rulebook.

### TMNT

- Product: https://restorationgames.com/shop/unmatched-adventures-teenage-mutant-ninja-turtles/
- Competitive Shredder/Krang Hero Decks: https://restorationgames.com/shop/unmatched-adventures-teenage-mutant-ninja-turtles-shredder-krang-hero-decks/

The Turtle heroes belong to the competitive corpus. Cooperative villain/minion/scenario behavior remains deferred. The Hero Deck expansion separately promotes Shredder and Krang into competitive fighter definitions.

### Nova High

Official battlefield existence is registered, but graph/topology/equivalence remains a Phase 5 verification item. Do not promote community reskin/equivalence claims into canonical data without authoritative map evidence.

### Deadpool

The historical released fighter is registered even though a current Restoration product/rulebook page is not publicly indexed. Release membership is cross-checked with archival publisher-version metadata and UmDb; exact deck semantics require Phase 4 verification.

### Hellboy

Officially announced, but no complete final playable roster/rulebook/card corpus is currently available. It remains `announced/blocked`; community-predicted characters or cards must not enter the canonical corpus.

## Maintenance rule

When a new official release, erratum or addendum appears:

1. update the relevant record in `../sets/registry.yaml`;
2. add the exact source to `../sets/source-bibliography.md`;
3. update `../sets/mechanics-index.md` if a new mechanic family appears;
4. preserve older edition/source provenance rather than overwriting it;
5. reopen the relevant later-phase validation if executable fighter/card/battlefield semantics change.