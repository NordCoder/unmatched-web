# Phase 4A validation

**Phase:** 4A — Representative fighter/deck stress-test corpus  
**Validation date:** 2026-07-26  
**Result:** PASS  
**Overall project readiness:** NOT YET `developer-ready`; Phase 4B and battlefield/cross-validation phases remain.

## Gate under test

Phase 4A requires the representative sample to be expressible without ad-hoc undocumented engine behavior. Any mechanic that does not fit the Phase 2 vocabulary must become an explicit, reusable semantic extension rather than an implicit character branch.

## Representative coverage

| Fighter | Primary stress dimension | Fighter manifest | Deck manifest | Result |
| --- | --- | --- | --- | --- |
| Achilles | sidekick-death state, participant replacement, bonus attack | `phase-4a/achilles.yaml` | `../cards/phase-4a/achilles.yaml` | PASS |
| Bloody Mary | action ordinal/history, parent-combat capture | `phase-4a/bloody-mary.yaml` | `../cards/phase-4a/bloody-mary.yaml` | PASS |
| Sun Wukong | summon reserve, prevention/redirection, bonus attack | `phase-4a/sun-wukong.yaml` | `../cards/phase-4a/sun-wukong.yaml` | PASS |
| Sherlock Holmes | cancellation protection, hidden/public commitment, printed value | `phase-4a/sherlock-holmes.yaml` | `../cards/phase-4a/sherlock-holmes.yaml` | PASS |
| Dracula | multiple sidekicks, return fighter, combat-card replacement | `phase-4a/dracula.yaml` | `../cards/phase-4a/dracula.yaml` | PASS |
| Raptors | multiple heroes / independent health | `phase-4a/raptors.yaml` | `../cards/phase-4a/raptors.yaml` | PASS |
| Wayward Sisters | multiple heroes, alternate card zone, ingredients/spells | `phase-4a/wayward-sisters.yaml` | `../cards/phase-4a/wayward-sisters.yaml` | PASS |
| Geralt of Rivia | 36→30 pre-game construction, gear, ongoing schemes | `phase-4a/geralt-of-rivia.yaml` | `../cards/phase-4a/geralt-of-rivia.yaml` | PASS |
| Yennefer & Triss | selected hero role, role-dependent ability, simultaneous choices | `phase-4a/yennefer-triss.yaml` | `../cards/phase-4a/yennefer-triss.yaml` | PASS |
| Black Panther | foreign-card ownership/storage, sequential BOOST/draw | `phase-4a/black-panther.yaml` | `../cards/phase-4a/black-panther.yaml` | PASS |

## Deck quantity reconciliation

| Fighter | Available action-card pool | Game deck | Quantity check |
| --- | ---: | ---: | --- |
| Achilles | 30 | 30 | 30 |
| Bloody Mary | 30 | 30 | 30 |
| Sun Wukong | 30 | 30 | 30 |
| Sherlock Holmes | 30 | 30 | 30 |
| Dracula | 30 | 30 | 30 |
| Raptors | 30 | 30 | 30 |
| Wayward Sisters | 30 | 30 | 30 + 4 external spells |
| Geralt of Rivia | 36 | 30 | 24 base + choose 6 of 12 gear cards |
| Yennefer & Triss | 30 | 30 | 30 |
| Black Panther | 30 | 30 | 30 |

Bonus-attack definitions and Wayward spells are explicitly `external_definitions`; they do not inflate action-deck quantity.

## Corpus-proven schema extensions

Phase 4A required the following generic extensions, all promoted to `docs/mechanics/effect-model.md`:

- `PREVENT_DAMAGE`;
- `REDIRECT_DAMAGE`;
- `PREVENT_OPERATION`;
- `SET_PRINTED_VALUE`;
- `ADD_BOOST_VALUE`;
- `REORDER`;
- explicit parent-context capture;
- `REPLACE_COMBAT_CARD` composite.

The validation criterion is therefore satisfied: new requirements are documented as reusable semantics, not hidden as fighter-specific implementation code.

## Deterministic fixture A — Achilles defender replacement

1. Opponent attacks Achilles.
2. Achilles defends with `Under Achilles' Helm` while Patroclus is undefeated.
3. At `IMMEDIATELY`, resolve the source-defined optional swap.
4. If the swap succeeds, mutate the current combat participant state so Patroclus is defender.
5. Continue the same combat; no new Attack action begins.
6. If Patroclus is defeated, the persistent `patroclus_defeated` state changes and later eligible Achilles effects observe current state.

**Result:** ordinary `SWAP` + explicit combat state is sufficient.

## Deterministic fixture B — Bloody Reprise parent snapshot

1. Bloody Mary attacks with `Bloody Requiem`.
2. Parent combat records whether a defense card existed and its printed value.
3. Main combat completes through After Combat.
4. Trigger the nested bonus attack.
5. `Bloody Reprise` reads captured parent context rather than attempting to inspect a discarded/replaced card later.
6. One top-level Attack boundary remains.

**Result:** bonus combat needs captured context, not arbitrary access to historical object graphs.

## Deterministic fixture C — Sun Wukong damage redirection

1. Sun Wukong/Clone defends with `Golden Chain Mail`.
2. Calculate combat damage normally.
3. Before applying it to the defender, the card redirects that combat-damage event to the opposing combat fighter.
4. Apply defeat/downstream rules to the actual recipient.
5. Do not synthesize new effect damage or treat this as post-damage healing.

`Bewilderment` similarly prevents scoped damage before application.

**Result:** prevention/redirection require first-class damage-event semantics.

## Deterministic fixture D — Sherlock prediction and printed value

1. Sherlock commits `Elementary` face-up according to its source rule.
2. Before ordinary reveal, lock Sherlock's predicted attack value.
3. Reveal/resolve combat.
4. A separate `Deduce Strategy` interaction can alter the effective printed-value layer without rewriting immutable card-definition metadata.
5. Current combat value remains a separate layer.

**Result:** visibility overrides and `printed` vs `current` values are representable without parsing card text at runtime.

## Deterministic fixture E — Dracula replaces attack card

1. Opponent has a committed attack card.
2. `Do My Bidding` resolves in the current combat.
3. Return/move the old attack card as directed, inspect the opponent hand, choose a legal replacement.
4. Replace `combat.attack_card` with the selected card under the source-defined timing rule.
5. Resume the same combat.

**Result:** `REPLACE_COMBAT_CARD` is a reusable composite; this is not a second Attack.

## Deterministic fixture F — multiple heroes

Raptors and Wayward Sisters each use multiple hero instances with independent health.

- Damage and defeat target one concrete runtime fighter instance.
- Generic `your hero` selection resolves to one applicable hero unless a source says otherwise.
- Character defeat checks the topology's explicit `all_heroes_defeated` loss rule.
- Ordinary card selectors/MOVE/PLACE operate over surviving hero instances.

**Result:** no parallel multi-hero rules engine is required.

## Deterministic fixture G — Wayward Cauldron and spell

1. A Wayward action card would enter discard during cleanup/effect/BOOST use.
2. Bubbling Brew replaces that destination with public `cauldron`.
3. Ingredient tags are derived from card definitions currently in that zone.
4. At the eligible post-attack hook, choose at most one castable external spell.
5. Resolve the spell as `TRIGGER_COMPOSITE(kind=external_spell)`.
6. Perform source-defined Cauldron cleanup to discard.

Special case: `Something Wicked This Way Comes` explicitly bypasses normal ingredient/consumption procedure for its two spells.

**Result:** alternate zones + tags + external definitions are sufficient.

## Deterministic fixture H — Geralt setup reconnect and ongoing scheme

1. Setup reaches deck construction before shuffle/draw.
2. Persist the pending choice: one Potion, one Armor, one Sword definition.
3. Reconnect reconstructs exactly that private selection domain.
4. Include both copies of each selected gear definition with the 24 base cards; assert game deck size 30.
5. Shuffle/draw only after legal construction completes.
6. An ongoing scheme later moves to public `ongoing-schemes` and installs its continuous/replacement behavior until source-defined discard.

**Result:** setup choice and ongoing source lifetime fit existing event/resume semantics.

## Deterministic fixture I — Yennefer & Triss and Lodge

1. Setup locks which fighter is hero; set hero to 14 health, other fighter to sidekick 6, install only selected ability.
2. Later `Lodge of Sorceresses` requests one hidden hand-card choice from every player.
3. Persist all choices as committed-hidden.
4. Reveal only after all required choices lock.
5. Resolve each player's draw count from the revealed BOOST value and damage the hero(es) tied for the maximum according to source rules.

**Result:** no information leak and no client-owned setup shortcut.

## Deterministic fixture J — Black Panther foreign card ownership

1. `Analyze and Adjust` reveals the top opponent card.
2. Move that same card instance to public `vibranium-suit`; `owner_player` remains the opponent.
3. Later Black Panther chooses it as a BOOST source if it has a BOOST value.
4. Apply BOOST value.
5. Immediately resolve Black Panther's ability draw.
6. Resolve applicable BOOST effect.
7. Move the stored card to its **owner's** discard pile.

For `Wakanda Forever!`, repeat the optional BOOST decision only after step 5/6 changes authoritative state; do not lock both BOOST cards before the first ability draw.

**Result:** immutable owner + mutable zone/use authority is required and sufficient.

## Cross-mechanic fixture — `End the turn` vs Wayward cleanup

If a Wayward Sister attacks Geralt and Geralt's Armor of the Forgotten Wolf ends the attacker's turn during After Combat:

1. global `END_TURN` stops unresolved effects and redirects to Cleanup;
2. Cleanup still discards the played Wayward attack card;
3. Bubbling Brew can replace that discard with the Cauldron because its replacement is still active;
4. post-Cleanup `after attacking` additional effects are skipped by `END_TURN`, so no ordinary spell opportunity follows;
5. proceed to End-of-Turn.

**Result:** Phase 2 control flow and Phase 4A replacement semantics compose deterministically.

## Source and fan-content validation

All representative manifests distinguish:

- official set/rulebook sources;
- published UmDb `/umdb/...` data;
- current ruling indexes;
- excluded community `/decks/...` balance patches.

This matters especially for Bloody Mary and Sun Wukong, where community patch pages reuse original card names but change published values/quantities/effects.

## Remaining ambiguities

There are **no known P0/P1 ambiguities required to execute the ten representative manifests under their documented source rules**.

The global ambiguity register still contains deferred questions for future Phase 4B content (for example Stars & Stripes specifics, large-fighter geometry, or other characters' dormant/revival behavior). Phase 4A does not pretend those are resolved.

## Gate decision

**PASS — Phase 4A is complete.**

The representative corpus demonstrates that:

- fighter topology, setup and resources have a stable schema;
- deck construction supports fixed and constructed decks;
- action cards and external definitions have a normalized representation;
- ownership, zones and hidden/public state are sufficient for online persistence/reconnect;
- new semantics discovered by real cards were promoted to generic documented operations/composites;
- no sample fighter depends on opaque undocumented character code.

Phase 4B may now expand this schema to every released competitive fighter in `docs/sets/registry.yaml`.