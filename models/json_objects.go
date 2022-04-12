package models

type ObjectDetectionJsonBaseObject struct {
	Id               string            `json:"id"`
	SourceId         string            `json:"source_id"`
	CreatedAt        string            `json:"created_at"`
	DetectedObjects  []*DetectedObject `json:"detected_objects"`
	VideoClipEnabled bool              `json:"video_clip_enabled"`
	ImageFileName    string            `json:"image_file_name"`
	DataFileName     string            `json:"data_file_name"`
}

type VideoClipJsonObject struct {
	FileName       string `json:"file_name"`
	CreatedAt      string `json:"created_at"`
	LastModifiedAt string `json:"last_modified_at"`
	Duration       int    `json:"duration"`
}

type ObjectDetectionJsonObject struct {
	ObjectDetection *ObjectDetectionJsonBaseObject `json:"object_detection"`
	Video           *VideoClipJsonObject           `json:"video"`
}
