package hub

import (
	"github.com/google/uuid"
)

// Player ID
type playerID string

// Create new player ID
func newPlayerID() playerID {
	return playerID(uuid.New().String())
}

// Conver player to string
func (p playerID) string() string {
	return string(p)
}
