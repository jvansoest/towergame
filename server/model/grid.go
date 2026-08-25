package model

import "fmt"

// The tower grid, measured in cells.
type Grid struct {
	Width     int      // cells per floor
	Floors    int      // floor count
	Built     [][]bool // [floor][col] base structure
	Rooms     []Room
	Stairs    []Stair
	Elevators []Elevator

	nextShaft ShaftID // ids stay stable as shafts come and go
}

func NewGrid(width, floors int) *Grid {
	built := make([][]bool, floors)
	for i := range built {
		built[i] = make([]bool, width)
	}
	return &Grid{
		Width:     width,
		Floors:    floors,
		Built:     built,
		Rooms:     []Room{},
		Stairs:    []Stair{},
		Elevators: []Elevator{},
	}
}

func (g *Grid) inRange(floor, col int) bool {
	return floor >= 0 && floor < g.Floors && col >= 0 && col < g.Width
}

// A cell needs ground or support below.
func (g *Grid) supported(floor, col int) bool {
	return floor == 0 || g.Built[floor-1][col]
}

// Builds one base cell.
func (g *Grid) BuildBase(floor, col int) error {
	if !g.inRange(floor, col) {
		return fmt.Errorf("cell (%d,%d) out of range", floor, col)
	}
	if g.Built[floor][col] {
		return fmt.Errorf("cell (%d,%d) already built", floor, col)
	}
	if !g.supported(floor, col) {
		return fmt.Errorf("no support below (%d,%d)", floor, col)
	}
	g.Built[floor][col] = true
	return nil
}

// Reports whether a room fits.
func (g *Grid) CanPlace(typeID string, floor, col int) error {
	rt, ok := RoomTypes[typeID]
	if !ok {
		return fmt.Errorf("unknown room type %q", typeID)
	}
	if floor < 0 || floor+rt.Height > g.Floors {
		return fmt.Errorf("floor %d out of range", floor)
	}
	if col < 0 || col+rt.Width > g.Width {
		return fmt.Errorf("column %d out of range", col)
	}
	if rt.Category == CategoryLobby && floor != 0 {
		return fmt.Errorf("lobby only on the ground floor")
	}
	for c := col; c < col+rt.Width; c++ {
		if !g.supported(floor, c) {
			return fmt.Errorf("no support below column %d", c)
		}
	}
	for _, r := range g.Rooms {
		rr := RoomTypes[r.Type]
		if overlaps(col, rt.Width, floor, rt.Height, r.Col, rr.Width, r.Floor, rr.Height) {
			return fmt.Errorf("overlaps a room at (%d,%d)", r.Floor, r.Col)
		}
	}
	return nil
}

// Places a room and its base.
func (g *Grid) Place(typeID string, floor, col int) error {
	if err := g.CanPlace(typeID, floor, col); err != nil {
		return err
	}
	rt := RoomTypes[typeID]
	for f := floor; f < floor+rt.Height; f++ {
		for c := col; c < col+rt.Width; c++ {
			g.Built[f][c] = true
		}
	}
	g.Rooms = append(g.Rooms, Room{Type: typeID, Floor: floor, Col: col})
	return nil
}

// Index of the room at a cell, or -1.
func (g *Grid) RoomAt(floor, col int) int {
	for i := range g.Rooms {
		r := g.Rooms[i]
		rt := RoomTypes[r.Type]
		if col >= r.Col && col < r.Col+rt.Width && floor >= r.Floor && floor < r.Floor+rt.Height {
			return i
		}
	}
	return -1
}

// Removes the room at a cell.
func (g *Grid) RemoveRoomAt(floor, col int) (Room, bool) {
	i := g.RoomAt(floor, col)
	if i < 0 {
		return Room{}, false
	}
	r := g.Rooms[i]
	g.Rooms = append(g.Rooms[:i], g.Rooms[i+1:]...)
	return r, true
}

// Reports whether two boxes overlap.
func overlaps(aCol, aW, aFloor, aH, bCol, bW, bFloor, bH int) bool {
	colsOverlap := aCol < bCol+bW && bCol < aCol+aW
	floorsOverlap := aFloor < bFloor+bH && bFloor < aFloor+aH
	return colsOverlap && floorsOverlap
}
