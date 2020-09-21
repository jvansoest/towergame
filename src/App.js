import React from "react";

/*
  THREE
*/
import { Canvas, extend } from "react-three-fiber";
import * as THREE from "three";
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

function matrix(m, n) {
  return Array.from({ length: m }, () => new Array(n).fill(0));
}

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      matrix: matrix(800, 400),
    };
  }

  handleClick = (e) => {
    // setActive(!active);
    const clickedX = (e.clientX / 800) * 2 - 1;
    const clickedY = (e.clientY / window.innerHeight) * 2 + 1;
    var vec = new THREE.Vector3(); // create once and reuse
    var pos = new THREE.Vector3(); // create once and reuse

    vec.set((e.clientX / 800) * 2 - 1, -(e.clientY / 400) * 2 + 1, 0.5);

    console.log(vec);
    this.setState((prevState, props) => ({
      matrix: prevState.matrix.map((arr, i) => {
        if (clickedX === i) {
          return arr.map((value, j) => {
            if (clickedY === j) {
              return 1;
            } else {
              return value;
            }
          });
        } else {
          return arr;
        }
      }),
    }));
  };

  render = () => {
    return (
      <Canvas
        onClick={this.handleClick}
        style={{ height: 400, width: 800 }}
        camera={{ fov: 50, position: [0, 0, 110], aspect: 70 }}
      >
        <CameraControls />
        <ambientLight />
        <pointLight position={[100, 100, 100]} />
        {/* <BackGround position={[-1, 0, -1]} /> */}
        {this.state.matrix.map((arr, i) => {
          return arr.map((el, j) => {
            return el === 1 ? (
              <Box key={i + "" + j} position={[i, j, 0]} />
            ) : null;
          });
        })}
      </Canvas>
    );
  };
}
export default App;
