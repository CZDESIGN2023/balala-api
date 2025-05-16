package main

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/domain/work_item/task"
	"go-cs/internal/test/dbt"
	"testing"
)

func Test_Check(t *testing.T) {
	task := task.NewWorkItemTask(klog.DefaultLogger, dbt.Data, dbt.R.SpaceRepo, dbt.R.WorkFlowRepo)

	task.CheckFlowNodeExpired()
}
