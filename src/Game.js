import React, { useState } from "react";

/*
  THREE
*/
import { extend, useThree, useFrame } from "react-three-fiber";
import * as THREE from "three";
import BackGround from "./BackGround.js";
/*
  MODELS
*/

import Box from "./Box.js";
import CameraControls from "./CameraControls.js";
/*
  ORBIT CONTROLS
*/
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
extend({ OrbitControls });

const Game = () => {
  const makeMatrix = (m, n) => {
    return Array.from({ length: m }, () => new Array(n).fill(0));
  };
  const [matrix, _setMatrix] = useState(makeMatrix(800, 400));
  const [lastCoords, setLastCoords] = useState({});
  const { camera } = useThree();

  const unSetMatrix = (matrix, x, y) => {
    _setMatrix(
      matrix.map((arr, i) => {
        if (x === i) {
          return arr.map((value, j) => {
            if (y === j) {
              setLastCoords({ x: x, y: y });
              return 0;
            } else {
              return value;
            }
          });
        } else {
          return arr;
        }
      })
    );
  };

  const setMatrix = (matrix, x, y) => {
    _setMatrix(
      matrix.map((arr, i) => {
        if (x === i) {
          return arr.map((value, j) => {
            if (y === j) {
              setLastCoords({ x: x, y: y });
              return 1;
            } else {
              return value;
            }
          });
        } else {
          return arr;
        }
      })
    );
  };

  const filterArrayToBoxes = (matrix) => {
    const boxes = [];
    for (let row = 0; row < matrix.length; row++) {
      const array = matrix[row];
      for (let col = 0; col < array.length; col++) {
        const element = array[col];
        if (element === 1) {
          boxes.push(
            <Box key={row + "" + col} position={[row - 100, col - 50, -10]} />
          );
        }
      }
    }
    return boxes;
  };

  const getPlaneXYfromMouse = (mouseX, mouseY, camera) => {
    var vec = new THREE.Vector3(); // create once and reuse
    var pos = new THREE.Vector3(); // create once and reuse

    const x = (mouseX / 800) * 2 - 1;
    const y = -(mouseY / 400) * 2 + 1;

    const z = -10;

    vec.set(x, y, z);
    vec.unproject(camera);

    vec.sub(camera.position).normalize();
    //var distance = -camera.position.z / vec.z;

    const targetZ = -11;

    var distance = (targetZ - camera.position.z) / vec.z;
    pos.copy(camera.position).add(vec.multiplyScalar(distance));
    return vec;
  };

  const makeGrid = (vec) => {
    //make grid
    const posX = Math.round((vec.x + 100) / 10) * 10;
    const posY = Math.round((vec.y + 50) / 5) * 5;
    console.log(posX, posY);
    return [posX, posY];
  };

  const handleClick = (e) => {
    const vec = getPlaneXYfromMouse(e.clientX, e.clientY, camera);
    const [posX, posY] = makeGrid(vec);
    setMatrix(matrix, posX, posY);
  };

  return (
    <>
      <CameraControls />
      <ambientLight />
      <pointLight position={[100, 100, 100]} />
      <BackGround
        matrix={matrix}
        lastCoords={lastCoords}
        makeGrid={makeGrid}
        unSetMatrix={unSetMatrix}
        setMatrix={setMatrix}
        getPlaneXYfromMouse={getPlaneXYfromMouse}
        // onMouseMove={handlePointerOver}
        onClick={handleClick}
        position={[0, 0, -10]}
      />
      {filterArrayToBoxes(matrix)}
    </>
  );
};
export default Game;
