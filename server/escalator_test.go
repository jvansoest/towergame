package main

import (
	"testing"

	"towergame/server/model"
)

// Escalators cost more than stairs and refund half.
func TestEscalatorCostsAndRefunds(t *testing.T) {
	w := newWorld()
	for c := 0; c < 12; c++ {
		w.grid.Built[0][c] = true
	}
	before := w.money
	if err := w.PlaceEscalator(0, 2); err != nil {
		t.Fatal(err)
	}
	if before-w.money != escalatorCost {
		t.Fatalf("cost %d, want %d", before-w.money, escalatorCost)
	}
	if !w.grid.Stairs[0].Escalator {
		t.Fatal("the run is not an escalator")
	}
	if err := w.Remove(0, 2); err != nil {
		t.Fatal(err)
	}
	if got := w.money; got != before-escalatorCost+escalatorCost/2 {
		t.Fatalf("money %d after removal, want a half refund", got)
	}
}

// Routes over an escalator use the gliding mode.
func TestEscalatorRouteGlides(t *testing.T) {
	g := model.NewGrid(20, 6)
	for f := 0; f < 3; f++ {
		for c := 0; c < 20; c++ {
			g.Built[f][c] = true
		}
	}
	if err := g.PlaceEscalator(0, 4); err != nil {
		t.Fatal(err)
	}
	path, ok := g.Path(0, 4, 1, 4)
	if !ok {
		t.Fatal("no route over the escalator")
	}
	found := false
	for _, wp := range path {
		if wp.Mode == model.ModeEscalator {
			found = true
		}
	}
	if !found {
		t.Fatalf("path %+v never glides", path)
	}
}

// A sim on an escalator reads as gliding, not walking.
func TestSimGlidesUpEscalator(t *testing.T) {
	w := newWorld()
	for f := 0; f < 3; f++ {
		for c := 3; c < 20; c++ {
			w.grid.Built[f][c] = true
		}
	}
	if err := w.grid.PlaceEscalator(0, 6); err != nil {
		t.Fatal(err)
	}
	p := &sim{id: 1, pace: 1, transient: true, homeF: 1, homeC: 6, state: stateOutside}
	w.sims = []*sim{p}
	seen := false
	for i := 0; i < 900 && !(p.state == statePresent && p.y > 0.99); i++ {
		w.Step(1.0 / 30)
		for _, v := range w.SimsUpdate().Sims {
			seen = seen || v.Gliding
		}
	}
	if !seen {
		t.Fatal("the sim never showed as gliding")
	}
	if p.y < 0.99 {
		t.Fatalf("the sim ended on floor %.2f", p.y)
	}
}
