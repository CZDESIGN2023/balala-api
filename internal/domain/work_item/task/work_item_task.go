package task

import (
	"go-cs/internal/data"
	space_repo "go-cs/internal/domain/space/repo"
	"go-cs/internal/domain/work_flow/repo"
	"go-cs/internal/utils"

	klog "github.com/go-kratos/kratos/v2/log"
)

type WorkItemTask struct {
	log          *klog.Helper
	data         *data.Data
	workFlowRepo repo.WorkFlowRepo
}

func NewWorkItemTask(
	logger klog.Logger,
	data *data.Data,

	spaceRepo space_repo.SpaceRepo,
	workFlowRepo repo.WorkFlowRepo,

) *WorkItemTask {

	moduleName := "WorkItemTask"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemTask{
		data:         data,
		log:          hlog,
		workFlowRepo: workFlowRepo,
	}
}
