from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

from scripts import validate_engine_bootstrap as validator


class EngineBootstrapValidatorTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary_directory = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary_directory.name)

    def tearDown(self) -> None:
        self.temporary_directory.cleanup()

    def write(self, relative: str, content: str) -> None:
        path = self.root / relative
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")

    def test_domain_rejects_adapter_stdlib_and_third_party_imports(self) -> None:
        self.write(
            "internal/domain/boundary.go",
            '''package domain

import (
    "database/sql"
    "github.com/example/websocket"
    "net/http"
    "os"
)
''',
        )

        errors: list[str] = []
        validator.validate_domain_boundary(self.root, errors)

        self.assertTrue(any("database/sql" in error for error in errors))
        self.assertTrue(any("github.com/example/websocket" in error for error in errors))
        self.assertTrue(any("net/http" in error for error in errors))
        self.assertTrue(any("'os'" in error for error in errors))

    def test_contract_boundary_rejects_named_interfaces_and_type_aliases(self) -> None:
        self.write(
            "packages/contracts/src/index.ts",
            '''export const contractGenerationState = "not-configured" as const;

export interface SubmitActionCommand {
  readonly actionId: string;
}

export type PlayerViewProjection = {
  readonly revision: number;
};

export const MatchEvent = { kind: "created" } as const;
''',
        )

        errors: list[str] = []
        validator.validate_contract_boundary(self.root, errors)

        self.assertTrue(any("interface SubmitActionCommand" in error for error in errors))
        self.assertTrue(any("type PlayerViewProjection" in error for error in errors))
        self.assertTrue(any("public export const MatchEvent" in error for error in errors))

    def test_gameplay_ids_are_data_derived_and_scripts_are_scanned(self) -> None:
        self.write(
            "docs/fighters/phase-4b/qa-sample.yaml",
            "id: qa-sample-fighter\n",
        )
        self.write(
            "docs/cards/phase-4b/qa-sample.yaml",
            "fighter_id: qa-sample-fighter\ncards:\n  - id: qa-sample-card\n",
        )
        self.write(
            "scripts/validate_engine_bootstrap.py",
            'FIGHTER = "qa-sample-fighter"\nCARD = "qa-sample-card"\n',
        )

        errors: list[str] = []
        validator.validate_no_gameplay_ids(self.root, errors)

        self.assertTrue(any("qa-sample-fighter" in error for error in errors))
        self.assertTrue(any("qa-sample-card" in error for error in errors))
        self.assertTrue(all("scripts/validate_engine_bootstrap.py" in error for error in errors))

    def test_nested_effect_ids_do_not_expand_the_identity_boundary(self) -> None:
        self.write(
            "docs/cards/phase-4b/qa-sample.yaml",
            "fighter_id: qa-sample-fighter\ncards:\n"
            "  - id: qa-sample-card\n"
            "    effects:\n"
            "      - id: command\n",
        )
        self.write("internal/domain/doc.go", "// command orchestration remains deferred\n")

        errors: list[str] = []
        validator.validate_no_gameplay_ids(self.root, errors)

        self.assertEqual([], errors)


if __name__ == "__main__":
    unittest.main()
