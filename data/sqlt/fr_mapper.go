package sqlt

import (
	"smcp/data"
	"smcp/models"
)

type FrMapper struct {
	Config *models.Config
}

func (f *FrMapper) Map(source *models.FaceRecognitionModel) []*FrEntity {
	ret := make([]*FrEntity, 0)
	for _, do := range source.DetectedFaces {
		entity := &FrEntity{}
		entity.GroupId = source.Id
		entity.SourceId = source.SourceId
		sio := &data.SaveImageOptions{Config: f.Config}
		sio.MapFromFr(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFileName = ""
		entity.SetupDates(source.CreatedAt)
		entity.PredClsName = do.PredClsName
		entity.PredClsIdx = do.PredClsIdx
		entity.PredScore = do.PredScore
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
