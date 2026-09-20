package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"towergame/server/model"
	"towergame/server/transport"
)

// Grid size in cells.
const (
	gridWidth  = 100
	gridFloors = 40
)

// Starting budget and costs.
const (
	startMoney        = 2_000_000
	baseCost          = 500
	stairCost         = 5000
	escalatorCost     = 25000
	escalatorStars    = 2     // rating needed to build one
	rampCost          = 30000 // one run of a car ramp
	rampStars         = 2
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
	bank     *transport.Bank // elevator cars

	// Room lighting, refreshed with each snapshot.
	litRooms  []bool
	occupants []int
	noise     []float64 // per room, from open neighbours

	vehicles  []*vehicle // cars in the garage
	train     *train     // the subway train, if running
	nextTrain float64    // sim-time of the next one

	books   ledger  // this quarter so far
	quarter int     // quarters closed
	report  *Report // waiting to be sent

	vipDay     int  // the day a VIP last came
	stars      int  // rating; tests start at the top
	profitable bool // the last quarter made money

	tick      uint64
	simTime   float64 // in-game seconds
	lastClock Clock   // detects minute changes
}

func newWorld() *World {
	return &World{
		grid:     model.NewGrid(gridWidth, gridFloors),
		messages: []string{},
		money:    startMoney,
		stars:    maxStars,
		bank:     &transport.Bank{Floors: gridFloors},
	}
}

// Builds a solid, fully-tiled starter tower.
// The block sits in the middle of the grid,
// leaving clear land to build on either side.
func (w *World) seed() {
	g := w.grid
	const (
		top  = 7  // roofline of the starter block
		left = 20 // grass kept clear each side
		band = 60 // starter block width
	)

	// Solid structure for the whole block.
	for f := 0; f <= top; f++ {
		for c := left; c < left+band; c++ {
			g.Built[f][c] = true
		}
	}

	// Each floor tiled edge to edge.
	for c := left; c < left+band; c++ {
		g.Place("lobby", 0, c)
	}
	// A restaurant row fronts floor one.
	g.Place("office", 1, left+0)
	g.Place("restaurant_indian", 1, left+9)
	g.Place("restaurant_french", 1, left+27)
	for _, x := range []int{45, 48, 51, 54, 57} {
		g.Place("hotel_single", 1, left+x)
	}
	for c := left; c < left+band; c += 5 {
		g.Place("condo", 2, c)
	}
	for c := left; c < left+48; c += 12 {
		g.Place("shop", 3, c)
	}
	for i, t := range []string{"burger", "pizza", "noodle", "burger"} {
		g.Place(t, 3, left+48+i*3)
	}
	for c := left; c < left+band; c += 10 {
		g.Place("hotel_suite", 4, c)
	}
	for c := left; c < left+54; c += 3 {
		g.Place("hotel_single", 5, c)
	}
	g.Place("housekeeping", 5, left+54)
	for _, x := range []int{0, 16, 32} {
		g.Place("fastfood", 6, left+x)
	}
	g.Place("izakaya", 6, left+48)
	g.Place("izakaya", 6, left+51)
	g.Place("izakaya", 6, left+54)
	g.Place("icecream", 6, left+57)
	for c := left; c < left+55; c += 5 {
		g.Place("condo", 7, c)
	}
	g.Place("security", 7, left+55)

	// A basement garage: a two-run ramp on the left, opening
	// on the street, parking rooms along both floors.
	for c := left - 18; c < left+30; c++ {
		g.SetBuilt(-1, c, true)
	}
	for c := left - 18; c < left-2; c++ {
		g.SetBuilt(-2, c, true)
	}
	g.PlaceRamp(-1, left-10)
	g.PlaceRamp(-2, left-10)
	for c := left - 2 + 0; c < left+28; c += 2 {
		g.Place("parking", -1, c)
	}
	for c := left - 18; c < left-10; c += 2 {
		g.Place("parking", -2, c)
	}
	g.PlaceEscalator(-1, left+30)

	// A subway station on the lowest floor, shops above it,
	// and escalators all the way down.
	for f := -1; f >= -model.Basement; f-- {
		for c := left + 30; c < left+52; c++ {
			g.SetBuilt(f, c, true)
		}
	}
	for f := -2; f >= -model.Basement; f-- {
		g.PlaceEscalator(f, left+30)
	}
	g.Place("subway", -model.Basement, left+34)
	g.Place("shop", -3, left+34)
	g.Place("icecream", -2, left+34)
	g.Place("izakaya", -2, left+37)
	for i, t := range []string{"burger", "pizza", "noodle", "burger"} {
		g.Place(t, -4, left+34+i*3)
	}

	// Escalators link the lowest floors.
	for f := 0; f < 3; f++ {
		g.PlaceEscalator(f, left+4)
	}

	// Transport sits in front of the rooms.
	for _, col := range []int{left + 16, left + 40} {
		g.PlaceElevator(top, col)
		shaft := g.Elevators[len(g.Elevators)-1].ID
		for i := 0; i < 3; i++ {
			w.addCar(shaft)
		}
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
	w.manageTriads()
	w.updateStars()
	w.stepCars(dt)
	w.stepTrains(dt)
	w.manageVIP()

	c := clockFromSimTime(w.simTime)
	if c == w.lastClock {
		return false
	}
	if c.Day != w.lastClock.Day {
		log.Printf("day %d begins", c.Day)
		if c.Day > 1 && (c.Day-1)%quarterDays == 0 {
			w.closeQuarter()
		}
	}
	w.lastClock = c
	return true
}

// Places a room if affordable.
func (w *World) Place(typeID string, floor, col int, align model.Alignment) error {
	rt, ok := model.RoomTypes[typeID]
	if !ok {
		return fmt.Errorf("unknown room type %q", typeID)
	}
	if !model.ValidAlignment(align) {
		return fmt.Errorf("unknown alignment %q", align)
	}
	if w.stars < rt.Stars {
		return fmt.Errorf("%s needs %s", rt.Name, starText(rt.Stars))
	}
	if w.money < rt.Cost {
		return fmt.Errorf("need $%d, have $%d", rt.Cost, w.money)
	}
	if err := w.grid.Place(typeID, floor, col); err != nil {
		return err
	}
	w.grid.Rooms[len(w.grid.Rooms)-1].Alignment = align
	w.money -= rt.Cost
	w.addSims(w.grid.Rooms[len(w.grid.Rooms)-1])
	return nil
}

// Sets who a room answers to.
func (w *World) SetAlign(floor, col int, align model.Alignment) error {
	if !model.ValidAlignment(align) {
		return fmt.Errorf("unknown alignment %q", align)
	}
	i := w.grid.RoomAt(floor, col)
	if i < 0 {
		return fmt.Errorf("no room here")
	}
	w.grid.Rooms[i].Alignment = align
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

// Places an escalator if affordable.
func (w *World) PlaceEscalator(floor, col int) error {
	if w.stars < escalatorStars {
		return fmt.Errorf("escalators need %s", starText(escalatorStars))
	}
	if w.money < escalatorCost {
		return fmt.Errorf("need $%d, have $%d", escalatorCost, w.money)
	}
	if err := w.grid.PlaceEscalator(floor, col); err != nil {
		return err
	}
	w.money -= escalatorCost
	return nil
}

// Places one run of a car ramp if affordable.
func (w *World) PlaceRamp(floor, col int) error {
	if w.stars < rampStars {
		return fmt.Errorf("car ramps need %s", starText(rampStars))
	}
	if w.money < rampCost {
		return fmt.Errorf("need $%d, have $%d", rampCost, w.money)
	}
	if err := w.grid.PlaceRamp(floor, col); err != nil {
		return err
	}
	w.money -= rampCost
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
	if !w.grid.InRange(floor, col) {
		return fmt.Errorf("cell out of range")
	}
	if r, ok, err := w.grid.RemoveRampAt(floor, col); err != nil {
		return err
	} else if ok {
		w.money += rampCost / 2
		_ = r
		return nil
	}
	if id, ok := w.grid.RemoveElevatorAt(floor, col); ok {
		w.bank.RemoveShaft(transport.ShaftID(id))
		w.dropRidersOf(transport.ShaftID(id))
		w.money += elevatorCost / 2
		return nil
	}
	if s, ok := w.grid.RemoveRunAt(floor, col); ok {
		if s.Escalator {
			w.money += escalatorCost / 2
		} else {
			w.money += stairCost / 2
		}
		return nil
	}
	if i := w.grid.RoomAt(floor, col); i >= 0 &&
		model.RoomTypes[w.grid.Rooms[i].Type].Unique {
		return fmt.Errorf("the %s cannot be demolished",
			model.RoomTypes[w.grid.Rooms[i].Type].Name)
	}
	if r, ok := w.grid.RemoveRoomAt(floor, col); ok {
		w.removeSimsOf(r)
		w.money += model.RoomTypes[r.Type].Cost / 2
		return nil
	}
	if !w.grid.IsBuilt(floor, col) {
		return fmt.Errorf("nothing to remove here")
	}
	// Carried concrete may not be stripped away. Above
	// ground the floor above rests on it; below ground
	// the floor below is reached through it.
	if floor >= 0 && w.grid.IsBuilt(floor+1, col) {
		return fmt.Errorf("structure above needs this floor")
	}
	if floor < 0 && w.grid.IsBuilt(floor-1, col) {
		return fmt.Errorf("structure below needs this floor")
	}
	w.grid.SetBuilt(floor, col, false)
	return nil
}

// Drops each sim whose home was here,
// residents and drop-in customers alike.
func (w *World) removeSimsOf(r model.Room) {
	rt := model.RoomTypes[r.Type]
	kept := w.sims[:0]
	for _, p := range w.sims {
		if p.homeF == r.Floor &&
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

// The buildable catalog, sorted once.
// Room types never change at runtime.
var typeViews = sortedTypes()

func sortedTypes() []model.RoomType {
	out := make([]model.RoomType, 0, len(model.RoomTypes))
	for _, rt := range model.RoomTypes {
		out = append(out, rt)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID < out[b].ID })
	return out
}

// Builds the client render view.
func (w *World) Snapshot() Snapshot {
	w.refreshLights()
	clk := clockFromSimTime(w.simTime)
	mod := clk.Hour*60 + clk.Minute

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
			Alignment: string(r.Alignment),
			// Same hours that gate visitor spawns.
			Open:   isOpen(r, mod),
			Dirty:  r.Dirty,
			Booked: r.Booked,
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
			Floor:     s.Floor,
			Col:       s.Col,
			Width:     model.StairWidth,
			Height:    model.StairRise,
			Escalator: s.Escalator,
		})
	}

	ramps := make([]RampView, 0, len(w.grid.Ramps))
	for _, r := range w.grid.Ramps {
		ramps = append(ramps, RampView{
			Floor: r.Floor, Col: r.Col, Width: model.RampWidth,
			DownRight: r.DownRight(),
		})
	}
	below := make([][]bool, len(w.grid.Below))
	for i, row := range w.grid.Below {
		below[i] = append([]bool(nil), row...)
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
		Type:       "snapshot",
		Tick:       w.tick,
		Clock:      clockFromSimTime(w.simTime),
		Money:      w.money,
		Stars:      w.stars,
		Population: w.population(),
		Grid: GridView{
			Width:     w.grid.Width,
			Floors:    w.grid.Floors,
			Built:     built,
			Rooms:     rooms,
			Stairs:    stairs,
			Ramps:     ramps,
			Below:     below,
			Basement:  model.Basement,
			Elevators: elevators,
			Types:     typeViews,
		},
	}
}

// Full chat log for clients.
func (w *World) ChatUpdate() ChatUpdate {
	chat := make([]string, len(w.messages))
	copy(chat, w.messages)
	return ChatUpdate{Type: "chatUpdate", Chat: chat}
}

// At-seat tolerance, in cells.
const atSeatEps = 0.15

// Whether a sim rests on their chair.
// Strollers about the room count as up.
func (w *World) seated(p *sim) bool {
	if p.state != statePresent {
		return false
	}
	r, ok := w.roomOf(p)
	if !ok {
		return false
	}
	rt := model.RoomTypes[r.Type]
	if rt.Habit != model.HabitSit {
		return false
	}
	// Parked means near their own seat column,
	// and resting rather than mid-stroll. Office
	// workers never stray, so they always count.
	nearSeat := math.Abs(p.x-float64(p.homeC)) <= atSeatEps
	if !nearSeat {
		return false
	}
	if p.category == model.CategoryOffice || rt.Line {
		return true
	}
	return p.pause > 0
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
			Riding:   p.state == stateRiding,
			Sitting:  w.seated(p),
			Spot:     p.spot,
			Cleaning: w.sweeping(p),
			Gliding:  gliding(p),
			Waiting:  p.state == stateWaiting,
			Dark:     p.dark,
			// Sims are drawn by their profession.
			Profession: string(p.prof),
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
	rides := make([]VehicleView, 0, len(w.vehicles))
	for _, v := range w.vehicles {
		rides = append(rides, VehicleView{
			ID: v.id, X: v.x, Y: v.y, Dir: v.dir, Slope: v.slope,
			Parked: v.phase == carParked,
		})
	}
	var tv *TrainView
	if t := w.train; t != nil {
		floor, _ := w.stationSpot()
		tv = &TrainView{Floor: floor, X: t.x, Stopped: t.phase == trainDwelling}
	}
	return SimsUpdate{Type: "sims", Sims: views, Cars: cars, Vehicles: rides, Train: tv}
}
