package data

import (
	"smcp/models"
)

type Repository interface {
	AiSave(ai *models.AiDetectionModel) error

	SetAiClipFields(groupId string, clip *AiClip) error

	SetVideoFields(params *SetVideoFileParams) error

	SetVideoFieldsMerged(params *SetVideoFileMergeParams) error
}
