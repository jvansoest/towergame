package main

import "towergame/server/model"

// Any message a client sends.
type Inbound struct {
	Type  string `json:"type"`
	Room  string `json:"room"`  // placeroom: RoomType ID
	Align string `json:"align"` // placeroom: political leaning
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
	Train    *TrainView    `json:"train"` // null with no train
	Vehicles []VehicleView `json:"vehicles"`
	Type     string        `json:"type"` // always "sims"
	Sims     []SimView     `json:"sims"`
	Cars     []CarView     `json:"cars"`
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
	ID         int     `json:"id"`
	X          float64 `json:"x"` // col
	Y          float64 `json:"y"` // floor
	Riding     bool    `json:"riding"`
	Sitting    bool    `json:"sitting"`
	Spot       int     `json:"spot"`     // seat index in the room
	Gliding    bool    `json:"gliding"`  // on an escalator
	Cleaning   bool    `json:"cleaning"` // a maid is sweeping
	Waiting    bool    `json:"waiting"`  // queued at a shaft
	Dark       bool    `json:"dark"`     // the room's lights are off
	Profession string  `json:"profession"`
}

// Per-tick view the client renders.
type Snapshot struct {
	Type       string   `json:"type"` // always "snapshot"
	Tick       uint64   `json:"tick"`
	Clock      Clock    `json:"clock"`
	Money      int      `json:"money"`
	Stars      int      `json:"stars"`      // tower rating, 1 to 5
	Population int      `json:"population"` // residents, staff, guests
	Grid       GridView `json:"grid"`
}

// The tower for the client.
type GridView struct {
	Width     int              `json:"width"`
	Floors    int              `json:"floors"`
	Built     [][]bool         `json:"built"`
	Rooms     []RoomView       `json:"rooms"`
	Stairs    []StairView      `json:"stairs"`
	Ramps     []RampView       `json:"ramps"`
	Below     [][]bool         `json:"below"`    // floors -1, -2, ...
	Basement  int              `json:"basement"` // floors below ground
	Elevators []ElevatorView   `json:"elevators"`
	Types     []model.RoomType `json:"types"` // the buildable catalog
}

// One staircase, with geometry.
type StairView struct {
	Floor     int  `json:"floor"`
	Col       int  `json:"col"`
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	Escalator bool `json:"escalator"` // moves by itself
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
	Open      bool   `json:"open"` // shops: within opening hours
	Alignment string `json:"alignment"`
	Dirty     bool   `json:"dirty"`  // hotel room awaits a maid
	Booked    bool   `json:"booked"` // hotel room has guests
}

// The books for one quarter, sent when it closes.
type Report struct {
	Type    string `json:"type"` // always "report"
	Quarter int    `json:"quarter"`
	Hotel   int    `json:"hotel"`  // guest rent
	Sales   int    `json:"sales"`  // customers
	Leases  int    `json:"leases"` // offices and homes
	Events  int    `json:"events"` // VIP rewards, less penalties
	Upkeep  int    `json:"upkeep"`
	Net     int    `json:"net"`
	Money   int    `json:"money"`
}

// One run of a car ramp, with geometry.
type RampView struct {
	Floor     int  `json:"floor"`
	Col       int  `json:"col"`
	Width     int  `json:"width"`
	DownRight bool `json:"downRight"` // slants toward the right
}

// One car in the garage or on its ramp.
type VehicleView struct {
	ID     int     `json:"id"`
	X      float64 `json:"x"`     // col
	Y      float64 `json:"y"`     // floor, may be fractional
	Dir    int     `json:"dir"`   // -1 faces left, 1 faces right
	Slope  float64 `json:"slope"` // floors per col, signed
	Parked bool    `json:"parked"`
}

// The subway train, while it runs.
type TrainView struct {
	Floor   int     `json:"floor"`
	X       float64 `json:"x"`       // centre column, may be off the grid
	Stopped bool    `json:"stopped"` // at the platform, doors open
}
