package main

import "towergame/server/model"

// Room lighting hours, around the norm.
const (
	officeOn    = 9 * 60
	officeOff   = 17 * 60
	bedtime     = 23 * 60
	wakeTime    = 6*60 + 30
	dayOfMins   = 24 * 60
	lightSpread = 45 // minutes either side
	lateSpread  = 90 // some offices work late
)

// A stable per-room offset in minutes.
// Rooms keep their habits between restarts.
func roomOffset(r model.Room, salt, span int) int {
	h := r.Col*7 + r.Floor*13 + salt*31
	return h%(2*span+1) - span
}

// Whether a room's lights are on now.
// A room with sims in it stays lit.
func lit(r model.Room, weekday bool, mod int, busy bool) bool {
	switch model.RoomTypes[r.Type].Category {
	case model.CategoryOffice:
		on := officeOn + roomOffset(r, 1, lightSpread)
		off := officeOff + roomOffset(r, 2, lateSpread)
		return busy || (weekday && mod >= on && mod < off)

	case model.CategoryResidential, model.CategoryHotel:
		// People here sleep, so busy keeps no light on.
		// The dark window wraps midnight.
		out := (bedtime + roomOffset(r, 3, lightSpread)) % dayOfMins
		up := (wakeTime + roomOffset(r, 4, lightSpread)) % dayOfMins
		return !(mod >= out || mod < up)

	case model.CategoryRetail, model.CategoryFood, model.CategoryMedical:
		// Lit while the doors are open.
		if model.RoomTypes[r.Type].Late {
			return busy || isOpen(r, mod)
		}
		open := mod >= openMinute(r) && mod < visitorClose+roomOffset(r, 5, 30)
		return busy || open
	}
	return true
}
