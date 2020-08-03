// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
)

// Room represents a single chat room
type Room struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	topic string
}

func newRoom(topic string) *Room {
	return &Room{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		topic:      topic,
	}
}

func (room *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("serving http failed ", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		room: room,
	}

	room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (room *Room) run() {
	log.Printf("running chat room %v", room.topic)
	for {
		select {
		case client := <-room.register:
			log.Printf("new client in room %v", room.topic)
			room.clients[client] = true
		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
			}
			log.Printf("client leaving room %v", room.topic)
		case msg := <-room.broadcast:
			for client := range room.clients {
				select {
				case client.send <- msg:
				default:
					delete(room.clients, client)
					close(client.send)
				}
			}
		}
	}
}
