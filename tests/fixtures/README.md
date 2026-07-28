# Deterministic engine fixtures

This directory is reserved for language-neutral command/event/replay fixtures.

The normative schema, canonical serialization rules, lifecycle/idempotency examples,
reconnect invariants, and snapshot-tail projection evidence are defined in
`docs/engine/deterministic-fixture-contract.md` on the `engine-foundation` base.
Fixture runners and data in this directory must conform to that contract and must
not introduce a worker-local schema variant.

Fixtures must be reviewable data, versioned explicitly, and runnable by the Go
engine without relying on wall-clock time, process-global randomness, network
state, or hidden local files. No gameplay fixture is added before its Phase 4C
requirements and battlefield evidence gates pass.
