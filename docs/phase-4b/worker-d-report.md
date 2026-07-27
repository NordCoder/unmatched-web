# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported externally after commit creation; a commit cannot embed its own final SHA.  
Assigned fighters: **16**  
Verified: **15** — `bruce-lee`, `muhammad-ali`, `blackbeard`, `chupacabra`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Partial: **1** — `loki`  
Blocked: **none at full fighter/deck level**; the remaining evidence/ruling retrieval gaps and shared-model integration items are explicit below.  
Verified action decks: **16/16**  
Quantity validation: **PASS — all 16 fixed decks reconcile to 30 cards.**  
Semantic/evidence audit: **PASS FOR CORPUS RECONCILIATION, NOT YET GLOBAL EXECUTION-READY** — confirmed Worker-D-owned transcription/normalization defects were corrected; remaining blockers are now narrowly classified instead of conflating source evidence with engine support.  
Files in Worker D scope: **33** — 16 fighter manifests, 16 deck manifests, this report.  
Shared files changed: **none**.  
Merge to `main`: **not performed**.

### Verification matrix

| Fighter | Fighter status | Deck status | Reconciliation | Evidence / integration note |
| --- | --- | --- | --- | --- |
| Bruce Lee | verified | verified | 17 / 30 | Single canonical 2019→2025/26 lineage. |
| Muhammad Ali | verified | verified | 13 / 30 | Float targeting and temporary space-lock semantics need shared support. |
| Blackbeard | verified | verified | 12 / 30 | Multiplayer source evidence is sufficient: any opponent may elect to pay the shared-Treasury ransom; generic multi-actor choice support remains. |
| Chupacabra | verified | verified | 11 / 30 | — |
| Loki | **partial** | verified | 11 / 30 | 1v1 is source-complete. Targeted research found no authoritative recipient-selection rule for a Loki-played non-combat TRICK in 2v2/FFA when multiple opponents are eligible. |
| Pandora | verified | verified | 12 / 30 + 7 Miseries | Controlled reveal/stop loop needs shared support. |
| Leonardo | verified | verified | 12 / 30 | `For Sensei` resolves the two optional fighter moves sequentially in controller-chosen order as project engine normalization; no evidence blocker remains. |
| Donatello | verified | verified | 14 / 30 | Inventions/out-of-play are explicit first-class zones; staged choice continuation needs shared support. |
| Michelangelo | verified | verified | 12 / 30 | Hand-size, choose-two-distinct fanout and authoritative per-turn card-play history need shared support. |
| Raphael | verified | verified | 12 / 30 | Break Something cross-window state is explicit. |
| Rosie the Riveter | verified | verified | 11 / 30 | All four physical upgrade identities/order/effects are now mapped; ordered-row, maneuver-only movement and runtime attack-type semantics need shared support. |
| John Henry | verified | verified | 12 / 30 | Official track placement/reuse is complete; reconnect-safe distinct-path identity needs shared representation. |
| Wyatt Earp | verified | verified | 12 / 30 | Showdown is explicitly own-turn-only; restricted immediate attack action needs shared support. |
| George Washington | verified | verified | 13 / 30 | Ruse lifecycle explicit; pre-defense attachment and staged choice continuation need shared support. |
| Shredder | verified | verified | 13 / 30 | Bebop & Rocksteady are one 7-health sidekick; Foot Soldiers use path_ref instances. A newest official Archive ruling exists for paths removed from the gamestate, but its verdict body still requires manual retrieval for Heorot integration. |
| Krang | verified | verified | 10 / 30 | Dynamic BOOST, face-up ranged commitment, die/reroll and temporary space-lock semantics need shared support. |

### Confirmed Worker-D-owned corrections applied

- **Bruce Lee:** canonical ability identity corrected to `Fleet of Foot`; one historical standalone → later-set fighter lineage retained.
- **Muhammad Ali:** combat-damage history updates when combat damage is applied, before `AFTER COMBAT`; `Stronger Than the Skill` looks at the opponent hand unconditionally and shuffles only on a win; deck randomization uses explicit `SHUFFLE`.
- **Loki:** TRICK cleanup distinguishes Loki-played versus opponent-played Loki-owned cards; `Sindri's Bet` is `DURING COMBAT`; `Svadilfari's Lure` uses opaque card-instance selection rather than fabricated uniform randomness. Current Reference also closes the `END_TURN` × TRICK cleanup question.
- **Pandora:** reconciled current Hindsight discard-pile choice, Guided by the Fates healing Pandora, Forged by Hephaestus ownership scope, Spite all-adjacent targeting, and current BOOST 2 values on Spite/Malice.
- **Leonardo:** `For Sensei` no longer carries an unresolved evidence flag. Both combat fighters may move up to 1 space; the effect controller chooses the sequential resolution order in the project engine.
- **Rosie the Riveter:** official set rules establish the row `Merlin Engine → Cavity Magnetron → Sedgley Fist Gun → Whizbang`. Physical-component transcription supplied by the project owner establishes the remaining mapping: Cavity Magnetron = optional draw after Rosie attacks; Sedgley Fist Gun = +1 to Rosie's attacks; Whizbang = Rosie becomes ranged. Merlin Engine remains maneuver-only optional +1 space, not a permanent move-value increase.
- **John Henry:** when all ten track tokens are deployed, the start-turn summon relocates an existing deployed token; opponent-controlled movement still does not receive the track benefit.
- **Wyatt Earp:** Showdown is restricted to attacks by Wyatt's owner during that owner's turn.
- **George Washington:** ruse tokens carry ready/attached/used lifecycle semantics; `Sympathetic Stain` retains target selection and conditional hand operation in one scheme-resolution context.
- **Donatello:** `inventions` and `out-of-play` declare ownership, visibility/use semantics and source lifetime explicitly.
- **Raphael:** `Break Something` cancellation result persists through the current combat; `Relentless` uses explicit `SHUFFLE`.
- **Blackbeard:** doubloon movement is explicit transfer; payer identity is not missing resource state because multiplayer ransom comes from the shared Treasury.
- **Michelangelo:** discard-to-deck recycling uses explicit `SHUFFLE`; Cowabunga's history dependency is deferred to shared event-history semantics rather than reconstructed from zones.
- **Shredder:** `Bebop & Rocksteady` corrected to one sidekick token / one 7-health fighter; exclusive cards and `Back to Work!` target that combined fighter. Foot Soldiers remain concrete `path_ref` token instances.
- **Krang:** `Pan-Dimensional Portal` scopes its both-player hand/deck operation to the two combat players; dynamic BOOST and face-up attack fields remain explicit provisional schema extensions.

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
12. Fighter attack-type modifier — Rosie `Whizbang`.
13. Dynamic/non-numeric BOOST value — Krang.
14. Face-up combat commitment / alternate ranged attack mode — Krang Missiles.
15. Ordered token-row runtime state — Rosie upgrades.
16. Authoritative per-turn card-play history — Michelangelo Cowabunga.
17. Delayed end-of-turn obligation tied to another player — Shredder `All According to Plan`.
18. Temporary cannot-leave-space lock with turn-scoped lifetime — Muhammad Ali / Krang.
19. Resumable staged effect continuation where a later choice domain depends on an earlier operation/knowledge result — Donatello, Wyatt Earp, George Washington.
20. Generic choose-N-distinct / alternating-controller / multi-actor choice orchestration — Michelangelo `Shell Insertion`, Wyatt Earp multi-choice cards, Blackbeard ransom.
21. Partial-information hand-card selection where the chooser selects a card instance without receiving its face identity — Loki `Svadilfari's Lure`.

### Remaining evidence / ruling blockers

These are intentionally **not guessed**:

#### Loki multiplayer TRICK recipient — fighter remains partial

`MISCHIEF-MONGER` says a Loki-played TRICK goes to `your opponent's hand`. A targeted web pass checked the current Loki Rules Hub, current Unmatched Reference, official Core Rules, current indexed Club rulings/disputes, indexed BGG material and current deck transcriptions. No authoritative ruling was found that selects the recipient for a Loki-played **non-combat** TRICK when multiple opponents are eligible. `TRICK: Svadilfari's Lure` played as a scheme is the concrete unresolved case.

The current official Free-for-All rule only says that `your opponent` on an **effect on a combat card** means the other player in that combat. That scope does not determine the recipient of Loki's fighter-ability cleanup after a scheme.

Relevant sources:
- https://www.the-unmatched.club/rules/loki
- https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Core-rules_EN_Light.pdf
- https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf

#### Shredder Foot Soldiers on removable paths — exact verdict retrieval

The current Club index confirms a newest **official / Archive-authority** ruling named `Foot Soldiers on Heorot, Point Pleasant — How do Foot Soldiers work on paths that can be removed from the gamestate?`, but the indexed page does not expose the ruling's verdict text to this worker.

Direct ruling page for manual inspection:
- https://www.the-unmatched.club/tools/disputes/2aa122db-53b9-4875-862e-766540449240

Underlying officially endorsed Rulings Archive:
- https://docs.google.com/document/d/13b-FbPq_vuqcc3IokeHvQ2ctJaDNZZuUaZmt4uft5h0/edit

Current Unmatched Reference v10 already establishes the surrounding map facts: on Heorot a closed door makes crossing lines cease to exist, which creates the real competitive path-lifecycle question. It also says that in **competitive Point Pleasant the bridge tokens are not placed at setup**, so Point Pleasant's bridge-removal case does not currently need implementation for ordinary competitive play. The remaining competitive issue is therefore primarily Heorot unless the newest ruling says otherwise.

### Resolved / reclassified items

- **Rosie mapping:** resolved by physical-component transcription supplied by the project owner; Rosie is now `verified`.
- **Leonardo `For Sensei`:** resolved as project engine semantics — controller chooses the order of the two sequential optional movements. No separate evidence blocker remains.
- **Loki `END_TURN` × TRICK:** resolved. Current Unmatched Reference states that TRICK cards still use the opponent-hand cleanup destination even when an effect ends the turn, because Cleanup still occurs.
- **Broad Rulings Archive freshness concern:** narrowed. The current Club dispute index and current Unmatched Reference expose many recent official rulings directly. The remaining problem is **per-ruling retrieval when a newest Club dispute confirms an Archive ruling but omits its verdict body**, as currently happens for the Shredder removable-path ruling. The Google Doc itself requires a JavaScript-capable browser and is not reliably readable by the web retrieval client.
- **Deck provenance/durable proof for newest sets:** deliberately deferred by the project owner to later evidence stages. It is no longer classified as a current Worker-D blocker.
- **John Henry distinct paths, Michelangelo Cowabunga, global endgame and other documented engine needs:** source rules are known; these are shared-runtime integration items rather than evidence gaps.

### Source discipline

Official Restoration Games / IELLO material remains canonical for set identity and special rules. Current official rulings and current rule-change indexes supersede old generic wording. Published physical components are authoritative facts even when a secondary index is used only as the transcription surface. The Unmatched Club / Unmatched Arena / current Reference are cross-check/transcription/ruling-discovery layers where first-party sources do not expose full machine-readable data. `unmatched.cards/decks/...` fan/community decks and provisional preview values are excluded.

### Scope

Worker D changed only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.