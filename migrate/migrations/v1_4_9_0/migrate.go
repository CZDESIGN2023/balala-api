package v1_4_9_0

import (
	_ "embed"
	"go-cs/migrate/pkg"
	"go-cs/pkg/sql_parser"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

const VERSION = "1.4.9.0"

var ProviderSet = wire.NewSet(
	NewMigrate,
	NewDML,
)

type Migrate struct {
	logger  *log.Helper
	version pkg.Version

	db *gorm.DB

	dml *DML
}

//go:embed ddl.sql
var ddl string

func NewMigrate(logger log.Logger, dml *DML) *Migrate {
	helper := log.NewHelper(logger, log.WithMessageKey("migrate/"+VERSION))
	return &Migrate{
		logger:  helper,
		version: pkg.NewVersion(VERSION),
		dml:     dml,
	}
}

func (m *Migrate) Migrate(tx *gorm.DB) error {
	m.logger.Info("migrate start")

	statements := sql_parser.ToStatements(ddl)

	for _, statement := range statements {
		if err := tx.Exec(statement).Error; err != nil {
			return err
		}
	}

	err := m.dml.HandleData()
	if err != nil {
		m.logger.Error(err)
	}

	m.logger.Info("migrate end")

	return nil
}

func (m *Migrate) MigrateNoTrans(tx *gorm.DB) error {
	return nil
}

func (m *Migrate) Version() pkg.Version {
	return m.version
}
