package main

import (
	"math/rand"

	"towergame/server/model"
)

// A VIP visits a clean suite. Happy, they pay well.
const (
	vipMinStars = 2
	vipFrom     = 9 * 60
	vipTo       = 15 * 60
	vipChance   = 0.002  // per check, in the window
	vipStay     = 20.0   // seconds in the suite
	vipReward   = 100000 // a delighted VIP
	vipPenalty  = 20000  // an unhappy one
)

// Now and then, sends a VIP to a free suite.
func (w *World) manageVIP() {
	if w.tick%visitorEvery != 0 || w.stars < vipMinStars {
		return
	}
	clk := clockFromSimTime(w.simTime)
	mod := clk.Hour*60 + clk.Minute
	if mod < vipFrom || mod >= vipTo || w.vipDay == clk.Day || w.vipActive() {
		return
	}
	if rand.Float64() > vipChance {
		return
	}
	w.vipDay = clk.Day
	j := w.vipSuite()
	if j < 0 {
		w.PostMessage("A VIP wanted a suite, but none was ready.")
		return
	}
	w.PostMessage("A VIP is on the way to a suite!")
	w.bookVIP(j)
}

// Whether a VIP is in the tower.
func (w *World) vipActive() bool {
	for _, p := range w.sims {
		if p.prof == profVIP {
			return true
		}
	}
	return false
}

// A clean, empty suite, or -1.
func (w *World) vipSuite() int {
	for j, r := range w.grid.Rooms {
		if r.Type == "hotel_suite" && !r.Booked && !r.Dirty {
			return j
		}
	}
	return -1
}

// Sends the VIP alone to a suite.
func (w *World) bookVIP(j int) {
	r := &w.grid.Rooms[j]
	r.Booked = true
	w.nextID++
	w.sims = append(w.sims, &sim{
		id:        w.nextID,
		category:  model.CategoryHotel,
		prof:      profVIP,
		pace:      2, // more patient than most
		homeF:     r.Floor,
		homeC:     model.SeatCol(*r, 0),
		transient: true,
		state:     stateOutside,
		bornAt:    w.simTime,
	})
}

// Judges the visit as the VIP leaves.
func (w *World) vipVerdict(p *sim) {
	switch {
	case !p.stayed:
		w.vipMoney(-vipPenalty)
		w.PostMessage("The VIP gave up waiting and left.")
	case p.disturbed || p.peak >= calmStress:
		w.vipMoney(-vipPenalty)
		w.PostMessage("The VIP left unhappy.")
	default:
		w.vipMoney(vipReward)
		w.PostMessage("The VIP was delighted with the tower!")
	}
}

func (w *World) vipMoney(n int) {
	w.money += n
	w.books.events += n
}
