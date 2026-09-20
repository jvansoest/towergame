import { Canvas } from "@react-three/fiber";

import Chat from "./Chat";
import Game from "./Game";
import Clock from "./ui/Clock";
import Budget from "./ui/Budget";
import CardLinks from "./ui/CardLinks";
import ErrorToast from "./ui/ErrorToast";
import ReportToast from "./ui/ReportToast";
import RoomInspector from "./ui/RoomInspector";
import RoomPalette from "./ui/RoomPalette";

const App = () => {
  return (
    <div className="app">
      <div className="hud">
        <Clock />
        <Budget />
      </div>
      <RoomPalette />
      <div className="canvas-wrap">
        <Canvas
          dpr={[1, 1.5]}
          camera={{ fov: 50, position: [0, 18, 60], near: 0.1, far: 2000 }}
        >
          <Game />
        </Canvas>
      </div>
      <Chat />
      <CardLinks />
      <ErrorToast />
      <ReportToast />
      <RoomInspector />
    </div>
  );
};

export default App;
