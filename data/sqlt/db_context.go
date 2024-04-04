package sqlt

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"path"
	"smcp/models"
)

type DbContext struct {
	Config *models.Config

	Ais *DbSet[AiEntity]
}

func (d *DbContext) Init() error {
	p := path.Join(d.Config.Db.ConnectionString, "feniks.db")
	db, _ := gorm.Open(sqlite.Open(p), &gorm.Config{})
	d.Ais = &DbSet[AiEntity]{db: db}

	err := d.Ais.Migrate()

	return err
}
