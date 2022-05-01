package disk

import (
	"smcp/models"
	"smcp/utils"
)

type AlprProvider struct {
	config *models.Config
	model  *models.AlprResponse
}

func (a *AlprProvider) SetConfigAndModel(config *models.Config, model interface{}) {
	a.config = config
	a.model = model.(*models.AlprResponse)
}

func (a *AlprProvider) GetImagesPathBySourceId() string {
	return utils.GetAlprImagesPathBySourceId(a.config, a.model.SourceId)
}

func (a *AlprProvider) GetDataPathBySourceId() string {
	return utils.GetAlprDataPathBySourceId(a.config, a.model.SourceId)
}

func (a *AlprProvider) GetBase64Image() *string {
	return &a.model.Base64Image
}

func (a *AlprProvider) GetCreatedAt() string {
	return a.model.CreatedAt
}

func (a *AlprProvider) GetSourceId() string {
	return a.model.SourceId
}

func (a *AlprProvider) GetFileName() string {
	return a.model.CreateFileName()
}

func (a *AlprProvider) SetImageFileName(jsonObj *models.AlprJsonObject, fullImageName string) {
	jsonObj.AlprResults.ImageFileName = fullImageName
}

func (a *AlprProvider) SetDataFileName(jsonObj *models.AlprJsonObject, fullDataName string) {
	jsonObj.AlprResults.DataFileName = fullDataName
}

func (a *AlprProvider) CreateJsonObject() *models.AlprJsonObject {
	ar := a.model
	baseObj := models.AlprJsonBaseObject{
		ImgWidth: ar.ImgWidth, ImgHeight: ar.ImgHeight, ProcessingTimeMs: ar.ProcessingTimeMs,
		Results: ar.Results, Id: ar.Id, SourceId: ar.SourceId, CreatedAt: ar.CreatedAt, AiClipEnabled: ar.AiClipEnabled,
	}
	jsonObj := &models.AlprJsonObject{AlprResults: &baseObj, Video: &models.VideoClipJsonObject{}}
	return jsonObj
}
