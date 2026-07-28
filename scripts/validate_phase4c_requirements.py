#!/usr/bin/env python3
"""Validate Phase 4C owner-requirement discovery, ownership and registry state."""

from __future__ import annotations

import json
import re
import sys
from collections import Counter
from pathlib import Path
from typing import Any

import yaml

ROOT = Path(__file__).resolve().parents[1]
REPORTS = [ROOT / f"docs/phase-4b/worker-{owner}-report.md" for owner in "abcd"]
OWNERSHIP_PATH = ROOT / "docs/mechanics/runtime-requirements/alias-ownership.yaml"
REGISTRY_PATH = ROOT / "docs/mechanics/runtime-requirements/registry.yaml"
OUTPUT_PATH = ROOT / "artifacts/phase-4c-requirement-validation.json"
REQ_PATTERN = re.compile(r"\b[ABCD]-REQ-\d{3}\b")
EXPECTED_DOMAINS = {
    "resolution",
    "fighter_presence",
    "movement_combat",
    "damage_modifiers",
    "cards_actions",
    "central_policy",
}
WORKER_DOMAINS = EXPECTED_DOMAINS - {"central_policy"}


def load_yaml(path: Path, errors: list[str]) -> dict[str, Any]:
    try:
        value = yaml.safe_load(path.read_text(encoding="utf-8"))
    except Exception as exc:  # noqa: BLE001
        errors.append(f"{path.relative_to(ROOT)}: YAML parse failed: {exc}")
        return {}
    if not isinstance(value, dict):
        errors.append(f"{path.relative_to(ROOT)}: expected top-level mapping")
        return {}
    return value


def source_aliases(errors: list[str]) -> set[str]:
    aliases: set[str] = set()
    per_report: dict[str, list[str]] = {}
    for report in REPORTS:
        if not report.exists():
            errors.append(f"missing owner report: {report.relative_to(ROOT)}")
            continue
        found = sorted(set(REQ_PATTERN.findall(report.read_text(encoding="utf-8"))))
        per_report[report.stem] = found
        aliases.update(found)
    expected_by_owner = {"A": 13, "B": 10, "C": 12, "D": 17}
    for owner, expected in expected_by_owner.items():
        actual = sum(1 for alias in aliases if alias.startswith(f"{owner}-REQ-"))
        if actual != expected:
            errors.append(f"owner {owner}: expected {expected} requirements, found {actual}")
    return aliases


def string_list(value: Any, label: str, errors: list[str]) -> list[str]:
    if not isinstance(value, list) or not all(isinstance(item, str) for item in value):
        errors.append(f"{label}: expected list of strings")
        return []
    return list(value)


def main() -> int:
    errors: list[str] = []
    warnings: list[str] = []

    defined = source_aliases(errors)
    ownership = load_yaml(OWNERSHIP_PATH, errors)
    registry = load_yaml(REGISTRY_PATH, errors)

    assignments = ownership.get("assignments", {})
    if not isinstance(assignments, dict):
        errors.append("alias-ownership.yaml: assignments must be a mapping")
        assignments = {}

    assignment_domains = set(assignments)
    if assignment_domains != EXPECTED_DOMAINS:
        errors.append(
            "alias-ownership.yaml: assignment domains mismatch: "
            f"missing={sorted(EXPECTED_DOMAINS - assignment_domains)}, "
            f"unexpected={sorted(assignment_domains - EXPECTED_DOMAINS)}"
        )

    primary_pairs: list[tuple[str, str]] = []
    per_domain_counts: dict[str, int] = {}
    for domain, record in assignments.items():
        if not isinstance(record, dict):
            errors.append(f"assignments.{domain}: expected mapping")
            continue
        aliases = string_list(record.get("aliases"), f"assignments.{domain}.aliases", errors)
        per_domain_counts[domain] = len(aliases)
        primary_pairs.extend((alias, domain) for alias in aliases)
        if domain in WORKER_DOMAINS:
            if not isinstance(record.get("issue"), int):
                errors.append(f"assignments.{domain}.issue: expected integer")
            if not isinstance(record.get("branch"), str):
                errors.append(f"assignments.{domain}.branch: expected string")

    primary_aliases = [alias for alias, _ in primary_pairs]
    duplicate_primary = sorted(alias for alias, count in Counter(primary_aliases).items() if count > 1)
    primary_set = set(primary_aliases)

    if duplicate_primary:
        errors.append(f"duplicate primary alias assignments: {duplicate_primary}")
    if primary_set != defined:
        errors.append(
            "primary assignment coverage mismatch: "
            f"unassigned={sorted(defined - primary_set)}, "
            f"unknown={sorted(primary_set - defined)}"
        )

    declared_count = ownership.get("source_alias_count")
    if declared_count != len(defined):
        errors.append(f"source_alias_count {declared_count!r} != discovered {len(defined)}")

    declared_counts = ownership.get("primary_ownership_counts", {})
    if not isinstance(declared_counts, dict):
        errors.append("primary_ownership_counts must be a mapping")
    else:
        normalized_declared = {str(key): value for key, value in declared_counts.items()}
        if normalized_declared != per_domain_counts:
            errors.append(
                "primary ownership counts mismatch: "
                f"declared={normalized_declared}, actual={per_domain_counts}"
            )

    cross_rules = ownership.get("cross_domain_rules", [])
    if not isinstance(cross_rules, list):
        errors.append("cross_domain_rules must be a list")
        cross_rules = []
    cross_aliases: set[str] = set()
    primary_lookup = dict(primary_pairs)
    for index, rule in enumerate(cross_rules):
        label = f"cross_domain_rules[{index}]"
        if not isinstance(rule, dict):
            errors.append(f"{label}: expected mapping")
            continue
        alias = rule.get("alias")
        primary = rule.get("primary")
        contributors = string_list(rule.get("contributors"), f"{label}.contributors", errors)
        if not isinstance(alias, str) or alias not in defined:
            errors.append(f"{label}.alias: unknown alias {alias!r}")
            continue
        if alias in cross_aliases:
            errors.append(f"duplicate cross-domain rule for {alias}")
        cross_aliases.add(alias)
        if primary_lookup.get(alias) != primary:
            errors.append(
                f"{label}: primary {primary!r} does not match assignment {primary_lookup.get(alias)!r}"
            )
        for contributor in contributors:
            if contributor not in WORKER_DOMAINS:
                errors.append(f"{label}: invalid contributor domain {contributor!r}")
            if contributor == primary:
                errors.append(f"{label}: contributor duplicates primary {primary!r}")

    validation = ownership.get("validation", {})
    if isinstance(validation, dict):
        if validation.get("primary_assignment_count") != len(primary_aliases):
            errors.append("ownership validation.primary_assignment_count is stale")
        if validation.get("duplicate_primary_assignments") != len(duplicate_primary):
            errors.append("ownership validation.duplicate_primary_assignments is stale")
        if validation.get("unassigned_aliases") != sorted(defined - primary_set):
            errors.append("ownership validation.unassigned_aliases is stale")
    else:
        errors.append("ownership validation must be a mapping")

    coverage = registry.get("coverage", {})
    if not isinstance(coverage, dict):
        errors.append("registry coverage must be a mapping")
        coverage = {}
    if coverage.get("expected_owner_aliases") != len(defined):
        errors.append("registry coverage.expected_owner_aliases is stale")

    fragments = registry.get("worker_fragments", {})
    if not isinstance(fragments, dict):
        errors.append("registry worker_fragments must be a mapping")
    elif set(fragments) != WORKER_DOMAINS:
        errors.append(
            "registry worker fragment domains mismatch: "
            f"missing={sorted(WORKER_DOMAINS - set(fragments))}, "
            f"unexpected={sorted(set(fragments) - WORKER_DOMAINS)}"
        )

    dispositions = registry.get("alias_dispositions", [])
    if not isinstance(dispositions, list):
        errors.append("registry alias_dispositions must be a list")
        dispositions = []
    disposition_aliases: list[str] = []
    for index, disposition in enumerate(dispositions):
        if not isinstance(disposition, dict) or not isinstance(disposition.get("alias"), str):
            errors.append(f"registry alias_dispositions[{index}]: missing string alias")
            continue
        disposition_aliases.append(disposition["alias"])
    duplicate_dispositions = sorted(
        alias for alias, count in Counter(disposition_aliases).items() if count > 1
    )
    if duplicate_dispositions:
        errors.append(f"duplicate registry alias dispositions: {duplicate_dispositions}")
    if dispositions:
        disposition_set = set(disposition_aliases)
        if disposition_set != defined:
            errors.append(
                "registry disposition coverage mismatch: "
                f"missing={sorted(defined - disposition_set)}, "
                f"unknown={sorted(disposition_set - defined)}"
            )
    else:
        warnings.append("canonical alias dispositions are pending worker output (0/52)")

    result = {
        "schema_version": 1,
        "phase": "phase-4c",
        "passed": not errors,
        "source_aliases": len(defined),
        "primary_assignments": len(primary_aliases),
        "primary_counts": dict(sorted(per_domain_counts.items())),
        "duplicate_primary_assignments": duplicate_primary,
        "unassigned_aliases": sorted(defined - primary_set),
        "unknown_assignments": sorted(primary_set - defined),
        "cross_domain_aliases": sorted(cross_aliases),
        "registry_dispositions": len(disposition_aliases),
        "errors": errors,
        "warnings": warnings,
    }
    OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)
    OUTPUT_PATH.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0 if not errors else 1


if __name__ == "__main__":
    sys.exit(main())
