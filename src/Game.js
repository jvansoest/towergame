import React, { Suspense, useState, useEffect } from "react";

import socket from "./socket.js";
/*
UTILITIES
*/
import makeMatrix, { setMatrixValueAtIndex } from "./Utilities.js";

/*
  THREE
*/
import { extend, useThree } from "react-three-fiber";
import * as THREE from "three";
import BackGround from "./BackGround.js";

/*
  MODELS
*/
import Box from "./Box.js";
import CameraControls from "./CameraControls.js";
import GroundFloor from "./GroundFloor.js";

import XBot from "./XBot.js";

/*
  ORBIT CONTROLS
*/
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
extend({ OrbitControls });

const Game = () => {
  const raycaster = new THREE.Raycaster();
  const [matrix, _setMatrix] = useState(makeMatrix(20, 20));
  const [t, setT] = useState(null);
  const { camera, mouse, scene } = useThree();

  const setMatrix = (matrix, x, y) => {
    _setMatrix(setMatrixValueAtIndex(matrix, x, y, 1));
  };

  const filterArrayToBoxes = (matrix) => {
    const boxes = [];
    for (let row = 0; row < matrix.length; row++) {
      const array = matrix[row];
      for (let col = 0; col < array.length; col++) {
        const element = array[col];
        if (element === 1) {
          boxes.push(
            <Box
              key={row + "" + col}
              position={[(row - 10) * 10, (col - 10) * 5, -10]}
            />
          );
        }
      }
    }
    return boxes;
  };

  const getPlaneXYfromMouse = () => {
    raycaster.setFromCamera(mouse, camera);
    var intersects = raycaster.intersectObjects(scene.children);
    if (intersects.length > 0) {
      var coordinates = intersects[0].point;
      return coordinates;
    } else {
      return new THREE.Vector3().set(1, 0, 0);
    }
  };

  const makeGrid = (vec) => {
    const posX = Math.round(vec.x / 10) + 10;
    const posY = Math.round(vec.y / 5) + 10;
    return [posX, posY];
  };

  const handleClick = (e) => {
    const coordinates = getPlaneXYfromMouse();
    const [posX, posY] = makeGrid(coordinates);
    setMatrix(matrix, posX, posY);
    sendCoords(posX, posY);
  };

  const sendCoords = (x, y) => {
    socket.emit("boxplaced", { x: x, y: y });
  };

  const setter = (data) => {
    setMatrix(matrix, data.x, data.y);
    console.log(data);
  };

  useEffect(() => {
    socket.on("boxupdate", setter);
    return () => {
      socket.off("boxupdate", setter);
    };
  });

  return (
    <>
      <CameraControls />
      <ambientLight />
      <pointLight position={[100, 100, 100]} />
      <BackGround
        // matrix={matrix}
        makeGrid={makeGrid}
        getPlaneXYfromMouse={getPlaneXYfromMouse}
        onClick={handleClick}
        position={[0, 0, -10]}
      />
      <GroundFloor position={[0, 2, -10]} />

      <Suspense fallback={<Box position={[10, 10, 0]} />}>
        {<XBot position={[0, 0, 65]} />}
      </Suspense>

      {filterArrayToBoxes(matrix)}
    </>
  );
};
export default Game;
