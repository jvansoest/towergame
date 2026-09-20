package main

import (
	"testing"

	"towergame/server/model"
)

// Runs the tower for whole quarters, returns reports.
func runQuarters(t *testing.T, w *World, n int) []Report {
	t.Helper()
	var out []Report
	for i := 0; i < n*quarterDays*180*30+10 && len(out) < n; i++ {
		w.Step(1.0 / 30)
		if r := w.TakeReport(); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

// A quarter's report adds up its own lines.
func TestQuarterAddsUp(t *testing.T) {
	w := newWorld()
	w.seed()
	reps := runQuarters(t, w, 2)
	if len(reps) != 2 {
		t.Fatalf("got %d reports, want 2", len(reps))
	}
	for _, r := range reps {
		if r.Net != r.Hotel+r.Sales+r.Leases+r.Events-r.Upkeep {
			t.Fatalf("report does not add up: %+v", r)
		}
		if r.Upkeep <= 0 {
			t.Fatalf("a tower with rooms costs nothing: %+v", r)
		}
		if r.Net <= 0 {
			t.Fatalf("the starter tower loses money: %+v", r)
		}
		t.Logf("quarter %d: hotel %d sales %d leases %d upkeep %d net %d",
			r.Quarter, r.Hotel, r.Sales, r.Leases, r.Upkeep, r.Net)
	}
}

// Customers pay for what they buy.
func TestCustomersPay(t *testing.T) {
	w, r := roomWorld(t, "shop")
	before := w.money
	_ = r
	for i := 0; i < 900; i++ {
		w.Step(1.0 / 30)
	}
	if w.money <= before {
		t.Fatal("the shop earned nothing in half a minute")
	}
	if w.books.sales != w.money-before {
		t.Fatalf("ledger %d, cash gain %d", w.books.sales, w.money-before)
	}
}

// Furious tenants pay no lease.
func TestFuriousTenantsPayNothing(t *testing.T) {
	w, _ := roomWorld(t, "office")
	room := w.grid.Rooms[0]
	if got := w.leaseShare(room); got != 1 {
		t.Fatalf("calm tenants pay %.1f, want 1", got)
	}
	for _, p := range w.sims {
		p.peak = 50
	}
	if got := w.leaseShare(room); got != 0.5 {
		t.Fatalf("annoyed tenants pay %.1f, want 0.5", got)
	}
	for _, p := range w.sims {
		p.peak = 80
	}
	if got := w.leaseShare(room); got != 0 {
		t.Fatalf("furious tenants pay %.1f, want 0", got)
	}
	_ = model.CategoryOffice
}
