package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"time"

	space_domain "go-cs/internal/domain/space"
	user_domain "go-cs/internal/domain/user"
)

func buildSpaceMsg(e notify.Event, operator *user_domain.User, space *space_domain.Space) *msg.Message {
	if operator == nil {
		operator = &user_domain.User{}
	}

	if space == nil {
		space = &space_domain.Space{}
	}

	return &msg.Message{
		Space: &msg.Space{
			SpaceId:   space.Id,
			SpaceName: space.SpaceName,
		},
		Relation: make([]msg.RelationType, 0),
		Type:     e,
		TypeDesc: e.String(),
		Notification: &msg.Notification{
			Action: msg.ActionType_edit,
			Subject: &msg.Subject{
				Type: msg.SubjectType_user,
				Data: &msg.UserData{
					Name:     operator.UserName,
					NickName: operator.UserName,
					Id:       operator.Id,
					Avatar:   operator.Avatar,
				},
			},
			Object: &msg.Object{
				Type: msg.ObjectType_space,
				Data: &msg.SpaceData{
					Id:   space.Id,
					Name: space.SpaceName,
				},
			},
			Date: time.Now(),
		},
	}
}

func (s *Notify) addMember(e *event.AddMember) {
	s.log.Infof("addMember: %+v", e)

	ctx := context.Background()
	space := e.Space
	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.TargetId})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]
	target := userMap[e.TargetId]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds). //管理员
		Concat(space.UserId). //加上空间创建人
		Unique(). //去重
		Diff(operator.Id). //排除操作人
		List()

	desc1 := func() string {
		roleName := consts.GetSpaceRoleName(int(e.RoleId))
		// [昵称(用户名)] 添加成员 [昵称(用户名)] ，分配权限为 “[权限名称]”
		return fmt.Sprintf("%v 添加成员 %v ，分配权限为 “%v", parseUserTmp(operator), parseUserTmp(target), roleName)
	}()

	desc2 := func() string {
		roleName := consts.GetSpaceRoleName(int(e.RoleId))
		//  [昵称(用户名)] 将你添加至项目中，当前权限 “[权限名称]”
		return fmt.Sprintf("%v 将你添加至项目中，当前权限 “%v”", parseUserTmp(operator), roleName)
	}()

	m := buildSpaceMsg(e.Event, operator, space)
	m.SetRedirectLink(s.makeRedirectLink(m))

	qlTarget := tea_im.NewRobotMessage().
		SetShowType(tea_im.ShowType_Text).
		SetTitle("项目提醒").
		SetTextRich(parseToImRich(MinorColorSpan(space.SpaceName+" ") + "\n" + desc2)).
		SetIcon(s.makeIconResLink("", IconRes_Project)).
		SetSVGIcon(IconSVGRes_Project).
		SetUrl(s.makeSpaceRedirectLink("", space.Id))

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(managerIds, []int64{e.TargetId}))

	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc1), managerIds...)
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc2).SetPopup(), e.TargetId)

		s.pushThirdPlatformMessage2(notifyCtx, qlTarget, []int64{e.TargetId})
	})

	s.CooperateAddMember(e)
}

func (s *Notify) removeMember(e *event.RemoveMember) {

	ctx := context.Background()

	s.log.Infof("removeMember: %+v", e)

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.TargetId})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]
	target := userMap[e.TargetId]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds). //管理员
		Concat(space.UserId). //加上空间创建人
		Unique(). //去重
		Diff(operator.Id). //排除操作人
		List()

	buildTemplate := func() string {
		//[昵称(用户名)] 移除成员 [昵称(用户名)]
		return fmt.Sprintf("%v 移除成员 %v", parseUserTmp(operator), parseUserTmp(target))
	}

	buildTemplate2 := func() string {
		// [昵称(用户名)] 已将你从项目中移除
		return fmt.Sprintf("%v 已将你从项目中移除", parseUserTmp(operator))
	}

	desc1 := buildTemplate()
	desc2 := buildTemplate2()

	m := buildSpaceMsg(e.Event, operator, space)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(managerIds, []int64{e.TargetId}))
	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc1), managerIds...)
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc2).SetPopup(), e.TargetId)
	})

	qlTarget := tea_im.NewRobotMessage().
		SetShowType(tea_im.ShowType_Text).
		SetTitle("项目提醒").
		SetTextRich(parseToImRich(MinorColorSpan(space.SpaceName+" ") + "\n" + desc2))
	qlTarget.SetIcon(s.makeIconResLink("", IconRes_Project)).
		SetSVGIcon(IconSVGRes_Project).
		SetUrl(s.makeSpaceRedirectLink("", space.Id))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlTarget, []int64{e.TargetId})
	})

	s.CooperateRemoveMember(e)
}

func (s *Notify) changeRole(e *event.ChangeRole) {
	s.log.Infof("changeRole: %+v", e)

	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.TargetId})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]
	target := userMap[e.TargetId]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds). //管理员
		Concat(space.UserId). //加上空间创建人
		Unique(). //去重
		Diff(operator.Id, e.TargetId). //排除操作人，被操作人
		List()

	buildTemplate := func() string {
		roleName := consts.GetSpaceRoleName(int(e.OldRoleId))
		newRoleName := consts.GetSpaceRoleName(int(e.NewRoleId))
		// [昵称(用户名)] 变更了 [昵称(用户名)] 权限：“[权限名称]” -> “[权限名称]”
		return fmt.Sprintf("%v 变更了 %v 权限：“%v” -> “%v”", parseUserTmp(operator), parseUserTmp(target), roleName, newRoleName)
	}

	buildTemplate2 := func() string {
		roleName := consts.GetSpaceRoleName(int(e.OldRoleId))
		newRoleName := consts.GetSpaceRoleName(int(e.NewRoleId))
		//  [昵称(用户名)] 变更了你的权限：“[权限名称]” -> “[权限名称]”
		return fmt.Sprintf("%v 变更了你的权限：“%v” -> “%v”", parseUserTmp(operator), roleName, newRoleName)
	}

	m := buildSpaceMsg(e.Event, operator, space)
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(managerIds, []int64{e.TargetId}))
	go func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(buildTemplate()), managerIds...)
		s.Send2(notifyCtx, m.Clone().SetDescribe(buildTemplate2()).SetPopup(), e.TargetId)
	}()

}

func (s *Notify) changeSpaceName(e *event.ChangeSpaceName) {
	s.log.Infof("changeSpaceName: %+v", e)

	if e.OldValue == e.NewValue {
		return
	}

	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	memberIds, _ := s.memberRepo.GetSpaceAllMemberIds(ctx, space.Id)
	memberIds = stream.Of(memberIds).Diff(operator.Id).List()

	buildTemplate := func() string {
		//   [昵称(用户名)]变更了项目名称： “[项目描述(变更前)]” -> “[项目描述]”
		return fmt.Sprintf("%v变更了项目名称： “%v” -> “%v”", parseUserTmp(operator), e.OldValue, e.NewValue)
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, memberIds)
	go func() {
		s.Send2(notifyCtx, m, memberIds...)
	}()
}

func (s *Notify) changeSpaceDescribe(e *event.ChangeSpaceDescribe) {
	s.log.Infof("ChangeSpaceDescribe: %+v", e)

	if e.OldValue == e.NewValue {
		return
	}

	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	memberIds, _ := s.memberRepo.GetSpaceAllMemberIds(ctx, space.Id)
	memberIds = stream.Of(memberIds).Diff(operator.Id).List()

	buildTemplate := func() string {
		oldValue := utils.ClearRichTextToPlanText(e.OldValue, false)
		newValue := utils.ClearRichTextToPlanText(e.NewValue, false)

		// [昵称(用户名)]变更了项目描述： “[项目描述(变更前)]” -> “[项目描述]”
		return fmt.Sprintf("%v 变更了项目描述： “%v” -> “%v”", parseUserTmp(operator), oldValue, newValue)
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, memberIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m, memberIds...)
	})

}

func (s *Notify) CreateSpace(e *event.CreateSpace) {
	s.log.Infof("CreateSpace: %+v", e)

	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	buildTemplate := func(roleId int) string {
		roleName := consts.GetSpaceRoleName(roleId)
		//  [昵称(用户名)] 将你添加至项目中，当前权限 “[权限名称]”
		return fmt.Sprintf("%v 将你添加至项目中，当前权限 “%v”", parseUserTmp(operator), roleName)
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetType(notify.Event_AddMember).
		SetPopup()
	m.SetRedirectLink(s.makeRedirectLink(m))

	utils.Go(func() {
		for _, member := range e.Members {
			// 跳过创建者
			if member.RoleId == consts.MEMBER_ROLE_SPACE_CREATOR {
				continue
			}

			// 前端通知
			desc := buildTemplate(int(member.RoleId))
			memberMsg := m.Clone().SetDescribe(desc)
			s.Send(memberMsg, member.UserId)

			// IM通知
			qlTarget := tea_im.NewRobotMessage().
				SetShowType(tea_im.ShowType_Text).
				SetTitle("项目提醒").
				SetTextRich(parseToImRich(MinorColorSpan(space.SpaceName+" ") + "\n" + desc)).
				SetIcon(s.makeIconResLink("", IconRes_Project)).
				SetSVGIcon(IconSVGRes_Project).
				SetUrl(s.makeSpaceRedirectLink("", space.Id))

			s.pushThirdPlatformMessage2(&notifyCtx{forceNotify: true}, qlTarget, []int64{member.UserId})

			// 协作通知
			s.CooperateAddMember(&event.AddMember{
				Operator: e.Operator,
				Space:    space,
				TargetId: member.UserId,
			})
		}
	})
}

func (s *Notify) deleteSpace(e *event.DeleteSpace) {
	s.log.Infof("deleteSpace: %+v", e)

	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	memberIds := stream.Diff(e.MemberIds, []int64{operator.Id})

	buildTemplate := func() string {
		// [昵称(用户名)]已删除项目：[项目名称]
		return fmt.Sprintf("%v 已删除项目：%v", parseUserTmp(operator), space.SpaceName)
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate()).
		SetPopup()

	notifyCtx := s.buildNotifyCtx(space, memberIds, WithNotifyCtxForceNotifyOpt())
	utils.Go(func() {
		s.Send2(notifyCtx, m, memberIds...)
	})

	for _, memberId := range memberIds {
		s.CooperateRemoveMember(&event.RemoveMember{
			Event:    notify.Event_RemoveMember,
			Space:    space,
			Operator: operator.Id,
			TargetId: memberId,
		})
	}
}

func (s *Notify) SetSpaceWorkingDay(e *event.SetSpaceWorkingDay) {
	s.log.Infof("SetSpaceWorkingDay: %+v", e)
	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	memberIds, _ := s.memberRepo.GetSpaceAllMemberIds(ctx, space.Id)
	memberIds = stream.Of(memberIds).Diff(operator.Id).List()

	buildTemplate := func() string {
		return fmt.Sprintf("%v在%v中设置了工作日： “%v”", parseUserTmp(operator), space.SpaceName, utils.ParseWorkingDay(e.WeekDays))
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, memberIds)
	s.Send2(notifyCtx, m, memberIds...)
}

func (s *Notify) SetCommentDeletable(e *event.SetCommentDeletable) {
	s.log.Infof("SetCommentDeletable: %+v", e)
	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.Of(managerIds).Diff(e.Operator).List()

	buildTemplate := func() string {
		action := "关闭"
		if e.Deletable == 1 {
			action = "开启"
		}
		return fmt.Sprintf("%v%v了评论可删除", parseUserTmp(operator), action)
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	s.Send2(notifyCtx, m, managerIds...)
}

func (s *Notify) SetCommentDeletableWhenArchived(e *event.SetCommentDeletableWhenArchived) {
	s.log.Infof("SetCommentDeletableWhenArchived: %+v", e)
	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.Of(managerIds).Diff(e.Operator).List()

	buildTemplate := func() string {
		if e.Value == 1 {
			return fmt.Sprintf("%v开启了评论归档可删除", parseUserTmp(operator))
		} else {
			return fmt.Sprintf("%v关闭了评论归档可删除", parseUserTmp(operator))
		}
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	s.Send2(notifyCtx, m, managerIds...)
}

func (s *Notify) SetCommentShowPos(e *event.SetCommentShowPos) {
	s.log.Infof("SetCommentShowPos: %+v", e)
	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.Of(managerIds).Diff(e.Operator).List()

	buildTemplate := func() string {
		if e.Value == 1 {
			return fmt.Sprintf("%v开启了评论独立显示", parseUserTmp(operator))
		} else {
			return fmt.Sprintf("%v关闭了评论独立显示", parseUserTmp(operator))
		}
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	s.Send2(notifyCtx, m, managerIds...)
}

func (s *Notify) SetSpaceNotify(e *event.SetSpaceNotify) {
	s.log.Infof("SetSpaceNotify: %+v", e)
	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.Of(managerIds).Diff(e.Operator).List()

	buildTemplate := func() string {
		if e.Notify == 1 {
			return fmt.Sprintf("%v开启了项目消息通知", parseUserTmp(operator))
		} else {
			return fmt.Sprintf("%v关闭了项目消息通知", parseUserTmp(operator))
		}
	}

	m := buildSpaceMsg(e.Event, operator, space).
		SetDescribe(buildTemplate())
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds, WithNotifyCtxForceNotifyOpt())
	s.Send2(notifyCtx, m, managerIds...)
}

func (s *Notify) SpaceAbnormal(e *event.SpaceAbnormal) {
	s.log.Infof("SpaceAbnormal: %+v", e)
	ctx := context.Background()

	space := e.Space

	// 仅在工作日推送
	spaceCfg, _ := s.spaceRepo.GetSpaceConfig(ctx, space.Id)
	if !spaceCfg.WorkingDay.IsWorkingDay(time.Now().Weekday()) {
		return
	}

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)

	// 项目进度异常，逾期任务数 [逾期任务数量]
	desc := fmt.Sprintf("项目进度异常，逾期任务数 [%v]", e.ExpiredNum)

	m := buildSpaceMsg(e.Event, nil, space).
		SetDescribe(desc)
	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	s.Send2(notifyCtx, m, managerIds...)
	imMsg := tea_im.NewRobotMessage().
		SetTitle("项目进度异常").
		SetTextRich(parseToImRich(MinorColorSpan(space.SpaceName+" ") + "\n" + desc)).
		SetShowType(tea_im.ShowType_Title).
		SetIcon(s.makeIconResLink("", IconRes_Project)).
		SetSVGIcon(IconSVGRes_Project).
		SetUrl(s.makeSpaceRedirectLink("", space.Id))
	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMsg, managerIds)
	})
}

func (s *Notify) addSpaceManager(e *event.AddSpaceManager) {
	s.log.Infof("addSpaceManager: %+v", e)

	ctx := context.Background()

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.TargetId})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds). //管理员
		Concat(space.UserId). //加上空间创建人
		Unique(). //去重
		Diff(operator.Id). //排除操作人
		List()

	buildTemplate2 := func() string {
		//  [昵称（用户名）] 将你添加至“系统用户组”
		return fmt.Sprintf("%v 将你设置为“项目管理员”", parseUserTmp(operator))
	}

	m := buildSpaceMsg(e.Event, operator, space)
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(managerIds, []int64{e.TargetId}))

	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(buildTemplate2()), e.TargetId)
	})
}

func (s *Notify) removeSpaceManager(e *event.RemoveSpaceManager) {

	ctx := context.Background()

	s.log.Infof("removeSpaceManager: %+v", e)

	space := e.Space

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.TargetId})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds). //管理员
		Concat(space.UserId). //加上空间创建人
		Unique(). //去重
		Diff(operator.Id). //排除操作人
		List()

	buildTemplate2 := func() string {
		//  [昵称（用户名）] 将你从用户组：“系统用户组” 中移除
		return fmt.Sprintf("%v 将你移除“项目管理员”", parseUserTmp(operator))
	}

	desc2 := buildTemplate2()

	m := buildSpaceMsg(e.Event, operator, space)
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(managerIds, []int64{e.TargetId}))
	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc2), e.TargetId)
	})
}
