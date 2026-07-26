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
| Blackbeard | verified | verified | 12 / 30 | Ransom semantics retained as explicit source-defined composite. |
| Chupacabra | verified | verified | 11 / 30 | Historical combat-win state persisted. |
| Loki | verified | verified | 11 / 30 | TRICK ownership/hand-zone semantics retained. |
| Pandora | verified | verified | 12 / 30 + 7 external Miseries | Miseries are external gameplay definitions, not action-card instances. |
| Leonardo | verified | verified | 12 / 30 | Competitive hero behavior only. |
| Donatello | verified | verified | 14 / 30 | Three Inventions reconciled against official IELLO text. |
| Michelangelo | verified | verified | 12 / 30 | Starting/max hand size 3; Nunchaku turn state explicit. |
| Raphael | verified | verified | 12 / 30 | Anger Issues plus lost-combat history explicit. |
| Rosie the Riveter | **partial** | **verified** | 11 / 30 | Action deck complete; four unique upgrade-token active abilities remain untranscribed. |
| John Henry | verified | verified | 12 / 30 | Track/path movement and distinct-path history explicit. |
| Wyatt Earp | verified | verified | 12 / 30 | Restricted immediate free-attack actions explicit. |
| George Washington | verified | verified | 13 / 30 | Ruse clauses and Culper Spy revival explicit. |
| Shredder | verified | verified | 13 / 30 | Competitive Hero Deck only; Foot Soldiers are path resources. |
| Krang | verified | verified | 10 / 30 | Competitive Hero Deck only; destruction-die rules and dynamic BOOST symbols preserved. |

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
| Rosie the Riveter | 30 | 11 | 0 action-deck externals | PASS |
| John Henry | 30 | 12 | 0 | PASS |
| Wyatt Earp | 30 | 12 | 0 | PASS |
| George Washington | 30 | 13 | 0 | PASS |
| Shredder | 30 | 13 | 0 | PASS |
| Krang | 30 | 10 | destruction-die rule | PASS |

### Final visual-corpus reconciliation

The final research pass used the current The Unmatched Club card-image corpus as a physical-card transcription source for the newest decks whose machine-readable databases were incomplete. The supplied card images exposed printed card names, quantities, fighter restrictions, card types, printed values, BOOSTs and exact effect text.

The visual transcription was not treated as standalone product authority. Canonical fighter/set identity and special-rule interpretation remain anchored to official Restoration Games / IELLO product pages, rulebooks and set rules. Quantities were reconciled back to the official fixed 30-card construction where available.

The final visual pass completed:

- Muhammad Ali — 13 definitions / 30;
- Michelangelo — 12 / 30;
- Raphael — 12 / 30;
- Rosie the Riveter action deck — 11 / 30;
- John Henry — 12 / 30;
- Wyatt Earp — 12 / 30;
- George Washington — 13 / 30;
- Shredder competitive Hero Deck — 13 / 30;
- Krang competitive Hero Deck — 10 / 30.

No `unmatched.cards/decks/...` fan/community deck was used to populate an official manifest.

### Rosie partial status

Rosie's **action deck is fully verified**: 11 unique definitions reconcile to 30 physical cards.

The fighter remains `partial` because the Stars & Stripes rules/components define four unique ordered upgrade tokens and a Mech Suit reference component. Active upgrades grant gameplay abilities in addition to the action-card effects that activate/deactivate them. The supplied screenshots covered all 30 action cards but did not expose the four individual active-side upgrade abilities.

Therefore:

- deck corpus: complete;
- upgrade activation/deactivation procedure: known;
- four upgrade identities/active abilities: still required;
- blocker severity: `partial`, not `blocked`.

The worker deliberately does **not** collapse these four identities into only `active_upgrade_count`, because that would lose required gameplay semantics.

### Source-discipline decisions

1. Official Restoration Games / IELLO product pages, rulebooks, set rules, addenda and printed components remain highest available authority.
2. The Unmatched Club visual card corpus is used as a lower-ranked transcription/reference layer where official sources establish the canonical product but do not publish the complete action-card text in machine-readable form.
3. Card-image transcription is accepted only when names/quantities reconcile to the canonical fixed deck and printed text is readable.
4. Conflicts resolve toward official rules/components and physical printed evidence.
5. `unmatched.cards/decks/...` fan/community decks remain excluded.
6. Adventures villain decks, initiative behavior, threat state and scenario AI are never substituted for Shredder/Krang competitive Hero Deck semantics.

### Important normalization results from the final pass

- **Muhammad Ali:** complete 13-card-definition corpus; Float/Bee clauses are conditional on persisted stance. `The Greatest` and `Louisville Lip` require historical same-turn combat state.
- **Michelangelo:** complete 12-definition corpus; starting/max hand size 3 remains a schema proposal; Nunchaku stacks a turn-scoped attack bonus.
- **Raphael:** complete 12-definition corpus; `Payback Time!` depends on historical combat loss this turn.
- **John Henry:** complete 12-definition corpus; HAMMER cards are tagged and several effects require the count of distinct paths used this turn.
- **Wyatt Earp:** complete 12-definition corpus; free actions to attack are restricted immediate action credits, and `You Die First` additionally binds the target.
- **George Washington:** complete 13-definition corpus; triangle clauses are modeled as effects conditional on a surviving ruse attachment. Culper Spies are summonable because `Recruit to the Ring` can return defeated Spies.
- **Shredder:** complete 13-definition competitive Hero Deck; Bebop/Rocksteady are summonable because `Back to Work!` revives defeated sidekicks. Foot Soldiers remain path/edge resources.
- **Krang:** complete 10-definition competitive Hero Deck. Each destruction-die symbol rolls independently; X is zero/no effect; active machines may be deactivated to reroll. Pan-Dimensional Portal and IQ of 968 use a dynamic destruction-die BOOST symbol rather than a numeric BOOST.

### Schema/effect-extension proposals for orchestrator

These remain proposals only. Worker D does not modify shared semantic files.

1. **Attack targeting/range modifier** — Muhammad Ali Float stance; Shredder Foot Soldier targeting network.
2. **Maneuver movement-value modifier** — Loki; Krang active Doomsday Machines.
3. **Controlled repeat/stop resolution loop** — Pandora's Box.
4. **Starting/max hand-size authority** — Michelangelo.
5. **Positioned path/edge resources + movement path-cost/history** — John Henry and Shredder.
6. **Restricted immediate free-action credits** — Wyatt Earp, John Henry, Rosie/Washington action cards.
7. **Pre-defense combat marker/attachment/veto** — George Washington ruses.
8. **Random die/table resolution + paid reroll authority** — Krang.
9. **Dynamic/non-numeric BOOST rule** — Krang cards whose BOOST is a destruction-die symbol.
10. **Third-party ransom/payment window** — Blackbeard.
11. **Ordered multi-fighter reposition** — Leonardo.
12. **Ordered unique upgrade-token row with active-source abilities** — Rosie the Riveter. This is the only remaining fighter-level data gap.

### TMNT scope boundary

For Leonardo, Donatello, Michelangelo, Raphael, Shredder and Krang, manifests model only competitive hero behavior. Explicitly excluded:

- Shredder/Krang Adventures villain decks;
- villain/minion initiative logic;
- threat/scenario state;
- scenario objectives and enemy AI;
- substitution of Adventures enemy cards for the competitive Hero Decks.

### Files and branch scope

Worker D owns exactly 33 paths:

- 16 `docs/fighters/phase-4b/<assigned-fighter-id>.yaml` manifests;
- 16 `docs/cards/phase-4b/<assigned-fighter-id>.yaml` manifests;
- `docs/phase-4b/worker-d-report.md`.

No shared schema, mechanics, rules, rulings, set-registry, roadmap or README file is modified. No merge to `main` is performed by this worker.
