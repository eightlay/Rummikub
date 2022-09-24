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
