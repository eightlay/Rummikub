// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

import (
	mapset "github.com/deckarep/golang-set/v2"
)

// Combination type
//
// Existing types: group, run
type combinationType string

const (
	// Group combination type
	group combinationType = "G"
	// Run combination type
	run combinationType = "R"
)

// Combination
//
// Contains information about used pieces and combination type
type Combination struct {
	Pieces pack            `json:"pieces"`
	Type   combinationType `json:"type"`
}

// Returns combination if provided pieces present valid initial meld
func validInitialMeld(pieces []*Piece) *Combination {
	newCombination := validCombination(pieces)

	if newCombination != nil {
		s := 0

		for _, p := range newCombination.Pieces {
			s += p.Number
		}

		correct := s >= InitialMeldSum

		if correct {
			return newCombination
		}
	}

	return nil
}

// Return combination if provided pieces present valid combination
func validCombination(pieces []*Piece) *Combination {
	validGroup := isValidGroup(pieces)

	if !validGroup {
		validRun := isValidRun(pieces)

		if !validRun {

			for _, p := range pieces {
				p.clearIfJoker()
			}

			return nil
		} else {
			return &Combination{
				Pieces: sortPieces(pieces),
				Type:   run,
			}
		}
	}

	return &Combination{
		Pieces: sortPieces(pieces),
		Type:   group,
	}
}

// Check if provided pieces present valid group
func isValidGroup(pieces []*Piece) bool {
	if len(pieces) < MinGroupSize || len(pieces) > MaxGroupSize {
		return false
	}

	usedColors := mapset.NewSet[color]()
	var number int = JokerNumber

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

// Check if provided pieces present valid run
func isValidRun(pieces_ []*Piece) bool {
	if len(pieces_) < MinRunSize {
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

	for i := 0; i < jokerCount; i++ {
		lastNumber += 1
		jokerValues = append(jokerValues, lastNumber)
	}

	for i, v := range jokerValues {
		pieces[i].Number = v
		pieces[i].Color = runColor
	}

	return true
}
