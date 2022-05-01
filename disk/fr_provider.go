package disk

import (
	"smcp/models"
	"smcp/utils"
)

type FrProvider struct {
	config *models.Config
	model  *models.FaceRecognitionModel
}

func (f *FrProvider) SetConfigAndModel(config *models.Config, model interface{}) {
	f.config = config
	f.model = model.(*models.FaceRecognitionModel)
}

func (f *FrProvider) GetImagesPathBySourceId() string {
	return utils.GetFrImagesPathBySourceId(f.config, f.model.SourceId)
}

func (f *FrProvider) GetDataPathBySourceId() string {
	return utils.GetFrDataPathBySourceId(f.config, f.model.SourceId)
}

func (f *FrProvider) GetBase64Image() *string {
	return &f.model.Base64Image
}

func (f *FrProvider) GetCreatedAt() string {
	return f.model.CreatedAt
}

func (f *FrProvider) GetSourceId() string {
	return f.model.SourceId
}

func (f *FrProvider) GetFileName() string {
	return f.model.CreateFileName()
}

func (f *FrProvider) SetImageFileName(jsonObj *models.FaceRecognitionJsonObject, fullImageName string) {
	jsonObj.FaceRecognition.ImageFileName = fullImageName
}

func (f *FrProvider) SetDataFileName(jsonObj *models.FaceRecognitionJsonObject, fullDataName string) {
	jsonObj.FaceRecognition.DataFileName = fullDataName
}

func (f *FrProvider) CreateJsonObject() *models.FaceRecognitionJsonObject {
	fr := f.model
	baseObj := models.FaceRecognitionJsonBaseObject{Id: fr.Id, SourceId: fr.SourceId, CreatedAt: fr.CreatedAt,
		DetectedFaces: fr.DetectedFaces, AiClipEnabled: fr.AiClipEnabled}
	jsonObj := &models.FaceRecognitionJsonObject{FaceRecognition: &baseObj, Video: &models.VideoClipJsonObject{}}
	return jsonObj
}
