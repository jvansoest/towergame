// Everything around the tower: ground, street, and doors.

import { CELL_W, DEPTH, FLOOR_T, WALL_T } from "./dims";
import type { GridView } from "./protocol";

// Pavement reaching past both ends of the tower.
const STREET_PAD = 12;
const STREET_COLOR = "#4b5563";

// Earth below the street. Basements come later.
const EARTH_COLOR = "#6b4f36";
const TOPSOIL_COLOR = "#4a3626";
const TOPSOIL_T = 0.5;
const EARTH_D = 60; // world units drawn below ground
const EARTH_W = 2000; // wide enough to fill a zoomed-out view

// Entrance doors in the lobby's end walls.
const DOOR_W = 1.8; // across the room depth
const DOOR_H = 2.4;
const FRAME_T = 0.12;
const FRAME_COLOR = "#3f4652";
const GLASS_COLOR = "#bfe3f5";

// Ground, street, and the doors people use.
export default function Scenery({
  towerW,
  grid,
}: {
  towerW: number;
  grid: GridView | null;
}) {
  return (
    <>
      <Ground width={towerW + STREET_PAD * 2} />
      {grid && <Doors grid={grid} />}
    </>
  );
}

// Street level, with solid earth under it.
function Ground({ width }: { width: number }) {
  return (
    <group>
      {/* People walk in and out along this. */}
      <mesh position={[0, FLOOR_T / 2 - 0.02, 0]}>
        <boxGeometry args={[width, FLOOR_T, DEPTH]} />
        <meshStandardMaterial color={STREET_COLOR} />
      </mesh>
      {/* The ground line itself. */}
      <mesh position={[0, -TOPSOIL_T / 2, 0]}>
        <boxGeometry args={[EARTH_W, TOPSOIL_T, DEPTH]} />
        <meshStandardMaterial color={TOPSOIL_COLOR} />
      </mesh>
      <mesh position={[0, -TOPSOIL_T - EARTH_D / 2, 0]}>
        <boxGeometry args={[EARTH_W, EARTH_D, DEPTH]} />
        <meshStandardMaterial color={EARTH_COLOR} />
      </mesh>
    </group>
  );
}

// A door in each end wall of the ground floor.
function Doors({ grid }: { grid: GridView }) {
  const row = grid.built?.[0] ?? [];
  let left = -1;
  let right = -1;
  for (let c = 0; c < grid.width; c++) {
    if (!row[c]) continue;
    if (left < 0) left = c;
    right = c;
  }
  if (left < 0) return null;
  // The inner face of each end wall, which is the one seen.
  const inset = WALL_T + 0.02;
  const leftX = (left - grid.width / 2) * CELL_W + inset;
  const rightX = (right + 1 - grid.width / 2) * CELL_W - inset;
  return (
    <>
      <LobbyDoor x={leftX} facing={1} />
      {right !== left && <LobbyDoor x={rightX} facing={-1} />}
    </>
  );
}

// Glass double door set in an end wall.
function LobbyDoor({ x, facing }: { x: number; facing: number }) {
  const half = DOOR_W / 2;
  const top = DOOR_H / 2;
  return (
    <group
      position={[x, FLOOR_T + DOOR_H / 2, 0]}
      rotation={[0, (facing * Math.PI) / 2, 0]}
    >
      <mesh>
        <boxGeometry args={[DOOR_W, DOOR_H, 0.06]} />
        <meshStandardMaterial
          color={GLASS_COLOR}
          transparent
          opacity={0.45}
        />
      </mesh>
      {/* posts, header, threshold, and the mullion */}
      <mesh position={[-half + FRAME_T / 2, 0, 0]}>
        <boxGeometry args={[FRAME_T, DOOR_H, 0.1]} />
        <meshStandardMaterial color={FRAME_COLOR} />
      </mesh>
      <mesh position={[half - FRAME_T / 2, 0, 0]}>
        <boxGeometry args={[FRAME_T, DOOR_H, 0.1]} />
        <meshStandardMaterial color={FRAME_COLOR} />
      </mesh>
      <mesh position={[0, top - FRAME_T / 2, 0]}>
        <boxGeometry args={[DOOR_W, FRAME_T, 0.1]} />
        <meshStandardMaterial color={FRAME_COLOR} />
      </mesh>
      <mesh position={[0, -top + FRAME_T / 2, 0]}>
        <boxGeometry args={[DOOR_W, FRAME_T, 0.1]} />
        <meshStandardMaterial color={FRAME_COLOR} />
      </mesh>
      <mesh>
        <boxGeometry args={[FRAME_T * 0.7, DOOR_H, 0.1]} />
        <meshStandardMaterial color={FRAME_COLOR} />
      </mesh>
      {[-1, 1].map((s) => (
        <mesh key={s} position={[s * 0.2, -0.15, 0.07]}>
          <boxGeometry args={[0.05, 0.5, 0.05]} />
          <meshStandardMaterial color="#cbd5e1" />
        </mesh>
      ))}
    </group>
  );
}
