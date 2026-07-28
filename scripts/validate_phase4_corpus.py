#!/usr/bin/env python3
"""Validate the integrated Phase 4 competitive fighter/card corpus.

This validator intentionally focuses on frozen corpus invariants that can be
checked without implementing the game engine: YAML integrity, canonical roster,
fighter/card pairing, action-deck construction, usable_by references, owner
requirement references, local IDs, and preservation of known exceptions.
"""

from __future__ import annotations

import json
import re
import sys
from collections import Counter
from pathlib import Path
from typing import Any, Iterable

import yaml

ROOT = Path(__file__).resolve().parents[1]
FIGHTER_DIRS = [ROOT / "docs/fighters/phase-4a", ROOT / "docs/fighters/phase-4b"]
CARD_DIRS = [ROOT / "docs/cards/phase-4a", ROOT / "docs/cards/phase-4b"]
REPORT_DIR = ROOT / "docs/phase-4b"
OUTPUT = ROOT / "artifacts/phase-4-validation.json"

EXPECTED_CANONICAL_IDS = {
    # Phase 4A
    "achilles", "bloody-mary", "sun-wukong", "sherlock-holmes", "dracula",
    "raptors", "wayward-sisters", "geralt-of-rivia", "yennefer-triss",
    "black-panther",
    # Worker A
    "alice", "king-arthur", "medusa", "sinbad", "robin-hood", "bigfoot",
    "robert-muldoon", "invisible-man", "jekyll-and-hyde", "buffy", "willow",
    "spike", "angel", "little-red-riding-hood", "beowulf", "deadpool",
    "yennenga",
    # Worker B
    "daredevil", "elektra", "bullseye", "ghost-rider", "luke-cage",
    "moon-knight", "dr-ellie-sattler", "t-rex", "houdini", "genie",
    "cloak-and-dagger", "ms-marvel", "squirrel-girl", "black-widow",
    "winter-soldier", "doctor-strange", "she-hulk", "spider-man",
    # Worker C
    "annie-christmas", "dr-jill-trent", "golden-bat", "nikola-tesla",
    "oda-nobunaga", "tomoe-gozen", "shakespeare", "hamlet", "titania",
    "ciri", "ancient-leshen", "eredin", "philippa",
    # Worker D
    "bruce-lee", "muhammad-ali", "blackbeard", "chupacabra", "loki",
    "pandora", "leonardo", "donatello", "michelangelo", "raphael",
    "rosie-the-riveter", "john-henry", "wyatt-earp", "george-washington",
    "shredder", "krang",
}

FORBIDDEN_ALIAS_IDS = {
    "yennefer-and-triss", "jekyll-hyde", "little-red", "dr-sattler"
}

EXPECTED_DECK_COUNTS = {
    "daredevil": (22, 22, "fixed"),
    "elektra": (20, 20, "fixed"),
    "black-widow": (31, 31, "fixed"),
    "geralt-of-rivia": (36, 30, "choose_groups"),
    "buffy": (35, 30, "choose_groups"),
}

ALLOWED_CARD_TYPES = {"attack", "defense", "versatile", "scheme"}
REQ_PATTERN = re.compile(r"^(?:[ \t]*- id:[ \t]*`?|###[ \t]+)([ABCD]-REQ-\d+)(?:`|\b)", re.MULTILINE)
REF_PATTERN = re.compile(r"^(?:choice|result|operation|bound)\.([A-Za-z0-9_-]+)$")


class UniqueKeyLoader(yaml.SafeLoader):
    pass


def construct_mapping(loader: UniqueKeyLoader, node: yaml.nodes.MappingNode, deep: bool = False) -> dict[Any, Any]:
    mapping: dict[Any, Any] = {}
    for key_node, value_node in node.value:
        key = loader.construct_object(key_node, deep=deep)
        if key in mapping:
            raise yaml.constructor.ConstructorError(
                "while constructing a mapping", node.start_mark,
                f"found duplicate key {key!r}", key_node.start_mark,
            )
        mapping[key] = loader.construct_object(value_node, deep=deep)
    return mapping


UniqueKeyLoader.add_constructor(
    yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG, construct_mapping
)


def yaml_files(directories: Iterable[Path]) -> list[Path]:
    return sorted(path for directory in directories for path in directory.glob("*.yaml"))


def load_yaml(path: Path, errors: list[str]) -> dict[str, Any] | None:
    try:
        value = yaml.load(path.read_text(encoding="utf-8"), Loader=UniqueKeyLoader)
    except Exception as exc:  # noqa: BLE001 - validation must report parser failures
        errors.append(f"{path.relative_to(ROOT)}: YAML parse failed: {exc}")
        return None
    if not isinstance(value, dict):
        errors.append(f"{path.relative_to(ROOT)}: top-level YAML value must be a mapping")
        return None
    return value


def walk(value: Any) -> Iterable[tuple[str | None, Any]]:
    if isinstance(value, dict):
        for key, child in value.items():
            yield str(key), child
            yield from walk(child)
    elif isinstance(value, list):
        for child in value:
            yield None, child
            yield from walk(child)


def as_string_list(value: Any) -> list[str]:
    if isinstance(value, str):
        return [value]
    if isinstance(value, list) and all(isinstance(item, str) for item in value):
        return list(value)
    return []


def collect_requirement_definitions(errors: list[str]) -> set[str]:
    found: list[str] = []
    for owner in "abcd":
        path = REPORT_DIR / f"worker-{owner}-report.md"
        if not path.exists():
            errors.append(f"missing owner requirement report: {path.relative_to(ROOT)}")
            continue
        found.extend(REQ_PATTERN.findall(path.read_text(encoding="utf-8")))
    duplicates = sorted(req for req, count in Counter(found).items() if count > 1)
    if duplicates:
        errors.append(f"duplicate owner requirement definitions: {duplicates}")
    return set(found)


def main() -> int:
    errors: list[str] = []
    warnings: list[str] = []

    fighter_paths = yaml_files(FIGHTER_DIRS)
    card_paths = yaml_files(CARD_DIRS)
    fighters: dict[str, tuple[Path, dict[str, Any]]] = {}
    cards: dict[str, tuple[Path, dict[str, Any]]] = {}

    for path in fighter_paths:
        data = load_yaml(path, errors)
        if data is None:
            continue
        fighter_id = data.get("id")
        if not isinstance(fighter_id, str) or not fighter_id:
            errors.append(f"{path.relative_to(ROOT)}: missing string id")
            continue
        if fighter_id in fighters:
            errors.append(f"duplicate fighter id {fighter_id}: {fighters[fighter_id][0]} and {path}")
        fighters[fighter_id] = (path, data)
        if path.stem != fighter_id:
            errors.append(f"{path.relative_to(ROOT)}: file stem does not match id {fighter_id}")
        if data.get("schema_version") != 1:
            errors.append(f"{path.relative_to(ROOT)}: schema_version must be 1")

    for path in card_paths:
        data = load_yaml(path, errors)
        if data is None:
            continue
        fighter_id = data.get("fighter_id")
        if not isinstance(fighter_id, str) or not fighter_id:
            errors.append(f"{path.relative_to(ROOT)}: missing string fighter_id")
            continue
        if fighter_id in cards:
            errors.append(f"duplicate card manifest for fighter {fighter_id}: {cards[fighter_id][0]} and {path}")
        cards[fighter_id] = (path, data)
        if path.stem != fighter_id:
            errors.append(f"{path.relative_to(ROOT)}: file stem does not match fighter_id {fighter_id}")
        if data.get("schema_version") != 1:
            errors.append(f"{path.relative_to(ROOT)}: schema_version must be 1")

    actual_ids = set(fighters)
    if len(EXPECTED_CANONICAL_IDS) != 74:
        errors.append(f"validator expected-ID table is not 74 entries: {len(EXPECTED_CANONICAL_IDS)}")
    if actual_ids != EXPECTED_CANONICAL_IDS:
        errors.append(
            "canonical fighter roster mismatch: "
            f"missing={sorted(EXPECTED_CANONICAL_IDS - actual_ids)}, "
            f"unexpected={sorted(actual_ids - EXPECTED_CANONICAL_IDS)}"
        )
    if set(cards) != actual_ids:
        errors.append(
            "fighter/card manifest pairing mismatch: "
            f"fighters_without_cards={sorted(actual_ids - set(cards))}, "
            f"cards_without_fighters={sorted(set(cards) - actual_ids)}"
        )
    aliases_present = sorted(FORBIDDEN_ALIAS_IDS & (actual_ids | set(cards)))
    if aliases_present:
        errors.append(f"evidence aliases present as canonical identities: {aliases_present}")
    if sum(1 for fighter_id in actual_ids if fighter_id == "bruce-lee") != 1:
        errors.append("Bruce Lee must resolve to exactly one canonical identity")

    requirement_definitions = collect_requirement_definitions(errors)
    requirement_refs: set[str] = set()
    status_counts: Counter[str] = Counter()
    deck_summary: dict[str, dict[str, Any]] = {}
    total_unique_action_cards = 0
    total_action_card_copies = 0

    for fighter_id in sorted(actual_ids & set(cards)):
        fighter_path, fighter = fighters[fighter_id]
        card_path, card_manifest = cards[fighter_id]
        status = fighter.get("status")
        status_counts[str(status)] += 1

        deck = fighter.get("deck")
        if not isinstance(deck, dict):
            errors.append(f"{fighter_path.relative_to(ROOT)}: deck must be a mapping")
        else:
            manifest_ref = deck.get("manifest")
            if not isinstance(manifest_ref, str):
                errors.append(f"{fighter_path.relative_to(ROOT)}: deck.manifest must be a string")
            else:
                resolved = (fighter_path.parent / manifest_ref).resolve()
                if resolved != card_path.resolve():
                    errors.append(
                        f"{fighter_path.relative_to(ROOT)}: deck.manifest resolves to "
                        f"{resolved.relative_to(ROOT) if resolved.is_relative_to(ROOT) else resolved}, "
                        f"expected {card_path.relative_to(ROOT)}"
                    )

        topology = fighter.get("topology", {})
        topology_fighters = topology.get("fighters", []) if isinstance(topology, dict) else []
        local_fighter_ids: list[str] = []
        if not isinstance(topology_fighters, list) or not topology_fighters:
            errors.append(f"{fighter_path.relative_to(ROOT)}: topology.fighters must be a non-empty list")
        else:
            for entry in topology_fighters:
                if not isinstance(entry, dict) or not isinstance(entry.get("id"), str):
                    errors.append(f"{fighter_path.relative_to(ROOT)}: invalid topology fighter entry {entry!r}")
                    continue
                local_fighter_ids.append(entry["id"])
                count = entry.get("count", 1)
                if not isinstance(count, int) or count < 1:
                    errors.append(f"{fighter_path.relative_to(ROOT)}: invalid count for topology fighter {entry['id']}")
            if len(local_fighter_ids) != len(set(local_fighter_ids)):
                errors.append(f"{fighter_path.relative_to(ROOT)}: duplicate topology fighter ids")

        construction = card_manifest.get("construction")
        if not isinstance(construction, dict):
            errors.append(f"{card_path.relative_to(ROOT)}: construction must be a mapping")
            continue
        kind = construction.get("kind")
        available = construction.get("available_pool_count")
        game_deck = construction.get("game_deck_count")
        if kind not in {"fixed", "choose_groups", "constructed"}:
            errors.append(f"{card_path.relative_to(ROOT)}: unsupported construction kind {kind!r}")
        if not isinstance(available, int) or available < 1:
            errors.append(f"{card_path.relative_to(ROOT)}: invalid available_pool_count {available!r}")
        if not isinstance(game_deck, int) or game_deck < 1:
            errors.append(f"{card_path.relative_to(ROOT)}: invalid game_deck_count {game_deck!r}")

        definitions = card_manifest.get("cards")
        if not isinstance(definitions, list) or not definitions:
            errors.append(f"{card_path.relative_to(ROOT)}: cards must be a non-empty list")
            continue
        card_ids: list[str] = []
        quantity_sum = 0
        for definition in definitions:
            if not isinstance(definition, dict):
                errors.append(f"{card_path.relative_to(ROOT)}: card entry must be a mapping")
                continue
            card_id = definition.get("id")
            if not isinstance(card_id, str) or not card_id:
                errors.append(f"{card_path.relative_to(ROOT)}: card entry missing string id")
                continue
            card_ids.append(card_id)
            quantity = definition.get("quantity")
            if not isinstance(quantity, int) or quantity < 1:
                errors.append(f"{card_path.relative_to(ROOT)}:{card_id}: quantity must be a positive integer")
            else:
                quantity_sum += quantity
            card_type = definition.get("type")
            if card_type not in ALLOWED_CARD_TYPES:
                errors.append(f"{card_path.relative_to(ROOT)}:{card_id}: invalid type {card_type!r}")
            usable_by = definition.get("usable_by")
            if usable_by != "any":
                allowed_users = as_string_list(usable_by)
                if not allowed_users:
                    errors.append(f"{card_path.relative_to(ROOT)}:{card_id}: invalid usable_by {usable_by!r}")
                else:
                    unknown = sorted(set(allowed_users) - set(local_fighter_ids))
                    if unknown:
                        errors.append(
                            f"{card_path.relative_to(ROOT)}:{card_id}: usable_by references unknown local fighters {unknown}"
                        )
        if len(card_ids) != len(set(card_ids)):
            duplicates = sorted(card_id for card_id, count in Counter(card_ids).items() if count > 1)
            errors.append(f"{card_path.relative_to(ROOT)}: duplicate card ids {duplicates}")
        if isinstance(available, int) and quantity_sum != available:
            errors.append(
                f"{card_path.relative_to(ROOT)}: quantity sum {quantity_sum} != available_pool_count {available}"
            )
        if kind == "fixed" and isinstance(available, int) and isinstance(game_deck, int) and available != game_deck:
            errors.append(f"{card_path.relative_to(ROOT)}: fixed construction must have pool == game deck")
        if kind in {"choose_groups", "constructed"} and isinstance(available, int) and isinstance(game_deck, int) and game_deck > available:
            errors.append(f"{card_path.relative_to(ROOT)}: constructed game deck exceeds available pool")

        expected_counts = EXPECTED_DECK_COUNTS.get(fighter_id)
        if expected_counts is not None and (available, game_deck, kind) != expected_counts:
            errors.append(
                f"{card_path.relative_to(ROOT)}: expected special construction {expected_counts}, "
                f"observed {(available, game_deck, kind)}"
            )

        for document in (fighter, card_manifest):
            for key, value in walk(document):
                if key == "requires":
                    refs = as_string_list(value)
                    if not refs and value not in ([], None):
                        errors.append(f"{fighter_id}: malformed requires value {value!r}")
                    requirement_refs.update(refs)

        # Binding/reference audit. Unresolved references are warnings because some
        # source-defined engine references are intentionally built-in rather than
        # bound by the local manifest.
        binds: set[str] = set()
        refs: set[str] = set()
        for document in (fighter, card_manifest):
            for key, value in walk(document):
                if key in {"bind", "bind_result"} and isinstance(value, str):
                    binds.add(value)
                if key == "ref" and isinstance(value, str):
                    match = REF_PATTERN.match(value)
                    if match:
                        refs.add(match.group(1))
        unresolved_local_refs = sorted(refs - binds)
        if unresolved_local_refs:
            warnings.append(f"{fighter_id}: locally unresolved choice/result refs {unresolved_local_refs}")

        total_unique_action_cards += len(card_ids)
        total_action_card_copies += quantity_sum
        deck_summary[fighter_id] = {
            "kind": kind,
            "available_pool_count": available,
            "game_deck_count": game_deck,
            "quantity_sum": quantity_sum,
            "unique_card_definitions": len(card_ids),
        }

    unresolved_requirements = sorted(requirement_refs - requirement_definitions)
    if unresolved_requirements:
        errors.append(f"unresolved owner requirement references: {unresolved_requirements}")

    blocked = sorted(fighter_id for fighter_id, (_, data) in fighters.items() if data.get("status") == "blocked")
    if blocked != ["deadpool"]:
        errors.append(f"expected only Deadpool to remain blocked; observed {blocked}")

    # Status/verification consistency for Phase 4B manifests.
    for fighter_id, (path, data) in fighters.items():
        if "phase-4b" not in path.parts:
            continue
        verification = data.get("verification")
        if not isinstance(verification, dict):
            errors.append(f"{path.relative_to(ROOT)}: Phase 4B verification must be a mapping")
            continue
        if data.get("status") == "verified" and verification.get("integration") != "ready":
            errors.append(f"{path.relative_to(ROOT)}: verified fighter must have verification.integration=ready")
        if data.get("status") == "blocked" and verification.get("policy") != "blocked":
            errors.append(f"{path.relative_to(ROOT)}: blocked fighter must retain blocked policy state")

    result = {
        "schema_version": 1,
        "phase": "phase-4-final-integration",
        "passed": not errors,
        "counts": {
            "canonical_fighter_identities": len(actual_ids),
            "fighter_manifests": len(fighter_paths),
            "card_manifests": len(card_paths),
            "unique_action_card_definitions": total_unique_action_cards,
            "available_action_card_copies": total_action_card_copies,
            "owner_requirement_definitions": len(requirement_definitions),
            "owner_requirement_references": len(requirement_refs),
        },
        "status_counts": dict(sorted(status_counts.items())),
        "blocked_fighters": blocked,
        "special_deck_constructions": {key: deck_summary.get(key) for key in sorted(EXPECTED_DECK_COUNTS)},
        "requirements": {
            "defined": sorted(requirement_definitions),
            "referenced": sorted(requirement_refs),
            "unresolved": unresolved_requirements,
        },
        "errors": errors,
        "warnings": warnings,
    }
    OUTPUT.parent.mkdir(parents=True, exist_ok=True)
    OUTPUT.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0 if not errors else 1


if __name__ == "__main__":
    sys.exit(main())
