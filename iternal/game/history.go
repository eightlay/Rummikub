package game

type history struct {
	combinations field
	firstStep    *step
	lastStep     *step
}

func createHistory() *history {
	s := &step{0, "", nil, nil}

	return &history{
		combinations: field{},
		firstStep:    s,
		lastStep:     s,
	}
}
