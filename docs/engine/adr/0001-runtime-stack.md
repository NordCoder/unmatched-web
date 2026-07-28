# ADR 0001 — Runtime Stack

## Status

```text
status: proposed
owner_decision_required: true
parent_issue: #19
```

## Context

The product requires:

- a server-authoritative deterministic game engine;
- persistent event history, snapshots and reconnect;
- bidirectional low-latency match updates;
- player-specific hidden-information projections;
- a portrait-oriented browser/PWA client;
- data-driven loading of fighter/card/battlefield definitions;
- strong tests around pure reducers, legal-action generation and replay;
- incremental deployment by a small/solo development operation.

The repository is currently documentation-first and has no existing application stack to preserve.

## Decision drivers

1. deterministic, explicit domain modeling;
2. simple deployment and operational recovery;
3. mature PostgreSQL and WebSocket support;
4. good test performance for replay/property fixtures;
5. strong static typing without excessive framework coupling;
6. mobile-web ecosystem quality;
7. maintainability for a small team using AI-assisted development;
8. ability to generate transport contracts rather than hand-copy them;
9. developer familiarity and debugging ergonomics.

## Options

### Option A — Go authoritative server + TypeScript PWA

```text
Go modular monolith
PostgreSQL event/snapshot store
WebSocket or equivalent server-push transport
HTTP command/query API
TypeScript mobile PWA
schema-generated transport types
```

Advantages:

- clear separation between authoritative server and untrusted client;
- simple deployment artifact and predictable runtime behavior;
- strong concurrency/networking/database ecosystem;
- fast deterministic test/replay loops;
- explicit domain types without requiring a large application framework;
- TypeScript remains where browser/PWA tooling is strongest;
- aligns with current developer Go experience while preserving web-client productivity.

Costs/risks:

- domain/transport types cannot be shared directly as source code;
- code generation/schema discipline is mandatory;
- two language toolchains;
- browser-side UI state must not duplicate authoritative rules.

### Option B — TypeScript full stack

```text
TypeScript server
PostgreSQL
WebSocket transport
TypeScript PWA
shared generated or source-level types
```

Advantages:

- one language across server/client;
- easiest UI/API type sharing;
- broad web ecosystem and rapid iteration.

Costs/risks:

- stronger discipline is needed to keep mutable application objects out of the deterministic domain core;
- server/client code sharing can accidentally move legality logic into the client;
- framework churn and runtime dependency surface may be larger;
- CPU-heavy replay/property tests may require more tuning.

### Option C — Rust authoritative server + TypeScript PWA

Advantages:

- strongest compile-time ownership/state guarantees;
- efficient deterministic core and replay;
- explicit serialization/versioning possible.

Costs/risks:

- highest implementation and onboarding cost;
- slower product iteration for the current solo workflow;
- more friction around common web/database operations;
- risk that language complexity distracts from unresolved gameplay contracts.

## Proposed decision

Adopt **Option A: Go authoritative server + TypeScript PWA**, implemented initially as a modular monolith.

Proposed boundaries:

```text
/apps/server            authoritative API, match host, persistence adapters
/apps/web               mobile PWA; commands + projections only
/internal/domain        pure engine state, reducers, commands/events, legality
/internal/application   match orchestration and transactions
/internal/persistence   PostgreSQL event/snapshot/idempotency adapters
/internal/transport     HTTP/WebSocket projection delivery
/packages/contracts     generated schemas/types for transport payloads
```

The exact directory layout may change after repository bootstrapping; the architectural dependency direction is normative:

```text
transport/persistence adapters
        ↓
application orchestration
        ↓
pure domain engine
```

The domain engine must not import HTTP, WebSocket, SQL, filesystem or framework packages.

## Persistence proposal

Use PostgreSQL as the first durable store for:

- matches and current heads;
- command idempotency;
- append-only event batches;
- snapshots;
- player/seat authority;
- definition-version references.

Do not introduce Redis or a separate event broker until measured concurrency/recovery requirements justify it. Single-process in-memory match caching is an optimization above PostgreSQL, not the source of truth.

## Transport proposal

- authenticated HTTP endpoints for match lifecycle and commands;
- WebSocket (or equivalent ordered server-push channel) for projection updates;
- every message carries match revision/event sequence;
- reconnect can always fall back to an authoritative full projection;
- transport payloads are generated from versioned schemas.

The transport cannot expose raw event private payloads or authoritative state.

## Contract generation

Choose one versioned schema source for public transport contracts. Acceptable implementation choices include JSON Schema/OpenAPI or Protobuf, provided:

- Go and TypeScript types are generated;
- runtime validation exists at trust boundaries;
- durable event schemas are versioned separately from transient UI payloads;
- generated code is not edited manually.

The concrete schema technology is a follow-up ADR after command/projection payload shapes stabilize.

## Consequences

If accepted:

- production server implementation starts in Go;
- browser/PWA implementation starts in TypeScript;
- no direct source-level sharing of domain types is assumed;
- pure engine packages receive strict import/dependency checks;
- PostgreSQL-backed event/snapshot persistence is part of the first online vertical slice;
- deployment begins as one server service plus static/PWA delivery and PostgreSQL;
- scaling to multiple match hosts requires an explicit lease/ownership design rather than shared in-memory state.

## Non-decisions

This ADR does not select:

- a Go HTTP router/framework;
- a TypeScript UI framework;
- CSS/component libraries;
- migration/query libraries;
- hosting provider;
- WebSocket library;
- contract schema technology;
- authentication provider.

Those choices require smaller ADRs and must not alter the server-authoritative boundary.

## Acceptance criteria

The Owner may:

- **ACCEPT** — authorize Go server + TypeScript PWA bootstrapping;
- **REQUEST CHANGES** — identify a concrete constraint that changes the decision drivers;
- **REJECT** — select Option B or C with the accepted operational tradeoff.

No production framework bootstrap should be merged before this ADR is accepted. Language-neutral engine fixtures/contracts may continue in parallel.
