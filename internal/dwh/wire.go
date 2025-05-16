//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package dwh

import (
	"go-cs/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(dwhConf *conf.Dwh, logger log.Logger) (*DwhApp, func(), error) {
	panic(wire.Build(
		DwhProviderSet,
		NewDwhApp,
	))
}
