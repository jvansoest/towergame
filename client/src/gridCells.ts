// Which cells hold base structure, above and below ground.

import type { GridView } from "./protocol";

// Floors -1, -2, ... live in `below`, the rest in `built`.
export function builtAt(grid: GridView, floor: number, col: number): boolean {
  if (floor < 0) return grid.below?.[-floor - 1]?.[col] ?? false;
  return grid.built[floor]?.[col] ?? false;
}
