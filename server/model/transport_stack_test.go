package model

import "testing"

// Stairs may stack directly on top of each other.
func TestStackedStairs(t *testing.T) {
	g := NewGrid(20, 6)
	for c := 0; c < 20; c++ {
		g.Built[0][c] = true
	}
	if err := g.PlaceStair(0, 2); err != nil {
		t.Fatalf("first stair: %v", err)
	}
	if err := g.PlaceStair(1, 2); err != nil {
		t.Fatalf("stacked stair: %v", err)
	}
	if err := g.PlaceStair(1, 3); err == nil {
		t.Fatal("want overlap on the same floor")
	}
}
