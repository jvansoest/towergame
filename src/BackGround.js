import React, { useRef, useState } from "react";
import Box from "./Box.js";
import { gridToCoords } from "./Utilities.js";

function BackGround(props) {
  const mesh = useRef();
  const [hovered, setHover] = useState(false);
  const [box, setBox] = useState({ x: -1, y: -1 });

  const handleHover = (e) => {
    const vec = props.getPlaneXYfromMouse();
    const [posX, posY] = props.makeGrid(vec);
    setBox({ x: posX, y: posY });
  };

  return (
    <>
      <Box
        key={box.y + "" + box.x}
        position={gridToCoords(box.x, box.y, -10)}

        //position={[(box.x - 10) * 10, (box.y - 10) * 5, -10]}
      />
      <mesh
        {...props}
        ref={mesh}
        onPointerMove={(e) => handleHover(e)}
        onPointerOut={(e) => setHover(false)}
      >
        <boxBufferGeometry attach="geometry" args={[200, 100, 0]} />
        <meshStandardMaterial
          attach="material"
          color={hovered ? "brown" : "cyan"}
        />
      </mesh>
    </>
  );
}
//<Box position={[10, 10, 0]} />
export default BackGround;
