# Set-specific mechanics index

**Last verified:** 2026-07-26

This index records mechanics that are not adequately described by the generic competitive Core Rules alone. It is a **discovery and provenance layer**: detailed fighter semantics belong in Phase 4 and exact battlefield graph semantics belong in Phase 5.

A mechanic is listed here only when the project has an authoritative official entry point. Secondary databases may be used to confirm which product/map contains it, but do not define its executable semantics.

## Battlefield mechanics

| Mechanic ID | Products / battlefield | Phase 3 authoritative source | Later work |
| --- | --- | --- | --- |
| `one-way-paths` | Jurassic Park InGen vs Raptors — Raptor Paddock | official InGen vs Raptors rulebook | Phase 5 encode directional traversal while preserving adjacency semantics |
| `secret-passages` | Cobble & Fog — Baskerville Manor | official Cobble & Fog rulebook | Phase 5 encode passage endpoints as movement-only special connections |
| `doors` | Little Red vs Beowulf — Heorot | official Little Red vs Beowulf rulebook | Phase 5 encode door state and its effect on edges/zones |
| `high-ground` | Battle of Legends Vol. 2 — Hanging Gardens | official BoL2 product/rulebook | Phase 5 encode high-ground markings and combat modifier conditions |
| `battlefield-items` | Marvel Hell's Kitchen / Redemption Row / Teen Spirit / For King and Country / Brains and Brawn | official Marvel product rulebooks/pages; REF10 battlefield-items reference | Phase 5 encode item slots/components; Phase 4 encode attached card-effect semantics |
| `secret-passages` | Houdini vs Genie — King Solomon's Mine | official set rulebook | Phase 5 encode passage graph |
| `one-way-paths` | Witcher Steel & Silver — Kaer Morhen | official Steel & Silver Set Rules | Phase 5 encode directional traversal |
| `multi-level-battlefield-adjacency` | Witcher Realms Fall — Naglfar | official Realms Fall Set Rules | Phase 5 encode above/below-deck adjacency relationship exactly |
| `one-way-paths` / `conveyor-belt` | Battle of Legends Vol. 3 — Santa's Workshop | official BoL3 material / REF10 | Phase 5 formal graph and conveyor behavior |
| `secret-passages` | Stars & Stripes — White House | **official 2026 Stars & Stripes addendum** | Phase 5 encode passage endpoints; addendum must not be lost behind original rulebook |

### White House provenance warning

The Secret Passages rules were omitted from the Stars & Stripes rulebook. Restoration Games published a separate addendum:

`https://restorationgames.com/wp-content/uploads/2026/03/UM-SnS-rules-addendum_color.pdf`

The addendum establishes that passage spaces may be traversed between each other for movement at a cost of one space, are **not** generally adjacent for attacks/effects, and cannot be traversed by large figures. This source has precedence over the omission in the original rulebook.

## Fighter / setup mechanics exposed at set level

The table intentionally records mechanic **families**, not complete fighter rules.

| Mechanic ID | Product / fighters | Official entry point | Phase 4 requirement |
| --- | --- | --- | --- |
| `pre-game-sidekick-choice` | Buffy — Buffy chooses Giles or Xander package | official Buffy rulebook | transcribe exact deck-construction/setup effect |
| `multiple-heroes` | Jurassic Park — Raptors | official InGen vs Raptors rulebook | model each Raptor as a hero and verify defeat/winner rulings |
| `traps` | Jurassic Park — Muldoon | official InGen vs Raptors rulebook | token placement/trigger/removal semantics |
| `fighter-resource` | Beowulf — Rage | official Little Red vs Beowulf rulebook | resource gain/spend semantics |
| `basket-state` | Little Red | official Little Red vs Beowulf rulebook | discard-top-derived state |
| `summonable-sidekicks` | Sun Wukong | official BoL2 material | defeated/off-board clone summon lifecycle |
| `battlefield-items` | Marvel fighters using item tokens | official Marvel set rulebooks / REF10 | attach selected item effect to combat card and cancellation scope |
| `nonstandard-deck-size` | Daredevil, Elektra, Black Widow and other explicit cases | official product component lists / rulebooks | deck construction must be data-driven, never `30` hard-coded |
| `large-fighter` | T. Rex | official Jurassic Park Sattler vs T. Rex material | footprint/orientation/movement interactions |
| `insight-tokens` | Dr. Sattler | official Jurassic Park Sattler vs T. Rex material | token generation/use |
| `multi-health-sidekicks` | Oda Nobunaga | official Sun's Origin product/Set Rules | independently tracked sidekick health |
| `multiple-heroes` / `ingredient-spell-system` | Wayward Sisters | official Slings & Arrows Set Rules | multiple heroes + cauldron/ingredient/spell state |
| `card-title-sequencing` | Shakespeare | official Slings & Arrows Set Rules | ordered title/syllable sequence state |
| `fighter-state-choice` | Hamlet | official Slings & Arrows Set Rules | persistent/selected state behavior |
| `glamour-deck` | Titania | official Slings & Arrows Set Rules | external/auxiliary card pool lifecycle |
| `pre-game-deck-construction` | Geralt | official Witcher Steel & Silver Set Rules | gear selection before shuffle/draw |
| `ongoing-schemes` | Witcher sets | official Witcher Set Rules | persistent scheme location/lifecycle |
| `setup-time-hero-selection` | Yennefer & Triss | official Witcher Realms Fall Set Rules | one becomes hero, other sidekick at setup |
| `external-card-pool` | Pandora | official BoL3 material / REF10 | Pandora's Box auxiliary card lifecycle |
| `fighter-resource` | Blackbeard | official BoL3 material / REF10 | doubloon/treasure resource lifecycle |
| `competitive-heroes-in-adventures-product` | Tales to Amaze and TMNT Turtles | official product pages | transcribe only hero-side competitive semantics in Phase 4; co-op enemy engine remains deferred |
| `competitive-conversion-of-adventures-villains` | Shredder & Krang Hero Decks | official Restoration accessory product page | treat as distinct competitive fighter definitions, not the co-op villain AI decks |
| `track-network` | John Henry | official Stars & Stripes product/release material | exact railroad-token path behavior from cards/rules |
| `mech-upgrades` | Rosie the Riveter | official Stars & Stripes product/release material | upgrade state/selection semantics |
| `ruse-network` | George Washington | official Stars & Stripes product/release material | exact ruse/private-information semantics |

## Adventures boundary

`Tales to Amaze` and `TMNT Adventures` contain both competitive-compatible heroes/battlefields and a separate cooperative enemy/scenario system.

Phase 3 registers both products, but the first implementation scope takes only:

- the competitive hero definitions;
- battlefield data required for standard Unmatched where supported;
- shared components required by those heroes.

The following remain Phase 8:

- villain/minion AI decks;
- initiative decks;
- threat/scenario objectives;
- cooperative targeting/priority rules;
- scenario-specific win/loss logic.

The Shredder & Krang Hero Decks are different: Restoration explicitly publishes them so those villains can be used as competitive heroes. They therefore belong in the competitive registry now.

## Not mechanics yet

The registry does not infer gameplay semantics from component names alone. For example:

- Nova High is a verified official battlefield accessory, but no graph equivalence is assumed until Phase 5 authoritative map inspection;
- a token appearing in a component list does not define its lifecycle;
- Hellboy has no mechanic entries until its authoritative rules/content are published.