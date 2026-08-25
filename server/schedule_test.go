package main

import (
	"testing"

	"towergame/server/model"
)

// Workers keep their own hours, not one shared clock.
func TestSchedulesVary(t *testing.T) {
	w := newWorld()
	w.seed()

	seen := map[int]bool{}
	workers := 0
	for _, p := range w.sims {
		if p.category != model.CategoryOffice {
			continue
		}
		workers++
		seen[p.shiftIn] = true
		if p.shiftIn > scheduleSpread || p.shiftIn < -scheduleSpread {
			t.Fatalf("worker %d starts %d minutes off", p.id, p.shiftIn)
		}
	}
	if workers == 0 {
		t.Fatal("the seeded tower has no office workers")
	}
	if len(seen) < workers/2 {
		t.Fatalf("%d workers share only %d start times", workers, len(seen))
	}
}

// The tower starts settled, not mid-commute.
func TestSeedStartsSettled(t *testing.T) {
	w := newWorld()
	w.seed()
	for _, p := range w.sims {
		if p.state != statePresent && p.state != stateOutside {
			t.Fatalf("sim %d starts in state %d", p.id, p.state)
		}
		if p.state == statePresent && int(p.y) != p.homeF {
			t.Fatalf("sim %d starts at floor %.0f, not home %d", p.id, p.y, p.homeF)
		}
	}
}
