package data

import (
	"go-cs/internal/data/test"
	"gorm.io/gorm"
)

func init() {
	test.Init("../../../configs/config.yaml", "../../../.env.local.test", true)

	data := &DwhData{
		db:   test.GetDB(),
		conf: test.GetConf().Dwh,
	}
	gdb = data.db

	repo = NewRptRepo(data)
}

var gdb *gorm.DB

var (
	repo *RptRepo
)
