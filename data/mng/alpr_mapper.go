package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/models"
)

type AlprMapper struct {
	Config *models.Config
}

func (a *AlprMapper) Map(source *models.AlprResponse) []interface{} {
	ret := make([]interface{}, 0)
	for _, r := range source.Results {
		entity := &AlprEntity{Id: primitive.NewObjectID(), GroupId: source.Id, SourceId: source.SourceId}
		sio := &data.SaveImageOptions{Config: a.Config}
		sio.MapFromAlpr(source)
		imageFileName, _ := sio.SaveImage()
		entity.ImageFileName = imageFileName
		entity.SetupDates(source.CreatedAt)
		entity.ImgWidth = source.ImgWidth
		entity.ImgHeight = source.ImgHeight
		entity.TotalProcessingTimeMs = source.ProcessingTimeMs
		entity.VideoFile = &VideoFile{}
		entity.DetectedPlate = &DetectedPlate{
			Plate:            r.Plate,
			Confidence:       r.Confidence,
			ProcessingTimeMs: r.ProcessingTimeMs,
			Coordinates: &Coordinates{
				X0: r.Coordinates.X0,
				Y0: r.Coordinates.Y0,
				X1: r.Coordinates.X1,
				Y1: r.Coordinates.Y1,
			},
		}
		entity.DetectedPlate.Candidates = make([]*Candidate, 0)
		for _, c := range r.Candidates {
			entityCandidate := &Candidate{Plate: c.Plate, Confidence: c.Confidence}
			entity.DetectedPlate.Candidates = append(entity.DetectedPlate.Candidates, entityCandidate)
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
