import React, { useRef, useState } from "react";
import * as THREE from "three";
import { useThree } from "react-three-fiber";
import Box from "./Box.js";

function Grid(props) {
  //const [box, setBox] = useState({ x: -1, y: -1 });
  const [showBox, setShowBox] = useState(false);
  const mesh = useRef();
  const gridHelper = new THREE.GridHelper(100, 20);
  const raycaster = new THREE.Raycaster();
  const { mouse, camera, scene } = useThree();

  var vec = new THREE.Vector3();
  vec.set(1, 0, 0);

  gridHelper.setRotationFromAxisAngle(vec, Math.PI / 2);

  const handleHover = (e) => {
    raycaster.setFromCamera(mouse, camera);
    var intersects = raycaster.intersectObjects(scene.children);
    console.log(intersects);
    if (intersects.length > 0) {
      //var intersect = intersects[0];
      //var coordinates = intersect[0].points;
      setShowBox(true);
    }
  };

  return (
    <>
      {showBox && <Box position={[1, 1, 1]} />}
      <primitive
        onPointerMove={(e) => handleHover(e)}
        ref={mesh}
        object={gridHelper}
        position={props.position}
      />
    </>
  );
}

export default Grid;
