package models

type AiDetectionModel struct {
	Id            string               `json:"id"`
	Module        string               `json:"module"`
	SourceId      string               `json:"source_id"`
	CreatedAt     string               `json:"created_at"`
	Detections    []*AiDetectionResult `json:"detections"`
	Base64Image   string               `json:"base64_image"`
	AiClipEnabled bool                 `json:"ai_clip_enabled"`
}

func (d *AiDetectionModel) CreateFileName() string {
	return d.SourceId + "_" + d.Module + "_" + d.CreatedAt + "_" + d.Id
}
