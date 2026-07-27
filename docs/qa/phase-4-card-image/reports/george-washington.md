# Phase 4 Card-Image QA Report

## Batch

- george-washington
- krang
- shredder
- wyatt-earp

## Verdict

**FAIL**

## QA Scope

Independent read-only card-corpus QA of the evidence batch containing `wyatt-earp`, `george-washington`, `shredder`, and `krang`.

Canonical repository: `NordCoder/unmatched-web`

Canonical branch used for all four fighters: `phase-4b-worker-d-latest`

Branch tip verified immediately before report persistence: `b9ae31c3b1a958e34bfb507d695cbd14650b9ed6`.

Canonical paths:

- fighter manifests: `docs/fighters/phase-4b/<fighter-id>.yaml`
- deck manifests: `docs/cards/phase-4b/<fighter-id>.yaml`

No canonical manifest, schema, source document, or gameplay file was modified during QA.

## Summary

The outer evidence archive contained exactly four fighter ZIPs: `wyatt-earp.zip`, `george-washington.zip`, `shredder.zip`, and `krang.zip`.

All four nested ZIPs opened successfully and passed archive integrity checks. Across the batch there were **48 unique card images**: 12 Wyatt Earp, 13 George Washington, 13 Shredder, and 10 Krang. All **48/48 images decoded and were visually inspected**. There were no zero-byte images, corrupt images, unreadable card faces, missing image files, extra image files, or duplicate binary images.

All four archive manifests reconciled to their canonical fixed deck construction counts. Each archive quantity sum was 30, matching Git `available_pool_count: 30` and `game_deck_count: 30`. No quantity failure was found.

Three gameplay-changing semantic discrepancies were found:

1. Wyatt Earp — `Gunfight at the O.K. Corral`: Git models only three of the four printed alternating damage assignments.
2. George Washington — `Gather Information`: printed `your fighters` is normalized as `each_friendly_fighter`, broadening the legal target domain in team play.
3. Shredder — `Back to Work!`: Git omits the printed closest-to-Shredder placement when Bebop & Rocksteady are still undefeated.

Three non-gameplay archive-to-Git card-ID mapping discrepancies were also found in Wyatt Earp evidence. Printed names make all three mappings unambiguous, so these are packaging/mapping issues rather than missing cards.

Krang's 10 accessible action-card images matched Git. Krang remains qualified because the physical Die of Ultimate Destruction was not included in this evidence ZIP; the card images confirm where die resolution is used but cannot independently prove the die face table or paid-reroll procedure.

Batch-level result: **FAIL** because at least one P1 discrepancy exists; specifically, three fighters contain P1 semantic mismatches.

---

## Fighter Results

### wyatt-earp

Branch: `phase-4b-worker-d-latest`

Deck manifest: `docs/cards/phase-4b/wyatt-earp.yaml`

Fighter manifest: `docs/fighters/phase-4b/wyatt-earp.yaml`

#### Archive integrity

**PASS**

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 12/12 |
| Unique images | 12 |
| Manifest card entries | 12 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |

#### Canonical manifest comparison / quantity validation

**PASS**

| Metric | Value |
|---|---:|
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 12 |
| Unique Git definitions | 12 |
| Construction | fixed |

The canonical fighter topology contains Wyatt Earp and Doc Holliday, so fighter restrictions referenced by the deck are resolvable.

#### Unique-card image completeness

**PASS — 12/12 unique card definitions have one corresponding readable image.**

Copies are represented through quantity rather than duplicate physical image files.

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Second Shot | PASS | PASS | PASS | PASS |
| You Just Gonna Stand There and Bleed? | PASS | PASS | PASS | PASS |
| In Vino Veritas | PASS | PASS | PASS | PASS |
| You Die First | PASS | PASS | PASS | PASS |
| I'm Your Huckleberry | PASS | P3 ID mapping | PASS | PASS* |
| I Have Two Guns, One for Each of You | PASS | PASS | PASS | PASS |
| Bring 'Em to Justice | PASS | PASS | PASS | PASS |
| Better Swear Me In | PASS | PASS | PASS | PASS |
| Fan the Hammer | PASS | PASS | PASS | PASS |
| You're a Daisy If You Do | PASS | P3 ID mapping | PASS | PASS* |
| A Marshal and an Outlaw | PASS | PASS | PASS | PASS |
| Gunfight at the O.K. Corral | PASS | P3 ID mapping | FAIL | FAIL |

`PASS*` means printed gameplay data matched Git, while the archive used a noncanonical `card_id` slug.

#### Discrepancies

##### WYATT-001

- fighter: `wyatt-earp`
- card: `Gunfight at the O.K. Corral`
- severity: **P1**
- Git location: `docs/cards/phase-4b/wyatt-earp.yaml`, `gunfight` effect
- expected / image fact: the printed AFTER COMBAT effect contains four sequential damage assignments: controller deals 1 damage to any fighter; opponent deals 1 damage to any fighter; controller deals 1 damage to any fighter; opponent deals 1 damage to any fighter.
- observed / Git representation: only three stages exist: `controller-first`, `opponent-second`, `controller-third`; the fourth opponent-controlled assignment is absent.
- evidence/reasoning: the missing fourth mandatory operation changes legal resolution and can change which fighters are defeated at the end of the effect. The printed wording establishes alternating controller authority across all four assignments.
- expected correction: preserve all four alternating stages; the fourth stage must be controlled by the combat opponent and select from the then-current legal fighter population after the previous three damage applications.

##### WYATT-002

- fighter: `wyatt-earp`
- cards: `I'm Your Huckleberry`, `You're a Daisy If You Do`, `Gunfight at the O.K. Corral`
- severity: **P3**
- expected: archive evidence should map to canonical Git card IDs, or explicitly declare aliases.
- observed:

| Archive `card_id` | Canonical Git ID |
|---|---|
| `i-m-your-huckleberry` | `im-your-huckleberry` |
| `you-re-a-daisy-if-you-do` | `youre-a-daisy-if-you-do` |
| `gunfight-at-the-o-k-corral` | `gunfight-at-the-ok-corral` |

- evidence/reasoning: printed names uniquely identify all three cards, so no gameplay definition or image is missing. The mismatch is evidence-package metadata only.
- expected correction: use canonical IDs or explicit alias mapping in future evidence packaging.

#### Integration requirements confirmed

- `D-REQ-001` — directly evidenced by staged/alternating choice behavior, especially `Gunfight at the O.K. Corral`.
- `D-REQ-003` — directly evidenced by restricted immediate attack actions on `In Vino Veritas` and `You Die First`.

#### Verdict

**FAIL**

Reason: `Gunfight at the O.K. Corral` loses one printed mandatory damage operation in Git.

---

### george-washington

Branch: `phase-4b-worker-d-latest`

Deck manifest: `docs/cards/phase-4b/george-washington.yaml`

Fighter manifest: `docs/fighters/phase-4b/george-washington.yaml`

#### Archive integrity

**PASS**

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13/13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |

#### Canonical manifest comparison / quantity validation

**PASS**

| Metric | Value |
|---|---:|
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 13 |
| Unique Git definitions | 13 |
| Construction | fixed |

The canonical fighter topology contains George Washington plus three distinct Culper Spy fighter instances.

#### Unique-card image completeness

**PASS — 13/13 unique card definitions have one corresponding readable image.**

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Feint | PASS | PASS | PASS | PASS |
| Sabotage | PASS | PASS | PASS | PASS |
| Circumvent | PASS | PASS | PASS | PASS |
| Lead Astray | PASS | PASS | PASS | PASS |
| Make a Stand | PASS | PASS | PASS | PASS |
| Network of Spies | PASS | PASS | PASS | PASS |
| Gather Information | PASS | PASS | FAIL | FAIL |
| Agent 711 | PASS | PASS | PASS | PASS |
| Sympathetic Stain | PASS | PASS | PASS | PASS |
| Allies Everywhere | PASS | PASS | PASS | PASS |
| Misinformation | PASS | PASS | PASS | PASS |
| Undercover Agent | PASS | PASS | PASS | PASS |
| Recruit to the Ring | PASS | PASS | PASS | PASS |

#### Discrepancies

##### GEORGE-001

- fighter: `george-washington`
- card: `Gather Information`
- severity: **P1**
- Git location: `docs/cards/phase-4b/george-washington.yaml`, `gather` scheme effect
- expected / image fact: printed effect instructs the player to move each of **your fighters** up to 3 spaces, with permission for those fighters to move through opposing fighters.
- observed / Git representation: movement target is `each_friendly_fighter`.
- evidence/reasoning: `your fighters` and `friendly fighters` are different target domains in team play. `friendly` can include teammate-controlled fighters, while the printed scheme restricts movement to fighters controlled by the resolving player. The same deck separately uses genuinely friendly-fighter wording on `Network of Spies`, confirming that this distinction is gameplay-significant rather than stylistic.
- expected correction: target only the controller's own fighters; retain maximum distance 3 and the printed permission to move through opposing fighters.

#### Integration requirements confirmed

- `D-REQ-001` — staged choice/rebinding flows such as `Undercover Agent`, `Recruit to the Ring`, and `Sympathetic Stain`.
- `D-REQ-003` — `Allies Everywhere` grants a restricted free attack action.
- `D-REQ-008` — Ruse lifecycle/effect interactions are directly referenced by the card corpus.

#### Verdict

**FAIL**

Reason: `Gather Information` broadens a printed own-fighter target domain to friendly fighters.

---

### shredder

Branch: `phase-4b-worker-d-latest`

Deck manifest: `docs/cards/phase-4b/shredder.yaml`

Fighter manifest: `docs/fighters/phase-4b/shredder.yaml`

#### Archive integrity

**PASS**

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 13/13 |
| Unique images | 13 |
| Manifest card entries | 13 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |

#### Canonical manifest comparison / quantity validation

**PASS**

| Metric | Value |
|---|---:|
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 13 |
| Unique Git definitions | 13 |
| Construction | fixed |

The canonical fighter topology correctly treats Bebop & Rocksteady as one independently damageable 7-health sidekick.

#### Unique-card image completeness

**PASS — 13/13 unique card definitions have one corresponding readable image.**

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Long Shot | PASS | PASS | PASS | PASS |
| Think Hard | PASS | PASS | PASS | PASS |
| Gruff Escort | PASS | PASS | PASS | PASS |
| Gang Up | PASS | PASS | PASS | PASS |
| Back to Work! | PASS | PASS | FAIL | FAIL |
| Masterful Defense | PASS | PASS | PASS | PASS |
| Savagery | PASS | PASS | PASS | PASS |
| Perplexing Tactics | PASS | PASS | PASS | PASS |
| All According to Plan | PASS | PASS | PASS | PASS |
| Disrupting Strike | PASS | PASS | PASS | PASS |
| Obedient Subjects | PASS | PASS | PASS | PASS |
| Swarming Strike | PASS | PASS | PASS | PASS |
| Master of the Foot | PASS | PASS | PASS | PASS |

#### Discrepancies

##### SHREDDER-001

- fighter: `shredder`
- card: `Back to Work!`
- severity: **P1**
- Git location: `docs/cards/phase-4b/shredder.yaml`, `back-to-work` effect
- expected / image fact: the printed scheme has two sequential instructions: Bebop & Rocksteady recover 3 health, **even if defeated**; then **place them in a space as close as possible to Shredder**. The placement instruction is not conditioned on whether they were defeated before resolution.
- observed / Git representation: the undefeated branch only performs `RECOVER` for 3. The defeated branch performs `RETURN_FIGHTER` with health 3 and closest-legal-space placement.
- evidence/reasoning: when Bebop & Rocksteady are alive, Git heals them but leaves them in their original space, omitting a mandatory printed placement. This can materially change adjacency, attack reach, zone relationships, and later effects.
- expected correction: preserve the alive/defeated health distinction but resolve closest-to-Shredder placement in both cases: undefeated → recover 3 then place as close as possible; defeated → return with 3 health in a closest legal space.

#### Integration requirements confirmed

- `D-REQ-001` — staged selection/ordered resolution, notably `Gruff Escort`.
- `D-REQ-007` — Foot Soldier deployment/state is directly referenced by multiple card faces.
- `D-REQ-015` — `All According to Plan` directly evidences a delayed end-of-turn obligation tied to the opponent.

`D-REQ-005` is not claimed as visually confirmed by this evidence batch because its core attack-network behavior is primarily Shredder's fighter ability rather than printed on the included action-card faces.

#### Verdict

**FAIL**

Reason: `Back to Work!` omits mandatory placement when Bebop & Rocksteady are undefeated.

---

### krang

Branch: `phase-4b-worker-d-latest`

Deck manifest: `docs/cards/phase-4b/krang.yaml`

Fighter manifest: `docs/fighters/phase-4b/krang.yaml`

#### Archive integrity

**PASS**

| Check | Result |
|---|---:|
| Nested ZIP readable | PASS |
| Manifest readable | PASS |
| Images readable | 10/10 |
| Unique images | 10 |
| Manifest card entries | 10 |
| Duplicate binary images | 0 |
| Missing images | 0 |
| Extra images | 0 |

#### Canonical manifest comparison / quantity validation

**PASS**

| Metric | Value |
|---|---:|
| Archive quantity sum | 30 |
| Git available pool | 30 |
| Git game deck | 30 |
| Unique archive cards | 10 |
| Unique Git action-card definitions | 10 |
| Construction | fixed |

The separate `ultimate-destruction-resolution` Git entry is an external gameplay definition rather than an ordinary action-card copy and was therefore correctly excluded from the action-card image count.

#### Unique-card image completeness

**PASS — 10/10 unique action-card definitions have one corresponding readable image.**

#### Printed card identity/content verification

| Card | Image | Metadata | Printed semantics | Result |
|---|---|---|---|---|
| Warlord of Dimension X | PASS | PASS | PASS | PASS |
| IQ of 968 | PASS | PASS | PASS | PASS |
| Pan-Dimensional Portal | PASS | PASS | PASS | PASS |
| Welcome to the Technodrome! | PASS | PASS | PASS | PASS |
| Minimizer | PASS | PASS | PASS | PASS |
| Android Arms: Wings | PASS | PASS | PASS | PASS |
| Android Arms: Missiles | PASS | PASS | PASS | PASS |
| Android Arms: Powerbomb | PASS | PASS | PASS | PASS |
| Android Arms: Chain Flail | PASS | PASS | PASS | PASS |
| Molecular Amplification Unit | PASS | PASS | PASS | PASS |

Additional checked details:

- `Molecular Amplification Unit` quantity is correctly 4.
- `IQ of 968` and `Pan-Dimensional Portal` use dynamic printed BOOST handling.
- `Android Arms: Missiles` correctly supports a face-up ranged commitment mode.
- `Android Arms: Powerbomb` keys post-combat movement distance to combat damage actually dealt.
- `Pan-Dimensional Portal` preserves PLACE semantics, a temporary cannot-leave-space restriction, one random-resolution result applied to both combat players' hand-bottom operation, then both combat players draw 2.
- `IQ of 968` preserves two independent Ultimate Destruction symbol resolutions for draw and recover.
- `Android Arms: Wings` preserves movement → damage to opposing fighters moved through → gain 1 action ordering.

#### Discrepancies

**None in the accessible action-card images.**

#### Qualification / warning

The ZIP contains the 10 action-card faces but not the physical Die of Ultimate Destruction. Therefore the card faces directly confirm the existence and placement of Ultimate Destruction resolution calls, dynamic BOOST usage, and related card semantics, but they do not independently prove the physical die face table or paid-reroll procedure.

Git already treats the die as an external physical definition and marks Krang evidence/source coverage as qualified. This is an evidence qualification, not a transcription discrepancy in the accessible card corpus.

#### Integration requirements confirmed

Card images directly confirm:

- `D-REQ-001` — staged/bound resolution.
- `D-REQ-006` — temporary cannot-leave-space state on `Pan-Dimensional Portal`.
- `D-REQ-009` — Die of Ultimate Destruction random-resolution integration.
- `D-REQ-010` — dynamic/non-numeric BOOST.
- `D-REQ-011` — face-up ranged attack commitment on `Android Arms: Missiles`.

#### Verdict

**PASS_WITH_QUALIFICATIONS**

All accessible action-card images match canonical Git card data and normalized gameplay behavior; qualification is limited to the absent physical die component.

---

## Findings

### Material findings

| ID | Severity | Fighter | Card | Finding |
|---|---|---|---|---|
| `WYATT-001` | P1 | wyatt-earp | Gunfight at the O.K. Corral | Git contains only 3 of the 4 printed alternating 1-damage assignments. |
| `GEORGE-001` | P1 | george-washington | Gather Information | Printed `your fighters` was broadened to `each_friendly_fighter`. |
| `SHREDDER-001` | P1 | shredder | Back to Work! | Git omits mandatory closest-to-Shredder placement when Bebop & Rocksteady are alive. |

No P2 findings were identified.

### Non-material / metadata findings

| ID | Severity | Fighter | Scope | Finding |
|---|---|---|---|---|
| `WYATT-002` | P3 | wyatt-earp | archive mapping | Three archive `card_id` slugs differ from canonical Git IDs, but printed names make all mappings unambiguous. |

### Quantity findings

None. All four decks reconcile to 30 action-card instances according to both archive quantity sums and canonical Git construction data.

### Image evidence findings

- unreadable images: 0
- corrupt images: 0
- missing images: 0
- extra images: 0
- binary duplicate images: 0
- successfully inspected unique card images: 48/48

---

## Corpus-Level Observations

1. The evidence packaging correctly follows the one-unique-definition-to-one-image model. No repeated physical copies were represented as duplicate image files; quantity remained in each archive manifest.
2. Archive manifests were useful as image-to-card/quantity mappings but were not treated as authoritative gameplay data. This distinction exposed the three Wyatt Earp slug differences without incorrectly treating them as new card definitions.
3. Fixed deck size happened to be 30 for all four fighters in this batch, but validation was performed against each Git construction record rather than assuming a universal 30-card rule.
4. External gameplay definitions were not counted as ordinary deck card images. Krang's Ultimate Destruction resolution/die contract is the relevant example.
5. Normalized Git structures were evaluated semantically rather than through literal prose equality. Staged effects, bindings, state operations, and shared integration requirements were accepted where they preserved printed legal gameplay behavior.
6. The three P1 findings are not formatting or wording differences. Each changes legal gameplay behavior: an omitted mandatory damage assignment, an expanded target domain, and an omitted mandatory placement.
7. Integration requirements were not treated as transcription failures when the printed behavior was correctly represented and Git explicitly deferred runtime support through `D-REQ-*` contracts.
8. Krang demonstrates the evidence-boundary rule: all included action cards can pass while a non-card physical component remains qualified because it was not present in the ZIP.

---

## Final Assessment

This batch is **technically complete and fully inspectable**, with 48/48 unique card images successfully reviewed and no archive corruption, missing images, duplicates, or quantity mismatches.

It is **not clean enough to accept as a fully verified Phase 4 card-image corpus without corrections**, because three canonical deck semantics conflict with the physical card evidence:

- Wyatt Earp / `Gunfight at the O.K. Corral` — missing fourth damage assignment;
- George Washington / `Gather Information` — wrong target domain;
- Shredder / `Back to Work!` — missing placement in the undefeated branch.

Krang's accessible action-card corpus is clean, subject only to the documented qualification that the physical Die of Ultimate Destruction was not included in this evidence archive.

Final batch verdict: **FAIL**.
