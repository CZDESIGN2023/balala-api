package dwh

import (
	"go-cs/internal/conf"
	"go-cs/internal/dwh/data"
	"go-cs/internal/dwh/job"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/dwh/service"
	"go-cs/internal/utils"

	s_log "log"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var DwhProviderSet = wire.NewSet(
	data.ProviderSet,
	service.ProviderSet,
)

type DwhApp struct {
	conf *conf.Dwh
	log  *log.Helper
	data *data.DwhData

	jobVariablesRepo *data.JobVariablesRepo
	rptRepo          *data.RptRepo

	service *service.DwhService
}

func NewDwhApp(
	logger log.Logger,
	conf *conf.Dwh,
	data *data.DwhData,

	jobVariablesRepo *data.JobVariablesRepo,
	rptRepo *data.RptRepo,

	service *service.DwhService,
) *DwhApp {

	moduleName := "Dwh"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	app := &DwhApp{
		conf:             conf,
		log:              hlog,
		data:             data,
		jobVariablesRepo: jobVariablesRepo,
		rptRepo:          rptRepo,
		service:          service,
	}

	return app
}

func (d *DwhApp) Service() *service.DwhService {
	return d.service
}

func (d *DwhApp) Stop() {
}

func (d *DwhApp) StartODS() error {

	//没有配置不启动
	if d.conf == nil || d.conf.Database == nil || d.conf.Database.Dsn == "" {
		return nil
	}

	jobCtx := &pkg.JobContext{
		Data:             d.data,
		JobVariablesRepo: d.jobVariablesRepo,
	}

	//一次型任务
	//日期纬度生成

	//服务型任务,常驻, 实时
	go job.NewSyncMysqlBinlogToODSJob("ods", jobCtx).Run()
	return nil
}

func (d *DwhApp) StartCron() error {
	//没有配置不启动
	if d.conf == nil || d.conf.Database == nil || d.conf.Database.Dsn == "" {
		return nil
	}

	jobCtx := &pkg.JobContext{
		Data:             d.data,
		JobVariablesRepo: d.jobVariablesRepo,
	}

	crond := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.PrintfLogger(s_log.Default())),
			cron.DelayIfStillRunning(cron.PrintfLogger(s_log.Default())),
		),
	)

	//周期型任务
	//从ODS逐行扫描生成 每10秒扫描一次 任务本身有状态记录 所以用担心重复执行问题
	crond.AddJob("*/10 * * * * *", job.NewOdsToDimJob("dim", jobCtx))
	crond.AddJob("*/10 * * * * *", job.NewOdsToDwdJob("dwd", jobCtx))
	//定时的汇总表生成 1小时增量
	crond.AddJob("0 10 */1 * * *", job.NewDwdToDws1hJob("dws1h", jobCtx))

	crond.Start()

	return nil
}
