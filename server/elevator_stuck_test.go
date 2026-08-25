package main

import (
	"math"
	"testing"

	"towergame/server/transport"
)

// A car must never park while someone still needs it.
func TestCarNeverStalls(t *testing.T) {
	w := newWorld()
	w.seed()

	const dt = 1.0 / 30
	const stallTicks = 30 * 90 // 90 seconds without moving

	prev := map[transport.CarID]float64{}
	lastMove := map[transport.CarID]int{}
	for _, c := range w.bank.Cars() {
		prev[c.ID] = c.Floor
	}

	for step := 0; step < 30*60*90; step++ { // 90 simulated minutes
		w.Step(dt)
		for _, c := range w.bank.Cars() {
			if math.Abs(c.Floor-prev[c.ID]) > transport.LevelEps {
				prev[c.ID], lastMove[c.ID] = c.Floor, step
				continue
			}
			if step-lastMove[c.ID] < stallTicks {
				continue
			}
			if !w.carHasWork(c.ID) {
				lastMove[c.ID] = step // idle with nothing to do is fine
				continue
			}
			t.Fatalf("car %d stalled %ds at floor %.2f with work pending: riders=%d dwell=%.2f dir=%d to=%.2f",
				c.ID, (step-lastMove[c.ID])/30, c.Floor, c.Riders(), c.Dwell(), c.Dir, c.To)
		}
	}
}

// Whether anyone is riding or waiting for this car.
func (w *World) carHasWork(car transport.CarID) bool {
	for _, p := range w.sims {
		if p.car != car {
			continue
		}
		if p.state == stateWaiting || p.holdsCarPlace() {
			return true
		}
	}
	return false
}
