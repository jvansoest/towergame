import React, { useRef } from "react";
import { useLoader } from "react-three-fiber";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
import model from "./models/donut2.glb";

function Duck() {
  const mesh = useRef();
  //useFrame(() => (mesh.current.rotation.x = mesh.current.rotation.y += 0.01));
  const gltf = useLoader(GLTFLoader, model);
  return (
    <primitive ref={mesh} object={gltf.scene} position={[0, -0.1, 19.7]} />
  );
}

//HOW TO USE: <Suspense fallback={<Box position={[10, 10, 0]} />}>{<Donut />}</Suspense>

export default Duck;
