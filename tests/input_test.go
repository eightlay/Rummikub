package tests

import (
	"encoding/json"
	"testing"

	"github.com/eightlay/rummikube/iternal/game"
)

func TestInitialMeldAction(t *testing.T) {
	g, err := game.NewTestGame()
	if err != nil {
		t.Errorf("test game creation error: %v", err)
	}

	// Valid action
	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: false,
	}

	j, err := json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	response, err := g.HandleAction(j)
	if err != nil {
		t.Errorf("can't handle action: %v", err)
	}

	var ar game.ActionResponse
	err = json.Unmarshal(response, &ar)
	if err != nil {
		t.Errorf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Error("initial meld was not placed to the game field")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Error("wrong hand size after initial meld")
	}

	// Invalid action 1
	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{0, 1, 2},
		TimerExceeded: false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 1 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Error("invalid action 1 changed hand size")
	}

	// Invalid action 2
	action = game.ActionRequest{
		Player:        "p2",
		Action:        game.InitialMeld,
		AddedPieces:   []int{0},
		TimerExceeded: false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 2 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Error("invalid action 2 changed hand size")
	}
}

func TestPassAction(t *testing.T) {
	g, err := game.NewTestGame()
	if err != nil {
		t.Errorf("test game creation error: %v", err)
	}

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.Pass,
		TimerExceeded: false,
	}

	j, err := json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err != nil {
		t.Errorf("can't handle action: %v", err)
	}

	if g.FieldSize() != 0 {
		t.Error("pass action placed pieces to the game field")
	}

	if g.HandSize("p1") != game.HandSize+game.PenaltySize {
		t.Error("wrong hand size after initial meld")
	}
}

func TestTimeExceededAction(t *testing.T) {
	g, err := game.NewTestGame()
	if err != nil {
		t.Errorf("test game creation error: %v", err)
	}

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: true,
	}

	j, err := json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err != nil {
		t.Errorf("can't handle action: %v", err)
	}

	if g.FieldSize() != 0 {
		t.Error("pass action placed pieces to the game field")
	}

	if g.HandSize("p1") != game.HandSize+game.PenaltySize {
		t.Error("wrong hand size after initial meld")
	}
}

func TestAddPieceAction(t *testing.T) {
	g, err := game.NewTestGame()
	if err != nil {
		t.Errorf("test game creation error: %v", err)
	}

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: false,
	}

	j, _ := json.Marshal(&action)
	g.HandleAction(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{8},
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	response, err := g.HandleAction(j)
	if err != nil {
		t.Errorf("can't handle action: %v", err)
	}

	var ar game.ActionResponse
	err = json.Unmarshal(response, &ar)
	if err != nil {
		t.Errorf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Error("adding piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-4 {
		t.Error("wrong hand size after adding piece")
	}

	// Invalid action 1: wrong run
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{0},
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 1 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-4 {
		t.Error("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong number of added pieces
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{7, 8},
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 2 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Error("invalid action 2 changed hand size")
	}

	// Invalid action 3: trying to add piece before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.AddPiece,
		AddedPieces:      []int{0},
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action 3 was handled")
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 3 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Error("invalid action 3 changed hand size")
	}
}

func TestRemovePieceAction(t *testing.T) {
	g, err := game.NewTestGame()
	if err != nil {
		t.Errorf("test game creation error: %v", err)
	}

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11, 12},
		TimerExceeded: false,
	}

	j, _ := json.Marshal(&action)
	g.HandleAction(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.RemovePiece,
		RemovedPiece:     0,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	response, err := g.HandleAction(j)
	if err != nil {
		t.Errorf("can't handle action: %v", err)
	}

	var ar game.ActionResponse
	err = json.Unmarshal(response, &ar)
	if err != nil {
		t.Errorf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Error("removing piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Error("wrong hand size after removing piece")
	}

	// Invalid action 1: wrong run
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{0},
		UsedCombinations: []int{3},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 1 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-4 {
		t.Error("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong number of added pieces
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{7, 8},
		UsedCombinations: []int{3},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 2 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Error("invalid action 2 changed hand size")
	}

	// Invalid action 3: trying to add piece before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.AddPiece,
		AddedPieces:      []int{0},
		UsedCombinations: []int{3},
		TimerExceeded:    false,
	}

	j, err = json.Marshal(&action)
	if err != nil {
		t.Errorf("can't convert action to json: %v", err)
	}

	_, err = g.HandleAction(j)
	if err == nil {
		t.Errorf("invalid action 3 was handled")
	}

	if g.FieldSize() != 1 {
		t.Error("invalid action 3 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Error("invalid action 3 changed hand size")
	}
}
