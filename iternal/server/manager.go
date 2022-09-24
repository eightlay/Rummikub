// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

// Hub manager
type Manager struct {
	hubs map[*Hub]bool
}

// Create new hub manager
func newManager() *Manager {
	return &Manager{
		hubs: map[*Hub]bool{},
	}
}

// Get open hub or create new one
func (m *Manager) getHub() *Hub {
	for hub := range m.hubs {
		if !hub.game.IsStarted() {
			return hub
		}
	}
	hub := newHub(m)
	m.hubs[hub] = true
	go hub.run()
	return hub
}

// Remove hub
func (m *Manager) removeHub(hub *Hub) {
	delete(m.hubs, hub)
}
