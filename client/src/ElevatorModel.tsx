// Two side rails; whatever is behind stays visible.
export default function ElevatorModel({
  width,
  height,
}: {
  width: number;
  height: number;
}) {
  const rail = 0.18;
  return (
    <group>
      <mesh position={[-width / 2 + rail / 2, 0, 0]}>
        <boxGeometry args={[rail, height, 0.4]} />
        <meshStandardMaterial color="#374151" />
      </mesh>
      <mesh position={[width / 2 - rail / 2, 0, 0]}>
        <boxGeometry args={[rail, height, 0.4]} />
        <meshStandardMaterial color="#374151" />
      </mesh>
    </group>
  );
}
