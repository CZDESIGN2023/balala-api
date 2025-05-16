package notify

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/event"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"slices"
	"time"
)

type buildWorkItemFlowNodeMessageVo struct {
	Event     notify.Event
	Describe  string
	SpaceId   int64
	SpaceName string

	WorkItemId   int64
	WorkItemName string

	Operator *db.User
}

func buildWorkItemFlowNodeMessage(in *buildWorkItemFlowNodeMessageVo) *msg.Message {

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
				},
			},
			Describe: in.Describe,
			Date:     time.Now(),
		},
	}

	return m
}

func (s *Notify) changeWorkItemFlowNode(e *event.ChangeWorkItemFlowNode) {

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

	flowNodeName := e.Data.ToNode.Name

	var (
		operatorId   = evData.Operator.Id
		rawCreator   = workItem.UserId
		rawDirectors = utils.StringArrToInt64Arr(workItem.Doc.Directors)
		rawFollowers = utils.StringArrToInt64Arr(workItem.Doc.Followers)
	)

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
		planStr := parsePlanTime(workItem.Doc.PlanStartAt, workItem.Doc.PlanCompleteAt)

		//"当前节点：“[节点名称]” ，总排期：[当前任务总排期]"
		switch workItem.WorkItemTypeKey {
		case consts.WorkItemTypeKey_Task:
			return fmt.Sprintf("当前节点：%v <br/>总排期：%v", flowNodeName, planStr)
		case consts.WorkItemTypeKey_StateTask:
			return fmt.Sprintf("当前状态：%v <br/>总排期：%v", flowNodeName, planStr)
		}
		return ""
	}

	buildImMessage := func(relation string) *tea_im.RobotMessage {

		planStr := parsePlanTime(workItem.Doc.PlanStartAt, workItem.Doc.PlanCompleteAt)

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
		switch workItem.WorkItemTypeKey {
		case consts.WorkItemTypeKey_Task:
			txt = fmt.Sprintf("当前节点：%v\n总排期：%v", flowNodeName, planStr)
		case consts.WorkItemTypeKey_StateTask:
			txt = fmt.Sprintf("当前状态：%v\n总排期：%v", flowNodeName, planStr)
		}
		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemFlowNodeMessage(&buildWorkItemFlowNodeMessageVo{
		Event:        e.Event,
		Describe:     buildTemplate(),
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemName: workItem.WorkItemName,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)
	toOper := m.Clone()

	if s.isTodoWorkItem(workItem) {
		if slices.Contains(rawDirectors, rawCreator) {
			toOwner.SetRelation(msg.Relation_workItemTodo)
		}
		toDirector.SetRelation(msg.Relation_workItemTodo)
	}

	space, _ := s.spaceRepo.GetSpace(ctx, evData.Space.Id)
	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send(toOper, operator.Id)
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})

	imMessage := buildImMessage("")
	utils.Go(func() {
		s.pushThirdPlatformMessage2(notifyCtx, imMessage, stream.Diff(rawDirectors, []int64{operatorId}))
	})
}

func (s *Notify) rollbackWorkItemFlowNode(e *event.RollbackWorkItemFlowNode) {

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

	var flowNodeName string
	if evData.ToNode != nil && evData.ToNode.Id != 0 {
		flowTplt, err := s.workFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, workItem.WorkFlowTemplateId)
		if err == nil && flowTplt.WorkFLowConfig != nil {
			nodeConf := flowTplt.WorkFLowConfig.GetNode(evData.ToNode.Code)
			if nodeConf != nil {
				flowNodeName = nodeConf.Name
			}
		}
	}

	if flowNodeName == "" {
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

		//"当前节点：“[节点名称]” ，总排期：[当前任务总排期]"
		return fmt.Sprintf("任务已回滚至 %v</br>操作人：%v</br>回滚原因：%v", flowNodeName, parseUserTmp(operator), reason)
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

		txt := fmt.Sprintf("任务已回滚至 %v\n操作人：%v\n回滚原因：%v", flowNodeName, parseUserTmp(operator), reason)
		m.SetTextRich(parseToImRich(txt))

		return m
	}

	m := buildWorkItemFlowNodeMessage(&buildWorkItemFlowNodeMessageVo{
		Event:        e.Event,
		Describe:     buildTemplate(),
		SpaceId:      evData.Space.Id,
		SpaceName:    evData.Space.Name,
		WorkItemId:   workItem.Id,
		WorkItemName: workItem.WorkItemName,
		Operator:     convert.UserEntityToPo(operator),
	})

	m.SetRedirectLink(s.makeRedirectLink(m))
	m.SetPopup()

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

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

func (s *Notify) ChangeWorkFlowNodePlanTime(e *event.ChangeWorkFlowNodePlanTime) {
	s.log.Infof("ChangeWorkFlowNodePlanTime: %+v", e)

	ctx := context.Background()

	space := e.Space
	workItem := e.WorkItem

	// 通知 任务创建人 / 当前节点负责人 / 关注人
	creator, directors, followers, _ := splitUser(
		e.Operator,
		workItem.UserId,
		utils.ToInt64Array(workItem.Doc.Directors),
		utils.ToInt64Array(workItem.Doc.Followers),
	)

	userMap, err := s.userRepo.UserMap(ctx, []int64{e.Operator})
	if err != nil {
		s.log.Error(err)
		return
	}

	operator := userMap[e.Operator]

	buildTemplate := func() string {
		return fmt.Sprintf(
			"%v 变更了 <br />排期：“%v” -> “%v”",
			parseUserTmp(operator),
			parsePlanTime(e.OldValues[0], e.OldValues[1]),
			parsePlanTime(e.NewValues[0], e.NewValues[1]),
		)
	}

	desc := buildTemplate()

	//组织推送内容
	m := buildWorkItemMsg(desc, e.Event, operator, space, workItem)
	m.SetRedirectLink(s.makeRedirectLink(m))

	toOwner := m.Clone().SetRelation(msg.Relation_workItemOwner)
	toFollower := m.Clone().SetRelation(msg.Relation_workItemFollower)
	toDirector := m.Clone().SetRelation(msg.Relation_workItemDirector)

	notifyCtx := s.buildNotifyCtx(space, stream.Concat(creator, directors, followers))
	utils.Go(func() {
		s.Send2(notifyCtx, toOwner, creator...)
		s.Send2(notifyCtx, toDirector, directors...)
		s.Send2(notifyCtx, toFollower, followers...)
	})
}
