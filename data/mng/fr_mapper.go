package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/models"
)

type FrMapper struct {
	Config *models.Config
}

func (f *FrMapper) Map(source *models.FaceRecognitionModel) []interface{} {
	ret := make([]interface{}, 0)
	for _, do := range source.DetectedFaces {
		entity := &FrEntity{Id: primitive.NewObjectID(), GroupId: source.Id, SourceId: source.SourceId}
		sio := &data.SaveImageOptions{Config: f.Config}
		sio.MapFromFr(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFile = &VideoFile{}
		entity.SetupDates(source.CreatedAt)
		entity.DetectedFace = &DetectedFace{
			PredClsName: do.PredClsName,
			PredClsIdx:  do.PredClsIdx,
			PredScore:   do.PredScore,
			X1:          do.X1,
			Y1:          do.Y1,
			X2:          do.X2,
			Y2:          do.Y2,
		}
		entity.AiClip = &data.AiClip{
			Enabled:        source.AiClipEnabled,
			FileName:       "",
			CreatedAt:      "",
			LastModifiedAt: "",
			Duration:       0,
		}
		ret = append(ret, entity)
	}

	return ret
}
