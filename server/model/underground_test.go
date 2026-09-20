package model

import "testing"

// Floors are dug from the surface downward.
func TestDigDown(t *testing.T) {
	g := NewGrid(30, 5)
	if err := g.BuildBase(-1, 4); err != nil {
		t.Fatalf("dig floor -1: %v", err)
	}
	if err := g.BuildBase(-3, 4); err == nil {
		t.Fatal("dug floor -3 with nothing above")
	}
	if err := g.BuildBase(-2, 4); err != nil {
		t.Fatalf("dig floor -2 under a dug cell: %v", err)
	}
	if !g.IsBuilt(-2, 4) || g.IsBuilt(-2, 5) {
		t.Fatal("the dug cells are wrong")
	}
	if err := g.BuildBase(-Basement-1, 4); err == nil {
		t.Fatal("dug below the last basement floor")
	}
}

// Parking goes underground only, and never across the ground.
func TestParkingUndergroundOnly(t *testing.T) {
	g := NewGrid(30, 5)
	for c := 0; c < 6; c++ {
		g.SetBuilt(-1, c, true)
		g.Built[0][c] = true
	}
	if err := g.Place("parking", 0, 0); err == nil {
		t.Fatal("parking was built above ground")
	}
	if err := g.Place("parking", -1, 0); err != nil {
		t.Fatalf("parking underground: %v", err)
	}
	if rt := RoomTypes["parking"]; rt.Width != 2 {
		t.Fatalf("parking is %d cells wide, want 2", rt.Width)
	}
}

// A ramp opens on the street and winds down a run at a time.
func TestRampWinds(t *testing.T) {
	g := NewGrid(40, 5)
	for c := 0; c < 20; c++ {
		g.SetBuilt(-1, c, true)
		g.SetBuilt(-2, c, true)
	}
	if err := g.PlaceRamp(-2, 4); err == nil {
		t.Fatal("a lower run stood with no run above")
	}
	g.Built[0][4] = true
	if err := g.PlaceRamp(-1, 4); err == nil {
		t.Fatal("the ramp opened under a built ground cell")
	}
	g.Built[0][4] = false
	if err := g.PlaceRamp(-1, 4); err != nil {
		t.Fatalf("first run: %v", err)
	}
	if err := g.PlaceRamp(-2, 4); err != nil {
		t.Fatalf("second run: %v", err)
	}
	r1, _ := g.RampAt(-1, 4)
	r2, _ := g.RampAt(-2, 4)
	if !r1.DownRight() || r2.DownRight() {
		t.Fatal("the runs do not turn back on each other")
	}
	if r1.BottomCol() != r2.TopCol() {
		t.Fatalf("run 1 ends at %d, run 2 starts at %d", r1.BottomCol(), r2.TopCol())
	}
	if _, ok, err := g.RemoveRampAt(-1, 4); ok || err == nil {
		t.Fatal("removed a run with another hanging below")
	}
	if _, ok, _ := g.RemoveRampAt(-2, 4); !ok {
		t.Fatal("could not remove the lowest run")
	}
}

// Rooms may not sit on the ramp.
func TestRampBlocksRooms(t *testing.T) {
	g := NewGrid(40, 5)
	for c := 0; c < 20; c++ {
		g.SetBuilt(-1, c, true)
	}
	if err := g.PlaceRamp(-1, 4); err != nil {
		t.Fatal(err)
	}
	if err := g.Place("parking", -1, 6); err == nil {
		t.Fatal("a parking room was built on the ramp")
	}
	if err := g.Place("parking", -1, 14); err != nil {
		t.Fatalf("parking beside the ramp: %v", err)
	}
}

// Routes work between floors below ground too.
func TestPathBelowGround(t *testing.T) {
	g := NewGrid(30, 5)
	for f := -1; f >= -3; f-- {
		for c := 0; c < 20; c++ {
			g.SetBuilt(f, c, true)
		}
	}
	if err := g.PlaceEscalator(-3, 8); err != nil {
		t.Fatal(err)
	}
	if err := g.PlaceStair(-2, 8); err != nil {
		t.Fatal(err)
	}
	if _, ok := g.Path(-3, 2, -1, 2); !ok {
		t.Fatal("no route from floor -3 to floor -1")
	}
	if _, ok := g.Path(-1, 2, -3, 12); !ok {
		t.Fatal("no route back down")
	}
}
