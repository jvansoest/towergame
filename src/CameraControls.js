import React, { useRef } from "react";

/*
  THREE
*/
import { useFrame, extend, useThree } from "react-three-fiber";
/*
  MODELS
*/

/*
  ORBIT CONTROLS
*/
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
extend({ OrbitControls });

const CameraControls = () => {
  // Get a reference to the Three.js Camera, and the canvas html element.
  const {
    //mouse,
    camera,
    gl: { domElement },
  } = useThree();

  // Ref to the controls, so that we can update them on every frame using useFrame
  const controls = useRef();

  useFrame((state) => {
    if (camera.position.x > 30) {
      controls.current.target.x = 30;
      camera.position.set(30, controls.current.target.y, camera.position.z);
    }
    if (camera.position.x < -30) {
      controls.current.target.x = -30;
      camera.position.set(-30, controls.current.target.y, camera.position.z);
    }
    if (camera.position.y > 30) {
      controls.current.target.y = 30;
      camera.position.set(controls.current.target.x, 30, camera.position.z);
    }
    if (camera.position.y < -30) {
      controls.current.target.y = -30;
      camera.position.set(controls.current.target.x, -30, camera.position.z);
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
