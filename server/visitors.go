package main

import "towergame/server/model"

// Customers of shops and eating places.
const (
	stayVariance  = 0.6  // share of the stay that varies
	visitStay     = 30.0 // sim-minutes a customer stays
	visitPatience = 20.0 // sim-minutes spent getting there
	inTransitPer  = 1    // customers a shop awaits at once
	visitorsOpen  = 10 * 60
	visitorClose  = 22 * 60
	visitorsPer   = 4  // active customers per shop
	visitorEvery  = 15 // ticks between customer top-ups
)

// Keeps each shop's customers near its target.
// A diner takes a free seat; a shopper takes the floor.
func (w *World) manageVisitors(mod int) {
	// Customers stay half a minute; a slow check does.
	if w.tick%visitorEvery != 0 {
		return
	}
	shops := append(
		w.roomsByCategory(model.CategoryRetail),
		w.roomsByCategory(model.CategoryFood)...,
	)
	for _, r := range shops {
		if mod < openMinute(r) || mod >= visitorClose {
			continue
		}
		rt := model.RoomTypes[r.Type]
		want := visitorsPer
		if rt.Habit == model.HabitSit {
			want = rt.Seats
		}
		taken, coming := w.spotsTaken(r)
		// Wait for the last one to arrive. Poor elevators
		// then hold a shop back, as they should.
		if coming >= inTransitPer {
			continue
		}
		// One at a time, so nobody arrives in a crowd.
		for i := 0; i < want; i++ {
			if taken[i] {
				continue
			}
			if len(taken) < want {
				w.addVisitor(r, i)
			}
			break
		}
	}
}

// The minute a shop opens. Not all on the hour.
func openMinute(r model.Room) int {
	return visitorsOpen + (r.Col*7+r.Floor*13)%61 - 30
}

// Which spots a room's customers hold, and how
// many of them are still on their way there.
func (w *World) spotsTaken(r model.Room) (map[int]bool, int) {
	rt := model.RoomTypes[r.Type]
	taken := map[int]bool{}
	coming := 0
	for _, p := range w.sims {
		if !p.transient || p.homeF != r.Floor {
			continue
		}
		if p.homeC < r.Col || p.homeC >= r.Col+rt.Width {
			continue
		}
		taken[p.spot] = true
		if p.leave == 0 {
			coming++
		}
	}
	return taken, coming
}

// Adds one customer, in the room's own style.
func (w *World) addVisitor(r model.Room, spot int) {
	rt := model.RoomTypes[r.Type]
	col := model.SeatCol(r, spot)
	if rt.Habit == model.HabitBrowse {
		col = r.Col + (rt.Width*(spot+1))/(visitorsPer+1)
	}
	w.nextID++
	w.sims = append(w.sims, &sim{
		id:        w.nextID,
		category:  rt.Category,
		homeF:     r.Floor,
		homeC:     col,
		spot:      spot,
		transient: true,
		state:     stateOutside,
		bornAt:    w.simTime,
	})
}

func (w *World) roomsByCategory(cat model.Category) []model.Room {
	var out []model.Room
	for _, r := range w.grid.Rooms {
		if model.RoomTypes[r.Type].Category == cat {
			out = append(out, r)
		}
	}
	return out
}
