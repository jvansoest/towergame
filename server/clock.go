package main

// In-game time scale.
const (
	gameMinutesPerRealSecond = 8.0
	startMinuteOfDay         = 7 * 60 // start at 07:00
	minutesPerDay            = 24 * 60
)

// In-game date and time.
type Clock struct {
	Day    int `json:"day"`    // 1-based
	Hour   int `json:"hour"`   // 0..23
	Minute int `json:"minute"` // 0..59
}

// Weekdays are days 1-2 of the 3-day week.
func isWeekday(day int) bool {
	return (day-1)%3 < 2
}

// Converts in-game seconds to a Clock.
func clockFromSimTime(simTime float64) Clock {
	total := startMinuteOfDay + int(simTime*gameMinutesPerRealSecond)
	dayMinutes := ((total % minutesPerDay) + minutesPerDay) % minutesPerDay
	return Clock{
		Day:    total/minutesPerDay + 1,
		Hour:   dayMinutes / 60,
		Minute: dayMinutes % 60,
	}
}
