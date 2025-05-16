package repo

import (
	"context"
	domain "go-cs/internal/domain/notify_snapshot"
)

type NotifySnapshotRepo interface {
	CreateNotify(ctx context.Context, in *domain.NotifySnapShot) error
	GetUserRelatedCommentIds(ctx context.Context, userId int64, pos, size int) (ids []int64, nextPos int64, hasNext bool, err error)
	GetNotifyByIds(ctx context.Context, userId int64, ids []int64) ([]*domain.NotifySnapShot, error)
	SaveOfflineNotify(ctx context.Context, userId int64, data []byte) error
	GetDelOfflineNotify(ctx context.Context, userId int64) ([]string, error)
}
