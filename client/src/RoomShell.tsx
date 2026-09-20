// The box a room sits in: back wall and end walls.

import { DIVIDER_CUT_Z, WALL_T } from "./dims";
import { CONCRETE } from "./roomTint";
import { tile } from "./layout";

// Back wall with real window openings and glass.
export function BackWall({
  width,
  height,
  depth,
  tint,
  windows,
}: {
  width: number;
  height: number;
  depth: number;
  tint: string;
  windows: boolean;
}) {
  const z = -depth / 2 + 0.15;

  if (!windows) {
    return (
      <mesh position={[0, 0, z]}>
        <boxGeometry args={[width, height, 0.3]} />
        <meshStandardMaterial color={tint} />
      </mesh>
    );
  }

  const winH = height * 0.42;
  const cy = height * 0.12; // window band center
  const top = cy + winH / 2;
  const bot = cy - winH / 2;
  const half = 0.7; // half window width
  const xs = tile(width, 4);

  // Solid piers filling the band between/around windows.
  const piers: { x: number; w: number }[] = [];
  let cursor = -width / 2;
  for (const wx of xs) {
    const left = wx - half;
    if (left > cursor) piers.push({ x: (cursor + left) / 2, w: left - cursor });
    cursor = wx + half;
  }
  if (width / 2 > cursor) {
    piers.push({ x: (cursor + width / 2) / 2, w: width / 2 - cursor });
  }

  return (
    <group>
      <mesh position={[0, (-height / 2 + bot) / 2, z]}>
        <boxGeometry args={[width, bot + height / 2, 0.3]} />
        <meshStandardMaterial color={tint} />
      </mesh>
      <mesh position={[0, (top + height / 2) / 2, z]}>
        <boxGeometry args={[width, height / 2 - top, 0.3]} />
        <meshStandardMaterial color={tint} />
      </mesh>
      {piers.map((p, i) => (
        <mesh key={`p${i}`} position={[p.x, cy, z]}>
          <boxGeometry args={[p.w, winH, 0.3]} />
          <meshStandardMaterial color={tint} />
        </mesh>
      ))}
      {xs.map((wx, i) => (
        <mesh key={`g${i}`} position={[wx, cy, z]}>
          <boxGeometry args={[2 * half, winH, 0.05]} />
          <meshStandardMaterial color="#bfe3f5" transparent opacity={0.25} />
        </mesh>
      ))}
    </group>
  );
}

// One concrete slab with a tinted inner face.
function WallBox({
  w,
  h,
  d,
  x,
  y,
  z,
  innerFace,
  tint,
}: {
  w: number;
  h: number;
  d: number;
  x: number;
  y: number;
  z: number;
  innerFace: number;
  tint: string;
}) {
  return (
    <mesh position={[x, y, z]}>
      <boxGeometry args={[w, h, d]} />
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <meshStandardMaterial
          key={i}
          attach={`material-${i}`}
          color={i === innerFace ? tint : CONCRETE}
        />
      ))}
    </mesh>
  );
}

// A side wall. Party walls to neighbouring rooms
// stop short of the front third, leaving the
// walkway open, so sims never cross concrete.
export function SideWall({
  x,
  innerFace,
  h,
  depth,
  tint,
  shared = false,
}: {
  x: number;
  innerFace: number;
  h: number;
  depth: number;
  tint: string;
  shared?: boolean;
}) {
  if (!shared) {
    return (
      <WallBox
        w={WALL_T}
        h={h}
        d={depth}
        x={x}
        y={0}
        z={0}
        innerFace={innerFace}
        tint={tint}
      />
    );
  }

  // Full height, but only up to the cut.
  const d = depth / 2 + DIVIDER_CUT_Z;
  return (
    <WallBox
      w={WALL_T}
      h={h}
      d={d}
      x={x}
      y={0}
      z={-depth / 2 + d / 2}
      innerFace={innerFace}
      tint={tint}
    />
  );
}
