package cmn

import (
	"smcp/data"
	"smcp/data/mng"
	"smcp/data/sqlt"
	"smcp/models"
	"strconv"
)

type Factory struct {
	Config *models.Config
	rep    data.Repository

	sdb *sqlt.DbContext
	mdb *mng.DbContext
}

func (f *Factory) Init() error {
	var err error
	switch f.Config.Db.Type {
	case 0:
		f.sdb = &sqlt.DbContext{Config: f.Config}
		err = f.sdb.Init()
		f.rep = &sqlt.SqliteRepository{Db: f.sdb}
		break
	case 1:
		f.mdb = &mng.DbContext{Config: f.Config}
		err = f.mdb.Init()
		f.rep = &mng.MongoRepository{Db: f.mdb}
		break
	default:
		panic("not supported: " + strconv.Itoa(f.Config.Db.Type))
	}
	return err
}

func (f *Factory) CreateRepository() data.Repository {
	return f.rep
}

func (f *Factory) Close() error {
	if f.mdb != nil {
		return f.mdb.Close()
	}

	return nil
}
