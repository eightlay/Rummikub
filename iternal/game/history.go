// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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
