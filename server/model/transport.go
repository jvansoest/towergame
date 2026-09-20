package model

import "fmt"

// A stair rises one floor and lands on the next.
const (
	StairWidth  = 4
	StairHeight = 2 // cells reserved, including the landing
	StairRise   = 1 // cells the run itself fills
)

// A staircase connecting two floors.
type Stair struct {
	Floor     int  `json:"floor"` // bottom floor
	Col       int  `json:"col"`
	Escalator bool `json:"escalator,omitempty"` // moves both ways by itself
}

// Center column (cell-center), used by pathfinding.
func (s Stair) CenterCol() float64 {
	return float64(s.Col) + float64(StairWidth)/2 - 0.5
}

// Bottom of the run, where the climb starts.
func (s Stair) FootCol() float64 {
	return float64(s.Col)
}

// Top of the run, on the floor above.
func (s Stair) HeadCol() float64 {
	return float64(s.Col + StairWidth - 1)
}

// Places a stair after validation.
func (g *Grid) PlaceStair(floor, col int) error {
	return g.placeRun(floor, col, false)
}

// Places an escalator; it takes the same room as a stair.
func (g *Grid) PlaceEscalator(floor, col int) error {
	return g.placeRun(floor, col, true)
}

func (g *Grid) placeRun(floor, col int, escalator bool) error {
	if floor < -Basement || floor+StairHeight > g.Floors {
		return fmt.Errorf("floor %d out of range", floor)
	}
	if col < 0 || col+StairWidth > g.Width {
		return fmt.Errorf("column %d out of range", col)
	}
	for c := col; c < col+StairWidth; c++ {
		if !g.supported(floor, c) {
			return fmt.Errorf("no support below column %d", c)
		}
	}
	// Stairs sit in front of rooms, so rooms behind them are allowed.
	for _, s := range g.Stairs {
		if overlaps(col, StairWidth, floor, StairRise, s.Col, StairWidth, s.Floor, StairRise) {
			return fmt.Errorf("overlaps a stair")
		}
	}
	for _, rp := range g.Ramps {
		if overlaps(col, StairWidth, floor, StairHeight, rp.Col, RampWidth, rp.Floor, 1) {
			return fmt.Errorf("overlaps a car ramp")
		}
	}
	for f := floor; f < floor+StairHeight; f++ {
		for c := col; c < col+StairWidth; c++ {
			g.setBuilt(f, c, true)
		}
	}
	g.Stairs = append(g.Stairs, Stair{Floor: floor, Col: col, Escalator: escalator})
	return nil
}

// An elevator serves floors 0..Top at its column.
const ElevatorWidth = 2

// A shaft's stable name. Never an index: shafts get removed,
// and everything holding a reference must survive that.
type ShaftID int

// A vertical elevator shaft from the ground up.
type Elevator struct {
	ID     ShaftID `json:"id"`
	Bottom int     `json:"bottom"`
	Top    int     `json:"top"`
	Col    int     `json:"col"`
}

// Center column (cell-center), used by pathfinding.
func (e Elevator) CenterCol() float64 {
	return float64(e.Col) + float64(ElevatorWidth)/2 - 0.5
}

// Index of the elevator covering a column, or -1.
func (g *Grid) ElevatorAt(col int) int {
	for i := range g.Elevators {
		e := g.Elevators[i]
		if col >= e.Col && col < e.Col+ElevatorWidth {
			return i
		}
	}
	return -1
}

// The shaft with this id, or false.
func (g *Grid) Shaft(id ShaftID) (Elevator, bool) {
	for _, e := range g.Elevators {
		if e.ID == id {
			return e, true
		}
	}
	return Elevator{}, false
}

// Places or extends an elevator to topFloor. Reports extension.
func (g *Grid) PlaceElevator(topFloor, col int) (bool, error) {
	top := topFloor
	if top < 1 {
		top = 1
	}
	if top >= g.Floors {
		return false, fmt.Errorf("elevator too tall")
	}
	if col < 0 || col+ElevatorWidth > g.Width {
		return false, fmt.Errorf("column %d out of range", col)
	}

	// Extend an existing shaft at this column.
	if idx := g.ElevatorAt(col); idx >= 0 {
		e := &g.Elevators[idx]
		if top <= e.Top {
			return false, fmt.Errorf("elevator already reaches floor %d", e.Top)
		}
		g.buildColumn(e.Col, e.Top+1, top)
		e.Top = top
		return true, nil
	}

	h := top + 1 // floors 0..top
	// Elevators sit in front of rooms, so rooms behind them are allowed.
	for _, s := range g.Stairs {
		if overlaps(col, ElevatorWidth, 0, h, s.Col, StairWidth, s.Floor, StairRise) {
			return false, fmt.Errorf("overlaps a stair")
		}
	}
	for _, e := range g.Elevators {
		if overlaps(col, ElevatorWidth, 0, h, e.Col, ElevatorWidth, e.Bottom, e.Top-e.Bottom+1) {
			return false, fmt.Errorf("overlaps an elevator")
		}
	}
	g.buildColumn(col, 0, top)
	g.nextShaft++
	g.Elevators = append(g.Elevators, Elevator{
		ID: g.nextShaft, Bottom: 0, Top: top, Col: col,
	})
	return false, nil
}

// Marks the shaft cells built between two floors.
func (g *Grid) buildColumn(col, from, to int) {
	for f := from; f <= to; f++ {
		for c := col; c < col+ElevatorWidth; c++ {
			g.Built[f][c] = true
		}
	}
}

// The elevator covering a cell, if any.
func (g *Grid) ElevatorCell(floor, col int) (Elevator, bool) {
	for _, e := range g.Elevators {
		if col >= e.Col && col < e.Col+ElevatorWidth &&
			floor >= e.Bottom && floor <= e.Top {
			return e, true
		}
	}
	return Elevator{}, false
}

// The stair covering a cell, if any.
func (g *Grid) StairCell(floor, col int) (Stair, bool) {
	for _, s := range g.Stairs {
		if col >= s.Col && col < s.Col+StairWidth &&
			floor >= s.Floor && floor < s.Floor+StairHeight {
			return s, true
		}
	}
	return Stair{}, false
}

// Removes the elevator at a cell. Reports the id removed.
func (g *Grid) RemoveElevatorAt(floor, col int) (ShaftID, bool) {
	for i := range g.Elevators {
		e := g.Elevators[i]
		if col >= e.Col && col < e.Col+ElevatorWidth && floor >= e.Bottom && floor <= e.Top {
			g.Elevators = append(g.Elevators[:i], g.Elevators[i+1:]...)
			return e.ID, true
		}
	}
	return 0, false
}

// Removes the stair covering a cell, so
// either row of the run counts as a hit.
func (g *Grid) RemoveStairAt(floor, col int) bool {
	_, ok := g.RemoveRunAt(floor, col)
	return ok
}

// Same, but says which run went.
func (g *Grid) RemoveRunAt(floor, col int) (Stair, bool) {
	for i := range g.Stairs {
		s := g.Stairs[i]
		if col >= s.Col && col < s.Col+StairWidth &&
			floor >= s.Floor && floor < s.Floor+StairHeight {
			g.Stairs = append(g.Stairs[:i], g.Stairs[i+1:]...)
			return s, true
		}
	}
	return Stair{}, false
}
