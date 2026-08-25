import type { ThreeEvent } from "@react-three/fiber";

import { askInspect } from "./inspect";
import { useCars } from "./simsFeed";
import {
  CAR_D,
  CAR_H,
  CAR_W,
  CAR_Z,
  CELL_H,
  FLOOR_T,
  FRONT_Z,
  colToX,
} from "./dims";
import { useSelectedRoom, useHover, setHover, clearHover } from "./store";
import Glow from "./Glow";

// Every moving car, one or more per shaft.
export default function ElevatorCars({ gridWidth }: { gridWidth: number }) {
  const cars = useCars();
  const hover = useHover();
  const inspecting = useSelectedRoom() === "inspect";

  // Cars sharing a column share a shaft.
  const rank = new Map<number, number>();
  return (
    <>
      {cars.map((car) => {
        const n = rank.get(car.col) ?? 0;
        rank.set(car.col, n + 1);
        // A hovered shaft lights up all of its cars.
        const hot =
          (hover?.kind === "car" && hover.id === car.id) ||
          (hover?.kind === "shaft" && hover.id === car.shaft);
        return (
          <group
            key={car.id}
            position={[
              colToX(car.col, gridWidth),
              // Rest the car floor on the riders' feet.
              car.floor * CELL_H + FLOOR_T + CAR_H / 2,
              // Stacked cars stay apart in depth.
              FRONT_Z + n * CAR_Z,
            ]}
            // Stop the move, or the plane behind wins the hover.
            onPointerMove={(e: ThreeEvent<PointerEvent>) => {
              if (!inspecting) return;
              e.stopPropagation();
              setHover({ kind: "car", id: car.id });
            }}
            onPointerOut={() => clearHover("car", car.id)}
            onPointerDown={(e: ThreeEvent<PointerEvent>) => {
              if (!inspecting || e.button !== 0) return;
              e.stopPropagation();
              askInspect("inspectcar", { id: car.id });
            }}
          >
            <mesh>
              <boxGeometry args={[CAR_W, CAR_H, CAR_D]} />
              <meshStandardMaterial
                color="#fbbf24"
                transparent
                opacity={0.55}
              />
            </mesh>
            {hot && <Glow width={CAR_W} height={CAR_H} depth={CAR_D} />}
          </group>
        );
      })}
    </>
  );
}
