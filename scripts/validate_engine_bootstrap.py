#!/usr/bin/env python3
"""Validate the bounded engine-bootstrap architecture contract."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

REQUIRED_PATHS = (
    "go.mod",
    "package.json",
    "apps/server/main.go",
    "apps/web/package.json",
    "apps/web/tsconfig.json",
    "apps/web/src/main.ts",
    "internal/domain/doc.go",
    "internal/application/doc.go",
    "internal/persistence/doc.go",
    "internal/transport/doc.go",
    "packages/contracts/package.json",
    "packages/contracts/tsconfig.json",
    "packages/contracts/src/index.ts",
    "compose.yaml",
    "tests/fixtures/README.md",
)

FORBIDDEN_DOMAIN_IMPORTS = (
    "database/sql",
    "io/fs",
    "net/http",
    "os",
    "path/filepath",
    "github.com/NordCoder/unmatched-web/internal/application",
    "github.com/NordCoder/unmatched-web/internal/persistence",
    "github.com/NordCoder/unmatched-web/internal/transport",
)

FORBIDDEN_BOOTSTRAP_IDENTIFIERS = (
    "robin-hood",
    "bigfoot",
    "sherwood-forest",
)

IMPORT_RE = re.compile(r'^\s*(?:[._A-Za-z][._A-Za-z0-9]*\s+)?"([^"]+)"\s*$', re.MULTILINE)
IMPORT_BLOCK_RE = re.compile(r"import\s*\((.*?)\)", re.DOTALL)
SINGLE_IMPORT_RE = re.compile(r'import\s+(?:[._A-Za-z][._A-Za-z0-9]*\s+)?"([^"]+)"')


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)


def go_imports(source: str) -> set[str]:
    imports = set(SINGLE_IMPORT_RE.findall(source))
    for block in IMPORT_BLOCK_RE.findall(source):
        imports.update(IMPORT_RE.findall(block))
    return imports


def validate_required_paths(errors: list[str]) -> None:
    for relative in REQUIRED_PATHS:
        if not (ROOT / relative).is_file():
            errors.append(f"required bootstrap file is missing: {relative}")


def validate_domain_boundary(errors: list[str]) -> None:
    domain_root = ROOT / "internal/domain"
    for path in domain_root.rglob("*.go"):
        source = path.read_text(encoding="utf-8")
        for imported in sorted(go_imports(source)):
            if imported in FORBIDDEN_DOMAIN_IMPORTS:
                errors.append(f"{path.relative_to(ROOT)} imports forbidden package {imported!r}")
            if imported.startswith("github.com/NordCoder/unmatched-web/internal/"):
                errors.append(f"{path.relative_to(ROOT)} imports another internal layer {imported!r}")


def validate_contract_boundary(errors: list[str]) -> None:
    contracts = ROOT / "packages/contracts/src"
    for path in contracts.rglob("*.ts"):
        text = path.read_text(encoding="utf-8").lower()
        hand_written_contract_tokens = ("interface command", "interface event", "interface projection")
        for token in hand_written_contract_tokens:
            if token in text:
                errors.append(f"{path.relative_to(ROOT)} contains hand-written transport contract token {token!r}")


def validate_no_gameplay_ids(errors: list[str]) -> None:
    roots = (ROOT / "apps", ROOT / "internal", ROOT / "packages")
    for source_root in roots:
        for path in source_root.rglob("*"):
            if not path.is_file() or path.suffix not in {".go", ".ts", ".tsx", ".js", ".json"}:
                continue
            text = path.read_text(encoding="utf-8").lower()
            for identifier in FORBIDDEN_BOOTSTRAP_IDENTIFIERS:
                if identifier in text:
                    errors.append(f"{path.relative_to(ROOT)} contains gameplay identifier {identifier!r}")


def main() -> int:
    errors: list[str] = []
    validate_required_paths(errors)
    validate_domain_boundary(errors)
    validate_contract_boundary(errors)
    validate_no_gameplay_ids(errors)

    if errors:
        for error in errors:
            fail(error)
        return 1

    print("engine bootstrap validation: PASS")
    print(f"required files: {len(REQUIRED_PATHS)}")
    print("domain forbidden imports: 0")
    print("hand-written public contracts: 0")
    print("gameplay identifiers in source: 0")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
