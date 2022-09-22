package game

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
func (p hand) largestPieceNumber() int {
	largest := MinNumber - 1
	if JokerNumber < largest {
		largest = JokerNumber
	}

	for _, v := range p {
		if v.Number > largest {
			largest = v.Number
		}
	}

	return largest
}
