package models

type FaceRecognitionJsonBaseObject struct {
	Id            string          `json:"id"`
	SourceId      string          `json:"source_id"`
	CreatedAt     string          `json:"created_at"`
	DetectedFaces []*DetectedFace `json:"detected_faces"`
	AiClipEnabled bool            `json:"ai_clip_enabled"`
	ImageFileName string          `json:"image_file_name"`
	DataFileName  string          `json:"data_file_name"`
}

type FaceRecognitionJsonObject struct {
	FaceRecognition *FaceRecognitionJsonBaseObject `json:"face_detection"`
	Video           *VideoClipJsonObject           `json:"video"`
}
