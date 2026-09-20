package main

import (
	"testing"

	"towergame/server/model"
)

func TestWorldMoney(t *testing.T) {
	w := newWorld()
	start := w.money
	if err := w.Place("office", 0, 0, model.AlignNeutral); err != nil {
		t.Fatalf("place: %v", err)
	}
	if want := start - model.RoomTypes["office"].Cost; w.money != want {
		t.Fatalf("money: want %d, got %d", want, w.money)
	}
}

func TestWorldBroke(t *testing.T) {
	w := newWorld()
	w.money = 100
	if err := w.Place("office", 0, 0, model.AlignNeutral); err == nil {
		t.Fatal("want not-enough-money error, got nil")
	}
	if len(w.grid.Rooms) != 0 {
		t.Fatal("a rejected room must not be placed")
	}
}

func TestWorldPlaceAndSnapshot(t *testing.T) {
	w := newWorld()
	if err := w.Place("office", 0, 0, model.AlignNeutral); err != nil {
		t.Fatalf("place: %v", err)
	}
	snap := w.Snapshot()
	if len(snap.Grid.Rooms) != 1 || snap.Grid.Rooms[0].Type != "office" {
		t.Fatalf("want one office, got %+v", snap.Grid.Rooms)
	}
}

// The snapshot carries the buildable catalog,
// sorted, so the client never guesses facts.
func TestSnapshotCarriesTypes(t *testing.T) {
	w := newWorld()
	types := w.Snapshot().Grid.Types
	if len(types) != len(model.RoomTypes) {
		t.Fatalf("want %d types, got %d", len(model.RoomTypes), len(types))
	}
	for i := 1; i < len(types); i++ {
		if types[i-1].ID >= types[i].ID {
			t.Fatalf("catalog not sorted at %q", types[i].ID)
		}
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

// Demolishing a room clears its visitors out,
// not only the residents. Left behind, they
// would block a rebuilt room's customer flow.
func TestRemoveRoomDropsVisitorsToo(t *testing.T) {
	w := newWorld()
	if err := w.Place("office", 0, 0, model.AlignNeutral); err != nil {
		t.Fatalf("place office: %v", err)
	}
	w.nextID++
	w.sims = append(w.sims, &sim{
		id: w.nextID, category: model.CategoryRetail,
		homeF: 0, homeC: 3, transient: true, spot: 1,
	})

	if err := w.Remove(0, 4); err != nil {
		t.Fatalf("remove office: %v", err)
	}
	if len(w.sims) != 0 {
		t.Fatalf("demolition left %d sims in the room", len(w.sims))
	}
}

// A floor carrying structure overhead may
// not be stripped; clear downward instead.
func TestRemoveKeepsSupportBelow(t *testing.T) {
	w := newWorld()
	if err := w.BuildBase(0, 20); err != nil {
		t.Fatalf("ground base: %v", err)
	}
	if err := w.BuildBase(1, 20); err != nil {
		t.Fatalf("upper base: %v", err)
	}

	if err := w.Remove(0, 20); err == nil {
		t.Fatal("clearing a carried floor must fail")
	}
	if !w.grid.Built[0][20] {
		t.Fatal("the support cell must survive")
	}

	if err := w.Remove(1, 20); err != nil {
		t.Fatalf("remove upper first: %v", err)
	}
	if err := w.Remove(0, 20); err != nil {
		t.Fatalf("then ground clears: %v", err)
	}
}
