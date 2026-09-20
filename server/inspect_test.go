package main

import (
	"strconv"
	"strings"
	"testing"

	"towergame/server/model"
	"towergame/server/transport"
)

// Joined lines, for checking a card's contents.
func cardText(i Inspection) string {
	return i.Title + "\n" + strings.Join(i.Lines, "\n")
}

// A sim reports where it came from and where it goes.
func TestInspectSimShowsTrip(t *testing.T) {
	w, _ := shaftWorld(t, 8)
	g := w.grid
	// Room for a condo beside the shaft.
	for f := 0; f <= 8; f++ {
		for c := 18; c <= 35; c++ {
			g.Built[f][c] = true
		}
	}
	if err := w.Place("condo", 7, 24, model.AlignNeutral); err != nil {
		t.Fatalf("condo: %v", err)
	}
	if len(w.sims) == 0 {
		t.Fatal("the condo has no residents")
	}
	p := w.sims[0]

	// At home, then sent out of the building.
	p.state = statePresent
	p.x, p.y = float64(p.homeC), float64(p.homeF)
	p.last = goal{true, p.homeF, p.homeC}
	w.applyGoal(p, goal{false, p.homeF, p.homeC})
	card, err := w.InspectSim(p.id)
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	text := cardText(card)
	if !strings.Contains(text, "Resident") {
		t.Fatalf("card does not name the kind:\n%s", text)
	}
	if !strings.Contains(text, "From: Condominium, floor 7") {
		t.Fatalf("card does not say where they came from:\n%s", text)
	}
	if !strings.Contains(text, "To: Outside") {
		t.Fatalf("card does not say where they are going:\n%s", text)
	}
	if !strings.Contains(text, "Now: walking") {
		t.Fatalf("card does not say what they are doing:\n%s", text)
	}

	if _, err := w.InspectSim(9999); err == nil {
		t.Fatal("inspecting nobody should fail")
	}
}

// A shaft card lists every one of its cars.
func TestInspectShaftListsCars(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	if err := w.AddCar(20); err != nil {
		t.Fatalf("second car: %v", err)
	}
	ids := w.bank.CarsIn(shaft)
	w.bank.Place(ids[0], 2)
	w.bank.Place(ids[1], 6)

	card, err := w.InspectCell(4, 20)
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	text := cardText(card)
	if !strings.Contains(text, "Elevator") {
		t.Fatalf("not an elevator card:\n%s", text)
	}
	for _, id := range ids {
		if !strings.Contains(text, "Car "+strconv.Itoa(int(id))+":") {
			t.Fatalf("car %d missing from the card:\n%s", id, text)
		}
	}
	if !strings.Contains(text, "floor 2.0") || !strings.Contains(text, "floor 6.0") {
		t.Fatalf("card does not place both cars:\n%s", text)
	}
}

// One car reports its load and heading.
func TestInspectCar(t *testing.T) {
	w, shaft := shaftWorld(t, 8)
	id := carIDOf(w, shaft)
	w.bank.Place(id, 3)
	w.bank.SyncRiders(map[transport.CarID]int{id: 2})

	card, err := w.InspectCar(int(id))
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	text := cardText(card)
	if !strings.Contains(text, "Riders: 2 / 4") {
		t.Fatalf("card does not show the load:\n%s", text)
	}
	if _, err := w.InspectCar(9999); err == nil {
		t.Fatal("inspecting no car should fail")
	}
}

// Rooms and empty cells still answer.
func TestInspectCellRoomAndEmpty(t *testing.T) {
	w, _ := roomWorld(t, "office")
	card, err := w.InspectCell(0, 2)
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	if !strings.Contains(cardText(card), "Office") {
		t.Fatalf("wrong card:\n%s", cardText(card))
	}
	if _, err := w.InspectCell(20, 50); err == nil {
		t.Fatal("empty sky should have nothing to inspect")
	}
}
