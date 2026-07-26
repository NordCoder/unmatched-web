# Canonical set and release registry

**Phase:** 3 — Canonical set and release registry  
**Status:** verified; see `phase-3-validation.md`  
**Last verified:** 2026-07-26

This directory is the canonical release/content inventory for competitive-compatible Unmatched material. It answers **what official content exists and which authoritative sources define its set-level behavior**. It does not yet transcribe every fighter card (Phase 4) or every battlefield edge/zone (Phase 5).

## Scope and counting policy

The registry deliberately separates three concepts:

1. **primary product records** — released boxed sets and historical single-fighter releases;
2. **competitive-content supplements** — official releases that add a fighter/deck or battlefield without being a normal standalone set;
3. **announced products** — tracked for freshness but blocked until an authoritative playable corpus exists.

Reprints and editions do **not** create duplicate gameplay records. They are attached to the same canonical product with edition metadata. Example: the 2025/2026 Cobble & Fog return from the Vault is the same canonical set as the 2020 release, while its clarified rulebook is retained as later provenance.

As of the verification date the registry contains:

- **25 released primary product records** through `stars-and-stripes`;
- **2 released competitive-content supplements**: `tmnt-shredder-krang-hero-decks` and `nova-high-battlefield-mat`;
- **1 announced/blocked product**: `hellboy`.

The 25-record number is a project registry count, not a claim about how Restoration Games chooses to market/count boxes. Restoration stated in November 2025 that it had published 24 Unmatched sets; Stars & Stripes is the next released primary set in 2026.

## Files

- [`registry.yaml`](registry.yaml) — machine-oriented canonical product/set inventory.
- [`mechanics-index.md`](mechanics-index.md) — set/battlefield mechanics that require rules beyond the generic core.
- [`source-bibliography.md`](source-bibliography.md) — official and secondary source map used to verify the registry.
- [`phase-3-validation.md`](phase-3-validation.md) — completeness and gate evidence.

## Required registry semantics

Every record has, where applicable:

```yaml
id: stable-project-id
kind: set | standalone_fighter | adventures_set | hero_deck_expansion | battlefield_accessory
release_status: released | announced
release_year: integer | null
competitive_compatible: true | false
license_kind: public_domain_or_historical | third_party | mixed_or_personality | original
availability: current | vaulted_or_oop | licensed_oop | announced | unknown
fighters: [canonical fighter ids]
battlefields: [canonical battlefield ids]
special_components: [semantic component names]
set_mechanics: [mechanic ids]
editions: []
sources: {}
verification: verified | partial | blocked
```

Exact release day is not required for engine behavior. When official sources disagree by channel (for example Kickstarter fulfillment versus retail availability), the registry preserves the known dates/notes rather than inventing one canonical date.

## Source discipline

The authority hierarchy in `docs/sources/source-policy.md` applies.

- Official Restoration/IELLO product pages, rulebooks, addenda and announcements establish released/announced status and special rules.
- UmDb and BoardGameGeek's Unmatched system index are used to cross-check historical roster/map completeness, especially for licensed products no longer present in the current Restoration catalog.
- A secondary source may identify a missing historical item but cannot silently define a gameplay rule when an official rulebook/addendum exists.
- `unmatched.cards/decks/...` fan decks are never part of this registry.

## Phase boundary

Phase 3 proves that the content universe is enumerated and that set-level mechanics have an authoritative entry point. It does **not** mean each fighter or battlefield is developer-ready.

Phase 4 must still verify fighter stats, deck construction, every action card, resources and fighter/card rulings. Phase 5 must convert every battlefield into a deterministic graph.