package models

type StreamModel struct {
	Id      string `json:"id" redis:"id"`
	Brand   string `json:"brand" redis:"brand"`
	Name    string `json:"name" redis:"name"`
	Address string `json:"address" redis:"address"`

	RtmpFeederPid  int    `json:"rtmp_feeder_pid" redis:"rtmp_feeder_pid"`
	RtmpFeederArgs string `json:"rtmp_feeder_args" redis:"rtmp_feeder_args"`
	HlsPid         int    `json:"hls_pid" redis:"hls_pid"`
	HlsArgs        string `json:"hls_args" redis:"hls_args"`
	CreatedAt      string `json:"created_at" redis:"created_at"`

	StreamType            int    `json:"stream_type" redis:"stream_type"`
	RtmpServerInitialized bool   `json:"rtmp_server_initialized" redis:"rtmp_server_initialized"`
	RtmpServerType        int    `json:"rtmp_server_type" redis:"rtmp_server_type"`
	RtmpImageName         string `json:"rtmp_image_name" redis:"rtmp_image_name"`
	RtmpContainerName     string `json:"rtmp_container_name" redis:"rtmp_container_name"`
	RtmpAddress           string `json:"rtmp_address" redis:"rtmp_address"`
	RtmpFlvAddress        string `json:"rtmp_flv_address" redis:"rtmp_flv_address"`
	RtmpContainerPorts    string `json:"rtmp_container_ports" redis:"rtmp_container_ports"`
	RtmpContainerCommands string `json:"rtmp_container_commands" redis:"rtmp_container_commands"`

	FfmpegReaderPid       int `json:"ffmpeg_reader_pid" redis:"ffmpeg_reader_pid"`
	FfmpegReaderFrameRate int `json:"ffmpeg_reader_frame_rate" redis:"ffmpeg_reader_frame_rate"`
	FfmpegReaderWidth     int `json:"ffmpeg_reader_width" redis:"ffmpeg_reader_width"`
	FfmpegReaderHeight    int `json:"ffmpeg_reader_height" redis:"ffmpeg_reader_height"`

	RecordEnabled  bool   `json:"record_enabled" redis:"record_enabled"`
	RecordPid      int    `json:"record_pid" redis:"record_pid"`
	RecordArgs     string `json:"record_args" redis:"record_args"`
	RecordDuration int    `json:"record_duration" redis:"record_duration"`

	SnapshotEnabled   bool `json:"snapshot_enabled" redis:"snapshot_enabled"`
	SnapshotPid       int  `json:"snapshot_pid" redis:"snapshot_pid"`
	SnapshotFrameRate int  `json:"snapshot_frame_rate" redis:"snapshot_frame_rate"`
	SnapshotWidth     int  `json:"snapshot_width" redis:"snapshot_width"`
	SnapshotHeight    int  `json:"snapshot_height" redis:"snapshot_height"`

	VideoClipEnabled bool `json:"video_clip_enabled" redis:"video_clip_enabled"`

	HlsOutputPath          string `json:"hls_output_path" redis:"hls_output_path"`
	RecordOutputFolderPath string `json:"record_output_folder_path" redis:"record_output_folder_path"`
}
