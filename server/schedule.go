package main

import (
	"math/rand"

	"towergame/server/model"
)

// A sim's day, kept loosely. Nobody keeps the clock exactly.
const (
	lunchStart     = 12 * 60
	lunchLen       = 45 // minutes at lunch
	scheduleSpread = 45 // minutes either side of the hour
	lunchSpread    = 40
)

// The target location for a sim right now.
func (w *World) goalFor(p *sim, weekday bool, mod int) goal {
	home := goal{true, p.homeF, p.homeC}
	out := goal{false, p.homeF, p.homeC}

	if p.transient {
		// Still on the way: keep going, but not forever.
		if p.leave == 0 {
			if w.simTime-p.bornAt > visitPatience {
				return out
			}
			return home
		}
		if w.simTime < p.leave {
			return home
		}
		return out
	}
	if !presentNow(p.category, weekday, mod, p.shiftIn, p.shiftOut) {
		return out
	}
	// Lunch runs in sittings, not one rush.
	eat := lunchStart + p.lunchAt
	if p.category == model.CategoryOffice && mod >= eat && mod < eat+lunchLen {
		if f, c, ok := w.lunchSpot(p); ok {
			return goal{true, f, c}
		}
	}
	return home
}

// Whether a resident/worker/guest is in the tower now.
// Each sim keeps their own hours, near the norm.
func presentNow(cat model.Category, weekday bool, mod, in, out int) bool {
	switch cat {
	case model.CategoryOffice:
		return weekday && mod >= 9*60+in && mod < 17*60+out
	case model.CategoryResidential:
		return !(weekday && mod >= 8*60+30+in && mod < 14*60+out)
	case model.CategoryHotel:
		return mod >= 17*60+in || mod < 6*60+30+out
	default:
		return false
	}
}

// A personal offset in minutes, either side.
func spreadMinutes(span int) int {
	return rand.Intn(2*span+1) - span
}

// A stable food room and seat for a worker's lunch.
func (w *World) lunchSpot(p *sim) (floor, col int, ok bool) {
	foods := w.roomsByCategory(model.CategoryFood)
	if len(foods) == 0 {
		return 0, 0, false
	}
	r := foods[p.id%len(foods)]
	return r.Floor, model.SeatCol(r, p.id), true
}
