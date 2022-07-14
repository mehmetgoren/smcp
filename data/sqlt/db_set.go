package sqlt

import (
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
)

type DbSet[T any] struct {
	db       *gorm.DB
	template *T
}

func (d *DbSet[T]) Migrate() error {
	entity := new(T)
	d.template = entity
	migrator := d.db.Migrator()
	err := migrator.AutoMigrate(entity)
	if err == nil {
		var tn = strings.ToLower(strings.Replace(reflect.TypeOf(entity).String(), "*sqlt.", "", -1))
		idxNames := []string{"idx_query"}
		for _, idxName := range idxNames {
			err = migrator.RenameIndex(entity, idxName, idxName+"_"+tn)
			if err != nil {
				log.Println(err.Error())
			}
			err = migrator.DropIndex(entity, idxName)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	return err
}

func (d *DbSet[T]) AddRange(models []*T) error {
	result := d.db.Create(models)
	err := result.Error
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func (d *DbSet[T]) GetGormDb() *gorm.DB {
	return d.db
}
