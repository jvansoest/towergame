package transport

import "math"

// Nearest Car dispatch, scored by suitability.
// See Barney, "Elevator Traffic Handbook".
//
//	idle or approaching, same way -> N + 2 - d
//	approaching, will reverse     -> N + 1 - d
//	moving away                   -> 1
//
// Cost inverts the score into walking cells.
const (
	fsFloorCost = 1.5  // cells per floor between car and call
	fsFullCar   = 20.0 // cells; a full car is a poor bet
)

// Cost of one car answering a call.
func (b *Bank) CarCost(id CarID, fromFloor, toFloor int) float64 {
	c, ok := b.Car(id)
	if !ok {
		return 0
	}
	n := float64(b.Floors)
	d := math.Abs(c.Floor - float64(fromFloor))

	carDir := float64(c.Dir)
	callDir := sign(float64(toFloor - fromFloor))
	// Idle or already here always suits.
	approaching := carDir == 0 || d < LevelEps ||
		sign(float64(fromFloor)-c.Floor) == carDir

	fs := n + 2 - d
	switch {
	case !approaching:
		fs = 1
	case carDir != 0 && callDir != carDir:
		fs = n + 1 - d
	}

	cost := (n + 2 - fs) * fsFloorCost
	if c.riders >= Capacity {
		cost += fsFullCar
	}
	return cost
}

// Best car in a shaft, and its cost.
// The group controller, when a shaft has several.
func (b *Bank) BestCar(shaft ShaftID, fromFloor, toFloor int) (CarID, float64, bool) {
	best, bestCost, found := CarID(0), math.MaxFloat64, false
	for _, c := range b.cars {
		if c.Shaft != shaft {
			continue
		}
		if cost := b.CarCost(c.ID, fromFloor, toFloor); cost < bestCost {
			best, bestCost, found = c.ID, cost, true
		}
	}
	return best, bestCost, found
}

// Route cost of a shaft's best car.
func (b *Bank) ShaftCost(shaft ShaftID, fromFloor, toFloor int) float64 {
	_, cost, ok := b.BestCar(shaft, fromFloor, toFloor)
	if !ok {
		return 0
	}
	return cost
}

func sign(v float64) float64 {
	switch {
	case v > LevelEps:
		return 1
	case v < -LevelEps:
		return -1
	}
	return 0
}
