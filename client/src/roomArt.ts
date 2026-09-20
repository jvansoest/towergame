// Painted interiors, one shared texture per room type.
// Drawn once, then reused by every room of that type.

import type * as THREE from "three";

import {
  acrossSpan,
  bed,
  box,
  cabinet,
  diningSet,
  homeShell,
  picture,
  plant,
  shelfUnit,
  sofa,
  wardrobe,
  nightstand,
  rug,
  kitchen,
  windowBand,
  glass,
  fireplace,
  tv,
  paintTexture,
  type Draw,
} from "./paint";
import {
  BURGER_WALL,
  ICECREAM_WALL,
  IZAKAYA_WALL,
  NOODLE_WALL,
  PIZZA_WALL,
  stampPixels,
} from "./pixelArt";

const cache = new Map<string, THREE.CanvasTexture | null>();

// How many looks a type has.
export function variantCount(type: string): number {
  const p = PAINTERS[type];
  return Array.isArray(p) ? p.length : 1;
}

// Back-wall art for a room type, or null. The
// variant picks one look; any integer works.
export function roomArt(
  type: string,
  width: number,
  height: number,
  variant = 0,
): THREE.CanvasTexture | null {
  const p = PAINTERS[type];
  const pick = Array.isArray(p) ? Math.abs(variant) % p.length : 0;
  const key = `${type}#${pick}`;
  const hit = cache.get(key);
  if (hit !== undefined) return hit;

  const draw = Array.isArray(p) ? p[pick] : p;
  const tex = draw ? paintTexture(width, height, draw) : null;
  cache.set(key, tex);
  return tex;
}

const PAINTERS: Record<string, Draw | Draw[]> = {
  // Five flats, each a different look.
  condo: [
    // Bed, wardrobe, TV over a fireplace.
    (g, w, h) => {
      const base = homeShell(g, w, h, "#efe3cd", "#a9764a");
      const s = h / 34;
      picture(g, 5 * s, h * 0.16, s);
      bed(g, 1 * s, base, s);
      nightstand(g, 27 * s, base, s);
      wardrobe(g, 36 * s, base, s);
      rug(g, 46 * s, base, s, "#a83246");
      tv(g, 57 * s, 4 * s, s);
      fireplace(g, 57 * s, base, s);
      plant(g, 73 * s, base, s);
    },
    // Mirrored: fireplace left, bed right.
    (g, w, h) => {
      const base = homeShell(g, w, h, "#dfe8ee", "#8f7b66");
      const s = h / 34;
      tv(g, 6 * s, 4 * s, s);
      fireplace(g, 6 * s, base, s);
      plant(g, 22 * s, base, s);
      rug(g, 4 * s, base, s, "#3f6f8f");
      wardrobe(g, 30 * s, base, s);
      picture(g, 50 * s, h * 0.16, s);
      bed(g, 56 * s, base, s);
      nightstand(g, 82 * s - 3 * s, base, s);
    },
    // A TV wall with a sofa; the bed sits apart.
    (g, w, h) => {
      const base = homeShell(g, w, h, "#e6efdc", "#a08466");
      const s = h / 34;
      bed(g, 1 * s, base, s);
      nightstand(g, 27 * s, base, s);
      picture(g, 6 * s, h * 0.15, s);
      tv(g, 47 * s, 5 * s, s);
      box(g, 46 * s, base - 5 * s, 16 * s, 5 * s, "#4a3a2c"); // TV unit
      sofa(g, 62 * s, base, s);
      plant(g, 39 * s, base, s);
    },
    // A cosy flat: fireplace, two pictures.
    (g, w, h) => {
      const base = homeShell(g, w, h, "#f0dcd8", "#9a6b50");
      const s = h / 34;
      picture(g, 6 * s, h * 0.14, s);
      picture(g, 20 * s, h * 0.2, s);
      bed(g, 1 * s, base, s);
      nightstand(g, 27 * s, base, s);
      fireplace(g, 46 * s, base, s);
      picture(g, 47 * s, 6 * s, s);
      rug(g, 40 * s, base, s, "#7a4a6a");
      cabinet(g, 66 * s, base, s);
      plant(g, 76 * s, base, s);
    },
    // A studio: kitchenette and a small TV.
    (g, w, h) => {
      const base = homeShell(g, w, h, "#e8e0f0", "#8a7f95");
      const s = h / 34;
      bed(g, 1 * s, base, s);
      nightstand(g, 27 * s, base, s);
      tv(g, 20 * s, 3 * s, s * 0.8);
      kitchen(g, 40 * s, base, s);
      plant(g, 76 * s, base, s);
    },
  ],

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

  // Cream wall, framed art, green wainscot
  // with red pilasters. Sets are props.
  restaurant_indian: (g, w, h) => {
    const base = homeShell(g, w, h, "#e8dcc0", "#8a2b2b");
    const s = h / 34;

    // Maroon frieze with small gold arches.
    box(g, 0, 0, w, 4 * s, "#4a1a14");
    acrossSpan(2 * s, w, 8 * s, (x) => box(g, x, s, 5 * s, 2 * s, "#c9a227"));

    // Framed artworks along the upper wall.
    acrossSpan(10 * s, w - 14 * s, 26 * s, (x) => {
      box(g, x, base - 26 * s, 14 * s, 12 * s, "#2f5d46");
      box(g, x + 1.5 * s, base - 24.5 * s, 11 * s, 9 * s, "#d9b23a");
    });

    // Green wainscot with red pilasters.
    box(g, 0, base - 13 * s, w, 13 * s, "#1e4d3a");
    acrossSpan(6 * s, w, 18 * s, (x) =>
      box(g, x, base - 13 * s, 4 * s, 13 * s, "#8a2b2b"),
    );

    // Service station left, door right.
    box(g, 4 * s, base - 12 * s, 18 * s, 12 * s, "#e8e4da");
    box(g, 4 * s, base - 14 * s, 18 * s, 2 * s, "#2f7d43");
    box(g, w - 12 * s, base - 30 * s, 9 * s, 30 * s, "#6b4226");
    box(g, w - 11 * s, base - 28 * s, 7 * s, 26 * s, "#8b5a33");
  },

  // Evening dining: glowing ceiling panels,
  // swags on pilasters. Tables are props.
  restaurant_french: (g, w, h) => {
    const base = homeShell(g, w, h, "#e3d9c2", "#41595c");
    const s = h / 34;

    // Dark ceiling with wide glowing panels.
    box(g, 0, 0, w, 7 * s, "#1d2420");
    acrossSpan(2 * s, w - 4 * s, 16 * s, (x) =>
      box(g, x, 1.5 * s, 13 * s, 4.5 * s, "#f7edc4"),
    );

    // Cornice, then swags with scalloped hems.
    box(g, 0, 7 * s, w, 2 * s, "#b9a77e");
    acrossSpan(0, w - 12 * s, 14 * s, (x) => {
      box(g, x - 2 * s, 9 * s, 17 * s, 4 * s, "#f6f1e4");
      g.fillStyle = "#f6f1e4";
      for (let k = 0; k < 4; k++) {
        g.beginPath();
        g.arc(x + 1.5 * s + k * 4 * s, 13 * s, 2 * s, 0, Math.PI);
        g.fill();
      }
      box(g, x - 2 * s, 13 * s, 17 * s, 1.2 * s, "#cfc4a6");
    });

    // Slim pilasters down to the floor.
    acrossSpan(0, w, 20 * s, (x) =>
      box(g, x, 9 * s, 2 * s, base - 9 * s, "#5c5240"),
    );

    // A pink sideboard at the right end.
    box(g, w - 14 * s, base - 16 * s, 12 * s, 16 * s, "#c96a7a");
    box(g, w - 15 * s, base - 18 * s, 14 * s, 2 * s, "#e8e4da");
  },

  // The wall is hand-authored pixel art;
  // counter and dining sets are props.
  icecream: (g, w, h) => {
    stampPixels(g, ICECREAM_WALL, w, h);
  },

  // Windows up top, a bar at one end; the
  // tables are props.
  restaurant: (g, w, h) => {
    const base = homeShell(g, w, h, "#7d2f38", "#5a2027");
    const s = h / 34;
    box(g, 0, h * 0.06, w, h * 0.05, "#c9a227"); // cornice

    windowBand(g, w, h, "#c3cbc8");

    const barW = 56 * s;
    const barX = w - barW - 2 * s;
    plant(g, barX - 12 * s, base, s);

    box(g, barX, base - 16 * s, barW, 16 * s, "#4a2a1c"); // bar
    box(g, barX, base - 18 * s, barW, 2 * s, "#6b4226");
    box(g, barX + 4 * s, base - 30 * s, barW - 8 * s, 9 * s, "#3a2016");
    acrossSpan(barX + 6 * s, barX + barW - 6 * s, 6 * s, (x) =>
      box(g, x, base - 28 * s, 3 * s, 6 * s, "#8fbf7f"),
    );
  },

  // A station hall: tiled walls, signs, a platform edge.
  subway: (g, w, h) => {
    box(g, 0, 0, w, h, "#d8d4c8"); // tiled wall
    for (let x = 0; x < w; x += 6) box(g, x, 0, 1, h, "#c4c0b3"); // grout
    box(g, 0, h * 0.02, w, h * 0.09, "#1f3a5f"); // sign band
    const s = h / 68;
    for (let x = 6 * s; x < w - 20 * s; x += 32 * s) {
      box(g, x, h * 0.035, 16 * s, 5 * s, "#f4f1ea"); // METRO plate
      box(g, x + 2 * s, h * 0.035 + 1.5 * s, 12 * s, 2 * s, "#c0392b");
    }
    box(g, 0, h * 0.5, w, h * 0.025, "#8f97a3"); // mezzanine edge
    box(g, 0, h * 0.93, w, h * 0.07, "#33373f"); // platform
    box(g, 0, h * 0.925, w, h * 0.012, "#f2d675"); // safety line
    for (const x of [0.12, 0.38, 0.62, 0.88]) {
      box(g, w * x - 2 * s, h * 0.11, 4 * s, h * 0.82, "#9aa1ac"); // pillars
    }
  },

  // A garage bay: concrete, stall lines, a P sign.
  parking: (g, w, h) => {
    box(g, 0, 0, w, h, "#565c66"); // wall
    const base = h * 0.86;
    box(g, 0, base, w, h - base, "#33373f"); // floor
    box(g, 0, h * 0.06, w, h * 0.05, "#2c3037"); // beam
    const s = h / 34;
    box(g, w * 0.4, h * 0.16, 9 * s, 9 * s, "#2f6fa8"); // sign
    box(g, w * 0.4 + 3 * s, h * 0.16 + 1.5 * s, 2 * s, 6 * s, "#f4f1ea");
    box(g, w * 0.4 + 3 * s, h * 0.16 + 1.5 * s, 4 * s, 2 * s, "#f4f1ea");
    box(g, w * 0.4 + 6 * s, h * 0.16 + 1.5 * s, 1 * s, 4 * s, "#f4f1ea");
    for (const x of [0.02, 0.5, 0.98]) {
      box(g, w * x - s, base - 1 * s, 2 * s, 1 * s, "#f2d675"); // stall lines
    }
    box(g, 0, h * 0.12, 3 * s, base - h * 0.12, "#3d434c"); // pillar
    box(g, w - 3 * s, h * 0.12, 3 * s, base - h * 0.12, "#3d434c");
  },

  // A clinic: a red cross, two beds, a cabinet.
  medical: (g, w, h) => {
    const base = homeShell(g, w, h, "#e6f2f0", "#7fa8a3");
    const s = h / 34;
    windowBand(g, w, h, "#c3cbc8");
    const cx = 44 * s;
    box(g, cx, h * 0.14, 9 * s, 9 * s, "#fbf6ee"); // sign
    box(g, cx + 3.5 * s, h * 0.14 + 1 * s, 2 * s, 7 * s, "#c0392b");
    box(g, cx + 1 * s, h * 0.14 + 3.5 * s, 7 * s, 2 * s, "#c0392b");
    bed(g, 2 * s, base, s);
    bed(g, 31 * s, base, s);
    cabinet(g, 62 * s, base, s);
    cabinet(g, 74 * s, base, s);
    box(g, 88 * s, base - 12 * s, 12 * s, 12 * s, "#dedad2"); // desk
  },

  // Towel shelves and a maid's cart.
  housekeeping: (g, w, h) => {
    const base = homeShell(g, w, h, "#e3ecef", "#6f8a93");
    const s = h / 34;
    box(g, 0, h * 0.05, w, h * 0.05, "#8fa9b2"); // rail
    for (const x of [6 * s, 30 * s, 54 * s]) {
      shelfUnit(g, x, base, s, ["#f4f1ea", "#8fc4e0", "#e8d9a0"]);
    }
    box(g, 80 * s, base - 14 * s, 16 * s, 10 * s, "#8f9aa5"); // cart
    box(g, 81 * s, base - 22 * s, 14 * s, 8 * s, "#f4f1ea"); // towels
    box(g, 82 * s, base - 4 * s, 3 * s, 4 * s, "#3a3f47"); // wheels
    box(g, 91 * s, base - 4 * s, 3 * s, 4 * s, "#3a3f47");
    box(g, 96 * s, base - 20 * s, 2 * s, 20 * s, "#6b4226"); // mop
  },

  // Counter joints share the ice cream layout.
  burger: (g, w, h) => stampPixels(g, BURGER_WALL, w, h),
  pizza: (g, w, h) => stampPixels(g, PIZZA_WALL, w, h),
  noodle: (g, w, h) => stampPixels(g, NOODLE_WALL, w, h),

  // Hand-painted wall; the bar and stools
  // are props.
  izakaya: (g, w, h) => {
    stampPixels(g, IZAKAYA_WALL, w, h);
  },

  // Windows along the wall; the service
  // counter is a prop out front.
  fastfood: (g, w, h) => {
    homeShell(g, w, h, "#f2e3c8", "#c96a2a");
    box(g, 0, h * 0.04, w, h * 0.06, "#d94f34"); // fascia

    windowBand(g, w, h, "#c3cbc8");
  },

  // Two small windows flank the shelves.
  shop: (g, w, h) => {
    const base = homeShell(g, w, h, "#efe6f6", "#6f6480");
    const s = h / 34;
    box(g, 0, h * 0.05, w, h * 0.07, "#3f6fa8"); // awning

    const winW = 16 * s;
    const winH = 12 * s;
    const winY = base - 24 * s;
    for (const x of [5 * s, w - 21 * s]) {
      box(g, x, winY, winW, winH, "#c3cbc8"); // frame
      box(
        g,
        x + 1.5 * s,
        winY + 1.5 * s,
        winW - 3 * s,
        winH - 3 * s,
        "#8fc4e0",
      ); // pane
      box(
        g,
        x + 1.5 * s,
        winY + 1.5 * s,
        winW - 3 * s,
        (winH - 3 * s) * 0.45,
        "#a9d7ef",
      ); // glint
    }

    const goods = ["#d4574e", "#4e8fd4", "#e8c65a", "#59a86b", "#b06fc4"];
    acrossSpan(26 * s, w - 26 * s, 22 * s, (x) =>
      shelfUnit(g, x, base, s, goods),
    );
  },

  // One bed, nightstand, and a window.
  hotel_single: (g, w, h) => {
    const base = homeShell(g, w, h, "#e6dced", "#8a6f9c");
    const s = h / 34;
    picture(g, w * 0.08, h * 0.16, s);

    bed(g, w * 0.06, base, s);
    box(g, w * 0.06 + 27 * s, base - 9 * s, 7 * s, 9 * s, "#6b4226");
    box(g, w * 0.06 + 29 * s, base - 13 * s, 3 * s, 4 * s, "#f2d675");
    glass(g, w * 0.72, h * 0.2, w * 0.22, h * 0.4);
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
    glass(g, w * 0.78, h * 0.2, w * 0.18, h * 0.4);
  },
};
