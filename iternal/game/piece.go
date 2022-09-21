package game

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Piece struct {
	Number int   `json:"number"`
	Color  color `json:"color"`
	Joker  bool  `json:"joker"`
}

func createPiece(number int, color_ color, joker bool) *Piece {
	return &Piece{number, color_, joker}
}

func (p *Piece) toJSON() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("can't conver combination to json: %v", err)
	}
	return b, nil
}

func sortPieces(pieces_ []*Piece) []*Piece {
	pieces := pieces_[:]

	sort.Slice(pieces, func(i, j int) bool {
		return pieces[i].Number < pieces[j].Number
	})

	return pieces
}

func (p *Piece) clearIfJoker() {
	if p.Joker {
		p.Number = JokerNumber
		p.Color = JokerColor
	}
}
