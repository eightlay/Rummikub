package game

type hand pack

func createHands(players []string) map[player]hand {
	hands := map[player]hand{}

	for _, p := range players {
		hands[player(p)] = hand{}
	}

	return hands
}
