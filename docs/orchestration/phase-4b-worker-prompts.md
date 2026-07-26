# Phase 4B worker prompts

All four workers use the same exact orchestration base:

`4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`

Branches:

- Worker A: `phase-4b-worker-a-classics`
- Worker B: `phase-4b-worker-b-licensed`
- Worker C: `phase-4b-worker-c-modern`
- Worker D: `phase-4b-worker-d-latest`

Each branch was created from the exact base above.

---

## Worker A prompt

Work as **Phase 4B Fighter/Deck Research Worker A** for `NordCoder/unmatched-web`.

Operating identity:

- Repository: `NordCoder/unmatched-web`
- Phase: `4B`
- Branch: `phase-4b-worker-a-classics`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: early/classic + retired licensed fighters
- Assigned fighters: 17

First, through the GitHub Connector, read:

1. `docs/orchestration/phase-4b-parallel-plan.md`
2. `docs/fighters/schema.md`
3. `docs/fighters/phase-4a-validation.md`
4. `docs/mechanics/effect-model.md`
5. `docs/sources/source-policy.md`
6. `docs/sets/registry.yaml`
7. representative Phase 4A fighter/card manifests as examples.

Confirm the branch is based on the Authorized Base before writing.

Your assigned fighter IDs are ONLY:

- `alice`
- `king-arthur`
- `medusa`
- `sinbad`
- `robin-hood`
- `bigfoot`
- `robert-muldoon`
- `invisible-man`
- `jekyll-and-hyde`
- `buffy`
- `willow`
- `spike`
- `angel`
- `little-red-riding-hood`
- `beowulf`
- `deadpool`
- `yennenga`

Do not rewrite Phase 4A fighters such as Raptors, Dracula, Sherlock Holmes, Achilles, Bloody Mary or Sun Wukong.

For every assigned fighter, create exactly the necessary manifests under:

- `docs/fighters/phase-4b/<fighter-id>.yaml`
- `docs/cards/phase-4b/<fighter-id>.yaml`

Follow the verified Phase 4A schema. Capture topology, movement, health, attack types, setup, ability, resources/state, deck construction, every unique published action card, quantity, user, type, printed value, BOOST, normalized effects, external definitions, rulings and provenance.

Source hierarchy:

1. official Restoration/IELLO rulebooks, set rules, addenda/errata and official rulings;
2. current Rulings Archive;
3. published UmDb paths under `/umdb/...` for normalized published deck facts;
4. current secondary indexes only for discovery/cross-checking.

Never use `unmatched.cards/decks/...` fan decks or balance patches as canonical data.

Do not copy long card prose. Store factual card metadata plus normalized semantics.

CRITICAL PARALLELISM RULE:
You may NOT edit shared files such as `docs/fighters/schema.md`, `docs/mechanics/**`, `docs/rules/**`, `docs/rulings/ambiguity-register.md`, `docs/sets/**`, `docs/research-plan.md`, or top-level READMEs.

If a card cannot be represented by the current vocabulary, do not invent a local opaque hero handler and do not patch `effect-model.md`. Record a reusable extension proposal in your worker report with the affected card, source, semantic requirement and proposed generic primitive/composite/state model. Mark affected content blocked if integration is required before it can be deterministic.

Validation per fighter:

- quantities reconcile with construction;
- every `usable_by` target exists;
- referenced resources/zones/state exist;
- hidden/public information is explicit;
- operations map to current vocabulary or an explicit extension proposal;
- known authoritative rulings are attached;
- no fan-patch data is mixed in.

Create `docs/phase-4b/worker-a-report.md` with:

`## Worker 4B-A Handoff`

and include branch, Authorized Base, exact final Head, assigned count, verified fighters, blocked fighters, quantity validation, schema-extension proposals, new ambiguities/blockers, source gaps and files created.

Do not merge to `main`. Do not modify another worker's files. Finish with the exact branch Head and handoff status.

---

## Worker B prompt

Work as **Phase 4B Fighter/Deck Research Worker B** for `NordCoder/unmatched-web`.

Operating identity:

- Repository: `NordCoder/unmatched-web`
- Phase: `4B`
- Branch: `phase-4b-worker-b-licensed`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: 2022–2023 licensed / Marvel-heavy corpus
- Assigned fighters: 18

Read the same authoritative repository contracts first:

1. `docs/orchestration/phase-4b-parallel-plan.md`
2. `docs/fighters/schema.md`
3. `docs/fighters/phase-4a-validation.md`
4. `docs/mechanics/effect-model.md`
5. `docs/sources/source-policy.md`
6. `docs/sets/registry.yaml`
7. representative Phase 4A manifests, especially Black Panther for Marvel card-zone semantics.

Confirm the branch is based on the Authorized Base before writing.

Assigned fighter IDs ONLY:

- `daredevil`
- `elektra`
- `bullseye`
- `ghost-rider`
- `luke-cage`
- `moon-knight`
- `dr-ellie-sattler`
- `t-rex`
- `houdini`
- `genie`
- `cloak-and-dagger`
- `ms-marvel`
- `squirrel-girl`
- `black-widow`
- `winter-soldier`
- `doctor-strange`
- `she-hulk`
- `spider-man`

Black Panther is already Phase 4A and must not be rewritten.

For each fighter create:

- `docs/fighters/phase-4b/<fighter-id>.yaml`
- `docs/cards/phase-4b/<fighter-id>.yaml`

Use official set/product/rulebook material and authoritative rulings to verify special resources, non-standard deck sizes, battlefield-item-facing card semantics, multiple heroes, summoned pieces, transformations/resources, large-fighter interactions and any non-standard card zones.

Use published UmDb `/umdb/...` for normalized published deck facts but never community `/decks/...` patches.

Do not copy long printed card text; normalize gameplay behavior.

Do NOT edit shared schema/mechanics/rules/rulings/set/roadmap files. Any required general semantic extension belongs only in `docs/phase-4b/worker-b-report.md` as a proposal for the orchestrator.

For every fighter validate quantities, `usable_by`, resources/zones/state, hidden/public information, effect vocabulary coverage, rulings and provenance. An exact `blocked` result is preferable to invented behavior.

Create `docs/phase-4b/worker-b-report.md` with `## Worker 4B-B Handoff` and include branch, Authorized Base, exact final Head, verified/blocked fighters, quantity checks, proposed reusable extensions, ambiguities, source gaps and created files.

Do not merge to `main` and do not touch another worker's paths.

---

## Worker C prompt

Work as **Phase 4B Fighter/Deck Research Worker C** for `NordCoder/unmatched-web`.

Operating identity:

- Repository: `NordCoder/unmatched-web`
- Phase: `4B`
- Branch: `phase-4b-worker-c-modern`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: mechanic-heavy modern sets
- Assigned fighters: 13

Read before writing:

1. `docs/orchestration/phase-4b-parallel-plan.md`
2. `docs/fighters/schema.md`
3. `docs/fighters/phase-4a-validation.md`
4. `docs/mechanics/effect-model.md`
5. `docs/mechanics/setup-hooks.md`
6. `docs/sources/source-policy.md`
7. `docs/sets/registry.yaml`
8. Phase 4A manifests for Wayward Sisters, Geralt and Yennefer & Triss.

Confirm branch/base identity first.

Assigned fighter IDs ONLY:

- `annie-christmas`
- `dr-jill-trent`
- `golden-bat`
- `nikola-tesla`
- `oda-nobunaga`
- `tomoe-gozen`
- `shakespeare`
- `hamlet`
- `titania`
- `ciri`
- `ancient-leshen`
- `eredin`
- `philippa`

Do not rewrite Wayward Sisters, Geralt or Yennefer & Triss.

Tales to Amaze scope is **competitive heroes only**. Do not model Mothman, Martian Invader, minions, initiative, threat/scenario objectives or other co-op enemy logic.

Create fighter/card manifests only under the assigned Phase 4B paths.

Pay particular attention to:

- multi-health sidekicks;
- token/resource trackers;
- transformation/state machines;
- card-title/sequencing systems;
- glamour/external card pools;
- ongoing schemes;
- summonable sidekicks;
- role/state-dependent card behavior;
- setup hooks and persistent state.

Source hierarchy and fan-deck prohibition are identical to the other workers.

Do NOT edit shared schema/mechanics/rules/rulings/set/roadmap files. Because this is the highest mechanic-density batch, explicitly record every apparent missing primitive in `docs/phase-4b/worker-c-report.md`; do not implement it globally yourself.

Validate all deck quantities and semantic references. Preserve uncertainty as `blocked` rather than guessing.

Create `docs/phase-4b/worker-c-report.md` with `## Worker 4B-C Handoff`, branch/base/head, verified/blocked list, quantity validation, extension proposals, ambiguities, source gaps and file list.

Do not merge to `main`.

---

## Worker D prompt

Work as **Phase 4B Fighter/Deck Research Worker D** for `NordCoder/unmatched-web`.

Operating identity:

- Repository: `NordCoder/unmatched-web`
- Phase: `4B`
- Branch: `phase-4b-worker-d-latest`
- Authorized Base: `4d259bc02a28d764b23ee1e6c50ebbad4f947ba9`
- Scope: newest / freshness-sensitive 2025–2026 content
- Assigned fighters: 16

Read before writing:

1. `docs/orchestration/phase-4b-parallel-plan.md`
2. `docs/fighters/schema.md`
3. `docs/fighters/phase-4a-validation.md`
4. `docs/mechanics/effect-model.md`
5. `docs/sources/source-policy.md`
6. `docs/sets/registry.yaml`
7. `docs/sets/source-bibliography.md`
8. current Phase 4A examples.

Confirm branch/base identity first.

Assigned fighter IDs ONLY:

- `bruce-lee`
- `muhammad-ali`
- `blackbeard`
- `chupacabra`
- `loki`
- `pandora`
- `leonardo`
- `donatello`
- `michelangelo`
- `raphael`
- `rosie-the-riveter`
- `john-henry`
- `wyatt-earp`
- `george-washington`
- `shredder`
- `krang`

Bruce Lee is one canonical fighter lineage spanning the historical standalone release and Bruce Lee vs. Muhammad Ali. Do not create duplicate canonical fighter IDs.

TMNT scope is competitive heroes/decks only. Do not substitute Shredder/Krang villain behavior for their official competitive hero decks and do not model Adventures co-op enemy/scenario logic.

Stars & Stripes is freshness-sensitive. Use current authoritative official/physical material where public normalized databases lag. Never reconstruct missing card data from guesses, previews of uncertain status or fan decks. A precise blocked record with the missing authoritative evidence is correct.

Create only:

- `docs/fighters/phase-4b/<fighter-id>.yaml`
- `docs/cards/phase-4b/<fighter-id>.yaml`
- `docs/phase-4b/worker-d-report.md`

Do NOT edit shared schema/mechanics/rules/rulings/set/roadmap files. Record proposed generic extensions in the worker report for orchestrator reconciliation.

Validate quantities, users, resources/zones/state, hidden information, effect vocabulary, rulings and provenance for every fighter.

Create `docs/phase-4b/worker-d-report.md` with `## Worker 4B-D Handoff`, branch/base/head, verified/blocked fighters, quantity validation, extension proposals, ambiguities, source gaps and files created.

Do not merge to `main`.