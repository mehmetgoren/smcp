package models

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type Metadata struct {
	Colors []Color `json:"colors"`
}

type DetectedObject struct {
	PredScore   float32       `json:"pred_score"`
	PredClsIdx  int           `json:"pred_cls_idx"`
	PredClsName string        `json:"pred_cls_name"`
	Box         *DetectionBox `json:"box"`
	Metadata    *Metadata     `json:"metadata" bson:"metadata"`
}

type ObjectDetectionModel struct {
	Id              string            `json:"id"`
	SourceId        string            `json:"source_id"`
	CreatedAt       string            `json:"created_at"`
	DetectedObjects []*DetectedObject `json:"detected_objects"`
	Base64Image     string            `json:"base64_image"`
	AiClipEnabled   bool              `json:"ai_clip_enabled"`
}

func (d *ObjectDetectionModel) CreateFileName() string {
	return d.SourceId + "_" + d.CreatedAt + "_" + d.Id
}
