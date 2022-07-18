package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/models"
)

type OdMapper struct {
	Config *models.Config
}

func (o *OdMapper) Map(source *models.ObjectDetectionModel) []interface{} {
	ret := make([]interface{}, 0)
	for _, do := range source.DetectedObjects {
		entity := &OdEntity{Id: primitive.NewObjectID(), GroupId: source.Id, SourceId: source.SourceId}
		sio := &data.SaveImageOptions{Config: o.Config}
		sio.MapFromOd(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFile = &VideoFile{}
		entity.SetupDates(source.CreatedAt)
		entity.DetectedObject = &DetectedObject{
			PredClsName: do.PredClsName,
			PredClsIdx:  do.PredClsIdx,
			PredScore:   do.PredScore,
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
