package ops_log

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) workVersionCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateWorkVersion)
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
			ModuleType:   oper.ModuleTypeSpaceWorkVersion,
			ModuleId:     int(opsLog.WorkVersionId),
			ModuleTitle:  opsLog.WorkVersionName,
		}

		result.OperMsg = "添加 " + W("版本") + Q(opsLog.WorkVersionName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workVersionDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteWorkVersion)
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

		result.OperMsg = "移除 " + W("版本") + Q(opsLog.WorkVersionName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workVersionModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifyWorkVersion)
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
			ModuleType:   oper.ModuleTypeSpaceWorkVersion,
			ModuleId:     int(opsLog.WorkVersionId),
			ModuleTitle:  opsLog.WorkVersionName,
		}

		var fileUpdates = make([]string, 0)
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

			fileUpdates = append(fileUpdates, fmt.Sprintf("%v: %v->%v", filedName, oldValue, newValue))

			result.OperMsg = "将 " + W("版本") + W(filedName) + " 由" + Q(oldValue) + "更新为" + Q(newValue)
			s.invokeOperLog(ctx, operLogger, result)
		}
	}
}
