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

// converts data from ODS to DIM space
type OdsToDwdWitemFlowNodeTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewOdsToDwdWitemFlowNodeTask(
	id string,
	ctx *pkg.TaskContext,
) *OdsToDwdWitemFlowNodeTask {
	return &OdsToDwdWitemFlowNodeTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *OdsToDwdWitemFlowNodeTask) Id() string {
	return t.id
}

func (t *OdsToDwdWitemFlowNodeTask) Name() string {
	return "ods_to_dwd_witem_flow_node_task"
}

func (t *OdsToDwdWitemFlowNodeTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *OdsToDwdWitemFlowNodeTask) Status() string {
	return t.status
}

func (t *OdsToDwdWitemFlowNodeTask) Run() {

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

	var list []*ods_model.OdsWitemFlowNode
	err = t.data.Db().Table("ods_witem_flow_node_d").Where("_id > ?", cast.ToInt64(lastIdVar.VariableValue)).Order("_id ASC").Limit(3000).Find(&list).Error
	if err != nil {
		return
	}

	odsIds := make([]int64, 0)
	for _, v := range list {
		odsIds = append(odsIds, v.Id)
	}
	odsIds = stream.Unique(odsIds)

	var dwdWitemFlowNodes []*dwd_model.DwdWitemFlowNode
	err = t.data.Db().Table("dwd_witem_flow_node").Where("node_id in ? and end_date = ?", odsIds, Endless).Find(&dwdWitemFlowNodes).Error
	if err != nil {
		return
	}

	dwdWitemFlowNodeMap := stream.ToMap(dwdWitemFlowNodes, func(_ int, t *dwd_model.DwdWitemFlowNode) (int64, *dwd_model.DwdWitemFlowNode) {
		return t.NodeId, t
	})

	for i := 0; i < len(list); i++ {

		OdsWitemFlowNode := list[i]

		lastIdVar.VariableValue = cast.ToString(OdsWitemFlowNode.OdsId)

		dwdWitemFlowNode := dwdWitemFlowNodeMap[OdsWitemFlowNode.Id]

		//不存在就新建
		if dwdWitemFlowNode == nil {
			dwdWitemFlowNode := t.convertToDwdWitemFlowNode(OdsWitemFlowNode)
			err = t.data.Db().Table("dwd_witem_flow_node").Create(dwdWitemFlowNode).Error
			if err != nil {
				fmt.Println(err)
				continue
			}

			dwdWitemFlowNodeMap[OdsWitemFlowNode.Id] = dwdWitemFlowNode
			continue
		}

		if OdsWitemFlowNode.DeletedAt > 0 {
			//删除操作, 让最后一个状态过期
			if dwdWitemFlowNode.EndDate.Year() == 9999 {
				dwdWitemFlowNode.EndDate = cast.ToTime(OdsWitemFlowNode.DeletedAt)

				err = t.data.Db().Table("dwd_witem_flow_node").
					Where("node_id = ? and end_date = ?", dwdWitemFlowNode.NodeId, Endless).
					UpdateColumns(map[string]interface{}{
						"end_date": dwdWitemFlowNode.EndDate,
					}).Error
				if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
					t.deleteItem(dwdWitemFlowNode.NodeId, Endless)
				}

				dwdWitemFlowNodeMap[OdsWitemFlowNode.Id] = dwdWitemFlowNode
			}
		} else {
			//更新操作
			if dwdWitemFlowNode.GmtModified.Compare(cast.ToTime(OdsWitemFlowNode.OdsOpTs)) <= 0 && dwdWitemFlowNode.EndDate.Year() == 9999 {

				newDwdWitemFlowNode := t.convertToDwdWitemFlowNode(OdsWitemFlowNode)
				//比较内容变化 相同变化不处理
				isSame := dwdWitemFlowNode.DeepEqual(newDwdWitemFlowNode)
				if isSame {
					continue
				}

				txErr := t.data.Db().Transaction(func(tx *gorm.DB) error {
					//然最后一个状态过期，然后写入新的纬度数据
					dwdWitemFlowNode.EndDate = cast.ToTime(OdsWitemFlowNode.OdsOpTs)
					err = t.data.Db().Table("dwd_witem_flow_node").
						Where("node_id = ? and end_date = ?", dwdWitemFlowNode.NodeId, Endless).
						UpdateColumns(map[string]interface{}{
							"end_date": dwdWitemFlowNode.EndDate,
						}).Error
					if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
						t.deleteItem(dwdWitemFlowNode.NodeId, Endless)
					}

					newDwdWitemFlowNode.StartDate = dwdWitemFlowNode.EndDate
					newDwdWitemFlowNode.EndDate = date.ParseInLocation("2006-01-02 15:04:05", Endless)
					err = t.data.Db().Table("dwd_witem_flow_node").Create(newDwdWitemFlowNode).Error
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

				dwdWitemFlowNodeMap[OdsWitemFlowNode.Id] = newDwdWitemFlowNode
			}
		}

	}

	err = t.variablesRepo.SaveVariables(lastIdVar)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *OdsToDwdWitemFlowNodeTask) convertToDwdWitemFlowNode(OdsWitemFlowNode *ods_model.OdsWitemFlowNode) *dwd_model.DwdWitemFlowNode {

	dwdWitemFlowNode := &dwd_model.DwdWitemFlowNode{}
	dwdWitemFlowNode.SpaceId = OdsWitemFlowNode.SpaceId
	dwdWitemFlowNode.WorkItemId = OdsWitemFlowNode.WorkItemId
	dwdWitemFlowNode.NodeId = OdsWitemFlowNode.Id
	dwdWitemFlowNode.NodeCode = OdsWitemFlowNode.FlowNodeCode
	dwdWitemFlowNode.NodeStatus = OdsWitemFlowNode.FlowNodeStatus

	dwdWitemFlowNode.PlanStartAt = OdsWitemFlowNode.PlanStartAt
	dwdWitemFlowNode.PlanCompleteAt = OdsWitemFlowNode.PlanCompleteAt
	dwdWitemFlowNode.Directors = OdsWitemFlowNode.Directors
	dwdWitemFlowNode.GmtCreate = cast.ToTime(OdsWitemFlowNode.CreatedAt)
	dwdWitemFlowNode.GmtModified = cast.ToTime(OdsWitemFlowNode.UpdatedAt)
	dwdWitemFlowNode.StartDate = cast.ToTime(OdsWitemFlowNode.CreatedAt)
	dwdWitemFlowNode.EndDate = date.ParseInLocation("2006-01-02 15:04:05", Endless)
	if OdsWitemFlowNode.DeletedAt > 0 {
		dwdWitemFlowNode.EndDate = cast.ToTime(OdsWitemFlowNode.DeletedAt)
	}

	return dwdWitemFlowNode
}

func (t *OdsToDwdWitemFlowNodeTask) Stop() {

}

func (t *OdsToDwdWitemFlowNodeTask) deleteItem(id int64, endDate string) error {
	err := t.data.Db().Exec("DELETE FROM dwd_witem_flow_node WHERE node_id = ? and end_date = ?", id, endDate).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
