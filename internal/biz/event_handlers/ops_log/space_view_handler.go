package ops_log

import (
	"context"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) spaceViewCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateSpaceView)
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

		var moduleType = oper.ModuleTypeSpaceUserView
		var viewTypeName = "个人视图"
		if opsLog.ViewType == consts.SpaceViewType_Public {
			moduleType = oper.ModuleTypeSpaceGlobalView
			viewTypeName = "公共视图"
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeAdd,
			ModuleType:   moduleType,
			ModuleId:     int(opsLog.ViewId),
			ModuleTitle:  opsLog.ViewName,
		}

		result.OperMsg = "另存 " + W(viewTypeName) + Q(opsLog.ViewName)

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceViewDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteSpaceView)
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

		var moduleType = oper.ModuleTypeSpaceUserView
		var viewTypeName = "个人视图"
		if opsLog.ViewType == consts.SpaceViewType_Public {
			moduleType = oper.ModuleTypeSpaceGlobalView
			viewTypeName = "公共视图"
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeDel,
			ModuleType:   moduleType,
			ModuleId:     int(opsLog.ViewId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "删除 " + W(viewTypeName) + Q(opsLog.ViewName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceViewSetNameHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetSpaceViewName)
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

		var moduleType = oper.ModuleTypeSpaceUserView
		if opsLog.ViewType == consts.SpaceViewType_Public {
			moduleType = oper.ModuleTypeSpaceGlobalView
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   moduleType,
			ModuleId:     int(opsLog.ViewId),
			ModuleTitle:  opsLog.ViewNewName,
		}

		s.invokeOperLog(ctx, operLogger, result)
		result.OperMsg = "将 " + W("视图名称") + " 由" + Q(opsLog.ViewOldName) + "更新为" + Q(opsLog.ViewNewName)
	}
}

func (s *OpsLogEventHandlers) spaceViewSetStatusHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetSpaceViewStatus)
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
			ModuleId:     int(opsLog.ViewId),
			ModuleTitle:  opsLog.ViewName,
		}

		action := "开启 "
		if opsLog.Status == 0 {
			action = "关闭 "
		}

		s.invokeOperLog(ctx, operLogger, result)
		result.OperMsg = action + W("查看视图") + Q(opsLog.ViewName)
	}
}

func (s *OpsLogEventHandlers) updateSpaceViewHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.UpdateSpaceView)
		if opsLog == nil || opsLog.Field == "name" {
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

		var moduleType = oper.ModuleTypeSpaceUserView
		if opsLog.ViewType == consts.SpaceViewType_Public {
			moduleType = oper.ModuleTypeSpaceGlobalView
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   moduleType,
			ModuleId:     int(opsLog.ViewId),
			ModuleTitle:  opsLog.ViewName,
		}

		s.invokeOperLog(ctx, operLogger, result)
		result.OperMsg = "更新了" + Q(opsLog.ViewName) + W("视图配置")
	}
}
