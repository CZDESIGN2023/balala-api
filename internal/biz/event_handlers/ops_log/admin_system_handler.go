package ops_log

import (
	"context"
	"github.com/tidwall/gjson"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) AdminChangeSystemLogoHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeSystemLogo)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeSystem,
			BusinessType: oper.BusinessTypeModify,
		}

		result.OperMsg = "更新了 " + W("网站LOGO")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeSystemTitleHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeSystemTitle)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeSystem,
			BusinessType: oper.BusinessTypeModify,
		}

		result.OperMsg = "将 " + W("网站标题") + " 由" + Q(opsLog.OldValue) + "更新为" + Q(opsLog.NewValue)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeSystemAccessUrlHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeSystemAccessUrl)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeSystem,
			BusinessType: oper.BusinessTypeModify,
		}

		result.OperMsg = "将 " + W("网站访问地址") + " 由" + Q(opsLog.OldValue) + "更新为" + Q(opsLog.NewValue)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeSystemLoginBgHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeSystemLoginBg)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeSystem,
			BusinessType: oper.BusinessTypeModify,
		}

		result.OperMsg = "更新了 " + W("登录背景")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeSystemRegisterEntryHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeSystemRegisterEntry)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeSystem,
			BusinessType: oper.BusinessTypeModify,
		}

		var action = "关闭"
		if opsLog.NewValue == "1" {
			action = "开启"
		}

		result.OperMsg = "将 " + W("注册功能") + " " + W(action)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeSystemAttachSizeHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeSystemAttachSize)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeSystem,
			BusinessType: oper.BusinessTypeModify,
		}

		// {"value":"5","unit":"GB"}
		oldVal := gjson.Get(opsLog.OldValue, "value").Str + gjson.Get(opsLog.OldValue, "unit").Str
		newVal := gjson.Get(opsLog.NewValue, "value").Str + gjson.Get(opsLog.NewValue, "unit").Str

		//oldVal := parseAttach(gjson.Get(opsLog.OldValue, "value").Int(), gjson.Get(opsLog.OldValue, "unit").Str)
		//newVal := parseAttach(gjson.Get(opsLog.NewValue, "value").Int(), gjson.Get(opsLog.NewValue, "unit").Str)

		result.OperMsg = "将 " + W("附件大小") + " 由" + Q(oldVal) + "更新为" + Q(newVal)
		s.invokeOperLog(ctx, operLogger, result)
	}
}
