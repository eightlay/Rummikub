// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GORILLA-LICENSE file.

package server

import (
	"encoding/json"

	"github.com/eightlay/rummikub-server/iternal/game"
	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Pointer to manager
	manager *Manager

	// Registered clients.
	clients map[*Client]uuid.UUID

	// Game
	game *game.Game

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(manager *Manager) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]uuid.UUID),
		game:       game.NewGame(),
		manager:    manager,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = uuid.New()
			h.game.AddPlayer(h.clients[client].String())
			client.conn.WriteJSON(game.Event{
				Type: game.EventTypeInit,
				Data: game.EventInit{
					Player: h.clients[client].String(),
				},
			})
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.game.RemovePlayer(h.clients[client].String())

				message, _ := json.Marshal(game.Event{
					Type: game.EventTypeDisconnect,
				})

				client.conn.WriteJSON(message)
				h.broadcast <- message

				delete(h.clients, client)
				close(client.send)

				if len(h.clients) == 0 {
					h.sendRemoveHub()
				}
			}
		case <-h.broadcast:
			for client, cid := range h.clients {
				select {
				case client.send <- h.game.State(cid.String()).ToJSON():
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) sendRemoveHub() {
	h.manager.removeHub(h)
}
