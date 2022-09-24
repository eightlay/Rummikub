// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

// Game
//
// Handles players' actions, provides game logic and tracks history
type Game struct {
	field        field
	history      *history
	bank         pack
	hands        map[player]hand
	stages       map[player]stage
	stepNumber   int
	turn         int
	players      []player
	readyPlayers map[player]bool
	finished     bool
	winner       player
	started      bool
}

// Create new game
func NewGame() *Game {
	return &Game{
		field:      field{},
		history:    createHistory(),
		bank:       createInitialPack(),
		hands:      map[player]hand{},
		stages:     map[player]stage{},
		stepNumber: 1,
		players:    []player{},
		finished:   false,
		started:    false,
	}
}

// Add player
func (g *Game) AddPlayer(p string) *Event {
	if len(g.players) > MaxPlayersNumber || len(g.players) < MinPlayersNumber {
		return &Event{
			EventTypeError,
			EventError{fmt.Sprintf(
				"must be from %v to %v players", MinPlayersNumber, MaxPlayersNumber,
			)},
		}
	}

	player_ := player(p)
	g.players = append(g.players, player_)
	g.readyPlayers[player_] = false
	g.hands[player_] = hand{}
	g.stages[player_] = systemStage

	return &Event{Type: EventTypeSuccess}
}

// Remove player
func (g *Game) RemovePlayer(p string) error {
	if len(g.players) <= 0 {
		return fmt.Errorf(
			"at least one player should be in the game to remove him",
		)
	}

	player_ := player(p)
	playerIndex := -1

	for i, v := range g.players {
		if v == player_ {
			playerIndex = i
		}
	}

	if playerIndex == -1 {
		return fmt.Errorf("no player with name %v", p)
	}

	g.players = append(g.players[:playerIndex], g.players[playerIndex+1:]...)

	if g.started {
		g.bank = append(g.bank, g.hands[player_]...)

		if playerIndex == g.turn {
			g.nextPlayer()
		}
	}

	delete(g.hands, player_)
	delete(g.stages, player_)

	return nil
}

// Try to start the game
func (g *Game) tryStart() {
	for _, r := range g.readyPlayers {
		if !r {
			return
		}
	}
	g.Start()
}

// Start game
func (g *Game) Start() {
	g.shuffleBank()
	g.firstPick()
	g.turnQueue()
	g.started = true
}

// Start game
func (g *Game) IsStarted() bool {
	return g.started
}

// Randomly shuffle bank
func (g *Game) shuffleBank() {
	for i := range g.bank {
		j := rand.Intn(i + 1)
		g.bank[i], g.bank[j] = g.bank[j], g.bank[i]
	}
}

// Deal pieces to players
func (g *Game) firstPick() {
	for p := range g.hands {
		g.hands[p] = hand(g.bank[:HandSize])
		g.bank = g.bank[HandSize:]
	}
}

// Create turn queue
func (g *Game) turnQueue() {
	firstPlayerIndex := -1
	firstPlayerValue := MinNumber - 1
	if firstPlayerValue > JokerNumber {
		firstPlayerValue = JokerNumber
	}

	for i, p := range g.players {
		largest := g.hands[p].largestPieceNumber()

		if largest > firstPlayerValue {
			firstPlayerIndex = i
			firstPlayerValue = largest
		}
	}

	g.turn = firstPlayerIndex
}

// Game state
func (g *Game) State(p string) *State {
	player_ := player(p)
	if _, ok := g.readyPlayers[player_]; !ok {
		return &State{Error: fmt.Sprintf("there is no player with id %v", p)}
	}

	turn := false
	if g.started {
		turn = g.players[g.turn] == player_
	}

	return &State{
		Turn: turn,
		// Field:           g.field,
		Hand:            g.hands[player_],
		AvailableEvents: g.stages[player_].availableEvents(),
		Started:         g.started,
		Finished:        g.finished,
		Winner:          g.winner,
		Error:           "",
	}
}

// Events handler
func (g *Game) HandleEvent(e *Event) *Event {
	// Handle system events
	if _, ok := systemEventsSet[e.Type]; ok {
		return &Event{
			Type: EventTypeSuccess,
		}
	}
	// Handle action
	err := g.handleAction(e)
	if err == nil {
		return &Event{
			Type: EventTypeSuccess,
		}
	}
	return &Event{
		Type: EventTypeError,
		Data: EventError{
			Error: err.Error(),
		},
	}
}

// Handle player's action
func (g *Game) handleAction(e *Event) error {
	data, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	if e.Type == EventTypeReady {
		return g.readyHandle(data)
	}

	err = fmt.Errorf("there is no event: %v", e.Type)

	if !g.started {
		return fmt.Errorf("game is not started yet")
	}

	switch e.Type {
	case EventTypeInitialMeld:
		err = g.initialMeldHandle(data)
	case EventTypeAddPiece:
		err = g.addPieceHandle(data)
	case EventTypeRemovePiece:
		err = g.removePieceHandle(data)
	case EventTypeReplacePiece:
		err = g.replacePieceHandle(data)
	case EventTypeAddCombination:
		err = g.addCombinationHandle(data)
	case EventTypeConcatCombinations:
		err = g.concatCombinations(data)
	case EventTypeSplitCombination:
		err = g.splitCombination(data)
	case EventTypePass:
		err = g.passHandle(data)
	}

	if err == nil {
		// Check if game is finished
		if !g.gameFinished() {
			g.nextPlayer()
		}
	}

	return err
}

// Check if game is finished
func (g *Game) gameFinished() bool {
	player_ := g.players[g.turn]
	finished := len(g.hands[player_]) == 0

	if finished {
		g.finished = true
		g.winner = g.players[g.turn]
	}

	return finished
}

// Next player
func (g *Game) nextPlayer() {
	g.turn += 1

	if g.turn == len(g.players) {
		g.turn = 0
	}
}

// Add penalty pieces to the player's hand
func (g *Game) passHandle(data []byte) error {
	var e EventPass
	json.Unmarshal(data, &e)

	bankLen := len(g.bank)

	var slicePos int
	if bankLen == 0 {
		return nil
	} else if bankLen >= PenaltySize {
		slicePos = PenaltySize
	} else {
		slicePos = bankLen
	}

	g.hands[e.Player] = append(g.hands[e.Player], g.bank[:slicePos]...)
	g.bank = g.bank[slicePos:]

	return nil
}

// Add penalty pieces to the player's hand
func (g *Game) readyHandle(data []byte) error {
	var e EventReady
	json.Unmarshal(data, &e)

	if _, ok := g.readyPlayers[player(e.Player)]; !ok {
		return fmt.Errorf("there is no player with name: %v", e.Player)
	}

	g.readyPlayers[player(e.Player)] = true

	g.tryStart()

	return nil
}

// Initial meld action handler
func (g *Game) initialMeldHandle(data []byte) error {
	var e EventInitialMeld
	json.Unmarshal(data, &e)

	if g.stages[e.Player] != initialMeldStage {
		return fmt.Errorf("wrong game stage for player: %v", e.Player)
	}

	pieces, notFoundIndex := g.gatherPieces(e.Player, e.AddedPieces)
	if pieces == nil {
		return fmt.Errorf("there is no piece with index %v", notFoundIndex)
	}

	combination := validInitialMeld(pieces)
	if combination == nil {
		return fmt.Errorf("invalid combination")
	}

	g.placeCombination(e.Player, combination)
	g.removePiecesFromHand(e.Player, e.AddedPieces)
	g.stages[e.Player] = mainGameStage
	return nil
}

// Add piece handler
func (g *Game) addPieceHandle(data []byte) error {
	var e EventAddPiece
	json.Unmarshal(data, &e)

	if g.stages[e.Player] == initialMeldStage {
		return fmt.Errorf("wrong stage action for player: %v", e.Player)
	}

	if len(e.AddedPieces) != 1 {
		return fmt.Errorf("exactly one piece per action can be added")
	}

	if len(e.UsedCombinations) != 1 {
		return fmt.Errorf(
			"excatly one combination per action can be used for addition",
		)
	}

	stepNumber := e.UsedCombinations[0]
	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return fmt.Errorf("there is no combination with index %v", stepNumber)
	}

	var pieces pack
	var pieceIndex int
	var piece *Piece

	pieceIndex = e.AddedPieces[0]
	piece = g.pieceByIndex(e.Player, pieceIndex)
	if piece == nil {
		return fmt.Errorf(
			"there is no piece with index %v", pieceIndex,
		)
	}

	pieces = append(combination.Pieces, piece)

	newCombination := validCombination(pieces)
	if newCombination == nil {
		return fmt.Errorf(
			"can't add the piece %v to the combination %v",
			pieceIndex, stepNumber,
		)
	}

	g.placeCombination(e.Player, newCombination)
	g.deleteCombinationByStepNumber(stepNumber)
	g.removePieceFromHand(e.Player, pieceIndex)

	return nil
}

// Remove piece handler
func (g *Game) removePieceHandle(data []byte) error {
	var e EventRemovePiece
	json.Unmarshal(data, &e)

	if g.stages[e.Player] == initialMeldStage {
		return fmt.Errorf("wrong stage action for player: %v", e.Player)
	}

	if len(e.UsedCombinations) != 1 {
		return fmt.Errorf(
			"excatly one combination per action can be used for removing",
		)
	}

	stepNumber := e.UsedCombinations[0]
	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return fmt.Errorf("there is no combination with index %v", stepNumber)
	}

	var pieces pack
	var pieceIndex int
	var piece *Piece

	pieceIndex = e.RemovedPiece
	piece = combination.Pieces[pieceIndex]

	pieces = append(
		combination.Pieces[:pieceIndex],
		combination.Pieces[pieceIndex+1:]...,
	)

	newCombination := validCombination(pieces)
	if newCombination == nil {
		return fmt.Errorf(
			"can't remove the piece %v to the combination %v",
			pieceIndex, stepNumber,
		)
	}

	g.placeCombination(e.Player, newCombination)
	g.deleteCombinationByStepNumber(stepNumber)
	g.addPieceToHand(e.Player, piece)

	return nil
}

// Repalce action handler
func (g *Game) replacePieceHandle(data []byte) error {
	var e EventReplacePiece
	json.Unmarshal(data, &e)

	if g.stages[e.Player] == initialMeldStage {
		return fmt.Errorf("wrong stage action for player: %v", e.Player)
	}

	if len(e.AddedPieces) != 1 {
		return fmt.Errorf("exactly one piece per action can be replaced")
	}

	if len(e.UsedCombinations) != 1 {
		return fmt.Errorf(
			"excatly one combination per action can be used for replacing",
		)
	}

	stepNumber := e.UsedCombinations[0]
	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return fmt.Errorf("there is no combination with index %v", stepNumber)
	}

	toAddPieceIndex := e.AddedPieces[0]
	toRemovePieceIndex := e.RemovedPiece

	toAddPiece := g.pieceByIndex(e.Player, toAddPieceIndex)
	pieceToRemove := combination.Pieces[toRemovePieceIndex]

	pieces := combination.Pieces[:]
	pieces[toRemovePieceIndex] = toAddPiece

	newCombination := validCombination(pieces)
	if newCombination == nil {
		return fmt.Errorf(
			"piece %v from hand can't replace piece %v from combination %v",
			toAddPieceIndex, toRemovePieceIndex, stepNumber,
		)
	}

	g.placeCombination(e.Player, newCombination)
	g.deleteCombinationByStepNumber(stepNumber)
	g.removePieceFromHand(e.Player, toAddPieceIndex)
	g.addPieceToHand(e.Player, pieceToRemove)
	pieceToRemove.clearIfJoker()

	return nil
}

// Add combination action handler
func (g *Game) addCombinationHandle(data []byte) error {
	var e EventAddCombination
	json.Unmarshal(data, &e)

	if g.stages[e.Player] == initialMeldStage {
		return fmt.Errorf("wrong game stage for player: %v", e.Player)
	}

	pieces, notFoundIndex := g.gatherPieces(e.Player, e.AddedPieces)
	if pieces == nil {
		return fmt.Errorf("there is no piece with index %v", notFoundIndex)
	}

	newCombination := validCombination(pieces)

	if newCombination != nil {
		g.placeCombination(e.Player, newCombination)
		g.removePiecesFromHand(e.Player, e.AddedPieces)
		return nil
	}

	return fmt.Errorf("invalid combination")
}

// Concat combinations action handler
func (g *Game) concatCombinations(data []byte) error {
	var e EventAddPiece
	json.Unmarshal(data, &e)

	if g.stages[e.Player] == initialMeldStage {
		return fmt.Errorf("wrong stage action for player: %v", e.Player)
	}

	if len(e.UsedCombinations) < 2 {
		return fmt.Errorf("at leat 2 combination can be concatenated")
	}

	pieces := []*Piece{}

	for _, stepNumber := range e.UsedCombinations {
		combination := g.combinationByStepNumber(stepNumber)
		if combination == nil {
			return fmt.Errorf(
				"there is no combination with index %v", stepNumber,
			)
		}

		for _, p := range combination.Pieces {
			pieces = append(pieces, p)
		}
	}

	newCombination := validCombination(pieces)
	if newCombination == nil {
		stepStrings := []string{}
		for _, stepNumber := range e.UsedCombinations {
			stepStrings = append(stepStrings, strconv.Itoa(stepNumber))
		}

		return fmt.Errorf(
			"combinations [%v] can't be concatenated to the valid one",
			strings.Join(stepStrings, ", "),
		)
	}

	g.placeCombination(e.Player, newCombination)

	for _, stepNumber := range e.UsedCombinations {
		g.deleteCombinationByStepNumber(stepNumber)
	}

	return nil
}

// Split combination action handler
func (g *Game) splitCombination(data []byte) error {
	var e EventSplitCombination
	json.Unmarshal(data, &e)

	if g.stages[e.Player] == initialMeldStage {
		return fmt.Errorf("wrong stage action for player: %v", e.Player)
	}

	if len(e.UsedCombinations) != 1 {
		return fmt.Errorf("exactly one combination can be splitted per action")
	}

	stepNumber := e.UsedCombinations[0]

	combination := g.combinationByStepNumber(stepNumber)
	if combination == nil {
		return fmt.Errorf("there is no combination with index %v", stepNumber)
	}

	if e.SplitBeforeIndex >= len(combination.Pieces) {
		return fmt.Errorf(
			"index %v out of range in combination %v",
			e.SplitBeforeIndex, stepNumber,
		)
	}

	pieces1 := combination.Pieces[:e.SplitBeforeIndex]
	pieces2 := combination.Pieces[e.SplitBeforeIndex:]

	newCombination1 := validCombination(pieces1)
	newCombination2 := validCombination(pieces2)

	if newCombination1 == nil || newCombination2 == nil {
		return fmt.Errorf(
			"can't create two valid combinations from splitting combination %v on index %v",
			stepNumber, e.SplitBeforeIndex,
		)
	}

	g.deleteCombinationByStepNumber(stepNumber)
	g.placeCombination(e.Player, newCombination1)
	g.placeCombination(e.Player, newCombination2)

	return nil
}

// Place combination on the field
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

// Find combination by its step number on the field
func (g *Game) combinationByStepNumber(stepNumber int) *Combination {
	for s, c := range g.field {
		if s.number == stepNumber {
			return c
		}
	}
	return nil
}

// Delete combination by its step number from the field
func (g *Game) deleteCombinationByStepNumber(stepNumber int) {
	for s := range g.field {
		if s.number == stepNumber {
			delete(g.field, s)
		}
	}
}

// Find piece by its index in the player's hand
func (g *Game) pieceByIndex(player_ player, pieceIndex int) *Piece {
	if pieceIndex >= len(g.hands[player_]) {
		return nil
	}
	return g.hands[player_][pieceIndex]
}

// Gather pieces together from player's hand by their indeces
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

// Add piece to the player's hand
func (g *Game) addPieceToHand(player_ player, piece *Piece) {
	g.hands[player_] = append(g.hands[player_], piece)
}

// Remove pieces from the player's hand by their indeces
func (g *Game) removePiecesFromHand(player_ player, pieceIndeces []int) {
	indeces := pieceIndeces[:]
	sort.Sort(sort.Reverse(sort.IntSlice(indeces)))

	for _, ind := range indeces {
		g.removePieceFromHand(player_, ind)
	}
}

// Remove piece by its index from the player's hand
func (g *Game) removePieceFromHand(player_ player, pieceIndex int) {
	g.hands[player_] = append(
		g.hands[player_][:pieceIndex],
		g.hands[player_][pieceIndex+1:]...,
	)
}
