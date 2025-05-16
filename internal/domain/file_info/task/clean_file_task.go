package task

import (
	"context"
	db "go-cs/internal/bean/biz"
	"os"
	"path/filepath"
	"time"
)

// CleanDeletedFile 清理已删除文件
func (task *FileInfoTask) CleanDeletedFile() {
	begin := time.Now()
	task.log.Info("clean file")
	defer task.log.Infof("clean file end %v", time.Now().Sub(begin))

	list := task.getDeletedFiles()

	basePath := task.conf.LocalPath
	var errCount int

	var successIds []int64
	for _, v := range list {
		path := filepath.Join(basePath, v.Uri)
		task.log.Infof("clean file: %v", path)
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			errCount++
		} else {
			successIds = append(successIds, v.Id)
		}
	}

	if len(successIds) != 0 {
		task.data.DB(context.Background()).Where("id in ?", successIds).Delete(&db.FileInfo{})
		task.data.DB(context.Background()).Where("file_info_id in ?", successIds).Delete(&db.SpaceFileInfo{})
	}

	task.log.Infof("clean file, success: %v, failed: %v", len(list)-errCount, errCount)
}

func (task *FileInfoTask) getDeletedFiles() []*db.FileInfo {
	const sql = `
SELECT
	* 
FROM
	file_info 
WHERE
	id IN (
	SELECT
		file_info_id
	FROM
		space_file_info 
	WHERE
		status = 2 AND deleted_at > 0
	);
`
	var list []*db.FileInfo

	err := task.data.RoDB(context.Background()).Model(&db.FileInfo{}).
		Raw(sql).
		Find(&list).Error
	if err != nil {
		task.log.Error(err)
		return nil
	}

	return list
}
