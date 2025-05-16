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
type OdsToDimObjectTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDimObjectTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDimObjectTask {
	return &OdsToDimObjectTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDimObjectTask) Id() string {
	return t.id
}

func (t *OdsToDimObjectTask) Name() string {
	return "ods_to_dim_object_task"
}

func (t *OdsToDimObjectTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDimObjectTask) Status() string {
	return t.status
}

func (t *OdsToDimObjectTask) Run() {

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

	var list []*ods_model.OdsObject
	err = t.data.Db().Table("ods_object_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dimObjects []*dim_model.DimObject
	err = t.data.Db().Table("dim_object").Where("object_id in ?", odsIds).Find(&dimObjects).Error
	if err != nil {
		return
	}

	dimObjectMap := stream.ToMap(dimObjects, func(_ int, t *dim_model.DimObject) (int64, *dim_model.DimObject) {
		return t.ObjectId, t
	})

	for i := 0; i < len(list); i++ {

		odsObject := list[i]
		dimObject := dimObjectMap[odsObject.Id]

		lastIdVar.VariableValue = cast.ToString(odsObject.OdsId)

		//不存在就新建
		if dimObject == nil {
			dimObject = &dim_model.DimObject{}
			dimObject.SpaceId = odsObject.SpaceId
			dimObject.ObjectId = odsObject.Id
			dimObject.ObjectName = odsObject.WorkObjectName
			dimObject.GmtCreate = cast.ToTime(odsObject.CreatedAt)
			dimObject.GmtModified = cast.ToTime(odsObject.UpdatedAt)
			dimObject.StartDate = cast.ToTime(odsObject.CreatedAt)
			dimObject.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
			if odsObject.DeletedAt > 0 {
				dimObject.EndDate = cast.ToTime(odsObject.DeletedAt)
			}

			err = t.data.Db().Table("dim_object").Where("object_id = ?", dimObject.ObjectId).Save(dimObject).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dimObjectMap[odsObject.Id] = dimObject
			continue
		}

		odsUpdateTime := cast.ToTime(odsObject.UpdatedAt)
		odsDeleteTime := cast.ToTime(odsObject.DeletedAt)

		if odsObject.DeletedAt > 0 {
			//删除操作
			if dimObject.EndDate.Year() == 9999 {
				dimObject.EndDate = odsDeleteTime
				err = t.data.Db().Table("dim_object").Where("object_id = ?", dimObject.ObjectId).Save(dimObject).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			//更新操作
			if dimObject.GmtModified.Compare(odsUpdateTime) <= 0 && dimObject.EndDate.Year() == 9999 {
				dimObject.SpaceId = odsObject.SpaceId
				dimObject.ObjectId = odsObject.Id
				dimObject.ObjectName = odsObject.WorkObjectName
				dimObject.GmtModified = cast.ToTime(odsObject.UpdatedAt)
				err = t.data.Db().Table("dim_object").Where("object_id = ?", dimObject.ObjectId).Save(dimObject).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		dimObjectMap[odsObject.Id] = dimObject
	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDimObjectTask) Stop() {
}
