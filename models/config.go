package models

type Config struct {
	Device struct {
		DeviceName string `json:"device_name"`
		DeviceArch int    `json:"device_arch"`
	} `json:"device"`
	General struct {
		DirPaths []string `json:"dir_paths"`
	} `json:"general"`
	Db struct {
		Type             int    `json:"type"` // 0 is SQLite, 1 is MongoDB
		ConnectionString string `json:"connection_string"`
	} `json:"db"`
	FFmpeg struct {
		UseDoubleQuotesForPath         bool    `json:"use_double_quotes_for_path"`
		MaxOperationRetryCount         int     `json:"max_operation_retry_count"`
		MsInitInterval                 float32 `json:"ms_init_interval"`
		WatchDogInterval               int     `json:"watch_dog_interval"`
		WatchDogFailedWaitInterval     float32 `json:"watch_dog_failed_wait_interval"`
		StartTaskWaitForInterval       float32 `json:"start_task_wait_for_interval"`
		RecordConcatLimit              int     `json:"record_concat_limit"`
		RecordVideoFileIndexerInterval int     `json:"record_video_file_indexer_interval"`
		MsPortStart                    int     `json:"ms_port_start"`
		MsPortEnd                      int     `json:"ms_port_end"`
	} `json:"ffmpeg"`
	Ai struct {
		VideoClipDuration int `json:"video_clip_duration"`
	} `json:"ai"`
	SenseAi struct {
		Image int    `json:"image"`
		Host  string `json:"host"`
		Port  int    `json:"port"`
	} `json:"sense_ai"`
	Jobs struct {
		MacIpMatchingEnabled       bool `json:"mac_ip_matching_enabled"`
		MacIpMatchingInterval      int  `json:"mac_ip_matching_interval"`
		BlackScreenMonitorEnabled  bool `json:"black_screen_monitor_enabled"`
		BlackScreenMonitorInterval int  `json:"black_screen_monitor_interval"`
	} `json:"jobs"`
	Archive struct {
		LimitPercent int    `json:"limit_percent"`
		ActionType   int    `json:"action_type"`
		MoveLocation string `json:"move_location"`
	} `json:"archive"`
	Snapshot struct {
		ProcessCount int  `json:"process_count"`
		Overlay      bool `json:"overlay"`
	} `json:"snapshot"`
	Hub struct {
		Enabled       bool   `json:"enabled"`
		Address       string `json:"address"`
		Token         string `json:"token"`
		WebAppAddress string `json:"web_app_address"`
		MaxRetry      int    `json:"max_retry"`
	} `json:"hub"`
}
