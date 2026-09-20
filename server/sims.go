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

// What a sim is, for the record and the client.
type Profession string

const (
	profVisitor  Profession = "visitor" // drinks, shops, eats, sightsees
	profResident Profession = "resident"
	profWorker   Profession = "worker"
	profSecurity Profession = "security"
	profTriad    Profession = "triad"
	profMaid     Profession = "maid"
	profDoctor   Profession = "doctor"
	profVIP      Profession = "vip"
)

// A sim: a resident, worker, guest, or customer.
type sim struct {
	id        int
	category  model.Category // which room kind anchors them
	prof      Profession
	homeF     int
	homeC     int
	transient bool    // customers leave and are removed
	spot      int     // seat or floor place a customer took
	pace      float64 // stretches a visitor's patience and stay

	// personal timekeeping, in minutes
	shiftIn  int // early or late to arrive
	shiftOut int // early or late to leave
	lunchAt  int // early or late to eat

	dark   bool    // the room around them is unlit
	asleep bool    // in bed, so not drawn
	stress float64 // 0..100, from waiting for elevators
	peak   float64 // worst stress this quarter

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

	// hotel guests and maids
	subway    bool    // came by train, and leaves by it
	stayed    bool    // a guest reached their room
	disturbed bool    // noise drove a guest out
	job       goal    // the dirty room a maid heads for
	cleanLeft float64 // seconds of scrubbing left

	// browsing a shop
	browseX float64 // spot being walked to
	pause   float64 // seconds left looking at it
}

// Adds a room's residents, per its capacity.
func (w *World) addSims(r model.Room) {
	rt := model.RoomTypes[r.Type]
	prof := profVisitor // hotel guests are visiting, loosely
	switch {
	case r.Type == "security":
		prof = profSecurity
	case rt.Category == model.CategoryOffice:
		prof = profWorker
	case rt.Category == model.CategoryResidential:
		prof = profResident
	case rt.Category == model.CategoryService:
		prof = profMaid
	case rt.Category == model.CategoryMedical:
		prof = profDoctor
	case rt.Category == model.CategoryHotel:
		return // guests rent rooms, they do not live here
	}
	for i := 0; i < rt.Capacity; i++ {
		w.nextID++
		w.sims = append(w.sims, &sim{
			id:       w.nextID,
			category: rt.Category,
			prof:     prof,
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
	w.manageGuests(mod)
	w.stepMaids(dt)
	w.refreshNoise(mod)

	kept := w.sims[:0]
	for _, p := range w.sims {
		g := w.goalFor(p, weekday, mod)
		w.applyGoal(p, g)
		w.move(p, dt)
		// Stress fades everywhere but the elevator queue.
		if p.state != stateWaiting {
			p.stress = math.Max(0, p.stress-stressFall*dt)
		}
		w.noiseStress(p, dt)
		p.peak = math.Max(p.peak, p.stress)
		if p.transient && p.state == stateOutside && p.last == (goal{false, p.homeF, p.homeC}) {
			w.checkOut(p)
			continue // customer has left; drop it
		}
		kept = append(kept, p)
	}
	w.sims = kept
	w.releaseGuests()
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
