// Room colours, and how the room's light changes them.

import { DAY } from "./dims";

// Multiplies a colour by the room's light.
export function shadeHex(hex: string, tint: string): string {
  if (tint === DAY) return hex;
  const a = parseInt(hex.slice(1), 16);
  const b = parseInt(tint.slice(1), 16);
  const ch = (s: number) =>
    Math.round((((a >> s) & 255) * ((b >> s) & 255)) / 255);
  const v = (ch(16) << 16) | (ch(8) << 8) | ch(0);
  return "#" + v.toString(16).padStart(6, "0");
}

// Back-wall tint per category.
export const WALL: Record<string, string> = {
  lobby: "#e8edf3",
  office: "#e2e8f0",
  residential: "#f1e9dc",
  hotel: "#b39ddb",
  food: "#f6e0dc",
  retail: "#efe6f6",
  entertainment: "#26263a",
};

// Floor color per category.
export const FLOOR: Record<string, string> = {
  lobby: "#c9ccd1",
};

// Concrete used for floors, ceilings, and side walls.
export const CONCRETE = "#6b7280";
