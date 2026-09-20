# Towergame — a SimTower-like building sim

Multiplayer tower sim. Players build rooms, stairs, and elevators on a grid.
Sims commute, work, shop, and ride elevators. Go server owns the world.
Vite + React + three.js client renders it.

## Rules for agents

- **Never start or restart a server.** The user runs them. Ask, then test.
- Comments <= 7 words. Function names <= 4 words. User-facing strings <= 10 words.
- ASD-STE100 style: short sentences, active voice, present tense.
  One term per concept. No hype.
- `client/src/protocol.ts` mirrors `server/messages.go`. Change both or neither.
- Tests are the spec. Run them after every change (see Commands).

## Commands

| Task | How |
|---|---|
| Server tests | `cd server && go test ./...` |
| Client typecheck + build | `cd client && npm run build` |
| Run the game | user starts: server `go run .` on :7777, client `npm run dev` on :3000 |

## Layout

    client/src/
      main.tsx          entry; splits /card routes from the game
      App.tsx           canvas + HUD
      Game.tsx          scene root: skyline, scenery, tower, build surface
      CameraControls.tsx  pan and zoom, flat tower view
      Glow.tsx          hover highlight box
      socket.ts         WebSocket wrapper, silence probe, resync
      protocol.ts       wire types, mirrors server/messages.go
      store.ts          hover + selected room type (useSyncExternalStore)
      catalog.ts        buildable types from each snapshot
      seating.ts        seat facings per snapshot
      simsFeed.ts       one subscription, split per reader
      dims.ts           ALL shared geometry constants. Add nothing elsewhere.
      useBuildInput.ts  pointer -> build intents -> socket
      inspect.ts        re-asks the last inspect on snapshot
      Tower.tsx         floors, rooms, stairs, shafts
      RoomShell.tsx     outer shell + glass per room
      RoomInterior.tsx  wall art, props, seats, occupants
      Furniture.tsx     desks, chairs, tables, beds
      paint.ts          canvas texture toolkit. PPU = 22 px per world unit.
      pixelArt.ts       hand-authored sprites: palette + rows of chars
      roomArt.ts        one back-wall texture per room type
      propArt.ts        flat props that stand before the wall
      lobbyArt.ts       repeating lobby wall strip
      roomTint.ts       colours + per-category light tint
      layout.ts         x positions across a room width
      roomCard.ts       flat card render for /card pages. Stamps the same art.
      Sims.tsx          sim positions from simsFeed
      SimModel.tsx      one sim: rig, walk clip, hover
      simRig.ts         GLTF bones + hand-posed sit pose
      simOutfit.ts      clothing colours per light and hover
      ElevatorModel.tsx shaft rails
      ElevatorCars.tsx  cars: position, riders, click
      ui/               Clock, Budget, RoomPalette, RoomInspector, ErrorToast, CardLinks

    server/
      main.go           HTTP + /ws on :7777
      hub.go            hub goroutine: 30 ticks/s, broadcasts every 3
      client.go         one connection: outbound queue, ping/pong
      messages.go       wire structs. Mirrors client protocol.ts.
      dispatch.go       inbound commands -> World calls
      world.go          World: grid, money, Place/Remove, snapshot
      schedule.go       daily goals per sim
      sims.go           sim states, stress meter
      walk.go           walking, stair, browse, wander speeds
      riders.go         elevator queues, boarding, riding
      visitors.go       shop and diner customers
      lights.go         per-room light hours
      clock.go          8 game-min per real-sec; 3-day week, days 1-2 work
      inspect.go        inspect answers
      factions_test.go  professions and alignment tests
      model/            pure state, no I/O
        grid.go         cells, rooms, stairs, elevators. Grid: 100 x 40.
        room.go         RoomType catalog, Category, Habit, Alignment
        transport.go    stair and elevator geometry
        pathfind.go     nav graph paths: walk, stair, elevator waypoints
      transport/        elevator physics
        bank.go         cars, shafts. Capacity 4, 8 cars per shaft.
        dispatch.go     FS car scoring, Barney's Elevator Traffic Handbook

## Architecture

Server: one game-loop goroutine owns all state. Connection goroutines send
commands on a channel. No locks. Never touch `World` from a connection.

Client: no game state. Snapshots replace the tower wholesale; the `sims`
stream (~10/s) carries positions only. Stores use `useSyncExternalStore`.

Money: start 2,000,000. Base 500, stair 5,000, elevator 100,000.

## Protocol (raw WebSocket, JSON, every frame is `{"type": ...}`)

| Direction | type | payload | meaning |
|---|---|---|---|
| c -> s | `placeroom` | `{room, align, floor, col}` | build a room |
| c -> s | `placebase` | `{floor, col}` | build one base cell |
| c -> s | `placestair` | `{floor, col}` | build a stair |
| c -> s | `placeescalator` | `{floor, col}` | build an escalator |
| c -> s | `placeramp` | `{floor, col}` | build one run of a car ramp |
| c -> s | `placeelevator` | `{floor, col}` | build or extend a shaft |
| c -> s | `remove` | `{floor, col}` | demolish at a cell |
| c -> s | `inspect` | `{floor, col, quiet}` | what is at this cell |
| c -> s | `inspectsim` | `{id}` | about one sim |
| c -> s | `inspectcar` | `{id}` | about one elevator car |
| c -> s | `resync` | `{}` | fresh snapshot; doubles as liveness probe |
| c -> s | `message` | `{text}` | chat |
| s -> c | `snapshot` | `{tick, clock, money, stars, population, grid, types}` | full state after each change |
| s -> c | `sims` | `{sims, cars, vehicles, train}` | positions, ~10/s |
| s -> c | `report` | `{quarter, hotel, sales, leases, events, upkeep, net, money}` | books when a quarter closes |
| s -> c | `chatUpdate` | `{chat}` | full chat log |
| s -> c | `inspect` | `{title, lines}` | one inspect answer |
| s -> c | `error` | `{reason}` | why an action failed |

On connect: `chatUpdate`, then `snapshot`. Reconnect resyncs by itself.

## Domain glossary

- **Cell** — one grid square. Floor = row, col = column.
- **Floor** — one row of cells. Ground floor is the lobby level.
- **Base** — built structure cells under a room.
- **Room** — cells on one floor, with type, cost, income, capacity, occupants.
- **Category** — room kind: lobby, office, residential, hotel, food, retail.
- **Alignment** — room politics: neutral, good, triad.
- **Shaft** — vertical run of cells: stair, escalator, or elevator shaft.
- **Car** — one elevator cab. A **bank** is the set of cars.
- **Sim** — an agent with home, work, schedule, needs, path, profession
  (worker, resident, visitor, officer, triad).
- **Nav graph** — floors are segments, shafts are edges. Pathfinding runs
  here, not on the raw grid.

## Caveats

- `dims.ts` is the only home for shared geometry. Two modules that must
  agree import from it.
- Cards must match rooms. `/card` pages stamp `roomArt`, `propArt`, and
  `paint` directly, so art fixes show up in game and on cards alike.
- Room facts come from the server catalog. Keep no local copy of costs.
- Server numbers are floats (cell positions). Round only at the edges.

## Card screenshot loop (style work)

1. Ask the user to start both servers. Vite listens on `[::1]:3000`.
2. Screenshot, then read the PNG:

       chromium --headless=new --disable-gpu --no-sandbox \
         --screenshot=/tmp/card.png --window-size=1200,800 \
         --virtual-time-budget=12000 "http://[::1]:3000/card/office"

3. Read the PNG. Adjust `roomArt`, `propArt`, or `paint`. Repeat.
   `?stair` adds a stair in front of the card.
