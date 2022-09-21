package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type Game struct {
	field      field
	history    *history
	bank       pack
	hands      map[player]hand
	stages     map[player]stage
	stepNumber int
}

func (g *Game) FieldSize() int {
	return len(g.field)
}

func (g *Game) HandSize(player_ string) int {
	return len(g.hands[player(player_)])
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

func NewTestGame() (*Game, error) {
	players := []string{"p1", "p2"}

	hands := createHands(players)
	stages := createStages(players)

	g := &Game{
		field:      field{},
		history:    createHistory(),
		bank:       createInitialPack(),
		hands:      hands,
		stages:     stages,
		stepNumber: 1,
	}

	g.testStart()

	return g, nil
}

func (g *Game) Start() {
	g.shuffleBank()
	g.firstPick()
}

func (g *Game) testStart() {
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
		g.hands[p] = hand(g.bank[:HandSize])
		g.bank = g.bank[HandSize:]
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
	if ar.TimerExceeded || ar.Action == Pass {
		g.applyPenalty(&ar)
		return actionSuccess()
	}

	switch ar.Action {
	case InitialMeld:
		response, err := g.initialMeldHandle(&ar)
		if err == nil {
			g.stages[player(ar.Player)] = mainGameStage
		}
		return response, err
	case AddPiece:
		return g.addRemovePieceHandle(&ar, true)
	case RemovePiece:
		return g.addRemovePieceHandle(&ar, false)
	case ReplacePiece:
		return g.replacePieceHandle(&ar)
	case AddCombination:
		return g.addCombinationHandle(&ar)
	case ConcatCombinations:
		return g.concatCombinations(&ar)
	case SplitCombination:
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
	} else if bankLen >= PenaltySize {
		slicePos = PenaltySize
	} else {
		slicePos = bankLen
	}

	g.hands[player_] = append(g.hands[player_], g.bank[:slicePos]...)
	g.bank = g.bank[slicePos:]
}

func (g *Game) initialMeldHandle(ar *ActionRequest) ([]byte, error) {
	player_ := player(ar.Player)

	if g.stages[player_] != initialMeldStage {
		return actionError(fmt.Errorf("wrong game stage for player: %v", player_))
	}

	pieces, notFoundIndex := g.gatherPieces(player_, ar.AddedPieces)
	if pieces == nil {
		return actionError(fmt.Errorf("there is no piece with index %v", notFoundIndex))
	}

	combination := validInitialMeld(pieces)

	if combination != nil {
		g.placeCombination(player_, combination)
		g.removePiecesFromHand(player_, ar.AddedPieces)
		return actionSuccess()
	}

	return actionError(fmt.Errorf("invalid combination"))
}

func (g *Game) placeCombination(player_ player, comb *Combination) {
	s := &step{
		number:   g.stepNumber,
		player:   player_,
		prevStep: g.history.lastStep,
		nextStep: nil,
	}

	g.stepNumber += 1

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

	if (len(ar.AddedPieces) > 1 || len(ar.AddedPieces) == 0) && addFlag {
		return actionError(
			fmt.Errorf("exactly one piece per action can be added/removed"),
		)
	}

	if len(ar.UsedCombinations) > 1 || len(ar.UsedCombinations) == 0 {
		return actionError(fmt.Errorf(
			"excatly one combination per action can be used for adding/removed",
		))
	}

	stepNumber := ar.UsedCombinations[0]
	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return actionError(fmt.Errorf("there is no combination with index %v", stepNumber))
	}

	var pieces pack
	var pieceIndex int
	var piece *Piece

	if addFlag {
		pieceIndex = ar.AddedPieces[0]
		piece = g.pieceByIndex(player_, pieceIndex)
		if piece == nil {
			return actionError(fmt.Errorf(
				"there is no piece with index %v", pieceIndex,
			))
		}

		pieces = append(combination.Pieces, piece)
	} else {
		pieceIndex = ar.RemovedPiece
		piece = combination.Pieces[pieceIndex]

		pieces = append(
			combination.Pieces[:pieceIndex],
			combination.Pieces[pieceIndex+1:]...,
		)
	}

	newCombination := validCombination(pieces)
	if newCombination == nil {
		return actionError(fmt.Errorf(
			"can't add/remove the piece %v to the combination %v",
			pieceIndex, stepNumber,
		))
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

func (g *Game) gatherPieces(player_ player, pieceIndeces []int) ([]*Piece, int) {
	pieces := []*Piece{}

	for _, pieceIndex := range pieceIndeces {
		p := g.pieceByIndex(player_, pieceIndex)
		if p != nil {
			pieces = append(pieces, g.pieceByIndex(player_, pieceIndex))
		} else {
			return nil, pieceIndex
		}
	}

	return pieces, -1
}

func (g *Game) addPieceToHand(player_ player, piece *Piece) {
	g.hands[player_] = append(g.hands[player_], piece)
}

func (g *Game) removePiecesFromHand(player_ player, pieceIndeces []int) {
	indeces := pieceIndeces[:]
	sort.Sort(sort.Reverse(sort.IntSlice(indeces)))

	for _, ind := range indeces {
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
	player_ := player(ar.Player)

	if g.stages[player_] == initialMeldStage {
		return actionError(fmt.Errorf("wrong stage action for player: %v", player_))
	}

	if len(ar.AddedPieces) > 1 || len(ar.AddedPieces) == 0 {
		return actionError(fmt.Errorf("exactly one piece per action can be replaced"))
	}

	if len(ar.UsedCombinations) > 1 || len(ar.UsedCombinations) == 0 {
		return actionError(fmt.Errorf(
			"excatly one combination per action can be used for adding/removing",
		))
	}

	stepNumber := ar.UsedCombinations[0]
	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return actionError(fmt.Errorf("there is no combination with index %v", stepNumber))
	}

	toAddPieceIndex := ar.AddedPieces[0]
	toRemovePieceIndex := ar.RemovedPiece

	toAddPiece := g.pieceByIndex(player_, toAddPieceIndex)
	toRemovePiece := combination.Pieces[toRemovePieceIndex]

	pieces := combination.Pieces[:]
	pieces[toRemovePieceIndex] = toAddPiece

	newCombination := validCombination(pieces)
	if newCombination == nil {
		return actionError(fmt.Errorf(
			"piece %v from hand can't replace piece %v from combination %v",
			toAddPieceIndex, toRemovePieceIndex, stepNumber,
		))
	}

	g.placeCombination(player_, newCombination)
	g.deleteCombinationByStepNumber(stepNumber)
	g.removePieceFromHand(player_, toAddPieceIndex)
	g.addPieceToHand(player_, toRemovePiece)
	toRemovePiece.clearIfJoker()

	return actionSuccess()
}

func (g *Game) addCombinationHandle(ar *ActionRequest) ([]byte, error) {
	player_ := player(ar.Player)

	if g.stages[player_] == initialMeldStage {
		return actionError(fmt.Errorf("wrong game stage for player: %v", player_))
	}

	pieces, notFoundIndex := g.gatherPieces(player_, ar.AddedPieces)
	if pieces == nil {
		return actionError(fmt.Errorf("there is no piece with index %v", notFoundIndex))
	}

	newCombination := validCombination(pieces)

	if newCombination != nil {
		g.placeCombination(player_, newCombination)
		g.removePiecesFromHand(player_, ar.AddedPieces)
		return actionSuccess()
	}

	return actionError(fmt.Errorf("invalid combination"))
}

func (g *Game) splitCombination(ar *ActionRequest) ([]byte, error) {
	player_ := player(ar.Player)

	if g.stages[player_] == initialMeldStage {
		return actionError(fmt.Errorf("wrong stage action for player: %v", player_))
	}

	if len(ar.UsedCombinations) != 1 {
		return actionError(
			fmt.Errorf("exactly one combination can be splitted per action"),
		)
	}

	stepNumber := ar.UsedCombinations[0]

	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return actionError(fmt.Errorf("there is no combination with index %v", stepNumber))
	}

	if ar.SplitAfterIndex >= len(combination.Pieces) {
		return actionError(fmt.Errorf(
			"index %v out of range in combination %v",
			ar.SplitAfterIndex, stepNumber,
		))
	}

	pieces1 := combination.Pieces[:ar.SplitAfterIndex]
	pieces2 := combination.Pieces[ar.SplitAfterIndex:]

	newCombination1 := validCombination(pieces1)
	newCombination2 := validCombination(pieces2)

	if newCombination1 == nil || newCombination2 == nil {
		return actionError(fmt.Errorf(
			"can't create two valid combinations from splitting combination %v on index %v",
			stepNumber, ar.SplitAfterIndex,
		))
	}

	g.deleteCombinationByStepNumber(stepNumber)
	g.placeCombination(player_, newCombination1)
	g.placeCombination(player_, newCombination2)

	return actionSuccess()
}

func (g *Game) concatCombinations(ar *ActionRequest) ([]byte, error) {
	player_ := player(ar.Player)

	if g.stages[player_] == initialMeldStage {
		return actionError(fmt.Errorf("wrong stage action for player: %v", player_))
	}

	if len(ar.UsedCombinations) < 2 {
		return actionError(
			fmt.Errorf("at leat 2 combination can be concatenated"),
		)
	}

	pieces := []*Piece{}

	for _, stepNumber := range ar.UsedCombinations {
		combination := g.combinationByStepNumber(stepNumber)
		if combination == nil {
			return actionError(fmt.Errorf(
				"there is no combination with index %v", stepNumber,
			))
		}

		for _, p := range combination.Pieces {
			pieces = append(pieces, p)
		}
	}

	newCombination := validCombination(pieces)
	if newCombination == nil {
		stepStrings := []string{}
		for _, stepNumber := range ar.UsedCombinations {
			stepStrings = append(stepStrings, strconv.Itoa(stepNumber))
		}

		return actionError(fmt.Errorf(
			"combinations [%v] can't be concatenated to the valid one",
			strings.Join(stepStrings, ", "),
		))
	}

	g.placeCombination(player_, newCombination)

	for _, stepNumber := range ar.UsedCombinations {
		g.deleteCombinationByStepNumber(stepNumber)
	}

	return actionSuccess()
}
