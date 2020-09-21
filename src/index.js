import ReactDOM from "react-dom";
import React, { useRef, useState, Suspense } from "react";
import { Canvas, useFrame, useLoader } from "react-three-fiber";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
import model from "./roombox2.glb";

function Duck() {
  const gltf = useLoader(GLTFLoader, model);
  return <primitive object={gltf.scene} position={[-2, 0, 18]} />;
}

function Box(props) {
  const mesh = useRef();
  const [hovered, setHover] = useState(false);
  const [active, setActive] = useState(false);
  useFrame(() => (mesh.current.rotation.x = mesh.current.rotation.y += 0.01));

  return (
    <mesh
      {...props}
      ref={mesh}
      onClick={(e) => setActive(!active)}
      onPointerOver={(e) => setHover(true)}
      onPointerOut={(e) => setHover(false)}
    >
      <boxBufferGeometry attach="geometry" args={[5, 5, 5]} />
      <meshStandardMaterial
        attach="material"
        color={hovered ? "hotpink" : "orange"}
      />
    </mesh>
  );
}

function App(props) {
  return (
    <Canvas
      style={{ height: 400, width: 800 }}
      camera={{ fov: 75, position: [0, 0, 20], aspect: 30 }}
    >
      <Suspense fallback={<Box position={[10, 10, 0]} />}>{<Duck />}</Suspense>
      <ambientLight />
      <pointLight position={[100, 100, 100]} />
      <Box position={[10, 10, 0]} />
      <Box position={[0, 0, 0]} />
      <Box position={[0, 10, 0]} />
      <Box position={[10, 0, 0]} />
    </Canvas>
  );
}

ReactDOM.render(<App />, document.getElementById("root"));

// import ReactDOM from "react-dom";
// import React, { Suspense, useState } from "react";
// import { Canvas, useLoader } from "react-three-fiber";
// import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";
// import duck from "./roombox2.glb";

// function Duck() {
//   const gltf = useLoader(GLTFLoader, duck);
//   return <primitive object={gltf.scene} position={[0, 0, 0]} />;
// }

// function Box() {
//   return (
//     <mesh>
//       <boxBufferGeometry attach="geometry" args={[1, 1, 1]} />
//       <meshStandardMaterial attach="material" transparent opacity={0.5} />
//     </mesh>
//   );
// }

// function App() {
//   const [clicked, set] = useState(false);
//   return (
//     <>
//       <Canvas camera={{ position: [0, 0, 5] }}>
//         <ambientLight intensity={0.5} />
//         <spotLight intensity={0.8} position={[300, 300, 400]} />
//         <Suspense fallback={<Box />}>{clicked && <Duck />}</Suspense>
//       </Canvas>
//       {!clicked && (
//         <button onClick={() => set(true)}>Load duck w/ 1s delay</button>
//       )}
//     </>
//   );
// }

// ReactDOM.render(<App />, document.getElementById("root"));
