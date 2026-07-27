# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported externally after commit creation; a commit cannot embed its own final SHA.  
Assigned fighters: **16**  
Verified: **15** — `bruce-lee`, `muhammad-ali`, `blackbeard`, `chupacabra`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Partial: **1** — `loki`  
Blocked: **none at full fighter/deck level and no remaining Worker-D evidence gap blocks implementation**. Loki has one explicitly marked uncertain project normalization for multiplayer non-combat TRICK recipient selection.  
Verified action decks: **16/16**  
Quantity validation: **PASS — all 16 fixed decks reconcile to 30 cards.**  
Semantic/evidence audit: **PASS FOR CORPUS RECONCILIATION, NOT YET GLOBAL EXECUTION-READY** — confirmed Worker-D-owned transcription/normalization defects were corrected; Loki's remaining source uncertainty has an explicit provisional engine rule and documented replacement point, while shared-runtime integration needs remain.  
Files in Worker D scope: **33** — 16 fighter manifests, 16 deck manifests, this report.  
Shared files changed: **none**.  
Merge to `main`: **not performed**.

### Verification matrix

| Fighter | Fighter status | Deck status | Reconciliation | Evidence / integration note |
| --- | --- | --- | --- | --- |
| Bruce Lee | verified | verified | 17 / 30 | Single canonical 2019→2025/26 lineage. |
| Muhammad Ali | verified | verified | 13 / 30 | Float targeting and temporary space-lock semantics need shared support. |
| Blackbeard | verified | verified | 12 / 30 | Any opponent may elect to pay the shared-Treasury ransom; generic multi-actor choice support remains. |
| Chupacabra | verified | verified | 11 / 30 | — |
| Loki | **partial** | verified | 11 / 30 | 1v1 is source-complete. For a Loki-played non-combat TRICK in multiplayer, the project provisionally lets Loki's controller choose the receiving eligible opponent; this chooser is explicitly `uncertain` because no authoritative ruling was found. |
| Pandora | verified | verified | 12 / 30 + 7 Miseries | Controlled reveal/stop loop needs shared support. |
| Leonardo | verified | verified | 12 / 30 | `For Sensei` resolves the two optional fighter moves sequentially in controller-chosen order as project engine normalization. |
| Donatello | verified | verified | 14 / 30 | Inventions/out-of-play are explicit first-class zones; staged choice continuation needs shared support. |
| Michelangelo | verified | verified | 12 / 30 | Hand-size, choose-two-distinct fanout and authoritative per-turn card-play history need shared support. |
| Raphael | verified | verified | 12 / 30 | Break Something cross-window state is explicit. |
| Rosie the Riveter | verified | verified | 11 / 30 | All four physical upgrade identities/order/effects are mapped; ordered-row, maneuver-only movement and runtime attack-type semantics need shared support. |
| John Henry | verified | verified | 12 / 30 | Official track placement/reuse is complete; reconnect-safe distinct-path identity needs shared representation. |
| Wyatt Earp | verified | verified | 12 / 30 | Showdown is explicitly own-turn-only; restricted immediate attack action needs shared support. |
| George Washington | verified | verified | 13 / 30 | Ruse lifecycle explicit; pre-defense attachment and staged choice continuation need shared support. |
| Shredder | verified | verified | 13 / 30 | Bebop & Rocksteady are one 7-health sidekick; Foot Soldiers use path_ref instances; removable-path lifecycle is source-complete from the official Archive ruling. |
| Krang | verified | verified | 10 / 30 | Dynamic BOOST, face-up ranged commitment, die/reroll and temporary space-lock semantics need shared support. |

### Confirmed Worker-D-owned corrections applied

- **Bruce Lee:** canonical ability identity corrected to `Fleet of Foot`; one historical standalone → later-set fighter lineage retained.
- **Muhammad Ali:** combat-damage history updates when combat damage is applied, before `AFTER COMBAT`; `Stronger Than the Skill` looks at the opponent hand unconditionally and shuffles only on a win; deck randomization uses explicit `SHUFFLE`.
- **Loki:** TRICK cleanup distinguishes Loki-played versus opponent-played Loki-owned cards; `Sindri's Bet` is `DURING COMBAT`; `Svadilfari's Lure` uses opaque card-instance selection rather than fabricated uniform randomness. Current Reference closes the `END_TURN` × TRICK cleanup question. For the remaining multiplayer non-combat recipient ambiguity, implementation now uses the explicit but uncertain project normalization: Loki's controller chooses the receiving eligible opponent.
- **Pandora:** reconciled current Hindsight discard-pile choice, Guided by the Fates healing Pandora, Forged by Hephaestus ownership scope, Spite all-adjacent targeting, and current BOOST 2 values on Spite/Malice.
- **Leonardo:** `For Sensei` no longer carries an unresolved evidence flag. Both combat fighters may move up to 1 space; the effect controller chooses sequential resolution order in the project engine.
- **Rosie the Riveter:** official set rules establish `Merlin Engine → Cavity Magnetron → Sedgley Fist Gun → Whizbang`. Physical-component transcription supplied by the project owner establishes: Cavity Magnetron = optional draw after Rosie attacks; Sedgley Fist Gun = +1 to Rosie's attacks; Whizbang = Rosie becomes ranged. Merlin Engine is maneuver-only optional +1 space.
- **John Henry:** when all ten track tokens are deployed, the start-turn summon relocates an existing deployed token; opponent-controlled movement does not receive the track benefit.
- **Wyatt Earp:** Showdown is restricted to attacks by Wyatt's owner during that owner's turn.
- **George Washington:** ruse tokens carry ready/attached/used lifecycle semantics; `Sympathetic Stain` retains target selection and conditional hand operation in one scheme-resolution context.
- **Donatello:** `inventions` and `out-of-play` declare ownership, visibility/use semantics and source lifetime explicitly.
- **Raphael:** `Break Something` cancellation result persists through the current combat; `Relentless` uses explicit `SHUFFLE`.
- **Blackbeard:** doubloon movement is explicit transfer; payer identity is not missing resource state because multiplayer ransom comes from the shared Treasury.
- **Michelangelo:** discard-to-deck recycling uses explicit `SHUFFLE`; Cowabunga history is deferred to shared event-history semantics rather than reconstructed from zones.
- **Shredder:** `Bebop & Rocksteady` corrected to one sidekick token / one 7-health fighter. The official Archive ruling for disappearing paths is transcribed: a Foot Soldier stays deployed when its path disappears; while absent it does not count for Shredder adjacency effects but still counts as deployed for `Perplexing Tactics`; if the Heorot path returns, normal semantics resume; Point Pleasant follows the same absence behavior but its removed path never returns.
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
9. Path/edge-resident token instances with `path_ref`, dormant-while-path-absent state, and path-return lifecycle — Shredder.
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
22. Multiplayer recipient choice for a destination expressed only as `an opponent's hand` — provisional Loki normalization currently assigns this choice to Loki's controller; authority remains uncertain and must be replaceable if an official ruling appears.

### Remaining evidence uncertainty — non-blocking project normalization

#### Loki multiplayer non-combat TRICK recipient

`MISCHIEF-MONGER` says a Loki-played TRICK goes to `your opponent's hand`. A targeted web pass checked the current Loki Rules Hub, current Unmatched Reference, official Core Rules, current indexed Club rulings/disputes, indexed BGG material and current deck transcriptions. No authoritative ruling was found that selects the recipient for a Loki-played **non-combat** TRICK when multiple opponents are eligible. `TRICK: Svadilfari's Lure` played as a scheme is the concrete evidence case.

For implementation, the project owner has selected the following **uncertain provisional normalization**:

> When multiple opponents are eligible to receive a Loki-played non-combat TRICK, Loki's controller chooses the receiving opponent.

This is **not represented as an official ruling**. It is tagged `uncertain_project_normalization` in the fighter manifest and should be changed if a later authoritative ruling specifies a different chooser. Thus it no longer blocks implementation, but Loki remains `partial` at the evidence-certainty level.

The official Free-for-All rule only says that `your opponent` on an **effect on a combat card** means the other player in that combat. That scope does not establish the chooser for Loki's fighter-ability cleanup after a scheme.

Relevant sources:
- https://www.the-unmatched.club/rules/loki
- https://iellogames.com/wp-content/uploads/2024/02/UN-Adventures_Core-rules_EN_Light.pdf
- https://how-to-play.s3.us-east-2.amazonaws.com/295564/rules/295564_rules.pdf

### Resolved / reclassified items

- **Rosie mapping:** resolved by physical-component transcription supplied by the project owner; Rosie is `verified`.
- **Leonardo `For Sensei`:** resolved as project engine semantics — controller chooses the order of the two sequential optional movements.
- **Loki `END_TURN` × TRICK:** resolved. Current Unmatched Reference states that TRICK cards still use the opponent-hand cleanup destination even when an effect ends the turn, because Cleanup still occurs.
- **Loki multiplayer recipient:** no longer an implementation blocker because an explicit project normalization now exists; it remains marked uncertain/partial until authoritative evidence confirms or replaces it.
- **Shredder removable paths:** resolved by manual transcription of the current official / Archive-authority verdict from `https://www.the-unmatched.club/tools/disputes/2aa122db-53b9-4875-862e-766540449240`. This is now an engine path-lifecycle requirement, not an evidence gap.
- **Broad Rulings Archive freshness concern:** no current Worker-D blocker remains from this alone. Where a newest indexed verdict body is not machine-readable, manual authoritative transcription is acceptable evidence; future gaps should be handled per ruling rather than blocking the whole Archive.
- **Deck provenance/durable proof for newest sets:** deliberately deferred by the project owner to later evidence stages.
- **John Henry distinct paths, Michelangelo Cowabunga, global endgame and other documented engine needs:** source rules are known; these are shared-runtime integration items rather than evidence gaps.

### Source discipline

Official Restoration Games / IELLO material remains canonical for set identity and special rules. Current official rulings and current rule-change indexes supersede old generic wording. Published physical components and manually transcribed official ruling bodies are authoritative facts even where the retrieval client cannot expose their text directly. Project-selected fallbacks for unresolved cases must be explicitly tagged as uncertain and must never be presented as publisher rulings. The Unmatched Club / Unmatched Arena / current Reference are cross-check/transcription/ruling-discovery layers where first-party sources do not expose full machine-readable data. `unmatched.cards/decks/...` fan/community decks and provisional preview values are excluded.

### Scope

Worker D changed only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.
