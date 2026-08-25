package main

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

const addr = "0.0.0.0:7777"

func main() {
	hub := newHub()
	hub.world.seed()
	go hub.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.serveWS)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("towergame server listening on %s (ws at /ws)", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// Upgrades HTTP to WebSocket.
func (h *Hub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Accept any origin (dev).
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("accept error: %v", err)
		return
	}

	client := &Client{hub: h, conn: conn, send: make(chan any, 16)}
	h.register <- client

	ctx := context.Background()
	go client.writePump(ctx)
	client.readPump(ctx) // blocks until disconnect
}
