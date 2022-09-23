package hub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eightlay/rummikub-server/iternal/game"
	"github.com/gorilla/mux"
)

// Room manager
type Manager struct {
	rooms map[roomID]*room
}

// Create new room manager
func CreateManager() *Manager {
	return &Manager{
		rooms: map[roomID]*room{},
	}
}

// Add player to a room with a desired room size
//
// Returns player id, room id, and error
func (m *Manager) AddPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomSize, err := strconv.Atoi(vars["roomSize"])
	if roomSize > game.MaxPlayersNumber || roomSize < game.MinPlayersNumber || err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"playerID": "", "roomID": "",
		})
		return
	}

	pid := newPlayerID()
	rid := m.findRoom(roomSize)
	m.rooms[rid].addPlayer(pid)

	json.NewEncoder(w).Encode(map[string]string{
		"playerID": pid.string(), "roomID": rid.string(),
	})
}

// Find not full room with a desired room size
func (m *Manager) findRoom(roomSize int) roomID {
	for i, r := range m.rooms {
		if !r.full() && r.roomSize == roomSize {
			return i
		}
	}
	return m.newRoom(roomSize)
}

// Create new room
func (m *Manager) newRoom(roomSize int) roomID {
	rid := newRoomID()

	m.rooms[rid] = &room{
		players:  []playerID{},
		roomSize: roomSize,
		game:     nil,
		Started:  false,
	}

	return rid
}

// Is the room ready to start game
func (m *Manager) RoomIsReady(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rid := roomID(vars["roomID"])

	if !m.roomExists(rid) {
		json.NewEncoder(w).Encode(map[string]bool{"ready": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"ready": m.rooms[rid].ready()})
}

// Start game in the room
func (m *Manager) StartRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rid := roomID(vars["roomID"])

	if !m.roomExists(rid) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	m.rooms[rid].startGame()

	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// Room is started
func (m *Manager) RoomStarted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rid := roomID(vars["roomID"])

	if !m.roomExists(rid) {
		json.NewEncoder(w).Encode(map[string]bool{"started": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"started": m.rooms[rid].Started,
	})
}

// Check if room exists
func (m *Manager) roomExists(rid roomID) bool {
	_, ok := m.rooms[rid]
	return ok
}

// Receive game input
func (m *Manager) ReceiveInput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rid := roomID(vars["roomID"])

	if !m.roomExists(rid) {
		json.NewEncoder(w).Encode(game.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("no room with id: %v", rid),
		})
		return
	}

	w.Write(m.rooms[rid].receiveInput(r))
}

// Get current game state
func (m *Manager) CurrentGameState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := playerID(vars["playerID"])
	rid := roomID(vars["roomID"])

	if !m.roomExists(rid) {
		json.NewEncoder(w).Encode(game.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("no room with id: %v", rid),
		})
		return
	}

	if !m.rooms[rid].playerExists(pid) {
		json.NewEncoder(w).Encode(game.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("no player with id: %v", pid),
		})
		return
	}

	w.Write(m.rooms[rid].currentGameState(pid))
}
