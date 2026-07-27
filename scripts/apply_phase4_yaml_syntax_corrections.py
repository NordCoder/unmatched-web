#!/usr/bin/env python3
"""Apply syntax-only corrections discovered by the integrated Phase 4 YAML gate.

The source branches contained several plain scalars that are invalid in YAML
flow mappings (`?` in names/URLs and indexed binding expressions) plus one
block scalar containing `: `. This migration quotes those existing scalar
values without changing gameplay data, IDs, quantities, or semantics.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
YAML_ROOTS = [ROOT / "docs/cards/phase-4a", ROOT / "docs/cards/phase-4b"]


def replace_exact(path: Path, old: str, new: str) -> bool:
    text = path.read_text(encoding="utf-8")
    count = text.count(old)
    if count != 1:
        raise RuntimeError(f"{path.relative_to(ROOT)}: expected exactly one occurrence of {old!r}, found {count}")
    path.write_text(text.replace(old, new), encoding="utf-8")
    return True


def main() -> int:
    changed: set[Path] = set()

    changed.add(ROOT / "docs/cards/phase-4a/geralt-of-rivia.yaml")
    replace_exact(
        ROOT / "docs/cards/phase-4a/geralt-of-rivia.yaml",
        "    name: GEAR: Sword of Steel",
        '    name: "GEAR: Sword of Steel"',
    )

    changed.add(ROOT / "docs/cards/phase-4b/wyatt-earp.yaml")
    replace_exact(
        ROOT / "docs/cards/phase-4b/wyatt-earp.yaml",
        "name: You Just Gonna Stand There and Bleed?",
        'name: "You Just Gonna Stand There and Bleed?"',
    )

    # Quote indexed binding expressions used as scalar values inside flow maps.
    indexed_scalar = re.compile(r"(?P<prefix>\b(?:ref|target): )(?P<value>[A-Za-z0-9_.-]+\[\d+\])")
    # Quote URL query strings inside flow maps. A plain '?' starts a YAML key.
    query_url = re.compile(r"(?P<prefix>\burl: )(?P<value>https?://[^,}\s]+\?[^,}\s]+)")

    for directory in YAML_ROOTS:
        for path in sorted(directory.glob("*.yaml")):
            text = path.read_text(encoding="utf-8")
            updated = indexed_scalar.sub(lambda match: f'{match.group("prefix")}\"{match.group("value")}\"', text)
            updated = query_url.sub(lambda match: f'{match.group("prefix")}\"{match.group("value")}\"', updated)
            if updated != text:
                path.write_text(updated, encoding="utf-8")
                changed.add(path)

    parse_errors: list[str] = []
    for directory in YAML_ROOTS:
        for path in sorted(directory.glob("*.yaml")):
            try:
                yaml.safe_load(path.read_text(encoding="utf-8"))
            except Exception as exc:  # noqa: BLE001 - migration must report every parser failure
                parse_errors.append(f"{path.relative_to(ROOT)}: {exc}")

    if parse_errors:
        print("YAML syntax correction did not produce a parse-clean card corpus:", file=sys.stderr)
        for error in parse_errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    print("Syntax-only YAML corrections applied:")
    for path in sorted(changed):
        print(f"- {path.relative_to(ROOT)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
