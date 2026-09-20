// Flat props that stand in front of the wall.
// Sized in world units, drawn with a clear background.

import type * as THREE from "three";

import { CHAIR_H, CHAIR_SEAT } from "./dims";
import { box, paintTexture, type Draw } from "./paint";
import { BAR_LOOK } from "./roomTint";

export type Prop = { w: number; h: number; tex: THREE.CanvasTexture };

const propCache = new Map<string, Prop | null>();

// One prop's art, or null if there is none.
export function propArt(kind: string): Prop | null {
  const hit = propCache.get(kind);
  if (hit !== undefined) return hit;

  const spec = PROPS[kind];
  const tex = spec ? paintTexture(spec.w, spec.h, spec.draw) : null;
  const prop = spec && tex ? { w: spec.w, h: spec.h, tex } : null;
  propCache.set(kind, prop);
  return prop;
}

type PropSpec = { w: number; h: number; draw: Draw };

const PROPS: Record<string, PropSpec> = {
  ...barProps(),
  // Used towels and litter on the floor.
  mess: {
    w: 1.4,
    h: 0.45,
    draw: (g, w, h) => {
      box(g, w * 0.05, h * 0.55, w * 0.3, h * 0.45, "#e8e2d0"); // towel
      box(g, w * 0.1, h * 0.4, w * 0.18, h * 0.2, "#d6cfba");
      box(g, w * 0.42, h * 0.7, w * 0.22, h * 0.3, "#6b4a2b"); // stain
      box(g, w * 0.7, h * 0.5, w * 0.2, h * 0.5, "#8c8a72"); // bin bag
      box(g, w * 0.76, h * 0.35, w * 0.08, h * 0.2, "#8c8a72");
      box(g, w * 0.55, h * 0.45, w * 0.08, h * 0.2, "#c9463d"); // can
    },
  },
  // A dining chair, seen from the side.
  chair: {
    w: 0.6,
    h: CHAIR_H,
    draw: (g, w, h) => {
      const seat = h * CHAIR_SEAT; // top of the seat
      box(g, 0, seat, w, h * 0.07, "#8b5a33"); // seat
      box(g, 0, h * 0.02, w * 0.18, seat, "#7a4a2b"); // back
      box(g, w * 0.18, h * 0.1, w * 0.12, seat * 0.7, "#a0663d");
      box(g, w * 0.06, seat + h * 0.07, w * 0.13, h * 0.43, "#6b4226");
      box(g, w * 0.8, seat + h * 0.07, w * 0.13, h * 0.43, "#6b4226");
    },
  },
  // The same chair, facing the other way.
  chairl: {
    w: 0.6,
    h: CHAIR_H,
    draw: (g, w, h) => {
      const seat = h * CHAIR_SEAT;
      box(g, 0, seat, w, h * 0.07, "#8b5a33"); // seat
      box(g, w * 0.82, h * 0.02, w * 0.18, seat, "#7a4a2b"); // back
      box(g, w * 0.7, h * 0.1, w * 0.12, seat * 0.7, "#a0663d");
      box(g, w * 0.81, seat + h * 0.07, w * 0.13, h * 0.43, "#6b4226");
      box(g, w * 0.07, seat + h * 0.07, w * 0.13, h * 0.43, "#6b4226");
    },
  },
  // A dining table, seen from the side.
  table: {
    w: 1.7,
    h: 0.78,
    draw: (g, w, h) => {
      box(g, 0, 0, w, h * 0.11, "#9c6a3f"); // top
      box(g, 0, h * 0.11, w, h * 0.05, "#7a4a2b");
      box(g, w * 0.1, h * 0.16, w * 0.09, h * 0.84, "#6b4226"); // legs
      box(g, w * 0.81, h * 0.16, w * 0.09, h * 0.84, "#6b4226");
    },
  },
  // An office desk with a screen on it.
  workdesk: {
    w: 1.9,
    h: 1.15,
    draw: (g, w, h) => {
      const top = h * 0.42; // desk surface
      box(g, w * 0.42, h * 0.04, w * 0.36, top * 0.62, "#2b3440"); // screen
      box(g, w * 0.45, h * 0.08, w * 0.3, top * 0.4, "#5b8fb0");
      box(g, w * 0.57, h * 0.3, w * 0.06, h * 0.08, "#8a94a3"); // stalk
      box(g, w * 0.52, h * 0.38, w * 0.16, h * 0.04, "#8a94a3");
      box(g, 0, top, w, h * 0.09, "#b9c2c0"); // top
      box(g, w * 0.03, top + h * 0.09, w * 0.28, h * 0.49, "#93a09c"); // pedestal
      box(g, w * 0.06, top + h * 0.2, w * 0.22, h * 0.03, "#7f8b88");
      box(g, w * 0.06, top + h * 0.36, w * 0.22, h * 0.03, "#7f8b88");
      box(g, w * 0.9, top + h * 0.09, w * 0.07, h * 0.49, "#7f8b88"); // leg
    },
  },
  // A shop checkout counter, seen from the side.
  // Canvas y grows downward: register up top.
  counter: {
    w: 2.6,
    h: 1.05,
    draw: (g, w, h) => {
      box(g, w * 0.6, h * 0.06, w * 0.3, h * 0.26, "#2b3440"); // register
      box(g, w * 0.64, h * 0.1, w * 0.22, h * 0.13, "#5b8fb0"); // screen
      box(g, 0, h * 0.32, w, h * 0.09, "#8b5a33"); // top slab
      box(g, 0, h * 0.41, w, h * 0.45, "#6b4226"); // front panel
      box(g, w * 0.04, h * 0.86, w * 0.92, h * 0.12, "#7a4a2b"); // kick
    },
  },
  // A fast food service counter, seen from
  // the side. Canvas y grows downward.
  fastcounter: {
    w: 4.2,
    h: 1.05,
    draw: (g, w, h) => {
      box(g, w * 0.86, h * 0.06, w * 0.11, h * 0.26, "#2b3440"); // till
      box(g, w * 0.88, h * 0.1, w * 0.07, h * 0.13, "#5b8fb0"); // screen
      box(g, 0, h * 0.32, w, h * 0.09, "#f2e3c8"); // top slab
      box(g, 0, h * 0.41, w, h * 0.45, "#d94f34"); // front panel
      box(g, w * 0.04, h * 0.86, w * 0.92, h * 0.12, "#b03a24"); // kick
    },
  },
  // A staircase, seen from the side, climbing
  // left to right like the runs in game.
  stair: {
    w: 6,
    h: 1.7,
    draw: (g, w, h) => {
      const steps = 8;
      const rise = h * 0.85;
      for (let i = 0; i < steps; i++) {
        const sw = w / steps;
        const sh = (rise * (i + 1)) / steps;
        box(g, i * sw, h * 0.9 - sh, sw + 1, sh, i % 2 ? "#8b5a33" : "#9c6a3f");
      }
      box(g, 0, h * 0.88, w, h * 0.07, "#6b4226"); // stringer
      box(g, w * 0.02, h * 0.28, w * 0.05, h * 0.6, "#6b4226"); // post
      box(g, w * 0.9, h * 0.04, w * 0.05, h * 0.84, "#6b4226"); // post
      box(g, 0, h * 0.28, w * 0.94, h * 0.06, "#8b5a33"); // handrail
    },
  },
  // Three z's drifting up from a sleeper.
  // The frames differ by one step, so they rise.
  zzz0: zzzFrame(0),
  zzz1: zzzFrame(1),
  zzz2: zzzFrame(2),
};

// Frames a sleeper's z's cycle through.
export const ZZZ_FRAMES = ["zzz0", "zzz1", "zzz2"] as const;

// One frame of the drifting z's, phase 0..2.
function zzzFrame(phase: number): PropSpec {
  return {
    w: 1.6,
    h: 1.4,
    draw: (g, w, h) => {
      g.fillStyle = "#f2f6ff";
      g.textBaseline = "alphabetic";
      for (let i = 0; i < 3; i++) {
        // Each z sits further along the same path.
        const t = (i + phase / 3) / 3;
        g.font = `bold ${Math.round(h * (0.3 + 0.26 * t))}px sans-serif`;
        g.fillText("z", w * (0.04 + 0.58 * t), h * (0.98 - 0.6 * t));
      }
    },
  };
}

// Bar ends and stools, one pair per room look.
function barProps(): Record<string, PropSpec> {
  const out: Record<string, PropSpec> = {};
  for (const [id, [front, top, seat]] of Object.entries(BAR_LOOK)) {
    out[`barend_${id}`] = {
      w: 0.9,
      h: 1,
      draw: (g, w, h) => {
        box(g, 0, 0, w, h * 0.1, top);
        box(g, w * 0.06, h * 0.1, w * 0.88, h * 0.9, front);
      },
    };
    out[`stool_${id}`] = {
      w: 0.4,
      h: 0.475,
      draw: (g, w, h) => {
        box(g, 0, 0, w, h * 0.15, seat);
        box(g, w * 0.4, h * 0.15, w * 0.2, h * 0.85, front);
      },
    };
  }
  return out;
}
