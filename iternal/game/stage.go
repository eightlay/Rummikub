package game

// Stage number
type stage int

const (
	// System stage
	systemStage stage = iota
	// Initial meld stage
	initialMeldStage
	// Main game stage
	mainGameStage
)

// Get available event for the stage
func (s stage) availableEvents() []EventType {
	if s == systemStage {
		return initialEvents[:]
	}
	if s == initialMeldStage {
		return initialMeldEvents[:]
	}
	return mainEvents[:]
}
