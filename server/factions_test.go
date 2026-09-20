package main

import (
	"testing"

	"towergame/server/model"
)

// The seeded tower staffs itself: offices send
// workers, homes send residents, and housekeeping
// sends maids. Guests only come in the evening.
func TestSeededProfessions(t *testing.T) {
	w := newWorld()
	w.seed()

	seen := map[Profession]int{}
	for _, p := range w.sims {
		seen[p.prof]++
	}
	if seen[profWorker] == 0 {
		t.Fatal("the seeded tower has no workers")
	}
	if seen[profResident] == 0 {
		t.Fatal("the seeded tower has no residents")
	}
	if seen[profMaid] == 0 {
		t.Fatal("the seeded tower has no maids")
	}
	if seen[profTriad] != 0 {
		t.Fatal("a triad member started in the tower")
	}
}

// A security office staffs four uniformed
// officers, like its capacity promises.
func TestSecurityOfficeSpawnsOfficers(t *testing.T) {
	w := newWorld()
	for c := 0; c < 11; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("security", 0, 1, model.AlignGood); err != nil {
		t.Fatalf("place security office: %v", err)
	}
	if len(w.sims) != model.RoomTypes["security"].Capacity {
		t.Fatalf("want %d officers, got %d",
			model.RoomTypes["security"].Capacity, len(w.sims))
	}
	for _, p := range w.sims {
		if p.prof != profSecurity {
			t.Fatalf("sim %d is a %s, not an officer", p.id, p.prof)
		}
		// On shift, so the stream carries them.
		p.state = statePresent
		p.x, p.y = float64(p.homeC), float64(p.homeF)
	}
	// The stream names their profession.
	up := w.SimsUpdate()
	if len(up.Sims) == 0 || up.Sims[0].Profession != "security" {
		t.Fatalf("stream profession: want security, got %+v", up.Sims)
	}
}

// Alignment sticks at placement, and the
// align command flips it afterwards.
func TestRoomAlignment(t *testing.T) {
	w := newWorld()
	for c := 0; c < 12; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("shop", 0, 1, model.AlignTriad); err != nil {
		t.Fatalf("place: %v", err)
	}
	if got := w.grid.Rooms[0].Alignment; got != model.AlignTriad {
		t.Fatalf("alignment: want triad, got %s", got)
	}

	if err := w.SetAlign(0, 2, model.AlignGood); err != nil {
		t.Fatalf("align: %v", err)
	}
	if got := w.grid.Rooms[0].Alignment; got != model.AlignGood {
		t.Fatalf("realigned to %s, want good", got)
	}
	if err := w.SetAlign(0, 2, "chaos"); err == nil {
		t.Fatal("a nonsense alignment must be refused")
	}
	if err := w.SetAlign(0, 20, model.AlignGood); err == nil {
		t.Fatal("an empty cell has no alignment to set")
	}
}

// Triads favour shops they run, and stay
// three times as long as any customer.
func TestTriadSpawnsInTriadRoom(t *testing.T) {
	w := newWorld()
	for c := 0; c < 30; c++ {
		w.grid.Built[0][c] = true
	}
	if err := w.Place("shop", 0, 1, model.AlignNeutral); err != nil {
		t.Fatalf("place shop: %v", err)
	}
	if err := w.Place("shop", 0, 14, model.AlignTriad); err != nil {
		t.Fatalf("place triad shop: %v", err)
	}
	// Midday: past every shop's staggered opening.
	w.simTime = float64(12*60-startMinuteOfDay) / gameMinutesPerRealSecond

	w.spawnTriad()
	if len(w.sims) != 1 {
		t.Fatalf("want one visitor spawned, got %d", len(w.sims))
	}
	p := w.sims[0]
	if p.prof != profTriad {
		t.Fatalf("spawned a %s, want a triad member", p.prof)
	}
	triad := model.RoomTypes["shop"]
	if p.homeC < 14 || p.homeC >= 14+triad.Width {
		t.Fatalf("triad sits at %d, outside the triad shop", p.homeC)
	}
	if p.pace != triadPace {
		t.Fatalf("triad pace %.1f, want %.1f", p.pace, triadPace)
	}
}

// The sim inspector names each sim's leaning.
func TestSimAlignmentLine(t *testing.T) {
	want := map[Profession]model.Alignment{
		profTriad:    model.AlignTriad,
		profSecurity: model.AlignGood,
		profWorker:   model.AlignNeutral,
		profVisitor:  model.AlignNeutral,
	}
	w := newWorld()
	for prof, align := range want {
		w.nextID++
		w.sims = append(w.sims, &sim{id: w.nextID, prof: prof})
		ins, err := w.InspectSim(w.nextID)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, l := range ins.Lines {
			if l == "Alignment: "+string(align) {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s: no %q line in %v", prof, align, ins.Lines)
		}
	}
}
