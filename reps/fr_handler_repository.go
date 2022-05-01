package reps

import (
	"smcp/disk"
	"smcp/models"
)

type FrHandlerRepository struct {
	Config *models.Config
}

func (f *FrHandlerRepository) Save(fr *models.FaceRecognitionModel) error {
	p := &disk.FrProvider{}
	s := disk.Storage[models.FaceRecognitionJsonObject]{}
	s.Provider = p
	return s.Save(f.Config, fr)
}
