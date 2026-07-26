# Phase 4B parallel orchestration plan

**Phase:** 4B — Complete fighter/deck corpus  
**Orchestration model:** 4 independent research/transcription workers + one integration/orchestration pass  
**Prepared:** 2026-07-26

## Goal

Expand the verified Phase 4A fighter/card schema to every released competitive fighter in `docs/sets/registry.yaml` without allowing parallel workers to race on shared semantic files.

Phase 4A already covers:

- Achilles;
- Bloody Mary;
- Sun Wukong;
- Sherlock Holmes;
- Dracula;
- Raptors;
- Wayward Sisters;
- Geralt of Rivia;
- Yennefer & Triss;
- Black Panther.

After removing those ten from the released registry, Phase 4B has **64 unique fighter IDs** remaining.

## Parallelism principle

Workers are partitioned primarily by **set/source family**, not alphabetically. Keeping a complete set with one worker avoids repeated rulebook/set-rule interpretation and reduces inconsistent normalization between fighters from the same product.

Workload is intentionally not 16/16/16/16 by fighter count. Older fighters are generally cheaper to normalize, while newer Witcher/Slings/TMNT/Stars & Stripes content has more setup, resource, state and ruling pressure.

The resulting split is:

| Worker | Scope profile | Fighters |
| --- | --- | ---: |
| A | early/classic + retired licensed | 17 |
| B | 2022–2023 licensed/Marvel | 18 |
| C | Tales / Sun's Origin / Slings / Witcher | 13 |
| D | 2025–2026 / newest & freshness-sensitive | 16 |

## Mandatory branch isolation

Every worker starts from the same orchestration base and writes only to its own branch.

Workers MUST NOT edit shared semantic/control files:

- `docs/fighters/schema.md`;
- `docs/mechanics/**`;
- `docs/rules/**`;
- `docs/rulings/ambiguity-register.md`;
- `docs/sets/**`;
- `docs/research-plan.md`;
- `docs/README.md`;
- `README.md`.

A worker may discover that the current schema/effect vocabulary is insufficient. It must **not** patch the shared model itself. Instead it records a proposed reusable extension in its worker report with:

1. affected fighter/card;
2. authoritative source;
3. why existing semantics are insufficient;
4. proposed generic primitive/composite/state extension;
5. whether the fighter is blocked until integration.

The orchestrator owns all shared-schema changes after comparing all four reports, preventing four workers from inventing overlapping primitives independently.

## Worker-owned paths

Each worker may create only:

```text
docs/fighters/phase-4b/<assigned-fighter-id>.yaml
docs/cards/phase-4b/<assigned-fighter-id>.yaml
docs/phase-4b/worker-<a|b|c|d>-report.md
```

No fighter ID is assigned to more than one worker.

## Worker A — early/classic + retired licensed

Assigned fighters (17):

### Battle of Legends, Volume One
- `alice`
- `king-arthur`
- `medusa`
- `sinbad`

### Robin Hood vs. Bigfoot
- `robin-hood`
- `bigfoot`

### Jurassic Park — InGen vs. Raptors
- `robert-muldoon`

`raptors` is already Phase 4A and must not be rewritten.

### Cobble & Fog
- `invisible-man`
- `jekyll-and-hyde`

`dracula` and `sherlock-holmes` are already Phase 4A.

### Buffy the Vampire Slayer
- `buffy`
- `willow`
- `spike`
- `angel`

### Little Red Riding Hood vs. Beowulf
- `little-red-riding-hood`
- `beowulf`

### Deadpool standalone
- `deadpool`

### Battle of Legends, Volume Two
- `yennenga`

Achilles, Bloody Mary and Sun Wukong are already Phase 4A.

## Worker B — 2022–2023 licensed/Marvel

Assigned fighters (18):

### Marvel — Hell's Kitchen
- `daredevil`
- `elektra`
- `bullseye`

### Marvel — Redemption Row
- `ghost-rider`
- `luke-cage`
- `moon-knight`

### Jurassic Park — Dr. Sattler vs. T. Rex
- `dr-ellie-sattler`
- `t-rex`

### Houdini vs. The Genie
- `houdini`
- `genie`

### Marvel — Teen Spirit
- `cloak-and-dagger`
- `ms-marvel`
- `squirrel-girl`

### Marvel — For King and Country
- `black-widow`
- `winter-soldier`

`black-panther` is already Phase 4A.

### Marvel — Brains and Brawn
- `doctor-strange`
- `she-hulk`
- `spider-man`

## Worker C — Tales / Sun's Origin / Slings / Witcher

Assigned fighters (13), intentionally smaller because this partition is mechanic-heavy:

### Unmatched Adventures: Tales to Amaze — competitive heroes only
- `annie-christmas`
- `dr-jill-trent`
- `golden-bat`
- `nikola-tesla`

Do not model cooperative villains/minions/initiative/scenario logic.

### Sun's Origin
- `oda-nobunaga`
- `tomoe-gozen`

### Slings and Arrows
- `shakespeare`
- `hamlet`
- `titania`

`wayward-sisters` is already Phase 4A.

### The Witcher — Steel & Silver
- `ciri`
- `ancient-leshen`

`geralt-of-rivia` is already Phase 4A.

### The Witcher — Realms Fall
- `eredin`
- `philippa`

`yennefer-and-triss` is already Phase 4A.

## Worker D — 2025–2026 / freshness-sensitive

Assigned fighters (16):

### Bruce Lee lineage / Bruce Lee vs. Muhammad Ali
- `bruce-lee`
- `muhammad-ali`

Bruce Lee must preserve the 2019 standalone → 2025 set membership lineage. Do not create two canonical fighter manifests.

### Battle of Legends, Volume Three
- `blackbeard`
- `chupacabra`
- `loki`
- `pandora`

### Unmatched Adventures: Teenage Mutant Ninja Turtles — competitive heroes
- `leonardo`
- `donatello`
- `michelangelo`
- `raphael`

Do not model cooperative enemy/scenario behavior.

### Stars and Stripes
- `rosie-the-riveter`
- `john-henry`
- `wyatt-earp`
- `george-washington`

This set is freshness-sensitive. The worker must not fabricate card data when current public normalized databases are incomplete; use authoritative physical/official material where available and mark unresolved records `blocked` with exact missing evidence.

### TMNT Shredder & Krang Hero Decks
- `shredder`
- `krang`

Only their official **competitive hero** rules/decks belong here. Do not substitute villain logic from Adventures.

## Per-fighter required output

Every assigned fighter must have a fighter manifest and deck manifest following `docs/fighters/schema.md`.

Fighter manifest must include:

- canonical ID and set membership;
- movement;
- topology and loss rule;
- every independently tracked fighter/sidekick with attack type and health semantics;
- setup hooks;
- ability;
- resources/tokens/custom card zones;
- persistent historical/state fields;
- deck construction;
- applicable authoritative rulings and provenance.

Deck manifest must include:

- available pool count and game deck count;
- every unique published action card;
- quantity;
- `usable_by`;
- card type;
- printed combat value where applicable;
- BOOST value;
- tags/resources/categories;
- normalized effects;
- external gameplay definitions where required;
- quantity reconciliation;
- sources and normalization notes.

## Source discipline

Use the hierarchy already established by the repository:

1. official Restoration/IELLO rules, set rules, addenda, errata and official rulings;
2. current official Rulings Archive;
3. published UmDb `/umdb/...` for normalized deck facts;
4. current secondary ruling indexes only for discovery/cross-checking.

Never import `unmatched.cards/decks/...` fan/community decks or balance patches.

Do not copy large blocks of card/rulebook prose. Store factual metadata plus normalized gameplay semantics and provenance.

## Validation required from each worker

For every fighter:

- quantities reconcile with deck construction;
- every `usable_by` fighter exists;
- all referenced resources/zones/state are defined;
- all card operations map to the current Phase 4A vocabulary or are explicitly proposed as a reusable extension in the report;
- known fighter/card rulings are attached;
- hidden/public information semantics are explicit where relevant;
- no fan-patch data is mixed into the official corpus.

The worker report must contain:

```text
## Worker 4B-X Handoff
Branch: ...
Base: ...
Head: ...
Assigned fighters: N
Verified: ...
Blocked: ...
Quantity validation: PASS/FAIL
Schema-extension proposals: ...
New ambiguity/blockers: ...
Source gaps: ...
Files created: ...
```

A worker must not mark an uncertain fighter `verified` merely to complete its batch. `blocked` with a precise evidence gap is a valid and preferred outcome.

## Integration phase owned by orchestrator

After all four handoffs:

1. verify branch scope and file ownership;
2. review every proposed schema/effect extension across all workers together;
3. deduplicate equivalent extension proposals;
4. update shared `schema.md`, `effect-model.md` and ambiguity/ruling documents once;
5. integrate fighter/card manifests;
6. run global uniqueness and quantity checks;
7. verify every released fighter in `sets/registry.yaml` is either Phase 4A, verified Phase 4B, or explicitly blocked;
8. only then mark Phase 4B complete.

Workers do not close the Phase 4B gate independently.