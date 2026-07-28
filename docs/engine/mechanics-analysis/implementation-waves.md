# Mechanics Implementation Waves

## Delivery principle

Implementation order follows runtime dependency, not product release order and not fighter-by-fighter completion.

A wave is complete only when:

- its generic contracts are frozen;
- domain tests and replay fixtures pass;
- at least two corpus examples exercise each reusable capability where available;
- reconnect at every pending interaction is tested;
- no core branch dispatches by fighter/card ID;
- later waves can build on the capability without changing its meaning.

## Wave 0 — Launch kernel

Purpose: produce the smallest correct server-authoritative game using the launch pair without inventing shortcuts that later prevent full-corpus support.

Scope:

- immutable definition loading/version identity;
- runtime fighter/card instances;
- two-player setup and initial hands;
- basic battlefield graph/occupancy after Sherwood Forest passes #18;
- turn and ordinary action permissions;
- Maneuver, ordinary BOOST, Scheme and Attack actions;
- base combat windows, optional defense and cleanup;
- baseline operations from `effect-model.md`;
- exhaustion, defeat, game end, event replay and projections;
- Robin Hood and Bigfoot loaded as data.

Important restriction: Wave 0 may implement only the generic subset used by the launch slice. It may not hardcode launch fighter/card identities.

Exit fixture:

```text
create match
join two seats
setup Robin Hood vs Bigfoot
play deterministic turns
perform maneuver + boost
play scheme
resolve defended and undefended attacks
pause/resume a legal choice
exhaust deck
reconnect and replay to identical state
```

## Wave 1 — Query, staged resolution and history

Capabilities:

- `CAP-001` staged resolution procedures;
- `CAP-004` history/provenance ledger;
- `CAP-018` derived query expressions.

Deliverables:

- closed query/expression AST;
- serializable `ProcedureInstance` and stage cursor;
- result/choice binding type system;
- independent-operation versus explicit-dependency behavior;
- empty-domain semantics;
- event provenance and rebuildable history indexes;
- deterministic definition validator.

Stress fixtures:

- Alice ordered choose-two effect;
- dependent discard-top then add BOOST value;
- operation cause distinction for Houdini;
- repeated fighter-instance history for Honor Guards/Wolves;
- Shakespeare/Ciri count and threshold expressions.

Dependency effect: nearly every later exotic mechanic depends on this wave. It should be implemented before specialized occupancy, damage or resource logic.

## Wave 2 — Interactions, visibility and advanced zones

Capabilities:

- `CAP-002` interactions/visibility;
- `CAP-003` reaction/cancellation windows;
- `CAP-015` ordered/auxiliary zones;
- `CAP-016` structured search/disclosure.

Deliverables:

- private/public/opaque interaction projections;
- legal-domain snapshots and opaque handles;
- scoped effect cancellation and private-zone reactions;
- typed ordered zones and auxiliary object/card instances;
- search predicate, selected-result binding, disclosure and post-search shuffle;
- cleanup against current live zone membership.

Stress fixtures:

- Spider-Man field-only numeric disclosure;
- Loki opaque hand selection followed by reveal/discard;
- Houdini private-hand reaction canceling a bound effect;
- Little Red Basket top-object/wild binding;
- Shakespeare Line exact/overflow completion and self-removal before cleanup;
- Titania Glamour deck lifecycle;
- Arthur/Ciri structured search and authoritative shuffle.

## Wave 3 — Presence, occupancy, components and relocation

Capabilities:

- `CAP-006` fighter presence/roster;
- `CAP-007` occupancy footprints/classes;
- `CAP-008` battlefield components/paths;
- `CAP-009` relocation transitions/results.

Deliverables:

- presence independent from defeat;
- active/out-of-play setup roster;
- one-space, multi-space and shared-space occupancy policies;
- component pools and stable space/path anchors;
- movement-step and entry events;
- persisted traversed spaces/paths;
- interruption preserving reached position;
- scheduled return and dormant-player integration.

Stress fixtures:

- T. Rex two-space movement, placement, adjacency and turn-start footprint;
- Squirrel shared-space capacity, pass-through and same-space adjacency;
- Muldoon trap interruption after entry;
- Fog/Shadows/Insight positioned component lifecycle;
- John Henry tracks and Shredder path-resident components;
- Invisible Man/Elektra/Ancient Leshen off-board undefeated return;
- Buffy unselected sidekick remains out of play;
- Yennenga/Krang damage fighters crossed by a bound movement path.

## Wave 4 — Combat legality, damage, modifiers and BOOST

Capabilities:

- `CAP-010` combat legality/play modes;
- `CAP-011` combat participant rebinding;
- `CAP-012` damage/health transaction;
- `CAP-013` continuous modifier layer;
- `CAP-014` BOOST resolution.

Deliverables:

- attack target generator using graph, path network and derived reach;
- alternate commitment modes and effective attack type;
- temporary movement/placement locks;
- same-combat participant/card replacement;
- damage proposal/prevention/allocation/application pipeline;
- exact health assignment, atomic return health and health floors;
- contextual modifiers with source lifetime and family-specific composition;
- ordinary, blind and dynamic BOOST through one pipeline.

Stress fixtures:

- Bullseye/Ms. Marvel graph-distance attacks;
- Muhammad Ali/Shredder derived reach;
- Krang face-up ranged mode and cannot-leave lock;
- Yennenga same-combat attacker replacement;
- Tomoe/Titania/Eredin team or source-driven participant replacement;
- Yennenga split damage with recipient caps/residual target;
- Squirrel propagated damage provenance;
- Annie health floor where damage still counts as dealt;
- Arthur/Philippa exact health assignment;
- Willow exact-health return;
- Krang dynamic BOOST and Spike resolved blind-BOOST transform.

## Wave 5 — Randomness, resources, actions and hand policy

Capabilities:

- `CAP-005` authoritative random procedures;
- `CAP-017` action/resource permissions;
- `CAP-019` hand policy.

Deliverables:

- random-table definitions and persisted outcomes;
- paid reroll/replacement loop;
- typed action credits with restrictions, targets and expiry;
- action permission spending as cost;
- declared resource tier independent from available quantity;
- delayed obligations with satisfaction/expiry events;
- multi-actor shared-resource payment arbitration;
- starting and maximum hand policy.

Stress fixtures:

- Krang independent rolls, X result and repeated machine-paid rerolls;
- Tesla legal two-coil declaration with insufficient charge and no tier fallback;
- Ms. Marvel/She-Hulk action permission spending;
- Stars and Stripes immediate restricted attacks;
- Shredder opponent-turn Scheme obligation;
- Blackbeard first-accepted ransom payment;
- Michelangelo starting/max hand size.

## Wave 6 — Full-corpus conformance

Purpose: validate that the composed capability kernel expresses all deterministic fighter/card manifests without arbitrary identity-based code.

Deliverables:

- compile/load every deterministic manifest against versioned runtime schemas;
- capability dependency report per fighter;
- operation/query/procedure usage inventory;
- fixture coverage matrix across all fighters;
- unsupported-definition failure report;
- extension-handler requests, expected to be empty unless specifically justified;
- balance-neutral regression suite comparing canonical examples and interpretations.

Required checks:

```text
74 fighter manifests load
926 unique action-card definitions validate
51 deterministic requirement aliases resolve
0 hidden fighter/card ID branches
0 unknown operation/query/procedure kinds
0 unserializable pending interactions
0 replay divergence
```

Deadpool remains excluded only by the explicit central policy alias.

## Wave 7 — External adaptation policy

Capability:

- `CAP-020` external adaptation policy.

This wave cannot start from engine inference. The Owner must choose a policy for published physical/social effects:

- automatically representable digital predicate;
- player-confirmed fact/action;
- approved deterministic digital substitute;
- unsupported effect/card;
- unsupported fighter for ranked/online play.

After policy acceptance, the engine may add generic confirmation/policy primitives. It must not embed one-off Deadpool code.

## Parallel work allowed

The following may proceed while earlier waves are being implemented:

- Phase 4C requirement consolidation;
- deterministic fixture authoring against accepted contracts;
- battlefield graph transcription/QA;
- transport schema ADR after command/projection shapes stabilize;
- UI projection prototypes using fixtures rather than client-side rules.

The following must not proceed ahead of dependencies:

- full gameplay handlers before the staged resolver/query model;
- special occupancy hacks before the footprint/occupancy contract;
- random client animations treated as authoritative results;
- direct fighter/card implementation before generic capability registration;
- full-roster merge before Phase 4C and architecture QA.
