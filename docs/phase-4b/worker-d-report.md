# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: reported externally after commit creation; a commit cannot embed its own final SHA.  
Assigned fighters: **16**  
Verified: **15** — `bruce-lee`, `muhammad-ali`, `blackbeard`, `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Partial: **1** — `rosie-the-riveter`  
Blocked: **none**  
Verified action decks: **16/16**  
Quantity validation: **PASS — all 16 fixed decks reconcile to 30 cards.**  
Semantic audit: **PASS WITH DEFERRED INTEGRATION ITEMS** — known transcription/normalization bugs found in audit were corrected; remaining items are explicit below.  
Source gap: Rosie's ordered upgrade procedure and four gameplay effects are known; the initial physical token identity/order mapping is not directly transcribed.  
Files in Worker D scope: **33** — 16 fighter manifests, 16 deck manifests, this report.  
Shared files changed: **none**.  
Merge to `main`: **not performed**.

### Verification matrix

| Fighter | Fighter status | Deck status | Reconciliation | Integration note |
| --- | --- | --- | --- | --- |
| Bruce Lee | verified | verified | 17 / 30 | Single canonical 2019→2025 lineage. |
| Muhammad Ali | verified | verified | 13 / 30 | Float targeting and temporary space-lock semantics need shared support. |
| Blackbeard | verified | verified | 12 / 30 | Resource transfer and multiplayer ransom resolution need shared support. |
| Chupacabra | verified | verified | 11 / 30 | — |
| Loki | verified | verified | 11 / 30 | Foreign-hand ownership and movement modifier need shared support; multiplayer TRICK destination remains a ruling question. |
| Pandora | verified | verified | 12 / 30 + 7 Miseries | Controlled reveal/stop loop needs shared support. |
| Leonardo | verified | verified | 12 / 30 | Sequential multi-fighter reposition needs shared support. |
| Donatello | verified | verified | 14 / 30 | Inventions/out-of-play are explicit first-class zones. |
| Michelangelo | verified | verified | 12 / 30 | Hand-size and authoritative per-turn card-play history need shared support. |
| Raphael | verified | verified | 12 / 30 | Break Something cross-window state is now explicit. |
| Rosie the Riveter | **partial** | verified | 11 / 30 | Ordered-row procedure known; initial token identity/order mapping unresolved. |
| John Henry | verified | verified | 12 / 30 | Track movement is controller-sensitive; persisted distinct-path identity needs shared representation. |
| Wyatt Earp | verified | verified | 12 / 30 | Showdown state transitions explicit; restricted immediate attack action needs shared support. |
| George Washington | verified | verified | 13 / 30 | Ruse ready/used lifecycle explicit; pre-defense attachment marker needs shared support. |
| Shredder | verified | verified | 13 / 30 | Foot Soldiers are concrete path_ref token instances; battlefield path lifecycle integration remains. |
| Krang | verified | verified | 10 / 30 | Dynamic BOOST, face-up ranged commitment, die/reroll and temporary space-lock semantics need shared support. |

### Semantic audit corrections applied

- **Muhammad Ali:** combat-damage history now updates when combat damage is applied, before `AFTER COMBAT`; `Stronger Than the Skill` looks at the opponent hand unconditionally and shuffles only on a win; deck randomization uses explicit `SHUFFLE`.
- **Loki:** TRICK cleanup now distinguishes Loki-played versus opponent-played Loki-owned cards; `Sindri's Bet` is `DURING COMBAT` and an accepted discard choice performs an actual `DISCARD` operation.
- **Rosie the Riveter:** ordered upgrade-row semantics are preserved as source-defined rules; `D-Day` now scales damage/recovery by the number of upgrades already active before all upgrades activate.
- **John Henry:** track-token movement discount now applies only when the movement is controlled by John Henry's owner, matching the current ruling.
- **Wyatt Earp:** all three Showdown branches and their once-per-turn state transitions are explicit rather than hidden in a fighter-specific handler.
- **George Washington:** ruse tokens now carry ready/attached/used lifecycle semantics; `Sympathetic Stain` keeps target selection and conditional hand operation in one scheme-resolution context.
- **Donatello:** `inventions` and `out-of-play` declare ownership, visibility/use semantics and source lifetime explicitly.
- **Raphael:** `Break Something` cancellation result is persisted through the current combat; `Relentless` uses explicit `SHUFFLE`.
- **Blackbeard:** doubloon movement is an explicit transfer composite; `SPEND_RESOURCE` is no longer overloaded with an undocumented destination field.
- **Michelangelo:** discard-to-deck recycling uses explicit `SHUFFLE`; Cowabunga's turn-history dependency is explicitly deferred rather than reconstructed from current zones.
- **Shredder:** every deployed Foot Soldier is modeled as a token instance anchored to a `path_ref`, with supply/deployed lifecycle and return-to-supply removal.
- **Krang:** `Pan-Dimensional Portal` applies its both-players hand/deck operation to the two combat players, not every player in multiplayer; dynamic BOOST and face-up attack fields are explicitly identified as provisional schema extensions.

### Schema-extension proposals

1. Attack targeting/range modifier — Muhammad Ali / Shredder.
2. Maneuver movement-value modifier — Loki / Krang / Rosie.
3. Controlled repeat/stop resolution loop — Pandora.
4. Starting/max hand-size authority — Michelangelo.
5. Positioned token anchors plus persisted path/space traversal identity — John Henry.
6. Restricted immediate free-action credit — Wyatt Earp / John Henry / Rosie / George Washington card effects.
7. Pre-defense combat marker attachment and token lifecycle — George Washington ruses.
8. Path/edge-resident token instances with `path_ref` — Shredder.
9. Random die/table resolution with paid reroll authority — Krang.
10. Generic resource transfer plus multiplayer eligible-payer resolution — Blackbeard.
11. Ordered multi-fighter reposition semantics — Leonardo.
12. Fighter attack-type modifier — Rosie ranged upgrade.
13. Dynamic/non-numeric BOOST value — Krang.
14. Face-up combat commitment / alternate ranged attack mode — Krang Missiles.
15. Ordered token-row runtime state — Rosie upgrades.
16. Authoritative per-turn card-play history — Michelangelo Cowabunga.
17. Delayed end-of-turn obligation tied to another player — Shredder `All According to Plan`.
18. Temporary cannot-leave-space lock with turn-scoped lifetime — Muhammad Ali / Krang.

### Remaining ambiguity / ruling reconciliation

These are intentionally **not guessed** in Worker D manifests:

- **Rosie:** exact initial mapping of the four physical upgrade tokens to their semantic abilities and row positions. This is the only remaining fighter-level source-data gap and is why the fighter remains `partial`.
- **Loki multiplayer:** when Loki plays a TRICK in 2v2/FFA, the published wording says it enters an opponent's hand; deterministic opponent-selection authority still needs an explicit ruling/model. The current Rulings Archive also has an `END_TURN` × TRICK cleanup ruling that should be incorporated during global end-turn reconciliation.
- **Blackbeard multiplayer:** the current Rulings Archive has a dedicated 2v2/FFA ransom-payer ruling; exact payer priority/choice sequencing should be imported before multiplayer implementation.
- **John Henry:** the corpus needs a reconnect-safe identity representation for the distinct paths used this turn; an integer total is not sufficient to prove uniqueness after replay/reconnect.
- **Michelangelo:** `Cowabunga!!` needs a canonical definition of a card being "played this turn" so maneuver BOOST/discard/bottom-deck operations cannot accidentally increment the count.
- **Shredder battlefield integration:** current rulings cover Foot Soldiers on paths that can disappear/change; those rulings belong in battlefield/path integration rather than a fighter-specific workaround.
- **Global current-rules dependency:** the 2026 Rulings Archive changed endgame winner-check timing around `AFTER COMBAT`; integration should reconcile the global core model before using Worker D as an executable fixture corpus.

### Source discipline

Official Restoration Games / IELLO material remains canonical for set identity and special rules. Current official Rulings Archive entries are treated as authoritative corrections/clarifications. The Unmatched Club card-image corpus is used only as a transcription layer where official sources do not expose complete machine-readable action-card text. `unmatched.cards/decks/...` fan/community decks are excluded.

### Scope

Worker D changed only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.
