// Every sim in the tower, sharing one loaded model.

import { useMemo } from "react";
import { useGLTF } from "@react-three/drei";
import * as THREE from "three";

import { useSimIds, simTargets } from "./simsFeed";
import { useRooms } from "./seating";
import { getTypes } from "./catalog";
import { useSelectedRoom, useHover } from "./store";
import { BAR_SERVER_Z, CELL_H, CELL_W, FLOOR_T, SIM_H } from "./dims";
import { barSpot, counterSpot } from "./Furniture";
import { OUTFITS, UNIFORM, SMOKING, VIP, dress, makeLooks } from "./simOutfit";
import type { LoadedGLTF } from "./simRig";
import SimModel from "./SimModel";
import Cashier from "./Cashier";
import { useSlimModel } from "./lowPoly";
import modelUrl from "./models/Xbot.glb";

useGLTF.preload(modelUrl);

// Who minds which counter, and their prop.
const STAFF: Record<string, { kind: string; face: number }> = {
  retail: { kind: "counter", face: 0 },
  food: { kind: "fastcounter", face: 0 },
};

export default function Sims({ gridWidth }: { gridWidth: number }) {
  const gltf = useGLTF(modelUrl) as unknown as LoadedGLTF;
  useSlimModel(gltf);
  const ids = useSimIds();
  const rooms = useRooms();
  const hover = useHover();
  const inspecting = useSelectedRoom() === "inspect";

  // Scale the model to ~2 units and rest its feet on the floor.
  const { scale, footY } = useMemo(() => {
    const box = new THREE.Box3().setFromObject(gltf.scene);
    const h = box.max.y - box.min.y || 1;
    const s = SIM_H / h;
    return { scale: s, footY: -box.min.y * s };
  }, [gltf.scene]);

  // One painted copy of each mesh per outfit,
  // plus the officers' uniform at the end.
  const wardrobe = useMemo(() => {
    const meshes: THREE.SkinnedMesh[] = [];
    gltf.scene.traverse((o) => {
      const m = o as THREE.SkinnedMesh;
      if (m.isSkinnedMesh) meshes.push(m);
    });
    const dressed = (outfit: (typeof OUTFITS)[number]) => {
      const byName = new Map<string, THREE.BufferGeometry>();
      for (const m of meshes) byName.set(m.name, dress(m, outfit));
      return byName;
    };
    return [
      ...OUTFITS.map(dressed),
      dressed(UNIFORM),
      dressed(SMOKING),
      dressed(VIP),
    ];
  }, [gltf.scene]);

  // Officers, triads, and VIPs dress apart;
  // everyone else mixes freely.
  const outfitFor = (id: number, profession: string | undefined) => {
    if (profession === "security") return wardrobe.length - 3;
    if (profession === "triad") return wardrobe.length - 2;
    if (profession === "vip") return wardrobe.length - 1;
    return id % OUTFITS.length;
  };

  // The four materials every sim shares.
  const looks = useMemo(makeLooks, []);

  return (
    <>
      {ids.map((id) => (
        <SimModel
          key={id}
          id={id}
          gltf={gltf}
          scale={scale}
          footY={footY}
          gridWidth={gridWidth}
          inspecting={inspecting}
          hot={hover?.kind === "sim" && hover.id === id}
          maid={simTargets.get(id)?.profession === "maid"}
          outfit={wardrobe[outfitFor(id, simTargets.get(id)?.profession)]}
          looks={looks}
        />
      ))}
      {rooms
        .filter((r) => STAFF[r.category] && r.open)
        .map((r) => {
          const bar = getTypes().find((t) => t.id === r.type)?.line;
          const job = STAFF[r.category];
          const spot = bar
            ? { x: barSpot(r.width * CELL_W).serverX, z: BAR_SERVER_Z }
            : counterSpot(job.kind, r.width * CELL_W);
          return (
            <Cashier
              key={`${r.floor}:${r.col}`}
              x={(r.col + r.width / 2 - gridWidth / 2) * CELL_W + spot.x}
              y={r.floor * CELL_H + FLOOR_T}
              // Half a unit behind the counter line.
              z={bar ? spot.z : spot.z - 0.35}
              face={bar ? -Math.PI / 2 : job.face}
              gltf={gltf}
              scale={scale}
              footY={footY}
              outfit={wardrobe[(r.floor * 31 + r.col) % OUTFITS.length]}
              looks={looks}
            />
          );
        })}
    </>
  );
}
