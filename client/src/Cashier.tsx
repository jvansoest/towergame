// The shop assistant behind a shop's counter.
// Pure scenery: it stands, faces the room,
// and appears only while the shop is open.

import { useEffect, useMemo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import * as THREE from "three";
import * as SkeletonUtils from "three/examples/jsm/utils/SkeletonUtils.js";

import { SIM_Z } from "./dims";
import type { LoadedGLTF } from "./simRig";
import type { Wardrobe } from "./simOutfit";

type CashierProps = {
  x: number; // world spot behind the counter
  y: number; // floor line the figure stands on
  z?: number; // depth lane; behind counters is less
  face?: number; // which way they look, radians
  gltf: LoadedGLTF;
  scale: number;
  footY: number;
  outfit: Map<string, THREE.BufferGeometry>;
  looks: Wardrobe;
};

export default function Cashier({
  x,
  y,
  z = SIM_Z,
  face = 0,
  gltf,
  scale,
  footY,
  outfit,
  looks,
}: CashierProps) {
  // One dressed copy in the day look.
  const clone = useMemo(() => {
    const c = SkeletonUtils.clone(gltf.scene);
    c.traverse((o) => {
      const m = o as THREE.SkinnedMesh;
      if (!m.isSkinnedMesh) return;
      const geo = outfit.get(m.name);
      if (!geo) return;
      m.geometry = geo;
      m.material = looks.day;
    });
    c.rotation.y = face;
    return c;
  }, [gltf.scene, outfit, looks, face]);

  const mixer = useMemo(() => new THREE.AnimationMixer(clone), [clone]);
  const idle = useMemo(
    () =>
      gltf.animations.find((c) => /idle/i.test(c.name)) ?? gltf.animations[0],
    [gltf.animations],
  );
  useEffect(() => {
    const action = idle ? mixer.clipAction(idle) : null;
    action?.play();
    return () => {
      mixer.stopAllAction();
    };
  }, [mixer, idle]);

  // Drive the clip at 10 Hz; idle needs no more.
  const owed = useRef(0);
  useFrame((_, dt) => {
    owed.current += dt;
    if (owed.current < 0.1) return;
    mixer.update(Math.min(owed.current, 0.2));
    owed.current = 0;
  });

  return (
    <group position={[x, y + footY, z]}>
      <primitive object={clone} scale={scale} />
    </group>
  );
}
