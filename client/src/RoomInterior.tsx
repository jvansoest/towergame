// One room, drawn open-front: shell, art, and fittings.

import { roomArt } from "./roomArt";
import { lobbyArt } from "./lobbyArt";
import { DAY, FLOOR_T, NIGHT, WALL_T } from "./dims";
import { CONCRETE, FLOOR, WALL, shadeHex } from "./roomTint";
import { BackWall, SideWall } from "./RoomShell";
import Furniture, { Seating } from "./Furniture";
import SleepSign, { BED_X } from "./SleepSign";

type Props = {
  width: number;
  height: number;
  depth: number;
  category: string;
  type: string;
  seats: number; // desks or chairs to draw
  lit: boolean;
  occupants: number; // people in the room now
};

// Open-front room shell with furniture.
export default function RoomInterior({
  width,
  height,
  depth,
  category,
  type,
  seats,
  lit,
  occupants,
}: Props) {
  // The tint multiplies into every surface.
  const shade = lit ? DAY : NIGHT;
  const wall = shadeHex(WALL[category] ?? "#e5e7eb", shade);
  const floor = shadeHex(FLOOR[category] ?? CONCRETE, shade);
  const floorTop = -height / 2 + FLOOR_T;
  const art = roomArt(type, width, height);
  const lobbyWall = category === "lobby" ? lobbyArt(width, height) : null;
  const sleeping = !lit && occupants > 0 && BED_X[type] !== undefined;

  // Painted rooms bring their own furniture, not their slab.
  if (art) {
    return (
      <group>
        <mesh position={[0, 0, -depth / 2 + 0.15]}>
          <planeGeometry args={[width, height]} />
          <meshBasicMaterial map={art} color={shade} toneMapped={false} />
        </mesh>
        <mesh position={[0, -height / 2 + FLOOR_T / 2, 0]}>
          <boxGeometry args={[width, FLOOR_T, depth]} />
          <meshStandardMaterial color={floor} />
        </mesh>
        <Seating
          seats={seats}
          width={width}
          category={category}
          floorTop={floorTop}
          shade={shade}
        />
        <SideWall
          x={-width / 2 + WALL_T / 2}
          innerFace={0}
          h={height}
          depth={depth}
          tint={wall}
        />
        <SideWall
          x={width / 2 - WALL_T / 2}
          innerFace={1}
          h={height}
          depth={depth}
          tint={wall}
        />
        {sleeping && (
          <SleepSign
            type={type}
            width={width}
            height={height}
            shade={shade}
          />
        )}
      </group>
    );
  }

  return (
    <group>
      {lobbyWall ? (
        <mesh position={[0, 0, -depth / 2 + 0.15]}>
          <planeGeometry args={[width, height]} />
          <meshBasicMaterial map={lobbyWall} color={shade} toneMapped={false} />
        </mesh>
      ) : (
        <BackWall
          width={width}
          height={height}
          depth={depth}
          tint={wall}
          windows={category !== "lobby"}
        />
      )}
      <mesh position={[0, -height / 2 + FLOOR_T / 2, 0]}>
        <boxGeometry args={[width, FLOOR_T, depth]} />
        <meshStandardMaterial color={floor} />
      </mesh>
      <SideWall x={-width / 2 + WALL_T / 2} innerFace={0} h={height} depth={depth} tint={wall} />
      <SideWall x={width / 2 - WALL_T / 2} innerFace={1} h={height} depth={depth} tint={wall} />
      <Furniture
        category={category}
        width={width}
        height={height}
        floorTop={floorTop}
      />
    </group>
  );
}
