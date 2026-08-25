package main

import (
	"log"
	"time"
)

// Simulation ticks per second.
const tickHz = 30

// Fixed timestep per tick.
const dtSeconds = 1.0 / float64(tickHz)

// Stream sim positions every few ticks.
const broadcastEvery = 3

// Inbound message plus its sender.
type command struct {
	client *Client
	msg    Inbound
}

// Networking hub. Routes commands to the World.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	commands   chan command
	world      *World
	tickCount  uint64
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		commands:   make(chan command),
		world:      newWorld(),
	}
}

// Runs the loop: connections, commands, ticks.
func (h *Hub) run() {
	ticker := time.NewTicker(time.Second / tickHz)
	defer ticker.Stop()

	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			// Send current state to newcomer.
			c.trySend(h.world.ChatUpdate())
			c.trySend(h.world.Snapshot())

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}

		case cmd := <-h.commands:
			h.apply(cmd)

		case <-ticker.C:
			if h.world.Step(dtSeconds) {
				h.broadcast(h.world.Snapshot())
			}
			h.tickCount++
			if h.tickCount%broadcastEvery == 0 {
				h.broadcast(h.world.SimsUpdate())
			}
		}
	}
}

// Routes a command to the World.
func (h *Hub) apply(cmd command) {
	switch cmd.msg.Type {
	case "placeroom":
		if err := h.world.Place(cmd.msg.Room, cmd.msg.Floor, cmd.msg.Col); err != nil {
			cmd.client.trySend(ServerError{Type: "error", Reason: err.Error()})
			return
		}
		h.broadcast(h.world.Snapshot())

	case "placebase":
		if err := h.world.BuildBase(cmd.msg.Floor, cmd.msg.Col); err != nil {
			cmd.client.trySend(ServerError{Type: "error", Reason: err.Error()})
			return
		}
		h.broadcast(h.world.Snapshot())

	case "placestair":
		if err := h.world.PlaceStair(cmd.msg.Floor, cmd.msg.Col); err != nil {
			cmd.client.trySend(ServerError{Type: "error", Reason: err.Error()})
			return
		}
		h.broadcast(h.world.Snapshot())

	case "placeelevator":
		if err := h.world.PlaceElevator(cmd.msg.Floor, cmd.msg.Col); err != nil {
			cmd.client.trySend(ServerError{Type: "error", Reason: err.Error()})
			return
		}
		h.broadcast(h.world.Snapshot())

	case "remove":
		if err := h.world.Remove(cmd.msg.Floor, cmd.msg.Col); err != nil {
			cmd.client.trySend(ServerError{Type: "error", Reason: err.Error()})
			return
		}
		h.broadcast(h.world.Snapshot())

	case "inspect":
		info, err := h.world.InspectCell(cmd.msg.Floor, cmd.msg.Col)
		cmd.client.reply(info, err, cmd.msg.Quiet)

	case "inspectsim":
		info, err := h.world.InspectSim(cmd.msg.ID)
		cmd.client.reply(info, err, cmd.msg.Quiet)

	case "inspectcar":
		info, err := h.world.InspectCar(cmd.msg.ID)
		cmd.client.reply(info, err, cmd.msg.Quiet)

	case "message":
		if h.world.PostMessage(cmd.msg.Text) {
			h.broadcast(h.world.ChatUpdate())
		}

	default:
		log.Printf("unknown message type: %q", cmd.msg.Type)
	}
}

// Answers one inspect request, or says why not.
// A refresh stays quiet, so no toast repeats.
func (c *Client) reply(info Inspection, err error, quiet bool) {
	if err != nil {
		if !quiet {
			c.trySend(ServerError{Type: "error", Reason: err.Error()})
		}
		return
	}
	c.trySend(info)
}

// Sends msg to all clients.
// Drops a client with a full buffer.
func (h *Hub) broadcast(msg any) {
	for c := range h.clients {
		select {
		case c.send <- msg:
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}
