import React, { useRef, useState } from "react";
import { extend, useThree, useFrame } from "react-three-fiber";
import Box from "./Box.js";

function BackGround(props) {
  const mesh = useRef();
  const [hovered, setHover] = useState(false);
  const { camera } = useThree();
  const [box, setBox] = useState({ x: -1, y: -1 });

  const handleHover = (e) => {
    const vec = props.getPlaneXYfromMouse(e.clientX, e.clientY, camera);
    const [posX, posY] = props.makeGrid(vec);
    setBox({ x: posX, y: posY });
    // console.log(props.lastCoords.x, props.lastCoords.y, " TO ", posX, posY);
    // if (props.lastCoords.x !== posX || props.lastCoords.y !== posY) {
    //   props.unSetMatrix(props.matrix, props.lastCoords.x, props.lastCoords.y);
    // }
    // props.setMatrix(props.matrix, posX, posY);
  };

  return (
    <>
      <Box
        key={box.y + "" + box.x}
        position={[(box.x - 10) * 10, (box.y - 10) * 5, -10]}
      />
      <mesh
        {...props}
        ref={mesh}
        onPointerMove={(e) => handleHover(e)}
        //onPointerOver={(e) => handleHover(true)}
        onPointerOut={(e) => setHover(false)}
      >
        <boxBufferGeometry attach="geometry" args={[200, 100, 0]} />
        <meshStandardMaterial
          attach="material"
          // color={"cyan"}
          color={hovered ? "brown" : "cyan"}
        />
      </mesh>
    </>
  );
}
//<Box position={[10, 10, 0]} />
export default BackGround;
