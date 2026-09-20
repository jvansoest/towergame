// The earth behind the underground floors, and the
// B1, B2, ... labels down the left side.

import { useMemo } from "react";
import * as THREE from "three";

import { CELL_H, CELL_W } from "./dims";
import { builtAt } from "./gridCells";
import type { GridView } from "./protocol";

const EARTH_W = 2000; // wide enough to fill a zoomed-out view
const EARTH_D = 90; // world units of earth below the street
const TILE_W = 16; // world width of one repeat
const EARTH_Z = -1.62; // behind the rooms' back walls
const TOPSOIL_H = 0.5;

// Speckled earth: light red-brown on top, dark below.
function earthTexture(): THREE.CanvasTexture {
  const cv = document.createElement("canvas");
  cv.width = 64;
  cv.height = 256;
  const g = cv.getContext("2d");
  if (g) {
    for (let y = 0; y < cv.height; y++) {
      const t = Math.min(1, y / 200);
      for (let x = 0; x < cv.width; x++) {
        // A little noise keeps it from looking flat.
        const n = (Math.sin(x * 12.9898 + y * 78.233) * 43758.5453) % 1;
        const j = (Math.abs(n) - 0.5) * 26;
        const r = 138 - t * 70 + j;
        const gr = 92 - t * 46 + j * 0.7;
        const b = 58 - t * 30 + j * 0.5;
        g.fillStyle = `rgb(${r | 0},${gr | 0},${b | 0})`;
        g.fillRect(x, y, 1, 1);
      }
    }
  }
  const tex = new THREE.CanvasTexture(cv);
  tex.wrapS = THREE.RepeatWrapping;
  tex.repeat.set(EARTH_W / TILE_W, 1);
  tex.magFilter = THREE.NearestFilter;
  tex.colorSpace = THREE.SRGBColorSpace;
  return tex;
}

// A small pixel label like "B2".
function labelTexture(text: string): THREE.CanvasTexture {
  const cv = document.createElement("canvas");
  cv.width = 48;
  cv.height = 32;
  const g = cv.getContext("2d");
  if (g) {
    g.fillStyle = "#23262d";
    g.fillRect(0, 0, cv.width, cv.height);
    g.fillStyle = "#e6e8ec";
    g.font = "bold 22px monospace";
    g.textAlign = "center";
    g.textBaseline = "middle";
    g.fillText(text, cv.width / 2, cv.height / 2 + 1);
  }
  const tex = new THREE.CanvasTexture(cv);
  tex.magFilter = THREE.NearestFilter;
  tex.colorSpace = THREE.SRGBColorSpace;
  return tex;
}

export default function Basement({ grid }: { grid: GridView | null }) {
  const earth = useMemo(earthTexture, []);
  const rows = grid?.basement ?? 0;

  // The leftmost dug cell of each floor, for its label.
  const labels = useMemo(() => {
    if (!grid) return [];
    const out: { floor: number; x: number }[] = [];
    for (let d = 1; d <= rows; d++) {
      const floor = -d;
      let left = -1;
      for (let c = 0; c < grid.width; c++) {
        if (builtAt(grid, floor, c)) {
          left = c;
          break;
        }
      }
      if (left >= 0) {
        out.push({ floor, x: (left - grid.width / 2) * CELL_W - 0.9 });
      }
    }
    return out;
  }, [grid, rows]);

  return (
    <group>
      <mesh position={[0, -EARTH_D / 2, EARTH_Z]}>
        <planeGeometry args={[EARTH_W, EARTH_D]} />
        <meshBasicMaterial map={earth} toneMapped={false} />
      </mesh>
      {/* the topsoil line at street level */}
      <mesh position={[0, -TOPSOIL_H / 2, EARTH_Z + 0.01]}>
        <planeGeometry args={[EARTH_W, TOPSOIL_H]} />
        <meshBasicMaterial color="#4a3626" toneMapped={false} />
      </mesh>
      {labels.map((l) => (
        <FloorLabel key={l.floor} floor={l.floor} x={l.x} />
      ))}
    </group>
  );
}

function FloorLabel({ floor, x }: { floor: number; x: number }) {
  const tex = useMemo(() => labelTexture(`B${-floor}`), [floor]);
  return (
    <mesh position={[x, (floor + 0.5) * CELL_H, -1.5]}>
      <planeGeometry args={[1.2, 0.8]} />
      <meshBasicMaterial map={tex} toneMapped={false} />
    </mesh>
  );
}
