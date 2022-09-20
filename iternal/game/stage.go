package game

type stage int

const (
	initialMeldStage stage = iota
	mainGameStage
)

func createStages(players []string) map[player]stage {
	hands := map[player]stage{}

	for _, p := range players {
		hands[player(p)] = initialMeldStage
	}

	return hands
}
