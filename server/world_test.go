package main

import (
	"testing"

	"towergame/server/model"
)

func TestWorldMoney(t *testing.T) {
	w := newWorld()
	start := w.money
	if err := w.Place("office", 0, 0); err != nil {
		t.Fatalf("place: %v", err)
	}
	if want := start - model.RoomTypes["office"].Cost; w.money != want {
		t.Fatalf("money: want %d, got %d", want, w.money)
	}
}

func TestWorldBroke(t *testing.T) {
	w := newWorld()
	w.money = 100
	if err := w.Place("office", 0, 0); err == nil {
		t.Fatal("want not-enough-money error, got nil")
	}
	if len(w.grid.Rooms) != 0 {
		t.Fatal("a rejected room must not be placed")
	}
}

func TestWorldPlaceAndSnapshot(t *testing.T) {
	w := newWorld()
	if err := w.Place("office", 0, 0); err != nil {
		t.Fatalf("place: %v", err)
	}
	snap := w.Snapshot()
	if len(snap.Grid.Rooms) != 1 || snap.Grid.Rooms[0].Type != "office" {
		t.Fatalf("want one office, got %+v", snap.Grid.Rooms)
	}
}

func TestWorldStepAdvancesClock(t *testing.T) {
	w := newWorld()
	if !w.Step(1) {
		t.Fatal("first step should change the clock")
	}
	if w.tick != 1 {
		t.Fatalf("tick should be 1, got %d", w.tick)
	}
}

func TestWorldPostMessage(t *testing.T) {
	w := newWorld()
	if w.PostMessage("   ") {
		t.Fatal("blank message should be rejected")
	}
	if !w.PostMessage("hi") {
		t.Fatal("message should be accepted")
	}
	if got := w.ChatUpdate(); len(got.Chat) != 1 {
		t.Fatalf("want 1 chat line, got %d", len(got.Chat))
	}
}
