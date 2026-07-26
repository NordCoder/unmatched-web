# Phase 3 validation

**Phase:** 3 — Canonical set and release registry  
**Validation date:** 2026-07-26  
**Result:** PASS for the declared competitive-content registry scope  
**Overall project readiness:** NOT YET `developer-ready`; fighter/card and battlefield-detail phases remain.

## Gate under test

Phase 3 requires:

1. every officially released competitive-compatible set/product to have a canonical record;
2. standalone fighters and competitive heroes originating from Adventures products to be included;
3. every released competitive fighter and battlefield to have a registry membership path;
4. set-level special mechanics to have at least one authoritative official entry point;
5. addenda/errata/reprints to preserve provenance instead of silently replacing history;
6. announced but unreleased/incomplete products to be explicitly blocked rather than guessed.

## Released primary-product coverage

The registry contains 25 released primary product records:

| # | Canonical ID | Release lineage | Fighters | Battlefields |
| ---: | --- | --- | --- | --- |
| 1 | `battle-of-legends-vol-1` | 2019 | Alice, King Arthur, Medusa, Sinbad | Marmoreal, Sarpedon |
| 2 | `robin-hood-vs-bigfoot` | 2019; Vault return 2026 | Robin Hood, Bigfoot | Sherwood Forest, Yukon |
| 3 | `bruce-lee-solo` | 2019 historical single | Bruce Lee | — |
| 4 | `jurassic-park-ingen-vs-raptors` | 2020 | Robert Muldoon/InGen, Raptors | Raptor Paddock |
| 5 | `cobble-and-fog` | 2020; clarified Vault reprint 2025 | Dracula, Invisible Man, Jekyll & Hyde, Sherlock Holmes | Baskerville Manor, Soho |
| 6 | `buffy-the-vampire-slayer` | 2020 | Buffy, Willow, Spike, Angel | Sunnydale High School, The Bronze |
| 7 | `little-red-vs-beowulf` | 2020 | Little Red, Beowulf | Heorot |
| 8 | `deadpool-solo` | 2021 historical single | Deadpool | — |
| 9 | `battle-of-legends-vol-2` | 2021 | Achilles, Bloody Mary, Sun Wukong, Yennenga | Hanging Gardens |
| 10 | `marvel-hells-kitchen` | 2022 | Daredevil, Elektra, Bullseye | Hell's Kitchen |
| 11 | `marvel-redemption-row` | 2022 | Ghost Rider, Luke Cage, Moon Knight | The Raft |
| 12 | `jurassic-park-sattler-vs-trex` | 2022 | Dr. Ellie Sattler, T. Rex | T. Rex Paddock |
| 13 | `houdini-vs-genie` | 2022 | Houdini, Genie | King Solomon's Mine |
| 14 | `marvel-teen-spirit` | 2023 | Cloak & Dagger, Ms. Marvel, Squirrel Girl | Navy Pier |
| 15 | `marvel-for-king-and-country` | 2023 | Black Panther, Black Widow, Winter Soldier | Helicarrier |
| 16 | `marvel-brains-and-brawn` | 2023 | Doctor Strange, She-Hulk, Spider-Man | Sanctum Sanctorum |
| 17 | `tales-to-amaze` | 2023 | Annie Christmas, Dr. Jill Trent, Golden Bat, Nikola Tesla | Point Pleasant, McMinnville |
| 18 | `suns-origin` | 2024 | Oda Nobunaga, Tomoe Gozen | Azuchi Castle |
| 19 | `slings-and-arrows` | 2024 | Shakespeare, Hamlet, Titania, Wayward Sisters | Globe Theatre |
| 20 | `witcher-steel-and-silver` | first public sale 2024 | Geralt, Ciri, Ancient Leshen | Kaer Morhen, Fayrlund Forest |
| 21 | `witcher-realms-fall` | first public sale 2024 | Yennefer & Triss, Eredin, Philippa | Naglfar, Streets of Novigrad |
| 22 | `battle-of-legends-vol-3` | 2025 | Blackbeard, Chupacabra, Loki, Pandora | Santa's Workshop, Venice |
| 23 | `bruce-lee-vs-muhammad-ali` | 2025 | Bruce Lee, Muhammad Ali | Thrilla in Manila, Tsing Shan Monastery |
| 24 | `tmnt-adventures` | 2025 crowdfunding lineage; general sale 2026 | Leonardo, Donatello, Michelangelo, Raphael | New York City, Technodrome |
| 25 | `stars-and-stripes` | 2026 | Rosie the Riveter, John Henry, Wyatt Earp, George Washington | White House, Alamo |

Bruce Lee appears in two product memberships but is one canonical fighter lineage. This prevents content duplication while preserving release provenance.

### External count reconciliation

Restoration Games stated on 2025-11-18 that 24 Unmatched sets had been published. The first 24 records above represent the release lineage through the 2025 catalog/crowdfunding period. Stars & Stripes was subsequently released in 2026 and is explicitly described by Restoration as its latest set release in June 2026.

This reconciliation is used only as an omission detector; the project does not depend on publisher marketing terminology for its internal IDs.

## Released supplements that add gameplay content

Two official products outside the normal primary-set model are required for competitive completeness:

| ID | Why it belongs in the gameplay corpus | Status |
| --- | --- | --- |
| `tmnt-shredder-krang-hero-decks` | Restoration explicitly enables Shredder and Krang as competitive heroes | verified |
| `nova-high-battlefield-mat` | official current Unmatched battlefield product | existence verified; graph intentionally deferred to Phase 5 |

Cosmetic-only foils, alternate art, alternate miniatures, sleeves, storage products and neoprene replacements of already registered boards do not create gameplay-content records.

## Announced-content containment

`hellboy` is present but `blocked`.

Official evidence currently establishes:

- a licensed Hellboy Unmatched product exists/has been announced;
- the announcement described a four-character box;
- the current product page exposes no final playable roster, rulebook, battlefield list or card corpus.

Therefore no community-predicted fighter names are entered into `fighters`, and Hellboy does not count as released-content coverage.

**Result:** PASS — future content is visible without contaminating the released corpus.

## Battlefield coverage

Every battlefield attached to the 25 primary product records is enumerated, plus the separate Nova High official battlefield accessory.

Special battlefield families are independently indexed in `mechanics-index.md`:

- one-way paths;
- secret passages;
- doors;
- high ground;
- battlefield items;
- Witcher multi-level adjacency;
- Santa's Workshop conveyor behavior.

Phase 5 remains responsible for proving every node, edge, zone, starting space and special connection. Phase 3 only proves the battlefield is not missing and the special-rule authority is discoverable.

## Set-mechanic authority coverage

The following high-risk set mechanics were traced to official product/rulebook/addendum material:

| Family | Authority checkpoint |
| --- | --- |
| Raptors multiple heroes, Muldoon traps, Raptor Paddock one-way paths | official InGen vs Raptors rulebook |
| Baskerville Manor Secret Passages / Invisible Man components | official Cobble & Fog rulebook |
| Buffy sidekick-selection setup | official Buffy rulebook |
| Little Red basket, Beowulf Rage, Heorot doors | official Little Red vs Beowulf rulebook |
| Hanging Gardens high ground / BoL2 fighter special components | official BoL2 rulebook/product |
| Marvel battlefield items and non-uniform deck/component sizes | official Marvel product/rulebook material |
| T. Rex large fighter / Sattler insight system | official Jurassic Park product/rulebook material |
| Oda independent-health sidekicks | official Sun's Origin product/Set Rules |
| Slings & Arrows multi-hero/cauldron/glamour/sequencing systems | official Slings & Arrows Set Rules |
| Geralt gear, Witcher ongoing schemes, Kaer Morhen one-way paths | official Steel & Silver Set Rules |
| Yennefer/Triss hero selection and Realms Fall map mechanic | official Realms Fall Set Rules |
| BoL3 special hero/map mechanics | official BoL3 product + current consolidated reference |
| TMNT competitive hero compatibility | official TMNT product/release material |
| Shredder/Krang competitive conversion | official Hero Decks product page |
| Stars & Stripes fighter systems and White House Secret Passages | official product/release material + separate addendum |

**Result:** PASS — no known set-specific mechanic in the Phase 3 inventory depends solely on an unsupported fan interpretation.

## Provenance / edition validation

### Cobble & Fog

The reprint is represented as an edition of the same canonical set. Restoration's Vault announcement states that the returned edition receives small rulebook clarifications but is the same game. Historical wording remains traceable.

### Robin Hood vs. Bigfoot

The 2026 return from the Vault does not create a duplicate fighter/map set. It is an edition event on `robin-hood-vs-bigfoot`.

### Bruce Lee

The original 2019 standalone release and the 2025 Lee vs. Ali membership are both retained. Restoration explicitly states that the returning Bruce deck could not change; Phase 4 must still compare exact published manifests before declaring byte/text identity.

### Stars & Stripes

The White House Secret Passages omission is represented as a separate addendum source. The project must never claim the original rulebook contained the omitted rule.

**Result:** PASS.

## Data gaps that do not fail Phase 3

| Gap | Containment |
| --- | --- |
| Complete Stars & Stripes normalized card database may lag physical release | Phase 4 must transcribe/verify against authoritative card material; no card text is invented here |
| Nova High exact graph/equivalence | Phase 5 blocker for supporting that battlefield; existence is already registered |
| Retired Deadpool current official product page/rulebook not publicly indexed | release membership is cross-checked with publisher-version metadata and UmDb; Phase 4 must independently verify its cards |
| Hellboy playable content not published | record remains `announced/blocked`; outside released roster |
| precise release-channel year/date for some crowdfunded/late-year products | edition notes preserve known channel facts; irrelevant to game-state semantics |

## Gate decision

**PASS — Phase 3 is complete for the competitive-content registry.**

The project now has a deterministic answer to:

- which released primary products exist;
- which official supplemental fighter/battlefield products add gameplay content;
- which fighters belong to which releases;
- which battlefields belong to which releases;
- which set-level mechanics require special authority;
- where the authoritative rulebook/product/addendum entry point is;
- which announced/incomplete content must remain blocked.

Phase 4 may now use `registry.yaml` as the exhaustive fighter/deck work queue. Phase 5 may use its battlefield memberships as the exhaustive graph work queue.

This does **not** make those fighters or battlefields `developer-ready` until their own later-phase gates pass.