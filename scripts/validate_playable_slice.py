#!/usr/bin/env python3
"""Validate the launch content and Sherwood Forest machine contract."""
from __future__ import annotations

import json
import re
import sys
from collections import deque
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
RUNTIME_MAP = ROOT / "internal/playableslice/content/data/sherwood-forest.json"
FIXTURE_MAP = ROOT / "tests/fixtures/battlefields/sherwood-forest.json"
YAML_MAP = ROOT / "docs/battlefields/sherwood-forest.yaml"

EXPECTED_STARTS = {"s20": 1, "s19": 2}
EXPECTED_EDGES = {
    tuple(sorted(edge))
    for edge in [
        ("s01", "s02"), ("s02", "s03"), ("s03", "s04"), ("s04", "s05"),
        ("s05", "s06"), ("s06", "s07"), ("s07", "s08"), ("s08", "s09"),
        ("s09", "s10"), ("s10", "s11"), ("s11", "s12"), ("s12", "s13"),
        ("s13", "s14"), ("s14", "s15"), ("s15", "s16"), ("s16", "s17"),
        ("s17", "s18"), ("s18", "s19"), ("s19", "s01"), ("s03", "s20"),
        ("s20", "s21"), ("s21", "s22"), ("s21", "s27"), ("s22", "s28"),
        ("s27", "s28"), ("s19", "s24"), ("s24", "s25"), ("s25", "s26"),
        ("s26", "s27"), ("s26", "s29"), ("s29", "s14"), ("s29", "s13"),
        ("s28", "s30"), ("s30", "s09"), ("s30", "s10"), ("s30", "s11"),
        ("s23", "s05"), ("s23", "s06"), ("s23", "s07"),
    ]
}

EXPECTED_CARDS = {
    "robin-hood": {
        "a-hunters-eye", "steal-from-the-rich", "disarming-shot", "piercing-shot",
        "highway-robbery", "defenders-of-sherwood", "feint", "regroup",
        "wily-fighting", "snark", "ambush",
    },
    "bigfoot": {
        "larger-than-life", "savagery", "crash-through-the-trees", "jackalope-horns",
        "hoax", "disengage", "its-just-your-imagination", "feint", "regroup",
        "skirmish", "momentous-shift",
    },
}

class ValidationError(ValueError):
    pass

def validate_battlefield(data: dict) -> None:
    if data.get("id") != "sherwood-forest":
        raise ValidationError("unexpected battlefield id")
    spaces = data.get("spaces", [])
    zones = data.get("zones", [])
    edges = data.get("edges", [])
    if len(spaces) != 30:
        raise ValidationError(f"space count is {len(spaces)}, want 30")
    if len(zones) != 7 or len(set(zones)) != 7:
        raise ValidationError("zone registry must contain seven unique zones")
    space_ids = [space.get("id") for space in spaces]
    if None in space_ids or len(set(space_ids)) != len(space_ids):
        raise ValidationError("space ids must be non-empty and unique")
    known = set(space_ids)
    starts = {space["id"]: space.get("starting_seat") for space in spaces if space.get("starting_seat")}
    if starts != EXPECTED_STARTS:
        raise ValidationError(f"starting markers are {starts}, want {EXPECTED_STARTS}")
    for space in spaces:
        if not space.get("zones") or not set(space["zones"]).issubset(set(zones)):
            raise ValidationError(f"space {space['id']} has invalid zones")
        if not (0 <= space.get("x", -1) <= 100 and 0 <= space.get("y", -1) <= 100):
            raise ValidationError(f"space {space['id']} has invalid coordinates")
    graph = {space_id: set() for space_id in known}
    seen = set()
    for edge in edges:
        a, b = edge.get("from"), edge.get("to")
        if a not in known or b not in known:
            raise ValidationError(f"edge references unknown space: {a}-{b}")
        if a == b:
            raise ValidationError(f"self edge at {a}")
        key = tuple(sorted((a, b)))
        if key in seen:
            raise ValidationError(f"duplicate undirected edge: {a}-{b}")
        seen.add(key)
        graph[a].add(b)
        graph[b].add(a)
    if seen != EXPECTED_EDGES:
        missing = sorted(EXPECTED_EDGES - seen)
        extra = sorted(seen - EXPECTED_EDGES)
        raise ValidationError(f"canonical edge set differs; missing={missing}, extra={extra}")
    reached = set()
    queue = deque([space_ids[0]])
    while queue:
        current = queue.popleft()
        if current in reached:
            continue
        reached.add(current)
        queue.extend(graph[current] - reached)
    if reached != known:
        raise ValidationError(f"battlefield disconnected: reached {len(reached)}/{len(known)}")

def top_level_card_ids(text: str) -> set[str]:
    ids = set()
    in_cards = False
    for line in text.splitlines():
        if line == "cards:":
            in_cards = True
            continue
        if in_cards and line and not line.startswith(" "):
            break
        if not in_cards or not line.startswith("  - "):
            continue
        match = re.match(r"\s*-\s*(?:\{\s*)?id:\s*([^,}\s]+)", line)
        if match:
            ids.add(match.group(1).strip('"'))
    return ids

def validate_deck(fighter_id: str) -> None:
    path = ROOT / "docs/cards/phase-4b" / f"{fighter_id}.yaml"
    text = path.read_text(encoding="utf-8")
    ids = top_level_card_ids(text)
    if ids != EXPECTED_CARDS[fighter_id]:
        raise ValidationError(f"{fighter_id} card ids differ: {sorted(ids ^ EXPECTED_CARDS[fighter_id])}")
    quantities = [int(value) for value in re.findall(r"\bquantity:\s*(\d+)", text)]
    if sum(quantities) != 30:
        raise ValidationError(f"{fighter_id} quantity sum is {sum(quantities)}, want 30")
    if "status: verified" not in text or "integration: ready" not in text:
        raise ValidationError(f"{fighter_id} manifest is not verified and integration-ready")

def validate() -> None:
    runtime_raw = RUNTIME_MAP.read_bytes()
    fixture_raw = FIXTURE_MAP.read_bytes()
    if runtime_raw != fixture_raw:
        raise ValidationError("runtime battlefield and test fixture differ")
    data = json.loads(runtime_raw)
    validate_battlefield(data)
    yaml_text = YAML_MAP.read_text(encoding="utf-8")
    for space in data["spaces"]:
        if f"id: {space['id']}" not in yaml_text:
            raise ValidationError(f"YAML mirror is missing {space['id']}")
    for edge in data["edges"]:
        marker = f"from: {edge['from']}, to: {edge['to']}"
        if marker not in yaml_text:
            raise ValidationError(f"YAML mirror is missing edge {marker}")
    for fighter_id in EXPECTED_CARDS:
        validate_deck(fighter_id)

if __name__ == "__main__":
    try:
        validate()
    except (OSError, json.JSONDecodeError, ValidationError) as error:
        print(f"playable-slice validation: FAIL: {error}", file=sys.stderr)
        raise SystemExit(1)
    print("playable-slice validation: PASS (2 decks, 60 cards, 30 spaces, 7 zones, 39 edges)")
