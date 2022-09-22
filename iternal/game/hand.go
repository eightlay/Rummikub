package game

import "fmt"

// Player's hand
type hand pack

// Create hands for players
func createHands(players []player) map[player]hand {
	hands := map[player]hand{}

	for _, p := range players {
		hands[p] = hand{}
	}

	return hands
}

// Find largest piece number in the hand
func (h hand) largestPieceNumber() int {
	largest := MinNumber - 1
	if JokerNumber < largest {
		largest = JokerNumber
	}

	for _, v := range h {
		if v.Number > largest {
			largest = v.Number
		}
	}

	return largest
}

// Find largest piece number in the hand
func (h hand) toJSON() (map[int][]byte, error) {
	result := map[int][]byte{}

	for i, piece := range h {
		j, err := piece.toJSON()
		if err != nil {
			return nil, fmt.Errorf("can't get current state: %v", err)
		}

		result[i] = j
	}

	return result, nil
}
