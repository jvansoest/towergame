// Flat, pixel-perfect room cards: the painted
// wall, its fittings and chairs, a counter, and
// a stair in front if asked. Everything stamps
// the shared sprites, so a card always matches
// the rooms in game. Facts come from the
// server's catalog, never a local copy.

import { getTypes } from "./catalog";
import { CELL_H, CELL_W } from "./dims";
import { barSpot, counterFor, counterSpot } from "./Furniture";
import { lobbyArt } from "./lobbyArt";
import { chairKind, seatCols, tableCols } from "./layout";
import { paintTexture, homeShell, windowBand, PPU } from "./paint";
import { propArt } from "./propArt";
import { roomArt } from "./roomArt";

// Local x of a cell centre, in pixels.
function xOf(col: number, cells: number): number {
  return ((col + 0.5 - cells / 2) * CELL_W * PPU) | 0;
}

// Stamps one sprite, feet on the card floor.
function stamp(
  g: CanvasRenderingContext2D,
  kind: string,
  x: number,
  floorY: number,
) {
  const p = propArt(kind);
  if (!p) return;
  const w = Math.round(p.w * PPU);
  const h = Math.round(p.h * PPU);
  g.drawImage(
    p.tex.image as CanvasImageSource,
    Math.round(x - w / 2),
    Math.round(floorY - h),
    w,
    h,
  );
}

// The wall art for a card, or a plain shell
// for the types nobody painted yet.
function backdrop(
  type: string,
  w: number,
  h: number,
  variant = 0,
): CanvasImageSource | null {
  const art = type === "lobby" ? lobbyArt(w, h) : roomArt(type, w, h, variant);
  if (art) return art.image as CanvasImageSource;
  return (
    paintTexture(w, h, (g, cw, ch) => {
      homeShell(g, cw, ch, "#dfe4e2", "#7b8a86");
      windowBand(g, cw, ch, "#c3cbc8");
    })?.image ?? null
  );
}

// Renders one room card. Returns the canvas,
// or null for types the composer cannot draw.
export function roomCardCanvas(
  type: string,
  opts?: { stair?: boolean; variant?: number },
): HTMLCanvasElement | null {
  const spec = getTypes().find((t) => t.id === type);
  if (!spec) return null;
  const w = spec.width * CELL_W;
  const desks = spec.category === "office";
  const cv = document.createElement("canvas");
  cv.width = Math.round(w * PPU);
  cv.height = Math.round(CELL_H * PPU);
  const g = cv.getContext("2d");
  if (!g) return null;
  g.imageSmoothingEnabled = false;

  const back = backdrop(type, w, CELL_H, opts?.variant);
  if (back) g.drawImage(back, 0, 0, cv.width, cv.height);

  // Tables behind, then chairs in front, laid
  // out exactly like the rooms draw them.
  // xOf is room-centred; the canvas is not.
  if (spec.line) {
    const bar = barSpot(w);
    stamp(g, `barend_${type}`, cv.width / 2 + bar.counterX * PPU, cv.height);
    stamp(g, `stool_${type}`, cv.width / 2 + bar.seatX * PPU, cv.height);
  } else if (spec.seats > 0) {
    const floorY = cv.height;
    for (const c of tableCols(spec.width, spec.seats, desks)) {
      stamp(
        g,
        desks ? "workdesk" : "table",
        cv.width / 2 + xOf(c, spec.width),
        floorY,
      );
    }
    seatCols(spec.width, spec.seats).forEach((c, i) => {
      stamp(g, chairKind(i, desks), cv.width / 2 + xOf(c, spec.width), floorY);
    });
  }

  const kind = counterFor(type);
  if (kind) {
    const spot = counterSpot(kind, w);
    stamp(g, kind, cv.width / 2 + spot.x * PPU, cv.height);
  }

  if (opts?.stair) {
    // Clear of the east wall, drawn over the lot.
    stamp(g, "stair", cv.width - Math.round(3.2 * PPU), cv.height);
  }
  return cv;
}

// Data URLs, made once and kept per card.
const urls = new Map<string, string>();

export function roomCardURL(
  type: string,
  opts?: { stair?: boolean; variant?: number },
): string | null {
  const key = `${type}#${opts?.variant ?? 0}${opts?.stair ? "+stair" : ""}`;
  const hit = urls.get(key);
  if (hit) return hit;
  const cv = roomCardCanvas(type, opts);
  if (!cv) return null;
  const url = cv.toDataURL("image/png");
  urls.set(key, url);
  return url;
}
