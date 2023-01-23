package models

type DetectedFace struct {
	PredScore   float32       `json:"pred_score"`
	PredClsIdx  int           `json:"pred_cls_idx"`
	PredClsName string        `json:"pred_cls_name"`
	Box         *DetectionBox `json:"box"`
}

type FaceRecognitionModel struct {
	Id            string          `json:"id"`
	SourceId      string          `json:"source_id"`
	CreatedAt     string          `json:"created_at"`
	DetectedFaces []*DetectedFace `json:"detected_faces"`
	Base64Image   string          `json:"base64_image"`
	AiClipEnabled bool            `json:"ai_clip_enabled"`
}

func (f *FaceRecognitionModel) CreateFileName() string {
	return f.SourceId + "_" + f.CreatedAt + "_" + f.Id
}
