// Lobby wall, one tile repeated across the strip.
// A merged strip can be any width, so the art tiles.

import * as THREE from "three";

import { box, paintTexture } from "./paint";

const LOBBY_TILE_W = 3; // world units per repeat

let lobbyTile: THREE.CanvasTexture | null | undefined;
const lobbyCache = new Map<number, THREE.CanvasTexture>();

// Repeating lobby wall, scaled to one strip.
export function lobbyArt(
  width: number,
  height: number
): THREE.CanvasTexture | null {
  const hit = lobbyCache.get(width);
  if (hit !== undefined) return hit;
  if (lobbyTile === undefined) lobbyTile = makeLobbyTile(height);
  if (!lobbyTile) return null;

  // Each strip needs its own repeat count.
  const tex = lobbyTile.clone();
  tex.needsUpdate = true;
  tex.repeat.set(Math.max(1, Math.round(width / LOBBY_TILE_W)), 1);
  lobbyCache.set(width, tex);
  return tex;
}

function makeLobbyTile(height: number): THREE.CanvasTexture | null {
  const tex = paintTexture(LOBBY_TILE_W, height, drawLobbyTile);
  if (tex) tex.wrapS = THREE.RepeatWrapping;
  return tex;
}

// Marble piers, a gold cornice, and panelling.
function drawLobbyTile(g: CanvasRenderingContext2D, w: number, h: number) {
  box(g, 0, 0, w, h, "#e9e2d2"); // wall
  box(g, 0, 0, w, h * 0.07, "#cdc4ae"); // ceiling
  box(g, 0, h * 0.07, w, h * 0.035, "#c9a227"); // cornice
  box(g, 0, h * 0.105, w, h * 0.012, "#8f7318");

  // Recessed panel between the piers.
  box(g, w * 0.34, h * 0.2, w * 0.5, h * 0.52, "#ded5c0");
  box(g, w * 0.36, h * 0.22, w * 0.46, h * 0.48, "#e6dfcd");
  box(g, w * 0.55, h * 0.3, w * 0.07, h * 0.1, "#f2d675"); // wall light

  // A marble pier at the tile edge.
  box(g, w * 0.06, h * 0.12, w * 0.15, h * 0.76, "#f4efe2");
  box(g, w * 0.19, h * 0.12, w * 0.03, h * 0.76, "#cfc6b1"); // shaded side
  box(g, w * 0.03, h * 0.12, w * 0.22, h * 0.05, "#f8f4e9"); // capital
  box(g, w * 0.03, h * 0.83, w * 0.22, h * 0.05, "#f8f4e9"); // base

  box(g, 0, h * 0.88, w, h * 0.12, "#7b6f57"); // skirting
}
