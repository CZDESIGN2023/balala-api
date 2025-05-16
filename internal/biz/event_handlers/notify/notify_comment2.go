package notify

import (
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	domain_message "go-cs/internal/domain/pkg/message"
)

func (s *Notify) addCommentByDomainMessage(e *domain_message.CreateComment) {
	s.addComment(&event.AddCommentEvent{
		EvData: &event.AddCommentEventData{
			Comment: &event.Comment{
				Id:             e.CommentId,
				UserId:         e.UserId,
				Content:        e.Content,
				ReferUserIds:   e.ReferUserIds,
				ReplyCommentId: e.ReplyCommentId,
			},
			OperUserId: e.Oper.GetId(),
			WorkItemId: e.WorkItemId,
		},
		EvType: notify.Event_AddCommentEvent,
	})
}

func (s *Notify) updateCommentByDomainMessage(e *domain_message.UpdateComment) {
	s.addComment(&event.AddCommentEvent{
		EvData: &event.AddCommentEventData{
			Comment: &event.Comment{
				Id:             e.CommentId,
				UserId:         e.UserId,
				Content:        e.NewContent,
				ReferUserIds:   e.ReferUserIds,
				ReplyCommentId: e.ReplyCommentId,
			},
			OperUserId: e.Oper.GetId(),
			WorkItemId: e.WorkItemId,
		},
		EvType: notify.Event_AddCommentEvent,
	})
}

func (s *Notify) deleteCommentByDomainMessage(e *domain_message.DeleteComment) {
	s.DeleteComment(&event.DeleteComment{
		Event:      notify.Event_DeleteComment,
		Operator:   e.Oper.GetId(),
		WorkItemId: e.WorkItemId,
		CommentId:  e.CommentId,
	})
}

func (s *Notify) addCommentEmojiByDomainMessage(e *domain_message.AddCommentEmoji) {
	s.addCommentEmoji(&event.AddCommentEmoji{
		Comment: &event.Comment{
			Id:             e.CommentId,
			UserId:         e.UserId,
			Content:        e.Content,
			ReferUserIds:   e.ReferUserIds,
			ReplyCommentId: e.ReplyCommentId,
		},
		Event:      notify.Event_AddCommentEmoji,
		Operator:   e.Oper.GetId(),
		WorkItemId: e.WorkItemId,
		Emoji:      e.Emoji,
	})
}
