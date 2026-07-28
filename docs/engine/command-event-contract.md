# Command and Event Contract

## Status

```text
status: draft-foundation
parent_issue: #19
```

This contract defines the language-neutral transaction boundary between authenticated clients, the authoritative command handler and persisted match history.

## 1. Command envelope

```yaml
command_id:
match_id:
actor_player_id:
expected_revision:
type:
payload: {}
client_metadata:
  client_instance_id:
  submitted_at:
```

Only `command_id`, `match_id`, `actor_player_id`, `expected_revision`, `type` and normalized `payload` participate in gameplay validation. Client timestamps never determine order.

## 2. Command result

```yaml
CommandResult:
  command_id:
  status: accepted | rejected | duplicate
  accepted_revision:
  event_sequence_range:
  rejection_code:
  projection_revision:
```

A duplicate command returns the same durable result as its first accepted/rejected idempotency record. It cannot emit a second event batch.

## 3. Validation phases

1. schema and payload normalization;
2. actor authentication and match membership;
3. command idempotency lookup;
4. expected revision policy;
5. match lifecycle;
6. active player / pending-choice authority;
7. current legal action/target domain;
8. source, zone, visibility and ownership constraints;
9. costs, resources and action permissions;
10. command-specific invariants.

A validation query must not consume RNG, mutate state or reveal private information through error detail.

## 4. Initial command catalog

### Match lifecycle

```text
CreateMatch
JoinMatch
SelectFighter
SelectBattlefield
ConfirmSetup
Concede
```

### Turn/action

```text
StartManeuver
ChooseManeuverBoost
ChooseMovementPath
PlayScheme
DeclareAttack
ChooseDefense
```

### Interaction

```text
SubmitChoice
SubmitCommittedChoice
DeclineOptionalChoice
```

Each command references runtime instance IDs or legal-action descriptor IDs. It does not send arbitrary rules expressions from the client.

## 5. Legal-action descriptors

```yaml
LegalActionDescriptor:
  legal_action_id:
  type:
  source_instance_ref:
  permission_id:
  required_inputs: []
  target_domain:
  cost_preview:
  ui_metadata:
```

`legal_action_id` is scoped to a state revision. A later revision requires regeneration/validation.

UI metadata contains only projection-safe labels, icons or definition references. It is not executable behavior.

## 6. Event envelope

```yaml
event_id:
match_id:
sequence:
revision:
type:
caused_by_command_id:
parent_event_id:
source_ref:
ruleset_version:
public_payload: {}
private_payloads: {}
```

Private payloads are keyed by authorized player IDs and are stored durably. A public replay/export must not concatenate private payloads by default.

## 7. Event batch

One accepted command produces an ordered batch:

```yaml
EventBatch:
  command_id:
  previous_revision:
  next_revision:
  events: []
  terminal_projection_hints: []
```

The batch is atomically persisted. Internal resolver continuation may emit further batches under a system-generated continuation command/transaction identity only when the architecture guarantees the match cannot be externally observed between logically atomic steps. The preferred foundation behavior is to resolve deterministic internal work in the same command transaction until input is required.

## 8. Event families

### Match/setup

```text
MatchCreated
PlayerJoined
FighterSelected
BattlefieldSelected
DeckConstructed
DeckShuffled
StartingHandDrawn
SetupPlacementEstablished
MatchStarted
```

### Turn/action

```text
TurnStarted
ActionPermissionGranted
ActionStarted
ActionCostPaid
ActionCompleted
TurnEndRequested
TurnEnded
```

### Card/zone

```text
CardMoved
CardsDrawn
CardsDiscarded
CardRevealed
CardVisibilityChanged
ZoneReordered
CardAttached
CardControlChanged
```

### Movement/battlefield

```text
MovementPathCommitted
FighterMovedStep
FighterPlaced
MovementInterrupted
ComponentDeployed
ComponentMoved
ConnectionStateChanged
SpaceLocked
```

### Combat

```text
AttackDeclared
DefenseCommitted
NoDefenseChosen
CombatCardsRevealed
CombatParticipantReplaced
CombatCardReplaced
CombatValueChanged
CombatDamageCalculated
CombatEnded
```

### Damage/health

```text
DamageProposed
DamagePrevented
DamageRedirected
DamageApplied
HealthAssigned
HealthRecovered
FighterDefeated
```

### Effect/choice

```text
EffectQueued
EffectStageStarted
InteractionOpened
ChoiceSubmitted
CommittedChoicesRevealed
InteractionClosed
EffectResolved
EffectCancelled
DelayedObligationCreated
DelayedObligationExpired
```

### Random/audit/result

```text
RandomResultEstablished
GameEnded
MatchQuarantined
```

Names are provisional; semantics and payload schemas are authoritative only after Phase 4C/Phase 7 validation.

## 9. Choice commands

A choice command contains:

```yaml
interaction_instance_id:
choice_id:
selected_values: []
```

The server checks:

- interaction ID is current;
- actor owns the choice;
- cardinality and distinctness;
- selected values are in the authoritative legal domain;
- ordering is supplied when semantically required;
- hidden/private data was visible to the actor;
- optional decline is allowed.

A choice result event persists selected runtime instance IDs and any authorized hidden payload.

## 10. Movement commands

The client submits an ordered path of stable space IDs. The server validates the complete path and may then resolve it stepwise to permit entry-trigger interruptions.

`MOVE` and `PLACE` are different command/effect semantics:

- MOVE consumes/traverses a legal path and records entered spaces;
- PLACE establishes a destination under placement legality and does not imply traversed spaces.

## 11. Attack commands

`DeclareAttack` references:

```yaml
attacker_fighter_instance_id:
defender_fighter_instance_id:
attack_card_instance_id:
legal_action_id:
```

The server derives attack type, current range policy, target legality, card usability and action cost. The client does not submit a damage value or claim that the target is in range.

`ChooseDefense` references a current combat interaction and either one legal defense card instance or an explicit no-defense choice.

## 12. Error model

Rejection codes are stable categories:

```text
INVALID_SCHEMA
NOT_AUTHENTICATED
NOT_MATCH_MEMBER
STALE_REVISION
MATCH_NOT_ACTIVE
NOT_COMMAND_OWNER
NO_PENDING_INTERACTION
ILLEGAL_ACTION
ILLEGAL_TARGET
ILLEGAL_CARD
INSUFFICIENT_COST
HIDDEN_INFORMATION_VIOLATION
DUPLICATE_CONFLICT
MATCH_QUARANTINED
```

Player-facing detail must not leak hidden state. Diagnostic detail belongs in protected operational logs.

## 13. Event application

Event application is a pure deterministic reducer:

```text
apply(GameState, Event) -> GameState
```

It cannot:

- query external services;
- consume random numbers;
- depend on wall-clock time;
- inspect binary card/battlefield assets;
- branch on fighter/card IDs;
- emit additional persisted events as an implicit side effect.

Event generation and effect resolution decide which events exist; event application only applies the recorded result.

## 14. Compatibility/versioning

- command schemas are versioned at the API boundary;
- event schemas are versioned for durable replay;
- definition/capability versions are stored with matches/events;
- migrations must preserve semantic replay or explicitly retain an old reader;
- renaming an event without a migration is not compatible.

## 15. Test requirements

The foundation test suite must prove:

- duplicate command IDs do not duplicate events;
- stale revisions are rejected or reconciled only by explicit policy;
- rejected commands produce no gameplay events;
- event sequences are total and gap-free;
- applying a persisted batch reproduces the expected state hash;
- private payloads never appear in unauthorized projections/errors;
- reconnect can resume the exact pending interaction;
- random results and shuffled order are replayed, not regenerated;
- Robin Hood/Bigfoot effects are executed through generic operations rather than identity checks.
