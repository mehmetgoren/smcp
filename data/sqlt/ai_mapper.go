package sqlt

import (
	"smcp/data"
	"smcp/models"
)

type AiMapper struct {
	Config *models.Config
}

func (o *AiMapper) Map(source *models.AiDetectionModel) []*AiEntity {
	ret := make([]*AiEntity, 0)
	for _, do := range source.Detections {
		entity := &AiEntity{}
		entity.Module = source.Module
		entity.GroupId = source.Id
		entity.SourceId = source.SourceId
		sio := &data.SaveImageOptions{Config: o.Config}
		sio.MapFromAi(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFileName = ""
		entity.SetupDates(source.CreatedAt)
		entity.Score = do.Score
		entity.Label = do.Label
		entity.X1 = do.Box.X1
		entity.Y1 = do.Box.Y1
		entity.X2 = do.Box.X2
		entity.Y2 = do.Box.Y2
		entity.AiClipEnabled = source.AiClipEnabled
		entity.AiClipFileName = ""
		entity.AiClipCreatedAtStr = ""
		entity.AiClipLastModifiedAtStr = ""
		entity.AiClipDuration = 0
		ret = append(ret, entity)
	}

	return ret
}
