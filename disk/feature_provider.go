package disk

import (
	"smcp/models"
)

type FeatureProvider[T any] interface {
	SetConfigAndModel(config *models.Config, model interface{})

	GetImagesPathBySourceId() string
	GetDataPathBySourceId() string

	GetBase64Image() *string
	GetCreatedAt() string
	GetSourceId() string
	GetFileName() string

	SetImageFileName(jsonObj *T, fullImageName string)
	SetDataFileName(jsonObj *T, fullDataName string)

	CreateJsonObject() *T
}
