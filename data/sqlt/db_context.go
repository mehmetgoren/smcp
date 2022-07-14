package sqlt

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"path"
	"smcp/models"
)

type DbContext struct {
	Config *models.Config

	Ods   *DbSet[OdEntity]
	Frs   *DbSet[FrEntity]
	Alprs *DbSet[AlprEntity]
}

func (d *DbContext) Init() error {
	p := path.Join(d.Config.Db.ConnectionString, "feniks.db")
	db, _ := gorm.Open(sqlite.Open(p), &gorm.Config{})
	d.Ods = &DbSet[OdEntity]{db: db}
	d.Frs = &DbSet[FrEntity]{db: db}
	d.Alprs = &DbSet[AlprEntity]{db: db}

	err := d.Ods.Migrate()
	if err != nil {
		return err
	}
	err = d.Frs.Migrate()
	if err != nil {
		return err
	}
	err = d.Alprs.Migrate()

	return err
}
