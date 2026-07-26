# unmatched-web

Mobile-first web implementation of the **Unmatched** game system.

The repository is intentionally **documentation-first**. Implementation should not begin until the competitive rules corpus, fighter/deck manifests, battlefield graphs, rulings, and engine-facing semantics reach the `developer-ready` gate defined in [`docs/specification-readiness.md`](docs/specification-readiness.md).

## Current phase

**Rules and content specification.** No application architecture or framework choice is authoritative yet.

## Documentation

- [`docs/README.md`](docs/README.md) — documentation map and status
- [`docs/research-plan.md`](docs/research-plan.md) — staged research and specification program
- [`docs/sources/source-policy.md`](docs/sources/source-policy.md) — source hierarchy and conflict-resolution rules
- [`docs/sources/source-registry.md`](docs/sources/source-registry.md) — canonical external sources and known gaps
- [`docs/specification-readiness.md`](docs/specification-readiness.md) — definition of done before implementation

## Scope boundary

The first specification target is **competitive Unmatched**. Unmatched Adventures cooperative rules are tracked as a later extension because villains, minions, objectives, initiative decks, and scenario-specific logic introduce a separate engine layer.

## Content policy

The specification should model rules and gameplay semantics in structured, implementation-oriented form. Do not treat artwork, graphic design, licensed character assets, or wholesale reproduction of copyrighted card/rulebook text as implementation dependencies.
