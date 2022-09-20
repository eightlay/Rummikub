package game

import (
	"encoding/json"
	"fmt"
)

type combinationType string

const (
	notCombination combinationType = ""
	group          combinationType = "G"
	run            combinationType = "R"
)

type Combination struct {
	Pieces []*Piece        `json:"pieces"`
	Type   combinationType `json:"type"`
}

func (c *Combination) toJSON() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("can't conver combination to json: %v", err)
	}
	return b, nil
}
