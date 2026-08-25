package model

import "testing"

func TestPlaceStair(t *testing.T) {
	g := NewGrid(40, 10)
	if err := g.PlaceStair(0, 0); err != nil {
		t.Fatalf("ground stair: %v", err)
	}
	if len(g.Stairs) != 1 {
		t.Fatalf("want 1 stair, got %d", len(g.Stairs))
	}
	// The stair built both floors.
	if !g.Built[0][0] || !g.Built[1][0] {
		t.Fatal("stair should build both floors")
	}
}

func TestStairNeedsSupport(t *testing.T) {
	g := NewGrid(40, 10)
	if err := g.PlaceStair(2, 0); err == nil {
		t.Fatal("want support error, got nil")
	}
}

func TestStairCoexistsWithRoom(t *testing.T) {
	g := NewGrid(40, 10)
	g.Place("office", 0, 0) // spans cols 0..8 on floor 0
	// Stairs are a front layer, so a room behind is fine.
	if err := g.PlaceStair(0, 2); err != nil {
		t.Fatalf("stair over a room should be allowed: %v", err)
	}
}

func TestStairOverlapsStair(t *testing.T) {
	g := NewGrid(40, 10)
	g.PlaceStair(0, 0)
	if err := g.PlaceStair(0, 2); err == nil {
		t.Fatal("want stair overlap error, got nil")
	}
}

func TestPlaceElevator(t *testing.T) {
	g := NewGrid(40, 20)
	if _, err := g.PlaceElevator(10, 5); err != nil {
		t.Fatalf("elevator: %v", err)
	}
	// It connects distant floors in one hop.
	wp, ok := g.Path(0, 0, 8, 30)
	if !ok {
		t.Fatal("elevator should connect floors 0 and 8")
	}
	if wp[len(wp)-1].Floor != 8 {
		t.Fatalf("route should reach floor 8, got %+v", wp[len(wp)-1])
	}
}

func TestExtendElevator(t *testing.T) {
	g := NewGrid(40, 20)
	g.PlaceElevator(5, 5)
	// Clicking the same shaft higher extends it, not a new one.
	extended, err := g.PlaceElevator(12, 6)
	if err != nil {
		t.Fatalf("extend: %v", err)
	}
	if !extended {
		t.Fatal("want extension, got a new elevator")
	}
	if len(g.Elevators) != 1 {
		t.Fatalf("want 1 elevator, got %d", len(g.Elevators))
	}
	if g.Elevators[0].Top != 12 {
		t.Fatalf("want top 12, got %d", g.Elevators[0].Top)
	}
}
