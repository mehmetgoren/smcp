package models

type DetectedObject struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
}

type ObjectDetectionModel struct {
	Id               string            `json:"id"`
	SourceId         string            `json:"source_id"`
	CreatedAt        string            `json:"created_at"`
	DetectedObjects  []*DetectedObject `json:"detected_objects"`
	Base64Image      string            `json:"base64_image"`
	VideoClipEnabled bool              `json:"video_clip_enabled"`
}

func (d *ObjectDetectionModel) CreateFileName() string {
	return d.SourceId + "_" + d.CreatedAt + "_" + d.Id
}
