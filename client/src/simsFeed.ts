// One subscription to the render-rate "sims" stream.
// Sim positions and car positions are split per reader.

import { useSyncExternalStore } from "react";
import socket from "./socket";
import type { CarView, SimView } from "./protocol";

export type SimTarget = {
  x: number;
  y: number;
  riding: boolean;
  sitting: boolean;
  dark: boolean; // the room's lights are off
};

// Latest server positions, read per frame.
export const simTargets = new Map<number, SimTarget>();

let simIds: number[] = [];
let cars: CarView[] = [];

const idListeners = new Set<() => void>();
const carListeners = new Set<() => void>();

function sameIds(a: number[], b: number[]) {
  return a.length === b.length && a.every((v, i) => v === b[i]);
}

socket.on("sims", (data) => {
  const list = (data.sims as SimView[]) ?? [];
  const seen = new Set<number>();
  for (const s of list) {
    simTargets.set(s.id, {
      x: s.x,
      y: s.y,
      riding: s.riding ?? false,
      sitting: s.sitting ?? false,
      dark: s.dark ?? false,
    });
    seen.add(s.id);
  }
  for (const id of [...simTargets.keys()]) {
    if (!seen.has(id)) simTargets.delete(id);
  }

  // Re-render only when the roster changes.
  const ids = list.map((s) => s.id);
  if (!sameIds(ids, simIds)) {
    simIds = ids;
    idListeners.forEach((l) => l());
  }

  cars = (data.cars as CarView[]) ?? [];
  carListeners.forEach((l) => l());
});

export function useSimIds(): number[] {
  return useSyncExternalStore(
    (cb) => {
      idListeners.add(cb);
      return () => idListeners.delete(cb);
    },
    () => simIds
  );
}

export function useCars(): CarView[] {
  return useSyncExternalStore(
    (cb) => {
      carListeners.add(cb);
      return () => carListeners.delete(cb);
    },
    () => cars
  );
}
