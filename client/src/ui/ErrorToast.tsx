import { useEffect, useState } from "react";
import socket from "../socket";

export default function ErrorToast() {
  const [msg, setMsg] = useState<string | null>(null);

  useEffect(() => {
    let timer = 0;
    const onError = (data: Record<string, unknown>) => {
      setMsg(String(data.reason ?? "Action failed"));
      clearTimeout(timer);
      timer = window.setTimeout(() => setMsg(null), 2500);
    };
    socket.on("error", onError);
    return () => {
      socket.off("error", onError);
      clearTimeout(timer);
    };
  }, []);

  if (!msg) return null;
  return <div className="toast">{msg}</div>;
}
