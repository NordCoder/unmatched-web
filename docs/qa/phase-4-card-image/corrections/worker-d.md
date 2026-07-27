# Phase 4 Card-Image Corrections — Worker D

## Identity

- Repository: `NordCoder/unmatched-web`
- Canonical source branch: `phase-4b-worker-d-latest`
- Correction branch: `phase-4-card-image-fixes-d`
- Base Head: `c154bf50acf259c95d0b7baf4877617556e38187`
- Pre-report correction content Head: `4fd457cfceed5511fe18fbdc081e8599710d998e`
- Final Head: current `phase-4-card-image-fixes-d` branch tip after this report commit; exact SHA is emitted in the external correction handoff because a commit cannot contain its own SHA.
- Scope: Worker D owned fighter/card manifests only, plus this correction report.

## Reports reviewed

Reports reviewed: **15/15**.

1. `docs/qa/phase-4-card-image/reports/alice.md`
2. `docs/qa/phase-4-card-image/reports/bigfoot.md`
3. `docs/qa/phase-4-card-image/reports/buffy.md`
4. `docs/qa/phase-4-card-image/reports/angel.md`
5. `docs/qa/phase-4-card-image/reports/achilles.md`
6. `docs/qa/phase-4-card-image/reports/bullseye.md`
7. `docs/qa/phase-4-card-image/reports/cloak-and-dagger.md`
8. `docs/qa/phase-4-card-image/reports/black-panther.md`
9. `docs/qa/phase-4-card-image/reports/annie-christmas.md`
10. `docs/qa/phase-4-card-image/reports/golden-bat.md`
11. `docs/qa/phase-4-card-image/reports/ciri.md`
12. `docs/qa/phase-4-card-image/reports/ancient-leshen.md`
13. `docs/qa/phase-4-card-image/reports/chupacabra.md`
14. `docs/qa/phase-4-card-image/reports/donatello.md`
15. `docs/qa/phase-4-card-image/reports/george-washington.md`

The ten reports for Alice, Buffy, Angel, Achilles, Bullseye, Cloak & Dagger, Black Panther, Annie Christmas, Golden Bat and Ciri contain no Worker D owned fighters. Mixed-report foreign findings were not modified.

Owned fighters searched: **16/16**:

`bruce-lee`, `muhammad-ali`, `blackbeard`, `chupacabra`, `loki`, `pandora`, `leonardo`, `donatello`, `michelangelo`, `raphael`, `rosie-the-riveter`, `john-henry`, `wyatt-earp`, `george-washington`, `shredder`, `krang`.

## Owned findings inventory

| Finding | Report | Fighter | Card / component | Severity | Classification | Disposition |
|---|---|---|---|---|---|---|
| `F-01` | `bigfoot.md` | bruce-lee | One-Inch Punch | P1 | target domain | applied |
| `F-04` | `bigfoot.md` | bruce-lee | "HOO! WHAAAAAA!" | P3 | archive ID mapping | evidence-package-only |
| `B12-BB-P3-01` | `ancient-leshen.md` | blackbeard | A Brace of Primed Pistols | P3 | archive ID mapping | evidence-package-only |
| `B12-BB-P3-02` | `ancient-leshen.md` | blackbeard | Queen Anne's Revenge | P3 | archive ID mapping | evidence-package-only |
| `CHUPACABRA-001` | `chupacabra.md` | chupacabra | Natural Toughness | P2 | timing window | applied |
| `LOKI-001` | `chupacabra.md` | loki | Looking for Trouble | P1 | condition and replacement ordering/domain | applied |
| `LOKI-002` | `chupacabra.md` | loki | five TRICK/archive mappings | P3 | archive IDs/name typo | evidence-package-only |
| `PANDORA-001` | `chupacabra.md` | pandora | three apostrophe-derived IDs | P3 | archive ID mapping | evidence-package-only |
| `LEONARDO-001` | `chupacabra.md` | leonardo | Eat, Sleep, and Breathe Ninjutsu | P3 | archive spelling/ID | evidence-package-only |
| `MIC-P1-01` | `donatello.md` | michelangelo | Shell Insertion | P1 | unsupported owner-selected ordering | applied |
| `MIC-P1-02` | `donatello.md` | michelangelo | Shell Insertion | P1 | prerequisite/cost dependency | applied |
| `ROS-P1-01` | `donatello.md` | rosie-the-riveter | D-Day | P1 | fixed conditional amount represented as scaling | applied |
| `JH-P1-01` | `donatello.md` | john-henry | Twelve-Pound Hammer | P1 | optionality | applied |
| `DON-P3-01` | `donatello.md` | donatello | The Future of Ninjutsu | P3 | archive ID mapping | evidence-package-only |
| `MIC-P3-01` | `donatello.md` | michelangelo | Let's Go!! | P3 | archive slug / title punctuation | evidence-package-only |
| `MIC-P3-02` | `donatello.md` | michelangelo | Turtle Power!! | P3 | title punctuation | evidence-package-only |
| `RAP-P3-01` | `donatello.md` | raphael | Let's Do This! | P3 | archive ID mapping | evidence-package-only |
| `WYATT-001` | `george-washington.md` | wyatt-earp | Gunfight at the O.K. Corral | P1 | exact effect count and alternating ownership | applied |
| `GEORGE-001` | `george-washington.md` | george-washington | Gather Information | P1 | target ownership domain | applied |
| `SHREDDER-001` | `george-washington.md` | shredder | Back to Work! | P1 | mandatory placement constraint | applied |
| `WYATT-002` | `george-washington.md` | wyatt-earp | three archive IDs | P3 | archive ID mapping | evidence-package-only |
| Krang die qualification | `george-washington.md` | krang | Die of Ultimate Destruction | qualification | missing non-card physical evidence | evidence qualification retained; no card correction |

Material accounting:

- Owned P1: **9/9 accounted for; 9 applied**.
- Owned P2: **1/1 accounted for; 1 applied**.
- P3 evidence-package findings: **11**.
- Non-card evidence qualification: **1**.
- Blocked findings: **0**.

## Applied corrections

1. **Bruce Lee — One-Inch Punch**
   - Restored the legal domain from opposing-only to any fighter adjacent to Bruce Lee.
   - Kept defeat-dependent card return behavior unchanged.

2. **Chupacabra — Natural Toughness**
   - Changed the return trigger from `cleanup_replacement` to printed `after_combat` timing.

3. **Loki — Looking for Trouble**
   - Added the printed gate requiring the opponent to have played a combat card.
   - Staged the flow as return current defense → inspect resulting hand → choose a legal replacement from that resulting hand → commit it into the same combat.
   - The returned defense is now in the hand choice domain and may be selected again when legal.
   - The Attack action is not restarted.

4. **Michelangelo — Shell Insertion ordering**
   - Removed the unsupported controller-selected permutation choice.
   - Selected effects resolve in printed/source order.

5. **Michelangelo — Shell Insertion action dependency**
   - Represented moving the selected hand card to the deck bottom as the cost for gaining the action.
   - Added an explicit dependency preventing the action if the cost is not paid/performed.

6. **Rosie the Riveter — D-Day**
   - Snapshots existing `count_active` before `activate_all`.
   - Uses the snapshot as a historical condition, not a scaling amount.
   - When all three upgrades were already active, applies fixed 1 damage and fixed 1 recovery.
   - No new D-REQ-012 query or transition was introduced.

7. **John Henry — Twelve-Pound Hammer**
   - Made the complete placement effect optional, preserving printed `may place` wording.

8. **Wyatt Earp — Gunfight at the O.K. Corral**
   - Added the fourth printed damage assignment.
   - Preserved the full controller/opponent/controller/opponent sequence and re-evaluation of current fighters after each stage.

9. **George Washington — Gather Information**
   - Narrowed movement from team-wide `each_friendly_fighter` to `controller_owned_fighters`, preserving printed `your fighters` ownership.

10. **Shredder — Back to Work!**
    - Undefeated Bebop & Rocksteady now recover 3 and are then placed in a closest legal space to Shredder.
    - Defeated behavior remains return with 3 health in a closest legal space.

## Already resolved

None of the ten owned material P1/P2 findings was already resolved at the correction branch base. All remained reproducible against `c154bf50acf259c95d0b7baf4877617556e38187` and received minimal corrections.

Clean owned manifests remained unchanged:

- `muhammad-ali`
- `blackbeard`
- `pandora`
- `leonardo`
- `donatello`
- `raphael`
- `krang` action-card manifest

## Evidence-package-only findings

No canonical card IDs were renamed to archive slugs. No clean manifest was changed for punctuation-only or package-only differences.

Evidence-package-only dispositions:

- Bruce Lee: `F-04`.
- Blackbeard: `B12-BB-P3-01`, `B12-BB-P3-02`.
- Loki: `LOKI-002`.
- Pandora: `PANDORA-001`.
- Leonardo: `LEONARDO-001`.
- Donatello: `DON-P3-01`.
- Michelangelo: `MIC-P3-01`, `MIC-P3-02`.
- Raphael: `RAP-P3-01`.
- Wyatt Earp: `WYATT-002`.

Krang's absent physical die remains an evidence/source-coverage qualification. The action-card manifest and fighter manifest continue to model the die as an external physical definition under `D-REQ-009`; no ordinary card definition was created or altered for the missing component.

## Integration requirements confirmed

Existing shared integration requirements were retained and were not edited:

- `D-REQ-001`: staged/bound resolution, replacement and source-order flows.
- `D-REQ-003`: restricted immediate attack action credits.
- `D-REQ-004`: dynamic Maneuver movement values.
- `D-REQ-005`: attack-network/reach behavior where already declared.
- `D-REQ-006`: temporary leave-space restrictions.
- `D-REQ-007`: path/token deployment state.
- `D-REQ-008`: Ruse token state.
- `D-REQ-009`: Die of Ultimate Destruction integration and rerolls.
- `D-REQ-010`: dynamic BOOST.
- `D-REQ-011`: face-up ranged commitment.
- `D-REQ-012`: upgrade activation/count state.
- `D-REQ-013`: Michelangelo hand-size dependency.
- `D-REQ-014`: card-play event counting.
- `D-REQ-015`: delayed turn obligations.
- `D-REQ-016`: Blackbeard ransom/shared treasury behavior.
- `D-REQ-017`: opaque opponent-hand instance selection.

Loki's fighter-level multiplayer non-combat TRICK recipient interpretation remains a medium-confidence `project_normalization`; `docs/fighters/phase-4b/loki.yaml` was not changed.

TMNT cooperative enemy/scenario logic was not modified. Only competitive Hero Deck card semantics in owned files were considered.

## Validation

- Branch ancestry: **PASS** — correction branch was created from current `phase-4b-worker-d-latest` tip `c154bf50acf259c95d0b7baf4877617556e38187`; no reset, rebase or merge.
- Current-state recheck: **PASS** — all ten material findings were re-read from the current correction branch after writes.
- YAML structural review: **PASS** — every changed file is readable through the GitHub contents endpoint; all changed YAML hunks were inspected for mapping/list indentation and flow-scalar balance. No GitHub Actions workflow is configured for this branch, so no CI run is claimed.
- Unique card IDs: **PASS** — no card definition ID was added, removed or renamed; no duplicate was introduced.
- Fixed deck totals: **PASS, 16/16** — all owned constructions remain fixed 30-card decks; quantities were unchanged and each manifest validation remains `quantity_sum: 30`.
- Card type/value/BOOST: **PASS** — no type, printed value or BOOST field was changed by the correction diff.
- Staged ordering: **PASS** — Loki replacement, Michelangelo source order and Wyatt alternating sequence explicitly preserve printed order.
- Cost dependencies: **PASS** — Shell Insertion action gain depends on successful bottom-card cost.
- Target ownership/domain: **PASS** — Bruce Lee and George Washington domains match printed relation/ownership wording.
- Placement selectors: **PASS** — Shredder uses `closest_legal_space_to_shredder` in both health branches.
- Optionality: **PASS** — Twelve-Pound Hammer whole effect is optional.
- Fixed versus scaled amounts: **PASS** — D-Day uses fixed 1/1 behind the pre-activation condition.
- `requires` references: **PASS** — existing top-level/shared references remain; no shared requirement document was modified.
- Loki unrelated interpretation: **PASS** — retained unchanged in fighter manifest.
- Krang die qualification: **PASS** — `evidence: qualified`, `source_coverage: QUALIFIED` and `D-REQ-009` remain intact.
- Foreign scope: **PASS** — pre-report branch diff contains only nine owned card manifests; no A/B/C manifest, report input, Phase 4A file, schema, mechanics, rules, rulings, set registry or cooperative scenario file changed.

Changed card manifests:

- `docs/cards/phase-4b/bruce-lee.yaml`
- `docs/cards/phase-4b/chupacabra.yaml`
- `docs/cards/phase-4b/george-washington.yaml`
- `docs/cards/phase-4b/john-henry.yaml`
- `docs/cards/phase-4b/loki.yaml`
- `docs/cards/phase-4b/michelangelo.yaml`
- `docs/cards/phase-4b/rosie-the-riveter.yaml`
- `docs/cards/phase-4b/shredder.yaml`
- `docs/cards/phase-4b/wyatt-earp.yaml`

No fighter manifest required correction.

## Remaining blockers

None.

Shared integration requirements remain explicit implementation dependencies, but they are not correction blockers and were not changed in this branch.

## Correction Handoff

Reports reviewed: 15/15  
Owned fighters searched: 16/16  
Owned P1 accounted for: 9/9  
Owned P2 accounted for: 1/1  
Applied: 10  
Already resolved: 0  
Evidence-only: 12 (11 P3 findings + 1 Krang non-card evidence qualification)  
Blocked: 0  
Base branch: `phase-4b-worker-d-latest`  
Correction branch: `phase-4-card-image-fixes-d`  
Base Head: `c154bf50acf259c95d0b7baf4877617556e38187`  
Final Head: current correction branch tip after persistence of this report; exact SHA is supplied by the external handoff  
Files changed: 10 (9 owned card manifests + this report)  
Merge performed: no
