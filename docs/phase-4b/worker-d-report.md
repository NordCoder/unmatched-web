# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported externally after commit creation; a commit cannot embed its own final SHA.  
Assigned fighters: **16**  
Verified: **14** — `bruce-lee`, `muhammad-ali`, `blackbeard`, `chupacabra`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Partial: **2** — `loki`, `rosie-the-riveter`  
Blocked: **none at full fighter/deck level**; explicit evidence and integration blockers are listed below.  
Verified action decks: **16/16**  
Quantity validation: **PASS — all 16 fixed decks reconcile to 30 cards.**  
Semantic/evidence audit: **PASS FOR CORPUS RECONCILIATION, NOT YET GLOBAL EXECUTION-READY** — all confirmed Worker-D-owned transcription/normalization defects found in this pass were corrected; unresolved source/ruling gaps and shared-model blockers are explicit below.  
Files in Worker D scope: **33** — 16 fighter manifests, 16 deck manifests, this report.  
Shared files changed: **none**.  
Merge to `main`: **not performed**.

### Verification matrix

| Fighter | Fighter status | Deck status | Reconciliation | Evidence / integration note |
| --- | --- | --- | --- | --- |
| Bruce Lee | verified | verified | 17 / 30 | Single canonical 2019→2025/26 lineage; current deck facts retain the historical deck behavior. |
| Muhammad Ali | verified | verified | 13 / 30 | Float targeting and temporary space-lock semantics need shared support. |
| Blackbeard | verified | verified | 12 / 30 | Multiplayer source evidence is sufficient: any opponent may elect to pay the shared-Treasury ransom; generic multi-actor choice support remains. |
| Chupacabra | verified | verified | 11 / 30 | — |
| Loki | **partial** | verified | 11 / 30 | 1v1 behavior is source-complete. In 2v2/FFA, no authoritative rule was found selecting which opponent receives a Loki-played non-combat TRICK. |
| Pandora | verified | verified | 12 / 30 + 7 Miseries | Current corpus rechecked; controlled reveal/stop loop still needs shared support. |
| Leonardo | verified | verified | 12 / 30 | `For Sensei` source wording is known, but mixed-ownership movement ordering authority is not explicitly sourced; no own-fighter-first guess remains in the manifest. |
| Donatello | verified | verified | 14 / 30 | Inventions/out-of-play are explicit first-class zones; staged choice continuation needs shared support. |
| Michelangelo | verified | verified | 12 / 30 | Hand-size, choose-two-distinct fanout and authoritative per-turn card-play history need shared support. |
| Raphael | verified | verified | 12 / 30 | Break Something cross-window state is explicit. |
| Rosie the Riveter | **partial** | verified | 11 / 30 | Official rules establish token names/order and Merlin Engine; mapping the other three printed tokens to their three known effects remains unverified. |
| John Henry | verified | verified | 12 / 30 | Official track placement/reuse is complete; reconnect-safe distinct-path identity needs shared representation. |
| Wyatt Earp | verified | verified | 12 / 30 | Showdown is explicitly own-turn-only; restricted immediate attack action needs shared support. |
| George Washington | verified | verified | 13 / 30 | Ruse ready/used lifecycle explicit; pre-defense attachment and staged choice continuation need shared support. |
| Shredder | verified | verified | 13 / 30 | Bebop & Rocksteady corrected to one 7-health sidekick; Foot Soldiers use path_ref instances; current removable-path ruling body still needs direct authoritative transcription for battlefield integration. |
| Krang | verified | verified | 10 / 30 | Dynamic BOOST, face-up ranged commitment, die/reroll and temporary space-lock semantics need shared support. |

### Confirmed Worker-D-owned corrections applied

- **Bruce Lee:** corrected the canonical ability identity to `Fleet of Foot`; retained one historical standalone → later-set fighter lineage and the current 14-health published corpus.
- **Muhammad Ali:** combat-damage history updates when combat damage is applied, before `AFTER COMBAT`; `Stronger Than the Skill` looks at the opponent hand unconditionally and shuffles only on a win; deck randomization uses explicit `SHUFFLE`.
- **Loki:** TRICK cleanup distinguishes Loki-played versus opponent-played Loki-owned cards; `Sindri's Bet` is `DURING COMBAT`; `Svadilfari's Lure` no longer fabricates a uniformly random discard and instead records Loki-owned opaque hand-card selection followed by reveal/discard.
- **Pandora:** reconciled current Hindsight discard-pile choice, Guided by the Fates healing Pandora, Forged by Hephaestus ownership scope, Spite all-adjacent targeting, and current BOOST 2 values on Spite/Malice. The stale/provisional values are not retained.
- **Leonardo:** removed an unsupported own-combat-fighter-first ordering from `For Sensei`; source wording is preserved as an explicit mixed-fighter movement composite pending ordering authority.
- **Rosie the Riveter:** official set rules close the old token identity/order gap: `Merlin Engine → Cavity Magnetron → Sedgley Fist Gun → Whizbang`. Merlin Engine is correctly maneuver-only optional +1 space, not a permanent move-value-3 modifier. The remaining three effect-to-token mappings are left unresolved rather than guessed.
- **John Henry:** official all-ten-deployed behavior is explicit: start-turn track summoning relocates an existing deployed token when none remain off-board; opponent-controlled movement still does not receive the track benefit.
- **Wyatt Earp:** Showdown is now explicitly restricted to attacks by Wyatt's owner during that owner's turn, matching the official set rules; all three once-per-turn branches remain explicit.
- **George Washington:** ruse tokens carry ready/attached/used lifecycle semantics; `Sympathetic Stain` retains target selection and conditional hand operation in one scheme-resolution context.
- **Donatello:** `inventions` and `out-of-play` declare ownership, visibility/use semantics and source lifetime explicitly.
- **Raphael:** `Break Something` cancellation result is persisted through the current combat; `Relentless` uses explicit `SHUFFLE`.
- **Blackbeard:** doubloon movement is an explicit transfer composite; current multiplayer ruling evidence means payer identity is not a missing gameplay fact because ransom comes from the shared Treasury.
- **Michelangelo:** discard-to-deck recycling uses explicit `SHUFFLE`; Cowabunga's turn-history dependency is explicitly deferred rather than reconstructed from current zones.
- **Shredder:** corrected a P0 topology error: `Bebop & Rocksteady` are one sidekick token / one 7-health fighter, not two independent 7-health fighters. Exclusive cards and `Back to Work!` now target that combined fighter. Foot Soldiers remain concrete `path_ref` token instances.
- **Krang:** `Pan-Dimensional Portal` scopes its both-player operation to the two combat players; dynamic BOOST and face-up attack fields remain explicit provisional schema extensions.

### Schema / generic-semantics proposals for the orchestrator

1. Attack targeting/range modifier — Muhammad Ali / Shredder.
2. Maneuver movement-value modifier — Loki / Krang.
3. Maneuver-only additional-space modifier — Rosie `Merlin Engine`; must not leak into effect/forced movement.
4. Controlled repeat/stop resolution loop — Pandora.
5. Starting/max hand-size authority — Michelangelo.
6. Positioned token anchors plus persisted `path_ref` / `path_ref_set` traversal identity — John Henry.
7. Restricted immediate free-action credit — Wyatt Earp / John Henry / Rosie / George Washington card effects.
8. Pre-defense combat marker attachment and token lifecycle — George Washington ruses.
9. Path/edge-resident token instances with `path_ref` and map-path lifecycle — Shredder.
10. Random die/table resolution with paid reroll authority — Krang.
11. Generic resource transfer plus multi-actor `any opponent may pay` choice aggregation — Blackbeard.
12. Mixed-ownership multi-fighter movement/reposition ordering authority — Leonardo.
13. Fighter attack-type modifier — one unresolved Rosie upgrade.
14. Dynamic/non-numeric BOOST value — Krang.
15. Face-up combat commitment / alternate ranged attack mode — Krang Missiles.
16. Ordered token-row runtime state — Rosie upgrades.
17. Authoritative per-turn card-play history — Michelangelo Cowabunga.
18. Delayed end-of-turn obligation tied to another player — Shredder `All According to Plan`.
19. Temporary cannot-leave-space lock with turn-scoped lifetime — Muhammad Ali / Krang.
20. Resumable staged effect continuation where a later choice domain depends on an earlier operation/knowledge result — Donatello, Wyatt Earp, George Washington.
21. Generic choose-N-distinct / alternating-controller / multi-actor choice orchestration — Michelangelo `Shell Insertion`, Wyatt Earp multi-choice cards, Blackbeard ransom.
22. Partial-information hand-card selection where the chooser can select a card instance without receiving its face identity — Loki `Svadilfari's Lure`.

### Evidence gaps / blockers

These are intentionally **not guessed**:

#### Fighter-level source gaps

- **Rosie — partial:** exact initial order is no longer missing. The official rulebook directly establishes `Merlin Engine`, `Cavity Magnetron`, `Sedgley Fist Gun`, `Whizbang` and directly transcribes Merlin Engine. The three remaining known effects are ranged attack type, +1 to Rosie's attacks, and optional draw 1 after Rosie attacks, but current first-party text evidence available in this pass does not establish which of the last three printed tokens grants which effect. This mapping is deterministic gameplay data, so Rosie remains `partial`.
- **Loki multiplayer — partial:** when Loki plays a TRICK in 2v2/FFA, the published ability sends it to an opponent's hand. No authoritative ruling was found that selects the recipient for a Loki-played **non-combat** TRICK when more than one opponent is eligible. 1v1 is unaffected.

#### Card/ruling execution ambiguity

- **Leonardo `For Sensei`:** the card says to move both combat fighters up to 1 space each. Core maneuver rules establish chosen ordering for a player's own multiple fighters, but Worker D found no authoritative text that explicitly supplies ordering authority for this mixed-ownership effect. The manifest therefore preserves one source-defined composite instead of imposing own-fighter-first sequencing.
- **Shredder removable paths:** the current rulings index contains an official Archive entry for Foot Soldiers on Heorot / Point Pleasant paths that can leave the gamestate. The index exposes the existence and authority of the ruling, but the exact verdict body was not available through the current indexed page. This is a battlefield/path integration evidence-retrieval blocker, not a reason to invent a fighter-local workaround.

#### Known rules, not evidence gaps

- **Blackbeard multiplayer ransom:** current ruling evidence says **any opponent may choose to pay**. Because payment transfers shared Treasury doubloons to Blackbeard, the identity of the opponent who elects payment is not persistent resource state. Remaining work is generic multi-actor choice orchestration, not missing gameplay evidence.
- **Loki `END_TURN` × TRICK cleanup:** the current rulings index marks an official Archive ruling for this interaction. It should be incorporated in global `END_TURN`/Cleanup reconciliation rather than treated as an unresolved Loki rule.
- **John Henry distinct paths:** the source rule is known; the blocker is reconnect-safe `path_ref_set` representation, not missing evidence.
- **Michelangelo Cowabunga:** the card dependency is known; the blocker is the engine's canonical `card played this turn` event/history definition.
- **Global endgame:** the current rules change delays winner checking until the end of `AFTER COMBAT`; shared core integration must reconcile this before the corpus is executable.

### Source freshness / provenance gap

- **Current Rulings Archive access is internally inconsistent:** the raw Google Doc currently obtainable through the connected Drive export identifies itself as `Last Updated: 6th Feb, 2026`, while the current Club rulings index exposes newer entries marked as official Archive rulings. For freshness-sensitive integration, existence/authority can be discovered from the current index, but exact newer ruling bodies must be obtained from a current authoritative export before encoding them as normative shared rules.
- **Stars & Stripes action cards:** official Restoration/IELLO rules and component pages are first-party for set/fighter mechanics and quantities. Complete action-card text was transcribed from physical/card-image evidence during reconciliation, but the current Club hero pages for Rosie, John Henry, Wyatt Earp and George Washington now expose an empty `Hero deck`, and the currently indexed UmDb coverage is incomplete. The retained URLs therefore do not presently let a third party reproduce the full card-text audit without obtaining a fresh physical/component scan or another complete current transcription source. This is a provenance/reproducibility gap, not evidence that the published physical decks do not exist; no `/decks/...` fan preview values are accepted as substitutes.
- **Shredder/Krang competitive Hero Deck action cards:** Restoration currently confirms the official competitive Hero Deck product, and the current Club pages confirm fighter/ability/component facts, but both Club hero pages presently expose `There are no cards yet`. The full 60-card physical corpus was transcribed during earlier reconciliation, yet the current retained public URLs do not independently reproduce all action-card definitions. Before the global executable-fixture gate, the orchestrator should retain a durable physical/component source or re-verify both 30-card decks against a complete current transcription. This is a source-durability blocker, not a basis for inventing or replacing cards with Adventures villain logic.

### Source discipline

Official Restoration Games / IELLO material remains canonical for set identity and special rules. Current official rulings and current rule-change indexes supersede old generic wording. Published physical components are authoritative facts even when a secondary index is used only as the transcription surface. The Unmatched Club / Unmatched Arena / similar indexes are cross-check/transcription layers where first-party sources do not expose full machine-readable card text. `unmatched.cards/decks/...` fan/community decks and provisional preview values are excluded.

### Scope

Worker D changed only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.