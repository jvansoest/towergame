import React, { useRef } from "react";
//import { useFrame } from "react-three-fiber";

function GroundFloor(props) {
  const mesh = useRef();
  //const [hovered, setHover] = useState(false);
  //const [active, setActive] = useState(false);

  return (
    <mesh {...props} ref={mesh}>
      <boxBufferGeometry attach="geometry" args={[150, 1, 5]} />
      <meshStandardMaterial attach="material" color={"grey"} />
    </mesh>
  );
}
export default GroundFloor;
