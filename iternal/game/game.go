package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

type Game struct {
	field      field
	history    *history
	bank       pack
	hands      map[player]hand
	stages     map[player]stage
	stepNumber int
}

func NewGame(players []string) (*Game, error) {
	if len(players) > 4 || len(players) < 2 {
		return nil, fmt.Errorf("must be from two to four players")
	}

	hands := createHands(players)
	stages := createStages(players)

	return &Game{
		field:      field{},
		history:    createHistory(),
		bank:       createInitialPack(),
		hands:      hands,
		stages:     stages,
		stepNumber: 1,
	}, nil
}

func (g *Game) Start() {
	g.shuffleBank()
	g.firstPick()
}

func (g *Game) shuffleBank() {
	for i := range g.bank {
		j := rand.Intn(i + 1)
		g.bank[i], g.bank[j] = g.bank[j], g.bank[i]
	}
}

func (g *Game) firstPick() {
	for p := range g.hands {
		g.hands[p] = hand(g.bank[:handSize])
		g.bank = g.bank[handSize:]
	}
}

func (g *Game) CurrentState(request []byte) ([]byte, error) {
	// Parse request
	sr := StateRequest{}

	err := json.Unmarshal(request, &sr)
	if err != nil {
		return nil, fmt.Errorf("can't get current state: %v", err)
	}

	player_ := player(sr.Player)

	// Create response
	state := StateResponse{
		Hand:     map[int][]byte{},
		Field:    map[int][]byte{},
		Actions:  []byte{},
		BankSize: len(g.bank),
	}

	// Player's hand
	for i, p := range g.hands[player_] {
		j, err := p.toJSON()
		if err != nil {
			return nil, fmt.Errorf("can't get current state: %v", err)
		}

		state.Hand[i] = j
	}

	// Game field
	for s, c := range g.field {
		j, err := c.toJSON()
		if err != nil {
			return nil, fmt.Errorf("can't get current state: %v", err)
		}

		state.Field[s.number] = j
	}

	// Action list
	var actions []action

	if g.stages[player_] == initialMeldStage {
		actions = initialActions[:]
	} else {
		actions = mainActions[:]
	}

	j, err := json.Marshal(actions)
	if err != nil {
		return nil, fmt.Errorf("can't get current state: %v", err)
	}
	state.Actions = j

	// Convert response to json
	j, err = json.Marshal(&state)
	if err != nil {
		return nil, fmt.Errorf("can't get current state: %v", err)
	}

	return j, nil
}

func (g *Game) HandleAction(request []byte) ([]byte, error) {
	// Parse request
	ar := ActionRequest{}

	err := json.Unmarshal(request, &ar)
	if err != nil {
		return nil, fmt.Errorf("can't get handle action: %v", err)
	}

	// Handle action
	if ar.TimerExceeded || ar.Action == pass {
		g.applyPenalty(&ar)
		return actionSuccess()
	}

	switch ar.Action {
	case initialMeld:
		return g.initialMeldHandle(&ar)
	case addPiece:
		return g.addRemovePieceHandle(&ar, true)
	case removePiece:
		return g.addRemovePieceHandle(&ar, false)
	case replacePiece:
		return g.replacePieceHandle(&ar)
	case addCombination:
		return g.addCombinationHandle(&ar)
	case splitCombination:
		return g.splitCombination(&ar)
	}

	return actionError(fmt.Errorf("unknown action"))
}

func (g *Game) applyPenalty(ar *ActionRequest) {
	player_ := player(ar.Player)
	bankLen := len(g.bank)

	var slicePos int
	if bankLen == 0 {
		return
	} else if bankLen >= penaltySize {
		slicePos = penaltySize
	} else {
		slicePos = bankLen
	}

	g.hands[player_] = append(g.hands[player_], g.bank[:slicePos]...)
	g.bank = g.bank[slicePos:]
}

func (g *Game) initialMeldHandle(ar *ActionRequest) ([]byte, error) {
	pieces := []*Piece{}
	player_ := player(ar.Player)

	for _, pieceIndex := range ar.UsedPieces {
		p := g.pieceByIndex(player_, pieceIndex)
		if p != nil {
			pieces = append(pieces, g.pieceByIndex(player_, pieceIndex))
		} else {
			return actionError(fmt.Errorf("there is no piece with index %v", pieceIndex))
		}
	}

	if g.stages[player_] == initialMeldStage {
		combination := validInitialMeld(pieces)

		if combination != nil {
			g.placeCombination(player_, combination)
			return actionSuccess()
		}
	}

	return actionError(fmt.Errorf("wrong game stage for player: %v", player_))
}

func (g *Game) placeCombination(player_ player, comb *Combination) {
	s := &step{
		number:   g.stepNumber,
		player:   player_,
		prevStep: g.history.lastStep,
		nextStep: nil,
	}

	g.history.lastStep.nextStep = s
	g.history.lastStep = s

	g.field[s] = comb
	g.history.combinations[s] = comb
}

func (g *Game) addRemovePieceHandle(ar *ActionRequest, addFlag bool) ([]byte, error) {
	player_ := player(ar.Player)

	if g.stages[player_] == initialMeldStage {
		return actionError(fmt.Errorf("wrong stage action for player: %v", player_))
	}

	if len(ar.UsedPieces) > 1 || len(ar.UsedPieces) == 0 {
		return actionError(fmt.Errorf("exactly one piece per action can be added\removed"))
	}

	if len(ar.UsedCombinations) > 1 || len(ar.UsedCombinations) == 0 {
		return actionError(
			fmt.Errorf("excatly one combination per action can be used for adding/removing"),
		)
	}

	pieceIndex := ar.UsedPieces[0]
	piece := g.pieceByIndex(player_, pieceIndex)
	if piece == nil {
		return actionError(fmt.Errorf("there is no piece with index %v", pieceIndex))
	}

	stepNumber := ar.UsedCombinations[0]
	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return actionError(fmt.Errorf("there is no combination with index %v", stepNumber))
	}

	var pieces pack

	if addFlag {
		pieces = append(combination.Pieces, piece)
	} else {
		pieces = append(
			combination.Pieces[:pieceIndex],
			combination.Pieces[pieceIndex+1:]...,
		)
	}

	newCombinationType := validCombination(pieces)
	if newCombinationType == notCombination {
		return actionError(fmt.Errorf(
			"can't add the piece %v to the combination %v",
			pieceIndex, stepNumber,
		))
	}

	newCombination := &Combination{
		Pieces: sortPieces(pieces),
		Type:   newCombinationType,
	}

	g.placeCombination(player_, newCombination)
	g.deleteCombinationByStepNumber(stepNumber)

	if addFlag {
		g.removePieceFromHand(player_, pieceIndex)
	} else {
		g.addPieceToHand(player_, piece)
	}

	return actionSuccess()
}

func (g *Game) combinationByStepNumber(stepNumber int) *Combination {
	for s, c := range g.field {
		if s.number == stepNumber {
			return c
		}
	}
	return nil
}

func (g *Game) deleteCombinationByStepNumber(stepNumber int) {
	for s := range g.field {
		if s.number == stepNumber {
			delete(g.field, s)
		}
	}
}

func (g *Game) pieceByIndex(player_ player, pieceIndex int) *Piece {
	if pieceIndex >= len(g.hands[player_]) {
		return nil
	}
	return g.hands[player_][pieceIndex]
}

func (g *Game) addPieceToHand(player_ player, piece *Piece) {
	g.hands[player_] = append(g.hands[player_], piece)
}

func (g *Game) removePiecesFromHand(player_ player, pieceIndeces []int) {
	for _, ind := range pieceIndeces {
		g.removePieceFromHand(player_, ind)
	}
}

func (g *Game) removePieceFromHand(player_ player, pieceIndex int) {
	g.hands[player_] = append(
		g.hands[player_][:pieceIndex],
		g.hands[player_][pieceIndex+1:]...,
	)
}

func (g *Game) replacePieceHandle(ar *ActionRequest) ([]byte, error) {
	// TODO
	return actionSuccess()
}

func (g *Game) addCombinationHandle(ar *ActionRequest) ([]byte, error) {
	// TODO
	return actionSuccess()
}

func (g *Game) splitCombination(ar *ActionRequest) ([]byte, error) {
	// TODO
	return actionSuccess()
}
