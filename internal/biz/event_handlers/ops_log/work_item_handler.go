package ops_log

import (
	"context"
	"fmt"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	wVersion_domain "go-cs/internal/domain/space_work_version"
	user_domain "go-cs/internal/domain/user"
	flow_config "go-cs/internal/domain/work_flow/flow_tplt_config"
	witem_domain "go-cs/internal/domain/work_item"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item_role"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/oper"
	"go-cs/pkg/stream"
	"strings"

	"github.com/spf13/cast"
)

func (s *OpsLogEventHandlers) workItemCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateWorkItem)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeAdd,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
		}

		result.OperMsg = "创建任务"
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemSubTaskCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateWorkItemSubTask)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		workItem, _ := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, &witem_repo.WithDocOption{
			PlanTime:    true,
			Directors:   true,
			ProcessRate: true,
		}, nil)

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeAdd,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
			ModuleFlags:  []oper.ModuleFlag{oper.ModuleFlag_subWorkItem},
		}

		result.OperMsg = "创建任务"
		s.invokeOperLog(ctx, operLogger, result)

		userMap, _ := s.userRepo.UserMap(ctx, workItem.Doc.Directors.ToInt64s())
		names := stream.Map(stream.Values(userMap), func(v *user_domain.User) string {
			return v.UserNickname
		})

		parentWorkItem, _ := s.wItemRepo.GetWorkItem(ctx, workItem.Pid, nil, nil)

		parentResult := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.ParentWorkItemId),
			ModuleTitle:  parentWorkItem.WorkItemName,
		}
		parentResult.OperMsg = "新增 " + W("子任务：") +
			"任务名称" + Q(opsLog.WorkItemName) +
			"，负责人" + Q(strings.Join(names, "，")) +
			"，排期" + Q(T(workItem.Doc.PlanStartAt, workItem.Doc.PlanCompleteAt)) +
			"，进度" + Q(cast.ToString(workItem.Doc.ProcessRate)+"%")
		s.invokeOperLog(ctx, operLogger, parentResult)
	}
}

func (s *OpsLogEventHandlers) workItemModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifyWorkItem)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
		}

		var oldValue, newValue, filedName string
		var fileUpdates = make([]string, 0)
		for _, v := range opsLog.Updates {

			switch v.Field {
			case "workItemName":
				filedName = "任务名称"
				oldValue = cast.ToString(v.OldValue)
				newValue = cast.ToString(v.NewValue)
			case "planTime":

				oldTimes := v.OldValue.([]any)
				newTimes := v.NewValue.([]any)

				if len(oldTimes) != 2 || len(newTimes) != 2 {
					continue
				}

				filedName = "总排期"
				newValue = T(cast.ToInt64(newTimes[0]), cast.ToInt64(newTimes[1]))
				oldValue = T(cast.ToInt64(oldTimes[0]), cast.ToInt64(oldTimes[1]))

			case "processRate":
				filedName = "进度"
				newValue = cast.ToString(v.NewValue) + "%"
				oldValue = cast.ToString(v.OldValue) + "%"

			case "priority":
				filedName = "优先级"
				newValue = consts.GetWorkItemPriorityName(cast.ToString(v.NewValue))
				oldValue = consts.GetWorkItemPriorityName(cast.ToString(v.OldValue))

			case "workObjectId":
				filedName = "所属模块"
				_oldVal := cast.ToInt64(v.OldValue)
				_newVal := cast.ToInt64(v.NewValue)
				objectMap, _ := s.wObjectRepo.SpaceWorkObjectMapByObjectIds(ctx, []int64{_oldVal, _newVal})

				oldValue = objectMap[_oldVal].WorkObjectName
				newValue = objectMap[_newVal].WorkObjectName

			case "versionId":
				filedName = "所属版本"
				_oldVal := cast.ToInt64(v.OldValue)
				_newVal := cast.ToInt64(v.NewValue)

				objectList, _ := s.wVersionRepo.GetSpaceWorkVersionByIds(ctx, []int64{_oldVal, _newVal})
				objectMap := stream.ToMap(objectList, func(i int, v *wVersion_domain.SpaceWorkVersion) (int64, *wVersion_domain.SpaceWorkVersion) {
					return v.Id, v
				})

				oldValue = objectMap[_oldVal].VersionName
				newValue = objectMap[_newVal].VersionName

			case "describe":
				filedName = "任务描述"

				newValue = utils.ClearRichTextToPlanText(cast.ToString(v.NewValue), true)
				oldValue = utils.ClearRichTextToPlanText(cast.ToString(v.OldValue), true)
			case "remark":
				filedName = "交付备注"

				newValue = utils.ClearRichTextToPlanText(cast.ToString(v.NewValue), true)
				oldValue = utils.ClearRichTextToPlanText(cast.ToString(v.OldValue), true)
			case "followers":
				filedName = "关注人"
				_oldVal := v.OldValue.([]int64)
				_newVal := v.NewValue.([]int64)

				userMap, _ := s.userRepo.UserMap(ctx, append(_oldVal, _newVal...))

				oldNames := stream.Map(_oldVal, func(v int64) string {
					return userMap[v].UserNickname
				})

				newNames := stream.Map(_newVal, func(v int64) string {
					return userMap[v].UserNickname
				})

				oldValue = strings.Join(oldNames, "，")
				newValue = strings.Join(newNames, "，")
			}

			fileUpdates = append(fileUpdates, "将 "+W(filedName)+" 由"+Q(oldValue)+"更新为"+Q(newValue))
		}

		result.OperMsg = strings.Join(fileUpdates, " ")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemStatusChangeHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeWorkItemStatus)
		if opsLog == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		oldStatus, err := s.wStatusRepo.GetWorkItemStatusItem(ctx, opsLog.OldWorkItemStatusId)
		if err != nil {
			continue
		}

		newStatus, err := s.wStatusRepo.GetWorkItemStatusItem(ctx, opsLog.NewWorkItemStatusId)
		if err != nil {
			continue
		}

		var msg string
		switch {
		case newStatus.IsCompleted():
			if opsLog.Pid != 0 {
				msg = "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
				if opsLog.AffectByParent {
					msg = "将父级任务 " + W("完成") + "，子任务 " + W("任务状态") + " 跟随更新为" + Q(newStatus.Name)
				}
			}
		case newStatus.IsClose():
			// 关闭
			msg = "将任务 " + W("关闭") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
			if opsLog.AffectByParent {
				msg = "将父级任务 " + W("关闭") + "，子任务 " + W("任务状态") + " 跟随更新为" + Q(newStatus.Name)
			}
		case newStatus.IsTerminated():
			// 终止
			switch opsLog.WorkItemTypeKey {
			case consts.WorkItemTypeKey_Task:
				msg = "将任务 " + W("终止") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
			case consts.WorkItemTypeKey_StateTask:
				msg = "将任务 " + W("终止") + "，原因" + Q(opsLog.Reason)
			case consts.WorkItemTypeKey_SubTask:
				msg = "将父级任务 " + W("终止") + "，子任务 " + W("任务状态") + " 跟随更新为" + Q(newStatus.Name)
			}
		case oldStatus.IsClose() && !newStatus.IsArchivedTypeState():
			// 关闭-重启
			msg = "将任务 " + W("重启") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
			if opsLog.AffectByParent {
				msg = "将父级任务 " + W("重启") + "，子任务 " + W("任务状态") + " 跟随更新为" + Q(newStatus.Name)
			}
		case oldStatus.IsTerminated() && !newStatus.IsArchivedTypeState():
			// 终止-恢复
			switch opsLog.WorkItemTypeKey {
			case consts.WorkItemTypeKey_Task:
				msg = "将任务 " + W("恢复") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
			case consts.WorkItemTypeKey_StateTask:
				msg = "将任务 " + W("恢复") + "，原因" + Q(opsLog.Reason)
			case consts.WorkItemTypeKey_SubTask:
				msg = "将父级任务 " + W("恢复") + "，子任务 " + W("任务状态") + " 跟随更新为" + Q(newStatus.Name)
			}
		case oldStatus.IsCompleted() && !newStatus.IsArchivedTypeState():
			//完成-重启
			if opsLog.Pid == 0 { //父级任务重启
				var flowNodeName string
				if opsLog.FlowNodeId != 0 {
					flowNode, _ := s.wItemRepo.GetWorkItemFlowNodeByNodeCode(ctx, opsLog.WorkItemId, opsLog.FlowNodeCode)
					flowTplt, err := s.wFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, flowNode.FlowTemplateId)
					if err == nil && flowTplt.WorkFLowConfig != nil {
						nodeConf := flowTplt.WorkFLowConfig.GetNode(opsLog.FlowNodeCode)
						if nodeConf != nil {
							flowNodeName = nodeConf.Name
						}
					}
				}

				if flowNodeName == "" {
					continue
				}
				msg = "将任务 " + W("重启") + " 至" + Q(flowNodeName) + W("节点") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
			} else { //子任务重启
				msg = "将任务 " + W("重启") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(newStatus.Name)
				if opsLog.AffectByParent {
					msg = "将父级任务 " + W("重启") + "，子任务 " + W("任务状态") + " 跟随更新为" + Q(newStatus.Name)
				}
			}
		case opsLog.WorkItemTypeKey == consts.WorkItemTypeKey_StateTask:
			msg = "将 " + W("任务状态") + " 由" + Q(oldStatus.Name) + "更新为" + Q(newStatus.Name)
			if opsLog.Reason != "" || opsLog.Remark != "" {
				msg = "将 " + W("任务状态") + " 由" + Q(oldStatus.Name) + "更新为" + Q(newStatus.Name)
				if opsLog.Reason != "" {
					msg += "，原因" + Q(opsLog.Reason)
				}
				if opsLog.Remark != "" {
					msg += "，备注" + Q(opsLog.Remark)
				}
			}

		}

		if msg != "" {
			operLogger.Operator = operUser
			result := &oper.OperResultInfo{
				SpaceId:      space.Id,
				SpaceName:    space.SpaceName,
				BusinessType: oper.BusinessTypeModify,
				ShowType:     oper.ShowTypeWorkItemNodeChange,
				ModuleType:   oper.ModuleTypeSpaceWorkItem,
				ModuleId:     int(opsLog.WorkItemId),
				ModuleTitle:  opsLog.WorkItemName,
				OperMsg:      msg,
			}

			s.invokeOperLog(ctx, operLogger, result)
		}
	}
}

func (s *OpsLogEventHandlers) workItemDirectorChangeHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeWorkItemDirector)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		allUsers := append(opsLog.OldDirectors, opsLog.NewDirectors...)
		for _, v := range opsLog.Nodes {
			allUsers = append(allUsers, v.OldDirectors...)
			allUsers = append(allUsers, v.NewDirectors...)
		}
		allUsers = stream.Unique(allUsers)
		users, _ := s.userRepo.GetUserByIds(ctx, utils.StringArrToInt64Arr(allUsers))
		userMap := stream.ToMap(users, func(idx int, item *user_domain.User) (string, *user_domain.User) {
			return cast.ToString(item.Id), item
		})

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
		}

		if len(opsLog.Nodes) == 0 { //子任务
			oldNames := stream.Map(opsLog.OldDirectors, func(v string) string {
				return userMap[v].UserNickname
			})
			newNames := stream.Map(opsLog.NewDirectors, func(v string) string {
				return userMap[v].UserNickname
			})

			result.OperMsg = "将 " + W("负责人") + " 由" + Q(strings.Join(oldNames, "，")) + "更新为" + Q(strings.Join(newNames, "，"))
			s.invokeOperLog(ctx, operLogger, result)
		} else { //-- 节点对应的负责人变化信息 -- -

			tpltConf, _ := s.wFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, opsLog.FlowTemplateId)
			if tpltConf == nil {
				continue
			}

			var nodeNameMap map[string]string

			switch tpltConf.FlowMode {
			case consts.FlowMode_WorkFlow:
				nodeNameMap = stream.ToMap(tpltConf.WorkFlowConf().Nodes, func(idx int, item *flow_config.WorkFlowNode) (string, string) {
					return item.Key, item.Name
				})
			case consts.FlowMode_StateFlow:
				nodeNameMap = stream.ToMap(tpltConf.StateFlowConf().StateFlowNodes, func(idx int, item *flow_config.StateFlowNode) (string, string) {
					return item.Key, item.Name
				})
			}

			for _, v := range opsLog.Nodes {
				result := *result // 复制，避免并发问题
				name := nodeNameMap[v.FlowNodeCode]
				if name == "" {
					continue
				}

				oldNames := stream.Map(v.OldDirectors, func(v string) string {
					return userMap[v].UserNickname
				})

				newNames := stream.Map(v.NewDirectors, func(v string) string {
					return userMap[v].UserNickname
				})

				result.OperMsg = "将" + Q(name) + W("负责人") + " 由" + Q(strings.Join(oldNames, "，")) + "更新为" + Q(strings.Join(newNames, "，"))
				s.invokeOperLog(ctx, operLogger, &result)
			}
		}

	}
}

func (s *OpsLogEventHandlers) workItemDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteWorkItem)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeDel,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
		}

		var parentWorkItem *witem_domain.WorkItem
		if opsLog.ParentWorkItemId != 0 {
			parentWorkItem, _ = s.wItemRepo.GetWorkItem(ctx, opsLog.ParentWorkItemId, &witem_repo.WithDocOption{
				PlanTime: true,
			}, nil)
		}

		userMap, _ := s.userRepo.UserMap(ctx, opsLog.Directors)
		names := stream.Map(opsLog.Directors, func(v int64) string {
			return userMap[v].UserNickname
		})

		if parentWorkItem != nil {
			result.ModuleId = int(parentWorkItem.Id)
			result.ModuleTitle = parentWorkItem.WorkItemName
			result.BusinessType = oper.BusinessTypeModify
			result.OperMsg = "删除 " + W("子任务：") +
				"任务名称" + Q(opsLog.WorkItemName) +
				"，负责人" + Q(strings.Join(names, "，")) +
				"，排期" + Q(T(opsLog.PlanStartAt, opsLog.PlanCompleteAt)) +
				"，进度" + Q(cast.ToString(opsLog.ProcessRate)+"%")
		} else {
			result.ModuleId = int(opsLog.WorkItemId)
			result.ModuleTitle = opsLog.WorkItemName
			result.OperMsg = fmt.Sprintf("删除了 任务「%v」 ", opsLog.WorkItemName)
			result.OperMsg = "删除了 " + W("任务") + Q(opsLog.WorkItemName)
		}

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemTagChangeHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeWorkItemTag)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleTitle:  opsLog.WorkItemName,
			ModuleId:     int(opsLog.WorkItemId),
		}

		allTags := append(opsLog.OldTags, opsLog.NewTags...)
		tagMap, _ := s.tagRepo.TagMap(ctx, utils.StringArrToInt64Arr(allTags))

		oldTagNames := stream.Map(utils.ToInt64Array(opsLog.OldTags), func(v int64) string {
			if tag, ok := tagMap[v]; ok {
				return tag.TagName
			}
			return ""
		})

		newTagNames := stream.Map(utils.ToInt64Array(opsLog.NewTags), func(v int64) string {
			if tag, ok := tagMap[v]; ok {
				return tag.TagName
			}
			return ""
		})

		result.OperMsg = "将 " + W("关联标签") + " 由" + Q(strings.Join(oldTagNames, "，")) + "更新为" + Q(strings.Join(newTagNames, "，"))

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemFlowNodeModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifyWorkItemFlowNode)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		witem, _ := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if witem == nil {
			continue
		}

		tpltConf, _ := s.wFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, witem.WorkFlowTemplateId)
		if tpltConf == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(witem.Id),
			ModuleTitle:  witem.WorkItemName,
		}

		var nodeName string
		switch {
		case witem.IsWorkFlowMainTask():
			nodeConf := tpltConf.WorkFlowConf().GetNode(opsLog.FlowNodeCode)
			if nodeConf != nil {
				nodeName = nodeConf.Name
			}
		case witem.IsStateFlowMainTask():
			nodeConf := tpltConf.StateFlowConf().GetNode(opsLog.FlowNodeCode)
			if nodeConf != nil {
				nodeName = nodeConf.Name
			}
		}

		var fileUpdates = make([]string, 0)
		var oldValue, newValue, filedName string
		for _, v := range opsLog.Updates {
			switch v.Field {

			case "planTime":

				oldTimes := v.OldValue.([]any)
				newTimes := v.NewValue.([]any)

				if len(oldTimes) != 2 || len(newTimes) != 2 {
					continue
				}

				filedName = "排期"
				oldValue = T(cast.ToInt64(oldTimes[0]), cast.ToInt64(oldTimes[1]))
				newValue = T(cast.ToInt64(newTimes[0]), cast.ToInt64(newTimes[1]))
			}

			fileUpdates = append(fileUpdates, "将"+Q(nodeName)+W(filedName)+" 由"+Q(oldValue)+"更新为"+Q(newValue))
		}

		result.OperMsg = strings.Join(fileUpdates, " ")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemFileChangeHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeWorkItemFile)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleTitle:  opsLog.WorkItemName,
			ModuleId:     int(opsLog.WorkItemId),
		}

		var msg []string
		if len(opsLog.AddFiles) > 0 {
			infos := stream.Map(opsLog.AddFiles, func(v domain_message.FileInfo) string {
				return Q(v.String())
			})

			msg = append(msg, "添加  "+W("资源附件")+strings.Join(infos, "，"))
		}

		if len(opsLog.RemoveFiles) > 0 {
			infos := stream.Map(opsLog.RemoveFiles, func(v domain_message.FileInfo) string {
				return Q(v.String())
			})

			msg = append(msg, "删除 "+W("资源附件")+strings.Join(infos, "，"))
		}

		result.OperMsg = strings.Join(msg, " ")

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemFlowNodeConfirmHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ConfirmWorkItemFlowNode)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		statusItem, err := s.wStatusRepo.GetWorkItemStatusItem(ctx, witem.WorkItemStatus.Id)
		if err != nil {
			continue
		}

		var flowNodeName string
		if opsLog.FlowNodeId != 0 {
			flowTplt, err := s.wFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, witem.WorkFlowTemplateId)
			if err == nil && flowTplt.WorkFLowConfig != nil {
				nodeConf := flowTplt.WorkFLowConfig.GetNode(opsLog.FlowNodeCode)
				if nodeConf != nil {
					flowNodeName = nodeConf.Name
				}

			}
		}

		if flowNodeName == "" {
			return
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ShowType:     oper.ShowTypeWorkItemNodeChange,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
		}

		result.OperMsg = "完成" + Q(flowNodeName) + W("节点") + "，" + "将 " + W("任务状态") + " 更新为" + Q(statusItem.Name)

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemFlowNodeRollbackHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.RollbackWorkItemFlowNode)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		var flowNodeName string
		if opsLog.FlowNodeId != 0 {
			flowTplt, err := s.wFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, witem.WorkFlowTemplateId)
			if err == nil && flowTplt.WorkFLowConfig != nil {
				nodeConf := flowTplt.WorkFLowConfig.GetNode(opsLog.FlowNodeCode)
				if nodeConf != nil {
					flowNodeName = nodeConf.Name
				}
			}
		}

		if flowNodeName == "" {
			return
		}

		statusItem, err := s.wStatusRepo.GetWorkItemStatusItem(ctx, witem.WorkItemStatus.Id)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ShowType:     oper.ShowTypeWorkItemNodeChange,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
		}

		result.OperMsg = "将任务 " + W("回滚") + " 至" + Q(flowNodeName) + W("节点") + "，原因" + Q(opsLog.Reason) + "。" + "将 " + W("任务状态") + " 更新为" + Q(statusItem.Name)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemFlowUpgradeHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.UpgradeTaskWorkFlow)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, opsLog.SpaceId)
		if err != nil {
			continue
		}

		roles, err := s.wRoleRepo.GetWorkItemRoles(ctx, space.Id)
		if err != nil {
			continue
		}

		roleMap := stream.ToMap(roles, func(idx int, item *work_item_role.WorkItemRole) (string, *work_item_role.WorkItemRole) {
			return cast.ToString(item.Id), item
		})

		allUid := make([]int64, 0)
		for _, v := range opsLog.NewRoles {
			allUid = append(allUid, utils.StringArrToInt64Arr(v.Directors)...)
		}
		for _, v := range opsLog.OldRoles {
			allUid = append(allUid, utils.StringArrToInt64Arr(v.Directors)...)
		}

		userMap, _ := s.userRepo.UserMap(ctx, allUid)

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ShowType:     oper.ShowTypeWorkItemNodeChange,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  opsLog.WorkItemName,
		}

		msg := "将 " + W("任务流程升级") + "，当前 " + W(opsLog.WorkFlowName) + " 由" + Q("v"+cast.ToString(opsLog.OldVersion)) + "升级至" + Q("v"+cast.ToString(opsLog.NewVersion)) + "。"

		rolesMsg := make([]string, 0)

		oldRoleDirectorsMap := stream.ToMap(opsLog.OldRoles, func(idx int, item domain_message.RoleDirector) (string, []string) {
			return item.RoleId, item.Directors
		})

		for _, v := range opsLog.NewRoles {
			if roleMap[v.RoleId] == nil {
				continue
			}

			newNames := stream.Map(utils.StringArrToInt64Arr(v.Directors), func(item int64) string {
				return userMap[item].UserNickname
			})

			oldNames := stream.Map(utils.StringArrToInt64Arr(oldRoleDirectorsMap[v.RoleId]), func(item int64) string {
				return userMap[item].UserNickname
			})

			msg := "将 " + W(fmt.Sprintf("“%v”负责人", roleMap[v.RoleId].Name)) + " 由" + Q(strings.Join(oldNames, "，")) + "变更为" + Q(strings.Join(newNames, "，"))

			rolesMsg = append(rolesMsg, msg)
		}

		result.OperMsg = msg + strings.Join(rolesMsg, "，")

		s.invokeOperLog(ctx, operLogger, result)
	}
}
