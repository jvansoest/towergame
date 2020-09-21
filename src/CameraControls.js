import React, { useRef } from "react";

/*
  THREE
*/
import { Canvas, useFrame, extend, useThree } from "react-three-fiber";
import * as THREE from "three";
/*
  MODELS
*/

import Box from "./Box.js";

/*
  ORBIT CONTROLS
*/
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
extend({ OrbitControls });

const CameraControls = () => {
  // Get a reference to the Three.js Camera, and the canvas html element.
  const {
    camera,
    gl: { domElement },
  } = useThree();

  // Ref to the controls, so that we can update them on every frame using useFrame
  const controls = useRef();

  useFrame((state) => {
    if (camera.position.x > 10) {
      controls.current.target.x = 10;
      camera.position.set(10, controls.current.target.y, camera.position.z);
    }
    if (camera.position.x < -1) {
      controls.current.target.x = -1;
      camera.position.set(-1, controls.current.target.y, camera.position.z);
    }
    if (camera.position.y > 10) {
      controls.current.target.y = 10;
      camera.position.set(controls.current.target.x, 10, camera.position.z);
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
      keyPanSpeed={10}
      enableDamping={true}
      dampingFactor={0.15}
    />
  );
};
export default CameraControls;
