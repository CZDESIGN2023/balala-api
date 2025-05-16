package pkg

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	"go-cs/pkg/stream"
	"gorm.io/gorm"
	"slices"
	"time"
)

type Migration interface {
	Migrate(tx *gorm.DB) error
	MigrateNoTrans(tx *gorm.DB) error
	Version() Version
}

type Manager struct {
	logger *log.Helper
	db     *gorm.DB

	migrations []Migration

	srcVersion Version
	dstVersion Version
}

func NewManager(logger log.Logger, db *gorm.DB, migrations []Migration, srcVersion, dstVersion Version) *Manager {
	if len(migrations) == 0 {
		panic("no migrations provided")
	}

	helper := log.NewHelper(logger, log.WithMessageKey("Migrate Manager"))

	// 版本排序
	slices.SortFunc(migrations, func(a, b Migration) int {
		return a.Version().CompareTo(b.Version())
	})

	return &Manager{
		logger:     helper,
		db:         db,
		migrations: migrations,
		srcVersion: srcVersion,
		dstVersion: dstVersion,
	}
}

func (m *Manager) Migrate() error {
	if len(m.migrations) == 0 {
		m.logger.Info("No migrations to apply")
		return nil
	}

	targetMigrations := stream.Filter(m.migrations, func(mg Migration) bool {
		return mg.Version().CompareTo(m.srcVersion) > 0 && mg.Version().CompareTo(m.dstVersion) <= 0
	})

	if len(targetMigrations) == 0 {
		m.logger.Info("No target migrations found")
		return nil
	}

	for _, migration := range targetMigrations {
		m.logger.Infof("Migrate %s start", migration.Version())

		// 迁移无事务
		err := migration.MigrateNoTrans(m.db)
		if err != nil {
			m.logger.Errorf("MigrateNoTrans %s error: %s", migration.Version(), err.Error())
			return err
		}

		// 迁移有事务
		err = m.db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Migrate(tx); err != nil {
				return err
			}

			// 更新系统版本
			if err := m.UpdateSystemVersion(tx, migration.Version()); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			m.logger.Errorf("Migrate %s error: %s", migration.Version(), err.Error())
			return err
		}

		m.logger.Infof("Migrate %s success", migration.Version())
	}

	// 更新dst系统版本
	if err := m.UpdateSystemVersion(m.db, m.GetDstVersion()); err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateSystemVersion(tx *gorm.DB, version Version) error {
	config := &db.Config{}
	err := tx.Where("config_key = ?", consts.CONFIG_BALALA_VERSION).Take(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = tx.Save(&db.Config{
		Id:           config.Id,
		ConfigName:   consts.CONFIG_BALALA_VERSION,
		ConfigKey:    consts.CONFIG_BALALA_VERSION,
		ConfigValue:  version.String(),
		ConfigStatus: "1",
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}).Error

	return err
}

func (m *Manager) GetSrcVersion() Version {
	return m.srcVersion
}

func (m *Manager) GetDstVersion() Version {
	return m.dstVersion
}

func (m *Manager) GetMigrations() []Migration {
	return m.migrations
}
