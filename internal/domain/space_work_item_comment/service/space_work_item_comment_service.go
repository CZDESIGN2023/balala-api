package service

import (
	"context"
	domain "go-cs/internal/domain/space_work_item_comment"
	repo "go-cs/internal/domain/space_work_item_comment/repo"
)

type SpaceWorkItemCommentService struct {
	repo repo.SpaceWorkItemCommentRepo
}

func NewSpaceWorkItemCommentService(
	repo repo.SpaceWorkItemCommentRepo,

) *SpaceWorkItemCommentService {
	return &SpaceWorkItemCommentService{
		repo: repo,
	}
}

type CreateCommentRequest struct {
	UserId         int64
	WorkItemId     int64
	Content        string
	ReferUserIds   domain.ReferUserIds
	ReplyCommentId int64
}

func (s *SpaceWorkItemCommentService) CreateComment(ctx context.Context, req *CreateCommentRequest) (*domain.SpaceWorkItemComment, error) {

	return domain.NewSpaceWorkItemComment(
		req.WorkItemId,
		req.UserId,
		req.Content,
		req.ReferUserIds,
		req.ReplyCommentId,
	), nil
}
