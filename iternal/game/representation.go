package game

import "encoding/json"

// Request game state for the player
type StateRequest struct {
	Player string `json:"player"`
}

// Game state for the player
type StateResponse struct {
	Hand     map[int][]byte `json:"hand"`
	Field    map[int][]byte `json:"field"`
	Actions  []byte         `json:"actions"`
	BankSize int            `json:"bankSize"`
}

// Request action handle
type ActionRequest struct {
	Player           string `json:"player"`
	Action           action `json:"action"`
	AddedPieces      []int  `json:"addedPieces"`
	RemovedPiece     int    `json:"removedPiece"`
	SplitBeforeIndex int    `json:"splitAfterIndex"`
	UsedCombinations []int  `json:"usedCombinations"`
	TimerExceeded    bool   `json:"timerExceeded"`
}

// Action handle result
type ActionResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

// Default success ActionResponse
var _actionSuccess ActionResponse = ActionResponse{true, nil}

// Get default success ActionResponse converted to JSON and nil error
func actionSuccess() ([]byte, error) {
	j, _ := json.Marshal(_actionSuccess)
	return j, nil
}

// Get error ActionResponse converted to JSON and error itself
func actionError(err error) ([]byte, error) {
	response := ActionResponse{false, err}
	j, e := json.Marshal(&response)
	if e != nil {
		return nil, e
	}
	return j, err
}
