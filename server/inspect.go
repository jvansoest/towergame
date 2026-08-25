package main

import (
	"fmt"

	"towergame/server/model"
	"towergame/server/transport"
)

// Describes whatever the player clicked on.

// A room, shaft, or stair at a cell.
func (w *World) InspectCell(floor, col int) (Inspection, error) {
	if floor < 0 || floor >= w.grid.Floors || col < 0 || col >= w.grid.Width {
		return Inspection{}, fmt.Errorf("cell out of range")
	}
	if i := w.grid.RoomAt(floor, col); i >= 0 {
		return w.roomCard(w.grid.Rooms[i]), nil
	}
	if e, ok := w.grid.ElevatorCell(floor, col); ok {
		return w.shaftCard(e), nil
	}
	if s, ok := w.grid.StairCell(floor, col); ok {
		return Inspection{
			Type:  "inspect",
			Title: "Stairs",
			Lines: []string{
				fmt.Sprintf("Floors %d to %d", s.Floor, s.Floor+model.StairRise),
			},
		}, nil
	}
	return Inspection{}, fmt.Errorf("nothing here")
}

// A room and who is in it.
func (w *World) roomCard(r model.Room) Inspection {
	rt := model.RoomTypes[r.Type]
	inside := 0
	for _, p := range w.sims {
		if p.state == statePresent && p.homeF == r.Floor &&
			p.homeC >= r.Col && p.homeC < r.Col+rt.Width {
			inside++
		}
	}
	lines := []string{fmt.Sprintf("Floor %d", r.Floor)}
	if rt.Capacity > 0 {
		lines = append(lines, fmt.Sprintf("People inside: %d / %d", inside, rt.Capacity))
	} else {
		lines = append(lines, fmt.Sprintf("People inside: %d", inside))
	}
	if rt.Seats > 0 {
		lines = append(lines, fmt.Sprintf("Seats: %d", rt.Seats))
	}
	return Inspection{Type: "inspect", Title: rt.Name, Lines: lines}
}

// A shaft and the cars running in it.
func (w *World) shaftCard(e model.Elevator) Inspection {
	cars := w.bank.CarsIn(transport.ShaftID(e.ID))
	waiting := 0
	for _, p := range w.sims {
		if p.state == stateWaiting && p.shaft == transport.ShaftID(e.ID) {
			waiting++
		}
	}
	lines := []string{
		fmt.Sprintf("Serves floors %d to %d", e.Bottom, e.Top),
		fmt.Sprintf("People waiting: %d", waiting),
	}
	// One line per car, in the order they were added.
	for _, id := range cars {
		c, ok := w.bank.Car(id)
		if !ok {
			continue
		}
		lines = append(lines, fmt.Sprintf("Car %d: floor %.1f, %d aboard, going %s",
			id, c.Floor, c.Riders(), headingName(c.Dir)))
	}
	if len(cars) == 0 {
		lines = append(lines, "No cars yet")
	}
	return Inspection{Type: "inspect", Title: "Elevator", Lines: lines}
}

// One car: where it is and what it carries.
func (w *World) InspectCar(id int) (Inspection, error) {
	c, ok := w.bank.Car(transport.CarID(id))
	if !ok {
		return Inspection{}, fmt.Errorf("no such car")
	}
	lines := []string{
		fmt.Sprintf("Floor %.1f", c.Floor),
		fmt.Sprintf("Riders: %d / %d", c.Riders(), transport.Capacity),
		"Going " + headingName(c.Dir),
	}
	if c.Dwell() > 0 {
		lines = append(lines, "Doors open")
	}
	return Inspection{
		Type:  "inspect",
		Title: fmt.Sprintf("Elevator car %d", id),
		Lines: lines,
	}, nil
}

// One sim: where from, where to, doing what.
func (w *World) InspectSim(id int) (Inspection, error) {
	for _, p := range w.sims {
		if p.id != id {
			continue
		}
		lines := []string{
			"From: " + w.placeName(p.from),
			"To: " + w.placeName(p.last),
			"Now: " + w.simDoing(p),
			fmt.Sprintf("Stress: %d%% (%s)", int(p.stress), stressName(p.stress)),
		}
		if p.state == stateWaiting || p.state == stateRiding {
			lines = append(lines, fmt.Sprintf("Elevator to floor %d", p.exitFloor))
		}
		return Inspection{Type: "inspect", Title: kindName(p), Lines: lines}, nil
	}
	return Inspection{}, fmt.Errorf("no such sim")
}

// Names the room at a goal, or the street.
func (w *World) placeName(g goal) string {
	if !g.present {
		return "Outside"
	}
	if i := w.grid.RoomAt(g.floor, g.col); i >= 0 {
		name := model.RoomTypes[w.grid.Rooms[i].Type].Name
		return fmt.Sprintf("%s, floor %d", name, g.floor)
	}
	return fmt.Sprintf("Floor %d", g.floor)
}

// What kind of sim this is.
func kindName(p *sim) string {
	if p.transient {
		return "Customer"
	}
	switch p.category {
	case model.CategoryOffice:
		return "Office worker"
	case model.CategoryResidential:
		return "Resident"
	case model.CategoryHotel:
		return "Hotel guest"
	}
	return "Visitor"
}

// How much elevator waiting has worn a sim down.
func stressName(stress float64) string {
	switch {
	case stress >= 70:
		return "furious"
	case stress >= 40:
		return "annoyed"
	case stress >= 15:
		return "impatient"
	}
	return "calm"
}

// What the sim is doing right now.
func (w *World) simDoing(p *sim) string {
	switch p.state {
	case stateOutside:
		return "away"
	case stateMoving:
		return "walking"
	case stateWaiting:
		return "waiting for an elevator"
	case stateBoarding:
		return "getting in"
	case stateRiding:
		return "riding"
	}
	if p.asleep {
		return "asleep"
	}
	if w.seated(p) {
		return "sitting down"
	}
	if r, ok := w.roomOf(p); ok &&
		model.RoomTypes[r.Type].Habit == model.HabitBrowse {
		return "looking around"
	}
	return "standing"
}

// Which way a car is travelling.
func headingName(dir int) string {
	switch {
	case dir > 0:
		return "up"
	case dir < 0:
		return "down"
	}
	return "nowhere"
}
