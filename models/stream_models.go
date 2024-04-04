package models

type StreamModel struct {
	Id      string `json:"id" redis:"id"`
	Brand   string `json:"brand" redis:"brand"`
	Name    string `json:"name" redis:"name"`
	Address string `json:"address" redis:"address"`

	MsFeederPid  int    `json:"ms_feeder_pid" redis:"ms_feeder_pid"`
	MSFeederArgs string `json:"ms_feeder_args" redis:"ms_feeder_args"`
	HlsPid       int    `json:"hls_pid" redis:"hls_pid"`
	HlsArgs      string `json:"hls_args" redis:"hls_args"`
	CreatedAt    string `json:"created_at" redis:"created_at"`

	MsType              int    `json:"ms_type" redis:"ms_type"`
	StreamType          int    `json:"stream_type" redis:"stream_type"`
	MsInitialized       bool   `json:"ms_initialized" redis:"ms_initialized"`
	MsImageName         string `json:"ms_image_name" redis:"ms_image_name"`
	MsContainerName     string `json:"ms_container_name" redis:"ms_container_name"`
	MsAddress           string `json:"ms_address" redis:"ms_address"`
	MsStreamAddress     string `json:"ms_stream_address" redis:"ms_stream_address"`
	MsContainerPorts    string `json:"ms_container_ports" redis:"ms_container_ports"`
	MsContainerCommands string `json:"ms_container_commands" redis:"ms_container_commands"`

	MpFfmpegReaderOwnerPid int `json:"mp_ffmpeg_reader_owner_pid" redis:"mp_ffmpeg_reader_owner_pid"`
	FfmpegReaderFrameRate  int `json:"ffmpeg_reader_frame_rate" redis:"ffmpeg_reader_frame_rate"`
	FfmpegReaderWidth      int `json:"ffmpeg_reader_width" redis:"ffmpeg_reader_width"`
	FfmpegReaderHeight     int `json:"ffmpeg_reader_height" redis:"ffmpeg_reader_height"`

	RecordEnabled  bool   `json:"record_enabled" redis:"record_enabled"`
	RecordPid      int    `json:"record_pid" redis:"record_pid"`
	RecordArgs     string `json:"record_args" redis:"record_args"`
	RecordDuration int    `json:"record_duration" redis:"record_duration"`

	SnapshotEnabled   bool `json:"snapshot_enabled" redis:"snapshot_enabled"`
	SnapshotPid       int  `json:"snapshot_pid" redis:"snapshot_pid"`
	SnapshotFrameRate int  `json:"snapshot_frame_rate" redis:"snapshot_frame_rate"`
	SnapshotWidth     int  `json:"snapshot_width" redis:"snapshot_width"`
	SnapshotHeight    int  `json:"snapshot_height" redis:"snapshot_height"`

	AiClipEnabled bool   `json:"ai_clip_enabled" redis:"ai_clip_enabled"`
	AiClipPid     int    `json:"ai_clip_pid" redis:"ai_clip_pid"`
	AiClipArgs    string `json:"ai_clip_args" redis:"ai_clip_args"`

	ConcatDemuxerPid  int    `json:"concat_demuxer_pid" redis:"concat_demuxer_pid"`
	ConcatDemuxerArgs string `json:"concat_demuxer_args" redis:"concat_demuxer_args"`

	RootDirPath string `json:"root_dir_path" redis:"root_dir_path"`

	FlvPlayerType            int  `json:"flv_player_type" redis:"flv_player_type"`
	BoosterEnabled           bool `json:"booster_enabled" redis:"booster_enabled"`
	LiveBufferLatencyChasing bool `json:"live_buffer_latency_chasing" redis:"live_buffer_latency_chasing"`

	Go2RtcPlayerMode int `json:"go2rtc_player_mode" redis:"go2rtc_player_mode"`
}

func (s *StreamModel) GetSourceId() string {
	return s.Id
}

func (s *StreamModel) GetDirPath() string {
	return s.RootDirPath
}
