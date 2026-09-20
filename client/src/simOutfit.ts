// What a sim wears, and how light and hover tint it.

import * as THREE from "three";

import { DAY, HIGHLIGHT, NIGHT } from "./dims";

// Hover shows as an emissive glow on the sim itself.
const COLD = "#000000";

// One look per light and hover state.
export type Look = "day" | "night" | "dayHot" | "nightHot";
export type Wardrobe = Record<Look, THREE.Material>;

// Vertex colours carry the outfit, so four materials serve all.
export function makeLooks(): Wardrobe {
  const make = (color: string, emissive: string) =>
    new THREE.MeshStandardMaterial({
      vertexColors: true,
      roughness: 0.85,
      metalness: 0,
      color,
      emissive,
      emissiveIntensity: 0.6,
    });
  return {
    day: make(DAY, COLD),
    night: make(NIGHT, COLD),
    dayHot: make(DAY, HIGHLIGHT),
    nightHot: make(NIGHT, HIGHLIGHT),
  };
}

// Clothes, painted onto the mesh per vertex.
// The rig has no separate parts, so the bone
// a vertex follows decides what it wears.
type Part = "shirt" | "trousers" | "skin" | "shoes";

type Outfit = Record<Part, string>;

export const OUTFITS: Outfit[] = [
  { shirt: "#3b6ea5", trousers: "#33384a", skin: "#c68642", shoes: "#26262b" },
  { shirt: "#b5432f", trousers: "#2f3a4a", skin: "#e0ac69", shoes: "#3a2a22" },
  { shirt: "#e4e6eb", trousers: "#4a4a52", skin: "#8d5524", shoes: "#26262b" },
  { shirt: "#4f8a5b", trousers: "#6b5a45", skin: "#f1c27d", shoes: "#40342a" },
  { shirt: "#7a4e9c", trousers: "#2b2b33", skin: "#c68642", shoes: "#26262b" },
  { shirt: "#d8a13a", trousers: "#3f4756", skin: "#ffdbac", shoes: "#3a2a22" },
  { shirt: "#2f7d8c", trousers: "#5a4636", skin: "#8d5524", shoes: "#26262b" },
  { shirt: "#c25b8a", trousers: "#333a3f", skin: "#e0ac69", shoes: "#40342a" },
];

// Security officers all wear this one kit.
// Kept out of the crowd's random cycle.
export const UNIFORM: Outfit = {
  shirt: "#2b3a55",
  trousers: "#1f2436",
  skin: "#c68642",
  shoes: "#15171f",
};

// Triad members wear a black smoking, so a
// sharp eye can pick them from the crowd.
export const SMOKING: Outfit = {
  shirt: "#17171c",
  trousers: "#101014",
  skin: "#c68642",
  shoes: "#0c0c10",
};

// Which garment covers the vertices of a bone.
function partOfBone(name: string): Part {
  const n = name.toLowerCase();
  if (n.includes("head") || n.includes("neck") || n.includes("hand")) {
    return "skin";
  }
  if (n.includes("foot") || n.includes("toe")) return "shoes";
  if (n.includes("leg") || n.includes("hips")) return "trousers";
  return "shirt";
}

// Copies a mesh, colouring each vertex by its bone.
export function dress(
  mesh: THREE.SkinnedMesh,
  outfit: Outfit,
): THREE.BufferGeometry {
  const geo = mesh.geometry.clone();
  const index = geo.getAttribute("skinIndex");
  const weight = geo.getAttribute("skinWeight");
  const count = geo.getAttribute("position").count;
  const colors = new Float32Array(count * 3);

  // One colour per bone, so the lookup runs once.
  const bones = mesh.skeleton.bones;
  const tint = bones.map((b) => new THREE.Color(outfit[partOfBone(b.name)]));

  for (let v = 0; v < count; v++) {
    let best = 0;
    let top = -1;
    for (let k = 0; k < 4; k++) {
      const w = weight.getComponent(v, k);
      if (w > top) {
        top = w;
        best = index.getComponent(v, k);
      }
    }
    (tint[best] ?? tint[0]).toArray(colors, v * 3);
  }
  geo.setAttribute("color", new THREE.BufferAttribute(colors, 3));
  return geo;
}

// Puts one material on every dressed mesh.
export function wear(
  root: THREE.Object3D,
  outfit: Map<string, THREE.BufferGeometry>,
  mat: THREE.Material,
) {
  root.traverse((o) => {
    const m = o as THREE.SkinnedMesh;
    if (m.isSkinnedMesh && outfit.has(m.name)) m.material = mat;
  });
}

// A VIP wears a cream suit with gold trim.
export const VIP: Outfit = {
  shirt: "#efe6c8",
  trousers: "#d8ccaa",
  skin: "#e0ac69",
  shoes: "#7a5a2a",
};
