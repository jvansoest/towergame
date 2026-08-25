import { useEffect, useState } from "react";
import socket from "../socket";

export default function Budget() {
  const [money, setMoney] = useState<number | null>(null);

  useEffect(() => {
    const onSnapshot = (data: Record<string, unknown>) => {
      if (typeof data.money === "number") setMoney(data.money);
    };
    socket.on("snapshot", onSnapshot);
    return () => socket.off("snapshot", onSnapshot);
  }, []);

  if (money === null) return null;
  return <div className="budget">${money.toLocaleString()}</div>;
}
