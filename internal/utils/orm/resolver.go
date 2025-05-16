package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ResolverType resolver type
type ResolverType string

const (
	Source  ResolverType = "source"
	Replica ResolverType = "replica"
)

// Resolver provide multiple database support
type Resolver struct {
	Type ResolverType
	DSN  string
}

// BuildDialector build gorm.Dialector
func BuildDialector(driver Driver, dsn string) (dialector gorm.Dialector, err error) {
	switch driver {
	case MySQL:
		dialector = mysql.New(mysql.Config{
			DriverName: driver.String(),
			DSN:        dsn,
		})
	case PostgresSQL:
		dialector = postgres.New(postgres.Config{
			DriverName: driver.String(),
			DSN:        dsn,
		})
	default:
		return nil, ErrUnsupportedType
	}

	return dialector, nil
}
