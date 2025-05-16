package ops_log

import (
	"context"
	"github.com/spf13/cast"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) workObjectCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateWorkObject)
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
			ModuleType:   oper.ModuleTypeSpaceWorkObject,
			ModuleId:     int(opsLog.WorkObjectId),
			ModuleTitle:  opsLog.WorkObjectName,
		}

		result.OperMsg = "添加 " + W("模块") + Q(opsLog.WorkObjectName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workObjectDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteWorkObject)
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
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
			ModuleTitle:  space.SpaceName,
		}

		result.OperMsg = "移除 " + W("模块") + Q(opsLog.WorkObjectName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workObjectModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifyWorkObject)
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
			ModuleType:   oper.ModuleTypeSpaceWorkObject,
			ModuleId:     int(opsLog.WorkObjectId),
			ModuleTitle:  opsLog.WorkObjectName,
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

			result.OperMsg = "将 " + W("模块") + W(filedName) + " 由" + Q(oldValue) + "更新为" + Q(newValue)
			s.invokeOperLog(ctx, operLogger, result)
		}
	}
}
