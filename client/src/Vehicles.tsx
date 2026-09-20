// Cars in the garage: they follow the server, tilt on the
// ramp, and stand still when parked.

import { useMemo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import type * as THREE from "three";

import { vehicleTargets, useVehicleIds } from "./simsFeed";
import { CELL_H, CELL_W, FLOOR_T, colToX } from "./dims";

const COLOURS = [
  "#c0392b",
  "#2f6f9f",
  "#e0b23a",
  "#3f8f5f",
  "#8a8f99",
  "#e8e2d0",
];
const CAR_Z = 0.5; // depth, in front of the wall
const WHEEL_R = 0.16;

export default function Vehicles({ gridWidth }: { gridWidth: number }) {
  const ids = useVehicleIds();
  return (
    <>
      {ids.map((id) => (
        <Vehicle key={id} id={id} gridWidth={gridWidth} />
      ))}
    </>
  );
}

function Vehicle({ id, gridWidth }: { id: number; gridWidth: number }) {
  const group = useRef<THREE.Group>(null);
  const cur = useRef<{ x: number; y: number } | null>(null);
  const colour = COLOURS[id % COLOURS.length];
  const wheels = useMemo(() => [-0.38, 0.38], []);

  useFrame((_, dt) => {
    const t = vehicleTargets.get(id);
    const g = group.current;
    if (!t || !g) return;
    if (!cur.current) cur.current = { x: t.x, y: t.y };
    const c = cur.current;
    const k = 1 - Math.pow(0.0001, Math.min(dt, 0.1));
    c.x += (t.x - c.x) * k;
    c.y += (t.y - c.y) * k;
    g.position.set(
      colToX(c.x, gridWidth),
      c.y * CELL_H + FLOOR_T + WHEEL_R,
      CAR_Z,
    );
    // Face the way it drives; lean into the slope.
    g.scale.x = t.dir < 0 ? -1 : 1;
    g.rotation.z = t.parked
      ? 0
      : Math.atan((t.slope * CELL_H) / CELL_W) * (t.dir < 0 ? -1 : 1) * 1;
  });

  return (
    <group ref={group}>
      <mesh position={[0, 0.2, 0]}>
        <boxGeometry args={[1.25, 0.3, 0.6]} />
        <meshStandardMaterial color={colour} />
      </mesh>
      <mesh position={[-0.1, 0.45, 0]}>
        <boxGeometry args={[0.62, 0.24, 0.52]} />
        <meshStandardMaterial color={colour} />
      </mesh>
      <mesh position={[-0.1, 0.46, 0.27]}>
        <boxGeometry args={[0.5, 0.16, 0.02]} />
        <meshStandardMaterial color="#9fd0e6" />
      </mesh>
      {wheels.map((x) => (
        <mesh key={x} position={[x, 0, 0.3]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[WHEEL_R, WHEEL_R, 0.12, 10]} />
          <meshStandardMaterial color="#1f2229" />
        </mesh>
      ))}
    </group>
  );
}
