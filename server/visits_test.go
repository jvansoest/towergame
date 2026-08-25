package main

import (
	"testing"

	"towergame/server/model"
)

// A visit starts when the customer arrives.
// Travel time must not eat into the stay.
func TestVisitStartsOnArrival(t *testing.T) {
	w, _ := roomWorld(t, "shop")

	var c *sim
	for i := 0; i < 900 && c == nil; i++ {
		w.Step(1.0 / 30)
		for _, p := range w.sims {
			if p.transient && p.state == statePresent {
				c = p
			}
		}
		// While still walking, no leaving time is set.
		for _, p := range w.sims {
			if p.transient && p.state == stateMoving && p.leave != 0 {
				t.Fatal("the stay started before the customer arrived")
			}
		}
	}
	if c == nil {
		t.Fatal("no customer reached the shop")
	}
	if c.leave <= w.simTime {
		t.Fatalf("customer arrived with %.1f minutes left", c.leave-w.simTime)
	}
}

// A shop waits for the customer already on the way.
// A tower with no elevators must not pile up a crowd.
func TestShopAwaitsOneCustomer(t *testing.T) {
	w := newWorld()
	g := w.grid
	for c := 0; c < 20; c++ {
		g.Built[0][c] = true
		g.Built[1][c] = true
	}
	// A shop upstairs with no way to reach it.
	if err := w.Place("shop", 1, 2); err != nil {
		t.Fatalf("shop: %v", err)
	}
	w.simTime = float64(12*60-startMinuteOfDay) / gameMinutesPerRealSecond

	for i := 0; i < 30*60*3; i++ {
		w.Step(1.0 / 30)
		if len(w.sims) > inTransitPer {
			t.Fatalf("tick %d: %d customers sent to an unreachable shop",
				i, len(w.sims))
		}
	}
}

// Customers give up if they cannot get there.
func TestCustomersGiveUpTravelling(t *testing.T) {
	r := model.Room{Type: "shop", Floor: 1, Col: 2}
	w := newWorld()
	for c := 0; c < 20; c++ {
		w.grid.Built[0][c] = true
		w.grid.Built[1][c] = true
	}
	if err := w.Place(r.Type, r.Floor, r.Col); err != nil {
		t.Fatalf("shop: %v", err)
	}
	w.simTime = float64(12*60-startMinuteOfDay) / gameMinutesPerRealSecond
	w.addVisitor(r, 0)
	p := w.sims[0]

	// Past the patience, they head back out.
	w.simTime += visitPatience + 1
	if g := w.goalFor(p, true, 12*60); g.present {
		t.Fatal("customer still heading for an unreachable shop")
	}
}
