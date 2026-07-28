# Engine bootstrap development

## Scope

This scaffold proves the accepted Go-server and TypeScript-PWA toolchains while
preserving the server-authoritative dependency boundary. It does not implement
gameplay, select application frameworks, or define public transport schemas.

## Toolchains

```text
Go: 1.26.5 (module language version 1.26.0)
Node.js: 24 LTS
TypeScript: 6.0.2
PostgreSQL development image: 18.4-alpine
Python: 3.13 for repository validators
```

## Local checks

```bash
python3 scripts/validate_engine_bootstrap.py
gofmt -w apps internal
go test ./...
npm ci --ignore-scripts
npm run check
npm run build
docker compose up -d postgres
```

The committed npm lockfile pins the compiler artifact and workspace graph used
by CI. Generated build output and `node_modules` remain untracked.

## Dependency direction

```text
transport / persistence adapters
              ↓
      application orchestration
              ↓
        pure domain engine
```

`internal/domain` must not import HTTP, WebSocket, SQL, filesystem, application,
persistence, or transport packages. The validator enforces this boundary before
compilation.

## Deferred decisions

The scaffold does not choose:

- HTTP router or WebSocket library;
- TypeScript UI framework or component library;
- schema technology and code generator;
- PostgreSQL query or migration library;
- authentication provider;
- hosting platform.

Each choice is introduced only when a bounded vertical slice requires it.
