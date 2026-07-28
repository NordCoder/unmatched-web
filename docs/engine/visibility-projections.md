# Visibility and Player Projection Contract

## Status

```text
status: draft-foundation
parent_issue: #19
```

The server stores one authoritative state and derives a separate projection for each authorized viewer. Projections contain presentation-ready facts and legal-action descriptors, not executable rules.

## 1. Visibility classes

The foundation recognizes:

```text
PUBLIC
OWNER_PRIVATE
CONTROLLER_PRIVATE
ACTOR_PRIVATE
COMMITTED_HIDDEN
REVEALED
ADMIN_PROTECTED
```

A concrete field/card/choice may transition between classes through events. Visibility is not inferred only from zone name.

## 2. Projection targets

```text
PLAYER(player_instance_id)
SPECTATOR(policy_id) optional
ADMIN_PROTECTED authorized tooling only
```

Spectator/admin behavior is outside the first vertical slice unless explicitly enabled. It must never be implemented by returning raw authoritative state.

## 3. Projection envelope

```yaml
PlayerProjection:
  match_id:
  revision:
  event_sequence:
  viewer_player_id:
  lifecycle:
  public_state:
  self_private_state:
  opponent_hidden_state_summary:
  pending_interaction:
  legal_actions: []
  public_history:
```

Every projection is tied to an authoritative revision. Commands based on a projection still undergo server validation.

## 4. Cards

### Public card state

A public card may expose:

- runtime instance ID;
- definition ID;
- current zone and position if public;
- controller/owner where rules make them public;
- face state;
- public instance state.

### Hidden card state

An unauthorized viewer receives only the permitted abstraction:

```yaml
HiddenCardPlaceholder:
  opaque_instance_ref_optional:
  zone:
  ordinal_or_count_visibility:
```

Definition ID, printed fields and private tags are omitted.

Stable opaque instance references are used only if rules/UI need continuity without identity; otherwise hidden cards are represented by counts to avoid correlation leaks.

## 5. Hands and decks

- a player sees their own hand identities/order if order is meaningful;
- the opponent sees only authorized hand count and public modifiers;
- deck identities/order are private unless revealed by an event;
- a search that permits viewing only one field returns a field projection, not the full card definition;
- shuffle invalidates any unauthorized continuity assumptions about hidden order.

## 6. Combat cards

Before reveal:

- each player sees their own committed card;
- the opponent sees only that a legal card/no-defense commitment exists;
- simultaneous commitments remain private until the reveal event.

After reveal, the public projection exposes the combat card definitions and public values according to timing rules. Historical public reveal remains in event history even after cards move to discard.

## 7. Choices

```yaml
ProjectedInteraction:
  interaction_instance_id:
  owner:
  prompt:
  cardinality:
  optional:
  legal_options:
  submitted_state:
```

Only an authorized choice owner receives private legal options. The opponent may receive a public “waiting for opponent” descriptor without option identities.

For committed-hidden groups:

- each actor sees their own submitted value;
- other submissions remain hidden;
- the group exposes only completion status allowed by policy;
- one reveal event publishes the committed values at the canonical timing point.

## 8. Random results

Random results may contain public and private parts. The persisted event stores the authoritative result; each projection receives only authorized fields.

Examples:

- public die result: all players receive value;
- shuffled deck order: neither player receives order;
- random card selected from own hand: affected owner may see identity while opponent receives only the resulting public consequence until reveal rules apply.

## 9. Legal actions

Legal-action generation runs against authoritative state and viewer authority. The projection includes only actions the viewer may currently command.

A legal action descriptor must not leak why an alternative is illegal when the reason depends on hidden information. Server rejection messages follow the same rule.

The client cannot synthesize an action absent from the server set and cannot expand a target domain.

## 10. History

Public history contains public event payloads and sanitized summaries of private events.

Player-private history may include that player's own private payloads. It must not include the opponent's private payloads merely because later effects reference them internally.

After a hidden fact becomes publicly revealed, a reveal event exposes it. The system does not rewrite older private events.

## 11. Projection derivation

```text
project(GameState, ViewerAuthority) -> PlayerProjection
```

The projector is:

- pure;
- deterministic;
- side-effect free;
- revision-scoped;
- incapable of consuming RNG;
- incapable of changing legal state;
- deny-by-default for fields without a visibility rule.

Projection code uses capability/field visibility metadata, not fighter/card ID checks.

## 12. Redaction model

Serialization occurs after projection. Raw authoritative objects must not be serialized and then redacted by deleting known fields, because new fields would default to visible.

Instead, the projector constructs an explicit output schema from allowed fields. New authoritative fields are invisible until projection policy deliberately includes them.

## 13. Reconnect

Reconnect derives a fresh projection at the current revision. It does not replay unauthorized raw events to reconstruct the client.

The reconnect response preserves:

- the viewer's own private hand/choice state;
- public battlefield/turn/combat state;
- current legal actions or owned pending interaction;
- authorized history cursor/delta.

## 14. Security invariants

- opponent hand/deck identities never appear in unauthorized JSON/log/error fields;
- private choice values are never included in public event payloads;
- server logs containing private payloads are protected and not client-accessible;
- definition lookup endpoints enforce projection authorization when queried by runtime instance;
- cache keys include viewer authority and revision;
- one player's projection cannot be served from another player's cache entry;
- spectator support cannot reuse a player-private projection.

## 15. Projection tests

The foundation must include golden fixtures proving:

- player A and player B projections differ only where authority requires;
- adding a new private authoritative field does not expose it by default;
- committed defense cards remain hidden until reveal;
- optional choices reveal no option identities to the opponent;
- reconnect projection equals live projection at the same revision;
- public event history is identical for both players while private history remains separated;
- legal-action descriptors reveal no hidden-information reason;
- card asset path/availability does not affect gameplay visibility.
