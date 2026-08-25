package model

// Room category, drives behavior later.
type Category string

const (
	CategoryLobby         Category = "lobby"
	CategoryOffice        Category = "office"
	CategoryResidential   Category = "residential"
	CategoryHotel         Category = "hotel"
	CategoryFood          Category = "food"
	CategoryRetail        Category = "retail"
	CategoryEntertainment Category = "entertainment"
)

// What sims do while they are in a room.
type Habit string

const (
	HabitStand  Habit = "stand"
	HabitSit    Habit = "sit"    // at a desk or table
	HabitBrowse Habit = "browse" // wanders the shop floor
)

// A buildable room type.
type RoomType struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Width    int      `json:"width"`  // cells
	Height   int      `json:"height"` // floors
	Cost     int      `json:"cost"`
	Capacity int      `json:"capacity"`
	Category Category `json:"category"`
	Seats    int      `json:"seats"` // desks or chairs in the room
	Habit    Habit    `json:"habit"`
}

// SimTower-style sizes, tunable.
var RoomTypes = map[string]RoomType{
	"lobby":        {ID: "lobby", Name: "Lobby", Width: 1, Height: 1, Cost: 5000, Category: CategoryLobby, Habit: HabitStand},
	"office":       {ID: "office", Name: "Office", Width: 9, Height: 1, Cost: 40000, Capacity: 6, Category: CategoryOffice, Seats: 6, Habit: HabitSit},
	"condo":        {ID: "condo", Name: "Condominium", Width: 10, Height: 1, Cost: 80000, Capacity: 3, Category: CategoryResidential, Seats: 3, Habit: HabitSit},
	"hotel_single": {ID: "hotel_single", Name: "Single Room", Width: 5, Height: 1, Cost: 20000, Capacity: 1, Category: CategoryHotel, Seats: 1, Habit: HabitSit},
	"hotel_suite":  {ID: "hotel_suite", Name: "Suite", Width: 10, Height: 1, Cost: 100000, Capacity: 3, Category: CategoryHotel, Seats: 3, Habit: HabitSit},
	"shop":         {ID: "shop", Name: "Shop", Width: 12, Height: 1, Cost: 100000, Category: CategoryRetail, Habit: HabitBrowse},
	"fastfood":     {ID: "fastfood", Name: "Fast Food", Width: 16, Height: 1, Cost: 100000, Category: CategoryFood, Seats: 6, Habit: HabitSit},
	"restaurant":   {ID: "restaurant", Name: "Restaurant", Width: 24, Height: 1, Cost: 200000, Category: CategoryFood, Seats: 8, Habit: HabitSit},
	"cinema":       {ID: "cinema", Name: "Cinema", Width: 32, Height: 2, Cost: 500000, Category: CategoryEntertainment, Habit: HabitStand},
}

// A placed room instance.
type Room struct {
	Type  string `json:"type"` // RoomType ID
	Floor int    `json:"floor"`
	Col   int    `json:"col"` // leftmost cell
}

// Column of one seat, spread left to right.
// The client draws its chairs on the same columns.
func SeatCol(r Room, i int) int {
	rt := RoomTypes[r.Type]
	if rt.Seats < 1 {
		return r.Col + rt.Width/2
	}
	return r.Col + (rt.Width*(i%rt.Seats+1))/(rt.Seats+1)
}
