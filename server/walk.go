package main

import (
	"math"
	"math/rand"

	"towergame/server/model"
	"towergame/server/transport"
)

// Walking speeds, in cells per second.
const (
	walkSpeed   = 6.0
	stairSpeed  = 2.5 // slower, one step at a time
	browseSpeed = 1.1 // on a shop floor
)

// Browsing a shop floor.
const (
	browseMin   = 1.5 // shortest walk to a new spot
	browsePause = 3.0 // longest look at the goods
)

// Replans a path when the goal changes.
func (w *World) applyGoal(p *sim, g goal) {
	if g == p.last {
		return
	}
	// Finish the ride, or they exit mid-shaft.
	if p.holdsCarPlace() {
		return
	}
	p.last = g
	p.idx = 0
	p.shaft = 0 // any queued call is abandoned
	// They are up and about again.
	p.dark, p.asleep = false, false

	fromF, fromC := int(math.Round(p.y)), int(math.Round(p.x))
	// Remember the start, so the inspector can show it.
	p.from = goal{true, fromF, fromC}
	if p.state == stateOutside {
		// They walk in off the street, not out of the lobby.
		fromF, fromC = 0, w.streetCol(g.col)
		p.x, p.y = float64(fromC), 0
		p.from = goal{present: false}
	}

	toF, toC := g.floor, g.col
	if !g.present {
		toF, toC = 0, w.streetCol(fromC)
	}

	if path, ok := w.grid.PathVia(fromF, fromC, toF, toC, w.shaftCost); ok {
		p.path = path
		p.state = stateMoving
	} else if !g.present {
		p.state = stateOutside
	}
}

// Advances a sim: walking, waiting, or riding.
func (w *World) move(p *sim, dt float64) {
	switch p.state {
	case statePresent:
		w.browse(p, dt)
		return
	case stateWaiting:
		w.waitForCar(p, dt)
		return
	case stateBoarding:
		w.walkIntoCar(p, dt)
		return
	case stateRiding:
		w.ride(p)
		return
	case stateMoving:
		// handled below
	default:
		return
	}

	if p.idx >= len(p.path) {
		if p.last.present {
			p.state = statePresent
			p.x, p.y = float64(p.last.col), float64(p.last.floor)
			w.markLight(p) // dark rooms hide sleepers at once
			// The visit starts on arrival, not on setting out.
			if p.transient && p.leave == 0 {
				p.leave = w.simTime + visitStay*(1+stayVariance*(rand.Float64()*2-1))
			}
		} else {
			p.state = stateOutside
		}
		return
	}

	target := p.path[p.idx]

	// Board an elevator rather than float up the shaft.
	if target.Mode == model.ModeElevator {
		p.shaft = transport.ShaftID(target.Shaft)
		p.boardFloor = int(math.Round(p.y))
		p.exitFloor = int(target.Floor)
		// The controller picks the car.
		car, _, ok := w.bank.BestCar(p.shaft, p.boardFloor, p.exitFloor)
		if !ok {
			p.shaft = 0
			return // no car serves this shaft yet
		}
		p.car = car
		p.waitTime = 0
		p.queueX = p.x // layoutQueues assigns the real spot
		p.y = float64(p.boardFloor)
		p.state = stateWaiting
		return
	}

	dx := target.Col - p.x
	dy := target.Floor - p.y
	dist := math.Hypot(dx, dy)
	// Climbing is slower than crossing a floor.
	speed := walkSpeed
	if target.Mode == model.ModeStair {
		speed = stairSpeed
	}
	step := speed * dt
	if dist <= step || dist == 0 {
		p.x, p.y = target.Col, target.Floor
		p.idx++
		return
	}
	p.x += dx / dist * step
	p.y += dy / dist * step
}

// Walks a shopper up and down the shop floor.
// People in other rooms stay at their seat.
func (w *World) browse(p *sim, dt float64) {
	// Only shoppers browse; skip the room lookup.
	if p.category != model.CategoryRetail {
		return
	}
	r, ok := w.roomOf(p)
	if !ok || model.RoomTypes[r.Type].Habit != model.HabitBrowse {
		return
	}
	if p.pause > 0 {
		p.pause -= dt
		return
	}
	if p.browseX != 0 && !stepToward(&p.x, p.browseX, browseSpeed*dt) {
		return
	}
	p.browseX = w.browseSpot(r, p.x)
	p.pause = rand.Float64() * browsePause
}

// A new spot to walk to, well away from here.
func (w *World) browseSpot(r model.Room, from float64) float64 {
	rt := model.RoomTypes[r.Type]
	left := float64(r.Col) + 0.5
	right := float64(r.Col+rt.Width) - 1.5
	if right-left < browseMin {
		return (left + right) / 2
	}
	for i := 0; i < 8; i++ {
		spot := left + rand.Float64()*(right-left)
		if math.Abs(spot-from) >= browseMin {
			return spot
		}
	}
	// Give up and cross the room instead.
	if from-left > right-from {
		return left
	}
	return right
}

// Moves v toward target. Reports arrival.
func stepToward(v *float64, target, step float64) bool {
	d := target - *v
	if math.Abs(d) <= step {
		*v = target
		return true
	}
	*v += math.Copysign(step, d)
	return false
}

// The built span of the ground floor, or -1.
func (w *World) groundSpan() (left, right int) {
	left, right = -1, -1
	for c := 0; c < w.grid.Width; c++ {
		if w.grid.Built[0][c] {
			if left < 0 {
				left = c
			}
			right = c
		}
	}
	return left, right
}

// The pavement outside the nearest end of the tower.
// People appear and vanish here, never inside the lobby.
func (w *World) streetCol(col int) int {
	left, right := w.groundSpan()
	if left < 0 {
		return clampCol(w, col)
	}
	if col-left <= right-col {
		return left - 1
	}
	return right + 1
}

func clampCol(w *World, col int) int {
	if col < 0 {
		return 0
	}
	if col >= w.grid.Width {
		return w.grid.Width - 1
	}
	return col
}
