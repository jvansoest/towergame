// One run of a car ramp: a slanted road across a floor's
// row, with kerbs and a dashed centre line. It starts at the
// top of the row and lands on the floor, turning the cars
// back the other way on the next run.

import { CELL_H, CELL_W, FLOOR_T } from "./dims";

const ROAD_T = 0.16;
const ROAD_D = 2.2; // front to back
const KERB_H = 0.22;

export default function RampModel({
  width,
  downRight,
}: {
  width: number;
  downRight: boolean;
}) {
  // The road joins the end cell centres; the ends overhang by half a cell.
  const run = width - CELL_W;
  const rise = CELL_H;
  const len = Math.hypot(run, rise) + CELL_W;
  const angle = Math.atan2(rise, run) * (downRight ? -1 : 1);
  const dashes = Math.floor(len / 0.9);
  // Centre line of the row's slab, so the road runs floor to floor.
  const y = FLOOR_T / 2 - 0.02;
  return (
    <group position={[0, y, 0]} rotation={[0, 0, angle]}>
      <mesh>
        <boxGeometry args={[len, ROAD_T, ROAD_D]} />
        <meshStandardMaterial color="#3a3f47" />
      </mesh>
      {/* kerbs at the front and back edges */}
      {[-1, 1].map((s) => (
        <mesh key={s} position={[0, KERB_H / 2, (s * (ROAD_D - 0.2)) / 2]}>
          <boxGeometry args={[len, KERB_H, 0.2]} />
          <meshStandardMaterial color="#8b929c" />
        </mesh>
      ))}
      {/* dashes along the middle */}
      {Array.from({ length: dashes }).map((_, i) => (
        <mesh
          key={i}
          position={[-len / 2 + (i + 0.5) * (len / dashes), 0.09, 0]}
        >
          <boxGeometry args={[0.45, 0.02, 0.12]} />
          <meshStandardMaterial color="#f2d675" />
        </mesh>
      ))}
    </group>
  );
}
