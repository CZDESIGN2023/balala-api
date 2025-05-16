//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"go-cs/internal/biz"
	"go-cs/internal/biz/event_handlers"
	"go-cs/internal/conf"
	"go-cs/internal/data"
	"go-cs/internal/domain"
	"go-cs/internal/dwh"
	dwh_service "go-cs/internal/dwh/service"
	"go-cs/internal/pkg"
	"go-cs/internal/server"
	"go-cs/internal/service"
	"go-cs/internal/utils/third_platform"
	"go-cs/migrate"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, *conf.Server, *conf.Data, *conf.Jwt, *conf.FileConfig, *conf.Dwh, *dwh.DwhApp, *dwh_service.DwhService, *third_platform.Client, log.Logger, *zap.Logger) (*server.App, func(), error) {

	panic(wire.Build(
		migrate.ProviderSet,
		server.ProviderSet,
		data.ProviderSet,
		domain.ProviderSet,
		service.ProviderSet,
		pkg.ProviderSet,
		event_handlers.ProviderSet,
		biz.ProviderSet,
		newApp,
		server.NewApp,
	))
}
