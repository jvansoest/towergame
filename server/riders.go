package main

import (
	"math"
	"sort"

	"towergame/server/model"
	"towergame/server/transport"
)

// Queueing for a car, and standing in one.
const (
	rethinkEvery = 15   // ticks between dispatch reviews
	rethinkGain  = 3.0  // cost a better car must beat
	riderSpacing = 0.45 // cells between riders in a car
	queueGap     = 1.5  // cells the queue head keeps off the shaft
	queueSpacing = 0.34 // cells between queued sims
)

// Longest wait first, id breaks ties.
// Reshuffled ties make sims swap places.
func sortQueue(ps []*sim) {
	sort.Slice(ps, func(a, b int) bool {
		if ps[a].waitTime != ps[b].waitTime {
			return ps[a].waitTime > ps[b].waitTime
		}
		return ps[a].id < ps[b].id
	})
}

// Lines up waiting sims beside each shaft, longest wait first.
func (w *World) layoutQueues() {
	type key struct {
		shaft transport.ShaftID
		floor int
	}
	groups := map[key][]*sim{}
	for _, p := range w.sims {
		if p.state == stateWaiting {
			k := key{p.shaft, p.boardFloor}
			groups[k] = append(groups[k], p)
		}
	}
	for k, ps := range groups {
		sortQueue(ps)
		col := w.shaftCol(k.shaft)
		for i, p := range ps {
			p.queueX = col - queueGap - float64(i)*queueSpacing
			p.y = float64(k.floor)
		}
	}
}

// Centre column of a shaft.
func (w *World) shaftCol(shaft transport.ShaftID) float64 {
	e, ok := w.grid.Shaft(model.ShaftID(shaft))
	if !ok {
		return 0
	}
	return e.CenterCol()
}

// Column of a fixed standing slot in a car.
func (w *World) slotCol(shaft transport.ShaftID, slot int) float64 {
	mid := float64(transport.Capacity-1) / 2
	return w.shaftCol(shaft) + (float64(slot)-mid)*riderSpacing
}

// Holds riders in their slot as the car travels.
func (w *World) layoutRiders() {
	for _, p := range w.sims {
		car, ok := w.bank.Car(p.car)
		if !ok {
			continue
		}
		switch p.state {
		case stateBoarding:
			p.boardX = w.slotCol(p.shaft, p.slot) // walked to, not jumped to
		case stateRiding:
			// Slots are fixed, so nobody shifts mid-ride.
			p.x = w.slotCol(p.shaft, p.slot)
			p.y = car.Floor
		}
	}
}

// Waits at the shaft, shuffling up as the queue moves.
func (w *World) waitForCar(p *sim, dt float64) {
	p.waitTime += dt
	if p.waitTime > stressGrace {
		p.stress = math.Min(stressMax, p.stress+stressRise*dt)
	}
	if _, ok := w.bank.Car(p.car); !ok {
		p.state = stateMoving // elevator gone
		return
	}
	w.rethinkCar(p)
	stepToward(&p.x, p.queueX, walkSpeed*dt)
}

// Moves a waiting sim to a clearly better car.
// A full or departing car must not hold them.
func (w *World) rethinkCar(p *sim) {
	if w.tick%rethinkEvery != 0 {
		return
	}
	best, cost, ok := w.bank.BestCar(p.shaft, p.boardFloor, p.exitFloor)
	if !ok || best == p.car {
		return
	}
	// Only a real gain, or the call bounces between cars.
	if cost+rethinkGain < w.bank.CarCost(p.car, p.boardFloor, p.exitFloor) {
		p.car = best
	}
}

// Boards waiting sims front of queue first.
func (w *World) boardQueues() {
	type key struct {
		shaft transport.ShaftID
		floor int
	}
	groups := map[key][]*sim{}
	for _, p := range w.sims {
		if p.state == stateWaiting {
			k := key{p.shaft, p.boardFloor}
			groups[k] = append(groups[k], p)
		}
	}
	for k, ps := range groups {
		sortQueue(ps) // same order the queue is drawn in
		for _, p := range ps {
			// Whatever opens its doors going your way.
			car, ok := w.openCar(p, k.floor)
			if !ok {
				continue
			}
			taken := w.takenSlots(car)
			// Fill from the far side, avoiding crossings.
			slot := -1
			for s := transport.Capacity - 1; s >= 0; s-- {
				if !taken[s] {
					slot = s
					break
				}
			}
			// Board enforces capacity, so may refuse.
			if slot < 0 || !w.bank.Board(car) {
				continue
			}
			p.car = car
			p.slot = slot
			p.boardX = w.slotCol(p.shaft, slot)
			p.state = stateBoarding
		}
	}
}

// A car with room, stopped here, going their way.
// Their own car first, then any other in the shaft.
func (w *World) openCar(p *sim, floor int) (transport.CarID, bool) {
	dir := callDir(p)
	if w.carTakes(p.car, floor, dir) {
		return p.car, true
	}
	for _, id := range w.bank.CarsIn(p.shaft) {
		if w.carTakes(id, floor, dir) {
			return id, true
		}
	}
	return 0, false
}

// Whether a car can take a rider here.
func (w *World) carTakes(car transport.CarID, floor, dir int) bool {
	return w.bank.StoppedAt(car, float64(floor)) &&
		w.bank.Serves(car, dir) &&
		w.bank.HasRoom(car)
}

// Whether a sim occupies a place in a car.
func (p *sim) holdsCarPlace() bool {
	if p.shaft == 0 {
		return false
	}
	return p.state == stateRiding || p.state == stateBoarding
}

// Restates occupancy from the sims.
// Derived, so no leak can strand a car.
func (w *World) recountRiders() {
	w.bank.SyncRiders(w.riderCounts())
}

// Slots already claimed in a car.
func (w *World) takenSlots(car transport.CarID) map[int]bool {
	taken := map[int]bool{}
	for _, p := range w.sims {
		if p.car == car && p.holdsCarPlace() {
			taken[p.slot] = true
		}
	}
	return taken
}

// Walks a boarder from the queue into the car.
func (w *World) walkIntoCar(p *sim, dt float64) {
	car, ok := w.bank.Car(p.car)
	if !ok {
		p.state = stateMoving
		return
	}
	// Hold the doors until everyone is inside.
	w.bank.HoldDoors(p.car)
	p.y = car.Floor

	if stepToward(&p.x, p.boardX, walkSpeed*dt) {
		p.state = stateRiding
	}
}

// Rides the car and alights at the exit floor.
func (w *World) ride(p *sim) {
	car, ok := w.bank.Car(p.car)
	if !ok {
		p.state = stateMoving
		return
	}
	// Step out only when level and stopped.
	if math.Abs(car.Floor-float64(p.exitFloor)) > transport.LevelEps {
		return
	}
	w.bank.HoldDoors(p.car)
	p.shaft = 0
	// Keep the column; the path walks them out.
	p.y = float64(p.exitFloor)
	p.idx++
	p.state = stateMoving
}
