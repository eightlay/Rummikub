package game

import (
	"encoding/json"
	"fmt"
)

// Request game state for the player
type StateRequest struct {
	Player string `json:"player"`
}

// Parse StateRequest from JSON
func ParseStateRequest(request []byte) (*StateRequest, error) {
	// Parse request
	sr := &StateRequest{}

	err := json.Unmarshal(request, sr)
	if err != nil {
		return nil, fmt.Errorf("can't parse state request: %v", err)
	}

	return sr, nil
}

// Convert StateRequest to JSON
func (sr *StateRequest) ToJSON() ([]byte, error) {
	j, err := json.Marshal(sr)
	if err != nil {
		return nil, fmt.Errorf("can't convert current state: %v", err)
	}
	return j, nil
}

// Game state for the player
type StateResponse struct {
	Hand     map[int][]byte `json:"hand"`
	Field    map[int][]byte `json:"field"`
	Actions  []byte         `json:"actions"`
	BankSize int            `json:"bankSize"`
}

// Parse StateResponse from JSON
func ParseStateResponse(request []byte) (*StateResponse, error) {
	sr := &StateResponse{}

	err := json.Unmarshal(request, sr)
	if err != nil {
		return nil, fmt.Errorf("can't parse state response: %v", err)
	}

	return sr, nil
}

// Convert StateResponse to JSON
func (sr *StateResponse) ToJSON() ([]byte, error) {
	j, err := json.Marshal(&sr)
	if err != nil {
		return nil, fmt.Errorf("can't convert  state response: %v", err)
	}

	return j, nil
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

// Parse ActionRequest from JSON
func ParseActionRequest(request []byte) (*ActionRequest, error) {
	sr := &ActionRequest{}

	err := json.Unmarshal(request, sr)
	if err != nil {
		return nil, fmt.Errorf("can't parse state response: %v", err)
	}

	return sr, nil
}

// Convert ActionRequest to JSON
func (sr *ActionRequest) ToJSON() ([]byte, error) {
	j, err := json.Marshal(&sr)
	if err != nil {
		return nil, fmt.Errorf("can't convert  state response: %v", err)
	}

	return j, nil
}

// Action handle result
type ActionResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

// Parse ActionResponse from JSON
func ParseActionResponse(request []byte) (*ActionResponse, error) {
	sr := &ActionResponse{}

	err := json.Unmarshal(request, sr)
	if err != nil {
		return nil, fmt.Errorf("can't parse state response: %v", err)
	}

	return sr, nil
}

// Convert ActionResponse to JSON
func (sr *ActionResponse) ToJSON() ([]byte, error) {
	j, err := json.Marshal(&sr)
	if err != nil {
		return nil, fmt.Errorf("can't convert  state response: %v", err)
	}

	return j, nil
}

// Default success ActionResponse
var _actionSuccess ActionResponse = ActionResponse{true, nil}

// Get default success ActionResponse converted to JSON and nil error
func actionSuccess() ([]byte, error) {
	return _actionSuccess.ToJSON()
}

// Get error ActionResponse converted to JSON and error itself
func actionError(err error) ([]byte, error) {
	response := ActionResponse{false, err}
	return response.ToJSON()
}
