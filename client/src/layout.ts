// Where things sit across a room's width.

import { BAR_Z_BACK, BAR_Z_FRONT } from "./dims";

// Evenly spaced x positions across the width.
export function tile(width: number, spacing: number): number[] {
  const n = Math.max(1, Math.floor(width / spacing));
  const step = width / n;
  return Array.from({ length: n }, (_, i) => -width / 2 + step * (i + 0.5));
}

export function seatCols(cells: number, seats: number): number[] {
  return Array.from({ length: seats }, (_, i) =>
    Math.floor((cells * (i + 1)) / (seats + 1)),
  );
}

// The one column all bar seats share.
export function barCol(cells: number): number {
  return Math.floor(cells / 2);
}

// Depth of one bar seat.
export function barZ(spot: number, seats: number): number {
  const t = seats > 1 ? Math.min(spot, seats - 1) / (seats - 1) : 1;
  return BAR_Z_BACK + t * (BAR_Z_FRONT - BAR_Z_BACK);
}

// Tables between chairs, in cell columns.
// Offices put a desk before each seat; other
// seats pair up, a table between each two.
export function tableCols(
  cells: number,
  seats: number,
  desks: boolean,
): number[] {
  const cols = seatCols(cells, seats);
  const tabs: number[] = [];
  if (desks) {
    for (const c of cols) tabs.push(c);
    return tabs;
  }
  for (let i = 0; i + 1 < cols.length; i += 2) {
    tabs.push((cols[i] + cols[i + 1]) / 2);
  }
  return tabs;
}

// Which chair a seat draws: even seats sit on
// the table's left, facing right; odd seats
// face left across it. Offices all match.
export function chairKind(i: number, desks: boolean): string {
  if (desks || i % 2 === 0) return "chair";
  return "chairl";
}

// A stable number per room spot, for picking a look.
export function variantOf(floor: number, col: number): number {
  let h = (floor * 73856093) ^ (col * 19349663);
  h = Math.imul(h ^ (h >>> 13), 0x5bd1e995);
  return (h ^ (h >>> 15)) >>> 0;
}
