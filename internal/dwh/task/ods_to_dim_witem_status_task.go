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
type OdsToDimWitemStatusTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDimWitemStatusTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDimWitemStatusTask {
	return &OdsToDimWitemStatusTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDimWitemStatusTask) Id() string {
	return t.id
}

func (t *OdsToDimWitemStatusTask) Name() string {
	return "ods_to_dim_witem_status_task"
}

func (t *OdsToDimWitemStatusTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.id
	}
	return t.Name() + ":" + t.id
}

func (t *OdsToDimWitemStatusTask) Status() string {
	return t.status
}

func (t *OdsToDimWitemStatusTask) Run() {

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

	var list []*ods_model.OdsWitemStatus
	err = t.data.Db().Table("ods_witem_status_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dimStatuses []*dim_model.DimWitemStatus
	err = t.data.Db().Table("dim_witem_status").Where("status_id in ?", odsIds).Find(&dimStatuses).Error
	if err != nil {
		return
	}

	dimStatusMap := stream.ToMap(dimStatuses, func(_ int, t *dim_model.DimWitemStatus) (int64, *dim_model.DimWitemStatus) {
		return t.StatusId, t
	})

	for _, v := range list {
		dimStatus := dimStatusMap[v.Id]

		lastIdVar.VariableValue = cast.ToString(v.OdsId)

		//不存在就新建
		if dimStatus == nil {
			dimStatus = &dim_model.DimWitemStatus{
				SpaceId:    v.SpaceId,
				StatusId:   v.Id,
				StatusName: v.Name,
				StatusKey:  v.Key,
				StatusVal:  v.Val,
				StatusType: v.StatusType,
				FlowScope:  v.FlowScope,
			}

			dimStatus.GmtCreate = cast.ToTime(v.CreatedAt)
			dimStatus.GmtModified = cast.ToTime(v.UpdatedAt)
			dimStatus.StartDate = cast.ToTime(v.CreatedAt)
			dimStatus.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")

			if v.DeletedAt > 0 {
				dimStatus.EndDate = cast.ToTime(v.DeletedAt)
			}

			err = t.data.Db().Table("dim_witem_status").Where("status_id = ?", dimStatus.StatusId).Save(dimStatus).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dimStatusMap[v.Id] = dimStatus
			continue
		}

		if v.DeletedAt > 0 {
			//删除操作
			odsDeleteTime := cast.ToTime(v.DeletedAt)
			if dimStatus.EndDate.Year() == 9999 {
				dimStatus.EndDate = odsDeleteTime
				err = t.data.Db().Table("dim_witem_status").Where("status_id = ?", dimStatus.StatusId).Save(dimStatus).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			//更新操作
			odsUpdateTime := cast.ToTime(v.OdsOpTs)

			if dimStatus.GmtModified.Compare(odsUpdateTime) <= 0 && dimStatus.EndDate.Year() == 9999 {
				dimStatus.SpaceId = v.SpaceId
				dimStatus.StatusId = v.Id
				dimStatus.StatusName = v.Name
				dimStatus.StatusKey = v.Key
				dimStatus.StatusVal = v.Val
				dimStatus.StatusType = v.StatusType
				dimStatus.FlowScope = v.FlowScope
				dimStatus.GmtModified = cast.ToTime(v.UpdatedAt)
				err = t.data.Db().Table("dim_witem_status").Where("status_id = ?", dimStatus.StatusId).Save(dimStatus).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		dimStatusMap[v.Id] = dimStatus

	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDimWitemStatusTask) Stop() {
}
