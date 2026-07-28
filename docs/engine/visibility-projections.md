# Visibility and Player Projection Contract

## Status

```text
status: draft-foundation
parent_issue: #19
correction_issue: #32
architecture: architecture-contract.md
state_model: state-model.md
fixture_contract: deterministic-fixture-contract.md
```

The server stores one deterministic authoritative `GameState` and derives a separate projection for each authorized viewer. External identity, seat authorization and operational presence are resolved outside the deterministic projector.

Projections contain presentation-ready facts and legal-action descriptors, not executable rules.

## 1. Visibility classes

```text
PUBLIC
OWNER_PRIVATE
CONTROLLER_PRIVATE
ACTOR_PRIVATE
COMMITTED_HIDDEN
REVEALED
ADMIN_PROTECTED
```

A field/card/choice may transition between classes through gameplay events. Visibility is not inferred only from zone name.

## 2. Viewer authority

```yaml
ViewerAuthority:
  match_id:
  player_instance_id:
  authority_version:
  policy_scope: PLAYER
```

`ViewerAuthority` is produced only after:

1. authenticating an external `principal_id`;
2. resolving an active durable `MatchAuthorityRecord`;
3. confirming the record references the requested match and player instance.

The deterministic projector never trusts a client-supplied player ID and never reads authentication tokens, sessions or connection state.

Optional spectator/admin policies require separate explicit authorization and must never return raw authoritative state.

## 3. Deterministic projection envelope

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

Every projection is tied to a deterministic authoritative revision. Commands based on a projection still undergo server authentication, binding, idempotency, revision and legality validation.

Principal IDs, session IDs, connection flags and last-seen timestamps do not appear in this envelope.

## 4. Operational presence projection

When product policy exposes connection presence, it is projected separately:

```yaml
PresenceProjection:
  generated_at_operational:
  viewer_player_id:
  participants:
    - player_instance_id:
      online:
      session_count_optional:
```

Composition is explicit:

```text
project(GameState, ViewerAuthority) -> PlayerProjection
projectPresence(OperationalPresenceRegistry, ViewerAuthority) -> PresenceProjection
compose(PlayerProjection, PresenceProjection) -> DeliveryEnvelope
```

Presence projection:

- has no gameplay revision or event sequence of its own;
- cannot alter legal actions or pending interaction;
- is excluded from deterministic state/projection hashes;
- may change without a gameplay event;
- may be omitted entirely after restart until sessions reconnect;
- must not reveal session/client identifiers to another player unless an explicit product policy authorizes a safe abstraction.

A disconnect cannot change the deterministic `PlayerProjection` at the same revision.

## 5. Cards

### Public card state

A public card may expose:

- runtime instance ID;
- definition ID;
- public zone and position;
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

Definition ID, printed fields and private tags are omitted. Stable opaque references are used only when continuity is required without identity; otherwise counts are preferred to avoid correlation leaks.

## 6. Hands and decks

- a player sees their own hand identities/order if order is meaningful;
- the opponent sees only authorized hand count and public modifiers;
- deck identities/order are private unless revealed by an event;
- field-limited look/search effects return only the authorized field projection;
- shuffle invalidates unauthorized assumptions about hidden order.

## 7. Combat cards

Before reveal:

- each player sees their own committed card;
- the opponent sees only that a legal card/no-defense commitment exists;
- simultaneous commitments remain private until the reveal event.

After reveal, public projection exposes card definitions and public values according to timing rules. Historical public reveal remains in event history after cards change zones.

## 8. Choices and pending interactions

```yaml
ProjectedInteraction:
  interaction_instance_id:
  owner:
  prompt:
  cardinality:
  optional:
  legal_options:
  submitted_state:
  resume_cursor_public_optional:
```

Only an authorized choice owner receives private legal options. The opponent may receive a public waiting descriptor without option identities.

For committed-hidden groups:

- each actor sees their own submitted value;
- other submissions remain hidden;
- the group exposes only policy-authorized completion state;
- one reveal event publishes committed values at the canonical timing point.

Reconnect returns the same authorized interaction ID, domain and submitted state at the same revision. Session loss cannot create a replacement interaction.

## 9. Random results

Persisted random-result events may contain public and private parts. Each projection receives only authorized fields.

Examples:

- public die result: all players receive value;
- shuffled deck order: neither player receives order;
- random card selected from own hand: affected owner may see identity while opponent receives only the public consequence until reveal rules apply.

Operational presence never changes random-result visibility.

## 10. Legal actions

Legal-action generation runs against deterministic state and validated `ViewerAuthority`. The projection includes only actions the bound player may command.

A descriptor must not leak why an alternative is illegal when the reason depends on hidden information. Server rejection details follow the same rule.

The client cannot synthesize an absent action or expand a target domain. Presence/connection status is not an input to competitive Unmatched legality unless a future explicit gameplay rule introduces a persisted state concept; no such rule exists in the foundation.

## 11. History

Public history contains public event payloads and sanitized summaries of private events.

Player-private history may include that player's own private payloads. It never includes the opponent's private payload merely because later effects reference it internally.

After a hidden fact becomes public, a reveal event exposes it. Older private events are not rewritten.

Principal authority changes and session presence are application/security or operational records, not gameplay history entries.

## 12. Projection derivation

```text
project(GameState, ViewerAuthority) -> PlayerProjection
```

The projector is:

- pure;
- deterministic;
- side-effect free;
- revision-scoped;
- unable to consume RNG;
- unable to mutate legal state;
- deny-by-default for fields without a visibility rule;
- unable to query authority or presence registries;
- unable to branch on fighter/card IDs.

The application layer resolves viewer authority before invoking it.

## 13. Redaction model

Serialization occurs after projection. Raw authoritative objects must not be serialized and then redacted by deleting known fields, because new fields would default to visible.

The projector constructs an explicit output schema from allowed fields. New authoritative fields remain invisible until projection policy deliberately includes them.

Operational delivery composition similarly constructs an explicit presence schema rather than attaching raw session registry objects.

## 14. Reconnect

Reconnect:

1. authenticates the external principal;
2. resolves the durable match/player binding;
3. restores deterministic state;
4. derives a fresh projection at the current revision;
5. optionally composes authorized operational presence.

It does not replay unauthorized raw events to reconstruct the client.

The deterministic reconnect response preserves:

- the viewer's private hand/choice state;
- public battlefield/turn/combat state;
- current legal actions or owned pending interaction;
- authorized history cursor/delta.

At the same revision, live and reconnect deterministic projections must be canonically equal even if presence differs.

## 15. Security invariants

- a client-supplied player ID never establishes viewer authority;
- opponent hand/deck identities never appear in unauthorized payloads/logs/errors;
- private choice values never appear in public event payloads;
- server logs containing private payloads are protected and not client-accessible;
- definition lookup by runtime instance enforces projection authorization;
- deterministic projection cache keys include match, player authority and revision;
- one player's projection cannot be served from another player's cache entry;
- presence cache/data cannot be treated as deterministic projection state;
- spectator support cannot reuse a player-private projection;
- duplicate-conflict errors reveal no prior request or private result.

## 16. Required projection evidence

Normative fixtures must prove:

- player A and player B projections differ only where authority requires;
- adding a private authoritative field does not expose it by default;
- committed defense cards remain hidden until reveal;
- optional choices reveal no option identities to the opponent;
- reconnect projection equals live projection at the same revision;
- disconnect/reconnect changes only separately declared presence output, not deterministic projection/hash;
- public history is identical while private history remains separated;
- legal-action descriptors reveal no hidden-information reason;
- unauthorized principal cannot select another `player_instance_id`;
- snapshot-plus-tail replay reconstructs the same viewer projections;
- card asset path/availability does not affect gameplay visibility.
