import React, { useRef, useState } from "react";

function Box(props) {
  const mesh = useRef();
  const [hovered, setHover] = useState(false);
  const [active, setActive] = useState(false);
  // useFrame(() => {
  //   if (mesh.current.rotation.x < 1) mesh.current.rotation.x += 0.01;
  // });

  return (
    <mesh
      {...props}
      ref={mesh}
      onClick={(e) => setActive(!active)}
      onPointerOver={(e) => setHover(true)}
      onPointerOut={(e) => setHover(false)}
    >
      <boxBufferGeometry attach="geometry" args={[10, 5, 5]} />
      <meshStandardMaterial
        attach="material"
        color={hovered ? "hotpink" : "grey"}
      />
    </mesh>
  );
}
//<Box position={[10, 10, 0]} />
export default Box;
