package message

import shared "go-cs/internal/pkg/domain"

const (
	Message_Type_Comment_Create      shared.MessageType = "Domain_Message.Comment.Create"
	Message_Type_Comment_Update      shared.MessageType = "Domain_Message.Comment.Update"
	Message_Type_Comment_Delete      shared.MessageType = "Domain_Message.Comment.Delete"
	Message_Type_Comment_EmojiAdd    shared.MessageType = "Domain_Message.Comment.EmojiAdd"
	Message_Type_Comment_EmojiRemove shared.MessageType = "Domain_Message.Comment.EmojiRemove"
)

type CreateComment struct {
	shared.DomainMessageBase

	WorkItemId int64
	CommentId  int64

	UserId         int64
	Content        string
	ReferUserIds   []int64
	ReplyCommentId int64
}

func (ops *CreateComment) MessageType() shared.MessageType {
	return Message_Type_Comment_Create
}

type UpdateComment struct {
	shared.DomainMessageBase

	WorkItemId int64
	CommentId  int64

	UserId         int64
	OldContent     string
	NewContent     string
	ReferUserIds   []int64
	ReplyCommentId int64
}

func (ops *UpdateComment) MessageType() shared.MessageType {
	return Message_Type_Comment_Update
}

type DeleteComment struct {
	shared.DomainMessageBase

	WorkItemId int64
	CommentId  int64

	UserId         int64
	Content        string
	ReferUserIds   []int64
	ReplyCommentId int64
}

func (ops *DeleteComment) MessageType() shared.MessageType {
	return Message_Type_Comment_Delete
}

type AddCommentEmoji struct {
	shared.DomainMessageBase

	SpaceId    int64
	WorkItemId int64
	CommentId  int64

	UserId         int64
	Emoji          string
	Content        string
	ReferUserIds   []int64
	ReplyCommentId int64
}

func (ops *AddCommentEmoji) MessageType() shared.MessageType {
	return Message_Type_Comment_EmojiAdd
}

type RemoveCommentEmoji struct {
	shared.DomainMessageBase

	SpaceId    int64
	WorkItemId int64
	CommentId  int64

	UserId         int64
	Emoji          string
	Content        string
	ReferUserIds   []int64
	ReplyCommentId int64
}

func (ops *RemoveCommentEmoji) MessageType() shared.MessageType {
	return Message_Type_Comment_EmojiRemove
}
