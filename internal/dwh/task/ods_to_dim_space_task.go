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
type OdsToDimSpaceTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDimSpaceTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDimSpaceTask {
	return &OdsToDimSpaceTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDimSpaceTask) Id() string {
	return t.id
}

func (t *OdsToDimSpaceTask) Name() string {
	return "ods_to_dim_space_task"
}

func (t *OdsToDimSpaceTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDimSpaceTask) Status() string {
	return t.status
}

func (t *OdsToDimSpaceTask) Run() {

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

	var list []*ods_model.OdsSpace
	err = t.data.Db().Table("ods_space_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		return
	}

	odsSpaceIds := make([]int64, 0)
	for _, v := range list {
		odsSpaceIds = append(odsSpaceIds, v.Id)
	}
	odsSpaceIds = stream.Unique(odsSpaceIds)

	var dimSpaces []*dim_model.DimSpace
	err = t.data.Db().Table("dim_space").Where("space_id in ?", odsSpaceIds).Find(&dimSpaces).Error
	if err != nil {
		return
	}

	dimSpaceMap := stream.ToMap(dimSpaces, func(_ int, t *dim_model.DimSpace) (int64, *dim_model.DimSpace) {
		return t.SpaceId, t
	})

	for i := 0; i < len(list); i++ {

		odsSpace := list[i]
		dimSpace := dimSpaceMap[odsSpace.Id]

		lastIdVar.VariableValue = cast.ToString(odsSpace.OdsId)

		//不存在就新建
		if dimSpace == nil {
			dimSpace = &dim_model.DimSpace{}
			dimSpace.SpaceId = odsSpace.Id
			dimSpace.SpaceName = odsSpace.SpaceName
			dimSpace.GmtCreate = cast.ToTime(odsSpace.CreatedAt)
			dimSpace.GmtModified = cast.ToTime(odsSpace.UpdatedAt)
			dimSpace.StartDate = cast.ToTime(odsSpace.CreatedAt)
			dimSpace.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")
			if odsSpace.DeletedAt > 0 {
				dimSpace.EndDate = cast.ToTime(odsSpace.DeletedAt)
			}

			err = t.data.Db().Table("dim_space").Where("space_id = ?", dimSpace.SpaceId).Save(dimSpace).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dimSpaceMap[odsSpace.Id] = dimSpace
			continue
		}

		odsUpdateTime := cast.ToTime(odsSpace.UpdatedAt)
		odsDeleteTime := cast.ToTime(odsSpace.DeletedAt)

		if odsSpace.DeletedAt > 0 {
			//删除操作
			if dimSpace.EndDate.Year() == 9999 {
				dimSpace.EndDate = odsDeleteTime
				err = t.data.Db().Table("dim_space").Where("space_id = ?", dimSpace.SpaceId).Save(dimSpace).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			//更新操作
			if dimSpace.GmtModified.Compare(odsUpdateTime) <= 0 && dimSpace.EndDate.Year() == 9999 {
				dimSpace.SpaceId = odsSpace.Id
				dimSpace.SpaceName = odsSpace.SpaceName
				dimSpace.GmtModified = cast.ToTime(odsSpace.UpdatedAt)
				err = t.data.Db().Table("dim_space").Where("space_id = ?", dimSpace.SpaceId).Save(dimSpace).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		dimSpaceMap[odsSpace.Id] = dimSpace
	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDimSpaceTask) Stop() {
}
