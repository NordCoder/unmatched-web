# Robin Hood vs Bigfoot playable slice

## Outcome

The local application supports one fixed two-player game:

```text
Robin Hood vs Bigfoot
battlefield: Sherwood Forest
```

One user creates a match in the first browser tab, copies the match code, and
joins as Bigfoot in a second tab. Each tab stores an independent opaque session
token in `sessionStorage` and polls the in-memory host for its player projection.

## Run locally

```bash
go run ./apps/server
```

Open `http://localhost:8080` in two separate tabs. In the first tab choose
**Create match**. In the second tab enter the displayed match code and choose
**Join as Bigfoot**.

A custom port may be supplied:

```bash
PORT=8090 go run ./apps/server
```

The approved PostgreSQL Core adapter remains in the repository and is unchanged.
This slice deliberately uses a process-local host because restart recovery,
reconnect hardening and production transport are deferred by Issue #51.

## Architecture

```text
verified fighter/card YAML manifests
        ↓ embedded source data
internal/playableslice/content
        ↓ closed typed definitions and semantic effect kinds
internal/playableslice/game
        ↓ private per-player projection
internal/playableslice/server
        ↓ HTTP + polling
apps/web
```

`internal/playableslice/content` is the only production registration boundary
that maps launch fighter/card IDs to typed behavior. The game kernel switches on
closed mechanics such as `cancel_opposing`, `regroup`, `horns` and
`momentous_shift`; it contains no Robin Hood, Bigfoot or concrete card-ID
branches. No generic JSON expression language is used.

## Ordinary rules implemented

- deterministic setup, shuffle and seat-1 starting player;
- hero and sidekick placement;
- health, deck, private hand, discard and exhaustion;
- two actions per turn;
- Maneuver draw, optional single-card BOOST and movement of each controlled
  fighter using the boosted movement value;
- Scheme actions and typed scheme procedures;
- melee adjacency and ranged shared-zone legality;
- hidden attack and optional defense selection followed by reveal;
- defender-first resolution within simultaneous timing windows;
- immediate, during-combat and after-combat typed effects used by both decks;
- effect damage, draw, discard, movement, placement and sidekick return;
- hand-limit discard prompts;
- hero defeat and winner declaration.

## Launch content coverage

The content registry embeds the existing verified manifests and requires exactly
30 cards for each deck.

Robin Hood coverage:

```text
A Hunter's Eye
Steal from the Rich
Disarming Shot
Piercing Shot
Highway Robbery
Defenders of Sherwood
Feint
Regroup
Wily Fighting
Snark
Ambush
post-attack retreat ability
```

Bigfoot coverage:

```text
Larger Than Life
Savagery
Crash Through the Trees
Jackalope Horns
Hoax
Disengage
It's Just Your Imagination
Feint
Regroup
Skirmish
Momentous Shift
isolated end-turn draw ability
```

## Privacy and authority

- opaque random tokens bind one browser session to one match and player;
- commands are authorized from the token, not a client-supplied player ID;
- only the viewing player receives their hand identities;
- only the choice owner receives legal prompt options;
- attack and defense identities remain hidden until combat reveal;
- public events contain only revealed/public information;
- every command carries an expected revision;
- unauthorized, stale and illegal commands return explicit errors without
  advancing the revision or consuming cards/actions.

## Automated evidence

```bash
gofmt -w apps internal docs/cards/phase-4b/embed.go docs/fighters/phase-4b/embed.go tests/integration/playableslice
go vet ./...
go test -race ./...
python scripts/validate_playable_slice.py
python -m unittest tests/architecture/test_playable_slice.py
npm ci
npm run check
npm run build
```

`TestDeterministicCompleteMatch` executes the same scripted match twice and
compares byte-identical final projections. The script uses the real deck and
battlefield definitions, moves multiple fighters in a boosted Maneuver, plays a
Scheme, resolves Attack/Defense and representative unique effects from both
fighters, defeats Robin Hood and declares Bigfoot the winner. Focused tests also
reject illegal movement, illegal attack, stale revision and unauthorized/private
choices without partial mutation.

## Practical two-tab checklist

1. Create Robin Hood in tab A and join Bigfoot in tab B.
2. Confirm each tab displays only its own five-card hand.
3. Execute a Maneuver, choose optional BOOST, move more than one fighter and
   finish the maneuver.
4. Play a Scheme and resolve any private/public prompt.
5. Attack from a legal melee or ranged position.
6. Confirm tab A cannot inspect or submit tab B's defense choice.
7. Resolve turns until a hero is defeated and confirm both tabs show the same
   winner.
8. Submit a stale command and confirm the UI reports the rejection and refreshes
   without a partial state change.

## Deferred limitations

- state is lost when the local server process stops;
- no production authentication, matchmaking, spectators, ranking or AI;
- no additional fighters, battlefields, team modes or Adventures;
- no final animation, branding or production network transport;
- independent edge-by-edge visual QA of the Sherwood Forest transcription and a
  final browser usability pass belong to the single fresh independent QA gate.
