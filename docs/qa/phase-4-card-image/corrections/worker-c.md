# Phase 4 Card-Image Corrections — Worker C

## Identity

- Repository: `NordCoder/unmatched-web`
- Canonical source branch: `phase-4b-worker-c-modern`
- Correction branch: `phase-4-card-image-fixes-c`
- Base Head: `4e521df76405e5ab65192eed057e43a80d561daa`
- Final Head: `8879b71ace3190987c286c45de91ba6fbd75edf1` (final canonical correction-content head before this report-persistence commit)
- Merge performed: no
- Source branch modified: no

Owned fighters:

`annie-christmas`, `dr-jill-trent`, `golden-bat`, `nikola-tesla`, `oda-nobunaga`, `tomoe-gozen`, `shakespeare`, `hamlet`, `titania`, `ciri`, `ancient-leshen`, `eredin`, `philippa`.

## Reports reviewed

Reports reviewed: **15/15**.

1. `alice.md`
2. `bigfoot.md`
3. `buffy.md`
4. `angel.md`
5. `achilles.md`
6. `bullseye.md`
7. `cloak-and-dagger.md`
8. `black-panther.md`
9. `annie-christmas.md`
10. `golden-bat.md`
11. `ciri.md`
12. `ancient-leshen.md`
13. `chupacabra.md`
14. `donatello.md`
15. `george-washington.md`

The non-Worker-C reports were reviewed for cross-owner findings. They contained no material finding for an owned fighter.

Owned fighters searched: **13/13**. Current card and fighter manifests were re-read from the correction branch before any disposition was assigned.

## Owned findings inventory

| Finding | Report | Fighter | Card | Severity | Classification | Disposition |
|---|---|---|---|---|---|---|
| `AC-P1-001` | `annie-christmas.md` | annie-christmas | Captain's Orders | P1 | canonical_manifest_error | Applied: both placement destination domains now require an empty space in Annie's zone and resolve in printed order. |
| `AC-P2-001` | `annie-christmas.md` | annie-christmas | Bottom Dealing | P2 | canonical_manifest_error | Applied: the publicly revealed card's top/bottom destination choice is public. |
| `JT-P1-001` | `annie-christmas.md` | dr-jill-trent | Utility Belt | P1 | canonical_manifest_error | Applied: Jill moves exactly 1 space, not up to 1. |
| `JT-P2-001` | `annie-christmas.md` | dr-jill-trent | Insightful Deduction | P2 | canonical_manifest_error | Applied: bottom-card and remaining-order choices preserve public reveal knowledge. |
| `SHAKESPEARE-001` | `golden-bat.md` | shakespeare | Et Tu, Brute? | P1 | canonical_manifest_error | Applied: both the adjacency count and Line-completion placement use `controller_fighters`, not team-wide `friendly_fighters`. |
| `SHAKESPEARE-002` | `golden-bat.md` | shakespeare | Horror | P1 | evidence_package_only | No Git correction: canonical `horror` quantity remains 5; five evidence entries are an archive construction defect. |
| `TITANIA-001` | `ciri.md` | titania | Gift Of The Fair Folk | P1 | canonical_manifest_error | Applied: movement mode binds exactly two distinct fighters, each moving up to 2 spaces. |
| `TITANIA-002` | `ciri.md` | titania | Glamour auxiliary deck | P2 | evidence_package_only | No Git deletion: six auxiliary Glamour definitions and their zones remain canonical. |
| `B12-PHILIPPA-001` | `ancient-leshen.md` | philippa | Do My Bidding | P1 | canonical_manifest_error | Applied: replacement selection is constrained to a card legally playable by the original attacking fighter in the current combat. |

Owned P1 accounted for: **6/6**.

Owned P2 accounted for: **3/3**.

## Applied corrections

Applied: **7** canonical corrections across five card manifests.

### Annie Christmas

- `Bottom Dealing`: changed the post-reveal destination choice from private to public.
- `Captain's Orders`: replaced unrestricted zone destinations with staged empty-space choices in Annie's zone; the second destination is evaluated after Annie's placement, then the action is gained.

### Dr. Jill Trent

- `Utility Belt`: changed Jill's movement branch from `max_distance: 1` to exact `distance: 1`.
- `Insightful Deduction`: changed both post-reveal choices to public visibility.

### Shakespeare

- `Et Tu, Brute?`: replaced all three `friendly_fighters` scopes with `controller_fighters`.
- `Horror` remains one canonical definition with `quantity: 5`.

### Titania

- `Gift Of The Fair Folk`: changed movement-target cardinality from `max: 2` to exact `cardinality: 2`.
- The action deck remains 30 cards and the six Glamours remain under `auxiliary_cards`.

### Philippa

- `Do My Bidding`: added original-attacker playability and current-combat legality to the replacement card filter.
- Ongoing-scheme structures were preserved unchanged.

## Already resolved

Already resolved: **0** owned P1/P2 findings.

The current source tip still contained each of the seven canonical discrepancies when rechecked. Clean or qualified corpora were not counted as “already resolved” findings:

- `golden-bat`: clean with integration qualification;
- `oda-nobunaga`: clean; controller-scoped `Battle Maneuvers` and intentionally friendly `Spring the Trap` preserved;
- `ciri`: no card transcription mismatch;
- `ancient-leshen`: clean with qualifications;
- `eredin`: no material card-image mismatch.

## Evidence-package-only findings

Evidence-only: **9** findings/qualifications. No canonical IDs or gameplay definitions were changed for these items.

- Annie Christmas — `captain-s-orders` archive slug versus canonical `captains-orders`.
- Nikola Tesla — `7-hertz` archive slug versus canonical `seven-hertz`.
- Tomoe Gozen — `a-warrior-s-way` archive slug versus canonical `a-warriors-way`.
- Tomoe Gozen — `lord-kiso-s-final-stand` archive slug versus canonical `lord-kisos-final-stand`.
- Shakespeare — five separate Horror evidence records versus one canonical definition with quantity 5.
- Hamlet — archive mapping/normalization qualification; `The Readiness Is All` project normalization was not rewritten as a printed mismatch.
- Titania — six Glamour images absent from the action-card evidence package; the auxiliary deck remains real canonical gameplay content.
- Eredin — `might-of-the-aen-elle` archive slug versus canonical `might-of-aen-elle`.
- Philippa — `spymaster-s-ruse` archive slug versus canonical `spymasters-ruse`.

## Integration requirements confirmed

No `C-REQ-*`, verification status, partial status, or integration status was changed.

Confirmed/preserved requirements across the owned manifests:

- Annie Christmas: `C-REQ-012`, `C-REQ-013`.
- Dr. Jill Trent: `C-REQ-013`.
- Golden Bat: `C-REQ-013`.
- Nikola Tesla: `C-REQ-004`, `C-REQ-013`.
- Oda Nobunaga: `C-REQ-003`, `C-REQ-010`, `C-REQ-013`.
- Tomoe Gozen: `C-REQ-003`, `C-REQ-013`.
- Shakespeare: `C-REQ-001`, `C-REQ-013`.
- Hamlet: `C-REQ-013`.
- Titania: `C-REQ-002`, `C-REQ-003`, `C-REQ-011`, `C-REQ-013`.
- Ciri: `C-REQ-005`, `C-REQ-006`, `C-REQ-013`.
- Ancient Leshen: `C-REQ-007`, `C-REQ-010`, `C-REQ-011`, `C-REQ-013`; the report's evidence qualification around `C-REQ-011` remains open rather than being promoted to full image confirmation.
- Eredin: `C-REQ-003`, `C-REQ-011`, `C-REQ-013`.
- Philippa: `C-REQ-008`, `C-REQ-011`, `C-REQ-013`.

## Validation

- YAML structural readback: **PASS** for all five modified manifests. Each file was re-fetched from the correction branch after persistence; modified indentation, inline collections, bindings, and tails were inspected. No branch CI/status check is configured for these push commits, so this is static manifest validation rather than a parser-backed CI claim.
- Unique canonical card IDs: **PASS**. No card ID was renamed or duplicated.
- Fixed deck counts: **PASS**.
  - Annie Christmas: 30.
  - Dr. Jill Trent: 30.
  - Shakespeare: 30, including `Horror quantity: 5`.
  - Titania action deck: 30; auxiliary Glamours: 6 definitions outside the action deck.
  - Philippa: 30.
- `usable_by`: **PASS**; no user restriction was broadened or removed.
- Target domains: **PASS** for Captain's Orders, Et Tu Brute, Gift Of The Fair Folk, and Do My Bidding.
- Staged bindings: **PASS**; all newly introduced destination and replacement filters resolve from prior-stage bindings.
- `requires` references: **PASS**; unchanged across all 13 owned fighter/card pairs.
- Source and fighter manifests: **PASS**; no fighter manifest required correction.
- Scope diff before report persistence: **PASS**. Only these five owned card files differed from the source branch:
  - `docs/cards/phase-4b/annie-christmas.yaml`
  - `docs/cards/phase-4b/dr-jill-trent.yaml`
  - `docs/cards/phase-4b/shakespeare.yaml`
  - `docs/cards/phase-4b/titania.yaml`
  - `docs/cards/phase-4b/philippa.yaml`
- Unrelated files changed: **none**.

## Remaining blockers

Blocked: **0** canonical correction findings.

Evidence-package corrections remain outside this branch's allowed scope. Integration requirements remain intentionally open and were not treated as transcription blockers.

## Correction Handoff

- Reports reviewed: **15/15**
- Owned fighters searched: **13/13**
- Owned P1 accounted for: **6/6**
- Owned P2 accounted for: **3/3**
- Applied: **7**
- Already resolved: **0**
- Evidence-only: **9**
- Blocked: **0**
- Base Head: `4e521df76405e5ab65192eed057e43a80d561daa`
- Final correction-content Head: `8879b71ace3190987c286c45de91ba6fbd75edf1`
- Files changed after report persistence: **6**
- Merge: **not performed**

Files changed:

1. `docs/cards/phase-4b/annie-christmas.yaml`
2. `docs/cards/phase-4b/dr-jill-trent.yaml`
3. `docs/cards/phase-4b/shakespeare.yaml`
4. `docs/cards/phase-4b/titania.yaml`
5. `docs/cards/phase-4b/philippa.yaml`
6. `docs/qa/phase-4-card-image/corrections/worker-c.md`
