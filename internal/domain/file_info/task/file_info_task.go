package task

import (
	"go-cs/internal/conf"
	"go-cs/internal/data"
	"go-cs/internal/utils"

	klog "github.com/go-kratos/kratos/v2/log"
)

type FileInfoTask struct {
	log  *klog.Helper
	conf *conf.FileConfig
	data *data.Data
}

func NewFileInfoTask(
	logger klog.Logger,
	conf *conf.FileConfig,
	data *data.Data,

) *FileInfoTask {

	moduleName := "FileInfoTask"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &FileInfoTask{
		conf: conf,
		log:  hlog,
		data: data,
	}
}
