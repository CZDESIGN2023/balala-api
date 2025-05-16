package pkg

import "go-cs/internal/dwh/data"

type Job interface {
	Id() string
	Name() string
	FullName() string
	Run()
}

type JobTask interface {
	Id() string
	Name() string
	FullName() string
	Run()
}

type JobContext struct {
	Data             *data.DwhData
	JobVariablesRepo *data.JobVariablesRepo
}

type TaskContext struct {
	Job              Job
	Data             *data.DwhData
	JobVariablesRepo *data.JobVariablesRepo
}

func NewTaskContextWithJob(job Job, jobCtx *JobContext) *TaskContext {
	return &TaskContext{
		Job:              job,
		JobVariablesRepo: jobCtx.JobVariablesRepo,
		Data:             jobCtx.Data,
	}
}
