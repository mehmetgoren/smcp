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
			X1:          do.Box.X1,
			Y1:          do.Box.Y1,
			X2:          do.Box.X2,
			Y2:          do.Box.Y2,
		}
		entity.AiClip = &data.AiClip{
			Enabled:        source.AiClipEnabled,
			FileName:       "",
			CreatedAt:      "",
			LastModifiedAt: "",
			Duration:       0,
		}
		if do.Metadata != nil {
			entity.DetectedObject.Metadata = &Metadata{}
			if do.Metadata.Colors != nil {
				entity.DetectedObject.Metadata.Colors = make([]Color, 0)
				for _, color := range do.Metadata.Colors {
					entity.DetectedObject.Metadata.Colors = append(entity.DetectedObject.Metadata.Colors, Color{
						R: color.R,
						G: color.G,
						B: color.B,
					})
				}
			}
		}
		ret = append(ret, entity)
	}

	return ret
}
