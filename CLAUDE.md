# Towergame — a SimTower-like building sim

Comment blocks are <= 7 words, function names <= 4 words. User-facing message strings should be <= 10 words. Use an active voice, no stage performances, and pick the most common word when choosing among alternatives.

## Agent rules

DO NOT START A SERVER I WILL START AND RESTART IT.
**Writing.** Use clear, standard technical English in the spirit of **ASD-STE100
(Simplified Technical English)**:
- Short sentences, one idea each. Active voice, present tense.
- Say what a thing *does*, not how impressive it is. No filler, no hype.
- Use one term per concept, consistently (a "room" is always a "room").

## What we want this game to be

A SimTower / Yoot Tower-style building sim. The player places **rooms** (offices,
apartments, restaurants, shops, hotel rooms, lobbies) on a grid and connects
floors with **stairs, escalators and elevators**. A living population of **sims**
walks the building — commuting to work, eating lunch, going home, riding
elevators. The player grows the tower from one star to a five-star TOWER by
managing money, traffic, and satisfaction. Multiplayer: many players share the
world in real time.

## Client ↔ server protocol (raw WebSocket, JSON)

Every message is `{"type": "...", ...}`.

| Direction | type | payload | meaning |
|---|---|---|---|
| client → server | `boxplaced` | `{x, y}` | intent: place a box/room |
| client → server | `message` | `{text}` | chat |
| server → client | `boxupdate` | `{x, y, matrix}` | grid changed |
| server → client | `chatUpdate` | `{chat: []}` | chat log changed |
| server → client | `snapshot` | `{...}` | (future) full sim state per tick |

The Go server owns state in a single **game-loop goroutine**; connection
goroutines send it commands over a channel (no locks). See `server/`.



## Domain glossary

- **Cell** — one grid square. The tower is a 2D array of cells (rows = floors,
  cols = horizontal position).
- **Floor** — one horizontal row of cells. Ground floor = lobby level.
- **Room** — a rectangular footprint of cells on a single floor. Has a type
  (office, apartment, …), cost, income, capacity, occupants.
- **Shaft** — a vertical run of cells used by a stair/escalator/elevator,
  connecting multiple floors.
- **Sim** — an agent. Has a home and/or workplace, a daily schedule, needs, and
  a current path. Rendered as a walking 3D character.
- **Nav graph** — floors are walkable segments; shafts are edges between them.
  Pathfinding runs on this graph, not the raw grid.
---

# Development roadmap (sprints)

TODO
 - stress level popup should update ever X ticks
 - the Z Z Z should have 3 different images or things or whatever it is and rotate between the 3 such that it looks like it is moving
 - the Z Z Z should also be a bit more shaded the same shade as night time
 - FIX sims sometimes seem to have their feet inside the floot or too low or something's
