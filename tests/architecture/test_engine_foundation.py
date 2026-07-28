import importlib.util
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SPEC = importlib.util.spec_from_file_location("foundation_validator", ROOT / "scripts/validate_engine_foundation.py")
VALIDATOR = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(VALIDATOR)


class FoundationValidationTests(unittest.TestCase):
    def setUp(self):
        self.schema = VALIDATOR.strict_load(VALIDATOR.SCHEMA)
        self.suite = VALIDATOR.strict_load(VALIDATOR.FIXTURE)
        self.audit = VALIDATOR.strict_load(VALIDATOR.AUDIT)

    def assertRejected(self, suite=None, audit=None):
        with self.assertRaises(VALIDATOR.ValidationError):
            VALIDATOR.validate_data(suite or self.suite, audit or self.audit, self.schema)

    def case(self, kind):
        return next(case for case in self.suite["cases"] if case["kind"] == kind)

    def test_valid_candidate_passes(self):
        self.assertEqual("PASS", VALIDATOR.validate_data(self.suite, self.audit, self.schema)["status"])

    def test_duplicate_json_key_rejected(self):
        with self.assertRaises(VALIDATOR.ValidationError):
            VALIDATOR.strict_loads('{"a":1,"a":2}')

    def test_nfc_created_duplicate_key_rejected(self):
        with self.assertRaises(VALIDATOR.ValidationError):
            VALIDATOR.canonical_bytes({"é": 1, "e\u0301": 2})

    def test_missing_join_case_rejected(self):
        self.suite["cases"] = [case for case in self.suite["cases"] if case["kind"] != "join"]
        self.assertRejected()

    def test_missing_reconnect_case_rejected(self):
        self.suite["cases"] = [case for case in self.suite["cases"] if case["kind"] != "reconnect"]
        self.assertRejected()

    def test_duplicate_case_id_rejected(self):
        self.suite["cases"][1]["id"] = self.suite["cases"][0]["id"]
        self.assertRejected()

    def test_unknown_nested_command_field_rejected(self):
        self.case("join")["command"]["payload"]["admin"] = True
        self.assertRejected()

    def test_missing_definition_pin_rejected(self):
        self.case("replay")["definition_id"] = "missing-definition"
        self.assertRejected()

    def test_conflict_cannot_be_durable_result(self):
        self.case("idempotency")["conflict"]["response"] = "durable_result"
        self.assertRejected()

    def test_conflict_cannot_create_second_record(self):
        self.case("idempotency")["conflict"]["record_delta"] = 1
        self.assertRejected()

    def test_join_cannot_supply_existing_actor(self):
        self.case("join")["identity"]["actor"] = "p2"
        self.assertRejected()

    def test_reconnect_hash_must_remain_invariant(self):
        self.case("reconnect")["state_hash_after"] = "sha256:" + "0" * 64
        self.assertRejected()

    def test_event_sequence_gap_rejected(self):
        self.case("replay")["events"][2]["sequence"] += 1
        self.assertRejected()

    def test_event_definition_pin_rejected(self):
        self.case("replay")["events"][0]["definition_id"] = "other-rules"
        self.assertRejected()

    def test_undeclared_state_mutation_rejected(self):
        self.case("replay")["events"][1]["payload"]["to_position"] = 3
        self.assertRejected()

    def test_transition_hash_rejected(self):
        self.audit["entries"][2]["post"] = "sha256:" + "f" * 64
        self.assertRejected(audit=self.audit)

    def test_hidden_value_exposure_rejected(self):
        projection = self.case("replay")["projections"][1]
        projection["value"]["pending"] = {"id": "i1", "prompt": "choose-card"}
        projection["sha256"] = VALIDATOR.digest(projection["value"])
        self.assertRejected()

    def test_operational_presence_in_state_rejected(self):
        state = self.suite["states"]["reconnect"]["value"]
        state["connection_status"] = "CONNECTED"
        self.suite["states"]["reconnect"]["sha256"] = VALIDATOR.digest(state)
        self.assertRejected()


if __name__ == "__main__":
    unittest.main()
