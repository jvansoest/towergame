// Pointer input for the build tools. Turns clicks and
// drags on the build surface into intents for the server.

import { useEffect, useRef, useState } from "react";
import type * as THREE from "three";
import type { ThreeEvent } from "@react-three/fiber";

import socket from "./socket";
import { askInspect } from "./inspect";
import { useSelectedRoom, setHover, type Hover } from "./store";
import { builtAt } from "./gridCells";
import { CELL_W, CELL_H } from "./dims";
import type { GridView } from "./protocol";

// Drag-paints these one-wide units.
const PAINTABLE = new Set(["base", "lobby"]);

// A shaft being dragged out, before it is sent.
export type ShaftDrag = { col: number; from: number; top: number };

// What sits at a cell, for the inspector's highlight.
function thingAt(
  grid: GridView | null,
  floor: number,
  col: number,
  bases = false,
): Hover {
  if (!grid) return null;
  const room = grid.rooms?.find(
    (r) =>
      floor >= r.floor &&
      floor < r.floor + r.height &&
      col >= r.col &&
      col < r.col + r.width,
  );
  if (room) return { kind: "room", key: `${room.floor},${room.col}` };

  const shaft = grid.elevators?.find(
    (e) =>
      floor >= e.bottom &&
      floor <= e.top &&
      col >= e.col &&
      col < e.col + e.width,
  );
  if (shaft) return { kind: "shaft", id: shaft.id };

  const ramp = grid.ramps?.find(
    (r) => floor === r.floor && col >= r.col && col < r.col + r.width,
  );
  if (ramp) return { kind: "ramp", key: `${ramp.floor},${ramp.col}` };

  const stair = grid.stairs?.find(
    (s) =>
      floor >= s.floor &&
      floor < s.floor + s.height + 1 &&
      col >= s.col &&
      col < s.col + s.width,
  );
  if (stair) return { kind: "stair", key: `${stair.floor},${stair.col}` };
  if (bases && builtAt(grid, floor, col)) {
    return { kind: "base", key: `${floor},${col}` };
  }
  return null;
}

export default function useBuildInput(grid: GridView | null) {
  const selectedRoom = useSelectedRoom();
  const inspecting = selectedRoom === "inspect";
  // Both tools mark the thing under the pointer.
  const marking = inspecting || selectedRoom === "bulldoze";

  const painting = useRef(false);
  const painted = useRef<Set<string>>(new Set());
  const lastCell = useRef<{ col: number; floor: number } | null>(null);
  const dragging = useRef<ShaftDrag | null>(null);
  const [preview, setPreview] = useState<ShaftDrag | null>(null);

  useEffect(() => {
    // Releasing sends the shaft the drag ended on.
    const stop = () => {
      painting.current = false;
      const d = dragging.current;
      if (!d) return;
      dragging.current = null;
      setPreview(null);
      socket.emit("placeelevator", { floor: d.top, col: d.col });
    };
    window.addEventListener("pointerup", stop);
    return () => window.removeEventListener("pointerup", stop);
  }, []);

  // Drop the highlight when the tool changes.
  useEffect(() => {
    if (!marking) setHover(null);
  }, [marking]);

  const cellFromPoint = (p: THREE.Vector3) => {
    if (!grid) return null;
    return {
      col: Math.floor(p.x / CELL_W + grid.width / 2),
      floor: Math.floor(p.y / CELL_H),
    };
  };

  const placeCell = (col: number, floor: number) => {
    const key = `${floor},${col}`;
    if (painted.current.has(key)) return;
    painted.current.add(key);
    if (selectedRoom === "base") socket.emit("placebase", { floor, col });
    else socket.emit("placeroom", { room: selectedRoom, floor, col });
  };

  // Fill the run so fast drags leave no gaps.
  const paintTo = (col: number, floor: number) => {
    const last = lastCell.current;
    if (last && last.floor === floor) {
      const step = col > last.col ? 1 : -1;
      for (let c = last.col; c !== col; c += step) placeCell(c, floor);
    }
    placeCell(col, floor);
    lastCell.current = { col, floor };
  };

  const onPointerDown = (e: ThreeEvent<PointerEvent>) => {
    if (e.button !== 0) return; // left button only
    const cell = cellFromPoint(e.point);
    if (!cell) return;
    if (PAINTABLE.has(selectedRoom)) {
      painting.current = true;
      painted.current = new Set();
      lastCell.current = null;
      paintTo(cell.col, cell.floor);
    } else if (selectedRoom === "ramp") {
      socket.emit("placeramp", { floor: cell.floor, col: cell.col });
    } else if (selectedRoom === "escalator") {
      socket.emit("placeescalator", { floor: cell.floor, col: cell.col });
    } else if (selectedRoom === "stairs") {
      socket.emit("placestair", { floor: cell.floor, col: cell.col });
    } else if (selectedRoom === "elevator") {
      // Hold and drag up to set how high it goes.
      const d = { col: cell.col, from: cell.floor, top: cell.floor };
      dragging.current = d;
      setPreview(d);
    } else if (selectedRoom === "bulldoze") {
      socket.emit("remove", { floor: cell.floor, col: cell.col });
    } else if (selectedRoom === "inspect") {
      askInspect("inspect", { floor: cell.floor, col: cell.col });
    } else {
      socket.emit("placeroom", {
        room: selectedRoom,
        floor: cell.floor,
        col: cell.col,
      });
    }
  };

  const onPointerMove = (e: ThreeEvent<PointerEvent>) => {
    const cell = cellFromPoint(e.point);
    if (marking) {
      const bases = selectedRoom === "bulldoze";
      setHover(cell ? thingAt(grid, cell.floor, cell.col, bases) : null);
      return;
    }
    if (dragging.current && cell) {
      const d = dragging.current;
      const top = Math.max(cell.floor, d.from);
      if (top !== d.top) {
        dragging.current = { ...d, top };
        setPreview(dragging.current);
      }
      return;
    }
    if (!painting.current) return;
    if (cell) paintTo(cell.col, cell.floor);
  };

  return { onPointerDown, onPointerMove, preview };
}
