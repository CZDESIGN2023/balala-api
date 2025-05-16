package migrate

import (
	"github.com/google/wire"
	"go-cs/migrate/migrations"
)

var ProviderSet = wire.NewSet(
	NewGlobal,
	migrations.ProviderSet,
)
