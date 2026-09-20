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
	if err := w.Place(typeID, 0, 1, model.AlignNeutral); err != nil {
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

// A home at daybreak: everyone is in.
func condoWorld(t *testing.T) (*World, model.Room) {
	t.Helper()
	w := newWorld()
	for c := 0; c < 12; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("condo", 0, 1, model.AlignNeutral); err != nil {
		t.Fatalf("place condo: %v", err)
	}
	return w, model.Room{Type: "condo", Floor: 0, Col: 1}
}

// Seated means parked on the chair: a sitter
// out for a stroll must report as up.
func TestSeatedMeansParked(t *testing.T) {
	w, r := roomWorld(t, "condo")
	col := model.SeatCol(r, 0)
	p := &sim{
		id: 991, category: model.CategoryResidential,
		homeF: r.Floor, homeC: col,
		state: statePresent,
		x:     float64(col), y: float64(r.Floor),
		pause: 2, // mid-rest at their seat
	}
	w.sims = []*sim{p}

	if !w.seated(p) {
		t.Fatal("a rested resident on their chair reads as up")
	}

	p.x = float64(col + 1)
	if w.seated(p) {
		t.Fatal("a resident away from their chair reads seated")
	}

	p.x = float64(col)
	p.pause = 0 // stirring, not resting yet
	if w.seated(p) {
		t.Fatal("the seat claims someone who is walking")
	}
}

// Residents rise now and then, wander the
// flat slowly, and never leave it mid-visit.
func TestResidentsStrollAtHome(t *testing.T) {
	w, r := condoWorld(t)
	for _, p := range w.sims {
		p.shiftIn = 999 // hold the day: nobody leaves for work
	}

	for i := 0; i < 1800; i++ { // a minute
		w.Step(1.0 / 30)
	}
	rt := model.RoomTypes[r.Type]
	moved := false
	for _, p := range w.sims {
		if math.Abs(p.x-float64(p.homeC)) > 0.25 {
			moved = true
		}
		if p.x < float64(r.Col) || p.x >= float64(r.Col+rt.Width) {
			t.Fatalf("resident %d left the flat: x=%.2f", p.id, p.x)
		}
	}
	if !moved {
		t.Fatal("nobody stirred at home in a minute")
	}
}

// Izakaya seats line up in one column, one per spot.
func TestIzakayaSeatsShareColumn(t *testing.T) {
	rt := model.RoomTypes["izakaya"]
	if !rt.Line || rt.Seats != 8 {
		t.Fatalf("izakaya is %+v", rt)
	}
	r := model.Room{Type: "izakaya", Floor: 0, Col: 1}
	want := model.SeatCol(r, 0)
	for i := 1; i < rt.Seats; i++ {
		if got := model.SeatCol(r, i); got != want {
			t.Fatalf("seat %d at column %d, want %d", i, got, want)
		}
	}
	if want < r.Col || want >= r.Col+rt.Width {
		t.Fatalf("seat column %d is outside the room", want)
	}
}

// Izakaya guests sit still on their stools.
func TestIzakayaGuestsStaySeated(t *testing.T) {
	w, r := roomWorld(t, "izakaya")
	seat := float64(model.SeatCol(r, 0))
	seen := 0
	for i := 0; i < 1500; i++ {
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			if !p.transient || p.state != statePresent {
				continue
			}
			seen++
			if p.x != seat {
				t.Fatalf("guest %d at column %.2f, want %.0f", p.id, p.x, seat)
			}
			if !w.seated(p) {
				t.Fatalf("guest %d is on the stool but not seated", p.id)
			}
		}
	}
	if seen == 0 {
		t.Fatal("no guest reached the izakaya")
	}
}
