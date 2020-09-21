import React, { useRef, Suspense } from "react";

/*
THREE
*/
import { Canvas, useFrame, extend, useThree } from "react-three-fiber";

/*
MODELS
*/
import Donut from "./duck.js";
import Box from "./Box.js";

/*
ORBIT CONTROLS
*/
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
extend({ OrbitControls });

function matrix(m, n) {
  return Array.from({ length: m }, () => new Array(n).fill(0));
}

const CameraControls = () => {
  // Get a reference to the Three.js Camera, and the canvas html element.
  const {
    camera,
    gl: { domElement },
  } = useThree();

  // Ref to the controls, so that we can update them on every frame using useFrame
  const controls = useRef();

  useFrame((state) => {
    /*
        I would like to only update the camera if the camera X or Y is not > 10 for example
    */
    // console.log(camera.position);
    if (camera.position.x > 1) {
      controls.current.target.x = 1;
      camera.position.set(1, controls.current.target.y, camera.position.z);
    }
    if (camera.position.x < -1) {
      controls.current.target.x = -1;
      camera.position.set(-1, controls.current.target.y, camera.position.z);
    }
    if (camera.position.y > 1) {
      controls.current.target.y = 1;
      camera.position.set(controls.current.target.x, 1, camera.position.z);
    }
    if (camera.position.y < -1) {
      controls.current.target.y = -1;
      camera.position.set(controls.current.target.x, -1, camera.position.z);
    }
    controls.current.update();
  });

  return (
    <orbitControls
      ref={controls}
      args={[camera, domElement]}
      enableZoom={false}
      enableRotate={false}
      keyPanSpeed={1}
      enableDamping={true}
      dampingFactor={0.15}
    />
  );
};

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = { Matrix: matrix(100, 100) };
  }
  render = () => {
    return (
      <Canvas
        style={{ height: 400, width: 800 }}
        camera={{ fov: 50, position: [0, 0, 22], aspect: 70 }}
      >
        <CameraControls />
        <Suspense fallback={<Box position={[10, 10, 0]} />}>
          {<Donut />}
        </Suspense>
        <ambientLight />
        <pointLight position={[100, 100, 100]} />
      </Canvas>
    );
  };
}

export default App;
