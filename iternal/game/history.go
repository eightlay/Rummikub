package game

// Game history
//
// Store information about game progress
type history struct {
	combinations field
	firstStep    *step
	lastStep     *step
}

// Create history
func createHistory() *history {
	s := &step{0, "", nil, nil}

	return &history{
		combinations: field{},
		firstStep:    s,
		lastStep:     s,
	}
}
