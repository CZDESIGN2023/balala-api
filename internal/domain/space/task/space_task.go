package task

import (
	"go-cs/internal/data"
	"go-cs/internal/utils"

	space_repo "go-cs/internal/domain/space/repo"

	klog "github.com/go-kratos/kratos/v2/log"
)

type SpaceTask struct {
	log       *klog.Helper
	data      *data.Data
	spaceRepo space_repo.SpaceRepo
}

func NewSpaceTask(
	logger klog.Logger,
	data *data.Data,
	spaceRepo space_repo.SpaceRepo,

) *SpaceTask {

	moduleName := "SpaceTask"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceTask{
		data:      data,
		log:       hlog,
		spaceRepo: spaceRepo,
	}
}
