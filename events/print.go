package events

import "encoding/json"

const PrintType = "print"

type PrintReport struct {
	Ams struct {
		Ams []struct {
			Humidity string `json:"humidity,omitempty"`
			Id       string `json:"id,omitempty"`
			Temp     string `json:"temp,omitempty"`
			Tray     []struct {
				Id            string   `json:"id,omitempty"`
				BedTemp       string   `json:"bed_temp,omitempty"`
				BedTempType   string   `json:"bed_temp_type,omitempty"`
				Cols          []string `json:"cols,omitempty"`
				DryingTemp    string   `json:"drying_temp,omitempty"`
				DryingTime    string   `json:"drying_time,omitempty"`
				NozzleTempMax string   `json:"nozzle_temp_max,omitempty"`
				NozzleTempMin string   `json:"nozzle_temp_min,omitempty"`
				Remain        int      `json:"remain,omitempty"`
				TagUid        string   `json:"tag_uid,omitempty"`
				TrayColor     string   `json:"tray_color,omitempty"`
				TrayDiameter  string   `json:"tray_diameter,omitempty"`
				TrayIdName    string   `json:"tray_id_name,omitempty"`
				TrayInfoIdx   string   `json:"tray_info_idx,omitempty"`
				TraySubBrands string   `json:"tray_sub_brands,omitempty"`
				TrayType      string   `json:"tray_type,omitempty"`
				TrayUuid      string   `json:"tray_uuid,omitempty"`
				TrayWeight    string   `json:"tray_weight,omitempty"`
				XcamInfo      string   `json:"xcam_info,omitempty"`
			} `json:"tray,omitempty"`
		} `json:"ams,omitempty"`
		AmsExistBits     string `json:"ams_exist_bits,omitempty"`
		InsertFlag       bool   `json:"insert_flag,omitempty"`
		PowerOnFlag      bool   `json:"power_on_flag,omitempty"`
		TrayExistBits    string `json:"tray_exist_bits,omitempty"`
		TrayIsBblBits    string `json:"tray_is_bbl_bits,omitempty"`
		TrayNow          string `json:"tray_now,omitempty"`
		TrayPre          string `json:"tray_pre,omitempty"`
		TrayReadDoneBits string `json:"tray_read_done_bits,omitempty"`
		TrayReadingBits  string `json:"tray_reading_bits,omitempty"`
		TrayTar          string `json:"tray_tar,omitempty"`
		Version          int    `json:"version,omitempty"`
	} `json:"ams"`
	AmsRfidStatus           int           `json:"ams_rfid_status,omitempty"`
	AmsStatus               int           `json:"ams_status,omitempty"`
	AuxPartFan              bool          `json:"aux_part_fan,omitempty"`
	BedTargetTemper         float64       `json:"bed_target_temper,omitempty"`
	BedTemper               float64       `json:"bed_temper,omitempty"`
	BigFan1Speed            string        `json:"big_fan1_speed,omitempty"`
	BigFan2Speed            string        `json:"big_fan2_speed,omitempty"`
	ChamberTemper           float64       `json:"chamber_temper,omitempty"`
	Command                 string        `json:"command,omitempty"`
	CoolingFanSpeed         string        `json:"cooling_fan_speed,omitempty"`
	FailReason              string        `json:"fail_reason,omitempty"`
	FanGear                 int           `json:"fan_gear,omitempty"`
	FilamBak                []interface{} `json:"filam_bak,omitempty"`
	ForceUpgrade            bool          `json:"force_upgrade,omitempty"`
	GcodeFile               string        `json:"gcode_file,omitempty"`
	GcodeFilePreparePercent string        `json:"gcode_file_prepare_percent,omitempty"`
	GcodeStartTime          string        `json:"gcode_start_time,omitempty"`
	GcodeState              string        `json:"gcode_state,omitempty"`
	HeatbreakFanSpeed       string        `json:"heatbreak_fan_speed,omitempty"`
	Hms                     []interface{} `json:"hms,omitempty"`
	HomeFlag                int           `json:"home_flag,omitempty"`
	HwSwitchState           int           `json:"hw_switch_state,omitempty"`
	Ipcam                   struct {
		IpcamDev    string `json:"ipcam_dev,omitempty"`
		IpcamRecord string `json:"ipcam_record,omitempty"`
		Resolution  string `json:"resolution,omitempty"`
		Timelapse   string `json:"timelapse,omitempty"`
	} `json:"ipcam,omitempty"`
	LayerNum     int    `json:"layer_num,omitempty"`
	Lifecycle    string `json:"lifecycle,omitempty"`
	LightsReport []struct {
		Mode string `json:"mode,omitempty"`
		Node string `json:"node,omitempty"`
	} `json:"lights_report,omitempty"`
	Maintain            int     `json:"maintain,omitempty"`
	McPercent           int     `json:"mc_percent,omitempty"`
	McPrintErrorCode    string  `json:"mc_print_error_code,omitempty"`
	McPrintStage        string  `json:"mc_print_stage,omitempty"`
	McPrintSubStage     int     `json:"mc_print_sub_stage,omitempty"`
	McRemainingTime     int     `json:"mc_remaining_time,omitempty"`
	MessProductionState string  `json:"mess_production_state,omitempty"`
	NozzleDiameter      string  `json:"nozzle_diameter,omitempty"`
	NozzleTargetTemper  float64 `json:"nozzle_target_temper,omitempty"`
	NozzleTemper        float64 `json:"nozzle_temper,omitempty"`
	Online              struct {
		Ahb     bool `json:"ahb,omitempty"`
		Rfid    bool `json:"rfid,omitempty"`
		Version int  `json:"version,omitempty"`
	} `json:"online,omitempty"`
	PrintError       int           `json:"print_error,omitempty"`
	PrintGcodeAction int           `json:"print_gcode_action,omitempty"`
	PrintRealAction  int           `json:"print_real_action,omitempty"`
	PrintType        string        `json:"print_type,omitempty"`
	ProfileId        string        `json:"profile_id,omitempty"`
	ProjectId        string        `json:"project_id,omitempty"`
	QueueNumber      int           `json:"queue_number,omitempty"`
	Sdcard           bool          `json:"sdcard,omitempty"`
	SequenceId       string        `json:"sequence_id,omitempty"`
	SpdLvl           int           `json:"spd_lvl,omitempty"`
	SpdMag           int           `json:"spd_mag,omitempty"`
	Stg              []interface{} `json:"stg,omitempty"`
	StgCur           int           `json:"stg_cur,omitempty"`
	SubtaskId        string        `json:"subtask_id,omitempty"`
	SubtaskName      string        `json:"subtask_name,omitempty"`
	TaskId           string        `json:"task_id,omitempty"`
	TotalLayerNum    int           `json:"total_layer_num,omitempty"`
	UpgradeState     struct {
		AhbNewVersionNumber string `json:"ahb_new_version_number,omitempty"`
		AmsNewVersionNumber string `json:"ams_new_version_number,omitempty"`
		ConsistencyRequest  bool   `json:"consistency_request,omitempty"`
		DisState            int    `json:"dis_state,omitempty"`
		ErrCode             int    `json:"err_code,omitempty"`
		ForceUpgrade        bool   `json:"force_upgrade,omitempty"`
		Message             string `json:"message,omitempty"`
		Module              string `json:"module,omitempty"`
		NewVersionState     int    `json:"new_version_state,omitempty"`
		OtaNewVersionNumber string `json:"ota_new_version_number,omitempty"`
		Progress            string `json:"progress,omitempty"`
		SequenceId          int    `json:"sequence_id,omitempty"`
		Status              string `json:"status,omitempty"`
	} `json:"upgrade_state,omitempty"`
	Upload struct {
		FileSize      int    `json:"file_size,omitempty"`
		FinishSize    int    `json:"finish_size,omitempty"`
		Message       string `json:"message,omitempty"`
		OssUrl        string `json:"oss_url,omitempty"`
		Progress      int    `json:"progress,omitempty"`
		SequenceId    string `json:"sequence_id,omitempty"`
		Speed         int    `json:"speed,omitempty"`
		Status        string `json:"status,omitempty"`
		TaskId        string `json:"task_id,omitempty"`
		TimeRemaining int    `json:"time_remaining,omitempty"`
		TroubleId     string `json:"trouble_id,omitempty"`
	} `json:"upload"`
	VtTray struct {
		BedTemp       string   `json:"bed_temp,omitempty"`
		BedTempType   string   `json:"bed_temp_type,omitempty"`
		Cols          []string `json:"cols,omitempty"`
		DryingTemp    string   `json:"drying_temp,omitempty"`
		DryingTime    string   `json:"drying_time,omitempty"`
		Id            string   `json:"id,omitempty"`
		NozzleTempMax string   `json:"nozzle_temp_max,omitempty"`
		NozzleTempMin string   `json:"nozzle_temp_min,omitempty"`
		Remain        int      `json:"remain,omitempty"`
		TagUid        string   `json:"tag_uid,omitempty"`
		TrayColor     string   `json:"tray_color,omitempty"`
		TrayDiameter  string   `json:"tray_diameter,omitempty"`
		TrayIdName    string   `json:"tray_id_name,omitempty"`
		TrayInfoIdx   string   `json:"tray_info_idx,omitempty"`
		TraySubBrands string   `json:"tray_sub_brands,omitempty"`
		TrayType      string   `json:"tray_type,omitempty"`
		TrayUuid      string   `json:"tray_uuid,omitempty"`
		TrayWeight    string   `json:"tray_weight,omitempty"`
		XcamInfo      string   `json:"xcam_info,omitempty"`
	} `json:"vt_tray,omitempty"`
	WifiSignal string `json:"wifi_signal,omitempty"`
	Xcam       struct {
		AllowSkipParts           bool   `json:"allow_skip_parts,omitempty"`
		BuildplateMarkerDetector bool   `json:"buildplate_marker_detector,omitempty"`
		FirstLayerInspector      bool   `json:"first_layer_inspector,omitempty"`
		HaltPrintSensitivity     string `json:"halt_print_sensitivity,omitempty"`
		PrintHalt                bool   `json:"print_halt,omitempty"`
		PrintingMonitor          bool   `json:"printing_monitor,omitempty"`
		SpaghettiDetector        bool   `json:"spaghetti_detector,omitempty"`
	} `json:"xcam,omitempty"`
	XcamStatus string `json:"xcam_status,omitempty"`
}

func (p *PrintReport) IsReportEvent() {
}

func (p *PrintReport) GetType() string {

	return PrintType
}

func (i *PrintReport) String() string {
	jsonBytes, _ := json.MarshalIndent(i, "", "  ")
	return string(jsonBytes)
}
