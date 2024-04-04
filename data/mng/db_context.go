package mng

import (
	"log"
	"smcp/models"
)

type DbContext struct {
	Config *models.Config

	Ais *DbSet[AiEntity]
}

func (d *DbContext) Init() error {
	cs := d.Config.Db.ConnectionString
	d.Ais = &DbSet[AiEntity]{Scheme: AiEntityScheme{}, ConnectionString: cs}
	err := d.Ais.Open()
	if err != nil {
		return err
	}
	_, err = d.Ais.CreateIndexes()
	if err != nil {
		return err
	}

	return err
}

func (d *DbContext) Close() error {
	err := d.Ais.Close()
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
