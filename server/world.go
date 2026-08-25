package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"towergame/server/model"
	"towergame/server/transport"
)

// Grid size in cells.
const (
	gridWidth  = 60
	gridFloors = 40
)

// Starting budget and costs.
const (
	startMoney        = 2_000_000
	baseCost          = 500
	stairCost         = 5000
	elevatorCost      = 100000
	elevatorFloorCost = 8000  // per floor when extending
	elevatorCarCost   = 20000 // an extra car in a shaft
)

// Owns the simulation state.
type World struct {
	grid     *model.Grid
	messages []string
	money    int
	sims     []*sim
	nextID   int
	bank     *transport.Bank         // elevator cars
	dropped  map[transport.CarID]int // riders who got out this tick

	// Room lighting, refreshed with each snapshot.
	litRooms  []bool
	occupants []int

	tick      uint64
	simTime   float64 // in-game seconds
	lastClock Clock   // detects minute changes
}

func newWorld() *World {
	return &World{
		grid:     model.NewGrid(gridWidth, gridFloors),
		messages: []string{},
		money:    startMoney,
		bank:     &transport.Bank{Floors: gridFloors},
		dropped:  map[transport.CarID]int{},
	}
}

// Builds a solid, fully-tiled starter tower for testing.
func (w *World) seed() {
	g := w.grid
	const top = 7

	// Solid structure for the whole block.
	for f := 0; f <= top; f++ {
		for c := 0; c < g.Width; c++ {
			g.Built[f][c] = true
		}
	}

	// Each floor tiled edge to edge.
	for c := 0; c < g.Width; c++ {
		g.Place("lobby", 0, c)
	}
	for _, c := range []int{0, 9, 18, 27} {
		g.Place("office", 1, c)
	}
	g.Place("restaurant", 1, 36)
	for c := 0; c < g.Width; c += 10 {
		g.Place("condo", 2, c)
	}
	for c := 0; c < g.Width; c += 12 {
		g.Place("shop", 3, c)
	}
	for c := 0; c < g.Width; c += 10 {
		g.Place("hotel_suite", 4, c)
	}
	for c := 0; c < g.Width; c += 5 {
		g.Place("hotel_single", 5, c)
	}
	for _, c := range []int{0, 16, 32} {
		g.Place("fastfood", 6, c)
	}
	g.Place("shop", 6, 48)
	for c := 0; c < g.Width; c += 10 {
		g.Place("condo", 7, c)
	}

	// Transport sits in front of the rooms.
	g.PlaceElevator(top, 40)
	shaft := g.Elevators[len(g.Elevators)-1].ID
	for i := 0; i < 3; i++ {
		w.addCar(shaft)
	}

	for _, r := range g.Rooms {
		w.addSims(r)
	}
	w.settle()
}

// Puts sims who are already home at home.
// Otherwise the whole tower walks in at once.
func (w *World) settle() {
	clk := clockFromSimTime(w.simTime)
	weekday := isWeekday(clk.Day)
	mod := clk.Hour*60 + clk.Minute
	for _, p := range w.sims {
		if !presentNow(p.category, weekday, mod, p.shiftIn, p.shiftOut) {
			// Away, and not walking out of the door either.
			p.state = stateOutside
			p.last = goal{false, p.homeF, p.homeC}
			continue
		}
		p.state = statePresent
		p.x, p.y = float64(p.homeC), float64(p.homeF)
		p.last = goal{true, p.homeF, p.homeC}
	}
}

// Gives a shaft another car.
func (w *World) addCar(shaft model.ShaftID) transport.CarID {
	return w.bank.Add(transport.ShaftID(shaft))
}

// Adds a car to the shaft at a column.
func (w *World) AddCar(col int) error {
	idx := w.grid.ElevatorAt(col)
	if idx < 0 {
		return fmt.Errorf("no elevator here")
	}
	e := w.grid.Elevators[idx]
	if len(w.bank.CarsIn(transport.ShaftID(e.ID))) >= transport.MaxPerShaft {
		return fmt.Errorf("this shaft holds no more cars")
	}
	if w.money < elevatorCarCost {
		return fmt.Errorf("need $%d, have $%d", elevatorCarCost, w.money)
	}
	w.money -= elevatorCarCost
	id := w.addCar(e.ID)
	log.Printf("shaft %d: car %d added (%d cars)", e.ID, id,
		len(w.bank.CarsIn(transport.ShaftID(e.ID))))
	return nil
}

// Advances one timestep. Reports a clock change.
func (w *World) Step(dt float64) bool {
	w.tick++
	w.simTime += dt

	w.recountRiders()
	w.stepElevators(dt)
	w.stepSims(dt)

	c := clockFromSimTime(w.simTime)
	if c == w.lastClock {
		return false
	}
	if c.Day != w.lastClock.Day {
		log.Printf("day %d begins", c.Day)
	}
	w.lastClock = c
	return true
}

// Places a room if affordable.
func (w *World) Place(typeID string, floor, col int) error {
	rt, ok := model.RoomTypes[typeID]
	if !ok {
		return fmt.Errorf("unknown room type %q", typeID)
	}
	if w.money < rt.Cost {
		return fmt.Errorf("need $%d, have $%d", rt.Cost, w.money)
	}
	if err := w.grid.Place(typeID, floor, col); err != nil {
		return err
	}
	w.money -= rt.Cost
	w.addSims(model.Room{Type: typeID, Floor: floor, Col: col})
	return nil
}

// Builds one base cell if affordable.
func (w *World) BuildBase(floor, col int) error {
	if w.money < baseCost {
		return fmt.Errorf("need $%d, have $%d", baseCost, w.money)
	}
	if err := w.grid.BuildBase(floor, col); err != nil {
		return err
	}
	w.money -= baseCost
	return nil
}

// Places a stair if affordable.
func (w *World) PlaceStair(floor, col int) error {
	if w.money < stairCost {
		return fmt.Errorf("need $%d, have $%d", stairCost, w.money)
	}
	if err := w.grid.PlaceStair(floor, col); err != nil {
		return err
	}
	w.money -= stairCost
	return nil
}

// Places or extends an elevator if affordable.
func (w *World) PlaceElevator(topFloor, col int) error {
	cost := elevatorCost
	if idx := w.grid.ElevatorAt(col); idx >= 0 {
		cur := w.grid.Elevators[idx].Top
		// A click on the shaft itself adds a car.
		if topFloor <= cur {
			return w.AddCar(col)
		}
		cost = (topFloor - cur) * elevatorFloorCost
	}
	if w.money < cost {
		return fmt.Errorf("need $%d, have $%d", cost, w.money)
	}
	extended, err := w.grid.PlaceElevator(topFloor, col)
	if err != nil {
		return err
	}
	w.money -= cost
	if !extended {
		w.addCar(w.grid.Elevators[len(w.grid.Elevators)-1].ID)
	}
	return nil
}

// Removes whatever is at a cell, with a partial refund.
func (w *World) Remove(floor, col int) error {
	if floor < 0 || floor >= w.grid.Floors || col < 0 || col >= w.grid.Width {
		return fmt.Errorf("cell out of range")
	}
	if id, ok := w.grid.RemoveElevatorAt(floor, col); ok {
		w.bank.RemoveShaft(transport.ShaftID(id))
		w.dropRidersOf(transport.ShaftID(id))
		w.money += elevatorCost / 2
		return nil
	}
	if w.grid.RemoveStairAt(floor, col) {
		w.money += stairCost / 2
		return nil
	}
	if r, ok := w.grid.RemoveRoomAt(floor, col); ok {
		w.removeSimsOf(r)
		w.money += model.RoomTypes[r.Type].Cost / 2
		return nil
	}
	if w.grid.Built[floor][col] {
		w.grid.Built[floor][col] = false
		return nil
	}
	return fmt.Errorf("nothing to remove here")
}

// Drops the residents of a removed room.
func (w *World) removeSimsOf(r model.Room) {
	rt := model.RoomTypes[r.Type]
	kept := w.sims[:0]
	for _, p := range w.sims {
		if !p.transient && p.homeF == r.Floor &&
			p.homeC >= r.Col && p.homeC < r.Col+rt.Width {
			continue
		}
		kept = append(kept, p)
	}
	w.sims = kept
}

// Sends users of a removed shaft walking.
// Ids are stable, so nothing renumbers.
func (w *World) dropRidersOf(shaft transport.ShaftID) {
	for _, p := range w.sims {
		if p.shaft != shaft {
			continue
		}
		p.shaft = 0
		p.state = stateMoving
		p.path = nil
		p.idx = 0
		p.last = goal{}
	}
}

// Hands the cars their stops.
func (w *World) stepElevators(dt float64) {
	w.bank.Step(dt, w.carCalls())
}

// Appends a chat line. Reports acceptance.
func (w *World) PostMessage(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	w.messages = append(w.messages, time.Now().Format("[15:04:05] ")+text)
	return true
}

// Builds the client render view.
func (w *World) Snapshot() Snapshot {
	w.refreshLights()

	rooms := make([]RoomView, 0, len(w.grid.Rooms))
	for j, r := range w.grid.Rooms {
		rt := model.RoomTypes[r.Type]
		rooms = append(rooms, RoomView{
			Type:      r.Type,
			Floor:     r.Floor,
			Col:       r.Col,
			Width:     rt.Width,
			Height:    rt.Height,
			Category:  string(rt.Category),
			Capacity:  rt.Capacity,
			Seats:     rt.Seats,
			Habit:     string(rt.Habit),
			Lit:       w.litRooms[j],
			Occupants: w.occupants[j],
		})
	}

	// Copy Built to avoid races.
	built := make([][]bool, len(w.grid.Built))
	for i, row := range w.grid.Built {
		built[i] = append([]bool(nil), row...)
	}

	stairs := make([]StairView, 0, len(w.grid.Stairs))
	for _, s := range w.grid.Stairs {
		stairs = append(stairs, StairView{
			Floor:  s.Floor,
			Col:    s.Col,
			Width:  model.StairWidth,
			Height: model.StairRise,
		})
	}

	elevators := make([]ElevatorView, 0, len(w.grid.Elevators))
	for _, e := range w.grid.Elevators {
		elevators = append(elevators, ElevatorView{
			ID:     int(e.ID),
			Bottom: e.Bottom,
			Top:    e.Top,
			Col:    e.Col,
			Width:  model.ElevatorWidth,
		})
	}

	return Snapshot{
		Type:  "snapshot",
		Tick:  w.tick,
		Clock: clockFromSimTime(w.simTime),
		Money: w.money,
		Grid: GridView{
			Width:     w.grid.Width,
			Floors:    w.grid.Floors,
			Built:     built,
			Rooms:     rooms,
			Stairs:    stairs,
			Elevators: elevators,
		},
	}
}

// Full chat log for clients.
func (w *World) ChatUpdate() ChatUpdate {
	chat := make([]string, len(w.messages))
	copy(chat, w.messages)
	return ChatUpdate{Type: "chatUpdate", Chat: chat}
}

// Whether a sim is sitting in their room.
func (w *World) seated(p *sim) bool {
	if p.state != statePresent {
		return false
	}
	r, ok := w.roomOf(p)
	if !ok {
		return false
	}
	return model.RoomTypes[r.Type].Habit == model.HabitSit
}

// Recomputes room lights, then who sits in the dark.
func (w *World) refreshLights() {
	clk := clockFromSimTime(w.simTime)
	weekday := isWeekday(clk.Day)
	mod := clk.Hour*60 + clk.Minute

	w.occupants = make([]int, len(w.grid.Rooms))
	for _, p := range w.sims {
		if p.state != statePresent {
			continue
		}
		if j := w.grid.RoomAt(p.homeF, p.homeC); j >= 0 {
			w.occupants[j]++
		}
	}
	w.litRooms = make([]bool, len(w.grid.Rooms))
	for j, r := range w.grid.Rooms {
		w.litRooms[j] = lit(r, weekday, mod, w.occupants[j] > 0)
	}
	for _, p := range w.sims {
		w.markLight(p)
	}
}

// Sets one sim's light state from their room.
func (w *World) markLight(p *sim) {
	p.dark, p.asleep = false, false
	if p.state != statePresent {
		return
	}
	j := w.grid.RoomAt(int(math.Round(p.y)), int(math.Round(p.x)))
	if j < 0 || j >= len(w.litRooms) || w.litRooms[j] {
		return
	}
	p.dark = true
	p.asleep = bedroom(model.RoomTypes[w.grid.Rooms[j].Type].Category)
}

// Whether sims sleep in this kind of room.
func bedroom(cat model.Category) bool {
	return cat == model.CategoryResidential || cat == model.CategoryHotel
}

// Current sim positions for clients.
func (w *World) SimsUpdate() SimsUpdate {
	views := make([]SimView, 0, len(w.sims))
	for _, p := range w.sims {
		// Sleepers are not drawn at all.
		if p.state == stateOutside || p.asleep {
			continue
		}
		views = append(views, SimView{
			ID: p.id, X: p.x, Y: p.y,
			Riding:  p.state == stateRiding,
			Sitting: w.seated(p),
			Dark:    p.dark,
		})
	}

	// One view per car, not per shaft.
	all := w.bank.Cars()
	cars := make([]CarView, 0, len(all))
	for _, c := range all {
		cars = append(cars, CarView{
			ID:    int(c.ID),
			Shaft: int(c.Shaft),
			Col:   w.shaftCol(c.Shaft),
			Floor: c.Floor,
		})
	}
	return SimsUpdate{Type: "sims", Sims: views, Cars: cars}
}
