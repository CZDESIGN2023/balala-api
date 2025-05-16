package task

import (
	"errors"
	"fmt"
	"go-cs/internal/dwh/data"
	dwd_model "go-cs/internal/dwh/model/dwd"
	ods_model "go-cs/internal/dwh/model/ods"
	"go-cs/internal/dwh/pkg"
	"go-cs/internal/utils/date"
	"go-cs/pkg/stream"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// converts data from ODS to dwd space
type OdsToDwdMemberTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDwdMemberTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDwdMemberTask {
	return &OdsToDwdMemberTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDwdMemberTask) Id() string {
	return t.id
}

func (t *OdsToDwdMemberTask) Name() string {
	return "ods_to_dwd_member_task"
}

func (t *OdsToDwdMemberTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDwdMemberTask) Status() string {
	return t.status
}

func (t *OdsToDwdMemberTask) Run() {

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

	var list []*ods_model.OdsMember
	err = t.data.Db().Table("ods_member_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dwdMembers []*dwd_model.DwdMember
	err = t.data.Db().Table("dwd_member").Where("member_id in ? and end_date = ?", odsIds, Endless).Find(&dwdMembers).Error
	if err != nil {
		return
	}

	dwdMemberMap := stream.ToMap(dwdMembers, func(_ int, t *dwd_model.DwdMember) (int64, *dwd_model.DwdMember) {
		return t.MemberId, t
	})

	for i := 0; i < len(list); i++ {

		odsMember := list[i]
		dwdMember := dwdMemberMap[odsMember.Id]

		lastIdVar.VariableValue = cast.ToString(odsMember.OdsId)

		//不存在就新建
		if dwdMember == nil {
			dwdMember := t.convertToDwdMember(odsMember)
			err = t.data.Db().Table("dwd_member").Create(dwdMember).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dwdMemberMap[odsMember.Id] = dwdMember
			continue
		}

		if odsMember.DeletedAt > 0 {
			//删除操作, 让最后一个状态过期
			if dwdMember.EndDate.Year() == 9999 {
				dwdMember.EndDate = cast.ToTime(odsMember.DeletedAt)

				err = t.data.Db().Table("dwd_member").Where("member_id = ? and end_date = ?", dwdMember.MemberId, Endless).
					UpdateColumns(map[string]interface{}{
						"end_date": dwdMember.EndDate,
					}).Error
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					t.deleteDwdItem(dwdMember.MemberId, Endless)
				}

				dwdMemberMap[odsMember.Id] = dwdMember
			}
		} else {
			//更新操作
			if dwdMember.GmtModified.Compare(cast.ToTime(odsMember.OdsOpTs)) <= 0 && dwdMember.EndDate.Year() == 9999 {

				newDwdMember := t.convertToDwdMember(odsMember)
				//比较内容变化 相同变化不处理
				isSame := dwdMember.DeepEqual(newDwdMember)
				if isSame {
					continue
				}

				txErr := t.data.Db().Transaction(func(tx *gorm.DB) error {
					//然最后一个状态过期，然后写入新的纬度数据
					dwdMember.EndDate = cast.ToTime(odsMember.OdsOpTs)
					err = t.data.Db().Table("dwd_member").Where("member_id = ? and end_date = ?", dwdMember.MemberId, Endless).
						UpdateColumns(map[string]interface{}{
							"end_date": dwdMember.EndDate,
						}).Error
					if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
						t.deleteDwdItem(dwdMember.MemberId, Endless)
					}

					newDwdMember.StartDate = dwdMember.EndDate
					newDwdMember.EndDate = date.ParseInLocation("2006-01-02 15:04:05", Endless)
					err = t.data.Db().Table("dwd_member").Create(newDwdMember).Error
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

				dwdMemberMap[odsMember.Id] = newDwdMember
			}
		}

	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDwdMemberTask) Stop() {
}

func (t *OdsToDwdMemberTask) convertToDwdMember(odsMember *ods_model.OdsMember) *dwd_model.DwdMember {

	dwdMember := &dwd_model.DwdMember{}
	dwdMember.SpaceId = odsMember.SpaceId
	dwdMember.MemberId = odsMember.Id
	dwdMember.UserId = odsMember.UserId
	dwdMember.RoleId = odsMember.RoleId
	dwdMember.GmtCreate = cast.ToTime(odsMember.CreatedAt)
	dwdMember.GmtModified = cast.ToTime(odsMember.UpdatedAt)
	dwdMember.StartDate = cast.ToTime(odsMember.CreatedAt)
	dwdMember.EndDate = date.ParseInLocation("2006-01-02 15:04:05", Endless)
	if odsMember.DeletedAt > 0 {
		dwdMember.EndDate = cast.ToTime(odsMember.DeletedAt)
	}
	return dwdMember
}

func (t *OdsToDwdMemberTask) deleteDwdItem(id int64, endDate string) error {
	err := t.data.Db().Exec("DELETE FROM dwd_member WHERE member_id = ? and end_date = ?", id, endDate).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
