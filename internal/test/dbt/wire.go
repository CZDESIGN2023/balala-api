//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package dbt

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go-cs/internal/biz"
	"go-cs/internal/conf"
	"go-cs/internal/data"
	"go-cs/internal/domain"
	"go-cs/internal/dwh"
	"go-cs/internal/pkg"
	"go-cs/internal/service"
	"go-cs/internal/utils/third_platform"

	"go.uber.org/zap"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, *conf.Data, *conf.Jwt, *conf.FileConfig, *conf.Dwh, log.Logger, *zap.Logger, *third_platform.Client) (*All, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		domain.ProviderSet,
		service.ProviderSet,
		pkg.ProviderSet,
		dwh.DwhProviderSet,
		NewRepo,
		NewUsecase,
		NewService,
		NewDomainService,
		NewAll,
	))
}
