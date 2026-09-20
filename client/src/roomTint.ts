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
  service: "#e3ecef",
  medical: "#e6f2f0",
  transit: "#dfe3ea",
  parking: "#565c66",
};

// Floor color per category.
export const FLOOR: Record<string, string> = {
  lobby: "#c9ccd1",
};

// Concrete used for floors, ceilings, and side walls.
export const CONCRETE = "#6b7280";

// Bar colours per room type: front, top, stool.
export const BAR_LOOK: Record<string, [string, string, string]> = {
  izakaya: ["#4a2a1c", "#9c6a3f", "#7a4a2b"],
  icecream: ["#f4a3bd", "#fbf6ee", "#e0709a"],
  burger: ["#d94f34", "#f2c230", "#3c3f46"],
  pizza: ["#2a6a3f", "#fbf6ee", "#c0392b"],
  noodle: ["#2f5d62", "#c79a5a", "#c0392b"],
};

// The look for a type; plain wood as fallback.
export function barLook(type: string): [string, string, string] {
  return BAR_LOOK[type] ?? BAR_LOOK.izakaya;
}
