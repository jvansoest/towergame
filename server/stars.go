package main

import "towergame/server/model"

// The tower's rating, one to five stars.
const (
	maxStars       = 5
	starCheckEvery = 30 // ticks between rating checks
)

// People needed for each star, from the second.
var starPeople = [maxStars]int{0, 60, 150, 300, 500}

// Rooms and results a star also asks for.
const (
	securityStar = 3 // a security office
	medicalStar  = 4 // a medical centre
	profitStar   = 5 // a profitable last quarter
)

// Residents, workers, staff, and hotel guests.
func (w *World) population() int {
	n := 0
	for _, p := range w.sims {
		if !p.transient || p.category == model.CategoryHotel {
			n++
		}
	}
	return n
}

// Whether any room of a type stands.
func (w *World) hasRoom(typeID string) bool {
	for _, r := range w.grid.Rooms {
		if r.Type == typeID {
			return true
		}
	}
	return false
}

// The stars the tower earns now.
func (w *World) earnedStars() int {
	pop := w.population()
	stars := 1
	for s := 2; s <= maxStars; s++ {
		if pop < starPeople[s-1] {
			break
		}
		if s >= securityStar && !w.hasRoom("security") {
			break
		}
		if s >= medicalStar && !w.hasRoom("medical") {
			break
		}
		if s >= profitStar && !w.profitable {
			break
		}
		stars = s
	}
	return stars
}

// Raises the rating when earned. It never falls.
func (w *World) updateStars() {
	if w.tick%starCheckEvery != 0 {
		return
	}
	if s := w.earnedStars(); s > w.stars {
		w.stars = s
		w.PostMessage("The tower reached " + starText(s) + "!")
	}
}

func starText(n int) string {
	if n == 1 {
		return "1 star"
	}
	return string(rune('0'+n)) + " stars"
}

// Starts a real game: one star, earned upward.
func (w *World) startCareer() {
	w.stars = 1
	w.updateStars()
}
