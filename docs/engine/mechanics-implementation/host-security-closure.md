# Rules Wave 1 host/security compatibility closure

Issue #46 / PR #40 continues reviewed Head `00eb337192bfce2a33d055325fb444cd02bf9035` with a bounded Lead-authorized closure.

## Hidden-information control flow

Wave 1 no longer registers generic `players/*/private_zones` state access. Private player zones require a future owner-bound typed selector and cannot be read through the generic expression language.

Stage conditions and prerequisites must infer as `public bool`. Values marked `owner_private` or `opaque` cannot decide publicly observable stage execution, skip events, operations, or interactions. This closes indirect disclosure where a secret was not emitted as payload but changed public behavior.

## Trusted host binding

`actor_player_id` is reserved transport metadata owned by Core:

```text
name: actor_player_id
type: player_ref
visibility: public
source: trusted Core host
persisted: ProcedureRef bindings
client-authoritative: no
definition-redefinable: no
```

Rules separates reserved host bindings from definition captured/input data. The actor receives its own closed validation and is canonical-copied into serialized procedure state. On resume Core may omit Context because the actor remains in the persisted procedure. `ResolutionInput.Context` may not supply the reserved actor key; Core authority is carried only by the procedure binding.

Definition data remains closed: ordinary unknown bindings, missing required values, wrong nested types, and undeclared fields are still rejected.

## Evidence

Focused tests cover:

- direct private-zone path rejection;
- private-tainted condition and prerequisite rejection;
- invalid and mismatched actor rejection;
- reserved-name redefinition rejection;
- actor persistence across serialized pause/resume;
- real Core `CreateMatch -> JoinMatch -> StartAction -> SubmitChoice` using the concrete Rules engine and actual injected binding map.

No Core-owned production path, shared contract, projection, reducer, operation disposition, or replay implementation changes in this closure.
