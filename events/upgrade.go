package events

import "encoding/json"

const UpgrateType = "upgrade"

type UpgradeReport struct {
	Command          string `json:"command"`
	FirmwareOptional []struct {
		Ams []struct {
			Address      int    `json:"address"`
			DevModelName string `json:"dev_model_name"`
			DeviceId     string `json:"device_id"`
			Firmware     []struct {
				Description string `json:"description"`
				ForceUpdate bool   `json:"force_update"`
				Url         string `json:"url"`
				Version     string `json:"version"`
			} `json:"firmware"`
			FirmwareCurrent interface{} `json:"firmware_current"`
		} `json:"ams"`
		Firmware struct {
			Description string `json:"description"`
			ForceUpdate bool   `json:"force_update"`
			Url         string `json:"url"`
			Version     string `json:"version"`
		} `json:"firmware"`
	} `json:"firmware_optional"`
	Reason     string `json:"reason"`
	Result     string `json:"result"`
	SequenceId string `json:"sequence_id"`
}

func (u *UpgradeReport) IsReportEvent() {
}

func (u *UpgradeReport) GetType() string {
	return UpgrateType
}

func (i *UpgradeReport) String() string {
	jsonBytes, _ := json.MarshalIndent(i, "", "  ")
	return string(jsonBytes)
}
