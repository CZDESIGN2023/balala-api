package test

import (
	"fmt"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/dwh/task"
	"testing"

	"github.com/spf13/cast"
	"gorm.io/gorm/logger"
)

func mockTestJob() pkg.Job {
	return MockTestJob("1", "test_job")
}

func TestOdsToDimObjectTask(t *testing.T) {

	ctx := NewMockDwhCtx()

	job := mockTestJob()
	task.NewOdsToDimObjectTask("dim_object", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDimUserTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDimUserTask("dim_user", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDimVersionTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDimVersionTask("dim_version", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDimWitemStatusTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDimWitemStatusTask("dim_status", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDimSpaceTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDimSpaceTask("dim_space", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDwdWitemTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDwdWitemTask("dwd_witem", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDwdMemberTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDwdMemberTask("dwd_member", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestOdsToDwdWitemFlowNodeTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewOdsToDwdWitemFlowNodeTask("dwd_witem_flow_node", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestDwdToDwsMbrWitem1hTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewDwdToDwsMemberWitem1hTask("dwsMbrWitem1h", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestDwdToDwsVersionWitem1hTask(t *testing.T) {
	ctx := NewMockDwhCtx()
	job := mockTestJob()
	task.NewDwdToDwsVersionWitem1hTask("dwsVerWitem1h", pkg.NewTaskContextWithJob(job, &pkg.JobContext{
		Data:             ctx.dwhData,
		JobVariablesRepo: ctx.JobVariablesRepo,
	})).Run()
}

func TestTemporaryTable(t *testing.T) {
	ctx := NewMockDwhCtx()

	db := ctx.dwhData.Db()
	db.Logger = db.Logger.LogMode(logger.Info)
	db.Debug()
	for i := 0; i < 10000; i++ {
		tbName := fmt.Sprintf("temp_test_%d", i)

		err0 := db.Exec("drop table if exists " + tbName).Error
		if err0 != nil {
			fmt.Println(err0)
		}

		err1 := db.Exec("CREATE TEMPORARY TABLE " + tbName + " AS SELECT " + cast.ToString(i) + " AS id").Error
		if err1 != nil {
			fmt.Println(err1)
		}

		res := make([]map[string]interface{}, 0)
		err2 := db.Raw("SELECT * FROM " + tbName).Find(&res).Error
		if err2 != nil {
			fmt.Println(err2)
		}

		fmt.Println(res)
	}
}
