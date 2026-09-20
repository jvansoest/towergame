import { useTypes } from "../catalog";
import { setSelectedRoom, useSelectedRoom } from "../store";
import { useRating } from "../rating";

// Tools first, rooms from the server catalog,
// then the rest of the tools.
type Entry = { id: string; name: string; stars?: number };
const HEAD: Entry[] = [
  { id: "base", name: "Base" },
  { id: "stairs", name: "Stairs" },
  { id: "escalator", name: "Escalator", stars: 2 },
  { id: "ramp", name: "Car ramp", stars: 2 },
  { id: "elevator", name: "Elevator" },
];
const TAIL: Entry[] = [
  { id: "inspect", name: "Inspect" },
  { id: "bulldoze", name: "Remove" },
];

export default function RoomPalette() {
  const selected = useSelectedRoom();
  const types = useTypes();
  const { stars } = useRating();
  const entries: Entry[] = [
    ...HEAD,
    ...types.map((t) => ({ id: t.id, name: t.name, stars: t.stars })),
    ...TAIL,
  ];
  return (
    <div className="palette">
      {entries.map((r) => (
        <button
          key={r.id}
          className={"palette__btn" + (selected === r.id ? " is-active" : "")}
          disabled={(r.stars ?? 1) > stars}
          title={(r.stars ?? 1) > stars ? `Needs ${r.stars} stars` : undefined}
          onClick={() => setSelectedRoom(r.id)}
        >
          <span>{r.name}</span>
        </button>
      ))}
    </div>
  );
}
