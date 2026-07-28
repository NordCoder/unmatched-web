#!/usr/bin/env python3
from __future__ import annotations
import argparse, copy, hashlib, json, re, unicodedata
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[1]
FIXTURE = ROOT / "docs/engine/fixtures/foundation-v2.json"
AUDIT = ROOT / "docs/engine/fixtures/foundation-v2-transition-audit.json"
SCHEMA = ROOT / "docs/engine/fixtures/schema-v2.json"
REQUIRED_KINDS = {"create", "join", "idempotency", "reconnect", "replay"}
OPERATIONAL_KEYS = {"connection", "connection_state", "connection_status", "session_id", "client_instance_id", "last_seen"}

class ValidationError(ValueError):
    pass

def _pairs(pairs):
    out = {}
    for key, value in pairs:
        if key in out:
            raise ValidationError(f"duplicate object key: {key}")
        out[key] = value
    return out

def strict_loads(text: str) -> Any:
    return json.loads(text, object_pairs_hook=_pairs)

def strict_load(path: Path) -> Any:
    return strict_loads(path.read_text(encoding="utf-8"))

def normalize(value: Any) -> Any:
    if isinstance(value, str):
        return unicodedata.normalize("NFC", value)
    if isinstance(value, list):
        return [normalize(item) for item in value]
    if isinstance(value, dict):
        out = {}
        for key, item in value.items():
            normalized_key = unicodedata.normalize("NFC", key)
            if normalized_key in out:
                raise ValidationError(f"NFC-created duplicate key: {normalized_key}")
            out[normalized_key] = normalize(item)
        return out
    if isinstance(value, (int, bool)) or value is None:
        return value
    raise ValidationError(f"unsupported canonical value: {type(value).__name__}")

def canonical_bytes(value: Any) -> bytes:
    return json.dumps(normalize(value), ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")

def digest(value: Any) -> str:
    return "sha256:" + hashlib.sha256(canonical_bytes(value)).hexdigest()

def assert_equal(actual: Any, expected: Any, message: str) -> None:
    if actual != expected:
        raise ValidationError(f"{message}: expected {expected!r}, got {actual!r}")

def _resolve_ref(root: dict, ref: str) -> dict:
    if not ref.startswith("#/"):
        raise ValidationError(f"external schema ref forbidden: {ref}")
    node: Any = root
    for part in ref[2:].split("/"):
        node = node[part.replace("~1", "/").replace("~0", "~")]
    return node

def validate_schema(instance: Any, schema: dict, root: dict, path: str = "$") -> None:
    if "$ref" in schema:
        validate_schema(instance, _resolve_ref(root, schema["$ref"]), root, path)
        return
    if "oneOf" in schema:
        matches = 0
        for branch in schema["oneOf"]:
            try:
                validate_schema(instance, branch, root, path)
                matches += 1
            except ValidationError:
                pass
        if matches != 1:
            raise ValidationError(f"{path}: oneOf matched {matches} branches")
        return
    if "const" in schema and instance != schema["const"]:
        raise ValidationError(f"{path}: const mismatch")
    if "enum" in schema and instance not in schema["enum"]:
        raise ValidationError(f"{path}: value not in enum")
    expected = schema.get("type")
    if expected:
        allowed = expected if isinstance(expected, list) else [expected]
        checks = {
            "object": lambda x: isinstance(x, dict),
            "array": lambda x: isinstance(x, list),
            "string": lambda x: isinstance(x, str),
            "integer": lambda x: isinstance(x, int) and not isinstance(x, bool),
            "boolean": lambda x: isinstance(x, bool),
            "null": lambda x: x is None,
        }
        if not any(checks[t](instance) for t in allowed):
            raise ValidationError(f"{path}: type mismatch")
    if isinstance(instance, str) and "pattern" in schema and not re.search(schema["pattern"], instance):
        raise ValidationError(f"{path}: pattern mismatch")
    if isinstance(instance, int) and "minimum" in schema and instance < schema["minimum"]:
        raise ValidationError(f"{path}: below minimum")
    if isinstance(instance, list) and "items" in schema:
        for index, item in enumerate(instance):
            validate_schema(item, schema["items"], root, f"{path}/{index}")
    if isinstance(instance, dict):
        properties = schema.get("properties", {})
        for required in schema.get("required", []):
            if required not in instance:
                raise ValidationError(f"{path}: missing {required}")
        additional = schema.get("additionalProperties", True)
        for key, value in instance.items():
            if key in properties:
                validate_schema(value, properties[key], root, f"{path}/{key}")
            elif additional is False:
                raise ValidationError(f"{path}: unknown field {key}")
            elif isinstance(additional, dict):
                validate_schema(value, additional, root, f"{path}/{key}")

def json_pointers(before: Any, after: Any, path: str = "") -> list[str]:
    if type(before) is not type(after):
        return [path or "/"]
    if isinstance(before, dict):
        changed = []
        for key in sorted(set(before) | set(after)):
            pointer = f"{path}/{key}"
            if key not in before or key not in after:
                changed.append(pointer)
            else:
                changed.extend(json_pointers(before[key], after[key], pointer))
        return changed
    if isinstance(before, list):
        return [] if before == after else [path or "/"]
    return [] if before == after else [path or "/"]

def apply_event(state: dict, event: dict) -> dict:
    next_state = copy.deepcopy(state)
    next_state["revision"] = event["revision"]
    next_state["event_sequence"] = event["sequence"]
    payload = event["payload"]
    kind = event["type"]
    if kind == "ChoiceSubmitted":
        next_state["pending"]["selected"] = next(iter(event["private"].values()))["selected"]
    elif kind == "CardMoved":
        card = next_state["cards"][payload["card_id"]]
        assert_equal(card["zone"], payload["from_zone"], "card from_zone")
        assert_equal(card["position"], payload["from_position"], "card from_position")
        card.update(zone=payload["to_zone"], position=payload["to_position"], face=payload["face"])
    elif kind == "InteractionClosed":
        assert_equal(next_state["pending"]["id"], payload["interaction_id"], "interaction id")
        next_state["pending"] = None
    elif kind == "EffectChanged":
        effect = next_state["effect"]
        assert_equal(effect["id"], payload["effect_id"], "effect id")
        assert_equal(effect["stage"], payload["from_stage"], "effect from_stage")
        assert_equal(effect["status"], payload["from_status"], "effect from_status")
        effect.update(stage=payload["to_stage"], status=payload["to_status"])
    elif kind == "EffectDequeued":
        next_state["queue"].remove(payload["effect_id"])
    elif kind == "ActionCompleted":
        assert_equal(next_state["action"]["id"], payload["action_id"], "action id")
        next_state["action"] = None
    else:
        raise ValidationError(f"unsupported event type: {kind}")
    return next_state

def projection(state: dict, viewer: str) -> dict:
    discard = [cid for cid, card in sorted(state["cards"].items(), key=lambda item: (item[1]["position"], item[0])) if card["zone"] == "DISCARD"]
    hand = [cid for cid, card in sorted(state["cards"].items(), key=lambda item: (item[1]["position"], item[0])) if card["zone"] == "HAND" and card["owner"] == viewer]
    pending = state["pending"]
    visible_pending = None if pending is None or pending["owner"] != viewer else {"id": pending["id"], "prompt": pending["prompt"]}
    return {
        "schema": "unmatched.projection/v1",
        "definition_id": state["definition_id"],
        "match_id": state["match_id"],
        "revision": state["revision"],
        "event_sequence": state["event_sequence"],
        "viewer": viewer,
        "public": {"lifecycle": state["lifecycle"], "action": None if state["action"] is None else state["action"]["type"], "discard": discard},
        "private": {"hand": hand},
        "pending": visible_pending,
    }

def walk_keys(value: Any):
    if isinstance(value, dict):
        for key, child in value.items():
            yield key
            yield from walk_keys(child)
    elif isinstance(value, list):
        for child in value:
            yield from walk_keys(child)

def validate_data(suite: dict, audit: dict, schema: dict) -> dict:
    validate_schema(suite, schema, schema)
    validate_schema(audit, schema, schema)
    ids = [case["id"] for case in suite["cases"]]
    if len(ids) != len(set(ids)):
        raise ValidationError("duplicate case id")
    by_kind = {case["kind"]: case for case in suite["cases"]}
    assert_equal(set(by_kind), REQUIRED_KINDS, "required case coverage")
    for definition_id, entry in suite["definitions"].items():
        assert_equal(entry["sha256"], digest(entry["content"]), f"definition hash {definition_id}")
        assert_equal(entry["ruleset"], entry["content"]["ruleset"], f"definition ruleset {definition_id}")
    for state_id, entry in suite["states"].items():
        assert_equal(entry["sha256"], digest(entry["value"]), f"state hash {state_id}")
        if OPERATIONAL_KEYS.intersection(walk_keys(entry["value"])):
            raise ValidationError(f"operational field in deterministic state {state_id}")
    vectors = suite["canonicalization"]
    assert_equal(canonical_bytes(vectors["nfc"]["a"]), canonical_bytes(vectors["nfc"]["b"]), "NFC bytes")
    assert_equal(vectors["nfc"]["bytes"].encode(), canonical_bytes(vectors["nfc"]["a"]), "NFC declared bytes")
    assert_equal(vectors["nfc"]["sha256"], digest(vectors["nfc"]["a"]), "NFC hash")
    assert_equal(vectors["escaping"]["bytes"].encode(), canonical_bytes(vectors["escaping"]["value"]), "escaping bytes")
    assert_equal(vectors["escaping"]["sha256"], digest(vectors["escaping"]["value"]), "escaping hash")
    try:
        strict_loads(vectors["duplicate_json"])
    except ValidationError:
        pass
    else:
        raise ValidationError("duplicate-key vector was accepted")
    for case in suite["cases"]:
        if case["definition_id"] not in suite["definitions"]:
            raise ValidationError(f"unpinned definition: {case['definition_id']}")
    create = by_kind["create"]
    assert_equal(create["command"]["type"], "CreateMatch", "create command")
    assert_equal(create["identity"]["scope"], "create", "create scope")
    assert_equal(create["fingerprint"], digest(create["identity"]), "create fingerprint")
    join = by_kind["join"]
    assert_equal(join["command"]["type"], "JoinMatch", "join command")
    assert_equal(join["identity"]["scope"], "join", "join scope")
    assert_equal(join["identity"]["actor"], {"absent": True}, "join actor absence")
    assert_equal(join["fingerprint"], digest(join["identity"]), "join fingerprint")
    initial = suite["states"][join["initial_state"]]["value"]
    final = suite["states"][join["expected_state"]]["value"]
    joining_player = join["providers"]["player_ids"][0]
    if any(player["player_id"] == joining_player for player in initial["players"]):
        raise ValidationError("JoinMatch player already existed")
    if not any(player["player_id"] == joining_player for player in final["players"]):
        raise ValidationError("JoinMatch player not established")
    idem = by_kind["idempotency"]
    assert_equal(idem["record"]["fingerprint"], idem["same"]["fingerprint"], "same fingerprint")
    assert_equal(idem["same"]["response"], "durable_result", "same response")
    assert_equal(idem["conflict"]["response"], "derived_collision", "conflict response")
    assert_equal(idem["conflict"]["code"], "DUPLICATE_CONFLICT", "conflict code")
    assert_equal(idem["conflict"]["record_delta"], 0, "conflict record delta")
    assert_equal(idem["conflict"]["event_delta"], 0, "conflict event delta")
    if idem["record"]["fingerprint"] == idem["conflict"]["fingerprint"]:
        raise ValidationError("conflict fingerprint did not differ")
    reconnect = by_kind["reconnect"]
    state_entry = suite["states"][reconnect["state"]]
    assert_equal(reconnect["state_hash_before"], state_entry["sha256"], "reconnect before hash")
    assert_equal(reconnect["state_hash_after"], state_entry["sha256"], "reconnect after hash")
    assert_equal(state_entry["value"]["pending"]["id"], reconnect["pending_id"], "reconnect pending id")
    replay = by_kind["replay"]
    state = copy.deepcopy(suite["states"][replay["snapshot"]]["value"])
    assert_equal(audit["fixture_id"], suite["fixture_id"], "audit fixture id")
    assert_equal(audit["case_id"], replay["id"], "audit case id")
    assert_equal(audit["snapshot_sha256"], digest(state), "audit snapshot hash")
    if len(audit["entries"]) != len(replay["events"]):
        raise ValidationError("audit/event length mismatch")
    expected_sequence = state["event_sequence"] + 1
    for event, recorded in zip(replay["events"], audit["entries"]):
        assert_equal(event["sequence"], expected_sequence, "contiguous event sequence")
        assert_equal(event["definition_id"], replay["definition_id"], "event definition pin")
        before = copy.deepcopy(state)
        state = apply_event(state, event)
        assert_equal(recorded["sequence"], event["sequence"], "audit sequence")
        assert_equal(recorded["type"], event["type"], "audit type")
        assert_equal(recorded["pre"], digest(before), "audit pre hash")
        assert_equal(recorded["paths"], json_pointers(before, state), "audit changed paths")
        assert_equal(recorded["post"], digest(state), "audit post hash")
        expected_sequence += 1
    expected_state = suite["states"][replay["expected_state"]]
    assert_equal(state, expected_state["value"], "replay final state")
    assert_equal(audit["final_sha256"], expected_state["sha256"], "audit final hash")
    for expected in replay["projections"]:
        actual = projection(state, expected["viewer"])
        assert_equal(expected["value"], actual, f"projection {expected['viewer']}")
        assert_equal(expected["sha256"], digest(actual), f"projection hash {expected['viewer']}")
        serialized = json.dumps(actual, ensure_ascii=False)
        for forbidden in expected["forbidden"]:
            if forbidden in serialized:
                raise ValidationError(f"hidden value leaked to {expected['viewer']}: {forbidden}")
    return {"fixture": suite["fixture_id"], "cases": sorted(ids), "final_hash": expected_state["sha256"], "status": "PASS"}

def validate_repository() -> dict:
    return validate_data(strict_load(FIXTURE), strict_load(AUDIT), strict_load(SCHEMA))

def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--report", default="engine-foundation-validation.json")
    args = parser.parse_args()
    try:
        report = validate_repository()
    except Exception as exc:
        report = {"status": "FAIL", "error": str(exc)}
        Path(args.report).write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
        print(json.dumps(report, indent=2))
        return 1
    Path(args.report).write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(report, indent=2))
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
