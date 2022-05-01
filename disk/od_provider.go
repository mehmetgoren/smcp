package disk

import (
	"smcp/models"
	"smcp/utils"
)

type OdProvider struct {
	config *models.Config
	model  *models.ObjectDetectionModel
}

func (o *OdProvider) SetConfigAndModel(config *models.Config, model interface{}) {
	o.config = config
	o.model = model.(*models.ObjectDetectionModel)
}

func (o *OdProvider) GetImagesPathBySourceId() string {
	return utils.GetOdImagesPathBySourceId(o.config, o.model.SourceId)
}

func (o *OdProvider) GetDataPathBySourceId() string {
	return utils.GetOdDataPathBySourceId(o.config, o.model.SourceId)
}

func (o *OdProvider) GetBase64Image() *string {
	return &o.model.Base64Image
}

func (o *OdProvider) GetCreatedAt() string {
	return o.model.CreatedAt
}

func (o *OdProvider) GetSourceId() string {
	return o.model.SourceId
}

func (o *OdProvider) GetFileName() string {
	return o.model.CreateFileName()
}

func (o *OdProvider) SetImageFileName(jsonObj *models.ObjectDetectionJsonObject, fullImageName string) {
	jsonObj.ObjectDetection.ImageFileName = fullImageName
}

func (o *OdProvider) SetDataFileName(jsonObj *models.ObjectDetectionJsonObject, fullDataName string) {
	jsonObj.ObjectDetection.DataFileName = fullDataName
}

func (o *OdProvider) CreateJsonObject() *models.ObjectDetectionJsonObject {
	de := o.model
	baseObj := models.ObjectDetectionJsonBaseObject{Id: de.Id, SourceId: de.SourceId, CreatedAt: de.CreatedAt,
		DetectedObjects: de.DetectedObjects, AiClipEnabled: de.AiClipEnabled}
	jsonObj := &models.ObjectDetectionJsonObject{ObjectDetection: &baseObj, Video: &models.VideoClipJsonObject{}}
	return jsonObj
}
