# Phase 4B Worker B report

## Worker 4B-B Handoff

Branch: `phase-4b-worker-b-licensed`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: branch tip; exact SHA is reported in the external final handoff because a commit cannot embed its own SHA.  
Assigned fighters: 18  
Verified: 12 — `daredevil`, `bullseye`, `ghost-rider`, `luke-cage`, `moon-knight`, `genie`, `cloak-and-dagger`, `ms-marvel`, `black-widow`, `winter-soldier`, `doctor-strange`, `she-hulk`  
Blocked: 6 — `elektra`, `dr-ellie-sattler`, `t-rex`, `houdini`, `squirrel-girl`, `spider-man`  
Quantity validation: PASS — 18/18 deck constructions reconcile.  
Schema-extension proposals: 6 generic integration requirements; none implemented by this worker.  
New ambiguity/blockers: none in game-rule research; the six blocked fighters require shared-model decisions only.  
Source gaps: none.  
Files created: 37 — 18 fighter manifests, 18 card manifests, this report.  

`blocked` below means the game behavior is established but the current shared Phase 4A model cannot encode it faithfully without a reusable extension.

### Orchestrator decisions

#### B-EXT-001 — Multi-space fighter footprint

Affected: `t-rex`  
Established rule: T. Rex is one fighter whose extended base can occupy two spaces; movement, pathing and placement operate on that footprint.  
Current gap: fighter occupancy and `MOVE`/`PLACE` assume one occupied space.  
Proposed extension: reusable fighter footprint/orientation state with occupied-space validation, rotation/path semantics and movement overrides.  
Evidence: `docs/fighters/phase-4b/t-rex.yaml` (`rulings`, `sources`); published UmDb/current rulings.  

#### B-EXT-002 — Small-fighter shared occupancy

Affected: `squirrel-girl`  
Established rule: multiple Squirrels may share a space under small-fighter capacity rules, coexist with an ordinary/opposing fighter, use shared-space/pass-through semantics, and propagate damage to same-type Squirrels in that space.  
Current gap: normal occupancy and single-target damage do not represent this behavior.  
Proposed extension: `occupancy_class`, per-space compatibility/capacity, shared-space adjacency/pass-through, and same-type co-located damage propagation with preserved source attribution.  
Evidence: official Teen Spirit rulebook plus current erratum/rulings in `docs/fighters/phase-4b/squirrel-girl.yaml`.  

#### B-EXT-003 — Off-board without defeat

Affected: `elektra`  
Established rule: on first would-be defeat, Elektra and all Hand pieces leave battlefield presence while Elektra is explicitly not defeated; she returns later under resurrection rules.  
Current gap: `DEFEAT` is incorrect and `RETURN_FIGHTER` models only the return transition.  
Proposed extension: generic removal/lifecycle state for `off_board_not_defeated`, preserving fighter identity and ownership for later return.  
Evidence: published Elektra data/current rulings in `docs/fighters/phase-4b/elektra.yaml`.  

#### B-EXT-004 — Card used as BOOST source

Affected: `houdini`  
Established rule: a card consumed as a BOOST source can resolve its own `BOOSTED WITH` effect, and later effects can refer to the specific card used to BOOST a specific target.  
Current gap: `BOOST` applies a value but exposes neither a source-card event nor reusable source-card history.  
Proposed extension: `card_used_as_boost`/equivalent event carrying source card instance, controller, boosted context and disposition, plus historical selectors where required.  
Evidence: official Houdini vs. The Genie material and published card data in `docs/cards/phase-4b/houdini.yaml`.  

#### B-EXT-005 — Field-only disclosure of a hidden card

Affected: `spider-man`  
Established rule: after an attack is committed against Spider-Man but before defense choice, only the attack card's printed numeric value is disclosed; card identity remains hidden.  
Current gap: `REVEAL` exposes the whole card and there is no reusable attack-commit/pre-defense-choice hook.  
Proposed extension: field-level disclosure such as `DISCLOSE_CARD_FIELD(card, printed_value)` plus the required timing hook, without changing card visibility.  
Evidence: official Brains and Brawn rulebook in `docs/fighters/phase-4b/spider-man.yaml`.  

#### B-EXT-006 — Battlefield token instances

Affected: `dr-ellie-sattler`  
Established rule: Sattler has five reusable Insight tokens that move between supply and battlefield spaces; multiple Insight tokens may occupy the same space.  
Current gap: `token_pool` tracks quantity but current persistent-state types cannot represent positioned token multiplicity/instances.  
Proposed extension: generic battlefield-token instance or equivalent multiset model with location (`supply`/`space_ref`), same-space multiplicity, selectors/counts and return-to-supply operations.  
Evidence: current Sattler rulings in `docs/fighters/phase-4b/dr-ellie-sattler.yaml`.  

### Integration-critical corpus notes

- Non-standard action-card pools: `daredevil` 22; `elektra` 20; `black-widow` 31, with `The Moscow Protocol` starting in hand before the remaining deck is shuffled and the normal starting draw occurs.
- Bullseye research is complete: direct component imagery confirms `Ricochet` = versatile, value 3, BOOST 2, x3.
- Set-level battlefield items in the assigned Marvel sets are not fighter-owned resources/card zones and were not duplicated into fighter manifests.
- No assigned fighter requires card ownership transfer; Black Widow mission acquisition changes location only, and Houdini BOOST cards retain ownership.

### Validation

- Assigned coverage: PASS — 18/18 fighter manifests and 18/18 card manifests.
- Rule/card evidence coverage: PASS — no unresolved factual/source gaps.
- Deck quantity reconciliation: PASS — 18/18.
- Worker-owned path scope: PASS.
- Shared semantic/control files changed: none.
- Black Panther changed: no.
- Fan/community `unmatched.cards/decks/...` data imported: no.
- Shared extensions implemented by worker: no.
- Merge to `main`: not performed.
