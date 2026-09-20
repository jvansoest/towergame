package main

import (
	"strings"
	"testing"
)

// The last chat line, or "".
func lastChat(w *World) string {
	if len(w.messages) == 0 {
		return ""
	}
	return w.messages[len(w.messages)-1]
}

// Runs until the VIP has gone. Optionally sours the stay.
func runVIP(w *World, sour bool) {
	for i := 0; i < 30*200 && (w.vipActive() || i < 5); i++ {
		w.Step(1.0 / 30)
		if sour {
			for _, p := range w.sims {
				if p.prof == profVIP && p.stayed {
					p.peak = 90
				}
			}
		}
	}
}

// A VIP who is kept calm pays a reward.
func TestVIPDelighted(t *testing.T) {
	w, _ := roomWorld(t, "hotel_suite")
	setClock(w, 10, 0)
	before := w.money
	w.bookVIP(0)
	runVIP(w, false)
	if w.books.events != vipReward {
		t.Fatalf("events %d, want the reward %d", w.books.events, vipReward)
	}
	if w.money-before < vipReward {
		t.Fatalf("money rose %d, want at least %d", w.money-before, vipReward)
	}
	if !strings.Contains(strings.Join(w.messages, "\n"), "delighted") {
		t.Fatalf("no delighted message: %v", w.messages)
	}
	if !w.grid.Rooms[0].Dirty {
		t.Fatal("the suite was not left dirty")
	}
}

// A VIP who waited too long costs money.
func TestVIPUnhappy(t *testing.T) {
	w, _ := roomWorld(t, "hotel_suite")
	setClock(w, 10, 0)
	w.bookVIP(0)
	runVIP(w, true)
	if w.books.events != -vipPenalty {
		t.Fatalf("events %d, want the penalty %d", w.books.events, -vipPenalty)
	}
}

// With no suite ready, the VIP turns away.
func TestVIPNoSuite(t *testing.T) {
	w, _ := roomWorld(t, "hotel_single")
	for i := 0; i < 40000 && w.vipDay == 0; i++ {
		setClock(w, 10, 0)
		w.tick = uint64(i * visitorEvery)
		w.manageVIP()
	}
	if w.vipDay == 0 {
		t.Fatal("no VIP ever came")
	}
	if !strings.Contains(lastChat(w), "none was ready") {
		t.Fatalf("chat says %q", lastChat(w))
	}
	if w.vipActive() {
		t.Fatal("a VIP came with no suite")
	}
}

// One VIP a day, and none in a one-star tower.
func TestVIPRules(t *testing.T) {
	w, _ := roomWorld(t, "hotel_suite")
	w.stars = 1
	for i := 0; i < 20000; i++ {
		setClock(w, 10, 0)
		w.tick = uint64(i * visitorEvery)
		w.manageVIP()
	}
	if w.vipActive() || w.vipDay != 0 {
		t.Fatal("a one-star tower got a VIP")
	}
}
