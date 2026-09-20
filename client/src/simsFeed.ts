// One subscription to the render-rate "sims" stream.
// Sim positions and car positions are split per reader.

import { useSyncExternalStore } from "react";
import socket from "./socket";
import type { CarView, SimView, TrainView, VehicleView } from "./protocol";

export type SimTarget = {
  x: number;
  y: number;
  riding: boolean;
  sitting: boolean;
  spot: number; // seat index in the room
  gliding: boolean; // riding an escalator
  cleaning: boolean; // a maid is sweeping
  waiting: boolean; // queued at a shaft
  dark: boolean; // the room's lights are off
  profession: string; // visitor, resident, worker, security, triad
};

// Latest server positions, read per frame.
export const simTargets = new Map<number, SimTarget>();

let simIds: number[] = [];
let cars: CarView[] = [];
let vehicleIds: number[] = [];
const vehicleListeners = new Set<() => void>();

// The subway train, read per frame.
export const trainState: { current: TrainView | null } = { current: null };

// Latest server view of each car, read per frame.
export const vehicleTargets = new Map<number, VehicleView>();

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
      spot: s.spot ?? 0,
      gliding: s.gliding ?? false,
      cleaning: s.cleaning ?? false,
      waiting: s.waiting ?? false,
      dark: s.dark ?? false,
      profession: s.profession ?? "",
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

  trainState.current = (data.train as TrainView | null) ?? null;

  const rides = (data.vehicles as VehicleView[]) ?? [];
  vehicleTargets.clear();
  for (const v of rides) vehicleTargets.set(v.id, v);
  const rideIds = rides.map((v) => v.id);
  if (!sameIds(rideIds, vehicleIds)) {
    vehicleIds = rideIds;
    vehicleListeners.forEach((l) => l());
  }
});

export function useSimIds(): number[] {
  return useSyncExternalStore(
    (cb) => {
      idListeners.add(cb);
      return () => idListeners.delete(cb);
    },
    () => simIds,
  );
}

export function useCars(): CarView[] {
  return useSyncExternalStore(
    (cb) => {
      carListeners.add(cb);
      return () => carListeners.delete(cb);
    },
    () => cars,
  );
}

// Ids of the cars in the garage.
export function useVehicleIds(): number[] {
  return useSyncExternalStore(
    (cb) => {
      vehicleListeners.add(cb);
      return () => vehicleListeners.delete(cb);
    },
    () => vehicleIds,
  );
}
