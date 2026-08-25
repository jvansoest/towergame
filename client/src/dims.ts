// One source for world geometry. Anything two modules
// must agree on lives here, not in either of them.

// Cell size in world units.
export const CELL_W = 1.5;
export const CELL_H = 3;
export const DEPTH = 3; // room depth, front to back

// Floor slab. Sims walk on top of it.
export const FLOOR_T = 0.3;
export const WALL_T = 0.15; // side wall, at a room's ends

// Depth layers, back to front.
export const PROP_Z = 0.5; // flat furniture
export const SIM_Z = 1; // people
export const FRONT_Z = 1.6; // stairs, shafts, cars

// Elevator shaft and its cars.
export const SHAFT_CELLS = 2; // matches model.ElevatorWidth
export const CAR_W = 2.2;
export const CAR_H = CELL_H * 0.85;
export const CAR_D = 0.5;
export const CAR_Z = 0.12; // depth between cars in one shaft

// A sim, and the chair it sits on.
export const SIM_H = 1.9;
export const HIP_STAND = 0.53; // hip height, share of SIM_H
export const SEAT_H = 0.25; // hip height once seated
export const SIT_FACING = Math.PI / 2; // chairs open to the right
export const CHAIR_H = 0.95;
// Seat height as a share of the chair, so sims land on it.
export const CHAIR_SEAT = (SIM_H * SEAT_H) / CHAIR_H;

// Shared tints. Unlit rooms keep a little moonlight.
export const DAY = "#ffffff";
export const NIGHT = "#4a5570";
export const HIGHLIGHT = "#fde68a"; // what the pointer is over

// Grid column to world x (cell center).
export function colToX(col: number, gridWidth: number): number {
  return (col + 0.5 - gridWidth / 2) * CELL_W;
}
