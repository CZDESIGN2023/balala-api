package ops_log

import (
	"context"
	"fmt"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) AdminAddUserHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminAddUser)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		user, err := s.userRepo.GetUserByUserId(ctx, opsLog.UserId)
		if err != nil {
			return
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeUser,
			BusinessType: oper.BusinessTypeAdd,
			ModuleTitle:  "添加用户",
			ModuleId:     int(user.Id),
		}

		result.OperMsg = "添加 " + W("账号") + Q(U(user))
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminCancelUserHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminCancelUser)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeUser,
			BusinessType: oper.BusinessTypeDel,
			ModuleTitle:  "注销用户",
			ModuleId:     int(opsLog.UserId),
		}

		result.OperMsg = W("注销用户") + Q(fmt.Sprintf("%s (%s)", opsLog.Nickname, opsLog.Username))
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeUserNicknameHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeUserNickname)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		user, err := s.userRepo.GetUserByUserId(ctx, opsLog.UserId)
		if err != nil {
			return
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeUser,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  "修改用户昵称",
			ModuleId:     int(opsLog.UserId),
		}

		user.UserNickname = opsLog.OldValue

		result.OperMsg = "将" + Q(U(user)) + W("昵称") + " 更新为" + Q(opsLog.NewValue)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminResetUserPasswordHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminResetUserPassword)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		user, err := s.userRepo.GetUserByUserId(ctx, opsLog.UserId)
		if err != nil {
			return
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeUser,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  "重置用户密码",
			ModuleId:     int(opsLog.UserId),
		}

		result.OperMsg = W("重置") + Q(U(user)) + W("登录密码")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) AdminChangeUserRoleHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AdminChangeUserRole)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		user, err := s.userRepo.GetUserByUserId(ctx, opsLog.UserId)
		if err != nil {
			return
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			ModuleType:   oper.ModuleTypeUser,
			BusinessType: oper.BusinessTypeModify,
			ModuleId:     int(opsLog.UserId),
			ModuleTitle:  "修改用户角色",
		}

		var roleName = opsLog.NewValue.String()

		result.OperMsg = "将" + Q(U(user)) + "设为 " + W(roleName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}
