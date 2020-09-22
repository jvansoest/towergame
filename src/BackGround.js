import React, { useRef, useState } from "react";
import { extend, useThree, useFrame } from "react-three-fiber";

function BackGround(props) {
  const mesh = useRef();
  const [hovered, setHover] = useState(false);
  const { camera } = useThree();

  const handleHover = (e) => {
    const vec = props.getPlaneXYfromMouse(e.clientX, e.clientY, camera);
    const [posX, posY] = props.makeGrid(vec);
    if (props.lastCoords.x !== posX) {
      if (props.lastCoords.y !== posY) {
        props.unSetMatrix(props.matrix, props.lastCoords.x, props.lastCoords.y);
      }
    }
    props.setMatrix(props.matrix, posX, posY);
  };
  //const [active, setActive] = useState(false);
  // useFrame(() => {
  //   if (mesh.current.rotation.x < 1) mesh.current.rotation.x += 0.01;
  // });
  // const handleClick = (e) => {
  //   // setActive(!active);
  //   const clickedX = e.clientX;
  //   const clickedY = e.clientY;
  //   console.log(clickedX, clickedY);
  // };

  return (
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
  );
}
//<Box position={[10, 10, 0]} />
export default BackGround;
