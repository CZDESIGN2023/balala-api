package space_work_item_comment

import (
	"time"
)

func NewSpaceWorkItemComment(
	workItemId int64,
	userId int64,
	comment string,
	referUserIds ReferUserIds,
	replyCommentId int64,
) *SpaceWorkItemComment {
	ts := time.Now().Unix()
	c := &SpaceWorkItemComment{
		UserId:         userId,
		WorkItemId:     workItemId,
		Content:        comment,
		ReferUserIds:   referUserIds,
		ReplyCommentId: replyCommentId,
		CreatedAt:      ts,
		UpdatedAt:      ts,
	}

	return c
}
