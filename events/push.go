package events

import "encoding/json"

type PushRequest struct {
	SequenceId string `json:"sequence_id"`
	Command    string `json:"command"`
	Version    int    `json:"version"`
	PushTarget int    `json:"push_target"`
}

func (p PushRequest) IsRequestEvent() {
}

func (p PushRequest) String() string {
	jsonBytes, _ := json.Marshal(p)
	return string(jsonBytes)
}
