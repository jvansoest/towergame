import React from "react";

import Chat from "./Chat";

/*
  THREE
*/
import { Canvas, extend } from "react-three-fiber";

/*
  MODELS
*/
import Game from "./Game.js";

/*
  ORBIT CONTROLS
*/
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls";
extend({ OrbitControls });

const App = () => {
  return (
    <>
      <Canvas
        style={{ height: 400, width: 800 }}
        camera={{ fov: 50, position: [0, 0, 70], aspect: 800 / 400 }}
      >
        <Game />
      </Canvas>
      <Chat />
    </>
  );
};
export default App;
