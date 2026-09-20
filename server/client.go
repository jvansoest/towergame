package main

import (
	"context"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// Client probes: cadence and reply timeout.
const (
	pingEvery = 30 * time.Second
	pingWait  = 10 * time.Second
)

// One WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan any // outbound JSON messages
}

// Queues a message, never blocks.
func (c *Client) trySend(msg any) {
	select {
	case c.send <- msg:
	default:
	}
}

// Reads until the connection closes.
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		var in Inbound
		if err := wsjson.Read(ctx, c.conn, &in); err != nil {
			return // closed or bad frame
		}
		c.hub.commands <- command{client: c, msg: in}
	}
}

// Writes queued messages and probes the link.
// One goroutine owns all writes to the socket.
func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingEvery)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return // hub hung up on a slow client
			}
			if err := wsjson.Write(ctx, c.conn, msg); err != nil {
				log.Printf("write error: %v", err)
				c.drop()
				return
			}
		case <-ticker.C:
			if err := c.ping(); err != nil {
				log.Printf("ping error: %v", err)
				c.drop()
				return
			}
		}
	}
}

// Probes the client once, bounded.
// A missed pong means a dead link:
// writes alone can sit unseen in buffers.
func (c *Client) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), pingWait)
	defer cancel()
	return c.conn.Ping(ctx)
}

// Cuts the link at once. Ends both pumps.
func (c *Client) sever() {
	c.conn.CloseNow()
}

// Hangs up on a faulty client. Unregister
// first, so broadcasts stop; then closing
// frees the reader. The hub skips repeats.
func (c *Client) drop() {
	c.hub.unregister <- c
	c.sever()
}
