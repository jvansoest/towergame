// The subway: a tunnel along the lowest floor, and the train
// that runs through it. Sims step off onto the platform in front.

import { useRef } from "react";
import { useFrame } from "@react-three/fiber";
import type * as THREE from "three";

import { trainState } from "./simsFeed";
import { CELL_H, CELL_W, FLOOR_T, colToX } from "./dims";
import type { GridView } from "./protocol";

const WAGONS = 4;
const WAGON_L = 4; // world units
const GAP = 0.25;
const TRAIN_Z = -0.45; // behind the platform, in front of the wall
const SPEEDY = 0.18; // smoothing

// The tunnel wall and rails, along the station's floor.
export function Tunnel({ grid }: { grid: GridView }) {
  const st = grid.rooms?.find((r) => r.type === "subway");
  if (!st) return null;
  const w = grid.width * CELL_W;
  const y = st.floor * CELL_H;
  const dashes = Math.floor(w / 1.2);
  return (
    <group position={[0, y, 0]}>
      {/* the dark tunnel wall, with brick courses */}
      <mesh position={[0, CELL_H / 2, -1.5]}>
        <planeGeometry args={[w, CELL_H]} />
        <meshBasicMaterial color="#2a2d34" toneMapped={false} />
      </mesh>
      {[0.25, 0.5, 0.75].map((t) => (
        <mesh key={t} position={[0, CELL_H * t, -1.49]}>
          <planeGeometry args={[w, 0.05]} />
          <meshBasicMaterial color="#3b3f48" toneMapped={false} />
        </mesh>
      ))}
      {/* rails and sleepers */}
      {[-0.25, 0.25].map((z) => (
        <mesh key={z} position={[0, FLOOR_T + 0.05, TRAIN_Z + z]}>
          <boxGeometry args={[w, 0.06, 0.06]} />
          <meshStandardMaterial color="#9aa1ac" />
        </mesh>
      ))}
      {Array.from({ length: dashes }).map((_, i) => (
        <mesh
          key={i}
          position={[
            -w / 2 + (i + 0.5) * (w / dashes),
            FLOOR_T + 0.02,
            TRAIN_Z,
          ]}
        >
          <boxGeometry args={[0.2, 0.04, 0.9]} />
          <meshStandardMaterial color="#4a3a2a" />
        </mesh>
      ))}
    </group>
  );
}

// The train itself, or nothing when it is away.
export default function Train({ grid }: { grid: GridView }) {
  const group = useRef<THREE.Group>(null);
  const cur = useRef<number | null>(null);
  const doors = useRef<THREE.Group[]>([]);

  useFrame(() => {
    const g = group.current;
    if (!g) return;
    const t = trainState.current;
    g.visible = t !== null;
    if (!t) {
      cur.current = null;
      return;
    }
    if (cur.current === null) cur.current = t.x;
    cur.current += (t.x - cur.current) * SPEEDY;
    g.position.set(
      colToX(cur.current, grid.width),
      t.floor * CELL_H + FLOOR_T + 0.35,
      TRAIN_Z,
    );
    for (const d of doors.current) if (d) d.visible = !t.stopped;
  });

  const total = WAGONS * WAGON_L + (WAGONS - 1) * GAP;
  return (
    <group ref={group} visible={false}>
      {Array.from({ length: WAGONS }).map((_, i) => {
        const x = -total / 2 + WAGON_L / 2 + i * (WAGON_L + GAP);
        return (
          <group key={i} position={[x, 0, 0]}>
            <mesh position={[0, 0.85, 0]}>
              <boxGeometry args={[WAGON_L, 1.7, 1.5]} />
              <meshStandardMaterial color="#c9ced6" />
            </mesh>
            <mesh position={[0, 0.5, 0.76]}>
              <boxGeometry args={[WAGON_L, 0.28, 0.02]} />
              <meshStandardMaterial color="#2f6fa8" />
            </mesh>
            {[-1.35, -0.45, 0.45, 1.35].map((wx) => (
              <mesh key={wx} position={[wx, 1.1, 0.76]}>
                <boxGeometry args={[0.6, 0.55, 0.02]} />
                <meshStandardMaterial color="#1f2a38" />
              </mesh>
            ))}
            {/* closed doors show as a darker panel; open ones hide it */}
            <group
              ref={(el) => {
                if (el) doors.current[i] = el;
              }}
            >
              <mesh position={[0, 0.8, 0.77]}>
                <boxGeometry args={[0.9, 1.3, 0.02]} />
                <meshStandardMaterial color="#8f97a3" />
              </mesh>
            </group>
          </group>
        );
      })}
    </group>
  );
}
