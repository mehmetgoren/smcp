package reps

import (
	"smcp/disk"
	"smcp/models"
)

type OdHandlerRepository struct {
	Config *models.Config
}

func (o *OdHandlerRepository) Save(od *models.ObjectDetectionModel) error {
	p := &disk.OdProvider{}
	s := disk.Storage[models.ObjectDetectionJsonObject]{}
	s.Provider = p
	return s.Save(o.Config, od)
}
