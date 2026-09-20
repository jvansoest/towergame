package main

import (
	"fmt"

	"towergame/server/model"
	"towergame/server/transport"
)

// Describes whatever the player clicked on.

// A room, shaft, or stair at a cell.
func (w *World) InspectCell(floor, col int) (Inspection, error) {
	if !w.grid.InRange(floor, col) {
		return Inspection{}, fmt.Errorf("cell out of range")
	}
	if i := w.grid.RoomAt(floor, col); i >= 0 {
		return w.roomCard(w.grid.Rooms[i]), nil
	}
	if e, ok := w.grid.ElevatorCell(floor, col); ok {
		return w.shaftCard(e), nil
	}
	if r, ok := w.grid.RampAt(floor, col); ok {
		return Inspection{
			Type:  "inspect",
			Title: "Car ramp",
			Lines: []string{
				fmt.Sprintf("Floor %d, sloping %s", r.Floor, rampSide(r)),
				fmt.Sprintf("Cars in the garage: %d", len(w.vehicles)),
			},
		}, nil
	}
	if s, ok := w.grid.StairCell(floor, col); ok {
		title := "Stairs"
		if s.Escalator {
			title = "Escalator"
		}
		return Inspection{
			Type:  "inspect",
			Title: title,
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
	lines = append(lines, "Alignment: "+string(r.Alignment))
	lines = append(lines, w.stayLines(r, rt)...)
	return Inspection{Type: "inspect", Title: rt.Name, Lines: lines}
}

// Which way a ramp run slants down.
func rampSide(r model.Ramp) string {
	if r.DownRight() {
		return "down to the right"
	}
	return "down to the left"
}

// Rent, cleanliness, and noise lines for a room.
func (w *World) stayLines(r model.Room, rt model.RoomType) []string {
	var lines []string
	if rt.Lease > 0 {
		lines = append(lines, fmt.Sprintf("Lease: $%d per quarter", rt.Lease))
	}
	if rt.Sale > 0 {
		lines = append(lines, fmt.Sprintf("Sale: $%d per customer", rt.Sale))
	}
	lines = append(lines, fmt.Sprintf("Upkeep: $%d per quarter", rt.Cost/upkeepShare))
	if rt.Noise > 0 {
		lines = append(lines, "Noise: "+noiseName(float64(rt.Noise)))
	}
	if rt.Category == model.CategoryHotel {
		state, clean := "vacant", "clean"
		if r.Booked {
			state = "rented"
		}
		if r.Dirty {
			clean = "dirty"
		}
		lines = append(lines, "Status: "+state, "Cleanliness: "+clean,
			fmt.Sprintf("Rent: $%d per guest", rt.Rent))
	}
	if bedroom(rt.Category) {
		if j := w.grid.RoomAt(r.Floor, r.Col); j >= 0 && j < len(w.noise) {
			lines = append(lines, "Noise next door: "+noiseName(w.noise[j]))
		}
	}
	if rt.Category == model.CategoryTransit {
		state := "no train"
		if t := w.train; t != nil {
			state = map[int]string{
				trainArriving: "a train is arriving",
				trainDwelling: "a train is at the platform",
				trainLeaving:  "a train is leaving",
			}[t.phase]
		}
		lines = append(lines, "Now: "+state)
	}
	if rt.Category == model.CategoryParking {
		cars := 0
		for _, v := range w.vehicles {
			if v.floor == r.Floor && v.col == r.Col {
				cars++
			}
		}
		lines = append(lines, fmt.Sprintf("Cars: %d / %d", cars, rt.Slots),
			fmt.Sprintf("Fee: $%d per car", parkFee))
	}
	if rt.Category == model.CategoryService {
		maids := 0
		for _, p := range w.sims {
			if p.prof == profMaid && p.homeF == r.Floor &&
				p.homeC >= r.Col && p.homeC < r.Col+rt.Width {
				maids++
			}
		}
		lines = append(lines, fmt.Sprintf("Maids: %d", maids))
	}
	return lines
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
			"Type: " + kindName(p),
			"Alignment: " + string(simAlign(p)),
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

// A sim's leaning, from its profession.
func simAlign(p *sim) model.Alignment {
	switch p.prof {
	case profTriad:
		return model.AlignTriad
	case profSecurity:
		return model.AlignGood
	}
	return model.AlignNeutral
}

// What kind of sim this is.
func kindName(p *sim) string {
	switch p.prof {
	case profWorker:
		return "Worker"
	case profResident:
		return "Resident"
	case profSecurity:
		return "Security officer"
	case profTriad:
		return "Triad member"
	case profMaid:
		return "Maid"
	case profDoctor:
		return "Doctor"
	case profVIP:
		return "VIP"
	}
	return "Visitor"
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
	if p.job.present {
		return "cleaning a room"
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
