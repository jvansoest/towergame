import { useSyncExternalStore } from "react";

// Currently selected room type to place.
let selectedRoom = "base";
const listeners = new Set<() => void>();

export function setSelectedRoom(id: string) {
  selectedRoom = id;
  listeners.forEach((l) => l());
}

export function useSelectedRoom() {
  return useSyncExternalStore(
    (cb) => {
      listeners.add(cb);
      return () => listeners.delete(cb);
    },
    () => selectedRoom
  );
}

// What the pointer is over, while inspecting.
export type Hover =
  | { kind: "room"; key: string }
  | { kind: "stair"; key: string }
  | { kind: "shaft"; id: number }
  | { kind: "car"; id: number }
  | { kind: "sim"; id: number }
  | null;

let hover: Hover = null;
const hoverListeners = new Set<() => void>();

export function setHover(next: Hover) {
  if (sameHover(hover, next)) return;
  hover = next;
  hoverListeners.forEach((l) => l());
}

// Clears only if that thing is still hovered.
export function clearHover(kind: string, id: number | string) {
  if (!hover) return;
  const key = "id" in hover ? hover.id : hover.key;
  if (hover.kind === kind && key === id) setHover(null);
}

function sameHover(a: Hover, b: Hover) {
  if (a === b) return true;
  if (!a || !b || a.kind !== b.kind) return false;
  return ("id" in a ? a.id : a.key) === ("id" in b ? b.id : b.key);
}

export function useHover(): Hover {
  return useSyncExternalStore(
    (cb) => {
      hoverListeners.add(cb);
      return () => hoverListeners.delete(cb);
    },
    () => hover
  );
}
