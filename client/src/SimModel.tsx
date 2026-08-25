// One sim: it follows the server's position, blends
// walk and idle, and answers the inspector.

import { useEffect, useMemo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import * as THREE from "three";
import * as SkeletonUtils from "three/examples/jsm/utils/SkeletonUtils.js";

import { askInspect } from "./inspect";
import { simTargets } from "./simsFeed";
import { setHover, clearHover } from "./store";
import {
  CELL_H,
  FLOOR_T,
  HIP_STAND,
  SEAT_H,
  SIM_H,
  SIM_Z,
  SIT_FACING,
  colToX,
} from "./dims";
import { findLegs, poseLegs, type LoadedGLTF } from "./simRig";
import { wear, type Look, type Wardrobe } from "./simOutfit";

type SimModelProps = {
  id: number;
  gltf: LoadedGLTF;
  scale: number;
  footY: number;
  gridWidth: number;
  inspecting: boolean;
  hot: boolean;
  outfit: Map<string, THREE.BufferGeometry>;
  looks: Wardrobe;
};

export default function SimModel({
  id,
  gltf,
  scale,
  footY,
  gridWidth,
  inspecting,
  hot,
  outfit,
  looks,
}: SimModelProps) {
  const group = useRef<THREE.Group>(null);
  const model = useRef<THREE.Object3D>(null);

  // Clones share geometry, so swap in the dressed copy.
  const clone = useMemo(() => {
    const c = SkeletonUtils.clone(gltf.scene);
    c.traverse((o) => {
      const m = o as THREE.SkinnedMesh;
      if (!m.isSkinnedMesh) return;
      const geo = outfit.get(m.name);
      // Without painted vertices the cloth material renders black.
      if (!geo) return;
      m.geometry = geo;
      m.material = looks.day;
    });
    return c;
  }, [gltf.scene, outfit, looks]);
  const look = useRef<Look>("day");
  const mixer = useMemo(() => new THREE.AnimationMixer(clone), [clone]);
  const cur = useRef<{ x: number; y: number } | null>(null);
  const last = useRef<{ x: number; y: number } | null>(null);
  const lastTarget = useRef<{ x: number; y: number } | null>(null);
  const lastX = useRef(0);
  const walkAction = useRef<THREE.AnimationAction | null>(null);
  const idleAction = useRef<THREE.AnimationAction | null>(null);
  const walkW = useRef(0);
  const legs = useMemo(() => findLegs(clone), [clone]);

  useEffect(() => {
    const walkClip =
      gltf.animations.find((c) => /walk/i.test(c.name)) ??
      gltf.animations.find((c) => /run/i.test(c.name)) ??
      gltf.animations[0];
    const idleClip = gltf.animations.find((c) => /idle/i.test(c.name));
    walkAction.current = walkClip ? mixer.clipAction(walkClip) : null;
    idleAction.current = idleClip ? mixer.clipAction(idleClip) : null;
    if (walkAction.current) walkAction.current.timeScale = 1.4; // brisk walk
    walkAction.current?.play();
    idleAction.current?.play();
    return () => {
      mixer.stopAllAction();
    };
  }, [mixer, gltf.animations]);

  useFrame((_, dt) => {
    mixer.update(dt);
    const t = simTargets.get(id);
    const g = group.current;
    if (!t || !g) return;

    // Dark rooms dim the sim; hover lights it up.
    const want: Look = t.dark
      ? hot
        ? "nightHot"
        : "night"
      : hot
        ? "dayHot"
        : "day";
    if (want !== look.current) {
      look.current = want;
      wear(clone, outfit, looks[want]);
    }
    if (!cur.current) cur.current = { ...t };

    // Snap on a server teleport (boarding, alighting, respawn).
    const lt = lastTarget.current;
    if (lt && Math.hypot(t.x - lt.x, t.y - lt.y) > 1) {
      cur.current.x = t.x;
      cur.current.y = t.y;
    }
    lastTarget.current = { x: t.x, y: t.y };

    if (t.riding) {
      // Track the car exactly, or the rider lags behind it.
      cur.current.x = t.x;
      cur.current.y = t.y;
    } else {
      const k = 1 - Math.pow(0.0001, dt); // exponential smoothing
      cur.current.x += (t.x - cur.current.x) * k;
      cur.current.y += (t.y - cur.current.y) * k;
    }

    // Blend walk/idle from movement speed.
    if (!last.current) last.current = { ...cur.current };
    const dist = Math.hypot(
      cur.current.x - last.current.x,
      cur.current.y - last.current.y
    );
    last.current = { ...cur.current };
    // Riders and sitters stand still.
    const moving =
      !t.riding && !t.sitting && dist / Math.max(dt, 1e-4) > 0.4;
    if (idleAction.current && walkAction.current) {
      walkW.current += ((moving ? 1 : 0) - walkW.current) * Math.min(1, dt * 12);
      walkAction.current.weight = walkW.current;
      idleAction.current.weight = 1 - walkW.current;
    } else if (walkAction.current) {
      walkAction.current.paused = !moving;
    }

    // Seated sims drop so their hips meet the seat.
    // Measured from the model's height, not its origin.
    const sitDrop = t.sitting ? SIM_H * (HIP_STAND - SEAT_H) : 0;
    // Feet rest on the slab, not on the floor line.
    g.position.set(
      colToX(cur.current.x, gridWidth),
      cur.current.y * CELL_H + FLOOR_T + footY - sitDrop,
      SIM_Z
    );
    poseLegs(legs, t.sitting);
    const dx = t.x - lastX.current;
    // Turn the model, not the group: the box stays square.
    const m = model.current;
    if (m && t.sitting) {
      // Face the chair, not the way they walked in.
      m.rotation.y = SIT_FACING;
    } else if (m && Math.abs(dx) > 0.001) {
      m.rotation.y = dx > 0 ? Math.PI / 2 : -Math.PI / 2;
    }
    lastX.current = t.x;
  });

  // The pick box is cheaper to hit than the rig.
  const midY = 0.95 - footY;
  return (
    <group ref={group}>
      <primitive ref={model} object={clone} scale={scale} />
      {inspecting && (
        <mesh
          position={[0, midY, 0]}
          // Stop the move, or the plane behind wins the hover.
          onPointerMove={(e) => {
            e.stopPropagation();
            setHover({ kind: "sim", id });
          }}
          onPointerOut={() => clearHover("sim", id)}
          onPointerDown={(e) => {
            if (e.button !== 0) return;
            e.stopPropagation();
            askInspect("inspectsim", { id });
          }}
        >
          <boxGeometry args={[0.9, 2, 0.9]} />
          <meshBasicMaterial transparent opacity={0} depthWrite={false} />
        </mesh>
      )}
    </group>
  );
}
