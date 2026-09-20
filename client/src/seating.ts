// Which way a seated sim looks: toward the
// table belonging to their chair. One facing
// per seat cell, rebuilt on each snapshot.
// Also feeds the room list to slow readers.

import { useSyncExternalStore } from "react";

import socket from "./socket";
import { getTypes } from "./catalog";
import { CELL_W, SIM_Z, SIT_FACING, TABLE_Z } from "./dims";
import { barCol, barZ, seatCols, tableCols } from "./layout";
import type { RoomView } from "./protocol";

// Seat cell key -> which way to look.
const facings = new Map<string, number>();

// Bar seat cell key -> seat count.
const lines = new Map<string, number>();

let rooms: RoomView[] = [];
const roomListeners = new Set<() => void>();

socket.on("snapshot", (data) => {
  const g = data.grid as { rooms?: RoomView[] } | undefined;
  if (!g?.rooms) return;
  rooms = g.rooms;
  roomListeners.forEach((l) => l());
  facings.clear();
  lines.clear();
  for (const r of g.rooms) indexRoom(r);
});

// The tower's rooms, refreshed per snapshot.
export function useRooms(): RoomView[] {
  return useSyncExternalStore(
    (cb) => {
      roomListeners.add(cb);
      return () => roomListeners.delete(cb);
    },
    () => rooms,
  );
}

function indexRoom(r: RoomView) {
  // Only painted rooms draw chairs at all.
  if (r.category === "lobby" || r.seats < 1) return;
  if (getTypes().find((t) => t.id === r.type)?.line) {
    lines.set(`${r.floor}:${r.col + barCol(r.width)}`, r.seats);
    return;
  }
  const desks = r.category === "office";
  const cols = seatCols(r.width, r.seats);
  for (let i = 0; i < cols.length; i++) {
    const tab = nearest(tableCols(r.width, r.seats, desks), cols[i]);
    if (tab === undefined) continue; // lone chair keeps walking-in look
    const dx = (tab - cols[i]) * CELL_W;
    facings.set(
      `${r.floor}:${r.col + cols[i]}`,
      Math.atan2(dx, TABLE_Z - SIM_Z),
    );
  }
}

// The closest table column, or undefined.
function nearest(tabs: number[], col: number): number | undefined {
  let best: number | undefined;
  for (const t of tabs) {
    if (best === undefined || Math.abs(t - col) < Math.abs(best - col)) {
      best = t;
    }
  }
  return best;
}

// Facing for a sim seated at this grid spot.
export function sitFacing(x: number, y: number): number {
  return facings.get(`${Math.round(y)}:${Math.round(x)}`) ?? SIT_FACING;
}

// Depth of a bar seat, if any.
export function sitLane(x: number, y: number, spot: number) {
  const n = lines.get(`${Math.round(y)}:${Math.round(x)}`);
  return n === undefined ? undefined : barZ(spot, n);
}
