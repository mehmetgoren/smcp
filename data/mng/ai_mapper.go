package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/models"
)

type AiMapper struct {
	Config *models.Config
}

func (o *AiMapper) Map(source *models.AiDetectionModel) []interface{} {
	ret := make([]interface{}, 0)
	for _, do := range source.Detections {
		entity := &AiEntity{Id: primitive.NewObjectID(), Module: source.Module, GroupId: source.Id, SourceId: source.SourceId}
		sio := &data.SaveImageOptions{Config: o.Config}
		sio.MapFromAi(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFile = &VideoFile{}
		entity.SetupDates(source.CreatedAt)
		entity.DetectedObject = &AiDetectedObject{
			Label: do.Label,
			Score: do.Score,
			X1:    do.Box.X1,
			Y1:    do.Box.Y1,
			X2:    do.Box.X2,
			Y2:    do.Box.Y2,
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
