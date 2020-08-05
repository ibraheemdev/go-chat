package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Hub maintains the set of active rooms and broadcasts messages to the
// rooms.
type Hub struct {
	// Registered rooms.
	rooms map[string]*Room

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// Room maintains the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	// Registered clients.
	clients map[*Client]bool
}

// Message represents a single chat message
type Message struct {
	// The name of the client who sent the message
	Sender string `json:"sender"`

	// The actual text of the message
	Body string `json:"body"`

	// The name of the room the message was sent to
	Room string `json:"room"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			room := h.rooms[client.room]
			if room == nil {
				room = &Room{clients: make(map[*Client]bool)}
				h.rooms[client.room] = room
				log.Print(fmt.Sprint("Registered Room"))
			}
			room.clients[client] = true
		case client := <-h.unregister:
			room := h.rooms[client.room]
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
				if len(room.clients) == 0 {
					delete(h.rooms, client.room)
					log.Print(fmt.Sprint("Closed Room"))
				}
			}
		case message := <-h.broadcast:
			room := h.rooms[message.Room]
			data, err := json.Marshal(message)
			if err != nil {
				log.Print(fmt.Sprint("Error while marshalling message:", err))
			}
			for client := range room.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.rooms, client.room)
					if len(room.clients) == 0 {
						log.Print(fmt.Sprint("Closed Room"))
						delete(h.rooms, client.room)
					}
				}
			}
		}
	}
}
