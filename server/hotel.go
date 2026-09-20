package main

import (
	"math/rand"

	"towergame/server/model"
)

// Renting rooms, and cleaning them after.
const (
	checkInFrom  = 16 * 60
	checkInTo    = 22 * 60
	rentChance   = 0.025 // per room, per check
	guestsComing = 3     // parties on the way at once
	checkOutAt   = 7 * 60
	cleanSeconds = 6.0 // a maid's time per room
)

// Books free, clean rooms for guests, now and then.
func (w *World) manageGuests(mod int) {
	if w.tick%visitorEvery != 0 || mod < checkInFrom || mod >= checkInTo {
		return
	}
	coming := 0
	for _, p := range w.sims {
		if p.transient && p.category == model.CategoryHotel && p.leave == 0 {
			coming++
		}
	}
	for j := range w.grid.Rooms {
		if coming >= guestsComing {
			return
		}
		r := &w.grid.Rooms[j]
		if model.RoomTypes[r.Type].Category != model.CategoryHotel ||
			r.Booked || r.Dirty || rand.Float64() > rentChance {
			continue
		}
		w.bookRoom(j)
		coming++
	}
}

// Sends a party of guests to a room.
func (w *World) bookRoom(j int) {
	r := &w.grid.Rooms[j]
	rt := model.RoomTypes[r.Type]
	r.Booked = true
	for i := 0; i < rt.Capacity; i++ {
		w.nextID++
		w.sims = append(w.sims, &sim{
			id:        w.nextID,
			category:  rt.Category,
			prof:      profVisitor,
			pace:      1,
			homeF:     r.Floor,
			homeC:     model.SeatCol(*r, i),
			spot:      i,
			transient: true,
			state:     stateOutside,
			bornAt:    w.simTime,
		})
	}
}

// Seconds until a guest's morning checkout.
func (w *World) untilMorning() float64 {
	clk := clockFromSimTime(w.simTime)
	mod := clk.Hour*60 + clk.Minute
	target := checkOutAt + rand.Intn(76) - 15
	delta := (target - mod + minutesPerDay) % minutesPerDay
	return float64(delta) / gameMinutesPerRealSecond
}

// Seconds a guest stays. A VIP is brief.
func (w *World) stayLength(p *sim) float64 {
	if p.prof == profVIP {
		return vipStay
	}
	return w.untilMorning()
}

// Notes that a guest has used the room.
func (w *World) markUsed(p *sim) {
	if j := w.grid.RoomAt(p.homeF, p.homeC); j >= 0 {
		w.grid.Rooms[j].Used = true
	}
}

// Rent is paid at checkout, unless noise won.
func (w *World) checkOut(p *sim) {
	if p.prof == profVIP {
		w.vipVerdict(p)
	}
	if p.category != model.CategoryHotel || !p.stayed || p.disturbed {
		return
	}
	if j := w.grid.RoomAt(p.homeF, p.homeC); j >= 0 {
		rt := model.RoomTypes[w.grid.Rooms[j].Type]
		pay := rt.Rent / max(rt.Capacity, 1)
		w.money += pay
		w.books.hotel += pay
	}
}

// Frees rooms whose guests have all gone. A room
// that was slept in is dirty until a maid comes.
func (w *World) releaseGuests() {
	live := map[int]bool{}
	for _, p := range w.sims {
		if p.transient && p.category == model.CategoryHotel {
			live[w.grid.RoomAt(p.homeF, p.homeC)] = true
		}
	}
	for j := range w.grid.Rooms {
		r := &w.grid.Rooms[j]
		if !r.Booked || live[j] {
			continue
		}
		r.Booked = false
		r.Dirty = r.Used
		r.Used = false
	}
}

// Maids pick a dirty room, walk there, and scrub.
func (w *World) stepMaids(dt float64) {
	clk := clockFromSimTime(w.simTime)
	weekday := isWeekday(clk.Day)
	mod := clk.Hour*60 + clk.Minute
	for _, p := range w.sims {
		if p.prof != profMaid {
			continue
		}
		w.checkJob(p)
		if !presentNow(p.category, weekday, mod, p.shiftIn, p.shiftOut) {
			w.dropJob(p)
			continue
		}
		if !p.job.present {
			w.claimRoom(p)
			continue
		}
		if !w.sweeping(p) {
			continue
		}
		j := w.grid.RoomAt(p.job.floor, p.job.col)
		p.cleanLeft -= dt
		if p.cleanLeft <= 0 {
			w.grid.Rooms[j].Dirty = false
			w.grid.Rooms[j].Cleaner = 0
			p.job = goal{}
		}
	}
}

// Whether a maid is at her dirty room.
func (w *World) sweeping(p *sim) bool {
	if p.prof != profMaid || !p.job.present || p.state != statePresent {
		return false
	}
	j := w.grid.RoomAt(p.job.floor, p.job.col)
	return j >= 0 && w.grid.RoomAt(int(p.y+0.5), int(p.x+0.5)) == j
}

// Drops a job whose room is gone or taken.
func (w *World) checkJob(p *sim) {
	if !p.job.present {
		return
	}
	j := w.grid.RoomAt(p.job.floor, p.job.col)
	if j < 0 || !w.grid.Rooms[j].Dirty || w.grid.Rooms[j].Cleaner != p.id {
		p.job = goal{}
	}
}

// Gives up the job, so another maid can take it.
func (w *World) dropJob(p *sim) {
	if !p.job.present {
		return
	}
	if j := w.grid.RoomAt(p.job.floor, p.job.col); j >= 0 &&
		w.grid.Rooms[j].Cleaner == p.id {
		w.grid.Rooms[j].Cleaner = 0
	}
	p.job = goal{}
}

// Takes the nearest dirty room nobody else has.
func (w *World) claimRoom(p *sim) {
	best, bestCost := -1, 0
	for j, r := range w.grid.Rooms {
		if !r.Dirty || r.Booked || (r.Cleaner != 0 && w.simAlive(r.Cleaner)) {
			continue
		}
		cost := abs(r.Floor-int(p.y+0.5))*10 + abs(r.Col-int(p.x+0.5))
		if best < 0 || cost < bestCost {
			best, bestCost = j, cost
		}
	}
	if best < 0 {
		return
	}
	r := &w.grid.Rooms[best]
	r.Cleaner = p.id
	p.job = goal{true, r.Floor, model.SeatCol(*r, 0)}
	p.cleanLeft = cleanSeconds
}

// Whether a sim id is still in the tower.
func (w *World) simAlive(id int) bool {
	for _, p := range w.sims {
		if p.id == id {
			return true
		}
	}
	return false
}
