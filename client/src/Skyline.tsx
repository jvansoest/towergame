import { useEffect, useMemo } from "react";
import * as THREE from "three";
import { useThree } from "@react-three/fiber";

// Sky gradient and city skyline backdrop.

// Skyline plane placement.
const PLANE_WIDTH = 1000;
const PLANE_HEIGHT = 90;
const PLANE_Y = 44; // center; base near ground
const PLANE_Z = -60; // behind the build plane

function makeSkyTexture(): THREE.Texture {
  const canvas = document.createElement("canvas");
  canvas.width = 2;
  canvas.height = 512;
  const ctx = canvas.getContext("2d")!;
  const g = ctx.createLinearGradient(0, 0, 0, 512);
  g.addColorStop(0, "#3b7fc4"); // deep blue
  g.addColorStop(0.55, "#8fbfe6"); // mid sky
  g.addColorStop(1, "#dcecf6"); // horizon haze
  ctx.fillStyle = g;
  ctx.fillRect(0, 0, 2, 512);
  const tex = new THREE.CanvasTexture(canvas);
  tex.colorSpace = THREE.SRGBColorSpace;
  return tex;
}

function makeSkylineTexture(): THREE.Texture {
  const w = 2048;
  const h = 512;
  const canvas = document.createElement("canvas");
  canvas.width = w;
  canvas.height = h;
  const ctx = canvas.getContext("2d")!;
  ctx.clearRect(0, 0, w, h); // transparent above the buildings

  const drawLayer = (
    color: string,
    minH: number,
    maxH: number,
    minW: number,
    maxW: number,
    lit: boolean
  ) => {
    let x = 0;
    while (x < w) {
      const bw = minW + Math.random() * (maxW - minW);
      const bh = minH + Math.random() * (maxH - minH);
      const top = h - bh;

      ctx.fillStyle = color;
      ctx.fillRect(x, top, bw, bh);

      // rooftop antenna
      if (Math.random() < 0.35) {
        ctx.fillRect(x + bw * 0.35, top - 12, bw * 0.14, 12);
      }

      // lit windows on near layer
      if (lit) {
        ctx.fillStyle = "rgba(255,236,170,0.55)";
        for (let wy = top + 12; wy < h - 8; wy += 16) {
          for (let wx = x + 6; wx < x + bw - 6; wx += 14) {
            if (Math.random() < 0.45) ctx.fillRect(wx, wy, 6, 9);
          }
        }
      }

      x += bw + 2;
    }
  };

  // Far layer, then near layer.
  drawLayer("#a7bdd4", 70, 190, 40, 80, false);
  drawLayer("#6f8dad", 150, 330, 55, 115, true);

  const tex = new THREE.CanvasTexture(canvas);
  tex.colorSpace = THREE.SRGBColorSpace;
  return tex;
}

export default function Skyline() {
  const { scene } = useThree();
  const skyTex = useMemo(makeSkyTexture, []);
  const skylineTex = useMemo(makeSkylineTexture, []);

  // Set sky as scene background.
  useEffect(() => {
    const prev = scene.background;
    scene.background = skyTex;
    return () => {
      scene.background = prev;
      skyTex.dispose();
    };
  }, [scene, skyTex]);

  useEffect(() => () => skylineTex.dispose(), [skylineTex]);

  return (
    <mesh
      position={[0, PLANE_Y, PLANE_Z]}
      raycast={() => null} // no pointer picking
    >
      <planeGeometry args={[PLANE_WIDTH, PLANE_HEIGHT]} />
      <meshBasicMaterial map={skylineTex} transparent depthWrite={false} />
    </mesh>
  );
}
