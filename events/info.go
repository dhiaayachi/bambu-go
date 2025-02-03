package events

import "encoding/json"

const InfoType = "info"

type InfoReport struct {
	Command string `json:"command"`
	Module  []struct {
		HwVer string `json:"hw_ver"`
		Name  string `json:"name"`
		Sn    string `json:"sn"`
		SwVer string `json:"sw_ver"`
	} `json:"module"`
	SequenceId string `json:"sequence_id"`
}

func (i *InfoReport) IsReportEvent() {
}

func (i *InfoReport) GetType() string {
	return InfoType
}

func (i *InfoReport) String() string {
	jsonBytes, _ := json.MarshalIndent(i, "", "  ")
	return string(jsonBytes)
}
