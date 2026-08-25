// What stands inside a room: seating, then the fittings
// each category brings with it.

import type { JSX } from "react";

import { propArt } from "./propArt";
import { CELL_W, PROP_Z } from "./dims";
import { seatCols, tile } from "./layout";

// Flat chairs and tables, set in front of the wall.
export function Seating({
  seats,
  width,
  category,
  floorTop,
  shade,
}: {
  seats: number;
  width: number;
  category: string;
  floorTop: number;
  shade: string;
}) {
  const desks = category === "office";
  const chair = propArt("chair");
  const table = propArt(desks ? "workdesk" : "table");
  if (seats < 1 || !chair || !table) return null;

  const cells = Math.round(width / CELL_W);
  const cols = seatCols(cells, seats);
  // Local x of a cell centre, relative to the room.
  const at = (c: number) => (c + 0.5 - cells / 2) * CELL_W;

  // A desk each, or a table shared by two seats.
  const tables: number[] = [];
  for (let i = 0; i < cols.length; i++) {
    if (desks) {
      tables.push(at(cols[i]));
    } else if (i > 0 && cols[i] - cols[i - 1] <= 3) {
      tables.push((at(cols[i - 1]) + at(cols[i])) / 2);
    }
  }

  return (
    <>
      {tables.map((x, i) => (
        <mesh key={`t${i}`} position={[x, floorTop + table.h / 2, PROP_Z - 0.1]}>
          <planeGeometry args={[table.w, table.h]} />
          <meshBasicMaterial
            map={table.tex}
            color={shade}
            transparent
            alphaTest={0.5}
            toneMapped={false}
          />
        </mesh>
      ))}
      {cols.map((c, i) => (
        <mesh key={`c${i}`} position={[at(c), floorTop + chair.h / 2, PROP_Z]}>
          <planeGeometry args={[chair.w, chair.h]} />
          <meshBasicMaterial
            map={chair.tex}
            color={shade}
            transparent
            alphaTest={0.5}
            toneMapped={false}
          />
        </mesh>
      ))}
    </>
  );
}

type FurnitureProps = {
  category: string;
  width: number;
  height: number;
  floorTop: number;
};

export default function Furniture({
  category,
  width,
  height,
  floorTop,
}: FurnitureProps): JSX.Element {
  switch (category) {
    case "lobby":
      return (
        <>
          {tile(width, 4).map((x, i) => (
            <group key={i} position={[x, floorTop, 0.3]}>
              <mesh position={[0, 0.3, 0]}>
                <cylinderGeometry args={[0.2, 0.28, 0.35, 8]} />
                <meshStandardMaterial color="#8b5e3c" />
              </mesh>
              <mesh position={[0, 0.7, 0]}>
                <sphereGeometry args={[0.38, 10, 10]} />
                <meshStandardMaterial color="#3f9142" />
              </mesh>
            </group>
          ))}
          {tile(width, 6).map((x, i) => (
            <mesh key={`b${i}`} position={[x, floorTop + 0.2, 0.9]}>
              <boxGeometry args={[1.6, 0.2, 0.5]} />
              <meshStandardMaterial color="#9aa5b1" />
            </mesh>
          ))}
        </>
      );
    case "office":
      return (
        <>
          {tile(width, 2.6).map((x, i) => (
            <group key={i} position={[x, floorTop, 0.4]}>
              <mesh position={[0, 0.28, 0]}>
                <boxGeometry args={[1.4, 0.15, 0.9]} />
                <meshStandardMaterial color="#8b5e3c" />
              </mesh>
              <mesh position={[0, 0.6, -0.3]}>
                <boxGeometry args={[0.7, 0.45, 0.08]} />
                <meshStandardMaterial color="#1f2937" />
              </mesh>
              <mesh position={[0, 0.2, 0.5]}>
                <boxGeometry args={[0.4, 0.4, 0.4]} />
                <meshStandardMaterial color="#334155" />
              </mesh>
            </group>
          ))}
        </>
      );
    case "hotel":
      return (
        <group position={[0, floorTop, 0.2]}>
          <mesh position={[0, 0.25, 0]}>
            <boxGeometry args={[2, 0.4, 1.4]} />
            <meshStandardMaterial color="#4f7cac" />
          </mesh>
          <mesh position={[-0.6, 0.55, -0.3]}>
            <boxGeometry args={[0.7, 0.25, 0.5]} />
            <meshStandardMaterial color="#eeeeee" />
          </mesh>
        </group>
      );
    case "residential":
      return (
        <>
          {tile(width, 3.2).map((x, i) => (
            <group key={i} position={[x, floorTop, 0.2]}>
              <mesh position={[0, 0.25, 0]}>
                <boxGeometry args={[2, 0.4, 1.4]} />
                <meshStandardMaterial color="#4f7cac" />
              </mesh>
              <mesh position={[-0.6, 0.55, -0.3]}>
                <boxGeometry args={[0.7, 0.25, 0.5]} />
                <meshStandardMaterial color="#eeeeee" />
              </mesh>
            </group>
          ))}
        </>
      );
    case "food":
      return (
        <>
          {tile(width, 2.2).map((x, i) => (
            <group key={i} position={[x, floorTop, 0.2]}>
              <mesh position={[0, 0.5, 0]}>
                <cylinderGeometry args={[0.5, 0.5, 0.1, 16]} />
                <meshStandardMaterial color="#c05621" />
              </mesh>
              <mesh position={[0, 0.25, 0]}>
                <cylinderGeometry args={[0.08, 0.08, 0.5, 8]} />
                <meshStandardMaterial color="#7b341e" />
              </mesh>
            </group>
          ))}
        </>
      );
    case "retail":
      return (
        <>
          {tile(width, 1.8).map((x, i) => (
            <mesh key={i} position={[x, floorTop + height * 0.22, -0.4]}>
              <boxGeometry args={[1.2, Math.min(height * 0.6, 2), 0.5]} />
              <meshStandardMaterial color="#a78bfa" />
            </mesh>
          ))}
        </>
      );
    case "entertainment":
      return (
        <>
          <mesh position={[0, height * 0.1, -0.1]}>
            <boxGeometry args={[width * 0.8, height * 0.55, 0.2]} />
            <meshStandardMaterial
              color="#111827"
              emissive="#1e3a8a"
              emissiveIntensity={0.5}
            />
          </mesh>
          {tile(width, 2).flatMap((x, i) =>
            [0, 1].map((row) => (
              <mesh
                key={`${i}-${row}`}
                position={[x, floorTop + 0.2, 0.6 + row * 0.7]}
              >
                <boxGeometry args={[1.2, 0.4, 0.5]} />
                <meshStandardMaterial color="#7f1d1d" />
              </mesh>
            ))
          )}
        </>
      );
    default:
      return <></>;
  }
}
