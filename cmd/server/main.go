package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/eightlay/rummikube/iternal/game"
)

func main() {
	file, err := os.OpenFile("runtime.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(file)

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
