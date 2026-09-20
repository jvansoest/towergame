package main

import (
	"strings"
	"testing"

	"towergame/server/model"
)

// Adds n resident sims, to raise the population.
func addPeople(w *World, n int) {
	for i := 0; i < n; i++ {
		w.nextID++
		w.sims = append(w.sims, &sim{id: w.nextID, prof: profResident})
	}
}

// Stars follow the population and needed rooms.
func TestStarsFollowPopulation(t *testing.T) {
	w := newWorld()
	if got := w.earnedStars(); got != 1 {
		t.Fatalf("an empty tower has %d stars", got)
	}
	addPeople(w, 60)
	if got := w.earnedStars(); got != 2 {
		t.Fatalf("60 people earn %d stars, want 2", got)
	}
	addPeople(w, 100)
	if got := w.earnedStars(); got != 2 {
		t.Fatalf("no security office, yet %d stars", got)
	}
	w.grid.Rooms = append(w.grid.Rooms, model.Room{Type: "security"})
	if got := w.earnedStars(); got != 3 {
		t.Fatalf("with security, want 3 stars, got %d", got)
	}
	addPeople(w, 200)
	if got := w.earnedStars(); got != 3 {
		t.Fatalf("no medical centre, yet %d stars", got)
	}
	w.grid.Rooms = append(w.grid.Rooms, model.Room{Type: "medical"})
	if got := w.earnedStars(); got != 4 {
		t.Fatalf("with medical, want 4 stars, got %d", got)
	}
	addPeople(w, 200)
	w.profitable = true
	if got := w.earnedStars(); got != 5 {
		t.Fatalf("with profit, want 5 stars, got %d", got)
	}
}

// The rating never falls, and unlocks rooms.
func TestStarsGateRooms(t *testing.T) {
	w := newWorld()
	for c := 0; c < 12; c++ {
		w.grid.Built[0][c] = true
	}
	w.stars = 1
	err := w.Place("izakaya", 0, 1, model.AlignNeutral)
	if err == nil || !strings.Contains(err.Error(), "3 stars") {
		t.Fatalf("a 1 star tower built an izakaya: %v", err)
	}
	if err := w.PlaceEscalator(0, 2); err == nil {
		t.Fatal("a 1 star tower built an escalator")
	}
	w.stars = 3
	if err := w.Place("izakaya", 0, 1, model.AlignNeutral); err != nil {
		t.Fatalf("a 3 star tower cannot build one: %v", err)
	}
	w.tick = starCheckEvery
	w.updateStars()
	if w.stars < 3 {
		t.Fatalf("the rating fell to %d", w.stars)
	}
}

// A medical centre has doctors and paying patients.
func TestMedicalCentreWorks(t *testing.T) {
	w, _ := roomWorld(t, "medical")
	doctors := 0
	for _, p := range w.sims {
		if p.prof == profDoctor {
			doctors++
		}
	}
	if doctors != model.RoomTypes["medical"].Capacity {
		t.Fatalf("got %d doctors", doctors)
	}
	before := w.books.sales
	for i := 0; i < 900; i++ {
		w.Step(1.0 / 30)
	}
	if w.books.sales <= before {
		t.Fatal("the clinic saw no paying patients")
	}
}

// Inspecting a hotel room says clean or dirty.
func TestRoomCardSaysDirty(t *testing.T) {
	w, r := roomWorld(t, "hotel_single")
	has := func(s string) bool {
		card := w.roomCard(w.grid.Rooms[0])
		for _, l := range card.Lines {
			if l == s {
				return true
			}
		}
		return false
	}
	if !has("Cleanliness: clean") {
		t.Fatal("a fresh room is not called clean")
	}
	w.grid.Rooms[0].Dirty = true
	if !has("Cleanliness: dirty") {
		t.Fatal("a dirty room is not called dirty")
	}
	_ = r
}
