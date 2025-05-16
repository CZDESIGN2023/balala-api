package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"go-cs/internal/dwh/data"
	dwd_model "go-cs/internal/dwh/model/dwd"
	ods_model "go-cs/internal/dwh/model/ods"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/utils/date"
	"go-cs/pkg/stream"
	"gorm.io/gorm"
)

// converts data from ODS to DIM space
type OdsToDwdWitemTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDwdWitemTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDwdWitemTask {
	return &OdsToDwdWitemTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDwdWitemTask) Id() string {
	return t.id
}

func (t *OdsToDwdWitemTask) Name() string {
	return "ods_to_dwd_witem_task"
}

func (t *OdsToDwdWitemTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDwdWitemTask) Status() string {
	return t.status
}

func (t *OdsToDwdWitemTask) Run() {

	if t.status == pkg.TASK_STATUS_RUNNING {
		return
	}

	defer func() {
		t.status = pkg.TASK_STATUS_READY
	}()

	t.status = pkg.TASK_STATUS_RUNNING

	//拉链表实现
	//每次拿1000条
	//需要知道从哪里开始重新拿-》获取最后一次的id？或者最后一个的时间？ 按时间段来拿？
	lastIdVar, err := t.variablesRepo.GetVariablesByName(t.FullName(), "last_id")
	if err != nil {
		fmt.Println(err)
		return
	}

	var list []*ods_model.OdsWitem
	err = t.data.Db().Table("ods_witem_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dwdWitems []*dwd_model.DwdWitem
	err = t.data.Db().Table("dwd_witem").Where("work_item_id in ? and end_date = ?", odsIds, Endless).Find(&dwdWitems).Error
	if err != nil {
		return
	}

	dwdWitemMap := stream.ToMap(dwdWitems, func(_ int, t *dwd_model.DwdWitem) (int64, *dwd_model.DwdWitem) {
		return t.WorkItemId, t
	})

	for i := 0; i < len(list); i++ {

		odsWitem := list[i]

		lastIdVar.VariableValue = cast.ToString(odsWitem.OdsId)

		//解析doc json
		var odsWitemDoc *ods_model.OdsWitemDoc
		err := json.Unmarshal([]byte(odsWitem.Doc), &odsWitemDoc)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dwdWitem := dwdWitemMap[odsWitem.Id]

		//不存在就新建
		if dwdWitem == nil {
			dwdWitem := t.convertToDwdWitem(odsWitem, odsWitemDoc)
			err = t.data.Db().Table("dwd_witem").Create(dwdWitem).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dwdWitemMap[odsWitem.Id] = dwdWitem
			continue
		}

		if odsWitem.DeletedAt > 0 {
			//删除操作, 让最后一个状态过期
			if dwdWitem.EndDate.Year() == 9999 {
				dwdWitem.EndDate = cast.ToTime(odsWitem.DeletedAt)
				err = t.data.Db().Table("dwd_witem").Where("work_item_id = ? and end_date = ?", dwdWitem.WorkItemId, Endless).
					UpdateColumns(map[string]interface{}{
						"end_date": dwdWitem.EndDate,
					}).Error
				if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
					t.deleteItem(dwdWitem.WorkItemId, Endless)
				}

				dwdWitemMap[odsWitem.Id] = dwdWitem
			}
		} else {
			//更新操作
			if dwdWitem.GmtModified.Compare(cast.ToTime(odsWitem.UpdatedAt)) <= 0 && dwdWitem.EndDate.Year() == 9999 {

				newDwdWitem := t.convertToDwdWitem(odsWitem, odsWitemDoc)
				//比较内容变化 相同变化不处理
				isSame := dwdWitem.DeepEqual(newDwdWitem)
				if isSame {
					continue
				}

				txErr := t.data.Db().Transaction(func(tx *gorm.DB) error {
					//然最后一个状态过期，然后写入新的纬度数据
					dwdWitem.EndDate = cast.ToTime(odsWitem.OdsOpTs)
					err = t.data.Db().Table("dwd_witem").Where("work_item_id = ? and end_date = ?", dwdWitem.WorkItemId, Endless).
						UpdateColumn("end_date", dwdWitem.EndDate).Error
					if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
						t.deleteItem(dwdWitem.WorkItemId, Endless)
					}

					newDwdWitem.StartDate = dwdWitem.EndDate
					newDwdWitem.EndDate = date.ParseInLocation("2006-01-02 15:04:05", Endless)
					err = t.data.Db().Table("dwd_witem").Create(newDwdWitem).Error
					if err != nil {
						fmt.Println(err)
						return err
					}

					return nil
				})

				if txErr != nil {
					fmt.Println(txErr)
					continue
				}

				dwdWitemMap[odsWitem.Id] = newDwdWitem
			}
		}

	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDwdWitemTask) Stop() {
}

func (t *OdsToDwdWitemTask) convertToDwdWitem(odsWitem *ods_model.OdsWitem, odsWitemDoc *ods_model.OdsWitemDoc) *dwd_model.DwdWitem {

	directors, _ := json.Marshal(odsWitemDoc.Directors)
	nodeDirectors, _ := json.Marshal(odsWitemDoc.NodeDirectors)
	participators, _ := json.Marshal(odsWitemDoc.Participators)

	dwdWitem := &dwd_model.DwdWitem{
		SpaceId:         odsWitem.SpaceId,
		WorkItemId:      odsWitem.Id,
		UserId:          odsWitem.UserId,
		StatusId:        odsWitem.WorkItemStatusId,
		ObjectId:        odsWitem.WorkObjectId,
		VersionId:       odsWitem.VersionId,
		WorkItemTypeKey: odsWitem.WorkItemTypeKey,
		LastStatusAt:    odsWitem.LastStatusAt,
		PlanStartAt:     odsWitemDoc.PlanStartAt,
		PlanCompleteAt:  odsWitemDoc.PlanCompleteAt,
		Priority:        odsWitemDoc.Priority,
		Directors:       string(directors),
		NodeDirectors:   string(nodeDirectors),
		Participators:   string(participators),
	}

	dwdWitem.GmtCreate = cast.ToTime(odsWitem.CreatedAt)
	dwdWitem.GmtModified = cast.ToTime(odsWitem.UpdatedAt)
	dwdWitem.StartDate = cast.ToTime(odsWitem.CreatedAt)
	dwdWitem.EndDate = date.ParseInLocation("2006-01-02 15:04:05", Endless)
	if odsWitem.DeletedAt > 0 {
		dwdWitem.EndDate = cast.ToTime(odsWitem.DeletedAt)
	}

	return dwdWitem
}

func (t *OdsToDwdWitemTask) deleteItem(id int64, endDate string) error {
	err := t.data.Db().Exec("DELETE FROM dwd_witem WHERE work_item_id = ? and end_date = ?", id, endDate).Error

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
