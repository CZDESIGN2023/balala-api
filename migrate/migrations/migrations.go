package migrations

import (
	"go-cs/migrate/migrations/v1_3_6_1"
	"go-cs/migrate/migrations/v1_3_6_4"
	"go-cs/migrate/migrations/v1_3_7_0"
	"go-cs/migrate/migrations/v1_3_7_2"
	"go-cs/migrate/migrations/v1_3_8_0"
	"go-cs/migrate/migrations/v1_3_8_1"
	"go-cs/migrate/migrations/v1_3_8_12"
	"go-cs/migrate/migrations/v1_3_9_0"
	"go-cs/migrate/migrations/v1_3_9_1"
	"go-cs/migrate/migrations/v1_4_10_0"
	"go-cs/migrate/migrations/v1_4_11_0"
	"go-cs/migrate/migrations/v1_4_12_0"
	"go-cs/migrate/migrations/v1_4_13_0"
	"go-cs/migrate/migrations/v1_4_14_0"
	"go-cs/migrate/migrations/v1_4_14_1"

	"go-cs/migrate/migrations/v1_4_1_0"
	"go-cs/migrate/migrations/v1_4_6_0"
	"go-cs/migrate/migrations/v1_4_7_0"
	"go-cs/migrate/migrations/v1_4_7_1"
	"go-cs/migrate/migrations/v1_4_8_0"
	"go-cs/migrate/migrations/v1_4_9_0"
	"go-cs/migrate/pkg"

	"github.com/google/wire"
)

type Migrations []pkg.Migration

// RequiredMinVersion is the minimum version of the database that can be migrated.
var RequiredMinVersion = pkg.NewVersion("1.3.6.0")

func NewMigrations(
	v1_3_6_1_ *v1_3_6_1.Migrate,
	v1_3_6_4_ *v1_3_6_4.Migrate,
	v1_3_7_0_ *v1_3_7_0.Migrate,
	v1_3_7_2_ *v1_3_7_2.Migrate,
	v1_3_8_0_ *v1_3_8_0.Migrate,
	v1_3_8_1_ *v1_3_8_1.Migrate,
	v1_3_8_12_ *v1_3_8_12.Migrate,
	v1_3_9_0_ *v1_3_9_0.Migrate,
	v1_3_9_1_ *v1_3_9_1.Migrate,
	v1_4_1_0_ *v1_4_1_0.Migrate,
	v1_4_6_0_ *v1_4_6_0.Migrate,
	v1_4_7_0_ *v1_4_7_0.Migrate,
	v1_4_7_1_ *v1_4_7_1.Migrate,
	v1_4_8_0_ *v1_4_8_0.Migrate,
	v1_4_9_0_ *v1_4_9_0.Migrate,
	v1_4_10_0_ *v1_4_10_0.Migrate,
	v1_4_11_0_ *v1_4_11_0.Migrate,
	v1_4_12_0_ *v1_4_12_0.Migrate,
	v1_4_13_0_ *v1_4_13_0.Migrate,
	v1_4_14_0_ *v1_4_14_0.Migrate,
	v1_4_14_1_ *v1_4_14_1.Migrate,

) Migrations {
	return []pkg.Migration{
		v1_3_6_1_,
		v1_3_6_4_,
		v1_3_7_0_,
		v1_3_7_2_,
		v1_3_8_0_,
		v1_3_8_1_,
		v1_3_8_12_,
		v1_3_9_0_,
		v1_3_9_1_,
		v1_4_1_0_,
		v1_4_6_0_,
		v1_4_7_0_,
		v1_4_7_1_,
		v1_4_8_0_,
		v1_4_9_0_,
		v1_4_10_0_,
		v1_4_11_0_,
		v1_4_12_0_,
		v1_4_13_0_,
		v1_4_14_0_,
		v1_4_14_1_,
	}
}

var ProviderSet = wire.NewSet(
	NewMigrations,
	v1_3_6_1.ProviderSet,
	v1_3_6_4.ProviderSet,
	v1_3_7_0.ProviderSet,
	v1_3_7_2.ProviderSet,
	v1_3_8_0.ProviderSet,
	v1_3_8_1.ProviderSet,
	v1_3_8_12.ProviderSet,
	v1_3_9_0.ProviderSet,
	v1_3_9_1.ProviderSet,
	v1_4_1_0.ProviderSet,
	v1_4_6_0.ProviderSet,
	v1_4_7_0.ProviderSet,
	v1_4_7_1.ProviderSet,
	v1_4_8_0.ProviderSet,
	v1_4_9_0.ProviderSet,
	v1_4_10_0.ProviderSet,
	v1_4_11_0.ProviderSet,
	v1_4_12_0.ProviderSet,
	v1_4_13_0.ProviderSet,
	v1_4_14_0.ProviderSet,
	v1_4_14_1.ProviderSet,
)
