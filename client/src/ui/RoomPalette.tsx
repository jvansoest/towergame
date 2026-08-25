import { setSelectedRoom, useSelectedRoom } from "../store";

const ROOMS = [
  { id: "base", name: "Base" },
  { id: "stairs", name: "Stairs" },
  { id: "elevator", name: "Elevator" },
  { id: "lobby", name: "Lobby" },
  { id: "office", name: "Office" },
  { id: "condo", name: "Condo" },
  { id: "hotel_single", name: "Hotel Room" },
  { id: "hotel_suite", name: "Suite" },
  { id: "shop", name: "Shop" },
  { id: "fastfood", name: "Fast Food" },
  { id: "restaurant", name: "Restaurant" },
  { id: "cinema", name: "Cinema" },
  { id: "inspect", name: "Inspect" },
  { id: "bulldoze", name: "Remove" },
];

export default function RoomPalette() {
  const selected = useSelectedRoom();
  return (
    <div className="palette">
      {ROOMS.map((r) => (
        <button
          key={r.id}
          className={"palette__btn" + (selected === r.id ? " is-active" : "")}
          onClick={() => setSelectedRoom(r.id)}
        >
          {r.name}
        </button>
      ))}
    </div>
  );
}
