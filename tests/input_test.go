package tests

import (
	"testing"

	"github.com/eightlay/rummikube/iternal/game"
)

func TestInitialMeldAction(t *testing.T) {
	g, _ := game.NewTestGame()

	// Valid action
	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("initial meld was not placed to the game field")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("wrong hand size after initial meld")
	}

	// Invalid action 1
	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{0, 1, 2},
		TimerExceeded: false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2
	action = game.ActionRequest{
		Player:        "p2",
		Action:        game.InitialMeld,
		AddedPieces:   []int{0},
		TimerExceeded: false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 2 changed hand size")
	}
}

func TestPassAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.Pass,
		TimerExceeded: false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 0 {
		t.Fatal("pass action placed pieces to the game field")
	}

	if g.HandSize("p1") != game.HandSize+game.PenaltySize {
		t.Fatal("wrong hand size after initial meld")
	}
}

func TestTimeExceededAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: true,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 0 {
		t.Fatal("pass action placed pieces to the game field")
	}

	if g.HandSize("p1") != game.HandSize+game.PenaltySize {
		t.Fatal("wrong hand size after initial meld")
	}
}

func TestAddPieceAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: false,
	}

	j, _ := action.ToJSON()
	g.ReceiveActionRequest(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{8},
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("adding piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-4 {
		t.Fatal("wrong hand size after adding piece")
	}

	// Invalid action 1: wrong run
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{1},
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-4 {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong number of added pieces
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.AddPiece,
		AddedPieces:      []int{7, 8},
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-4 {
		t.Fatal("invalid action 2 changed hand size")
	}

	// Invalid action 3: trying to add piece before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.AddPiece,
		AddedPieces:      []int{0},
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 3 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 3 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 3 changed hand size")
	}
}

func TestRemovePieceAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11, 12},
		TimerExceeded: false,
	}

	j, _ := action.ToJSON()
	g.ReceiveActionRequest(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.RemovePiece,
		RemovedPiece:     0,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("removing piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("wrong hand size after removing piece")
	}

	// Invalid action 1: trying to add piece before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.RemovePiece,
		RemovedPiece:     0,
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong combination number
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.RemovePiece,
		RemovedPiece:     0,
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("invalid action 2 changed hand size")
	}

	// Invalid action 3: wrong resulting combination
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.RemovePiece,
		RemovedPiece:     0,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 3 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 3 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("invalid action 3 changed hand size")
	}
}

func TestReplacePieceAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{9, 10, 11},
		TimerExceeded: false,
	}

	j, _ := action.ToJSON()
	g.ReceiveActionRequest(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.ReplacePiece,
		AddedPieces:      []int{0},
		RemovedPiece:     0,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("removing piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("wrong hand size after removing piece")
	}

	// Invalid action 1: trying to add piece before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.ReplacePiece,
		AddedPieces:      []int{0},
		RemovedPiece:     0,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong combination number
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.ReplacePiece,
		AddedPieces:      []int{0},
		RemovedPiece:     0,
		UsedCombinations: []int{10},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("invalid action 2 changed hand size")
	}

	// Invalid action 3: wrong resulting combination
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.ReplacePiece,
		AddedPieces:      []int{1},
		RemovedPiece:     0,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 3 was handled")
	}

	if g.FieldSize() != 1 {
		t.Fatal("invalid action 3 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-3 {
		t.Fatal("invalid action 3 changed hand size")
	}
}

func TestAddCombinationAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{10, 11, 12},
		TimerExceeded: false,
	}

	j, _ := action.ToJSON()
	g.ReceiveActionRequest(j)

	// Valid action
	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.AddCombination,
		AddedPieces:   []int{0, 1, 2},
		TimerExceeded: false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 2 {
		t.Fatal("removing piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-6 {
		t.Fatal("wrong hand size after removing piece")
	}

	// Invalid action 1: trying to add piece before initial meld made
	action = game.ActionRequest{
		Player:        "p2",
		Action:        game.AddCombination,
		AddedPieces:   []int{0, 1, 2},
		TimerExceeded: false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 2 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong combination
	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.AddCombination,
		AddedPieces:   []int{0, 3, 5},
		TimerExceeded: false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 2 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-6 {
		t.Fatal("invalid action 2 changed hand size")
	}
}

func TestConcatCombinationsAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{10, 11, 12},
		TimerExceeded: false,
	}

	j, _ := action.ToJSON()
	g.ReceiveActionRequest(j)

	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.AddCombination,
		AddedPieces:   []int{7, 8, 9},
		TimerExceeded: false,
	}

	j, _ = action.ToJSON()
	g.ReceiveActionRequest(j)

	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.AddCombination,
		AddedPieces:   []int{4, 5, 6},
		TimerExceeded: false,
	}

	j, _ = action.ToJSON()
	g.ReceiveActionRequest(j)

	action = game.ActionRequest{
		Player:        "p1",
		Action:        game.AddCombination,
		AddedPieces:   []int{1, 2, 3},
		TimerExceeded: false,
	}

	j, _ = action.ToJSON()
	g.ReceiveActionRequest(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.ConcatCombinations,
		UsedCombinations: []int{1, 2},
		TimerExceeded:    false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 3 {
		t.Fatal("removing piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-12 {
		t.Fatal("wrong hand size after removing piece")
	}

	// Invalid action 1: trying to concat combinations before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.ConcatCombinations,
		UsedCombinations: []int{3, 4},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 3 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong combination
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.ConcatCombinations,
		UsedCombinations: []int{4, 5},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 3 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-12 {
		t.Fatal("invalid action 2 changed hand size")
	}
}

func TestSplitCombinationsAction(t *testing.T) {
	g, _ := game.NewTestGame()

	action := game.ActionRequest{
		Player:        "p1",
		Action:        game.InitialMeld,
		AddedPieces:   []int{7, 8, 9, 10, 11, 12},
		TimerExceeded: false,
	}

	j, _ := action.ToJSON()
	g.ReceiveActionRequest(j)

	// Valid action
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.SplitCombination,
		SplitBeforeIndex: 3,
		UsedCombinations: []int{1},
		TimerExceeded:    false,
	}

	j, err := action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	response, err := g.ReceiveActionRequest(j)
	if err != nil {
		t.Fatalf("can't handle action: %v", err)
	}

	_, err = game.ParseActionResponse(response)
	if err != nil {
		t.Fatalf("can't parse action response: %v", err)
	}

	if g.FieldSize() != 2 {
		t.Fatal("removing piece changed field size")
	}

	if g.HandSize("p1") != game.HandSize-6 {
		t.Fatal("wrong hand size after removing piece")
	}

	// Invalid action 1: trying to split combinations before initial meld made
	action = game.ActionRequest{
		Player:           "p2",
		Action:           game.SplitCombination,
		SplitBeforeIndex: 2,
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action was handled: %v", err)
	}

	if g.FieldSize() != 2 {
		t.Fatal("invalid action 1 changed field size")
	}

	if g.HandSize("p2") != game.HandSize {
		t.Fatal("invalid action 1 changed hand size")
	}

	// Invalid action 2: wrong combination
	action = game.ActionRequest{
		Player:           "p1",
		Action:           game.SplitCombination,
		SplitBeforeIndex: 1,
		UsedCombinations: []int{2},
		TimerExceeded:    false,
	}

	j, err = action.ToJSON()
	if err != nil {
		t.Fatalf("can't convert action to json: %v", err)
	}

	_, err = g.ReceiveActionRequest(j)
	if err == nil {
		t.Fatalf("invalid action 2 was handled")
	}

	if g.FieldSize() != 2 {
		t.Fatal("invalid action 2 changed field size")
	}

	if g.HandSize("p1") != game.HandSize-6 {
		t.Fatal("invalid action 2 changed hand size")
	}
}
