// What stands inside a room: seating, then the fittings
// each category brings with it.

import type { JSX } from "react";

import { propArt } from "./propArt";
import {
  BAR_H,
  BAR_W,
  BAR_Z_BACK,
  BAR_Z_FRONT,
  CELL_W,
  PROP_Z,
  SIM_H,
  SEAT_H,
  TABLE_Z,
} from "./dims";
import { barCol, barZ, chairKind, seatCols, tableCols, tile } from "./layout";
import { barLook, shadeHex } from "./roomTint";

// A staffed service counter: flat art pulled
// well forward of the wall line, like the
// other fittings. Shops keep it east, fast
// food west. A worker stands behind each one.
export function ServiceCounter({
  kind,
  width,
  floorTop,
  shade,
}: {
  kind: string;
  width: number;
  floorTop: number;
  shade: string;
}) {
  const counter = propArt(kind);
  if (!counter) return null;
  const { x, z } = counterSpot(kind, width);
  return (
    <mesh position={[x, floorTop + counter.h / 2, z]}>
      <planeGeometry args={[counter.w, counter.h]} />
      <meshBasicMaterial
        map={counter.tex}
        color={shade}
        transparent
        alphaTest={0.5}
        toneMapped={false}
      />
    </mesh>
  );
}

// Where a counter sits, relative to the room's
// centre. The worker behind it shares this spot.
export function counterSpot(
  kind: string,
  width: number,
): { x: number; z: number } {
  const counter = propArt(kind);
  if (!counter) return { x: 0, z: 0 };
  // Clear of the side wall by half a unit.
  const edge = width / 2 - 0.5 - counter.w / 2;
  const x = kind === "fastcounter" ? -edge : edge;
  // Further forward than tables, not into walkers.
  return { x, z: 0.8 };
}

// Litter on a room's floor, until it is cleaned.
export function Mess({ floorTop, shade }: { floorTop: number; shade: string }) {
  const mess = propArt("mess");
  if (!mess) return null;
  return (
    <mesh position={[0, floorTop + mess.h / 2, PROP_Z + 0.2]}>
      <planeGeometry args={[mess.w, mess.h]} />
      <meshBasicMaterial
        map={mess.tex}
        color={shade}
        transparent
        alphaTest={0.5}
        toneMapped={false}
      />
    </mesh>
  );
}

// Local x of seats, counter, server.
export function barSpot(width: number) {
  const cells = Math.round(width / CELL_W);
  const at = (c: number) => (c + 0.5 - cells / 2) * CELL_W;
  const col = barCol(cells);
  return {
    seatX: at(col),
    counterX: at(col) + CELL_W / 2,
    serverX: at(col + 1),
  };
}

// A counter running front to back,
// with a stool per seat beside it.
export function Bar({
  type,
  seats,
  width,
  floorTop,
  shade,
}: {
  type: string;
  seats: number;
  width: number;
  floorTop: number;
  shade: string;
}) {
  const [front, top, seat] = barLook(type);
  const { seatX, counterX } = barSpot(width);
  const len = BAR_Z_FRONT - BAR_Z_BACK + 0.3;
  const zMid = (BAR_Z_FRONT + BAR_Z_BACK) / 2;
  const seatH = SIM_H * SEAT_H;
  return (
    <>
      <mesh position={[counterX, floorTop + (BAR_H - 0.1) / 2, zMid]}>
        <boxGeometry args={[BAR_W, BAR_H - 0.1, len]} />
        <meshBasicMaterial color={shadeHex(front, shade)} toneMapped={false} />
      </mesh>
      <mesh position={[counterX, floorTop + BAR_H - 0.05, zMid]}>
        <boxGeometry args={[BAR_W + 0.1, 0.1, len + 0.05]} />
        <meshBasicMaterial color={shadeHex(top, shade)} toneMapped={false} />
      </mesh>
      {Array.from({ length: seats }, (_, i) => (
        <mesh key={i} position={[seatX, floorTop + seatH / 2, barZ(i, seats)]}>
          <cylinderGeometry args={[0.15, 0.15, seatH, 10]} />
          <meshBasicMaterial color={shadeHex(seat, shade)} toneMapped={false} />
        </mesh>
      ))}
    </>
  );
}

// Which service counter a room type hosts,
// if any. Rooms and cards share this.
export function counterFor(type: string): string | undefined {
  if (type === "shop") return "counter";
  if (type === "fastfood") return "fastcounter";
  return undefined;
}

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

  const tables = tableCols(cells, seats, desks).map((c) => at(c));

  return (
    <>
      {tables.map((x, i) => (
        <mesh key={`t${i}`} position={[x, floorTop + table.h / 2, TABLE_Z]}>
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
      {cols.map((c, i) => {
        const seat = propArt(chairKind(i, desks));
        if (!seat) return null;
        return (
          <mesh key={`c${i}`} position={[at(c), floorTop + seat.h / 2, PROP_Z]}>
            <planeGeometry args={[seat.w, seat.h]} />
            <meshBasicMaterial
              map={seat.tex}
              color={shade}
              transparent
              alphaTest={0.5}
              toneMapped={false}
            />
          </mesh>
        );
      })}
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
            <mesh key={i} position={[x, floorTop + height * 0.22, TABLE_Z]}>
              <boxGeometry args={[1.2, Math.min(height * 0.6, 2), 0.5]} />
              <meshStandardMaterial color="#a78bfa" />
            </mesh>
          ))}
        </>
      );
    default:
      return <></>;
  }
}
