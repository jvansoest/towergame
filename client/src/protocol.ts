// The wire format, mirroring server/messages.go.
// Both sides must agree, so the client states it once.

export type RoomView = {
  type: string;
  floor: number;
  col: number;
  width: number;
  height: number;
  category: string;
  capacity: number;
  seats: number;
  habit: string;
  lit: boolean;
  occupants: number;
};

export type StairView = {
  floor: number;
  col: number;
  width: number;
  height: number;
};

export type ElevatorView = {
  id: number;
  bottom: number;
  top: number;
  col: number;
  width: number;
};

export type GridView = {
  width: number;
  floors: number;
  built: boolean[][];
  rooms: RoomView[];
  stairs: StairView[];
  elevators: ElevatorView[];
};

export type GameClock = { day: number; hour: number; minute: number };

export type Snapshot = {
  tick: number;
  clock: GameClock;
  money: number;
  grid: GridView;
};

// One sim's position in grid coordinates.
export type SimView = {
  id: number;
  x: number; // col
  y: number; // floor
  riding?: boolean;
  sitting?: boolean;
  dark?: boolean; // the room's lights are off
};

// One elevator car's position.
export type CarView = {
  id: number;
  shaft: number;
  col: number;
  floor: number;
};

export type SimsUpdate = {
  sims: SimView[];
  cars: CarView[];
};

// What the inspector shows about one thing.
export type Inspection = {
  title: string;
  lines: string[];
};
