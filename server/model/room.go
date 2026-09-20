package model

// Room category, drives behavior later.
type Category string

const (
	CategoryLobby       Category = "lobby"
	CategoryOffice      Category = "office"
	CategoryResidential Category = "residential"
	CategoryHotel       Category = "hotel"
	CategoryFood        Category = "food"
	CategoryRetail      Category = "retail"
	CategoryService     Category = "service"
	CategoryMedical     Category = "medical"
	CategoryParking     Category = "parking"
	CategoryTransit     Category = "transit"
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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Width       int      `json:"width"`  // cells
	Height      int      `json:"height"` // floors
	Cost        int      `json:"cost"`
	Capacity    int      `json:"capacity"`
	Category    Category `json:"category"`
	Seats       int      `json:"seats"` // desks or chairs in the room
	Habit       Habit    `json:"habit"`
	Line        bool     `json:"line,omitempty"`        // all seats in one column
	Rent        int      `json:"rent,omitempty"`        // per guest, per stay
	Stars       int      `json:"stars,omitempty"`       // rating needed to build
	Slots       int      `json:"slots,omitempty"`       // cars it holds
	Underground bool     `json:"underground,omitempty"` // built below ground only
	Deepest     bool     `json:"deepest,omitempty"`     // only on the lowest floor
	Unique      bool     `json:"unique,omitempty"`      // one per tower, never demolished
	Lease       int      `json:"lease,omitempty"`       // per quarter, tenants
	Sale        int      `json:"sale,omitempty"`        // per customer served
	Noise       int      `json:"noise,omitempty"`       // 0 quiet .. 3 loud
	Late        bool     `json:"late,omitempty"`        // open into the night
}

// SimTower-style sizes, tunable.
var RoomTypes = map[string]RoomType{
	"lobby":             {ID: "lobby", Name: "Lobby", Width: 1, Height: 1, Cost: 5000, Category: CategoryLobby, Habit: HabitStand},
	"office":            {ID: "office", Name: "Office", Width: 9, Height: 1, Cost: 40000, Capacity: 6, Category: CategoryOffice, Seats: 6, Habit: HabitSit, Lease: 9000},
	"condo":             {ID: "condo", Name: "Condominium", Width: 5, Height: 1, Cost: 80000, Capacity: 3, Category: CategoryResidential, Seats: 3, Habit: HabitSit, Lease: 12000},
	"hotel_single":      {ID: "hotel_single", Name: "Single Room", Width: 3, Height: 1, Cost: 20000, Capacity: 1, Category: CategoryHotel, Seats: 1, Habit: HabitSit, Rent: 3000},
	"hotel_suite":       {ID: "hotel_suite", Name: "Suite", Width: 10, Height: 1, Cost: 100000, Capacity: 3, Category: CategoryHotel, Seats: 3, Habit: HabitSit, Rent: 8000, Stars: 2},
	"shop":              {ID: "shop", Name: "Shop", Width: 12, Height: 1, Cost: 100000, Category: CategoryRetail, Habit: HabitBrowse, Sale: 150},
	"fastfood":          {ID: "fastfood", Name: "Fast Food", Width: 16, Height: 1, Cost: 100000, Category: CategoryFood, Seats: 6, Habit: HabitSit, Noise: 1, Sale: 100},
	"icecream":          {ID: "icecream", Name: "Ice Cream Parlour", Width: 3, Height: 1, Cost: 80000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Line: true, Noise: 1, Sale: 70, Stars: 2},
	"izakaya":           {ID: "izakaya", Name: "Izakaya", Width: 3, Height: 1, Cost: 60000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Line: true, Noise: 3, Late: true, Sale: 250, Stars: 3},
	"burger":            {ID: "burger", Name: "Burger Joint", Width: 3, Height: 1, Cost: 60000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Line: true, Noise: 1, Sale: 90, Stars: 2},
	"pizza":             {ID: "pizza", Name: "Pizza Place", Width: 3, Height: 1, Cost: 60000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Line: true, Noise: 1, Sale: 110, Stars: 2},
	"noodle":            {ID: "noodle", Name: "Noodle Bar", Width: 3, Height: 1, Cost: 60000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Line: true, Noise: 1, Sale: 90, Stars: 2},
	"housekeeping":      {ID: "housekeeping", Name: "Housekeeping", Width: 6, Height: 1, Cost: 50000, Capacity: 2, Category: CategoryService, Habit: HabitStand, Stars: 2},
	"medical":           {ID: "medical", Name: "Medical Centre", Width: 6, Height: 1, Cost: 400000, Capacity: 2, Category: CategoryMedical, Habit: HabitBrowse, Sale: 400, Stars: 4},
	"parking":           {ID: "parking", Name: "Parking", Width: 2, Height: 1, Cost: 20000, Category: CategoryParking, Habit: HabitStand, Slots: 2, Underground: true, Stars: 2},
	"subway":            {ID: "subway", Name: "Subway Station", Width: 14, Height: 2, Cost: 400000, Category: CategoryTransit, Habit: HabitStand, Underground: true, Deepest: true, Unique: true, Stars: 4},
	"restaurant":        {ID: "restaurant", Name: "Restaurant", Width: 18, Height: 1, Cost: 200000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Noise: 1, Sale: 300, Stars: 2},
	"restaurant_indian": {ID: "restaurant_indian", Name: "Indian Restaurant", Width: 18, Height: 1, Cost: 200000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Noise: 1, Sale: 300, Stars: 3},
	"restaurant_french": {ID: "restaurant_french", Name: "French Restaurant", Width: 18, Height: 1, Cost: 200000, Category: CategoryFood, Seats: 8, Habit: HabitSit, Noise: 1, Sale: 350, Stars: 3},
	"security":          {ID: "security", Name: "Security Office", Width: 5, Height: 1, Cost: 60000, Category: CategoryOffice, Capacity: 4, Seats: 4, Habit: HabitSit, Stars: 3},
}

// A room's political leaning.
type Alignment string

const (
	AlignNeutral Alignment = "neutral"
	AlignGood    Alignment = "good"
	AlignTriad   Alignment = "triad"
)

// Whether a string names a real alignment.
func ValidAlignment(a Alignment) bool {
	switch a {
	case AlignNeutral, AlignGood, AlignTriad:
		return true
	}
	return false
}

// A placed room instance.
type Room struct {
	Type      string    `json:"type"` // RoomType ID
	Floor     int       `json:"floor"`
	Col       int       `json:"col"` // leftmost cell
	Alignment Alignment `json:"alignment"`

	// Hotel state, kept by the server only.
	Booked  bool `json:"-"` // guests are coming or here
	Used    bool `json:"-"` // a guest has slept here
	Dirty   bool `json:"-"` // needs housekeeping
	Cleaner int  `json:"-"` // maid sim on the way
}

// Column of one seat, spread left to right.
// The client draws its chairs on the same columns.
func SeatCol(r Room, i int) int {
	rt := RoomTypes[r.Type]
	if rt.Seats < 1 {
		return r.Col + rt.Width/2
	}
	if rt.Line {
		return r.Col + rt.Width/2
	}
	return r.Col + (rt.Width*(i%rt.Seats+1))/(rt.Seats+1)
}
