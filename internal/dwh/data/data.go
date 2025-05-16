package data

import (
	"go-cs/internal/conf"
	"sync"

	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	NewDwhData,
	NewJobVariablesRepo,
	NewRptRepo,
)

type DwhData struct {
	conf *conf.Dwh

	dbInitOnce sync.Once
	db         *gorm.DB
}

func NewDwhData(conf *conf.Dwh) *DwhData {
	return &DwhData{
		conf: conf,
	}
}

func (d *DwhData) Db() *gorm.DB {

	d.dbInitOnce.Do(func() {
		db, err := gorm.Open(mysql.Open(d.conf.Database.Dsn), &gorm.Config{
			TranslateError: true,
		})
		if err != nil {
			panic(err)
		}
		d.db = db
	})

	return d.db
}

func (d *DwhData) DbConf() *conf.Dwh_Database {
	return d.conf.Database
}

func (d *DwhData) ExternalDataSourceConf() *conf.Dwh_ExternalDataSource {
	return d.conf.ExternalDataSource
}
