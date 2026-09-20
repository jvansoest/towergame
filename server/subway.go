package main

import (
	"math/rand"

	"towergame/server/model"
)

// Trains bring visitors to the shops underground.
const (
	trainEvery  = 20.0 // seconds between trains, roughly
	trainFrom   = 6 * 60
	trainTo     = 23 * 60
	trainSpeed  = 8.0 // cells per second
	trainDwell  = 5.0 // seconds at the platform
	trainPeople = 10  // most passengers a train unloads
	trainReach  = 14  // cells beyond the grid, out of sight
)

const (
	trainArriving = iota + 1
	trainDwelling
	trainLeaving
)

// The train in the tunnel, if one is running.
type train struct {
	x     float64 // centre column
	phase int
	wait  float64
}

// The subway station room, if built.
func (w *World) station() (model.Room, bool) {
	for _, r := range w.grid.Rooms {
		if r.Type == "subway" {
			return r, true
		}
	}
	return model.Room{}, false
}

// Where passengers step on and off: the platform.
func (w *World) stationSpot() (floor, col int) {
	r, _ := w.station()
	return r.Floor, r.Col + model.RoomTypes[r.Type].Width/2
}

// Runs the timetable.
func (w *World) stepTrains(dt float64) {
	st, ok := w.station()
	if !ok {
		w.train = nil
		return
	}
	clk := clockFromSimTime(w.simTime)
	mod := clk.Hour*60 + clk.Minute
	stop := float64(st.Col + model.RoomTypes[st.Type].Width/2)
	t := w.train
	if t == nil {
		if w.simTime < w.nextTrain || mod < trainFrom || mod >= trainTo {
			return
		}
		w.train = &train{x: -trainReach, phase: trainArriving}
		return
	}
	switch t.phase {
	case trainArriving:
		t.x += trainSpeed * dt
		if t.x >= stop {
			t.x, t.phase, t.wait = stop, trainDwelling, trainDwell
			w.unloadTrain(mod)
		}
	case trainDwelling:
		if t.wait -= dt; t.wait <= 0 {
			t.phase = trainLeaving
		}
	case trainLeaving:
		t.x += trainSpeed * dt
		if t.x > float64(w.grid.Width+trainReach) {
			w.train = nil
			w.nextTrain = w.simTime + trainEvery*(0.7+0.6*rand.Float64())
		}
	}
}

// Sends passengers to open shops below ground. They shop
// and come back to the station, and never go up to the street.
func (w *World) unloadTrain(mod int) {
	floor, col := w.stationSpot()
	var open []model.Room
	for _, r := range w.grid.Rooms {
		rt := model.RoomTypes[r.Type]
		if r.Floor >= 0 || rt.Sale == 0 || !isOpen(r, mod) {
			continue
		}
		if _, ok := w.grid.PathVia(floor, col, r.Floor, model.SeatCol(r, 0), nil); ok {
			open = append(open, r)
		}
	}
	for i := 0; i < trainPeople && len(open) > 0; i++ {
		r := open[rand.Intn(len(open))]
		taken, _ := w.spotsTaken(r)
		rt := model.RoomTypes[r.Type]
		want := visitorsPer
		if rt.Habit == model.HabitSit {
			want = rt.Seats
		}
		spot := -1
		for s := 0; s < want; s++ {
			if !taken[s] {
				spot = s
				break
			}
		}
		if spot < 0 {
			continue
		}
		w.addVisitor(r, spot)
		p := w.sims[len(w.sims)-1]
		p.subway = true
		p.pace = 2 // a long way to walk
	}
}
