// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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
