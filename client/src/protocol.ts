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
  open?: boolean; // shops: within opening hours
  alignment?: string; // neutral, good, or triad
  dirty?: boolean; // hotel room awaits a maid
  booked?: boolean; // hotel room has guests
};

export type StairView = {
  floor: number;
  col: number;
  width: number;
  height: number;
  escalator?: boolean; // moves by itself, drawn hidden
};

// One run of a car ramp.
export type RampView = {
  floor: number; // the row it crosses, below ground
  col: number;
  width: number;
  downRight: boolean; // slants toward the right
};

// The subway train, while it runs.
export type TrainView = {
  floor: number;
  x: number; // centre column, may be off the grid
  stopped: boolean; // at the platform, doors open
};

// One car in the garage or on its ramp.
export type VehicleView = {
  id: number;
  x: number; // col
  y: number; // floor, may be fractional
  dir: number; // -1 faces left, 1 faces right
  slope: number; // floors per col, signed
  parked: boolean;
};

export type ElevatorView = {
  id: number;
  bottom: number;
  top: number;
  col: number;
  width: number;
};

// One buildable room type, from the server's
// catalog. The single source for seats, size,
// and category — cards and palette read it.
export type RoomTypeView = {
  id: string;
  name: string;
  width: number;
  height: number;
  cost: number;
  capacity: number;
  category: string;
  seats: number;
  habit: string;
  line?: boolean; // seats share one column
  rent?: number; // per guest, per stay
  noise?: number; // 0 quiet .. 3 loud
  late?: boolean; // open into the night
  stars?: number; // rating needed to build
  slots?: number; // cars a garage holds
  underground?: boolean; // built below ground only
};

export type GridView = {
  width: number;
  floors: number;
  built: boolean[][];
  rooms: RoomView[];
  stairs: StairView[];
  ramps: RampView[];
  below: boolean[][]; // floors -1, -2, ...
  basement: number; // floors below ground
  elevators: ElevatorView[];
  types?: RoomTypeView[];
};

export type GameClock = { day: number; hour: number; minute: number };

export type Snapshot = {
  tick: number;
  clock: GameClock;
  money: number;
  stars: number; // tower rating, 1 to 5
  population: number; // residents, staff, guests
  grid: GridView;
};

// One sim's position in grid coordinates.
export type SimView = {
  id: number;
  x: number; // col
  y: number; // floor
  riding?: boolean;
  sitting?: boolean;
  spot?: number; // seat index in the room
  gliding?: boolean; // riding an escalator
  cleaning?: boolean; // a maid is sweeping
  waiting?: boolean; // queued at a shaft
  dark?: boolean; // the room's lights are off
  profession?: string; // visitor, resident, worker, security, triad
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
  train: TrainView | null;
  vehicles: VehicleView[];
  cars: CarView[];
};

// What the inspector shows about one thing.
export type Inspection = {
  title: string;
  lines: string[];
};

// The books for one quarter, sent when it closes.
export type Report = {
  quarter: number;
  hotel: number; // guest rent
  sales: number; // customers
  leases: number; // offices and homes
  events: number; // VIP rewards, less penalties
  upkeep: number;
  net: number;
  money: number;
};
