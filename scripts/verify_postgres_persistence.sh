#!/usr/bin/env bash
set -euo pipefail

project_name="unmatched-bootstrap-persistence-${GITHUB_RUN_ID:-$$}"
export POSTGRES_DB="${POSTGRES_DB:-unmatched}"
export POSTGRES_USER="${POSTGRES_USER:-unmatched}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-unmatched-local-only}"
export POSTGRES_PORT="${POSTGRES_PORT:-55432}"

compose() {
  docker compose -p "$project_name" "$@"
}

cleanup() {
  compose down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

cleanup
compose up -d --wait postgres
compose exec -T postgres psql \
  -v ON_ERROR_STOP=1 \
  -U "$POSTGRES_USER" \
  -d "$POSTGRES_DB" <<'SQL'
CREATE TABLE IF NOT EXISTS bootstrap_persistence_probe (
  id integer PRIMARY KEY,
  value text NOT NULL
);
INSERT INTO bootstrap_persistence_probe (id, value)
VALUES (1, 'survives-recreate')
ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value;
SQL

compose down --remove-orphans
compose up -d --wait --force-recreate postgres

probe_value="$(
  compose exec -T postgres psql \
    -v ON_ERROR_STOP=1 \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB" \
    -Atc "SELECT value FROM bootstrap_persistence_probe WHERE id = 1"
)"

if [[ "$probe_value" != "survives-recreate" ]]; then
  echo "PostgreSQL named-volume persistence check failed: got '$probe_value'" >&2
  exit 1
fi

echo "PostgreSQL health and named-volume persistence: PASS"
