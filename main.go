package main

import (
	"encoding/json"
	"fmt"

	"github.com/eightlay/rummikube/iternal/game"
)

func main() {
	players := []string{"me", "opponent"}
	g, err := game.NewGame(players)
	if err != nil {
		panic(err)
	}
	g.Start()

	request := game.StateRequest{Player: "me"}
	j, _ := json.Marshal(&request)

	state, _ := g.CurrentState(j)

	response := game.StateResponse{}
	json.Unmarshal(state, &response)

	fmt.Println(response)
}
