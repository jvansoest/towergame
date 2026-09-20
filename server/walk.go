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
	glideSpeed  = 3.0 // an escalator carries them
	browseSpeed = 1.1 // on a shop floor
)

// Browsing a shop floor.
const (
	browseMin   = 1.5 // shortest walk to a new spot
	browsePause = 3.0 // longest look at the goods
)

// Sitters sometimes rise for a slow stroll.
const (
	wanderSpeed = 0.8 // cells per second, unhurried
	wanderRest  = 6.0 // least rest between strolls
)

// Whether a sim is riding an escalator now.
func gliding(p *sim) bool {
	return p.state == stateMoving && p.idx < len(p.path) &&
		p.path[p.idx].Mode == model.ModeEscalator
}

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
		if p.subway {
			// They step off the train.
			fromF, fromC = w.stationSpot()
		} else {
			// They walk in off the street, not out of the lobby.
			fromF, fromC = 0, w.streetCol(g.col)
		}
		p.x, p.y = float64(fromC), float64(fromF)
		p.from = goal{present: false}
	}

	toF, toC := g.floor, g.col
	if !g.present {
		toF, toC = 0, w.streetCol(fromC)
		if p.subway {
			toF, toC = w.stationSpot()
		}
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
				p.stayed = true
				if p.category == model.CategoryHotel {
					w.markUsed(p)
					p.leave = w.simTime + w.stayLength(p)
				} else {
					w.sell(p)
					p.leave = w.simTime + visitStay*p.pace*(1+stayVariance*(rand.Float64()*2-1))
				}
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
	if target.Mode == model.ModeEscalator {
		speed = glideSpeed
	} else if target.Mode == model.ModeStair {
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

// Paces a room. Shoppers roam all visit long;
// other sitters rise now and then for a slow
// stroll. Desk work keeps workers on chairs.
func (w *World) browse(p *sim, dt float64) {
	// A maid at work keeps to one spot.
	if p.prof == profMaid && p.job.present {
		return
	}
	// A desk job chains its worker to the chair.
	if p.category == model.CategoryOffice {
		return
	}
	r, ok := w.roomOf(p)
	if !ok {
		return
	}
	rt := model.RoomTypes[r.Type]
	// Bar guests stay on their stools.
	if rt.Line {
		return
	}
	habit := rt.Habit
	if habit != model.HabitBrowse && habit != model.HabitSit {
		return
	}
	speed, rest := wanderSpeed, wanderRest
	if habit == model.HabitBrowse { // shoppers keep at it
		speed, rest = browseSpeed, 0
	}
	if p.pause > 0 {
		p.pause -= dt
		return
	}
	if p.browseX != 0 && !stepToward(&p.x, p.browseX, speed*dt) {
		return
	}
	p.browseX = w.browseSpot(r, p.x)
	p.pause = rand.Float64()*browsePause + rest
}

// A new spot to walk to, well away from here.
func (w *World) browseSpot(r model.Room, from float64) float64 {
	rt := model.RoomTypes[r.Type]
	left := float64(r.Col) + 0.5
	right := float64(r.Col+rt.Width) - 1.5
	// Small rooms allow shorter walks.
	least := math.Min(browseMin, (right-left)*0.6)
	if right <= left {
		return (left + right) / 2
	}
	for i := 0; i < 8; i++ {
		spot := left + rand.Float64()*(right-left)
		if math.Abs(spot-from) >= least {
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
