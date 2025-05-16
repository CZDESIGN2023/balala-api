package ops_log

import (
	"context"
	"github.com/spf13/cast"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) workItemStatusCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateWorkItemStatus)
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
			ModuleType:   oper.ModuleTypeWorkItemStatus,
			ModuleId:     int(opsLog.WorkItemStatusId),
			ModuleTitle:  opsLog.WorkItemStatusName,
		}

		var typeStr string
		switch opsLog.FlowScope {
		case consts.FlowScope_Stateflow:
			typeStr = "任务状态-状态模式"
		case consts.FlowScope_Workflow:
			typeStr = "任务状态-节点模式"
		}

		result.OperMsg = "新建 " + W(typeStr) + Q(opsLog.WorkItemStatusName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemStatusDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteWorkItemStatus)
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
			ModuleType:   oper.ModuleTypeWorkItemStatus,
			ModuleId:     int(opsLog.WorkItemStatusId),
			ModuleTitle:  opsLog.WorkItemStatusName,
		}

		var typeStr string
		switch opsLog.FlowScope {
		case consts.FlowScope_Stateflow:
			typeStr = "任务状态-状态模式"
		case consts.FlowScope_Workflow:
			typeStr = "任务状态-节点模式"
		}

		result.OperMsg = "删除 " + W(typeStr) + Q(opsLog.WorkItemStatusName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workItemStatusModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifyWorkItemStatus)
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

		var oldValue, newValue, filedName string
		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeWorkItemStatus,
			ModuleId:     int(opsLog.WorkItemStatusId),
			ModuleTitle:  opsLog.WorkItemStatusName,
		}

		for _, v := range opsLog.Updates {
			result := &(*result)
			switch v.Field {
			case "name":
				filedName = "名称"
				oldValue = cast.ToString(v.OldValue)
				newValue = cast.ToString(v.NewValue)

			//case "ranking":
			//	filedName = "排序"
			//	oldValue = cast.ToString(v.OldValue)
			//	newValue = cast.ToString(v.NewValue)

			case "status":
				filedName = "状态"
				statusName := func(status int) string {
					if status == 1 {
						return "启用"
					} else {
						return "禁用"
					}
				}

				oldValue = statusName(cast.ToInt(v.OldValue))
				newValue = statusName(cast.ToInt(v.NewValue))
			}

			result.OperMsg = "将 " + W("任务状态") + W(filedName) + " 由" + Q(oldValue) + "更新为" + Q(newValue)
			s.invokeOperLog(ctx, operLogger, result)
		}

	}
}
