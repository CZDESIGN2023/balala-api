package job

import (
	"github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/utils"
)

type Example struct {
	log *log.Helper
}

func NewExample(logger log.Logger) *Example {
	moduleName := "JobExample"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &Example{
		log: hlog,
	}
}

func (s *Example) Run() {
	s.log.Info("WorkItemExpired 任务执行成功")
}
