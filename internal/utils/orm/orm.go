package orm

import (
	"errors"
	"go-cs/internal/conf"
	"go-cs/internal/utils/orm/mysql"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	zapgorm "moul.io/zapgorm2"
)

var (
	ErrUnsupportedType         = errors.New("unsupported database type")
	ErrUnsupportedResolverType = errors.New("unsupported resolver type")
)

// Driver database driver type
type Driver string

func (d Driver) String() string {
	return string(d)
}

const (
	MySQL       Driver = "mysql"
	PostgresSQL Driver = "postgres"
)

// LogLevel database logger level
type LogLevel string

const (
	Silent LogLevel = "silent"
	Error  LogLevel = "error"
	Warn   LogLevel = "warn"
	Info   LogLevel = "info"
)

// Convert convert to gorm logger level
func (l LogLevel) Convert() logger.LogLevel {
	switch l {
	case Silent:
		return logger.Silent
	case Error:
		return logger.Error
	case Warn:
		return logger.Warn
	case Info:
		return logger.Info
	default:
		return logger.Silent
	}
}

type Config struct {
	DSN             string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
	LogLevel        LogLevel
	Plugins         func(db *gorm.DB) ([]gorm.Plugin, error)
	Resolvers       []Resolver
	Debug           bool
}

func NewConfig(data *conf.Data) *Config {
	return &Config{
		DSN:             data.Database.Dsn,
		Debug:           data.Database.Debug,
		MaxOpenConn:     80,
		MaxIdleConn:     10,
		ConnMaxLifeTime: 4 * time.Hour,
		ConnMaxIdleTime: 3 * time.Minute,
	}
}

// New initialize orm
func NewOrm(c *conf.Data, kLogger klog.Logger, zLogger *zap.Logger) (db *gorm.DB, cleanup func(), err error) {
	if c == nil {
		return nil, func() {}, nil
	}

	config := NewConfig(c)
	gLogger := zapgorm.New(zLogger.WithOptions(zap.AddCallerSkip(3)))
	gLogger.SetAsDefault()

	db, err = mysql.New(mysql.Config{
		Driver:                    MySQL.String(),
		DSN:                       config.DSN,
		MaxIdleConn:               config.MaxIdleConn,
		MaxOpenConn:               config.MaxOpenConn,
		ConnMaxIdleTime:           config.ConnMaxIdleTime,
		ConnMaxLifeTime:           config.ConnMaxLifeTime,
		Logger:                    gLogger.LogMode(config.LogLevel.Convert()),
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         0,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
	})

	if config.Debug {
		db = db.Debug()
	}

	if err != nil {
		return nil, nil, err

	}

	if len(config.Resolvers) > 0 {
		if err = registerResolver(db, MySQL, config.Resolvers); err != nil {
			return nil, nil, err
		}
	}

	cleanup = func() {
		klog.NewHelper(kLogger).Info("closing the database resources")

		sqlDB, err := db.DB()
		if err != nil {
			klog.NewHelper(kLogger).Error(err)
		}

		if err := sqlDB.Close(); err != nil {
			klog.NewHelper(kLogger).Error(err)
		}
	}

	return db, cleanup, nil
}

func registerResolver(db *gorm.DB, driver Driver, resolvers []Resolver) error {
	if len(resolvers) > 0 {
		var (
			sources  = make([]gorm.Dialector, 0, len(resolvers))
			replicas = make([]gorm.Dialector, 0, len(resolvers))
		)

		for _, resolver := range resolvers {
			dial, err := BuildDialector(driver, resolver.DSN)
			if err != nil {
				return err
			}
			switch resolver.Type {
			case Source:
				sources = append(sources, dial)
			case Replica:
				replicas = append(replicas, dial)
			default:
				return ErrUnsupportedResolverType
			}
		}

		return db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  sources,
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	return nil
}
