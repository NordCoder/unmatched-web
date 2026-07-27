# Phase 4 Card-Image QA Report

## Batch

- golden-bat
- nikola-tesla
- oda-nobunaga
- shakespeare
- tomoe-gozen

Canonical branch inspected: `phase-4b-worker-c-modern`

Branch tip observed immediately before persistence: `1f4fd973de66ab7995c0b2213e7dcda4163ff7c8`

## Verdict

**FAIL**

## Summary

The batch contained five fighter ZIP archives and 62 card-image files. All five nested ZIPs were readable, all image files decoded successfully, no zero-byte/corrupted files were found, and no binary-duplicate images were found.

Four fighters (`golden-bat`, `nikola-tesla`, `oda-nobunaga`, `tomoe-gozen`) have no material P1/P2 gameplay discrepancies against the canonical Git manifests. `nikola-tesla` and `tomoe-gozen` contain only archive-side P3 card-ID normalization issues.

`shakespeare` fails for two material reasons:

1. `Et Tu, Brute?` is normalized in Git with `friendly_fighters` where the printed card says **your fighters**, which broadens the legal fighter domain in team play and changes gameplay behavior.
2. The evidence archive models five separate Horror image entries at quantity 1 instead of one unique `Horror` card definition at quantity 5. The canonical Git deck representation is correct; the archive corpus construction is not.

Batch-level counts:

| Metric | Result |
|---|---:|
| Fighters received | 5 |
| Fighters fully checked | 5 |
| PASS | 0 |
| PASS_WITH_QUALIFICATIONS | 4 |
| FAIL | 1 |
| BLOCKED | 0 |
| Archive image files | 62 |
| Images successfully inspected | 62 / 62 |
| Unreadable images | 0 |
| Missing images | 0 |
| Binary duplicate images | 0 |
| Redundant definition-level images | 4 |
| Quantity-validation failures | 1 fighter (`shakespeare`) |
| P3 metadata findings | 3 |
| Gameplay-semantic discrepancy cards | 1 (`Et Tu, Brute?`) |
| P1 findings | 2 |

## Fighter Results

### golden-bat

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/golden-bat.yaml`
- fighter manifest: `docs/fighters/phase-4b/golden-bat.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS**
- printed card identity/content verification: **PASS**
- discrepancies: none
- verdict: **PASS_WITH_QUALIFICATIONS**

Archive / construction evidence:

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12 / 12 |
| Unique binary images | 12 |
| Manifest card entries | 12 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 12 |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| A Punch to Shake the Earth | PASS | PASS | PASS | PASS |
| Alpine Fortress | PASS | PASS | PASS | PASS |
| Arrive Just in Time | PASS | PASS | PASS | PASS |
| He Laughs at Your Feebleness | PASS | PASS | PASS | PASS |
| Imposing Presence | PASS | PASS | PASS | PASS |
| Insight of the Ancients | PASS | PASS | PASS | PASS |
| Like a Flash of Golden Light | PASS | PASS | PASS | PASS |
| Sight Beyond Sight | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Super Strength | PASS | PASS | PASS | PASS |
| Terrifying Roar | PASS | PASS | PASS | PASS |
| Vaporizing Eyebeams | PASS | PASS | PASS | PASS |

Important semantics checked included the conditional damage replacement on `A Punch to Shake the Earth`, ordered discard/shuffle/draw/move resolution on `Alpine Fortress`, turn-start-space condition on `Like a Flash of Golden Light`, ordered private top-deck inspection/reordering on `Sight Beyond Sight`, and the global cannot-leave restriction on `Terrifying Roar`.

Integration requirement directly confirmed by card images:

- `C-REQ-013`: ordered/dependent staged resolution, especially `Alpine Fortress` and `Sight Beyond Sight`.

Qualification is integration-only; no transcription/gameplay discrepancy was found.

### nikola-tesla

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/nikola-tesla.yaml`
- fighter manifest: `docs/fighters/phase-4b/nikola-tesla.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS_WITH_P3_MAPPING_FINDING**
- unique-card image completeness: **PASS**
- printed card identity/content verification: **PASS**
- discrepancies: one P3 archive card-ID mismatch
- verdict: **PASS_WITH_QUALIFICATIONS**

Archive / construction evidence:

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 11 / 11 |
| Unique binary images | 11 |
| Manifest card entries | 11 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 11 |
| Unique Git definitions | 11 |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| 7 Hertz | PASS | FAIL — P3 ID | PASS | QUALIFIED |
| Death Ray | PASS | PASS | PASS | PASS |
| Fully Charged | PASS | PASS | PASS | PASS |
| Intense Experimentation | PASS | PASS | PASS | PASS |
| Kinetic Induction | PASS | PASS | PASS | PASS |
| Lightning Storm | PASS | PASS | PASS | PASS |
| Polyphase Coils | PASS | PASS | PASS | PASS |
| Remote Control | PASS | PASS | PASS | PASS |
| Repulsion Blast | PASS | PASS | PASS | PASS |
| The Alternating Current | PASS | PASS | PASS | PASS |
| X-Ray Radiation | PASS | PASS | PASS | PASS |

Important semantics checked included one-coil/two-coil tier behavior, `instead`/`also` dependencies, reveal-before-discharge on `X-Ray Radiation`, and move-before-discharge on `Repulsion Blast`. The printed cards support the tiered discharge structure, while the special underfunded-declaration legality remains ruling-derived rather than visually proven by these card images.

Integration requirements directly supported:

- `C-REQ-004`: tiered declared resource discharge and dependent outcomes; the underfunded-declaration edge case itself remains ruling-derived.
- `C-REQ-013`: ordered effect stages and dependent continuation.

### oda-nobunaga

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/oda-nobunaga.yaml`
- fighter manifest: `docs/fighters/phase-4b/oda-nobunaga.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS**
- unique-card image completeness: **PASS**
- printed card identity/content verification: **PASS**
- discrepancies: none
- verdict: **PASS_WITH_QUALIFICATIONS**

Archive / construction evidence:

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 11 / 11 |
| Unique binary images | 11 |
| Manifest card entries | 11 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 11 |
| Unique Git definitions | 11 |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Battle Maneuvers | PASS | PASS | PASS | PASS |
| Demon King of the Sixth Heaven | PASS | PASS | PASS | PASS |
| Fire and Flames | PASS | PASS | PASS | PASS |
| Lightning and Thunder | PASS | PASS | PASS | PASS |
| Momentous Shift | PASS | PASS | PASS | PASS |
| Patience and Strategy | PASS | PASS | PASS | PASS |
| Pragmatism | PASS | PASS | PASS | PASS |
| Reinforce | PASS | PASS | PASS | PASS |
| Spring the Trap | PASS | PASS | PASS | PASS |
| Student of War | PASS | PASS | PASS | PASS |
| Sun and Moon | PASS | PASS | PASS | PASS |

Important semantics checked:

- `Battle Maneuvers` says **your fighters** and is correctly normalized to fighters controlled by the player rather than all friendly fighters.
- `Spring the Trap` says **adjacent friendly fighter** and intentionally uses the broader friendly domain.
- `Reinforce` requires exactly two different effects.
- `Student of War` draws according to combat damage actually dealt.
- flanking conditions and fighter restrictions match the images.

Integration requirements directly supported:

- `C-REQ-003`: same-combat defender replacement on `Spring the Trap`.
- `C-REQ-010`: historical turn-start space semantics on `Momentous Shift`, with runtime-instance identity additionally relevant because Oda has two Honor Guards.
- `C-REQ-013`: optional replacement and dependent continuation on `Spring the Trap`.

### shakespeare

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/shakespeare.yaml`
- fighter manifest: `docs/fighters/phase-4b/shakespeare.yaml`
- archive integrity: **FAIL at corpus-definition level**
- canonical manifest comparison: **FAIL**
- unique-card image completeness: **FAIL**
- printed card identity/content verification: **FAIL**
- discrepancies: two P1 material findings
- verdict: **FAIL**

Technical file integrity is clean: all image files open and decode. The failure is structural/semantic, not file corruption.

Archive / construction evidence:

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 16 / 16 |
| Unique binary images | 16 |
| Manifest card entries | 16 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images vs archive manifest | 0 |
| Excess images vs canonical definitions | 4 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Archive manifest entries | 16 |
| Canonical unique definitions | 12 |

Five separate images/manifest entries represent the same printed `Horror` definition. All five carry the same printed card name, restriction, ATTACK 4, BOOST 3, two syllables, ordinary effect, Line completion effect, and printed quantity `×5`. They are artwork/copy variants of one card definition, not five independent gameplay definitions.

Archive effectively models:

```text
horror          ×1
horror-539      ×1
horror-539-8    ×1
horror-539-9    ×1
horror-539-10   ×1
```

Canonical Git correctly models:

```text
horror ×5
```

Per-image/card verification:

| Card / image | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Again | PASS | PASS | PASS | PASS |
| Alas | PASS | PASS | PASS | PASS |
| All Are Punished | PASS | PASS | PASS | PASS |
| Deceive | PASS | PASS | PASS | PASS |
| Et Tu, Brute? | PASS | PASS | FAIL — P1 | FAIL |
| Horror — `horror.webp` | PASS | FAIL — P1 archive qty/identity | PASS | FAIL |
| Horror — `horror-539.webp` | PASS | FAIL — P1 archive qty/identity | PASS | FAIL |
| Horror — `horror-539-8.webp` | PASS | FAIL — P1 archive qty/identity | PASS | FAIL |
| Horror — `horror-539-9.webp` | PASS | FAIL — P1 archive qty/identity | PASS | FAIL |
| Horror — `horror-539-10.webp` | PASS | FAIL — P1 archive qty/identity | PASS | FAIL |
| My Kingdom For a Horse | PASS | PASS | PASS | PASS |
| Once More Unto The Breach | PASS | PASS | PASS | PASS |
| Places, Places! | PASS | PASS | PASS | PASS |
| Revise | PASS | PASS | PASS | PASS |
| Such Sweet Sorrow | PASS | PASS | PASS | PASS |
| The Ides Of March | PASS | PASS | PASS | PASS |

Integration requirements directly supported:

- `C-REQ-001`: the images directly demonstrate syllable counts, ordered Line membership, and separate Line-completion effects.
- `C-REQ-013`: ordered interactions such as `Again` resolving before the card is appended to the Line.

### tomoe-gozen

- branch: `phase-4b-worker-c-modern`
- deck manifest: `docs/cards/phase-4b/tomoe-gozen.yaml`
- fighter manifest: `docs/fighters/phase-4b/tomoe-gozen.yaml`
- archive integrity: **PASS**
- canonical manifest comparison: **PASS_WITH_P3_MAPPING_FINDINGS**
- unique-card image completeness: **PASS**
- printed card identity/content verification: **PASS**
- discrepancies: two P3 archive card-ID mismatches
- verdict: **PASS_WITH_QUALIFICATIONS**

Archive / construction evidence:

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12 / 12 |
| Unique binary images | 12 |
| Manifest card entries | 12 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 12 |

Per-card verification:

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| A Warrior's Way | PASS | FAIL — P3 ID | PASS | QUALIFIED |
| A Worthy Opponent | PASS | PASS | PASS | PASS |
| Confront Any Demon or God | PASS | PASS | PASS | PASS |
| Deeds of Valor | PASS | PASS | PASS | PASS |
| Fearsome Strength | PASS | PASS | PASS | PASS |
| Five Against Thousands | PASS | PASS | PASS | PASS |
| Flash of Steel | PASS | PASS | PASS | PASS |
| Lord Kiso's Final Stand | PASS | FAIL — P3 ID | PASS | QUALIFIED |
| Piercing Shot | PASS | PASS | PASS | PASS |
| Refuse to Retreat | PASS | PASS | PASS | PASS |
| Skirmish | PASS | PASS | PASS | PASS |
| Witness My Last Battle | PASS | PASS | PASS | PASS |

Important semantics checked included the ordered `then` relationship on `A Warrior's Way`, hidden-information handling on `Confront Any Demon or God`, optional discard cost on `Five Against Thousands`, sidekick traversal followed by a mandatory benefit choice on `Lord Kiso's Final Stand`, and face-up/adjacent-target declaration behavior on `Witness My Last Battle`.

Integration requirements directly supported:

- `C-REQ-003`: declaration-time target restriction and face-up commitment on `Witness My Last Battle`.
- `C-REQ-013`: ordered/dependent stages on several cards.

## Findings

### SHAKESPEARE-001 — `Et Tu, Brute?` target-domain mismatch

- fighter: `shakespeare`
- card: `Et Tu, Brute?`
- severity: **P1**
- Git location: `docs/cards/phase-4b/shakespeare.yaml`, `brute-during` and `brute-completion`
- expected: printed **your fighters** must mean fighters controlled by Shakespeare's player.
- observed: Git uses `friendly_fighters` for both the during-combat adjacency count and Line-completion placement domain/target.
- evidence/reasoning: in team play, `friendly_fighters` can include teammate-controlled fighters, while the printed card says **your fighters**. The repository itself already distinguishes these concepts on Oda: `Battle Maneuvers` is scoped to controller fighters for printed **your fighters**, while `Spring the Trap` intentionally uses friendly fighters for printed **friendly fighter**.
- gameplay impact: Git can over-count adjacent teammate fighters for combat value and can permit placement of teammate fighters during Line completion.
- expected correction: scope both selectors to fighters controlled by Shakespeare's player (`controller_fighters` or schema-equivalent owner/controller-scoped selector). Do not broaden to all friendly fighters.

### SHAKESPEARE-002 — Horror definition/quantity corpus mismatch

- fighter: `shakespeare`
- card: `Horror`
- severity: **P1**
- Git location: canonical `horror`, `quantity: 5`
- expected: one unique `Horror` definition with quantity 5 and one representative evidence image for that unique definition.
- observed: five archive manifest/image entries are treated as separate definitions at quantity 1 (`horror`, `horror-539`, `horror-539-8`, `horror-539-9`, `horror-539-10`).
- evidence/reasoning: all five images visibly print the same card identity, ATTACK 4, BOOST 3, two syllables, identical normal/Line effects, and `×5`. `Horror` mechanics also explicitly count Horror cards in the Line, making definition identity gameplay-relevant rather than a presentation-only detail.
- expected correction: evidence corpus should retain one canonical `horror` definition/image mapping with quantity 5; artwork/copy variants must not be represented as independent gameplay definitions. Canonical Git deck quantity is already correct.

### TESLA-META-001 — archive card-ID normalization

- fighter: `nikola-tesla`
- card: `7 Hertz`
- severity: **P3**
- expected: canonical card ID `seven-hertz`
- observed: archive uses `7-hertz`
- evidence/reasoning: printed identity, quantity, type, value, BOOST, restriction and gameplay semantics all match the canonical `seven-hertz` definition.
- expected correction: normalize the evidence/archive mapping to `seven-hertz` or define an evidence-only alias. No canonical gameplay change is indicated.

### TOMOE-META-001 — archive card-ID normalization

- fighter: `tomoe-gozen`
- card: `A Warrior's Way`
- severity: **P3**
- expected: canonical card ID `a-warriors-way`
- observed: archive uses `a-warrior-s-way`
- evidence/reasoning: printed card and canonical gameplay content otherwise match.
- expected correction: normalize archive mapping only.

### TOMOE-META-002 — archive card-ID normalization

- fighter: `tomoe-gozen`
- card: `Lord Kiso's Final Stand`
- severity: **P3**
- expected: canonical card ID `lord-kisos-final-stand`
- observed: archive uses `lord-kiso-s-final-stand`
- evidence/reasoning: printed card and canonical gameplay content otherwise match.
- expected correction: normalize archive mapping only.

## Corpus-Level Observations

1. All five nested fighter ZIPs were technically healthy and all 62 supplied images were visually inspectable.
2. No zero-byte, truncated, corrupted, or binary-duplicate image files were found.
3. Total quantity alone is insufficient to validate a corpus: Shakespeare still sums to 30 while incorrectly splitting a single `Horror ×5` definition across five pseudo-definitions.
4. Filename or archive manifest identity was not treated as authoritative. Printed images and canonical Git definitions were used to resolve card identity.
5. Printed semantic verification was performed against normalized Git behavior, not literal prose equality.
6. The distinction between **your fighters** and **friendly fighters** is material in team play and must be preserved consistently across normalized manifests.
7. Shared integration requirements (`C-REQ-*`) were treated as valid documented dependencies rather than transcription failures when the printed card behavior itself matched the Git representation.
8. Tesla's underfunded two-coil declaration remains a ruling/normalization qualification; the card images directly confirm tier structure but not that special declaration legality by themselves.
9. No fighter was blocked by unreadable evidence or missing canonical Git source.

## Final Assessment

This batch is **not fully suitable for acceptance as a clean Phase 4 card-image corpus** in its current state because of the Shakespeare material findings.

The four non-Shakespeare fighters are suitable from a card-image/transcription perspective, subject only to their documented integration qualifications and the three P3 archive-ID normalization issues.

For Shakespeare, the canonical Git `Horror ×5` deck construction is supported by the physical images, while the archive corpus must be de-duplicated at the definition level. Separately, the canonical Shakespeare manifest requires semantic correction for `Et Tu, Brute?` so that printed **your fighters** does not expand to teammate-controlled friendly fighters.

Batch verdict remains **FAIL** until the material Shakespeare discrepancies are resolved and re-QA'd.
