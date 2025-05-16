package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_work_item_comment"
)

func SpaceWorkItemCommentEntityToPo(comment *domain.SpaceWorkItemComment) *db.SpaceWorkItemComment {
	return &db.SpaceWorkItemComment{
		Id:             comment.Id,
		UserId:         comment.UserId,
		WorkItemId:     comment.WorkItemId,
		Content:        comment.Content,
		ReferUserIds:   comment.ReferUserIds.ToJsonString(),
		ReplyCommentId: comment.ReplyCommentId,
		Emojis:         comment.Emojis.ToJson(),
		CreatedAt:      comment.CreatedAt,
		UpdatedAt:      comment.UpdatedAt,
		DeletedAt:      comment.DeletedAt,
	}
}

func SpaceWorkItemCommentEntityToPos(comments domain.SpaceWorkItemComments) []*db.SpaceWorkItemComment {
	entities := make([]*db.SpaceWorkItemComment, 0)
	for _, comment := range comments {
		entities = append(entities, SpaceWorkItemCommentEntityToPo(comment))
	}
	return entities
}

func SpaceWorkItemCommentPoToEntity(comment *db.SpaceWorkItemComment) *domain.SpaceWorkItemComment {
	ent := &domain.SpaceWorkItemComment{
		Id:             comment.Id,
		UserId:         comment.UserId,
		WorkItemId:     comment.WorkItemId,
		Content:        comment.Content,
		ReferUserIds:   domain.ReferUserIds{}.FormJsonString(comment.ReferUserIds),
		ReplyCommentId: comment.ReplyCommentId,
		Emojis:         domain.Emojis{}.FromJson(comment.Emojis),
		CreatedAt:      comment.CreatedAt,
		UpdatedAt:      comment.UpdatedAt,
		DeletedAt:      comment.DeletedAt,
	}

	ent.ReferUserIds.FormJsonString(comment.ReferUserIds)
	return ent
}

func SpaceWorkItemCommentPoToEntities(comments []*db.SpaceWorkItemComment) domain.SpaceWorkItemComments {
	entities := make([]*domain.SpaceWorkItemComment, 0)
	for _, comment := range comments {
		entities = append(entities, SpaceWorkItemCommentPoToEntity(comment))
	}
	return entities
}
