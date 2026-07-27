# Phase 4B Worker B report

## Worker 4B-B Handoff

Branch: `phase-4b-worker-b-licensed`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Assigned fighters: 18  
Verified: 10 — `daredevil`, `ghost-rider`, `luke-cage`, `moon-knight`, `genie`, `cloak-and-dagger`, `black-widow`, `winter-soldier`, `doctor-strange`, `she-hulk`  
Blocked: 8 — `elektra`, `bullseye`, `dr-ellie-sattler`, `t-rex`, `houdini`, `ms-marvel`, `squirrel-girl`, `spider-man`  
Quantity validation: PASS — 18/18 deck constructions reconcile.  
Shared-model extension proposals: 8 generic integration requirements; none implemented by this worker.  
Evidence status: no unresolved deck-list/card-text gaps; two narrow interpretations remain explicitly `uncertain` — T. Rex `Momentous Shift` footprint history and Squirrel Girl Go Nuts `empty` correction. Squirrel shared-damage propagation is confirmed by an official ruling and its downstream damage-count normalization is an explicit project decision.  
Worker-owned artifacts: 37 — 18 fighter manifests, 18 card manifests, this report.  
Shared semantic/control files changed: none.  
Black Panther changed: no.  
Merge to `main`: not performed.

`blocked` primarily means the current shared Phase 4A model cannot encode or execute the published behavior faithfully without a reusable extension. A blocked fighter may additionally carry an explicitly documented evidence uncertainty; semantic blocking and evidence confidence are tracked separately rather than conflated.

## Adversarial consistency audit

A second full pass was performed after the initial research pass. The audit compared published behavior, fighter manifests, card manifests and the Phase 4A schema/effect-model contracts, with special attention to `FIGHTER-SCHEMA-005`, historical state, `FX-032`/`FX-034` dependency semantics, explicit choices, cancellation, hidden information, lifecycle and occupancy.

The initial research corpus was materially complete, but it was not yet execution-complete. The audit found and corrected false-verified states, missing printed effects, lost IF/THEN dependencies, unbound choices, undeclared historical state and several ad-hoc placeholders that looked structured but had no shared runtime meaning.

### High-impact corrections

#### Status / model integrity

- `bullseye` changed from `verified` to `blocked`: its five-space graph-distance attack rule had been represented by an undocumented structured `attack_range_rule` object stored in enum state. That is neither schema-valid state modeling nor an executable attack-legality policy.
- `ms-marvel` changed from `verified` to `blocked` for the same generic graph-distance attack-legality gap at range two.
- `t-rex` card manifest changed from `verified` to `blocked`: Momentous Shift, Closer Than She Appears and Ripples in the Water themselves require the unresolved two-space footprint model.
- `squirrel-girl` card manifest changed from `verified` to `blocked`: multiple cards require small-fighter shared-space/capacity/path semantics, not merely the fighter ability.
- `dr-ellie-sattler` card manifest changed from `verified` to `blocked`: positioned Insight-token selectors/transitions depend directly on B-EXT-006.
- `houdini` fighter status changed from `verified` to `blocked` to match the deck's first-class BOOST-source/history and discard-provenance requirements.

#### Factual transcription / omitted rule corrections

- Daredevil `Grappling Hook`: corrected from `MOVE 2` to `PLACE Daredevil in any space in his zone`.
- Houdini `All Part of the Show`: added explicit cancellation protection.
- Houdini `And the Beautiful Bess!`: restored its draw when that card is discarded by an opponent's effect.
- Houdini `Flourish`: restored its ordinary optional DURING COMBAT BOOST effect in addition to BOOSTED WITH.
- Luke Cage: restored the fighter-level rule that, while defending, Luke Cage wins a combat when he takes no combat damage, including an undefended attack.
- Ghost Rider `Penance Stare`: corrected so the first opposing-card BOOST addition is mandatory; only the second application is optional and costs 2 Hellfire.
- Winter Soldier: BRAINWASHED suppression is now scoped to the negative/star clauses rather than swallowing ordinary portions of the same cards.

#### Costs and IF/THEN dependencies

Effects whose printed consequence depends on a successful payment/action are no longer modeled as unrelated sequential operations:

- Ghost Rider Hellfire spends (`Spirit of Vengeance`, `Penance Stare`, `I Finally Escaped Hell`, `I Brought the Devil With Me`, `Deal With the Devil`).
- Elektra `The Fist` self-damage before return-to-hand/action gain.
- Genie `Your Wish Is My Command` and `Imprisoned Wrath` two-card discard costs.
- She-Hulk `The Savage She-Hulk` additional-action spend.
- Ms. Marvel `Friends and Family` additional-action spend.
- Cloak & Dagger `Channel the Dark`: no action is gained if the required placement cannot occur.
- Cloak & Dagger `Living Shadow`: defender/value changes do not occur if the required swap cannot occur.
- Luke Cage `Got My Back?`: defender changes only if the swap occurs.

#### Explicit choice / resolution-order corrections

Choice domains that were previously only implied by names such as `chosen_card`, `chosen_opponent` or `chosen_fighter` are now bound explicitly. Important corrections include:

- Genie `Three Wishes`: all three branches now contain executable operations rather than prose labels; the value-4 branch sets explicit turn-scoped state.
- She-Hulk `Legalese`: concrete draw/discard branches with the correct chooser based on combat result.
- Ms. Marvel `Shrink! Shrink! Shrink!`: executes both effects when no zones are shared, otherwise exposes the actual one-of-two choice.
- Houdini `For My Next Trick`, `Vanishing Act`, `Set the Stage`, `Sleight of Hand` and BOOSTED WITH hand-selection effects.
- Black Widow `The Kinshasa Directive`: binds the chosen opponent before that opponent chooses the discarded cards.
- Doctor Strange `No Really, I'm a Doctor`: binds both the revealed hand card and healing target.
- Moon Knight `That's the Part I Like`: binds the selected looked-at card, then reorders the captured remainder.
- Elektra `Hands of Red`, `Snakeroot Clan`, `Mesmerize`.
- Squirrel Girl `Get 'Em Tippy-Toe!`, `Horde of Squirrels`, `Nutwork of Spies`.
- She-Hulk `Sensational` and T. Rex `65 Million Years of Gut Instinct`.
- Cloak & Dagger `Into Darkness` and `Chosen Fate`.

#### Historical-state corrections

The initial manifests sometimes evaluated past-turn predicates without storing the past state. The audit added explicit history where current board state cannot reconstruct the fact:

- Bullseye: turn-start space and whether he already won an earlier combat this turn.
- Ghost Rider: turn-start space for `The Wicked Will Burn`.
- Moon Knight: turn-start space for `Madness Will Keep You Alive`.
- Ms. Marvel and Spider-Man: turn-start space for `Momentous Shift`.
- Black Widow: whether an opposing fighter took damage during the current turn for `The Moscow Protocol`.
- Winter Soldier: declared turn-scoped `ignore_brainwashed_effects` state used by `A Boy Named Bucky`.

#### Placeholder / pseudo-model removal

- Removed the undocumented `attack_range_rule` structured state from Bullseye, Ms. Marvel and T. Rex. Bullseye/Ms. Marvel now expose the generic missing attack-legality model; T. Rex folds range/origin into its already-required multi-space footprint model.
- Removed Sattler card references to undeclared `insight_token_locations`. Insight movement is now represented only as explicit blocked token-instance composites tied to B-EXT-006.
- Ghost Rider's Hellfire maneuver is represented as the maneuver's replacement movement, up to four spaces, after optional BOOST processing; spending Hellfire fixes the effective movement at four rather than adding to a BOOST.

## Evidence decisions

### B-EVID-001 — T. Rex `Momentous Shift` footprint history

- **Status:** `uncertain`; not represented as an established official ruling.
- **Higher-authority evidence:** published large-fighter rules establish that T. Rex is one fighter occupying two individual spaces, but no captured authoritative ruling resolves the exact `started this turn in a different space` predicate when old and new footprints overlap.
- **Project interpretation:** conservative. If the turn-start footprint is `{A, B}` and the current footprint is `{B, C}`, T. Rex is still in starting space `B`, so `Momentous Shift` is treated as not satisfied. The condition is true only when the current footprint shares no space with the turn-start footprint.
- **Conflicting evidence:** current Gridbeast secondary FAQ says occupying at least one new space is sufficient and both occupied spaces need not change.
- **Disposition:** keep the project interpretation explicit and `uncertain` until an exact authoritative ruling is captured. B-EXT-001 still provides the required runtime footprint/history model independently of this evidence uncertainty.

### B-EVID-002 — Squirrel Girl Go Nuts `empty` correction

- **Status:** `uncertain`.
- **Official printed evidence:** the Teen Spirit rulebook says the start-turn placement uses an `empty` adjacent space.
- **Current evidence:** modern reference/index wording omits `empty`, and current secondary FAQ explicitly treats the printed restriction as an error and permits an otherwise legal shared-occupancy destination.
- **Project interpretation:** follow the modern non-empty-restriction behavior: Go Nuts may use an adjacent space allowed by the small-fighter occupancy rules, including a compatible occupied space.
- **Disposition:** retain the behavior but mark it `uncertain`; do not claim an official erratum until an exact authoritative correcting ruling is captured.

### B-EVID-003 — Squirrel shared damage

- **Official evidence status:** resolved for propagation. An official ruling confirms that when a small fighter takes damage, all small fighters of the same type in that space take an equal amount; with 1-health Squirrels, any positive propagated damage defeats every Squirrel in that space. The same ruling confirms the four-small-fighter capacity and shared-space adjacency.
- **Project normalization decision:** every Squirrel that actually receives the propagated damage counts as a damaged fighter. The originating effect remains the source of every resulting damage application while provenance records the primary target separately from recipients introduced by the small-fighter propagation rule.
- **Example:** an effect dealing 1 damage to one of four co-located Squirrels results in four damaged fighters and 4 total damage caused by that effect.
- **Disposition:** no remaining evidence blocker for shared-damage propagation/counting under the project normalization. Runtime support remains part of B-EXT-002.

## Orchestrator decisions

### B-EXT-001 — Multi-space fighter footprint

Affected: `t-rex`  
Established rule: T. Rex is one fighter whose extended base can occupy two spaces; movement, placement and attack origin/range operate on that footprint. The exact `Momentous Shift` overlap predicate is tracked separately as B-EVID-001 and is not claimed as an established official ruling.  
Current gap: fighter occupancy and generic MOVE/PLACE/history assume one occupied space.  
Proposed extension: reusable fighter footprint/orientation state with occupied-space validation, rotation/path semantics, attack origin/range, turn-start footprint snapshot and movement overrides.  
Integration status: blocks both fighter and footprint-dependent card verification.

### B-EXT-002 — Small-fighter shared occupancy

Affected: `squirrel-girl`  
Established rule: multiple Squirrels may share a space under small-fighter capacity rules, coexist with an ordinary/opposing fighter where the small-fighter rules permit it, use shared-space/pass-through semantics, and receive same-type propagated damage. Go Nuts' use of an already occupied destination is tracked separately as B-EVID-002 and remains uncertain.  
Current gap: normal fighter occupancy/path/damage semantics assume ordinary single-occupancy fighters and do not encode same-type shared-damage propagation/provenance.  
Proposed extension: `occupancy_class`, per-space compatibility/capacity, shared-space adjacency/pass-through and same-type co-located damage propagation with recipient provenance.  
Integration status: blocks both fighter and several card effects.

### B-EXT-003 — Off-board without defeat

Affected: `elektra`  
Established rule: on first would-be defeat, Elektra and all Hand pieces leave battlefield presence while Elektra is explicitly not defeated; she returns later under resurrection rules.  
Current gap: `DEFEAT` is incorrect and `RETURN_FIGHTER` models only the return transition.  
Proposed extension: generic lifecycle state/operation for `off_board_not_defeated`, preserving fighter identity and ownership for later return.  
Integration status: blocks Elektra verification.

### B-EXT-004 — Card used as BOOST source

Affected: `houdini`  
Established rule: a card consumed as a BOOST source can resolve its own BOOSTED WITH effect, and later effects can refer to the specific card used to BOOST a specific target.  
Current gap: `BOOST` applies value but exposes neither a source-card event nor reusable source-card history.  
Proposed extension: `card_used_as_boost`/equivalent event carrying source card instance, controller, boosted context and disposition, plus historical selectors where required.  
Integration status: blocks Houdini verification.

### B-EXT-005 — Field-only disclosure of a hidden card

Affected: `spider-man`  
Established rule: after an attack is committed against Spider-Man but before defense choice, only the attack card's printed numeric value is disclosed; card identity remains hidden.  
Current gap: `REVEAL` exposes the whole card and there is no reusable attack-commit/pre-defense-choice scalar-disclosure primitive.  
Proposed extension: field-level disclosure such as `DISCLOSE_CARD_FIELD(card, printed_value)` plus the required timing hook, without changing card visibility.  
Integration status: blocks Spider-Man verification.

### B-EXT-006 — Battlefield token instances

Affected: `dr-ellie-sattler`  
Established rule: Sattler has five reusable Insight tokens that move between supply and battlefield spaces; multiple Insight tokens may occupy the same space.  
Current gap: `token_pool` tracks quantity but current persistent-state types cannot represent positioned token multiplicity/instances.  
Proposed extension: generic battlefield-token instance/multiset model with location (`supply`/`space_ref`), same-space multiplicity, selectors/counts and return-to-supply operations.  
Integration status: blocks Sattler fighter and card effects.

### B-EXT-007 — Operation-cause / discard provenance event

Affected: `houdini` (`And the Beautiful Bess!`)  
Established rule: this particular hand card reacts only when it is discarded as the result of an opponent's effect.  
Current gap: `DISCARD` is an operation, but the shared model exposes no downstream event carrying the discarded card instance together with cause/source/controller attribution. Costs, cleanup, own effects and opponent effects therefore cannot be distinguished reliably by a later trigger.  
Proposed extension: generic operation-result provenance or a reusable discard event containing card instance, affected owner, causing effect/source, actor and disposition.  
Integration status: additional Houdini blocker independent of B-EXT-004.

### B-EXT-008 — Graph-distance attack legality

Affected: `bullseye`, `ms-marvel`  
Established rules: Bullseye may attack fighters up to five spaces away ignoring zones; Ms. Marvel may attack fighters up to two spaces away ignoring zones.  
Current gap: normal melee/ranged-zone legality has no reusable per-fighter graph-distance override. The former `attack_range_rule` objects were ad-hoc state with no schema or runtime contract.  
Proposed extension: reusable attack-legality policy supporting maximum graph distance, zone bypass and legal-origin semantics.  
Integration status: blocks Bullseye and Ms. Marvel verification. T. Rex's superficially similar range is not folded into this extension because its legal origins depend on its two-space footprint and remain part of B-EXT-001.

## Integration-critical corpus notes

- Non-standard action-card pools remain intentional: `daredevil` 22; `elektra` 20; `black-widow` 31, with `The Moscow Protocol` starting in hand before the remaining deck is shuffled and the normal starting draw occurs.
- Bullseye `Ricochet` remains reconciled as versatile, value 3, BOOST 2, x3.
- Set-level battlefield items in the assigned Marvel sets are not fighter-owned resources/card zones and were not duplicated into fighter manifests.
- No assigned fighter requires card ownership transfer; Black Widow mission acquisition changes location only, and Houdini BOOST cards retain ownership.
- Metadata/quantity verification, evidence confidence and runtime integration verification are deliberately separate: a deck may have verified printed data while a fighter remains semantic-blocked, and a specific interpretation may be marked `uncertain` without obscuring otherwise established rules.

## Validation

- Assigned coverage: PASS — 18/18 fighter manifests and 18/18 card manifests.
- Deck quantity reconciliation: PASS — 18/18.
- Published rule/card evidence coverage: QUALIFIED — ordinary printed data are covered; B-EVID-001 and B-EVID-002 remain explicitly uncertain rather than being silently promoted to official rules; B-EVID-003 propagation is official-ruling-backed and its damage-count treatment is an explicit project normalization.
- Internal status consistency: PASS after corrections — no fighter remains `verified` when its required core runtime semantics are known to be missing.
- Historical-state audit: PASS for declared Worker B turn-history predicates; T. Rex's exact Momentous Shift overlap interpretation remains B-EVID-001 `uncertain`.
- Choice/dependency audit: second-pass corrections remain recorded above; this evidence update does not make a new claim that no further semantic-hardening issues exist.
- Worker-owned path scope: PASS.
- Shared semantic/control files changed: none.
- Black Panther changed: no.
- Fan/community `unmatched.cards/decks/...` data imported: no.
- Shared extensions implemented by worker: no.
- Merge to `main`: not performed.
