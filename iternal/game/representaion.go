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
	Error error `json:"error"`
}

var _actionSuccess ActionResponse = ActionResponse{nil}

func actionSuccess() ([]byte, error) {
	j, _ := json.Marshal(_actionSuccess)
	return j, nil
}

func actionError(err error) ([]byte, error) {
	j, e := json.Marshal(_actionSuccess)
	if e != nil {
		return nil, e
	}
	return j, nil
}
