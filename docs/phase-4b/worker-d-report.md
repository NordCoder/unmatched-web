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
Schema-extension proposals: **12**, listed below; no shared semantic file changed.  
New ambiguity/blockers: **no gameplay blocker**; Rosie has one printed-identity gap only.  
Source gaps: Rosie upgrade effects are known, but exact printed token labels/order were not directly transcribed from the first-party upgrade image.  
Files created: **33** — 16 fighter manifests, 16 deck manifests, this report.  
Shared files changed: **none**.  
Merge to `main`: **not performed**.

### Verification matrix

| Fighter | Fighter status | Deck status | Reconciliation | Remaining note |
| --- | --- | --- | --- | --- |
| Bruce Lee | verified | verified | 17 / 30 | Single canonical 2019→2025 lineage. |
| Muhammad Ali | verified | verified | 13 / 30 | — |
| Blackbeard | verified | verified | 12 / 30 | Ransom requires generic third-party payment semantics. |
| Chupacabra | verified | verified | 11 / 30 | — |
| Loki | verified | verified | 11 / 30 | TRICK ownership semantics preserved. |
| Pandora | verified | verified | 12 / 30 + 7 Miseries | Miseries are external definitions. |
| Leonardo | verified | verified | 12 / 30 | Competitive behavior only. |
| Donatello | verified | verified | 14 / 30 | Inventions modeled as a public card zone. |
| Michelangelo | verified | verified | 12 / 30 | Hand-size semantics require shared support. |
| Raphael | verified | verified | 12 / 30 | — |
| Rosie the Riveter | **partial** | verified | 11 / 30 | Gameplay complete; printed upgrade labels/order unverified. |
| John Henry | verified | verified | 12 / 30 | Track/path semantics require shared support. |
| Wyatt Earp | verified | verified | 12 / 30 | Restricted free-action semantics require shared support. |
| George Washington | verified | verified | 13 / 30 | Ruse attachment semantics require shared support. |
| Shredder | verified | verified | 13 / 30 | Competitive Hero Deck only; Foot Soldiers are path resources. |
| Krang | verified | verified | 10 / 30 | Competitive Hero Deck only; destruction-die semantics preserved. |

### Rosie residual gap

All four upgrade gameplay effects are recorded in `docs/fighters/phase-4b/rosie-the-riveter.yaml`:

- ranged attack type;
- movement value 3;
- +1 to Rosie's attacks;
- optional draw 1 after Rosie attacks.

`rosie-the-riveter` remains `partial` only because the exact printed token labels and physical order/identity mapping were not directly transcribed. Semantic IDs in the manifest are descriptive and are not claimed printed names.

### Schema-extension proposals

1. Attack targeting/range modifier — Muhammad Ali / Shredder.
2. Maneuver movement-value modifier — Loki / Krang / Rosie.
3. Controlled repeat/stop resolution loop — Pandora.
4. Starting/max hand-size authority — Michelangelo.
5. Positioned resource/path-cost semantics — John Henry.
6. Restricted free-action credit — Wyatt Earp / Rosie card effects.
7. Pre-defense marker attachment — George Washington.
8. Path/edge-resident tokens — Shredder.
9. Random die/table resolution with paid reroll authority — Krang.
10. Third-party ransom/payment window — Blackbeard.
11. Ordered multi-fighter reposition semantics — Leonardo.
12. Fighter attack-type modifier — Rosie.

### Source discipline

Official Restoration Games / IELLO material remains canonical for set identity and special rules. The Unmatched Club card-image corpus is used only as a transcription layer where official sources do not expose complete machine-readable action-card text. `unmatched.cards/decks/...` fan/community decks are excluded.

### Scope

Worker D changed only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.
