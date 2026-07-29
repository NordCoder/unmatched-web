import copy
import importlib.util
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SPEC = importlib.util.spec_from_file_location("validate_playable_slice", ROOT / "scripts/validate_playable_slice.py")
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)

class PlayableSliceValidatorTest(unittest.TestCase):
    def setUp(self):
        import json
        self.data = json.loads(MODULE.RUNTIME_MAP.read_text(encoding="utf-8"))

    def test_repository_contract(self):
        MODULE.validate()

    def test_rejects_unknown_endpoint(self):
        changed = copy.deepcopy(self.data)
        changed["edges"][0]["to"] = "missing"
        with self.assertRaises(MODULE.ValidationError):
            MODULE.validate_battlefield(changed)

    def test_rejects_duplicate_undirected_edge(self):
        changed = copy.deepcopy(self.data)
        edge = changed["edges"][0]
        changed["edges"].append({"from": edge["to"], "to": edge["from"]})
        with self.assertRaises(MODULE.ValidationError):
            MODULE.validate_battlefield(changed)

    def test_rejects_disconnected_graph(self):
        changed = copy.deepcopy(self.data)
        isolated = changed["spaces"][-1]["id"]
        changed["edges"] = [edge for edge in changed["edges"] if isolated not in edge.values()]
        with self.assertRaises(MODULE.ValidationError):
            MODULE.validate_battlefield(changed)

    def test_rejects_valid_but_noncanonical_edge(self):
        changed = copy.deepcopy(self.data)
        changed["edges"][0] = {"from": "s01", "to": "s03"}
        with self.assertRaises(MODULE.ValidationError):
            MODULE.validate_battlefield(changed)

    def test_rejects_shifted_start_marker(self):
        changed = copy.deepcopy(self.data)
        for space in changed["spaces"]:
            if space["id"] == "s19":
                space.pop("starting_seat")
            if space["id"] == "s18":
                space["starting_seat"] = 2
        with self.assertRaises(MODULE.ValidationError):
            MODULE.validate_battlefield(changed)

if __name__ == "__main__":
    unittest.main()
