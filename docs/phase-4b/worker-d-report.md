# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head ref: `refs/heads/phase-4b-worker-d-latest`  
Head SHA: **reported in the external handoff immediately after this commit is created.** A Git commit cannot truthfully embed its own SHA in a file inside that same commit because changing the embedded SHA changes the commit object ID.  
Assigned fighters: 16  
Verified: `bruce-lee`  
Blocked: `muhammad-ali`, `blackbeard`, `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Quantity validation: **PARTIAL — PASS for Bruce Lee (30/30); BLOCKED for 15 fighters because a complete eligible per-card corpus is not publicly verified.**  
Schema-extension proposals: see below.  
New ambiguity/blockers: complete newest-release action-card corpus; source-defined range/movement/path/free-action/die semantics.  
Source gaps: see per-family matrix below.  
Files created: 33 (16 fighter manifests, 16 deck manifests, this report).

### Branch/base verification

At worker start and immediately before staging the worker tree, `phase-4b-worker-d-latest` compared **identical** to Authorized Base `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9` (`ahead_by=0`, `behind_by=0`). The worker writes only the Phase 4B fighter/card paths assigned to Worker D plus this report. No shared schema, mechanics, rules, rulings, set-registry, roadmap, or README file is changed.

### Verification result

| Fighter | Fighter-level evidence | Deck result | Phase 4B status |
| --- | --- | --- | --- |
| Bruce Lee | Official 2019→2025 lineage + published current deck corpus | 17 unique definitions, quantities sum 30 | **verified** |
| Muhammad Ali | Current stats/stance ability indexed | Complete published per-card corpus unavailable | **blocked** |
| Blackbeard | Current stats/doubloon/ransom behavior indexed | Complete published per-card corpus unavailable | **blocked** |
| Chupacabra | Current stats/ability indexed | Secondary index can list names/quantities, but eligible complete type/value/BOOST/effect corpus unavailable | **blocked** |
| Loki | Current stats/TRICK global rules indexed | Complete eligible card metadata/effect corpus unavailable | **blocked** |
| Pandora | Current stats + seven Misery external definitions indexed | Current public UmDb action-card data is visibly partial | **blocked** |
| Leonardo | Competitive hero stats/ability indexed | Official set exposes aggregate 120 hero cards only; complete Leonardo corpus unavailable | **blocked** |
| Donatello | Competitive hero stats/ability indexed | Invention/action-card corpus unavailable | **blocked** |
| Michelangelo | Competitive hero stats/ability and hand-size rule indexed | Complete action-card corpus unavailable | **blocked** |
| Raphael | Competitive hero stats/ability indexed | Complete action-card corpus unavailable | **blocked** |
| Rosie the Riveter | Official physical inventory + current stats/upgrade rule indexed | Official 30-card allocation known; per-card transcript unavailable | **blocked** |
| John Henry | Official physical inventory + current stats/track rule indexed | Official 30-card allocation known; per-card transcript unavailable | **blocked** |
| Wyatt Earp | Official physical inventory + current stats/Showdown rule indexed | Official 30-card allocation known; per-card transcript unavailable | **blocked** |
| George Washington | Official physical inventory + current stats/Espionage rule indexed | Official 30-card allocation known; current normalized database incomplete | **blocked** |
| Shredder | Official product proves separate competitive Hero Deck; current competitive stats/Foot Soldier ability indexed | Complete Hero Deck transcript unavailable | **blocked** |
| Krang | Official product proves separate competitive Hero Deck; current competitive stats/Doomsday ability indexed | Complete Hero Deck transcript unavailable | **blocked** |

### Quantity validation

- `bruce-lee`: fixed deck 30; 17 unique card definitions; quantity sum **30** — PASS.
- Stars & Stripes: official rulebook/physical component inventory establishes **30 action cards for each fighter**, but per-card quantities/metadata cannot be reconciled; each deck remains BLOCKED.
- TMNT base product: official product establishes **120 hero action cards across the four Turtle heroes**. The worker does not infer per-fighter allocation or per-card facts from the aggregate total.
- Battle of Legends Volume Three, Muhammad Ali, Shredder, and Krang: no per-fighter deck count is promoted into the manifest where the checked eligible corpus did not independently expose it.
- `quantity_sum: null` in a blocked deck means **unknown/not validated**, not zero cards.

### Source-discipline decisions

1. Official Restoration product/rulebook/addendum material is retained as highest available authority.
2. Only `unmatched.cards/umdb/...` is eligible as a published normalized database.
3. `unmatched.cards/decks/...` fan/community decks were deliberately excluded.
4. Community/fan transcriptions discovered for recent decks were not used to populate official card records unless independently verified by an eligible higher-ranked source.
5. The Unmatched Club is used as a current rules/roster/ruling discovery index, with explicit lower authority.
6. Missing action-card metadata is represented as `blocked`, never reconstructed from previews, balance patches, or guesses.

### Bruce Lee lineage

There is exactly one `bruce-lee` fighter manifest. Its `set_ids` preserve both:
- `bruce-lee-solo` (2019 standalone lineage);
- `bruce-lee-vs-muhammad-ali` (2025 set membership).

Official Restoration history says the returning Bruce Lee deck is the out-of-print deck from the 2019 solo pack. No duplicate fighter/deck identity is created.

### TMNT scope boundary

For Leonardo, Donatello, Michelangelo, Raphael, Shredder, and Krang, manifests contain only competitive hero behavior. The following are deliberately excluded:
- Shredder/Krang Adventures villain decks;
- villain/minion initiative behavior;
- threat/scenario state;
- scenario objectives and enemy AI;
- any attempt to substitute Adventures enemy cards for the 2026 Shredder/Krang Hero Decks.

### Schema/effect extension proposals for orchestrator

These are **proposals only**. No shared semantic file is modified.

1. **Attack targeting/range modifier**
   - Affected: Muhammad Ali (Float stance), Shredder (Foot Soldier attack network).
   - Gap: current effect vocabulary changes combat value but has no generic way to change who a melee/ranged fighter may legally target based on distance or a path token.
   - Proposal: generic attack-eligibility/targeting modifier evaluated before combat-card commitment.
   - Blocking: fighter-level behavior can be recorded with a source-defined composite; generic integration should occur before engine implementation.

2. **Movement-value modifier**
   - Affected: Loki, Krang.
   - Gap: `ADD_VALUE` is a combat-value operation; these abilities alter maneuver movement value.
   - Proposal: generic movement-value modifier/layer with source lifetime and recomputation.
   - Blocking: engine-level implementation needs integration.

3. **Controlled repeat/stop loop**
   - Affected: Pandora's Box.
   - Gap: reveal/resolve another Misery repeatedly, controller may stop between iterations, but a feather threshold forces termination and damage.
   - Proposal: generic resumable repeated-resolution composite with optional stop predicate and forced-stop predicate.
   - Blocking: Box can remain a source-defined composite until promoted.

4. **Starting/max hand-size rules**
   - Affected: Michelangelo.
   - Gap: fighter schema has no explicit `starting_hand_size` / `maximum_hand_size`.
   - Proposal: generic fighter-level hand-size rule fields, with setup draw and later hand-limit consumers reading the same authoritative values.
   - Blocking: schema integration required for complete normalization.

5. **Positioned resource tokens / movement-space cost**
   - Affected: John Henry.
   - Gap: track tokens live on spaces and alter movement counting without being fighters.
   - Proposal: generic positioned-resource instance + movement path-cost override.
   - Blocking: engine integration required.

6. **Restricted free action**
   - Affected: Wyatt Earp.
   - Gap: `GAIN_ACTION` grants a generic action; Showdown grants an action usable only to attack.
   - Proposal: action credit with allowed action-kind constraint.
   - Blocking: source-defined Showdown composite used in the fighter manifest.

7. **Pre-defense marker attachment/veto**
   - Affected: George Washington.
   - Gap: ruse is committed after attack declaration but before defender decides whether to defend, is attached to the unresolved attack, and can be removed by a random discard.
   - Proposal: explicit pre-defense combat hook plus temporary marker/resource attachment to combat/card context.
   - Blocking: source-defined Espionage composite used pending integration.

8. **Path/edge-resident tokens**
   - Affected: Shredder.
   - Gap: Foot Soldiers occupy paths (edges), while current persistent-state references provide spaces but no `path_ref`/edge anchor.
   - Proposal: `path_ref` + positioned token instances on battlefield edges; removal when traversed can then be ordinary resource-state mutation.
   - Blocking: Shredder engine implementation requires integration.

9. **Random die/table roll with reroll authority**
   - Affected: Krang.
   - Gap: competitive ability references the physical Die of Ultimate Destruction and permits spending/deactivating a machine to reroll.
   - Proposal: generic random-table/die resolution composite with captured result and reroll/replace authority.
   - Blocking: do not import Adventures villain die-result logic as a substitute for competitive card semantics.

10. **Third-party ransom/payment window**
    - Affected: Blackbeard.
    - Gap: any opponent may transfer source-defined Treasury doubloons to Blackbeard to ignore a marked card effect.
    - Proposal: generic eligible-payer choice window + resource transfer + scoped effect prevention.
    - Blocking: exact ransom annotations remain deck-corpus blocked.

### Ambiguities and source gaps

- **Muhammad Ali:** current normalized public database is incomplete; current rules index exposes the stance ability but no hero cards.
- **Battle of Legends Volume Three:** current public normalized/index sources expose substantial fighter mechanics, but not a complete eligible action-card transcript for all four fighters. Chupacabra/Loki card-name lists discovered in lower-tier indexes are insufficient to establish all printed values, BOOST values, types, and effects.
- **Pandora:** seven Misery definitions are preserved as external definitions. The current public action-card corpus is partial, so action cards are not guessed.
- **TMNT turtles:** official product confirms competitive compatibility and 120 aggregate hero action cards. Current public rules index reports no hero cards for each Turtle, so no community transcript is promoted.
- **Stars & Stripes:** official sources establish release, physical components and 30-card fighter allocations. The public normalized database is incomplete and current rules index reports no hero cards. No preview/community deck is used to fill the gap.
- **Shredder/Krang:** official product explicitly supplies competitive Hero Decks, but checked eligible sources do not expose their complete card corpus. Adventures villain behavior is not substituted.

### Files created

Fighter manifests:
- `docs/fighters/phase-4b/bruce-lee.yaml`
- `docs/fighters/phase-4b/muhammad-ali.yaml`
- `docs/fighters/phase-4b/blackbeard.yaml`
- `docs/fighters/phase-4b/chupacabra.yaml`
- `docs/fighters/phase-4b/loki.yaml`
- `docs/fighters/phase-4b/pandora.yaml`
- `docs/fighters/phase-4b/leonardo.yaml`
- `docs/fighters/phase-4b/donatello.yaml`
- `docs/fighters/phase-4b/michelangelo.yaml`
- `docs/fighters/phase-4b/raphael.yaml`
- `docs/fighters/phase-4b/rosie-the-riveter.yaml`
- `docs/fighters/phase-4b/john-henry.yaml`
- `docs/fighters/phase-4b/wyatt-earp.yaml`
- `docs/fighters/phase-4b/george-washington.yaml`
- `docs/fighters/phase-4b/shredder.yaml`
- `docs/fighters/phase-4b/krang.yaml`

Deck manifests:
- `docs/cards/phase-4b/bruce-lee.yaml`
- `docs/cards/phase-4b/muhammad-ali.yaml`
- `docs/cards/phase-4b/blackbeard.yaml`
- `docs/cards/phase-4b/chupacabra.yaml`
- `docs/cards/phase-4b/loki.yaml`
- `docs/cards/phase-4b/pandora.yaml`
- `docs/cards/phase-4b/leonardo.yaml`
- `docs/cards/phase-4b/donatello.yaml`
- `docs/cards/phase-4b/michelangelo.yaml`
- `docs/cards/phase-4b/raphael.yaml`
- `docs/cards/phase-4b/rosie-the-riveter.yaml`
- `docs/cards/phase-4b/john-henry.yaml`
- `docs/cards/phase-4b/wyatt-earp.yaml`
- `docs/cards/phase-4b/george-washington.yaml`
- `docs/cards/phase-4b/shredder.yaml`
- `docs/cards/phase-4b/krang.yaml`

Report:
- `docs/phase-4b/worker-d-report.md`
