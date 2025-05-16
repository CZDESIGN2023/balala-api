package job

import (
	"github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/dwh/task"
	"sync"
)

type OdsToDwdJob struct {
	id    string
	ctx   *pkg.JobContext
	tasks []pkg.JobTask

	initOnce sync.Once
}

func NewOdsToDwdJob(
	id string,
	ctx *pkg.JobContext,
) *OdsToDwdJob {
	return &OdsToDwdJob{
		id:  id,
		ctx: ctx,
	}
}

func (job *OdsToDwdJob) Name() string {
	return "ods_to_dwd_job"
}

func (job *OdsToDwdJob) Id() string {
	return job.id
}

func (job *OdsToDwdJob) FullName() string {
	return job.Name() + ":" + job.Id()
}

func (job *OdsToDwdJob) init() {
	job.initOnce.Do(func() {

		taskCtx := pkg.NewTaskContextWithJob(job, job.ctx)
		job.tasks = make([]pkg.JobTask, 0)
		job.tasks = append(job.tasks, task.NewOdsToDwdWitemTask("witemTask", taskCtx))
		job.tasks = append(job.tasks, task.NewOdsToDwdMemberTask("memberTask", taskCtx))
		job.tasks = append(job.tasks, task.NewOdsToDwdWitemFlowNodeTask("witemFlowNodeTask", taskCtx))
	})
}

func (job *OdsToDwdJob) Run() {

	job.init()

	//TODO 需要获取子任务状态，判断是否继续执行
	log.Debug("dwh job run ----> " + job.FullName())

	//TODO 需要引入工作流概念，控制任务流程
	//并行处理
	for _, v := range job.tasks {
		go v.Run()
	}

}

func (job *OdsToDwdJob) Stop() {

}
