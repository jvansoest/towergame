// The tower itself: floors, rooms, stairs, and shafts.

import { useMemo } from "react";

import { useHover, useSelectedRoom } from "./store";
import Glow from "./Glow";
import { variantOf } from "./layout";
import RoomInterior from "./RoomInterior";
import StairModel from "./StairModel";
import EscalatorModel from "./EscalatorModel";
import ElevatorModel from "./ElevatorModel";
import {
  CELL_W,
  CELL_H,
  colToX,
  DEPTH,
  DOOMED,
  FLOOR_T,
  FRONT_Z,
  HIGHLIGHT,
  SHAFT_CELLS,
} from "./dims";
import { builtAt } from "./gridCells";
import RampModel from "./RampModel";
import type { GridView } from "./protocol";
import type { ShaftDrag } from "./useBuildInput";

const BASE_COLOR = "#9ca3af";
const SLAB_COLOR = "#6b7280";
// Concrete cut into the earth is darker.
const UNDER_COLOR = "#414852";
const UNDER_SLAB = "#2b2f38";

// A run of same-type cells, drawn as one piece.
type Kind = "base" | "lobby";
type Strip = {
  x: number;
  y: number;
  w: number;
  c0: number;
  f: number;
  kind: Kind;
};

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
  for (let f = -(grid.basement ?? 0); f < grid.floors; f++) {
    let start = 0;
    let kind: Kind | null = null;
    for (let c = 0; c <= grid.width; c++) {
      const key = `${f},${c}`;
      let cell: Kind | null = null;
      if (c < grid.width) {
        if (lobbyCells.has(key)) cell = "lobby";
        else if (builtAt(grid, f, c) && !covered.has(key)) cell = "base";
      }
      if (cell !== kind) {
        if (kind !== null) {
          const len = c - start;
          strips.push({
            x: (start + len / 2 - grid.width / 2) * CELL_W,
            y: (f + 0.5) * CELL_H,
            w: len * CELL_W,
            c0: start,
            f,
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
  // Red while removing, yellow while inspecting.
  const glow = useSelectedRoom() === "bulldoze" ? DOOMED : HIGHLIGHT;
  // Strips rebuild only when the tower does,
  // not on every hover-driven re-render.
  const strips = useMemo(() => stripsOf(grid), [grid]);

  // Columns each floor's rooms cover, so walls
  // between touching rooms can be told apart
  // from the tower's outer edges.
  const roomCols = useMemo(() => {
    const m = new Map<number, Set<number>>();
    grid.rooms?.forEach((r) => {
      for (let f = r.floor; f < r.floor + r.height; f++) {
        let cols = m.get(f);
        if (!cols) {
          cols = new Set<number>();
          m.set(f, cols);
        }
        for (let c = r.col; c < r.col + r.width; c++) cols.add(c);
      }
    });
    return m;
  }, [grid.rooms]);
  const besideAt = (f: number, c: number) => !!roomCols.get(f)?.has(c);

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
              roomLeft={besideAt(s.f, s.c0 - 1)}
              roomRight={besideAt(s.f, s.c0 + Math.round(s.w / CELL_W))}
            />
          </group>
        ) : (
          <group key={`s${i}`} position={[s.x, s.y, 0]}>
            {/* back wall */}
            <mesh position={[0, 0, -DEPTH / 2 + 0.1]}>
              <boxGeometry args={[s.w, CELL_H, 0.2]} />
              <meshStandardMaterial
                color={s.f < 0 ? UNDER_COLOR : BASE_COLOR}
              />
            </mesh>
            {/* floor slab */}
            <mesh position={[0, -CELL_H / 2 + FLOOR_T / 2, 0]}>
              <boxGeometry args={[s.w, FLOOR_T, DEPTH]} />
              <meshStandardMaterial color={s.f < 0 ? UNDER_SLAB : SLAB_COLOR} />
            </mesh>
          </group>
        ),
      )}

      {hover?.kind === "base" &&
        (() => {
          const [f, c] = hover.key.split(",").map(Number);
          return (
            <Glow
              width={CELL_W}
              height={CELL_H}
              depth={DEPTH}
              position={[colToX(c, grid.width), (f + 0.5) * CELL_H, 0]}
              color={glow}
            />
          );
        })()}

      {grid.ramps?.map((r, i) => {
        const x = (r.col + r.width / 2 - grid.width / 2) * CELL_W;
        const y = (r.floor + 0.5) * CELL_H;
        return (
          <group key={`rp${i}`} position={[x, y, 0]}>
            <RampModel width={r.width * CELL_W} downRight={r.downRight} />
            {hover?.kind === "ramp" && hover.key === `${r.floor},${r.col}` && (
              <Glow
                width={r.width * CELL_W}
                height={CELL_H}
                depth={DEPTH}
                color={glow}
              />
            )}
          </group>
        );
      })}

      {grid.stairs?.map((s, i) => {
        const rise = s.height * CELL_H;
        const x = (s.col + s.width / 2 - grid.width / 2) * CELL_W;
        // The run climbs from one slab top to the next.
        const y = s.floor * CELL_H + rise / 2 + FLOOR_T;
        return (
          <group key={`st${i}`} position={[x, y, FRONT_Z]}>
            {s.escalator ? (
              <EscalatorModel width={s.width * CELL_W} height={rise} />
            ) : (
              <StairModel width={s.width * CELL_W} height={rise} />
            )}
            {hover?.kind === "stair" && hover.key === `${s.floor},${s.col}` && (
              <Glow
                width={s.width * CELL_W}
                height={rise}
                depth={DEPTH / 2}
                color={glow}
              />
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
                color={glow}
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
          let sharedL = false;
          let sharedR = false;
          for (let f = r.floor; f < r.floor + r.height; f++) {
            sharedL ||= besideAt(f, r.col - 1);
            sharedR ||= besideAt(f, r.col + r.width);
          }
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
                roomLeft={sharedL}
                roomRight={sharedR}
                variant={variantOf(r.floor, r.col)}
                dirty={r.dirty}
              />
              {hover?.kind === "room" &&
                hover.key === `${r.floor},${r.col}` && (
                  <Glow
                    width={r.width * CELL_W}
                    height={r.height * CELL_H}
                    depth={DEPTH}
                    color={glow}
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
    (e) => drag.col >= e.col && drag.col < e.col + e.width,
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
