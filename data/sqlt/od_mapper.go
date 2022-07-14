package sqlt

import (
	"smcp/data"
	"smcp/models"
)

type OdMapper struct {
	Config *models.Config
}

func (o *OdMapper) Map(source *models.ObjectDetectionModel) []*OdEntity {
	ret := make([]*OdEntity, 0)
	for _, do := range source.DetectedObjects {
		entity := &OdEntity{}
		entity.GroupId = source.Id
		entity.SourceId = source.SourceId
		sio := &data.SaveImageOptions{Config: o.Config}
		sio.MapFromOd(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFileName = ""
		entity.SetupDates(source.CreatedAt)
		entity.PredScore = do.PredScore
		entity.PredClsIdx = do.PredClsIdx
		entity.PredClsName = do.PredClsName
		entity.AiClipEnabled = source.AiClipEnabled
		entity.AiClipFileName = ""
		entity.AiClipCreatedAtStr = ""
		entity.AiClipLastModifiedAtStr = ""
		entity.AiClipDuration = 0
		ret = append(ret, entity)
	}

	return ret
}
