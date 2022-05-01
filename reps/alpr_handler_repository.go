package reps

import (
	"smcp/disk"
	"smcp/models"
)

type AlprHandlerRepository struct {
	Config *models.Config
}

func (a *AlprHandlerRepository) Save(ar *models.AlprResponse) error {
	p := &disk.AlprProvider{}
	s := disk.Storage[models.AlprJsonObject]{}
	s.Provider = p
	return s.Save(a.Config, ar)
}
