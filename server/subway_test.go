package main

import (
	"strings"
	"testing"

	"towergame/server/model"
)

// A dug basement with a station, escalators, and a shop.
func metroWorld(t *testing.T) *World {
	t.Helper()
	w := newWorld()
	for f := -1; f >= -model.Basement; f-- {
		for c := 0; c < 24; c++ {
			w.grid.SetBuilt(f, c, true)
		}
	}
	if err := w.Place("subway", -model.Basement, 2, model.AlignNeutral); err != nil {
		t.Fatalf("subway: %v", err)
	}
	for _, f := range []int{-6, -5} {
		if err := w.grid.PlaceEscalator(f, 16); err != nil {
			t.Fatalf("escalator %d: %v", f, err)
		}
	}
	if err := w.Place("shop", -4, 2, model.AlignNeutral); err != nil {
		t.Fatalf("shop: %v", err)
	}
	setClock(w, 10, 0)
	return w
}

// The station sits on the lowest floor, once per tower.
func TestSubwayRules(t *testing.T) {
	w := newWorld()
	for f := -1; f >= -model.Basement; f-- {
		for c := 0; c < 40; c++ {
			w.grid.SetBuilt(f, c, true)
		}
	}
	if err := w.Place("subway", -4, 2, model.AlignNeutral); err == nil {
		t.Fatal("the subway was built above the lowest floor")
	}
	if err := w.Place("subway", -model.Basement, 2, model.AlignNeutral); err != nil {
		t.Fatalf("the subway on the lowest floor: %v", err)
	}
	err := w.Place("subway", -model.Basement, 20, model.AlignNeutral)
	if err == nil || !strings.Contains(err.Error(), "only one") {
		t.Fatalf("a second subway was built: %v", err)
	}
	if err := w.Remove(-model.Basement, 5); err == nil {
		t.Fatal("the subway was demolished")
	}
}

// A train comes, unloads shoppers, and they go back to it.
func TestTrainsBringShoppers(t *testing.T) {
	w := metroWorld(t)
	before := w.books.sales
	sawTrain, sawStopped, sawRider := false, false, false
	for i := 0; i < 30*120; i++ {
		w.Step(1.0 / 30)
		if tr := w.SimsUpdate().Train; tr != nil {
			sawTrain = true
			sawStopped = sawStopped || tr.Stopped
		}
		for _, p := range w.sims {
			if p.subway {
				sawRider = true
			}
		}
	}
	if !sawTrain || !sawStopped {
		t.Fatalf("train seen %v, stopped at the platform %v", sawTrain, sawStopped)
	}
	if !sawRider {
		t.Fatal("no passengers stepped off")
	}
	if w.books.sales <= before {
		t.Fatal("the passengers bought nothing")
	}
}

// Passengers never surface at the street.
func TestPassengersStayUnderground(t *testing.T) {
	w := metroWorld(t)
	for i := 0; i < 30*120; i++ {
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			if p.subway && p.state != stateOutside && p.y > -1 {
				t.Fatalf("a passenger reached floor %.1f", p.y)
			}
		}
	}
}

// No station, no trains.
func TestNoStationNoTrain(t *testing.T) {
	w := newWorld()
	setClock(w, 10, 0)
	for i := 0; i < 30*60; i++ {
		w.Step(1.0 / 30)
	}
	if w.SimsUpdate().Train != nil {
		t.Fatal("a train ran with no station")
	}
}

// The starter tower has a station that works.
func TestSeededSubwayWorks(t *testing.T) {
	w := newWorld()
	w.seed()
	setClock(w, 10, 0)
	riders := false
	for i := 0; i < 30*90 && !riders; i++ {
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			riders = riders || p.subway
		}
	}
	if !riders {
		t.Fatal("no train passengers in the starter tower")
	}
}
