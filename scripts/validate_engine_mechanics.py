#!/usr/bin/env python3
"""Validate the full-corpus engine mechanics capability map."""

from __future__ import annotations

import json
import re
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MAP_PATH = ROOT / "docs/engine/mechanics-analysis/capability-map.json"
REQUIRED_DOCS = (
    ROOT / "docs/engine/mechanics-analysis/README.md",
    ROOT / "docs/engine/mechanics-analysis/runtime-blueprint.md",
    ROOT / "docs/engine/mechanics-analysis/implementation-waves.md",
)

EXPECTED_ALIASES = {
    *(f"A-REQ-{number:03d}" for number in range(3, 16)),
    *(f"B-REQ-{number:03d}" for number in range(1, 11)),
    *(f"C-REQ-{number:03d}" for number in (1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13)),
    *(f"D-REQ-{number:03d}" for number in range(1, 18)),
}

CAPABILITY_ID_RE = re.compile(r"^CAP-[0-9]{3}$")
SLUG_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
ALLOWED_STATUSES = {"deterministic", "blocked"}


def error(errors: list[str], message: str) -> None:
    errors.append(message)


def validate_cycle(capabilities: dict[str, dict[str, object]], errors: list[str]) -> None:
    visiting: set[str] = set()
    visited: set[str] = set()

    def visit(capability_id: str, chain: list[str]) -> None:
        if capability_id in visited:
            return
        if capability_id in visiting:
            error(errors, f"capability dependency cycle: {' -> '.join([*chain, capability_id])}")
            return
        visiting.add(capability_id)
        capability = capabilities[capability_id]
        for dependency in capability["depends_on"]:
            if dependency in capabilities:
                visit(dependency, [*chain, capability_id])
        visiting.remove(capability_id)
        visited.add(capability_id)

    for capability_id in capabilities:
        visit(capability_id, [])


def main() -> int:
    errors: list[str] = []

    for path in REQUIRED_DOCS:
        if not path.is_file():
            error(errors, f"required mechanics document is missing: {path.relative_to(ROOT)}")

    if not MAP_PATH.is_file():
        error(errors, f"capability map is missing: {MAP_PATH.relative_to(ROOT)}")
        for message in errors:
            print(f"ERROR: {message}", file=sys.stderr)
        return 1

    try:
        document = json.loads(MAP_PATH.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        print(f"ERROR: cannot parse {MAP_PATH.relative_to(ROOT)}: {exc}", file=sys.stderr)
        return 1

    sources = document.get("sources", {})
    expected_source_values = {
        "fighters": 74,
        "unique_action_card_definitions": 926,
        "owner_requirement_aliases": 52,
    }
    for key, expected in expected_source_values.items():
        if sources.get(key) != expected:
            error(errors, f"sources.{key} must be {expected}, got {sources.get(key)!r}")

    raw_capabilities = document.get("capabilities")
    if not isinstance(raw_capabilities, list):
        error(errors, "capabilities must be a list")
        raw_capabilities = []

    ids: list[str] = []
    slugs: list[str] = []
    primary_aliases: list[str] = []
    capability_by_id: dict[str, dict[str, object]] = {}

    required_fields = {
        "id",
        "slug",
        "kernel",
        "primary_aliases",
        "cross_domain_aliases",
        "depends_on",
        "wave",
        "status",
    }

    for index, capability in enumerate(raw_capabilities):
        location = f"capabilities[{index}]"
        if not isinstance(capability, dict):
            error(errors, f"{location} must be an object")
            continue

        missing = sorted(required_fields - capability.keys())
        if missing:
            error(errors, f"{location} missing fields: {', '.join(missing)}")
            continue

        capability_id = capability["id"]
        slug = capability["slug"]
        status = capability["status"]
        wave = capability["wave"]

        if not isinstance(capability_id, str) or not CAPABILITY_ID_RE.fullmatch(capability_id):
            error(errors, f"{location}.id is invalid: {capability_id!r}")
            continue
        if not isinstance(slug, str) or not SLUG_RE.fullmatch(slug):
            error(errors, f"{location}.slug is invalid: {slug!r}")
        if status not in ALLOWED_STATUSES:
            error(errors, f"{location}.status is invalid: {status!r}")
        if not isinstance(wave, int) or not 1 <= wave <= 7:
            error(errors, f"{location}.wave must be an integer from 1 through 7")

        ids.append(capability_id)
        slugs.append(slug)
        capability_by_id[capability_id] = capability

        for field in ("primary_aliases", "cross_domain_aliases", "depends_on"):
            value = capability[field]
            if not isinstance(value, list) or not all(isinstance(item, str) for item in value):
                error(errors, f"{location}.{field} must be a list of strings")

        if isinstance(capability["primary_aliases"], list):
            primary_aliases.extend(capability["primary_aliases"])

    for name, values in (("capability IDs", ids), ("capability slugs", slugs)):
        duplicates = sorted(value for value, count in Counter(values).items() if count > 1)
        if duplicates:
            error(errors, f"duplicate {name}: {duplicates}")

    alias_duplicates = sorted(alias for alias, count in Counter(primary_aliases).items() if count > 1)
    missing_aliases = sorted(EXPECTED_ALIASES - set(primary_aliases))
    unknown_aliases = sorted(set(primary_aliases) - EXPECTED_ALIASES)
    if alias_duplicates:
        error(errors, f"duplicate primary aliases: {alias_duplicates}")
    if missing_aliases:
        error(errors, f"unmapped primary aliases: {missing_aliases}")
    if unknown_aliases:
        error(errors, f"unknown primary aliases: {unknown_aliases}")

    coverage = document.get("coverage", {})
    if coverage.get("expected_primary_aliases") != len(EXPECTED_ALIASES):
        error(errors, "coverage.expected_primary_aliases is stale")
    if coverage.get("mapped_primary_aliases") != len(primary_aliases):
        error(errors, "coverage.mapped_primary_aliases is stale")
    if coverage.get("duplicate_primary_aliases") != alias_duplicates:
        error(errors, "coverage.duplicate_primary_aliases is stale")
    if coverage.get("unmapped_primary_aliases") != missing_aliases:
        error(errors, "coverage.unmapped_primary_aliases is stale")

    for capability_id, capability in capability_by_id.items():
        capability_wave = capability["wave"]
        for dependency in capability["depends_on"]:
            if dependency == capability_id:
                error(errors, f"{capability_id} depends on itself")
                continue
            dependency_capability = capability_by_id.get(dependency)
            if dependency_capability is None:
                error(errors, f"{capability_id} has unknown dependency {dependency}")
                continue
            if dependency_capability["wave"] > capability_wave:
                error(
                    errors,
                    f"{capability_id} wave {capability_wave} precedes dependency "
                    f"{dependency} wave {dependency_capability['wave']}",
                )

        unknown_cross_domain = sorted(set(capability["cross_domain_aliases"]) - EXPECTED_ALIASES)
        if unknown_cross_domain:
            error(errors, f"{capability_id} has unknown cross-domain aliases: {unknown_cross_domain}")

    validate_cycle(capability_by_id, errors)

    blocked = [
        (capability_id, capability["primary_aliases"])
        for capability_id, capability in capability_by_id.items()
        if capability["status"] == "blocked"
    ]
    if blocked != [("CAP-020", ["A-REQ-015"])]:
        error(errors, f"blocked capability policy must be only CAP-020/A-REQ-015, got {blocked!r}")

    if errors:
        for message in errors:
            print(f"ERROR: {message}", file=sys.stderr)
        return 1

    print(
        "Engine mechanics capability map validated: "
        f"{len(capability_by_id)} capabilities, "
        f"{len(primary_aliases)}/{len(EXPECTED_ALIASES)} primary aliases, "
        "dependency graph acyclic."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
