package server

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/biz/event_handlers"
	"go-cs/internal/dwh"
	"go-cs/internal/server/job"
	"go-cs/internal/server/river"
	"go-cs/internal/utils"
	"go-cs/migrate"
	"gorm.io/gorm"
)

var GApp *App

type App struct {
	log          *log.Helper
	db           *gorm.DB
	server       *kratos.App
	cron         *job.Cron
	river        *river.River
	eventHandler *event_handlers.AppEventHandlers
	migrate      *migrate.Global
	dwh          *dwh.DwhApp
	version      string
}

func NewApp(logger log.Logger,
	db *gorm.DB,
	cron *job.Cron,
	river *river.River,
	ka *kratos.App,
	eventHandler *event_handlers.AppEventHandlers,
	global *migrate.Global,
	dwh *dwh.DwhApp,
) *App {

	moduleName := "APP"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	if GApp == nil {
		GApp = &App{
			log:          hlog,
			db:           db,
			server:       ka,
			cron:         cron,
			river:        river,
			eventHandler: eventHandler,
			migrate:      global,
			dwh:          dwh,
		}
	}

	return GApp
}

// Start 启动应用
func (a *App) Start(ctx context.Context) (stop chan error, err error) {
	//数据仓服务 在 river服务 前启动，避免漏接binlog
	a.river.Start()
	// 启动 数仓应用 ods服务
	a.dwh.StartODS()

	// 等待river服务 启动完成，否则会遇到表结构错误
	a.river.WaitUntilInitialized()

	//数据迁移
	if a.version != "" {
		if err = a.migrate.MigrateTo(a.version); err != nil {
			return
		}

		a.log.Info("database migration completed")
	}

	// 启动 数仓应用 定时服务
	a.dwh.StartCron()

	// 启动 cron 服务
	if err = a.cron.Start(); err != nil {
		return
	}

	stop = make(chan error, 1)

	// 启动 transport 服务
	go func() {
		a.log.Info("transport server starting ...")

		err := a.server.Run()

		if err != nil {
			stop <- err
			return
		}
	}()

	return stop, err
}

// Stop 停止应用
func (a *App) Stop(ctx context.Context) error {
	//关闭 cron 服务
	if err := a.cron.Stop(ctx); err != nil {
		a.log.Error(err)
	}

	// 关闭ipc服务
	// a.ipc.Close()
	a.dwh.Stop()
	a.river.Stop()

	if err := a.server.Stop(); err != nil {
		return err
	}

	a.log.Info("transport server stopping ...")

	return nil
}

func (a *App) SetVersion(ver string) {
	a.version = ver
}

func (a *App) Version() string {
	return a.version
}
