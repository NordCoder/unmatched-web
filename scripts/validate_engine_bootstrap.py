#!/usr/bin/env python3
"""Validate the bounded engine-bootstrap architecture contract."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MODULE_PATH = "github.com/NordCoder/unmatched-web"

REQUIRED_PATHS = (
    "go.mod",
    "package.json",
    "package-lock.json",
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
    "docs/engine/bootstrap-development.md",
    "scripts/verify_postgres_persistence.sh",
    "tests/architecture/test_validate_engine_bootstrap.py",
    "tests/fixtures/README.md",
)

FORBIDDEN_DOMAIN_STDLIB_PREFIXES = (
    "database",
    "embed",
    "io/fs",
    "net",
    "os",
    "path/filepath",
    "syscall",
)

GAMEPLAY_CORPUS_ROOTS = (
    "docs/fighters",
    "docs/cards",
    "docs/battlefields",
)

IMPLEMENTATION_SUFFIXES = {
    ".go",
    ".js",
    ".json",
    ".md",
    ".mod",
    ".py",
    ".sh",
    ".ts",
    ".tsx",
    ".yaml",
    ".yml",
}
IMPLEMENTATION_ROOTS = (
    "apps",
    "internal",
    "packages",
    "tests/architecture",
)
IMPLEMENTATION_EXPLICIT_FILES = {
    ".env.example",
    ".gitignore",
    ".github/workflows/engine-bootstrap.yml",
    "compose.yaml",
    "docs/engine/bootstrap-development.md",
    "go.mod",
    "package-lock.json",
    "package.json",
    "scripts/validate_engine_bootstrap.py",
    "scripts/verify_postgres_persistence.sh",
}

IMPORT_RE = re.compile(r'^\s*(?:[._A-Za-z][._A-Za-z0-9]*\s+)?"([^"]+)"\s*$', re.MULTILINE)
IMPORT_BLOCK_RE = re.compile(r"\bimport\s*\((.*?)\)", re.DOTALL)
SINGLE_IMPORT_RE = re.compile(r'\bimport\s+(?:[._A-Za-z][._A-Za-z0-9]*\s+)?"([^"]+)"')
SCALAR_GAMEPLAY_ID_RE = re.compile(
    r"^[\"']?([a-z0-9][a-z0-9-]*)[\"']?\s*$"
)
TOP_LEVEL_GAMEPLAY_ID_RE = re.compile(
    r"^(?:id|fighter_id|card_id|battlefield_id):\s*(.+?)\s*$",
    re.MULTILINE,
)
CARD_DEFINITION_ID_RE = re.compile(r"^  - id:\s*(.+?)\s*$", re.MULTILINE)
FIGHTER_LIST_START_RE = re.compile(r"^  fighters:\s*$")
FIGHTER_LIST_ID_RE = re.compile(r"^    -\s*(?:\{\s*)?id:\s*([^,}\s]+)")

TYPE_DECLARATION_RE = re.compile(
    r"^\s*(?:export\s+)?(?:declare\s+|abstract\s+)?"
    r"(interface|type|class|enum|namespace)\s+([A-Za-z_$][\w$]*)",
    re.MULTILINE,
)
EXPORT_DECLARATION_RE = re.compile(
    r"^\s*export\s+(?:default\s+)?(?:declare\s+|abstract\s+)?"
    r"(interface|type|class|enum|namespace|function|const|let|var)\s+"
    r"([A-Za-z_$][\w$]*)",
    re.MULTILINE,
)
EXPORT_LIST_RE = re.compile(r"^\s*export\s*{", re.MULTILINE)


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)


def go_imports(source: str) -> set[str]:
    imports = set(SINGLE_IMPORT_RE.findall(source))
    for block in IMPORT_BLOCK_RE.findall(source):
        imports.update(IMPORT_RE.findall(block))
    return imports


def validate_required_paths(root: Path, errors: list[str]) -> None:
    for relative in REQUIRED_PATHS:
        if not (root / relative).is_file():
            errors.append(f"required bootstrap file is missing: {relative}")


def is_standard_library(imported: str) -> bool:
    return "." not in imported.split("/", 1)[0]


def has_forbidden_stdlib_prefix(imported: str) -> bool:
    return any(
        imported == prefix or imported.startswith(f"{prefix}/")
        for prefix in FORBIDDEN_DOMAIN_STDLIB_PREFIXES
    )


def validate_domain_boundary(root: Path, errors: list[str]) -> None:
    domain_root = root / "internal/domain"
    if not domain_root.exists():
        return

    for path in domain_root.rglob("*.go"):
        source = path.read_text(encoding="utf-8")
        for imported in sorted(go_imports(source)):
            relative = path.relative_to(root)
            domain_path = f"{MODULE_PATH}/internal/domain"
            if imported == domain_path or imported.startswith(f"{domain_path}/"):
                continue
            if imported.startswith(f"{MODULE_PATH}/internal/"):
                errors.append(f"{relative} imports another internal layer {imported!r}")
                continue
            if not is_standard_library(imported):
                errors.append(f"{relative} imports third-party package {imported!r}")
                continue
            if has_forbidden_stdlib_prefix(imported):
                errors.append(f"{relative} imports forbidden standard-library package {imported!r}")


def is_generated_contract(path: Path, source: str) -> bool:
    relative_parts = path.parts
    header = "\n".join(source.splitlines()[:10]).lower()
    return (
        "generated" in relative_parts
        and "@generated" in header
        and "do not edit" in header
    )


def validate_contract_boundary(root: Path, errors: list[str]) -> None:
    contracts = root / "packages/contracts/src"
    if not contracts.exists():
        return

    for path in contracts.rglob("*.ts"):
        source = path.read_text(encoding="utf-8")
        if is_generated_contract(path.relative_to(contracts), source):
            continue

        relative = path.relative_to(root)
        for kind, name in TYPE_DECLARATION_RE.findall(source):
            errors.append(
                f"{relative} contains hand-written contract declaration {kind} {name}"
            )

        for kind, name in EXPORT_DECLARATION_RE.findall(source):
            if path == contracts / "index.ts" and (kind, name) == (
                "const",
                "contractGenerationState",
            ):
                continue
            if kind not in {"interface", "type", "class", "enum", "namespace"}:
                errors.append(
                    f"{relative} contains hand-written public export {kind} {name}"
                )

        if EXPORT_LIST_RE.search(source):
            errors.append(f"{relative} contains a hand-written public export list")


def scalar_gameplay_id(raw_value: str) -> str | None:
    match = SCALAR_GAMEPLAY_ID_RE.match(raw_value.strip())
    return match.group(1) if match else None


def fighter_manifest_identifiers(source: str) -> set[str]:
    identifiers = {
        identifier
        for raw_value in TOP_LEVEL_GAMEPLAY_ID_RE.findall(source)
        if (identifier := scalar_gameplay_id(raw_value)) is not None
    }

    in_fighter_list = False
    for line in source.splitlines():
        if FIGHTER_LIST_START_RE.match(line):
            in_fighter_list = True
            continue
        if in_fighter_list and line.strip() and len(line) - len(line.lstrip()) <= 2:
            in_fighter_list = False
        if not in_fighter_list:
            continue
        match = FIGHTER_LIST_ID_RE.match(line)
        if match and (identifier := scalar_gameplay_id(match.group(1))) is not None:
            identifiers.add(identifier)
    return identifiers


def card_manifest_identifiers(source: str) -> set[str]:
    identifiers = {
        identifier
        for raw_value in TOP_LEVEL_GAMEPLAY_ID_RE.findall(source)
        if (identifier := scalar_gameplay_id(raw_value)) is not None
    }
    identifiers.update(
        identifier
        for raw_value in CARD_DEFINITION_ID_RE.findall(source)
        if (identifier := scalar_gameplay_id(raw_value)) is not None
    )
    return identifiers


def battlefield_manifest_identifiers(source: str) -> set[str]:
    return {
        identifier
        for raw_value in TOP_LEVEL_GAMEPLAY_ID_RE.findall(source)
        if (identifier := scalar_gameplay_id(raw_value)) is not None
    }


def canonical_gameplay_identifiers(root: Path) -> set[str]:
    identifiers: set[str] = set()
    extractors = {
        "docs/fighters": fighter_manifest_identifiers,
        "docs/cards": card_manifest_identifiers,
        "docs/battlefields": battlefield_manifest_identifiers,
    }
    for relative_root, extractor in extractors.items():
        corpus_root = root / relative_root
        if not corpus_root.exists():
            continue
        for path in corpus_root.rglob("*"):
            if not path.is_file() or path.suffix not in {".yaml", ".yml"}:
                continue
            identifiers.update(extractor(path.read_text(encoding="utf-8").lower()))
    return identifiers


def implementation_source_files(root: Path):
    yielded: set[Path] = set()
    for relative_text in sorted(IMPLEMENTATION_EXPLICIT_FILES):
        path = root / relative_text
        if path.is_file():
            yielded.add(path)
            yield path

    for relative_root in IMPLEMENTATION_ROOTS:
        source_root = root / relative_root
        if not source_root.exists():
            continue
        for path in source_root.rglob("*"):
            if not path.is_file() or path.suffix not in IMPLEMENTATION_SUFFIXES:
                continue
            if path in yielded:
                continue
            yielded.add(path)
            yield path


def contains_identifier(text: str, identifier: str) -> bool:
    pattern = rf"(?<![a-z0-9_-]){re.escape(identifier)}(?![a-z0-9_-])"
    return re.search(pattern, text) is not None


def validate_no_gameplay_ids(root: Path, errors: list[str]) -> None:
    identifiers = canonical_gameplay_identifiers(root)
    if not identifiers:
        errors.append("canonical gameplay corpus yielded no identifiers")
        return

    for path in implementation_source_files(root):
        text = path.read_text(encoding="utf-8").lower()
        for identifier in sorted(identifiers):
            if contains_identifier(text, identifier):
                errors.append(
                    f"{path.relative_to(root)} contains canonical gameplay identifier "
                    f"{identifier!r}"
                )


def validate_node_lock(root: Path, errors: list[str]) -> None:
    package = json.loads((root / "package.json").read_text(encoding="utf-8"))
    lock = json.loads((root / "package-lock.json").read_text(encoding="utf-8"))

    expected_version = package["devDependencies"]["typescript"]
    locked_package = lock["packages"].get("node_modules/typescript", {})
    if locked_package.get("version") != expected_version:
        errors.append("package-lock TypeScript version does not match package.json")
    if not locked_package.get("integrity"):
        errors.append("package-lock TypeScript entry has no integrity digest")

    expected_workspaces = sorted(package["workspaces"])
    locked_workspaces = sorted(lock["packages"][""]["workspaces"])
    if locked_workspaces != expected_workspaces:
        errors.append("package-lock workspace graph does not match package.json")


def validate(root: Path) -> list[str]:
    errors: list[str] = []
    validate_required_paths(root, errors)
    if not errors:
        validate_domain_boundary(root, errors)
        validate_contract_boundary(root, errors)
        validate_no_gameplay_ids(root, errors)
        validate_node_lock(root, errors)
    return errors


def main() -> int:
    errors = validate(ROOT)
    if errors:
        for error in errors:
            fail(error)
        return 1

    print("engine bootstrap validation: PASS")
    print(f"required files: {len(REQUIRED_PATHS)}")
    print("domain forbidden imports: 0")
    print("hand-written public contracts: 0")
    print("canonical gameplay identifiers in implementation: 0")
    print("npm lock consistency: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
