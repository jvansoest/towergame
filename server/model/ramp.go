package model

import "fmt"

// A car ramp winds down through the basement, one floor
// a run. Each run crosses its own floor's row, slanting
// from the top of the row to the floor, and the next run
// turns back the other way.
const RampWidth = 8

// One run of a car ramp.
type Ramp struct {
	Floor int `json:"floor"` // the row it crosses, below ground
	Col   int `json:"col"`   // leftmost cell
}

// Whether the run slants down toward the right.
func (r Ramp) DownRight() bool { return r.Floor%2 != 0 }

// Column where the run starts, at the top of its row.
func (r Ramp) TopCol() int {
	if r.DownRight() {
		return r.Col
	}
	return r.Col + RampWidth - 1
}

// Column where the run ends, on the floor of its row.
func (r Ramp) BottomCol() int {
	if r.DownRight() {
		return r.Col + RampWidth - 1
	}
	return r.Col
}

// The run crossing a row at a column, if any.
func (g *Grid) RampAt(floor, col int) (Ramp, bool) {
	for _, r := range g.Ramps {
		if r.Floor == floor && col >= r.Col && col < r.Col+RampWidth {
			return r, true
		}
	}
	return Ramp{}, false
}

// Whether a run sits below this one, in the same chain.
func (g *Grid) RampBelow(r Ramp) bool {
	for _, o := range g.Ramps {
		if o.Floor == r.Floor-1 && o.Col == r.Col {
			return true
		}
	}
	return false
}

// Places one run after validation. The first run opens
// onto the street; the rest stack below it.
func (g *Grid) PlaceRamp(floor, col int) error {
	if floor >= 0 || floor < -Basement {
		return fmt.Errorf("a car ramp goes underground")
	}
	if col < 0 || col+RampWidth > g.Width {
		return fmt.Errorf("column %d out of range", col)
	}
	for c := col; c < col+RampWidth; c++ {
		if !g.IsBuilt(floor, c) {
			return fmt.Errorf("dig the floor under the ramp first")
		}
	}
	for _, r := range g.Rooms {
		rt := RoomTypes[r.Type]
		if overlaps(col, RampWidth, floor, 1, r.Col, rt.Width, r.Floor, rt.Height) {
			return fmt.Errorf("overlaps a room")
		}
	}
	for _, s := range g.Stairs {
		if overlaps(col, RampWidth, floor, 1, s.Col, StairWidth, s.Floor, StairHeight) {
			return fmt.Errorf("overlaps a stair")
		}
	}
	for _, o := range g.Ramps {
		if o.Floor == floor && overlaps(col, RampWidth, 0, 1, o.Col, RampWidth, 0, 1) {
			return fmt.Errorf("overlaps a car ramp")
		}
	}
	next := Ramp{Floor: floor, Col: col}
	if floor == -1 {
		if g.Built[0][next.TopCol()] {
			return fmt.Errorf("the ramp must open onto the street")
		}
	} else if _, ok := g.RampAt(floor+1, col); !ok || !g.hasRun(floor+1, col) {
		return fmt.Errorf("continue the ramp below the last run")
	}
	g.Ramps = append(g.Ramps, next)
	return nil
}

// Whether a run of this chain sits on a row.
func (g *Grid) hasRun(floor, col int) bool {
	for _, o := range g.Ramps {
		if o.Floor == floor && o.Col == col {
			return true
		}
	}
	return false
}

// Removes the run at a cell, unless another hangs below it.
func (g *Grid) RemoveRampAt(floor, col int) (Ramp, bool, error) {
	r, ok := g.RampAt(floor, col)
	if !ok {
		return Ramp{}, false, nil
	}
	if g.RampBelow(r) {
		return r, false, fmt.Errorf("remove the ramp below first")
	}
	for i, o := range g.Ramps {
		if o == r {
			g.Ramps = append(g.Ramps[:i], g.Ramps[i+1:]...)
			break
		}
	}
	return r, true, nil
}
