package game

import "encoding/json"

type State struct {
	Turn bool `jso:"turn"`
	// Field           field       `json:"field"`
	Hand            hand        `json:"hand"`
	AvailableEvents []EventType `json:"availableEvents"`
	Started         bool        `json:"started"`
	Finished        bool        `json:"finished"`
	Winner          player      `json:"winner"`
	Error           string      `json:"error"`
}

func (s State) ToJSON() []byte {
	bytes, _ := json.Marshal(s)
	return bytes
}
