package main

import (
	"math"
	"sort"
	"testing"

	"towergame/server/model"
	"towergame/server/transport"
)

func TestElevatorRide(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 6; f++ {
		for c := 18; c <= 34; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(6, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)

	p := &sim{id: 1, category: model.CategoryOffice, homeF: 6, homeC: 27, state: stateOutside}
	w.sims = []*sim{p}
	w.simTime = float64(10*60-startMinuteOfDay) / gameMinutesPerRealSecond // 10:00 weekday

	for i := 0; i < 900; i++ {
		w.Step(1.0 / 30)
	}
	if p.state != statePresent || int(math.Round(p.y)) != 6 {
		t.Fatalf("worker should reach the office on floor 6, got state=%d y=%.1f", p.state, p.y)
	}
}

func TestPresentNow(t *testing.T) {
	// Office: weekday 9-17 only.
	if !presentNow(model.CategoryOffice, true, 10*60, 0, 0) {
		t.Fatal("office worker should be in at 10:00 weekday")
	}
	if presentNow(model.CategoryOffice, true, 20*60, 0, 0) {
		t.Fatal("office worker should be gone at 20:00")
	}
	if presentNow(model.CategoryOffice, false, 10*60, 0, 0) {
		t.Fatal("office worker should not work weekends")
	}
	// Resident: home except the 8:30-14:00 weekday window.
	if presentNow(model.CategoryResidential, true, 10*60, 0, 0) {
		t.Fatal("resident should be out mid-morning weekday")
	}
	if !presentNow(model.CategoryResidential, true, 20*60, 0, 0) {
		t.Fatal("resident should be home in the evening")
	}
	// Hotel guest: overnight.
	if !presentNow(model.CategoryHotel, true, 23*60, 0, 0) {
		t.Fatal("hotel guest should be in at 23:00")
	}
	if presentNow(model.CategoryHotel, true, 12*60, 0, 0) {
		t.Fatal("hotel guest should be checked out at noon")
	}
}

// Personal hours shift the start of the day.
func TestScheduleShifts(t *testing.T) {
	// An early bird is in before nine.
	if !presentNow(model.CategoryOffice, true, 8*60+30, -40, 0) {
		t.Fatal("early worker should be in at 08:30")
	}
	// A late riser is not.
	if presentNow(model.CategoryOffice, true, 9*60+20, 30, 0) {
		t.Fatal("late worker should still be out at 09:20")
	}
	// One stays past five, another leaves before.
	if !presentNow(model.CategoryOffice, true, 17*60+20, 0, 40) {
		t.Fatal("worker on late hours should still be in")
	}
	if presentNow(model.CategoryOffice, true, 16*60+50, 0, -30) {
		t.Fatal("worker on early hours should have gone")
	}
}

func TestAddOccupants(t *testing.T) {
	w := newWorld()
	w.money = startMoney
	if err := w.Place("office", 0, 0); err != nil {
		t.Fatalf("place: %v", err)
	}
	if len(w.sims) != model.RoomTypes["office"].Capacity {
		t.Fatalf("want %d occupants, got %d",
			model.RoomTypes["office"].Capacity, len(w.sims))
	}
}

func TestWeekday(t *testing.T) {
	if !isWeekday(1) || !isWeekday(2) {
		t.Fatal("days 1-2 are weekdays")
	}
	if isWeekday(3) {
		t.Fatal("day 3 is the weekend")
	}
}

// Riders stand side by side inside the car.
func TestRidersSpread(t *testing.T) {
	w := newWorld()
	g := w.grid
	for c := 18; c <= 24; c++ {
		g.Built[0][c] = true
	}
	if _, err := g.PlaceElevator(3, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)
	setCarFloor(w, shaft, 1.5)

	for i := 0; i < transport.Capacity; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, state: stateRiding, shaft: shaft,
			car: carIDOf(w, shaft), slot: i,
		})
	}
	w.layoutRiders()

	center := w.shaftCol(shaft)
	seen := map[float64]bool{}
	for _, p := range w.sims {
		if p.state != stateRiding {
			continue
		}
		if seen[p.x] {
			t.Fatalf("two riders share column %.2f", p.x)
		}
		seen[p.x] = true
		if math.Abs(p.x-center) > float64(model.ElevatorWidth)/2 {
			t.Fatalf("rider %.2f outside the car", p.x)
		}
		if p.y != 1.5 {
			t.Fatalf("rider y = %.2f, want the car floor", p.y)
		}
	}
	if len(seen) != transport.Capacity {
		t.Fatalf("placed %d riders, want %d", len(seen), transport.Capacity)
	}
}

// Boarding is a walk, not a jump.
func TestBoardingWalksIn(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 6; f++ {
		for c := 18; c <= 34; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(6, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)

	// A queue, so later boarders start far from the car.
	for i := 0; i < transport.Capacity; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, category: model.CategoryOffice,
			homeF: 6, homeC: 27, state: stateOutside,
		})
	}
	w.simTime = float64(10*60-startMinuteOfDay) / gameMinutesPerRealSecond

	const dt = 1.0 / 30
	maxStep := walkSpeed*dt + 1e-6
	prev := map[int]float64{}
	boarded := 0
	for i := 0; i < 900; i++ {
		for _, p := range w.sims {
			if p.state == stateBoarding {
				prev[p.id] = p.x
			}
		}
		w.Step(dt)
		for _, p := range w.sims {
			x0, ok := prev[p.id]
			if !ok || (p.state != stateBoarding && p.state != stateRiding) {
				continue
			}
			if d := math.Abs(p.x - x0); d > maxStep {
				t.Fatalf("sim %d jumped %.2f cells while boarding", p.id, d)
			}
			delete(prev, p.id)
			if p.state == stateRiding {
				boarded++
			}
		}
	}
	if boarded == 0 {
		t.Fatal("no sim ever boarded")
	}
}

// The longest waiter boards first.
func TestQueueBoardsInOrder(t *testing.T) {
	w := newWorld()
	g := w.grid
	for c := 18; c <= 24; c++ {
		g.Built[0][c] = true
	}
	if _, err := g.PlaceElevator(3, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)

	// Shortest wait first, so list order fights queue order.
	for i := 0; i < transport.Capacity; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, state: stateWaiting, shaft: shaft,
			car:        carIDOf(w, shaft),
			boardFloor: 0, exitFloor: 3, waitTime: float64(i),
		})
	}
	w.boardQueues()

	for _, p := range w.sims {
		if p.state != stateBoarding {
			t.Fatalf("sim %d did not board", p.id)
		}
	}
	// Slots fill from the far side, longest waiter first.
	sorted := append([]*sim(nil), w.sims...)
	sort.Slice(sorted, func(a, b int) bool { return sorted[a].slot > sorted[b].slot })
	for i := 1; i < len(sorted); i++ {
		if sorted[i-1].waitTime < sorted[i].waitTime {
			t.Fatalf("sim %d boarded before longer waiter %d", sorted[i-1].id, sorted[i].id)
		}
	}
}

// A full car turns the rest away.
func TestCarCapacityLimit(t *testing.T) {
	w := newWorld()
	g := w.grid
	for c := 18; c <= 24; c++ {
		g.Built[0][c] = true
	}
	if _, err := g.PlaceElevator(3, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)

	for i := 0; i < transport.Capacity+2; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, state: stateWaiting, shaft: shaft,
			car:        carIDOf(w, shaft),
			boardFloor: 0, exitFloor: 3, waitTime: float64(i),
		})
	}
	w.boardQueues()

	left := 0
	for _, p := range w.sims {
		if p.state == stateWaiting {
			left++
		}
	}
	if carOf(w, shaft).Riders() != transport.Capacity {
		t.Fatalf("car holds %d, want %d", carOf(w, shaft).Riders(), transport.Capacity)
	}
	if left != 2 {
		t.Fatalf("%d sims still waiting, want 2", left)
	}
}

// Riders leave only when the car is level and still.
func TestAlightsOnlyWhenLevel(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 6; f++ {
		for c := 18; c <= 34; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(6, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)

	for i := 0; i < 3; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, category: model.CategoryOffice,
			homeF: 6, homeC: 27, state: stateOutside,
		})
	}
	w.simTime = float64(10*60-startMinuteOfDay) / gameMinutesPerRealSecond

	const dt = 1.0 / 30
	wasRiding := map[int]bool{}
	for i := 0; i < 900; i++ {
		w.Step(dt)
		for _, p := range w.sims {
			if p.state == stateRiding {
				wasRiding[p.id] = true
				continue
			}
			if !wasRiding[p.id] {
				continue
			}
			wasRiding[p.id] = false
			// Just stepped out: must be on a whole floor.
			if math.Abs(p.y-math.Round(p.y)) > transport.LevelEps {
				t.Fatalf("sim %d left the car at y=%.3f", p.id, p.y)
			}
			if carOf(w, shaft).Floor != math.Round(p.y) {
				t.Fatalf("sim %d left while the car was at %.3f", p.id, carOf(w, shaft).Floor)
			}
		}
	}
	if len(wasRiding) == 0 {
		t.Fatal("no sim ever rode the elevator")
	}
}

// Riders hold their slot for the whole ride.
func TestRiderSlotIsFixed(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 6; f++ {
		for c := 18; c <= 34; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(6, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)

	// Different floors, so riders alight one at a time.
	for i := 0; i < 3; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, category: model.CategoryOffice,
			homeF: 4 + i, homeC: 27, state: stateOutside,
		})
	}
	w.simTime = float64(10*60-startMinuteOfDay) / gameMinutesPerRealSecond

	rideX := map[int]float64{}
	rode := 0
	for i := 0; i < 900; i++ {
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			if p.state != stateRiding {
				continue
			}
			x0, seen := rideX[p.id]
			if !seen {
				rideX[p.id] = p.x
				rode++
				continue
			}
			if math.Abs(p.x-x0) > transport.LevelEps {
				t.Fatalf("rider %d slid from %.3f to %.3f mid-ride", p.id, x0, p.x)
			}
		}
	}
	if rode == 0 {
		t.Fatal("no sim ever rode the elevator")
	}
}

// Removing a room frees the boarder's place in the car.
func TestRemoveRoomFreesCarPlace(t *testing.T) {
	w := newWorld()
	g := w.grid
	for c := 0; c < g.Width; c++ {
		g.Built[0][c] = true
		g.Built[1][c] = true
	}
	if err := w.Place("office", 1, 0); err != nil {
		t.Fatalf("place: %v", err)
	}
	if _, err := g.PlaceElevator(1, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)

	// Every occupant is mid-board, holding a place.
	for _, p := range w.sims {
		p.state = stateBoarding
		p.shaft = shaft
		p.car = carIDOf(w, shaft)
	}
	w.recountRiders()
	if carOf(w, shaft).Riders() == 0 {
		t.Fatal("setup: nobody is holding a place")
	}
	if err := w.Remove(1, 0); err != nil {
		t.Fatalf("remove: %v", err)
	}
	w.Step(1.0 / 30)

	if carOf(w, shaft).Riders() != 0 {
		t.Fatalf("car still holds %d places after the room went", carOf(w, shaft).Riders())
	}
}

// A schedule change must not strand a car.
func TestGoalChangeFreesCarPlace(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 3; f++ {
		for c := 0; c < g.Width; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(3, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)

	// Riders whose goal is about to change under them.
	for i := 0; i < transport.Capacity; i++ {
		w.sims = append(w.sims, &sim{
			id: i + 1, category: model.CategoryOffice,
			homeF: 3, homeC: 27, state: stateRiding, shaft: shaft,
			car: carIDOf(w, shaft), slot: i, exitFloor: 3,
			last: goal{true, 3, 27},
		})
	}
	w.recountRiders()
	// 17:00: every office worker goes home.
	w.simTime = float64(17*60-startMinuteOfDay) / gameMinutesPerRealSecond

	for i := 0; i < 600; i++ {
		w.Step(1.0 / 30)
	}
	if carOf(w, shaft).Riders() >= transport.Capacity {
		t.Fatalf("car stuck holding %d places", carOf(w, shaft).Riders())
	}
}

// Equal waits keep a stable queue order.
func TestQueueOrderIsStable(t *testing.T) {
	ps := []*sim{{id: 3}, {id: 1}, {id: 2}}
	sortQueue(ps)
	first := []int{ps[0].id, ps[1].id, ps[2].id}

	ps = []*sim{{id: 2}, {id: 3}, {id: 1}}
	sortQueue(ps)
	for i, p := range ps {
		if p.id != first[i] {
			t.Fatalf("queue order changed: %v then %d at %d", first, p.id, i)
		}
	}
}

// A busy shaft sends new sims to the free one.
func TestDispatchAvoidsBusyCar(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 5; f++ {
		for c := 0; c < g.Width; c++ {
			g.Built[f][c] = true
		}
	}
	// Two shafts, the near one at 24, the far one at 34.
	if _, err := g.PlaceElevator(5, 24); err != nil {
		t.Fatalf("near elevator: %v", err)
	}
	if _, err := g.PlaceElevator(5, 34); err != nil {
		t.Fatalf("far elevator: %v", err)
	}
	near := transport.ShaftID(g.Elevators[0].ID)
	far := transport.ShaftID(g.Elevators[1].ID)
	w.addCar(g.Elevators[0].ID)
	w.addCar(g.Elevators[1].ID)

	// Both idle: the near shaft wins.
	if got := chosenShaft(t, w, 0, 26, 5, 26); got != near {
		t.Fatalf("idle: took shaft %d, want %d", got, near)
	}

	// Near car full, climbing away from the call.
	fillCar(w, near)
	w.bank.Place(carIDOf(w, near), 3)
	if got := chosenShaft(t, w, 0, 26, 5, 26); got != far {
		t.Fatalf("busy: took shaft %d, want %d", got, far)
	}
}

// Shaft a route picks, or zero.
func chosenShaft(t *testing.T, w *World, ff, fc, tf, tc int) transport.ShaftID {
	t.Helper()
	wp, ok := w.grid.PathVia(ff, fc, tf, tc, w.shaftCost)
	if !ok {
		t.Fatal("no path")
	}
	for _, p := range wp {
		if p.Mode == model.ModeElevator {
			return transport.ShaftID(p.Shaft)
		}
	}
	return 0
}

// A car with room stops for same-way calls on its route.
func TestPicksUpOnTheWay(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 8; f++ {
		for c := 18; c <= 26; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(8, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)
	car := carIDOf(w, shaft)

	// A rider going 0 -> 8, and someone waiting at 4 going up.
	rider := &sim{
		id: 1, state: stateRiding, shaft: shaft, car: car, slot: 0,
		boardFloor: 0, exitFloor: 8,
	}
	hall := &sim{
		id: 2, state: stateWaiting, shaft: shaft, car: car,
		boardFloor: 4, exitFloor: 8, queueX: 19,
		x: 19, y: 4,
	}
	w.sims = []*sim{rider, hall}

	for i := 0; i < 600 && hall.state == stateWaiting; i++ {
		w.Step(1.0 / 30)
	}
	if hall.state == stateWaiting {
		t.Fatal("car passed a waiting sim going its way")
	}
	// Picked up mid-trip, not after finishing the first ride.
	if rider.state != stateRiding {
		t.Fatal("pickup happened only after the first rider got off")
	}
	if math.Abs(carOf(w, shaft).Floor-4) > transport.LevelEps {
		t.Fatalf("car was at floor %.2f, want the call at 4", carOf(w, shaft).Floor)
	}
}

// A full car does not stop for more passengers.
func TestFullCarPassesBy(t *testing.T) {
	w := newWorld()
	g := w.grid
	for f := 0; f <= 8; f++ {
		for c := 18; c <= 26; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(8, 20); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	w.addCar(g.Elevators[len(g.Elevators)-1].ID)
	shaft := transport.ShaftID(g.Elevators[len(g.Elevators)-1].ID)
	car := carIDOf(w, shaft)

	var sims []*sim
	for i := 0; i < transport.Capacity; i++ {
		sims = append(sims, &sim{
			id: i + 1, state: stateRiding, shaft: shaft, car: car, slot: i,
			boardFloor: 0, exitFloor: 8,
		})
	}
	hall := &sim{
		id: 99, state: stateWaiting, shaft: shaft, car: car,
		boardFloor: 4, exitFloor: 8, queueX: 19,
		x: 19, y: 4,
	}
	w.sims = append(sims, hall)

	for i := 0; i < 300; i++ {
		w.Step(1.0 / 30)
		if carOf(w, shaft).Floor > 6 {
			break
		}
	}
	if hall.state != stateWaiting {
		t.Fatal("a full car should not have picked anyone up")
	}
}

// The one car serving a shaft, for tests.
func carOf(w *World, shaft transport.ShaftID) transport.Car {
	ids := w.bank.CarsIn(shaft)
	if len(ids) == 0 {
		return transport.Car{}
	}
	c, _ := w.bank.Car(ids[0])
	return c
}

// Puts the shaft's car at a floor.
func setCarFloor(w *World, shaft transport.ShaftID, floor float64) {
	w.bank.Place(carIDOf(w, shaft), floor)
}

// Fills the shaft's car to capacity.
func fillCar(w *World, shaft transport.ShaftID) {
	w.bank.SyncRiders(map[transport.CarID]int{
		carIDOf(w, shaft): transport.Capacity,
	})
}

func carIDOf(w *World, shaft transport.ShaftID) transport.CarID {
	ids := w.bank.CarsIn(shaft)
	if len(ids) == 0 {
		return 0
	}
	return ids[0]
}
