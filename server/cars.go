package main

import (
	"math"
	"math/rand"

	"towergame/server/model"
)

// Cars drive down a ramp, park, and drive out.
const (
	carSpeed  = 5.0 // cells per second along the road
	carsMax   = 8   // in the garage at once
	carChance = 0.06
	carFrom   = 7 * 60
	carTo     = 20 * 60
	parkFee   = 60
	parkStay  = 15.0 // least seconds parked
	parkSpan  = 30.0 // extra seconds, at random
	floorCost = 2.0  // a floor of travel, in cells
)

type vec struct{ x, y float64 }

const (
	carArriving = iota
	carParked
	carLeaving
)

// One car: its route, and the slot it holds.
type vehicle struct {
	id      int
	x, y    float64
	dir     int     // -1 faces left, 1 right
	slope   float64 // floors per cell, signed
	path    []vec
	idx     int
	phase   int
	floor   int // parking room row
	col     int // parking room column
	slot    int
	leaveAt float64
}

// Route from the street to a slot, or false.
func (w *World) routeTo(pr model.Room, slot int) ([]vec, bool) {
	rt := model.RoomTypes[pr.Type]
	for _, first := range w.grid.Ramps {
		if first.Floor != -1 {
			continue
		}
		cc := first.Col
		pts := []vec{{float64(first.TopCol()), 0}}
		var last model.Ramp
		ok := true
		for f := -1; f >= pr.Floor; f-- {
			r, found := w.grid.RampAt(f, cc)
			if !found || r.Col != cc {
				ok = false
				break
			}
			pts = append(pts, vec{float64(r.BottomCol()), float64(f)})
			last = r
		}
		if !ok {
			continue
		}
		// The room lies beyond the run's low end.
		if last.DownRight() && pr.Col < cc+model.RampWidth ||
			!last.DownRight() && pr.Col+rt.Width > cc {
			continue
		}
		slotX := float64(pr.Col + slot)
		lo, hi := int(math.Min(float64(last.BottomCol()), slotX)),
			int(math.Max(float64(last.BottomCol()), slotX))
		for c := lo; c <= hi; c++ {
			if !w.grid.IsBuilt(pr.Floor, c) {
				ok = false
			}
		}
		if !ok {
			continue
		}
		return append(pts, vec{slotX, float64(pr.Floor)}), true
	}
	return nil, false
}

// Whether a slot has a car reserved or parked in it.
func (w *World) slotTaken(pr model.Room, slot int) bool {
	for _, v := range w.vehicles {
		if v.floor == pr.Floor && v.col == pr.Col && v.slot == slot &&
			v.phase != carLeaving {
			return true
		}
	}
	return false
}

// Sends a car to a free, reachable slot now and then.
func (w *World) manageCars(mod int) {
	if w.tick%visitorEvery != 0 || len(w.vehicles) >= carsMax ||
		mod < carFrom || mod >= carTo || rand.Float64() > carChance {
		return
	}
	for _, j := range rand.Perm(len(w.grid.Rooms)) {
		r := w.grid.Rooms[j]
		rt := model.RoomTypes[r.Type]
		if rt.Category != model.CategoryParking {
			continue
		}
		for s := 0; s < rt.Slots; s++ {
			if w.slotTaken(r, s) {
				continue
			}
			path, ok := w.routeTo(r, s)
			if !ok {
				continue
			}
			w.nextID++
			w.vehicles = append(w.vehicles, &vehicle{
				id: w.nextID, x: path[0].x, y: path[0].y, dir: 1,
				path: path, idx: 1, phase: carArriving,
				floor: r.Floor, col: r.Col, slot: s,
			})
			return
		}
	}
}

// Drives every car along its route.
func (w *World) stepCars(dt float64) {
	clk := clockFromSimTime(w.simTime)
	w.manageCars(clk.Hour*60 + clk.Minute)
	kept := w.vehicles[:0]
	for _, v := range w.vehicles {
		if i := w.grid.RoomAt(v.floor, v.col); i < 0 ||
			w.grid.Rooms[i].Type != "parking" {
			continue // the garage is gone
		}
		if v.phase == carParked {
			if w.simTime >= v.leaveAt && w.leave(v) {
				v.phase = carLeaving
			}
			kept = append(kept, v)
			continue
		}
		if w.drive(v, dt) && v.phase == carLeaving {
			w.money += parkFee
			w.books.sales += parkFee
			continue // out of the garage
		}
		kept = append(kept, v)
	}
	w.vehicles = kept
}

// Plans the way out: the way in, reversed.
func (w *World) leave(v *vehicle) bool {
	i := w.grid.RoomAt(v.floor, v.col)
	path, ok := w.routeTo(w.grid.Rooms[i], v.slot)
	if !ok {
		return false
	}
	for a, b := 0, len(path)-1; a < b; a, b = a+1, b-1 {
		path[a], path[b] = path[b], path[a]
	}
	v.path, v.idx = path, 1
	return true
}

// Moves along the route. Reports the end, and parks there.
func (w *World) drive(v *vehicle, dt float64) bool {
	step := carSpeed * dt
	for step > 0 && v.idx < len(v.path) {
		to := v.path[v.idx]
		dx, dy := to.x-v.x, (to.y-v.y)*floorCost
		dist := math.Hypot(dx, dy)
		if dist <= step {
			v.x, v.y = to.x, to.y
			step -= dist
			v.idx++
			continue
		}
		v.x += (to.x - v.x) * step / dist
		v.y += (to.y - v.y) * step / dist
		if dx != 0 {
			v.dir = 1
			if dx < 0 {
				v.dir = -1
			}
			v.slope = (to.y - v.y) / (to.x - v.x)
		} else {
			v.slope = 0
		}
		step = 0
	}
	if v.idx < len(v.path) {
		return false
	}
	if v.phase == carArriving {
		v.phase = carParked
		v.slope = 0
		v.leaveAt = w.simTime + parkStay + rand.Float64()*parkSpan
		return false
	}
	return true
}
