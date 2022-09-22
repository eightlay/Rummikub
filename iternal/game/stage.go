package game

// Stage number
type stage int

const (
	// Initial meld stage
	initialMeldStage stage = iota
	// Main game stage
	mainGameStage
)

// Create stage storage for players
func createStages(players []string) map[player]stage {
	hands := map[player]stage{}

	for _, p := range players {
		hands[player(p)] = initialMeldStage
	}

	return hands
}
