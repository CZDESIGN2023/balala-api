package repo

import (
	"context"
	domain "go-cs/internal/domain/space_work_item_comment"
	"time"
)

type SpaceWorkItemCommentRepo interface {
	SpaceWorkItemCommentQueryRepo
	SpaceWorkItemCommentCacheRepo

	CreateComment(ctx context.Context, comment *domain.SpaceWorkItemComment) error
	SaveComment(ctx context.Context, comment *domain.SpaceWorkItemComment) error

	GetComment(ctx context.Context, id int64) (*domain.SpaceWorkItemComment, error)
	DelComment(ctx context.Context, id int64) (int64, error)
	GetCommentByIds(ctx context.Context, ids []int64) (domain.SpaceWorkItemComments, error)
	DelCommentByIds(ctx context.Context, ids []int64) (int64, error)
	CommentMap(ctx context.Context, ids []int64) (map[int64]*domain.SpaceWorkItemComment, error)
	GetCommentByWorkItemId(ctx context.Context, workItemId int64) (domain.SpaceWorkItemComments, error)

	DelCommentByWorkItemIds(ctx context.Context, workItemIds []int64) (int64, error)
	CountCommentByWorkItemId(ctx context.Context, workItemId int64) (int64, error)
	CountCommentByWorkItemIds(ctx context.Context, workItemIds []int64) (map[int64]int64, error)
	CountWorkItemCommentNumByTime(ctx context.Context, workItemId int64, t time.Time) (int64, error)
	SetUserReadTime(ctx context.Context, userId int64, workItemId int64, t time.Time) error
	GetUserReadTime(ctx context.Context, userId int64, workItemId int64) (time.Time, error)
}

type SpaceWorkItemCommentQueryRepo interface {
	QCommentPagination(ctx context.Context, workItemId int64, pos, size int, order string) (domain.SpaceWorkItemComments, error)
}

type SpaceWorkItemCommentCacheRepo interface {
	IncrUnreadNumForUser(ctx context.Context, workItemId int64, userIds []int64) error
	RemoveUnreadNumForUser(ctx context.Context, userId int64, workItemIds []int64) error
	GetUserUnreadNum(ctx context.Context, userId int64, workItemId int64) (int64, error)
	UserUnreadNumMapByWorkItemIds(ctx context.Context, userId int64, workItemIds []int64) (map[int64]int64, error)
	SetUserUnreadNum(ctx context.Context, userId int64, numMap map[int64]int64) error
}
