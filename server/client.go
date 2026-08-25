package main

import (
	"context"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
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

// Writes queued messages until closed.
func (c *Client) writePump(ctx context.Context) {
	for msg := range c.send {
		if err := wsjson.Write(ctx, c.conn, msg); err != nil {
			log.Printf("write error: %v", err)
			return
		}
	}
}
