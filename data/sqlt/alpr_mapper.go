package sqlt

import (
	"encoding/json"
	"smcp/data"
	"smcp/models"
)

type AlprMapper struct {
	Config *models.Config
}

func (a *AlprMapper) Map(source *models.AlprResponse) []*AlprEntity {
	ret := make([]*AlprEntity, 0)
	for _, r := range source.Results {
		entity := &AlprEntity{}
		entity.GroupId = source.Id
		entity.SourceId = source.SourceId
		sio := &data.SaveImageOptions{Config: a.Config}
		sio.MapFromAlpr(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.VideoFileName = ""
		entity.SetupDates(source.CreatedAt)
		entity.ImgWidth = source.ImgWidth
		entity.ImgHeight = source.ImgHeight
		entity.TotalProcessingTimeMs = source.ProcessingTimeMs
		entity.Plate = r.Plate
		entity.Confidence = r.Confidence
		entity.ProcessingTimeMs = r.ProcessingTimeMs
		entity.X0 = r.Coordinates.X0
		entity.Y0 = r.Coordinates.Y0
		entity.X1 = r.Coordinates.X1
		entity.Y1 = r.Coordinates.Y1

		js, _ := json.Marshal(r.Candidates)
		if len(js) > 0 {
			entity.CandidatesJson = string(js)
		}

		entity.AiClipEnabled = source.AiClipEnabled
		entity.AiClipFileName = ""
		entity.AiClipCreatedAtStr = ""
		entity.AiClipLastModifiedAtStr = ""
		entity.AiClipDuration = 0
		ret = append(ret, entity)
	}

	return ret
}
