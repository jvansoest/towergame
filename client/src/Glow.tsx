// A translucent box marking what the pointer is over.

import { HIGHLIGHT } from "./dims";

type Props = {
  width: number;
  height: number;
  depth: number;
  position?: [number, number, number];
  color?: string;
};

const PAD = 0.12; // stands slightly proud of the thing

export default function Glow({
  width,
  height,
  depth,
  position,
  color = HIGHLIGHT,
}: Props) {
  return (
    <mesh position={position ?? [0, 0, 0]} raycast={() => null}>
      <boxGeometry args={[width + PAD, height + PAD, depth + PAD]} />
      <meshBasicMaterial
        color={color}
        transparent
        opacity={0.28}
        depthWrite={false}
        toneMapped={false}
      />
    </mesh>
  );
}
