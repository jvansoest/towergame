package main

import (
	"testing"

	"towergame/server/model"
)

func officeRoom(col int) model.Room { return model.Room{Type: "office", Floor: 1, Col: col} }
func condoRoom(col int) model.Room  { return model.Room{Type: "condo", Floor: 2, Col: col} }

// Offices go dark outside working hours.
func TestOfficeLightsFollowHours(t *testing.T) {
	r := officeRoom(0)
	if !lit(r, true, 13*60, false) {
		t.Fatal("office is dark at one in the afternoon")
	}
	if lit(r, true, 3*60, false) {
		t.Fatal("office is lit at three in the morning")
	}
	if lit(r, false, 13*60, false) {
		t.Fatal("office is lit on a weekend afternoon")
	}
}

// Homes go dark through the small hours.
func TestHomesDarkAtNight(t *testing.T) {
	r := condoRoom(0)
	if lit(r, true, 3*60, false) {
		t.Fatal("condo is lit at three in the morning")
	}
	if !lit(r, true, 20*60, false) {
		t.Fatal("condo is dark at eight in the evening")
	}
	// Shops follow their opening hours instead.
	shop := model.Room{Type: "shop", Floor: 3, Col: 0}
	if lit(shop, true, 3*60, false) {
		t.Fatal("shop is lit at three in the morning")
	}
	if !lit(shop, true, 15*60, false) {
		t.Fatal("shop is dark in the afternoon")
	}
	// The lobby is always lit.
	if !lit(model.Room{Type: "lobby", Floor: 0, Col: 0}, true, 3*60, false) {
		t.Fatal("the lobby went dark")
	}
}

// Rooms do not all switch together.
func TestLightsVaryByRoom(t *testing.T) {
	late := map[bool]int{}
	for c := 0; c < 40; c++ {
		late[lit(officeRoom(c), true, 17*60+30, false)]++
	}
	if late[true] == 0 || late[false] == 0 {
		t.Fatalf("every office switched together at 17:30: %v", late)
	}

	night := map[bool]int{}
	for c := 0; c < 40; c++ {
		night[lit(condoRoom(c), true, 23*60+15, false)]++
	}
	if night[true] == 0 || night[false] == 0 {
		t.Fatalf("every condo switched together at 23:15: %v", night)
	}
}
