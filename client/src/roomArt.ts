// Painted interiors, one shared texture per room type.
// Drawn once, then reused by every room of that type.

import type * as THREE from "three";

import {
  acrossSpan,
  bed,
  box,
  cabinet,
  counter,
  diningSet,
  homeShell,
  kitchen,
  picture,
  plant,
  shelfUnit,
  sofa,
  wardrobe,
  windowBand,
  paintTexture,
  type Draw,
} from "./paint";

const cache = new Map<string, THREE.CanvasTexture | null>();

// Back-wall art for a room type, or null.
export function roomArt(
  type: string,
  width: number,
  height: number
): THREE.CanvasTexture | null {
  const hit = cache.get(type);
  if (hit !== undefined) return hit;

  const draw = PAINTERS[type];
  const tex = draw ? paintTexture(width, height, draw) : null;
  cache.set(type, tex);
  return tex;
}

const PAINTERS: Record<string, Draw> = {
  // Bedroom, living area, then kitchen.
  condo: (g, w, h) => {
    const base = homeShell(g, w, h, "#efe3cd", "#a9764a");
    const s = h / 34;

    bed(g, w * 0.03, base, s);
    box(g, w * 0.03 + 27 * s, base - 9 * s, 7 * s, 9 * s, "#6b4226"); // nightstand
    box(g, w * 0.03 + 29 * s, base - 13 * s, 3 * s, 4 * s, "#f2d675"); // lamp
    picture(g, w * 0.06, h * 0.18, s);

    wardrobe(g, w * 0.27, base, s);

    sofa(g, w * 0.4, base, s);
    plant(g, w * 0.37, base, s);
    picture(g, w * 0.44, h * 0.16, s);

    kitchen(g, w * 0.74, base, s);
  },

  // A break corner; the desks are props.
  office: (g, w, h) => {
    const base = homeShell(g, w, h, "#dfe4e2", "#7b8a86");
    const s = h / 34;
    windowBand(g, w, h, "#c3cbc8");

    const breakW = 44 * s;
    const breakX = w - breakW;
    cabinet(g, breakX, base, s);
    plant(g, breakX + 13 * s, base, s);
    sofa(g, breakX + 22 * s, base, s);
  },

  // A bar at one end; the tables are props.
  restaurant: (g, w, h) => {
    const base = homeShell(g, w, h, "#7d2f38", "#5a2027");
    const s = h / 34;
    box(g, 0, h * 0.06, w, h * 0.05, "#c9a227"); // cornice

    const barW = 56 * s;
    const barX = w - barW - 2 * s;
    plant(g, barX - 12 * s, base, s);

    box(g, barX, base - 16 * s, barW, 16 * s, "#4a2a1c"); // bar
    box(g, barX, base - 18 * s, barW, 2 * s, "#6b4226");
    box(g, barX + 4 * s, base - 30 * s, barW - 8 * s, 9 * s, "#3a2016");
    acrossSpan(barX + 6 * s, barX + barW - 6 * s, 6 * s, (x) =>
      box(g, x, base - 28 * s, 3 * s, 6 * s, "#8fbf7f")
    );
  },

  // Counter one end; the tables are props.
  fastfood: (g, w, h) => {
    const base = homeShell(g, w, h, "#f2e3c8", "#c96a2a");
    const s = h / 34;
    box(g, 0, h * 0.04, w, h * 0.06, "#d94f34"); // fascia

    const barW = Math.min(w * 0.28, 60 * s);
    counter(g, 3 * s, base, s, barW, "#d94f34");
  },

  // Shelves of goods and a till.
  shop: (g, w, h) => {
    const base = homeShell(g, w, h, "#efe6f6", "#6f6480");
    const s = h / 34;
    box(g, 0, h * 0.05, w, h * 0.07, "#3f6fa8"); // awning

    const goods = ["#d4574e", "#4e8fd4", "#e8c65a", "#59a86b", "#b06fc4"];
    const tillW = 26 * s;
    const tillX = w - tillW - 4 * s;
    acrossSpan(4 * s, tillX - 6 * s, 22 * s, (x) =>
      shelfUnit(g, x, base, s, goods)
    );
    box(g, tillX, base - 13 * s, tillW, 13 * s, "#6b4226"); // till
    box(g, tillX, base - 15 * s, tillW, 2 * s, "#8b5a33");
    box(g, tillX + 8 * s, base - 20 * s, 8 * s, 5 * s, "#2b3440");
  },

  // One bed, nightstand, and a window.
  hotel_single: (g, w, h) => {
    const base = homeShell(g, w, h, "#e6dced", "#8a6f9c");
    const s = h / 34;
    picture(g, w * 0.08, h * 0.16, s);

    bed(g, w * 0.06, base, s);
    box(g, w * 0.06 + 27 * s, base - 9 * s, 7 * s, 9 * s, "#6b4226");
    box(g, w * 0.06 + 29 * s, base - 13 * s, 3 * s, 4 * s, "#f2d675");
    box(g, w * 0.72, h * 0.2, w * 0.22, h * 0.4, "#8fc4e0"); // window
    box(g, w * 0.72, h * 0.2, w * 0.22, h * 0.14, "#a9d7ef");
  },

  // Bed, seating, and a table.
  hotel_suite: (g, w, h) => {
    const base = homeShell(g, w, h, "#e6dced", "#8a6f9c");
    const s = h / 34;
    picture(g, w * 0.06, h * 0.16, s);

    bed(g, w * 0.04, base, s);
    box(g, w * 0.04 + 27 * s, base - 9 * s, 7 * s, 9 * s, "#6b4226");
    box(g, w * 0.04 + 29 * s, base - 13 * s, 3 * s, 4 * s, "#f2d675");
    wardrobe(g, w * 0.3, base, s);
    sofa(g, w * 0.43, base, s);
    plant(g, w * 0.4, base, s);
    diningSet(g, w * 0.6, base, s);
    box(g, w * 0.78, h * 0.2, w * 0.18, h * 0.4, "#8fc4e0"); // window
    box(g, w * 0.78, h * 0.2, w * 0.18, h * 0.14, "#a9d7ef");
  },

};
