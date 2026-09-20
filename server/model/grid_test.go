package model

import "testing"

func TestPlaceOffice(t *testing.T) {
	g := NewGrid(40, 10)
	if err := g.Place("office", 0, 0); err != nil {
		t.Fatalf("place office: %v", err)
	}
	if len(g.Rooms) != 1 {
		t.Fatalf("want 1 room, got %d", len(g.Rooms))
	}
}

func TestOverlapRejected(t *testing.T) {
	g := NewGrid(40, 10)
	g.Place("office", 0, 0) // width 9 spans cols 0..8

	if err := g.Place("shop", 0, 5); err == nil {
		t.Fatal("want overlap error, got nil")
	}
	// Column 9 is clear.
	if err := g.Place("shop", 0, 9); err != nil {
		t.Fatalf("place shop: %v", err)
	}
}

func TestSameColsDifferentFloors(t *testing.T) {
	g := NewGrid(40, 10)
	g.Place("office", 0, 0)
	if err := g.Place("office", 1, 0); err != nil {
		t.Fatalf("stacking should be allowed: %v", err)
	}
}

func TestSupportRequired(t *testing.T) {
	g := NewGrid(40, 10)
	// Floor 1 has nothing below it.
	if err := g.Place("office", 1, 0); err == nil {
		t.Fatal("want support error, got nil")
	}
	// Ground floor builds the base below.
	g.Place("office", 0, 0)
	if err := g.Place("office", 1, 0); err != nil {
		t.Fatalf("floor 1 should be supported: %v", err)
	}
}

func TestLobbyGroundOnly(t *testing.T) {
	g := NewGrid(40, 10)
	if err := g.Place("lobby", 0, 5); err != nil {
		t.Fatalf("ground lobby: %v", err)
	}
	// Support exists, so only the ground rule can reject this.
	g.BuildBase(0, 8)
	if err := g.Place("lobby", 1, 8); err == nil {
		t.Fatal("want ground-only error, got nil")
	}
}

func TestBuildBase(t *testing.T) {
	g := NewGrid(40, 10)
	if err := g.BuildBase(1, 0); err == nil {
		t.Fatal("want support error, got nil")
	}
	if err := g.BuildBase(0, 0); err != nil {
		t.Fatalf("ground base: %v", err)
	}
	if err := g.BuildBase(1, 0); err != nil {
		t.Fatalf("supported base: %v", err)
	}
}

func TestHeightExceedsFloors(t *testing.T) {
	g := NewGrid(40, 5)
	// Cinema is 2 floors tall; the top floor leaves no room.
	if err := g.Place("cinema", 4, 0); err == nil {
		t.Fatal("want height error, got nil")
	}
}

func TestOutOfRange(t *testing.T) {
	g := NewGrid(10, 5)
	if err := g.Place("office", 0, 5); err == nil {
		t.Fatal("want column error, got nil") // width 9 from col 5 exceeds width 10
	}
	if err := g.Place("office", 9, 0); err == nil {
		t.Fatal("want floor error, got nil")
	}
}

func TestUnknownType(t *testing.T) {
	g := NewGrid(40, 10)
	if err := g.Place("penthouse", 0, 0); err == nil {
		t.Fatal("want unknown-type error, got nil")
	}
}
