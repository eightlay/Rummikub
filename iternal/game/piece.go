// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

import (
	"sort"
)

// Piece
//
// Contains information about piece's nubmer,
// color and flag is it joker or not
type Piece struct {
	Number int   `json:"number"`
	Color  color `json:"color"`
	Joker  bool  `json:"joker"`
}

// Create new piece
func createPiece(number int, color_ color, joker bool) *Piece {
	return &Piece{number, color_, joker}
}

// Sort the given pieces
func sortPieces(pieces_ []*Piece) []*Piece {
	pieces := pieces_[:]

	sort.Slice(pieces, func(i, j int) bool {
		return pieces[i].Number < pieces[j].Number
	})

	return pieces
}

// Set default joker parameters to the joker piece
func (p *Piece) clearIfJoker() {
	if p.Joker {
		p.Number = JokerNumber
		p.Color = JokerColor
	}
}
