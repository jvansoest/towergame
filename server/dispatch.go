package main

import (
	"towergame/server/model"
	"towergame/server/transport"
)

// Adapts sims to transport: sims in, calls out.

// Which way a waiting sim travels.
func callDir(p *sim) int {
	if p.exitFloor > p.boardFloor {
		return 1
	}
	return -1
}

// The stops every car owes its sims.
// Pickups count only while a car has room.
func (w *World) carCalls() map[transport.CarID][]transport.Call {
	calls := map[transport.CarID][]transport.Call{}
	for _, p := range w.sims {
		if p.shaft == 0 {
			continue
		}
		switch p.state {
		case stateRiding:
			calls[p.car] = append(calls[p.car],
				transport.Call{Floor: float64(p.exitFloor)})
		case stateWaiting:
			if w.bank.HasRoom(p.car) {
				calls[p.car] = append(calls[p.car], transport.Call{
					Floor: float64(p.boardFloor),
					Dir:   callDir(p),
				})
			}
		}
	}
	return calls
}

// Who is aboard each car.
func (w *World) riderCounts() map[transport.CarID]int {
	aboard := map[transport.CarID]int{}
	for _, p := range w.sims {
		if p.holdsCarPlace() {
			aboard[p.car]++
		}
	}
	return aboard
}

// Route cost of a shaft, for pathfinding.
func (w *World) shaftCost(shaft model.ShaftID, fromFloor, toFloor int) float64 {
	return w.bank.ShaftCost(transport.ShaftID(shaft), fromFloor, toFloor)
}
