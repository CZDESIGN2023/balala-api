package ops_log

import (
	"context"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) changeWorkObjectOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeWorkObjectOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "变更了 " + W("模块") + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)

	}
}

func (s *OpsLogEventHandlers) changeVersionOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeVersionOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "变更了 " + W("版本") + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)

	}
}

func (s *OpsLogEventHandlers) changeWorkFlowOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeWorkFlowOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "变更了 " + W("任务流程") + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)

	}
}

func (s *OpsLogEventHandlers) changeRoleOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeRoleOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		var typeStr string
		switch opsLog.FlowScope {
		case consts.FlowScope_Stateflow:
			typeStr = "流程角色-状态模式"
		case consts.FlowScope_Workflow:
			typeStr = "流程角色-节点模式"
		}
		result.OperMsg = "变更了 " + W(typeStr) + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) changeStatusOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeStatusOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		var typeStr string
		switch opsLog.FlowScope {
		case consts.FlowScope_Stateflow:
			typeStr = "任务状态-状态模式"
		case consts.FlowScope_Workflow:
			typeStr = "任务状态-节点模式"
		}

		result.OperMsg = "变更了 " + W(typeStr) + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)

	}
}

func (s *OpsLogEventHandlers) changeOverviewDataItemOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeOverviewDataItemOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "变更了 " + W("概览数据项") + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)

	}
}

func (s *OpsLogEventHandlers) changeOverviewBlockOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeOverviewBlockOrder)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "变更了 " + W("概览信息卡") + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) changeViewOrderHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ChangeSpaceViewOrder)
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
			ModuleType:   oper.ModuleTypeSpaceUserView,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "变更了 " + W("视图") + " 的排序"
		s.invokeOperLog(ctx, operLogger, result)
	}
}
