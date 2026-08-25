package model

import "testing"

func TestPathSameFloor(t *testing.T) {
	g := NewGrid(40, 5)
	wp, ok := g.Path(0, 2, 0, 10)
	if !ok {
		t.Fatal("same-floor path should exist")
	}
	if len(wp) != 2 {
		t.Fatalf("want 2 waypoints, got %d", len(wp))
	}
}

func TestPathNoStairs(t *testing.T) {
	g := NewGrid(40, 5)
	if _, ok := g.Path(0, 0, 2, 0); ok {
		t.Fatal("no stairs: path should not exist")
	}
}

func TestPathUpStairs(t *testing.T) {
	g := NewGrid(40, 5)
	for c := 0; c < 40; c++ {
		g.BuildBase(0, c) // ground floor to support upper stairs
	}
	g.PlaceStair(0, 10) // connects floor 0 and 1
	g.PlaceStair(1, 20) // connects floor 1 and 2
	wp, ok := g.Path(0, 0, 2, 30)
	if !ok {
		t.Fatal("two stairs should connect floors 0 and 2")
	}
	if wp[len(wp)-1].Floor != 2 {
		t.Fatalf("route should end on floor 2, got %+v", wp[len(wp)-1])
	}
}

// Sims use whichever shaft is nearest, not the first built.
func TestPicksNearestElevator(t *testing.T) {
	g := NewGrid(60, 8)
	for f := 0; f <= 5; f++ {
		for c := 0; c < g.Width; c++ {
			g.Built[f][c] = true
		}
	}
	if _, err := g.PlaceElevator(5, 4); err != nil {
		t.Fatalf("first elevator: %v", err)
	}
	if _, err := g.PlaceElevator(5, 50); err != nil {
		t.Fatalf("second elevator: %v", err)
	}
	near := g.Elevators[0].ID // column 4
	far := g.Elevators[1].ID  // column 50

	// Near the far shaft, so it must be chosen.
	wp, ok := g.Path(0, 52, 5, 52)
	if !ok {
		t.Fatal("no path")
	}
	var used ShaftID
	for _, p := range wp {
		if p.Mode == ModeElevator {
			used = p.Shaft
		}
	}
	if used != far {
		t.Fatalf("used shaft %d, want the near one %d", used, far)
	}

	// And the other way round.
	wp, ok = g.Path(0, 2, 5, 2)
	if !ok {
		t.Fatal("no path")
	}
	used = 0
	for _, p := range wp {
		if p.Mode == ModeElevator {
			used = p.Shaft
		}
	}
	if used != near {
		t.Fatalf("used shaft %d, want the near one %d", used, near)
	}
}
