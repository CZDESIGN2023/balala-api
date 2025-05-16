package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/data/convert"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"slices"
	"time"
)

type buildWorkItemStatsMessageVo struct {
	Event     notify.Event
	Describe  string
	SpaceId   int64
	SpaceName string

	WorkItemPid  int64
	WorkItemId   int64
	WorkItemName string

	Operator *db.User
}

func buildWorkItemStatsMessage(in *buildWorkItemStatsMessageVo) *msg.Message {
	m := &msg.Message{
		Space: &msg.Space{
			SpaceId:   in.SpaceId,
			SpaceName: in.SpaceName,
		},
		Relation: make([]msg.RelationType, 0),
		Type:     in.Event,
		TypeDesc: in.Event.String(),
		Notification: &msg.Notification{
			Action: msg.ActionType_edit,
			Subject: &msg.Subject{
				Type: msg.SubjectType_user,
				Data: &msg.UserData{
					Name:     in.Operator.UserName,
					NickName: in.Operator.UserName,
					Id:       in.Operator.Id,
					Avatar:   in.Operator.Avatar,
				},
			},
			Object: &msg.Object{
				Type: msg.ObjectType_workItem,
				Data: &msg.WorkItemData{
					Id:   in.WorkItemId,
					Name: in.WorkItemName,
					Pid:  in.WorkItemPid,
				},
			},
			Describe: in.Describe,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) closeWorkItem(e *event.CloseWorkItem) {

	ctx := context.Background()
	evData := e.Data

	workItem, err := s.workItemRepo.GetWorkItem(ctx, evData.WorkItem.Id, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)

	if err != nil {
		s.log.Error(err)
		return
	}

	creator, directors, followers, _ := splitUser(
		evData.Operator.Id,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	operator, err := s.userRepo.GetUserByUserId(ctx, evData.Operator.Id)
	if err != nil {
		s.log.Error(err)
		return
	}

	buildTemplate := func() string {
		reason := evData.Reason
		//任务已关闭，关闭原因“[关闭原因]”，操作人：[昵称（用户名）]
		return fmt.Sprintf("任务已关闭</br>操作人：%v</br>关闭原因: %v", parseUserTmp(operator), reason)
	}

	buildImMessage := func(relation string) *tea_im.RobotMessage {

		reason := evData.Reason

		m := tea_im.NewRobotMessage(
			tea_im.WithShowSubContentTypeOption(),
		)

		if relation == "follow" {
			m.SetIcon(s.makeIconResLink("", IconRes_Follow)).
				SetSVGIcon(IconSVGRes_Follow)
			m.SetTitle("关注任务提醒")
		} else {
			m.SetIcon(s.makeIconResLink("", IconRes_Task)).
				SetSVGIcon(IconSVGRes_Task)
			m.SetTitle("任务提醒")
		}

		m.SetSubTitle(evData.Space.Name)
		m.SetSubContent(workItem.WorkItemName)
		m.SetUrl(s.makeWorkItemRedirectLink("", evData.Space.Id, workItem.Id, workItem.Pid))

		txt := fmt.Sprintf("任务已关闭\n操作人：%v\n关闭原因: %v", parseUserTmp(operator), reason)
		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemStatsMessage(&buildWorkItemStatsMessageVo{
		Event:        e.Event,
		Describe:     buildTemplate(),
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemName: workItem.WorkItemName,
		WorkItemPid:  workItem.Pid,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetPopup()
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	space, _ := s.spaceRepo.GetSpace(ctx, evData.Space.Id)
	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	imMessage := buildImMessage("")
	imMessageWithFollow := buildImMessage("follow")
	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, creator)
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, directors)
		s.pushThirdPlatformMessage2(notifyCtx, imMessageWithFollow, followers)
	})
}

func (s *Notify) terminateWorkItem(event *event.TerminateWorkItem) {

	ctx := context.Background()
	evData := event.Data

	workItem, err := s.workItemRepo.GetWorkItem(ctx, evData.WorkItem.Id, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)

	if err != nil {
		return
	}

	creator, directors, followers, _ := splitUser(
		evData.Operator.Id,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	operator, err := s.userRepo.GetUserByUserId(ctx, evData.Operator.Id)
	if err != nil {
		s.log.Error(err)
		return
	}

	buildTemplate := func() string {
		reason := evData.Reason
		//任务已终止，终止原因“[终止原因]”，操作人：[昵称（用户名）]
		return fmt.Sprintf("任务已终止</br>操作人：%v</br>终止原因：%v", parseUserTmp(operator), reason)
	}

	buildImMessage := func(relation string) *tea_im.RobotMessage {

		reason := evData.Reason

		m := tea_im.NewRobotMessage(
			tea_im.WithShowSubContentTypeOption(),
		)

		if relation == "follow" {
			m.SetIcon(s.makeIconResLink("", IconRes_Follow)).
				SetSVGIcon(IconSVGRes_Follow)
			m.SetTitle("关注任务提醒")
		} else {
			m.SetIcon(s.makeIconResLink("", IconRes_Task)).
				SetSVGIcon(IconSVGRes_Task)
			m.SetTitle("任务提醒")
		}

		m.SetSubTitle(evData.Space.Name)
		m.SetSubContent(workItem.WorkItemName)
		m.SetUrl(s.makeWorkItemRedirectLink("", evData.Space.Id, workItem.Id, workItem.Pid))

		txt := fmt.Sprintf("任务已终止\n操作人：%v\n终止原因：%v", parseUserTmp(operator), reason)
		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemStatsMessage(&buildWorkItemStatsMessageVo{
		Event:        event.Event,
		Describe:     buildTemplate(),
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemPid:  workItem.Pid,
		WorkItemName: workItem.WorkItemName,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	space, _ := s.spaceRepo.GetSpace(ctx, evData.Space.Id)
	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	imMessage := buildImMessage("")
	imMessageWithFollow := buildImMessage("follow")
	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, creator)
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, directors)
		s.pushThirdPlatformMessage2(notifyCtx, imMessageWithFollow, followers)
	})

}

func (s *Notify) completeWorkItem(e *event.CompleteWorkItem) {

	ctx := context.Background()
	evData := e.Data

	workItem, err := s.workItemRepo.GetWorkItem(ctx, evData.WorkItem.Id, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)

	if err != nil {
		return
	}

	creator, directors, followers, _ := splitUser(
		evData.Operator.Id,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	operator, err := s.userRepo.GetUserByUserId(ctx, evData.Operator.Id)
	if err != nil {
		s.log.Error(err)
		return
	}

	buildTemplate := func() string {
		//任务已完成，操作人：[昵称（用户名）]
		return fmt.Sprintf("任务已完成</br>操作人：%v", parseUserTmp(operator))
	}

	buildImMessage := func(relation string) *tea_im.RobotMessage {

		m := tea_im.NewRobotMessage(
			tea_im.WithShowSubContentTypeOption(),
		)

		if relation == "follow" {
			m.SetIcon(s.makeIconResLink("", IconRes_Follow)).
				SetSVGIcon(IconSVGRes_Follow)
			m.SetTitle("关注任务提醒")
		} else {
			m.SetIcon(s.makeIconResLink("", IconRes_Task)).
				SetSVGIcon(IconSVGRes_Task)
			m.SetTitle("任务提醒")
		}

		m.SetSubTitle(evData.Space.Name)
		m.SetSubContent(workItem.WorkItemName)
		m.SetUrl(s.makeWorkItemRedirectLink("", evData.Space.Id, workItem.Id, workItem.Pid))

		txt := fmt.Sprintf("任务已完成\n操作人：%v", parseUserTmp(operator))
		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemStatsMessage(&buildWorkItemStatsMessageVo{
		Event:        e.Event,
		Describe:     buildTemplate(),
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemName: workItem.WorkItemName,
		WorkItemPid:  workItem.Pid,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetPopup()
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	// toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	space, _ := s.spaceRepo.GetSpace(ctx, evData.Space.Id)
	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		// s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	imMessage := buildImMessage("")
	imMessageWithFollow := buildImMessage("follow")
	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, creator)
		s.pushThirdPlatformMessage2(notifyCtx, imMessageWithFollow, followers)
	})
}

func (s *Notify) restartWorkItem(e *event.RestartWorkItem) {

	ctx := context.Background()
	evData := e.Data

	workItem, err := s.workItemRepo.GetWorkItem(ctx, evData.WorkItem.Id, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)
	if err != nil {
		return
	}

	var nodeName string
	if e.Data.ToNode != nil && e.Data.ToNode.Code != "" {
		tplt, err := s.workFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, workItem.WorkFlowTemplateId)
		if err != nil || tplt.WorkFLowConfig == nil {
			return
		}

		nodeConf := tplt.WorkFLowConfig.GetNode(e.Data.ToNode.Code)
		if nodeConf == nil {
			return
		}
		nodeName = nodeConf.Name
	}

	rawCreator, operatorId, rawDirectors, rawFollowers := workItem.UserId, evData.Operator.Id, utils.ToInt64Array(workItem.Doc.Directors), utils.ToInt64Array(workItem.Doc.Followers)
	creators, directors, followers, _ := splitUser(
		operatorId,
		rawCreator,
		rawDirectors,
		rawFollowers,
	)

	operator, err := s.userRepo.GetUserByUserId(ctx, evData.Operator.Id)
	if err != nil {
		s.log.Error(err)
		return
	}

	buildTemplate := func() string {

		reason := evData.Reason

		var nodeNamePart string
		var reasonPart string
		if nodeName != "" {
			nodeNamePart = fmt.Sprintf("至 %v", nodeName)
		}
		if reason != "" {
			reasonPart = fmt.Sprintf("</br>重启原因：%v", reason)
		}

		//任务已重启至 [节点名称]，重启原因“[重启原因]”，操作人：[昵称（用户名）]
		return fmt.Sprintf("任务已重启%v</br>操作人：%v%v", nodeNamePart, parseUserTmp(operator), reasonPart)
	}

	buildImMessage := func(relation string) *tea_im.RobotMessage {

		reason := evData.Reason

		m := tea_im.NewRobotMessage(
			tea_im.WithShowSubContentTypeOption(),
		)

		if relation == "follow" {
			m.SetIcon(s.makeIconResLink("", IconRes_Follow)).
				SetSVGIcon(IconSVGRes_Follow)
			m.SetTitle("关注任务提醒")
		} else {
			m.SetIcon(s.makeIconResLink("", IconRes_Task)).
				SetSVGIcon(IconSVGRes_Task)
			m.SetTitle("任务提醒")
		}

		m.SetSubTitle(evData.Space.Name)
		m.SetSubContent(workItem.WorkItemName)
		m.SetUrl(s.makeWorkItemRedirectLink("", evData.Space.Id, workItem.Id, workItem.Pid))

		var nodeNamePart string
		var reasonPart string
		if nodeName != "" {
			nodeNamePart = fmt.Sprintf("至 %v", nodeName)
		}
		if reason != "" {
			reasonPart = fmt.Sprintf("</br>重启原因：%v", reason)
		}

		txt := fmt.Sprintf("任务已重启%v\n操作人：%v\n%v", nodeNamePart, parseUserTmp(operator), reasonPart)
		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemStatsMessage(&buildWorkItemStatsMessageVo{
		Event:        e.Event,
		Describe:     "",
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemName: workItem.WorkItemName,
		WorkItemPid:  workItem.Pid,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetPopup()
	m.SetRedirectLink(s.makeRedirectLink(m))

	desc := buildTemplate()
	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toOwner.SetDescribe(desc)

	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toDirector.SetDescribe(desc)

	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toFollower.SetDescribe(desc)

	if s.isTodoWorkItem(workItem) {
		if slices.Contains(rawDirectors, rawCreator) {
			toOwner.SetRelation(msg.Relation_workItemTodo)
		}
		toDirector.SetRelation(msg.Relation_workItemTodo)
	}

	space, _ := s.spaceRepo.GetSpace(ctx, evData.Space.Id)
	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creators, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creators...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	imMessage := buildImMessage("")
	imMessageWithFollow := buildImMessage("follow")
	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, creators)
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, directors)
		s.pushThirdPlatformMessage2(notifyCtx, imMessageWithFollow, followers)
	})
}

func (s *Notify) resumeWorkItem(e *event.ResumeWorkItem) {

	ctx := context.Background()
	evData := e.Data

	workItem, err := s.workItemRepo.GetWorkItem(ctx, evData.WorkItem.Id, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)
	if err != nil {
		return
	}

	rawCreator, operatorId, rawDirectors, rawFollowers := workItem.UserId, evData.Operator.Id, utils.ToInt64Array(workItem.Doc.Directors), utils.ToInt64Array(workItem.Doc.Followers)
	creator, directors, followers, _ := splitUser(
		operatorId,
		rawCreator,
		rawDirectors,
		rawFollowers,
	)

	operator, err := s.userRepo.GetUserByUserId(ctx, evData.Operator.Id)
	if err != nil {
		s.log.Error(err)
		return
	}

	buildTemplate := func() string {
		reason := evData.Reason
		if evData.Scene == event.ResumeWorkItemScene_formClosed {
			return fmt.Sprintf("任务已重启<br/>操作人：%v<br/>重启原因: %v", parseUserTmp(operator), reason)
		}

		//任务恢复，恢复原因“[恢复原因]”，操作人：[昵称（用户名）]
		return fmt.Sprintf("任务已恢复<br/>操作人：%v<br/>恢复原因: %v", parseUserTmp(operator), reason)
	}

	buildImMessage := func(relation string) *tea_im.RobotMessage {

		reason := evData.Reason

		m := tea_im.NewRobotMessage(
			tea_im.WithShowSubContentTypeOption(),
		)
		if relation == "follow" {
			m.SetIcon(s.makeIconResLink("", IconRes_Follow)).
				SetSVGIcon(IconSVGRes_Follow)
			m.SetTitle("关注任务提醒")
		} else {
			m.SetIcon(s.makeIconResLink("", IconRes_Task)).
				SetSVGIcon(IconSVGRes_Task)
			m.SetTitle("任务提醒")
		}

		m.SetSubTitle(evData.Space.Name)
		m.SetSubContent(workItem.WorkItemName)
		m.SetUrl(s.makeWorkItemRedirectLink("", evData.Space.Id, workItem.Id, workItem.Pid))

		var txt string
		if evData.Scene == event.ResumeWorkItemScene_formClosed {
			txt = fmt.Sprintf("任务已重启\n操作人：%v\n重启原因: %v", parseUserTmp(operator), reason)
		} else {
			txt = fmt.Sprintf("任务已恢复\n操作人：%v\n恢复原因: %v", parseUserTmp(operator), reason)
		}

		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemStatsMessage(&buildWorkItemStatsMessageVo{
		Event:        e.Event,
		Describe:     buildTemplate(),
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemName: workItem.WorkItemName,
		WorkItemPid:  workItem.Pid,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()
	// if evData.Scene == event.ResumeWorkItemScene_formClosed {
	// 	m.SetPopup()
	// }

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)

	if s.isTodoWorkItem(workItem) {
		if slices.Contains(rawDirectors, rawCreator) {
			toOwner.SetRelation(msg.Relation_workItemTodo)
		}
		toDirector.SetRelation(msg.Relation_workItemTodo)
	}

	space, _ := s.spaceRepo.GetSpace(ctx, evData.Space.Id)
	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	imMessage := buildImMessage("")
	imMessageWithFollow := buildImMessage("follow")

	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, creator)
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, directors)
		s.pushThirdPlatformMessage2(notifyCtx, imMessageWithFollow, followers)
	})
}
