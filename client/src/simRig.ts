// The character model: its bones, and the sit pose.
// The rig has no sit clip, so the legs are posed by hand.

import type * as THREE from "three";

// The model as loaded, with its clips.
export type LoadedGLTF = { scene: THREE.Group; animations: THREE.AnimationClip[] };

// Seated pose, in radians. The rig has no sit clip,
// so the legs are posed by hand after the mixer runs.
const SIT_THIGH = -Math.PI / 2; // hip swings forward
const SIT_KNEE = Math.PI / 2; // shin drops down
const SIT_LEAN = 0.12; // slight recline

// Leg bones of a Mixamo rig, if present.
export type Legs = {
  lUp: THREE.Object3D;
  rUp: THREE.Object3D;
  lLow: THREE.Object3D;
  rLow: THREE.Object3D;
  hips: THREE.Object3D;
};

// Matches bone names, ignoring prefix punctuation.
function findBone(root: THREE.Object3D, want: string): THREE.Object3D | null {
  let hit: THREE.Object3D | null = null;
  root.traverse((o) => {
    if (hit) return;
    if (o.name.replace(/[^a-zA-Z]/g, "").toLowerCase().endsWith(want)) hit = o;
  });
  return hit;
}

export function findLegs(root: THREE.Object3D): Legs | null {
  const lUp = findBone(root, "leftupleg");
  const rUp = findBone(root, "rightupleg");
  const lLow = findBone(root, "leftleg");
  const rLow = findBone(root, "rightleg");
  const hips = findBone(root, "hips");
  if (!lUp || !rUp || !lLow || !rLow || !hips) return null;
  return { lUp, rUp, lLow, rLow, hips };
}

// Bends the legs into a sit, or lets the clip play.
export function poseLegs(legs: Legs | null, sitting: boolean) {
  if (!legs) return;
  if (!sitting) {
    legs.hips.rotation.x = 0;
    return;
  }
  // Runs after mixer.update, so it wins for this frame.
  legs.lUp.rotation.set(SIT_THIGH, 0, 0);
  legs.rUp.rotation.set(SIT_THIGH, 0, 0);
  legs.lLow.rotation.set(SIT_KNEE, 0, 0);
  legs.rLow.rotation.set(SIT_KNEE, 0, 0);
  legs.hips.rotation.x = SIT_LEAN;
}
