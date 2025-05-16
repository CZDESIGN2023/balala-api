package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	"time"

	user_domain "go-cs/internal/domain/user"
)

func buildAdminMsg(desc string, event notify.Event, user *user_domain.User) *msg.Message {
	m := &msg.Message{
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
			Describe: desc,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) adminChangeUserNicknameByDomainMessage(e *domain_message.AdminChangeUserNickname) {
	s.adminChangeUserNickname(&event.AdminChangeUserNickname{
		Event:    notify.Event_AdminChangeUserNickname,
		Operator: e.Oper.GetId(),
		UserId:   e.UserId,
		OldValue: e.OldValue,
		NewValue: e.NewValue,
	})
}

func (s *Notify) adminChangeUserRoleByDomainMessage(e *domain_message.AdminChangeUserRole) {
	s.adminChangeUserRole(&event.AdminChangeUserRole{
		Event:    notify.Event_AdminChangeUserRole,
		Operator: e.Oper.GetId(),
		UserId:   e.UserId,
		OldValue: e.OldValue,
		NewValue: e.NewValue,
	})
}

func (s *Notify) adminChangeUserNickname(e *event.AdminChangeUserNickname) {
	s.log.Infof("adminChangeUserNicknameByDomainMessage: %+v", e)

	ctx := context.Background()

	if e.Operator == e.UserId {
		return
	}

	operator, err := s.userRepo.GetUserByUserId(ctx, e.Operator)
	if err != nil {
		s.log.Error(err)
		return
	}

	desc := fmt.Sprintf("系统管理员已将你的昵称变更<br/>“%v” -> “%v”", e.OldValue, e.NewValue)
	msg := buildAdminMsg(desc, e.Event, operator)

	if e.Operator == e.UserId {
		return
	}

	s.Send(msg, e.UserId)
}

func (s *Notify) adminChangeUserRole(e *event.AdminChangeUserRole) {
	s.log.Infof("adminChangeUserRole: %+v", e)

	ctx := context.Background()

	operator, err := s.userRepo.GetUserByUserId(ctx, e.Operator)
	if err != nil {
		s.log.Error(err)
		return
	}

	desc := fmt.Sprintf("admin已将你设置为%s", e.NewValue)
	if e.NewValue == consts.SystemRole_Normal {
		desc = fmt.Sprintf("admin已将你%s身份撤销", e.OldValue)

	}

	message := buildAdminMsg(desc, e.Event, operator)

	s.Send(message, e.UserId)
}

func (s *Notify) adminCancelUserByDomainMessage(e *domain_message.AdminCancelUser) {
	s.CooperateAdminCancelUser(&event.AdminCancelUser{
		Event:    notify.Event_AdminCancelUser,
		Operator: e.Oper.GetId(),
		UserId:   e.UserId,
	})
}
