package main

import (
	"log"
	"math"

	"towergame/server/model"
	"towergame/server/transport"
)

// Income and costs. A quarter is one game week.
const (
	quarterDays      = 3
	upkeepShare      = 20   // room upkeep: 1/20 of cost
	carUpkeep        = 1500 // per elevator car
	shaftFloorUpkeep = 250  // per floor a shaft serves
	stairUpkeep      = 100
	escalatorUpkeep  = 600
	rampUpkeep       = 400
	calmStress       = 40.0 // tenants above this pay half
	furiousStress    = 70.0 // tenants above this pay nothing
)

// Money moved so far this quarter, by kind.
type ledger struct {
	hotel  int // guest rent
	sales  int // customers
	leases int
	events int // VIP rewards and penalties
	upkeep int
}

// Pays for a customer who has just arrived.
func (w *World) sell(p *sim) {
	j := w.grid.RoomAt(p.homeF, p.homeC)
	if j < 0 {
		return
	}
	sale := model.RoomTypes[w.grid.Rooms[j].Type].Sale
	w.money += sale
	w.books.sales += sale
}

// How much of its lease a room pays. Tenants who
// waited too long for the elevator pay less.
func (w *World) leaseShare(r model.Room) float64 {
	rt := model.RoomTypes[r.Type]
	worst := 0.0
	for _, p := range w.sims {
		if p.homeF == r.Floor && p.homeC >= r.Col && p.homeC < r.Col+rt.Width {
			worst = math.Max(worst, p.peak)
		}
	}
	switch {
	case worst >= furiousStress:
		return 0
	case worst >= calmStress:
		return 0.5
	}
	return 1
}

// What the tower costs to keep for one quarter.
func (w *World) upkeep() int {
	total := 0
	for _, r := range w.grid.Rooms {
		total += model.RoomTypes[r.Type].Cost / upkeepShare
	}
	for _, e := range w.grid.Elevators {
		floors := e.Top - e.Bottom + 1
		total += floors*shaftFloorUpkeep +
			len(w.bank.CarsIn(transport.ShaftID(e.ID)))*carUpkeep
	}
	for _, s := range w.grid.Stairs {
		if s.Escalator {
			total += escalatorUpkeep
		} else {
			total += stairUpkeep
		}
	}
	return total + len(w.grid.Ramps)*rampUpkeep
}

// Collects leases, pays upkeep, and files the report.
func (w *World) closeQuarter() {
	for _, r := range w.grid.Rooms {
		if lease := model.RoomTypes[r.Type].Lease; lease > 0 {
			w.books.leases += int(float64(lease) * w.leaseShare(r))
		}
	}
	w.books.upkeep = w.upkeep()
	net := w.books.leases - w.books.upkeep
	w.money += net
	w.quarter++
	w.profitable = net+w.books.hotel+w.books.sales+w.books.events > 0
	w.report = &Report{
		Type:    "report",
		Quarter: w.quarter,
		Hotel:   w.books.hotel,
		Sales:   w.books.sales,
		Leases:  w.books.leases,
		Events:  w.books.events,
		Upkeep:  w.books.upkeep,
		Net:     net + w.books.hotel + w.books.sales + w.books.events,
		Money:   w.money,
	}
	log.Printf("quarter %d: net $%d, money $%d", w.quarter, w.report.Net, w.money)
	w.books = ledger{}
	for _, p := range w.sims {
		p.peak = p.stress
	}
}

// Hands over the finished report, once.
func (w *World) TakeReport() *Report {
	r := w.report
	w.report = nil
	return r
}
