package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	msg "go-cs/internal/bean/vo/message/notify_message"
	domain_message "go-cs/internal/domain/pkg/message"
	space_domain "go-cs/internal/domain/space"
	user_domain "go-cs/internal/domain/user"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"time"
)

func buildSpaceViewMsg(desc string, event notify.Event, user *user_domain.User, space *space_domain.Space, data *msg.ViewData) *msg.Message {
	m := &msg.Message{
		Space: &msg.Space{
			SpaceId:   space.Id,
			SpaceName: space.SpaceName,
		},
		Relation: make([]msg.RelationType, 0),
		Type:     event,
		TypeDesc: event.String(),
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
				Type: msg.ObjectType_View,
				Data: data,
			},
			Describe: desc,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) updateSpaceViewByDomainMessage(e *domain_message.UpdateSpaceView) {
	s.log.Infof("updateSpaceViewByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Oper.GetId()})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Oper.GetId()]

	memberIds, _ := s.memberRepo.GetSpaceAllMemberIds(ctx, space.Id)
	memberIds = stream.
		Of(memberIds).
		Diff(operator.Id). //排除操作人
		List()

	buildTemplate := func() string {
		// [昵称（用户名）] 新建任务流程“XXX”
		return fmt.Sprintf("%v 已变更视图配置：%v", parseUserTmp(operator), e.ViewName)
	}

	desc := buildTemplate()

	m := buildSpaceViewMsg(desc, notify.Event_UpdateSpaceView, operator, space, &msg.ViewData{Id: e.ViewId, Name: e.ViewName, Key: e.ViewKey, Field: e.Field})
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, memberIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc), memberIds...)
	})
}
