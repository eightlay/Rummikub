package game

import "encoding/json"

type StateRequest struct {
	Player string `json:"player"`
}

type StateResponse struct {
	Hand     map[int][]byte `json:"hand"`
	Field    map[int][]byte `json:"field"`
	Actions  []byte         `json:"actions"`
	BankSize int            `json:"bankSize"`
}

type ActionRequest struct {
	Player           string `json:"player"`
	Action           action `json:"action"`
	UsedPieces       []int  `json:"usedPieces"`
	UsedCombinations []int  `json:"usedCombinations"`
	TimerExceeded    bool   `json:"timerExceeded"`
}

type ActionResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

var _actionSuccess ActionResponse = ActionResponse{true, nil}

func actionSuccess() ([]byte, error) {
	j, _ := json.Marshal(_actionSuccess)
	return j, nil
}

func actionError(err error) ([]byte, error) {
	response := ActionResponse{false, err}
	j, e := json.Marshal(&response)
	if e != nil {
		return nil, e
	}
	return j, nil
}
