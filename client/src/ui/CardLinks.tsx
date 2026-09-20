import { useTypes } from "../catalog";

// Very small links to every room's card page:
// the style sheets we tune room art against.
export default function CardLinks() {
  const types = useTypes();
  return (
    <div className="cardlinks">
      {types.map((t) => (
        <a key={t.id} href={`/card/${t.id}`}>
          {t.id}
        </a>
      ))}
    </div>
  );
}
