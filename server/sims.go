package main

import (
	"math"

	"towergame/server/model"
	"towergame/server/transport"
)

// Stress, in points per real second.
const (
	stressMax   = 100.0
	stressRise  = 2.0 // while waiting for an elevator
	stressFall  = 0.4 // anywhere else
	stressGrace = 3.0 // seconds of waiting nobody minds
)

type simState int

const (
	stateOutside simState = iota
	stateMoving
	statePresent
	stateWaiting  // queued for an elevator
	stateBoarding // walking into a waiting car
	stateRiding   // inside an elevator car
)

// A desired location, or outside.
type goal struct {
	present bool
	floor   int
	col     int
}

// A sim: a resident, worker, guest, or customer.
type sim struct {
	id        int
	category  model.Category
	homeF     int
	homeC     int
	transient bool // customers leave and are removed
	spot      int  // seat or floor place a customer took

	// personal timekeeping, in minutes
	shiftIn  int // early or late to arrive
	shiftOut int // early or late to leave
	lunchAt  int // early or late to eat

	dark   bool    // the room around them is unlit
	asleep bool    // in bed, so not drawn
	stress float64 // 0..100, from waiting for elevators

	x, y   float64
	path   []model.Waypoint
	idx    int
	state  simState
	last   goal    // where they are headed
	from   goal    // where this trip started
	bornAt float64 // sim-time a customer set out
	leave  float64 // sim-time a customer departs, 0 until it arrives

	// elevator ride
	shaft      transport.ShaftID // 0 when not using an elevator
	car        transport.CarID   // the car assigned to them
	boardFloor int
	exitFloor  int
	queueX     float64 // column to walk to while queueing
	boardX     float64 // column to walk to when boarding
	slot       int     // standing place inside the car
	waitTime   float64

	// browsing a shop
	browseX float64 // spot being walked to
	pause   float64 // seconds left looking at it
}

// Adds a room's residents, per its capacity.
func (w *World) addSims(r model.Room) {
	rt := model.RoomTypes[r.Type]
	for i := 0; i < rt.Capacity; i++ {
		w.nextID++
		w.sims = append(w.sims, &sim{
			id:       w.nextID,
			category: rt.Category,
			homeF:    r.Floor,
			homeC:    model.SeatCol(r, i),
			state:    stateOutside,
			shiftIn:  spreadMinutes(scheduleSpread),
			shiftOut: spreadMinutes(scheduleSpread),
			lunchAt:  spreadMinutes(lunchSpread),
		})
	}
}

// Schedules and moves every sim.
func (w *World) stepSims(dt float64) {
	clk := clockFromSimTime(w.simTime)
	weekday := isWeekday(clk.Day)
	mod := clk.Hour*60 + clk.Minute

	w.manageVisitors(mod)

	kept := w.sims[:0]
	for _, p := range w.sims {
		g := w.goalFor(p, weekday, mod)
		w.applyGoal(p, g)
		w.move(p, dt)
		// Stress fades everywhere but the elevator queue.
		if p.state != stateWaiting {
			p.stress = math.Max(0, p.stress-stressFall*dt)
		}
		if p.transient && p.state == stateOutside && p.last == (goal{false, p.homeF, p.homeC}) {
			continue // customer has left; drop it
		}
		kept = append(kept, p)
	}
	w.sims = kept
	w.logDropOffs()
	w.boardQueues()
	w.layoutQueues()
	w.layoutRiders()
}

// The room a sim stands in, if any.
func (w *World) roomOf(p *sim) (model.Room, bool) {
	i := w.grid.RoomAt(int(math.Round(p.y)), int(math.Round(p.x)))
	if i < 0 {
		return model.Room{}, false
	}
	return w.grid.Rooms[i], true
}
