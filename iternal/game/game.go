package game

import (
	"encoding/json"
	"fmt"
	"math/rand"

	mapset "github.com/deckarep/golang-set/v2"
)

type Game struct {
	field      field
	history    *history
	bank       bag
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
		bank:       createInitialBag(),
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

	// TODO: implement all actions
	switch ar.Action {
	case initialMeld:
		return g.initialMeldHandle(&ar)
	case addPiece:
		return g.addPieceHandle(&ar)
	case removePiece:
		return g.removePiece(&ar)
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
	handLen := len(g.hands[player_])

	for _, pieceIndex := range ar.UsedPieces {
		if pieceIndex >= handLen {
			return actionError(
				fmt.Errorf("player doesn't have piece with key %v", pieceIndex),
			)
		}
		pieces = append(pieces, g.hands[player_][pieceIndex])
	}

	if g.stages[player_] == initialMeldStage {
		combination := g.isValidInitialMeld(pieces)

		if combination != nil {
			g.placeCombination(player_, combination)
			return actionSuccess()
		}
	}

	return actionError(fmt.Errorf("wrong game stage for player: %v", player_))
}

func (g *Game) isValidInitialMeld(pieces []*Piece) *Combination {
	ct := g.isValidCombination(pieces)

	if ct != notCombination {
		s := 0

		for _, p := range pieces {
			s += p.Number
		}

		correct := s >= initialMeldSum

		if correct {
			return &Combination{sortPieces(pieces), ct}
		}
	}

	return nil
}

func (g *Game) isValidCombination(pieces []*Piece) combinationType {
	validGroup := g.isValidGroup(pieces)

	if !validGroup {
		validRun := g.isValidRun(pieces)

		if !validRun {

			for _, p := range pieces {
				if p.Joker {
					p.Number = jokerNumber
					p.Color = jokerColor
				}
			}

			return notCombination
		} else {
			return run
		}
	}

	return group
}

func (g *Game) isValidGroup(pieces []*Piece) bool {
	if len(pieces) < minGroupSize || len(pieces) > maxGroupSize {
		return false
	}

	usedColors := mapset.NewSet[color]()
	var number int = jokerNumber

	for _, p := range pieces {
		if p.Joker {
			continue
		}

		if number == 0 {
			number = p.Number
		} else if number != p.Number {
			return false
		} else if usedColors.Contains(p.Color) {
			return false
		}

		usedColors.Add(p.Color)
	}

	for _, p := range pieces {
		if p.Joker {
			c, _ := colorsSet.Difference(usedColors).Pop()
			p.Number = number
			p.Color = c
		}
	}

	return true
}

func (g *Game) isValidRun(pieces_ []*Piece) bool {
	if len(pieces_) < minRunSize {
		return false
	}

	pieces := sortPieces(pieces_)

	runColor := pieces[len(pieces)-1].Color

	jokerCount := 0
	jokerValues := []int{}
	startIndex := 1

	if pieces[0].Joker {
		startIndex += 1
		jokerCount += 1
	}

	if pieces[1].Joker {
		startIndex += 1
		jokerCount += 1
	}

	var lastNumber = pieces[startIndex-1].Number

	for i := startIndex; i < len(pieces); i++ {
		if pieces[i].Color != runColor {
			return false
		}

		diff := pieces[i].Number - lastNumber

		if diff <= 0 {
			return false
		}

		if diff != 1 {
			if jokerCount == 0 {
				return false
			}

			jokersNeeded := diff - 1

			if jokersNeeded > jokerCount {
				return false
			}

			jokerCount -= jokersNeeded

			for j := 1; j <= jokersNeeded; j++ {
				lastNumber += 1
				jokerValues = append(jokerValues, lastNumber)
			}
		}

		lastNumber = pieces[i].Number
	}

	for i := 0; i <= jokerCount; i++ {
		lastNumber += 1
		jokerValues = append(jokerValues, lastNumber)
	}

	for i, v := range jokerValues {
		pieces[i].Number = v
		pieces[i].Color = runColor
	}

	return true
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

func (g *Game) addPieceHandle(ar *ActionRequest) ([]byte, error) {
	return actionSuccess()
}

func (g *Game) removePiece(ar *ActionRequest) ([]byte, error) {
	return actionSuccess()
}

func (g *Game) replacePieceHandle(ar *ActionRequest) ([]byte, error) {
	return actionSuccess()
}

func (g *Game) addCombinationHandle(ar *ActionRequest) ([]byte, error) {
	return actionSuccess()
}

func (g *Game) splitCombination(ar *ActionRequest) ([]byte, error) {
	return actionSuccess()
}
