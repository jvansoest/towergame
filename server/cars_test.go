package main

import (
	"math"
	"testing"

	"towergame/server/transport"
)

// Builds a solid block with one shaft to the top.
func shaftWorld(t *testing.T, top int) (*World, transport.ShaftID) {
	t.Helper()
	w := newWorld()
	g := w.grid
	for f := 0; f <= top; f++ {
		for c := 18; c <= 26; c++ {
			g.Built[f][c] = true
		}
	}
	if err := w.PlaceElevator(top, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	return w, transport.ShaftID(g.Elevators[0].ID)
}

// A click on the shaft itself buys another car.
func TestClickOnShaftAddsCar(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	if got := len(w.bank.CarsIn(shaft)); got != 1 {
		t.Fatalf("new shaft has %d cars, want 1", got)
	}

	before := w.money
	if err := w.PlaceElevator(4, 20); err != nil { // inside the shaft
		t.Fatalf("add car: %v", err)
	}
	if got := len(w.bank.CarsIn(shaft)); got != 2 {
		t.Fatalf("shaft has %d cars, want 2", got)
	}
	if w.money != before-elevatorCarCost {
		t.Fatalf("money %d, want %d", w.money, before-elevatorCarCost)
	}

	// A click above the top still extends the shaft.
	if err := w.PlaceElevator(10, 20); err != nil {
		t.Fatalf("extend: %v", err)
	}
	if got := w.grid.Elevators[0].Top; got != 10 {
		t.Fatalf("shaft tops out at %d, want 10", got)
	}
	if got := len(w.bank.CarsIn(shaft)); got != 2 {
		t.Fatalf("extending changed the car count to %d", got)
	}
}

// A shaft takes only so many cars.
func TestShaftCarLimit(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	for len(w.bank.CarsIn(shaft)) < transport.MaxPerShaft {
		if err := w.AddCar(20); err != nil {
			t.Fatalf("car %d: %v", len(w.bank.CarsIn(shaft))+1, err)
		}
	}
	if err := w.AddCar(20); err == nil {
		t.Fatalf("shaft took more than %d cars", transport.MaxPerShaft)
	}
	if err := w.AddCar(40); err == nil {
		t.Fatal("added a car to a column with no elevator")
	}
}

// Two cars in one shaft each carry their own rider.
func TestBothCarsCarryRiders(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	if err := w.AddCar(20); err != nil {
		t.Fatalf("second car: %v", err)
	}
	ids := w.bank.CarsIn(shaft)
	w.bank.Place(ids[0], 0)
	w.bank.Place(ids[1], 8)

	// The near car answers each call.
	up, _, _ := w.bank.BestCar(shaft, 1, 4)
	down, _, _ := w.bank.BestCar(shaft, 7, 5)
	if up == down {
		t.Fatalf("both calls went to car %d", up)
	}

	a := &sim{
		id: 1, state: stateWaiting, shaft: shaft, car: up,
		boardFloor: 1, exitFloor: 4, x: 19, queueX: 19, y: 1,
	}
	bb := &sim{
		id: 2, state: stateWaiting, shaft: shaft, car: down,
		boardFloor: 7, exitFloor: 5, x: 19, queueX: 19, y: 7,
	}
	w.sims = []*sim{a, bb}

	for i := 0; i < 3000; i++ {
		w.Step(1.0 / 30)
		if a.state == stateMoving && bb.state == stateMoving {
			break
		}
	}
	if int(math.Round(a.y)) != 4 {
		t.Fatalf("rider A left at floor %.2f, want 4", a.y)
	}
	if int(math.Round(bb.y)) != 5 {
		t.Fatalf("rider B left at floor %.2f, want 5", bb.y)
	}
}

// After a drop-off the car takes the queue there.
// It turns round first, so it holds the doors once.
func TestPicksUpAfterDropOff(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	car := carIDOf(w, shaft)
	w.bank.Place(car, 3)

	rider := &sim{
		id: 1, state: stateRiding, shaft: shaft, car: car, slot: 0,
		boardFloor: 3, exitFloor: 0,
	}
	hall := &sim{
		id: 2, state: stateWaiting, shaft: shaft, car: car,
		boardFloor: 0, exitFloor: 5, x: 19, queueX: 19, y: 0,
	}
	w.sims = []*sim{rider, hall}

	for i := 0; i < 400 && hall.state == stateWaiting; i++ {
		w.Step(1.0 / 30)
	}
	if hall.state == stateWaiting {
		c, _ := w.bank.Car(car)
		t.Fatalf("nobody boarded: car at %.2f to %.2f dir %d", c.Floor, c.To, c.Dir)
	}
	if rider.state == stateRiding {
		t.Fatal("the car took the queue before dropping its rider")
	}
}

// Any car that opens its doors takes the queue.
func TestBoardsWhicheverCarOpens(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	if err := w.AddCar(20); err != nil {
		t.Fatalf("second car: %v", err)
	}
	ids := w.bank.CarsIn(shaft)
	mine, other := ids[0], ids[1]
	w.bank.Place(mine, 8)  // the car they were given
	w.bank.Place(other, 3) // the one actually here

	p := &sim{
		id: 1, state: stateWaiting, shaft: shaft, car: mine,
		boardFloor: 3, exitFloor: 6, x: 19, queueX: 19, y: 3,
	}
	w.sims = []*sim{p}

	w.boardQueues()
	if p.state != stateBoarding {
		t.Fatalf("sim stayed waiting beside an open car (state %d)", p.state)
	}
	if p.car != other {
		t.Fatalf("boarded car %d, want the one at their floor (%d)", p.car, other)
	}
}

// A full car must not keep a queue to itself.
func TestFullCarHandsOverTheQueue(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	if err := w.AddCar(20); err != nil {
		t.Fatalf("second car: %v", err)
	}
	ids := w.bank.CarsIn(shaft)
	full, free := ids[0], ids[1]
	w.bank.Place(full, 8)
	w.bank.Place(free, 0)

	// Four riders fill one car; it is busy going down.
	w.sims = nil
	for i := 0; i < transport.Capacity; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, state: stateRiding, shaft: shaft, car: full, slot: i,
			boardFloor: 8, exitFloor: 0,
		})
	}
	p := &sim{
		id: 99, state: stateWaiting, shaft: shaft, car: full,
		boardFloor: 3, exitFloor: 6, x: 19, queueX: 19, y: 3,
	}
	w.sims = append(w.sims, p)

	for i := 0; i < 60 && p.car == full; i++ {
		w.Step(1.0 / 30)
	}
	if p.car == full {
		t.Fatal("sim still waits on a car with no room")
	}
	if p.car != free {
		t.Fatalf("sim moved to car %d, want the free one (%d)", p.car, free)
	}
}

// A busy tower puts both cars to work.
func TestSecondCarInSeededTower(t *testing.T) {
	w := newWorld()
	w.seed()
	shaft := transport.ShaftID(w.grid.Elevators[0].ID)
	if err := w.AddCar(w.grid.Elevators[0].Col); err != nil {
		t.Fatalf("second car: %v", err)
	}
	ids := w.bank.CarsIn(shaft)

	carried := map[transport.CarID]bool{}
	for step := 0; step < 30*60*90; step++ { // 90 simulated minutes
		w.Step(1.0 / 30)
		for _, id := range ids {
			if c, ok := w.bank.Car(id); ok && c.Riders() > 0 {
				carried[id] = true
			}
		}
	}
	for _, id := range ids {
		if !carried[id] {
			t.Fatalf("car %d never carried anyone", id)
		}
	}
}
