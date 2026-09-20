package main

import (
	"testing"

	"towergame/server/model"
)

// A garage: a two-run ramp at column 10, rooms beside it.
func garageWorld(t *testing.T) *World {
	t.Helper()
	w := newWorld()
	for c := 0; c < 32; c++ {
		w.grid.SetBuilt(-1, c, true)
	}
	for c := 4; c < 18; c++ {
		w.grid.SetBuilt(-2, c, true)
	}
	for _, f := range []int{-1, -2} {
		if err := w.grid.PlaceRamp(f, 10); err != nil {
			t.Fatalf("ramp %d: %v", f, err)
		}
	}
	setClock(w, 10, 0)
	return w
}

// Cars drive in, park, leave, and pay.
func TestCarsParkAndPay(t *testing.T) {
	w := garageWorld(t)
	if err := w.Place("parking", -1, 18, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	before := w.money
	parked, left := false, false
	for i := 0; i < 30*200 && !left; i++ {
		w.Step(1.0 / 30)
		for _, v := range w.SimsUpdate().Vehicles {
			parked = parked || v.Parked
		}
		if parked && len(w.vehicles) == 0 {
			left = true
		}
	}
	if !parked {
		t.Fatal("no car ever parked")
	}
	if !left {
		t.Fatal("the car never left")
	}
	if w.money <= before {
		t.Fatalf("no fee paid: %d to %d", before, w.money)
	}
}

// Both rows work, on the side each run points to.
func TestRoutesFollowRunDirection(t *testing.T) {
	w := garageWorld(t)
	if err := w.Place("parking", -2, 6, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	deep := w.grid.Rooms[0]
	path, ok := w.routeTo(deep, 0)
	if !ok {
		t.Fatal("no route to the lower garage")
	}
	// Street, the two run ends, then the slot.
	if len(path) != 4 || path[0].y != 0 || path[len(path)-1].y != -2 {
		t.Fatalf("route %+v", path)
	}
	if err := w.Place("parking", -1, 2, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	wrong := w.grid.Rooms[1]
	if _, ok := w.routeTo(wrong, 0); ok {
		t.Fatal("a route to the wrong side of the ramp")
	}
}

// Two slots hold two cars, never three.
func TestSlotsFillOnce(t *testing.T) {
	w := garageWorld(t)
	if err := w.Place("parking", -1, 18, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 30*80; i++ {
		w.Step(1.0 / 30)
		// A car driving out has given its slot up.
		held := 0
		for _, v := range w.vehicles {
			if v.phase != carLeaving {
				held++
			}
		}
		if held > model.RoomTypes["parking"].Slots {
			t.Fatalf("%d cars hold a %d slot garage", held, model.RoomTypes["parking"].Slots)
		}
	}
}

// Removing the garage clears its cars.
func TestGarageGoneClearsCars(t *testing.T) {
	w := garageWorld(t)
	if err := w.Place("parking", -1, 18, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 30*30 && len(w.vehicles) == 0; i++ {
		w.Step(1.0 / 30)
	}
	if len(w.vehicles) == 0 {
		t.Fatal("no car came")
	}
	if err := w.Remove(-1, 18); err != nil {
		t.Fatal(err)
	}
	w.Step(1.0 / 30)
	if len(w.vehicles) != 0 {
		t.Fatalf("%d cars remain with no garage", len(w.vehicles))
	}
}

// Digging out from under a floor is refused.
func TestDugFloorsKeepTheirAccess(t *testing.T) {
	w := garageWorld(t)
	if err := w.Remove(-1, 6); err == nil {
		t.Fatal("removed floor -1 over a dug floor -2")
	}
	if err := w.Remove(-2, 6); err != nil {
		t.Fatalf("could not remove the lowest floor: %v", err)
	}
}
