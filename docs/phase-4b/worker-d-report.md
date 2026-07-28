# Phase 4B Worker D — reconciliation current state

## 1. Identity

Branch: `phase-4b-worker-d-latest`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Scope: newest / freshness-sensitive competitive fighter corpus  
Assigned: **16**  
Shared files changed: **none**  
Merge to `main`: **not performed**.

Bruce Lee remains one canonical fighter lineage spanning the standalone release and later Bruce Lee vs. Muhammad Ali membership. TMNT manifests contain competitive hero behavior only. Shredder/Krang contain official competitive Hero Deck behavior only; Adventures villain/scenario AI is excluded.

## 2. Final status matrix

The unified status contract is applied to both fighter and deck manifests. `verified` requires evidence + semantics completeness **and** sufficient existing shared runtime; missing D-REQ integration therefore produces `partial` even when the published behavior is known.

| Fighter | Pair status | Evidence | Semantics | Integration | Main reason |
| --- | --- | --- | --- | --- | --- |
| Bruce Lee | **verified** | verified | verified | ready | Existing generic operations are sufficient. |
| Muhammad Ali | **partial** | verified | verified | requires_shared_extension | Float reach + temporary cannot-leave-space. |
| Blackbeard | **partial** | verified | verified | requires_shared_extension | Multi-actor shared-Treasury ransom arbitration. |
| Chupacabra | **partial** | verified | verified | requires_shared_extension | Choose-distinct + result binding on staged effects. |
| Loki | **partial** | qualified | qualified | requires_shared_extension | Explicit uncertain recipient normalization + contextual movement + opaque selection. |
| Pandora | **partial** | verified | verified | requires_shared_extension | Resumable reveal/resolve/stop loop and staged result binding. |
| Leonardo | **partial** | verified | qualified | requires_shared_extension | Staged multi-fighter repositioning; For Sensei uses project ordering normalization. |
| Donatello | **partial** | verified | verified | requires_shared_extension | Movement/look-dependent staged choices and repeated per-player choices. |
| Michelangelo | **partial** | verified | verified | requires_shared_extension | Hand-size policy, play-event history, choose-two-distinct fanout. |
| Raphael | **partial** | verified | verified | requires_shared_extension | Cross-window cancellation-result binding. |
| Rosie the Riveter | **partial** | qualified | verified | requires_shared_extension | Ordered upgrades, maneuver-only modifier, runtime ranged state, restricted attack action. |
| John Henry | **partial** | verified | verified | requires_shared_extension | `path_ref_set`, positioned tracks, contextual movement, restricted attack action. |
| Wyatt Earp | **partial** | verified | verified | requires_shared_extension | Alternating staged choices + restricted immediate attacks. |
| George Washington | **partial** | verified | verified | requires_shared_extension | Pre-defense Ruse state, staged private choices, restricted attack action. |
| Shredder | **partial** | verified | verified | requires_shared_extension | Path-resident Foot Soldiers + attack network + delayed opponent-turn obligation. |
| Krang | **partial** | qualified | verified | requires_shared_extension | Random die/reroll, dynamic BOOST, alternate commitment, locks, movement modifier. |

Final counts: **1 verified / 15 partial / 0 blocked**.

## 3. D-REQ requirements

### D-REQ-001
- id: `D-REQ-001`
- family: `resolution_and_choices`
- severity: `required_for_verified_integration`
- affects: `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`
- established_rule: Dependent effects resolve in printed order; later choice domains may depend on earlier movement, disclosure, cancellation, random or damage results; choose-N choices may require distinct selections and actor-specific ordering.
- missing_capability: Persisted/resumable stages with typed bindings, operation-result bindings, domain reevaluation, repeated/alternating actors and combat-participant rebinding.
- proposed_generic_contract: `effects[].stages[]`; every choice has `id/owner/visibility/optional/domain/bind`; operations may `bind` results; later stages reference `bound.*`; `repeat_for_each` supports ordered actor/item iteration; continuation state survives reconnect.
- evidence_status: `verified`

### D-REQ-002
- id: `D-REQ-002`
- family: `resolution_and_choices`
- severity: `required_for_verified_integration`
- affects: `pandora`
- established_rule: Pandora repeatedly reveals and resolves Miseries, may stop between iterations, must stop at the feather threshold, then takes damage based on revealed Miseries.
- missing_capability: Resumable repeat/stop resolution with a forced-stop condition and post-loop effect.
- proposed_generic_contract: `procedure.kind: resumable_repeat_stop` with source zone, per-iteration operations, optional continue decision, forced-stop predicate, persisted loop state and `after_stop` operations.
- evidence_status: `verified`

### D-REQ-003
- id: `D-REQ-003`
- family: `resources_actions_and_turn_control`
- severity: `required_for_verified_integration`
- affects: `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`
- established_rule: Some gained actions must be used immediately as the next action, may permit only Attack, and may additionally bind a required target.
- missing_capability: Restricted expiring action-credit semantics distinct from ordinary `GAIN_ACTION`.
- proposed_generic_contract: `GAIN_ACTION.credit` with `must_be_next`, `allowed_actions`, optional `required_target`, and `expires_if_not_immediately_used`.
- evidence_status: `verified`

### D-REQ-004
- id: `D-REQ-004`
- family: `derived_attributes_and_modifiers`
- severity: `required_for_verified_integration`
- affects: `loki`, `rosie-the-riveter`, `john-henry`, `krang`
- established_rule: Movement can be dynamically modified by current game state and may apply only to a Maneuver, only to owner-controlled movement, or only to marked spaces.
- missing_capability: Context-sensitive movement-value/distance/space-cost modifiers that do not leak into forced/effect movement.
- proposed_generic_contract: Derived movement modifiers declare `kind`, target, value/amount rule, movement context and exclusions (`maneuver`, controller, `effect_move`, `place`, forced movement).
- evidence_status: `verified`

### D-REQ-005
- id: `D-REQ-005`
- family: `movement_targeting_and_combat_legality`
- severity: `required_for_verified_integration`
- affects: `muhammad-ali`, `rosie-the-riveter`, `shredder`
- established_rule: A source can change legal attack reach/network or runtime fighter attack type without changing immutable printed identity.
- missing_capability: Runtime attack-legality and attack-type modifiers.
- proposed_generic_contract: Derived combat-legality modifiers for reach/path-network eligibility plus runtime attack-type layer (`melee`/`ranged`) with source lifetime.
- evidence_status: `verified`

### D-REQ-006
- id: `D-REQ-006`
- family: `movement_targeting_and_combat_legality`
- severity: `required_for_verified_integration`
- affects: `muhammad-ali`, `krang`
- established_rule: A fighter may be forbidden to leave its current space until turn end while still being removable by defeat.
- missing_capability: Turn-scoped movement/placement legality lock.
- proposed_generic_contract: Persisted `temporary_space_locks` entry with target, captured/current space, lifetime, blocked operations (`MOVE`/`PLACE` as applicable) and explicit defeat/removal exclusion.
- evidence_status: `verified`

### D-REQ-007
- id: `D-REQ-007`
- family: `battlefield_components_and_paths`
- severity: `required_for_verified_integration`
- affects: `john-henry`, `shredder`
- established_rule: Track/Foot Soldier tokens have stable instance identity and board anchors; John Henry needs distinct traversed `path_ref` history; Shredder tokens live on paths and can remain deployed while a path is absent.
- missing_capability: First-class `space_ref`/`path_ref` component instances, path traversal history, path disappearance/dormancy and path-return lifecycle.
- proposed_generic_contract: Component pools with per-instance state/location; `path_ref_set` history; occupancy rules; movement-result traversed-path bindings; absent-path state leaves component deployed but dormant for path/adjacency effects when source rules require it.
- evidence_status: `verified`

### D-REQ-008
- id: `D-REQ-008`
- family: `resources_actions_and_turn_control`
- severity: `required_for_verified_integration`
- affects: `george-washington`
- established_rule: A Ruse is attached after attack declaration but before defense selection; defender may pay the random-discard veto; token then follows ready/attached/used/regained lifecycle.
- missing_capability: Pre-defense combat marker window and source-owned token attachment lifecycle.
- proposed_generic_contract: Explicit `attack_declared_pre_defense` stage; token transitions `ready -> attached -> used`, defender veto stage before defense commitment, combat-scoped attachment query, and `regain_one` transition.
- evidence_status: `verified`

### D-REQ-009
- id: `D-REQ-009`
- family: `search_randomness_and_disclosure`
- severity: `required_for_verified_integration`
- affects: `krang`
- established_rule: Each Ultimate Destruction symbol rolls independently; number supplies the value; X is zero and suppresses that symbol's effect; Krang may deactivate a machine after a roll to reroll repeatedly while able to pay.
- missing_capability: Authoritative random-table result with post-roll paid reroll/replacement loop.
- proposed_generic_contract: `random_resolution` names a table and binds an authoritative result; optional reroll window consumes declared state/resource, records replacement RNG result, and repeats while legal; replay persists all accepted random results.
- evidence_status: `verified`

### D-REQ-010
- id: `D-REQ-010`
- family: `boost_pipeline`
- severity: `required_for_verified_integration`
- affects: `krang`
- established_rule: Some Krang cards have a BOOST resolved by the Ultimate Destruction die rather than a printed numeric BOOST.
- missing_capability: Dynamic/non-numeric BOOST source in all pipelines that consume BOOST.
- proposed_generic_contract: `boost: dynamic` plus `boost_resolution` that returns authoritative numeric/zero result; Maneuver boost, blind BOOST and effects reading BOOST consume the resolved value through the same pipeline.
- evidence_status: `verified`

### D-REQ-011
- id: `D-REQ-011`
- family: `movement_targeting_and_combat_legality`
- severity: `required_for_verified_integration`
- affects: `krang`
- established_rule: Android Arms: Missiles may be committed face up as an alternate ranged attack mode.
- missing_capability: Source-defined combat commitment mode that changes visibility and attack type without mutating printed card metadata.
- proposed_generic_contract: Declarative `play_modes[]` with commitment visibility, attack type, legality predicate and whether the mode replaces normal face-down commitment.
- evidence_status: `verified`

### D-REQ-012
- id: `D-REQ-012`
- family: `card_zones_and_auxiliary_systems`
- severity: `required_for_verified_integration`
- affects: `rosie-the-riveter`
- established_rule: Four named upgrades form an ordered row; effects activate/deactivate the next token, activating/deactivating can reorder the row, cards query active count, and D-Day snapshots count before activating all.
- missing_capability: Ordered token-row identity/state with reusable transitions and queries.
- proposed_generic_contract: `ordered_token_row` with stable token IDs and generic transitions `activate_next`, `deactivate_next`, `activate_all`, `deactivate_all_in_order`; queries `count_active` and `is_active`; operations use `CHANGE_STATE` and bindings rather than fighter-specific handlers.
- evidence_status: `qualified`

### D-REQ-013
- id: `D-REQ-013`
- family: `card_zones_and_auxiliary_systems`
- severity: `required_for_verified_integration`
- affects: `michelangelo`
- established_rule: Michelangelo has source-defined starting and maximum hand size of 3.
- missing_capability: Fighter-level starting-hand and end-turn maximum-hand policy layers.
- proposed_generic_contract: Declarative `hand_size_policy` with `starting_hand_size` and `maximum_hand_size`, consumed by setup and cleanup/discard-limit logic.
- evidence_status: `verified`

### D-REQ-014
- id: `D-REQ-014`
- family: `identity_history_and_provenance`
- severity: `required_for_verified_integration`
- affects: `michelangelo`
- established_rule: Cowabunga counts other cards actually played by its controller this turn, not cards discarded, revealed, boosted, moved between zones or otherwise manipulated.
- missing_capability: Canonical per-turn card-play event history.
- proposed_generic_contract: Persist immutable `PLAY_CARD` events with card instance, controller, action/combat context and turn; expose count/query selectors distinct from zone-transition events.
- evidence_status: `verified`

### D-REQ-015
- id: `D-REQ-015`
- family: `resources_actions_and_turn_control`
- severity: `required_for_verified_integration`
- affects: `shredder`
- established_rule: All According to Plan grants an opponent an action and creates an obligation lasting through that player's turn: if they do not play a Scheme, they discard at turn end.
- missing_capability: Persisted delayed cross-player obligation with satisfaction event and expiry consequence.
- proposed_generic_contract: `delayed_turn_obligations` entries name source owner, subject player, lifetime/expiry, satisfaction event and unsatisfied consequence; replay/reconnect preserves obligation state.
- evidence_status: `verified`

### D-REQ-016
- id: `D-REQ-016`
- family: `resources_actions_and_turn_control`
- severity: `required_for_verified_integration`
- affects: `blackbeard`
- established_rule: Ransom clauses may be ignored when any opponent elects to pay the stated number of shared Treasury doubloons to Blackbeard; only one accepted payment resolves the window.
- missing_capability: Multi-actor payment arbitration over shared resource plus deterministic resource transfer/decline branch.
- proposed_generic_contract: Staged payment window with eligible actor set, amount/resource, first accepted payment closing the window, `SPEND_RESOURCE` + `GAIN_RESOURCE` atomic transfer, and explicit unpaid branch. Ordinary deterministic transfers use the existing resource operations directly.
- evidence_status: `verified`

### D-REQ-017
- id: `D-REQ-017`
- family: `search_randomness_and_disclosure`
- severity: `required_for_verified_integration`
- affects: `loki`
- established_rule: For Svadilfari's Lure, Loki selects a concrete card instance from another player's hand without learning non-TRICK face identity from the selection UI; the chosen card is then revealed and discarded.
- missing_capability: Partial-information selector over hidden card instances.
- proposed_generic_contract: Choice domain exposes stable opaque instance handles/card backs only according to chooser knowledge; server validates selected handle; identity is disclosed only by subsequent `REVEAL`; hidden metadata never reaches chooser client before disclosure.
- evidence_status: `verified`

D-REQ count: **17**.

## 4. Evidence and project interpretations

### Loki multiplayer non-combat TRICK recipient

```yaml
id: loki-noncombat-trick-recipient
status: project_normalization
confidence: medium
behavior: If multiple opponents are eligible, Loki controller chooses recipient.
replacement_condition: Replace if authoritative multiplayer ruling specifies chooser.
```

The current rules establish that a Loki-played TRICK goes to an opponent's hand but no authoritative source found in the targeted pass identifies the chooser for a non-combat TRICK with multiple eligible opponents. This is deterministic for implementation but remains evidence-qualified and is not represented as an official ruling.

### Leonardo — For Sensei ordering

Project normalization, confidence **medium**: the controller chooses the sequential order in which the two combat fighters make their optional up-to-1 moves. Replace if an authoritative ruling later specifies a different ordering authority.

### Rosie upgrade provenance by fact

| Fact | Evidence |
| --- | --- |
| Initial row `Merlin Engine -> Cavity Magnetron -> Sedgley Fist Gun -> Whizbang` | `official_set_rules`, durable |
| Merlin Engine maneuver-only optional +1 space | `official_set_rules`, durable |
| Upgrade component identities/layout | `official_component_image`, durable |
| Cavity Magnetron = optional draw after Rosie attacks | `physical_component_transcription`, transient |
| Sedgley Fist Gun = Rosie's attacks +1 | `physical_component_transcription`, transient |
| Whizbang = Rosie becomes ranged | `physical_component_transcription`, transient |

The three owner-supplied physical-card readings are preserved as transcription evidence and are **not** relabeled as publisher-hosted rulebook text.

### Shredder disappearing paths

The manually read current official Archive verdict is stored as `official_ruling_transcription`, durability `transient`:

- Heorot: Foot Soldier remains deployed while the path is absent; it does not count for Shredder adjacency effects, but it does count as deployed for Perplexing Tactics; normal behavior resumes if the path returns.
- Point Pleasant: same absent-path behavior, except the removed path does not return.

This is now deterministic battlefield lifecycle data under D-REQ-007, not an evidence blocker.

## 5. Validation

- Assigned fighter manifests: **16/16 present**.
- Assigned deck manifests: **16/16 present**.
- Fixed-deck quantity: **PASS — 16/16 reconcile to 30 cards**.
- `usable_by`: **PASS** across the reconciled corpus.
- References: **PASS**; all runtime refs used by reconciled manifests are declared or explicitly bound by stages/requirements.
- Source coverage: **PASS/QUALIFIED**, never blocked. Qualified corpus items are Loki recipient authority, Rosie transient physical mapping provenance, and Krang's current rules-mirror provenance.
- Semantic structure: **no blocked behavior**. Loki and Leonardo retain explicit project-normalization qualifications.
- Integration: **1 ready / 15 requires_shared_extension / 0 blocked**.
- Fan content: **PASS — no fan `/decks/...` source is accepted as canonical data**.
- Reconciliation lint: no remaining `source_defined_*`, `followup`, `choices_after`, `nested_choice`, or ordinary `REQUEST_CHOICE` construct in Worker D manifests.

## 6. Provenance limitations

Durable proof preservation is intentionally deferred to later project evidence stages, but current strength/durability is explicit in manifests.

Unique transient evidence artifacts currently retained:

1. Battle of Legends Vol. 3 physical/video component transcription used as a cross-check for Blackbeard, Chupacabra, Loki and Pandora.
2. TMNT physical/video component transcription used as a cross-check for Leonardo and Donatello.
3. Owner-supplied Rosie upgrade-component photograph transcription; materially supplies the three non-Merlin effect mappings.
4. Manually read Shredder removable-path official Archive verdict; materially supplies disappearing-path lifecycle text.

Krang is evidence-qualified because the current Ultimate Destruction procedure is retained from a published rules mirror rather than a publisher-hosted current textual rules URL. Current newest-set card corpora may also rely on visual/reference transcription where publisher pages do not expose complete machine-readable card text; this is recorded provenance weakness, not permission to substitute fan data.

## 7. Files

Worker D scope remains exactly:

- `docs/fighters/phase-4b/<16 assigned ids>.yaml`;
- `docs/cards/phase-4b/<16 assigned ids>.yaml`;
- `docs/phase-4b/worker-d-report.md`.

No shared schema/mechanics/rules/set file is modified by this worker.

## Worker 4B-D Reconciliation Handoff

Branch: `phase-4b-worker-d-latest`  
Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Exact final Head: **reported externally after this report commit; a commit cannot contain its own resulting SHA**.  
Assigned: **16**  
Verified: **1 — `bruce-lee`**  
Partial: **15 — `muhammad-ali`, `blackbeard`, `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`**  
Blocked: **0**  
Quantity validation: **PASS — all 16 fixed action decks are 30/30**  
D-REQ count: **17**  
Loki evidence status: **qualified; multiplayer non-combat recipient is medium-confidence `project_normalization`, not official ruling**  
Transient provenance items: **4 unique evidence artifacts** (BoL3 physical cross-check, TMNT physical cross-check, Rosie owner-supplied component transcription, Shredder official-ruling transcription)  
Unresolved behavior choices: **2 deterministic project normalizations** — Loki multiplayer recipient chooser and Leonardo `For Sensei` move ordering; both have replacement conditions if authoritative rulings appear  
Files changed: **33 Worker-D-owned files only**  
Shared files changed: **none**  
Merge: **not performed**.