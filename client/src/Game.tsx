import { Suspense, useEffect, useRef, useState } from "react";

import socket from "./socket";
import { setHover } from "./store";
import CameraControls from "./CameraControls";
import Skyline from "./Skyline";
import Scenery from "./Scenery";
import Tower, { ShaftGhost } from "./Tower";
import ElevatorCars from "./ElevatorCars";
import Vehicles from "./Vehicles";
import Train, { Tunnel } from "./Train";
import Sims from "./Sims";
import useBuildInput from "./useBuildInput";
import { CELL_W, CELL_H } from "./dims";
import type { GridView } from "./protocol";

// Size of the build surface before the first snapshot.
const PLANE_W = 200;
const PLANE_H = 120;

const Game = () => {
  const [grid, setGrid] = useState<GridView | null>(null);
  const { onPointerDown, onPointerMove, preview } = useBuildInput(grid);

  const gridSig = useRef("");
  useEffect(() => {
    const onSnapshot = (data: Record<string, unknown>) => {
      const g = data.grid as GridView | undefined;
      if (!g) return;
      // Rebuild only when the tower actually changed, not every clock tick.
      const sig = JSON.stringify(g);
      if (sig === gridSig.current) return;
      gridSig.current = sig;
      setGrid(g);
    };
    socket.on("snapshot", onSnapshot);
    return () => socket.off("snapshot", onSnapshot);
  }, []);

  const planeW = grid ? grid.width * CELL_W : PLANE_W;
  // From the deepest basement floor to the top.
  const deep = grid?.basement ?? 0;
  const planeH = grid ? (grid.floors + deep) * CELL_H : PLANE_H;
  const planeY = grid ? ((grid.floors - deep) * CELL_H) / 2 : PLANE_H / 2;

  return (
    <>
      <CameraControls />
      <Skyline />
      {/* Sun light, no distance falloff. */}
      <ambientLight intensity={0.7} />
      <directionalLight position={[8, 15, 10]} intensity={2.2} />

      {/* Invisible build surface. */}
      <mesh
        position={[0, planeY, 0]}
        onPointerDown={onPointerDown}
        onPointerMove={onPointerMove}
        onPointerOut={() => setHover(null)}
      >
        <planeGeometry args={[planeW, planeH]} />
        <meshBasicMaterial transparent opacity={0} depthWrite={false} />
      </mesh>

      <Scenery towerW={planeW} grid={grid} />

      {grid && <Tower grid={grid} />}
      {grid && preview && <ShaftGhost grid={grid} drag={preview} />}
      {grid && <ElevatorCars gridWidth={grid.width} />}
      {grid && <Vehicles gridWidth={grid.width} />}
      {grid && <Tunnel grid={grid} />}
      {grid && <Train grid={grid} />}

      {grid && (
        <Suspense fallback={null}>
          <Sims gridWidth={grid.width} />
        </Suspense>
      )}
    </>
  );
};

export default Game;
