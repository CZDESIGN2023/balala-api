package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	domain_message "go-cs/internal/domain/pkg/message"
	space_domain "go-cs/internal/domain/space"
	user_domain "go-cs/internal/domain/user"
	"go-cs/internal/domain/work_flow"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"time"

	"github.com/spf13/cast"
)

func buildWorkFlowMsg(desc string, event notify.Event, user *user_domain.User, space *space_domain.Space, workFlowData *msg.WorkFlowData) *msg.Message {
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
				Type: msg.ObjectType_workFlow,
				Data: workFlowData,
			},
			Describe: desc,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) createWorkFlowByDomainMessage(e *domain_message.CreateWorkFlow) {
	s.log.Infof("createWorkFlowByDomainMessage: %+v", e)

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

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds).       //管理员
		Concat(space.UserId). //加上空间创建人
		Unique().             //去重
		Diff(operator.Id).    //排除操作人
		List()

	buildTemplate := func() string {
		// [昵称（用户名）] 新建任务流程“XXX”
		return fmt.Sprintf("%v 新建任务流程“%v”", parseUserTmp(operator), e.WorkFlowName)
	}

	desc := buildTemplate()

	m := buildWorkFlowMsg(desc, notify.Event_CreateWorkFlow, operator, space, &msg.WorkFlowData{Id: e.WorkFlowId, Name: e.WorkFlowName})
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc), managerIds...)
	})
}

func (s *Notify) modifyWorkFlowByDomainMessage(e *domain_message.ModifyWorkFlow) {
	s.log.Infof("modifyWorkFlowByDomainMessage: %+v", e)

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

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds).       //管理员
		Concat(space.UserId). //加上空间创建人
		Unique().             //去重
		Diff(operator.Id).    //排除操作人
		List()

	buildStatusTemplate := func(newStatus, oldStatus work_flow.WorkFlowStatus) string {
		var action string
		switch {
		case newStatus == work_flow.WorkFlowStatus_Enable:
			action = "启用"
		case oldStatus == work_flow.WorkFlowStatus_Hide && newStatus == work_flow.WorkFlowStatus_Disable:
			action = "显示"
		case newStatus == work_flow.WorkFlowStatus_Hide && oldStatus == work_flow.WorkFlowStatus_Hide:
			action = "隐藏"
		case newStatus == work_flow.WorkFlowStatus_Disable:
			action = "禁用"
		}

		//  [昵称（用户名）] 已启用任务流程“XXX”
		return fmt.Sprintf("%v %v任务流程“%v”", parseUserTmp(operator), action, e.WorkFlowName)
	}

	for _, v := range e.Updates {
		var desc string
		switch v.Field {
		case "status":
			desc = buildStatusTemplate(v.NewValue.(work_flow.WorkFlowStatus), v.OldValue.(work_flow.WorkFlowStatus))
		}

		if desc == "" {
			continue
		}

		m := buildWorkFlowMsg(desc, notify.Event_ChangeWorkFlowField, operator, space, &msg.WorkFlowData{Id: e.WorkFlowId, Name: e.WorkFlowName})
		m.SetRedirectLink(s.makeRedirectLink(m))

		notifyCtx := s.buildNotifyCtx(space, managerIds)
		utils.Go(func() {
			s.Send2(notifyCtx, m.Clone().SetDescribe(desc), managerIds...)
		})

	}

	//协作
	for _, v := range e.Updates {
		v := v
		if v.Field == "status" && cast.ToInt64(v.NewValue) == 0 {
			utils.Go(func() {
				s.CooperateDisableWorkFlow(&event.CooperateDisableWorkFlow{
					Event:        notify.Event_DisableWorkFlow,
					Operator:     e.Oper.GetId(),
					SpaceId:      e.SpaceId,
					WorkFlowId:   e.WorkFlowId,
					WorkFlowName: e.WorkFlowName,
				})
			})
			break
		}
	}
}

func (s *Notify) deleteWorkFlowByDomainMessage(e *domain_message.DeleteWorkFlow) {
	s.log.Infof("deleteWorkFlowByDomainMessage: %+v", e)

	//发送协作消息
	utils.Go(func() {
		s.CooperateDeleteWorkFlow(&event.CooperateDeleteWorkFlow{
			Event:        notify.Event_DeleteWorkFlow,
			Operator:     e.Oper.GetId(),
			SpaceId:      e.SpaceId,
			WorkFlowId:   e.WorkFlowId,
			WorkFlowName: e.WorkFlowName,
		})
	})

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

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds).       //管理员
		Concat(space.UserId). //加上空间创建人
		Unique().             //去重
		Diff(operator.Id).    //排除操作人
		List()

	buildTemplate := func() string {
		//  [昵称（用户名）] 已删除任务流程“XXX”
		return fmt.Sprintf("%v 已删除任务流程“%v”", parseUserTmp(operator), e.WorkFlowName)
	}

	desc := buildTemplate()

	m := buildWorkFlowMsg(desc, notify.Event_CreateWorkFlow, operator, space, &msg.WorkFlowData{
		Id:   e.WorkFlowId,
		Name: e.WorkFlowName,
	})
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc), managerIds...)
	})

}

func (s *Notify) saveWorkFlowTemplateByDomainMessage(e *domain_message.SaveWorkFlowTemplate) {
	s.log.Infof("saveWorkFlowTemplateByDomainMessage: %+v", e)

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

	managerIds, _ := s.memberRepo.GetSuperManagerIds(ctx, space.Id)
	managerIds = stream.
		Of(managerIds).       //管理员
		Concat(space.UserId). //加上空间创建人
		Unique().             //去重
		Diff(operator.Id).    //排除操作人
		List()

	buildTemplate := func() string {
		// "[昵称（用户名）] 变更了任务流程“XXX”配置，请注意查看
		return fmt.Sprintf("%v 变更了任务流程“%v”配置，请注意查看", parseUserTmp(operator), e.WorkFlowName)
	}

	desc := buildTemplate()

	m := buildWorkFlowMsg(desc, notify.Event_UpdateWorkFlow, operator, space, &msg.WorkFlowData{Id: e.WorkFlowId, Name: e.WorkFlowName})
	m.SetRedirectLink(s.makeRedirectLink(m))

	notifyCtx := s.buildNotifyCtx(space, managerIds)
	utils.Go(func() {
		s.Send2(notifyCtx, m.Clone().SetDescribe(desc), managerIds...)
	})
}

func (s *Notify) upgradeWorkFlowByDomainMessage(e *domain_message.UpgradeTaskWorkFlow) {
	s.log.Infof("upgradeWorkFlowByDomainMessage: %+v", e)

	//协作消息
	s.CooperateUpgradeWorkFlow(&event.CooperateUpgradeWorkFlow{
		Event:        notify.Event_UpgradeTaskWorkFlow,
		Operator:     e.Oper.GetId(),
		SpaceId:      e.SpaceId,
		WorkFlowId:   e.WorkFlowId,
		WorkItemId:   e.WorkItemId,
		WorkFlowName: e.WorkFlowName,
	})
}
