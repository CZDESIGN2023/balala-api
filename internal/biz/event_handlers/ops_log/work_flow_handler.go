package ops_log

import (
	"context"
	"github.com/spf13/cast"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/domain/work_flow"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) workFlowCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateWorkFlow)
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
			ModuleType:   oper.ModuleTypeWorkFlow,
			ModuleId:     int(opsLog.WorkFlowId),
			ModuleTitle:  opsLog.WorkFlowName,
		}

		var flowType string
		switch opsLog.FlowMode {
		case consts.FlowMode_StateFlow:
			flowType = "任务流程-状态模式"
		case consts.FlowMode_WorkFlow:
			flowType = "任务流程-节点模式"
		}

		result.OperMsg = "新建 " + W(flowType) + Q(opsLog.WorkFlowName)
		if opsLog.SrcSpaceName != "" {
			result.OperMsg = "通过复制项目" + Q(opsLog.SrcSpaceName) + "中流程" + Q(opsLog.SrcWorkFlowName) + "，创建了 " + W(flowType) + Q(opsLog.WorkFlowName)
		}
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workFlowDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteWorkFlow)
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
			ModuleType:   oper.ModuleTypeWorkFlow,
			ModuleId:     int(opsLog.WorkFlowId),
			ModuleTitle:  opsLog.WorkFlowName,
		}

		result.OperMsg = "删除 " + W("任务流程") + Q(opsLog.WorkFlowName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) workFlowModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifyWorkFlow)
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
			ModuleType:   oper.ModuleTypeWorkFlow,
			ModuleId:     int(opsLog.WorkFlowId),
			ModuleTitle:  opsLog.WorkFlowName,
		}

		for _, v := range opsLog.Updates {
			result := &(*result)
			msg := ""

			switch v.Field {
			case "name":
				filedName = "名称"
				oldValue = cast.ToString(v.OldValue)
				newValue = cast.ToString(v.NewValue)

				msg = "将 " + W("任务流程") + W(filedName) + "，由 " + Q(oldValue) + "更新为" + Q(newValue)

			case "status":
				var action string
				switch {
				case v.NewValue == work_flow.WorkFlowStatus_Enable:
					action = "启用"
				case v.OldValue == work_flow.WorkFlowStatus_Hide && v.NewValue == work_flow.WorkFlowStatus_Disable:
					action = "显示"
				case v.OldValue == work_flow.WorkFlowStatus_Disable && v.NewValue == work_flow.WorkFlowStatus_Hide:
					action = "隐藏"
				case v.NewValue == work_flow.WorkFlowStatus_Disable:
					action = "禁用"
				}

				msg = action + " " + W("任务流程") + Q(opsLog.WorkFlowName)
			}

			result.OperMsg = msg
			s.invokeOperLog(ctx, operLogger, result)
		}

	}
}

func (s *OpsLogEventHandlers) workFlowTemplateSaveHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SaveWorkFlowTemplate)
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
			ModuleType:   oper.ModuleTypeWorkFlow,
			ModuleId:     int(opsLog.WorkFlowId),
			ModuleTitle:  opsLog.WorkFlowName,
		}

		result.OperMsg = "修改 " + W("任务流程") + Q(opsLog.WorkFlowName) + "，" + W("版本") + "更新为" + Q("v"+cast.ToString(opsLog.TemplateVersion))
		s.invokeOperLog(ctx, operLogger, result)
	}
}
