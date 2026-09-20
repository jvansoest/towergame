// A moving stair in the SimTower style: a sloped belt
// with scrolling treads, glass sides, and a handrail.
// Sims glide on it, behind the glass.

import { useRef } from "react";
import { useFrame } from "@react-three/fiber";
import type * as THREE from "three";

import { CELL_W } from "./dims";

const TREAD_GAP = 0.36; // spacing of the ridges
const BELT_SPEED = 0.7; // world units per second
const DEPTH = 2; // front to back
const GLASS_Z = 0.95; // the balustrade, in front of the sims
const RAIL_H = 1;

export default function EscalatorModel({
  width,
  height,
}: {
  width: number;
  height: number;
}) {
  // Sims ride between the end cell centres.
  const run = width - CELL_W;
  const len = Math.hypot(run, height);
  const angle = Math.atan2(height, run);
  const count = Math.max(1, Math.floor(len / TREAD_GAP));
  const gap = len / count;
  const treads = useRef<THREE.Group>(null);

  useFrame(({ clock }) => {
    const g = treads.current;
    if (!g) return;
    const off = (clock.elapsedTime * BELT_SPEED) % gap;
    g.children.forEach((m, i) => {
      m.position.x = -len / 2 + i * gap + off;
    });
  });

  const landW = CELL_W / 2 + 0.1;
  return (
    <group>
      {/* Landings at each floor. */}
      <mesh position={[-run / 2 - landW / 2 + 0.05, -height / 2 - 0.05, 0]}>
        <boxGeometry args={[landW, 0.1, DEPTH]} />
        <meshStandardMaterial color="#6b7280" />
      </mesh>
      <mesh position={[run / 2 + landW / 2 - 0.05, height / 2 - 0.05, 0]}>
        <boxGeometry args={[landW, 0.1, DEPTH]} />
        <meshStandardMaterial color="#6b7280" />
      </mesh>

      <group rotation={[0, 0, angle]}>
        {/* The truss under the belt. */}
        <mesh position={[0, -0.34, 0]}>
          <boxGeometry args={[len + 0.3, 0.5, DEPTH * 0.9]} />
          <meshStandardMaterial color="#374151" />
        </mesh>
        <mesh position={[0, -0.05, 0]}>
          <boxGeometry args={[len, 0.1, DEPTH]} />
          <meshStandardMaterial color="#4b5563" />
        </mesh>
        {/* Ridges that crawl along the belt. */}
        <group ref={treads}>
          {Array.from({ length: count }).map((_, i) => (
            <mesh key={i} position={[-len / 2 + i * gap, 0.02, 0]}>
              <boxGeometry args={[gap * 0.35, 0.05, DEPTH * 0.92]} />
              <meshStandardMaterial color="#9ca3af" />
            </mesh>
          ))}
        </group>
        {/* Glass side and handrail, in front of the riders. */}
        <mesh position={[0, RAIL_H / 2, GLASS_Z]}>
          <planeGeometry args={[len, RAIL_H]} />
          <meshStandardMaterial
            color="#bfe3f5"
            transparent
            opacity={0.22}
            depthWrite={false}
          />
        </mesh>
        <mesh position={[0, RAIL_H, GLASS_Z]}>
          <boxGeometry args={[len + 0.3, 0.09, 0.09]} />
          <meshStandardMaterial color="#1f2937" />
        </mesh>
      </group>
    </group>
  );
}
