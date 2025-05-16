package job

import (
	"context"
	"go-cs/internal/utils"
	"log"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"

	file_task "go-cs/internal/domain/file_info/task"
	space_task "go-cs/internal/domain/space/task"
	witem_task "go-cs/internal/domain/work_item/task"
	temp_file_task "go-cs/internal/server/file/task"
)

type Cron struct {
	log    *klog.Helper
	server *cron.Cron

	spaceTask    *space_task.SpaceTask
	witemTask    *witem_task.WorkItemTask
	fileTask     *file_task.FileInfoTask
	tempFileTask *temp_file_task.Task
}

func NewCron(
	logger klog.Logger,
	spaceTask *space_task.SpaceTask,
	witemTask *witem_task.WorkItemTask,
	fileTask *file_task.FileInfoTask,
	temp_file_task *temp_file_task.Task,
) (*Cron, error) {

	moduleName := "Cron"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	server := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.PrintfLogger(log.Default())),
			cron.DelayIfStillRunning(cron.PrintfLogger(log.Default())),
		),
	)

	return &Cron{
		log:          hlog,
		server:       server,
		spaceTask:    spaceTask,
		witemTask:    witemTask,
		fileTask:     fileTask,
		tempFileTask: temp_file_task,
	}, nil
}

func (c *Cron) Start() (err error) {
	// 每天0点10秒执行
	if _, err = c.server.AddFunc("10 0 0 * * *", func() {
		c.witemTask.CheckWorkItemExpired()
		c.witemTask.CheckFlowNodeExpired()
	}); err != nil {
		return err
	}

	// 每3天0点1分执行
	if _, err = c.server.AddFunc("0 1 0 */3 * *", func() {
		c.spaceTask.CheckSpaceAbnormal()
	}); err != nil {
		return err
	}

	// 每周六2点0分执行 已删除文件清理
	if _, err = c.server.AddFunc("0 0 2 * * 6", func() {
		c.fileTask.CleanDeletedFile()
	}); err != nil {
		return err
	}

	// 每天3点0分执行 临时文件清理
	if _, err = c.server.AddFunc("0 0 3 * * *", func() {
		c.tempFileTask.CleanTempFiles()
	}); err != nil {
		return err
	}

	// 启动 cron 服务
	c.server.Start()

	c.log.Info("cron server started")
	return nil
}

// Stop cron 服务关闭
func (c *Cron) Stop(ctx context.Context) (err error) {
	c.server.Stop()

	c.log.Info("cron server has been stop")
	return nil
}
