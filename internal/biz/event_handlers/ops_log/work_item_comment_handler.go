package ops_log

import (
	"context"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/oper"
)

func (s *OpsLogEventHandlers) addCommentHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.CreateComment)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, witem.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  witem.WorkItemName,
		}

		content := utils.ClearRichTextToPlanText(opsLog.Content, true)
		var msg string
		if opsLog.ReplyCommentId != 0 {
			replyComment, _ := s.commentRepo.GetComment(ctx, opsLog.ReplyCommentId)
			replyCommentUser, _ := s.userRepo.GetUserByUserId(ctx, replyComment.UserId)
			replyCommentContent := utils.ClearRichTextToPlanText(replyComment.Content, true)

			msg = "引用了 " + W(U(replyCommentUser)+"评论") + Q(replyCommentContent) + "，回复" + Q(content)
		} else {
			msg = "新增 " + W("评论") + Q(content)
		}

		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) deleteCommentHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.DeleteComment)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, witem.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  witem.WorkItemName,
		}

		content := utils.ClearRichTextToPlanText(opsLog.Content, true)
		msg := "删除 " + W("评论") + Q(content)

		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}

func (s *OpsLogEventHandlers) updateCommentHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.UpdateComment)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, witem.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  witem.WorkItemName,
		}

		oldContent := utils.ClearRichTextToPlanText(opsLog.OldContent, true)
		newContent := utils.ClearRichTextToPlanText(opsLog.NewContent, true)
		msg := "将 " + W("评论") + Q(oldContent) + "更新为" + Q(newContent)

		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)

	}
}

func (s *OpsLogEventHandlers) addCommentEmojiHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.AddCommentEmoji)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, witem.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  witem.WorkItemName,
		}

		content := utils.ClearRichTextToPlanText(opsLog.Content, true)

		msg := "对 " + W("评论") + Q(content) + "，回复" + Q(opsLog.Emoji)

		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}
func (s *OpsLogEventHandlers) removeCommentEmojiHandler(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages) {

	for _, v := range opsLogs {
		opsLog := v.(*domain_message.RemoveCommentEmoji)
		if opsLog == nil {
			continue
		}

		operUser := s.getOperUserInfo(ctx, opsLog.GetOper())
		if operUser == nil {
			continue
		}

		witem, err := s.wItemRepo.GetWorkItem(ctx, opsLog.WorkItemId, nil, nil)
		if err != nil {
			continue
		}

		space, err := s.spaceRepo.GetSpace(ctx, witem.SpaceId)
		if err != nil {
			continue
		}

		operLogger.Operator = operUser
		result := &oper.OperResultInfo{
			SpaceId:      space.Id,
			SpaceName:    space.SpaceName,
			BusinessType: oper.BusinessTypeModify,
			ModuleType:   oper.ModuleTypeSpaceWorkItem,
			ModuleId:     int(opsLog.WorkItemId),
			ModuleTitle:  witem.WorkItemName,
		}

		content := utils.ClearRichTextToPlanText(opsLog.Content, true)

		msg := "对 " + W("评论") + Q(content) + "，取消" + Q(opsLog.Emoji)

		result.OperMsg = msg
		s.invokeOperLog(ctx, operLogger, result)
	}
}
