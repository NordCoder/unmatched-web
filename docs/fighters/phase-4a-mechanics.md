# Phase 4A mechanics findings

**Scope:** Achilles, Bloody Mary, Sun Wukong, Sherlock Holmes, Dracula, Raptors, Wayward Sisters, Geralt of Rivia, Yennefer & Triss, Black Panther.

The representative corpus was chosen to stress different semantic dimensions rather than to maximize the number of cards transcribed.

## Result

The sample is expressible using the Phase 2 event/choice/effect framework plus a small set of **generic corpus-proven extensions**. No fighter requires an opaque character-specific script as the canonical rules model.

## Schema capabilities proven necessary

### 1. Topology and loss rules

The fighter model must distinguish:

- hero + sidekick(s): Achilles, Sherlock, Dracula, Geralt, Black Panther;
- multiple heroes with independent health: Raptors, Wayward Sisters;
- setup-selected hero/sidekick roles: Yennefer & Triss;
- summonable off-board sidekicks/reserve: Sun Wukong.

`fighter_count` alone cannot determine defeat/loss semantics. `loss_rule` is explicit.

### 2. Persisted historical state

Current board state is insufficient for effects such as:

- Bloody Mary's action ordinal;
- Bloody Mary's turn-start space;
- whether Geralt attacked this turn for ongoing-scheme expiry;
- captured parent-combat defense value for Bloody Reprise.

Historical snapshots are first-class authoritative state.

### 3. Available pool vs game deck

Geralt proves these are separate quantities:

```text
available pool = 36 action cards
base cards     = 24
game setup     = choose 1 Potion + 1 Armor + 1 Sword definition
copies chosen  = 2 of each
game deck      = 30
```

A fixed 30-card deck invariant is invalid globally.

### 4. External definitions

Not every gameplay card-like object is an action-card instance:

- Achilles / Bloody Mary / Sun Wukong bonus attacks are external combat definitions;
- Wayward Sisters spells are external effect definitions with ingredient requirements.

They have IDs and semantics but do not live in the normal action deck/hand/discard lifecycle.

### 5. First-class card zones and ownership

Wayward Sisters and Black Panther prove that arbitrary card lists attached to a hero are insufficient.

A card instance separates:

```text
immutable owner
current zone
zone controller / use authority
visibility
```

Examples:

- Cauldron: public alternate destination for the Wayward Sisters' own discarded action cards.
- Vibranium Suit: public storage for **opponent-owned** cards; Black Panther may use them only as BOOST sources, after which they return to their owner's discard.

### 6. Setup is resumable gameplay state

Geralt and Yennefer & Triss prove pre-game configuration must use the same persisted pending-choice infrastructure as normal play. A reconnect during gear selection or hero selection must reconstruct the same legal options and stage.

## Generic effect extensions proven by the sample

The following were added to `docs/mechanics/effect-model.md` because published cards require them.

| Semantic | Corpus evidence | Why existing primitive was insufficient |
| --- | --- | --- |
| `PREVENT_DAMAGE` | Sun Wukong — Bewilderment | prevention changes whether damage was dealt and downstream triggers |
| `REDIRECT_DAMAGE` | Sun Wukong — Golden Chain Mail | recipient changes while preserving the underlying damage event/type |
| `PREVENT_OPERATION` | Geralt — Tawny Owl | a replacement can stop a pending discard/effect operation rather than modify damage/value |
| `SET_PRINTED_VALUE` | Sherlock — Deduce Strategy | printed value and current combat value are distinct rule inputs |
| `ADD_BOOST_VALUE` | Yennefer & Triss — Advisor to the King | effective BOOST can be modified independently of base printed BOOST |
| `REORDER` | Wayward Sisters — Prophecy | looking at/reordering deck cards is not drawing/discarding |
| parent-context capture | Bloody Mary — Bloody Reprise | nested combat may depend on immutable data from the parent combat |
| `REPLACE_COMBAT_CARD` | Dracula — Do My Bidding | committed combat card identity changes inside one existing combat |

These are generic names. None is named after a fighter.

## Existing primitives proven sufficient for high-risk cases

### Combat participant replacement

Achilles' `Under Achilles' Helm` is representable as:

```text
SWAP(Achilles, Patroclus)
CHANGE_STATE(combat.defender, Patroclus)
```

No custom Achilles combat resolver is required.

### Multiple heroes

Raptor and Wayward Sister action cards use ordinary selectors over multiple hero instances. Multiple-hero semantics belong to topology/loss rules, not a separate card execution engine.

### Summoning

Sun Wukong uses the existing `SUMMON` composite plus a reserve token pool. `Sly Monkey` adds an explicit source-defined temporary occupancy override, not a universal exception to occupancy.

### Ongoing schemes

Geralt and Yennefer & Triss use a public `ongoing-schemes` zone whose resident cards install static/replacement/turn-end effects. This is ordinary source lifetime plus first-class zones.

### Spells

Wayward spells are `TRIGGER_COMPOSITE(kind=external_spell)` definitions. Ingredient requirements are card tags queried from the Cauldron; no parser for printed prose is required at runtime.

### Foreign cards

Black Panther uses ordinary card instances with immutable ownership and mutable zone/use authority. No copied/converted Black Panther card instance is created when an opponent card enters the Suit.

## Value model refined by the sample

The engine must keep at least these distinct values:

```text
printed_value_base        # immutable card-definition fact
effective_printed_value   # what "printed value" rules currently read
current_combat_value      # value used during combat resolution
boost_base                # immutable printed BOOST fact
effective_boost_value     # BOOST after static/source modifiers
```

This distinction is required by Sherlock's Deduce Strategy, Geralt's Sword of Silver, Yennefer & Triss' Paralyzing Fetters/Advisor to the King and other printed/current value interactions.

## Control-flow findings

### Bonus attacks

Bonus attacks remain nested combat inside one Attack action. The sample adds a requirement to carry source-defined captured context from the parent combat when needed.

### Combat-card replacement

`Do My Bidding` changes the attack card in the current combat. The replacement card does not create a second Attack action. Exact entry into remaining timing windows is source/ruling-defined and must be represented explicitly by the composite.

### Sequential repeated choices

Black Panther's `Wakanda Forever!` demonstrates that `up to two BOOSTs` cannot be modeled by locking both choices up front: his ability draws immediately after each BOOST value, changing the legal card pool before the second optional BOOST.

### End the turn

Geralt's Armor of the Forgotten Wolf reuses the existing `END_TURN` composite. If it fires during an opponent's Attack, the Phase 2 rule still redirects to Cleanup and skips unresolved post-Cleanup additional effects rather than behaving as `actions_remaining = 0`.

## No-DSL conclusion

Phase 4A supports the project's hybrid design direction:

- common semantics are structured data;
- composites encode global procedures with non-trivial control flow;
- source-defined overrides remain explicit;
- genuinely exceptional future mechanics may use a documented custom handler only after the corpus proves the generic model insufficient.

The representative corpus provides no evidence that a fully general embedded scripting language is necessary.