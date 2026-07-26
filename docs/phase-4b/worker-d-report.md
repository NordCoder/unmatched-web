# Phase 4B Worker D report

## Worker 4B-D Handoff

Branch: `phase-4b-worker-d-latest`  
Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`  
Head: **reported in the external handoff after this report commit is created.** A commit cannot truthfully embed its own SHA in a file inside that same commit because changing the embedded SHA changes the commit object ID.  
Assigned fighters: **16**  
Verified: `bruce-lee`, `blackbeard`, `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`  
Blocked: `muhammad-ali`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`  
Quantity validation: **PASS for all 7 verified decks; construction size known but corpus incomplete for 7 additional blocked decks; Shredder/Krang remain transcript-blocked.**  
Files created on Worker D branch: **33** (16 fighter manifests, 16 deck manifests, this report).  
Shared schema/mechanics/rules/rulings/set/roadmap files changed: **none**.

### Branch/base and scope

The branch was originally created from Authorized Base `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`. Worker D writes only:

- `docs/fighters/phase-4b/<assigned-id>.yaml`;
- `docs/cards/phase-4b/<assigned-id>.yaml`;
- `docs/phase-4b/worker-d-report.md`.

No merge to `main` is performed by this worker.

### Revised verification matrix

| Fighter | Result | Quantity state | Evidence/result |
| --- | --- | --- | --- |
| Bruce Lee | **verified** | 30/30; 17 definitions | Published UmDb corpus + official Restoration 2019→2025 lineage; one canonical fighter identity. |
| Blackbeard | **verified** | 30/30; 12 definitions | Official BOL3 membership + physical deck review + complete current transcription; ransom semantics explicit. |
| Chupacabra | **verified** | 30/30; 11 definitions | Official BOL3 membership + physical/current transcript cross-check; `Tooth and Tail` spelling reconciled. |
| Loki | **verified** | 30/30; 11 definitions | Official BOL3 membership + physical/current transcript cross-check; TRICK ownership/BOOST semantics retained. |
| Pandora | **verified** | 30/30; 12 definitions + 7 external Miseries | Official BOL3 membership + physical/current cross-check; Miseries are external definitions, not action cards. |
| Leonardo | **verified** | 30/30; 12 definitions | Official IELLO TMNT rules establish stats/ability/30-card allocation and print `For Sensei`; complete transcript cross-checked against physical components. |
| Donatello | **verified** | 30/30; 14 definitions | Official IELLO TMNT rules establish 30-card allocation and authoritative Invention text; complete transcript independently cross-checked. |
| Muhammad Ali | **blocked** | construction 30/30; quantity sum unknown | Official IELLO rules establish deck size, stance mechanics and selected cards, but full per-card transcript is not yet reconciled. |
| Michelangelo | **blocked** | construction 30/30; quantity sum not promoted | Official IELLO rules establish 30 cards; current tracker reconciles names/quantities, but full type/value/BOOST/effect layer is incomplete. |
| Raphael | **blocked** | construction 30/30; quantity sum unknown | Official IELLO rules establish 30 cards and Anger Issues; complete per-card transcript remains missing. |
| Rosie the Riveter | **blocked** | construction 30/30 | Official rulebook proves 30 cards; first-party/physical card material exists, but full frame transcription and reconciliation remain pending. |
| John Henry | **blocked** | construction 30/30 | Same physical-corpus blocker; track mechanics are known independently. |
| Wyatt Earp | **blocked** | construction 30/30 | Same physical-corpus blocker; Showdown mechanics are known independently. |
| George Washington | **blocked** | construction 30/30 | Same physical-corpus blocker; ruse/Espionage mechanics are known independently. |
| Shredder | **blocked** | full deck transcript unavailable | Restoration explicitly confirms a separate competitive Hero Deck; Adventures villain deck is excluded. |
| Krang | **blocked** | full deck transcript unavailable | Restoration explicitly confirms a separate competitive Hero Deck; Adventures villain deck is excluded. |

### Quantity validation

Verified deck reconciliation:

| Fighter | Game deck | Unique action-card definitions | External definitions | Result |
| --- | ---: | ---: | ---: | --- |
| Bruce Lee | 30 | 17 | 0 | PASS |
| Blackbeard | 30 | 12 | 0 | PASS |
| Chupacabra | 30 | 11 | 0 | PASS |
| Loki | 30 | 11 | 0 | PASS |
| Pandora | 30 | 12 | 7 Miseries | PASS |
| Leonardo | 30 | 12 | 0 | PASS |
| Donatello | 30 | 14 | 0 | PASS |

For blocked decks, `quantity_sum: null` means the card-definition corpus has **not** been fully validated; it does not mean zero cards.

Official construction-size improvements from the second research pass:

- Muhammad Ali: IELLO contents page explicitly lists **30 action cards**.
- Michelangelo: IELLO TMNT contents page explicitly lists **30 action cards**.
- Raphael: IELLO TMNT contents page explicitly lists **30 action cards**.
- Rosie, John Henry, Wyatt Earp, George Washington: official Stars & Stripes contents page explicitly lists **30 action cards each**.

### Source-discipline decisions

1. Official Restoration/IELLO product pages, rulebooks, set rules, addenda and printed components remain highest authority.
2. A lower-ranked current transcript is used only where an official source proves the canonical product/fighter/deck but does not publish a machine-readable card transcript. Such use is explicitly identified and cross-checked against physical-component material or an independent current index.
3. Conflicts are resolved toward official/physical component evidence. Examples include `Electro Grenade`, `Tooth and Tail`, `Baldr's Downfall`, `Malicious Flyting`, `Freyja's Rescue`, `Svadilfari's Lure`, `Hindsight`, and `Celestial Raiments`.
4. `unmatched.cards/decks/...` fan/community decks are never used to populate the official corpus.
5. No card is reconstructed from remembered wording, speculative previews or fan balance patches.
6. A public physical reveal is evidence that a corpus exists, but does **not** by itself justify `verified` until the worker completes a reliable card-by-card transcription and reconciliation.

### Bruce Lee lineage

There is exactly one `bruce-lee` fighter manifest. Its `set_ids` preserve both `bruce-lee-solo` and `bruce-lee-vs-muhammad-ali`. Restoration's release history states that the returning deck is the original 2019 standalone deck; no duplicate fighter identity is created.

### TMNT scope boundary

For Leonardo, Donatello, Michelangelo, Raphael, Shredder and Krang, only competitive hero behavior belongs in these manifests. Explicitly excluded:

- Shredder/Krang Adventures villain decks;
- villain/minion initiative logic;
- threat/scenario state;
- scenario objectives and enemy AI;
- any substitution of Adventures enemy cards for the 2026 competitive Hero Decks.

### Schema/effect-extension proposals for orchestrator

These are proposals only; this worker does not modify shared semantics.

1. **Attack targeting/range modifier**
   - Affected: Muhammad Ali, Shredder.
   - Gap: combat value operations do not express a temporary expansion of legal attack targets/range.
   - Proposal: generic attack-eligibility/targeting modifier evaluated before combat-card commitment.

2. **Movement-value modifier**
   - Affected: Loki, Krang.
   - Gap: `ADD_VALUE` is combat-value semantics, while these mechanics modify maneuver movement.
   - Proposal: generic maneuver movement-value layer/modifier.
   - Loki corpus is verified; engine integration remains orchestrator-owned.

3. **Controlled repeat/stop resolution loop**
   - Affected: Pandora's Box.
   - Gap: repeated external-definition reveal/resolve with optional stop and forced feather threshold.
   - Proposal: resumable repeated-resolution composite with explicit continue/stop and forced-stop predicates.
   - Pandora corpus is verified using a source-defined composite pending integration.

4. **Starting/max hand-size rules**
   - Affected: Michelangelo.
   - Proposal: fighter-level `starting_hand_size` / `maximum_hand_size` authority shared by setup and cleanup.

5. **Positioned resource tokens / movement path cost**
   - Affected: John Henry.
   - Proposal: positioned-resource instances plus movement path-cost override.

6. **Restricted free action**
   - Affected: Wyatt Earp.
   - Gap: `GAIN_ACTION` is unrestricted while Showdown can grant an attack-only action.
   - Proposal: action credit carrying an allowed action-kind constraint.

7. **Pre-defense combat marker/veto**
   - Affected: George Washington.
   - Proposal: explicit pre-defense hook plus temporary marker attachment to combat/card context.

8. **Path/edge-resident tokens**
   - Affected: Shredder.
   - Proposal: `path_ref`/edge anchors for Foot Soldier tokens and source-defined traversal/removal semantics.

9. **Random die/table resolution with reroll authority**
   - Affected: Krang.
   - Proposal: generic random-table/die composite with captured result and explicit reroll/replace authority.
   - Do not substitute Adventures villain result logic unless the competitive source explicitly shares it.

10. **Third-party ransom/payment window**
    - Affected: Blackbeard.
    - Proposal: eligible-payer choice + resource transfer + scoped effect prevention.
    - Blackbeard corpus is verified using an explicit source-defined ransom composite pending integration.

11. **Ordered multi-fighter reposition composite**
    - Affected: Leonardo (`Eat, Sleep, and Breathe Ninjutsu`).
    - Proposal: orchestrator may normalize as repeated `REQUEST_CHOICE` + `PLACE` operations or promote a reusable ordered multi-target reposition composite. No new primitive is required to verify the deck.

Donatello's Inventions do not require a new card-zone primitive: public `inventions` + ordinary card ownership + source-lifetime static effects fit the Phase 4A model. The multi-player optional discard recovery on `Donatello Does Machines` can use the existing pending-choice/fanout model or remain an explicit composite.

### Remaining ambiguities / source gaps

- **Muhammad Ali:** primary rules now prove 30-card construction and selected printed cards; missing evidence is specifically the complete per-card transcript.
- **Michelangelo:** 30-card construction is official and names/quantities can be reconciled by a current tracker; full type/value/BOOST/effect metadata still needs independent verification.
- **Raphael:** 30-card construction and fighter rules are official; complete action-card transcript remains missing.
- **Stars & Stripes:** this is **not** an absent-corpus case. The official rulebook establishes 30 cards per fighter, Restoration publishes gameplay/card-layout material, and a public physical all-cards reveal exists. The blocker is manual frame/image transcription plus independent reconciliation before promotion into `cards:`.
- **Shredder/Krang:** official Restoration product explicitly establishes separate competitive Hero Decks, but the complete competitive card transcripts remain unavailable in the checked authoritative/current corpus. Adventures villain logic is not substituted.

### Key authoritative/current sources added in the second pass

- BOL3 official product: `https://restorationgames.com/shop/unmatched-battle-of-legends-vol-3/`
- BOL3 physical deck review: `https://www.youtube.com/watch?v=6fP2amMnphI`
- TMNT official IELLO set rules: `https://iellogames.com/wp-content/uploads/2026/01/UMA-TMNT_Set-rules_EN_Light.pdf`
- Bruce Lee vs Muhammad Ali official IELLO set rules: `https://iellogames.com/wp-content/uploads/2026/01/Unmatched-LvsA_Set-rules_EN_Light.pdf`
- Stars & Stripes official rules: `https://restorationgames.com/wp-content/uploads/2025/08/UM-Stars-Stripes_Rulebook-FLAT.pdf`
- Stars & Stripes official gameplay article: `https://restorationgames.com/gameplay-of-unmatched-stars-and-stripes/`
- Shredder/Krang competitive Hero Deck product: `https://restorationgames.com/shop/unmatched-adventures-teenage-mutant-ninja-turtles-shredder-krang-hero-decks/`

### Files created

The Worker D branch still owns exactly 33 created paths:

- 16 `docs/fighters/phase-4b/<assigned-fighter-id>.yaml` manifests;
- 16 `docs/cards/phase-4b/<assigned-fighter-id>.yaml` manifests;
- `docs/phase-4b/worker-d-report.md`.

The second research pass updates existing Worker D paths only; it creates no shared semantic/control file and performs no merge.
