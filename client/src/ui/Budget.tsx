import { useEffect, useState } from "react";
import socket from "../socket";
import { useRating } from "../rating";

export default function Budget() {
  const [money, setMoney] = useState<number | null>(null);
  const { stars, population } = useRating();

  useEffect(() => {
    const onSnapshot = (data: Record<string, unknown>) => {
      if (typeof data.money === "number") setMoney(data.money);
    };
    socket.on("snapshot", onSnapshot);
    return () => socket.off("snapshot", onSnapshot);
  }, []);

  if (money === null) return null;
  return (
    <>
      <div className="budget">${money.toLocaleString()}</div>
      <div className="rating" title={`${population} people`}>
        {"★".repeat(stars)}
        {"☆".repeat(Math.max(0, 5 - stars))}
        <span className="rating__pop">{population}</span>
      </div>
    </>
  );
}
