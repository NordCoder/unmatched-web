# Sherwood Forest validation

## Machine contract

Authoritative runtime data:

```text
internal/playableslice/content/data/sherwood-forest.json
```

Byte-identical acceptance fixture:

```text
tests/fixtures/battlefields/sherwood-forest.json
```

Human-readable mirror:

```text
docs/battlefields/sherwood-forest.yaml
```

Run:

```bash
python scripts/validate_playable_slice.py
python -m unittest tests/architecture/test_playable_slice.py
```

The validator requires the exact launch contract, not merely a connected graph:
30 stable spaces, 7 registered zones, seat 1 at `s20`, seat 2 at `s19`, and the
canonical set of 39 undirected edges. It also rejects unknown endpoints,
self-edges, duplicate reverse edges, invalid zone references, shifted start
markers and disconnected graphs.

## Runtime use

The typed game kernel derives:

- movement paths from undirected adjacency;
- melee legality from direct adjacency;
- ranged legality from shared zone membership;
- initial hero and sidekick placement from the two start markers and graph
  distance.

Artwork coordinates are used only by the browser SVG.

## Current verdict

```text
space inventory: pass
zone registry and membership: pass
starting markers: pass
unknown endpoint validation: pass
duplicate-edge validation: pass
connectivity: pass
runtime/fixture identity: pass
source metadata: pass
independent visual QA: pending final playable-slice review
```
