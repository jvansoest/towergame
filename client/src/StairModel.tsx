// Steps per floor. A step is one riser high.
const STEPS = 16;

const STEP_D = 2; // depth of the run
const RAIL_T = 0.09;
const TREAD_T = 0.09; // a tread, not a solid block
const STRINGER_T = 0.12;
const RAIL_Y = 0.9; // handrail above the treads

// A diagonal run of steps with a handrail.
export default function StairModel({
  width,
  height,
}: {
  width: number;
  height: number;
}) {
  const stepW = width / STEPS;
  const stepH = height / STEPS;
  const rail = Math.hypot(width, height);
  return (
    <group>
      {Array.from({ length: STEPS }).map((_, i) => (
        // Treads float on the stringer, with a gap below.
        <mesh
          key={i}
          position={[
            -width / 2 + stepW * (i + 0.5),
            -height / 2 + stepH * (i + 1) - TREAD_T / 2,
            0,
          ]}
        >
          <boxGeometry args={[stepW, TREAD_T, STEP_D]} />
          <meshStandardMaterial color="#94a3b8" />
        </mesh>
      ))}
      {/* The stringer the treads sit on. */}
      <mesh
        position={[0, -TREAD_T, 0]}
        rotation={[0, 0, Math.atan2(height, width)]}
      >
        <boxGeometry args={[rail, STRINGER_T, STEP_D * 0.8]} />
        <meshStandardMaterial color="#64748b" />
      </mesh>
      <mesh
        position={[0, RAIL_Y, STEP_D / 2]}
        rotation={[0, 0, Math.atan2(height, width)]}
      >
        <boxGeometry args={[rail, RAIL_T, RAIL_T]} />
        <meshStandardMaterial color="#64748b" />
      </mesh>
    </group>
  );
}
