# Source policy

## Purpose

`unmatched-web` needs rules precise enough for deterministic software. A convenient database entry is not automatically authoritative, and an old printed rulebook may be superseded by errata.

Every normative rule, non-obvious card interpretation, fighter special rule, and battlefield exception must retain provenance.

## Authority hierarchy

Use the highest applicable source, while respecting explicit later corrections:

1. **Official publisher/designer errata, addenda and rulings** — corrections and authoritative clarifications.
2. **Current official Core Rules** — general competitive rules.
3. **Official Set Rules / official rulebooks** — set, fighter and battlefield special rules.
4. **Published component content** — character cards, action cards and battlefield/component markings.
5. **UmDb published database** (`unmatched.cards/umdb/...`) — normalized published fighter/deck data and discovery.
6. **Unmatched Reference** — consolidated index of special rules/errata/mechanics; useful for discovery and cross-checking, but some entries contain editorial synthesis.
7. **The Unmatched Club / similar community indexes** — useful indexes to official rulings and current roster information.
8. **Other community sources** — discovery only unless independently verified above.

A lower-ranked source may expose information that a higher-ranked source does not contain. In that case it can be used with explicit provenance and confidence, but it must not silently override a higher-ranked source.

## Specific rules

### Current rules beat old generic wording

Older set rulebooks remain important for their set-specific content, but generic rules should be reconciled against the current Core Rules and later official errata.

### Card/special effects may override general rules

When an effect explicitly contradicts a general rule, model the effect as an override. Do not “correct” the card back to the general rule.

### Errata and addenda are first-class records

Do not edit history away. Store:

- original source/edition;
- correcting source;
- effective normalized rule;
- date/version when known.

Example: a later battlefield addendum should be linked as a distinct source rather than pretending the original set rulebook contained the missing rule.

### Rulings are not card text

A ruling should be represented separately and linked to the affected rule/card/fighter/battlefield. This preserves the distinction between published component semantics and later clarification.

### Fan decks are forbidden in the official corpus

On unmatched.cards:

- allowed discovery namespace: `https://unmatched.cards/umdb/...`;
- fan-deck namespace: `https://unmatched.cards/decks/...`.

The latter must never be imported as an official fighter merely because search results look plausible.

### Missing data is a blocker, not permission to invent

Use a structured status such as:

```yaml
status: blocked
reason: public_normalized_card_corpus_incomplete
sources_checked:
  - ...
```

An implementation assumption may be proposed later, but it must be labeled as an assumption and must not masquerade as an official rule.

## Provenance fields

Normative records should support at least:

```yaml
sources:
  - authority: official_core | official_set | official_component | official_ruling | secondary_database | reference
    url: https://...
    title: ...
    version: ...
    checked_at: YYYY-MM-DD
    note: ...
```

For derived engine semantics, also record why the normalization is equivalent to the source rule.

## Copyright and asset boundary

Documentation should capture facts, structured gameplay data, rule semantics, short necessary references, and our own implementation-oriented paraphrases. Do not make the repository depend on wholesale reproduction of copyrighted rulebook prose, card artwork, graphic layouts, or licensed media assets.

This is both cleaner engineering and important for licensed properties such as Marvel, The Witcher, Jurassic Park, Buffy, and TMNT.
