package events

import (
	"encoding/json"
	"fmt"
)

type Pushing struct {
	SequenceId string `json:"sequence_id"`
	Command    string `json:"command"`
	Version    int    `json:"version"`
	PushTarget int    `json:"push_target"`
}
type PushRequest struct {
	Pushing Pushing `json:"pushing"`
}

func (p *PushRequest) SetSeq(id uint64) {
	p.Pushing.SequenceId = fmt.Sprintf("%5d", id)
}

func (p *PushRequest) IsRequestEvent() {
}

func (p *PushRequest) String() string {
	jsonBytes, _ := json.Marshal(p)
	return string(jsonBytes)
}
