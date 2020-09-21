import ReactDOM from "react-dom";
import React, { Suspense, useState } from "react";
import { Canvas, useLoader } from "react-three-fiber";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
import duck from "./duck.glb";
import "./styles.css";

function Duck() {
  const gltf = useLoader(GLTFLoader, duck);
  return <primitive object={gltf.scene} position={[0, 0, 0]} />;
}

function Box() {
  return (
    <mesh>
      <boxBufferGeometry attach="geometry" args={[1, 1, 1]} />
      <meshStandardMaterial attach="material" transparent opacity={0.5} />
    </mesh>
  );
}

function App() {
  const [clicked, set] = useState(false);
  return (
    <>
      <Canvas camera={{ position: [0, 0, 50] }}>
        <ambientLight intensity={0.5} />
        <spotLight intensity={0.8} position={[300, 300, 400]} />
        <Suspense fallback={<Box />}>{clicked && <Duck />}</Suspense>
      </Canvas>
      {!clicked && (
        <button onClick={() => set(true)}>Load duck w/ 1s delay</button>
      )}
    </>
  );
}

ReactDOM.render(<App />, document.getElementById("root"));
