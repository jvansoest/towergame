// The painting toolkit: canvas textures, and the
// pieces the room art is drawn from.

import * as THREE from "three";

// Pixels per world unit in the painted art.
export const PPU = 22;

export type Draw = (g: CanvasRenderingContext2D, w: number, h: number) => void;

// Paints one texture, sized in world units.
export function paintTexture(
  width: number,
  height: number,
  draw: Draw
): THREE.CanvasTexture | null {
  const cv = document.createElement("canvas");
  cv.width = Math.round(width * PPU);
  cv.height = Math.round(height * PPU);
  const g = cv.getContext("2d");
  if (!g) return null;
  g.imageSmoothingEnabled = false;
  draw(g, cv.width, cv.height);

  const tex = new THREE.CanvasTexture(cv);
  tex.magFilter = THREE.NearestFilter; // keep the pixel edges
  tex.minFilter = THREE.LinearMipmapLinearFilter;
  tex.colorSpace = THREE.SRGBColorSpace;
  return tex;
}

// Fills a rectangle in pixels.
export function box(
  g: CanvasRenderingContext2D,
  x: number,
  y: number,
  w: number,
  h: number,
  color: string
) {
  g.fillStyle = color;
  g.fillRect(Math.round(x), Math.round(y), Math.round(w), Math.round(h));
}

// A bed seen from the side, headboard left.
export function bed(g: CanvasRenderingContext2D, x: number, base: number, s: number) {
  box(g, x, base - 13 * s, 4 * s, 13 * s, "#7a4a2b"); // headboard
  box(g, x, base - 7 * s, 26 * s, 3 * s, "#d9d2c4"); // mattress
  box(g, x + 4 * s, base - 9 * s, 7 * s, 2 * s, "#f4f1ea"); // pillow
  box(g, x + 11 * s, base - 8 * s, 15 * s, 4 * s, "#a83246"); // blanket
  box(g, x, base - 4 * s, 26 * s, 4 * s, "#5c3620"); // frame
}

// A wardrobe with two doors.
export function wardrobe(
  g: CanvasRenderingContext2D,
  x: number,
  base: number,
  s: number
) {
  box(g, x, base - 30 * s, 16 * s, 30 * s, "#6b4226");
  box(g, x + 1 * s, base - 28 * s, 6 * s, 26 * s, "#8b5a33");
  box(g, x + 9 * s, base - 28 * s, 6 * s, 26 * s, "#8b5a33");
  box(g, x + 7 * s, base - 17 * s, 2 * s, 2 * s, "#e8d9a0");
}

// A sofa facing the viewer.
export function sofa(g: CanvasRenderingContext2D, x: number, base: number, s: number) {
  box(g, x, base - 13 * s, 22 * s, 8 * s, "#3f6f5e");
  box(g, x, base - 7 * s, 22 * s, 7 * s, "#4e876f");
  box(g, x, base - 15 * s, 3 * s, 15 * s, "#35604f");
  box(g, x + 19 * s, base - 15 * s, 3 * s, 15 * s, "#35604f");
}

// A low table with two chairs.
export function diningSet(
  g: CanvasRenderingContext2D,
  x: number,
  base: number,
  s: number
) {
  box(g, x + 4 * s, base - 12 * s, 18 * s, 2 * s, "#8b5a33");
  box(g, x + 6 * s, base - 10 * s, 2 * s, 10 * s, "#6b4226");
  box(g, x + 18 * s, base - 10 * s, 2 * s, 10 * s, "#6b4226");
  box(g, x, base - 16 * s, 3 * s, 16 * s, "#7a4a2b"); // chair back
  box(g, x, base - 9 * s, 6 * s, 2 * s, "#8b5a33");
  box(g, x + 23 * s, base - 16 * s, 3 * s, 16 * s, "#7a4a2b");
  box(g, x + 20 * s, base - 9 * s, 6 * s, 2 * s, "#8b5a33");
}

// Counter, cabinets, and a fridge.
export function kitchen(
  g: CanvasRenderingContext2D,
  x: number,
  base: number,
  s: number
) {
  box(g, x, base - 14 * s, 26 * s, 14 * s, "#b9b2a4"); // units
  box(g, x, base - 15 * s, 26 * s, 2 * s, "#5c5750"); // worktop
  box(g, x + 2 * s, base - 11 * s, 9 * s, 9 * s, "#cfc8ba");
  box(g, x + 14 * s, base - 11 * s, 9 * s, 9 * s, "#cfc8ba");
  box(g, x + 4 * s, base - 25 * s, 18 * s, 7 * s, "#a89f8f"); // wall units
  box(g, x + 30 * s, base - 26 * s, 12 * s, 26 * s, "#dedad2"); // fridge
  box(g, x + 40 * s, base - 18 * s, 1 * s, 5 * s, "#8a857c");
}

// A framed picture on the wall.
export function picture(
  g: CanvasRenderingContext2D,
  x: number,
  y: number,
  s: number
) {
  box(g, x, y, 10 * s, 8 * s, "#8b5a33");
  box(g, x + 1 * s, y + 1 * s, 8 * s, 6 * s, "#7fa7c4");
}

// A pot plant.
export function plant(g: CanvasRenderingContext2D, x: number, base: number, s: number) {
  box(g, x + 2 * s, base - 6 * s, 6 * s, 6 * s, "#8c5a3c");
  box(g, x + 3 * s, base - 14 * s, 4 * s, 8 * s, "#2f7d43");
  box(g, x, base - 12 * s, 10 * s, 4 * s, "#3f9142");
}

// A filing cabinet.
export function cabinet(
  g: CanvasRenderingContext2D,
  x: number,
  base: number,
  s: number
) {
  box(g, x, base - 20 * s, 11 * s, 20 * s, "#8f9aa5");
  for (let i = 0; i < 3; i++) {
    box(g, x + 1 * s, base - (18 - i * 6) * s, 9 * s, 5 * s, "#a7b2bd");
  }
}

// A band of windows near the ceiling.
export function windowBand(
  g: CanvasRenderingContext2D,
  w: number,
  h: number,
  frame: string
) {
  const top = h * 0.1;
  const tall = h * 0.32;
  box(g, 0, top - h * 0.03, w, h * 0.03, frame);
  const unit = Math.max(1, Math.round(w / Math.round(w / (h * 0.9))));
  for (let x = unit * 0.15; x < w - unit * 0.2; x += unit) {
    box(g, x, top, unit * 0.7, tall, "#8fc4e0");
    box(g, x, top, unit * 0.7, tall * 0.35, "#a9d7ef");
    box(g, x + unit * 0.7, top, unit * 0.3, tall, frame);
  }
}

// A serving counter with a menu board.
export function counter(
  g: CanvasRenderingContext2D,
  x: number,
  base: number,
  s: number,
  wide: number,
  face: string
) {
  box(g, x, base - 14 * s, wide, 14 * s, face);
  box(g, x, base - 16 * s, wide, 2 * s, "#e8e3d8"); // top
  box(g, x + 2 * s, base - 27 * s, wide - 4 * s, 8 * s, "#2b3440"); // menu
  for (let i = 0; i < 3; i++) {
    box(g, x + 4 * s, base - (25 - i * 3) * s, wide - 8 * s, 1 * s, "#e8c65a");
  }
}

// Shelving stacked with goods.
export function shelfUnit(
  g: CanvasRenderingContext2D,
  x: number,
  base: number,
  s: number,
  goods: string[]
) {
  box(g, x, base - 26 * s, 18 * s, 26 * s, "#8d8577");
  for (let row = 0; row < 3; row++) {
    const y = base - (24 - row * 8) * s;
    box(g, x + 1 * s, y, 16 * s, 6 * s, "#b3ab9c");
    for (let i = 0; i < 4; i++) {
      box(g, x + 2 * s + i * 4 * s, y + 1 * s, 3 * s, 5 * s, goods[(row + i) % goods.length]);
    }
  }
}

// Calls fn at a fixed step across a span.
export function acrossSpan(
  from: number,
  to: number,
  step: number,
  fn: (x: number) => void
) {
  for (let x = from; x + step <= to; x += step) fn(x);
}

// Wall, skirting, and floor for a home.
export function homeShell(
  g: CanvasRenderingContext2D,
  w: number,
  h: number,
  wall: string,
  floor: string
) {
  box(g, 0, 0, w, h, wall);
  const base = h * 0.86;
  box(g, 0, base, w, h - base, floor);
  box(g, 0, base - h * 0.03, w, h * 0.03, "#8a7f6d"); // skirting
  return base;
}
