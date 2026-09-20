package model

import "math"

// How a waypoint is reached.
const (
	ModeWalk = iota
	ModeStair
	ModeElevator
	ModeEscalator
)

// A step in a route.
type Waypoint struct {
	Col   float64 `json:"col"`
	Floor float64 `json:"floor"`
	Mode  int     `json:"mode"`
	Shaft ShaftID `json:"shaft"` // which shaft, for ModeElevator
}

// One connector traversal in a route.
// A stair run enters at col and leaves at exit.
type stairStep struct {
	col    float64
	exit   float64
	from   int
	to     int
	shaft  ShaftID // 0 for a stair
	moving bool    // an escalator carries them
	hop    float64 // cost of using this connector
}

// Cost of one connector, in cells of walking.
// Keeps a route from hopping between connectors for nothing.
const transferCost = 4.0

// An escalator is easier than stairs.
const escalatorCost = 1.0

// Floor connections, keyed by the floor they leave from.
func (g *Grid) floorEdges() map[int][]stairStep {
	adj := map[int][]stairStep{}
	for _, s := range g.Stairs {
		// The run climbs left to right, so the ends differ.
		lo, hi := s.FootCol(), s.HeadCol()
		hop := transferCost
		if s.Escalator {
			hop = escalatorCost
		}
		adj[s.Floor] = append(adj[s.Floor], stairStep{
			col: lo, exit: hi, from: s.Floor, to: s.Floor + 1,
			moving: s.Escalator, hop: hop})
		adj[s.Floor+1] = append(adj[s.Floor+1], stairStep{
			col: hi, exit: lo, from: s.Floor + 1, to: s.Floor,
			moving: s.Escalator, hop: hop})
	}
	for _, e := range g.Elevators {
		c := e.CenterCol()
		for a := e.Bottom; a <= e.Top; a++ {
			for b := e.Bottom; b <= e.Top; b++ {
				if a != b {
					adj[a] = append(adj[a],
						stairStep{col: c, exit: c, from: a, to: b, shaft: e.ID, hop: transferCost})
				}
			}
		}
	}
	return adj
}

// Extra cost of boarding a car, in cells of walking.
// The dispatcher supplies it; nil means no preference.
type ElevatorCost func(shaft ShaftID, fromFloor, toFloor int) float64

// Cheapest connector route, measured in walking distance.
// Picking by walk keeps sims on the nearest shaft or stair.
func (g *Grid) routeFloors(fromFloor, fromCol, toFloor, toCol int, elevCost ElevatorCost) ([]stairStep, bool) {
	if fromFloor == toFloor {
		return nil, true
	}
	adj := g.floorEdges()

	dist := map[int]float64{fromFloor: 0}
	at := map[int]float64{fromFloor: float64(fromCol)} // column on arrival
	parent := map[int]stairStep{}
	done := map[int]bool{}

	for {
		cur, best, found := 0, math.MaxFloat64, false
		for f, d := range dist {
			if !done[f] && d < best {
				cur, best, found = f, d, true
			}
		}
		if !found || cur == toFloor {
			break
		}
		done[cur] = true

		for _, st := range adj[cur] {
			cost := best + math.Abs(st.col-at[cur]) + st.hop
			if st.shaft != 0 && elevCost != nil {
				cost += elevCost(st.shaft, st.from, st.to)
			}
			if st.to == toFloor {
				cost += math.Abs(float64(toCol) - st.exit)
			}
			if d, seen := dist[st.to]; seen && d <= cost {
				continue
			}
			dist[st.to] = cost
			at[st.to] = st.exit
			parent[st.to] = st
		}
	}

	if _, ok := dist[toFloor]; !ok {
		return nil, false
	}
	steps := []stairStep{}
	for cur := toFloor; cur != fromFloor; {
		st := parent[cur]
		steps = append([]stairStep{st}, steps...)
		cur = st.from
	}
	return steps, true
}

// Waypoints from one cell to another, via connectors.
func (g *Grid) Path(fromFloor, fromCol, toFloor, toCol int) ([]Waypoint, bool) {
	return g.PathVia(fromFloor, fromCol, toFloor, toCol, nil)
}

// Waypoints, letting the dispatcher rate each car.
func (g *Grid) PathVia(fromFloor, fromCol, toFloor, toCol int, elevCost ElevatorCost) ([]Waypoint, bool) {
	steps, ok := g.routeFloors(fromFloor, fromCol, toFloor, toCol, elevCost)
	if !ok {
		return nil, false
	}
	wp := []Waypoint{{Col: float64(fromCol), Floor: float64(fromFloor), Mode: ModeWalk}}
	for _, st := range steps {
		wp = append(wp, Waypoint{Col: st.col, Floor: float64(st.from), Mode: ModeWalk})
		mode := ModeStair
		if st.moving {
			mode = ModeEscalator
		}
		if st.shaft != 0 {
			mode = ModeElevator
		}
		// Stairs land at the far end of the run.
		wp = append(wp, Waypoint{Col: st.exit, Floor: float64(st.to), Mode: mode, Shaft: st.shaft})
	}
	wp = append(wp, Waypoint{Col: float64(toCol), Floor: float64(toFloor), Mode: ModeWalk})
	return wp, true
}
