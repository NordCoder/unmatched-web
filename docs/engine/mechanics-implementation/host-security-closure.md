# Rules Wave 1 trusted actor ownership

Issue #48 / PR #40 replaces the final generic interaction-owner surface with a Lead-directed Wave 1 scope reduction.

## Closed owner selector

Wave 1 definitions serialize interaction ownership only as:

```json
{"owner_selector":"actor"}
```

`actor` is the only supported selector. The implementation maps it to one closed internal reference to the trusted `actor_player_id` binding. Legacy `owner` expressions remain readable only so definition loading can reject them explicitly; they are never evaluated as ownership policy.

The following owner sources are not supported in Wave 1:

- state paths;
- captured definition data;
- input context;
- operation results;
- prior choices;
- private or opaque player references;
- literals naming a specific seat or player.

Future selectors such as opponent, attacker, defender, or controller are deferred to Rules Wave 2 as separate typed contracts.

## Trusted actor path

`actor_player_id` is reserved transport metadata owned by Core:

```text
name: actor_player_id
type: player_ref
visibility: public
source: trusted Core host
required: every definition using owner_selector=actor
persisted: ProcedureRef bindings
client-authoritative: no
definition-redefinable: no
```

The Rules validator recognizes only the closed actor selector. For compatibility with the older internal validator slot it substitutes a validation-only public player literal; the trusted actor binding is never added to the general definition expression environment. Runtime binding validation requires the real host binding whenever an actor-owned choice exists, rejects missing, empty, or wrong-typed actors, and rejects attempts to provide it through `ResolutionInput.Context`.

Because the binding is included in the serialized procedure state, the exact interaction owner survives pause/resume without reinjection.

## Security boundaries retained

Wave 1 still rejects generic `players/*/private_zones` access. Stage conditions and prerequisites must be public, so private values cannot affect publicly observable execution.

Choice domain and prompt may remain private or opaque. Only the owner identity is public routing metadata derived from the trusted actor selector.

## Evidence

Focused tests prove:

- the same definition opens actor-owned interactions for both dynamically allocated Core seats;
- projected and persisted owner identity equals the current actor;
- no hardcoded player identifier is used by the definition;
- missing and wrong actor bindings are rejected;
- Context injection and definition redeclaration are rejected;
- legacy, owner-private, and opaque owner expressions are rejected;
- actor identity remains exact through serialized pause/resume;
- the full real Core Host + concrete Rules flow succeeds for both players.

No Core-owned production path, shared contract, reducer, operation disposition, projection, or replay implementation changes in this scope reduction.
