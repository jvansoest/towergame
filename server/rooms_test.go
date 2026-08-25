package main

import (
	"math"
	"testing"

	"towergame/server/model"
)

// A one-floor world holding a single room.
func roomWorld(t *testing.T, typeID string) (*World, model.Room) {
	t.Helper()
	w := newWorld()
	rt := model.RoomTypes[typeID]
	for c := 0; c < rt.Width+2; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place(typeID, 0, 1); err != nil {
		t.Fatalf("place %s: %v", typeID, err)
	}
	// Noon on a weekday: everyone is in.
	w.simTime = float64(12*60-startMinuteOfDay) / gameMinutesPerRealSecond
	return w, model.Room{Type: typeID, Floor: 0, Col: 1}
}

// Workers sit at the desks the client draws.
func TestWorkersSitAtDesks(t *testing.T) {
	w, r := roomWorld(t, "office")
	for i := 0; i < 300; i++ {
		w.Step(1.0 / 30)
	}
	if len(w.sims) == 0 {
		t.Fatal("the office has no workers")
	}
	seats := map[int]bool{}
	for i := 0; i < model.RoomTypes["office"].Seats; i++ {
		seats[model.SeatCol(r, i)] = true
	}
	for _, p := range w.sims {
		if !w.seated(p) {
			t.Fatalf("worker %d is not seated (state %d)", p.id, p.state)
		}
		if !seats[int(math.Round(p.x))] {
			t.Fatalf("worker %d sits at column %.2f, off any desk", p.id, p.x)
		}
	}
}

// Diners take a seat each, never sharing one.
// They come and go, so the room only fills at times.
func TestDinersTakeSeparateSeats(t *testing.T) {
	w, _ := roomWorld(t, "restaurant")
	seats := model.RoomTypes["restaurant"].Seats

	fullest := 0
	for i := 0; i < 1800; i++ {
		w.Step(1.0 / 30)
		used := map[int]bool{}
		for _, p := range w.sims {
			if used[p.homeC] {
				t.Fatalf("tick %d: two diners share column %d", i, p.homeC)
			}
			used[p.homeC] = true
		}
		if len(w.sims) > seats {
			t.Fatalf("tick %d: %d diners for %d seats", i, len(w.sims), seats)
		}
		if len(used) > fullest {
			fullest = len(used)
		}
	}
	if fullest < seats {
		t.Fatalf("the restaurant only ever filled %d of %d seats", fullest, seats)
	}
}

// Shoppers walk the floor instead of sitting.
func TestShoppersWalkTheFloor(t *testing.T) {
	w, r := roomWorld(t, "shop")
	rt := model.RoomTypes["shop"]

	lo := map[int]float64{}
	hi := map[int]float64{}
	for i := 0; i < 30*120; i++ { // two minutes
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			if p.state != statePresent {
				continue
			}
			if w.seated(p) {
				t.Fatalf("shopper %d sat down in a shop", p.id)
			}
			if _, seen := lo[p.id]; !seen {
				lo[p.id], hi[p.id] = p.x, p.x
			}
			lo[p.id] = math.Min(lo[p.id], p.x)
			hi[p.id] = math.Max(hi[p.id], p.x)
		}
	}
	if len(lo) == 0 {
		t.Fatal("the shop drew no customers")
	}
	for id, low := range lo {
		if hi[id]-low < browseMin {
			t.Fatalf("shopper %d only covered %.2f cells", id, hi[id]-low)
		}
		if low < float64(r.Col) || hi[id] > float64(r.Col+rt.Width-1) {
			t.Fatalf("shopper %d walked out of the shop: %.2f..%.2f", id, low, hi[id])
		}
	}
}
