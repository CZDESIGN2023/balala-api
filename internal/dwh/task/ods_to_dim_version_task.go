package task

import (
	"fmt"
	"go-cs/internal/dwh/data"
	dim_model "go-cs/internal/dwh/model/dim"
	ods_model "go-cs/internal/dwh/model/ods"
	"go-cs/internal/utils/date"

	"go-cs/internal/dwh/pkg"
	"go-cs/pkg/stream"

	"github.com/spf13/cast"
)

// converts data from ODS to DIM space
type OdsToDimVersionTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDimVersionTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDimVersionTask {
	return &OdsToDimVersionTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDimVersionTask) Id() string {
	return t.id
}

func (t *OdsToDimVersionTask) Name() string {
	return "ods_to_dim_version_task"
}

func (t *OdsToDimVersionTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDimVersionTask) Status() string {
	return t.status
}

func (t *OdsToDimVersionTask) Run() {

	if t.status == pkg.TASK_STATUS_RUNNING {
		return
	}

	defer func() {
		t.status = pkg.TASK_STATUS_READY
	}()

	t.status = pkg.TASK_STATUS_RUNNING

	//每次拿1000条
	//需要知道从哪里开始重新拿-》获取最后一次的id？或者最后一个的时间？ 按时间段来拿？
	lastIdVar, err := t.variablesRepo.GetVariablesByName(t.FullName(), "last_id")
	if err != nil {
		fmt.Println(err)
		return
	}

	var list []*ods_model.OdsVersion
	err = t.data.Db().Table("ods_version_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dimVersions []*dim_model.DimVersion
	err = t.data.Db().Table("dim_version").Where("version_id in ?", odsIds).Find(&dimVersions).Error
	if err != nil {
		return
	}

	dimVersionMap := stream.ToMap(dimVersions, func(_ int, t *dim_model.DimVersion) (int64, *dim_model.DimVersion) {
		return t.VersionId, t
	})

	for i := 0; i < len(list); i++ {

		odsObject := list[i]
		dimVersion := dimVersionMap[odsObject.Id]

		lastIdVar.VariableValue = cast.ToString(odsObject.OdsId)

		//不存在就新建
		if dimVersion == nil {
			dimVersion = &dim_model.DimVersion{}
			dimVersion.SpaceId = odsObject.SpaceId
			dimVersion.VersionId = odsObject.Id
			dimVersion.VersionName = odsObject.VersionName
			dimVersion.GmtCreate = cast.ToTime(odsObject.CreatedAt)
			dimVersion.GmtModified = cast.ToTime(odsObject.UpdatedAt)
			dimVersion.StartDate = cast.ToTime(odsObject.CreatedAt)
			dimVersion.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
			if odsObject.DeletedAt > 0 {
				dimVersion.EndDate = cast.ToTime(odsObject.DeletedAt)
			}

			err = t.data.Db().Table("dim_version").Where("version_id = ?", dimVersion.VersionId).Save(dimVersion).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dimVersionMap[odsObject.Id] = dimVersion
			continue
		}

		odsUpdateTime := cast.ToTime(odsObject.UpdatedAt)
		odsDeleteTime := cast.ToTime(odsObject.DeletedAt)

		if odsObject.DeletedAt > 0 {
			//删除操作
			if dimVersion.EndDate.Year() == 9999 {
				dimVersion.EndDate = odsDeleteTime
				err = t.data.Db().Table("dim_version").Where("version_id = ?", dimVersion.VersionId).Save(dimVersion).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			//更新操作
			if dimVersion.GmtModified.Compare(odsUpdateTime) <= 0 && dimVersion.EndDate.Year() == 9999 {
				dimVersion.SpaceId = odsObject.SpaceId
				dimVersion.VersionId = odsObject.Id
				dimVersion.VersionName = odsObject.VersionName
				dimVersion.GmtModified = cast.ToTime(odsObject.UpdatedAt)
				err = t.data.Db().Table("dim_version").Where("version_id = ?", dimVersion.VersionId).Save(dimVersion).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		dimVersionMap[odsObject.Id] = dimVersion

	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDimVersionTask) Stop() {
}
