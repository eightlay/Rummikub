package tests

import (
	"testing"

	"github.com/eightlay/rummikub-server/iternal/game"
)

func TestGameLogic(t *testing.T) {
	g, _ := game.NewTestGame()

	jsonState, err := g.CurrentState("p1")
	if err != nil {
		t.Fatalf("can't get current state: %v", err)
	}

	state, err := game.ParseStateResponse(jsonState)
	if err != nil {
		t.Fatalf("can't parse current state: %v", err)
	}

	futureWinner := state.Turn

	piecesToAdd := []int{}
	for i := 0; i < game.HandSize; i++ {
		piecesToAdd = append(piecesToAdd, i)
	}

	action, _ := (&game.ActionRequest{
		Player:      string(futureWinner),
		Action:      game.InitialMeld,
		AddedPieces: piecesToAdd,
	}).ToJSON()

	g.ReceiveActionRequest(action)

	jsonState, _ = g.CurrentState(string(futureWinner))
	state, _ = game.ParseStateResponse(jsonState)

	if !state.Finished {
		t.Fatal("game is not finished")
	}

	if state.Winner != futureWinner {
		t.Fatal("wrong winner")
	}
}
