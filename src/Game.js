import React, { useState } from "react";

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
  const { camera, scene } = useThree();

  const setMatrix = (matrix, x, y) => {
    _setMatrix(
      matrix.map((arr, i) => {
        if (x === i) {
          return arr.map((value, j) => {
            if (y === j) {
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

  const handleClick = (e) => {
    const mouse = {};

    mouse.x = (e.clientX / window.innerWidth) * 2 - 1;
    mouse.y = -(e.clientY / window.innerHeight) * 2 + 1;
    const raycaster = new THREE.Raycaster();

    var a = new THREE.Vector2(mouse.x, mouse.y);

    raycaster.setFromCamera(a, camera);

    var objects = raycaster.intersectObjects(scene.children, true);
    console.log(objects[0].point);

    var vec = new THREE.Vector3(); // create once and reuse
    var pos = new THREE.Vector3(); // create once and reuse

    vec.set((e.clientX / 800) * 2 - 1, -(e.clientY / 400) * 2 + 1, -10);
    vec.unproject(camera);
    vec.sub(camera.position).normalize();
    //var distance = -camera.position.z / vec.z;
    const targetZ = -11;
    var distance = (targetZ - camera.position.z) / vec.z;
    pos.copy(camera.position).add(vec.multiplyScalar(distance));

    const posX = Math.round((vec.x + 100) / 10) * 10;
    const posY = Math.round((vec.y + 50) / 5) * 5;

    // console.log(posX, posY);
    setMatrix(matrix, posX, posY);
  };

  return (
    <>
      <CameraControls />
      <ambientLight />
      <pointLight position={[100, 100, 100]} />
      <BackGround onClick={handleClick} position={[-1, -1, -10]} />
      {matrix.map((arr, i) => {
        return arr.map((el, j) => {
          return el === 1 ? (
            <Box key={i + "" + j} position={[i - 100, j - 50, -10]} />
          ) : null;
        });
      })}
    </>
  );
};
export default Game;
