// Slims the character once, before any sim clones it:
// fewer triangles, fewer tracks, only the clips in use.

import { use } from "react";
import * as THREE from "three";
import { MeshoptSimplifier } from "meshoptimizer";

import type { LoadedGLTF } from "./simRig";

const KEEP_TRIS = 0.08; // share of triangles left
const MAX_ERROR = 0.05; // of the mesh size

// Rebuilds a skinned mesh with fewer triangles.
function simplify(mesh: THREE.SkinnedMesh) {
  const geo = mesh.geometry;
  const index = geo.getIndex();
  const pos = geo.getAttribute("position");
  const skin = geo.getAttribute("skinWeight");
  if (!index || !skin || !(pos.array instanceof Float32Array)) return;

  const [kept] = MeshoptSimplifier.simplifyWithAttributes(
    Uint32Array.from(index.array),
    pos.array,
    pos.itemSize,
    skin.array as Float32Array,
    skin.itemSize,
    [1, 1, 1, 1],
    null,
    Math.floor((index.count * KEEP_TRIS) / 3) * 3,
    MAX_ERROR,
  );

  // Renumber the vertices that survive.
  const remap = new Map<number, number>();
  const order: number[] = [];
  const next = new Uint32Array(kept.length);
  kept.forEach((v, i) => {
    let n = remap.get(v);
    if (n === undefined) {
      n = order.length;
      remap.set(v, n);
      order.push(v);
    }
    next[i] = n;
  });

  const slim = new THREE.BufferGeometry();
  for (const [name, attr] of Object.entries(geo.attributes)) {
    const a = attr as THREE.BufferAttribute;
    const Kind = a.array.constructor as new (n: number) => THREE.TypedArray;
    const out = new Kind(order.length * a.itemSize);
    order.forEach((v, i) => {
      for (let k = 0; k < a.itemSize; k++) {
        out[i * a.itemSize + k] = a.array[v * a.itemSize + k];
      }
    });
    slim.setAttribute(
      name,
      new THREE.BufferAttribute(out, a.itemSize, a.normalized),
    );
  }
  slim.setIndex(new THREE.BufferAttribute(next, 1));
  mesh.geometry = slim;
  geo.dispose();
}

// Keeps rotations, plus the hips' travel. Drops
// fingers, and the scale and position of the rest.
function slimClip(clip: THREE.AnimationClip): THREE.AnimationClip {
  const tracks = clip.tracks.filter((t) => {
    if (/(Thumb|Index|Middle|Ring|Pinky)/.test(t.name)) return false;
    if (t.name.endsWith(".quaternion")) return true;
    return t.name.endsWith(".position") && /Hips\./.test(t.name);
  });
  return new THREE.AnimationClip(clip.name, clip.duration, tracks);
}

const done = new WeakMap<object, Promise<void>>();

async function slimDown(gltf: LoadedGLTF) {
  await MeshoptSimplifier.ready;
  if (MeshoptSimplifier.supported) {
    gltf.scene.traverse((o) => {
      const m = o as THREE.SkinnedMesh;
      if (m.isSkinnedMesh) simplify(m);
    });
  }
  gltf.animations = gltf.animations
    .filter((c) => /idle|walk/i.test(c.name))
    .map(slimClip);
}

// Suspends until the model is slim; safe to call each render.
export function useSlimModel(gltf: LoadedGLTF) {
  let p = done.get(gltf.scene);
  if (!p) {
    p = slimDown(gltf);
    done.set(gltf.scene, p);
  }
  use(p);
}
