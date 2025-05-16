package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"time"

	space_domain "go-cs/internal/domain/space"
	user_domain "go-cs/internal/domain/user"
)

func buildUserMsg(desc string, event notify.Event, user *user_domain.User, space *space_domain.Space) *msg.Message {
	m := &msg.Message{
		Space: &msg.Space{
			SpaceId:   space.Id,
			SpaceName: space.SpaceName,
		},
		Type:     event,
		TypeDesc: event.String(),
		Relation: make([]msg.RelationType, 0),
		Notification: &msg.Notification{
			Action: msg.ActionType_edit,
			Subject: &msg.Subject{
				Type: msg.SubjectType_user,
				Data: &msg.UserData{
					Name:     user.UserName,
					NickName: user.UserName,
					Id:       user.Id,
					Avatar:   user.Avatar,
				},
			},
			Object: &msg.Object{
				Type: msg.ObjectType_space,
				Data: &msg.SpaceData{
					Id:   space.Id,
					Name: space.SpaceName,
				},
			},
			Describe: desc,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) QuitSpace(e *event.QuitSpace) {
	s.log.Infof("QuitSpace: %+v", e)

	ctx := context.Background()
	space := e.Space

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.Remove(managerIds, e.Operator)

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	desc := fmt.Sprintf("%s 已退出项目", parseUserTmp(operator))

	m := buildUserMsg(desc, e.Event, operator, space)
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	s.Send2(notifyCtx, m, managerIds...)
}

func (s *Notify) TransferSpace(e *event.TransferSpace) {
	s.log.Infof("TransferSpace: %+v", e)

	ctx := context.Background()

	space := e.Space

	allMemberIds, _ := s.memberRepo.GetSpaceAllMemberIds(ctx, space.Id)
	allMemberIds = stream.Diff(allMemberIds, []int64{e.Operator, e.TargetId})

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator, e.TargetId})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]
	target := userMap[e.TargetId]

	desc1 := fmt.Sprintf("%v 将项目创建者转移给 %v", parseUserTmp(operator), parseUserTmp(target))
	desc2 := fmt.Sprintf("%v 将项目创建者转移给你", parseUserTmp(operator))

	m := buildUserMsg("", e.Event, operator, space)
	m.SetRedirectLink(s.makeTableRedirectLink("", space.Id))

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(allMemberIds, []int64{e.TargetId}))
	s.Send2(notifyCtx, m.Clone().SetDescribe(desc1), allMemberIds...)
	s.Send2(notifyCtx, m.Clone().SetDescribe(desc2).SetPopup(), e.TargetId)

	qlTargetMsg := tea_im.NewRobotMessage().SetShowType(tea_im.ShowType_Text).
		SetTitle("项目提醒").
		SetTextRich(parseToImRich(MinorColorSpan(space.SpaceName+" ") + "\n" + desc2)).
		SetIcon(s.makeIconResLink("", IconRes_Project)).
		SetSVGIcon(IconSVGRes_Project).
		SetUrl(s.makeTableRedirectLink("", space.Id))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, qlTargetMsg, []int64{e.TargetId})
	})
}

func (s *Notify) TransferWorkItem(e *event.TransferWorkItem) {
	s.log.Infof("TransferWorkItem: %+v", e)

	space := e.Space

	userMap, _ := s.userRepo.UserMap(context.Background(), []int64{e.Operator})
	operator := userMap[e.Operator]

	buildTemplate := func() string {
		return fmt.Sprintf("%s 将%v条任务单转移给你", parseUserTmp(operator), e.Num)
	}

	//组织推送内容

	desc := buildTemplate()

	m := buildUserMsg(desc, e.Event, operator, space)
	m.SetRelation(msg.Relation_workItemTodo).
		SetRedirectLink(s.makeTableRedirectLink("", space.Id)).
		SetPopup()

	userIds := []int64{e.TargetId}

	notifyCtx := s.buildNotifyCtx(space, userIds)
	s.Send2(notifyCtx, m, userIds...)

	imMsg := tea_im.NewRobotMessage().SetShowType(tea_im.ShowType_Text).
		SetTitle("任务提醒").SetSubTitle(space.SpaceName).
		SetTextRich(parseToImRich(desc)).
		SetIcon(s.makeIconResLink("", IconRes_Task)).
		SetSVGIcon(IconSVGRes_Task).
		SetUrl(s.makeTableRedirectLink("", space.Id))

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMsg, userIds)
	})
}
