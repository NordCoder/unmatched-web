# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: **reported in the external handoff after this report commit is created.** A Git commit cannot truthfully embed its own SHA in a file inside that same commit because changing the embedded SHA changes the commit object ID.  
Assigned fighters: **16**  
Verified fighters: `bruce-lee`, `muhammad-ali`, `blackbeard`, `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Partial fighters: `rosie-the-riveter`  
Blocked fighters: **none**  
Verified action decks: **16/16**  
Quantity validation: **PASS for all 16 action decks; every fixed action deck reconciles to 30 physical cards.**  
Files created on Worker D branch: **33** (16 fighter manifests, 16 deck manifests, this report).  
Shared schema/mechanics/rules/rulings/set/roadmap files changed: **none**.

### Final verification matrix

| Fighter | Fighter status | Deck status | Deck reconciliation | Notes |
| --- | --- | --- | --- | --- |
| Bruce Lee | verified | verified | 17 definitions / 30 cards | One canonical 2019→2025 lineage. |
| Muhammad Ali | verified | verified | 13 / 30 | Full visual corpus reconciled with official IELLO stance/set rules. |
| Blackbeard | verified | verified | 12 / 30 | Ransom semantics explicit. |
| Chupacabra | verified | verified | 11 / 30 | Historical combat-win state explicit. |
| Loki | verified | verified | 11 / 30 | TRICK ownership/hand-zone semantics retained. |
| Pandora | verified | verified | 12 / 30 + 7 Miseries | Miseries are external gameplay definitions. |
| Leonardo | verified | verified | 12 / 30 | Competitive hero behavior only. |
| Donatello | verified | verified | 14 / 30 | Inventions reconciled against official IELLO text. |
| Michelangelo | verified | verified | 12 / 30 | Starting/max hand size 3; Nunchaku state explicit. |
| Raphael | verified | verified | 12 / 30 | Anger Issues and lost-combat history explicit. |
| Rosie the Riveter | **partial** | **verified** | 11 / 30 | All four upgrade gameplay semantics identified; exact printed token labels/order remain untranscribed. |
| John Henry | verified | verified | 12 / 30 | Track/path movement and distinct-path history explicit. |
| Wyatt Earp | verified | verified | 12 / 30 | Restricted free-attack actions explicit. |
| George Washington | verified | verified | 13 / 30 | Ruse clauses and Culper Spy revival explicit. |
| Shredder | verified | verified | 13 / 30 | Competitive Hero Deck only; Foot Soldiers are path resources. |
| Krang | verified | verified | 10 / 30 | Competitive Hero Deck only; destruction-die rules preserved. |

### Quantity validation

| Fighter | Game deck | Unique action-card definitions | External definitions | Result |
| --- | ---: | ---: | ---: | --- |
| Bruce Lee | 30 | 17 | 0 | PASS |
| Muhammad Ali | 30 | 13 | 0 | PASS |
| Blackbeard | 30 | 12 | 0 | PASS |
| Chupacabra | 30 | 11 | 0 | PASS |
| Loki | 30 | 11 | 0 | PASS |
| Pandora | 30 | 12 | 7 Miseries | PASS |
| Leonardo | 30 | 12 | 0 | PASS |
| Donatello | 30 | 14 | 0 | PASS |
| Michelangelo | 30 | 12 | 0 | PASS |
| Raphael | 30 | 12 | 0 | PASS |
| Rosie the Riveter | 30 | 11 | 4 fighter-level upgrades | PASS |
| John Henry | 30 | 12 | 0 | PASS |
| Wyatt Earp | 30 | 12 | 0 | PASS |
| George Washington | 30 | 13 | 0 | PASS |
| Shredder | 30 | 13 | 0 | PASS |
| Krang | 30 | 10 | destruction-die rule | PASS |

### Rosie upgrade reconciliation

The official Restoration Games Stars & Stripes material establishes four upgrade tokens, all inactive at setup, and states that active upgrades grant Rosie additional abilities. A first-party image asset for the upgrade set is published at `Rosie-Upgrades-1024x354.png`.

The four gameplay semantics are now identified and recorded in `docs/fighters/phase-4b/rosie-the-riveter.yaml`:

1. **Ranged upgrade** — Rosie becomes a ranged fighter while active.
2. **Movement upgrade** — Rosie's movement value becomes 3 instead of 2 while active.
3. **Attack upgrade** — Rosie's attacks receive +1 value while active.
4. **Draw upgrade** — after Rosie attacks, she may draw 1 card while active.

Evidence stack:

- official Stars & Stripes rules/components establish the four-upgrade system;
- official Restoration gameplay article establishes active-upgrade additional abilities;
- official `Rosie-Upgrades-1024x354.png` asset establishes a first-party visual reference exists;
- independent physical-owner reports corroborate all four semantics, with a second owner report independently corroborating the draw-after-attack upgrade.

Rosie remains `partial`, not `verified`, for one narrowly defined reason: **the exact printed token labels and physical order/identity mapping have not been directly transcribed from the first-party upgrade image**. The manifest therefore uses descriptive semantic ids and does not claim them as printed names.

This is no longer a gameplay-semantics gap. It is a printed-identity/order verification gap.

### Final visual-corpus reconciliation

The final research pass used the current The Unmatched Club card-image corpus as a physical-card transcription layer for newest decks whose machine-readable databases were incomplete. Supplied images exposed printed names, quantities, fighter restrictions, card types, printed values, BOOSTs and exact effects.

The visual transcription is not standalone product authority. Canonical fighter/set identity and special-rule interpretation remain anchored to official Restoration Games / IELLO material.

No `unmatched.cards/decks/...` fan/community deck is used to populate an official manifest.

### Schema/effect-extension proposals for orchestrator

Proposals only; Worker D does not modify shared semantics.

1. Attack targeting/range modifier — Muhammad Ali / Shredder.
2. Maneuver movement-value modifier — Loki / Krang / Rosie.
3. Controlled repeat/stop resolution loop — Pandora.
4. Starting/max hand-size rules — Michelangelo.
5. Positioned resource/path-cost semantics — John Henry.
6. Restricted free action — Wyatt Earp / Rosie action cards.
7. Pre-defense marker attachment — George Washington.
8. Path/edge-resident tokens — Shredder.
9. Random die/table resolution with reroll authority — Krang.
10. Third-party ransom/payment window — Blackbeard.
11. Ordered/multi-fighter reposition semantics — Leonardo.
12. Fighter attack-type modifier — Rosie ranged upgrade.

### Scope

Worker D writes only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.

No shared semantic/control file is modified. No merge to `main` is performed by this worker.
