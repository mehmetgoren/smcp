package data

import (
	"smcp/models"
)

type Repository interface {
	OdSave(od *models.ObjectDetectionModel) error
	FrSave(fr *models.FaceRecognitionModel) error
	AlprSave(alpr *models.AlprResponse) error

	SetAiClipFields(groupId string, clip *AiClip) error

	SetVideoFields(params *SetVideoFileParams) error

	SetVideoFieldsMerged(params *SetVideoFileMergeParams) error
}
