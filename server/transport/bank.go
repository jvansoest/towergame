// Package transport moves elevator cars.
// Cars and shafts are named by id.
package transport

import "math"

// Car behaviour, in floors and seconds.
const (
	Capacity = 4    // riders per car
	Speed    = 8.0  // floors per second
	DwellFor = 1.2  // seconds held with doors open
	LevelEps = 1e-6 // a car this close counts as level
)

// Cars one shaft may hold. They share the
// shaft freely: they pass and stack.
const MaxPerShaft = 8

// A shaft's stable name.
type ShaftID int

// A car's stable name.
type CarID int

// A stop a car owes.
type Call struct {
	Floor float64
	Dir   int // -1 or 1 for a hall call, 0 for a drop-off
}

// One car in a shaft.
type Car struct {
	ID    CarID
	Shaft ShaftID
	Floor float64 // where it is now
	To    float64 // where it is heading
	Dir   int     // -1, 0 or 1

	riders int
	dwell  float64
}

// Seconds left with the doors open.
func (c Car) Dwell() float64 { return c.dwell }

// How many riders are aboard.
func (c Car) Riders() int { return c.riders }

// The cars of one building.
type Bank struct {
	Floors int // building height, for the suitability score

	cars  []Car
	lastI CarID
}

// Adds a car to a shaft.
func (b *Bank) Add(shaft ShaftID) CarID {
	b.lastI++
	b.cars = append(b.cars, Car{ID: b.lastI, Shaft: shaft})
	return b.lastI
}

// Removes every car serving a shaft.
func (b *Bank) RemoveShaft(shaft ShaftID) {
	kept := b.cars[:0]
	for _, c := range b.cars {
		if c.Shaft != shaft {
			kept = append(kept, c)
		}
	}
	b.cars = kept
}

func (b *Bank) Len() int { return len(b.cars) }

// Every car, copied so callers cannot steer.
func (b *Bank) Cars() []Car {
	return append([]Car(nil), b.cars...)
}

// Puts a car at a floor, stopped.
func (b *Bank) Place(id CarID, floor float64) {
	if i := b.at(id); i >= 0 {
		b.cars[i].Floor, b.cars[i].To, b.cars[i].Dir = floor, floor, 0
	}
}

// The car with this id.
func (b *Bank) Car(id CarID) (Car, bool) {
	if i := b.at(id); i >= 0 {
		return b.cars[i], true
	}
	return Car{}, false
}

// The cars serving a shaft.
func (b *Bank) CarsIn(shaft ShaftID) []CarID {
	var ids []CarID
	for _, c := range b.cars {
		if c.Shaft == shaft {
			ids = append(ids, c.ID)
		}
	}
	return ids
}

// Slot of a car id, or -1.
func (b *Bank) at(id CarID) int {
	for i := range b.cars {
		if b.cars[i].ID == id {
			return i
		}
	}
	return -1
}

// Holds a car while riders move.
func (b *Bank) HoldDoors(id CarID) {
	if i := b.at(id); i >= 0 {
		b.cars[i].dwell = DwellFor
	}
}

// Whether a car waits at this floor.
func (b *Bank) StoppedAt(id CarID, floor float64) bool {
	i := b.at(id)
	if i < 0 {
		return false
	}
	c := b.cars[i]
	// Level, and meaning to stop here.
	return math.Abs(c.Floor-floor) <= LevelEps &&
		math.Abs(c.To-floor) <= LevelEps
}

// Whether a car goes this way.
func (b *Bank) Serves(id CarID, dir int) bool {
	i := b.at(id)
	return i >= 0 && (b.cars[i].Dir == 0 || b.cars[i].Dir == dir)
}

// Whether a car has spare room.
func (b *Bank) HasRoom(id CarID) bool {
	i := b.at(id)
	return i >= 0 && b.cars[i].riders < Capacity
}

// Takes one passenger aboard. Enforces capacity.
func (b *Bank) Board(id CarID) bool {
	i := b.at(id)
	if i < 0 || b.cars[i].riders >= Capacity {
		return false
	}
	b.cars[i].riders++
	b.cars[i].dwell = DwellFor
	return true
}

// Restates occupancy; never adjusts it.
// A vanished rider cannot strand a car.
func (b *Bank) SyncRiders(aboard map[CarID]int) {
	for i := range b.cars {
		b.cars[i].riders = aboard[b.cars[i].ID]
	}
}

// Runs every car for one timestep.
// Collective control: sweep one way, then reverse.
// Calls belong to a car, not to its shaft.
func (b *Bank) Step(dt float64, calls map[CarID][]Call) {
	for i := range b.cars {
		c := &b.cars[i]

		// Hold still while riders move.
		if c.dwell > 0 {
			c.dwell -= dt
			continue
		}
		c.To = c.pickStop(calls[c.ID])

		diff := c.To - c.Floor
		if step := Speed * dt; math.Abs(diff) <= step {
			c.Floor = c.To
		} else {
			c.Floor += math.Copysign(step, diff)
		}
	}
}

// The car's next stop and direction.
func (c *Car) pickStop(calls []Call) float64 {
	if len(calls) == 0 {
		c.Dir = 0
		return c.Floor
	}

	dir := c.Dir
	if dir == 0 {
		dir = c.dirOfNearest(calls)
	}
	stop, ok := c.nextStop(calls, dir)
	if !ok {
		dir = -dir // nothing left this way, turn round
		stop, ok = c.nextStop(calls, dir)
	}
	if !ok {
		// All calls oppose us: fetch nearest.
		stop, dir, ok = c.fetchStop(calls)
	}
	if !ok {
		c.Dir = 0
		return c.Floor
	}
	c.Dir = dir
	return stop
}

// Nearest stop ahead, our way only.
func (c *Car) nextStop(calls []Call, dir int) (float64, bool) {
	stop, bestD, found := 0.0, math.MaxFloat64, false
	for _, call := range calls {
		ahead := (call.Floor - c.Floor) * float64(dir)
		if ahead < -LevelEps {
			continue // already passed it
		}
		if call.Dir != 0 && call.Dir != dir {
			continue // they want the other way
		}
		if ahead < bestD {
			stop, bestD, found = call.Floor, ahead, true
		}
	}
	return stop, found
}

// Nearest call when the sweep is done.
// The car fetches it, then follows them.
func (c *Car) fetchStop(calls []Call) (float64, int, bool) {
	best, bestD, found := Call{}, math.MaxFloat64, false
	for _, call := range calls {
		if d := math.Abs(call.Floor - c.Floor); d < bestD {
			best, bestD, found = call, d, true
		}
	}
	if !found {
		return 0, 0, false
	}
	// Take their direction on arrival.
	dir := best.Dir
	if bestD > LevelEps {
		dir = 1
		if best.Floor < c.Floor {
			dir = -1
		}
	}
	if dir == 0 {
		dir = 1
	}
	return best.Floor, dir, true
}

// Which way an idle car starts.
func (c *Car) dirOfNearest(calls []Call) int {
	dir, bestD := 1, math.MaxFloat64
	for _, call := range calls {
		if d := math.Abs(call.Floor - c.Floor); d < bestD {
			bestD = d
			dir = 1
			if call.Floor < c.Floor {
				dir = -1
			}
		}
	}
	return dir
}
