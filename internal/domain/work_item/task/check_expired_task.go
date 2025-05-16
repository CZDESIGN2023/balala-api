package task

import (
	"context"
	"go-cs/api/notify"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/consts"
	search22 "go-cs/internal/domain/search/search2"
	"go-cs/internal/utils"
	"go-cs/internal/utils/date"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"
)

// 项目异常检查，逾期任务太多
func (task *WorkItemTask) CheckWorkItemExpired() {
	begin := time.Now()
	task.log.Info("CheckWorkItemExpired start")
	defer task.log.Infof("CheckWorkItemExpired end %v", time.Now().Sub(begin))

	ctx := context.Background()

	today := date.TodayBegin()
	//yesterday := today.Add(time.Hour * -24)

	var list []*search22.Model
	res := task.data.RoDB(ctx).Table((&db.SpaceWorkItemV2{}).TableName()+" t").
		Joins("INNER JOIN work_item_status s ON t.work_item_status_id = s.id").
		Where("s.status_type != ?", consts.WorkItemStatusType_Archived).
		Where("doc->>'$.plan_complete_at' > 0").
		Where("doc->>'$.plan_complete_at' < ?", today.Unix()).
		Select("t.`id`, t.`space_id`, t.`work_item_name`, t.`user_id`, t.`doc` -> '$.directors' AS `directors`, t.`doc` -> '$.followers' AS `followers`").
		Find(&list)

	if res.Error != nil {
		task.log.Error(res.Error)
		return
	}

	for _, workItem := range list {

		bus.Emit(notify.Event_WorkItemExpired, &event.WorkItemExpired{
			Event:      notify.Event_WorkItemExpired,
			WorkItemId: workItem.Id,
		})
	}
}

func (task *WorkItemTask) CheckFlowNodeExpired() {
	begin := time.Now()
	task.log.Info("CheckExpiredNode start")
	defer task.log.Infof("CheckExpiredNode end %v", time.Now().Sub(begin))

	ctx := context.Background()

	today := date.TodayBegin()

	var list []*db.SpaceWorkItemFlowV2
	res := task.data.RoDB(ctx).Table((&db.SpaceWorkItemFlowV2{}).TableName()+" f").
		Joins("INNER JOIN space_work_item_v2 t ON t.id = f.work_item_id").
		Joins("INNER JOIN work_item_status s ON t.work_item_status_id = s.id").
		Where("s.status_type != ?", consts.WorkItemStatusType_Archived).
		Where("t.pid = 0").
		Where("t.work_item_type_key = 'task' AND f.flow_node_status != 3 OR t.work_item_type_key = 'state_task' AND f.flow_node_code = t.work_item_status_key").
		Where("f.plan_complete_at > 0").
		Where("f.plan_complete_at < ?", today.Unix()).
		Find(&list)

	if res.Error != nil {
		task.log.Error(res.Error)
		return
	}

	templateIds := stream.Map(list, func(v *db.SpaceWorkItemFlowV2) int64 {
		return v.FlowTemplateId
	})

	templateMap, err := task.workFlowRepo.FlowTemplateMap(ctx, templateIds)
	if err != nil {
		task.log.Error(err)
		return
	}

	for _, v := range list {

		template := templateMap[v.FlowTemplateId]
		if template == nil {
			continue
		}

		var nodeName string
		switch template.FlowMode {
		case consts.FlowMode_WorkFlow:
			node := template.WorkFLowConfig.GetNode(v.FlowNodeCode)
			if node != nil {
				nodeName = node.Name
			}
		case consts.FlowMode_StateFlow:
			node := template.StateFlowConfig.GetNode(v.FlowNodeCode)
			if node != nil {
				nodeName = node.Name
			}
		default:
			continue
		}

		bus.Emit(notify.Event_WorkItemFlowNodeExpired, &event.WorkItemFlowNodeExpired{
			Event:         notify.Event_WorkItemFlowNodeExpired,
			WorkItemId:    v.WorkItemId,
			NodeName:      nodeName,
			NodeDirectors: utils.ToInt64Array(utils.StrToStrArray(v.Directors)),
		})
	}
}
