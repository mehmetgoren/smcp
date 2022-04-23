package models

type DetectedFace struct {
	PredScore       float32 `json:"pred_score"`
	PredClsIdx      int     `json:"pred_cls_idx"`
	PredClsName     string  `json:"pred_cls_name"`
	CropBase64Image string  `json:"crop_base64_image"`
	X1              float32 `json:"x1"`
	Y1              float32 `json:"y1"`
	X2              float32 `json:"x2"`
	Y2              float32 `json:"y2"`
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
