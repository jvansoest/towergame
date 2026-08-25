// Where things sit across a room's width.

// Evenly spaced x positions across the width.
export function tile(width: number, spacing: number): number[] {
  const n = Math.max(1, Math.floor(width / spacing));
  const step = width / n;
  return Array.from({ length: n }, (_, i) => -width / 2 + step * (i + 0.5));
}

export function seatCols(cells: number, seats: number): number[] {
  return Array.from({ length: seats }, (_, i) =>
    Math.floor((cells * (i + 1)) / (seats + 1))
  );
}
