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
python3 -m unittest discover -s tests/architecture -p "test_*.py"
python3 scripts/validate_engine_bootstrap.py
gofmt -w apps internal
go test ./...
npm ci --ignore-scripts
npm run check
npm run build
docker compose up -d --wait postgres
bash scripts/verify_postgres_persistence.sh
```

The committed npm lockfile pins the compiler artifact and workspace graph used
by CI. Generated build output and `node_modules` remain untracked.

The PostgreSQL 18 service mounts its named volume at `/var/lib/postgresql` and
uses `/var/lib/postgresql/18/docker` as `PGDATA`. The persistence script runs in
an isolated Compose project, writes a probe row, recreates the container without
deleting the named volume, verifies the row, and removes the isolated resources.

## Dependency direction

```text
transport / persistence adapters
              ↓
      application orchestration
              ↓
        pure domain engine
```

`internal/domain` may import its own subpackages and suitable standard-library
packages. It must not import another internal layer, third-party packages, or
HTTP, WebSocket, SQL, filesystem, and related adapter packages. The validator
enforces this boundary before compilation.

The public contracts workspace accepts only the bootstrap generation sentinel
or files under `generated/` carrying both `@generated` and `DO NOT EDIT` markers.
Hand-written exported declarations are rejected.

Gameplay identifiers are loaded from the canonical fighter, card, and
battlefield YAML corpus and checked against implementation, validation, test,
workflow, and configuration sources. The validator does not embed roster IDs.

## Deferred decisions

The scaffold does not choose:

- HTTP router or WebSocket library;
- TypeScript UI framework or component library;
- schema technology and code generator;
- PostgreSQL query or migration library;
- authentication provider;
- hosting platform.

Each choice is introduced only when a bounded vertical slice requires it.
