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
