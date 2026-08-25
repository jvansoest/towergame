import { useEffect, useState } from "react";
import socket from "../socket";
import { repeatInspect, stopInspect } from "../inspect";
import type { Inspection } from "../protocol";

// How often an open card asks for fresh numbers.
const REFRESH_MS = 500;

export default function RoomInspector() {
  const [info, setInfo] = useState<Inspection | null>(null);

  useEffect(() => {
    const onInfo = (data: Record<string, unknown>) => {
      setInfo({
        title: String(data.title ?? ""),
        lines: (data.lines as string[]) ?? [],
      });
    };
    socket.on("inspect", onInfo);
    return () => socket.off("inspect", onInfo);
  }, []);

  const open = info !== null;
  useEffect(() => {
    if (!open) return;
    const timer = window.setInterval(repeatInspect, REFRESH_MS);
    return () => window.clearInterval(timer);
  }, [open]);

  if (!info) return null;
  return (
    <div className="inspector">
      <button
        className="inspector__close"
        onClick={() => {
          stopInspect();
          setInfo(null);
        }}
      >
        ×
      </button>
      <div className="inspector__title">{info.title}</div>
      {info.lines.map((line, i) => (
        <div key={i}>{line}</div>
      ))}
    </div>
  );
}
