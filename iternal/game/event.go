// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package game

// Game event
type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

// Event error
type EventError struct {
	Error string `json:"error"`
}

// Event InitialMeld
type EventInit struct {
	Player string `json:"player"`
}

// Event InitialMeld
type EventInitialMeld struct {
	Player      player `json:"player"`
	AddedPieces []int  `json:"addedPieces"`
}

// Event AddPiece
type EventAddPiece struct {
	Player           player `json:"player"`
	AddedPieces      []int  `json:"addedPieces"`
	UsedCombinations []int  `json:"usedCombinations"`
}

// Event RemovePiece
type EventRemovePiece struct {
	Player           player `json:"player"`
	RemovedPiece     int    `json:"removedPiece"`
	UsedCombinations []int  `json:"usedCombinations"`
}

// Event ReplacePiece
type EventReplacePiece struct {
	Player           player `json:"player"`
	AddedPieces      []int  `json:"addedPieces"`
	RemovedPiece     int    `json:"removedPiece"`
	UsedCombinations []int  `json:"usedCombinations"`
}

// Event AddCombination
type EventAddCombination struct {
	Player      player `json:"player"`
	AddedPieces []int  `json:"addedPieces"`
}

// Event ConcatCombinations
type EventConcatCombinations struct {
	Player           player `json:"player"`
	UsedCombinations []int  `json:"usedCombinations"`
}

// Event SplitCombination
type EventSplitCombination struct {
	Player           player `json:"player"`
	SplitBeforeIndex int    `json:"splitAfterIndex"`
	UsedCombinations []int  `json:"usedCombinations"`
}

// Event Pass
type EventPass struct {
	Player player `json:"player"`
}

// Event ready
type EventReady struct {
	Player string `json:"player"`
}

// Event type
type EventType string

const (
	// Successfull event handling
	EventTypeSuccess EventType = "success"
	// Error during handling event
	EventTypeError EventType = "error"
	// Player joined the game
	EventTypeInit EventType = "init"
	// Other player connected to the game
	EventTypeConnect EventType = "connect"
	// Player disconnected from the game
	EventTypeDisconnect EventType = "disconnect"
	// Initial meld
	EventTypeInitialMeld EventType = "initialMeld"
	// Add piece to the existing combination
	EventTypeAddPiece EventType = "addPiece"
	// Remove piece from the existing combination
	EventTypeRemovePiece EventType = "removePiece"
	// Replace a piece from the hands with a piece from the combination
	EventTypeReplacePiece EventType = "replacePiece"
	// Add new combination to the game field
	EventTypeAddCombination EventType = "addCombination"
	// Concat two existing combinations
	EventTypeConcatCombinations EventType = "concatCombinations"
	// Split the existing combination
	EventTypeSplitCombination EventType = "splitCombination"
	// Pass
	EventTypePass EventType = "pass"
	// Ready to start
	EventTypeReady EventType = "ready"
)

// System events set
var systemEventsSet map[EventType]bool = map[EventType]bool{
	EventTypeInit:       true,
	EventTypeConnect:    true,
	EventTypeDisconnect: true,
}

// Initial events
var initialEvents [3]EventType = [3]EventType{
	EventTypeReady,
}

// Events available for initial meld stage
var initialMeldEvents [2]EventType = [2]EventType{EventTypeInitialMeld, EventTypePass}

// Events available for main game stage
var mainEvents [7]EventType = [7]EventType{
	EventTypeAddPiece, EventTypeRemovePiece, EventTypeReplacePiece,
	EventTypeAddCombination, EventTypeConcatCombinations,
	EventTypeSplitCombination, EventTypePass,
}
