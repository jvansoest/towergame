import React, { useRef } from "react";
import { useLoader } from "react-three-fiber";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
import model from "./models/roombox2.glb";

function Duck(props) {
  const mesh = useRef();
  //useFrame(() => (mesh.current.rotation.x = mesh.current.rotation.y += 0.01));
  const gltf = useLoader(GLTFLoader, model);
  return <primitive ref={mesh} object={gltf.scene} position={props.position} />;
}

//HOW TO USE: <Suspense fallback={<Box position={[10, 10, 0]} />}>{<Donut />}</Suspense>

export default Duck;
