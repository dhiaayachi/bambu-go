package events

import "encoding/json"

const SystemType = "system"

type LedCtlReport struct {
	SequenceId string `json:"sequence_id"`
	Command    string `json:"command"`
	LedNode    string `json:"led_node"`
	LedMode    string `json:"led_mode"`
	Reason     string `json:"reason"`
	Result     string `json:"result"`
}

func (l LedCtlReport) GetType() string {
	return SystemType
}

func (l LedCtlReport) IsReportEvent() {
}

func (l LedCtlReport) String() string {
	jsonBytes, _ := json.MarshalIndent(l, "", "  ")
	return string(jsonBytes)
}

type System struct {
	SequenceId   string `json:"sequence_id"`
	Command      string `json:"command"`
	LedNode      string `json:"led_node"`
	LedMode      string `json:"led_mode"`
	LedOnTime    int    `json:"led_on_time"`
	LedOffTime   int    `json:"led_off_time"`
	LoopTimes    int    `json:"loop_times"`
	IntervalTime int    `json:"interval_time"`
}

type LedCtlRequest struct {
	System System `json:"system"`
}

func (l *LedCtlRequest) String() string {
	jsonBytes, _ := json.Marshal(l)
	return string(jsonBytes)
}

func (l *LedCtlRequest) IsRequestEvent() {
}
