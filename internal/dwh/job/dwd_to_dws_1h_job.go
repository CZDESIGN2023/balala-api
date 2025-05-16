package job

import (
	"fmt"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/dwh/task"
	"sync"
)

// 版本主题报表生成任务
type DwdToDws1hJob struct {
	id    string
	ctx   *pkg.JobContext
	tasks []pkg.JobTask

	initOnce sync.Once
}

func NewDwdToDws1hJob(
	id string,
	ctx *pkg.JobContext,
) *DwdToDws1hJob {
	return &DwdToDws1hJob{
		id:  id,
		ctx: ctx,
	}
}

func (job *DwdToDws1hJob) Name() string {
	return "dwd_to_dws_1h_job"
}

func (job *DwdToDws1hJob) Id() string {
	return job.id
}

func (job *DwdToDws1hJob) FullName() string {
	return job.Name() + ":" + job.Id()
}

func (job *DwdToDws1hJob) init() {
	job.initOnce.Do(func() {
		taskCtx := pkg.NewTaskContextWithJob(job, job.ctx)
		job.tasks = make([]pkg.JobTask, 0)
		job.tasks = append(job.tasks, task.NewDwdToDwsVersionWitem1hTask("dwsVerWitem1h", taskCtx))
		job.tasks = append(job.tasks, task.NewDwdToDwsMemberWitem1hTask("dwsMbrWitem1h", taskCtx))
	})
}

func (job *DwdToDws1hJob) Run() {

	job.init()

	//TODO 需要获取子任务状态，判断是否继续执行
	fmt.Println("rpt version job run ----> " + job.FullName())

	//TODO 需要引入工作流概念，控制任务流程
	//并行处理
	for _, v := range job.tasks {
		go v.Run()
	}

}

func (job *DwdToDws1hJob) Stop() {
	//TODO 需要获取子任务状态，判断是否继续执行
}
