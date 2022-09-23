package hub

import (
	"net/http"

	"github.com/eightlay/rummikub-server/iternal/game"
	"github.com/google/uuid"
)

// ID of the room
type roomID string

// Create new room ID
func newRoomID() roomID {
	return roomID(uuid.New().String())
}

// Conver player to string
func (r roomID) string() string {
	return string(r)
}

// Game room
type room struct {
	players  []playerID
	roomSize int
	game     *game.Game
	Started  bool
}

// Start game in the room
func (r *room) startGame() {
	players := []string{}
	for _, p := range r.players {
		players = append(players, p.string())
	}

	r.game, _ = game.NewGame(players)
	r.game.Start()
	r.Started = true
}

// Add player to room
func (r *room) addPlayer(id playerID) {
	r.players = append(r.players, id)
}

// Check if room is full
func (r *room) full() bool {
	return len(r.players) >= game.MaxPlayersNumber
}

// Check if room is ready to start
func (r *room) ready() bool {
	return len(r.players) == r.roomSize
}

// Receive game input
func (r *room) receiveInput(request *http.Request) []byte {
	response, _ := r.game.ReceiveActionRequest(request)
	return response
}

// Get current game state
func (r *room) currentGameState(pid playerID) []byte {
	response, _ := r.game.CurrentState(pid.string())
	return response
}

// Check if player exists
func (r *room) playerExists(pid playerID) bool {
	for _, p := range r.players {
		if p == pid {
			return true
		}
	}
	return false
}
