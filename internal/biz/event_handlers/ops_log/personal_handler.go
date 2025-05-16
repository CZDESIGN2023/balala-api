package ops_log

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go-cs/api/comm"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/oper"
	"go-cs/pkg/stream"
	"slices"
	"strings"
)

func (s *OpsLogEventHandlers) personalChangeAvatarHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalChangeAvatar)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(opsLog.UserId),
			ModuleTitle:  "修改个人头像",
		}

		result.OperMsg = "更新了 " + W("头像")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalChangeNicknameHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalChangeNickName)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(opsLog.UserId),
			ModuleTitle:  "修改个人昵称",
		}

		result.OperMsg = "将 " + W("昵称") + " 由" + Q(opsLog.OldNickName) + " 更新为" + Q(opsLog.NewNickName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalResetPwdHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalResetPwd)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(opsLog.UserId),
			ModuleTitle:  "修改个人密码",
		}

		result.OperMsg = W("修改密码")
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalSetSpaceNotifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {
	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalSetSpaceNotify)
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
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleTitle:  "修改项目通知",
		}

		action := "开启"
		if opsLog.Notify == 0 {
			action = "关闭"
		}

		result.OperMsg = W(action+"接收") + Q(space.SpaceName) + "项目通知"
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalBindThirdPlatformHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalBindThirdPlatform)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(operLogger.Operator.OperUid),
			ModuleTitle:  "绑定第三方",
		}

		result.OperMsg = W("绑定") + " 三方平台" + Q(opsLog.PlatformName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalUnBindThirdPlatformHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalUnBindThirdPlatform)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(operLogger.Operator.OperUid),
			ModuleTitle:  "解绑第三方",
		}

		result.OperMsg = W("解绑") + " 三方平台" + Q(opsLog.PlatformName)
		s.invokeOperLog(ctx, operLogger, result)
	}
}

type syncTableColumn struct {
	Version string `json:"version"`
	Data    struct {
		Source []string   `json:"source"`
		Target [][]string `json:"target"`
	}
}

func (s *OpsLogEventHandlers) personalSetTempConfigHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalSetTempConfig)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		var msg string

		for k, v := range opsLog.NewValues {
			switch k {
			case "autoRefreshTaskFlow":
				action := "关闭"
				if v == "true" {
					action = "开启"
				}
				msg = action + " " + W("任务列表自动刷新")
				break
			case "syncTableColumn":
				conf := syncTableColumn{}
				json.Unmarshal([]byte(v), &conf)
				obj := conf.Data

				if len(obj.Source) == 0 {
					msg = "将 " + W("跨空间表格列字段同步设置 取消同步")
				} else {
					sourceSpaceId := obj.Source[0]
					sourceViewId := obj.Source[1]

					targetSpaceIds := stream.Map(obj.Target, func(e []string) string {
						return e[0]
					})

					var allSpaceIds = []string{sourceSpaceId}
					var allViewIds = []string{sourceViewId}

					var spaceId2ViewId = map[string][]string{}
					var targetRankingMap = map[string]int{}

					if slices.Contains(targetSpaceIds, "all") {
						allSpaceIds = []string{sourceSpaceId}
						allViewIds = []string{sourceViewId}

						spaceId2ViewId["all"] = nil
					} else {
						for i, v := range obj.Target {
							spaceId2ViewId[v[0]] = append(spaceId2ViewId[v[0]], v[1])

							if _, ok := targetRankingMap[v[0]]; !ok {
								targetRankingMap[v[0]] = i
							}
						}

						for spaceId, viewIds := range spaceId2ViewId {
							allSpaceIds = append(allSpaceIds, spaceId)

							validViewIds := stream.Remove(viewIds, "all")

							allViewIds = append(allViewIds, validViewIds...)
							spaceId2ViewId[spaceId] = validViewIds

							//if slices.Contains(viewIds, "all") {
							//	spaceId2ViewId[spaceId] = []string{"all"}
							//} else {
							//	allViewIds = append(allViewIds, viewIds...)
							//}
						}
					}

					space, _ := s.spaceRepo.SpaceMap(ctx, utils.ToInt64Array(allSpaceIds))
					views, _ := s.viewRepo.GetUserViewMapByIds(ctx, utils.StringArrToInt64Arr(allViewIds))

					f := func(spaceId string, e []string) string {
						if spaceId == "all" {
							return "全部项目空间视图"
						}

						spaceName := space[cast.ToInt64(spaceId)].SpaceName
						if slices.Contains(e, "all") {
							return fmt.Sprintf("%s / 项目所有视图", spaceName)
						}

						viewNames := stream.Map(e, func(e string) string {
							return views[cast.ToInt64(e)].Name
						})

						return fmt.Sprintf("%s / [%s]", spaceName, strings.Join(viewNames, "、"))
					}

					targetEntries := stream.ToEntries(spaceId2ViewId)

					slices.SortedFunc(slices.Values(targetEntries), func(a, b stream.Entry[string, []string]) int {
						return targetRankingMap[a.Key] - targetRankingMap[b.Key]
					})

					sourceLog := f(sourceSpaceId, []string{sourceViewId})
					targetLogs := stream.Map(targetEntries, func(e stream.Entry[string, []string]) string {
						return f(e.Key, e.Val)
					})

					msg = "将项目" + Q(sourceLog) + "的表格列字段设置，自动同步至" + Q(strings.Join(targetLogs, "、"))

					break
				}

			default:
				continue
			}
		}

		if msg == "" {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(operLogger.Operator.OperUid),
			ModuleTitle:  "",
		}

		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalSetThirdPlatformNotifyHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalSetThirdPlatformNotify)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		platformName := consts.GetPlatformName(comm.ThirdPlatformCode(opsLog.PlatformCode))

		action := "开启"
		if opsLog.Notify == 0 {
			action = "关闭"
		}

		msg := W(action+"接收") + Q(platformName) + "通知"

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(operLogger.Operator.OperUid),
			ModuleTitle:  "第三方通知",
		}
		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) personalSetUserConfigHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.PersonalSetUserConfig)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		var msg string

		keyName := consts.GetUserConfigKeyName(opsLog.Key)
		switch opsLog.Key {
		case consts.UserConfigKey_NotifySwitchGlobal, consts.UserConfigKey_NotifySwitchThirdPlatform, consts.UserConfigKey_NotifySwitchSpace:
			msg = W("开启接收") + keyName
			if opsLog.NewValue == "0" {
				msg = W("关闭接收") + keyName
			}
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      0,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeUser,
			ModuleId:     int(operLogger.Operator.OperUid),
			ModuleTitle:  "通知",
		}
		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}
