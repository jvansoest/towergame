// The tower itself: floors, rooms, stairs, and shafts.

import { useHover } from "./store";
import Glow from "./Glow";
import RoomInterior from "./RoomInterior";
import StairModel from "./StairModel";
import ElevatorModel from "./ElevatorModel";
import { CELL_W, CELL_H, DEPTH, FLOOR_T, FRONT_Z, SHAFT_CELLS } from "./dims";
import type { GridView } from "./protocol";
import type { ShaftDrag } from "./useBuildInput";

const BASE_COLOR = "#9ca3af";
const SLAB_COLOR = "#6b7280";

// A run of same-type cells, drawn as one piece.
type Kind = "base" | "lobby";
type Strip = { x: number; y: number; w: number; kind: Kind };

// Merges same-type cells into seamless strips.
function stripsOf(grid: GridView): Strip[] {
  // Rooms and stairs hide the base behind them.
  const covered = new Set<string>();
  const lobbyCells = new Set<string>();
  grid.rooms?.forEach((r) => {
    for (let f = r.floor; f < r.floor + r.height; f++) {
      for (let c = r.col; c < r.col + r.width; c++) {
        covered.add(`${f},${c}`);
        if (r.category === "lobby") lobbyCells.add(`${f},${c}`);
      }
    }
  });
  grid.stairs?.forEach((s) => {
    for (let f = s.floor; f < s.floor + s.height; f++) {
      for (let c = s.col; c < s.col + s.width; c++) covered.add(`${f},${c}`);
    }
  });

  const strips: Strip[] = [];
  for (let f = 0; f < grid.floors; f++) {
    let start = 0;
    let kind: Kind | null = null;
    for (let c = 0; c <= grid.width; c++) {
      const key = `${f},${c}`;
      let cell: Kind | null = null;
      if (c < grid.width) {
        if (lobbyCells.has(key)) cell = "lobby";
        else if (grid.built[f]?.[c] && !covered.has(key)) cell = "base";
      }
      if (cell !== kind) {
        if (kind !== null) {
          const len = c - start;
          strips.push({
            x: (start + len / 2 - grid.width / 2) * CELL_W,
            y: (f + 0.5) * CELL_H,
            w: len * CELL_W,
            kind,
          });
        }
        start = c;
        kind = cell;
      }
    }
  }
  return strips;
}

export default function Tower({ grid }: { grid: GridView }) {
  const hover = useHover();
  const strips = stripsOf(grid);

  return (
    <>
      {strips.map((s, i) =>
        s.kind === "lobby" ? (
          <group key={`s${i}`} position={[s.x, s.y, 0]}>
            <RoomInterior
              width={s.w}
              height={CELL_H}
              depth={DEPTH}
              category="lobby"
              type="lobby"
              seats={0}
              lit
              occupants={0}
            />
          </group>
        ) : (
          <group key={`s${i}`} position={[s.x, s.y, 0]}>
            {/* back wall */}
            <mesh position={[0, 0, -DEPTH / 2 + 0.1]}>
              <boxGeometry args={[s.w, CELL_H, 0.2]} />
              <meshStandardMaterial color={BASE_COLOR} />
            </mesh>
            {/* floor slab */}
            <mesh position={[0, -CELL_H / 2 + FLOOR_T / 2, 0]}>
              <boxGeometry args={[s.w, FLOOR_T, DEPTH]} />
              <meshStandardMaterial color={SLAB_COLOR} />
            </mesh>
          </group>
        )
      )}

      {grid.stairs?.map((s, i) => {
        const rise = s.height * CELL_H;
        const x = (s.col + s.width / 2 - grid.width / 2) * CELL_W;
        // The run climbs from one slab top to the next.
        const y = s.floor * CELL_H + rise / 2 + FLOOR_T;
        return (
          <group key={`st${i}`} position={[x, y, FRONT_Z]}>
            <StairModel width={s.width * CELL_W} height={rise} />
            {hover?.kind === "stair" && hover.key === `${s.floor},${s.col}` && (
              <Glow width={s.width * CELL_W} height={rise} depth={DEPTH / 2} />
            )}
          </group>
        );
      })}

      {grid.elevators?.map((e, i) => {
        const floors = e.top - e.bottom + 1;
        const x = (e.col + e.width / 2 - grid.width / 2) * CELL_W;
        const y = (e.bottom + floors / 2) * CELL_H;
        return (
          <group key={`el${i}`} position={[x, y, FRONT_Z]}>
            <ElevatorModel width={e.width * CELL_W} height={floors * CELL_H} />
            {hover?.kind === "shaft" && hover.id === e.id && (
              <Glow
                width={e.width * CELL_W}
                height={floors * CELL_H}
                depth={DEPTH / 2}
              />
            )}
          </group>
        );
      })}

      {grid.rooms
        ?.filter((r) => r.category !== "lobby")
        .map((r, i) => {
          const x = (r.col + r.width / 2 - grid.width / 2) * CELL_W;
          const y = (r.floor + r.height / 2) * CELL_H;
          return (
            <group key={i} position={[x, y, 0]}>
              <RoomInterior
                width={r.width * CELL_W}
                height={r.height * CELL_H}
                depth={DEPTH}
                category={r.category}
                type={r.type}
                seats={r.seats}
                lit={r.lit}
                occupants={r.occupants}
              />
              {hover?.kind === "room" && hover.key === `${r.floor},${r.col}` && (
                <Glow
                  width={r.width * CELL_W}
                  height={r.height * CELL_H}
                  depth={DEPTH}
                />
              )}
            </group>
          );
        })}
    </>
  );
}

// Shows how far a dragged shaft will reach.
export function ShaftGhost({
  grid,
  drag,
}: {
  grid: GridView;
  drag: ShaftDrag;
}) {
  const shaft = grid.elevators?.find(
    (e) => drag.col >= e.col && drag.col < e.col + e.width
  );
  const bottom = shaft ? shaft.bottom : 0;
  const col = shaft ? shaft.col : drag.col;
  const wide = (shaft?.width ?? SHAFT_CELLS) * CELL_W;
  const floors = drag.top - bottom + 1;
  const x = (col + (shaft?.width ?? SHAFT_CELLS) / 2 - grid.width / 2) * CELL_W;
  const y = (bottom + floors / 2) * CELL_H;
  return (
    <group position={[x, y, FRONT_Z]}>
      <Glow width={wide} height={floors * CELL_H} depth={DEPTH / 2} />
    </group>
  );
}
