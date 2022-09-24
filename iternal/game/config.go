// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

const (
	// Number of pieces a player has at the beginning
	HandSize int = 14

	// Minimal number on a piece
	MinNumber int = 1
	// Maximal number on a piece
	MaxNumber int = 13

	// Number on a joker piece
	JokerNumber int = 0
	// Color of a joker piece
	JokerColor color = "jokerColor"

	// Number of decks in the bank at the beginning
	DecksNumber int = 2

	// Time limit for a move
	TimeLimitSeconds int = 60
	// Penalty size (in pieces) for exceeding time limit for a move
	PenaltySize int = 3

	// Initial meld sum minimal value
	InitialMeldSum int = 30

	// Minimal size of the group combination type
	MinGroupSize int = 3
	// Maximal size of the group combination type
	MaxGroupSize int = 4
	// Minimal size of the run combination type
	MinRunSize int = 3

	// Minimal number of players in the game
	MinPlayersNumber = 2
	// Maximal number of players in the game
	MaxPlayersNumber = 4
)
