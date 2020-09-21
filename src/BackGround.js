import React, { useRef, useState } from "react";

function BackGround(props) {
  const mesh = useRef();
  const [hovered, setHover] = useState(false);
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
      // onClick={handleClick}
      onPointerOver={(e) => setHover(true)}
      onPointerOut={(e) => setHover(false)}
    >
      <boxBufferGeometry attach="geometry" args={[200, 100, 1]} />
      <meshStandardMaterial
        attach="material"
        color={hovered ? "brown" : "cyan"}
      />
    </mesh>
  );
}
//<Box position={[10, 10, 0]} />
export default BackGround;
