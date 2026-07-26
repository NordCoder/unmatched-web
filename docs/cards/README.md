# Card corpus

Action-card manifests are stored separately from fighter topology so deck construction, card instances, alternate card zones and external definitions can be validated independently.

## Phase 4A

Representative manifests are written under `phase-4a/` and follow `docs/fighters/schema.md`.

The corpus records:

- card name and stable ID;
- quantity;
- authorized fighter(s);
- card type;
- printed combat value, when present;
- BOOST value;
- mechanic tags such as ingredients or gear categories;
- normalized effect semantics;
- external definitions such as bonus attacks/spells;
- provenance.

## Source policy

Published UmDb paths under `https://unmatched.cards/umdb/...` are a normalized secondary index. Set rules, component rules and official rulings override it when they conflict.

Fan-deck paths under `https://unmatched.cards/decks/...` are explicitly excluded. This is particularly important for community balance-patch projects that preserve original card names while changing quantities, values or effects.

## Copyright boundary

This corpus does not reproduce full card prose. It stores structured gameplay facts and normalized operations needed by the rules engine. Original card art, trade dress and long-form printed wording are not implementation dependencies.