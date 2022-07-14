package mng

import (
	"log"
	"smcp/models"
)

type DbContext struct {
	Config *models.Config

	Ods   *DbSet[OdEntity]
	Frs   *DbSet[FrEntity]
	Alprs *DbSet[AlprEntity]
}

func (d *DbContext) Init() error {
	cs := d.Config.Db.ConnectionString
	d.Ods = &DbSet[OdEntity]{Scheme: OdEntityScheme{}, ConnectionString: cs}
	err := d.Ods.Open()
	if err != nil {
		return err
	}
	_, err = d.Ods.CreateIndexes()
	if err != nil {
		return err
	}

	if err == nil {
		d.Frs = &DbSet[FrEntity]{Scheme: FrEntityScheme{}, ConnectionString: cs}
		err = d.Frs.Open()
		if err != nil {
			return err
		}
		_, err = d.Frs.CreateIndexes()
		if err != nil {
			return err
		}

		d.Alprs = &DbSet[AlprEntity]{Scheme: AlprEntityScheme{}, ConnectionString: cs}
		err = d.Alprs.Open()
		if err != nil {
			return err
		}
		_, err = d.Alprs.CreateIndexes()
	}

	return err
}

func (d *DbContext) Close() error {
	err := d.Ods.Close()
	if err != nil {
		log.Println(err.Error())
	}
	err = d.Frs.Close()
	if err != nil {
		log.Println(err.Error())
	}
	err = d.Alprs.Close()
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
