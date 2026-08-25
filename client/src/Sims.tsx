// Every sim in the tower, sharing one loaded model.

import { useMemo } from "react";
import { useGLTF } from "@react-three/drei";
import * as THREE from "three";

import { useSimIds } from "./simsFeed";
import { useSelectedRoom, useHover } from "./store";
import { SIM_H } from "./dims";
import { OUTFITS, dress, makeLooks } from "./simOutfit";
import type { LoadedGLTF } from "./simRig";
import SimModel from "./SimModel";
import modelUrl from "./models/Xbot.glb";

useGLTF.preload(modelUrl);

export default function Sims({ gridWidth }: { gridWidth: number }) {
  const gltf = useGLTF(modelUrl) as unknown as LoadedGLTF;
  const ids = useSimIds();
  const hover = useHover();
  const inspecting = useSelectedRoom() === "inspect";

  // Scale the model to ~2 units and rest its feet on the floor.
  const { scale, footY } = useMemo(() => {
    const box = new THREE.Box3().setFromObject(gltf.scene);
    const h = box.max.y - box.min.y || 1;
    const s = SIM_H / h;
    return { scale: s, footY: -box.min.y * s };
  }, [gltf.scene]);

  // One painted copy of each mesh per outfit.
  const wardrobe = useMemo(() => {
    const meshes: THREE.SkinnedMesh[] = [];
    gltf.scene.traverse((o) => {
      const m = o as THREE.SkinnedMesh;
      if (m.isSkinnedMesh) meshes.push(m);
    });
    return OUTFITS.map((outfit) => {
      const byName = new Map<string, THREE.BufferGeometry>();
      for (const m of meshes) byName.set(m.name, dress(m, outfit));
      return byName;
    });
  }, [gltf.scene]);

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
          outfit={wardrobe[id % wardrobe.length]}
          looks={looks}
        />
      ))}
    </>
  );
}
