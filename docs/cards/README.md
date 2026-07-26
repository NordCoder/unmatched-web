# Card corpus

Action-card manifests are stored separately from fighter topology so deck construction, card instances, alternate card zones and external definitions can be validated independently.

## Phase 4A

**Status:** complete — gate passed 2026-07-26.

Full representative manifests under `phase-4a/`:

- `achilles.yaml`;
- `bloody-mary.yaml`;
- `sun-wukong.yaml`;
- `sherlock-holmes.yaml`;
- `dracula.yaml`;
- `raptors.yaml`;
- `wayward-sisters.yaml`;
- `geralt-of-rivia.yaml`;
- `yennefer-triss.yaml`;
- `black-panther.yaml`.

They follow `docs/fighters/schema.md` and are validated in `docs/fighters/phase-4a-validation.md`.

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

## Phase 4B

The same schema must now be expanded across the remaining released fighter roster from `docs/sets/registry.yaml`. New semantic operations are allowed only when the full corpus proves the existing generic vocabulary insufficient.

## Source policy

Published UmDb paths under `https://unmatched.cards/umdb/...` are a normalized secondary index. Set rules, component rules and official rulings override it when they conflict.

Fan-deck paths under `https://unmatched.cards/decks/...` are explicitly excluded. This is particularly important for community balance-patch projects that preserve original card names while changing quantities, values or effects.

## Copyright boundary

This corpus does not reproduce full card prose. It stores structured gameplay facts and normalized operations needed by the rules engine. Original card art, trade dress and long-form printed wording are not implementation dependencies.