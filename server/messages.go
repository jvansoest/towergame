package main

// Any message a client sends.
type Inbound struct {
	Type  string `json:"type"`
	Room  string `json:"room"` // placeroom: RoomType ID
	Floor int    `json:"floor"`
	Col   int    `json:"col"`
	ID    int    `json:"id"`    // inspectsim/inspectcar: which one
	Text  string `json:"text"`  // message: chat text
	Quiet bool   `json:"quiet"` // inspect: a refresh, so no error
}

// Carries the full chat log.
type ChatUpdate struct {
	Type string   `json:"type"` // always "chatUpdate"
	Chat []string `json:"chat"`
}

// Tells one client why an action failed.
type ServerError struct {
	Type   string `json:"type"` // always "error"
	Reason string `json:"reason"`
}

// What the inspector shows about one thing.
type Inspection struct {
	Type  string   `json:"type"` // always "inspect"
	Title string   `json:"title"`
	Lines []string `json:"lines"`
}

// Sim positions, streamed at render rate.
type SimsUpdate struct {
	Type string    `json:"type"` // always "sims"
	Sims []SimView `json:"sims"`
	Cars []CarView `json:"cars"`
}

// One elevator car's position.
type CarView struct {
	ID    int     `json:"id"`
	Shaft int     `json:"shaft"`
	Col   float64 `json:"col"`
	Floor float64 `json:"floor"`
}

// One sim's position in grid coordinates.
type SimView struct {
	ID      int     `json:"id"`
	X       float64 `json:"x"` // col
	Y       float64 `json:"y"` // floor
	Riding  bool    `json:"riding"`
	Sitting bool    `json:"sitting"`
	Dark    bool    `json:"dark"` // the room's lights are off
}

// Per-tick view the client renders.
type Snapshot struct {
	Type  string   `json:"type"` // always "snapshot"
	Tick  uint64   `json:"tick"`
	Clock Clock    `json:"clock"`
	Money int      `json:"money"`
	Grid  GridView `json:"grid"`
}

// The tower for the client.
type GridView struct {
	Width     int            `json:"width"`
	Floors    int            `json:"floors"`
	Built     [][]bool       `json:"built"`
	Rooms     []RoomView     `json:"rooms"`
	Stairs    []StairView    `json:"stairs"`
	Elevators []ElevatorView `json:"elevators"`
}

// One staircase, with geometry.
type StairView struct {
	Floor  int `json:"floor"`
	Col    int `json:"col"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// One elevator shaft, with geometry.
type ElevatorView struct {
	ID     int `json:"id"`
	Bottom int `json:"bottom"`
	Top    int `json:"top"`
	Col    int `json:"col"`
	Width  int `json:"width"`
}

// One placed room, with geometry.
type RoomView struct {
	Type      string `json:"type"`
	Floor     int    `json:"floor"`
	Col       int    `json:"col"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Category  string `json:"category"`
	Capacity  int    `json:"capacity"`
	Seats     int    `json:"seats"`
	Habit     string `json:"habit"`
	Lit       bool   `json:"lit"`
	Occupants int    `json:"occupants"`
}
