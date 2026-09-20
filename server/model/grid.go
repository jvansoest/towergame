package model

import "fmt"

// The tower grid, measured in cells.
type Grid struct {
	Width     int      // cells per floor
	Floors    int      // floor count, from the ground up
	Built     [][]bool // [floor][col] base structure
	Below     [][]bool // [depth-1][col], floors -1, -2, ...
	Rooms     []Room
	Stairs    []Stair
	Elevators []Elevator
	Ramps     []Ramp

	nextShaft ShaftID // ids stay stable as shafts come and go
}

// Floors dug below the ground, numbered -1 downward.
const Basement = 6

func NewGrid(width, floors int) *Grid {
	built := make([][]bool, floors)
	for i := range built {
		built[i] = make([]bool, width)
	}
	below := make([][]bool, Basement)
	for i := range below {
		below[i] = make([]bool, width)
	}
	return &Grid{
		Width:     width,
		Floors:    floors,
		Built:     built,
		Below:     below,
		Rooms:     []Room{},
		Stairs:    []Stair{},
		Elevators: []Elevator{},
	}
}

func (g *Grid) inRange(floor, col int) bool {
	return floor >= -Basement && floor < g.Floors && col >= 0 && col < g.Width
}

// Whether a cell is on the grid, above or below ground.
func (g *Grid) InRange(floor, col int) bool { return g.inRange(floor, col) }

// Whether a cell has base structure. Off the grid: no.
func (g *Grid) IsBuilt(floor, col int) bool {
	if !g.inRange(floor, col) {
		return false
	}
	if floor < 0 {
		return g.Below[-floor-1][col]
	}
	return g.Built[floor][col]
}

// Sets a cell built or not; the caller checks the rules.
func (g *Grid) SetBuilt(floor, col int, v bool) { g.setBuilt(floor, col, v) }

func (g *Grid) setBuilt(floor, col int, v bool) {
	if floor < 0 {
		g.Below[-floor-1][col] = v
		return
	}
	g.Built[floor][col] = v
}

// A cell needs ground or support beside its access.
// Above ground it rests on the floor below. Underground
// it is dug from the surface, or from the floor above.
func (g *Grid) supported(floor, col int) bool {
	switch {
	case floor == 0 || floor == -1:
		return true
	case floor < 0:
		return g.IsBuilt(floor+1, col)
	}
	return g.Built[floor-1][col]
}

// Builds one base cell.
func (g *Grid) BuildBase(floor, col int) error {
	if !g.inRange(floor, col) {
		return fmt.Errorf("cell (%d,%d) out of range", floor, col)
	}
	if g.IsBuilt(floor, col) {
		return fmt.Errorf("cell (%d,%d) already built", floor, col)
	}
	if !g.supported(floor, col) {
		return fmt.Errorf("no support below (%d,%d)", floor, col)
	}
	g.setBuilt(floor, col, true)
	return nil
}

// Reports whether a room fits.
func (g *Grid) CanPlace(typeID string, floor, col int) error {
	rt, ok := RoomTypes[typeID]
	if !ok {
		return fmt.Errorf("unknown room type %q", typeID)
	}
	if floor < -Basement || floor+rt.Height > g.Floors {
		return fmt.Errorf("floor %d out of range", floor)
	}
	if floor < 0 && floor+rt.Height > 0 {
		return fmt.Errorf("a room may not cross the ground")
	}
	if rt.Underground && floor >= 0 {
		return fmt.Errorf("%s goes underground", rt.Name)
	}
	if rt.Deepest && floor != -Basement {
		return fmt.Errorf("%s goes on the lowest floor, %d", rt.Name, -Basement)
	}
	if rt.Unique {
		for _, r := range g.Rooms {
			if r.Type == typeID {
				return fmt.Errorf("only one %s", rt.Name)
			}
		}
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
	for _, rp := range g.Ramps {
		if overlaps(col, rt.Width, floor, rt.Height, rp.Col, RampWidth, rp.Floor, 1) {
			return fmt.Errorf("overlaps a car ramp")
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
// A new room starts politically neutral.
func (g *Grid) Place(typeID string, floor, col int) error {
	if err := g.CanPlace(typeID, floor, col); err != nil {
		return err
	}
	rt := RoomTypes[typeID]
	for f := floor; f < floor+rt.Height; f++ {
		for c := col; c < col+rt.Width; c++ {
			g.setBuilt(f, c, true)
		}
	}
	g.Rooms = append(g.Rooms, Room{
		Type: typeID, Floor: floor, Col: col, Alignment: AlignNeutral,
	})
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
