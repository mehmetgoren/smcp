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
		ModelName string `json:"model_name"`
	} `json:"jetson"`
	Torch struct {
		ModelName         string `json:"model_name"`
		ModelNameSpecific string `json:"model_name_specific"`
	} `json:"torch"`
	Tensorflow struct {
		ModelName   string `json:"model_name"`
		CacheFolder string `json:"cache_folder"`
	} `json:"tensorflow"`
	OnceDetector struct {
		ImagehashThreshold int     `json:"imagehash_threshold"`
		PsnrThreshold      float64 `json:"psnr_threshold"`
		SsimThreshold      float64 `json:"ssim_threshold"`
	} `json:"once_detector"`
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
		UseDoubleQuotesForPath     bool    `json:"use_double_quotes_for_path"`
		MaxOperationRetryCount     int     `json:"max_operation_retry_count"`
		RtmpServerInitInterval     float32 `json:"rtmp_server_init_interval"`
		WatchDogInterval           int     `json:"watch_dog_interval"`
		WatchDogFailedWaitInterval float32 `json:"watch_dog_failed_wait_interval"`
		StartTaskWaitForInterval   float32 `json:"start_task_wait_for_interval"`
	} `json:"ffmpeg"`
	AiConfig struct {
		ReadServiceOverlay bool   `json:"read_service_overlay"`
		DetectedFolder     string `json:"detected_folder"`
	} `json:"ai"`
}
