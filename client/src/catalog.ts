// The server's buildable room types, fed
// once per snapshot like the rooms are.
// The one source for what can be built.

import { useSyncExternalStore } from "react";

import socket from "./socket";
import type { RoomTypeView } from "./protocol";

let types: RoomTypeView[] = [];
const listeners = new Set<() => void>();

socket.on("snapshot", (data) => {
  const g = data.grid as { types?: RoomTypeView[] } | undefined;
  if (!g?.types) return;
  types = g.types;
  listeners.forEach((l) => l());
});

// The catalog, refreshing per snapshot.
export function useTypes(): RoomTypeView[] {
  return useSyncExternalStore(
    (cb) => {
      listeners.add(cb);
      return () => listeners.delete(cb);
    },
    () => types,
  );
}

// Raw read for callers outside React.
export function getTypes(): RoomTypeView[] {
  return types;
}

// Calls back now, or when the first catalog
// lands. One-shot, page-lifetime.
export function whenTypes(cb: () => void) {
  if (types.length) cb();
  else listeners.add(cb);
}
