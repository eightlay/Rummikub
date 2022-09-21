package game

import (
	"encoding/json"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
)

type combinationType string

const (
	notCombination combinationType = ""
	group          combinationType = "G"
	run            combinationType = "R"
)

type Combination struct {
	Pieces pack            `json:"pieces"`
	Type   combinationType `json:"type"`
}

func (c *Combination) toJSON() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("can't conver combination to json: %v", err)
	}
	return b, nil
}

func validInitialMeld(pieces []*Piece) *Combination {
	ct := validCombination(pieces)

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

func validCombination(pieces []*Piece) combinationType {
	validGroup := validGroup(pieces)

	if !validGroup {
		validRun := validRun(pieces)

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

func validGroup(pieces []*Piece) bool {
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

func validRun(pieces_ []*Piece) bool {
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
