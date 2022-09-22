package game

// Stage number
type stage int

const (
	// Initial meld stage
	initialMeldStage stage = iota
	// Main game stage
	mainGameStage
)

func (s stage) availableActions() []action {
	if s == initialMeldStage {
		return initialActions[:]
	}
	return mainActions[:]
}

// Create stage storage for players
func createStages(players []player) map[player]stage {
	hands := map[player]stage{}

	for _, p := range players {
		hands[p] = initialMeldStage
	}

	return hands
}
