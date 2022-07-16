package data

import (
	"smcp/models"
)

type Repository interface {
	OdSave(od *models.ObjectDetectionModel) error
	FrSave(fr *models.FaceRecognitionModel) error
	AlprSave(alpr *models.AlprResponse) error

	SetOdVideoClipFields(groupId string, clip *AiClip) error

	SetVideoFileNames(params *SetVideoFileNameParams) error
}
