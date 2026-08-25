// The z's that drift over a bed while sims sleep.

import { useMemo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import type * as THREE from "three";

import { propArt, ZZZ_FRAMES } from "./propArt";
import type { Prop } from "./propArt";
import { PROP_Z } from "./dims";

// Bed head, as a share of room width, per room type.
export const BED_X: Record<string, number> = {
  condo: 0.03,
  hotel_single: 0.06,
  hotel_suite: 0.04,
};

// Seconds one frame of the z's is shown.
const ZZZ_STEP = 0.45;

// Drifting z's above the bed while people sleep.
export default function SleepSign({
  type,
  width,
  height,
  shade,
}: {
  type: string;
  width: number;
  height: number;
  shade: string;
}) {
  const frames = useMemo(
    () => ZZZ_FRAMES.map((f) => propArt(f)).filter((p): p is Prop => !!p),
    []
  );
  const mat = useRef<THREE.MeshBasicMaterial>(null);
  const shown = useRef(-1);

  useFrame(({ clock }) => {
    if (frames.length === 0 || !mat.current) return;
    const i = Math.floor(clock.elapsedTime / ZZZ_STEP) % frames.length;
    if (i === shown.current) return;
    shown.current = i;
    mat.current.map = frames[i].tex;
  });

  if (frames.length === 0) return null;
  const art = frames[0];
  const x = -width / 2 + width * BED_X[type] + (7 * height) / 34;
  const y = -height * 0.095 + art.h / 2 + 0.35;
  return (
    <mesh position={[x, y, PROP_Z]}>
      <planeGeometry args={[art.w, art.h]} />
      <meshBasicMaterial
        ref={mat}
        map={art.tex}
        color={shade}
        transparent
        alphaTest={0.5}
        toneMapped={false}
      />
    </mesh>
  );
}
