package game

// Player's hand
type hand pack

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
