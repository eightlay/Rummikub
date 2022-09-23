package tests

import (
	"testing"

	"github.com/eightlay/rummikub-server/iternal/game"
)

func TestGameCreation(t *testing.T) {
	_, err := game.NewGame([]string{"p1", "p2"})
	if err != nil {
		t.Fatalf("game creation error: %v", err)
	}
}

func TestGameStart(t *testing.T) {
	g, _ := game.NewGame([]string{"p1", "p2"})
	g.Start()
}
