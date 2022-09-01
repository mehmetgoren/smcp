package models

const (
	Od   = 0
	Fr   = 1
	Alpr = 2
)

type AiClipQueueModel struct {
	AiType int                   `json:"ai_type"`
	Od     *ObjectDetectionModel `json:"od"`
	Fr     *FaceRecognitionModel `json:"fr"`
	Alpr   *AlprResponse         `json:"alpr"`
}
