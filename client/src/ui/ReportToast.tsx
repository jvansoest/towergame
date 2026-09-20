import { useEffect, useState } from "react";
import socket from "../socket";
import type { Report } from "../protocol";

const SHOW_MS = 20000;

const money = (n: number) =>
  (n < 0 ? "-$" : "$") + Math.abs(n).toLocaleString();

// The quarterly books, shown for a while.
export default function ReportToast() {
  const [rep, setRep] = useState<Report | null>(null);

  useEffect(() => {
    let timer = 0;
    const onReport = (data: Record<string, unknown>) => {
      setRep(data as unknown as Report);
      clearTimeout(timer);
      timer = window.setTimeout(() => setRep(null), SHOW_MS);
    };
    socket.on("report", onReport);
    return () => {
      socket.off("report", onReport);
      clearTimeout(timer);
    };
  }, []);

  if (!rep) return null;
  const rows: [string, number][] = [
    ["Leases", rep.leases],
    ["Sales", rep.sales],
    ["Hotel", rep.hotel],
    ["Events", rep.events],
    ["Upkeep", -rep.upkeep],
  ];
  return (
    <div className="report" onClick={() => setRep(null)}>
      <div className="report__title">Quarter {rep.quarter}</div>
      {rows.map(([name, v]) => (
        <div className="report__row" key={name}>
          <span>{name}</span>
          <span>{money(v)}</span>
        </div>
      ))}
      <div className="report__row report__net">
        <span>Net</span>
        <span>{money(rep.net)}</span>
      </div>
    </div>
  );
}
