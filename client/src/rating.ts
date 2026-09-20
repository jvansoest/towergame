// The tower's rating and population, fed by
// each snapshot.

import { useSyncExternalStore } from "react";

import socket from "./socket";

export type Rating = { stars: number; population: number };

let rating: Rating = { stars: 1, population: 0 };
const listeners = new Set<() => void>();

socket.on("snapshot", (data) => {
  const stars = Number(data.stars ?? 1);
  const population = Number(data.population ?? 0);
  if (stars === rating.stars && population === rating.population) return;
  rating = { stars, population };
  listeners.forEach((l) => l());
});

// The rating, refreshing per snapshot.
export function useRating(): Rating {
  return useSyncExternalStore(
    (cb) => {
      listeners.add(cb);
      return () => listeners.delete(cb);
    },
    () => rating,
  );
}
