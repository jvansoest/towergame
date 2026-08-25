import { useEffect, useState } from "react";
import socket from "../socket";
import type { GameClock } from "../protocol";

const pad = (n: number) => n.toString().padStart(2, "0");

export default function Clock() {
  const [clock, setClock] = useState<GameClock | null>(null);

  useEffect(() => {
    const onSnapshot = (data: Record<string, unknown>) => {
      const c = data.clock as GameClock | undefined;
      if (c) setClock(c);
    };
    socket.on("snapshot", onSnapshot);
    return () => socket.off("snapshot", onSnapshot);
  }, []);

  if (!clock) return null;

  return (
    <div className="clock">
      Day {clock.day} · {pad(clock.hour)}:{pad(clock.minute)}
    </div>
  );
}
