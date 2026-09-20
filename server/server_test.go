package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// Decodes any server message.
type anyMsg struct {
	Type  string   `json:"type"`
	Chat  []string `json:"chat"`
	Clock Clock    `json:"clock"`
	Money int      `json:"money"`
	Grid  struct {
		Width  int      `json:"width"`
		Floors int      `json:"floors"`
		Built  [][]bool `json:"built"`
		Rooms  []struct {
			Type  string `json:"type"`
			Floor int    `json:"floor"`
			Col   int    `json:"col"`
		} `json:"rooms"`
	} `json:"grid"`
}

func dialTestServer(t *testing.T) (*websocket.Conn, context.Context) {
	t.Helper()
	hub := newHub()
	go hub.run()
	srv := httptest.NewServer(http.HandlerFunc(hub.serveWS))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close(websocket.StatusNormalClosure, "") })
	return conn, ctx
}

func readRaw(t *testing.T, ctx context.Context, conn *websocket.Conn) anyMsg {
	t.Helper()
	var m anyMsg
	if err := wsjson.Read(ctx, conn, &m); err != nil {
		t.Fatalf("read: %v", err)
	}
	return m
}

// Reads the next non-snapshot message.
func read(t *testing.T, ctx context.Context, conn *websocket.Conn) anyMsg {
	t.Helper()
	for {
		m := readRaw(t, ctx, conn)
		if m.Type == "snapshot" {
			continue
		}
		return m
	}
}

func TestInitialSync(t *testing.T) {
	conn, ctx := dialTestServer(t)

	// Connect sends chat, then a snapshot.
	if got := readRaw(t, ctx, conn); got.Type != "chatUpdate" {
		t.Fatalf("first message: want chatUpdate, got %q", got.Type)
	}
	snap := readRaw(t, ctx, conn)
	if snap.Type != "snapshot" {
		t.Fatalf("second message: want snapshot, got %q", snap.Type)
	}
	if snap.Grid.Width != gridWidth || snap.Grid.Floors != gridFloors {
		t.Fatalf("grid: want %dx%d, got %dx%d",
			gridWidth, gridFloors, snap.Grid.Width, snap.Grid.Floors)
	}
}

func TestPlaceRoom(t *testing.T) {
	conn, ctx := dialTestServer(t)

	if err := wsjson.Write(ctx, conn, map[string]any{
		"type": "placeroom", "room": "office", "floor": 0, "col": 3,
	}); err != nil {
		t.Fatalf("write: %v", err)
	}

	for {
		m := readRaw(t, ctx, conn)
		if m.Type != "snapshot" || len(m.Grid.Rooms) == 0 {
			continue
		}
		r := m.Grid.Rooms[0]
		if r.Type != "office" || r.Floor != 0 || r.Col != 3 {
			t.Fatalf("want office at (0,3), got %+v", r)
		}
		return
	}
}

func TestPlaceBase(t *testing.T) {
	conn, ctx := dialTestServer(t)

	if err := wsjson.Write(ctx, conn, map[string]any{
		"type": "placebase", "floor": 0, "col": 2,
	}); err != nil {
		t.Fatalf("write: %v", err)
	}

	for {
		m := readRaw(t, ctx, conn)
		if m.Type == "snapshot" && len(m.Grid.Built) > 0 && m.Grid.Built[0][2] {
			return
		}
	}
}

func TestChatMessage(t *testing.T) {
	conn, ctx := dialTestServer(t)
	read(t, ctx, conn) // initial chatUpdate

	if err := wsjson.Write(ctx, conn, map[string]any{
		"type": "message", "text": "hello tower",
	}); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := read(t, ctx, conn)
	if got.Type != "chatUpdate" || len(got.Chat) != 1 {
		t.Fatalf("want chatUpdate with 1 message, got %+v", got)
	}
	if !strings.Contains(got.Chat[0], "hello tower") {
		t.Fatalf("chat should contain the text, got %q", got.Chat[0])
	}
}

func TestSnapshotBroadcast(t *testing.T) {
	conn, ctx := dialTestServer(t)
	for {
		m := readRaw(t, ctx, conn)
		if m.Type == "snapshot" {
			if m.Clock.Day != 1 || m.Clock.Hour != 7 {
				t.Fatalf("first snapshot: want day 1 07:xx, got %+v", m.Clock)
			}
			return
		}
	}
}

// A client cut off mid-session leaves
// nothing registered behind.

// A resync answers promptly, whole.
func TestResync(t *testing.T) {
	conn, ctx := dialTestServer(t)
	readRaw(t, ctx, conn) // chatUpdate
	readRaw(t, ctx, conn) // snapshot

	if err := wsjson.Write(ctx, conn, map[string]any{
		"type": "resync",
	}); err != nil {
		t.Fatalf("write: %v", err)
	}
	for {
		m := readRaw(t, ctx, conn)
		if m.Type != "snapshot" {
			continue // sims and other chatter
		}
		if m.Money != startMoney {
			t.Fatalf("resync money: want %d, got %d", startMoney, m.Money)
		}
		return
	}
}

func TestDroppedClientLeaves(t *testing.T) {
	hub := newHub()
	go hub.run()
	srv := httptest.NewServer(http.HandlerFunc(hub.serveWS))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	readRaw(t, ctx, conn) // snapshot, so the client is known

	for i := 0; i < 100 && len(hub.clients) != 1; i++ {
		time.Sleep(20 * time.Millisecond)
	}
	if len(hub.clients) != 1 {
		t.Fatalf("want 1 registered client, got %d", len(hub.clients))
	}

	conn.CloseNow() // crash, no goodbye frame

	for i := 0; i < 100 && len(hub.clients) != 0; i++ {
		time.Sleep(20 * time.Millisecond)
	}
	if len(hub.clients) != 0 {
		t.Fatal("client stayed registered after its link died")
	}
}

func TestClockFromSimTime(t *testing.T) {
	if c := clockFromSimTime(0); c.Day != 1 || c.Hour != 7 || c.Minute != 0 {
		t.Fatalf("t=0: want day 1 07:00, got %+v", c)
	}
	// 60 in-game minutes → 08:00.
	if c := clockFromSimTime(60.0 / gameMinutesPerRealSecond); c.Hour != 8 || c.Minute != 0 {
		t.Fatalf("want 08:00, got %+v", c)
	}
	// A full day wraps to day 2, 07:00.
	if c := clockFromSimTime(float64(minutesPerDay) / gameMinutesPerRealSecond); c.Day != 2 || c.Hour != 7 {
		t.Fatalf("want day 2 07:00, got %+v", c)
	}
}
