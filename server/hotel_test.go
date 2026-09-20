package main

import (
	"testing"

	"towergame/server/model"
)

// Sets the clock to a minute of the first day.
func setClock(w *World, hour, minute int) {
	w.simTime = float64(hour*60+minute-startMinuteOfDay) / gameMinutesPerRealSecond
}

// A guest arrives, sleeps, leaves, pays, and the
// room stays dirty until cleaned.
func TestGuestStayThenDirty(t *testing.T) {
	w, _ := roomWorld(t, "hotel_single")
	setClock(w, 16, 30)
	before := w.money
	w.bookRoom(0)
	if !w.grid.Rooms[0].Booked {
		t.Fatal("the room is not booked")
	}
	arrived := false
	for i := 0; i < 9000; i++ {
		w.Step(1.0 / 30)
		if w.grid.Rooms[0].Used {
			arrived = true
		}
		if arrived && !w.grid.Rooms[0].Booked {
			break
		}
	}
	r := w.grid.Rooms[0]
	if !arrived || r.Booked || !r.Dirty {
		t.Fatalf("after the stay: arrived=%v booked=%v dirty=%v", arrived, r.Booked, r.Dirty)
	}
	if got := w.money - before; got != model.RoomTypes["hotel_single"].Rent {
		t.Fatalf("rent paid %d, want %d", got, model.RoomTypes["hotel_single"].Rent)
	}
}

// Dirty rooms wait; clean ones fill in the evening.
func TestOnlyCleanRoomsRent(t *testing.T) {
	w, _ := roomWorld(t, "hotel_single")
	w.grid.Rooms[0].Dirty = true
	for i := 0; i < 3000; i++ {
		setClock(w, 17, 0)
		w.tick = uint64(i * visitorEvery)
		w.manageGuests(17 * 60)
	}
	if w.grid.Rooms[0].Booked {
		t.Fatal("a dirty room was rented")
	}
	w.grid.Rooms[0].Dirty = false
	for i := 0; i < 3000 && !w.grid.Rooms[0].Booked; i++ {
		w.tick = uint64(i * visitorEvery)
		w.manageGuests(17 * 60)
	}
	if !w.grid.Rooms[0].Booked {
		t.Fatal("a clean room never rented")
	}
}

// A maid walks to a dirty room and cleans it.
func TestMaidCleansRoom(t *testing.T) {
	w := newWorld()
	for c := 0; c < 16; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("housekeeping", 0, 1, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	if err := w.Place("hotel_single", 0, 8, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	if len(w.sims) != model.RoomTypes["housekeeping"].Capacity {
		t.Fatalf("want %d maids, got %d sims", model.RoomTypes["housekeeping"].Capacity, len(w.sims))
	}
	setClock(w, 12, 0)
	w.grid.Rooms[1].Dirty = true
	swept := false
	for i := 0; i < 1500 && w.grid.Rooms[1].Dirty; i++ {
		w.Step(1.0 / 30)
		for _, v := range w.SimsUpdate().Sims {
			swept = swept || v.Cleaning
		}
	}
	if w.grid.Rooms[1].Dirty {
		t.Fatal("the room is still dirty")
	}
	if !swept {
		t.Fatal("no maid was ever shown sweeping")
	}
}

// Loud neighbours drive sleeping guests out unpaid.
func TestNoiseDrivesGuestsOut(t *testing.T) {
	w := newWorld()
	for c := 0; c < 10; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("hotel_single", 0, 1, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	if err := w.Place("izakaya", 0, 4, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	setClock(w, 0, 0)
	w.refreshNoise(0)
	if w.noise[0] < 3 {
		t.Fatalf("the single hears %.1f, want the izakaya's 3", w.noise[0])
	}
	w.bookRoom(0)
	disturbed := false
	for i := 0; i < 4000 && !disturbed; i++ {
		w.refreshLights()
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			if p.category == model.CategoryHotel && p.disturbed {
				disturbed = true
			}
		}
	}
	if !disturbed {
		t.Fatal("the noise never woke the guest")
	}
	for i := 0; i < 900; i++ {
		w.Step(1.0 / 30)
	}
	if w.books.hotel != 0 {
		t.Fatalf("a disturbed guest paid %d", w.books.hotel)
	}
}

// Quiet neighbours cause no noise.
func TestOfficeIsQuiet(t *testing.T) {
	w := newWorld()
	for c := 0; c < 14; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("hotel_single", 0, 1, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	if err := w.Place("office", 0, 4, model.AlignNeutral); err != nil {
		t.Fatal(err)
	}
	w.refreshNoise(12 * 60)
	if w.noise[0] != 0 {
		t.Fatalf("an office made noise %.1f", w.noise[0])
	}
}
