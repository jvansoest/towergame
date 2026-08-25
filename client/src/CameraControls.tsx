import { useEffect, useRef, type ComponentRef } from "react";
import { OrbitControls } from "@react-three/drei";

// Pan and zoom a flat tower view.
export default function CameraControls() {
  const controls = useRef<ComponentRef<typeof OrbitControls>>(null);

  useEffect(() => {
    const c = controls.current;
    if (!c) return;
    c.target.set(0, 18, 0);
    c.update();
  }, []);

  return (
    <OrbitControls
      ref={controls}
      makeDefault
      enableRotate={false}
      screenSpacePanning
      enableDamping
      dampingFactor={0.12}
      minDistance={20}
      maxDistance={260}
    />
  );
}
