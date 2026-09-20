package main

import (
	"math"

	"towergame/server/model"
)

// Noise from open rooms, felt by sleepers next door.
const (
	noiseRise = 2.0 // stress per second, per noise point
	lateOpen  = 17 * 60
	lateClose = 2 * 60
)

// Whether a shop or diner takes customers now.
func isOpen(r model.Room, mod int) bool {
	if model.RoomTypes[r.Type].Late {
		return mod >= lateOpen || mod < lateClose
	}
	return mod >= openMinute(r) && mod < visitorClose
}

// Noise a room hears from open neighbours.
// Side walls count in full, floors and ceilings half.
func (w *World) refreshNoise(mod int) {
	rooms := w.grid.Rooms
	if w.tick%visitorEvery != 0 && len(w.noise) == len(rooms) {
		return
	}
	w.noise = make([]float64, len(rooms))
	for j, r := range rooms {
		cat := model.RoomTypes[r.Type].Category
		if !bedroom(cat) {
			continue
		}
		rw := model.RoomTypes[r.Type].Width
		for k, o := range rooms {
			ot := model.RoomTypes[o.Type]
			if k == j || ot.Noise == 0 || !isOpen(o, mod) {
				continue
			}
			switch {
			case o.Floor == r.Floor &&
				(o.Col+ot.Width == r.Col || r.Col+rw == o.Col):
				w.noise[j] += float64(ot.Noise)
			case abs(o.Floor-r.Floor) == 1 &&
				o.Col < r.Col+rw && r.Col < o.Col+ot.Width:
				w.noise[j] += float64(ot.Noise) / 2
			}
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Noise a sim's own room suffers.
func (w *World) noiseOf(p *sim) float64 {
	j := w.grid.RoomAt(p.homeF, p.homeC)
	if j < 0 || j >= len(w.noise) {
		return 0
	}
	return w.noise[j]
}

// Sleepers fret in noise; guests give up and leave.
func (w *World) noiseStress(p *sim, dt float64) {
	if !p.asleep || p.state != statePresent {
		return
	}
	n := w.noiseOf(p)
	if n == 0 {
		return
	}
	p.stress = math.Min(stressMax, p.stress+noiseRise*n*dt)
	if p.transient && p.category == model.CategoryHotel &&
		!p.disturbed && p.stress >= angryStress {
		p.disturbed = true
		p.leave = w.simTime
	}
}

// A plain word for a noise level.
func noiseName(n float64) string {
	switch {
	case n >= 3:
		return "loud"
	case n >= 1:
		return "noisy"
	case n > 0:
		return "faint"
	}
	return "quiet"
}
