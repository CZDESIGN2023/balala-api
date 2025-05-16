package task

import (
	"github.com/go-kratos/kratos/v2/log"
	"go-cs/internal/conf"
	"go-cs/internal/utils"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Task struct {
	log  *log.Helper
	conf *conf.FileConfig
}

func New(fileCfg *conf.FileConfig, log log.Logger) *Task {
	_, helper := utils.InitModuleLogger(log, "temp_file_task")
	return &Task{
		log:  helper,
		conf: fileCfg,
	}
}

// CleanTempFiles 清理临时文件
func (t *Task) CleanTempFiles() {
	maxAge := time.Hour * 24 * 7 // 设置最大年龄为7天

	tempPath := path.Join(t.conf.LocalPath, ".temp")

	err := t.deleteFilesOlderThan(tempPath, maxAge)
	if err != nil {
		t.log.Errorf("Failed to delete .temp files: %v", err)
	}

	mergedPath := path.Join(t.conf.LocalPath, ".merged")
	err = t.deleteFilesOlderThan(mergedPath, maxAge)
	if err != nil {
		t.log.Errorf("Failed to delete .merged files: %v", err)
	}
}

// deleteFilesOlderThan 删除指定目录下创建时间超过maxAge的文件
func (t *Task) deleteFilesOlderThan(dir string, maxAge time.Duration) error {
	// 获取当前时间
	now := time.Now()
	// 遍历目录
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录
		if info.IsDir() {
			return nil
		}

		// 检查文件是否超过最大年龄
		if now.Sub(info.ModTime()) > maxAge {
			// 删除文件
			err := os.Remove(path)
			if err != nil {
				return err
			}
			t.log.Infof("Deleted file: %s", path)
		}
		return nil
	})
}
