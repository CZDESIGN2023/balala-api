package notify

import (
	"context"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	n_message "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/domain/space_work_item_comment"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/stream"
	"slices"
	"time"
)

func (s *Notify) addCommentEmoji(e *event.AddCommentEmoji) {
	ctx := context.Background()

	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.WorkItemId, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return
	}

	operator, err := s.userRepo.GetUserByUserId(ctx, e.Operator)
	if err != nil {
		s.log.Error(err)
		return
	}

	msg := n_message.NewMessage()
	msg.SetType(e.Event)
	msg.Space.SpaceId = space.Id
	msg.Space.SpaceName = space.SpaceName
	//故事内容
	//某人
	msg.Notification.Subject = n_message.NewCommentOper(&n_message.UserData{
		Id:       operator.Id,
		NickName: operator.UserNickname,
		Avatar:   operator.Avatar,
		Name:     operator.UserName,
	})
	//某任务某玩意
	msg.Notification.Object = n_message.NewCommentWorkItemObject(&n_message.WorkItemData{
		Id:   workItem.Id,
		Pid:  workItem.Pid,
		Name: workItem.WorkItemName,
	})
	//添加了
	msg.Notification.Action = n_message.ActionType_add
	//某评论
	msg.Notification.SubObject = n_message.NewCommentObject(&n_message.WorkItemCommentData{
		Id:    e.Comment.Id,
		Emoji: e.Emoji,
	})
	//评论内容
	describe := utils.ClearRichTextToPlanText(e.Comment.Content, false)
	msg.Notification.Describe = describe
	//故事发生时间
	msg.Notification.Date = time.Now()

	msg.SetRedirectLink(s.makeWorkItemCommentRedirectLink("", space.Id, workItem.Id, workItem.Pid)).
		SetPopup()

	userIds := stream.Diff([]int64{e.Comment.UserId}, []int64{operator.Id})

	notifyCtx := s.buildNotifyCtx(space, userIds)

	emitCommentEvent := func(userIds []int64, msg *n_message.Message) {
		for _, userId := range userIds {
			userId := userId // copy
			member := notifyCtx.memberMap[userId]

			globalNotify := notifyCtx.userNotifySwitchGlobalMap[userId]
			spaceNotify := notifyCtx.userNotifySwitchSpaceMap[userId]

			if globalNotify.Value == "1" &&
				spaceNotify.Value == "1" &&
				(space == nil || space.Notify == 1) &&
				(member == nil || member.Notify == 1) {
				//保存评论通知快照
				utils.Go(func() {
					nSnapShot := s.notifySnapShotService.CreateNotifySnapShot(context.Background(), msg.Space.SpaceId, userId, int64(notify.Event_AddCommentEvent), msg.String())
					s.notifySnapShotRepo.CreateNotify(context.Background(), nSnapShot)
				})
			}
		}
	}

	emitCommentEvent(userIds, msg)

	// IM通知
	imMsg := tea_im.NewRobotMessage().
		SetTitle("评论提醒").SetSubTitle(space.SpaceName).SetSubContent(workItem.WorkItemName).
		SetShowType(tea_im.ShowType_Text).
		SetIcon(s.makeIconResLink("", IconRes_Comment)).
		SetSVGIcon(IconSVGRes_Comment).
		SetUrl(s.makeWorkItemCommentRedirectLink("", space.Id, workItem.Id, workItem.Pid)).
		SetTextRich(parseToImRich(parseUserTmp(operator) + "回复 " + e.Emoji))

	utils.Go(func() {
		s.Send2(notifyCtx, msg, userIds...)

		s.pushThirdPlatformMessage2(notifyCtx, imMsg, userIds)
	})
}

func (s *Notify) addComment(e *event.AddCommentEvent) {

	ctx := context.Background()
	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.EvData.WorkItemId, &repo.WithDocOption{
		PlanTime:      true,
		Participators: true,
		Followers:     true,
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return
	}

	operator, err := s.userRepo.GetUserByUserId(ctx, e.EvData.OperUserId)
	if err != nil {
		s.log.Error(err)
		return
	}

	operators := []int64{e.EvData.OperUserId}

	atUserIds := stream.Of(e.EvData.Comment.ReferUserIds).
		Diff(operators...).
		List()
	participators := stream.Of(utils.ToInt64Array(workItem.Doc.Participators)).
		Diff(operators...).
		Diff(atUserIds...).
		List()
	followers := stream.Of(utils.ToInt64Array(workItem.Doc.Followers)).
		Diff(operators...).
		Diff(atUserIds...).
		Diff(participators...).
		List()

	msg := n_message.NewMessage()
	msg.SetType(e.EvType)
	msg.Space.SpaceId = space.Id
	msg.Space.SpaceName = space.SpaceName
	//故事内容
	//某人
	msg.Notification.Subject = n_message.NewCommentOper(&n_message.UserData{
		Id:       operator.Id,
		NickName: operator.UserNickname,
		Avatar:   operator.Avatar,
		Name:     operator.UserName,
	})
	//某任务某玩意
	msg.Notification.Object = n_message.NewCommentWorkItemObject(&n_message.WorkItemData{
		Id:   workItem.Id,
		Pid:  workItem.Pid,
		Name: workItem.WorkItemName,
	})
	//添加了
	msg.Notification.Action = n_message.ActionType_add
	//某评论
	msg.Notification.SubObject = n_message.NewCommentObject(&n_message.WorkItemCommentData{
		Id:             e.EvData.Comment.Id,
		ReplyCommentId: e.EvData.Comment.ReplyCommentId,
	})
	//评论内容
	describe := utils.ClearRichTextToPlanText(e.EvData.Comment.Content, false)
	msg.Notification.Describe = describe
	//故事发生时间
	msg.Notification.Date = time.Now()

	msg.SetRedirectLink(s.makeWorkItemCommentRedirectLink("", space.Id, workItem.Id, workItem.Pid)).
		SetPopup()

	toFollower := msg.Clone().SetRelation(n_message.Relation_workItemFollower)
	toDirector := msg.Clone().SetRelation(n_message.Relation_workItemDirector)
	toAt := msg.Clone().SetRelation(n_message.Relation_workItemCommentAt)
	toReply := msg.Clone().SetRelation(n_message.Relation_workItemCommentRefer)

	// IM通知
	imMsgTemplate := tea_im.NewRobotMessage().
		SetTitle("评论提醒").SetSubTitle(space.SpaceName).SetSubContent(workItem.WorkItemName).
		SetShowType(tea_im.ShowType_Text).
		SetIcon(s.makeIconResLink("", IconRes_Comment)).
		SetSVGIcon(IconSVGRes_Comment).
		SetUrl(s.makeWorkItemCommentRedirectLink("", space.Id, workItem.Id, workItem.Pid))

	operatorStr := parseUserTmp(operator)
	pubDesc := operatorStr + "发布评论\n" + describe
	atDesc := operatorStr + "@了你\n" + describe
	pubMsg := imMsgTemplate.Clone().SetTextRich(parseToImRich(pubDesc))
	atMsg := imMsgTemplate.Clone().SetTextRich(parseToImRich(atDesc))

	var replyComment *space_work_item_comment.SpaceWorkItemComment
	var replyDesc string
	var replyMsg *tea_im.RobotMessage
	var replyCommentUserId int64
	if e.EvData.Comment.ReplyCommentId > 0 {
		replyComment, _ = s.commentRepo.GetComment(ctx, e.EvData.Comment.ReplyCommentId)
		replyCommentUserId = replyComment.UserId
		replyDesc = operatorStr + "引用了你的评论，回复\n" + describe
		replyMsg = imMsgTemplate.Clone().SetTextRich(parseToImRich(replyDesc))
	}

	allUserIds := stream.Concat(participators, followers, atUserIds)
	if replyComment != nil {
		allUserIds = stream.Concat(allUserIds, []int64{replyComment.UserId})
	}

	notifyCtx := s.buildNotifyCtx(space, allUserIds)

	// 去除被回复人
	participators = stream.Remove(participators, replyCommentUserId)
	followers = stream.Remove(followers, replyCommentUserId)

	utils.Go(func() {
		s.Send2(notifyCtx, toDirector, participators...)
		s.Send2(notifyCtx, toFollower, followers...)
		s.Send2(notifyCtx, toAt, atUserIds...)

		s.pushThirdPlatformMessage2(notifyCtx, pubMsg, participators)
		s.pushThirdPlatformMessage2(notifyCtx, pubMsg, followers)
		s.pushThirdPlatformMessage2(notifyCtx, atMsg, atUserIds)

		if replyComment != nil && !slices.Contains(atUserIds, replyComment.UserId) && operator.Id != replyComment.UserId {
			s.Send2(notifyCtx, toReply, replyComment.UserId)
			s.pushThirdPlatformMessage2(notifyCtx, replyMsg, []int64{replyComment.UserId})
		}
	})

	emitCommentEvent := func(userIds []int64, msg *n_message.Message) {
		for _, userId := range userIds {
			userId := userId // copy
			member := notifyCtx.memberMap[userId]
			globalNotify := notifyCtx.userNotifySwitchGlobalMap[userId]
			spaceNotify := notifyCtx.userNotifySwitchSpaceMap[userId]

			if globalNotify.Value == "1" &&
				spaceNotify.Value == "1" &&
				space.Notify == 1 &&
				(member == nil || member.Notify == 1) {
				//保存评论通知快照
				utils.Go(func() {
					nSnapShot := s.notifySnapShotService.CreateNotifySnapShot(context.Background(), msg.Space.SpaceId, userId, int64(msg.Type), msg.String())
					s.notifySnapShotRepo.CreateNotify(context.Background(), nSnapShot)
				})
			}
		}
	}

	emitCommentEvent(participators, toDirector)
	emitCommentEvent(followers, toFollower)
	emitCommentEvent(atUserIds, toAt)
	if replyComment != nil && !slices.Contains(atUserIds, replyComment.UserId) && operator.Id != replyComment.UserId {
		emitCommentEvent([]int64{replyComment.UserId}, toReply)
	}

	s.CooperateComment(&event.CooperateComment{
		Event:     e.EvType,
		Operator:  e.EvData.OperUserId,
		CommentId: e.EvData.Comment.Id,
		Space:     space,
		WorkItem:  workItem,
	})
}

func (s *Notify) DeleteComment(e *event.DeleteComment) {
	ctx := context.Background()
	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.WorkItemId, &repo.WithDocOption{
		PlanTime:      true,
		Participators: true,
		Followers:     true,
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return
	}

	s.CooperateComment(&event.CooperateComment{
		Event:     e.Event,
		Operator:  e.Operator,
		CommentId: e.CommentId,
		Space:     space,
		WorkItem:  workItem,
	})
}
