package transport

import (
	"math"
	"testing"
)

const dt = 1.0 / 30

// Runs a bank until its cars settle, or the limit is hit.
// Calls are keyed by car: the first car added is 1.
func run(b *Bank, calls map[CarID][]Call, ticks int) {
	for i := 0; i < ticks; i++ {
		b.Step(dt, calls)
	}
}

func TestCarGoesToItsCall(t *testing.T) {
	b := &Bank{Floors: 10}
	id := b.Add(1)
	run(b, map[CarID][]Call{1: {{Floor: 5, Dir: -1}}}, 60)
	_ = id
	if math.Abs(b.cars[0].Floor-5) > LevelEps {
		t.Fatalf("car at %.2f, want 5", b.cars[0].Floor)
	}
}

// The classic collective rule: serve calls on the way.
func TestStopsOnTheWayUp(t *testing.T) {
	b := &Bank{Floors: 10}
	b.Add(1)
	calls := map[CarID][]Call{1: {
		{Floor: 8},          // a rider going to 8
		{Floor: 4, Dir: 1},  // someone at 4 going up
		{Floor: 6, Dir: -1}, // someone at 6 going down
	}}
	b.cars[0].Dir = 1

	stops := []float64{}
	for i := 0; i < 200; i++ {
		before := b.cars[0].To
		b.Step(dt, calls)
		if b.cars[0].To != before {
			stops = append(stops, b.cars[0].To)
		}
	}
	if len(stops) == 0 || stops[0] != 4 {
		t.Fatalf("first stop %v, want floor 4 on the way up", stops)
	}
	for _, s := range stops {
		if s == 6 {
			t.Fatal("stopped for a down call while going up")
		}
	}
}

// A full car runs its drop-offs and ignores hall calls.
func TestFullCarSkipsPickups(t *testing.T) {
	b := &Bank{Floors: 10}
	b.Add(1)
	for i := 0; i < Capacity; i++ {
		if !b.Board(b.cars[0].ID) {
			t.Fatalf("boarding %d of %d was refused", i+1, Capacity)
		}
	}
	if b.Board(b.cars[0].ID) {
		t.Fatal("an overfull car accepted another rider")
	}
	if b.HasRoom(b.cars[0].ID) {
		t.Fatal("a full car should report no room")
	}
	// The caller is expected to drop pickups; drop-offs still run.
	b.cars[0].dwell = 0 // boarding held the doors
	run(b, map[CarID][]Call{1: {{Floor: 3}}}, 60)
	if math.Abs(b.cars[0].Floor-3) > LevelEps {
		t.Fatalf("car at %.2f, want its drop-off at 3", b.cars[0].Floor)
	}
}

// Everyone below wanting up must not strand a car at the top.
func TestFetchesOppositeCalls(t *testing.T) {
	b := &Bank{Floors: 10}
	b.Add(1)
	b.cars[0].Floor = 7
	b.cars[0].Dir = 1
	calls := map[CarID][]Call{1: {{Floor: 2, Dir: 1}}}

	run(b, calls, 200)
	if math.Abs(b.cars[0].Floor-2) > LevelEps {
		t.Fatalf("car stranded at %.2f, want it to fetch floor 2", b.cars[0].Floor)
	}
	if b.cars[0].Dir != 1 {
		t.Fatalf("car dir %d on arrival, want their up direction", b.cars[0].Dir)
	}
}

func TestDoorsHoldTheCar(t *testing.T) {
	b := &Bank{Floors: 10}
	b.Add(1)
	b.HoldDoors(b.cars[0].ID)
	run(b, map[CarID][]Call{1: {{Floor: 9}}}, 10) // a third of a second
	if b.cars[0].Floor != 0 {
		t.Fatalf("car moved to %.2f with its doors open", b.cars[0].Floor)
	}
}

func TestStoppedAtOnlyWhenStopping(t *testing.T) {
	b := &Bank{Floors: 10}
	b.Add(1)
	b.cars[0].Floor = 4
	b.cars[0].To = 8
	if b.StoppedAt(b.cars[0].ID, 4) {
		t.Fatal("a car passing floor 4 must not count as stopped")
	}
	b.cars[0].To = 4
	if !b.StoppedAt(b.cars[0].ID, 4) {
		t.Fatal("a car targeting floor 4 should count as stopped")
	}
}

// Nearest Car: a near idle car beats a far or departing one.
func TestCostPrefersNearIdleCar(t *testing.T) {
	b := &Bank{Floors: 10}
	b.Add(1)
	b.Add(1)
	b.cars[0].Floor = 1 // idle, close
	b.cars[1].Floor = 9 // idle, far
	if b.CarCost(b.cars[0].ID, 0, 5) >= b.CarCost(b.cars[1].ID, 0, 5) {
		t.Fatal("the near car should cost less")
	}

	b.cars[0].riders = Capacity // now full
	if b.CarCost(b.cars[0].ID, 0, 5) <= b.CarCost(b.cars[1].ID, 0, 5) {
		t.Fatal("a full car should cost more than a far empty one")
	}
}

// Removing a shaft leaves every other car's id valid.
func TestRemoveShaftKeepsOtherIDs(t *testing.T) {
	b := &Bank{Floors: 10}
	a := b.Add(1)
	doomed := b.Add(2)
	c := b.Add(3)
	b.cars[b.at(c)].Floor = 5

	b.RemoveShaft(2)

	if _, ok := b.Car(doomed); ok {
		t.Fatal("the removed shaft's car is still here")
	}
	// Ids are names, not positions: both survivors still resolve.
	if _, ok := b.Car(a); !ok {
		t.Fatal("car a lost its id when another shaft went")
	}
	got, ok := b.Car(c)
	if !ok || got.Floor != 5 {
		t.Fatalf("car c resolved to %+v, want floor 5", got)
	}
	b.RemoveShaft(99) // unknown shaft must not panic
}

// A shaft may hold more than one car.
func TestShaftHoldsSeveralCars(t *testing.T) {
	b := &Bank{Floors: 20}
	low := b.Add(1)
	high := b.Add(1)
	if got := b.CarsIn(1); len(got) != 2 {
		t.Fatalf("shaft 1 has %d cars, want 2", len(got))
	}
	b.cars[b.at(low)].Floor = 0
	b.cars[b.at(high)].Floor = 18

	// The group controller sends the nearer car.
	pick, _, ok := b.BestCar(1, 17, 19)
	if !ok || pick != high {
		t.Fatalf("picked car %d for a call at 17, want the high car %d", pick, high)
	}
	pick, _, ok = b.BestCar(1, 1, 5)
	if !ok || pick != low {
		t.Fatalf("picked car %d for a call at 1, want the low car %d", pick, low)
	}

	// A call belongs to one car, so the other stays put.
	run(b, map[CarID][]Call{low: {{Floor: 9, Dir: 1}}}, 300)
	if math.Abs(b.cars[b.at(low)].Floor-9) > LevelEps {
		t.Fatalf("called car at %.2f, want floor 9", b.cars[b.at(low)].Floor)
	}
	if b.cars[b.at(high)].Floor != 18 {
		t.Fatalf("uncalled car moved to %.2f", b.cars[b.at(high)].Floor)
	}
}

// Cars share a shaft: they stack and pass.
func TestCarsPassEachOther(t *testing.T) {
	b := &Bank{Floors: 20}
	low := b.Add(1)
	high := b.Add(1)
	b.cars[b.at(low)].Floor = 2
	b.cars[b.at(high)].Floor = 12

	// Each is sent to the other's floor.
	calls := map[CarID][]Call{low: {{Floor: 12}}, high: {{Floor: 2}}}
	run(b, calls, 300)
	if !b.StoppedAt(low, 12) || !b.StoppedAt(high, 2) {
		t.Fatalf("cars at %.2f and %.2f, want them swapped",
			b.cars[b.at(low)].Floor, b.cars[b.at(high)].Floor)
	}
}

// Two cars may rest on the same floor.
func TestCarsShareAFloor(t *testing.T) {
	b := &Bank{Floors: 20}
	one := b.Add(1)
	two := b.Add(1)
	b.cars[b.at(two)].Floor = 9

	run(b, map[CarID][]Call{one: {{Floor: 9}}, two: {{Floor: 9}}}, 300)
	if !b.StoppedAt(one, 9) || !b.StoppedAt(two, 9) {
		t.Fatalf("cars at %.2f and %.2f, want both at 9",
			b.cars[b.at(one)].Floor, b.cars[b.at(two)].Floor)
	}
}
