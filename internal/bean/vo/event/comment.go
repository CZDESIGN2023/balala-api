package event

import (
	"go-cs/api/notify"
)

type Comment struct {
	Id             int64
	UserId         int64
	Content        string
	ReferUserIds   []int64
	ReplyCommentId int64
}

// 添加评论事件
type AddCommentEvent struct {
	EvType notify.Event //AddCommentEvent
	EvData *AddCommentEventData
}

type AddCommentEventData struct {
	OperUserId int64
	WorkItemId int64
	Comment    *Comment
}

type DeleteComment struct {
	Event      notify.Event
	Operator   int64
	WorkItemId int64
	CommentId  int64
}

type UpdateComment struct {
	Event      notify.Event
	Operator   int64
	WorkItemId int64
	CommentId  int64
	Comment    *Comment
}

type AddCommentEmoji struct {
	Event      notify.Event
	Operator   int64
	WorkItemId int64
	Comment    *Comment
	Emoji      string
}
