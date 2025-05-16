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
type OdsToDimUserTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDimUserTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDimUserTask {
	return &OdsToDimUserTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDimUserTask) Id() string {
	return t.id
}

func (t *OdsToDimUserTask) Name() string {
	return "ods_to_dim_user_task"
}

func (t *OdsToDimUserTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDimUserTask) Status() string {
	return t.status
}

func (t *OdsToDimUserTask) Run() {

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

	var list []*ods_model.OdsUser
	err = t.data.Db().Table("ods_user_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dimUsers []*dim_model.DimUser
	err = t.data.Db().Table("dim_user").Where("user_id in ?", odsIds).Find(&dimUsers).Error
	if err != nil {
		return
	}

	dimUserMap := stream.ToMap(dimUsers, func(_ int, t *dim_model.DimUser) (int64, *dim_model.DimUser) {
		return t.UserId, t
	})

	for i := 0; i < len(list); i++ {

		odsUser := list[i]
		dimUser := dimUserMap[odsUser.Id]

		lastIdVar.VariableValue = cast.ToString(odsUser.OdsId)

		//不存在就新建
		if dimUser == nil {
			dimUser = &dim_model.DimUser{}
			dimUser.UserId = odsUser.Id
			dimUser.UserName = odsUser.UserName
			dimUser.UserNickName = odsUser.UserNickname
			dimUser.UserPinyin = odsUser.UserPinyin
			dimUser.GmtCreate = cast.ToTime(odsUser.CreatedAt)
			dimUser.GmtModified = cast.ToTime(odsUser.UpdatedAt)
			dimUser.StartDate = cast.ToTime(odsUser.CreatedAt)
			dimUser.EndDate = date.ParseInLocation("2006-01-02 15:04:05", "9999-12-31 00:00:00")

			//如果是删除的信息
			if odsUser.DeletedAt > 0 {
				dimUser.EndDate = cast.ToTime(odsUser.DeletedAt)
			}

			if odsUser.UserName == "" {
				if dimUser.UserName == "" {
					dimUser.UserName = "anonymous"
				}
				dimUser.EndDate = cast.ToTime(odsUser.UpdatedAt)
			}

			err = t.data.Db().Table("dim_user").Where("user_id = ?", dimUser.UserId).Save(dimUser).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dimUserMap[odsUser.Id] = dimUser
			continue
		}

		odsUpdateTime := cast.ToTime(odsUser.UpdatedAt)
		odsDeleteTime := cast.ToTime(odsUser.DeletedAt)

		if odsUser.DeletedAt > 0 {
			//删除操作
			if dimUser.EndDate.Year() == 9999 || dimUser.UserName == "" {
				dimUser.EndDate = odsDeleteTime
				err = t.data.Db().Table("dim_user").Where("user_id = ?", dimUser.UserId).Save(dimUser).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			//更新操作
			if dimUser.GmtModified.Compare(odsUpdateTime) <= 0 && dimUser.EndDate.Year() == 9999 {
				dimUser.UserId = odsUser.Id
				dimUser.UserName = odsUser.UserName
				dimUser.GmtModified = cast.ToTime(odsUser.UpdatedAt)
				err = t.data.Db().Table("dim_user").Where("user_id = ?", dimUser.UserId).Save(dimUser).Error
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		dimUserMap[odsUser.Id] = dimUser

	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDimUserTask) Stop() {
}
