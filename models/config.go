package models

type Config struct {
	DeviceConfig struct {
		DeviceName     string `json:"device_name"`
		DeviceType     int    `json:"device_type"`
		DeviceServices []int  `json:"device_services"`
	} `json:"device"`
	HeartbeatConfig struct {
		Interval int `json:"interval"`
	} `json:"heartbeat"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Jetson struct {
		ModelName string  `json:"model_name"`
		Threshold float64 `json:"threshold"`
		WhiteList []int   `json:"white_list"`
	} `json:"jetson"`
	Torch struct {
		ModelName         string  `json:"model_name"`
		ModelNameSpecific string  `json:"model_name_specific"`
		Threshold         float64 `json:"threshold"`
		WhiteList         []int   `json:"white_list"`
	} `json:"torch"`
	OnceDetector struct {
		ImagehashThreshold int     `json:"imagehash_threshold"`
		PsnrThreshold      float64 `json:"psnr_threshold"`
		SsimThreshold      float64 `json:"ssim_threshold"`
	} `json:"once_detector"`
	Handler struct {
		SaveImageFolderPath string `json:"save_image_folder_path"`
		SaveImageExtension  string `json:"save_image_extension"`
		ShowImageWaitKey    int    `json:"show_image_wait_key"`
		ShowImageCaption    bool   `json:"show_image_caption"`
		ShowImageFullscreen bool   `json:"show_image_fullscreen"`
		ReadServiceOverlay  bool   `json:"read_service_overlay"`
	} `json:"handler"`
	SourceReader struct {
		Fps        int `json:"fps"`
		BufferSize int `json:"buffer_size"`
		MaxRetry   int `json:"max_retry"`
		MaxRetryIn int `json:"max_retry_in"`
	} `json:"source_reader"`
	Path struct {
		Stream string `json:"stream"`
		Record string `json:"record"`
		Read   string `json:"read"`
	} `json:"path"`
	FFmpeg struct {
		UseDoubleQuotesForPath                  bool    `json:"use_double_quotes_for_path"`
		MaxOperationRetryCount                  int     `json:"max_operation_retry_count"`
		CheckLeakyFfmpegProcessesInterval       int     `json:"check_leaky_ffmpeg_processes_interval"`
		CheckUnstoppedContainersInterval        int     `json:"check_unstopped_containers_interval"`
		CheckFfmpegStreamRunningProcessInterval int     `json:"check_ffmpeg_stream_running_process_interval"`
		CheckFfmpegRecordRunningProcessInterval int     `json:"check_ffmpeg_record_running_process_interval"`
		StartTaskWaitForInterval                float32 `json:"start_task_wait_for_interval"`
		EventListenerHandlerType                int     `json:"event_listener_handler_type"`
	} `json:"ffmpeg"`
	AiConfig struct {
		DetectedFolder string `json:"detected_folder"`
	} `json:"ai"`
}
