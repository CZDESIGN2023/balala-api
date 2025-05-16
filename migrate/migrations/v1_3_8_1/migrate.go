package v1_3_8_1

import (
	_ "embed"
	"go-cs/migrate/pkg"
	"go-cs/pkg/sql_parser"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

const VERSION = "1.3.8.1"

var ProviderSet = wire.NewSet(
	NewMigrate,
)

type Migrate struct {
	logger  *log.Helper
	version pkg.Version

	db *gorm.DB
}

//go:embed ddl.sql
var ddl string

func NewMigrate(logger log.Logger) *Migrate {
	helper := log.NewHelper(logger, log.WithMessageKey("migrate/"+VERSION))
	return &Migrate{
		logger:  helper,
		version: pkg.NewVersion(VERSION),
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

	m.logger.Info("migrate end")

	return nil
}

func (m *Migrate) MigrateNoTrans(tx *gorm.DB) error {
	return nil
}

func (m *Migrate) Version() pkg.Version {
	return m.version
}
