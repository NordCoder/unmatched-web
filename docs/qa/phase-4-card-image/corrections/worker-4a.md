# Phase 4 Card-Image Corrections — Worker 4A

## Identity

Repository: `NordCoder/unmatched-web`  
Base branch: `main`  
Correction branch: `phase-4-card-image-fixes-4a`  
Base Head: `3a4196f0a1596d9411971da47bd652a69613f060`  
Starting correction Head: `3a4196f0a1596d9411971da47bd652a69613f060`  
Final Head: commit containing this report; the exact SHA is emitted by the GitHub Connector and repeated in the external correction handoff because a Git commit cannot embed its own SHA.  
Owned fighters: `achilles`, `bloody-mary`, `sun-wukong`, `sherlock-holmes`, `dracula`, `raptors`, `wayward-sisters`, `geralt-of-rivia`, `yennefer-and-triss`, `black-panther`.

Canonical identity note: the QA contract calls the selectable Witcher pair `yennefer-and-triss`, while the current Phase 4A canonical file stem and `fighter_id` remain `yennefer-triss`. No `yennefer-and-triss.yaml` file exists on the correction branch, and no gameplay or ID rewrite was made from the P3 archive/orchestration alias finding.

## Reports reviewed

- `docs/qa/phase-4-card-image/reports/alice.md`
- `docs/qa/phase-4-card-image/reports/bigfoot.md`
- `docs/qa/phase-4-card-image/reports/buffy.md`
- `docs/qa/phase-4-card-image/reports/angel.md`
- `docs/qa/phase-4-card-image/reports/achilles.md`
- `docs/qa/phase-4-card-image/reports/bullseye.md`
- `docs/qa/phase-4-card-image/reports/cloak-and-dagger.md`
- `docs/qa/phase-4-card-image/reports/black-panther.md`
- `docs/qa/phase-4-card-image/reports/annie-christmas.md`
- `docs/qa/phase-4-card-image/reports/golden-bat.md`
- `docs/qa/phase-4-card-image/reports/ciri.md`
- `docs/qa/phase-4-card-image/reports/ancient-leshen.md`
- `docs/qa/phase-4-card-image/reports/chupacabra.md`
- `docs/qa/phase-4-card-image/reports/donatello.md`
- `docs/qa/phase-4-card-image/reports/george-washington.md`

Reports reviewed: **15/15**

## Owned findings inventory

| Finding ID | Report | Fighter | Card | Severity | Classification | Disposition |
|---|---|---|---|---|---|---|
| `RAPTORS-P1-AMBUSH` (`F-02`) | `bigfoot.md` | raptors | Ambush | P1 | canonical_manifest_error | Applied: `versatile` → `attack`. |
| `RAPTORS-P1-EATEN-ALIVE` (`F-03`) | `bigfoot.md` | raptors | Eaten Alive | P1 | canonical_manifest_error | Applied: `attack` → `versatile`. |
| `RAPTORS-QUAL-DISENGAGE` (`W-02`) | `bigfoot.md` | raptors | Disengage | qualification | already_resolved | No change: current PLACE normalization is explicitly errata-aware; the image-only MOVE wording does not prove a canonical transcription error. |
| `SW-P2-001` | `achilles.md` | sun-wukong | Bewilderment | P2 | canonical_manifest_error | Applied: damage prevention narrowed from all damage in the combat window to combat damage only. |
| `SW-P3-001` | `achilles.md` | sun-wukong | 72 Transformations | P3 | evidence_package_only | No Git change: archive `72-transformations` maps to canonical `seventy-two-transformations`. |
| `SUN-WUKONG-QUAL-TRICKED-YOU` | `achilles.md` | sun-wukong | Tricked You | qualification | already_resolved | No change: remains an external bonus-attack definition, not an action-card copy. |
| `ACHILLES-QUAL-RELENTLESS-ASSAULT` | `achilles.md` | achilles | Relentless Assault | qualification | already_resolved | No change: remains an external bonus-attack definition, not an action-card copy. |
| `BLOODY-MARY-QUAL-BLOODY-REPRISE` | `achilles.md` | bloody-mary | Bloody Reprise | qualification | already_resolved | No change: remains an external bonus-attack definition with captured parent-combat context. |
| `B08-BP-001` | `black-panther.md` | black-panther | Analyze and Adjust | P1 | canonical_manifest_error | Applied: `versatile` → `attack`. |
| `B08-BP-002` | `black-panther.md` | black-panther | Ancestral Insight | P1 | canonical_manifest_error | Applied: `attack` → `versatile`. |
| `B08-BP-003` | `black-panther.md` | black-panther | Anti-Metal Claws | P1 | canonical_manifest_error | Applied: `attack` → `versatile`. |
| `B08-BP-004` | `black-panther.md` | black-panther | Evade | P1 | canonical_manifest_error | Applied: `versatile` → `defense`. |
| `B08-BP-005` | `black-panther.md` | black-panther | Nanotriage Processor | P1 | canonical_manifest_error | Applied: `defense` → `versatile`. |
| `B08-BP-006` | `black-panther.md` | black-panther | Vibranium Shockwave | P1 | canonical_manifest_error | Applied: `versatile` → `attack`. |
| `WAYWARD-SISTERS-QUAL-EXTERNAL-SPELLS` | `ciri.md` | wayward-sisters | four external spells | qualification | already_resolved | No change: spells remain external gameplay definitions and are not ordinary action-card instances. |
| `GERALT-001` | `ciri.md` | geralt-of-rivia | Damn, You're Ugly | P3 | evidence_package_only | No Git change: archive `damn-you-re-ugly` maps to canonical `damn-youre-ugly`. |
| `GERALT-QUAL-CONSTRUCTION` | `ciri.md` | geralt-of-rivia | deck construction | qualification | already_resolved | No change: available pool 36 and constructed game deck 30 remain distinct and correct. |
| `B12-YT-P3-01` | `ancient-leshen.md` | yennefer-and-triss / canonical `yennefer-triss` | fighter/path alias | P3 | evidence_package_only | No Git rename: archive/orchestration alias differs from the unambiguous canonical file stem. |
| `B12-YT-P3-02` | `ancient-leshen.md` | yennefer-and-triss / canonical `yennefer-triss` | Merigold's Hailstorm | P3 | evidence_package_only | No Git change: archive `merigold-s-hailstorm` maps to canonical `merigolds-hailstorm`. |

No owned P1/P2/P3/qualification finding was found in the remaining reviewed fighter sections for `sherlock-holmes` or `dracula`. Their card-image comparisons were clean.

## Applied corrections

### Raptors / Ambush

Finding: `RAPTORS-P1-AMBUSH` (`F-02`)  
Fighter: `raptors`  
Card: `Ambush`  
Severity: P1  
Report: `bigfoot.md`  
Before: `type: versatile`  
After: `type: attack`  
Files: `docs/cards/phase-4a/raptors.yaml`  
Reason: the physical card has the red ATTACK type; defense-side use was incorrectly permitted.

### Raptors / Eaten Alive

Finding: `RAPTORS-P1-EATEN-ALIVE` (`F-03`)  
Fighter: `raptors`  
Card: `Eaten Alive`  
Severity: P1  
Report: `bigfoot.md`  
Before: `type: attack`  
After: `type: versatile`  
Files: `docs/cards/phase-4a/raptors.yaml`  
Reason: the physical card has the purple VERSATILE type; defense-side use was incorrectly removed.

### Sun Wukong / Bewilderment

Finding: `SW-P2-001`  
Fighter: `sun-wukong`  
Card: `Bewilderment`  
Severity: P2  
Report: `achilles.md`  
Before: `scope: all_damage_during_combat_window`  
After: `scope: combat_damage`  
Files: `docs/cards/phase-4a/sun-wukong.yaml`  
Reason: the printed effect prevents combat damage, not unrelated effect damage that happens during the broader combat window.

### Black Panther / Analyze and Adjust

Finding: `B08-BP-001`  
Fighter: `black-panther`  
Card: `Analyze and Adjust`  
Severity: P1  
Report: `black-panther.md`  
Before: `type: versatile`  
After: `type: attack`  
Files: `docs/cards/phase-4a/black-panther.yaml`  
Reason: printed red ATTACK icon.

### Black Panther / Ancestral Insight

Finding: `B08-BP-002`  
Fighter: `black-panther`  
Card: `Ancestral Insight`  
Severity: P1  
Report: `black-panther.md`  
Before: `type: attack`  
After: `type: versatile`  
Files: `docs/cards/phase-4a/black-panther.yaml`  
Reason: printed purple VERSATILE icon.

### Black Panther / Anti-Metal Claws

Finding: `B08-BP-003`  
Fighter: `black-panther`  
Card: `Anti-Metal Claws`  
Severity: P1  
Report: `black-panther.md`  
Before: `type: attack`  
After: `type: versatile`  
Files: `docs/cards/phase-4a/black-panther.yaml`  
Reason: printed purple VERSATILE icon.

### Black Panther / Evade

Finding: `B08-BP-004`  
Fighter: `black-panther`  
Card: `Evade`  
Severity: P1  
Report: `black-panther.md`  
Before: `type: versatile`  
After: `type: defense`  
Files: `docs/cards/phase-4a/black-panther.yaml`  
Reason: printed blue DEFENSE icon.

### Black Panther / Nanotriage Processor

Finding: `B08-BP-005`  
Fighter: `black-panther`  
Card: `Nanotriage Processor`  
Severity: P1  
Report: `black-panther.md`  
Before: `type: defense`  
After: `type: versatile`  
Files: `docs/cards/phase-4a/black-panther.yaml`  
Reason: printed purple VERSATILE icon.

### Black Panther / Vibranium Shockwave

Finding: `B08-BP-006`  
Fighter: `black-panther`  
Card: `Vibranium Shockwave`  
Severity: P1  
Report: `black-panther.md`  
Before: `type: versatile`  
After: `type: attack`  
Files: `docs/cards/phase-4a/black-panther.yaml`  
Reason: printed red ATTACK icon.

### YAML scalar quoting

Finding: validation-only syntax defect exposed by standard YAML parsing  
Fighter: `raptors`, `black-panther`  
Card: none  
Severity: validation  
Report: none; current-state parser recheck  
Before: normalization-note scalar began with an unquoted backtick  
After: the same note text is double-quoted  
Files: `docs/cards/phase-4a/raptors.yaml`, `docs/cards/phase-4a/black-panther.yaml`  
Reason: YAML 1.2 plain scalars cannot begin with a backtick. This is a syntax-only correction with no data or semantic change.

## Already resolved

- `raptors / Disengage`: current PLACE representation remains the documented errata-aware normalization; no image-only reversal was applied.
- `achilles / Relentless Assault`: current external bonus-attack definition is correct and preserved.
- `bloody-mary / Bloody Reprise`: current external bonus-attack definition and captured parent-defense value are correct and preserved.
- `sun-wukong / Tricked You`: current external bonus-attack definition is correct and preserved.
- `wayward-sisters`: all four spell definitions remain external and outside the 30-card action deck.
- `geralt-of-rivia`: `available_pool_count: 36`, base quantity 24, selectable Gear quantity 12, selected Gear quantity 6, and `game_deck_count: 30` remain correct.

## Evidence-package-only findings

- `SW-P3-001`: archive `72-transformations` should map to canonical `seventy-two-transformations`.
- `GERALT-001`: archive `damn-you-re-ugly` should map to canonical `damn-youre-ugly`.
- `B12-YT-P3-01`: archive/orchestration fighter alias `yennefer-and-triss` should resolve to canonical `yennefer-triss`; no canonical fighter duplicate or rename was created.
- `B12-YT-P3-02`: archive `merigold-s-hailstorm` should map to canonical `merigolds-hailstorm`.

## Integration requirements confirmed

No owned finding required a new shared semantic extension in this correction scope. Existing external bonus-attack, external-spell, fixed-deck, and choose-group construction representations were preserved rather than rewritten as card transcription fixes.

## Validation

YAML: **PASS** — standard safe parsing succeeded for all three corrected deck manifests after syntax-only normalization-note quoting.  
Quantities: **PASS** — `raptors`, `sun-wukong`, and `black-panther` each retain quantity sum 30, available pool 30, and game deck 30.  
References: **PASS** — corrected manifests retain valid `usable_by` references; Sun Wukong's `tricked-you` composite reference resolves to the preserved external definition; no external definitions were lost.  
Construction: **PASS** — all three corrected fighters remain fixed 30-card decks; the reviewed Geralt 36-card available pool / 30-card constructed deck distinction remains unchanged.  
Operation and choice bindings: **PASS** — no operation, choice, capture, dependency, zone, resource, or target binding was changed by the corrections; current references remain structurally present.  
Affected card recheck: **PASS** — printed value, BOOST, quantity, `usable_by`, effects, ordering, and optionality for all nine corrected definitions are unchanged except for the proven card-type or damage-scope field.  
Scope: **PASS** — branch comparison against `main` before this report showed only the three allowed owned Phase 4A card manifests; this report adds only the allowed correction-report path.  
Unexpected files: **none**.

## Remaining blockers

None.

The P3 archive/orchestration mapping items require evidence-package maintenance outside the allowed repository paths, but they do not block the canonical Phase 4A manifest corrections.

## Correction Handoff

Final Head: commit containing this report; exact Connector-returned SHA is supplied in the external handoff.  
Owned fighters searched: **10/10**  
Owned P1 findings accounted for: **8/8**  
Owned P2 findings accounted for: **1/1**  
Applied: **9**  
Already resolved: **6**  
Evidence-only: **4**  
Blocked: **0**  
Files changed: **4** (`docs/cards/phase-4a/raptors.yaml`, `docs/cards/phase-4a/sun-wukong.yaml`, `docs/cards/phase-4a/black-panther.yaml`, `docs/qa/phase-4-card-image/corrections/worker-4a.md`)
