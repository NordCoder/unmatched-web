# Setup extension hooks

Published fighters increasingly modify setup before the first turn: fighter-role selection, deck construction, equipment, forms, placement choices, and post-placement state. The engine needs stable setup stages rather than fighter-name conditionals.

## Canonical generalized order

`HOOK-001` — Use the current generalized setup order established by the official rulings/reference when set/fighter rules do not explicitly override it:

```text
1. determine player order
2. choose battlefield
3. choose characters
4. resolve character pre-game configuration in player order
5. construct/shuffle decks and draw starting hands
6. place heroes in player order
7. place sidekicks / remaining starting fighters as instructed
8. resolve post-placement character setup
9. begin first turn
```

Source: REF10 Major Rulings: Setup Order, plus CORE standard setup.

`HOOK-002` — Tournament/draft procedures may replace the first character/battlefield selection stages without changing the later gameplay setup hooks unless their rules explicitly do so.  
Source: REF10 Setup Order notes; source-specific precedence.

## Hook stages

The documentation model defines the following semantic hooks:

```yaml
setup_hooks:
  - PRE_CHARACTER_SELECTION
  - CHARACTER_CONFIGURATION
  - DECK_CONSTRUCTION
  - BEFORE_STARTING_HAND
  - HERO_PLACEMENT
  - SIDEKICK_PLACEMENT
  - POST_PLACEMENT_CONFIGURATION
  - BEFORE_FIRST_TURN
```

`HOOK-010` — A fighter/set manifest declares only the hooks it actually uses and must cite the authoritative rule/card that defines each hook.  
Source: normalized extension requirement.

`HOOK-011` — A hook may create one or more pending player choices and therefore uses the same `CHOICE-*` pause/resume semantics as in-game choices. Setup is not allowed to be an unpersisted client-only wizard.  
Source: normalized online-play requirement.

`HOOK-012` — If multiple players have character-configuration choices in the same generalized stage, resolve them in the authoritative player order unless a set-specific rule states otherwise.  
Source: REF10 Major Rulings: Setup Order.

## Character configuration

`HOOK-020` — Choices that determine the character's identity/topology for the match occur before ordinary deck shuffle/draw when the authoritative setup rule requires it. Examples include choosing which fighter is the hero or selecting a legal pre-game package.  
Source: REF10 Setup Order examples (Buffy, Geralt, Yennefer & Triss).

`HOOK-021` — Pre-game deck construction must finish before shuffle and starting-hand draw. The resulting legal game deck is then the deck used for ordinary game rules.  
Source: REF10 Setup Order; Witcher set rules for Geralt as representative content.

`HOOK-022` — The engine must store both the **available construction pool/rules** and the **chosen game configuration** when a fighter supports pre-game construction. Do not overwrite the source definition with only the chosen 30-card runtime deck.  
Source: normalized content/provenance requirement.

## Placement hooks

`HOOK-030` — Standard hero placement uses the numbered starting-space procedure from Phase 1 unless a fighter/set rule replaces it.  
Source: CORE setup; `docs/rules/setup.md`.

`HOOK-031` — Sidekick/other starting-fighter placement follows the source-defined legal spaces and ordering after hero placement under the generalized setup order.  
Source: CORE setup; REF10 Setup Order.

`HOOK-032` — A character's post-placement configuration occurs only after the placements it depends on. Alice's starting size is the canonical generalized-order example.  
Source: REF10 Major Rulings: Setup Order.

## Setup state and auditability

`HOOK-040` — Setup decisions that affect later legality must persist in `GameState` or equivalent immutable match configuration: selected hero role, selected equipment/cards, initial form/state, starting placement, and any consumed setup-only components.  
Source: normalized deterministic-engine requirement.

`HOOK-041` — Secret setup choices remain private until the relevant source says they become public; public setup choices are recorded in public event history.  
Source: information-visibility semantics + source-specific rules.

## No fighter-name branching

The setup engine should dispatch by declared hook definitions/capabilities, not by hard-coded names such as:

```text
if fighter == Geralt ...
if fighter == Alice ...
```

A genuinely unique setup mechanic may use a custom hook handler, but its inputs, outputs, choices, visibility, ordering and provenance must still be documented explicitly.
