package job

import (
	"github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/dwh/task"
	"sync"
)

type OdsToDimJob struct {
	id    string
	ctx   *pkg.JobContext
	tasks []pkg.JobTask

	initOnce sync.Once
}

func NewOdsToDimJob(
	id string,
	ctx *pkg.JobContext,
) *OdsToDimJob {
	return &OdsToDimJob{
		id:  id,
		ctx: ctx,
	}
}

func (job *OdsToDimJob) Name() string {
	return "ods_to_dim_job"
}

func (job *OdsToDimJob) Id() string {
	return job.id
}

func (job *OdsToDimJob) FullName() string {
	return job.Name() + ":" + job.Id()
}

func (job *OdsToDimJob) init() {
	job.initOnce.Do(func() {

		taskCtx := pkg.NewTaskContextWithJob(job, job.ctx)

		job.tasks = make([]pkg.JobTask, 0)
		job.tasks = append(job.tasks, task.NewOdsToDimSpaceTask("space", taskCtx))
		job.tasks = append(job.tasks, task.NewOdsToDimObjectTask("object", taskCtx))
		job.tasks = append(job.tasks, task.NewOdsToDimUserTask("user", taskCtx))
		job.tasks = append(job.tasks, task.NewOdsToDimVersionTask("version", taskCtx))
		job.tasks = append(job.tasks, task.NewOdsToDimWitemStatusTask("witemStatus", taskCtx))

	})
}

func (job *OdsToDimJob) Run() {

	job.init()

	//TODO 需要获取子任务状态，判断是否继续执行
	log.Debug("dwh job run ----> ods_to_dim_job")

	//TODO 需要引入工作流概念，控制任务流程
	//并行处理
	for _, v := range job.tasks {
		go v.Run()
	}

}

func (job *OdsToDimJob) Stop() {

}
