package service

import (
	"context"
	domain "go-cs/internal/domain/notify_snapshot"
	"go-cs/internal/domain/notify_snapshot/repo"
)

type NotifySnapShotService struct {
	repo repo.NotifySnapshotRepo
}

func NewNotifySnapShotService(
	repo repo.NotifySnapshotRepo,
) *NotifySnapShotService {
	return &NotifySnapShotService{
		repo: repo,
	}
}

func (s *NotifySnapShotService) CreateNotifySnapShot(ctx context.Context, spaceId int64, uid int64, typ int64, doc string) *domain.NotifySnapShot {
	return domain.NewNotifySnapShot(
		spaceId,
		uid,
		domain.NotifyEventType(typ),
		doc,
	)
}
