package ops_log

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/domain/search/human"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/oper"
	"go-cs/pkg/stream"
	"strings"
)

func (s *OpsLogEventHandlers) spaceQuitLogHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.QuitSpace)
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

		userMap, err := s.userRepo.UserMap(ctx, []int64{opsLog.MemberUid, opsLog.TransferTargetId})
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		result.OperMsg = W("离开项目")
		if opsLog.TransferTargetId != 0 {
			result.OperMsg = W("离开项目") + "，将 " + W(fmt.Sprintf("%v 条任务单", opsLog.WorkItemNum)) + " 转移给" + Q(userMap[opsLog.TransferTargetId].UserNickname)
		}
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceAddMemberHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AddSpaceMember)
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

		userMap, err := s.userRepo.UserMap(ctx, []int64{opsLog.MemberUid})
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(space.Id),
		}

		result.OperMsg = "添加 " + W("成员") + Q(userMap[opsLog.MemberUid].UserNickname) + "，分配权限" + Q(consts.GetSpaceRoleName(int(opsLog.RoleId)))
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceRemoveMemberHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.RemoveSpaceMember)
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

		userMap, err := s.userRepo.UserMap(ctx, []int64{opsLog.MemberUid, opsLog.TransferTargetId})
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeDel,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(space.Id),
		}

		result.OperMsg = "移除 " + W("成员") + Q(userMap[opsLog.MemberUid].UserNickname)
		if opsLog.TransferTargetId != 0 {
			result.OperMsg = "移除 " + W("成员") + Q(userMap[opsLog.MemberUid].UserNickname) + "，将 " + W(fmt.Sprintf("%v 条任务单", opsLog.WorkItemNum)) + " 转移给" + Q(userMap[opsLog.TransferTargetId].UserNickname)
		}
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetMemberRoleHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetSpaceMemberRole)
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

		user, err := s.userRepo.GetUserByUserId(ctx, opsLog.MemeberUid)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(space.Id),
		}

		result.OperMsg = "将" + Q(user.UserNickname) + "的权限" + "，由" + Q(consts.GetSpaceRoleName(int(opsLog.OldRole))) + "更新为" + Q(consts.GetSpaceRoleName(int(opsLog.NewRole)))
		if opsLog.OldRole == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
			result.OperMsg = "将" + Q(user.UserNickname) + "撤销 " + W(consts.GetSpaceRoleName(int(opsLog.OldRole)))
		}

		if opsLog.NewRole == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
			result.OperMsg = "将" + Q(user.UserNickname) + "设为 " + W(consts.GetSpaceRoleName(int(opsLog.NewRole)))
		}

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) TransferSpaceHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.TransferSpace)
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

		target, err := s.userRepo.GetUserByUserId(ctx, opsLog.TargetUserId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(space.Id),
		}

		result.OperMsg = "将 " + W("项目创建者") + " 转移给" + Q(target.UserNickname) + "，" + Q(operUser.OperUname) + "权限更新为 " + Q(consts.GetSpaceRoleName(consts.MEMBER_ROLE_MANAGER))
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceCreateHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateSpace)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      opsLog.SpaceId,
			SpaceName:    opsLog.SpaceName,
			BusinessType: oper.BusinessTypeAdd,
			ModuleTitle:  opsLog.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		result.OperMsg = "创建 " + W("项目") + Q(opsLog.SpaceName)
		if opsLog.SrcSpaceId != 0 {
			result.OperMsg = "通过复制项目" + Q(opsLog.SrcSpaceName) + "创建了 " + W("项目") + Q(opsLog.SpaceName)
		}

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceModifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.ModifySpace)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      opsLog.SpaceId,
			SpaceName:    opsLog.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  opsLog.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		var oldValue, newValue, filedName string
		for _, v := range opsLog.Updates {
			if v.Field == "spaceName" {
				filedName = "项目名称"
				oldValue = cast.ToString(v.OldValue)
				newValue = cast.ToString(v.NewValue)
			}

			if v.Field == "describe" {
				filedName = "项目描述"
				oldValue = utils.ClearRichTextToPlanText(cast.ToString(v.OldValue), true)
				newValue = utils.ClearRichTextToPlanText(cast.ToString(v.NewValue), true)
			}

			result.OperMsg = "将 " + W(filedName) + " 由" + Q(oldValue) + "更新为" + Q(newValue)
			s.invokeOperLog(ctx, operLogger, result)
		}
	}
}

func (s *OpsLogEventHandlers) spaceDeleteHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DelSpace)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      opsLog.SpaceId,
			SpaceName:    opsLog.SpaceName,
			BusinessType: oper.BusinessTypeDel,
			ModuleTitle:  opsLog.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		result.OperMsg = "删除 " + W("项目") + Q(opsLog.SpaceName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetCommentDeletableHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetCommentDeletable)
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

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		action := "开启"
		if opsLog.Deletable == 0 {
			action = "关闭"
		}

		result.OperMsg = action + " " + W("评论可删除")

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetCommentDeletableWhenArchivedHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetCommentDeletableWhenArchived)
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

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		action := "开启"
		if opsLog.Value == 0 {
			action = "关闭"
		}

		result.OperMsg = action + " " + W("评论归档可删除")

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetCommentShowPosHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetCommentShowPos)
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

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		action := "开启"
		if opsLog.Value == 0 {
			action = "关闭"
		}

		result.OperMsg = action + " " + W("评论独立显示")

		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetWorkingDayHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetWorkingDay)
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

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		result.OperMsg = "将 " + W("工作日周期") + " 由" + Q(utils.ParseWorkingDay(opsLog.OldWeekDays)) + "更新为" + Q(utils.ParseWorkingDay(opsLog.WeekDays))
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetSpaceNotifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetSpaceNotify)
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

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleTitle:  space.SpaceName,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleId:     int(opsLog.SpaceId),
		}

		action := "开启"
		if opsLog.Notify == 0 {
			action = "关闭"
		}

		result.OperMsg = action + " " + W("消息推送")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) spaceSetTempConfigHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.SetSpaceTempConfig)
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

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpace,
			ModuleTitle:  space.SpaceName,
			ModuleId:     int(opsLog.SpaceId),
		}

		oldValues := opsLog.OldValues
		newValues := opsLog.NewValues

		keys := stream.Keys(newValues)

		for _, key := range keys {
			oldVal := oldValues[key]
			newVal := newValues[key]

			if oldVal == newVal {
				continue
			}

			var msg string

			switch key {
			case consts.SpaceTempConfigKey_overviewOrder:
				if oldVal != newVal {
					msg = "变更了 " + W("概览信息卡") + " 的排序"
				}
			case consts.SpaceTempConfigKey_spaceWorkbenchCountConditions:
				parser := human.NewConditionGroupParser(space.Id, s.wStatusRepo, s.wFlowRepo)

				itemList, orderChanged := human.GetUpdatedConditionGroup(oldVal, newVal)
				if orderChanged {
					msg = "变更了 " + W("概览数据项") + " 的排序"
				} else if len(itemList) > 0 {
					updates := human.ParseToChangeLog(parser, itemList[0])
					oldDesc := itemList[0].First.Desc

					var updateList []string
					for _, v := range updates {
						switch v.Field {
						case "desc":
							updateList = append(updateList, W("统计名称")+" 由"+Q(v.OldValue)+"更新为"+Q(v.NewValue))
						case "condition":
							updateList = append(updateList, W("汇总条件")+" 由"+Q(v.OldValue)+"更新为"+Q(v.NewValue))
						}
					}

					msg = "将 " + W("自定义指标") + Q(oldDesc) + "的 " + strings.Join(updateList, "，")
				} else {
					continue
				}
			}

			result.OperMsg = msg
			s.invokeOperLog(ctx, operLogger, result)
		}

	}
}
